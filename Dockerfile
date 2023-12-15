FROM golang:1.21
WORKDIR /app
ADD . ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o v2ex_server

FROM alpine:latest
WORKDIR /app
ADD --from=builder /app/v2ex_server ./
ADD --from=builder /app/static ./static
EXPOSE 80
EXPOSE 443
CMD [ "/app/v2ex_server","-c","/app/config.yaml"]