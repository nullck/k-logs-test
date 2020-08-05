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
        chmod +x ./kind && mv ./kind /usr/local/bin/kind
      fi
      if [ "${i}" == "kubectl" ]; then
        curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
        chmod +x kubectl && mv ./kubectl /usr/local/bin/kubectl
      fi
    else
      echo "please check your installation"
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
    kind create cluster --name ${CLUSTER_NAME};
    sleep 15;
  fi
  kubectl apply -f fluentbit
  kubectl apply -f elastic
  sleep 5
  while ! kubectl get pods/elasticsearch-0 | grep "Running"; do
    sleep 2
  done
  kubectl port-forward svc/elasticsearch 9200:9200 &
fi

if [ "$1" == "destroy" ]; then
  kind get clusters | grep "${CLUSTER_NAME}"
  if [ $? == 0 ]; then
    kind delete cluster --name ${CLUSTER_NAME}
  fi
fi

exit 0
