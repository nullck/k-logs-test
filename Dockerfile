ARG GO_VERSION=1.15
ARG VERSION=dev

FROM --platform=${BUILDPLATFORM:-linux/amd64} tonistiigi/xx:golang AS xgo

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:${GO_VERSION}-alpine

ADD . /go/src/app
WORKDIR /go/src/app
ENV ELASTIC_SEARCH_ADDR="http://localhost"
ENV ELASTIC_SEARCH_PORT=9200
ENV ELASTIC_INDEX_NAME="test_logs"
ENV LOGS_HITS=10
ENV THRESHOULD_MS=8000
ENV SLACK_CHANNEL="#k-logs"
ENV SLACK_WEBHOOK=""
ENV SLACK_ENABLED="true"
ENV PROM_ENABLED="true"
ENV PROM_GW_URL="localhost"

CMD ["go", "run", "main.go", "run", "-e", "${ELASTIC_SEARCH_ADDR}:${ELASTIC_SEARCH_PORT}/${ELASTIC_INDEX_NAME}", "--logs-hits", "${LOGS_HITS}", "--threshold", "${THRESHOULD_MS}", "--channel", "${SLACK_CHANNEL}", "--webhook-url", "${SLACK_WEBHOOK}", "--slack-alert-enabled", "${SLACK_ENABLED}", "--prom-enabled", "${PROM_ENABLED}", "--prom-endpoint", "${PROM_GW_URL}"]
