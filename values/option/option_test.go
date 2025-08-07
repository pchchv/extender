package optionext

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"

	. "github.com/pchchv/go-assert"
)

type testStruct struct{}

type customStringType string

type testStructType struct {
	Name string
}

type valueTest struct {
}

func (valueTest) Value() (driver.Value, error) {
	return "value", nil
}

type customScanner struct {
	S string
}

func (c *customScanner) Scan(src interface{}) error {
	if src != nil {
		c.S = src.(string)
	}

	return nil
}

func TestAndXXX(t *testing.T) {
	s := Some(1)
	Equal(t, Some(3), s.And(func(i int) int { return 3 }))
	Equal(t, Some(3), s.AndThen(func(i int) Option[int] { return Some(3) }))
	Equal(t, None[int](), s.AndThen(func(i int) Option[int] { return None[int]() }))

	n := None[int]()
	Equal(t, None[int](), n.And(func(i int) int { return 3 }))
	Equal(t, None[int](), n.AndThen(func(i int) Option[int] { return Some(3) }))
	Equal(t, None[int](), n.AndThen(func(i int) Option[int] { return None[int]() }))
	Equal(t, None[int](), s.AndThen(func(i int) Option[int] { return None[int]() }))
}

func TestUnwraps(t *testing.T) {
	none := None[int]()
	PanicMatches(t, func() { none.Unwrap() }, "Option.Unwrap: option is None")

	v := none.UnwrapOr(3)
	Equal(t, 3, v)

	v = none.UnwrapOrElse(func() int { return 2 })
	Equal(t, 2, v)

	v = none.UnwrapOrDefault()
	Equal(t, 0, v)

	// now test with a pointer type.
	type testStruct struct {
		S string
	}

	sNone := None[*testStruct]()
	PanicMatches(t, func() { sNone.Unwrap() }, "Option.Unwrap: option is None")

	v2 := sNone.UnwrapOr(&testStruct{S: "blah"})
	Equal(t, &testStruct{S: "blah"}, v2)

	v2 = sNone.UnwrapOrElse(func() *testStruct { return &testStruct{S: "blah 2"} })
	Equal(t, &testStruct{S: "blah 2"}, v2)

	v2 = sNone.UnwrapOrDefault()
	Equal(t, nil, v2)
}

func TestNilOption(t *testing.T) {
	value := Some[any](nil)
	Equal(t, false, value.IsNone())
	Equal(t, true, value.IsSome())
	Equal(t, nil, value.Unwrap())

	ret := returnTypedNoneOption()
	Equal(t, true, ret.IsNone())
	Equal(t, false, ret.IsSome())
	PanicMatches(t, func() {
		ret.Unwrap()
	}, "Option.Unwrap: option is None")

	ret = returnTypedSomeOption()
	Equal(t, false, ret.IsNone())
	Equal(t, true, ret.IsSome())
	Equal(t, testStruct{}, ret.Unwrap())

	retPtr := returnTypedNoneOptionPtr()
	Equal(t, true, retPtr.IsNone())
	Equal(t, false, retPtr.IsSome())

	retPtr = returnTypedSomeOptionPtr()
	Equal(t, false, retPtr.IsNone())
	Equal(t, true, retPtr.IsSome())
	Equal(t, new(testStruct), retPtr.Unwrap())
}

func TestOptionJSON(t *testing.T) {
	type s struct {
		Timestamp Option[time.Time] `json:"ts"`
	}

	now := time.Now().UTC().Truncate(time.Minute)
	tv := s{Timestamp: Some(now)}

	b, err := json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":"`+now.Format(time.RFC3339)+`"}`, string(b))

	tv = s{}
	b, err = json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":null}`, string(b))
}

func TestOptionJSONOmitempty(t *testing.T) {
	type s struct {
		Timestamp Option[time.Time] `json:"ts,omitempty"`
	}

	now := time.Now().UTC().Truncate(time.Minute)
	tv := s{Timestamp: Some(now)}

	b, err := json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":"`+now.Format(time.RFC3339)+`"}`, string(b))

	type s2 struct {
		Timestamp *Option[time.Time] `json:"ts,omitempty"`
	}
	tv2 := &s2{}
	b, err = json.Marshal(tv2)
	Equal(t, nil, err)
	Equal(t, `{}`, string(b))
}

