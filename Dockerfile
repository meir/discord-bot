FROM golang:1.17.2-alpine3.14 AS go

ARG PROJECT="executable"

WORKDIR /src
COPY . .

ARG PROJECT="executable"
ARG VERSION=???

RUN go mod download
RUN go build -o /src/${PROJECT} /src/app/${PROJECT}

#

FROM golang:1.17.2-alpine3.14

LABEL image=${PROJECT}
LABEL maintainer="github.com/meir"
LABEL madew="love"

ARG PROJECT="executable"
ARG VERSION=???

ENV VERSION=$VERSION
ENV WEB="/app/website"

ENV DEBUG_WEBHOOK=
ENV DEBUG=false

RUN apk add --no-cache mysql-client

RUN mkdir /app

WORKDIR /app/
COPY --from=go /src/${PROJECT} /app/program
COPY --from=go /src/assets /app/assets
COPY --from=go /src/web /app/web
RUN chmod +x /app/program

CMD /app/program
