FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0 
ENV GOPROXY https://goproxy.cn,direct
ENV GOCACHE /build/.cache/go-build

WORKDIR /build

COPY go.mod .
COPY . .
RUN go mod tidy
RUN go build -ldflags="-s -w" -o /app/main ./main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/main ./
COPY --from=builder /build/login.html ./

CMD ["./main"]



