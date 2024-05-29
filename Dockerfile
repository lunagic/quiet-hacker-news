FROM golang:1.22
WORKDIR /app
COPY . .
RUN go build -o var/quiet-hacker-news
CMD [ "var/quiet-hacker-news" ]
EXPOSE 8080
