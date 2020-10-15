package enrollment

import (
	"encoding/json"
	"fmt"
)

type DBType string

const (
	TypeString    DBType = "g:String"
	TypeInteger   DBType = "g:Int32"
	TypeBoolean   DBType = "g:Boolean"
	TypeList      DBType = "g:List"
	TypeMap       DBType = "g:Map"
	TypeFloat     DBType = "g:Float"
	TypeTimestamp DBType = "g:Timestamp"
	TypeLong      DBType = "g:Long"
	TypeDate      DBType = "g:Date"
	TypeDouble    DBType = "g:Double"
)

var dbTypes = map[string]DBType{
	"g:String":  TypeString,
	"g:Int32":   TypeInteger,
	"g:Boolean": TypeBoolean,
	"g:List":    TypeList,
	"g:Map":     TypeMap,
	"g:Float":   TypeFloat,
}

type unmarshal func(raw []byte, v *interface{}) error

var dbUnmarshals = map[DBType]unmarshal{
	TypeString:  toString,
	TypeInteger: toInt32,
	TypeFloat:   toFloat64,
	TypeBoolean: toBool,
	TypeList:    toList,
	TypeMap:     toMap,
}

func toString(raw []byte, v *interface{}) error {
	var val string
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*v = val
	return nil
}

func toInt32(raw []byte, v *interface{}) error {
	var val int32
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*v = val
	return nil
}

func toFloat64(raw []byte, v *interface{}) error {
	var val float64
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*v = val
	return nil
}

func toBool(raw []byte, v *interface{}) error {
	var val bool
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*v = val
	return nil
}

func toList(raw []byte, v *interface{}) error {
	var val List
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*v = val
	return nil
}

func toMap(raw []byte, v *interface{}) error {
	var val Map
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*v = val
	return nil
}

type ListOfStrings []string

func (l *ListOfStrings) UnmarshalJSON(b []byte) error {
	*l = nil
	var a Attribute
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	if a.Type != TypeList {
		return fmt.Errorf("got %s where %s is expected", a.Type, TypeList)
	}
	values := a.ListValue()
	for _, v := range values {
		*l = append(*l, v.ToString())
	}
	return nil
}

func UnmarshalStringList(recs [][]byte) ([]string, error) {
	var err error
	var items []string
	var list ListOfStrings
	for _, r := range recs {
		err = json.Unmarshal(r, &list)
		if err != nil {
			return nil, err
		}
		items = append(items, list...)
	}
	return items, nil
}

type ListOfInt32 []int32

func (l *ListOfInt32) UnmarshalJSON(b []byte) error {
	*l = nil
	var a Attribute
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	if a.Type != TypeList {
		return fmt.Errorf("got %s where %s is expected", a.Type, TypeList)
	}
	values := a.ListValue()
	for _, v := range values {
		*l = append(*l, v.Int32Value())
	}
	return nil
}

func UnmarshalInt32List(recs [][]byte) ([]int32, error) {
	var err error
	var items []int32
	var list ListOfInt32
	for _, r := range recs {
		err = json.Unmarshal(r, &list)
		if err != nil {
			return nil, err
		}
		items = append(items, list...)
	}
	return items, nil
}

type dbRecord struct {
	Type string          `json:"@type"`
	Raw  json.RawMessage `json:"@value"`
}

type List []Attribute

type Map map[string]Attribute

type Attribute struct {
	Type  DBType
	Value interface{}
}

func (a *Attribute) UnmarshalJSON(b []byte) error {
	var rec dbRecord
	if err := json.Unmarshal(b, &rec); err != nil {
		return err
	}
	t, ok := dbTypes[rec.Type]
	if !ok {
		return fmt.Errorf("unknown type: %s", a.Type)
	}
	a.Type = t
	unmarshalValue, ok := dbUnmarshals[a.Type]
	if !ok {
		unmarshalValue = toString
	}
	if err := unmarshalValue(rec.Raw, &a.Value); err != nil {
		return err
	}
	return nil
}

func (a Attribute) ToString() string {
	v, ok := a.Value.(string)
	if !ok {
		return fmt.Sprintf("%v", a.Value)
	}
	return v
}

func (a Attribute) StringValue() string {
	v, ok := a.Value.(string)
	if !ok {
		return ""
	}
	return v
}

func (a Attribute) BoolValue() bool {
	v, ok := a.Value.(bool)
	if !ok {
		return false
	}
	return v
}

func (a Attribute) Int32Value() int32 {
	v, ok := a.Value.(int32)
	if !ok {
		return 0
	}
	return v
}

func (a Attribute) Float64Value() float64 {
	v, ok := a.Value.(float64)
	if !ok {
		return 0
	}
	return v
}

func (a Attribute) ListValue() List {
	v, ok := a.Value.(List)
	if !ok {
		return List{}
	}
	return v
}

func (a Attribute) MapValue() Map {
	v, ok := a.Value.(Map)
	if !ok {
		return make(Map)
	}
	return v
}
