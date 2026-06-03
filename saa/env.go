// Ported to Go from Dave Hooper's SAASound (https://github.com/stripwax/SAASound),
// BSD licence — see SAASound-LICENCE.txt.
//
// env.go: one of the SAA1099's two envelope generators (envelope shapes /
// resolution). Ported from src/SAAEnv.{cpp,h} and the ENVDATA tables therein.

package saa

// envData is the per-shape envelope description.
// nLevels is indexed [resolution][phase][withinphase], where resolution index
// 0 == 4-bit resolution and 1 == 3-bit resolution (matching the C++ literal
// layout; SetLevels selects index m_nResolution-1).
type envData struct {
	numberOfPhases int
	looping        bool
	levels         [2][2][16]int
}

// csEnvData are the 8 envelope shapes, ported verbatim from CSAAEnv::cs_EnvData.
var csEnvData = [8]envData{
	{1, false, [2][2][16]int{
		{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{1, true, [2][2][16]int{
		{{15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}, {15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}},
		{{14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14}, {14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14}}}},
	{1, false, [2][2][16]int{
		{{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{{14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{1, true, [2][2][16]int{
		{{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{{14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{2, false, [2][2][16]int{
		{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
		{{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}}}},
	{2, true, [2][2][16]int{
		{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
		{{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}}}},
	{1, false, [2][2][16]int{
		{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{1, true, [2][2][16]int{
		{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
}

// saaEnv is one of the chip's two envelope generators.
type saaEnv struct {
	leftLevel  int
	rightLevel int
	envDataPtr *envData

	enabled            bool
	invertRightChannel bool
	phase              byte
	phasePosition      byte
	envelopeEnded      bool
	looping            bool
	numberOfPhases     int
	resolution         int // 1 == 4-bit; 2 == 3-bit
	newData            bool
	nextData           byte
	clockExternally    bool
}

func newSaaEnv() *saaEnv {
	e := &saaEnv{
		enabled:       false,
		phase:         0,
		phasePosition: 0,
		envelopeEnded: true,
		resolution:    1,
		newData:       false,
		nextData:      0,
	}
	e.setEnvControl(0)
	return e
}

// internalClock advances the envelope when configured for internal clocking.
func (e *saaEnv) internalClock() {
	if e.enabled && !e.clockExternally {
		e.tick()
	}
}

// externalClock advances the envelope when configured for external clocking
// (driven from the address-write to env registers 24/25).
func (e *saaEnv) externalClock() {
	if e.clockExternally && e.enabled {
		e.tick()
	}
}

func (e *saaEnv) setEnvControl(data int) {
	enabled := (data & 0x80) == 0x80
	if !enabled && !e.enabled {
		return
	}
	e.enabled = enabled
	if !e.enabled {
		// was enabled, now disabled. Subsequent changes are immediate.
		e.envelopeEnded = true
		return
	}

	// Resolution (3bit/4bit) is processed immediately.
	newResolution := 1
	if (data & 0x10) == 0x10 {
		newResolution = 2
	}
	// Undocumented behaviour when changing resolution mid-waveform (see
	// SAAEnv.cpp / test case envext_34b).
	if e.resolution == 1 && newResolution == 2 {
		// 4-bit -> 3-bit
		e.phasePosition &= 0xe
	} else if e.resolution == 2 && newResolution == 1 {
		// 3-bit -> 4-bit
		e.phasePosition |= 0x1
	}
	e.resolution = newResolution

	if e.envelopeEnded {
		e.setNewEnvData(data) // also calls setLevels
		e.newData = false
	} else {
		// resolution changes arrive unbuffered, so the current level may need
		// updating:
		e.setLevels()
		e.newData = true
		e.nextData = byte(data)
	}
}

func (e *saaEnv) leftLevelVal() int  { return e.leftLevel }
func (e *saaEnv) rightLevelVal() int { return e.rightLevel }
func (e *saaEnv) isActive() bool     { return e.enabled }

func (e *saaEnv) tick() {
	if !e.enabled {
		e.envelopeEnded = true
		e.phase = 0
		e.phasePosition = 0
		return
	}

	if e.envelopeEnded {
		// do nothing (leave phase/phasePosition for SetLevels)
		return
	}

	// Continue playing the same envelope.
	e.phasePosition += byte(e.resolution)

	processNewDataIfAvailable := false
	if e.phasePosition >= 16 {
		e.phase++
		if int(e.phase) == e.numberOfPhases {
			if !e.looping {
				// position (3): non-looping waveform ended; sustain is zero.
				e.envelopeEnded = true
				processNewDataIfAvailable = true
			} else {
				// position (4): latched data acted upon ONLY here.
				e.envelopeEnded = false
				e.phase = 0
				e.phasePosition -= 16
				processNewDataIfAvailable = true
			}
		} else {
			// middle of a multi-phase envelope; buffer any commands.
			e.envelopeEnded = false
			e.phasePosition -= 16
		}
	} else {
		// still within the same phase but no longer at its start, so new data
		// must be buffered.
		e.envelopeEnded = false
	}

	if e.newData && processNewDataIfAvailable {
		e.newData = false
		e.setNewEnvData(int(e.nextData))
	} else {
		e.setLevels()
	}
}

func (e *saaEnv) setLevels() {
	switch e.resolution {
	case 2: // 3-bit res waveforms
		if e.envelopeEnded && !e.looping {
			e.leftLevel = 0
		} else {
			e.leftLevel = e.envDataPtr.levels[1][e.phase][e.phasePosition]
		}
		if e.invertRightChannel {
			e.rightLevel = 14 - e.leftLevel
		} else {
			e.rightLevel = e.leftLevel
		}
	default: // case 1: 4-bit res waveforms
		if e.envelopeEnded && !e.looping {
			e.leftLevel = 0
		} else {
			e.leftLevel = e.envDataPtr.levels[0][e.phase][e.phasePosition]
		}
		if e.invertRightChannel {
			e.rightLevel = 15 - e.leftLevel
		} else {
			e.rightLevel = e.leftLevel
		}
	}
}

func (e *saaEnv) setNewEnvData(data int) {
	e.phase = 0
	e.phasePosition = 0
	e.envDataPtr = &csEnvData[(data>>1)&0x07]
	e.invertRightChannel = (data & 0x01) == 0x01
	e.clockExternally = (data & 0x20) == 0x20
	e.numberOfPhases = e.envDataPtr.numberOfPhases
	e.looping = e.envDataPtr.looping
	if (data & 0x10) == 0x10 {
		e.resolution = 2
	} else {
		e.resolution = 1
	}
	e.enabled = (data & 0x80) == 0x80
	if e.enabled {
		e.envelopeEnded = false
	} else {
		e.envelopeEnded = true
		e.phase = 0
		e.phasePosition = 0
	}
	e.setLevels()
}
