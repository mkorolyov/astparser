package astparser

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_validFile(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		include *regexp.Regexp
		exclude *regexp.Regexp
		want    bool
	}{
		{
			name:    "include ok",
			s:       "event.go",
			include: regexp.MustCompile("event"),
			want:    true,
		},
		{
			name:    "include dont match",
			s:       "event.go",
			include: regexp.MustCompile("type"),
			want:    false,
		},
		{
			name: "valid go file",
			s:    "event.go",
			want: true,
		},
		{
			name: "test file",
			s:    "event_test.go",
			want: false,
		},
		{
			name:    "exclude ok",
			s:       "event.go",
			exclude: regexp.MustCompile("event"),
			want:    false,
		},
		{
			name:    "exclude dont match",
			s:       "event.go",
			exclude: regexp.MustCompile("type"),
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validFile(tt.s, tt.include, tt.exclude); got != tt.want {
				t.Errorf("validFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     ParsedFile
		wantErr  bool
	}{
		{
			name:     "struct with primitives",
			filename: "fixtures_test/struct_with_primitives.go",
			want: ParsedFile{
				Structs: []StructDef{
					{
						Name: "Primitives",
						Fields: []FieldDef{
							{
								FieldName: "Int",
								JsonName:  "int",
								Comments:  []string{"comment here"},
								FieldType: TypeSimple{Name: "int"},
								AllTags:   map[string]string{"json": "int"},
							},
							{
								FieldName: "Int64",
								JsonName:  "int_64",
								FieldType: TypeSimple{Name: "int64"},
								AllTags:   map[string]string{"json": "int_64"},
							},
							{
								FieldName: "Float32",
								JsonName:  "float_32",
								FieldType: TypeSimple{Name: "float32"},
								AllTags:   map[string]string{"json": "float_32"},
							},
							{
								FieldName: "Float64",
								JsonName:  "float_64",
								FieldType: TypeSimple{Name: "float64"},
								AllTags:   map[string]string{"json": "float_64"},
							},
							{
								FieldName: "Bool",
								JsonName:  "bool",
								FieldType: TypeSimple{Name: "bool"},
								AllTags:   map[string]string{"json": "bool"},
							},
							{
								FieldName: "String",
								JsonName:  "string",
								FieldType: TypeSimple{Name: "string"},
								AllTags:   map[string]string{"json": "string"},
							},
							{
								FieldName: "Bytes",
								JsonName:  "bytes",
								FieldType: TypeArray{
									InnerType: TypeSimple{Name: "byte"}},
								AllTags: map[string]string{"json": "bytes"},
							},
							{
								FieldName: "Map",
								JsonName:  "map",
								FieldType: TypeMap{
									KeyType:   TypeSimple{Name: "string"},
									ValueType: TypeSimple{Name: "string"}},
								AllTags: map[string]string{"json": "map"},
							},
							{
								FieldName: "MapInterface",
								JsonName:  "map_interface",
								FieldType: TypeMap{
									KeyType:   TypeSimple{Name: "string"},
									ValueType: TypeInterfaceValue{}},
								AllTags: map[string]string{"json": "map_interface"},
							},
							{
								FieldName: "Slice",
								JsonName:  "slice",
								FieldType: TypeArray{InnerType: TypeSimple{Name: "int"}},
								AllTags:   map[string]string{"json": "slice"},
							},
							{
								FieldName: "Omitempty",
								JsonName:  "omitempty",
								FieldType: TypeSimple{Name: "int"},
								Nullable:  true,
								AllTags:   map[string]string{"json": "omitempty,omitempty"},
							},
							{
								FieldName: "Required",
								JsonName:  "some_int",
								FieldType: TypeSimple{Name: "int"},
								AllTags:   map[string]string{"json": "some_int,required"},
							},
							{
								FieldName: "Ptr",
								JsonName:  "ptr",
								FieldType: TypePointer{
									InnerType: TypeSimple{Name: "int"}},
								AllTags: map[string]string{"json": "ptr"},
							},
							{
								FieldName: "NullableBool",
								JsonName:  "nullable_bool",
								FieldType: TypeSimple{Name: "bool"},
								Nullable:  true,
								AllTags:   map[string]string{"json": "nullable_bool", "nullable": "true"},
							},
							{
								FieldName: "NullableBoolOmitempty",
								JsonName:  "nullable_bool_omitempty",
								FieldType: TypeSimple{Name: "bool"},
								Nullable:  true,
								AllTags:   map[string]string{"json": "nullable_bool_omitempty,omitempty", "nullable": "true"},
							},
						},
					},
				},
				Package: "fixtures_test",
			},
		},
		{
			name:     "struct with dep",
			filename: "fixtures_test/struct_with_dep.go",
			want: ParsedFile{
				Structs: []StructDef{
					{
						Name: "Dep",
						Fields: []FieldDef{
							{
								FieldType: TypeSimple{Name: "int"},
								JsonName:  "int",
								FieldName: "Int",
								AllTags:   map[string]string{"json": "int"},
							},
						},
					},
					{
						Name: "Dep2",
						Fields: []FieldDef{
							{
								FieldType: TypeSimple{Name: "string"},
								JsonName:  "string",
								FieldName: "String",
								AllTags:   map[string]string{"json": "string"},
							},
						},
					},
					{
						Name: "Struct",
						Fields: []FieldDef{
							{
								FieldType: TypeCustom{Name: "Dep"},
								JsonName:  "dep",
								FieldName: "Dep",
								AllTags:   map[string]string{"json": "dep"},
							},
							{
								FieldType:        TypeCustom{Name: "Dep2"},
								CompositionField: true,
							},
							{
								CompositionField: false,
								FieldName:        "Constant",
								FieldType:        TypeCustom{Name: "MyEnum", Alias: true},
								Nullable:         false,
							},
							{
								CompositionField: false,
								FieldName:        "Constant2",
								FieldType:        TypeCustom{Name: "MyEnum2", AliasType: TypeSimple{Name: "string"}},
								Nullable:         false,
							},
						},
					},
				},
				Constants: []ConstantDef{
					{
						Name:  "MyEnum21",
						Value: "1",
					},
					{
						Name:  "MyEnum22",
						Value: "2",
					},
				},
				Package: "fixtures_test",
			},
		},
		{
			name:     "constants",
			filename: "fixtures_test/constants.go",
			want: ParsedFile{Constants: []ConstantDef{
				{
					Name:  "PublicConst",
					Value: "public",
				},
				{
					Name:  "privateConst",
					Value: "private",
				},
				{
					Name:  "MyEnumValue1",
					Value: "enum-1",
				},
				{
					Name:  "MyEnumValue2",
					Value: "enum-2",
				},
			},
				Package: "fixtures_test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nhave %+v, \nwant %+v", got, tt.want)
			}
		})
	}
}
