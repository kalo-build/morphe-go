package clone

func Slice[TValue any](original []TValue) []TValue {
	if original == nil {
		return nil
	}
	fieldsCopy := make([]TValue, len(original))
	copy(fieldsCopy, original)
	return fieldsCopy
}
