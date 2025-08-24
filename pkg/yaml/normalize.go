package yaml

import "strings"

// NormalizeEntity trims whitespace from string fields after unmarshaling
func NormalizeEntity(e *Entity) {
	// Normalize entity name
	e.Name = strings.TrimSpace(e.Name)

	// Normalize fields
	normalizedFields := make(map[string]EntityField)
	for fieldName, field := range e.Fields {
		normalizedFieldName := strings.TrimSpace(fieldName)
		field.Type = ModelFieldPath(strings.TrimSpace(string(field.Type)))

		// Normalize attributes
		normalizedAttributes := make([]string, len(field.Attributes))
		for i, attr := range field.Attributes {
			normalizedAttributes[i] = strings.TrimSpace(attr)
		}
		field.Attributes = normalizedAttributes

		normalizedFields[normalizedFieldName] = field
	}
	e.Fields = normalizedFields

	// Normalize identifiers
	normalizedIdentifiers := make(map[string]EntityIdentifier)
	for idName, identifier := range e.Identifiers {
		normalizedIdName := strings.TrimSpace(idName)

		// Normalize identifier fields
		normalizedIdFields := make([]string, len(identifier.Fields))
		for i, field := range identifier.Fields {
			normalizedIdFields[i] = strings.TrimSpace(field)
		}
		identifier.Fields = normalizedIdFields

		normalizedIdentifiers[normalizedIdName] = identifier
	}
	e.Identifiers = normalizedIdentifiers

	// Normalize relations
	normalizedRelations := make(map[string]EntityRelation)
	for relationName, relation := range e.Related {
		normalizedRelationName := strings.TrimSpace(relationName)
		relation.Aliased = strings.TrimSpace(relation.Aliased)
		relation.Through = strings.TrimSpace(relation.Through)

		// Normalize For field
		normalizedFor := make([]string, len(relation.For))
		for i, forItem := range relation.For {
			normalizedFor[i] = strings.TrimSpace(forItem)
		}
		relation.For = normalizedFor

		normalizedRelations[normalizedRelationName] = relation
	}
	e.Related = normalizedRelations
}

// NormalizeModel trims whitespace from string fields after unmarshaling
func NormalizeModel(m *Model) {
	// Normalize model name
	m.Name = strings.TrimSpace(m.Name)

	// Normalize fields
	normalizedFields := make(map[string]ModelField)
	for fieldName, field := range m.Fields {
		normalizedFieldName := strings.TrimSpace(fieldName)
		field.Type = ModelFieldType(strings.TrimSpace(string(field.Type)))

		// Normalize attributes
		normalizedAttributes := make([]string, len(field.Attributes))
		for i, attr := range field.Attributes {
			normalizedAttributes[i] = strings.TrimSpace(attr)
		}
		field.Attributes = normalizedAttributes

		normalizedFields[normalizedFieldName] = field
	}
	m.Fields = normalizedFields

	// Normalize identifiers
	normalizedIdentifiers := make(map[string]ModelIdentifier)
	for idName, identifier := range m.Identifiers {
		normalizedIdName := strings.TrimSpace(idName)

		// Normalize identifier fields
		normalizedIdFields := make([]string, len(identifier.Fields))
		for i, field := range identifier.Fields {
			normalizedIdFields[i] = strings.TrimSpace(field)
		}
		identifier.Fields = normalizedIdFields

		normalizedIdentifiers[normalizedIdName] = identifier
	}
	m.Identifiers = normalizedIdentifiers

	// Normalize relations
	normalizedRelations := make(map[string]ModelRelation)
	for relationName, relation := range m.Related {
		normalizedRelationName := strings.TrimSpace(relationName)
		relation.Aliased = strings.TrimSpace(relation.Aliased)
		relation.Through = strings.TrimSpace(relation.Through)

		// Normalize For field
		normalizedFor := make([]string, len(relation.For))
		for i, forItem := range relation.For {
			normalizedFor[i] = strings.TrimSpace(forItem)
		}
		relation.For = normalizedFor

		normalizedRelations[normalizedRelationName] = relation
	}
	m.Related = normalizedRelations
}

// NormalizeAllEntities applies normalization to all entities
func NormalizeAllEntities(entities map[string]Entity) {
	for entityName, entity := range entities {
		NormalizeEntity(&entity)
		entities[entityName] = entity
	}
}

// NormalizeAllModels applies normalization to all models
func NormalizeAllModels(models map[string]Model) {
	for modelName, model := range models {
		NormalizeModel(&model)
		models[modelName] = model
	}
}

// NormalizeEnum trims whitespace from string fields after unmarshaling
func NormalizeEnum(e *Enum) {
	// Normalize enum name
	e.Name = strings.TrimSpace(e.Name)

	// Normalize enum type
	e.Type = EnumType(strings.TrimSpace(string(e.Type)))

	// Normalize enum entries (map keys)
	normalizedEntries := make(map[string]any)
	for entryName, entryValue := range e.Entries {
		normalizedEntryName := strings.TrimSpace(entryName)
		normalizedEntries[normalizedEntryName] = entryValue
	}
	e.Entries = normalizedEntries
}

// NormalizeStructure trims whitespace from string fields after unmarshaling
func NormalizeStructure(s *Structure) {
	// Normalize structure name
	s.Name = strings.TrimSpace(s.Name)

	// Normalize fields
	normalizedFields := make(map[string]StructureField)
	for fieldName, field := range s.Fields {
		normalizedFieldName := strings.TrimSpace(fieldName)
		field.Type = StructureFieldType(strings.TrimSpace(string(field.Type)))

		// Normalize attributes
		normalizedAttributes := make([]string, len(field.Attributes))
		for i, attr := range field.Attributes {
			normalizedAttributes[i] = strings.TrimSpace(attr)
		}
		field.Attributes = normalizedAttributes

		normalizedFields[normalizedFieldName] = field
	}
	s.Fields = normalizedFields
}

// NormalizeAllEnums applies normalization to all enums
func NormalizeAllEnums(enums map[string]Enum) {
	for enumName, enum := range enums {
		NormalizeEnum(&enum)
		enums[enumName] = enum
	}
}

// NormalizeAllStructures applies normalization to all structures
func NormalizeAllStructures(structures map[string]Structure) {
	for structureName, structure := range structures {
		NormalizeStructure(&structure)
		structures[structureName] = structure
	}
}
