CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "hash_password" varchar NOT NULL,
  "name" varchar UNIQUE NOT NULL,
  "power" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "phone" varchar(20) NOT NULL
);

CREATE TABLE "post" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "titles" varchar NOT NULL,
  "content" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "comment" (
  "id" bigserial PRIMARY KEY,
  "content" varchar NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "reply" (
  "id" bigserial PRIMARY KEY,
  "content" varchar NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "comment_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "post" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "comment" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "comment" ADD FOREIGN KEY ("post_id") REFERENCES "post" ("id");

ALTER TABLE "reply" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "reply" ADD FOREIGN KEY ("post_id") REFERENCES "post" ("id");

ALTER TABLE "reply" ADD FOREIGN KEY ("comment_id") REFERENCES "comment" ("id");