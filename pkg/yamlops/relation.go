package yamlops

import "strings"

func IsRelationFor(relationType string) bool {
	return strings.HasPrefix(strings.ToLower(relationType), "for")
}

func IsRelationHas(relationType string) bool {
	return strings.HasPrefix(strings.ToLower(relationType), "has")
}

func IsRelationMany(relationType string) bool {
	return strings.Contains(strings.ToLower(relationType), "many")
}

func IsRelationOne(relationType string) bool {
	return strings.Contains(strings.ToLower(relationType), "one")
}

func IsRelationPoly(relationType string) bool {
	lowerType := strings.ToLower(relationType)
	return (IsRelationFor(relationType) || IsRelationHas(relationType)) &&
		strings.HasSuffix(lowerType, "poly")
}

func IsRelationPolyFor(relationType string) bool {
	return IsRelationPoly(relationType) && IsRelationFor(relationType)
}

func IsRelationPolyHas(relationType string) bool {
	return IsRelationPoly(relationType) && IsRelationHas(relationType)
}

func IsRelationPolyOne(relationType string) bool {
	return IsRelationPoly(relationType) && IsRelationOne(relationType)
}

func IsRelationPolyMany(relationType string) bool {
	return IsRelationPoly(relationType) && IsRelationMany(relationType)
}
