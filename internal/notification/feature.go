package notification

import (
	"context"
	"sync"

	"github.com/flohansen/dasher/internal/sqlc"
	"github.com/flohansen/dasher/proto"
	"google.golang.org/grpc"
)

type FeatureStore interface {
	GetAll(ctx context.Context) ([]sqlc.Feature, error)
}

type FeatureNotifier struct {
	proto.UnimplementedFeatureStateServiceServer
	store   FeatureStore
	streams map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}
	mu      *sync.Mutex
}

func NewFeatureNotifier(grpcServer grpc.ServiceRegistrar, store FeatureStore) *FeatureNotifier {
	notifier := FeatureNotifier{
		streams: make(map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}),
		store:   store,
		mu:      &sync.Mutex{},
	}
	proto.RegisterFeatureStateServiceServer(grpcServer, &notifier)
	return &notifier
}

func (n *FeatureNotifier) SubscribeFeatureChanges(_ *proto.FeatureSubscription, stream proto.FeatureStateService_SubscribeFeatureChangesServer) error {
	n.mu.Lock()
	n.streams[stream] = struct{}{}
	n.mu.Unlock()

	features, err := n.store.GetAll(stream.Context())
	if err != nil {
		return err
	}

	for _, feature := range features {
		n.Notify(feature)
	}

	<-stream.Context().Done()

	n.mu.Lock()
	delete(n.streams, stream)
	n.mu.Unlock()

	return nil
}

func (n *FeatureNotifier) Notify(feature sqlc.Feature) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	for stream := range n.streams {
		err := stream.Send(&proto.Feature{
			FeatureId: feature.FeatureID,
			Enabled:   feature.Enabled,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
