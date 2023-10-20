ALTER TABLE verify_emails
ALTER COLUMN email DROP DEFAULT;
CREATE UNIQUE INDEX idx_verify_emails_email ON verify_emails (email);

