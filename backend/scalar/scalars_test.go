package scalar

import (
	"bytes"
	"testing"
)

func TestInt64(t *testing.T) {
	//64bit max value
	//         0x7FFF-FFFF FFFF-FFFF
	//           9223372036854775807
	strInt64 := "1234567890121234567"
	int64Dst := int64(0x112210F47DC79887)
	var i Int64
	i.UnmarshalGQL(strInt64)
	if int64(i) != int64Dst {
		t.Errorf("%d should be %d", i, int64Dst)
	}

	var buf bytes.Buffer
	i.MarshalGQL(&buf)
	d := buf.String()
	if d != strInt64 {
		t.Errorf("%s should be %s", d, strInt64)
	}
}

func TestBinary(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x99, 0x9a, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
	}

	s := "AAECAwQFBgcICQoLDA0ODxAiM0RVZneImZqqu8zd7v8="
	b := Binary(data)

	var buf bytes.Buffer
	b.MarshalGQL(&buf)
	d := buf.String()
	if d != s {
		t.Errorf("%s should be %s", d, s)
		return
	}

	b.UnmarshalGQL(s)
	for k, v := range data {
		if b[k] != v {
			t.Errorf("%d should be %d", b[k], v)
			return
		}
	}
}
