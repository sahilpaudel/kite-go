package kite_go

import (
"fmt"
goVersion "github.com/hashicorp/go-version"
"github.com/sirupsen/logrus"
)

type ExpressionVariableMeta struct {
	Name       string
	Comparator string
	Operator   string
	Value      interface{}
}

type VersionResult struct {
	Name   string
	Result bool
}

func IsTypeVersion(variableName string, ExpressionVariables []Variable) bool  {
	for _, variable := range ExpressionVariables {
		if variable.Type == Version && variable.Name == variableName {
			return true
		}
	}
	return false
}

func CompareVersion(comparableVersions []VersionEvaluationPayload, variables []ExpressionVariableMeta) ([]VersionResult, error) {

	var nameVersionMap = make(map[string]*goVersion.Version, len(comparableVersions))

	for _, comparableVersion := range comparableVersions{
		nameVersionMap[comparableVersion.name] = comparableVersion.parsedVersion
	}

	var results []VersionResult
	for _, variable := range variables {
		newVariable, err := goVersion.NewVersion(fmt.Sprint(variable.Value))
		if err != nil {
			logrus.Error("Unable to parse version in the expression")
			return results, err
		}

		if nameVersionMap[variable.Name] == nil {
			results = append(results, VersionResult{
				Name:   variable.Name,
				Result: false,
			})
			continue
		}

		switch variable.Comparator {
		case "==":
			results = append(results, VersionResult{
				Name:   variable.Name,
				Result: nameVersionMap[variable.Name].Equal(newVariable),
			})
		case ">=":
			results = append(results, VersionResult{
				Name:   variable.Name,
				Result: nameVersionMap[variable.Name].GreaterThanOrEqual(newVariable),
			})
		case "<=":
			results = append(results, VersionResult{
				Name:   variable.Name,
				Result: nameVersionMap[variable.Name].LessThanOrEqual(newVariable),
			})
		case ">":
			results = append(results, VersionResult{
				Name:   variable.Name,
				Result: nameVersionMap[variable.Name].GreaterThan(newVariable),
			})
		case "<":
			results = append(results, VersionResult{
				Name:   variable.Name,
				Result: nameVersionMap[variable.Name].LessThan(newVariable),
			})
		case "!=":
			results = append(results, VersionResult{
				Name:   variable.Name,
				Result: !nameVersionMap[variable.Name].Equal(newVariable),
			})
		}
	}
	return results, nil
}
