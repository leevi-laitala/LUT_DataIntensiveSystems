from db_postgres import PostgresDatabase
from db_mongodb import MongoDatabase
import json
import sys

# Connected databases
gDatabases = []

# Append new database to connected databases upon successful connection
def connectToDatabase(db):
    success = db.connect()
    if success:
        gDatabases.append(db)

# Close all open database connections
def closeConnections():
    for db in gDatabases:
        db.close()

# List connected servers
def listServers() -> list:
    if len(gDatabases) == 0:
        print("No connected databases")

    l = []

    for db in gDatabases:
        l.append(f"{db.dbType} @ {db.host}:{db.port}")

    return l

# Match and return database instance by type (mongodb, postgres, ...)
def getServer(server: str):
    # Iterate through all connected databases
    for db in gDatabases:
        if db.dbType == server:
            return db

    print("Could not find server")
    return None

# List all tables/collections from specific connected server
def listTablesFromServer(server: str) -> list:
    db = getServer(server)
    if not db:
        return []

    return db.list()

# List all tables/collecitons from all connected servers without duplicates
def listTables() -> list:
    allTables = set() # Disallow duplicates -> store in set instead of list

    for db in gDatabases:
        allTables |= set(listTablesFromServer(db.dbType))

    # Return list instead of set
    return list(allTables)

# List data from specific table/collection from specific connected server
def fetchAllFromServer(table: str, server: str) -> list:
    results = []

    db = getServer(server)
    if not db:
        return results

    # Sanitize user input query. Only allow letters
    sanitizedTable = ''.join(char for char in table if char.isalpha())
    data = db.fetch(sanitizedTable)

    return list(data)

# List data from specific table/collection from all connected servers
# Won't return duplicate items
def fetchAll(table: str) -> list:
    results = []

    for db in gDatabases:
        # Fetch data from single database
        newdata = fetchAllFromServer(table, db.dbType)

        # Check whether the data items already exist in return list
        for data in newdata:
            alreadyExists = False

            for res in results:
                if res["_id"] == data["_id"]:
                    alreadyExists = True

            # If item is unique, append to return list
            if not alreadyExists:
                results.append(data)

    return list(results)

def insertData(table: str):
    # Find out the first database instance where the table/collection exist
    insertionDb = None
    for db in gDatabases:
        tables = db.list()

        if table in tables:
            insertionDb = db

    if not insertionDb:
        print(f"Table '{table}' does not exist. Data not inserted")
        return

    # Wait for user input
    print("Insert new data item in JSON format. Finish with Ctrl+D:")
    lines = sys.stdin.read()

    data = {}
    try:
        data = json.loads(lines)
    except Exception as e:
        print("Failed to parse user input: ", e)
        return

    # Execute data insertion
    insertionDb.insert(table, data)

def deleteData(table: str, itemId: int):
    # Find out the first database instance where the table/collection exist
    deletionDb = None
    for db in gDatabases:
        tables = db.list()

        if table in tables:
            deletionDb = db

    if not deletionDb:
        print(f"Table '{table}' does not exist. Data not inserted")
        return
    
    deletionDb.delete(table, itemId)

