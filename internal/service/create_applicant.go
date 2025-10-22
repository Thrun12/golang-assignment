package service

import (
	"context"
	"strings"

	"github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
	"github.com/Thrun12/golang-assignment/internal/util"
)

// CreateApplicant creates a new applicant with calculated overall score
func (s *ApplicantService) CreateApplicant(ctx context.Context, req *applicantsv1.CreateApplicantRequest) (*applicantsv1.CreateApplicantResponse, error) {
	// Validate input
	if err := util.ValidateApplicant(req.Name, req.Email, req.Position, req.YearsExperience, req.GithubStars, req.InterviewScore, req.CulturalFitScore, req.TechnicalScore, false, 0); err != nil {
		s.logger.Debug("validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	s.logger.Debug("creating applicant", zap.String("email", req.Email))

	// Calculate overall score using our sophisticated (totally unbiased) algorithm
	overallScore := util.CalculateOverallScore(
		req.Name,
		req.Skills,
		req.YearsExperience,
		req.InterviewScore,
		req.CulturalFitScore,
		req.TechnicalScore,
		req.CanExitVim,
		req.KnowsGo,
		req.DebugsInProduction,
	)

	// Prepare optional fields
	var funFact, availability, salaryExpectation *string
	if req.FunFact != "" {
		funFact = &req.FunFact
	}
	if req.Availability != "" {
		availability = &req.Availability
	}
	if req.SalaryExpectation != "" {
		salaryExpectation = &req.SalaryExpectation
	}

	// Create applicant
	applicant, err := s.queries.CreateApplicant(ctx, sqlc.CreateApplicantParams{
		Name:               req.Name,
		Email:              req.Email,
		Position:           req.Position,
		YearsExperience:    req.YearsExperience,
		Skills:             req.Skills,
		GithubStars:        req.GithubStars,
		CanExitVim:         req.CanExitVim,
		KnowsGo:            req.KnowsGo,
		DebugsInProduction: req.DebugsInProduction,
		InterviewScore:     req.InterviewScore,
		CulturalFitScore:   req.CulturalFitScore,
		TechnicalScore:     req.TechnicalScore,
		OverallScore:       overallScore,
		Status:             int32(req.Status),
		FunFact:            util.ToNullString(funFact),
		Availability:       util.ToNullString(availability),
		SalaryExpectation:  util.ToNullString(salaryExpectation),
	})

	if err != nil {
		// Check for unique constraint violation on email
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && strings.Contains(pqErr.Constraint, "email") {
				s.logger.Error("email already exists", zap.String("email", req.Email))
				return nil, status.Errorf(codes.AlreadyExists, "email address already exists: %s", req.Email)
			}
		}
		s.logger.Error("failed to create applicant", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create applicant: %v", err)
	}

	s.logger.Info("applicant created with calculated score",
		zap.String("name", applicant.Name),
		zap.Float64("overall_score", overallScore),
		zap.Int32("status", applicant.Status),
	)

	return &applicantsv1.CreateApplicantResponse{
		Applicant: util.DbApplicantToProto(&applicant),
	}, nil
}
