package service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/util"
)

// GetApplicant retrieves an applicant by ID
func (s *ApplicantService) GetApplicant(ctx context.Context, req *applicantsv1.GetApplicantRequest) (*applicantsv1.GetApplicantResponse, error) {
	// Validate input
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id must be positive")
	}

	s.logger.Debug("getting applicant", zap.Int64("id", req.Id))

	applicant, err := s.queries.GetApplicant(ctx, req.Id)
	if err != nil {
		s.logger.Error("failed to get applicant", zap.Int64("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "applicant not found: %v", err)
	}

	return &applicantsv1.GetApplicantResponse{
		Applicant: util.DbApplicantToProto(&applicant),
	}, nil
}
