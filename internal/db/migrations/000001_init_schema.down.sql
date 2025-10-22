-- Drop trigger
DROP TRIGGER IF EXISTS update_applicants_updated_at ON applicants;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_applicants_created_at;
DROP INDEX IF EXISTS idx_applicants_overall_score;
DROP INDEX IF EXISTS idx_applicants_status;
DROP INDEX IF EXISTS idx_applicants_position;
DROP INDEX IF EXISTS idx_applicants_email;

-- Drop table
DROP TABLE IF EXISTS applicants;
