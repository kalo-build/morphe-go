package clone

func DeepCloneSlice[TCloneable DeepCloneable[TCloneable]](original []TCloneable) []TCloneable {
	if original == nil {
		return nil
	}
	copiedSlice := make([]TCloneable, len(original))
	for copyIdx, copyable := range original {
		copiedSlice[copyIdx] = copyable.DeepClone()
	}
	return copiedSlice
}

func DeepCloneSlicePointers[TCloneable DeepCloneable[TCloneable]](original []*TCloneable) []*TCloneable {
	if original == nil {
		return nil
	}
	copiedSlice := make([]*TCloneable, len(original))
	for copyIdx, copyablePtr := range original {
		if copyablePtr == nil {
			copiedSlice[copyIdx] = nil
			continue
		}
		copyable := *copyablePtr
		copyableClone := copyable.DeepClone()
		copiedSlice[copyIdx] = &copyableClone
	}
	return copiedSlice
}
