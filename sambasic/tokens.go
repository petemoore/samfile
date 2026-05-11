package sambasic

import "strconv"

type Token interface {
	Bytes() []byte
}

type SingleByteKeyword byte

func (k SingleByteKeyword) Bytes() []byte {
	return []byte{byte(k)}
}

type TwoByteKeyword byte

func (k TwoByteKeyword) Bytes() []byte {
	return []byte{0xFF, byte(k)}
}

type Num struct {
	Display string
	Value   [5]byte
}

func (n *Num) Bytes() []byte {
	result := []byte(n.Display)
	result = append(result, 0x0E, n.Value[0], n.Value[1], n.Value[2], n.Value[3], n.Value[4])
	return result
}

func Number(n uint16) *Num {
	return &Num{
		Display: strconv.Itoa(int(n)),
		Value:   [5]byte{0x00, 0x00, byte(n & 0xFF), byte(n >> 8), 0x00},
	}
}

type Str []byte

func (s *Str) Bytes() []byte {
	return []byte(*s)
}

func String(s string) *Str {
	v := Str(s)
	return &v
}

type Literal byte

func (l Literal) Bytes() []byte {
	return []byte{byte(l)}
}
