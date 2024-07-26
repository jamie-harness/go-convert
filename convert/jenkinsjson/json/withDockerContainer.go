package json

import (
	"encoding/json"
	"log"

	harness "github.com/drone/spec/dist/go"
)

func ConvertWithDockerContainer(node Node, variables map[string]string) *harness.Step {
	var script, image string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if s, ok := attrMap["script"].(string); ok {
				script = s
			}
			if i, ok := attrMap["image"].(string); ok {
				image = i
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
		return nil
	}

	if script == "" {
		log.Printf("no script found for node %s", node.SpanName)
		return nil
	}

	if image == "" {
		image = "default-image:latest"
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "Run",
		Spec: &harness.StepExec{
			Image: image,
			Run:   script,
		},
	}

	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}

	return step
}