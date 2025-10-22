package service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
)

// DeleteApplicant deletes an applicant by ID
func (s *ApplicantService) DeleteApplicant(ctx context.Context, req *applicantsv1.DeleteApplicantRequest) (*applicantsv1.DeleteApplicantResponse, error) {
	// Validate input
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id must be positive")
	}

	s.logger.Debug("deleting applicant", zap.Int64("id", req.Id))

	err := s.queries.DeleteApplicant(ctx, req.Id)
	if err != nil {
		s.logger.Error("failed to delete applicant", zap.Int64("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete applicant: %v", err)
	}

	return &applicantsv1.DeleteApplicantResponse{
		Success: true,
	}, nil
}
