package main

import "github.com/spf13/pflag"

var (
	slackWebhook = pflag.String("slack-webhook", "", "slack webhook URL to post notifications")
	namespace = pflag.String("namespace", "pod-crash", "namespace to be watched by the controller")
)
