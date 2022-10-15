package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "strings"

    cloudevents "github.com/cloudevents/sdk-go/v2"
    "github.com/google/go-github/v47/github"
    "github.com/slack-go/slack"
)

var (
    slackapi *slack.Client
    channel  string
)

func threadTitle(number int, title, name string) string {
    return fmt.Sprintf("%s #%d by %s", title, number, name)
}

func sendToSlack(ceevent cloudevents.Event) error {
    messageType := strings.TrimPrefix(ceevent.Type(), "dev.knative.source.github.")
    ghevent, err := github.ParseWebHook(messageType, ceevent.Data())
    if err != nil {
        log.Printf("failed to parse github event: %v", err)
        return err
    }

    issueNumber := 0
    header := ""
    name := ""
    comment := ""
    iconURL := ""
    switch event := ghevent.(type) {
    case *github.IssuesEvent:
        issueNumber = event.GetIssue().GetNumber()
        title := event.GetIssue().GetTitle()
        name = event.GetSender().GetLogin()
        header = threadTitle(issueNumber, title, name)
        comment = event.GetIssue().GetBody()
        iconURL = event.GetSender().GetAvatarURL()
    case *github.IssueEvent:
        issueNumber = event.GetIssue().GetNumber()
        title := event.GetIssue().GetTitle()
        name = event.GetActor().GetLogin()
        header = threadTitle(issueNumber, title, name)
        comment = event.GetIssue().GetBody()
        iconURL = event.GetActor().GetAvatarURL()
    }

    // Ensure slack thread header exists and up-to-date
    thread, err := getThread(header)
    if err != nil {
        log.Printf("failed to get slack threads: %v", err)
        return err
    }

    if thread == "" {
        thread, err = createThread(name, header)
        if err != nil {
            log.Printf("failed to create slack thread: %v", err)
            return err
        }
    }

    // Add (TODO: update) comment in thread
    if comment != "" {
        options := []slack.MsgOption{
            slack.MsgOptionText(comment, false),
            slack.MsgOptionUsername(name),
            slack.MsgOptionIconURL(iconURL),
            slack.MsgOptionTS(thread),
        }

        _, _, err = slackapi.PostMessage(channel, options...)

        if err != nil {
            log.Printf("failed to post slack message: %v", err)
            return err
        }
    }

    return nil
}

func getThread(text string) (string, error) {
    // TODO: iterate
    resp, err := slackapi.GetConversationHistory(&slack.GetConversationHistoryParameters{
        ChannelID: channel,
    })

    if err != nil {
        return "", err
    }

    for _, msg := range resp.Messages {
        if msg.Text == text {
            return msg.Timestamp, nil
        }
    }

    return "", nil
}

func createThread(user string, header string) (string, error) {
    options := []slack.MsgOption{
        slack.MsgOptionText(header, false),
        slack.MsgOptionAsUser(true),
        slack.MsgOptionUsername(user),
    }

    _, ts, err := slackapi.PostMessage(channel, options...)
    return ts, err
}

func main() {
    run(context.Background())
}

func run(ctx context.Context) {
    slackToken := os.Getenv("SLACK_TOKEN")
    if slackToken == "" {
        log.Fatal("missing SLACK_TOKEN")
    }
    channel = os.Getenv("SLACK_CHANNEL_ID")
    if channel == "" {
        log.Fatal("missing SLACK_CHANNEL_ID")
    }

    slackapi = slack.New(slackToken)

    c, err := cloudevents.NewClientHTTP()
    if err != nil {
        log.Fatal("Failed to create client: ", err)
    }
    if err := c.StartReceiver(ctx, sendToSlack); err != nil {
        log.Fatal("Error during receiver's runtime: ", err)
    }
}
