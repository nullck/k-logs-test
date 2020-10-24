#!/bin/bash

if [ -z "$1" ]; then
  echo "use $0 install or clean"
  exit 1
fi

export KUBECONFIG="/Users/nullck/.kube/config-rasp"

if [ "$1" == "install" ]; then
  kubectl apply -f scripts/fluentbit
  kubectl apply -f scripts/elastic-arm/elastic.yaml
  kubectl apply -f scripts/prometheus-pushgateway
  sleep 5
  while ! kubectl get pods/elasticsearch-0 | grep "Running"; do
    sleep 2
  done
  # check if the connection through ports 9200 and 9091 are stablished. In case not, start the proxies
  nc -z -v -w 2 127.0.0.1 9200
  if [ $? != 0 ]; then
    kubectl port-forward svc/elasticsearch 9200:9200 &
  fi
  nc -z -v -w 2 127.0.0.1 9091
  if [ $? != 0 ]; then
    kubectl port-forward svc/prometheus-pushgateway 9091:9091 &
  fi
fi

if [ "$1" == "clean" ]; then
  kubectl delete -f scripts/fluentbit
  kubectl delete -f scripts/elastic-arm
  kubectl delete -f scripts/prometheus-pushgateway
fi

exit 0
