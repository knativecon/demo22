#!/usr/bin/env bash

export KO_DOCKER_REPO=kind.local
export KIND_CLUSTER_NAME=demo1

echo Installing the GitHub source add-on

kubectl apply -f https://github.com/knative-sandbox/eventing-github/releases/download/knative-v1.7.0/mt-github.yaml

echo Waiting fot GitHub source to be ready
kubectl wait -n knative-sources --for=condition=Ready ksvc/github-adapter

echo Installing Apache Kafka
kubectl create namespace kafka

# operator
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
kubectl wait -n kafka --for=condition=Available deployment/strimzi-cluster-operator

# kafka
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n kafka
kubectl wait -n kafka --for=condition=Ready kafka/my-cluster --timeout=300s

echo Installing the Kafka channel and broker add-on
kubectl apply -f https://github.com/knative-sandbox/eventing-kafka-broker/releases/download/knative-v1.7.4/eventing-kafka.yaml


