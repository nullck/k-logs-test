#!/bin/bash

if ! command -v kind &> /dev/null; then
  echo "I cannot find the kind binary; please check your installation"
  exit 1
fi

# check if thet kube-logs-test cluster exists
CLUSTER_NAME="kube-logs-test"

if [ -z "$1" ]; then
  echo "use $0 start or destroy"
  exit 1
fi

if [ "$1" == "start" ]; then
  kind get clusters | grep "${CLUSTER_NAME}"
  if [ $? != 0 ]; then
    kind create cluster --name ${CLUSTER_NAME};
    sleep 15;
  fi
  kubectl apply -f test-pod.yaml
  kubectl apply -f fluentbit
  kubectl apply -f elastic
fi

if [ "$1" == "destroy" ]; then
  kind get clusters | grep "${CLUSTER_NAME}"
  if [ $? == 0 ]; then
    kind delete cluster --name ${CLUSTER_NAME}
  fi
fi

exit 0
