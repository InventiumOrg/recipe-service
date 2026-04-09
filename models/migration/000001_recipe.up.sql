CREATE TABLE "recipe" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(128) NOT NULL,
  "ingredients" text[] NOT NULL,
  "cost" int NOT NULL
);
