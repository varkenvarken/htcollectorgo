FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod ./
# now we have external deps in our app, to wit: Prometheus
COPY go.sum ./

RUN go mod download

COPY collector/* ./collector/
COPY main.go ./

RUN go build -o /htcollector

## Deploy to a rather minimal image
FROM alpine

WORKDIR /

COPY --from=build /htcollector /htcollector

# expose prometheus metrics
EXPOSE 2112

ENTRYPOINT ["/htcollector"]
