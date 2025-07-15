package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelValidateWithModels_Success(t *testing.T) {
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
		Related: map[string]ModelRelation{
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
	err := personModel.ValidateWithModels(allModels, allEnums)
	assert.NoError(t, err)
}

func TestModelValidateWithModels_UnknownAliasedTarget(t *testing.T) {
	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
		},
		Related: map[string]ModelRelation{
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
	err := personModel.ValidateWithModels(allModels, allEnums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown aliased target: UnknownTarget")
}

func TestModelValidateWithModels_MultipleAliasedRelations(t *testing.T) {
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
		Related: map[string]ModelRelation{
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
	err := personModel.ValidateWithModels(allModels, allEnums)
	assert.NoError(t, err)
}

func TestModelValidateWithModels_MixedAliasedAndNonAliased(t *testing.T) {
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

	companyModel := Model{
		Name: "Company",
		Fields: map[string]ModelField{
			"name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"name"}},
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
		Related: map[string]ModelRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "ContactInfo",
			},
			"Company": {
				Type: "ForOne",
				// No aliased field - should use relationship name
			},
		},
	}

	allModels := map[string]Model{
		"ContactInfo": contactModel,
		"Company":     companyModel,
		"Person":      personModel,
	}
	allEnums := map[string]Enum{}

	// Test successful validation with mixed aliased/non-aliased relationships
	err := personModel.ValidateWithModels(allModels, allEnums)
	assert.NoError(t, err)
}

func TestModelValidateWithModels_PolymorphicInverseAliasing_Success(t *testing.T) {
	// Setup test models for polymorphic inverse aliasing
	commentModel := Model{
		Name: "Comment",
		Fields: map[string]ModelField{
			"text": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"text"}},
		},
		Related: map[string]ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	postModel := Model{
		Name: "Post",
		Fields: map[string]ModelField{
			"title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"title"}},
		},
		Related: map[string]ModelRelation{
			"Note": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
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
		Related: map[string]ModelRelation{
			"Review": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	allModels := map[string]Model{
		"Comment": commentModel,
		"Post":    postModel,
		"Article": articleModel,
	}
	allEnums := map[string]Enum{}

	// Test successful validation with polymorphic inverse aliasing
	err := postModel.ValidateWithModels(allModels, allEnums)
	assert.NoError(t, err)

	err = articleModel.ValidateWithModels(allModels, allEnums)
	assert.NoError(t, err)
}

func TestModelValidateWithModels_PolymorphicInverseAliasing_MissingThroughRelation(t *testing.T) {
	// Comment model without the expected "Commentable" relationship
	commentModel := Model{
		Name: "Comment",
		Fields: map[string]ModelField{
			"text": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"text"}},
		},
		Related: map[string]ModelRelation{
			"SomeOtherRelation": {
				Type: "ForOne",
			},
		},
	}

	postModel := Model{
		Name: "Post",
		Fields: map[string]ModelField{
			"title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"title"}},
		},
		Related: map[string]ModelRelation{
			"Note": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	allModels := map[string]Model{
		"Comment": commentModel,
		"Post":    postModel,
	}
	allEnums := map[string]Enum{}

	// Test validation failure when through relationship is missing
	err := postModel.ValidateWithModels(allModels, allEnums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "aliased model 'Comment' does not have relationship 'Commentable'")
}

func TestModelValidateWithModels_PolymorphicInverseAliasing_InvalidThroughRelationType(t *testing.T) {
	// Comment model with "Commentable" but not polymorphic
	commentModel := Model{
		Name: "Comment",
		Fields: map[string]ModelField{
			"text": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"text"}},
		},
		Related: map[string]ModelRelation{
			"Commentable": {
				Type: "ForOne", // Not polymorphic!
			},
		},
	}

	postModel := Model{
		Name: "Post",
		Fields: map[string]ModelField{
			"title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"title"}},
		},
		Related: map[string]ModelRelation{
			"Note": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	allModels := map[string]Model{
		"Comment": commentModel,
		"Post":    postModel,
	}
	allEnums := map[string]Enum{}

	// Test validation failure when through relationship is not polymorphic
	err := postModel.ValidateWithModels(allModels, allEnums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not a polymorphic 'For' relationship")
}

func TestModelValidateWithModels_PolymorphicInverseAliasing_ModelNotInForList(t *testing.T) {
	// Comment model with polymorphic relationship but Post not in 'for' list
	commentModel := Model{
		Name: "Comment",
		Fields: map[string]ModelField{
			"text": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"text"}},
		},
		Related: map[string]ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Article", "Task"}, // Post is missing!
			},
		},
	}

	postModel := Model{
		Name: "Post",
		Fields: map[string]ModelField{
			"title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"title"}},
		},
		Related: map[string]ModelRelation{
			"Note": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	allModels := map[string]Model{
		"Comment": commentModel,
		"Post":    postModel,
	}
	allEnums := map[string]Enum{}

	// Test validation failure when current model is not in 'for' list
	err := postModel.ValidateWithModels(allModels, allEnums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "current model 'Post' is not in the 'for' list")
}
