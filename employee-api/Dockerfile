FROM golang:1.19 AS build-go
 
WORKDIR /src
COPY go.* ./
RUN go mod download
 
COPY . .
RUN go build -o /go/bin/server github.com/pheely/employee-api/server

FROM gcr.io/distroless/base-debian10:nonroot AS run

COPY --from=build-go /go/bin/server /app/server
 
ENTRYPOINT ["/app/server"]
