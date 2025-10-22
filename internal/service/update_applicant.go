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

// UpdateApplicant updates an existing applicant and recalculates score
func (s *ApplicantService) UpdateApplicant(ctx context.Context, req *applicantsv1.UpdateApplicantRequest) (*applicantsv1.UpdateApplicantResponse, error) {
	// Validate input
	if err := util.ValidateApplicant(req.Name, req.Email, req.Position, req.YearsExperience, req.GithubStars, req.InterviewScore, req.CulturalFitScore, req.TechnicalScore, true, req.Id); err != nil {
		s.logger.Debug("validation failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	s.logger.Debug("updating applicant", zap.Int64("id", req.Id))

	// Recalculate overall score
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

	applicant, err := s.queries.UpdateApplicant(ctx, sqlc.UpdateApplicantParams{
		ID:                 req.Id,
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
		s.logger.Error("failed to update applicant", zap.Int64("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update applicant: %v", err)
	}

	return &applicantsv1.UpdateApplicantResponse{
		Applicant: util.DbApplicantToProto(&applicant),
	}, nil
}
