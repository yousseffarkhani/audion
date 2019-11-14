FROM golang AS builder

WORKDIR /go/src/github.com/yousseffarkhani/audion
ADD . .

RUN go get ./
RUN CGO_ENABLED=0 GOOS=linux go build -o audion

FROM alpine:latest AS production
COPY --from=builder /go/src/github.com/yousseffarkhani/audion .
CMD ["./audion"]