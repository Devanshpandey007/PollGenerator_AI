# Use the official Golang image as the base for the build stage
FROM golang:1.23.6 AS build


# Set the working directory for the build stage
WORKDIR /app

# Copy the entire Shared module code
COPY . /app

# Download dependencies
RUN go mod download

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o /app/main /app/Service/main.go

# Final stage - Lambda runtime
FROM public.ecr.aws/lambda/provided:al2023

# Copy the built binary from the build stage
COPY --from=build /app/main /var/runtime/bootstrap

# Run the Lambda function
ENTRYPOINT [ "/var/runtime/bootstrap" ]
