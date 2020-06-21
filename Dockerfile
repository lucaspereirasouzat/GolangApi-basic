FROM golang:latest

# ENV HOME=./
# COPY . $HOME/

# USER root
# WORKDIR $HOME

# # ADD . /

# FROM golang

# RUN go get -v

# # RUN go get -u github.com/gin-gonic/gin
# # RUN go get golang.org/x/net/html
# # RUN go get github.com/go-redis/redis

# # RUN go install wscraper


# # ENTRYPOINT $HOME
# RUN go run main.go

EXPOSE 5430
# EXPOSE 2020

# FROM golang:1.11.1-alpine3.8 as build-env
# # All these steps will be cached
# RUN mkdir /hello
# WORKDIR /hello
# COPY go.mod . 
# # <- COPY go.mod and go.sum files to the workspace
# COPY go.sum .

# RUN go run main.go

# COPY . .

# # Build the binary
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/helloFROM scratch 
# # <- Second step to build minimal image
# COPY --from=build-env /go/bin/hello /go/bin/hello
# ENTRYPOINT ["/go/bin/hello"]

