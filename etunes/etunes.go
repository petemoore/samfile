package etunes

import (
	"bytes"
	"fmt"

	"github.com/petemoore/samfile/v3/saa"
)

const (
	sampleRate      = 44100
	frameHz         = 50
	samplesPerFrame = sampleRate / frameHz // 882

	initEntry    = 0x8000
	playEntry    = 0x8006
	shadowBase   = 0x83D3 // SAA regs 0x00..0x19 (26 bytes)
	shadowLen    = 26
	orderPtrAddr = 0x8462 // self-modified operand: live order-list pointer
	songDataOff  = 0x4B3  // song data begins at 0x8000+0x4B3
)

// etrackerEntry is the engine's fixed entry prologue: LD HL,0x84B3 ; JP 0x83EF.
var etrackerEntry = []byte{0x21, 0xB3, 0x84, 0xC3, 0xEF, 0x83}

// IsETrackerModule reports whether body looks like a SAM E-Tracker module: the
// fixed replay-engine entry prologue, enough length to hold song data, and the
// embedded "ETracker" signature.
func IsETrackerModule(body []byte) bool {
	return len(body) > songDataOff &&
		bytes.HasPrefix(body, etrackerEntry) &&
		bytes.Contains(body, []byte("ETracker"))
}

// Frames runs the module's replay engine and returns the per-frame SAA register
// shadow (regs 0x00..0x19) for one full pass (intro + one loop body), plus the
// intro and loop lengths in frames. introFrames is 0 when the tune loops from
// the start.
func Frames(body []byte) (shadows [][shadowLen]byte, introFrames, loopFrames int, err error) {
	if !IsETrackerModule(body) {
		return nil, 0, 0, fmt.Errorf("not an E-Tracker module")
	}
	cpu, s := newCPU(body)
	call(cpu, s, initEntry)

	const maxFrames = frameHz * 600
	var wraps []int
	prev := s.Get16(orderPtrAddr)
	for f := 0; f < maxFrames && len(wraps) < 2; f++ {
		call(cpu, s, playEntry)
		var sh [shadowLen]byte
		for i := 0; i < shadowLen; i++ {
			sh[i] = s.Get(uint16(shadowBase + i))
		}
		order := s.Get16(orderPtrAddr)
		if f > 0 && order < prev {
			wraps = append(wraps, f)
		}
		prev = order
		shadows = append(shadows, sh)
	}
	if len(wraps) >= 2 {
		w1, w2 := wraps[0], wraps[1]
		loopFrames = w2 - w1
		introFrames = w1 - loopFrames
		if introFrames < 0 {
			introFrames = 0
		}
		shadows = shadows[:w1]
	} else {
		loopFrames = len(shadows) // no wrap found: treat all as the loop
	}
	return shadows, introFrames, loopFrames, nil
}

// Render decodes a module to 44.1 kHz interleaved 16-bit stereo PCM: the intro
// once, then the loop body repeated loops times. loops must be >= 1.
func Render(body []byte, loops int) (pcm []int16, sr int, err error) {
	if loops < 1 {
		return nil, 0, fmt.Errorf("loops must be >= 1, got %d", loops)
	}
	shadows, introFrames, loopFrames, err := Frames(body)
	if err != nil {
		return nil, 0, err
	}
	chip := saa.New(0, 0)             // SAM defaults: 8 MHz / 44100 Hz
	chip.WriteAddressData(0x1C, 0x01) // sound enable
	buf := make([]int16, samplesPerFrame*2)
	emit := func(lo, hi int) {
		for f := lo; f < hi; f++ {
			for r := 0; r < shadowLen; r++ {
				chip.WriteAddressData(byte(r), shadows[f][r])
			}
			chip.GenerateMany(buf, samplesPerFrame)
			pcm = append(pcm, buf...)
		}
	}
	emit(0, introFrames)
	for rep := 0; rep < loops; rep++ {
		emit(introFrames, introFrames+loopFrames)
	}
	return pcm, sampleRate, nil
}
