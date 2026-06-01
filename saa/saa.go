// Ported to Go from Dave Hooper's SAASound (https://github.com/stripwax/SAASound),
// BSD licence — see SAASound-LICENCE.txt.
//
// saa.go: top-level chip implementation (register decode, _WriteData dispatch,
// GenerateMany sample loop, combining the 6 amp outputs). Ported from
// src/SAAImpl.{cpp,h} and include/SAASound.h.
//
// Package saa is a faithful Go port of the SAA1099 sound-chip emulator used by
// SimCoupe for the SAM Coupe.
package saa

// Compile-time configuration, ported from src/defns.h.
const (
	// externalClkHz is the default SAA1099 crystal clock in Hz.
	//
	// SAM Coupe default: 8 MHz. Evidence: SimCoupe's SAADevice constructor
	// (Base/Sound.h) calls CreateCSAASound() and SetSampleRate(SAMPLE_FREQ)
	// but never calls SetClockRate, so the chip runs at SAASound's compiled
	// default EXTERNAL_CLK_HZ, which defns.h sets to 8000000.
	externalClkHz = 8000000

	// sampleRateHz is the default audio sample rate. SimCoupe uses
	// SAMPLE_FREQ = 44100 (Base/Sound.h).
	sampleRateHz = 44100

	// defaultOversample is the audio-quality oversample exponent (64x), per
	// DEFAULT_OVERSAMPLE in defns.h.
	defaultOversample = 6

	// defaultUnboostedMultiplier scales the summed amp output into a good
	// 16-bit range (DEFAULT_UNBOOSTED_MULTIPLIER in defns.h).
	defaultUnboostedMultiplier = 11.35

	// defaultBoost is the post-scale boost. DEFAULT_BOOST==1 means DO_BOOST is
	// not compiled in, i.e. no extra boost is applied.
	defaultBoost = 1.0
)

// SAASound is the public chip object, equivalent to CSAASoundInternal.
type SAASound struct {
	chip       *saaDevice
	clockRate  int
	sampleRate int
	oversample uint32
	highpass   bool

	// One-pole high-pass filter state (persistent across GenerateMany calls,
	// matching the C++ static filterout_z1_* locals).
	filterZ1Left  float64
	filterZ1Right float64

	// outputBitmask selects which of the 6 channels contribute to the mix
	// (default 0x3f = all channels), matching m_output_bitmask.
	outputBitmask byte
}

// New constructs a chip at the given clock and sample rates. Pass clockHz==0
// and/or sampleHz==0 to use the SAM Coupe defaults (8 MHz / 44100 Hz).
func New(clockHz, sampleHz int) *SAASound {
	if clockHz == 0 {
		clockHz = externalClkHz
	}
	if sampleHz == 0 {
		sampleHz = sampleRateHz
	}
	s := &SAASound{
		chip:          newSaaDevice(),
		clockRate:     externalClkHz,
		sampleRate:    sampleRateHz,
		oversample:    defaultOversample,
		highpass:      false, // matches CSAASoundInternal default ctor
		outputBitmask: 0x3f,
	}
	s.chip.setClockRate(s.clockRate)
	s.chip.setOversample(s.oversample)
	s.SetClockRate(clockHz)
	s.SetSampleRate(sampleHz)
	return s
}

// SetClockRate changes the SAA crystal clock rate (Hz) and rebuilds the
// frequency table.
func (s *SAASound) SetClockRate(hz int) {
	s.clockRate = hz
	s.chip.setClockRate(hz)
}

// SetSampleRate changes the output audio sample rate (Hz).
func (s *SAASound) SetSampleRate(hz int) {
	if hz != s.sampleRate {
		s.sampleRate = hz
		s.chip.setSampleRate(hz)
	}
}

// SetOversample changes the oversample exponent (0..6; 6 == 64x).
func (s *SAASound) SetOversample(oversample uint32) {
	if oversample != s.oversample {
		s.oversample = oversample
		s.chip.setOversample(oversample)
	}
}

// SetHighpass enables/disables the output high-pass filter.
func (s *SAASound) SetHighpass(highpass bool) {
	s.highpass = highpass
}

// SetOutputMixerBitmask selects which channels (bits 0..5) are mixed.
func (s *SAASound) SetOutputMixerBitmask(bitmask byte) {
	s.outputBitmask = bitmask
}

// WriteAddress selects the current register (OUT 511,r on a SAM).
func (s *SAASound) WriteAddress(reg byte) {
	s.chip.writeAddress(reg)
}

// WriteData writes data to the currently-selected register (OUT 255,d).
func (s *SAASound) WriteData(data byte) {
	s.chip.writeData(data)
}

// WriteAddressData selects a register and writes data in one call.
func (s *SAASound) WriteAddressData(reg, data byte) {
	s.chip.writeAddress(reg)
	s.chip.writeData(data)
}

// Clear reinitialises the virtual SAA, matching CSAASoundInternal::Clear.
func (s *SAASound) Clear() {
	s.WriteAddressData(28, 2)
	for i := 31; i >= 0; i-- {
		if i != 28 {
			s.WriteAddressData(byte(i), 0)
		}
	}
	s.WriteAddressData(28, 0)
	s.WriteAddress(0)
}

// GenerateMany fills buf with numSamples stereo frames (interleaved L,R as
// signed 16-bit). buf must hold at least 2*numSamples int16s.
//
// This mirrors CSAASoundInternal::GenerateMany + scale_for_output, except that
// it writes int16 values directly rather than little-endian bytes.
func (s *SAASound) GenerateMany(buf []int16, numSamples int) {
	oversampleScalar := float64(int(1) << s.oversample)
	for n := 0; n < numSamples; n++ {
		leftMixed, rightMixed := s.chip.tickAndOutputStereo(s.outputBitmask)
		l, r := s.scaleForOutput(leftMixed, rightMixed, oversampleScalar)
		buf[2*n] = l
		buf[2*n+1] = r
	}
}

// scaleForOutput replicates the C++ scale_for_output: divide by the oversample
// scalar, apply the unboosted multiplier, optional 5 Hz one-pole high-pass,
// optional boost, then hard-clip to int16.
func (s *SAASound) scaleForOutput(leftInput, rightInput uint32, oversampleScalar float64) (int16, int16) {
	floatLeft := float64(leftInput) / oversampleScalar
	floatRight := float64(rightInput) / oversampleScalar

	floatLeft *= defaultUnboostedMultiplier
	floatRight *= defaultUnboostedMultiplier

	if s.highpass {
		// cutoff ~5 Hz: b1 = exp(-2*pi*Fc/Fs), a0 = 1 - b1
		const b1 = 0.99928787
		const a0 = 1.0 - b1
		s.filterZ1Left = floatLeft*a0 + s.filterZ1Left*b1
		s.filterZ1Right = floatRight*a0 + s.filterZ1Right*b1
		floatLeft -= s.filterZ1Left
		floatRight -= s.filterZ1Right
	}

	// defaultBoost == 1 (DO_BOOST not active), so no boost multiply is applied.

	return clip16(floatLeft), clip16(floatRight)
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
