FROM golang:latest as builder

WORKDIR /auth_svc/

COPY . .

RUN CGO_ENABLED=0 go build -o auth_svc /auth_svc/cmd/main.go

FROM alpine:latest

WORKDIR /auth_svc

COPY --from=builder /auth_svc/ /auth_svc/

EXPOSE 8001

CMD ./auth_svc