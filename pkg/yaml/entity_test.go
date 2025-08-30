package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestEntityValidate_AliasedFieldPath_Success(t *testing.T) {
	// Setup test models
	contactModel := Model{
		Name: "Contact",
		Fields: map[string]ModelField{
			"ID":    {Type: "AutoIncrement"},
			"Email": {Type: "String"},
			"Phone": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
			"PersonalContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
		},
	}

	// Entity that references fields through aliased relationships
	personProfileEntity := Entity{
		Name: "PersonProfile",
		Fields: map[string]EntityField{
			"id":            {Type: "Person.ID"},
			"name":          {Type: "Person.Name"},
			"workEmail":     {Type: "Person.WorkContact.Email"},     // Should resolve to Contact.Email
			"personalPhone": {Type: "Person.PersonalContact.Phone"}, // Should resolve to Contact.Phone
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"Contact": contactModel,
		"Person":  personModel,
	}
	allEntities := map[string]Entity{
		"PersonProfile": personProfileEntity,
	}
	allEnums := map[string]Enum{}

	// Test successful validation with aliased field paths
	err := personProfileEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_AliasedFieldPath_NonExistentAlias(t *testing.T) {
	// Setup test models
	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"BadAlias": {
				Type:    "ForOne",
				Aliased: "NonExistentModel", // This model doesn't exist
			},
		},
	}

	// Entity that tries to reference through non-existent aliased model
	testEntity := Entity{
		Name: "TestEntity",
		Fields: map[string]EntityField{
			"id":       {Type: "Person.ID"},
			"badField": {Type: "Person.BadAlias.SomeField"}, // Should fail
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"Person": personModel,
	}
	allEntities := map[string]Entity{
		"TestEntity": testEntity,
	}
	allEnums := map[string]Enum{}

	// Test that validation fails with appropriate error
	err := testEntity.Validate(allEntities, allModels, allEnums)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "aliased target model NonExistentModel")
	assert.Contains(t, err.Error(), "via relationship BadAlias")
}

