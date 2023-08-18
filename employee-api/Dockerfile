FROM golang:1.19 AS build-go
 
WORKDIR /src
COPY go.* ./
RUN go mod download
 
COPY . .
RUN go build -o /go/bin/main 

FROM gcr.io/distroless/base-debian10:nonroot AS run

COPY --from=build-go /go/bin/main /app/main
 
ENTRYPOINT ["/app/main"]