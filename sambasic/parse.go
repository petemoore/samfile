package sambasic

import "fmt"

func Parse(body []byte) (*File, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("parse: empty body")
	}

	f := &File{}
	pos := 0
	n := len(body)

	for pos < n {
		if body[pos] == 0xFF {
			pos++
			break
		}
		if pos+3 >= n {
			return nil, fmt.Errorf("parse: truncated line header at offset %d", pos)
		}
		lineNum := uint16(body[pos])<<8 | uint16(body[pos+1])
		lineLen := int(body[pos+2]) | int(body[pos+3])<<8
		pos += 4

		if pos+lineLen > n {
			return nil, fmt.Errorf("parse: line %d body extends past input", lineNum)
		}

		line := Line{Number: lineNum}
		end := pos + lineLen
		i := pos
		for i < end {
			b := body[i]
			switch {
			case b == 0x0D && i == end-1:
				// Line terminator — Line.Bytes() adds it back
				i++
			case b == 0xFF && i+1 < end:
				line.Tokens = append(line.Tokens, TwoByteKeyword(body[i+1]))
				i += 2
			case b >= 0x85 && b <= 0xF6:
				line.Tokens = append(line.Tokens, SingleByteKeyword(b))
				i++
			case b == 0x0E:
				if i+6 > end {
					return nil, fmt.Errorf("parse: truncated numeric form at offset %d", i)
				}
				display := []byte{}
				for len(line.Tokens) > 0 {
					last, ok := line.Tokens[len(line.Tokens)-1].(literal)
					if !ok {
						break
					}
					if last >= '0' && last <= '9' || last == '.' || last == '-' || last == 'E' || last == 'e' {
						display = append([]byte{byte(last)}, display...)
						line.Tokens = line.Tokens[:len(line.Tokens)-1]
					} else {
						break
					}
				}
				num := &Num{
					Display: string(display),
				}
				copy(num.Value[:], body[i+1:i+6])
				line.Tokens = append(line.Tokens, num)
				i += 6
			default:
				line.Tokens = append(line.Tokens, literal(b))
				i++
			}
		}
		f.Lines = append(f.Lines, line)
		pos = end
	}

	if pos < n {
		trailer := make([]byte, n-pos)
		copy(trailer, body[pos:])
		total := len(trailer)
		if total == 604 {
			f.NumericVars = trailer[:92]
			f.Gap = trailer[92:]
		} else if total > 0 {
			f.NumericVars = trailer
			f.Gap = []byte{}
		}
	} else {
		f.NumericVars = []byte{}
		f.Gap = []byte{}
	}

	return f, nil
}
