FROM golang

RUN mkdir -p go-ping

WORKDIR /home/go-ping/

RUN go mod tidy

RUN chmod +x build.sh

CMD ["./build.sh"]
