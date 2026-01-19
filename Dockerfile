# ===== build stage =====
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 先拷貝 go.mod / go.sum（利用 cache）
COPY go.mod go.sum ./
RUN go mod download

# 再拷貝全部程式碼
COPY . .

# ⭐⭐⭐ 關鍵：指定 main.go 所在資料夾
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api

# ===== runtime stage =====
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
