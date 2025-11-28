#!/usr/bin/env bash

set -e

PG_DB_NAME="database"
PG_DB_USER="user"
PG_DB_HOST="127.0.0.1"
PG_DB_PORT="27018"

MONGODB_HOST="127.0.0.1"
MONGODB_PORT="27017"

psql -U "$PG_DB_USER" --host "$PG_DB_HOST" --port "$PG_DB_PORT" -c "DROP DATABASE IF EXISTS $PG_DB_NAME;"
psql -U "$PG_DB_USER" --host "$PG_DB_HOST" --port "$PG_DB_PORT" -c "CREATE DATABASE $PG_DB_NAME;"
psql -U "$PG_DB_USER" --host "$PG_DB_HOST" --port "$PG_DB_PORT" -d "$PG_DB_NAME" -f "setPostgresDummyData.sql"

MONGO_URI="mongodb://$MONGODB_HOST:$MONGODB_PORT/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.9" \
    mongosh --port "$MONGODB_PORT" --file setMongoDummyData.js
