package fielder

import (
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
	"time"
)

type Parent interface {
	GetResultItemFieldFromKey(f FieldKey) Field
	GetFieldTypeFromKey(f FieldKey) reflect.Type
	GetReflectValueOfKey(f FieldKey) reflect.Value
	CheckKeyExists(f FieldKey) bool
}

type FieldKey string

const FieldKeyNil FieldKey = "nil"

var FieldNil = CreateFieldFromType((&EmptyField{}).Type(), nil, FieldKeyNil)

// field interface
type Field interface {
	Value() any         // value of the field
	Key() FieldKey      // name of the field
	Type() reflect.Type // type of field == string, time.Time, decimal.Decimal
	LessThan(in2 any) bool
	GreaterThan(in2 any) bool
	Equal(in2 any) bool
	ToString() string
	FromString(st string)
}

type StringField struct {
	ValueField string
	KeyField   FieldKey
}

func (s *StringField) Value() any {
	return s.ValueField
}

func (s *StringField) Key() FieldKey {
	return s.KeyField
}

func (s *StringField) Type() reflect.Type {
	return reflect.TypeOf("")
}

func (s *StringField) LessThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, LT); out != nil {
		return *out
	}
	return s.ValueField < in2.(*StringField).ValueField
}

func (s *StringField) GreaterThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, GT); out != nil {
		return *out
	}
	return s.ValueField > in2.(*StringField).ValueField
}

func (s *StringField) Equal(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, EQ); out != nil {
		return *out
	}
	return s.ValueField == in2.(*StringField).ValueField
}

func (s *StringField) ToString() string {
	return s.ValueField
}

func (s *StringField) FromString(st string) {
	s.ValueField = st
}

type TimeField struct {
	ValueField time.Time
	KeyField   FieldKey
}

func (s *TimeField) Value() any {
	return s.ValueField
}

func (s *TimeField) Key() FieldKey {
	return s.KeyField
}

func (s *TimeField) Type() reflect.Type {
	return reflect.TypeOf(time.Time{})
}

func (s *TimeField) LessThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, LT); out != nil {
		return *out
	}
	return s.ValueField.Before(in2.(*TimeField).ValueField)
}

func (s *TimeField) GreaterThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, GT); out != nil {
		return *out
	}
	return s.ValueField.After(in2.(*TimeField).ValueField)
}

func (s *TimeField) Equal(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, EQ); out != nil {
		return *out
	}
	return s.ValueField.Equal(in2.(*TimeField).ValueField)
}

func (s *TimeField) ToString() string {
	return s.ValueField.Format(time.RFC3339)
}

func (s *TimeField) FromString(st string) {
	t, err := time.Parse(time.RFC3339, st)
	if err != nil {
		return
	}
	s.ValueField = t
}

type DecimalField struct {
	ValueField decimal.Decimal
	KeyField   FieldKey
}

func (s *DecimalField) Value() any {
	return s.ValueField
}

func (s *DecimalField) Key() FieldKey {
	return s.KeyField
}

func (s *DecimalField) Type() reflect.Type {
	return reflect.TypeOf(decimal.Decimal{})
}

func (s *DecimalField) LessThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, LT); out != nil {
		return *out
	}
	return s.ValueField.LessThan(in2.(*DecimalField).ValueField)
}

func (s *DecimalField) GreaterThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, GT); out != nil {
		return *out
	}
	return s.ValueField.GreaterThan(in2.(*DecimalField).ValueField)
}

func (s *DecimalField) Equal(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, EQ); out != nil {
		return *out
	}
	return s.ValueField.Equal(in2.(*DecimalField).ValueField)
}

func (s *DecimalField) ToString() string {
	return s.ValueField.String()
}

func (s *DecimalField) FromString(st string) {
	d, err := decimal.NewFromString(st)
	if err != nil {
		return
	}
	s.ValueField = d
}

type IntegerField struct {
	ValueField int
	KeyField   FieldKey
}

func (s *IntegerField) Value() any {
	return s.ValueField
}

func (s *IntegerField) Key() FieldKey {
	return s.KeyField
}

func (s *IntegerField) Type() reflect.Type {
	return reflect.TypeOf(int(0))
}

func (s *IntegerField) LessThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, LT); out != nil {
		return *out
	}
	return s.ValueField < in2.(*IntegerField).ValueField
}

func (s *IntegerField) GreaterThan(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, GT); out != nil {
		return *out
	}
	return s.ValueField > in2.(*IntegerField).ValueField
}

func (s *IntegerField) Equal(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, EQ); out != nil {
		return *out
	}
	return s.ValueField == in2.(*IntegerField).ValueField
}

func (s *IntegerField) ToString() string {
	return strconv.Itoa(s.ValueField)
}

func (s *IntegerField) FromString(st string) {
	it, err := strconv.Atoi(st)
	if err != nil {
		s.ValueField = 0
	}
	s.ValueField = it
}

type BoolField struct {
	ValueField bool
	KeyField   FieldKey
}

func (s *BoolField) Value() any {
	return s.ValueField
}

func (s *BoolField) Key() FieldKey {
	return s.KeyField
}

func (s *BoolField) Type() reflect.Type {
	return reflect.TypeOf(true)
}

