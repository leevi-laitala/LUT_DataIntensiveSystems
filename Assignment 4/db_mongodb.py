from db_base import Database
import pymongo

# Implement connect, close, list, ... functions for mongodb database
class MongoDatabase(Database):
    # Override baseclass constructor to define database type as mongodb
    def __init__(self, name: str, user: str, password: str, host: str, port: int):
        super().__init__(name, user, password, host, port)
        self.dbType = "MongoDB"

    # Database client instance
    connection = None
    
    def connect(self) -> bool:
        uri = f"mongodb://{self.host}:{str(self.port)}/{self.name}"

        try:
            self.connection = pymongo.MongoClient(uri)
        except:
            print("MongoDB connection failed")
            return False

        return True

    def close(self):
        self.connection.close()
        print("MongoDB connection closed")

    def list(self) -> list:
        collectionList = self.connection[self.name].list_collections()
        ret = []

        # Only get the collection names
        for c in collectionList:
            ret.append(c["name"])

        return ret

    def fetch(self, collection: str) -> list:
        ret = []

        # Check whether the table exists in the database
        collections = self.list()
        if collection not in collections:
            return ret

        db = self.connection[self.name]
        ret = list(db[collection].find({}))

        return ret

    def insert(self, collection: str, data: dict):
        try:
            coll = self.connection[self.name][collection].insert_one(data)
        except Exception as e:
            print("Insertion error: ", e)
            return

    def delete(self, collection: str, itemId: int):
        try:
            query = { "_id": itemId }
            self.connection[self.name][collection].delete_one(query)
        except Exception as e:
            print("Delete error: ", e)
            return
