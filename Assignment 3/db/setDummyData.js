db = connect(process.env["URI"]);

// Drop all collections
db.getCollectionNames().forEach(function(collection) {
    db[collection].drop();
});

db.createCollection("authors");
db.createCollection("books");
db.createCollection("loans");
db.createCollection("members");
db.createCollection("publishers");

num = process.env["DB_NUM"]

db.authors.insertMany(require("./dummydata/db" + num + "_authors.json"));
db.books.insertMany(require("./dummydata/db" + num + "_books.json"));
db.loans.insertMany(require("./dummydata/db" + num + "_loans.json"));
db.members.insertMany(require("./dummydata/db" + num + "_members.json"));
db.publishers.insertMany(require("./dummydata/db" + num + "_publishers.json"));

