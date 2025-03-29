package clone

type DeepCloneable[TType any] interface {
	DeepClone() TType
}
