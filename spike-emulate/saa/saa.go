// Package saa is a Go port of Dave Hooper's SAASound — a portable Philips
// SAA1099 sound-chip emulator — used here to render the SAA register writes
// captured from the emulated E-Tracker player into PCM.
//
// Ported from the C++ source at https://github.com/stripwax/SAASound
// (canonical read-only checkout: ~/git/SAASound). The port is deliberately
// faithful: class-by-class, preserving the tick/oversample/PDM-amplitude
// behaviour so the output matches SimCoupé's SAA model byte-for-byte where
// practical.
//
// ----------------------------------------------------------------------------
// SAASound - a portable Philips SAA 1099 sound chip emulator
//
// Copyright (c) 1998-2004, Dave Hooper <dave@beermex.com>
// Copyright (c) 2004-2025, Dave Hooper <dave@beermex.com> + Simon Owen <simon@simonowen.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the conditions of the BSD-3-Clause
// licence are met. See saa/LICENCE (copied verbatim from the SAASound repo).
// THIS SOFTWARE IS PROVIDED "AS IS" — see LICENCE for the full disclaimer.
// ----------------------------------------------------------------------------
package saa

// Compile-time defaults (defns.h).
const (
	externalClkHz    = 8000000 // SAA1099 crystal clock; SimCoupé leaves this default
	sampleRateHz     = 44100
	defaultOversample = 6 // 64x
	unboostedMult    = 11.35
)

// ---------------------------------------------------------------------------
// Noise generator (CSAANoise)
// ---------------------------------------------------------------------------

type noise struct {
	counter      uint64
	add          uint64
	counterLow   uint64
	oversample   uint
	counterLimit uint64
	sync         bool
	sampleRate   uint64
	sourceMode   int
	addBase      uint64
	rand         uint64
}

func newNoise(seed uint64) *noise {
	n := &noise{
		counterLimit: 1,
		sampleRate:   sampleRateHz,
		rand:         seed,
	}
	n.setClockRate(externalClkHz)
	n.add = n.addBase
	return n
}

func (n *noise) setClockRate(clock uint64) {
	// at 8MHz the noise clock is 31.25kHz = clock/256; shift <<12 for accuracy,
	// so shift by (12-8).
	n.addBase = clock << (12 - 8)
}

func (n *noise) setSampleRate(sr uint64) { n.sampleRate = sr }

func (n *noise) setOversample(os uint) {
	if os < n.oversample {
		n.counterLow <<= (n.oversample - os)
	} else {
		n.counterLow >>= (os - n.oversample)
	}
	n.counterLimit = 1 << os
	n.oversample = os
}

func (n *noise) setSource(src int) {
	n.sourceMode = src
	n.add = n.addBase >> uint(n.sourceMode)
}

func (n *noise) trigger() {
	if n.sourceMode == 3 {
		n.changeLevel()
	}
}

func (n *noise) tick() {
	if !n.sync && n.sourceMode != 3 {
		n.counter += n.add
		for n.counter >= (n.sampleRate << 12) {
			n.counter -= (n.sampleRate << 12)
			n.counterLow++
			if n.counterLow >= n.counterLimit {
				n.counterLow = 0
				n.changeLevel()
			}
		}
	}
}

func (n *noise) doSync(s bool) {
	if s {
		n.counter = 0
		n.counterLow = 0
	}
	n.sync = s
}

func (n *noise) changeLevel() {
	// 18-bit Galois LFSR, poly x^18 + x^11 + x^1 (Jepael).
	if n.rand&1 != 0 {
		n.rand = (n.rand >> 1) ^ 0x20400
	} else {
		n.rand >>= 1
	}
}

func (n *noise) level() uint { return uint(n.rand & 1) }

// ---------------------------------------------------------------------------
// Envelope generator (CSAAEnv)
// ---------------------------------------------------------------------------

type envData struct {
	numberOfPhases int
	looping        bool
	levels         [2][2][16]int // [resolution][phase][withinphase]
}

