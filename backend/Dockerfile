FROM golang:1.22.2
WORKDIR /app
# ホットリロード用モジュール
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]