// Ported to Go from Dave Hooper's SAASound (https://github.com/stripwax/SAASound),
// BSD licence — see SAASound-LICENCE.txt.
//
// amp.go: per-channel amplitude / mixing (tone+noise enable, L/R levels,
// envelope multiply). Ported from src/SAAAmp.{cpp,h}.

package saa

// pdmTable models the low-pass-filtered result of the logical AND of the
// amplitude PDM and envelope PDM patterns — a more accurate evaluation of the
// SAA than amp*env. Indexed [amp_div2][env]. Ported from CSAAAmp pdm[8][16].
var pdmTable = [8][16]uint32{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	{0, 2, 3, 5, 6, 8, 9, 11, 12, 14, 15, 17, 18, 20, 21, 23},
	{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30},
	{0, 3, 5, 8, 10, 13, 15, 18, 20, 23, 25, 28, 30, 33, 35, 38},
	{0, 3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36, 39, 42, 45},
	{0, 4, 7, 11, 14, 18, 21, 25, 28, 32, 35, 39, 42, 46, 49, 53},
}

// saaAmp is one of the chip's six amplitude/mixing stages.
type saaAmp struct {
	leftLevel          uint32
	leftLevelDiv2      uint32
	rightLevel         uint32
	rightLevelDiv2     uint32
	outputIntermediate uint32
	mixMode            uint32

	connectedTone  *saaFreq
	connectedNoise *saaNoise
	connectedEnv   *saaEnv
	useEnvelope    bool

	mute          bool
	sync          bool
	lastLevelByte byte
}

func newSaaAmp(toneGen *saaFreq, noiseGen *saaNoise, envGen *saaEnv) *saaAmp {
	a := &saaAmp{
		connectedTone:  toneGen,
		connectedNoise: noiseGen,
		connectedEnv:   envGen,
		useEnvelope:    envGen != nil,
		mute:           true,
		sync:           false,
	}
	a.setAmpLevel(0x00)
	return a
}

func (a *saaAmp) setAmpLevel(levelByte byte) {
	if levelByte != a.lastLevelByte {
		a.lastLevelByte = levelByte
		a.leftLevel = uint32(levelByte & 0x0f)
		a.leftLevelDiv2 = a.leftLevel >> 1
		a.rightLevel = uint32((levelByte >> 4) & 0x0f)
		a.rightLevelDiv2 = a.rightLevel >> 1
	}
}

func (a *saaAmp) setToneMixer(enabled byte) {
	if enabled == 0 {
		a.mixMode &^= 0x01
	} else {
		a.mixMode |= 0x01
	}
}

func (a *saaAmp) setNoiseMixer(enabled byte) {
	if enabled == 0 {
		a.mixMode &^= 0x02
	} else {
		a.mixMode |= 0x02
	}
}

// setMute controls the GLOBAL mute setting (register 28 bit 0), not the
// per-channel mixer settings.
func (a *saaAmp) setMute(m bool) {
	a.mute = m
}

// doSync controls the GLOBAL sync setting (register 28 bit 1).
func (a *saaAmp) doSync(s bool) {
	a.sync = s
}

func (a *saaAmp) tick() {
	// connected oscillator always ticks
	level := a.connectedTone.tick()

	switch a.mixMode {
	case 0:
		a.outputIntermediate = 0
	case 1:
		// tone only (tone generator returns 0 or 1)
		a.outputIntermediate = uint32(level) * 2
	case 2:
		// noise only (noise level returns 0 or 1)
		a.outputIntermediate = uint32(a.connectedNoise.level()) * 2
	case 3:
		// tone+noise: tone * (2 - noise)
		a.outputIntermediate = uint32(level) * (2 - uint32(a.connectedNoise.level()))
	}
}

// effectiveAmplitude returns the effective amplitude of the low-pass-filtered
// AND of amplitude and envelope PDM patterns.
func (a *saaAmp) effectiveAmplitude(ampDiv2, env uint32) uint32 {
	return pdmTable[ampDiv2][env] * 4
}

// tickAndOutputStereo returns left/right amplitudes, each 0..480 inclusive at
// full volume without envelopes (just over 88% of that, ~424, with envelopes).
func (a *saaAmp) tickAndOutputStereo() (left, right uint32) {
	if a.sync {
		return 0, 0
	}

	a.tick()

	if a.mute {
		return 0, 0
	} else if a.useEnvelope && a.connectedEnv.isActive() {
		left = a.effectiveAmplitude(a.leftLevelDiv2, uint32(a.connectedEnv.leftLevelVal())) * (2 - a.outputIntermediate)
		right = a.effectiveAmplitude(a.rightLevelDiv2, uint32(a.connectedEnv.rightLevelVal())) * (2 - a.outputIntermediate)
		return left, right
	}
	left = a.leftLevel * a.outputIntermediate * 16
	right = a.rightLevel * a.outputIntermediate * 16
	return left, right
}
