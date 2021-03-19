package fixtures_test

type Dep struct {
	Int int `json:"int"`
}

type Dep2 struct {
	String string `json:"string"`
}

type StructSlice []Dep

type Struct struct {
	Dep Dep `json:"dep"`
	Dep2
	Constant MyEnum
}
