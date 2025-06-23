package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityRelationDeepClone(t *testing.T) {
	original := EntityRelation{
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

func TestEntityRelationPolymorphicFields_ForOnePoly(t *testing.T) {
	forOnePolyRelation := EntityRelation{
		Type: "ForOnePoly",
		For:  []string{"Post", "Article"},
	}
	assert.Equal(t, "ForOnePoly", forOnePolyRelation.Type)
	assert.Equal(t, []string{"Post", "Article"}, forOnePolyRelation.For)
	assert.Empty(t, forOnePolyRelation.Through)
}

func TestEntityRelationPolymorphicFields_HasOnePoly(t *testing.T) {
	hasOnePolyRelation := EntityRelation{
		Type:    "HasOnePoly",
		Through: "Commentable",
	}
	assert.Equal(t, "HasOnePoly", hasOnePolyRelation.Type)
	assert.Empty(t, hasOnePolyRelation.For)
	assert.Equal(t, "Commentable", hasOnePolyRelation.Through)
}
