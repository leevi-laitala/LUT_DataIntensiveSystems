CREATE TABLE authors (
    _id INT PRIMARY KEY,
    name TEXT,
    country TEXT
);

\set json `cat ./dummydata/db2_authors.json`

INSERT INTO authors
SELECT * FROM jsonb_populate_recordset(NULL::authors, :'json'::jsonb);



CREATE TABLE books (
    _id INT PRIMARY KEY,
    title TEXT,
    authorId INT,
    publisherId INT
);

\set json `cat ./dummydata/db2_books.json`

INSERT INTO books
SELECT 
    (json->>'_id')::int,
    (json->>'title')::text,
    (json->>'authorId')::int,
    (json->>'publisherId')::int
FROM jsonb_array_elements(:'json'::jsonb) AS t(json);



CREATE TABLE loans (
    _id INT PRIMARY KEY,
    memberId INT,
    bookId INT,
    due TEXT
);

\set json `cat ./dummydata/db2_loans.json`

INSERT INTO loans
SELECT 
    (json->>'_id')::int,
    (json->>'memberId')::int,
    (json->>'bookId')::int,
    (json->>'due')::text
FROM jsonb_array_elements(:'json'::jsonb) AS t(json);



CREATE TABLE members (
    _id INT PRIMARY KEY,
    name TEXT,
    email TEXT
);

\set json `cat ./dummydata/db2_members.json`

INSERT INTO members
SELECT * FROM jsonb_populate_recordset(NULL::members, :'json'::jsonb);



CREATE TABLE publishers (
    _id INT PRIMARY KEY,
    name TEXT,
    address TEXT
);

\set json `cat ./dummydata/db2_publishers.json`

INSERT INTO publishers
SELECT * FROM jsonb_populate_recordset(NULL::publishers, :'json'::jsonb);
