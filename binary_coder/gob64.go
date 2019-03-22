package binary_coder

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"gitlab.com/artilligence/http-db-saver/domain"
)

type Gob64Coder struct {}

func NewGob64Coder() domain.BinaryCoder {
	gob.Register(&domain.Entity{})

	return &Gob64Coder{}
}

func (g *Gob64Coder) Encode(m interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		return nil, fmt.Errorf("failed gob Encode: %s", err)
	}

	return b.Bytes(), nil
}

func (g *Gob64Coder) Decode(str string) (interface{}, error) {
	var m interface{}

	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("failed base64 Decode: %s", err)
	}

	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)

	err = d.Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("failed gob Decode: %s", err)
	}

	return m, nil
}