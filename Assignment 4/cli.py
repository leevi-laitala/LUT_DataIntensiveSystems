import db_api
import inspect
import json

# Shell runs until false
gRunning = True

# Each command has the suffix "cmd"

# Help menu
# "help"
def cmdHelp(*args) -> None:
    l = []

    for cmd, d in gCmds.items():
        arginfo = f" {d[0]}" if d[0] else ""
        l.append(f"{cmd}{arginfo} :: {d[1]}")

    printList("Available commands", l)

# Exit shell
# "exit"
def cmdExit(*args) -> None:
    global gRunning
    gRunning = False

# List servers
# "servers"
def cmdServers(*args) -> None:
    printList("Connected servers", db_api.listServers())

# List tables. Optionally from specific connected server
# "tables <?server-name>"
def cmdTables(*args) -> None:
    detail = ""
    values = []

    # If "server-name" arg is provided
    if len(args) > 0:
        server = args[0]
        detail = f" in server '{server}'"
        values = db_api.listTablesFromServer(server) # From specific server
    else:
        values = db_api.listTables() # From all connected server

    printList(f"Available tables{detail}", values)

# Print data from table/collection. Optionally from specific connected server
# "fetch <table> <?server-name>"
def cmdFetch(*args) -> None: 
    # Enforce mandatory args
    if len(args) == 0:
        print("Must provide table")
        return

    detail = ""
    values = []

    # If "server-name" arg is provided
    if len(args) > 1:
        server = args[1]
        detail = f" in '{server}' server"
        values = db_api.fetchAllFromServer(args[0], server) # Specific server
    else:
        values = db_api.fetchAll(args[0]) # All connected servers

    printList(f"Data of '{args[0]}'{detail}", values)

# Insert data to table/collection
# "insert <table>"
def cmdInsert(*args) -> None:
    # Enforce mandatory args
    if len(args) == 0:
        print("Must provide table")
        return

    db_api.insertData(args[0])

# Delete data item via id
# "delete <table> <id>"
def cmdDelete(*args) -> None:
    # Enforce mandatory args
    if len(args) < 2:
        print("Must provide table and id")
        return

    # Validate whether mandatory "id" argument is an integer
    itemId = int()
    try:
        itemId = int(args[1])
    except Exception as e:
        print("Invalid id: ", e)

    db_api.deleteData(args[0], itemId)

# Pretty print lists with provided title
def printList(t: str, l: list) -> None:
    if len(l) == 0:
        print("Nothing to list")
        return

    print(f"{t}:")
    for line in l:
        print(f" - {line}")

# Available commands
#  Key:   Cmd 
#  Value: (args info, description, callback func)
gCmds = {
#   Key:        Value:
#   Command     Arg info                  Description                   Callback
    "exit":    ("",                       "exit application",           cmdExit),
    "help":    ("",                       "list available commands",    cmdHelp),
    "servers": ("",                       "list connected servers",     cmdServers),
    "tables":  ("<?server-name>",         "list available tables",      cmdTables),
    "fetch":   ("<table> <?server-name>", "list data from a table",     cmdFetch),
    "insert":  ("<table>",                "insert new data to a table", cmdInsert),
    "delete":  ("<table> <id>",           "dalete data from a table",   cmdDelete),
}

# Split input cmd to tokens (basically split from spaces)
def tokenize(s: str) -> list(str):
    tok = s.split(" ")

    if tok[0] == "":
        tok = []

    return tok

def runCli() -> None:
    while gRunning:
        tok = tokenize(input("$ "))

        if len(tok) == 0:
            continue

        # Separate tokens. First token is command and all others are args
        command = tok[0]
        args = tok[1:]

        # Check if command exists
        if not command in gCmds:
            print(f"Unknown command '{command}'. See 'help' for available commands")
            continue

        # Execute command
        gCmds[command][2](*args)
