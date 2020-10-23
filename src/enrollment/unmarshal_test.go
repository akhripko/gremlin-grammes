package enrollment

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalInt32List(t *testing.T) {
	js1 := []byte(`{"@type":"g:List","@value":[{"@type":"g:Int32","@value":1},{"@type":"g:Int32","@value":2},{"@type":"g:Int32","@value":3}]}`)
	js2 := []byte(`{"@type":"g:List","@value":[{"@type":"g:Int32","@value":123},{"@type":"g:Int32","@value":456},{"@type":"g:Int32","@value":789}]}`)

	recs := [][]byte{js1, js2}

	ids, err := UnmarshalInt32List(recs)
	require.NoError(t, err)
	assert.Equal(t, int32(1), ids[0])
	assert.Equal(t, int32(2), ids[1])
	assert.Equal(t, int32(3), ids[2])
	assert.Equal(t, int32(123), ids[3])
	assert.Equal(t, int32(456), ids[4])
	assert.Equal(t, int32(789), ids[5])
}

func TestUnmarshalStringList(t *testing.T) {
	js1 := []byte(`{"@type":"g:List","@value":[{"@type":"g:Int32","@value":1},{"@type":"g:Int32","@value":2},{"@type":"g:Int32","@value":3}]}`)
	js2 := []byte(`{"@type":"g:List","@value":[{"@type":"g:Int32","@value":123},{"@type":"g:Int32","@value":456},{"@type":"g:Int32","@value":789}]}`)

	recs := [][]byte{js1, js2}

	ids, err := UnmarshalStringList(recs)
	require.NoError(t, err)
	assert.Equal(t, "1", ids[0])
	assert.Equal(t, "2", ids[1])
	assert.Equal(t, "3", ids[2])
	assert.Equal(t, "123", ids[3])
	assert.Equal(t, "456", ids[4])
	assert.Equal(t, "789", ids[5])
}

func TestAttribute_UnmarshalJSON_String_Property(t *testing.T) {
	js := []byte(`{"@type":"g:Property","@value":{"key":"service","value":"childCare"}}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeProperty, a.Type)
	v := a.PropertyValue()
	assert.Equal(t, "service", v.Key)
	assert.Equal(t, TypeString, v.Value.Type)
	assert.Equal(t, "childCare", v.Value.StringValue())
}

func TestAttribute_UnmarshalJSON_Property(t *testing.T) {
	js := []byte(`{"@type":"g:Property","@value":{"key":"min_rate","value":{"@type":"g:Int32","@value":15}}}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeProperty, a.Type)
	v := a.PropertyValue()
	assert.Equal(t, "min_rate", v.Key)
	assert.Equal(t, TypeInteger, v.Value.Type)
	assert.Equal(t, int32(15), v.Value.Int32Value())
}

func TestAttribute_UnmarshalJSON_VertexProperty(t *testing.T) {
	js := []byte(`{"@type":"g:VertexProperty","@value":{"id":{"@type":"g:Int32","@value":-123},"label":"min_rate","value":{"@type":"g:Int32","@value":5}}}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeVertexProperty, a.Type)
	v := a.VertexPropertyValue()
	assert.Equal(t, int32(-123), v.ID.Int32Value())
	assert.Equal(t, "min_rate", v.Label)
	assert.Equal(t, TypeInteger, v.Value.Type)
	assert.Equal(t, int32(5), v.Value.Int32Value())
}

func TestAttribute_UnmarshalJSON_VertexProperty_AsString(t *testing.T) {
	js := []byte(`{"@type":"g:VertexProperty","@value":{"id":{"@type":"g:Int32","@value":123},"label":"service","value":"petCare"}}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeVertexProperty, a.Type)
	v := a.VertexPropertyValue()
	assert.Equal(t, int32(123), v.ID.Int32Value())
	assert.Equal(t, "service", v.Label)
	assert.Equal(t, TypeString, v.Value.Type)
	assert.Equal(t, "petCare", v.Value.StringValue())
}

func TestAttribute_UnmarshalJSON_Properties(t *testing.T) {
	js := []byte(`{"@type":"g:List","@value":[{"@type":"g:VertexProperty","@value":{"id":{"@type":"g:Int32","@value":67337261},"label":"service","value":"petCare"}},{"@type":"g:VertexProperty","@value":{"id":{"@type":"g:Int32","@value":-737478256},"label":"min_rate","value":{"@type":"g:Int32","@value":5}}},{"@type":"g:VertexProperty","@value":{"id":{"@type":"g:Int32","@value":1085490811},"label":"max_rate","value":{"@type":"g:Int32","@value":40}}},{"@type":"g:VertexProperty","@value":{"id":{"@type":"g:Int32","@value":-1866499204},"label":"years_of_exp","value":{"@type":"g:Int32","@value":10}}},{"@type":"g:VertexProperty","@value":{"id":{"@type":"g:Int32","@value":-465014591},"label":"avg_rank","value":{"@type":"g:Float","@value":5}}}]}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeList, a.Type)
	v := a.ListValue()
	assert.Equal(t, TypeVertexProperty, v[0].Type)
	vp := v[0].VertexPropertyValue()
	assert.Equal(t, "service", vp.Label)
	assert.Equal(t, "petCare", vp.Value.StringValue())
}

func TestAttribute_UnmarshalJSON_PropertyValues(t *testing.T) {
	js := []byte(`{"@type":"g:List","@value":["petCare",{"@type":"g:Int32","@value":5},{"@type":"g:Int32","@value":40},{"@type":"g:Int32","@value":10},{"@type":"g:Float","@value":5}]}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeList, a.Type)
	v := a.ListValue()
	assert.Equal(t, TypeString, v[0].Type)
	assert.Equal(t, "petCare", v[0].StringValue())
	assert.Equal(t, TypeInteger, v[1].Type)
	assert.Equal(t, int32(5), v[1].Int32Value())
}

func TestAttribute_UnmarshalJSON_Int32(t *testing.T) {
	js := []byte(`{"@type":"g:Int32","@value":123}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeInteger, a.Type)
	assert.Equal(t, int32(123), a.Int32Value())
}

func TestAttribute_String_UnmarshalJSON(t *testing.T) {
	js := []byte(`"abc"`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeString, a.Type)
	assert.Equal(t, "abc", a.StringValue())
}

func TestAttribute_UnmarshalJSON_String(t *testing.T) {
	js := []byte(`{"@type":"g:String","@value":"abc"}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeString, a.Type)
	assert.Equal(t, "abc", a.StringValue())
}

func TestAttribute_UnmarshalJSON_Bool(t *testing.T) {
	js := []byte(`{"@type":"g:Boolean","@value":true}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeBoolean, a.Type)
	assert.True(t, a.BoolValue())
}

func TestAttribute_UnmarshalJSON_Float(t *testing.T) {
	js := []byte(`{"@type":"g:Float","@value":1.23}`)
	var a Attribute
	err := json.Unmarshal(js, &a)
	require.NoError(t, err)
	assert.Equal(t, TypeFloat, a.Type)
	assert.Equal(t, float64(1.23), a.Float64Value())
}
