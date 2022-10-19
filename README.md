# Basic Eventing Patterns

This projects explore several basic [Knative Eventing](https://knative.dev/docs/eventing/) patterns.


## Installing

Run these scripts, in order:
1. **0-setup.sh**: create a `kind` cluster and install core Knative Serving and Eventing using
   [`kn quickstart`](https://knative.dev/docs/install/quickstart-install/)
2. **1-addons.sh**: install the required add-ons to run the demo:
   * GitHubSource
   * Strimzi (this is not a Knative add-on)
   * KafkaChannel
   * KafkaBroker
3. **3-deploy.sh**: deploy the demo. This scripts looks for 2 environment variables:
   * GITHUB_TOKEN: your GitHub [personal access token](https://github.com/settings/tokens)
   * SLACK_TOKEN: the [slack token](https://api.slack.com/authentication/oauth-v2) for your slack app.
     See [Slack 101](./doc/slack.md) for more information.

The default configuration listens for events coming from [knativecon/demo22](https://github.com/knativecon/demo22).

## Patterns

### Pattern 1: direct delivery

GitHub comments are directly mirrored to a slack channel.

#### Topology 

```
GitHub -- (filter: issues, issue_comment) -> GitHub Adapter   
                                          -> Slack Sink App 
                                          -> Slack
```

#### Steps

1. Create a GitHub issue. 
   * Title: `There is a bug.` 
   * Body: `It's not working`
2. Observe Slack notifications in `knativecon22-direct`
 

### Pattern 2: queue, ordered

#### Topology

```
GitHub -- (filter: issues, issue_comment) -> GitHub Adapter 
                                          -> Kafka Channel
                                          -> Slack Sink App 
                                          -> Slack
```

### Steps

1. Add a comment to the previous created GitHub issue
   * Body: `first comment`
2. Add another comment (don't wait too long):
   * Body: `second comment`
3. In `knativecon22-direct` slack channel, observe the comments being out-of-order
4. In `knativecon22-channel` slack channel, observe the comments being in-order

**Note**: Both GitHub and the GitHub adapter don't guarantee ordering. 

### Pattern 3: retries (TODO)


GH -> GH Adapter -> SlackSinkAPP (send error) -> Slack

- the slack sink throw an error (cannot connect to backend database)

GH -> GH Adapter -> Kafka (with retry) -> SlackSinkAPP (send error) -> Slack

### Pattern 4: not losing events - using DLS (TODO)

### Pattern 5: using broker vs channel (TODO)

GH1 -(filter: issues, issue_comment) > ... GH Adapter -> Broker - (issues) > aggregations - does not care where it's coming from.
GH2 -(filter: issues, issue_comment) > ... GH Adapter -> Broker - (issues) > collect metrics?

