package xstring

import (
	"fmt"
	"testing"
)

func TestJSON_OmitDefault(t *testing.T) {
	type CC struct {
		D string `json:",omitempty" toml:",omitempty" `
	}
	type AA struct {
		A string `json:"a,omitempty"`
		B int    `json:",omitempty"`
		C CC     `json:",omitempty"`
		D *CC    `json:",omitempty" toml:",omitempty" `
		E *CC    `json:"e" toml:"e" `
	}

	aa := AA{A: "11", D: &CC{}}

	bs, err := OmitDefaultAPI.Marshal(aa)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("string(bs) = %+v\n", string(bs))

	as, err := _jsonAPI.Marshal(aa)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("as = %+v\n", string(as))
}
