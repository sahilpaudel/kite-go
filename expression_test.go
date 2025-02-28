package kite_go

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/Knetic/govaluate.v3"
	"testing"
)

func TestGetExpressionToken(t *testing.T) {
	expression, err := govaluate.NewEvaluableExpression("appVersion>='9.0.1'&&appVersion<='10.0.1'")
	assert.Nil(t, err)
	tokens := GetExpressionToken(expression)
	assert.Equal(t, 2, len(tokens))
	assert.Equal(t, tokens[0].Name, "appVersion")
	assert.Equal(t, tokens[0].Comparator, ">=")
	assert.Equal(t, tokens[0].Operator, "&&")
	assert.Equal(t, tokens[0].Value, "9.0.1")
	assert.Equal(t, tokens[1].Name, "appVersion")
	assert.Equal(t, tokens[1].Comparator, "<=")
	assert.Equal(t, tokens[1].Operator, "")
	assert.Equal(t, tokens[1].Value, "10.0.1")
}

func TestGetVersionsFromRequest(t *testing.T) {

	requestParam := map[string]interface{}{
		"appVersion":  "10.9.1",
		"homeVersion": "10.9.1",
		"platform":    "android",
	}

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "10.9.1",
	}
	variable2 := Variable{
		Name:  "homeVersion",
		Type:  Version,
		Value: "10.9.1",
	}

	var variables []Variable

	variables = append(variables, variable1)
	variables = append(variables, variable2)

	versions := GetVersionsFromRequest(requestParam, variables)
	assert.Len(t, versions, 2)
}

func TestPopulateParameters(t *testing.T) {
	expression, err := govaluate.NewEvaluableExpression("Price>=100&&Age==32&&Bulls==3&&Swine==3.0&&IsValid==true")
	assert.Nil(t, err)

	requestParam := map[string]interface{}{
		"Price": 10.9,
		"Age":   32,
	}

	variable1 := Variable{
		Name:  "Price",
		Type:  Float,
		Value: 10.9,
	}
	variable2 := Variable{
		Name:  "Age",
		Type:  Number,
		Value: 32,
	}
	variable3 := Variable{
		Name: "Bulls",
		Type: Number,
	}
	variable4 := Variable{
		Name: "Swine",
		Type: Float,
	}
	variable5 := Variable{
		Name: "IsValid",
		Type: Boolean,
	}

	var variables []Variable
	variables = append(variables, variable1)
	variables = append(variables, variable2)
	variables = append(variables, variable3)
	variables = append(variables, variable4)
	variables = append(variables, variable5)

	payload := PopulateParamPayload{
		Expression:    expression,
		RequestParams: requestParam,
		Variables:     variables,
	}

	parameters := PopulateParameters(payload)
	assert.NotNil(t, parameters)
	assert.Equal(t, parameters["IsValid"], false)
	assert.Equal(t, parameters["Price"], 10.9)
}

func TestEvaluateVersionExpression(t *testing.T) {
	var versionPayload []VersionEvaluationPayload
	expression, err := EvaluateVersionExpression(versionPayload, nil, nil)
	assert.Nil(t, expression)
	assert.Nil(t, err)
}
