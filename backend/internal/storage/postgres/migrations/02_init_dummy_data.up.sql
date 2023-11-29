-- Inserting a word
INSERT INTO "words" ("language", "word", "example", "image_url", "link")
VALUES ('eng', 'example_word', 'This is an example sentence.', 'example_image.jpg', 'https://example.com');

-- Inserting definitions in different languages
INSERT INTO "definitions" ("language", "definition", "word_fk")
VALUES ('eng', 'This is the English definition of the word.', (SELECT "id" FROM "words" WHERE "word" = 'example_word'));

INSERT INTO "definitions" ("language", "definition", "word_fk")
VALUES ('spa', 'Esta es la definición en español de la palabra.', (SELECT "id" FROM "words" WHERE "word" = 'example_word'));

-- Repeat the above process for other languages as needed

-- Inserting a user
INSERT INTO "users" ("name", "email")
VALUES ('John Doe', 'john@example.com');

-- Inserting a userword
INSERT INTO "userwords" ("user_fk", "word_fk", "status", "note")
VALUES (
  (SELECT "id" FROM "users" WHERE "email" = 'john@example.com'),
  (SELECT "id" FROM "words" WHERE "word" = 'example_word'),
  1, -- Replace with the actual status value
  'This is a note about the word.'
);