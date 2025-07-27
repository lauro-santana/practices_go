package grpc

import (
	"context"

	"movie.com/gen"
	"movie.com/internal/grpcutil"
	"movie.com/pkg/discovery"
	"movie.com/rating/pkg/model"
)

// Gateway defines an gRPC gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})
	if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}

// PutRating writes a rating.
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	_, err = client.PutRating(ctx,
		&gen.PutRatingRequest{
			UserId:      string(rating.UserID),
			RecordId:    string(recordID),
			RecordType:  string(recordType),
			RatingValue: int32(rating.Value),
		})
	if err != nil {
		return err
	}
	return nil
}
