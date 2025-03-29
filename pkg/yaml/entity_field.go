package yaml

import (
	"github.com/kalo-build/clone"
)

type EntityField struct {
	Type       ModelFieldPath `yaml:"type"`
	Attributes []string       `yaml:"attributes"`
}

func (f EntityField) DeepClone() EntityField {
	return EntityField{
		Type:       f.Type,
		Attributes: clone.Slice(f.Attributes),
	}
}