func TestSQLDriverValue(t *testing.T) {
	var v valueTest
	Equal(t, reflect.TypeOf(v).Implements(valuerType), true)

	// none
	nOpt := None[string]()
	nVal, err := nOpt.Value()
	Equal(t, err, nil)
	Equal(t, nVal, nil)

	// string + convert custom string type
	sOpt := Some("myString")
	sVal, err := sOpt.Value()
	Equal(t, err, nil)

	_, ok := sVal.(string)
	Equal(t, ok, true)
	Equal(t, sVal, "myString")

	sCustOpt := Some(customStringType("string"))
	sCustVal, err := sCustOpt.Value()
	Equal(t, err, nil)
	Equal(t, sCustVal, "string")

	_, ok = sCustVal.(string)
	Equal(t, ok, true)

	// bool
	bOpt := Some(true)
	bVal, err := bOpt.Value()
	Equal(t, err, nil)

	_, ok = bVal.(bool)
	Equal(t, ok, true)
	Equal(t, bVal, true)

	// int64
	iOpt := Some(int64(2))
	iVal, err := iOpt.Value()
	Equal(t, err, nil)

	_, ok = iVal.(int64)
	Equal(t, ok, true)
	Equal(t, iVal, int64(2))

	// float64
	fOpt := Some(1.1)
	fVal, err := fOpt.Value()
	Equal(t, err, nil)

	_, ok = fVal.(float64)
	Equal(t, ok, true)
	Equal(t, fVal, 1.1)

	// time.Time
	dt := time.Now().UTC()
	dtOpt := Some(dt)
	dtVal, err := dtOpt.Value()
	Equal(t, err, nil)

	_, ok = dtVal.(time.Time)
	Equal(t, ok, true)
	Equal(t, dtVal, dt)

	// Slice []byte
	b := []byte("myBytes")
	bytesOpt := Some(b)
	bytesVal, err := bytesOpt.Value()
	Equal(t, err, nil)

	_, ok = bytesVal.([]byte)
	Equal(t, ok, true)
	Equal(t, bytesVal, b)

	// Slice []uint8
	b2 := []uint8("myBytes")
	bytes2Opt := Some(b2)
	bytes2Val, err := bytes2Opt.Value()
	Equal(t, err, nil)

	_, ok = bytes2Val.([]byte)
	Equal(t, ok, true)
	Equal(t, bytes2Val, b2)

	// Array []byte
	a := []byte{'1', '2', '3'}
	arrayOpt := Some(a)
	arrayVal, err := arrayOpt.Value()
	Equal(t, err, nil)

	_, ok = arrayVal.([]byte)
	Equal(t, ok, true)
	Equal(t, arrayVal, a)

	// Slice []byte
	data := []testStructType{{Name: "test"}}
	b, err = json.Marshal(data)
	Equal(t, err, nil)

	dataOpt := Some(data)
	dataVal, err := dataOpt.Value()
	Equal(t, err, nil)

	_, ok = dataVal.([]byte)
	Equal(t, ok, true)
	Equal(t, dataVal, b)

	// Map
	data2 := map[string]int{"test": 1}
	b, err = json.Marshal(data2)
	Equal(t, err, nil)

	data2Opt := Some(data2)
	data2Val, err := data2Opt.Value()
	Equal(t, err, nil)

	_, ok = data2Val.([]byte)
	Equal(t, ok, true)
	Equal(t, data2Val, b)

	// Struct
	data3 := testStructType{Name: "test"}
	b, err = json.Marshal(data3)
	Equal(t, err, nil)

	data3Opt := Some(data3)
	data3Val, err := data3Opt.Value()
	Equal(t, err, nil)

	_, ok = data3Val.([]byte)
	Equal(t, ok, true)
	Equal(t, data3Val, b)
}

