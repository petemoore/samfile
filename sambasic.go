package samfile

import (
	"fmt"
	"os"

	"github.com/petemoore/samfile/v3/sambasic"
)

type (
	// SAMBasic wraps the body bytes of a SAM BASIC (FT_SAM_BASIC,
	// type 16) file so they can be detokenised back into a plain
	// text listing. Data must be the file body without its 9-byte
	// FileHeader prefix — i.e. the Body field of a File returned
	// by DiskImage.File, or equivalent.
	//
	// The body is a sequence of tokenised lines terminated by a
	// 0xFF end-of-program sentinel. Each line has a 2-byte
	// big-endian line number, a 2-byte little-endian body length,
	// the tokenised body and a 0x0D line terminator. Keyword
	// tokens are single bytes in the range 0x85..0xF6 or the
	// two-byte sequence 0xFF, <idx>; numeric literals carry a
	// 5-byte "invisible" floating-point representation introduced
	// by 0x0E (which Output silently skips).
	SAMBasic struct {
		Data []byte
	}
)

// NewSAMBasic wraps a SAM BASIC body for detokenisation. data is
// taken by reference, not copied.
func NewSAMBasic(data []byte) *SAMBasic {
	return &SAMBasic{
		Data: data,
	}
}

// Output writes basic.Data as a plain-text BASIC listing to stdout:
// each line is prefixed with a 5-space-padded decimal line number,
// keyword tokens are expanded via the v3 SAM BASIC keyword table,
// the invisible 5-byte numeric form after each 0x0E byte is
// skipped, and 0x0D becomes a newline. Control characters below
// 0x20 (other than 0x0D and 0x0E) are rendered as "{N}". Returns
// an error if the input is empty, truncated, or contains an
// out-of-range keyword index.
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
				name, ok := sambasic.KeywordName(b, true)
				if !ok {
					return fmt.Errorf("basic-to-text: invalid keyword byte 0x%02x after 0xff escape at offset %d", b, index+uint32(c))
				}
				if !spaceBefore {
					fmt.Print(" ")
				}
				fmt.Print(name + " ")
				spaceBefore = true
			case b == 0x0e:
				c += 5
			case b == 0x0d:
				fmt.Println("")
				spaceBefore = false
			case b < 0x20:
				fmt.Printf("{%v}", int(b))
			case b >= 0x85 && b <= 0xf6:
				name, ok := sambasic.KeywordName(b, false)
				if !ok {
					return fmt.Errorf("basic-to-text: keyword index %d out of range", b-0x3b)
				}
				if !spaceBefore {
					fmt.Print(" ")
				}
				fmt.Print(name + " ")
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
