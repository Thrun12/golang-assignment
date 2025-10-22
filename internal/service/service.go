package service

import (
	"go.uber.org/zap"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
)

// ApplicantService provides business logic for applicant operations and implements the gRPC service
type ApplicantService struct {
	applicantsv1.UnimplementedApplicantsServiceServer
	queries sqlc.Querier
	logger  *zap.Logger
}

// NewApplicantService creates a new applicant service
func NewApplicantService(queries sqlc.Querier, logger *zap.Logger) *ApplicantService {
	return &ApplicantService{
		queries: queries,
		logger:  logger,
	}
}
