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

func TestCompareVersionFails(t *testing.T) {
	parsedVersion, _ := goVersion.NewVersion(fmt.Sprint("0.0.0"))
	versions := VersionEvaluationPayload{
		name:          "appVersion",
		parsedVersion: parsedVersion,
	}
	var requestVersions []VersionEvaluationPayload

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "10.9.7",
	}
	var variables []Variable

	variables = append(variables, variable1)

	data := ExpressionVariableMeta{
		Name:       "Version", // name invalid
		Comparator: "<=",
		Operator:   "&&",
		Value:      "0.0.0",
	}

	var versionVariable []ExpressionVariableMeta
	versionVariable = append(versionVariable, data)

	requestVersions = append(requestVersions, versions)
	_, _ = CompareVersion(requestVersions, versionVariable)
}

func TestCompareVersionOperatorEquals(t *testing.T) {
	parsedVersion, _ := goVersion.NewVersion(fmt.Sprint("1.0.0"))
	versions := VersionEvaluationPayload{
		name:          "appVersion",
		parsedVersion: parsedVersion,
	}
	var requestVersions []VersionEvaluationPayload

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "1.0.0",
	}
	var variables []Variable

	variables = append(variables, variable1)

	data := ExpressionVariableMeta{
		Name:       "appVersion",
		Comparator: "==",
		Operator:   "&&",
		Value:      "1.0.0",
	}

	var versionVariable []ExpressionVariableMeta
	versionVariable = append(versionVariable, data)

	requestVersions = append(requestVersions, versions)
	_, _ = CompareVersion(requestVersions, versionVariable)
}

func TestCompareVersionOperatorLessThanEquals(t *testing.T) {
	parsedVersion, _ := goVersion.NewVersion(fmt.Sprint("1.0.0"))
	versions := VersionEvaluationPayload{
		name:          "appVersion",
		parsedVersion: parsedVersion,
	}
	var requestVersions []VersionEvaluationPayload

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "1.0.0",
	}
	var variables []Variable

	variables = append(variables, variable1)

	data := ExpressionVariableMeta{
		Name:       "appVersion",
		Comparator: "<=",
		Operator:   "&&",
		Value:      "1.0.0",
	}

	var versionVariable []ExpressionVariableMeta
	versionVariable = append(versionVariable, data)

	requestVersions = append(requestVersions, versions)
	_, _ = CompareVersion(requestVersions, versionVariable)
}

func TestCompareVersionOperatorGreaterThan(t *testing.T) {
	parsedVersion, _ := goVersion.NewVersion(fmt.Sprint("1.0.0"))
	versions := VersionEvaluationPayload{
		name:          "appVersion",
		parsedVersion: parsedVersion,
	}
	var requestVersions []VersionEvaluationPayload

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "1.0.0",
	}
	var variables []Variable

	variables = append(variables, variable1)

	data := ExpressionVariableMeta{
		Name:       "appVersion",
		Comparator: ">",
		Operator:   "&&",
		Value:      "1.0.0",
	}

	var versionVariable []ExpressionVariableMeta
	versionVariable = append(versionVariable, data)

	requestVersions = append(requestVersions, versions)
	_, _ = CompareVersion(requestVersions, versionVariable)
}

func TestCompareVersionOperatorLessThan(t *testing.T) {
	parsedVersion, _ := goVersion.NewVersion(fmt.Sprint("1.0.0"))
	versions := VersionEvaluationPayload{
		name:          "appVersion",
		parsedVersion: parsedVersion,
	}
	var requestVersions []VersionEvaluationPayload

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "1.0.0",
	}
	var variables []Variable

	variables = append(variables, variable1)

	data := ExpressionVariableMeta{
		Name:       "appVersion",
		Comparator: "<",
		Operator:   "&&",
		Value:      "1.0.0",
	}

	var versionVariable []ExpressionVariableMeta
	versionVariable = append(versionVariable, data)

	requestVersions = append(requestVersions, versions)
	_, _ = CompareVersion(requestVersions, versionVariable)
}

func TestCompareVersionOperatorNotEquals(t *testing.T) {
	parsedVersion, _ := goVersion.NewVersion(fmt.Sprint("1.0.0"))
	versions := VersionEvaluationPayload{
		name:          "appVersion",
		parsedVersion: parsedVersion,
	}
	var requestVersions []VersionEvaluationPayload

	variable1 := Variable{
		Name:  "appVersion",
		Type:  Version,
		Value: "1.0.0",
	}
	var variables []Variable

	variables = append(variables, variable1)

	data := ExpressionVariableMeta{
		Name:       "appVersion",
		Comparator: "!=",
		Operator:   "&&",
		Value:      "1.0.0",
	}

	var versionVariable []ExpressionVariableMeta
	versionVariable = append(versionVariable, data)

	requestVersions = append(requestVersions, versions)
	_, _ = CompareVersion(requestVersions, versionVariable)
}
