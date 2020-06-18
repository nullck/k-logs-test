#!/bin/bash

# check if thet kube-logs-test cluster exists
CLUSTER_NAME="kube-logs-test"

if [ "$1" == "start" ]; then
  kind get clusters | grep "${CLUSTER_NAME}"
  if [ $? != 0 ]; then
    kind create cluster --name ${CLUSTER_NAME};
  fi
  kubectl apply -f test-pod.yaml
fi

if [ "$1" == "destroy" ]; then
  kind get clusters | grep "${CLUSTER_NAME}"
  if [ $? == 0 ]; then
    kind delete cluster --name ${CLUSTER_NAME}
  fi
fi

exit 0
