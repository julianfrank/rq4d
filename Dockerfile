FROM golang:latest

# Download RQLite and install in /usr/local/bin
WORKDIR /go/bin
COPY rqlited-v3.9.1-linux-amd64.tar.gz /go/bin
RUN ["tar","xvfz","rqlited-v3.9.1-linux-amd64.tar.gz"]
RUN ["mv","rqlited-v3.9.1-linux-amd64","rqlite"]

# copy my app to go src folder and compile
ADD / /go/src
RUN ["go","build","/go/src/rq4d.go"]

ENTRYPOINT ["./rq4d"]