package yamlops

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

func GetModelPrimaryIdentifierFieldName(modelDef yaml.Model) (string, error) {
	primaryID, hasPrimaryID := modelDef.Identifiers["primary"]
	if !hasPrimaryID {
		return "", fmt.Errorf("model '%s' has no defined primary identifier", modelDef.Name)
	}
	if len(primaryID.Fields) != 1 {
		return "", fmt.Errorf("model '%s' primary identifier has %v fields, 1 expected", modelDef.Name, len(primaryID.Fields))
	}
	primaryIDFieldName := primaryID.Fields[0]
	return primaryIDFieldName, nil
}

func GetModelFieldDefinitionByName(modelDef yaml.Model, fieldName string) (yaml.ModelField, error) {
	fieldDef, hasField := modelDef.Fields[fieldName]
	if !hasField {
		return yaml.ModelField{}, fmt.Errorf("model '%s' has no defined field '%s'", modelDef.Name, fieldName)
	}
	return fieldDef, nil
}
