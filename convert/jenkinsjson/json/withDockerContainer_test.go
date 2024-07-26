package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertWithDockerContainer(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/withDockerContainer/withDockerContainerSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json Node
		want Node
	}{
		{
			json: node1,
			want: Node{
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "12",
					"jenkins.pipeline.step.name":           "Shell Script",
					"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
					"jenkins.pipeline.step.plugin.version": "1360.v82d13453da_a_f",
					"jenkins.pipeline.step.type":           "sh",
					"harness-attribute":                    "{\n  \"script\" : \"docker pull \\\"$JD_TO_PULL\\\"\"\n}",
					"harness-others":                       "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
				},
				Name:         "docker declarative #4",
				Parent:       "docker declarative",
				ParentSpanId: "1f73619e05fc8256",
				SpanId:       "17209b052752dad3",
				SpanName:     "sh",
				TraceId:      "9bbba51ab06f5084ea5020af759c1301",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"script": "docker pull \"$JD_TO_PULL\""},
			},
		},
	}

	for i, test := range tests {
		got := test.json
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}
