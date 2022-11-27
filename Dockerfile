FROM golang:1.19.3-bullseye

LABEL Name="streamx" Version=0.1.0

EXPOSE 8001

ENV PORT = 8001
ENV HOST = 0.0.0.0

RUN go mod tidy

RUN mkdir -p /base/app


COPY . /base/app

CMD [ "go", "run", "main.go" ]