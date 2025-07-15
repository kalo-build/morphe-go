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
	err := personEntity.Validate(allModels, allEnums)
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
	allEnums := map[string]Enum{}

	// Test validation failure with unknown aliased target
	err := personEntity.Validate(allModels, allEnums)
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
	allEnums := map[string]Enum{}

	// Test successful validation with multiple aliased relationships
	err := personEntity.Validate(allModels, allEnums)
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
	allEnums := map[string]Enum{}

	// Test successful validation with polymorphic aliased relationship
	err := commentEntity.Validate(allModels, allEnums)
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
	allEnums := map[string]Enum{}

	// Test successful validation with mixed aliased and non-aliased relationships
	err := personEntity.Validate(allModels, allEnums)
	assert.NoError(t, err)
}
