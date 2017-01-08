# rq4d

Wrapper for RQLite Distributed Database to run inside docker compose based environment

## How to use / check it out

### Pre-Requisites

Its Simple if you have the following already Installed and in working condition

- Docker 1.12
- Docker-compose
- Git

### Clone this repository

    git clone https://github.com/julianfrank/rq4d

    cd rq4d

### Build the Container

This step would create the seed Instance that detects that it was the only container in the rdb network and hence spin up rqlite as seed

    docker-compose up --build -d rdb

Watch the log using

    docker-compose logs -f

If you want to change the port addresses or other parameters then change the ENTRYPOINT parameter to use your custom parameter as below

    # ./rq4d --help
    Usage of ./rq4d:
      -db string
            DB Directory as applicable in your Container environment (default "/db")
      -exec string
            Executable as applicable in your Container environment (default "./rqlite/rqlited")
      -http string
            Port for use with the http parameter (default ":4001")
      -raft string
            Port for use with the raft parameter (default ":4002")
      -sername string
            This should be the service Name used by the rqLite Cluster (default "rdb")

for example if you like the raft port to be 4040 instead of the default 4002 then change the ENTRYPOINT as follows and then rebuild

    ENTRYPOINT ["./rq4d","-raft",":4040"]

Once happy with build in future you can startup the seed just using the regular

    docker-compose up -d rdb

### Scale the Database

This is where this wrapper makes life easier ... If you want 5 Instances use the following

    docker-compose scale rdb=5

## Known Issues

- The DB Directory does not work with Volumes that directly end up in the host disk. If you use a Volume 'Container' then this problem does not showup...
