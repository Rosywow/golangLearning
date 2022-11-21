package scalar

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx/types"
	"io"
	"strconv"
	"strings"
)

type Int64 int64

//UnmarshalGQL implements the graphql.UnMarshaller interface
func (i *Int64) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("value must be integer in string")
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*i = Int64(n)
	return nil
}

//MarshalGQL implements the graphql.Marshaller interface
func (i *Int64) MarshalGQL(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("%d", *i)))
	return
}

type Binary []byte

//UnmarshalGQL implements the graphql.UnMarshaller interface
func (b *Binary) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("binary must be base64 string")
	}
	var buf []byte
	buf, err := base64.StdEncoding.DecodeString(s)
	*b = buf
	return err
}

//MarshalGQL implements the graphql.Marshaller interface
func (b Binary) MarshalGQL(w io.Writer) {
	buf, err := json.Marshal(b)
	if err != nil {
		return
	}

	w.Write([]byte(strings.ReplaceAll(string(buf), "\"", "")))
	return
}

// https://github.com/99designs/gqlgen/issues/597

type Raw types.JSONText

//MarshalGQL implements the graphql.Marshaller interface
func (b Raw) MarshalGQL(w io.Writer) {
	if b == nil {
		return
	}

	w.Write(b)
	return
}

//UnmarshalGQL implements the graphql.UnMarshaller interface
func (b *Raw) UnmarshalGQL(v interface{}) (err error) {
	buf, ok := v.([]byte)
	if ok {
		*b = buf
		return
	}

	s, ok := v.(string)
	if !ok {
		return errors.New("binary must be string")
	}

	*b = []byte(s)
	return
}
