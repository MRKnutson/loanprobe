FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/michaelrknutson/loanprogo
RUN cd /build && git clone https://github.com/michaelrknutson/loanprobe.git

RUN cd /build/loanprogo && go build
EXPOSE 8080

ENTRYPOINT [ "/build/loanprogo/main" ]