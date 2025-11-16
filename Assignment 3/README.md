# Assignment 3

CLI tool to access three MongoDB databases, view their contents and delete/insert new data.

This repo contains docker compose file to launch three MongoDB instances, filler data and script to import them to the databases.

### Building CLI

```
$ git clone https://github.com/leevi-laitala/LUT_DataIntensiveSystems.git
$ go build
```

To be able to run the CLI, three MongoDB URIs must be set as env vars.

### Setup databases

The CLI attempts to access three MongoDB instances, which are defined in `MONGO_URI_1`, `MONGO_URI_2` and `MONGO_URI_3` environment variables. These variables must be defined in order for the program to start. These can be also set in `.env` file, which is provided in the repo.

As is the docker compose will run the databases with ports 27017, 27018 and 27019, for which localhost URIs are already defined in provided `.env` file.

To setup three databases with dummy data:

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

Running built binary launches shell:

```
$ ./assignment3
```

The shell will launch only if the three env variables are set.

Available commands in the shell can be viewed with `help` command:

```
$ help
Available commands:
- exit : Exits CLI
- help : Shows this help menu
- select <server-name> : Selects server to become active
- deselect : Deselects active server
- show <collections | servers> : List collection or servers
- find <collection> <?id> : Print collection contents, or specific record
- insert <collection> : Insert new data to collection from stdin
- delete <collection> <id> : Delete specific record from collection
```
