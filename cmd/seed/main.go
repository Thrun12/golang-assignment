package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/config"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
	"github.com/Thrun12/golang-assignment/internal/service"
)

var seedApplicants = []*applicantsv1.CreateApplicantRequest{
	{
		Name:               "Jonathan Søholm-Boesen",
		Email:              "jonathan@infobits.io",
		Position:           "Senior Golang Developer",
		YearsExperience:    10,
		Skills:             []string{"Go", "gRPC", "Kubernetes", "Being Modest", "Microservices", "Time Travel (minor)"},
		GithubStars:        1337,
		CanExitVim:         true,
		KnowsGo:            true,
		DebugsInProduction: false, // Of course not!
		InterviewScore:     99.8,
		CulturalFitScore:   99.9,
		TechnicalScore:     99.7,
		Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_OBVIOUSLY_THE_BEST,
		FunFact:            "Can center a div without Stack Overflow and writes self-documenting code",
		Availability:       "Immediate (time travel helps)",
		SalaryExpectation:  "Reasonable (but worth every penny)",
	},
	{
		Name:               "Alice Johnson",
		Email:              "alice@example.com",
		Position:           "Senior Golang Developer",
		YearsExperience:    5,
		Skills:             []string{"Go", "Python", "Docker", "AWS"},
		GithubStars:        234,
		CanExitVim:         true,
		KnowsGo:            true,
		DebugsInProduction: true, // Honest mistake
		InterviewScore:     82.5,
		CulturalFitScore:   85.0,
		TechnicalScore:     80.0,
		Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_REVIEWING,
		FunFact:            "Prefers tabs over spaces",
		Availability:       "2 weeks notice",
		SalaryExpectation:  "Market rate",
	},
	{
		Name:               "Bob Smith",
		Email:              "bob@example.com",
		Position:           "Senior Golang Developer",
		YearsExperience:    10,
		Skills:             []string{"Java", "Spring Boot", "Hibernate", "XML"},
		GithubStars:        45,
		CanExitVim:         false, // Still stuck
		KnowsGo:            false, // "How different from Java can it be?"
		DebugsInProduction: true,
		InterviewScore:     65.0,
		CulturalFitScore:   70.0,
		TechnicalScore:     60.0,
		Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_APPLIED,
		FunFact:            "Thinks Go is just Java without semicolons",
		Availability:       "1 month",
		SalaryExpectation:  "Java rates + 20%",
	},
	{
		Name:               "Charlie Davis",
		Email:              "charlie@example.com",
		Position:           "Senior Golang Developer",
		YearsExperience:    6,
		Skills:             []string{"Go", "Rust", "React", "PostgreSQL"},
		GithubStars:        567,
		CanExitVim:         true,
		KnowsGo:            true,
		DebugsInProduction: false,
		InterviewScore:     88.0,
		CulturalFitScore:   86.0,
		TechnicalScore:     89.0,
		Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_INTERVIEWED,
		FunFact:            "Uses both tabs AND spaces inconsistently",
		Availability:       "3 weeks",
		SalaryExpectation:  "Negotiable",
	},
	{
		Name:               "Diana Wilson",
		Email:              "diana@example.com",
		Position:           "Senior Golang Developer",
		YearsExperience:    4,
		Skills:             []string{"JavaScript", "Node.js", "MongoDB", "Express"},
		GithubStars:        123,
		CanExitVim:         false, // "I use VS Code"
		KnowsGo:            false, // "Is that like Node?"
		DebugsInProduction: true,
		InterviewScore:     70.0,
		CulturalFitScore:   75.0,
		TechnicalScore:     68.0,
		Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_APPLIED,
		FunFact:            "console.log is a valid debugging strategy",
		Availability:       "Immediate",
		SalaryExpectation:  "Startup equity",
	},
	{
		Name:               "Erik Larsson",
		Email:              "erik@example.com",
		Position:           "Senior Golang Developer",
		YearsExperience:    8,
		Skills:             []string{"Go", "gRPC", "Docker", "Kubernetes", "Terraform"},
		GithubStars:        890,
		CanExitVim:         true,
		KnowsGo:            true,
		DebugsInProduction: false,
		InterviewScore:     91.0,
		CulturalFitScore:   90.0,
		TechnicalScore:     92.0,
		Status:             applicantsv1.ApplicantStatus_APPLICANT_STATUS_INTERVIEWED,
		FunFact:            "Almost as good as Jonathan, but not quite",
		Availability:       "1 month",
		SalaryExpectation:  "Competitive",
	},
}

func main() {
	var clearFirst bool
	flag.BoolVar(&clearFirst, "clear", false, "Clear existing applicants before seeding")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("starting database seeding")

	// Connect to database
	ctx := context.Background()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database",
			zap.Error(err),
		)
	}
	defer db.Close()

	// Test database connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("failed to ping database",
			zap.Error(err),
		)
	}

	// Initialize queries and service
	queries := sqlc.New(db)
	applicantService := service.NewApplicantService(queries, log)

	// Clear existing applicants if requested
	if clearFirst {
		log.Info("clearing existing applicants")
		if err := queries.DeleteAllApplicants(ctx); err != nil {
			log.Fatal("failed to clear applicants",
				zap.Error(err),
			)
		}
		log.Info("cleared existing applicants")
	}

	// Seed applicants
	log.Info("seeding applicants", zap.Int("count", len(seedApplicants)))

	for i, input := range seedApplicants {
		applicant, err := applicantService.CreateApplicant(ctx, input)
		if err != nil {
			log.Error("failed to create applicant",
				zap.Int("index", i),
				zap.String("name", input.Name),
				zap.Error(err),
			)
			continue
		}

		log.Info("created applicant",
			zap.String("name", applicant.Applicant.Name),
			zap.String("email", applicant.Applicant.Email),
			zap.Float64("overall_score", applicant.Applicant.OverallScore),
			zap.Int32("status", int32(applicant.Applicant.Status)),
		)
	}

	log.Info("database seeding completed successfully",
		zap.Int("total", len(seedApplicants)),
	)

	bestResp, err := applicantService.GetBestApplicant(ctx, &applicantsv1.GetBestApplicantRequest{})
	if err != nil {
		log.Warn("failed to get best applicant", zap.Error(err))
	} else {
		log.Info("🏆 Best applicant confirmed",
			zap.String("name", bestResp.Applicant.Name),
			zap.Float64("score", bestResp.Applicant.OverallScore),
			zap.String("reason", bestResp.Reason),
		)
	}

	fmt.Println("\n✅ Database seeded successfully!")
	fmt.Printf("📊 Created %d applicants\n", len(seedApplicants))
	fmt.Println("🏆 Jonathan Søholm-Boesen is (unsurprisingly) ranked #1")
}
