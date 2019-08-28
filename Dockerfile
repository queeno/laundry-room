FROM golang:1.13beta1-alpine AS builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") && \
  go build -a -installsuffix cgo \
  -ldflags "-extldflags -static" \
  -o main .

FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /build/laundry.html /app/
COPY --from=builder /build/static/lilla.png /app/static/
WORKDIR /app
CMD ["./main"]
EXPOSE 8080