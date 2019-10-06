package utils

import (
	"log"
	"testing"
	"io/ioutil"
	"path/filepath"

	"github.com/infracloudio/botkube/pkg/config"
	"github.com/infracloudio/botkube/pkg/notify"
	"github.com/infracloudio/botkube/pkg/utils"
	"github.com/nlopes/slack"
	v1 "k8s.io/api/core/v1"
	extV1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// TESTDATA contains path to the dir where test data is stored
	TESTDATA = "testdata"
	// SAMPLELOGS contains sample logs for test cases
	SAMPLELOGS = "sample-logs.txt"
)

// SlackMessage structure
type SlackMessage struct {
	Text        string
	Attachments []slack.Attachment
}

// WebhookPayload structure
type WebhookPayload struct {
	Summary     string             `json:"summary"`
	EventMeta   notify.EventMeta   `json:"meta"`
	EventStatus notify.EventStatus `json:"status"`
}

// CreateObjects stores specs for creating a k8s fake object and expected Slack response
type CreateObjects struct {
	Kind                   string
	Namespace              string
	Specs                  runtime.Object
	NotifType              config.NotifType
	ExpectedWebhookPayload WebhookPayload
	ExpectedSlackMessage   SlackMessage
}

// CreateResource with fake client
func CreateResource(t *testing.T, obj CreateObjects) {
	switch obj.Kind {
	case "pod":
		s := obj.Specs.(*v1.Pod)
		_, err := utils.KubeClient.CoreV1().Pods(obj.Namespace).Create(s)
		if err != nil {
			t.Fatalf("Failed to create pod: %v", err)
		}
	case "service":
		s := obj.Specs.(*v1.Service)
		_, err := utils.KubeClient.CoreV1().Services(obj.Namespace).Create(s)
		if err != nil {
			t.Fatalf("Failed to create service: %v", err)
		}
	case "ingress":
		s := obj.Specs.(*extV1beta1.Ingress)
		_, err := utils.KubeClient.ExtensionsV1beta1().Ingresses(obj.Namespace).Create(s)
		if err != nil {
			t.Fatalf("Failed to create service: %v", err)
		}
	}
}

// ReadTestData reads content of files in testdata and returns output in string format
func ReadTestData(filename string) string {
	filePath, _ := filepath.Abs(TESTDATA + "/" + filename)
	logs, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read testdata file %s: %v", filePath, err)
	}
	return string(logs)
}
