package types

import (
	"encoding/hex"
	"encoding/json"
)

type HexBytes []byte

func (t HexBytes) Size() int {
	return len([]byte(t))
}

func (t HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *HexBytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
}
