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

// IsRelationAliased checks if a relationship has an alias defined
func IsRelationAliased(aliasedField string) bool {
	return strings.TrimSpace(aliasedField) != ""
}

// GetRelationTargetName returns the actual target name for a relationship
// If aliased is provided, returns the aliased target; otherwise returns the relationship name
func GetRelationTargetName(relationshipName, aliasedField string) string {
	if IsRelationAliased(aliasedField) {
		return strings.TrimSpace(aliasedField)
	}
	return relationshipName
}