var csEnvData = [8]envData{
	{1, false, [2][2][16]int{{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, {{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{1, true, [2][2][16]int{{{15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}, {15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}}, {{14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14}, {14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14}}}},
	{1, false, [2][2][16]int{{{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, {{14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{1, true, [2][2][16]int{{{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, {{14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{2, false, [2][2][16]int{{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}}, {{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}}}},
	{2, true, [2][2][16]int{{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}}, {{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {14, 14, 12, 12, 10, 10, 8, 8, 6, 6, 4, 4, 2, 2, 0, 0}}}},
	{1, false, [2][2][16]int{{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, {{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
	{1, true, [2][2][16]int{{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, {{0, 0, 2, 2, 4, 4, 6, 6, 8, 8, 10, 10, 12, 12, 14, 14}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}}},
}

type env struct {
	leftLevel, rightLevel int
	data                  *envData
	enabled               bool
	invertRight           bool
	phase                 int
	phasePosition         int
	envelopeEnded         bool
	looping               bool
	numberOfPhases        int
	resolution            int
	newData               bool
	nextData              int
	clockExternally       bool
}

func newEnv() *env {
	e := &env{envelopeEnded: true, resolution: 1}
	e.setEnvControl(0)
	return e
}

func (e *env) internalClock() {
	if e.enabled && !e.clockExternally {
		e.tick()
	}
}

func (e *env) externalClock() {
	if e.clockExternally && e.enabled {
		e.tick()
	}
}

func (e *env) setEnvControl(nData int) {
	bEnabled := (nData & 0x80) == 0x80
	if !bEnabled && !e.enabled {
		return
	}
	e.enabled = bEnabled
	if !e.enabled {
		e.envelopeEnded = true
		return
	}

	newResolution := 1
	if nData&0x10 == 0x10 {
		newResolution = 2
	}
	if e.resolution == 1 && newResolution == 2 {
		e.phasePosition &= 0xe
	} else if e.resolution == 2 && newResolution == 1 {
		e.phasePosition |= 0x1
	}
	e.resolution = newResolution

	if e.envelopeEnded {
		e.setNewEnvData(nData)
		e.newData = false
	} else {
		e.setLevels()
		e.newData = true
		e.nextData = nData
	}
}

func (e *env) leftLvl() int  { return e.leftLevel }
func (e *env) rightLvl() int { return e.rightLevel }
func (e *env) isActive() bool { return e.enabled }

func (e *env) tick() {
	if !e.enabled {
		e.envelopeEnded = true
		e.phase = 0
		e.phasePosition = 0
		return
	}
	if e.envelopeEnded {
		return
	}

	e.phasePosition += e.resolution

	bProcessNewData := false
	if e.phasePosition >= 16 {
		e.phase++
		if e.phase == e.numberOfPhases {
			if !e.looping {
				e.envelopeEnded = true
				bProcessNewData = true
			} else {
				e.envelopeEnded = false
				e.phase = 0
				e.phasePosition -= 16
				bProcessNewData = true
			}
		} else {
			e.envelopeEnded = false
			e.phasePosition -= 16
		}
	} else {
		e.envelopeEnded = false
	}

	if e.newData && bProcessNewData {
		e.newData = false
		e.setNewEnvData(e.nextData)
	} else {
		e.setLevels()
	}
}

func (e *env) setLevels() {
	switch e.resolution {
	case 2:
		if e.envelopeEnded && !e.looping {
			e.leftLevel = 0
		} else {
			e.leftLevel = e.data.levels[1][e.phase][e.phasePosition]
		}
		if e.invertRight {
			e.rightLevel = 14 - e.leftLevel
		} else {
			e.rightLevel = e.leftLevel
		}
	default: // case 1
		if e.envelopeEnded && !e.looping {
			e.leftLevel = 0
		} else {
			e.leftLevel = e.data.levels[0][e.phase][e.phasePosition]
		}
		if e.invertRight {
			e.rightLevel = 15 - e.leftLevel
		} else {
			e.rightLevel = e.leftLevel
		}
	}
}

func (e *env) setNewEnvData(nData int) {
	e.phase = 0
	e.phasePosition = 0
	e.data = &csEnvData[(nData>>1)&0x07]
	e.invertRight = (nData & 0x01) == 0x01
	e.clockExternally = (nData & 0x20) == 0x20
	e.numberOfPhases = e.data.numberOfPhases
	e.looping = e.data.looping
	if nData&0x10 == 0x10 {
		e.resolution = 2
	} else {
		e.resolution = 1
	}
	e.enabled = (nData & 0x80) == 0x80
	if e.enabled {
		e.envelopeEnded = false
	} else {
		e.envelopeEnded = true
		e.phase = 0
		e.phasePosition = 0
	}
	e.setLevels()
}

// ---------------------------------------------------------------------------
// Frequency / tone generator (CSAAFreq)
// ---------------------------------------------------------------------------

const initialLevel = 1

type freq struct {
	freqTable     *[2048]uint64
	counter       uint64
	add           uint64
	counterLow    uint64
	oversample    uint
	counterLimit  uint64
	level         int
	currentOffset int
	currentOctave int
	nextOffset    int
	nextOctave    int
	ignoreOffset  bool
	newData       bool
	sync          bool
	sampleRate    uint64

	noiseGen *noise
	envGen   *env
	mode     int // 0 nothing, 1 env, 2 noise
}

func newFreq(table *[2048]uint64, ng *noise, eg *env) *freq {
	mode := 0
	if ng != nil {
		mode = 2
	} else if eg != nil {
		mode = 1
	}
	f := &freq{
		freqTable:    table,
		counterLimit: 1,
		level:        initialLevel,
		sampleRate:   sampleRateHz,
		noiseGen:     ng,
		envGen:       eg,
		mode:         mode,
	}
	f.setAdd()
	return f
}

func (f *freq) setFreqOffset(nOffset int) {
	if !f.sync {
		f.nextOffset = nOffset
		f.newData = true
		if f.nextOctave == f.currentOctave {
			f.ignoreOffset = true
		}
	} else {
		f.newData = false
		f.ignoreOffset = false
		f.currentOffset = nOffset
		f.nextOffset = nOffset
		f.currentOctave = f.nextOctave
		f.setAdd()
	}
}

func (f *freq) setFreqOctave(nOctave int) {
	if !f.sync {
		f.nextOctave = nOctave
		f.newData = true
		f.ignoreOffset = false
	} else {
		f.newData = false
		f.ignoreOffset = false
		f.currentOctave = nOctave
		f.nextOctave = nOctave
		f.currentOffset = f.nextOffset
		f.setAdd()
	}
}

func (f *freq) updateOctaveOffsetData() {
	if !f.newData {
		return
	}
	f.currentOctave = f.nextOctave
	if !f.ignoreOffset {
		f.currentOffset = f.nextOffset
		f.newData = false
	}
	f.ignoreOffset = false
	f.setAdd()
}

func (f *freq) setSampleRate(sr uint64) { f.sampleRate = sr }

func (f *freq) setOversample(os uint) {
	if os < f.oversample {
		f.counterLow <<= (f.oversample - os)
	} else {
		f.counterLow >>= (os - f.oversample)
	}
	f.counterLimit = 1 << os
	f.oversample = os
}

func (f *freq) tick() int {
	if f.sync {
		return 1
	}
	f.counter += f.add
	for f.counter >= (f.sampleRate << 12) {
		f.counter -= (f.sampleRate << 12)
		f.counterLow++
		if f.counterLow >= f.counterLimit {
			f.counterLow = 0
			f.level = 1 - f.level
			switch f.mode {
			case 1:
				f.envGen.internalClock()
			case 2:
				f.noiseGen.trigger()
			}
			f.updateOctaveOffsetData()
		}
	}
	return f.level
}

func (f *freq) setAdd() {
	f.add = f.freqTable[f.currentOctave<<8|f.currentOffset]
}

func (f *freq) doSync(s bool) {
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

// ---------------------------------------------------------------------------
// Amplitude / mixer (CSAAAmp)
// ---------------------------------------------------------------------------

var pdm = [8][16]uint{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8},
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	{0, 2, 3, 5, 6, 8, 9, 11, 12, 14, 15, 17, 18, 20, 21, 23},
	{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30},
	{0, 3, 5, 8, 10, 13, 15, 18, 20, 23, 25, 28, 30, 33, 35, 38},
	{0, 3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36, 39, 42, 45},
	{0, 4, 7, 11, 14, 18, 21, 25, 28, 32, 35, 39, 42, 46, 49, 53},
}

type amp struct {
	leftLevel     uint
	leftLevelDiv2 uint
	rightLevel    uint
	rightLevelDiv2 uint
	outputInter   uint
	mixMode       uint
	tone          *freq
	noiseGen      *noise
	envGen        *env
	useEnvelope   bool
	mute          bool
	sync          bool
	lastLevelByte uint8
}

func newAmp(tone *freq, ng *noise, eg *env) *amp {
	a := &amp{
		tone:        tone,
		noiseGen:    ng,
		envGen:      eg,
		useEnvelope: eg != nil,
		mute:        true,
	}
	a.setAmpLevel(0x00)
	return a
}

func (a *amp) setAmpLevel(b uint8) {
	if b != a.lastLevelByte {
		a.lastLevelByte = b
		a.leftLevel = uint(b & 0x0f)
		a.leftLevelDiv2 = a.leftLevel >> 1
		a.rightLevel = uint((b >> 4) & 0x0f)
		a.rightLevelDiv2 = a.rightLevel >> 1
	}
}

func (a *amp) setToneMixer(enabled byte) {
	if enabled == 0 {
		a.mixMode &^= 0x01
	} else {
		a.mixMode |= 0x01
	}
}

func (a *amp) setNoiseMixer(enabled byte) {
	if enabled == 0 {
		a.mixMode &^= 0x02
	} else {
		a.mixMode |= 0x02
	}
}

func (a *amp) doMute(m bool) { a.mute = m }
func (a *amp) doSync(s bool) { a.sync = s }

func (a *amp) tick() {
	level := a.tone.tick()
	switch a.mixMode {
	case 0:
		a.outputInter = 0
	case 1:
		a.outputInter = uint(level) * 2
	case 2:
		a.outputInter = a.noiseGen.level() * 2
	case 3:
		a.outputInter = uint(level) * (2 - a.noiseGen.level())
	}
}

func effectiveAmplitude(ampDiv2, envv uint) uint {
	return pdm[ampDiv2][envv] * 4
}

func (a *amp) tickAndOutputStereo() (left, right uint) {
	if a.sync {
		return 0, 0
	}
	a.tick()
	if a.mute {
		return 0, 0
	}
	if a.useEnvelope && a.envGen.isActive() {
		left = effectiveAmplitude(a.leftLevelDiv2, uint(a.envGen.leftLvl())) * (2 - a.outputInter)
		right = effectiveAmplitude(a.rightLevelDiv2, uint(a.envGen.rightLvl())) * (2 - a.outputInter)
		return left, right
	}
	left = a.leftLevel * a.outputInter * 16
	right = a.rightLevel * a.outputInter * 16
	return left, right
}

// ---------------------------------------------------------------------------
// Device (CSAADevice) + public chip (CSAASoundInternal)
// ---------------------------------------------------------------------------

// SAA is a virtual SAA1099 (CSAASound + CSAADevice merged). Drive it with
// WriteAddress/WriteData and pull PCM with GenerateMany.
type SAA struct {
	freqTable [2048]uint64
	clockRate uint64

	noise [2]*noise
	env   [2]*env
	osc   [6]*freq
	amp   [6]*amp

	currentReg    uint8
	outputEnabled bool
	syncFlag      bool
	oversample    uint
	outputBitmask byte

	highpass bool
	// persistent highpass filter state (was function-static in C++).
	z1Left, z1Right float64
}

// New creates a SAA1099 with SimCoupé's defaults: 8MHz clock, 44100Hz,
// 64x oversample, highpass on, full output mixer.
func New() *SAA {
	s := &SAA{
		oversample:    defaultOversample,
		outputBitmask: 0x3f,
		highpass:      true,
		clockRate:     externalClkHz,
	}
	s.computeFreqTable(externalClkHz)

	s.noise[0] = newNoise(0xffffffff)
	s.noise[1] = newNoise(0xffffffff)
	s.env[0] = newEnv()
	s.env[1] = newEnv()

	s.osc[0] = newFreq(&s.freqTable, s.noise[0], nil)
	s.osc[1] = newFreq(&s.freqTable, nil, s.env[0])
	s.osc[2] = newFreq(&s.freqTable, nil, nil)
	s.osc[3] = newFreq(&s.freqTable, s.noise[1], nil)
	s.osc[4] = newFreq(&s.freqTable, nil, s.env[1])
	s.osc[5] = newFreq(&s.freqTable, nil, nil)

	s.amp[0] = newAmp(s.osc[0], s.noise[0], nil)
	s.amp[1] = newAmp(s.osc[1], s.noise[0], nil)
	s.amp[2] = newAmp(s.osc[2], s.noise[0], s.env[0])
	s.amp[3] = newAmp(s.osc[3], s.noise[1], nil)
	s.amp[4] = newAmp(s.osc[4], s.noise[1], nil)
	s.amp[5] = newAmp(s.osc[5], s.noise[1], s.env[1])

	s.setClockRate(externalClkHz)
	s.setOversample(defaultOversample)
	return s
}

func (s *SAA) computeFreqTable(clock uint64) {
	if clock == s.clockRate && s.freqTable[0x100] != 0 {
		// already populated for this clock
	}
	s.clockRate = clock
	ix := 0
	for octave := 0; octave < 8; octave++ {
		for offset := 0; offset < 256; offset++ {
			s.freqTable[ix] = uint64((8192.0 * 15625.0 * float64(uint64(1)<<uint(octave)) * (float64(clock) / 8000000.0)) / (511.0 - float64(offset)))
			ix++
		}
	}
}

func (s *SAA) setClockRate(clock uint64) {
	s.computeFreqTable(clock)
	for i := 0; i < 2; i++ {
		s.noise[i].setClockRate(clock)
	}
}

// SetClockRate lets callers match a non-default SAA crystal (Hz).
func (s *SAA) SetClockRate(clock uint64) { s.setClockRate(clock) }

func (s *SAA) setSampleRate(sr uint64) {
	for i := 0; i < 6; i++ {
		s.osc[i].setSampleRate(sr)
	}
	for i := 0; i < 2; i++ {
		s.noise[i].setSampleRate(sr)
	}
}

func (s *SAA) setOversample(os uint) {
	if os != s.oversample {
		s.oversample = os
	}
	for i := 0; i < 6; i++ {
		s.osc[i].setOversample(os)
	}
	for i := 0; i < 2; i++ {
		s.noise[i].setOversample(os)
	}
}

// SetHighpass toggles the 5Hz DC-blocking high-pass on the output stage.
func (s *SAA) SetHighpass(b bool) { s.highpass = b }

// WriteAddress selects the current register (OUT 511,r).
func (s *SAA) WriteAddress(nReg uint8) {
	s.currentReg = nReg & 31
	if s.currentReg == 24 {
		s.env[0].externalClock()
	} else if s.currentReg == 25 {
		s.env[1].externalClock()
	}
}

// WriteData writes to the currently-selected register (OUT 255,d).
func (s *SAA) WriteData(nData uint8) {
	switch s.currentReg {
	case 0, 1, 2, 3, 4, 5:
		s.amp[s.currentReg].setAmpLevel(nData)
	case 8, 9, 10, 11, 12, 13:
		s.osc[s.currentReg-8].setFreqOffset(int(nData))
	case 16:
		s.osc[0].setFreqOctave(int(nData & 0x07))
		s.osc[1].setFreqOctave(int((nData >> 4) & 0x07))
	case 17:
		s.osc[2].setFreqOctave(int(nData & 0x07))
		s.osc[3].setFreqOctave(int((nData >> 4) & 0x07))
	case 18:
		s.osc[4].setFreqOctave(int(nData & 0x07))
		s.osc[5].setFreqOctave(int((nData >> 4) & 0x07))
	case 20:
		s.amp[0].setToneMixer(nData & 0x01)
		s.amp[1].setToneMixer(nData & 0x02)
		s.amp[2].setToneMixer(nData & 0x04)
		s.amp[3].setToneMixer(nData & 0x08)
		s.amp[4].setToneMixer(nData & 0x10)
		s.amp[5].setToneMixer(nData & 0x20)
	case 21:
		s.amp[0].setNoiseMixer(nData & 0x01)
		s.amp[1].setNoiseMixer(nData & 0x02)
		s.amp[2].setNoiseMixer(nData & 0x04)
		s.amp[3].setNoiseMixer(nData & 0x08)
		s.amp[4].setNoiseMixer(nData & 0x10)
		s.amp[5].setNoiseMixer(nData & 0x20)
	case 22:
		s.noise[0].setSource(int(nData & 0x03))
		s.noise[1].setSource(int((nData >> 4) & 0x03))
	case 24:
		s.env[0].setEnvControl(int(nData))
	case 25:
		s.env[1].setEnvControl(int(nData))
	case 28:
		bSync := nData&0x02 != 0
		if bSync != s.syncFlag {
			for i := 0; i < 6; i++ {
				s.osc[i].doSync(bSync)
				s.amp[i].doSync(bSync)
			}
			s.noise[0].doSync(bSync)
			s.noise[1].doSync(bSync)
			s.syncFlag = bSync
		}
		bOut := nData&0x01 != 0
		if bOut != s.outputEnabled {
			for i := 0; i < 6; i++ {
				s.amp[i].doMute(!bOut)
			}
			s.outputEnabled = bOut
		}
	}
}

// WriteAddressData is WriteAddress followed by WriteData.
func (s *SAA) WriteAddressData(nReg, nData uint8) {
	s.WriteAddress(nReg)
	s.WriteData(nData)
}

// Clear reinitialises the chip exactly as CSAASound::Clear does.
func (s *SAA) Clear() {
	s.WriteAddressData(28, 2)
	for i := 31; i >= 0; i-- {
		if i != 28 {
			s.WriteAddressData(uint8(i), 0)
		}
	}
	s.WriteAddressData(28, 0)
	s.WriteAddress(0)
}

func (s *SAA) tickAndOutputStereo() (leftMixed, rightMixed uint) {
	var accumL, accumR uint
	for i := 1 << s.oversample; i > 0; i-- {
		s.noise[0].tick()
		s.noise[1].tick()
		for c := 0; c < 6; c++ {
			l, r := s.amp[c].tickAndOutputStereo()
			if s.outputBitmask&(1<<uint(c)) != 0 {
				accumL += l
				accumR += r
			}
		}
	}
	return accumL, accumR
}

// GenerateMany renders nSamples of 16-bit signed stereo PCM (interleaved
// L,R little-endian) into buf, which must hold nSamples*4 bytes. Mirrors
// CSAASoundInternal::GenerateMany + scale_for_output.
func (s *SAA) GenerateMany(buf []byte, nSamples int) {
	oversampleScalar := float64(uint64(1) << s.oversample)
	const b1 = 0.99928787
	const a0 = 1.0 - b1
	p := 0
	for ; nSamples > 0; nSamples-- {
		lIn, rIn := s.tickAndOutputStereo()
		fl := float64(lIn) / oversampleScalar
		fr := float64(rIn) / oversampleScalar
		fl *= unboostedMult
		fr *= unboostedMult
		if s.highpass {
			s.z1Left = fl*a0 + s.z1Left*b1
			s.z1Right = fr*a0 + s.z1Right*b1
			fl -= s.z1Left
			fr -= s.z1Right
		}
		lo := clip16(fl)
		ro := clip16(fr)
		buf[p] = byte(lo & 0xff)
		buf[p+1] = byte((lo >> 8) & 0xff)
		buf[p+2] = byte(ro & 0xff)
		buf[p+3] = byte((ro >> 8) & 0xff)
		p += 4
	}
}

func clip16(v float64) int16 {
	if v > 32767 {
		return 32767
	}
	if v < -32768 {
		return -32768
	}
	return int16(v)
}
