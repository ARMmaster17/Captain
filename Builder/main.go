package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ARMmaster17/Captain/Builder/BuilderSvc"
	"github.com/ARMmaster17/Captain/Shared"
	"github.com/ARMmaster17/Captain/Shared/DB"
	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/transport/amqp"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log2 "github.com/rs/zerolog/log"
	amqp2 "github.com/streadway/amqp"
	"net/http"
	"os"
	"time"
)

type BuilderService interface {
	BuildPlane(planeId int) (string, error)
}

type builderService struct{
	ProvisionEndpoint endpoint.Endpoint
}

func (b builderService) BuildPlane(planeId int) (string, error) {
	db, err := DB.ConnectToDB()
	if err != nil {
		log2.Error().Msgf("unable to connect to database:\n%w", err)
		return "", err
	}
	err = BuilderSvc.BuildPlane(db, int64(planeId))
	if err != nil {
		log2.Error().Err(err).Msgf("unable to build plane with ID %d", planeId)
		return "", err
	}
	response, err := b.ProvisionEndpoint(context.Background(), Shared.ProvisionPlaneRequest{S: string(planeId)})
	if err != nil {
		log2.Error().Err(err).Msgf("unable to send provisioning request")
		return "", err
	} else {
		log2.Info().Msgf("provisioning service returned the message %v", response)
	}
	return "", nil
}

func makeBuildPlaneEndpoint(svc BuilderService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(Shared.BuildPlaneRequest)
		v, err := svc.BuildPlane(req.S)
		if err != nil {
			return Shared.BuildPlaneResponse{V: v, Err: err.Error()}, nil
		}
		log2.Warn().Msgf("RESPONSE: %s", v)
		return Shared.BuildPlaneResponse{V: v}, nil
	}
}

func decodeBuildPlaneRequest(_ context.Context, a *amqp2.Delivery) (interface{}, error) {
	var request Shared.BuildPlaneRequest
	if err := json.NewDecoder(bytes.NewReader(a.Body)).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func main() {
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "captain",
		Subsystem: "builder",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "captain",
		Subsystem: "builder",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "captain",
		Subsystem: "builder",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:9090", nil)

	var connection *amqp2.Connection
	var err error
	for {
		connection, err = amqp2.Dial(os.Getenv("AMQP_URL"))
		if err == nil {
			break
		}
		log2.Warn().Err(err).Msgf("unable to connect to AMQP server %s, retrying in 1s", os.Getenv("AMQP_URL"))
		time.Sleep(1 * time.Second)
	}

	defer connection.Close()
	ch, err := connection.Channel()
	if err != nil {
		log2.Fatal().Err(err).Msgf("unable to bind to AMQP channel")
	}
	defer ch.Close()

	//////////////////////////////////////////
	// Publisher stuff
	ch2, err := connection.Channel()
	if err != nil {
		log2.Fatal().Err(err).Msgf("unable to bind to AMQP channel 2")
	}
	defer ch2.Close()

	provisionQueue, err := ch2.QueueDeclare("preflight_provision", false, false, false, false, nil)
	provisionPublisher := amqp.NewPublisher(ch2, &provisionQueue, encodeProvisionRequest, decodeProvisionResponse, amqp.PublisherBefore(amqp.SetPublishKey(provisionQueue.Name)))
	provisionEndpoint := provisionPublisher.Endpoint()
	//////////////////////////////////////////

	var svc BuilderService
	svc = builderService{
		ProvisionEndpoint: provisionEndpoint,
	}
	svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)
	buildQueue, err := ch.QueueDeclare("builder_build", false, false, false, false, nil)
	buildMsgs, err := ch.Consume(buildQueue.Name, "", false, false, false, false, nil)
	buildAMQPHandler := amqp.NewSubscriber(makeBuildPlaneEndpoint(svc), decodeBuildPlaneRequest, amqp.EncodeJSONResponse)
	buildListener := buildAMQPHandler.ServeDelivery(ch)

	forever := make(chan bool)

	go func() {
		for true {
			select {
			case buildDeliv := <-buildMsgs:
				log2.Info().Msg("received build job")
				buildListener(&buildDeliv)
				buildDeliv.Ack(false)
			}
		}
	}()

	<-forever
}

func encodeProvisionRequest(ctx context.Context, publishing *amqp2.Publishing, request interface{}) error {
	provisionRequest, ok := request.(Shared.ProvisionPlaneRequest)
	if !ok {
		return fmt.Errorf("unable to encode provisioning request")
	}
	b, err := json.Marshal(provisionRequest)
	if err != nil {
		return err
	}
	publishing.Body = b
	return nil
}

func decodeProvisionResponse(ctx context.Context, delivery *amqp2.Delivery) (interface{}, error) {
	var response Shared.ProvisionPlaneResponse
	err := json.Unmarshal(delivery.Body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}