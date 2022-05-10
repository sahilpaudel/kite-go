package kite_go

import (
	"fmt"
	goVersion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"gopkg.in/Knetic/govaluate.v3"
	"testing"
)

func TestIsTypeVersion(t *testing.T) {
	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "10.9.7",
	}
	variable2 := Variable{
		Name:  "platform",
		Type:  String,
		Value: "12.1.10",
	}
	var variables []Variable

	variables = append(variables, variable1)
	variables = append(variables, variable2)

	assert.True(t, IsTypeVersion("appVersion", variables))
	assert.False(t, IsTypeVersion("platform", variables))
}

func TestCompareVersion(t *testing.T) {
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

	var variables []Variable

	variables = append(variables, variable1)
	variables = append(variables, variable2)

	expression, err := govaluate.NewEvaluableExpression("appVersion>='10.9.7'")
	assert.Nil(t, err)

	expressionVariables := GetExpressionToken(expression)

	var versionVariable []ExpressionVariableMeta
	for _, newVariable := range expressionVariables {
		if IsTypeVersion(newVariable.Name, variables) {
			versionVariable = append(versionVariable, newVariable)
		}
	}

	parsedVersion, err := goVersion.NewVersion(fmt.Sprint(variable1.Value))
	assert.Nil(t, err)

	versions := VersionEvaluationPayload{
		name:          "appVersion",
		parsedVersion: parsedVersion,
	}
	var requestVersions []VersionEvaluationPayload
	requestVersions = append(requestVersions, versions)

	result, err := CompareVersion(requestVersions, versionVariable)
	assert.Nil(t, err)
	assert.Equal(t, "appVersion", result[0].Name)
	assert.Equal(t, true, result[0].Result)
}
