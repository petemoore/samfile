// Ported to Go from Dave Hooper's SAASound (https://github.com/stripwax/SAASound),
// BSD licence — see SAASound-LICENCE.txt.
//
// freq.go: one SAA1099 tone-generator frequency divider (note→freq+octave).
// Ported from src/SAAFreq.{cpp,h}. The SAASound build used by SimCoupe
// computes the frequency lookup table at runtime (SAAFREQ_FIXED_CLOCKRATE is
// undefined), so we do the same and do not need SAAFreq.dat.

package saa

const initialLevel = 1

// connected-generator modes for an oscillator.
const (
	connectedNone  = 0 // nothing
	connectedEnv   = 1 // env generator
	connectedNoise = 2 // noise generator
)

// freqTable is the shared frequency lookup table, indexed by octave<<8|offset.
// In the C++ this is a static member shared across all CSAAFreq instances,
// recomputed only when the clock rate changes; freqTableClock guards that.
//
// Each entry is (15625 << octave) / (511 - offset), the standard SAA formula
// for an 8 MHz base clock, multiplied by 8192 (12 fractional bits, doubled so
// the counter toggles on half-waves), then rescaled by clockRate/8000000.
var (
	freqTable      [2048]uint32
	freqTableClock int
)

func setFreqClockRate(clockRate int) {
	if clockRate != freqTableClock {
		freqTableClock = clockRate
		ix := 0
		for octave := 0; octave < 8; octave++ {
			for offset := 0; offset < 256; offset++ {
				freqTable[ix] = uint32((8192.0 * 15625.0 * float64(int(1)<<octave) *
					(float64(clockRate) / 8000000.0)) / (511.0 - float64(offset)))
				ix++
			}
		}
	}
}

// saaFreq is one of the chip's six tone generators.
type saaFreq struct {
	counter      uint32
	add          uint32
	counterLow   uint32
	oversample   uint32
	counterLimit uint32 // counterLimit_low
	level        int

	currentOffset int
	currentOctave int
	nextOffset    int
	nextOctave    int

	ignoreOffsetData bool
	newData          bool
	sync             bool

	sampleRate uint32

	connectedNoise *saaNoise
	connectedEnv   *saaEnv
	connectedMode  int
}

// newSaaFreq constructs a tone generator optionally connected to a noise or
// envelope generator. At most one of noiseGen/envGen is non-nil, matching the
// device wiring.
func newSaaFreq(noiseGen *saaNoise, envGen *saaEnv) *saaFreq {
	mode := connectedNone
	if noiseGen != nil {
		mode = connectedNoise
	} else if envGen != nil {
		mode = connectedEnv
	}
	f := &saaFreq{
		counter:        0,
		add:            0,
		counterLow:     0,
		oversample:     0,
		counterLimit:   1,
		level:          initialLevel,
		currentOffset:  0,
		currentOctave:  0,
		nextOffset:     0,
		nextOctave:     0,
		sampleRate:     sampleRateHz,
		connectedNoise: noiseGen,
		connectedEnv:   envGen,
		connectedMode:  mode,
	}
	f.setClockRate(externalClkHz)
	f.setAdd()
	return f
}

func (f *saaFreq) setFreqOffset(offset byte) {
	if !f.sync {
		f.nextOffset = int(offset)
		f.newData = true
		if f.nextOctave == f.currentOctave {
			// Per Philips: if new octave then new offset are sent in that
			// order, on the next half-cycle ONLY the octave data is acted
			// upon; the offset data is acted upon next time.
			f.ignoreOffsetData = true
		}
	} else {
		// updates straightaway if sync
		f.newData = false
		f.ignoreOffsetData = false
		f.currentOffset = int(offset)
		f.nextOffset = int(offset)
		f.currentOctave = f.nextOctave
		f.setAdd()
	}
}

func (f *saaFreq) setFreqOctave(octave byte) {
	if !f.sync {
		f.nextOctave = int(octave)
		f.newData = true
		f.ignoreOffsetData = false
	} else {
		f.newData = false
		f.ignoreOffsetData = false
		f.currentOctave = int(octave)
		f.nextOctave = int(octave)
		f.currentOffset = f.nextOffset
		f.setAdd()
	}
}

func (f *saaFreq) updateOctaveOffsetData() {
	if !f.newData {
		// optimise for the most common case! No new data!
		return
	}
	f.currentOctave = f.nextOctave
	if !f.ignoreOffsetData {
		f.currentOffset = f.nextOffset
		f.newData = false
	}
	f.ignoreOffsetData = false
	f.setAdd()
}

func (f *saaFreq) setSampleRate(sampleRate int) {
	f.sampleRate = uint32(sampleRate)
}

func (f *saaFreq) setOversample(oversample uint32) {
	if oversample < f.oversample {
		f.counterLow <<= (f.oversample - oversample)
	} else {
		f.counterLow >>= (oversample - f.oversample)
	}
	f.counterLimit = 1 << oversample
	f.oversample = oversample
}

func (f *saaFreq) setClockRate(clockRate int) {
	setFreqClockRate(clockRate)
}

func (f *saaFreq) tick() int {
	if f.sync {
		return 1
	}

	f.counter += f.add
	for f.counter >= (f.sampleRate << 12) {
		f.counter -= f.sampleRate << 12
		f.counterLow++
		if f.counterLow >= f.counterLimit {
			// period elapsed for (at least) one half-cycle
			f.counterLow = 0
			// flip state
			f.level = 1 - f.level

			// trigger any connected devices
			switch f.connectedMode {
			case connectedEnv:
				f.connectedEnv.internalClock()
			case connectedNoise:
				f.connectedNoise.trigger()
			}

			// get new frequency if new data is waiting
			f.updateOctaveOffsetData()
		}
	}

	return f.level
}

// levelNow mirrors CSAAFreq::Level (the const accessor).
func (f *saaFreq) levelNow() int {
	if f.sync {
		return 1
	}
	return f.level
}

func (f *saaFreq) setAdd() {
	f.add = freqTable[f.currentOctave<<8|f.currentOffset]
}

func (f *saaFreq) doSync(s bool) {
	f.sync = s
	if f.sync {
		f.counter = 0
		f.counterLow = 0
		f.level = initialLevel
		f.currentOctave = f.nextOctave
		f.currentOffset = f.nextOffset
		f.setAdd()
	}
}
