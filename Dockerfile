FROM golang:1.17.2-alpine3.14 AS go

ARG PROJECT="executable"

WORKDIR /go
COPY . .

RUN go mod download
RUN go build -o /go/${PROJECT} ./go/app

#

FROM alpine:latest

LABEL image=${PROJECT}
LABEL maintainer="github.com/meir"
LABEL madew="love"

ARG PROJECT="executable"
ARG VERSION=???

ENV VERSION=$VERSION
ENV WEB="/root/website"

ENV DEBUG_WEBHOOK=
ENV DEBUG=false

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=go /go/${PROJECT} ./
COPY --from=go /go/assets ./assets
COPY --from=go /go/web ./web

RUN chmod +x /root/${PROJECT}

CMD /root/${PROJECT}
