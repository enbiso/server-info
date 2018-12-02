FROM golang AS builder
WORKDIR /app
COPY . .
RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM scratch
COPY --from=builder /app/server-info /app
ENTRYPOINT [ "/app/server-info" ]