// THIS FILE IS AUTOMATICALLY GENERATED; DO NOT EDIT

package my_test

import (
	encoding "encoding"

	uuid "github.com/google/uuid"
)

type ID struct {
	Value uuid.UUID `json:"value"`
}

func (i ID)  MarshalText() ([]byte, error)   { return i.Value.MarshalText() }
func (i *ID) UnmarshalText(raw []byte) error { return i.Value.UnmarshalText(raw) }

var (
	_ encoding.TextMarshaler   = ID{}
	_ encoding.TextUnmarshaler = (*ID)(nil)
)

