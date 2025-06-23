package yamlops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRelationFor(t *testing.T) {
	assert.True(t, IsRelationFor("ForOne"))
	assert.True(t, IsRelationFor("ForMany"))
	assert.True(t, IsRelationFor("forone"))
	assert.False(t, IsRelationFor("HasOne"))
	assert.False(t, IsRelationFor("HasMany"))
	assert.False(t, IsRelationFor("Invalid"))
	assert.False(t, IsRelationFor(""))
}

func TestIsRelationHas(t *testing.T) {
	assert.True(t, IsRelationHas("HasOne"))
	assert.True(t, IsRelationHas("HasMany"))
	assert.True(t, IsRelationHas("hasone"))
	assert.False(t, IsRelationHas("ForOne"))
	assert.False(t, IsRelationHas("ForMany"))
	assert.False(t, IsRelationHas("Invalid"))
	assert.False(t, IsRelationHas(""))
}

func TestIsRelationMany(t *testing.T) {
	assert.True(t, IsRelationMany("ForMany"))
	assert.True(t, IsRelationMany("HasMany"))
	assert.True(t, IsRelationMany("formany"))
	assert.False(t, IsRelationMany("ForOne"))
	assert.False(t, IsRelationMany("HasOne"))
	assert.False(t, IsRelationMany("Invalid"))
	assert.False(t, IsRelationMany(""))
}

func TestIsRelationOne(t *testing.T) {
	assert.True(t, IsRelationOne("ForOne"))
	assert.True(t, IsRelationOne("HasOne"))
	assert.True(t, IsRelationOne("forone"))
	assert.False(t, IsRelationOne("ForMany"))
	assert.False(t, IsRelationOne("HasMany"))
	assert.False(t, IsRelationOne("Invalid"))
	assert.False(t, IsRelationOne(""))
}

func TestIsRelationPoly(t *testing.T) {
	assert.True(t, IsRelationPoly("ForOnePoly"))
	assert.True(t, IsRelationPoly("ForManyPoly"))
	assert.True(t, IsRelationPoly("HasOnePoly"))
	assert.True(t, IsRelationPoly("HasManyPoly"))
	assert.True(t, IsRelationPoly("foronepoly"))
	assert.False(t, IsRelationPoly("ForOne"))
	assert.False(t, IsRelationPoly("ForMany"))
	assert.False(t, IsRelationPoly("HasOne"))
	assert.False(t, IsRelationPoly("HasMany"))
	assert.False(t, IsRelationPoly("SomethingPoly"))
	assert.False(t, IsRelationPoly("Poly"))
	assert.False(t, IsRelationPoly(""))
}

func TestIsRelationPolyFor(t *testing.T) {
	assert.True(t, IsRelationPolyFor("ForOnePoly"))
	assert.True(t, IsRelationPolyFor("ForManyPoly"))
	assert.True(t, IsRelationPolyFor("foronepoly"))
	assert.False(t, IsRelationPolyFor("HasOnePoly"))
	assert.False(t, IsRelationPolyFor("HasManyPoly"))
	assert.False(t, IsRelationPolyFor("ForOne"))
	assert.False(t, IsRelationPolyFor("ForMany"))
	assert.False(t, IsRelationPolyFor("Invalid"))
	assert.False(t, IsRelationPolyFor(""))
}

func TestIsRelationPolyHas(t *testing.T) {
	assert.True(t, IsRelationPolyHas("HasOnePoly"))
	assert.True(t, IsRelationPolyHas("HasManyPoly"))
	assert.True(t, IsRelationPolyHas("hasonepoly"))
	assert.False(t, IsRelationPolyHas("ForOnePoly"))
	assert.False(t, IsRelationPolyHas("ForManyPoly"))
	assert.False(t, IsRelationPolyHas("HasOne"))
	assert.False(t, IsRelationPolyHas("HasMany"))
	assert.False(t, IsRelationPolyHas("Invalid"))
	assert.False(t, IsRelationPolyHas(""))
}

func TestIsRelationPolyOne(t *testing.T) {
	assert.True(t, IsRelationPolyOne("ForOnePoly"))
	assert.True(t, IsRelationPolyOne("HasOnePoly"))
	assert.True(t, IsRelationPolyOne("foronepoly"))
	assert.False(t, IsRelationPolyOne("ForManyPoly"))
	assert.False(t, IsRelationPolyOne("HasManyPoly"))
	assert.False(t, IsRelationPolyOne("ForOne"))
	assert.False(t, IsRelationPolyOne("HasOne"))
	assert.False(t, IsRelationPolyOne("Invalid"))
	assert.False(t, IsRelationPolyOne(""))
}

func TestIsRelationPolyMany(t *testing.T) {
	assert.True(t, IsRelationPolyMany("ForManyPoly"))
	assert.True(t, IsRelationPolyMany("HasManyPoly"))
	assert.True(t, IsRelationPolyMany("formanypoly"))
	assert.False(t, IsRelationPolyMany("ForOnePoly"))
	assert.False(t, IsRelationPolyMany("HasOnePoly"))
	assert.False(t, IsRelationPolyMany("ForMany"))
	assert.False(t, IsRelationPolyMany("HasMany"))
	assert.False(t, IsRelationPolyMany("Invalid"))
	assert.False(t, IsRelationPolyMany(""))
}
