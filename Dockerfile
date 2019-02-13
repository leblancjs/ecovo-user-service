FROM golang:latest AS build

ENV PROJECT_NAME=azure.com/ecovo/user-service
ENV BINARY_NAME=user-service

# Force the project to use Go mod for dependencies
ENV GO111MODULE=on

# Copy the project files
COPY . $GOPATH/src/${PROJECT_NAME}
WORKDIR $GOPATH/src/${PROJECT_NAME}

# Manage dependencies
RUN go mod download

# Build the project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o /bin/${BINARY_NAME}

# Expose port
EXPOSE 8080/tcp

# Start a new container from scratch to keep only the compiled binary
FROM scratch
COPY --from=build /bin/${BINARY_NAME} /bin/${BINARY_NAME}
ENTRYPOINT ["/bin/user-service"]
