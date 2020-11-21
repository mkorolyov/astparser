package astparser

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"github.com/pkg/errors"
)

// TODO parse type comments
// Walker implements go/ast.Visitor to walk through golang
// structs and constants to parse them.
type Walker struct {
	Structs   []StructDef
	Constants []ConstantDef
	Package   string
}

// A Walkers's Visit method is invoked for each node encountered by go/ast.Walk.
// If the result visitor w is not nil, go/ast.Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
func (w *Walker) Visit(node ast.Node) ast.Visitor {
	switch spec := node.(type) {
	case *ast.TypeSpec:
		w.visitTypeSpec(spec)
		return nil
	case *ast.ValueSpec:
		w.visitConstant(spec)
	case *ast.File:
		w.Package = spec.Name.String()
	}

	return w
}

func (w *Walker) visitConstant(astValueSpec *ast.ValueSpec) {
	if len(astValueSpec.Names) < 1 || len(astValueSpec.Values) < 1 {
		return
	}
	name := astValueSpec.Names[0].Name
	var value string
	switch v := astValueSpec.Values[0].(type) {
	case *ast.BasicLit:
		value = removeQuotes(v.Value)
	default:
		return
	}

	w.Constants = append(w.Constants, ConstantDef{
		Name: name, Value: value,
	})
}

func (w *Walker) visitTypeSpec(astTypeSpec *ast.TypeSpec) {
	structName := astTypeSpec.Name.Name

	switch v := astTypeSpec.Type.(type) {
	case *ast.StructType:
		astFields := v.Fields.List

		s := StructDef{
			Name:     structName,
			Comments: parseComments(astTypeSpec.Doc)}

		for _, astField := range astFields {
			field, err := parseField(astField)
			if err != nil {
				log.Fatalf("failed to parse struct %s: %v", structName, err)
			}
			if field != nil {
				s.Fields = append(s.Fields, *field)
			}
		}

		w.Structs = append(w.Structs, s)

	case *ast.ArrayType:
		// skip array aliases for now
		return

	case *ast.Ident:
		// *ast.TypeSpec can also be a type alias
		return
	case *ast.InterfaceType:
		return
	default:
		log.Fatalf("unexpected type for typeSpec: %s, %+v: %T", structName, astTypeSpec, astTypeSpec.Type)
	}

}

func parseField(astField *ast.Field) (*FieldDef, error) {
	fieldName := parseFieldName(astField.Names)

	tag, err := parseTags(astField.Tag)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse field %q tags", fieldName)
	}

	// if field marked json:"-" we skip it.
	if tag.JsonName == "-" {
		return nil, nil
	}

	fieldType, err := parseFieldType(astField.Type)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse field %s type", fieldName)
	}

	fieldDef := &FieldDef{
		FieldName: fieldName,
		FieldType: fieldType,
		Nullable:  tag.Omitempty || tag.Nullable,
		JsonName:  tag.JsonName,
		Comments:  parseComments(astField.Doc),
		AllTags:   tag.AllTags,
	}

	if fieldDef.FieldName == "" {
		fieldDef.CompositionField = true
	}

	return fieldDef, nil
}

func parseComments(group *ast.CommentGroup) []string {
	if group == nil || len(group.List) == 0 {
		return nil
	}

	var comments []string
	for _, c := range group.List {
		if c == nil {
			continue
		}
		comments = append(comments, strings.TrimSpace(strings.Replace(c.Text, "//", "", -1)))
	}
	return comments
}

func removeQuotes(s string) string {
	if len(s) < 2 {
		//panic(fmt.Sprintf("bad input for removing quotes: %s", s))
		return s
	}

	return s[1 : len(s)-1]
}

func parseTags(astTag *ast.BasicLit) (Tag, error) {
	if astTag == nil {
		return Tag{}, nil
	}
	tagString := astTag.Value
	if tagString == "" {
		return Tag{}, nil
	}
	tagString = removeQuotes(tagString) //clean from `json:"place_type,omitempty"` to  json:"place_type,omitempty"
	splittedTags := strings.Split(tagString, " ")
	t := Tag{AllTags: make(map[string]string, len(splittedTags))}
	for _, tagWithName := range splittedTags {
		if tagWithName == "" {
			continue
		}
		v := strings.SplitN(tagWithName, ":", 2)
		if len(v) != 2 {
			return Tag{}, fmt.Errorf("invalid tag %s", tagWithName)
		}
		tagName := strings.Trim(v[0], " ")
		t.AllTags[tagName] = removeQuotes(v[1])

		switch tagName {
		case "nullable":
			if removeQuotes(v[1]) == "true" {
				t.Nullable = true
			}
		case "json":
			tagValues := strings.Split(strings.TrimSpace(v[1]), ",")
			if len(tagValues) == 1 { // e.g. json:"field_name"
				t.JsonName = removeQuotes(tagValues[0])
			} else { // e.g. handle json:"field_name,omitempty" where
				// tagValues[0] = '"field_name', so we need to cut first char "
				t.JsonName = tagValues[0][1:len(tagValues[0])]
				if strings.Index(tagValues[1], "omitempty") != -1 {
					t.Omitempty = true
				}
			}
		default:
			continue

		}
	}
	return t, nil
}

func parseFieldType(t ast.Expr) (Type, error) {
	switch v := t.(type) {
	case *ast.InterfaceType:
		return TypeInterfaceValue{}, nil
	case *ast.Ident:
		if st := simpleType(v.Name); st != nil {
			return st, nil
		}
		return TypeCustom{Name: v.Name}, nil
	case *ast.SelectorExpr:
		return TypeCustom{Name: v.Sel.Name, Expr: t}, nil
	case *ast.ArrayType:
		t, err := parseFieldType(v.Elt)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse array nested type %+v", t)
		}
		return TypeArray{InnerType: t}, nil
	case *ast.StarExpr:
		t, err := parseFieldType(v.X)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse star expr type %+v", t)
		}
		return TypePointer{InnerType: t}, nil
	case *ast.MapType:
		kt, err := parseFieldType(v.Key)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse map key type %+v", v.Key)
		}
		vt, er := parseFieldType(v.Value)
		if er != nil {
			return nil, errors.Wrapf(er, "failed to parse map value type %+v", v.Value)
		}
		return TypeMap{KeyType: kt, ValueType: vt}, nil
	default:
		return nil, fmt.Errorf("unexpected %+[1]v with type %[1]T", t)
	}
}

func parseFieldName(fieldNames []*ast.Ident) string {
	if len(fieldNames) == 0 {
		return ""
	}
	return fieldNames[0].Name
}

func simpleType(fieldType string) Type {
	switch fieldType {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "bool", "byte":
		return TypeSimple{Name: fieldType}
	default:
		return nil
	}
}
