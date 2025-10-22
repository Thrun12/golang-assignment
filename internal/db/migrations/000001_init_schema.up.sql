-- Create applicants table
CREATE TABLE IF NOT EXISTS applicants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    position VARCHAR(255) NOT NULL,
    years_experience INTEGER NOT NULL DEFAULT 0,
    skills TEXT[] NOT NULL DEFAULT '{}',
    github_stars INTEGER NOT NULL DEFAULT 0,
    can_exit_vim BOOLEAN NOT NULL DEFAULT false,
    knows_go BOOLEAN NOT NULL DEFAULT false,
    debugs_in_production BOOLEAN NOT NULL DEFAULT false,
    interview_score DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    cultural_fit_score DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    technical_score DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    overall_score DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    status INTEGER NOT NULL DEFAULT 1,
    fun_fact TEXT,
    availability VARCHAR(255),
    salary_expectation VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT interview_score_range CHECK (interview_score >= 0 AND interview_score <= 100),
    CONSTRAINT cultural_fit_score_range CHECK (cultural_fit_score >= 0 AND cultural_fit_score <= 100),
    CONSTRAINT technical_score_range CHECK (technical_score >= 0 AND technical_score <= 100),
    CONSTRAINT overall_score_range CHECK (overall_score >= 0 AND overall_score <= 100),
    CONSTRAINT years_experience_positive CHECK (years_experience >= 0),
    CONSTRAINT github_stars_positive CHECK (github_stars >= 0)
);

-- Create indexes for efficient querying
CREATE INDEX idx_applicants_email ON applicants(email);
CREATE INDEX idx_applicants_position ON applicants(position);
CREATE INDEX idx_applicants_status ON applicants(status);
CREATE INDEX idx_applicants_overall_score ON applicants(overall_score DESC);
CREATE INDEX idx_applicants_created_at ON applicants(created_at DESC);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_applicants_updated_at
    BEFORE UPDATE ON applicants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

