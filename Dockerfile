FROM golang:1.8

ENV GOPATH /opt/CSE/cse-autotest

RUN mkdir -p $GOPATH

ADD source.tar.gz $GOPATH

WORKDIR $GOPATH

CMD ["sh", "/opt/CSE/cse-autotest/src/code.huawei.com/cse/test-case/test.sh"]
