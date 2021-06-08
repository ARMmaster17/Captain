package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/transport/amqp"
	log2 "github.com/rs/zerolog/log"
	"strings"

	"github.com/go-kit/kit/endpoint"
	amqp2 "github.com/streadway/amqp"
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
	svc := builderService{}
	connection, err := amqp2.Dial("testURL")
	if err != nil {
		log2.Fatal().Err(err).Msgf("unable to connect to AMQP server %s", "testURL")
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