// less than and greater than are not relevant for bool
func (s *BoolField) LessThan(in2 any) bool {
	return false
}

func (s *BoolField) GreaterThan(in2 any) bool {
	return false
}

func (s *BoolField) Equal(in2 any) bool {
	if out := checkAndDoSafeCompare(s, in2, EQ); out != nil {
		return *out
	}
	return s.ValueField == in2.(*BoolField).ValueField
}

func (s *BoolField) ToString() string {
	return strconv.FormatBool(s.ValueField)
}

func (s *BoolField) FromString(st string) {
	it, err := strconv.ParseBool(st)
	if err != nil {
		s.ValueField = false
	}
	s.ValueField = it
}

type EmptyField struct {
	KeyField FieldKey
}

func (s *EmptyField) Value() any {
	return nil
}

func (s *EmptyField) Key() FieldKey {
	return s.KeyField
}

func (s *EmptyField) Type() reflect.Type {
	return reflect.TypeOf(&EmptyField{})
}

// less than and greater than are not relevant for bool
func (s *EmptyField) LessThan(in2 any) bool {
	return false
}

func (s *EmptyField) GreaterThan(in2 any) bool {
	return false
}

func (s *EmptyField) Equal(in2 any) bool {
	if in2 == nil {
		return false
	}
	f2 := in2.(Field)
	if noSafeCheck := sameCompareTypes(s, f2); !noSafeCheck {
		return false
	}
	return s.Key() == in2.(*EmptyField).Key()
}

func (s *EmptyField) ToString() string {
	return ""
}

func (s *EmptyField) FromString(st string) {
	s.KeyField = FieldKey(st)
}

func IsFieldKey(s string, keySet []FieldKey) bool {
	return s != "" && SliceContains[FieldKey](keySet, FieldKey(s), func(s1, s2 FieldKey) bool {
		return s1 == s2
	})
}

func FullKeySet[inType any](tag string) []FieldKey {
	keySet := []FieldKey{}
	reflectType := reflect.TypeOf(*new(inType))
	for i := 0; i < reflectType.NumField(); i++ {
		keySet = append(keySet, FieldKey(reflectType.Field(i).Tag.Get(tag)))
	}
	return keySet
}

func CreateFieldFromType(ty reflect.Type, va any, fk FieldKey) Field {
	s := ""
	t := time.Time{}
	d := decimal.Decimal{}
	i := int(0)
	b := true
	if ty == nil {
		return nil
	}
	switch ty {
	case reflect.TypeOf(s):
		if va == nil {
			return &StringField{
				KeyField: fk,
			}
		}
		return &StringField{
			ValueField: va.(string),
			KeyField:   fk,
		}
	case reflect.TypeOf(t):
		if va == nil {
			return &TimeField{
				KeyField: fk,
			}
		}
		return &TimeField{
			ValueField: va.(time.Time),
			KeyField:   fk,
		}
	case reflect.TypeOf(d):
		if va == nil {
			return &DecimalField{
				KeyField: fk,
			}
		}
		return &DecimalField{
			ValueField: va.(decimal.Decimal),
			KeyField:   fk,
		}
	case reflect.TypeOf(i):
		if va == nil {
			return &IntegerField{
				KeyField: fk,
			}
		}
		return &IntegerField{
			ValueField: va.(int),
			KeyField:   fk,
		}
	case reflect.TypeOf(b):
		if va == nil {
			return &BoolField{
				KeyField: fk,
			}
		}
		return &BoolField{
			ValueField: va.(bool),
			KeyField:   fk,
		}
	case reflect.TypeOf(&EmptyField{}):
		return &EmptyField{KeyField: fk}
	default:
		// THIS SHOULD NEVER HAPPEN
		return nil
	}
}

func SliceContains[valuetype any](s []valuetype, val valuetype, equals func(valuetype, valuetype) bool) bool {
	for _, v := range s {
		if equals(v, val) {
			return true
		}
	}
	return false
}

func checkAndDoSafeCompare(f1 Field, f2 any, o safeOp) *bool {
	pointerTo := func(in bool) *bool {
		return &in
	}
	// we know f1 isnt nil because we got here from its receiver function
	if f2 == nil {
		return pointerTo(false)
	}
	f2f := f2.(Field)
	if noSafeCheck := sameCompareTypes(f1, f2f); !noSafeCheck {
		return pointerTo(safeCompare(f1.ToString(), f2f.ToString(), o))
	}
	// return nil if we need to do a same type comparison
	return nil
}

// for safe operations between different types
func sameCompareTypes(f1, f2 Field) bool {
	return f1.Type() == f2.Type()
}

func safeCompare(f1 string, f2 string, o safeOp) bool {
	return safeOps[o](f1, f2)
}

type safeOp string
type Op func(s1, s2 string) bool

var (
	EQ safeOp = "EQ"
	LT safeOp = "LT"
	GT safeOp = "GT"

	safeOps = map[safeOp]Op{
		EQ: func(s1, s2 string) bool {
			return s1 == s2
		},
		LT: func(s1, s2 string) bool {
			return s1 < s2
		},
		GT: func(s1, s2 string) bool {
			return s1 > s2
		},
	}
)
