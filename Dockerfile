FROM golang:1.26.2 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /api_server ./cmd/api_server

FROM gcr.io/distroless/base-debian13

WORKDIR /

COPY --from=build /api_server /api_server

EXPOSE 8080

USER nonroot:nonroot

CMD ["/api_server"]