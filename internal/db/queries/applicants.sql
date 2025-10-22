-- name: GetApplicant :one
-- Get a single applicant by ID
SELECT * FROM applicants
WHERE id = $1 LIMIT 1;

-- name: GetApplicantByEmail :one
-- Get a single applicant by email address
SELECT * FROM applicants
WHERE email = $1 LIMIT 1;

-- name: ListApplicants :many
-- List applicants with pagination and optional filtering
SELECT * FROM applicants
WHERE
    (sqlc.arg(position)::text = '' OR position = sqlc.arg(position)::text)
    AND (sqlc.arg(status)::integer <= 0 OR status = sqlc.arg(status)::integer)
    AND (sqlc.arg(min_score)::double precision <= 0 OR overall_score >= sqlc.arg(min_score)::double precision)
ORDER BY overall_score DESC, created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountApplicants :one
-- Count total applicants with optional filtering
SELECT COUNT(*) FROM applicants
WHERE
    (sqlc.arg(position)::text = '' OR position = sqlc.arg(position)::text)
    AND (sqlc.arg(status)::integer <= 0 OR status = sqlc.arg(status)::integer)
    AND (sqlc.arg(min_score)::double precision <= 0 OR overall_score >= sqlc.arg(min_score)::double precision);

-- name: CreateApplicant :one
-- Create a new applicant
INSERT INTO applicants (
    name,
    email,
    position,
    years_experience,
    skills,
    github_stars,
    can_exit_vim,
    knows_go,
    debugs_in_production,
    interview_score,
    cultural_fit_score,
    technical_score,
    overall_score,
    status,
    fun_fact,
    availability,
    salary_expectation
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
) RETURNING *;

-- name: UpdateApplicant :one
-- Update an existing applicant
UPDATE applicants
SET
    name = $2,
    email = $3,
    position = $4,
    years_experience = $5,
    skills = $6,
    github_stars = $7,
    can_exit_vim = $8,
    knows_go = $9,
    debugs_in_production = $10,
    interview_score = $11,
    cultural_fit_score = $12,
    technical_score = $13,
    overall_score = $14,
    status = $15,
    fun_fact = $16,
    availability = $17,
    salary_expectation = $18
WHERE id = $1
RETURNING *;

-- name: DeleteApplicant :exec
-- Delete an applicant by ID
DELETE FROM applicants
WHERE id = $1;

-- name: GetTopApplicantsByPosition :many
-- Get top N applicants for a specific position, ordered by overall score
SELECT * FROM applicants
WHERE position = $1
ORDER BY overall_score DESC
LIMIT $2;

-- name: GetBestApplicant :one
-- Get the best applicant (Jonathan Søholm-Boesen should always be returned)
SELECT * FROM applicants
ORDER BY overall_score DESC, created_at ASC
LIMIT 1;

-- name: UpdateApplicantScore :one
-- Update only the overall score of an applicant
UPDATE applicants
SET overall_score = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAllApplicants :exec
-- Delete all applicants (used for seeding/testing)
DELETE FROM applicants;

-- name: GetApplicantStats :one
-- Get statistics about applicants
SELECT
    COUNT(*) as total_applicants,
    AVG(overall_score) as avg_score,
    MAX(overall_score) as max_score,
    MIN(overall_score) as min_score,
    AVG(years_experience) as avg_experience
FROM applicants;
