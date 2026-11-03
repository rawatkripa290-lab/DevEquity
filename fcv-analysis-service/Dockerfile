
FROM golang:1.23-alpine AS build

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod ./

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

RUN CGO_ENABLED=0 go build -o /out/worker ./cmd/worker


FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /out/api /app/api
COPY --from=build /out/worker /app/worker

EXPOSE 8080

CMD ["/app/api"]
