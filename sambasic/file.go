package sambasic

const (
	defaultNumericVarsSize = 92
	defaultGapSize         = 512
)

type File struct {
	Lines           []Line
	NumericVars     []byte
	Gap             []byte
	StringArrayVars []byte
	StartLine       uint16
}

type Line struct {
	Number uint16
	Tokens []Token
}

func (l *Line) Bytes() []byte {
	data := []byte{}
	for _, t := range l.Tokens {
		data = append(data, t.Bytes()...)
	}
	data = append(data, 0x0D)
	result := []byte{
		byte(l.Number >> 8),
		byte(l.Number & 0xFF),
		byte(len(data) & 0xFF),
		byte(len(data) >> 8),
	}
	return append(result, data...)
}

func (f *File) ProgBytes() []byte {
	result := []byte{}
	for _, line := range f.Lines {
		result = append(result, line.Bytes()...)
	}
	result = append(result, 0xFF)
	return result
}

func (f *File) numericVars() []byte {
	if f.NumericVars != nil {
		return f.NumericVars
	}
	return make([]byte, defaultNumericVarsSize)
}

func (f *File) gap() []byte {
	if f.Gap != nil {
		return f.Gap
	}
	return make([]byte, defaultGapSize)
}

func (f *File) stringArrayVars() []byte {
	if f.StringArrayVars != nil {
		return f.StringArrayVars
	}
	return nil
}

func (f *File) Bytes() []byte {
	result := f.ProgBytes()
	result = append(result, f.numericVars()...)
	result = append(result, f.gap()...)
	result = append(result, f.stringArrayVars()...)
	return result
}

func (f *File) NVARSOffset() uint32 {
	return uint32(len(f.ProgBytes()))
}

func (f *File) NUMENDOffset() uint32 {
	return f.NVARSOffset() + uint32(len(f.numericVars()))
}

func (f *File) SAVARSOffset() uint32 {
	return f.NUMENDOffset() + uint32(len(f.gap()))
}
