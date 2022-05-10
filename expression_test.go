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
