FROM golang:alpine

COPY ["database.accdb", "public", "html", "main.go", "go.mod", "go.sum", "/"]

RUN go build -o /main && go clean \
    && rm -rf main.go go.mod go.sum html

EXPOSE 8080

CMD ["/main"]