FROM golang:1.26.2 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY api ./api
COPY dao ./dao

RUN go build -o /crud-example

FROM gcr.io/distroless/base-debian13

WORKDIR /

COPY --from=build /crud-example /crud-example

EXPOSE 8080

USER nonroot:nonroot

CMD ["/crud-example"]