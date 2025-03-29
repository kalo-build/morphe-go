package yaml

import (
	"github.com/kalo-build/clone"
	"gopkg.in/yaml.v3"
)

type EntityIdentifier struct {
	Fields []string `yaml:"fields"`
}

func (id EntityIdentifier) DeepClone() EntityIdentifier {
	return EntityIdentifier{
		Fields: clone.Slice(id.Fields),
	}
}

func (id *EntityIdentifier) UnmarshalYAML(value *yaml.Node) error {
	var fieldName string
	unmarshalErr := value.Decode(&fieldName)
	if unmarshalErr == nil {
		id.Fields = []string{fieldName}
		return nil
	}
	var fieldNames []string
	unmarshalErr = value.Decode(&fieldNames)
	if unmarshalErr == nil {
		id.Fields = fieldNames
		return nil
	}

	return unmarshalErr
}
