db = connect(process.env["MONGO_URI"]);

// Drop all collections
db.getCollectionNames().forEach(function(collection) {
    db[collection].drop();
});

db.createCollection("authors");
db.createCollection("books");
db.createCollection("events");
db.createCollection("loans");
db.createCollection("reviews");

db.authors.insertMany(require("./dummydata/db1_authors.json"));
db.books.insertMany(require("./dummydata/db1_books.json"));
db.events.insertMany(require("./dummydata/db1_events.json"));
db.loans.insertMany(require("./dummydata/db1_loans.json"));
db.reviews.insertMany(require("./dummydata/db1_reviews.json"));

