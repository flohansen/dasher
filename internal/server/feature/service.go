//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
//go:generate protoc --proto_path=../../../pkg/proto --go_out=../../../pkg/proto --go_opt=paths=source_relative --go-grpc_out=../../../pkg/proto --go-grpc_opt=paths=source_relative feature.proto

package feature

import (
	"context"
	"sync"

	"github.com/flohansen/dasher/internal/sqlc"
	"github.com/flohansen/dasher/pkg/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type FeatureStore interface {
	GetAll(ctx context.Context) ([]sqlc.Feature, error)
	Upsert(ctx context.Context, feature sqlc.Feature) error
}

type Service struct {
	proto.UnimplementedFeatureStateServiceServer
	store         FeatureStore
	subscriptions map[string]map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}
	mu            sync.Mutex
}

func NewService(grpcServer grpc.ServiceRegistrar, store FeatureStore) *Service {
	notifier := Service{
		subscriptions: make(map[string]map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}),
		store:         store,
	}
	proto.RegisterFeatureStateServiceServer(grpcServer, &notifier)
	return &notifier
}

func (n *Service) registerSubscriptions(subscription *proto.FeatureSubscription, stream proto.FeatureStateService_SubscribeFeatureChangesServer) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	for _, f := range subscription.FeatureToggles {
		if _, ok := n.subscriptions[f.FeatureId]; !ok {
			n.subscriptions[f.FeatureId] = make(map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{})
		}

		n.subscriptions[f.FeatureId][stream] = struct{}{}
	}

	features, err := n.store.GetAll(stream.Context())
	if err != nil {
		return err
	}

	featuresAdded, err := n.createNonExistingFeatures(stream.Context(), subscription, features)
	if err != nil {
		return err
	}

	// Initially send states of all features
	for _, feature := range append(features, featuresAdded...) {
		err := stream.Send(&proto.FeatureToggleChange{
			FeatureId: feature.FeatureID,
			Enabled:   feature.Enabled,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Service) unregisterSubscriptions(subscription *proto.FeatureSubscription, stream proto.FeatureStateService_SubscribeFeatureChangesServer) {
	n.mu.Lock()
	defer n.mu.Unlock()

	for _, f := range subscription.FeatureToggles {
		delete(n.subscriptions[f.FeatureId], stream)
	}
}

func (n *Service) createNonExistingFeatures(ctx context.Context, subscription *proto.FeatureSubscription, features []sqlc.Feature) ([]sqlc.Feature, error) {
	featureLookup := make(map[string]sqlc.Feature)
	for _, f := range features {
		featureLookup[f.FeatureID] = f
	}

	var featuresToAdd []sqlc.Feature
	for _, f := range subscription.FeatureToggles {
		if _, ok := featureLookup[f.FeatureId]; !ok {
			featuresToAdd = append(featuresToAdd, sqlc.Feature{
				FeatureID:   f.FeatureId,
				Description: f.Description,
				Enabled:     false,
			})
		}
	}

	for _, f := range featuresToAdd {
		if err := n.store.Upsert(ctx, f); err != nil {
			return nil, errors.Wrap(err, "feature store upsert")
		}
	}

	return featuresToAdd, nil
}

func (n *Service) SubscribeFeatureChanges(subscription *proto.FeatureSubscription, stream proto.FeatureStateService_SubscribeFeatureChangesServer) error {
	n.registerSubscriptions(subscription, stream)

	// Wait for client closing the connection
	<-stream.Context().Done()

	n.unregisterSubscriptions(subscription, stream)
	return nil
}

func (n *Service) Notify(feature sqlc.Feature) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	for stream := range n.subscriptions[feature.FeatureID] {
		err := stream.Send(&proto.FeatureToggleChange{
			FeatureId: feature.FeatureID,
			Enabled:   feature.Enabled,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
