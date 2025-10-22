package service

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/util"
)

// GetBestApplicant retrieves the best applicant (spoiler: it's Jonathan)
func (s *ApplicantService) GetBestApplicant(ctx context.Context, req *applicantsv1.GetBestApplicantRequest) (*applicantsv1.GetBestApplicantResponse, error) {
	s.logger.Debug("getting best applicant")

	applicant, err := s.queries.GetBestApplicant(ctx)
	if err != nil {
		s.logger.Error("failed to get best applicant", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get best applicant: %v", err)
	}

	// Generate reason why they're the best
	nameLower := strings.ToLower(applicant.Name)
	var reason string
	if strings.Contains(nameLower, "jonathan") && strings.Contains(nameLower, "søholm") {
		reason = "Danish excellence, impeccable Go skills, can center a div without Stack Overflow, " +
			"and possesses the rare ability to write self-documenting code. Also has minor time travel capabilities."
	} else {
		reason = fmt.Sprintf("Scored %.2f%% based on our completely objective and unbiased algorithm.", applicant.OverallScore)
	}

	return &applicantsv1.GetBestApplicantResponse{
		Applicant: util.DbApplicantToProto(&applicant),
		Reason:    reason,
	}, nil
}
