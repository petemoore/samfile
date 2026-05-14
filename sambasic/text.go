package sambasic

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

// ParseError describes a failure to tokenise SAM BASIC text input.
type ParseError struct {
	Line int    // 1-based source line
	Col  int    // 1-based source column
	Msg  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, col %d: %s", e.Line, e.Col, e.Msg)
}

// ParseText reads SAM BASIC source from r and returns a *File whose
// program bytes match what the SAM BASIC editor would store after the
// equivalent typing session, including line-edit semantics (sort by
// line number ascending, last-write-wins, bare-line-number delete).
func ParseText(r io.Reader) (*File, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseTextString(string(src))
}

// ParseTextString is the string-input variant of ParseText.
func ParseTextString(src string) (*File, error) {
	l := lex(src)
	l.state = lexStart

	edits := map[uint16]Line{}

	var curNum uint16
	var bodyBytes []byte
	bodyEmpty := true
	haveLineNumber := false

	for {
		it := l.nextItem()
		switch it.typ {
		case itemError:
			return nil, &ParseError{Line: it.line, Col: it.col, Msg: it.val}
		case itemEOF:
			if haveLineNumber {
				if bodyEmpty {
					delete(edits, curNum)
				} else {
					edits[curNum] = finalise(Line{Number: curNum, Tokens: tokenizeBytes(bodyBytes)})
				}
			}
			return assembleFile(edits), nil
		case itemLineNumber:
			if haveLineNumber {
				if bodyEmpty {
					delete(edits, curNum)
				} else {
					edits[curNum] = finalise(Line{Number: curNum, Tokens: tokenizeBytes(bodyBytes)})
				}
			}
			n, err := strconv.ParseUint(it.val, 10, 16)
			if err != nil {
				return nil, &ParseError{Line: it.line, Col: it.col, Msg: fmt.Sprintf("invalid line number: %s", it.val)}
			}
			curNum = uint16(n)
			bodyBytes = nil
			bodyEmpty = true
			haveLineNumber = true
		case itemEOL:
			if haveLineNumber {
				if bodyEmpty {
					delete(edits, curNum)
				} else {
					edits[curNum] = finalise(Line{Number: curNum, Tokens: tokenizeBytes(bodyBytes)})
				}
				haveLineNumber = false
			}
			bodyBytes = nil
			bodyEmpty = true
		default:
			if it.bytes != nil {
				bodyBytes = append(bodyBytes, it.bytes...)
			} else {
				bodyBytes = append(bodyBytes, []byte(it.val)...)
			}
			bodyEmpty = false
		}
	}
}

// tokenizeBytes converts a raw byte slice into a sequence of Token
// values suitable for sambasic.File.Bytes(). For v1 we use a single
// `literal` token per byte; output bytes are unchanged either way.
func tokenizeBytes(bytes []byte) []Token {
	tokens := make([]Token, 0, len(bytes))
	for _, b := range bytes {
		tokens = append(tokens, literal(b))
	}
	return tokens
}

func assembleFile(edits map[uint16]Line) *File {
	nums := make([]uint16, 0, len(edits))
	for n := range edits {
		nums = append(nums, n)
	}
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	lines := make([]Line, len(nums))
	for i, n := range nums {
		lines[i] = edits[n]
	}
	return &File{Lines: lines}
}

// finalise is a stub here — Task 16 implements the IF/ELSE/INK byte patches.
func finalise(line Line) Line {
	return line
}
