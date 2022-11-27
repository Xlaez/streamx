FROM golang:1.19.3-bullseye

LABEL Name="streamx" Version=0.1.0

EXPOSE 8001

ENV PORT = 8001
ENV HOST = 0.0.0.0

WORKDIR /

COPY . .

RUN go mod download && go mod verify

# COPY . /base/app
# RUN go build -v -o /usr/local/bin/app ./...

CMD ["go", "run", "main.go" ]