from mongodb import MongoDatabase
from postgres import PostgresDatabase
import db_api
import cli

if __name__ == "__main__":
    db_api.connectToDatabase(PostgresDatabase("database", "user", "password", "localhost", 27018))
    db_api.connectToDatabase(MongoDatabase("test", "", "", "localhost", 27017))

    cli.runCli()

    db_api.closeConnections()