func TestEntityValidate_AliasedFieldPath_NonExistentField(t *testing.T) {
	// Setup test models
	contactModel := Model{
		Name: "Contact",
		Fields: map[string]ModelField{
			"ID":    {Type: "AutoIncrement"},
			"Email": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"ID": {Type: "AutoIncrement"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
		},
	}

	// Entity that references non-existent field in aliased target
	testEntity := Entity{
		Name: "TestEntity",
		Fields: map[string]EntityField{
			"id":       {Type: "Person.ID"},
			"badField": {Type: "Person.WorkContact.NonExistentField"}, // Field doesn't exist in Contact
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"Contact": contactModel,
		"Person":  personModel,
	}
	allEntities := map[string]Entity{
		"TestEntity": testEntity,
	}
	allEnums := map[string]Enum{}

	// Test that validation fails with appropriate error
	err := testEntity.Validate(allEntities, allModels, allEnums)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "NonExistentField")
}

func TestEntityValidate_AliasedFieldPath_ComplexTraversal(t *testing.T) {
	// Setup test models for multi-hop traversal
	contactModel := Model{
		Name: "Contact",
		Fields: map[string]ModelField{
			"ID":    {Type: "AutoIncrement"},
			"Email": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"WorkContact": {
				Type:    "ForOne",
				Aliased: "Contact",
			},
		},
	}

	companyModel := Model{
		Name: "Company",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"CEO": {
				Type:    "ForOne",
				Aliased: "Person",
			},
		},
	}

	// Entity that uses multi-hop traversal through aliases
	companyProfileEntity := Entity{
		Name: "CompanyProfile",
		Fields: map[string]EntityField{
			"id":           {Type: "Company.ID"},
			"companyName":  {Type: "Company.Name"},
			"ceoName":      {Type: "Company.CEO.Name"},              // Company -> Person
			"ceoWorkEmail": {Type: "Company.CEO.WorkContact.Email"}, // Company -> Person -> Contact
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"Contact": contactModel,
		"Person":  personModel,
		"Company": companyModel,
	}
	allEntities := map[string]Entity{
		"CompanyProfile": companyProfileEntity,
	}
	allEnums := map[string]Enum{}

	// Test successful validation with complex multi-hop aliased paths
	err := companyProfileEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_AliasedFieldPath_BackwardCompatibility(t *testing.T) {
	// Test that non-aliased paths continue to work correctly
	contactModel := Model{
		Name: "ContactInfo",
		Fields: map[string]ModelField{
			"ID":    {Type: "AutoIncrement"},
			"Email": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	personModel := Model{
		Name: "Person",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"ContactInfo": {
				Type: "ForOne",
				// No aliasing - direct relationship name matches model name
			},
		},
	}

	// Entity using non-aliased paths
	personEntity := Entity{
		Name: "PersonEntity",
		Fields: map[string]EntityField{
			"id":    {Type: "Person.ID"},
			"name":  {Type: "Person.Name"},
			"email": {Type: "Person.ContactInfo.Email"}, // Non-aliased path
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"ContactInfo": contactModel,
		"Person":      personModel,
	}
	allEntities := map[string]Entity{
		"PersonEntity": personEntity,
	}
	allEnums := map[string]Enum{}

	// Test that non-aliased paths still work
	err := personEntity.Validate(allEntities, allModels, allEnums)
	assert.NoError(t, err)
}

func TestEntityValidate_AliasedFieldPath_PolymorphicTraversal(t *testing.T) {
	// Setup test models with polymorphic relationships
	commentModel := Model{
		Name: "Comment",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Text": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
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
			"ID":    {Type: "AutoIncrement"},
			"Title": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"Comments": {
				Type:    "HasManyPoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	// Entity that tries to traverse through polymorphic relationship
	testEntity := Entity{
		Name: "TestEntity",
		Fields: map[string]EntityField{
			"id":       {Type: "Post.ID"},
			"badField": {Type: "Post.Comments.Text"}, // Should fail - cannot traverse through polymorphic
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"Comment": commentModel,
		"Post":    postModel,
	}
	allEntities := map[string]Entity{
		"TestEntity": testEntity,
		"Comment":    {Name: "Comment"},
		"Post":       {Name: "Post"},
	}
	allEnums := map[string]Enum{}

	// Test that validation fails with appropriate error
	err := testEntity.Validate(allEntities, allModels, allEnums)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot traverse through polymorphic relationship")
	assert.Contains(t, err.Error(), "Comments")
	assert.Contains(t, err.Error(), "Post.Comments")
}

func TestEntityValidate_AliasedFieldPath_TraversalThroughForPoly(t *testing.T) {
	// Test that we can't traverse through ForOnePoly/ForManyPoly relationships in field paths
	authorModel := Model{
		Name: "Author",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	commentModel := Model{
		Name: "Comment",
		Fields: map[string]ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Text": {Type: "String"},
		},
		Identifiers: map[string]ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
			"Author": {
				Type:    "ForOne",
				Aliased: "Author",
			},
		},
	}

	// Entity that tries to traverse through a ForOnePoly relationship
	testEntity := Entity{
		Name: "TestEntity",
		Fields: map[string]EntityField{
			"id":            {Type: "Comment.ID"},
			"authorName":    {Type: "Comment.Author.Name"},    // This should work - regular ForOne
			"commentableId": {Type: "Comment.Commentable.ID"}, // Should fail - cannot traverse through ForOnePoly
		},
		Identifiers: map[string]EntityIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	allModels := map[string]Model{
		"Comment": commentModel,
		"Author":  authorModel,
	}
	allEntities := map[string]Entity{
		"TestEntity": testEntity,
		"Comment":    {Name: "Comment"},
		"Author":     {Name: "Author"},
	}
	allEnums := map[string]Enum{}

	// Test that validation fails for the polymorphic traversal
	err := testEntity.Validate(allEntities, allModels, allEnums)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot traverse through polymorphic relationship")
	assert.Contains(t, err.Error(), "Commentable")
}
