FROM golang:1.13-alpine

LABEL maintainer="jamie@jsmth.co.uk"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/jsmithedin/overmyhouse
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package and create test binary
RUN go install -v ./... && \
    CGO_ENABLED=0 GOOS=linux go test -c

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the executable
CMD ["overmyhouse"]
