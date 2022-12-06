package state

type VM interface {
	Accepted(*Block) error
	Verified(*Block) error
}
