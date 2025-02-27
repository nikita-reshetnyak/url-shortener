package analyticsgrpc

import (
	"context"
	"fmt"
	"log/slog"

	v1 "github.com/nikita-reshetnyak/analytics-protos/gen/go/analytics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	api v1.AnalyticsClient
	log *slog.Logger
}

func New(
	log *slog.Logger,
	addr string,
) (*Client, error) {
	const op = "clients.analytics.grpc.New"
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	grpcClient := v1.NewAnalyticsClient(conn)
	return &Client{api: grpcClient, log: log}, nil
}
func (c *Client) SendEvent(ctx context.Context, name string, date *timestamppb.Timestamp) error {
	const op = "clients.analytics.grpc.SendEvent"
	_, err := c.api.SendEvent(ctx, &v1.Event{
		Name: name,
		Date: date,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
