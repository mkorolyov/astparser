package fixtures_test

type Dep struct {
	Int int `json:"int"`
}

type StructSlice []Dep

type Struct struct {
	Dep Dep `json:"dep"`
}
