package notification

import (
	"context"
	"sync"

	"github.com/flohansen/dasher/internal/sqlc"
	"github.com/flohansen/dasher/pkg/proto"
	"google.golang.org/grpc"
)

type FeatureStore interface {
	GetAll(ctx context.Context) ([]sqlc.Feature, error)
	Upsert(ctx context.Context, feature sqlc.Feature) error
}

type FeatureNotifier struct {
	proto.UnimplementedFeatureStateServiceServer
	store         FeatureStore
	streams       map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}
	subscriptions map[string]map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}
	mu            *sync.Mutex
}

func NewFeatureNotifier(grpcServer grpc.ServiceRegistrar, store FeatureStore) *FeatureNotifier {
	notifier := FeatureNotifier{
		streams:       make(map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}),
		subscriptions: make(map[string]map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}),
		store:         store,
		mu:            &sync.Mutex{},
	}
	proto.RegisterFeatureStateServiceServer(grpcServer, &notifier)
	return &notifier
}

func (n *FeatureNotifier) SubscribeFeatureChanges(subscription *proto.FeatureSubscription, stream proto.FeatureStateService_SubscribeFeatureChangesServer) error {
	n.mu.Lock()
	for _, f := range subscription.FeatureToggles {
		if _, ok := n.subscriptions[f.FeatureId]; !ok {
			n.subscriptions[f.FeatureId] = make(map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{})
		}

		n.subscriptions[f.FeatureId][stream] = struct{}{}
	}
	n.mu.Unlock()

	features, err := n.store.GetAll(stream.Context())
	if err != nil {
		return err
	}

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
		n.store.Upsert(stream.Context(), f)
	}

	for _, feature := range features {
		n.Notify(feature)
	}

	<-stream.Context().Done()

	n.mu.Lock()
	for _, f := range subscription.FeatureToggles {
		delete(n.subscriptions[f.FeatureId], stream)
	}
	n.mu.Unlock()

	return nil
}

func (n *FeatureNotifier) Notify(feature sqlc.Feature) error {
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
