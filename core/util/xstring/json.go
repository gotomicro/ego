package xstring

import (
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/modern-go/reflect2"
)

var _jsonPrettyAPI = jsoniter.Config{
	IndentionStep:                 4,
	MarshalFloatWith6Digits:       false,
	EscapeHTML:                    true,
	SortMapKeys:                   false,
	UseNumber:                     false,
	DisallowUnknownFields:         false,
	TagKey:                        "",
	OnlyTaggedField:               false,
	ValidateJsonRawMessage:        false,
	ObjectFieldMustBeSimpleString: false,
	CaseSensitive:                 false,
}.Froze()

var _jsonAPI = jsoniter.Config{
	SortMapKeys:            true,
	UseNumber:              true,
	CaseSensitive:          true,
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

// OmitDefaultAPI ...
var OmitDefaultAPI = jsoniter.Config{
	SortMapKeys:            true,
	UseNumber:              true,
	CaseSensitive:          true,
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

func init() {
	OmitDefaultAPI.RegisterExtension(new(emitDefaultExtension))
}

// JSON ...
func JSON(obj interface{}) string {
	aa, _ := _jsonAPI.Marshal(obj)
	return string(aa)
}

// JSONBytes ...
func JSONBytes(obj interface{}) []byte {
	aa, _ := _jsonAPI.Marshal(obj)
	return aa
}

// PrettyJSON ...
func PrettyJSON(obj interface{}) string {
	aa, _ := _jsonPrettyAPI.MarshalIndent(obj, "", "    ")
	return string(aa)
}

// PrettyJSONBytes ...
func PrettyJSONBytes(obj interface{}) []byte {
	aa, _ := _jsonPrettyAPI.MarshalIndent(obj, "", "    ")
	return aa
}

type emitDefaultExtension struct {
	jsoniter.DummyExtension
}

// UpdateStructDescriptor ...
func (ed emitDefaultExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, field := range structDescriptor.Fields {
		var hasOmitEmpty bool
		tagParts := strings.Split(field.Field.Tag().Get("json"), ",")
		for _, tagPart := range tagParts[1:] {
			if tagPart == "omitempty" {
				hasOmitEmpty = true
				break
			}
		}
		if hasOmitEmpty {
			oldField := field.Field
			field.Field = &myfield{oldField}
		}
	}
}

type myfield struct{ reflect2.StructField }

// Tag 不得不用这么骚的操作
func (mf *myfield) Tag() reflect.StructTag {
	return reflect.StructTag(strings.Replace(string(mf.StructField.Tag()), ",omitempty", "", -1))
}
