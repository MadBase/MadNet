package interfaces

type Filter interface {
	Insert(input []byte) error
	Delete(needle []byte)
	Lookup(needle []byte) bool
}
