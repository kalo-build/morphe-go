package yaml

import (
	"fmt"

	"github.com/kalo-build/go-util/core"
)

type Enum struct {
	Name    string         `yaml:"name"`
	Type    EnumType       `yaml:"type"`
	Entries map[string]any `yaml:"entries"`
}

func (e Enum) Validate() error {
	if e.Name == "" {
		return ErrNoMorpheEnumName
	}
	if e.Type == "" {
		return ErrNoMorpheEnumType
	}
	if len(e.Entries) == 0 {
		return ErrNoMorpheEnumEntries
	}

	entryTypesErr := e.validateAllEntryTypes()
	if entryTypesErr != nil {
		return entryTypesErr
	}

	return nil
}

func (e Enum) DeepClone() Enum {
	enumCopy := Enum{
		Name: e.Name,
		Type: e.Type,
	}

	entriesCopy := make(map[string]any, len(e.Entries))
	for key, primitive := range e.Entries {
		entriesCopy[key] = primitive
	}

	enumCopy.Entries = entriesCopy

	return enumCopy
}

func (e Enum) validateAllEntryTypes() error {
	entryNames := core.MapKeysSorted(e.Entries)
	for _, entryName := range entryNames {
		entryValue := e.Entries[entryName]
		validateErr := e.validateEnumEntryValueType(entryName, entryValue)
		if validateErr != nil {
			return validateErr
		}
	}
	return nil
}

func (e Enum) validateEnumEntryValueType(entryName string, entryValue any) error {
	if e.Type != EnumTypeString && e.Type != EnumTypeInteger && e.Type != EnumTypeFloat {
		return fmt.Errorf("enum type '%s' is not supported", e.Type)
	}

	isString := false
	isNumber := false
	switch entryValue.(type) {
	case string:
		isString = true
	case int:
		isNumber = true
	case int8:
		isNumber = true
	case int16:
		isNumber = true
	case int32:
		isNumber = true
	case int64:
		isNumber = true
	case uint:
		isNumber = true
	case uint8:
		isNumber = true
	case uint16:
		isNumber = true
	case uint32:
		isNumber = true
	case uint64:
		isNumber = true
	case float32:
		isNumber = true
	case float64:
		isNumber = true
	default:
		return ErrMorpheEnumUnsupportedEntryType(entryName, entryValue)
	}

	if e.Type == EnumTypeString && !isString {
		return ErrMorpheEnumEntryTypeMismatch(e.Type, entryName, entryValue)
	}

	if (e.Type == EnumTypeInteger || e.Type == EnumTypeFloat) && !isNumber {
		return ErrMorpheEnumEntryTypeMismatch(e.Type, entryName, entryValue)
	}

	return nil
}
