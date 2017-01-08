
FROM golang:latest

# Download RQLite and install in /usr/local/bin
WORKDIR /usr/local/bin
COPY rqlited-v3.6.0-linux-amd64.tar.gz /usr/local/bin
RUN ["tar","xvfz","rqlited-v3.6.0-linux-amd64.tar.gz"]
RUN ["mv","rqlited-v3.6.0-linux-amd64","rqlite"]
#RUN ["chmod", "+x","rqlite"]

# copy my app to go src folder and compile
ADD / /go/src
RUN ["go","build","/go/src/rq4d.go"]

ENTRYPOINT ["./rq4d"]