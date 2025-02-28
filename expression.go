package kite_go

import (
	"fmt"
	goVersion "github.com/hashicorp/go-version"
	"gopkg.in/Knetic/govaluate.v3"
	"reflect"
)

type VersionEvaluationPayload struct {
	name          string
	parsedVersion *goVersion.Version
}

type PopulateParamPayload struct {
	Expression    *govaluate.EvaluableExpression
	RequestParams map[string]interface{}
	Variables     []Variable
}

func GetVersionsFromRequest(requestParam map[string]interface{}, variables []Variable) []string {
	var versions []string
	var versionsDB []string

	for _, variable := range variables {
		if variable.Type == Version {
			versionsDB = append(versionsDB, variable.Name)
		}
	}

	for key := range requestParam {
		if Contains(versionsDB, key) {
			versions = append(versions, key)
		}
	}
	return versions
}

func Contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// GetExpressionToken for every complete expression there is 4 attributes variable, comparator, operator, value
func GetExpressionToken(expression *govaluate.EvaluableExpression) []ExpressionVariableMeta {
	var variables []ExpressionVariableMeta
	variable := ExpressionVariableMeta{}
	tokens := expression.Tokens()
	for k, data := range tokens {
		if k%4 == 0 {
			variable.Name = fmt.Sprint(data.Value)
			continue
		}
		if k%4 == 1 {
			variable.Comparator = data.Value.(string)
			continue
		}
		if k%4 == 2 {
			variable.Value = data.Value
			if k != len(tokens)-1 {
				continue
			}
		}
		if k%4 == 3 {
			variable.Operator = fmt.Sprint(data.Value)
		}
		variables = append(variables, variable)
		variable = ExpressionVariableMeta{}
	}
	return variables
}

func EvaluateVersionExpression(requestVersions []VersionEvaluationPayload, expression *govaluate.EvaluableExpression, variables []Variable) (*govaluate.EvaluableExpression, error) {

	if len(requestVersions) == 0 {
		return expression, nil
	}

	expressionVariables := GetExpressionToken(expression)
	var versionVariable []ExpressionVariableMeta
	for _, newVariable := range expressionVariables {
		if IsTypeVersion(newVariable.Name, variables) {
			versionVariable = append(versionVariable, newVariable)
		}
	}
	result, err := CompareVersion(requestVersions, versionVariable)

	if err != nil {
		return nil, err
	}

	versionVariableIndex := 0
	var finalExpression string
	for _, newVariable := range expressionVariables {
		valueType := reflect.TypeOf(newVariable.Value)
		value := fmt.Sprint(newVariable.Value)
		if valueType.String() == "string" {
			value = "'" + value + "'"
		}

		// not version type
		if !IsTypeVersion(newVariable.Name, variables) {
			finalExpression += newVariable.Name + newVariable.Comparator + value + newVariable.Operator
		}

		// version type
		if IsTypeVersion(newVariable.Name, variables) {
			if result[versionVariableIndex].Name == newVariable.Name {
				finalExpression += fmt.Sprint(result[versionVariableIndex].Result) + newVariable.Operator
			}
			versionVariableIndex++
		}
	}
	return govaluate.NewEvaluableExpression(finalExpression)
}

// PopulateParameters To add missing variables from the request params and populate a default value for those
func PopulateParameters(payload PopulateParamPayload) map[string]interface{} {
	// a map that hold variable as key and value as the type of variable
	variableTypeMap := make(map[string]interface{}, len(payload.Variables))

	// populate the variable key map
	for _, variable := range payload.Variables {
		variableTypeMap[variable.Name] = variable.Type
	}

	// get all the variable from the expression
	vars := payload.Expression.Vars()
	keysFromParams := make([]string, 0, len(payload.RequestParams))

	// iterate a map to fetch all the variable from request param
	for k := range payload.RequestParams {
		keysFromParams = append(keysFromParams, fmt.Sprint(k))
	}

	// prefill the variables not present in the request param but present in expression
	var keysNotPresentInParam []string
	for _, variable := range vars {
		if !Contains(keysFromParams, variable) {
			keysNotPresentInParam = append(keysNotPresentInParam, variable)
		}
	}

	// add default dummy value for missing variables
	for _, key := range keysNotPresentInParam {
		if variableTypeMap[key] == Type(Number) {
			payload.RequestParams[key] = -1
		} else if variableTypeMap[key] == Type(Float) {
			payload.RequestParams[key] = -1.0
		} else if variableTypeMap[key] == Type(Boolean) {
			payload.RequestParams[key] = false
		} else {
			payload.RequestParams[key] = ""
		}
	}

	return payload.RequestParams
}
