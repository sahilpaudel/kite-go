package kite_go

import (
	"fmt"
	goVersion "github.com/hashicorp/go-version"
)
import "gopkg.in/Knetic/govaluate.v3"
import "github.com/sirupsen/logrus"

func EvaluateExpression(experimentExpression string, requestExpression map[string]interface{}, variables []Variable) (*bool, error) {
	expression, err := govaluate.NewEvaluableExpression(experimentExpression)

	if err != nil {
		logrus.Error("error setting expression ", err)
		return nil, err
	}

	// it can return multiple version types
	versionNames := GetVersionsFromRequest(requestExpression, variables)
	var parsedVersion *goVersion.Version
	var versionsToEvaluate []VersionEvaluationPayload

	for _, versionName := range versionNames {
		parsedVersion, err = goVersion.NewVersion(fmt.Sprint(requestExpression[versionName]))
		if err != nil {
			logrus.Error("error parsing version from request ", err)
			return nil, err
		}
		payload := VersionEvaluationPayload{
			name:          versionName,
			parsedVersion: parsedVersion,
		}
		versionsToEvaluate = append(versionsToEvaluate, payload)
	}

	newExpression, err := EvaluateVersionExpression(versionsToEvaluate, expression, variables)

	if err != nil {
		logrus.Error("error forming new expression ", err)
		return nil, err
	}

	populateParamPayload := PopulateParamPayload{
		expression:    expression,
		requestParams: requestExpression,
		variables:     variables,
	}

	parameters := PopulateParameters(populateParamPayload)
	result, err := newExpression.Evaluate(parameters)
	if err != nil {
		logrus.Error("error evaluating expression ", err)
		return nil, err
	}

	booleanResult := new(bool)
	*booleanResult = result.(bool)
	return booleanResult, nil
}
