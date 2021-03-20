package fixtures_test

type Dep struct {
	Int int `json:"int"`
}

type Dep2 struct {
	String string `json:"string"`
}

type StructSlice []Dep

type MyEnum2 string

const (
	MyEnum21 MyEnum2 = "1"
	MyEnum22 MyEnum2 = "2"
)

type Struct struct {
	Dep Dep `json:"dep"`
	Dep2
	Constant  MyEnum
	Constant2 MyEnum2
}
