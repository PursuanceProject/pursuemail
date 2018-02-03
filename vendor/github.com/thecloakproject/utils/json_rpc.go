// Steve Phillips / elimisteve
// 2013.02.07

package utils

import (
	"fmt"
)

// SumEmptyInterfaceSlice sums the numbers in params. If params
// includes non-numbers, or a number that cannot be parsed as a
// float64, the resulting sum will be 0 and a non-nil error will be
// returned. Often used to sum values given in JSON-RPC requests.
func SumEmptyInterfaceSlice(params []interface{}) (sum float64, err error) {
	// Parse params as a slice of float64s to add
	for _, n := range params {
		num, ok := n.(float64)
		if ok {
			sum += num
			continue
		}
		sum = 0
		err = fmt.Errorf("Couldn't parse params '%+v' as float64s", params)
		break
	}
	return
}

// ErrToEmptyInterface mostly just makes up for the fact that Go
// marshals interface{} values to '{}' even when its underlying value
// is of type error
func ErrToEmptyInterface(err error) (errVal interface{}) {
	// Set `errStr = err.Error()` if an error occurred. Simply setting
	// `Error: err.Error()` below isn't safe because `err` could be
	// nil, plus the RFC says if there's no error, `Error` should
	// equal nil, not empty string.
	if err != nil {
		errVal = err.Error()
	} else {
		errVal = nil
	}
	return
}