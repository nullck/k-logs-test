ARG GO_VERSION=1.15
ARG VERSION=dev

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS build

WORKDIR /src
ENV CGO_ENABLED=1

ADD . /go/src/app
WORKDIR /go/src/app

ARG TARGETOS
ARG TARGETARCH
RUN apk add build-base
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/src/k-logs .

ENV K_LOGS_LOGS_HITS=30
ENV K_LOGS_CHANNEL="#k-logs"
ENV K_LOGS_ELASTIC_ENDPOINT="https=//localhost=9200/fluentd"
ENV K_LOGS_NAMESPACE="monitoring"
ENV K_LOGS_POD_NAME="test-logs"
ENV K_LOGS_SLACK_ALERT_ENABLED="true"
ENV K_LOGS_THRESHOLD="2000"
ENV K_LOGS_WEBHOOK_URL="https=//hooks.slack.com/services/XXXXXXX/YYYYYY/WWWWW"
ENV K_LOGS_PROM_ENABLED="true"
ENV K_LOGS_PROM_ENDPOINT="localhost"

WORKDIR /go/src

CMD ["./k-logs"]

