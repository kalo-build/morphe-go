package yaml

import (
	"github.com/kalo-build/clone"
)

type Structure struct {
	Name   string                    `yaml:"name"`
	Fields map[string]StructureField `yaml:"fields"`
}

func (s Structure) Validate(allEnums map[string]Enum) error {
	if s.Name == "" {
		return ErrNoMorpheStructureName
	}
	if len(s.Fields) == 0 {
		return ErrNoMorpheStructureFields
	}
	if len(allEnums) == 0 {
		return nil
	}

	fieldTypesErr := s.validateFieldTypes(allEnums)
	if fieldTypesErr != nil {
		return fieldTypesErr
	}

	return nil
}

func (s Structure) DeepClone() Structure {
	structureCopy := Structure{
		Name:   s.Name,
		Fields: clone.DeepCloneMap(s.Fields),
	}

	return structureCopy
}

func (s Structure) validateFieldTypes(allEnums map[string]Enum) error {
	if len(allEnums) == 0 {
		return nil
	}
	for fieldName, fieldDef := range s.Fields {
		fieldType := fieldDef.Type
		if IsStructureFieldTypePrimitive(fieldType) {
			continue
		}

		fieldTypeString := string(fieldType)
		_, enumTypeExists := allEnums[fieldTypeString]
		if !enumTypeExists {
			return ErrMorpheStructureUnknownFieldType(fieldName, fieldTypeString)
		}
	}
	return nil
}
