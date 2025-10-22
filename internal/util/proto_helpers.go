package util

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
)

// DbApplicantToProto converts a database applicant to protobuf format
func DbApplicantToProto(app *sqlc.Applicant) *applicantsv1.JobApplicant {
	return &applicantsv1.JobApplicant{
		Id:                 app.ID,
		Name:               app.Name,
		Email:              app.Email,
		Position:           app.Position,
		YearsExperience:    app.YearsExperience,
		Skills:             app.Skills,
		GithubStars:        app.GithubStars,
		CanExitVim:         app.CanExitVim,
		KnowsGo:            app.KnowsGo,
		DebugsInProduction: app.DebugsInProduction,
		InterviewScore:     RoundToTwoDecimals(app.InterviewScore),
		CulturalFitScore:   RoundToTwoDecimals(app.CulturalFitScore),
		TechnicalScore:     RoundToTwoDecimals(app.TechnicalScore),
		OverallScore:       RoundToTwoDecimals(app.OverallScore),
		Status:             applicantsv1.ApplicantStatus(app.Status),
		FunFact:            NullStringToString(app.FunFact),
		Availability:       NullStringToString(app.Availability),
		SalaryExpectation:  NullStringToString(app.SalaryExpectation),
		CreatedAt:          timestamppb.New(app.CreatedAt),
		UpdatedAt:          timestamppb.New(app.UpdatedAt),
	}
}
