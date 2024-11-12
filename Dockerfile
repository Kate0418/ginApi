FROM golang:1.23
WORKDIR /back
COPY ./back .

RUN go mod download
RUN go install github.com/air-verse/air@latest
ENV PATH="$PATH:$(go env GOPATH)/bin"

CMD ["air"]
EXPOSE 8080
