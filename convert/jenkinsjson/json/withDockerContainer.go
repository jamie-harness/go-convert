package json

import (
	"encoding/json"
	"log"
	harness "github.com/drone/spec/dist/go"
)

func ConvertWithDockerContainer(node Node, variables map[string]string) *harness.Step {
	var args, image string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if a, ok := attrMap["args"].(string); ok {
				args = a
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

	if args == "" {
		log.Printf("no args found for node %s", node.SpanName)
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
			Run:   args,
		},
	}

	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}

	return step
}