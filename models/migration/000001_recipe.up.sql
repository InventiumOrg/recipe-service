CREATE TABLE "recipe" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(128) NOT NULL,
  "ingredients" jsonb NOT NULL DEFAULT '[]'::jsonb,
  "cost" int NOT NULL
);
