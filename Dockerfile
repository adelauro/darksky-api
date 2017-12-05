FROM alpine
RUN apk add --no-cache ca-certificates
ADD darksky /
EXPOSE 80
CMD ["/darksky"]
