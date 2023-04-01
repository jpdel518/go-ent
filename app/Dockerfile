FROM golang:alpine
#FROM golang:latest

WORKDIR /go/src/app
COPY ./ ./

RUN apk update && apk add git
#RUN go mod init flowers # just first time
#COPY go.mod go.sum ./
RUN go build
#RUN go mod download

# hot reload library (for dev)
RUN go get -u github.com/cosmtrek/air
RUN go get -u github.com/go-delve/delve/cmd/dlv

EXPOSE 8080

# 簡易実行
#CMD ["go", "run", "main.go"]
# コンパイル
#CMD ["go", "build", "-o", "main", "main.go"]
# hot reload
#CMD ["air", "-c", ".air.toml"]
