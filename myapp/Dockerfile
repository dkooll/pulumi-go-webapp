FROM golang:1.12.0-alpine3.9 AS builder
RUN mkdir /src/
WORKDIR /src
ADD . /src/
RUN go build -o main .

# Multistage
FROM golang:1.12.0-alpine3.9
WORKDIR /src
COPY --from=builder /src/main /src
CMD ["/src/main"]