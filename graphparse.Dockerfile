FROM golang:1.10.3

RUN apt-get install bash wget curl git

RUN mkdir -p $$GOPATH/bin && \
    go get github.com/pilu/fresh && \
    go get -u github.com/govend/govend

CMD fresh -c runner.conf main.go