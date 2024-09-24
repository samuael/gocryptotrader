package solsha3

import (
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var (
	errInvalidInput = errors.New("invalid input")
)

func pack(typ string, value interface{}, _isArray bool) ([]byte, error) {
	switch typ {
	case "address":
		if _isArray {
			return common.LeftPadBytes(Address(value), 32), nil
		}

		return Address(value), nil
	case "string":
		return String(value), nil
	case "bool":
		if _isArray {
			return common.LeftPadBytes(Bool(value), 32), nil
		}

		return Bool(value), nil
	}

	regexNumber := regexp.MustCompile(`^(u?int)(\d*)$`)
	matches := regexNumber.FindAllStringSubmatch(typ, -1)
	if len(matches) > 0 {
		match := matches[0]
		var err error
		size := 256
		if len(match) > 2 {
			size, err = strconv.Atoi(match[2])
			if err != nil {
				return nil, err
			}
		}

		_ = size
		if (size%8 != 0) || size == 0 || size > 256 {
			return nil, errors.New("invalid number type " + typ)
		}

		if _isArray {
			size = 256
		}

		var v []byte
		if strings.HasPrefix(typ, "uint256") {
			v = Uint256(value)
		} else {
			return nil, errors.New("type not supported")
		}
		return common.LeftPadBytes(v, size/8), nil
	}

	regexBytes := regexp.MustCompile(`^bytes(\d+)$`)
	matches = regexBytes.FindAllStringSubmatch(typ, -1)
	if len(matches) > 0 {
		match := matches[0]

		byteLen, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}

		strSize := strconv.Itoa(byteLen)
		if strSize != match[1] || byteLen == 0 || byteLen > 32 {
			panic("invalid number type " + typ)
		}

		if _isArray {
			s := reflect.ValueOf(value)
			v := s.Index(0).Bytes()
			z := make([]byte, 64)
			copy(z, v)
			return z, nil
		}

		str, isString := value.(string)
		if isString && strings.HasPrefix(str, "0x") {
			s := strings.TrimPrefix(str, "0x")
			if len(s)%2 == 1 {
				s = "0" + s
			}
			hexb, err := hex.DecodeString(s)
			if err != nil {
				panic(err)
			}
			z := make([]byte, byteLen)
			copy(z, hexb)
			return z, nil
		} else if isString {
			s := reflect.ValueOf(value)
			z := make([]byte, byteLen)
			copy(z, s.Bytes())
			return z, nil
		}

		s := reflect.ValueOf(value)
		z := make([]byte, byteLen)
		b := make([]byte, s.Len())
		for i := range s.Len() {
			ifc := s.Index(i).Interface()
			v, ok := ifc.(byte)
			if ok {
				b[i] = v
			} else {
				v, ok := ifc.(string)
				if ok {
					v = strings.TrimPrefix(v, "0x")
					if len(v)%2 == 1 {
						v = "0" + v
					}
					decoded, err := hex.DecodeString(v)
					if err != nil {
						panic(err)
					}
					b[i] = decoded[0]
				}
			}
		}
		copy(z, b)
		return z, nil
	}
	regexArray := regexp.MustCompile(`^(.*)\[(\d*)\]$`)
	matches = regexArray.FindAllStringSubmatch(typ, -1)
	if len(matches) > 0 {
		match := matches[0]

		_ = match
		if reflect.TypeOf(value).Kind() == reflect.Array || reflect.TypeOf(value).Kind() == reflect.Slice {
			baseType := match[1]
			k := reflect.ValueOf(value)
			count := k.Len()
			var err error
			if len(match) > 1 && match[2] != "" {
				count, err = strconv.Atoi(match[2])
				if err != nil {
					return nil, err
				}
			}
			if count != k.Len() {
				return nil, errors.New("invalid value for " + typ)
			}

			var result [][]byte
			for i := range k.Len() {
				val := k.Index(i).Interface()
				data, err := pack(baseType, val, true)
				if err != nil {
					return nil, err
				}
				result = append(result, data)
			}
			var array []byte
			for _, b := range result {
				array = append(array, b...)
			}
			return array, nil
		}
	}
	return nil, errInvalidInput
}

// Address address
func Address(input interface{}) []byte {
	switch v := input.(type) {
	case common.Address:
		return v.Bytes()
	case string:
		v = strings.TrimPrefix(v, "0x")
		if v == "" || v == "0" {
			return []byte{0}
		}

		v = func(val string) string {
			if len(v)%2 == 1 {
				val = "0" + val
			}
			return val
		}(v)
		decoded, err := hex.DecodeString(v)
		if err != nil {
			panic(err)
		}

		return decoded
	case []byte:
		return v
	}

	if reflect.TypeOf(input).Kind() == reflect.Array ||
		reflect.TypeOf(input).Kind() == reflect.Slice {
		return AddressArray(input)
	}

	return common.HexToAddress("").Bytes()
}

// AddressArray address
func AddressArray(input interface{}) []byte {
	var values []byte
	s := reflect.ValueOf(input)
	for i := range s.Len() {
		val := s.Index(i).Interface()
		result := common.LeftPadBytes(Address(val), 32)
		values = append(values, result...)
	}
	return values
}

// String string
func String(input interface{}) []byte {
	switch v := input.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	}

	if reflect.TypeOf(input).Kind() == reflect.Array ||
		reflect.TypeOf(input).Kind() == reflect.Slice {
		return StringArray(input)
	}

	return []byte("")
}

// StringArray string
func StringArray(input interface{}) []byte {
	var values []byte
	s := reflect.ValueOf(input)
	for i := range s.Len() {
		val := s.Index(i).Interface()
		result := String(val)
		values = append(values, result...)
	}
	return values
}

// Uint256 uint256
func Uint256(input interface{}) []byte {
	switch v := input.(type) {
	case *big.Int:
		return math.U256Bytes(v)
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		return math.U256Bytes(bn)
	}

	if reflect.TypeOf(input).Kind() == reflect.Array ||
		reflect.TypeOf(input).Kind() == reflect.Slice {
		return Uint256Array(input)
	}

	return common.RightPadBytes([]byte(""), 32)
}

// Uint256Array uint256 array
func Uint256Array(input interface{}) []byte {
	var values []byte
	s := reflect.ValueOf(input)
	for i := range s.Len() {
		val := s.Index(i).Interface()
		result := common.LeftPadBytes(Uint256(val), 32)
		values = append(values, result...)
	}
	return values
}

// Bool bool
func Bool(input interface{}) []byte {
	if v, ok := input.(bool); ok {
		if v {
			return []byte{0x1}
		}
		return []byte{0x0}
	}

	if reflect.TypeOf(input).Kind() == reflect.Array ||
		reflect.TypeOf(input).Kind() == reflect.Slice {
		return BoolArray(input)
	}

	return []byte{0x0}
}

// BoolArray bool array
func BoolArray(input interface{}) []byte {
	var values []byte
	s := reflect.ValueOf(input)
	for i := range s.Len() {
		val := s.Index(i).Interface()
		result := common.LeftPadBytes(Bool(val), 32)
		values = append(values, result...)
	}
	return values
}
