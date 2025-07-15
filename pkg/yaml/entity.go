package yaml

import (
	"fmt"
	"strings"

	"github.com/kalo-build/clone"
)

type Entity struct {
	Name        string                      `yaml:"name"`
	Fields      map[string]EntityField      `yaml:"fields"`
	Identifiers map[string]EntityIdentifier `yaml:"identifiers"`
	Related     map[string]EntityRelation   `yaml:"related"`
}

func (e Entity) DeepClone() Entity {
	entityCopy := Entity{
		Name:        e.Name,
		Fields:      clone.DeepCloneMap(e.Fields),
		Identifiers: clone.DeepCloneMap(e.Identifiers),
		Related:     clone.DeepCloneMap(e.Related),
	}

	return entityCopy
}

func (e Entity) Validate(allEntities map[string]Entity, allModels map[string]Model, allEnums map[string]Enum) error {
	if e.Name == "" {
		return ErrNoMorpheEntityName
	}

	if len(e.Fields) == 0 {
		return ErrNoMorpheEntityFields(e.Name)
	}

	if len(e.Identifiers) == 0 {
		return ErrNoMorpheEntityIdentifiers(e.Name)
	}

	if err := e.validateAllFieldTypes(allModels, allEnums); err != nil {
		return err
	}

	if err := e.validateAllIdentifiers(); err != nil {
		return err
	}

	if err := e.validateAllRelations(allEntities, allModels); err != nil {
		return err
	}

	return nil
}

func (e Entity) validateAllIdentifiers() error {
	for identifierName, identifier := range e.Identifiers {
		if len(identifier.Fields) == 0 {
			return ErrNoMorpheEntityIdentifierFields(e.Name, identifierName)
		}
		for _, fieldName := range identifier.Fields {
			if _, exists := e.Fields[fieldName]; !exists {
				return ErrUnknownMorpheEntityIdentifierField(e.Name, identifierName, fieldName)
			}
		}
	}
	return nil
}

func (e Entity) validateAllFieldTypes(allModels map[string]Model, allEnums map[string]Enum) error {
	for fieldName, field := range e.Fields {
		if err := e.validateFieldType(fieldName, field, allModels, allEnums); err != nil {
			return err
		}
	}
	return nil
}

func (e Entity) validateAllRelations(allEntities map[string]Entity, allModels map[string]Model) error {
	for relatedName, relation := range e.Related {
		if err := e.validateRelation(relatedName, relation, allEntities, allModels); err != nil {
			return err
		}
	}
	return nil
}

func (e Entity) validateFieldType(fieldName string, field EntityField, allModels map[string]Model, allEnums map[string]Enum) error {
	if field.Type == "" {
		return ErrNoMorpheEntityFieldType(e.Name, fieldName)
	}

	fieldPath := e.parseFieldTypePath(field.Type)
	if pathValidationErr := e.validateFieldTypePath(fieldPath, fieldName); pathValidationErr != nil {
		return pathValidationErr
	}

	rootModel, rootModelErr := e.resolveRootModel(fieldPath[0], fieldName, allModels)
	if rootModelErr != nil {
		return rootModelErr
	}

	currentModel, modelPathErr := e.resolveModelFieldPath(rootModel, fieldPath[1:len(fieldPath)-1], fieldName, field.Type, allModels)
	if modelPathErr != nil {
		return modelPathErr
	}

	if terminalFieldErr := e.validateTerminalField(currentModel, fieldPath[len(fieldPath)-1], fieldName, field.Type, allEnums); terminalFieldErr != nil {
		return terminalFieldErr
	}

	return nil
}

func (e Entity) parseFieldTypePath(fieldType ModelFieldPath) []string {
	return strings.Split(string(fieldType), ".")
}

func (e Entity) validateFieldTypePath(fieldPath []string, fieldName string) error {
	if len(fieldPath) < 2 {
		return ErrInvalidMorpheEntityFieldTypePath(e.Name, fieldName, strings.Join(fieldPath, "."))
	}
	return nil
}

func (e Entity) resolveRootModel(rootModelName string, fieldName string, allModels map[string]Model) (Model, error) {
	rootModel, exists := allModels[rootModelName]
	if !exists {
		return Model{}, ErrUnknownMorpheEntityFieldRootModel(e.Name, fieldName, rootModelName)
	}
	return rootModel, nil
}

func (e Entity) resolveModelFieldPath(startModel Model, pathSegments []string, fieldName string, fieldType ModelFieldPath, allModels map[string]Model) (Model, error) {
	currentModel := startModel
	for _, relatedName := range pathSegments {
		if relationValidationErr := e.validateModelRelation(currentModel, relatedName, fieldName, fieldType); relationValidationErr != nil {
			return Model{}, relationValidationErr
		}

		nextModel, relatedModelErr := e.resolveRelatedModel(relatedName, fieldName, fieldType, allModels)
		if relatedModelErr != nil {
			return Model{}, relatedModelErr
		}
		currentModel = nextModel
	}
	return currentModel, nil
}

func (e Entity) validateModelRelation(model Model, relatedName string, fieldName string, fieldType ModelFieldPath) error {
	if _, exists := model.Related[relatedName]; !exists {
		return ErrUnknownMorpheEntityFieldRelatedModel(e.Name, fieldName, relatedName, fieldType)
	}
	return nil
}

