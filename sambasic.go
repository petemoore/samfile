package samfile

import (
	"fmt"
	"os"
)

type (
	SAMBasic struct {
		Data []byte
	}
)

func NewSAMBasic(data []byte) *SAMBasic {
	return &SAMBasic{
		Data: data,
	}
}

func (basic *SAMBasic) Output() error {
	if len(basic.Data) == 0 {
		return fmt.Errorf("basic-to-text: empty input; expected SAM BASIC bytes on stdin")
	}
	n := uint32(len(basic.Data))
	index := uint32(0)
	for {
		if index >= n {
			return fmt.Errorf("basic-to-text: truncated input: missing 0xff end-of-program sentinel after offset %d", index)
		}
		if basic.Data[index] == 0xff {
			break
		}
		if index+3 >= n {
			return fmt.Errorf("basic-to-text: truncated input: incomplete line header at offset %d (need 4 bytes, have %d)", index, n-index)
		}
		lineNo := uint16(basic.Data[index])<<8 | uint16(basic.Data[index+1])
		lineLen := uint16(basic.Data[index+2]) | uint16(basic.Data[index+3])<<8
		index += 4
		fmt.Printf("%5d ", lineNo)
		spaceBefore := true
		for c := uint16(0); c < lineLen; c++ {
			if index+uint32(c) >= n {
				return fmt.Errorf("basic-to-text: truncated input: line body for line %d extends past input (offset %d, length %d)", lineNo, index+uint32(c), n)
			}
			b := basic.Data[index+uint32(c)]
			switch {
			case b == 0xff:
				c++
				if index+uint32(c) >= n {
					return fmt.Errorf("basic-to-text: truncated input: 0xff keyword escape at end of input (offset %d)", index+uint32(c))
				}
				b := basic.Data[index+uint32(c)]
				if b < 0x3b {
					return fmt.Errorf("basic-to-text: invalid keyword byte 0x%02x after 0xff escape at offset %d", b, index+uint32(c))
				}
				if int(b-0x3b) >= len(keywords) {
					return fmt.Errorf("basic-to-text: keyword index %d out of range (table has %d entries)", b-0x3b, len(keywords))
				}
				if !spaceBefore {
					fmt.Print(" ")
				}
				fmt.Print(keywords[b-0x3b] + " ")
				spaceBefore = true
			case b == 0x0e:
				c += 5
			case b == 0x0d:
				fmt.Println("")
				spaceBefore = false
			case b < 0x20:
				fmt.Printf("{%v}", int(b))
			case b >= 0x85 && b <= 0xf6:
				if int(b-0x3b) >= len(keywords) {
					return fmt.Errorf("basic-to-text: keyword index %d out of range (table has %d entries)", b-0x3b, len(keywords))
				}
				if !spaceBefore {
					fmt.Print(" ")
				}
				fmt.Print(keywords[b-0x3b] + " ")
				spaceBefore = true
			default:
				_, _ = os.Stdout.Write(basic.Data[index+uint32(c) : index+uint32(c)+1])
				spaceBefore = false
			}
		}
		index += uint32(lineLen)
	}
	return nil
}
