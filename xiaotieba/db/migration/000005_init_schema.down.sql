ALTER TABLE verify_emails
ALTER COLUMN email SET DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_verify_emails_email 