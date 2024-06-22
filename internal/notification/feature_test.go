package notification

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/flohansen/dasher/internal/notification/mocks"
	"github.com/flohansen/dasher/internal/sqlc"
	"github.com/flohansen/dasher/pkg/proto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestFeatureNotifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockFeatureStore(ctrl)
	lis := bufconn.Listen(2048)
	s := grpc.NewServer()
	notifier := NewFeatureNotifier(s, store)

	go func() {
		s.Serve(lis)
	}()

	t.Run("Notify", func(t *testing.T) {
		t.Run("should broadcast the change", func(t *testing.T) {
			// given
			store.EXPECT().
				GetAll(gomock.Any()).
				Return([]sqlc.Feature{}, nil).
				Times(2)

			ctx := context.Background()
			client1 := createClient(t, ctx, bufDialer(lis))
			stream1, err := client1.SubscribeFeatureChanges(ctx, &proto.FeatureSubscription{})
			if err != nil {
				t.Fatal(err)
			}

			client2 := createClient(t, ctx, bufDialer(lis))
			stream2, err := client2.SubscribeFeatureChanges(ctx, &proto.FeatureSubscription{})
			if err != nil {
				t.Fatal(err)
			}

			// when
			time.Sleep(10 * time.Millisecond) // workaround to wait until both clients subscribed
			notifier.Notify(sqlc.Feature{FeatureID: "SOME_FEATURE_ID"})

			// then
			feature1, err := stream1.Recv()
			assert.NoError(t, err)
			assert.Equal(t, "SOME_FEATURE_ID", feature1.FeatureId)

			feature2, err := stream2.Recv()
			assert.NoError(t, err)
			assert.Equal(t, "SOME_FEATURE_ID", feature2.FeatureId)
		})
	})
}

func createClient(t *testing.T, ctx context.Context, dialer func(context.Context, string) (net.Conn, error)) proto.FeatureStateServiceClient {
	conn, err := grpc.DialContext(
		ctx, "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatal(err)
	}

	return proto.NewFeatureStateServiceClient(conn)
}
