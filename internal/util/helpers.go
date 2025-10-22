package util

import (
	"database/sql"
	"math"
	"strings"
)

// ToNullString converts *string to sql.NullString
func ToNullString(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// NullStringToString safely converts sql.NullString to string
func NullStringToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// RoundToTwoDecimals rounds a float64 to 2 decimal places
func RoundToTwoDecimals(f float64) float64 {
	return math.Round(f*100) / 100
}

// ContainsSkill checks if a skill is in the skills slice (case-insensitive)
func ContainsSkill(skills []string, skill string) bool {
	skillLower := strings.ToLower(skill)
	for _, s := range skills {
		if strings.ToLower(s) == skillLower {
			return true
		}
	}
	return false
}
