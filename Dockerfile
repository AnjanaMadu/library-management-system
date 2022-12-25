FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o main . && chmod +x main && go clean

CMD ["/app/main"]