package types

type Extension struct {
	Type uint16   `serialize:"true" json:"type"`
	Data HexBytes `serialize:"true" json:"data"`
}
