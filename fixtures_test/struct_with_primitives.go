package fixtures_test

type Primitives struct {
	// comment here
	Int       int               `json:"int"`
	Int64     int64             `json:"int_64"`
	Float32   float32           `json:"float_32"`
	Float64   float64           `json:"float_64"`
	Bool      bool              `json:"bool"`
	String    string            `json:"string"`
	Bytes     []byte            `json:"bytes"`
	Map       map[string]string `json:"map"`
	Slice     []int             `json:"slice"`
	Omitempty int               `json:"omitempty,omitempty"`
	Ptr       *int              `json:"ptr"`
}
