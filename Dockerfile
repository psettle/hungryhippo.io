FROM golang:latest

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -o main .
COPY wrapperscript.sh wrapperscript.sh

ENTRYPOINT ["sh", "wrapperscript.sh"]
