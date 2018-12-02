FROM golang AS builder
WORKDIR /app
COPY . .
RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o "/app/serverinfo"

FROM scratch
WORKDIR /app
COPY --from=builder /app/serverinfo /app
ENTRYPOINT [ "./serverinfo" ]