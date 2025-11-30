from db_base import Database
from psycopg.rows import dict_row
import psycopg

# Implement connect, close, list, ... functions for postgres database
class PostgresDatabase(Database):
    # Override baseclass constructor to define database type as postgres
    def __init__(self, name: str, user: str, password: str, host: str, port: int):
        super().__init__(name, user, password, host, port)
        self.dbType = "Postgres"

    # Database client instance
    connection = None

    def connect(self) -> bool:
        uri = f"postgresql://{self.user}:{self.password}@{self.host}:{str(self.port)}/{self.name}"

        try:
            self.connection = psycopg.connect(uri)
        except:
            print("Postgres connection failed")
        
        # Return true on success
        return not (self.connection.closed or self.connection.broken)

    def close(self):
        self.connection.close()

        if self.connection.closed:
            print("Postgres connection closed")
        else:
            print("Error closing postgres connection")
    
    def list(self) -> list:
        # Get tables from database
        with self.connection.cursor() as cursor:
            cursor.execute("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
            tables = cursor.fetchall()

        ret = []

        for table in tables:
            ret.append(table[0])

        return ret

    # Fetch data from given table
    def fetch(self, table: str) -> list:
        ret = []

        # Check whether the table exists in the database
        tables = self.list()
        if table not in tables:
            return ret

        # Create cursor that returns data as a dictionary
        with self.connection.cursor(row_factory = dict_row) as cursor:
            try:
                cursor.execute(f"SELECT * FROM {table};")
            except:
                return ret

            ret = cursor.fetchall()

        return ret

    def insert(self, table: str, data: dict):
        # Separate dictionary keys and values to use as columns and rows
        cols = list(data.keys())
        vals = list(data.values())

        # Create placeholder for "cursor.execute" function
        # E.g. (id, data1, data2) would result placeholder string "%s %s %s"
        placeholders = ", ".join(["%s"] * len(vals))

        # Convert dictionary keys to a comma separated string
        colnames = ", ".join(cols)

        query = f"INSERT INTO {table} ({colnames}) VALUES ({placeholders});"

        try:
            with self.connection.cursor() as cursor:
                cursor.execute(query, vals) # 'vals' populate the placeholders
            self.connection.commit()
        except Exception as e:
            self.connection.rollback()
            print("Insert failed: ", e)

    def delete(self, table: str, itemId: int):
        query = f"DELETE FROM {table} WHERE _id = %s" # %s is placeholder

        try:
            with self.connection.cursor() as cursor:
                cursor.execute(query, (itemId,)) # Provide id for placeholder
            self.connection.commit()
        except Exception as e:
            self.connection.rollback()
            print("Delete failed: ", e)
