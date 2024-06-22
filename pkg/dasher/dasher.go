package dasher

import (
	"context"
	"log"

	"github.com/flohansen/dasher/pkg/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var registeredFeatures = make(map[string]*Feature)

func MustRegister(f *Feature) {
	registeredFeatures[f.FeatureID] = f
}

type Feature struct {
	FeatureID   string
	Description string
	Enabled     bool
}

type FeatureParams struct {
	Name        string
	Description string
}

func NewFeature(opts FeatureParams) *Feature {
	return &Feature{
		FeatureID:   opts.Name,
		Description: opts.Description,
	}
}

func Connect(ctx context.Context, addr string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errors.Wrap(err, "create new grpc client")
	}

	client := proto.NewFeatureStateServiceClient(conn)
	stream, err := client.SubscribeFeatureChanges(ctx, &proto.FeatureSubscription{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "subscribe to feature changes"))
	}

	for {
		feature, err := stream.Recv()
		if err != nil {
			errors.Wrap(err, "receiving feature message")
		}

		if registeredFeature, ok := registeredFeatures[feature.FeatureId]; ok {
			registeredFeature.Enabled = feature.Enabled
		}
	}
}
