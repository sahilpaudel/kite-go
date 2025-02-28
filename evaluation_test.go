package kite_go

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/Knetic/govaluate.v3"
	"testing"
)

type MockPopulator struct {
	mock.Mock
}

func (m *MockPopulator) EvaluateVersionExpression(requestVersions []VersionEvaluationPayload, expression *govaluate.EvaluableExpression, variables []Variable) (*govaluate.EvaluableExpression, error) {
	args := m.Called(requestVersions, expression, variables)
	var evalExpr *govaluate.EvaluableExpression
	var err error
	if arg0 := args.Get(0); arg0 != nil {
		evalExpr, _ = arg0.(*govaluate.EvaluableExpression)
	}

	if len(args) == 2 {
		err = args.Error(1)
	}

	return evalExpr, err
}

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
		"segment":    98,
	}

	result, err := EvaluateExpression("appVersion>='10.9.8'&&segment==98", request, variables)
	assert.Nil(t, err)
	assert.True(t, *result)
}

func TestEvaluateExpressionFails(t *testing.T) {
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
		"segment":    98,
		"bullocks":   32524,
	}

	_, err := EvaluateExpression("appVersion>>='10.9.8'", request, variables)
	assert.NotNil(t, err)

	_, _ = EvaluateExpression("appVersion>='invalid'", request, variables)

	_, _ = EvaluateExpression("ThisVersionKeyIsNotValid>='invalid'", request, variables)

	_, _ = EvaluateExpression("appVersion>=0.0.0", request, variables)

	request = map[string]interface{}{
		"appVersion": "-",
	}
	_, _ = EvaluateExpression("appVersion>='10.9.8'", request, variables)
}
