package yaml

import (
	"fmt"
	"log"
	"strings"

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

// ValidateWithModels validates a model with access to all models for aliasing validation
func (m Model) ValidateWithModels(allModels map[string]Model, allEnums map[string]Enum) error {
	// First run the basic validation
	if err := m.Validate(allEnums); err != nil {
		return err
	}

	// Validate aliased relationships
	if err := m.validateAliasedRelations(allModels); err != nil {
		return err
	}

	return nil
}

func (m Model) validateAliasedRelations(allModels map[string]Model) error {
	for relationName, relation := range m.Related {
		if relation.Aliased != "" {
			// Check if the aliased target exists in the registry
			if _, exists := allModels[relation.Aliased]; !exists {
				return ErrMorpheModelUnknownAliasedTarget(m.Name, relationName, relation.Aliased)
			}

			// Enhanced validation for polymorphic inverse relationships
			if m.isRelationPolyHas(relation.Type) && relation.Through != "" {
				// This is a HasOnePoly/HasManyPoly with through + aliased pattern
				if err := m.validatePolymorphicInverseAliasing(relationName, relation, allModels); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m Model) validatePolymorphicInverseAliasing(relationName string, relation ModelRelation, allModels map[string]Model) error {
	aliasedModel := allModels[relation.Aliased]

	// Check if aliased model has the 'through' relationship
	throughRelation, exists := aliasedModel.Related[relation.Through]
	if !exists {
		return ErrMorpheModelPolymorphicInverseValidation(m.Name, relationName, relation.Aliased, relation.Through,
			fmt.Sprintf("aliased model '%s' does not have relationship '%s'", relation.Aliased, relation.Through))
	}

	// Check if the 'through' relationship is polymorphic ForOnePoly/ForManyPoly
	if !m.isRelationPolyFor(throughRelation.Type) {
		return ErrMorpheModelPolymorphicInverseValidation(m.Name, relationName, relation.Aliased, relation.Through,
			fmt.Sprintf("relationship '%s' in model '%s' is not a polymorphic 'For' relationship (type: %s)", relation.Through, relation.Aliased, throughRelation.Type))
	}

	// Check if current model is in the 'for' list of the polymorphic relationship
	found := false
	for _, forModel := range throughRelation.For {
		if forModel == m.Name {
			found = true
			break
		}
	}
	if !found {
		return ErrMorpheModelPolymorphicInverseValidation(m.Name, relationName, relation.Aliased, relation.Through,
			fmt.Sprintf("current model '%s' is not in the 'for' list of polymorphic relationship '%s' in model '%s'", m.Name, relation.Through, relation.Aliased))
	}

	return nil
}

// Helper functions copied from yamlops to avoid import cycle
func (m Model) isRelationFor(relationType string) bool {
	return strings.HasPrefix(strings.ToLower(relationType), "for")
}

func (m Model) isRelationHas(relationType string) bool {
	return strings.HasPrefix(strings.ToLower(relationType), "has")
}

func (m Model) isRelationPoly(relationType string) bool {
	lowerType := strings.ToLower(relationType)
	return (m.isRelationFor(relationType) || m.isRelationHas(relationType)) &&
		strings.HasSuffix(lowerType, "poly")
}

func (m Model) isRelationPolyFor(relationType string) bool {
	return m.isRelationPoly(relationType) && m.isRelationFor(relationType)
}

func (m Model) isRelationPolyHas(relationType string) bool {
	return m.isRelationPoly(relationType) && m.isRelationHas(relationType)
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
