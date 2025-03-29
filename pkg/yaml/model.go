package yaml

import (
	"log"

	"github.com/kalo-build/clone"
)

type Model struct {
	Name        string                     `yaml:"name"`
	Fields      map[string]ModelField      `yaml:"fields"`
	Identifiers map[string]ModelIdentifier `yaml:"identifiers"`
	Related     map[string]ModelRelation   `yaml:"related"`
}

func (m Model) Validate(allEnums map[string]Enum) error {
	if m.Name == "" {
		return ErrNoMorpheModelName
	}
	if len(m.Fields) == 0 {
		return ErrNoMorpheModelFields
	}
	if len(m.Identifiers) == 0 {
		return ErrNoMorpheModelIdentifiers
	}
	if len(allEnums) == 0 {
		return nil
	}

	fieldTypesErr := m.validateFieldTypes(allEnums)
	if fieldTypesErr != nil {
		return fieldTypesErr
	}

	return nil
}

func (m Model) DeepClone() Model {
	modelCopy := Model{
		Name:        m.Name,
		Fields:      clone.DeepCloneMap(m.Fields),
		Identifiers: clone.DeepCloneMap(m.Identifiers),
		Related:     clone.DeepCloneMap(m.Related),
	}

	return modelCopy
}

func (m Model) GetIdentifierFields() []ModelField {
	var fields []ModelField
	for _, identifier := range m.Identifiers {
		for _, fieldName := range identifier.Fields {
			idField, fieldExists := m.Fields[fieldName]
			if !fieldExists {
				log.Printf("identifier field '%s' does not exist in model '%s'", fieldName, m.Name)
				continue
			}
			fields = append(fields, idField)
		}
	}
	return fields
}

func (m Model) validateFieldTypes(allEnums map[string]Enum) error {
	if len(allEnums) == 0 {
		return nil
	}
	for fieldName, fieldDef := range m.Fields {
		fieldType := fieldDef.Type
		if IsModelFieldTypePrimitive(fieldType) {
			continue
		}

		fieldTypeString := string(fieldType)
		_, enumTypeExists := allEnums[fieldTypeString]
		if !enumTypeExists {
			return ErrMorpheModelUnknownFieldType(fieldName, fieldTypeString)
		}
	}
	return nil
}
