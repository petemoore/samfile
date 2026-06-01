// Ported to Go from Dave Hooper's SAASound (https://github.com/stripwax/SAASound),
// BSD licence — see SAASound-LICENCE.txt.
//
// noise.go: one SAA1099 noise generator (LFSR + noise clock).
// Ported from src/SAANoise.{cpp,h}.

package saa

// saaNoise is one of the chip's two pseudo-random noise generators.
//
// It is clocked either from its own internal noise clock (SourceMode 0/1/2,
// derived from the chip clock) or from a connected frequency generator
// (SourceMode 3, driven via Trigger). The noise itself is an 18-bit Galois
// LFSR (feedback polynomial x^18 + x^11 + x^1), per Jepael's documentation of
// the real SAA1099P.
type saaNoise struct {
	counter      uint32
	add          uint32
	counterLow   uint32
	oversample   uint32
	counterLimit uint32 // counterLimit_low
	sync         bool   // see "SYNC" bit of register 28
	sampleRate   uint32
	sourceMode   int
	addBase      uint32 // add for 31.25 kHz noise at 44.1 kHz samplerate

	rand uint32 // pseudo-random number generator state
}

// newSaaNoise constructs a noise generator with the given LFSR seed.
func newSaaNoise(seed uint32) *saaNoise {
	n := &saaNoise{
		counter:      0,
		counterLow:   0,
		oversample:   0,
		counterLimit: 1,
		sync:         false,
		sampleRate:   sampleRateHz,
		sourceMode:   0,
		rand:         seed,
	}
	n.setClockRate(externalClkHz)
	n.add = n.addBase
	return n
}

func (n *saaNoise) setClockRate(clockRate int) {
	// at 8MHz the clock rate is 31.250kHz: clock rate divided by 256 (2^8).
	// We then shift by 2^12 (like Freq) for better period accuracy, i.e. shift
	// by (12-8).
	n.addBase = uint32(clockRate) << (12 - 8)
}

func (n *saaNoise) seed(seed uint32) {
	n.rand = seed
}

func (n *saaNoise) setSource(source int) {
	n.sourceMode = source
	n.add = n.addBase >> uint(n.sourceMode)
}

// trigger advances the LFSR; only meaningful when clocking from a frequency
// generator (SourceMode 3).
func (n *saaNoise) trigger() {
	if n.sourceMode == 3 {
		n.changeLevel()
	}
}

func (n *saaNoise) tick() {
	// Tick only does anything when clocking from the noise clock (SourceMode
	// 0/1/2). If SourceMode==3 we are clocked by a frequency generator, so do
	// nothing here.
	if !n.sync && n.sourceMode != 3 {
		n.counter += n.add
		for n.counter >= (n.sampleRate << 12) {
			n.counter -= n.sampleRate << 12
			n.counterLow++
			if n.counterLow >= n.counterLimit {
				n.counterLow = 0
				n.changeLevel()
			}
		}
	}
}

func (n *saaNoise) doSync(s bool) {
	if s {
		n.counter = 0
		n.counterLow = 0
	}
	n.sync = s
}

func (n *saaNoise) setSampleRate(sampleRate int) {
	n.sampleRate = uint32(sampleRate)
}

func (n *saaNoise) setOversample(oversample uint32) {
	if oversample < n.oversample {
		n.counterLow <<= (n.oversample - oversample)
	} else {
		n.counterLow >>= (oversample - n.oversample)
	}
	n.counterLimit = 1 << oversample
	n.oversample = oversample
}

// level returns the current LFSR output bit (0 or 1).
func (n *saaNoise) level() int {
	return int(n.rand & 0x00000001)
}

func (n *saaNoise) changeLevel() {
	// 18-bit Galois LFSR, feedback polynomial x^18 + x^11 + x^1,
	// period 2^18-1. Tap mask 0x20400.
	if n.rand&1 != 0 {
		n.rand = (n.rand >> 1) ^ 0x20400
	} else {
		n.rand >>= 1
	}
}
