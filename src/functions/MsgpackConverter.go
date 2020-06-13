package functions

import (
	"encoding/base64"

	"github.com/vmihailenco/msgpack"
)

// ToMSGPACK Converte para msgpack
func ToMSGPACK(data interface{}) []byte {
	b, err := msgpack.Marshal(&data)
	if err != nil {
		return nil
	} else {
		return b
	}
}

// FromMSGPACK Convert de msgpack para item
func FromMSGPACK(dataString string, item interface{}) error {
	data, err := base64.StdEncoding.DecodeString(dataString)
	if err != nil {
		return err
	}
	err = msgpack.Unmarshal(data, &item)
	if err != nil {
		return err
	}
	return nil
}
