FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o stress .

FROM scratch

COPY --from=build /app/stress /usr/local/bin/stress

ENTRYPOINT ["/usr/local/bin/stress"]
