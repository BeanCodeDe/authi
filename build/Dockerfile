FROM golang:1.20.3 as development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN make app.build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/authi ./
CMD ["./authi"]  