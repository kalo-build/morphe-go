package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelRelationDeepClone(t *testing.T) {
	original := ModelRelation{
		Type:    "ForOnePoly",
		For:     []string{"Post", "Article", "Video"},
		Through: "Commentable",
	}

	cloned := original.DeepClone()

	assert.Equal(t, original.Type, cloned.Type)
	assert.Equal(t, original.For, cloned.For)
	assert.Equal(t, original.Through, cloned.Through)

	cloned.For[0] = "Modified"
	assert.NotEqual(t, original.For[0], cloned.For[0])
	assert.Equal(t, "Post", original.For[0])
	assert.Equal(t, "Modified", cloned.For[0])
}

func TestModelRelationDeepCloneWithAliased(t *testing.T) {
	original := ModelRelation{
		Type:    "ForOne",
		Aliased: "ContactInfo",
	}

	cloned := original.DeepClone()

	assert.Equal(t, original.Type, cloned.Type)
	assert.Equal(t, original.Aliased, cloned.Aliased)

	cloned.Aliased = "Modified"
	assert.NotEqual(t, original.Aliased, cloned.Aliased)
	assert.Equal(t, "ContactInfo", original.Aliased)
	assert.Equal(t, "Modified", cloned.Aliased)
}

func TestModelRelationPolymorphicFields_ForOnePoly(t *testing.T) {
	forOnePolyRelation := ModelRelation{
		Type: "ForOnePoly",
		For:  []string{"Post", "Article"},
	}
	assert.Equal(t, "ForOnePoly", forOnePolyRelation.Type)
	assert.Equal(t, []string{"Post", "Article"}, forOnePolyRelation.For)
	assert.Empty(t, forOnePolyRelation.Through)
	assert.Empty(t, forOnePolyRelation.Aliased)
}

func TestModelRelationPolymorphicFields_HasOnePoly(t *testing.T) {
	hasOnePolyRelation := ModelRelation{
		Type:    "HasOnePoly",
		Through: "Commentable",
	}
	assert.Equal(t, "HasOnePoly", hasOnePolyRelation.Type)
	assert.Empty(t, hasOnePolyRelation.For)
	assert.Equal(t, "Commentable", hasOnePolyRelation.Through)
	assert.Empty(t, hasOnePolyRelation.Aliased)
}

func TestModelRelationAliased_ForOne(t *testing.T) {
	aliasedRelation := ModelRelation{
		Type:    "ForOne",
		Aliased: "ContactInfo",
	}
	assert.Equal(t, "ForOne", aliasedRelation.Type)
	assert.Equal(t, "ContactInfo", aliasedRelation.Aliased)
	assert.Empty(t, aliasedRelation.For)
	assert.Empty(t, aliasedRelation.Through)
}

func TestModelRelationAliased_PolymorphicWithAlias(t *testing.T) {
	aliasedPolyRelation := ModelRelation{
		Type:    "ForOnePoly",
		For:     []string{"Post", "Article"},
		Aliased: "Commentable",
	}
	assert.Equal(t, "ForOnePoly", aliasedPolyRelation.Type)
	assert.Equal(t, []string{"Post", "Article"}, aliasedPolyRelation.For)
	assert.Equal(t, "Commentable", aliasedPolyRelation.Aliased)
	assert.Empty(t, aliasedPolyRelation.Through)
}
