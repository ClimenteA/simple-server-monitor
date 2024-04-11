FROM ubuntu:latest

WORKDIR /home

COPY .env .env

RUN apt-get update && apt-get install -y curl
RUN curl -L -o /home/server https://github.com/ClimenteA/simple-server-monitor/releases/download/v0.0.1/server
RUN chmod +x /home/server

CMD ["/home/server"]