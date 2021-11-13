package gateway

import "errors"

type RawCodec struct{}

func (r RawCodec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*[]byte)
	if !ok {
		return nil, errors.New("no bytes")
	}
	return *out, nil
}

func (r RawCodec) Unmarshal(data []byte, v interface{}) error {
	dest, ok := v.(*[]byte)
	if !ok {
		return errors.New("invalid cast type")
	}
	for _, datum := range data {
		*dest = append(*dest, datum)
	}
	return nil
}

func (r RawCodec) Name() string {
	return "myCodec"
}
