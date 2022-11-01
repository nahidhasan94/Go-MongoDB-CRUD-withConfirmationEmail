FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/bin/Go-MongoDB-CRUD-withConfirmationEmail .

FROM golang:latest

WORKDIR /app

COPY --from=builder /app/bin /app

EXPOSE 8080

# Run the executable
CMD ["./Go-MongoDB-CRUD-withConfirmationEmail"]
