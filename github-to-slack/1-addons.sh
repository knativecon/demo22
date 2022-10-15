#!/usr/bin/env bash

echo Installing the GitHub Source add-on

kubectl apply -f https://github.com/knative-sandbox/eventing-github/releases/download/knative-v1.7.0/mt-github.yaml

kubectl wait -n knative-sources --for=condition=Ready ksvc/github-adapter
