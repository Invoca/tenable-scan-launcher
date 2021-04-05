package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Invoca/tenable-scan-launcher/pkg/config"
	"github.com/Invoca/tenable-scan-launcher/pkg/tenable"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type SlackInterface interface {
	PrintAlerts(tenable.Alerts) error
}

type markdownText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type block struct {
	BlockType string        `json:"type"`
	BlockText *markdownText `json:"text,omitempty"`
}

type slackBody struct {
	Blocks []block `json:"blocks"`
}

type slack struct {
	slackUrl  string
	rateLimit *rateLimitedHTTPClient
}

type rateLimitedHTTPClient struct {
	client   *http.Client
	rlClient *rate.Limiter
}

func (c *rateLimitedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.rlClient.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func New(config config.BaseConfig) (*slack, error) {
	if config.SlackConfig == nil {
		return nil, fmt.Errorf("Error: SlackConfig cannot be nil")
	}

	if config.SlackConfig.SlackURL == "" {
		return nil, fmt.Errorf("Error: SlackURL cannot be empty")
	}

	s := slack{}
	s.slackUrl = config.SlackConfig.SlackURL
	s.rateLimit = &rateLimitedHTTPClient{
		client:   http.DefaultClient,
		rlClient: rate.NewLimiter(rate.Every(10*time.Second), 10),
	}

	return &s, nil
}

func (s *slack) createBlockSlackPost(text string, additionalText string) error {
	var blockSlice []block

	divider := block{
		BlockType: "divider",
	}

	mainMarkdown := markdownText{
		Type: "mrkdwn",
		Text: text,
	}

	additionalInfoMarkdown := markdownText{
		Type: "mrkdwn",
		Text: additionalText,
	}

	mainBlock := block{
		BlockType: "section",
		BlockText: &mainMarkdown,
	}

	additionalInfoBlock := block{
		BlockType: "section",
		BlockText: &additionalInfoMarkdown,
	}

	blockSlice = append(blockSlice, divider, mainBlock, additionalInfoBlock)

	body := slackBody{Blocks: blockSlice}

	data, _ := json.Marshal(body)

	log.Debug(string(data))

	req, _ := http.NewRequest("POST", s.slackUrl, bytes.NewBuffer(data))
	resp, err := s.rateLimit.Do(req)

	if err != nil {
		return fmt.Errorf("createBlockSlackPost: Error Posting Request %s", err)
	}

	if resp.Status != "200 OK" {
		return fmt.Errorf("createBlockSlackPost: Received non 200 Status Code %s", err)
	}
	log.Debug("Received Status: " + resp.Status)
	return nil
}

func (s *slack) formatLabels(labels map[string]string) string {
	baseString := "Labels:\t"
	formatSpacing := "\n\t\t\t\t"
	firstLabel := true
	for name, value := range labels {
		if firstLabel {
			baseString = baseString + name + ": " + value
			firstLabel = false
			continue
		}
		baseString = baseString + formatSpacing + name + ": " + value
	}
	return baseString
}

func (s *slack) PrintAlerts(alerts tenable.Alerts) error {
	if s.slackUrl == "" {
		return fmt.Errorf("PrintOpenedPorts: slackUrl cannot be empty")
	}

	for _, v := range alerts.Vulnerabilities {
		v.PluginName = strings.ReplaceAll(v.PluginName, "<", " ")
		v.PluginName = strings.ReplaceAll(v.PluginName, ">", " ")
		title := ":red_circle: " + "*" + strconv.FormatUint(uint64(v.Count), 10) + "*" + " instance(s) currently contain(s) the `" + v.VulnerabilityState + "` vulnerability: \n"
		attachmentText := "<https://www.tenable.com/plugins/nessus/" + strconv.FormatUint(uint64(v.PluginID), 10) + "|" + v.PluginName + ">\n"
		err := s.createBlockSlackPost(title, attachmentText)
		if err != nil {
			return fmt.Errorf("PrintAlerts: Error posting message to slack %s", err)
		}
	}
	return nil
}
