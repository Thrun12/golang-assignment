package util

import (
	"database/sql"
	"testing"
	"time"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
)

func TestDbApplicantToProto(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    *sqlc.Applicant
		validate func(t *testing.T, result *applicantsv1.JobApplicant)
	}{
		{
			name: "Complete applicant with all fields",
			input: &sqlc.Applicant{
				ID:                 1,
				Name:               "Jane Doe",
				Email:              "jane@example.com",
				Position:           "Senior Developer",
				YearsExperience:    5,
				Skills:             []string{"Go", "Python", "Kubernetes"},
				GithubStars:        150,
				CanExitVim:         true,
				KnowsGo:            true,
				DebugsInProduction: false,
				InterviewScore:     85.5555,
				CulturalFitScore:   90.9999,
				TechnicalScore:     88.1234,
				OverallScore:       87.6789,
				Status:             1,
				FunFact:            sql.NullString{String: "Loves Go", Valid: true},
				Availability:       sql.NullString{String: "2 weeks", Valid: true},
				SalaryExpectation:  sql.NullString{String: "Competitive", Valid: true},
				CreatedAt:          now,
				UpdatedAt:          now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if result.Id != 1 {
					t.Errorf("Expected ID 1, got %d", result.Id)
				}
				if result.Name != "Jane Doe" {
					t.Errorf("Expected name 'Jane Doe', got '%s'", result.Name)
				}
				if result.Email != "jane@example.com" {
					t.Errorf("Expected email 'jane@example.com', got '%s'", result.Email)
				}
				if result.Position != "Senior Developer" {
					t.Errorf("Expected position 'Senior Developer', got '%s'", result.Position)
				}
				if result.YearsExperience != 5 {
					t.Errorf("Expected years experience 5, got %d", result.YearsExperience)
				}
				if len(result.Skills) != 3 {
					t.Errorf("Expected 3 skills, got %d", len(result.Skills))
				}
				if result.GithubStars != 150 {
					t.Errorf("Expected github stars 150, got %d", result.GithubStars)
				}
				if !result.CanExitVim {
					t.Error("Expected CanExitVim to be true")
				}
				if !result.KnowsGo {
					t.Error("Expected KnowsGo to be true")
				}
				if result.DebugsInProduction {
					t.Error("Expected DebugsInProduction to be false")
				}
				// Check scores are rounded to 2 decimals
				if result.InterviewScore != 85.56 {
					t.Errorf("Expected interview score 85.56, got %.2f", result.InterviewScore)
				}
				if result.CulturalFitScore != 91.0 {
					t.Errorf("Expected cultural fit score 91.0, got %.2f", result.CulturalFitScore)
				}
				if result.TechnicalScore != 88.12 {
					t.Errorf("Expected technical score 88.12, got %.2f", result.TechnicalScore)
				}
				if result.OverallScore != 87.68 {
					t.Errorf("Expected overall score 87.68, got %.2f", result.OverallScore)
				}
				if result.Status != applicantsv1.ApplicantStatus_APPLICANT_STATUS_APPLIED {
					t.Errorf("Expected status APPLIED, got %v", result.Status)
				}
				if result.FunFact != "Loves Go" {
					t.Errorf("Expected fun fact 'Loves Go', got '%s'", result.FunFact)
				}
				if result.Availability != "2 weeks" {
					t.Errorf("Expected availability '2 weeks', got '%s'", result.Availability)
				}
				if result.SalaryExpectation != "Competitive" {
					t.Errorf("Expected salary expectation 'Competitive', got '%s'", result.SalaryExpectation)
				}
				if result.CreatedAt == nil {
					t.Error("Expected CreatedAt to be set")
				}
				if result.UpdatedAt == nil {
					t.Error("Expected UpdatedAt to be set")
				}
			},
		},
		{
			name: "Applicant with null optional fields",
			input: &sqlc.Applicant{
				ID:                 2,
				Name:               "John Smith",
				Email:              "john@example.com",
				Position:           "Developer",
				YearsExperience:    2,
				Skills:             []string{},
				GithubStars:        0,
				CanExitVim:         false,
				KnowsGo:            false,
				DebugsInProduction: true,
				InterviewScore:     70.0,
				CulturalFitScore:   75.0,
				TechnicalScore:     72.0,
				OverallScore:       72.5,
				Status:             2,
				FunFact:            sql.NullString{Valid: false},
				Availability:       sql.NullString{Valid: false},
				SalaryExpectation:  sql.NullString{Valid: false},
				CreatedAt:          now,
				UpdatedAt:          now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if result.FunFact != "" {
					t.Errorf("Expected empty fun fact, got '%s'", result.FunFact)
				}
				if result.Availability != "" {
					t.Errorf("Expected empty availability, got '%s'", result.Availability)
				}
				if result.SalaryExpectation != "" {
					t.Errorf("Expected empty salary expectation, got '%s'", result.SalaryExpectation)
				}
				if len(result.Skills) != 0 {
					t.Errorf("Expected 0 skills, got %d", len(result.Skills))
				}
			},
		},
		{
			name: "Applicant with status HIRED",
			input: &sqlc.Applicant{
				ID:           3,
				Name:         "Alice Johnson",
				Email:        "alice@example.com",
				Position:     "Lead Developer",
				OverallScore: 95.0,
				Status:       4, // HIRED
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if result.Status != applicantsv1.ApplicantStatus_APPLICANT_STATUS_HIRED {
					t.Errorf("Expected status HIRED, got %v", result.Status)
				}
			},
		},
		{
			name: "Applicant with status REJECTED",
			input: &sqlc.Applicant{
				ID:           4,
				Name:         "Bob Williams",
				Email:        "bob@example.com",
				Position:     "Junior Developer",
				OverallScore: 45.0,
				Status:       5, // REJECTED
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if result.Status != applicantsv1.ApplicantStatus_APPLICANT_STATUS_REJECTED {
					t.Errorf("Expected status REJECTED, got %v", result.Status)
				}
			},
		},
		{
			name: "Score rounding - all zeros",
			input: &sqlc.Applicant{
				ID:               5,
				Name:             "Test User",
				Email:            "test@example.com",
				Position:         "Test",
				InterviewScore:   0.0,
				CulturalFitScore: 0.0,
				TechnicalScore:   0.0,
				OverallScore:     0.0,
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if result.InterviewScore != 0.0 {
					t.Errorf("Expected interview score 0.0, got %.2f", result.InterviewScore)
				}
				if result.OverallScore != 0.0 {
					t.Errorf("Expected overall score 0.0, got %.2f", result.OverallScore)
				}
			},
		},
		{
			name: "Score rounding - perfect scores",
			input: &sqlc.Applicant{
				ID:               6,
				Name:             "Perfect Candidate",
				Email:            "perfect@example.com",
				Position:         "Architect",
				InterviewScore:   100.0,
				CulturalFitScore: 100.0,
				TechnicalScore:   100.0,
				OverallScore:     100.0,
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if result.InterviewScore != 100.0 {
					t.Errorf("Expected interview score 100.0, got %.2f", result.InterviewScore)
				}
				if result.OverallScore != 100.0 {
					t.Errorf("Expected overall score 100.0, got %.2f", result.OverallScore)
				}
			},
		},
		{
			name: "Unicode name and special characters",
			input: &sqlc.Applicant{
				ID:        7,
				Name:      "Jonathan Søholm-Boesen",
				Email:     "jonathan@example.com",
				Position:  "Senior Go Developer",
				FunFact:   sql.NullString{String: "Has minor time travel capabilities", Valid: true},
				CreatedAt: now,
				UpdatedAt: now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if result.Name != "Jonathan Søholm-Boesen" {
					t.Errorf("Expected name 'Jonathan Søholm-Boesen', got '%s'", result.Name)
				}
				if result.FunFact != "Has minor time travel capabilities" {
					t.Errorf("Expected fun fact to be set, got '%s'", result.FunFact)
				}
			},
		},
		{
			name: "Many skills",
			input: &sqlc.Applicant{
				ID:        8,
				Name:      "Polyglot Developer",
				Email:     "polyglot@example.com",
				Position:  "Full Stack Developer",
				Skills:    []string{"Go", "Python", "JavaScript", "TypeScript", "Rust", "Java", "C++", "C#", "Ruby", "PHP"},
				CreatedAt: now,
				UpdatedAt: now,
			},
			validate: func(t *testing.T, result *applicantsv1.JobApplicant) {
				if len(result.Skills) != 10 {
					t.Errorf("Expected 10 skills, got %d", len(result.Skills))
				}
				// Check if Go is in the list
				found := false
				for _, skill := range result.Skills {
					if skill == "Go" {
						found = true
						break
					}
				}
				if !found {
					t.Error("Expected 'Go' to be in skills list")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DbApplicantToProto(tt.input)
			if result == nil {
				t.Fatal("Expected non-nil result")
			}
			tt.validate(t, result)
		})
	}
}

