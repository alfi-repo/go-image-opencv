FROM docker.io/gocv/opencv:4.8.1

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o /goapp

EXPOSE 3000

CMD [ "/goapp" ]