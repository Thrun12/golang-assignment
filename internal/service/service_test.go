package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
)

func TestGetApplicant(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Successful retrieval", func(t *testing.T) {
		expectedApplicant := sqlc.Applicant{
			ID:                 1,
			Name:               "Jane Doe",
			Email:              "jane@example.com",
			Position:           "Senior Developer",
			YearsExperience:    5,
			Skills:             []string{"Go", "Kubernetes"},
			GithubStars:        200,
			CanExitVim:         true,
			KnowsGo:            true,
			DebugsInProduction: false,
			InterviewScore:     85.0,
			CulturalFitScore:   90.0,
			TechnicalScore:     88.0,
			OverallScore:       87.5,
			Status:             1,
			FunFact:            sql.NullString{},
			Availability:       sql.NullString{},
			SalaryExpectation:  sql.NullString{},
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		mockQ := &mockQuerier{
			getFunc: func(ctx context.Context, id int64) (sqlc.Applicant, error) {
				if id != 1 {
					t.Errorf("Expected ID 1, got %d", id)
				}
				return expectedApplicant, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.GetApplicant(ctx, &applicantsv1.GetApplicantRequest{Id: 1})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Applicant.Id != expectedApplicant.ID {
			t.Errorf("Expected ID %d, got %d", expectedApplicant.ID, resp.Applicant.Id)
		}
		if resp.Applicant.Name != expectedApplicant.Name {
			t.Errorf("Expected name %s, got %s", expectedApplicant.Name, resp.Applicant.Name)
		}
	})

	t.Run("Applicant not found", func(t *testing.T) {
		mockQ := &mockQuerier{
			getFunc: func(ctx context.Context, id int64) (sqlc.Applicant, error) {
				return sqlc.Applicant{}, errors.New("applicant not found")
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.GetApplicant(ctx, &applicantsv1.GetApplicantRequest{Id: 999})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on error")
		}
	})
}

func TestListApplicants(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Successful listing with defaults", func(t *testing.T) {
		expectedApplicants := []sqlc.Applicant{
			{
				ID:           1,
				Name:         "Jane Doe",
				Email:        "jane@example.com",
				Position:     "Senior Developer",
				OverallScore: 87.5,
			},
			{
				ID:           2,
				Name:         "John Smith",
				Email:        "john@example.com",
				Position:     "Developer",
				OverallScore: 75.0,
			},
		}

		mockQ := &mockQuerier{
			listFunc: func(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error) {
				// Check that limit was set to default (10)
				if params.Limit != 10 {
					t.Errorf("Expected default limit 10, got %d", params.Limit)
				}
				return expectedApplicants, nil
			},
			countFunc: func(ctx context.Context, params sqlc.CountApplicantsParams) (int64, error) {
				return 2, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.ListApplicants(ctx, &applicantsv1.ListApplicantsRequest{
			Limit: 0, // Should default to 10
		})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(resp.Applicants) != 2 {
			t.Errorf("Expected 2 applicants, got %d", len(resp.Applicants))
		}
		if resp.TotalCount != 2 {
			t.Errorf("Expected total count 2, got %d", resp.TotalCount)
		}
		if resp.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", resp.Limit)
		}
	})

	t.Run("Limit capping at 100", func(t *testing.T) {
		mockQ := &mockQuerier{
			listFunc: func(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error) {
				// Check that limit was capped at 100
				if params.Limit != 100 {
					t.Errorf("Expected capped limit 100, got %d", params.Limit)
				}
				return []sqlc.Applicant{}, nil
			},
			countFunc: func(ctx context.Context, params sqlc.CountApplicantsParams) (int64, error) {
				return 0, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.ListApplicants(ctx, &applicantsv1.ListApplicantsRequest{
			Limit: 500, // Should be capped at 100
		})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Limit !=100 {
			t.Errorf("Expected limit capped at 100, got %d", resp.Limit)
		}
	})

	t.Run("Negative offset handling", func(t *testing.T) {
		mockQ := &mockQuerier{
			listFunc: func(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error) {
				// Check that negative offset was corrected to 0
				if params.Offset != 0 {
					t.Errorf("Expected offset 0, got %d", params.Offset)
				}
				return []sqlc.Applicant{}, nil
			},
			countFunc: func(ctx context.Context, params sqlc.CountApplicantsParams) (int64, error) {
				return 0, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.ListApplicants(ctx, &applicantsv1.ListApplicantsRequest{
			Limit:  10,
			Offset: -5, // Should be corrected to 0
		})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Offset !=0 {
			t.Errorf("Expected offset 0, got %d", resp.Offset)
		}
	})

	t.Run("Filtering by position", func(t *testing.T) {
		mockQ := &mockQuerier{
			listFunc: func(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error) {
				if params.Position != "Senior Developer" {
					t.Errorf("Expected position filter 'Senior Developer', got '%s'", params.Position)
				}
				return []sqlc.Applicant{
					{ID: 1, Position: "Senior Developer"},
				}, nil
			},
			countFunc: func(ctx context.Context, params sqlc.CountApplicantsParams) (int64, error) {
				if params.Position != "Senior Developer" {
					t.Errorf("Expected position filter 'Senior Developer', got '%s'", params.Position)
				}
				return 1, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.ListApplicants(ctx, &applicantsv1.ListApplicantsRequest{
			Limit:    10,
			Position: "Senior Developer",
		})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(resp.Applicants) !=1 {
			t.Errorf("Expected 1 applicant, got %d", len(resp.Applicants))
		}
	})

	t.Run("Repository error on list", func(t *testing.T) {
		mockQ := &mockQuerier{
			listFunc: func(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error) {
				return nil, errors.New("database error")
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.ListApplicants(ctx, &applicantsv1.ListApplicantsRequest{
			Limit: 10,
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on error")
		}
	})

	t.Run("Repository error on count", func(t *testing.T) {
		mockQ := &mockQuerier{
			listFunc: func(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error) {
				return []sqlc.Applicant{}, nil
			},
			countFunc: func(ctx context.Context, params sqlc.CountApplicantsParams) (int64, error) {
				return 0, errors.New("database error")
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.ListApplicants(ctx, &applicantsv1.ListApplicantsRequest{
			Limit: 10,
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on error")
		}
	})
}

func TestUpdateApplicant(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Successful update", func(t *testing.T) {
		mockQ := &mockQuerier{
			updateFunc: func(ctx context.Context, params sqlc.UpdateApplicantParams) (sqlc.Applicant, error) {
				if params.ID != 1 {
					t.Errorf("Expected ID 1, got %d", params.ID)
				}
				if params.Name != "Jane Updated" {
					t.Errorf("Expected name 'Jane Updated', got '%s'", params.Name)
				}

				return sqlc.Applicant{
					ID:                 params.ID,
					Name:               params.Name,
					Email:              params.Email,
					Position:           params.Position,
					YearsExperience:    params.YearsExperience,
					Skills:             params.Skills,
					GithubStars:        params.GithubStars,
					CanExitVim:         params.CanExitVim,
					KnowsGo:            params.KnowsGo,
					DebugsInProduction: params.DebugsInProduction,
					InterviewScore:     params.InterviewScore,
					CulturalFitScore:   params.CulturalFitScore,
					TechnicalScore:     params.TechnicalScore,
					OverallScore:       params.OverallScore,
					Status:             params.Status,
					FunFact:            params.FunFact,
					Availability:       params.Availability,
					SalaryExpectation:  params.SalaryExpectation,
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				}, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:                 1,
			Name:               "Jane Updated",
			Email:              "jane.updated@example.com",
			Position:           "Lead Developer",
			YearsExperience:    7,
			Skills:             []string{"Go", "Kubernetes", "gRPC"},
			GithubStars:        300,
			CanExitVim:         true,
			KnowsGo:            true,
			DebugsInProduction: true,
			InterviewScore:     90.0,
			CulturalFitScore:   92.0,
			TechnicalScore:     91.0,
			Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_INTERVIEWED,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Applicant.Name != "Jane Updated" {
			t.Errorf("Expected name 'Jane Updated', got '%s'", resp.Applicant.Name)
		}
		if resp.Applicant.OverallScore <= 0 {
			t.Errorf("Expected overall score to be > 0, got %f", resp.Applicant.OverallScore)
		}
	})

	t.Run("Validation failure - invalid ID", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:    0,
			Name:  "Jane Doe",
			Email: "jane@example.com",
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Validation failure - invalid email", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "not-an-email",
			Position:         "Developer",
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Validation failure - empty position", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "",
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Validation failure - position too short", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "D",
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Validation failure - negative github stars", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "Developer",
			GithubStars:      -10,
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Validation failure - interview score out of range", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "Developer",
			InterviewScore:   150.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Validation failure - cultural fit score out of range", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "Developer",
			InterviewScore:   85.0,
			CulturalFitScore: -10.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Validation failure - technical score out of range", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "Developer",
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   110.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockQ := &mockQuerier{
			updateFunc: func(ctx context.Context, params sqlc.UpdateApplicantParams) (sqlc.Applicant, error) {
				return sqlc.Applicant{}, errors.New("database error")
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		req := &applicantsv1.UpdateApplicantRequest{
			Id:               1,
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "Developer",
			YearsExperience:  5,
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.UpdateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected error from repository, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on repository error")
		}
	})
}

func TestDeleteApplicant(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Successful deletion", func(t *testing.T) {
		mockQ := &mockQuerier{
			deleteFunc: func(ctx context.Context, id int64) error {
				if id != 1 {
					t.Errorf("Expected ID 1, got %d", id)
				}
				return nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.DeleteApplicant(ctx, &applicantsv1.DeleteApplicantRequest{Id: 1})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if !resp.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("Validation failure - invalid ID", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		resp, err := service.DeleteApplicant(ctx, &applicantsv1.DeleteApplicantRequest{Id: 0})
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockQ := &mockQuerier{
			deleteFunc: func(ctx context.Context, id int64) error {
				return errors.New("database error")
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.DeleteApplicant(ctx, &applicantsv1.DeleteApplicantRequest{Id: 1})
		if err == nil {
			t.Fatal("Expected error from repository, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on repository error")
		}
	})
}

func TestGetBestApplicant(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Successful retrieval", func(t *testing.T) {
		mockQ := &mockQuerier{
			getBestFunc: func(ctx context.Context) (sqlc.Applicant, error) {
				return sqlc.Applicant{
					ID:           1,
					Name:         "Jonathan Søholm-Boesen",
					Email:        "jonathan@example.com",
					Position:     "Senior Go Developer",
					OverallScore: 99.5,
					Status:       1,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.GetBestApplicant(ctx, &applicantsv1.GetBestApplicantRequest{})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Applicant == nil {
			t.Fatal("Expected applicant, got nil")
		}

		if resp.Reason == "" {
			t.Error("Expected reason to be populated")
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockQ := &mockQuerier{
			getBestFunc: func(ctx context.Context) (sqlc.Applicant, error) {
				return sqlc.Applicant{}, errors.New("database error")
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.GetBestApplicant(ctx, &applicantsv1.GetBestApplicantRequest{})
		if err == nil {
			t.Fatal("Expected error from repository, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on repository error")
		}
	})

	t.Run("GenerateBestApplicantReason - Jonathan special case", func(t *testing.T) {
		mockQ := &mockQuerier{
			getBestFunc: func(ctx context.Context) (sqlc.Applicant, error) {
				return sqlc.Applicant{
					ID:           1,
					Name:         "Jonathan Søholm-Boesen",
					OverallScore: 99.5,
				}, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.GetBestApplicant(ctx, &applicantsv1.GetBestApplicantRequest{})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should contain special keywords
		if !containsAny(resp.Reason, []string{"Danish", "time travel"}) {
			t.Errorf("Expected special reason for Jonathan, got: %s", resp.Reason)
		}
	})

	t.Run("GenerateBestApplicantReason - generic case", func(t *testing.T) {
		mockQ := &mockQuerier{
			getBestFunc: func(ctx context.Context) (sqlc.Applicant, error) {
				return sqlc.Applicant{
					ID:           2,
					Name:         "Jane Doe",
					OverallScore: 95.5,
				}, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		resp, err := service.GetBestApplicant(ctx, &applicantsv1.GetBestApplicantRequest{})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should contain the score
		if !containsAny(resp.Reason, []string{"95.50"}) {
			t.Errorf("Expected reason to contain score, got: %s", resp.Reason)
		}
	})
}

// Helper function to check if string contains any of the substrings
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
