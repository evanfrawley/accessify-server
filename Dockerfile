FROM alpine
RUN apk add --no-cache ca-certificates
COPY accessify-server /accessify-server
EXPOSE 80
ENTRYPOINT ["/accessify-server"]
