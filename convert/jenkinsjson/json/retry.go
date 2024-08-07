package json

import (
	"encoding/json"
)

func ExtractingRetryCount(node Node) int {
	var count int
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["count"].(int); ok {
				count = f
			}
		}
	}
	return count
}
