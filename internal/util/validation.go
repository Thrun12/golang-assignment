package util

import (
	"fmt"
	"net/mail"
	"strings"
)

// ValidateApplicant validates applicant fields for both create and update requests
func ValidateApplicant(name, email, position string, yearsExperience, githubStars int32, interviewScore, culturalFitScore, technicalScore float64, isUpdate bool, id int64) error {
	// Validate ID for update requests
	if isUpdate && id <= 0 {
		return fmt.Errorf("id must be positive")
	}

	// Validate name
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(name) > 255 {
		return fmt.Errorf("name must be at most 255 characters")
	}

	// Validate email
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email must be a valid email address")
	}

	// Validate position
	if strings.TrimSpace(position) == "" {
		return fmt.Errorf("position is required")
	}
	if len(position) < 2 {
		return fmt.Errorf("position must be at least 2 characters")
	}

	// Validate years experience and github stars
	if yearsExperience < 0 {
		return fmt.Errorf("years_experience must be positive")
	}
	if githubStars < 0 {
		return fmt.Errorf("github_stars must be positive")
	}

	// Validate scores
	if interviewScore < 0 || interviewScore > 100 {
		return fmt.Errorf("interview_score must be between 0 and 100")
	}
	if culturalFitScore < 0 || culturalFitScore > 100 {
		return fmt.Errorf("cultural_fit_score must be between 0 and 100")
	}
	if technicalScore < 0 || technicalScore > 100 {
		return fmt.Errorf("technical_score must be between 0 and 100")
	}

	return nil
}
