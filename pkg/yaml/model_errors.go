package yaml

import (
	"errors"
	"fmt"
)

var ErrNoMorpheModelName = errors.New("morphe model has no name")
var ErrNoMorpheModelFields = errors.New("morphe model has no fields")
var ErrNoMorpheModelIdentifiers = errors.New("morphe model has no identifiers")

func ErrMorpheModelUnknownFieldType(fieldName string, typeName string) error {
	return fmt.Errorf("morphe model field '%s' has unknown non-primitive type '%s'", fieldName, typeName)
}

func ErrMorpheModelUnknownAliasedTarget(modelName string, relationName string, aliasedTarget string) error {
	return fmt.Errorf("morphe model '%s' relation '%s' has unknown aliased target: %s", modelName, relationName, aliasedTarget)
}

func ErrMorpheModelPolymorphicInverseValidation(modelName string, relationName string, aliasedTarget string, through string, reason string) error {
	return fmt.Errorf("morphe model '%s' polymorphic inverse relation '%s' (aliased: %s, through: %s): %s", modelName, relationName, aliasedTarget, through, reason)
}
