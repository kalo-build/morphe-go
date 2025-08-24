package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityValidate_AliasedRelation_Success(t *testing.T) {
	// Setup test models
	contactModel := Model{
		Name: "ContactInfo",
		Fields: map[string]ModelField{
			"email": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"email"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
	}

	// Entity with aliased relationships
	personEntity := Entity{
		Name: "PersonEntity",
		Fields: map[string]EntityField{
			"full_name": {Type: "Person.name"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"full_name"}},
		},
		Related: map[string]EntityRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "ContactInfo",
			},
			"HomeContact": {
				Type:    "ForOne",
				Aliased: "ContactInfo",
			},
		},
	}

	allModels := map[string]Model{
		"ContactInfo": contactModel,
		"Person":      personModel,
	}
	allEnums := map[string]Enum{}

	// Test successful validation with aliased relationships
	allEntities := map[string]Entity{
		"ContactInfo": {Name: "ContactInfo"},
	}
	err := personEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_AliasedRelation_UnknownTarget(t *testing.T) {
	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
	}

	// Entity with unknown aliased target
	personEntity := Entity{
		Name: "PersonEntity",
		Fields: map[string]EntityField{
			"full_name": {Type: "Person.name"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"full_name"}},
		},
		Related: map[string]EntityRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "UnknownTarget",
			},
		},
	}

	allModels := map[string]Model{
		"Person": personModel,
	}
	allEntities := map[string]Entity{
		// No UnknownTarget entity - should cause validation error
	}
	allEnums := map[string]Enum{}

	// Test validation failure with unknown aliased target
	err := personEntity.Validate(allEntities, allModels, allEnums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown aliased target: UnknownTarget")
}

func TestEntityValidate_MultipleAliasedRelations_Success(t *testing.T) {
	// Setup test models
	contactModel := Model{
		Name: "ContactInfo",
		Fields: map[string]ModelField{
			"email": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"email"}},
		},
	}

	projectModel := Model{
		Name: "Project",
		Fields: map[string]ModelField{
			"title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"title"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
	}

	// Entity with multiple different aliased relationships
	personEntity := Entity{
		Name: "PersonEntity",
		Fields: map[string]EntityField{
			"full_name": {Type: "Person.name"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"full_name"}},
		},
		Related: map[string]EntityRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "ContactInfo",
			},
			"HomeContact": {
				Type:    "ForOne",
				Aliased: "ContactInfo",
			},
			"WorkProjects": {
				Type:    "ForMany",
				Aliased: "Project",
			},
			"HobbyProjects": {
				Type:    "ForMany",
				Aliased: "Project",
			},
		},
	}

	allModels := map[string]Model{
		"ContactInfo": contactModel,
		"Project":     projectModel,
		"Person":      personModel,
	}
	allEntities := map[string]Entity{
		"ContactInfo": {Name: "ContactInfo"},
		"Project":     {Name: "Project"},
	}
	allEnums := map[string]Enum{}

	// Test successful validation with multiple aliased relationships
	err := personEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_PolymorphicAliasedRelation_Success(t *testing.T) {
	// Setup test models
	postModel := Model{
		Name: "Post",
		Fields: map[string]ModelField{
			"title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"title"}},
		},
	}

	articleModel := Model{
		Name: "Article",
		Fields: map[string]ModelField{
			"title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"title"}},
		},
	}

	commentModel := Model{
		Name: "Comment",
		Fields: map[string]ModelField{
			"text": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"text"}},
		},
	}

	// Entity with polymorphic aliased relationship
	commentEntity := Entity{
		Name: "CommentEntity",
		Fields: map[string]EntityField{
			"comment_text": {Type: "Comment.text"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"comment_text"}},
		},
		Related: map[string]EntityRelation{
			"WorkCommentable": {
				Type:    "ForOnePoly",
				For:     []string{"Post", "Article"},
				Aliased: "Commentable",
			},
		},
	}

	commentableModel := Model{
		Name: "Commentable",
		Fields: map[string]ModelField{
			"id": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"Post":        postModel,
		"Article":     articleModel,
		"Comment":     commentModel,
		"Commentable": commentableModel,
	}
	allEntities := map[string]Entity{
		"Commentable": {Name: "Commentable"},
	}
	allEnums := map[string]Enum{}

	// Test successful validation with polymorphic aliased relationship
	err := commentEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_MixedAliasedAndNonAliased_Success(t *testing.T) {
	// Setup test models
	contactModel := Model{
		Name: "ContactInfo",
		Fields: map[string]ModelField{
			"email": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"email"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
	}

	companyModel := Model{
		Name: "Company",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
	}

	// Entity with both aliased and non-aliased relationships
	personEntity := Entity{
		Name: "PersonEntity",
		Fields: map[string]EntityField{
			"full_name": {Type: "Person.name"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"full_name"}},
		},
		Related: map[string]EntityRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "ContactInfo", // Aliased relationship
			},
			"Company": {
				Type: "ForOne", // Non-aliased relationship
			},
		},
	}

	allModels := map[string]Model{
		"ContactInfo": contactModel,
		"Person":      personModel,
		"Company":     companyModel,
	}
	allEntities := map[string]Entity{
		"ContactInfo": {Name: "ContactInfo"},
		"Company":     {Name: "Company"},
	}
	allEnums := map[string]Enum{}

	// Test successful validation with mixed aliased and non-aliased relationships
	err := personEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_PolymorphicAliasedRelation_EntityNotInForList(t *testing.T) {
	// Setup test entities for polymorphic relationship
	commentEntity := Entity{
		Name: "Comment",
		Fields: map[string]EntityField{
			"comment_id": {Type: "Comment.id"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"comment_id"}},
		},
		Related: map[string]EntityRelation{
			"commentable": {
				Type: "ForOnePoly",
				For:  []string{"Person", "Article"}, // Entity names in 'for' list
			},
		},
	}

	// Entity with polymorphic inverse relationship
	postEntity := Entity{
		Name: "Post",
		Fields: map[string]EntityField{
			"post_id": {Type: "Post.id"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"post_id"}},
		},
		Related: map[string]EntityRelation{
			"comments": {
				Type:    "HasOnePoly",
				Through: "commentable",
				Aliased: "Comment",
			},
		},
	}

	allModels := map[string]Model{
		"Comment": {
			Name: "Comment",
			Fields: map[string]ModelField{
				"id": {Type: "String"},
			},
			Identifiers: map[string]ModelIdentifier{
				"primary": {Fields: []string{"id"}},
			},
		},
		"Post": {
			Name: "Post",
			Fields: map[string]ModelField{
				"id": {Type: "String"},
			},
			Identifiers: map[string]ModelIdentifier{
				"primary": {Fields: []string{"id"}},
			},
		},
	}
	allEntities := map[string]Entity{
		"Comment": commentEntity,
		"Post":    postEntity,
		"Person":  {Name: "Person"},
		"Article": {Name: "Article"},
	}
	allEnums := map[string]Enum{}

	// This should fail because "Post" is not in the 'for' list ["Person", "Article"]
	err := postEntity.Validate(allEntities, allModels, allEnums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "entity 'Post' is not in the 'for' list")
	assert.Contains(t, err.Error(), "polymorphic relationship 'commentable' in entity 'Comment'")
}

