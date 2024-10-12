FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd

# ------------------------------------------------------------------------------

FROM gcr.io/distroless/static-debian12:nonroot AS final

WORKDIR /app
COPY --from=build /app/main .

ENTRYPOINT ["/app/main"]
