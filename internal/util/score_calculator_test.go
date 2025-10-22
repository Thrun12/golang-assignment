package util

import (
	"testing"
)

func TestCalculateOverallScore(t *testing.T) {
	tests := []struct {
		name               string
		applicantName      string
		skills             []string
		yearsExperience    int32
		interviewScore     float64
		culturalFitScore   float64
		technicalScore     float64
		canExitVim         bool
		knowsGo            bool
		debugsInProduction bool
		expectedScore      float64
		description        string
	}{
		{
			name:               "High scores with all bonuses",
			applicantName:      "Jonathan Søholm-Boesen",
			skills:             []string{"Go", "gRPC", "Kubernetes", "Being Modest", "Microservices", "Time Travel"},
			yearsExperience:    10,
			interviewScore:     99.8,
			culturalFitScore:   99.9,
			technicalScore:     99.7,
			canExitVim:         true,
			knowsGo:            true,
			debugsInProduction: false,
			expectedScore:      100.00, // Will be capped at 100
			description:        "High input scores with bonuses should result in top score",
		},
		{
			name:               "Base score calculation (40/30/30 split)",
			applicantName:      "Alice Johnson",
			skills:             []string{"Go", "Kubernetes"},
			yearsExperience:    0,
			interviewScore:     60.0,
			culturalFitScore:   70.0,
			technicalScore:     80.0,
			canExitVim:         false,
			knowsGo:            true,
			debugsInProduction: false,
			expectedScore:      71.0, // (80*0.4 + 60*0.3 + 70*0.3) = 32+18+21 = 71
			description:        "Should correctly calculate weighted base score",
		},
		{
			name:               "Vim exit bonus applied",
			applicantName:      "Bob Developer",
			skills:             []string{"Go", "Docker"},
			yearsExperience:    5,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         true,
			knowsGo:            true,
			debugsInProduction: false,
			expectedScore:      74.0, // Base: 70, + 2 (vim) + 2.5 (exp) - 0.5 (honesty penalty) = 74.0
			description:        "Should add 2 points for Vim exit capability",
		},
		{
			name:               "Debugging in production honesty bonus",
			applicantName:      "Charlie Honest",
			skills:             []string{"Go", "Python"},
			yearsExperience:    5,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            true,
			debugsInProduction: true,
			expectedScore:      73.5, // Base: 70, + 1 (honesty) + 2.5 (exp) = 73.5
			description:        "Should add 1 point for admitting to debugging in production",
		},
		{
			name:               "Honesty penalty for claiming never to debug in production (experienced dev)",
			applicantName:      "Dave Liar",
			skills:             []string{"Go"},
			yearsExperience:    5,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            true,
			debugsInProduction: false,
			expectedScore:      72.0, // Base: 70, - 0.5 (lying) + 2.5 (exp) = 72.0
			description:        "Should subtract 0.5 points for likely lying about not debugging in production",
		},
		{
			name:               "No honesty penalty for junior dev",
			applicantName:      "Emma Junior",
			skills:             []string{"Go"},
			yearsExperience:    1,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            true,
			debugsInProduction: false,
			expectedScore:      70.5, // Base: 70, + 0.5 (exp), no penalty for junior
			description:        "Juniors (<=2 years) shouldn't get honesty penalty",
		},
		{
			name:               "Experience boost (diminishing returns)",
			applicantName:      "Frank Senior",
			skills:             []string{"Go"},
			yearsExperience:    10,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            true,
			debugsInProduction: true,
			expectedScore:      74.5, // Base: 70, + 1 (honesty) + 3.5 (exp capped at 7*0.5) = 74.5
			description:        "Experience bonus should cap at 3.5 points (7 years)",
		},
		{
			name:               "Skill diversity bonus",
			applicantName:      "Grace Polyglot",
			skills:             []string{"Go", "Python", "JavaScript", "Rust", "Java", "C++", "TypeScript", "Kotlin"},
			yearsExperience:    5,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            true,
			debugsInProduction: true,
			expectedScore:      74.1, // Base: 70, + 1 (honesty) + 2.5 (exp) + 0.6 (3 skills * 0.2) = 74.1
			description:        "Should add bonus for skill diversity (0.2 per skill beyond 5, max 2.0)",
		},
		{
			name:               "Java developer without Go knowledge penalty",
			applicantName:      "Harry Java",
			skills:             []string{"Java", "Spring"},
			yearsExperience:    5,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            false,
			debugsInProduction: true,
			expectedScore:      52.5, // Base: 70 * 0.7 (Java penalty) = 49, + 1 (honesty) + 2.5 (exp) = 52.5
			description:        "Java devs without Go should get 0.7 multiplier penalty",
		},
		{
			name:               "JavaScript developer without TypeScript or Go penalty",
			applicantName:      "Iris JavaScript",
			skills:             []string{"JavaScript", "React"},
			yearsExperience:    5,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            false,
			debugsInProduction: true,
			expectedScore:      55.12, // Base: 70, + 1 (honesty) + 2.5 (exp) = 73.5 * 0.75 (JS penalty) = 55.125
			description:        "JS devs without TS or Go should get 0.75 multiplier penalty",
		},
		{
			name:               "JavaScript developer with TypeScript (no penalty)",
			applicantName:      "Jack TypeScript",
			skills:             []string{"JavaScript", "TypeScript", "React"},
			yearsExperience:    5,
			interviewScore:     70.0,
			culturalFitScore:   70.0,
			technicalScore:     70.0,
			canExitVim:         false,
			knowsGo:            false,
			debugsInProduction: true,
			expectedScore:      73.5, // Base: 70, + 1 (honesty) + 2.5 (exp) = 73.5
			description:        "JS devs with TypeScript should not get penalty",
		},
		{
			name:               "Perfect score capped at 100",
			applicantName:      "Karen Perfect",
			skills:             []string{"Go", "Kubernetes", "gRPC", "PostgreSQL", "Docker", "AWS", "Terraform", "React"},
			yearsExperience:    7,
			interviewScore:     100.0,
			culturalFitScore:   100.0,
			technicalScore:     100.0,
			canExitVim:         true,
			knowsGo:            true,
			debugsInProduction: true,
			expectedScore:      100.0, // Should cap at 100
			description:        "Score should cap at 100",
		},
		{
			name:               "All bonuses combined",
			applicantName:      "Linda AllBonuses",
			skills:             []string{"Go", "Python", "JavaScript", "TypeScript", "Rust", "Java", "C++", "Kotlin", "Ruby", "PHP", "Swift", "Kotlin"},
			yearsExperience:    10,
			interviewScore:     90.0,
			culturalFitScore:   90.0,
			technicalScore:     90.0,
			canExitVim:         true,
			knowsGo:            true,
			debugsInProduction: true,
			expectedScore:      97.9, // Base: 90, + 2 (vim) + 1 (honesty) + 3.5 (exp) + 1.4 (7 skills * 0.2) = 97.9
			description:        "Multiple bonuses should stack correctly",
		},
		{
			name:               "Minimum score floor (negative scores should be capped)",
			applicantName:      "Mike Terrible",
			skills:             []string{"Java"},
			yearsExperience:    0,
			interviewScore:     0.0,
			culturalFitScore:   0.0,
			technicalScore:     0.0,
			canExitVim:         false,
			knowsGo:            false,
			debugsInProduction: false,
			expectedScore:      0.0, // Base: 0 * 0.7 (Java penalty) = 0, floored at 0
			description:        "Score should never go below 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateOverallScore(
				tt.applicantName,
				tt.skills,
				tt.yearsExperience,
				tt.interviewScore,
				tt.culturalFitScore,
				tt.technicalScore,
				tt.canExitVim,
				tt.knowsGo,
				tt.debugsInProduction,
			)

			// Use a small delta for floating point comparison
			delta := 0.1
			if score < tt.expectedScore-delta || score > tt.expectedScore+delta {
				t.Errorf("%s\nExpected score: %.2f, got: %.2f\nDifference: %.2f",
					tt.description, tt.expectedScore, score, score-tt.expectedScore)
			}
		})
	}
}

