#!/usr/bin/env bash

echo Deploying application

if [ -z ${GITHUB_TOKEN+x} ]; then
  echo "error: undefined \$GITHUB_TOKEN"
  exit 1
fi

if [ -z ${GITHUB_SECRET+x} ]; then
  echo "error: undefined \$GITHUB_SECRET"
  exit 1
fi

if [ -z ${SLACK_TOKEN+x} ]; then
  echo "error: undefined \$SLACK_TOKEN"
  exit 1
fi

export KO_DOCKER_REPO=kind.local
export KIND_CLUSTER_NAME=demo1

cat config/templates/github-secret.yaml\
 | sed "s/ACCESS_TOKEN/${GITHUB_TOKEN}/"\
 | sed "s/SECRET_TOKEN/${GITHUB_SECRET}/"\
 | kubectl apply --filename -

cat config/templates/slack-secret.yaml\
 | sed "s/OAUTH_TOKEN/${SLACK_TOKEN}/"\
 | kubectl apply --filename -


ko apply -f config
