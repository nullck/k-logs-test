#!/bin/bash

for i in kind kubectl; do
  if ! command -v "${i}" &> /dev/null; then
    echo "I cannot find the kind binary"
    echo "Trying to install if the OS is Linux"

    uname | grep "Linux"

    if [ $? == 0  ]; then
      echo "it is Linux"
      if [ "${i}" == "kind" ]; then
        curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.8.1/kind-linux-amd64
        chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind
      fi
      if [ "${i}" == "kubectl" ]; then
        curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
        chmod +x kubectl && sudo mv ./kubectl /usr/local/bin/kubectl
      fi
    else
      echo "please check your kind AND/OR kubectl installation"
      exit 1
    fi
  fi
done
# check if thet kube-logs-test cluster exists
CLUSTER_NAME="kube-logs-test"

if [ -z "$1" ]; then
  echo "use $0 start or destroy"
  exit 1
fi

if [ "$1" == "start" ]; then
  kind get clusters | grep "${CLUSTER_NAME}"
  if [ $? != 0 ]; then
    kind create cluster --name ${CLUSTER_NAME} --config scripts/kind-1-18.yaml;
    sleep 15;
    kubectl get pods
  fi
  kubectl apply -f scripts/fluentbit
  kubectl apply -f scripts/elastic
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

if [ "$1" == "destroy" ]; then
  kind get clusters | grep "${CLUSTER_NAME}"
  if [ $? == 0 ]; then
    kind delete cluster --name ${CLUSTER_NAME}
  fi
fi

exit 0
