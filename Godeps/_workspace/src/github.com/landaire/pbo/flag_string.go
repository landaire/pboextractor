// generated by stringer -type=Flag; DO NOT EDIT

package pbo

import "fmt"

const (
	_Flag_name_0 = "Uncompressed"
	_Flag_name_1 = "Packed"
	_Flag_name_2 = "ProductEntry"
)

var (
	_Flag_index_0 = [...]uint8{12}
	_Flag_index_1 = [...]uint8{6}
	_Flag_index_2 = [...]uint8{12}
)

func (i Flag) String() string {
	switch {
	case i == 0:
		return _Flag_name_0
	case i == 1131442803:
		return _Flag_name_1
	case i == 1449489011:
		return _Flag_name_2
	default:
		return fmt.Sprintf("Flag(%d)", i)
	}
}