func TestCalculateOverallScore_EdgeCases(t *testing.T) {
	t.Run("Empty skills array", func(t *testing.T) {
		score := CalculateOverallScore(
			"Test User",
			[]string{},
			5,
			70.0,
			70.0,
			70.0,
			false,
			true,
			true,
		)

		if score < 0 {
			t.Errorf("Score should not be negative with empty skills, got: %.2f", score)
		}
	})

	t.Run("Nil skills array", func(t *testing.T) {
		score := CalculateOverallScore(
			"Test User",
			nil,
			5,
			70.0,
			70.0,
			70.0,
			false,
			true,
			true,
		)

		if score < 0 {
			t.Errorf("Score should not be negative with nil skills, got: %.2f", score)
		}
	})

	t.Run("Negative years of experience", func(t *testing.T) {
		score := CalculateOverallScore(
			"Test User",
			[]string{"Go"},
			-5,
			70.0,
			70.0,
			70.0,
			false,
			true,
			true,
		)

		// Should handle gracefully (negative experience would give negative boost, but still valid)
		if score > 100 || score < 0 {
			t.Errorf("Score should be in valid range [0, 100], got: %.2f", score)
		}
	})

	t.Run("Score components above 100", func(t *testing.T) {
		score := CalculateOverallScore(
			"Test User",
			[]string{"Go"},
			5,
			150.0, // Invalid but should be handled
			150.0,
			150.0,
			false,
			true,
			true,
		)

		// Score should still be capped at 100
		if score > 100 {
			t.Errorf("Score should be capped at 100, got: %.2f", score)
		}
	})

	t.Run("Unicode in name (Danish characters)", func(t *testing.T) {
		score := CalculateOverallScore(
			"Jørgen Søholm-Hansen",
			[]string{"Go"},
			5,
			70.0,
			70.0,
			70.0,
			false,
			true,
			true,
		)

		// Should not panic and should handle Unicode correctly
		if score < 0 || score > 100 {
			t.Errorf("Should handle Unicode characters, got score: %.2f", score)
		}
	})
}
