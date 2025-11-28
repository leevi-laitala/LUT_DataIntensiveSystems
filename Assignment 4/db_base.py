# Base class for databases
class Database:
    dbType: str
    name: str
    user: str
    password: str
    host: str
    port: int
    connected: bool

    def __init__(self, name: str, user: str, password: str, host: str, port: int):
        self.name = name
        self.user = user
        self.password = password
        self.host = host
        self.port = port

        self.connected = False

        self.dbType = "invalid"

    def connect(self) -> bool:
        print("invalid connected")
        return connected

    def close(self):
        print("invalid closed")

    def list(self) -> list:
        print("invalid list")
        return []

    def fetch(self):
        print("invalid fetch")

    def insert(self, table: str, data: str):
        print("invalid insert")

    def delete(self):
        print("invalid delete")

