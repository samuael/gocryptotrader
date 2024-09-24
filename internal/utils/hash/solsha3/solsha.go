package solsha3

import (
	"errors"

	"golang.org/x/crypto/sha3"
)

// solsha3 solidity sha3
func solsha3(types []string, values ...interface{}) []byte {
	b := [][]byte{}
	for i, typ := range types {
		data, err := pack(typ, values[i], false)
		if err != nil {
			return nil
		}
		b = append(b, data)
	}

	hash := sha3.NewLegacyKeccak256()

	var bs []byte
	for _, bi := range b {
		bs = append(bs, bi...)
	}
	hash.Write(bs)
	return hash.Sum(nil)
}

// SoliditySHA3 solidity sha3
func SoliditySHA3(data ...interface{}) ([]byte, error) {
	types, ok := data[0].([]string)
	if !ok {
		return nil, errors.New("invalid data types")
	}
	rest := data[1:]
	if len(rest) == len(types) {
		return solsha3(types, data[1:]...), nil
	}
	iface, ok := data[1].([]interface{})
	if ok {
		return solsha3(types, iface...), nil
	}
	return nil, errors.New("invalid input")
}
