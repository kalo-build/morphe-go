package yaml

import "fmt"

var ErrNoMorpheEntityName = fmt.Errorf("morphe entity has no name")

func ErrNoMorpheEntityFields(entityName string) error {
	return fmt.Errorf("morphe entity %s has no fields", entityName)
}

func ErrNoMorpheEntityFieldType(entityName string, fieldName string) error {
	return fmt.Errorf("morphe entity %s field %s has no type", entityName, fieldName)
}

func ErrInvalidMorpheEntityFieldTypePath(entityName string, fieldName string, fieldPath string) error {
	return fmt.Errorf("morphe entity %s field %s has invalid type path: %s", entityName, fieldName, fieldPath)
}

func ErrUnknownMorpheEntityFieldRootModel(entityName string, fieldName string, rootModelName string) error {
	return fmt.Errorf("morphe entity %s field %s references unknown root model: %s", entityName, fieldName, rootModelName)
}

func ErrUnknownMorpheEntityFieldRelatedModel(entityName string, fieldName string, relatedName string, fieldType ModelFieldPath) error {
	return fmt.Errorf("morphe entity %s field %s references unknown related model: %s in path %s", entityName, fieldName, relatedName, fieldType)
}

func ErrUnknownMorpheEntityFieldModel(entityName string, fieldName string, modelName string, fieldType ModelFieldPath) error {
	return fmt.Errorf("morphe entity %s field %s references unknown model: %s in path %s", entityName, fieldName, modelName, fieldType)
}

func ErrUnknownMorpheEntityFieldTerminalField(entityName string, fieldName string, terminalFieldName string, fieldType ModelFieldPath) error {
	return fmt.Errorf("morphe entity %s field %s references unknown terminal field: %s in path %s", entityName, fieldName, terminalFieldName, fieldType)
}

func ErrNoMorpheEntityRelationType(entityName string, relatedName string) error {
	return fmt.Errorf("morphe entity %s relation %s has no type", entityName, relatedName)
}

func ErrInvalidMorpheEntityRelationType(entityName string, relatedName string, relationType string) error {
	return fmt.Errorf("morphe entity %s relation %s has invalid type: %s", entityName, relatedName, relationType)
}

func ErrUnknownMorpheEntityFieldType(entityName string, fieldName string, typeName string) error {
	return fmt.Errorf("morphe entity '%s' field '%s' has unknown non-primitive type '%s'", entityName, fieldName, typeName)
}

func ErrNoMorpheEntityIdentifiers(entityName string) error {
	return fmt.Errorf("entity '%s' has no identifiers", entityName)
}

func ErrNoMorpheEntityIdentifierFields(entityName string, identifierName string) error {
	return fmt.Errorf("entity '%s' identifier '%s' has no fields", entityName, identifierName)
}

func ErrUnknownMorpheEntityIdentifierField(entityName string, identifierName string, fieldName string) error {
	return fmt.Errorf("entity '%s' identifier '%s' references unknown field '%s'", entityName, identifierName, fieldName)
}

func ErrMorpheEntityPolyRelationMissingFor(entityName string, relatedName string, relationType string) error {
	return fmt.Errorf("morphe entity %s polymorphic relation %s of type %s is missing required 'for' property", entityName, relatedName, relationType)
}

func ErrMorpheEntityPolyRelationMissingThrough(entityName string, relatedName string, relationType string) error {
	return fmt.Errorf("morphe entity %s polymorphic relation %s of type %s is missing required 'through' property", entityName, relatedName, relationType)
}

func ErrMorpheEntityUnknownAliasedTarget(entityName string, relationName string, aliasedTarget string) error {
	return fmt.Errorf("morphe entity '%s' relation '%s' has unknown aliased target: %s", entityName, relationName, aliasedTarget)
}

func ErrMorpheEntityPolymorphicInverseValidation(entityName string, relationName string, aliasedTarget string, through string, reason string) error {
	return fmt.Errorf("morphe entity '%s' polymorphic inverse relation '%s' (aliased: %s, through: %s): %s", entityName, relationName, aliasedTarget, through, reason)
}
