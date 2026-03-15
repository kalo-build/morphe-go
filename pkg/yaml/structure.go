package yaml

import (
	"github.com/kalo-build/clone"
)

type Structure struct {
	Name   string                    `yaml:"name"`
	Fields map[string]StructureField `yaml:"fields"`
}

func (s Structure) Validate(allEnums map[string]Enum, allStructures ...map[string]Structure) error {
	if s.Name == "" {
		return ErrNoMorpheStructureName
	}
	if len(s.Fields) == 0 {
		return ErrNoMorpheStructureFields
	}
	if len(allEnums) == 0 && len(allStructures) == 0 {
		return nil
	}

	var structures map[string]Structure
	if len(allStructures) > 0 {
		structures = allStructures[0]
	}
	fieldTypesErr := s.validateFieldTypes(allEnums, structures)
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

func (s Structure) validateFieldTypes(allEnums map[string]Enum, allStructures map[string]Structure) error {
	if len(allEnums) == 0 && len(allStructures) == 0 {
		return nil
	}
	for fieldName, fieldDef := range s.Fields {
		fieldType := fieldDef.Type
		if IsStructureFieldTypePrimitive(fieldType) {
			continue
		}

		fieldTypeString := string(fieldType)
		if _, ok := allEnums[fieldTypeString]; ok {
			continue
		}
		if allStructures != nil {
			if _, ok := allStructures[fieldTypeString]; ok {
				continue
			}
		}
		return ErrMorpheStructureUnknownFieldType(fieldName, fieldTypeString)
	}
	return nil
}
