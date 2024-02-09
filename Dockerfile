FROM golang:alpine

WORKDIR /app 

COPY . .

EXPOSE 1234

CMD [ "go","run","main.go" ]
