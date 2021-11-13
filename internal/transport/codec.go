package transport

import "errors"

type rawCodec struct{}

func (r rawCodec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*[]byte)
	if !ok {
		return nil, errors.New("no bytes")
	}
	return *out, nil
}

func (r rawCodec) Unmarshal(data []byte, v interface{}) error {
	dest, ok := v.(*[]byte)
	if !ok {
		return errors.New("invalid cast type")
	}
	for _, datum := range data {
		*dest = append(*dest, datum)
	}
	return nil
}

func (r rawCodec) Name() string {
	return "myCodec"
}
