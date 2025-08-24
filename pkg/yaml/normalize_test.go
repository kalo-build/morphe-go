package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeEntity_TrimsWhitespace(t *testing.T) {
	entity := Entity{
		Name: "  TestEntity  ", // whitespace around name
		Fields: map[string]EntityField{
			"  Email  ": { // whitespace around field name
				Type:       "  ContactInfo.Email  ",                // whitespace around field type
				Attributes: []string{"  required  ", "  unique  "}, // whitespace around attributes
			},
		},
		Identifiers: map[string]EntityIdentifier{
			"  Primary  ": { // whitespace around identifier name
				Fields: []string{"  ID  ", "  UUID  "}, // whitespace around identifier fields
			},
		},
		Related: map[string]EntityRelation{
			"  ContactInfo  ": { // whitespace around relation name
				Type:    "ForOne",
				Aliased: "  ContactInfo  ", // whitespace around aliased
				Through: "  parent  ",      // whitespace around through
			},
		},
	}

	// Apply normalization
	NormalizeEntity(&entity)

	// Verify entity name is trimmed
	assert.Equal(t, "TestEntity", entity.Name)

	// Verify field names and types are trimmed
	field, exists := entity.Fields["Email"]
	assert.True(t, exists)
	assert.Equal(t, "ContactInfo.Email", string(field.Type))
	assert.Equal(t, []string{"required", "unique"}, field.Attributes)

	// Verify identifier names and fields are trimmed
	identifier, exists := entity.Identifiers["Primary"]
	assert.True(t, exists)
	assert.Equal(t, []string{"ID", "UUID"}, identifier.Fields)

	// Verify relation names and fields are trimmed
	relation, exists := entity.Related["ContactInfo"]
	assert.True(t, exists)
	assert.Equal(t, "ContactInfo", relation.Aliased)
	assert.Equal(t, "parent", relation.Through)
}

func TestNormalizeModel_TrimsWhitespace(t *testing.T) {
	model := Model{
		Name: "  TestModel  ", // whitespace around name
		Fields: map[string]ModelField{
			"  Email  ": { // whitespace around field name
				Type:       "  String  ",                           // whitespace around field type
				Attributes: []string{"  required  ", "  unique  "}, // whitespace around attributes
			},
		},
		Identifiers: map[string]ModelIdentifier{
			"  Primary  ": { // whitespace around identifier name
				Fields: []string{"  ID  ", "  UUID  "}, // whitespace around identifier fields
			},
		},
		Related: map[string]ModelRelation{
			"  ContactInfo  ": { // whitespace around relation name
				Type:    "ForOne",
				Aliased: "  ContactInfo  ", // whitespace around aliased
				Through: "  parent  ",      // whitespace around through
			},
		},
	}

	// Apply normalization
	NormalizeModel(&model)

	// Verify model name is trimmed
	assert.Equal(t, "TestModel", model.Name)

	// Verify field names and types are trimmed
	field, exists := model.Fields["Email"]
	assert.True(t, exists)
	assert.Equal(t, "String", string(field.Type))
	assert.Equal(t, []string{"required", "unique"}, field.Attributes)

	// Verify identifier names and fields are trimmed
	identifier, exists := model.Identifiers["Primary"]
	assert.True(t, exists)
	assert.Equal(t, []string{"ID", "UUID"}, identifier.Fields)

	// Verify relation names and fields are trimmed
	relation, exists := model.Related["ContactInfo"]
	assert.True(t, exists)
	assert.Equal(t, "ContactInfo", relation.Aliased)
	assert.Equal(t, "parent", relation.Through)
}

func TestNormalizeEntity_TrimsRelationForField(t *testing.T) {
	entity := Entity{
		Name: "TestEntity",
		Related: map[string]EntityRelation{
			"Polymorphic": {
				Type:    "ForOnePoly",
				For:     []string{"  Person  ", "  Company  "}, // whitespace around For items
				Aliased: "Target",
				Through: "parent",
			},
		},
	}

	// Apply normalization
	NormalizeEntity(&entity)

	// Verify For field items are trimmed
	relation := entity.Related["Polymorphic"]
	assert.Equal(t, []string{"Person", "Company"}, relation.For)
}

func TestNormalizeModel_TrimsRelationForField(t *testing.T) {
	model := Model{
		Name: "TestModel",
		Related: map[string]ModelRelation{
			"Polymorphic": {
				Type:    "ForOnePoly",
				For:     []string{"  Person  ", "  Company  "}, // whitespace around For items
				Aliased: "Target",
				Through: "parent",
			},
		},
	}

	// Apply normalization
	NormalizeModel(&model)

	// Verify For field items are trimmed
	relation := model.Related["Polymorphic"]
	assert.Equal(t, []string{"Person", "Company"}, relation.For)
}

