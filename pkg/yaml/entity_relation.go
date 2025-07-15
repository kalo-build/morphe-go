package yaml

import "github.com/kalo-build/clone"

type EntityRelation struct {
	Type    string   `yaml:"type"`
	For     []string `yaml:"for,omitempty"`
	Through string   `yaml:"through,omitempty"`
	Aliased string   `yaml:"aliased,omitempty"`
}

func (f EntityRelation) DeepClone() EntityRelation {
	return EntityRelation{
		Type:    f.Type,
		For:     clone.Slice(f.For),
		Through: f.Through,
		Aliased: f.Aliased,
	}
}