func TestSQLScanner(t *testing.T) {
	var optionI Option[int]
	var optionI8 Option[int8]
	var optionI16 Option[int16]
	var optionI32 Option[int32]
	var optionI64 Option[int64]
	var optionString Option[string]
	var optionBool Option[bool]
	var optionF32 Option[float32]
	var optionF64 Option[float64]
	var optionByte Option[byte]
	var optionTime Option[time.Time]
	var optionInterface Option[any]
	var optionArrBytes Option[[]byte]
	var optionRawMessage Option[json.RawMessage]
	var optionUint Option[uint]
	var optionUint8 Option[uint8]
	var optionUint16 Option[uint16]
	var optionUint32 Option[uint32]
	var optionUint64 Option[uint64]
	err := optionInterface.Scan(1)
	value := int64(123)
	Equal(t, err, nil)
	Equal(t, optionInterface, Some(any(1)))

	err = optionInterface.Scan("blah")
	Equal(t, err, nil)
	Equal(t, optionInterface, Some(any("blah")))

	err = optionUint64.Scan(uint64(200))
	Equal(t, err, nil)
	Equal(t, optionUint64, Some(uint64(200)))

	err = optionUint32.Scan(uint32(200))
	Equal(t, err, nil)
	Equal(t, optionUint32, Some(uint32(200)))

	err = optionUint16.Scan(uint16(200))
	Equal(t, err, nil)
	Equal(t, optionUint16, Some(uint16(200)))

	err = optionUint8.Scan(uint8(200))
	Equal(t, err, nil)
	Equal(t, optionUint8, Some(uint8(200)))

	err = optionUint.Scan(uint(200))
	Equal(t, err, nil)
	Equal(t, optionUint, Some(uint(200)))

	err = optionUint64.Scan("200")
	Equal(t, err.Error(), "value string not convertable to uint64")

	err = optionI64.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionI64, Some(value))

	err = optionI32.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionI32, Some(int32(value)))

	err = optionI16.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionI16, Some(int16(value)))

	err = optionI8.Scan(math.MaxInt32)
	Equal(t, err.Error(), "value 2147483647 out of range for int8")
	Equal(t, optionI8, None[int8]())

	err = optionI8.Scan(int8(3))
	Equal(t, err, nil)
	Equal(t, optionI8, Some(int8(3)))

	err = optionI.Scan(3)
	Equal(t, err, nil)
	Equal(t, optionI, Some(3))

	err = optionBool.Scan(1)
	Equal(t, err, nil)
	Equal(t, optionBool, Some(true))

	err = optionString.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionString, Some("123"))

	err = optionF32.Scan(float32(2.0))
	Equal(t, err, nil)
	Equal(t, optionF32, Some(float32(2.0)))

	err = optionF32.Scan(math.MaxFloat64)
	Equal(t, err, nil)
	Equal(t, optionF32, Some(float32(math.Inf(1))))

	err = optionF64.Scan(2.0)
	Equal(t, err, nil)
	Equal(t, optionF64, Some(2.0))

	err = optionByte.Scan(uint8('1'))
	Equal(t, err, nil)
	Equal(t, optionByte, Some(uint8('1')))

	err = optionTime.Scan("2023-06-13T06:34:32Z")
	Equal(t, err, nil)
	Equal(t, optionTime, Some(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC)))

	err = optionTime.Scan([]byte("2023-06-13T06:34:32Z"))
	Equal(t, err, nil)
	Equal(t, optionTime, Some(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC)))

	err = optionTime.Scan(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC))
	Equal(t, err, nil)
	Equal(t, optionTime, Some(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC)))

	// Test nil
	var nullableOption Option[int64]
	err = nullableOption.Scan(nil)
	Equal(t, err, nil)
	Equal(t, nullableOption, None[int64]())

	// custom scanner
	var custom Option[customScanner]
	err = custom.Scan("GOT HERE")
	Equal(t, err, nil)
	Equal(t, custom, Some(customScanner{S: "GOT HERE"}))

	// custom scanner scan nil
	var customNil Option[customScanner]
	err = customNil.Scan(nil)
	Equal(t, err, nil)
	Equal(t, customNil, None[customScanner]())

	// test unmarshal to struct
	type testStruct struct {
		Name string `json:"name"`
	}

	var optiontestStruct Option[testStruct]
	err = optiontestStruct.Scan([]byte(`{"name":"test"}`))
	Equal(t, err, nil)
	Equal(t, optiontestStruct, Some(testStruct{Name: "test"}))

	err = optiontestStruct.Scan(json.RawMessage(`{"name":"test2"}`))
	Equal(t, err, nil)
	Equal(t, optiontestStruct, Some(testStruct{Name: "test2"}))

	var optionArrayOftestStruct Option[[]testStruct]
	err = optionArrayOftestStruct.Scan([]byte(`[{"name":"test"}]`))
	Equal(t, err, nil)
	Equal(t, optionArrayOftestStruct, Some([]testStruct{{Name: "test"}}))

	var optionMap Option[map[string]any]
	err = optionMap.Scan([]byte(`{"name":"test"}`))
	Equal(t, err, nil)
	Equal(t, optionMap, Some(map[string]any{"name": "test"}))

	// test custom types
	var ct Option[customStringType]
	err = ct.Scan("test")
	Equal(t, err, nil)
	Equal(t, ct, Some(customStringType("test")))

	err = optionArrBytes.Scan([]byte(`[1,2,3]`))
	Equal(t, err, nil)
	Equal(t, optionArrBytes, Some([]byte(`[1,2,3]`)))

	err = optionArrBytes.Scan([]byte{4, 5, 6})
	Equal(t, err, nil)
	Equal(t, optionArrBytes, Some([]byte{4, 5, 6}))

	err = optionRawMessage.Scan([]byte(`[1,2,3]`))
	Equal(t, err, nil)
	Equal(t, true, string(optionRawMessage.Unwrap()) == "[1,2,3]")

	err = optionRawMessage.Scan([]byte{4, 5, 6})
	Equal(t, err, nil)
	Equal(t, true, string(optionRawMessage.Unwrap()) == string([]byte{4, 5, 6}))
}

func BenchmarkOption(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := returnTypedSomeOption()
		if opt.IsSome() {
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkOptionPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if opt := returnTypedSomeOptionPtr(); opt.IsSome() {
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNoOptionPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if result := returnTypedNoOption(); result != nil {
			_ = result
		}
	}
}

func BenchmarkOptionNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if opt := returnTypedSomeOptionNil(); opt.IsSome() {
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNoOptionNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if result, found := returnNoOptionNil(); found {
			_ = result
		}
	}
}

func returnTypedNoneOption() Option[testStruct] {
	return None[testStruct]()
}

func returnTypedSomeOption() Option[testStruct] {
	return Some(testStruct{})
}

func returnTypedNoneOptionPtr() Option[*testStruct] {
	return None[*testStruct]()
}

func returnTypedSomeOptionPtr() Option[*testStruct] {
	return Some(new(testStruct))
}

func returnTypedSomeOptionNil() Option[any] {
	return Some[any](nil)
}

func returnTypedNoOption() *testStruct {
	return new(testStruct)
}

func returnNoOptionNil() (any, bool) {
	return nil, true
}
