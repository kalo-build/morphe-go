package yaml

import "github.com/kaloseia/clone"

type ModelField struct {
	Type       ModelFieldType `yaml:"type"`
	Attributes []string       `yaml:"attributes"`
}

func (f ModelField) DeepClone() ModelField {
	return ModelField{
		Type:       f.Type,
		Attributes: clone.Slice(f.Attributes),
	}
}
