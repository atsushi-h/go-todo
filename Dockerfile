# Go公式イメージを使用
FROM golang:1.24

# 作業ディレクトリを設定
WORKDIR /app

# go.mod と go.sum をコピー
COPY go.mod go.sum ./

# 依存パッケージを取得
RUN go mod download

# アプリ全体をコピー
COPY . .

# ポートを開放
EXPOSE 8080

# アプリをビルドして実行（ビルドは後で書き換える予定）
CMD ["go", "run", "cmd/main.go"]
