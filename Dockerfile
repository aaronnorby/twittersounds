FROM golang:1.6

ADD . /go/src/github.com/aaronnorby/twittersounds

RUN go get github.com/ChimeraCoder/anaconda && go get golang.org/x/net/html

RUN go install github.com/aaronnorby/twittersounds/tweetout

ENTRYPOINT /go/bin/tweetout
