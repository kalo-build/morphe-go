package yaml

import "gopkg.in/yaml.v3"

type ModelIdentifier struct {
	Fields []string
}

func (id *ModelIdentifier) UnmarshalYAML(value *yaml.Node) error {
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
