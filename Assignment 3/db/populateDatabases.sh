#!/usr/bin/env bash

URI="mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.9" \
DB_NUM="1" \
    mongosh --port 27017 --file setDummyData.js

URI="mongodb://127.0.0.1:27018/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.9" \
DB_NUM="2" \
    mongosh --port 27018 --file setDummyData.js

URI="mongodb://127.0.0.1:27019/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.9" \
DB_NUM="3" \
    mongosh --port 27019 --file setDummyData.js