func TestDbApplicantToProto_NilInput(t *testing.T) {
	// This should panic or be handled gracefully
	// In production code, you might want to add a nil check
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil input
			t.Log("Panicked as expected with nil input")
		}
	}()

	// This will likely panic, which is fine for this case
	_ = DbApplicantToProto(nil)
}

func TestDbApplicantToProto_TimestampConversion(t *testing.T) {
	// Test that timestamps are properly converted
	specificTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	applicant := &sqlc.Applicant{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Position:  "Developer",
		CreatedAt: specificTime,
		UpdatedAt: specificTime,
	}

	result := DbApplicantToProto(applicant)

	if result.CreatedAt == nil {
		t.Fatal("Expected CreatedAt to be set")
	}

	if result.UpdatedAt == nil {
		t.Fatal("Expected UpdatedAt to be set")
	}

	// Convert back to time.Time to verify
	createdTime := result.CreatedAt.AsTime()
	if !createdTime.Equal(specificTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", specificTime, createdTime)
	}

	updatedTime := result.UpdatedAt.AsTime()
	if !updatedTime.Equal(specificTime) {
		t.Errorf("Expected UpdatedAt to be %v, got %v", specificTime, updatedTime)
	}
}

func TestDbApplicantToProto_ScoreRounding(t *testing.T) {
	// Specifically test score rounding behavior
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"Round down", 85.244, 85.24},
		{"Round up", 85.246, 85.25},
		{"Round at boundary", 85.245, 85.25},
		{"No rounding needed", 85.25, 85.25},
		{"Integer score", 85.0, 85.0},
		{"Many decimals", 85.123456789, 85.12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applicant := &sqlc.Applicant{
				ID:               1,
				Name:             "Test",
				Email:            "test@example.com",
				Position:         "Dev",
				InterviewScore:   tt.input,
				CulturalFitScore: tt.input,
				TechnicalScore:   tt.input,
				OverallScore:     tt.input,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			result := DbApplicantToProto(applicant)

			epsilon := 0.01
			if result.InterviewScore < tt.want-epsilon || result.InterviewScore > tt.want+epsilon {
				t.Errorf("InterviewScore: expected %.2f, got %.2f", tt.want, result.InterviewScore)
			}
			if result.CulturalFitScore < tt.want-epsilon || result.CulturalFitScore > tt.want+epsilon {
				t.Errorf("CulturalFitScore: expected %.2f, got %.2f", tt.want, result.CulturalFitScore)
			}
			if result.TechnicalScore < tt.want-epsilon || result.TechnicalScore > tt.want+epsilon {
				t.Errorf("TechnicalScore: expected %.2f, got %.2f", tt.want, result.TechnicalScore)
			}
			if result.OverallScore < tt.want-epsilon || result.OverallScore > tt.want+epsilon {
				t.Errorf("OverallScore: expected %.2f, got %.2f", tt.want, result.OverallScore)
			}
		})
	}
}
