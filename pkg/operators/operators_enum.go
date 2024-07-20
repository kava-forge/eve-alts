// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package operators

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

const (
	// OperatorAll is a Operator of type all.
	OperatorAll Operator = "all"
	// OperatorAny is a Operator of type any.
	OperatorAny Operator = "any"
	// OperatorNone is a Operator of type none.
	OperatorNone Operator = "none"
)

var ErrInvalidOperator = fmt.Errorf("not a valid Operator, try [%s]", strings.Join(_OperatorNames, ", "))

var _OperatorNames = []string{
	string(OperatorAll),
	string(OperatorAny),
	string(OperatorNone),
}

// OperatorNames returns a list of possible string values of Operator.
func OperatorNames() []string {
	tmp := make([]string, len(_OperatorNames))
	copy(tmp, _OperatorNames)
	return tmp
}

// OperatorValues returns a list of the values for Operator
func OperatorValues() []Operator {
	return []Operator{
		OperatorAll,
		OperatorAny,
		OperatorNone,
	}
}

// String implements the Stringer interface.
func (x Operator) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Operator) IsValid() bool {
	_, err := ParseOperator(string(x))
	return err == nil
}

var _OperatorValue = map[string]Operator{
	"all":  OperatorAll,
	"any":  OperatorAny,
	"none": OperatorNone,
}

// ParseOperator attempts to convert a string to a Operator.
func ParseOperator(name string) (Operator, error) {
	if x, ok := _OperatorValue[name]; ok {
		return x, nil
	}
	return Operator(""), fmt.Errorf("%s is %w", name, ErrInvalidOperator)
}

// MarshalText implements the text marshaller method.
func (x Operator) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Operator) UnmarshalText(text []byte) error {
	tmp, err := ParseOperator(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errOperatorNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *Operator) Scan(value interface{}) (err error) {
	if value == nil {
		*x = Operator("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseOperator(v)
	case []byte:
		*x, err = ParseOperator(string(v))
	case Operator:
		*x = v
	case *Operator:
		if v == nil {
			return errOperatorNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errOperatorNilPtr
		}
		*x, err = ParseOperator(*v)
	default:
		return errors.New("invalid type for Operator")
	}

	return
}

// Value implements the driver Valuer interface.
func (x Operator) Value() (driver.Value, error) {
	return x.String(), nil
}
