CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "words" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "language" CHAR(3), -- ISO 639-3
  "word" VARCHAR(255),
  "example" TEXT, -- example sentences and use cases
  "image_url" VARCHAR(255), -- image to define the word
  "link" VARCHAR(255), -- link to page for analyzing purpose
  "created_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "definitions" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "language" CHAR(3), -- ISO 639-3
  "definition" TEXT,
  "word_fk" UUID REFERENCES "words"("id"),
  "created_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" VARCHAR(255),
  "email" VARCHAR(255),
  "created_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "userwords" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_fk" UUID REFERENCES "users"("id"),
  "word_fk" UUID REFERENCES "words"("id"),
  "status" INTEGER,
  "note" TEXT,
  "created_at" TIMESTAMP DEFAULT now()
);
