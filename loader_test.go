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
			want: ParsedFile{Structs: []StructDef{
				{
					Name: "Primitives",
					Fields: []FieldDef{
						{
							FieldName: "Int",
							JsonName:  "int",
							Comments:  []string{"comment here"},
							FieldType: TypeSimple{Name: "int"},
						},
						{
							FieldName: "Int64",
							JsonName:  "int_64",
							FieldType: TypeSimple{Name: "int64"},
						},
						{
							FieldName: "Float32",
							JsonName:  "float_32",
							FieldType: TypeSimple{Name: "float32"},
						},
						{
							FieldName: "Float64",
							JsonName:  "float_64",
							FieldType: TypeSimple{Name: "float64"},
						},
						{
							FieldName: "Bool",
							JsonName:  "bool",
							FieldType: TypeSimple{Name: "bool"},
						},
						{
							FieldName: "String",
							JsonName:  "string",
							FieldType: TypeSimple{Name: "string"},
						},
						{
							FieldName: "Bytes",
							JsonName:  "bytes",
							FieldType: TypeArray{
								InnerType: TypeSimple{Name: "byte"}},
						},
						{
							FieldName: "Map",
							JsonName:  "map",
							FieldType: TypeMap{
								KeyType:   TypeSimple{Name: "string"},
								ValueType: TypeSimple{Name: "string"}},
						},
						{
							FieldName: "Slice",
							JsonName:  "slice",
							FieldType: TypeArray{InnerType: TypeSimple{Name: "int"}},
						},
						{
							FieldName: "Omitempty",
							JsonName:  "omitempty",
							FieldType: TypeSimple{Name: "int"},
							Nullable:  true,
						},
						{
							FieldName: "Required",
							JsonName:  "some_int",
							FieldType: TypeSimple{Name: "int"},
						},
						{
							FieldName: "Ptr",
							JsonName:  "ptr",
							FieldType: TypePointer{
								InnerType: TypeSimple{Name: "int"}},
						},
						{
							FieldName: "NullableBool",
							JsonName:  "nullable_bool",
							FieldType: TypeSimple{Name: "bool"},
							Nullable:  true,
						},
						{
							FieldName: "NullableBoolOmitempty",
							JsonName:  "nullable_bool_omitempty",
							FieldType: TypeSimple{Name: "bool"},
							Nullable:  true,
						},
					},
				},
			},
			},
		},
		{
			name:     "struct with dep",
			filename: "fixtures_test/struct_with_dep.go",
			want: ParsedFile{Structs: []StructDef{
				{
					Name: "Dep",
					Fields: []FieldDef{
						{
							FieldType: TypeSimple{Name: "int"},
							JsonName:  "int",
							FieldName: "Int",
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
						},
					},
				},
			},
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
			}},
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
				t.Errorf("parseFile() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
