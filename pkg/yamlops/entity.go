package yamlops

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

func GetEntityPrimaryIdentifierFieldName(entityDef yaml.Entity) (string, error) {
	primaryID, hasPrimaryID := entityDef.Identifiers["primary"]
	if !hasPrimaryID {
		return "", fmt.Errorf("entity '%s' has no defined primary identifier", entityDef.Name)

	}
	if len(primaryID.Fields) != 1 {
		return "", fmt.Errorf("entity '%s' primary identifier has %v fields, 1 expected", entityDef.Name, len(primaryID.Fields))
	}
	primaryIDFieldName := primaryID.Fields[0]
	return primaryIDFieldName, nil
}

func GetEntityFieldDefinitionByName(entityDef yaml.Entity, fieldName string) (yaml.EntityField, error) {
	fieldDef, hasField := entityDef.Fields[fieldName]
	if !hasField {
		return yaml.EntityField{}, fmt.Errorf("entity '%s' has no defined field '%s'", entityDef.Name, fieldName)
	}
	return fieldDef, nil
}
