package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/transport/amqp"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log2 "github.com/rs/zerolog/log"
	amqp2 "github.com/streadway/amqp"
	"net/http"
	"os"
	"strings"
	"time"
)

type BuilderService interface {
	BuildPlane(s string) (string, error)
}

type builderService struct{}

func (builderService) BuildPlane(s string) (string, error) {
	return strings.ToUpper(s), nil
}

type buildPlaneRequest struct {
	S string `json:"s"`
}

type buildPlaneResponse struct {
	V string `json:"v"`
	Err string `json:"err,omitempty"`
}

func makeBuildPlaneEndpoint(svc BuilderService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(buildPlaneRequest)
		v, err := svc.BuildPlane(req.S)
		if err != nil {
			return buildPlaneResponse{v, err.Error()}, nil
		}
		log2.Warn().Msgf("RESPONSE: %s", v)
		return buildPlaneResponse{v, ""}, nil
	}
}

func decodeBuildPlaneRequest(_ context.Context, a *amqp2.Delivery) (interface{}, error) {
	var request buildPlaneRequest
	if err := json.NewDecoder(bytes.NewReader(a.Body)).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func main() {
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:9090", nil)

	var svc BuilderService
	svc = builderService{}
	svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)
	var connection *amqp2.Connection
	var err error
	for {
		connection, err = amqp2.Dial(os.Getenv("AMQP_URL"))
		if err == nil {
			break
		}
		log2.Warn().Err(err).Msgf("unable to connect to AMQP server %s", os.Getenv("AMQP_URL"))
		time.Sleep(1 * time.Second)
	}

	defer connection.Close()
	ch, err := connection.Channel()
	if err != nil {
		log2.Fatal().Err(err).Msgf("unable to bind to AMQP channel")
	}
	defer ch.Close()
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