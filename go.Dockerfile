FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o server .

FROM alpine:latest AS server

WORKDIR /app
COPY --from=build /app/server .

EXPOSE 3000
CMD ["./server"]