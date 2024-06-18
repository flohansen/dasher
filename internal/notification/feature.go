package notification

import (
	"sync"

	"github.com/flohansen/dasher-server/internal/sqlc"
	"github.com/flohansen/dasher-server/proto"
	"google.golang.org/grpc"
)

type FeatureNotifier struct {
	proto.UnimplementedFeatureStateServiceServer
	streams map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}
	mu      *sync.Mutex
}

func NewFeatureNotifier(grpcServer grpc.ServiceRegistrar) *FeatureNotifier {
	notifier := FeatureNotifier{
		streams: make(map[proto.FeatureStateService_SubscribeFeatureChangesServer]struct{}),
		mu:      &sync.Mutex{},
	}
	proto.RegisterFeatureStateServiceServer(grpcServer, &notifier)
	return &notifier
}

func (n *FeatureNotifier) SubscribeFeatureChanges(_ *proto.FeatureSubscription, stream proto.FeatureStateService_SubscribeFeatureChangesServer) error {
	n.mu.Lock()
	n.streams[stream] = struct{}{}
	n.mu.Unlock()

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
