FROM goalng:1.20 alpine:3.11 AS build
# Use the official Golang image to build the Go application
RUN apk --no-cache add gcc g++ make ca-certificates 
#Install necessary packages
WORKDIR D:programming/go/DevOps/Projects/go-microservices-project
# Set the working directory inside the container
COPY go.mod go.sum ./
# Copy the Go module files
COPY vendor vendor
# Copy the vendor directory
COPY account account    
# Copy the account directory
RUN  GO111MODULE=on go build -mod vendor -o /go/bin/app ./account/cmd/account
# Build the Go application and output it to /go/bin/app 
FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
