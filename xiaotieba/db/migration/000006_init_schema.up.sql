CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY,
  "role_name" varchar NOT NULL
);

CREATE TABLE "users_roles" (
  "ur_id" integer PRIMARY KEY,
  "user_id" integer NOT NULL,
  "user_name" varchar NOT NULL,
  "role_id" integer NOT NULL,
  "role_name" varchar NOT NULL
);

CREATE TABLE "power" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "roles_power" (
  "rp_id" bigserial PRIMARY KEY,
  "role_id" integer NOT NULL,
  "role_name" varchar NOT NULL,
  "power_id" integer NOT NULL,
  "power_name" varchar NOT NULL
);

ALTER TABLE "users_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

