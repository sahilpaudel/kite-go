package kite_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluateExpression(t *testing.T) {

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "10.9.7",
	}
	variable2 := Variable{
		Name:  "platform",
		Type:  String,
		Value: "android",
	}
	variable3 := Variable{
		Name:  "segment",
		Type:  Number,
		Value: 98,
	}

	var variables []Variable

	variables = append(variables, variable1)
	variables = append(variables, variable2)
	variables = append(variables, variable3)
	request := map[string]interface{}{
		"appVersion": "10.9.8",
		"segment": 98,
	}

	result, err := EvaluateExpression("appVersion>='10.9.8'&&segment==98", request, variables)
	assert.Nil(t, err)
	assert.True(t, *result)
}