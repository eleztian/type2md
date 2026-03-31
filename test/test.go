package test

import "github.com/eleztian/type2md/test/ext"

//go:generate type2md -f ../docs/doc_config.md github.com/eleztian/type2md/test Config

// Config doc.
type Config struct {
	Pre     Pre
	Post    *ext.Post
	Servers map[string]struct {
		Host string `json:"host,omitempty"`
		Port int    `json:"port" enums:"22,65522" require:"false"`
	} `json:"servers"` // server list
	InlineStruct `json:",inline"` // inline struct
	Slice        []string         // sss
	MapData      map[string]map[int]*OtherStruct
	ArrayData    [][2]string   `json:"array_data"`
	C            []interface{} // slice interface{}
}

// InlineStruct inline struct.
type InlineStruct struct {
	A string `json:"a"` // inline struct field a
}

// OtherStruct other struct
// this is use for test.
type OtherStruct struct {
	A string                 `json:"a" require:"true" default:"default value"` // +optional
	B [][2]ext.Mode          `json:"b"`                                        // array string
	C map[string]interface{} `json:"c"`                                        // map[string]interface{}
	D *OtherStruct           `json:"d"`                                        // nested struct
}