func TestEntityValidate_AliasedRelation_WhitespaceHandling(t *testing.T) {
	// Setup test models
	contactModel := Model{
		Name: "ContactInfo",
		Fields: map[string]ModelField{
			"email": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"email"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
	}

	// Entity with aliased relationships containing whitespace
	personEntity := Entity{
		Name: "PersonEntity",
		Fields: map[string]EntityField{
			"full_name": {Type: "Person.name"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"full_name"}},
		},
		Related: map[string]EntityRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "  ContactInfo  ", // Leading and trailing whitespace
			},
			"HomeContact": {
				Type:    "ForOne",
				Aliased: "\tContactInfo\n", // Tab and newline whitespace
			},
		},
	}

	allModels := map[string]Model{
		"ContactInfo": contactModel,
		"Person":      personModel,
	}
	allEntities := map[string]Entity{
		"ContactInfo": {Name: "ContactInfo"},
	}
	allEnums := map[string]Enum{}

	// Apply normalization before validation (simulating registry loading behavior)
	NormalizeEntity(&personEntity)

	// Test successful validation with whitespace in aliased relationships
	err := personEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_AliasedRelation_WhitespaceOnlyTreatedAsEmpty(t *testing.T) {
	// Setup test models
	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
	}

	// Entity with relationship having whitespace-only aliased value
	personEntity := Entity{
		Name: "PersonEntity",
		Fields: map[string]EntityField{
			"full_name": {Type: "Person.name"},
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"full_name"}},
		},
		Related: map[string]EntityRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "   \t\n   ", // Only whitespace - should be treated as non-aliased
			},
		},
	}

	allModels := map[string]Model{
		"Person": personModel,
	}
	allEntities := map[string]Entity{
		"WorkContact": {Name: "WorkContact"}, // Non-aliased relationships reference by name
	}
	allEnums := map[string]Enum{}

	// Apply normalization before validation (simulating registry loading behavior)
	NormalizeEntity(&personEntity)

	// Test successful validation where whitespace-only aliased is treated as non-aliased
	err := personEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}
