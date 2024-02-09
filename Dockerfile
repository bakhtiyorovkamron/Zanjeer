FROM golang:alpine

WORKDIR /app 

COPY . .

EXPOSE 7777

CMD [ "go","run","main.go" ]
