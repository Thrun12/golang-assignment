package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
	"github.com/Thrun12/golang-assignment/internal/util"
)

// mockQuerier is a mock implementation of sqlc.Querier for testing
type mockQuerier struct {
	createFunc  func(ctx context.Context, params sqlc.CreateApplicantParams) (sqlc.Applicant, error)
	getFunc     func(ctx context.Context, id int64) (sqlc.Applicant, error)
	listFunc    func(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error)
	countFunc   func(ctx context.Context, params sqlc.CountApplicantsParams) (int64, error)
	updateFunc  func(ctx context.Context, params sqlc.UpdateApplicantParams) (sqlc.Applicant, error)
	deleteFunc  func(ctx context.Context, id int64) error
	getBestFunc func(ctx context.Context) (sqlc.Applicant, error)
}

func (m *mockQuerier) CreateApplicant(ctx context.Context, params sqlc.CreateApplicantParams) (sqlc.Applicant, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, params)
	}
	return sqlc.Applicant{}, errors.New("createFunc not implemented")
}

func (m *mockQuerier) GetApplicant(ctx context.Context, id int64) (sqlc.Applicant, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return sqlc.Applicant{}, errors.New("getFunc not implemented")
}

func (m *mockQuerier) GetApplicantByEmail(ctx context.Context, email string) (sqlc.Applicant, error) {
	return sqlc.Applicant{}, errors.New("not implemented")
}

func (m *mockQuerier) ListApplicants(ctx context.Context, params sqlc.ListApplicantsParams) ([]sqlc.Applicant, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, params)
	}
	return nil, errors.New("listFunc not implemented")
}

func (m *mockQuerier) CountApplicants(ctx context.Context, params sqlc.CountApplicantsParams) (int64, error) {
	if m.countFunc != nil {
		return m.countFunc(ctx, params)
	}
	return 0, errors.New("countFunc not implemented")
}

func (m *mockQuerier) UpdateApplicant(ctx context.Context, params sqlc.UpdateApplicantParams) (sqlc.Applicant, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, params)
	}
	return sqlc.Applicant{}, errors.New("updateFunc not implemented")
}

func (m *mockQuerier) DeleteApplicant(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("deleteFunc not implemented")
}

func (m *mockQuerier) GetTopApplicantsByPosition(ctx context.Context, params sqlc.GetTopApplicantsByPositionParams) ([]sqlc.Applicant, error) {
	return nil, errors.New("not implemented")
}

func (m *mockQuerier) GetBestApplicant(ctx context.Context) (sqlc.Applicant, error) {
	if m.getBestFunc != nil {
		return m.getBestFunc(ctx)
	}
	return sqlc.Applicant{}, errors.New("getBestFunc not implemented")
}

func (m *mockQuerier) UpdateApplicantScore(ctx context.Context, params sqlc.UpdateApplicantScoreParams) (sqlc.Applicant, error) {
	return sqlc.Applicant{}, errors.New("not implemented")
}

func (m *mockQuerier) DeleteAllApplicants(ctx context.Context) error {
	return errors.New("not implemented")
}

func (m *mockQuerier) GetApplicantStats(ctx context.Context) (sqlc.GetApplicantStatsRow, error) {
	return sqlc.GetApplicantStatsRow{}, errors.New("not implemented")
}

func TestCreateApplicantRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		input       *applicantsv1.CreateApplicantRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid input",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "Senior Developer",
				YearsExperience:  5,
				GithubStars:      100,
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: false,
		},
		{
			name: "Empty name",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "",
				Email:            "john@example.com",
				Position:         "Developer",
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "Name too short",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "J",
				Email:            "john@example.com",
				Position:         "Developer",
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "name must be at least 2 characters",
		},
		{
			name: "Empty email",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "",
				Position:         "Developer",
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "email is required",
		},
		{
			name: "Invalid email format",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "not-an-email",
				Position:         "Developer",
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "email must be a valid email address",
		},
		{
			name: "Empty position",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "",
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "position is required",
		},
		{
			name: "Negative years of experience",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "Developer",
				YearsExperience:  -5,
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "years_experience must be positive",
		},
		{
			name: "Negative github stars",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "Developer",
				GithubStars:      -100,
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "github_stars must be positive",
		},
		{
			name: "Interview score below 0",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "Developer",
				InterviewScore:   -10.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "interview_score must be between 0 and 100",
		},
		{
			name: "Interview score above 100",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "Developer",
				InterviewScore:   110.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "interview_score must be between 0 and 100",
		},
		{
			name: "Cultural fit score out of range",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "Developer",
				InterviewScore:   85.0,
				CulturalFitScore: 150.0,
				TechnicalScore:   88.0,
			},
			expectError: true,
			errorMsg:    "cultural_fit_score must be between 0 and 100",
		},
		{
			name: "Technical score out of range",
			input: &applicantsv1.CreateApplicantRequest{
				Name:             "John Doe",
				Email:            "john@example.com",
				Position:         "Developer",
				InterviewScore:   85.0,
				CulturalFitScore: 90.0,
				TechnicalScore:   -5.0,
			},
			expectError: true,
			errorMsg:    "technical_score must be between 0 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.ValidateApplicant(
				tt.input.Name,
				tt.input.Email,
				tt.input.Position,
				tt.input.YearsExperience,
				tt.input.GithubStars,
				tt.input.InterviewScore,
				tt.input.CulturalFitScore,
				tt.input.TechnicalScore,
				false, // not an update
				0,     // id not needed for create
			)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestCreateApplicant(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Successful creation", func(t *testing.T) {
		mockQ := &mockQuerier{
			createFunc: func(ctx context.Context, params sqlc.CreateApplicantParams) (sqlc.Applicant, error) {
				// Verify the params are correct
				if params.Name != "Jane Doe" {
					t.Errorf("Expected name 'Jane Doe', got '%s'", params.Name)
				}
				if params.Email != "jane@example.com" {
					t.Errorf("Expected email 'jane@example.com', got '%s'", params.Email)
				}
				// Verify overall score was calculated
				if params.OverallScore <= 0 {
					t.Errorf("Expected overall score to be calculated, got %f", params.OverallScore)
				}

				return sqlc.Applicant{
					ID:                 1,
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
					FunFact:            sql.NullString{},
					Availability:       sql.NullString{},
					SalaryExpectation:  sql.NullString{},
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				}, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		req := &applicantsv1.CreateApplicantRequest{
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
			Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_APPLIED,
		}

		resp, err := service.CreateApplicant(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Applicant.Name != "Jane Doe" {
			t.Errorf("Expected name 'Jane Doe', got '%s'", resp.Applicant.Name)
		}
		if resp.Applicant.OverallScore <= 0 {
			t.Errorf("Expected overall score to be > 0, got %f", resp.Applicant.OverallScore)
		}
	})

	t.Run("Validation failure", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		req := &applicantsv1.CreateApplicantRequest{
			Name:  "", // Invalid: empty name
			Email: "jane@example.com",
		}

		resp, err := service.CreateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.InvalidArgument {
			t.Error("Expected InvalidArgument status code")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockQ := &mockQuerier{
			createFunc: func(ctx context.Context, params sqlc.CreateApplicantParams) (sqlc.Applicant, error) {
				return sqlc.Applicant{}, errors.New("database error")
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		req := &applicantsv1.CreateApplicantRequest{
			Name:             "Jane Doe",
			Email:            "jane@example.com",
			Position:         "Developer",
			YearsExperience:  5,
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.CreateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected error from repository, got nil")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Internal {
			t.Errorf("Expected Internal status code, got %v", st.Code())
		}
		if resp != nil {
			t.Error("Expected nil response on repository error")
		}
	})

	t.Run("Email already exists", func(t *testing.T) {
		mockQ := &mockQuerier{
			createFunc: func(ctx context.Context, params sqlc.CreateApplicantParams) (sqlc.Applicant, error) {
				// Simulate unique constraint violation
				pqErr := &pq.Error{
					Code:       "23505",
					Constraint: "applicants_email_key",
				}
				return sqlc.Applicant{}, pqErr
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		req := &applicantsv1.CreateApplicantRequest{
			Name:             "Jane Doe",
			Email:            "existing@example.com",
			Position:         "Developer",
			YearsExperience:  5,
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.CreateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected error for duplicate email, got nil")
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.AlreadyExists {
			t.Errorf("Expected AlreadyExists status code, got %v", st.Code())
		}
		if resp != nil {
			t.Error("Expected nil response on email conflict")
		}
	})

	t.Run("Validation - name too long", func(t *testing.T) {
		service := &ApplicantService{
			queries: &mockQuerier{},
			logger:  logger,
		}

		// Name with more than 255 characters
		longName := strings.Repeat("a", 256)

		req := &applicantsv1.CreateApplicantRequest{
			Name:             longName,
			Email:            "jane@example.com",
			Position:         "Developer",
			InterviewScore:   85.0,
			CulturalFitScore: 90.0,
			TechnicalScore:   88.0,
		}

		resp, err := service.CreateApplicant(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if resp != nil {
			t.Error("Expected nil response on validation error")
		}
	})

	t.Run("With optional fields populated", func(t *testing.T) {
		mockQ := &mockQuerier{
			createFunc: func(ctx context.Context, params sqlc.CreateApplicantParams) (sqlc.Applicant, error) {
				// Verify optional fields are set
				if !params.FunFact.Valid || params.FunFact.String != "Loves coding" {
					t.Errorf("Expected fun fact to be set")
				}
				if !params.Availability.Valid || params.Availability.String != "Immediate" {
					t.Errorf("Expected availability to be set")
				}
				if !params.SalaryExpectation.Valid || params.SalaryExpectation.String != "100k" {
					t.Errorf("Expected salary expectation to be set")
				}

				return sqlc.Applicant{
					ID:                1,
					Name:              params.Name,
					Email:             params.Email,
					FunFact:           params.FunFact,
					Availability:      params.Availability,
					SalaryExpectation: params.SalaryExpectation,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}, nil
			},
		}

		service := &ApplicantService{
			queries: mockQ,
			logger:  logger,
		}

		req := &applicantsv1.CreateApplicantRequest{
			Name:              "Jane Doe",
			Email:             "jane@example.com",
			Position:          "Developer",
			YearsExperience:   5,
			InterviewScore:    85.0,
			CulturalFitScore:  90.0,
			TechnicalScore:    88.0,
			FunFact:           "Loves coding",
			Availability:      "Immediate",
			SalaryExpectation: "100k",
		}

		resp, err := service.CreateApplicant(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if resp.Applicant.FunFact != "Loves coding" {
			t.Errorf("Expected fun fact to be preserved")
		}
	})
}
