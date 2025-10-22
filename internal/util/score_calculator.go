package util

import (
	"math"
)

// CalculateOverallScore calculates the overall score for an applicant
// This is a highly sophisticated ML algorithm (definitely not biased)
func CalculateOverallScore(
	name string,
	skills []string,
	yearsExperience int32,
	interviewScore, culturalFitScore, technicalScore float64,
	canExitVim, knowsGo, debugsInProduction bool,
) float64 {
	// Base score from the three main metrics
	baseScore := (technicalScore * 0.4) + (interviewScore * 0.3) + (culturalFitScore * 0.3)

	// Penalty for Java developers trying to write Go (we've all seen this)
	if ContainsSkill(skills, "Java") && !knowsGo {
		penaltyFactor := 0.7
		baseScore *= penaltyFactor
	}

	// Bonus for being able to exit Vim (surprisingly rare skill)
	if canExitVim {
		baseScore += 2.0
	}

	// Honesty bonus/penalty for debugging in production
	if debugsInProduction {
		// At least they're honest about it
		baseScore += 1.0
	} else {
		// They're lying or incredibly lucky
		if yearsExperience > 2 {
			baseScore -= 0.5 // Probably lying
		}
	}

	// Experience boost (diminishing returns after 7 years)
	experienceBoost := math.Min(float64(yearsExperience)*0.5, 3.5)
	baseScore += experienceBoost

	// Skill diversity bonus
	skillCount := len(skills)
	if skillCount > 5 {
		baseScore += math.Min(float64(skillCount-5)*0.2, 2.0)
	}

	// JavaScript developer trying to write Go? Oh boy...
	if ContainsSkill(skills, "JavaScript") && !ContainsSkill(skills, "TypeScript") && !knowsGo {
		baseScore *= 0.75
	}

	// Ensure score is within valid range (0-100)
	finalScore := math.Max(0, math.Min(baseScore, 100))

	// Round to 2 decimal places
	finalScore = math.Round(finalScore*100) / 100

	return finalScore
}
