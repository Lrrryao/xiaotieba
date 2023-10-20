CREATE TABLE "session" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "expire_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "client_ip" varchar NOT NULL,
  "useragent" varchar NOT NULL,
  "isblocked" boolean NOT NULL DEFAULT false,
  "refresh_token" varchar NOT NULL
);
