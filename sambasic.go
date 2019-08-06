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

func (basic *SAMBasic) Output() {
	index := uint32(0)
	for {
		if basic.Data[index] == 0xff {
			break
		}
		lineNo := uint16(basic.Data[index])<<8 | uint16(basic.Data[index+1])
		lineLen := uint16(basic.Data[index+2]) | uint16(basic.Data[index+3])<<8
		index += 4
		fmt.Printf("%5d ", lineNo)
		spaceBefore := true
		for c := uint16(0); c < lineLen; c++ {
			b := basic.Data[index+uint32(c)]
			switch {
			case b == 0xff:
				c++
				b := basic.Data[index+uint32(c)]
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
}
