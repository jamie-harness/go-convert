package json

import (
	"encoding/json"
	"fmt"
	harness "github.com/drone/spec/dist/go"
	"strconv"
)

type anchoreArguments struct {
	BailOnFail              string `json:"bailOnFail"`
	ForceAnalyze            bool   `json:"forceAnalyze"`
	Name                    string `json:"name"`
	PolicyBundleId          string `json:"policyBundleId"`
	EngineCredentialsId     string `json:"engineCredentialsId"`
	EngineRetries           string `json:"engineRetries"`
	Engineurl               string `json:"engineurl"`
	Engineverify            bool   `json:"engineverify"`
	ExcludeFromBaseImage    bool   `json:"excludeFromBaseImage"`
	BailOnPluginFail        bool   `json:"bailOnPluginFail"`
	AutoSubscribeTagUpdates bool   `json:"autoSubscribeTagUpdates"`
	Engineaccount           string `json:"engineaccount"`
}

func (a *anchoreArguments) UnmarshalJSON(data []byte) error {
	type Alias anchoreArguments
	aux := &struct {
		ForceAnalyze            string `json:"forceAnalyze"`
		Engineverify            bool   `json:"engineverify"`
		ExcludeFromBaseImage    bool   `json:"excludeFromBaseImage"`
		BailOnPluginFail        bool   `json:"bailOnPluginFail"`
		AutoSubscribeTagUpdates bool   `json:"autoSubscribeTagUpdates"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	a.ForceAnalyze, err = strconv.ParseBool(aux.ForceAnalyze)
	if err != nil {
		return err
	}
	return nil
}

func ConvertAnchore(node Node, variables map[string]string) *harness.Step {
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var anchoreAttr struct {
			Delegate struct {
				Arguments anchoreArguments `json:"arguments"`
			} `json:"delegate"`
		}

		err := json.Unmarshal([]byte(attr), &anchoreAttr)
		if err != nil {
			fmt.Println("Error parsing anchore attribute:", err)
			return nil
		}

		envs := map[string]string{
			"ANCHORECTL_FAIL_BASED_ON_RESULTS": anchoreAttr.Delegate.Arguments.BailOnFail,
			"ANCHORECTL_FORCE":                 strconv.FormatBool(anchoreAttr.Delegate.Arguments.ForceAnalyze),
			"ANCHORE_FILE_NAME":                anchoreAttr.Delegate.Arguments.Name,
			"ANCHORECTL_POLICY":                anchoreAttr.Delegate.Arguments.PolicyBundleId,
		}

		// Merge variables into envs
		for k, v := range variables {
			envs[k] = v
		}

		step := &harness.Step{
			Name: node.SpanName,
			Id:   SanitizeForId(node.SpanName, node.SpanId),
			Type: "script",
			Spec: &harness.StepExec{
				Shell: "sh",
				Run:   fmt.Sprintf("curl -sSfL https://anchorectl-releases.anchore.io/anchorectl/install.sh | sh -s -- -b ${HOME}/.local/bin\n                  export PATH=\"${HOME}/.local/bin/:${PATH}\"\n                  anchorectl --version\n                  ANCHORE_IMAGE=$(cat $ANCHORE_FILE_NAME)\n                  anchorectl image add --wait $ANCHORE_IMAGE\n                  anchorectl image vulnerabilities $ANCHORE_IMAGE\n                  anchorectl image check --detail $ANCHORE_IMAGE"),
				Envs:  envs,
			},
		}

		return step
	}
	fmt.Println("Error: harness-attribute not found for anchore step")
	return nil
}
