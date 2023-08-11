FROM golang:1.19 AS build-go
 
WORKDIR /src
COPY go.* ./
RUN go mod download
 
COPY . .
RUN go build -o /go/bin/event event-trigger

FROM gcr.io/distroless/base-debian10:nonroot AS run

COPY --from=build-go /go/bin/event /app/event
 
ENTRYPOINT ["/app/event"]