func TestNormalizeAllEntities_TrimsWhitespace(t *testing.T) {
	entities := map[string]Entity{
		"TestEntity": {
			Name: "  TestEntity  ",
			Fields: map[string]EntityField{
				"  Email  ": {
					Type: "  ContactInfo.Email  ",
				},
			},
			Related: map[string]EntityRelation{
				"  ContactInfo  ": {
					Type:    "ForOne",
					Aliased: "  ContactInfo  ",
				},
			},
		},
	}

	// Apply normalization
	NormalizeAllEntities(entities)

	// Verify entity name is trimmed
	entity := entities["TestEntity"]
	assert.Equal(t, "TestEntity", entity.Name)

	// Verify field names are trimmed
	field, exists := entity.Fields["Email"]
	assert.True(t, exists)
	assert.Equal(t, "ContactInfo.Email", string(field.Type))

	// Verify relation names are trimmed
	relation, exists := entity.Related["ContactInfo"]
	assert.True(t, exists)
	assert.Equal(t, "ContactInfo", relation.Aliased)
}

func TestNormalizeAllModels_TrimsWhitespace(t *testing.T) {
	models := map[string]Model{
		"TestModel": {
			Name: "  TestModel  ",
			Fields: map[string]ModelField{
				"  Email  ": {
					Type: "  String  ",
				},
			},
			Related: map[string]ModelRelation{
				"  ContactInfo  ": {
					Type:    "ForOne",
					Aliased: "  ContactInfo  ",
				},
			},
		},
	}

	// Apply normalization
	NormalizeAllModels(models)

	// Verify model name is trimmed
	model := models["TestModel"]
	assert.Equal(t, "TestModel", model.Name)

	// Verify field names are trimmed
	field, exists := model.Fields["Email"]
	assert.True(t, exists)
	assert.Equal(t, "String", string(field.Type))

	// Verify relation names are trimmed
	relation, exists := model.Related["ContactInfo"]
	assert.True(t, exists)
	assert.Equal(t, "ContactInfo", relation.Aliased)
}

func TestNormalizeEnum_TrimsWhitespace(t *testing.T) {
	enum := Enum{
		Name: "  TestEnum  ", // whitespace around name
		Type: "  String  ",   // whitespace around type
		Entries: map[string]any{
			"  ACTIVE  ":   "active", // whitespace around entry keys
			"  INACTIVE  ": "inactive",
		},
	}

	// Apply normalization
	NormalizeEnum(&enum)

	// Verify enum name is trimmed
	assert.Equal(t, "TestEnum", enum.Name)

	// Verify enum type is trimmed
	assert.Equal(t, "String", string(enum.Type))

	// Verify entry keys are trimmed
	assert.Equal(t, "active", enum.Entries["ACTIVE"])
	assert.Equal(t, "inactive", enum.Entries["INACTIVE"])

	// Verify old keys with whitespace are removed
	_, exists := enum.Entries["  ACTIVE  "]
	assert.False(t, exists)
}

func TestNormalizeStructure_TrimsWhitespace(t *testing.T) {
	structure := Structure{
		Name: "  TestStructure  ", // whitespace around name
		Fields: map[string]StructureField{
			"  Name  ": { // whitespace around field name
				Type:       "  String  ",                          // whitespace around field type
				Attributes: []string{"  required  ", "  min:1  "}, // whitespace around attributes
			},
		},
	}

	// Apply normalization
	NormalizeStructure(&structure)

	// Verify structure name is trimmed
	assert.Equal(t, "TestStructure", structure.Name)

	// Verify field names and types are trimmed
	field, exists := structure.Fields["Name"]
	assert.True(t, exists)
	assert.Equal(t, "String", string(field.Type))
	assert.Equal(t, []string{"required", "min:1"}, field.Attributes)

	// Verify old field key with whitespace is removed
	_, exists = structure.Fields["  Name  "]
	assert.False(t, exists)
}

func TestNormalizeAllEnums_TrimsWhitespace(t *testing.T) {
	enums := map[string]Enum{
		"TestEnum": {
			Name: "  TestEnum  ",
			Type: "  String  ",
			Entries: map[string]any{
				"  ACTIVE  ": "active",
			},
		},
	}

	// Apply normalization
	NormalizeAllEnums(enums)

	// Verify enum name is trimmed
	enum := enums["TestEnum"]
	assert.Equal(t, "TestEnum", enum.Name)
	assert.Equal(t, "String", string(enum.Type))

	// Verify entry keys are trimmed
	assert.Equal(t, "active", enum.Entries["ACTIVE"])
}

func TestNormalizeAllStructures_TrimsWhitespace(t *testing.T) {
	structures := map[string]Structure{
		"TestStructure": {
			Name: "  TestStructure  ",
			Fields: map[string]StructureField{
				"  Name  ": {
					Type: "  String  ",
				},
			},
		},
	}

	// Apply normalization
	NormalizeAllStructures(structures)

	// Verify structure name is trimmed
	structure := structures["TestStructure"]
	assert.Equal(t, "TestStructure", structure.Name)

	// Verify field names are trimmed
	field, exists := structure.Fields["Name"]
	assert.True(t, exists)
	assert.Equal(t, "String", string(field.Type))
}
