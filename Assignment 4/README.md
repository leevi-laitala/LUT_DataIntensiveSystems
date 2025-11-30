# Assignment 4

CLI tool to access databases, MongoDB and ProstgeSQL, view their contents and delete/insert new data.

This repo contains docker compose file to launch the database instances, filler data and a script to import them to the databases.

### Setup

```
$ git clone https://github.com/leevi-laitala/LUT_DataIntensiveSystems.git
$ cd LUT_DataIntensiveSystems/Assignment\ 4
$ python -m venv runner
$ ./runner/bin/pip install -r requirements.txt
$ ./runner/bin/python main.py
```

To be able to run the CLI, three MongoDB URIs must be set as env vars.

### Setup databases

The CLI attempts to access two database instances (MongoDB & Postgres), connections of which are defined in `main.py`

As is the docker compose will run the databases with ports 27017 and 27018

To setup three databases with dummy data navigate to `db` dir inside the repo:

```
$ cd db
$ docker compose up -d
$ ./populateDatabases.sh
```

Teardown with:

```
$ docker compose stop
$ docker compose rm
```

### Running CLI

Running the following command will start a shell interface that waits for user commands:

```
$ ./runner/bin/python main.py
```

Available commands in the shell can be viewed with `help` command:

```
$ help
Available commands:
 - exit :: exit application
 - help :: list available commands
 - servers :: list connected servers
 - tables <?server-name> :: list available tables
 - fetch <table> <?server-name> :: list data from a table
 - insert <table> :: insert new data to a table
 - delete <table> <id> :: dalete data from a table
```

