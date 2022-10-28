package main

type StructInfo struct {
	Name     string      `json:"name"`
	Describe string      `json:"describe"`
	Fields   []FieldInfo `json:"fields"`
	Enums    *EnumInfo   `json:"enums"`
}

type FieldInfo struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Default   string   `json:"default"`
	Require   bool     `json:"require"`
	Enums     EnumInfo `json:"enums"`
	Describe  string   `json:"describe"`
	Reference string   `json:"reference"`
	skipNum   int
}

func (fi FieldInfo) Copy() FieldInfo {
	n := fi
	enums := make([][2]string, 0, len(fi.Enums.Names))
	for _, d := range enums {
		enums = append(enums, d)
	}
	n.Enums.Names = enums
	return n
}

type EnumInfo struct {
	Type  string
	Names [][2]string // key: desc
}
