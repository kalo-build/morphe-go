package yaml

import (
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

func (e Entity) Validate(allModels map[string]Model, allEnums map[string]Enum) error {
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

	if err := e.validateAllRelations(); err != nil {
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

func (e Entity) validateAllRelations() error {
	for relatedName, relation := range e.Related {
		if err := e.validateRelation(relatedName, relation); err != nil {
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

func (e Entity) validateRelation(relatedName string, relation EntityRelation) error {
	if relation.Type == "" {
		return ErrNoMorpheEntityRelationType(e.Name, relatedName)
	}

	validTypes := map[string]bool{
		"ForOne":  true,
		"ForMany": true,
		"HasOne":  true,
		"HasMany": true,
	}

	if !validTypes[relation.Type] {
		return ErrInvalidMorpheEntityRelationType(e.Name, relatedName, relation.Type)
	}

	return nil
}
