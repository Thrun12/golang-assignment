package service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
	"github.com/Thrun12/golang-assignment/internal/util"
)

// ListApplicants retrieves a list of applicants with pagination
func (s *ApplicantService) ListApplicants(ctx context.Context, req *applicantsv1.ListApplicantsRequest) (*applicantsv1.ListApplicantsResponse, error) {
	s.logger.Debug("listing applicants",
		zap.Int32("limit", req.Limit),
		zap.Int32("offset", req.Offset),
		zap.String("position", req.Position),
	)

	// Default limit
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Ensure offset is non-negative
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Get applicants
	applicants, err := s.queries.ListApplicants(ctx, sqlc.ListApplicantsParams{
		Limit:    limit,
		Offset:   offset,
		Position: req.Position,
		Status:   int32(req.Status),
		MinScore: req.MinScore,
	})
	if err != nil {
		s.logger.Error("failed to list applicants", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list applicants: %v", err)
	}

	// Get total count
	totalCount, err := s.queries.CountApplicants(ctx, sqlc.CountApplicantsParams{
		Position: req.Position,
		Status:   int32(req.Status),
		MinScore: req.MinScore,
	})
	if err != nil {
		s.logger.Error("failed to count applicants", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to count applicants: %v", err)
	}

	// Convert to proto
	protoApplicants := make([]*applicantsv1.JobApplicant, len(applicants))
	for i, app := range applicants {
		protoApplicants[i] = util.DbApplicantToProto(&app)
	}

	return &applicantsv1.ListApplicantsResponse{
		Applicants: protoApplicants,
		TotalCount: int32(totalCount),
		Limit:      limit,
		Offset:     offset,
	}, nil
}
