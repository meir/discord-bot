FROM golang:1.17.2-alpine3.14 AS go

ARG PROJECT="executable"

WORKDIR /src
COPY . .

ARG PROJECT="executable"
ARG VERSION=???

RUN go mod download
RUN go build -o /src/${PROJECT} /src/app/${PROJECT}

#

FROM alpine:latest

LABEL image=${PROJECT}
LABEL maintainer="github.com/meir"
LABEL madew="love"

ARG PROJECT="executable"
ARG VERSION=???

ENV VERSION=$VERSION
ENV WEB="/app/website"

ENV DEBUG_WEBHOOK=
ENV DEBUG=false

RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=go /src/${PROJECT} ./
COPY --from=go /src/assets ./assets
COPY --from=go /src/web ./web

RUN chmod +x /app/${PROJECT}

CMD /app/${PROJECT}
