package yaml

import "github.com/kalo-build/clone"

type ModelRelation struct {
	Type    string   `yaml:"type"`
	For     []string `yaml:"for,omitempty"`
	Through string   `yaml:"through,omitempty"`
}

func (r ModelRelation) DeepClone() ModelRelation {
	return ModelRelation{
		Type:    r.Type,
		For:     clone.Slice(r.For),
		Through: r.Through,
	}
}
