package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ARMmaster17/Captain/Shared"
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

type PreflightService interface {
	ProvisionPlane(s string) (string, error)
}

type preflightService struct{}

func (preflightService) ProvisionPlane(s string) (string, error) {
	log2.Info().Msgf("INPUT LINE: %s", s)
	return strings.ToLower(s), nil
}

func makeProvisionPlaneEndpoint(svc PreflightService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(Shared.ProvisionPlaneRequest)
		v, err := svc.ProvisionPlane(req.S)
		if err != nil {
			return Shared.ProvisionPlaneResponse{v, err.Error()}, nil
		}
		log2.Warn().Msgf("RESPONSE: %s", v)
		return Shared.ProvisionPlaneResponse{v, ""}, nil
	}
}

func decodeProvisionPlaneRequest(_ context.Context, a *amqp2.Delivery) (interface{}, error) {
	var request Shared.ProvisionPlaneRequest
	log2.Info().Msgf("RAW BODY: %s", string(a.Body))
	if err := json.NewDecoder(bytes.NewReader(a.Body)).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func main() {
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "captain",
		Subsystem: "preflight",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "captain",
		Subsystem: "preflight",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "captain",
		Subsystem: "preflight",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe("0.0.0.0:9090", nil)

	var svc PreflightService
	svc = preflightService{}
	svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)
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
	provisionQueue, err := ch.QueueDeclare("preflight_provision", false, false, false, false, nil)
	provisionMsgs, err := ch.Consume(provisionQueue.Name, "", false, false, false, false, nil)
	provisionAMQPHandler := amqp.NewSubscriber(makeProvisionPlaneEndpoint(svc), decodeProvisionPlaneRequest, amqp.EncodeJSONResponse)
	provisionListener := provisionAMQPHandler.ServeDelivery(ch)

	forever := make(chan bool)

	go func() {
		for true {
			select {
			case provisionDeliv := <- provisionMsgs:
				log2.Info().Msg("received provision job")
				provisionListener(&provisionDeliv)
				provisionDeliv.Ack(false)
			}
		}
	}()

	<- forever
}