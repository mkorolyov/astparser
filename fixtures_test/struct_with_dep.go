package fixtures_test

type Dep struct {
	Int int `json:"int"`
}

type Struct struct {
	Dep Dep `json:"dep"`
}