func (e Entity) resolveRelatedModel(relatedName string, fieldName string, fieldType ModelFieldPath, allModels map[string]Model) (Model, error) {
	relatedModel, exists := allModels[relatedName]
	if !exists {
		return Model{}, ErrUnknownMorpheEntityFieldModel(e.Name, fieldName, relatedName, fieldType)
	}
	return relatedModel, nil
}

func (e Entity) validateTerminalField(model Model, fieldName string, originalFieldName string, fieldType ModelFieldPath, allEnums map[string]Enum) error {
	terminalField, exists := model.Fields[fieldName]
	if !exists {
		return ErrUnknownMorpheEntityFieldTerminalField(e.Name, originalFieldName, fieldName, fieldType)
	}
	if IsModelFieldTypePrimitive(terminalField.Type) {
		return nil
	}

	terminalFieldTypeString := string(terminalField.Type)
	_, enumExists := allEnums[terminalFieldTypeString]
	if !enumExists {
		return ErrUnknownMorpheEntityFieldType(e.Name, fieldName, terminalFieldTypeString)
	}

	return nil
}

func (e Entity) validateRelation(relatedName string, relation EntityRelation, allEntities map[string]Entity, allModels map[string]Model) error {
	if relation.Type == "" {
		return ErrNoMorpheEntityRelationType(e.Name, relatedName)
	}

	validTypes := map[string]bool{
		"ForOne":      true,
		"ForMany":     true,
		"HasOne":      true,
		"HasMany":     true,
		"ForOnePoly":  true,
		"ForManyPoly": true,
		"HasOnePoly":  true,
		"HasManyPoly": true,
	}

	if !validTypes[relation.Type] {
		return ErrInvalidMorpheEntityRelationType(e.Name, relatedName, relation.Type)
	}

	// Validate aliased relationships
	if relation.Aliased != "" {
		// Check if the aliased target exists in the entity registry
		if _, entityExists := allEntities[relation.Aliased]; !entityExists {
			return ErrMorpheEntityUnknownAliasedTarget(e.Name, relatedName, relation.Aliased)
		}

		// Enhanced validation for polymorphic inverse relationships
		if e.isRelationPolyHas(relation.Type) && relation.Through != "" {
			// This is a HasOnePoly/HasManyPoly with through + aliased pattern
			if err := e.validatePolymorphicInverseAliasing(relatedName, relation, allModels, allEntities); err != nil {
				return err
			}
		}
	}

	// Validate polymorphic relationships
	lowerType := strings.ToLower(relation.Type)
	isPoly := strings.HasSuffix(lowerType, "poly")

	if !isPoly {
		return nil
	}

	isFor := strings.HasPrefix(lowerType, "for")
	if isFor && len(relation.For) == 0 {
		return ErrMorpheEntityPolyRelationMissingFor(e.Name, relatedName, relation.Type)
	}

	isHas := strings.HasPrefix(lowerType, "has")
	if isHas && relation.Through == "" {
		return ErrMorpheEntityPolyRelationMissingThrough(e.Name, relatedName, relation.Type)
	}

	return nil
}

func (e Entity) validatePolymorphicInverseAliasing(relationName string, relation EntityRelation, allModels map[string]Model, allEntities map[string]Entity) error {
	aliasedEntity := allEntities[relation.Aliased]

	// Check if aliased entity has the 'through' relationship
	throughRelation, exists := aliasedEntity.Related[relation.Through]
	if !exists {
		return ErrMorpheEntityPolymorphicInverseValidation(e.Name, relationName, relation.Aliased, relation.Through,
			fmt.Sprintf("aliased entity '%s' does not have relationship '%s'", relation.Aliased, relation.Through))
	}

	// Check if the 'through' relationship is polymorphic ForOnePoly/ForManyPoly
	if !e.isRelationPolyFor(throughRelation.Type) {
		return ErrMorpheEntityPolymorphicInverseValidation(e.Name, relationName, relation.Aliased, relation.Through,
			fmt.Sprintf("relationship '%s' in entity '%s' is not a polymorphic 'For' relationship (type: %s)", relation.Through, relation.Aliased, throughRelation.Type))
	}

	// Check if the current entity is in the 'for' list of the polymorphic relationship
	found := false
	for _, forEntity := range throughRelation.For {
		if forEntity == e.Name {
			found = true
			break
		}
	}
	if !found {
		return ErrMorpheEntityPolymorphicInverseValidation(e.Name, relationName, relation.Aliased, relation.Through,
			fmt.Sprintf("entity '%s' is not in the 'for' list of polymorphic relationship '%s' in entity '%s'", e.Name, relation.Through, relation.Aliased))
	}

	return nil
}

// Helper functions for entity relation validation
func (e Entity) isRelationFor(relationType string) bool {
	return strings.HasPrefix(strings.ToLower(relationType), "for")
}

func (e Entity) isRelationHas(relationType string) bool {
	return strings.HasPrefix(strings.ToLower(relationType), "has")
}

func (e Entity) isRelationPoly(relationType string) bool {
	lowerType := strings.ToLower(relationType)
	return (e.isRelationFor(relationType) || e.isRelationHas(relationType)) &&
		strings.HasSuffix(lowerType, "poly")
}

func (e Entity) isRelationPolyFor(relationType string) bool {
	return e.isRelationPoly(relationType) && e.isRelationFor(relationType)
}

func (e Entity) isRelationPolyHas(relationType string) bool {
	return e.isRelationPoly(relationType) && e.isRelationHas(relationType)
}
