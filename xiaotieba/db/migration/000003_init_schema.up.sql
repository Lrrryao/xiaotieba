CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "username" varchar,
  "email" varchar DEFAULT (now()),
  "secret_code" varchar NOT NULL,
  "is_secret_used" boolean NOT NULL DEFAULT false,
  "secret_created_at" timestamptz NOT NULL DEFAULT (now()),
  "secret_expired_at" timestamptz NOT NULL DEFAULT (now()+interval '15 minutes')
);

