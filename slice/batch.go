package slice

import (
	"errors"
	"reflect"
)

func GetChunks(slice interface{}, size uint) ([][]interface{}, error) {
	interfaceSlice, err := GetInterfaceSlice(slice)
	if err != nil {
		return nil, err
	}

	return GetInterfaceSliceBatch(interfaceSlice, size), nil
}

func GetInterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func GetInterfaceSliceBatch(values []interface{}, sizeUint uint) (chunks [][]interface{}) {
	size := int(sizeUint)
	resultLength := (len(values) + size - 1) / size
	result := make([][]interface{}, resultLength)
	lo, hi := 0, size
	for i := range result {
		if hi > len(values) {
			hi = len(values)
		}
		result[i] = values[lo:hi:hi]
		lo, hi = hi, hi+size
	}
	return result
	//shorter version (https://gist.github.com/mustafaturan/7a29e8251a7369645fb6c2965f8c2daf)
	//for int(sizeUint) < len(values) {
	//	values, chunks = values[sizeUint:], append(chunks, values[0:sizeUint:sizeUint])
	//}
	//return append(chunks, values)
}
