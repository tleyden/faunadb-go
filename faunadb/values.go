package faunadb

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

/*
Value represents valid FaunaDB values returned from the server. Values also implement the Expr interface.
They can be sent back and forth to FaunaDB with no extra escaping needed.

Then Get method is used to decode a FaunaDB value into a Go type. For example:

	var t time.Time

	faunaTime, _ := client.Query(Time("now"))
	_ := faunaTime.Get(&t)

The At method uses field extractors to transverse the data and reach to a more specific field.

	var firstEmail string

	profile, _ := client.Query(Ref("classes/profile/43"))
	profile.At(ObjKey("emails").AtIndex(0)).Get(&firstEmail)

For more information, check https://fauna.com/documentation/queries#values.
*/
type Value interface {
	Expr
	Get(interface{}) error // Decode a FaunaDB value into a native Go type
	At(Field) FieldValue   // Transverse the value using the field extractor informed
}

// StringV represents a valid JSON string.
type StringV string

// Get implements the Value interface by decoding the underlying value to either a StringV or a string type.
func (str StringV) Get(i interface{}) error { return newValueDecoder(i).assign(str) }

// At implements the Value interface by returning an invalid field since StringV is not transversable.
func (str StringV) At(field Field) FieldValue { return field.get(str) }

// LongV represents a valid JSON number.
type LongV int64

// Get implements the Value interface by decoding the underlying value to either a LongV or a numeric type.
func (num LongV) Get(i interface{}) error { return newValueDecoder(i).assign(num) }

// At implements the Value interface by returning an invalid field since LongV is not transversable.
func (num LongV) At(field Field) FieldValue { return field.get(num) }

// DoubleV represents a valid JSON double.
type DoubleV float64

// Get implements the Value interface by decoding the underlying value to either a DoubleV or a float type.
func (num DoubleV) Get(i interface{}) error { return newValueDecoder(i).assign(num) }

// At implements the Value interface by returning an invalid field since DoubleV is not transversable.
func (num DoubleV) At(field Field) FieldValue { return field.get(num) }

// BooleanV represents a valid JSON boolean.
type BooleanV bool

// Get implements the Value interface by decoding the underlying value to either a BooleanV or a boolean type.
func (boolean BooleanV) Get(i interface{}) error { return newValueDecoder(i).assign(boolean) }

// At implements the Value interface by returning an invalid field since BooleanV is not transversable.
func (boolean BooleanV) At(field Field) FieldValue { return field.get(boolean) }

// DateV represents a FaunaDB date type.
type DateV time.Time

// Get implements the Value interface by decoding the underlying value to either a DateV or a time.Time type.
func (date DateV) Get(i interface{}) error { return newValueDecoder(i).assign(date) }

// At implements the Value interface by returning an invalid field since DateV is not transversable.
func (date DateV) At(field Field) FieldValue { return field.get(date) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB date representation.
func (date DateV) MarshalJSON() ([]byte, error) {
	return escape("@date", time.Time(date).Format("2006-01-02"))
}

// TimeV represents a FaunaDB time type.
type TimeV time.Time

// Get implements the Value interface by decoding the underlying value to either a TimeV or a time.Time type.
func (localTime TimeV) Get(i interface{}) error { return newValueDecoder(i).assign(localTime) }

// At implements the Value interface by returning an invalid field since TimeV is not transversable.
func (localTime TimeV) At(field Field) FieldValue { return field.get(localTime) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB time representation.
func (localTime TimeV) MarshalJSON() ([]byte, error) {
	return escape("@ts", time.Time(localTime).Format("2006-01-02T15:04:05.999999999Z"))
}

// RefV represents a FaunaDB ref type.
type RefV struct {
	ID string
}

// Get implements the Value interface by decoding the underlying ref to a RefV.
func (ref RefV) Get(i interface{}) error { return newValueDecoder(i).assign(ref) }

// At implements the Value interface by returning an invalid field since RefV is not transversable.
func (ref RefV) At(field Field) FieldValue { return field.get(ref) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB ref representation.
func (ref RefV) MarshalJSON() ([]byte, error) { return escape("@ref", ref.ID) }

// SetRefV represents a FaunaDB setref type.
type SetRefV struct {
	Parameters map[string]Value
}

// Get implements the Value interface by decoding the underlying value to a SetRefV.
func (set SetRefV) Get(i interface{}) error { return newValueDecoder(i).assign(set) }

// At implements the Value interface by returning an invalid field since SetRefV is not transversable.
func (set SetRefV) At(field Field) FieldValue { return field.get(set) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB setref representation.
func (set SetRefV) MarshalJSON() ([]byte, error) { return escape("@set", set.Parameters) }

// ObjectV represents a FaunaDB object type.
type ObjectV map[string]Value

// Get implements the Value interface by decoding the underlying value to either a ObjectV or a native map type.
func (obj ObjectV) Get(i interface{}) error { return newValueDecoder(i).decodeMap(obj) }

// At implements the Value interface by transversing the object and extracting the field informed.
func (obj ObjectV) At(field Field) FieldValue { return field.get(obj) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB object representation.
func (obj ObjectV) MarshalJSON() ([]byte, error) { return escape("object", map[string]Value(obj)) }

// ArrayV represents a FaunaDB array type.
type ArrayV []Value

// Get implements the Value interface by decoding the underlying value to either an ArrayV or a native slice type.
func (arr ArrayV) Get(i interface{}) error { return newValueDecoder(i).decodeArray(arr) }

// At implements the Value interface by transversing the array and extracting the field informed.
func (arr ArrayV) At(field Field) FieldValue { return field.get(arr) }

// NullV represents a valid JSON null.
type NullV struct{}

// Get implements the Value interface by decoding the underlying value to a either a NullV or a nil pointer.
func (null NullV) Get(i interface{}) error { return nil }

// At implements the Value interface by returning an invalid field since NullV is not transversable.
func (null NullV) At(field Field) FieldValue { return field.get(null) }

// MarshalJSON implements json.Marshaler by escaping its value according to JSON null representation.
func (null NullV) MarshalJSON() ([]byte, error) { return []byte("null"), nil }

// BytesV represents a FaunaDB binary blob type.
type BytesV []byte

// Get implements the Value interface by decoding the underlying value to either a ByteV or a []byte type.
func (bytes BytesV) Get(i interface{}) error { return newValueDecoder(i).assign(bytes) }

// At implements the Value interface by returning an invalid field since BytesV is not transversable.
func (bytes BytesV) At(field Field) FieldValue { return field.get(bytes) }

// MarshalJSON implements json.Marshaler by escaping its value according to FaunaDB bytes representation.
func (bytes BytesV) MarshalJSON() ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return escape("@bytes", encoded)
}

// Implement Expr for all values

func (str StringV) expr()      {}
func (num LongV) expr()        {}
func (num DoubleV) expr()      {}
func (boolean BooleanV) expr() {}
func (date DateV) expr()       {}
func (localTime TimeV) expr()  {}
func (ref RefV) expr()         {}
func (set SetRefV) expr()      {}
func (obj ObjectV) expr()      {}
func (arr ArrayV) expr()       {}
func (null NullV) expr()       {}
func (bytes BytesV) expr()     {}

func escape(key string, value interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{key: value})
}
