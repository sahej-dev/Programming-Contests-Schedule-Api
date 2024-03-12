FROM golang:1.21

WORKDIR /app

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY src/ ./

RUN GOOS=linux go build -o /executable

EXPOSE 4000
