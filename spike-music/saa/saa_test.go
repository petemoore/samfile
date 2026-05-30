// Ported to Go from Dave Hooper's SAASound (https://github.com/stripwax/SAASound),
// BSD licence — see SAASound-LICENCE.txt.

package saa

import (
	"math"
	"testing"
)

// TestSteadyToneIsNonSilent programs a single steady tone on channel 0 and
// checks the rendered output is non-trivial (both positive and negative
// excursions, non-trivial RMS). This proves the signal path works end-to-end;
// it is NOT a fidelity check.
func TestSteadyToneIsNonSilent(t *testing.T) {
	s := New(0, 44100) // 8 MHz default clock, 44.1 kHz sample rate
	// The raw amp output is a non-negative amplitude (0..480), i.e. a square
	// wave biased above zero. Enable the high-pass filter (as real playback
	// does) to centre it around zero so we get both polarities.
	s.SetHighpass(true)

	// Program channel 0: full amplitude both channels.
	s.WriteAddressData(0x00, 0xff) // amp ch0: left=0xf, right=0xf
	// Frequency: offset + octave. ~mid range.
	s.WriteAddressData(0x08, 0x80) // freq offset ch0
	s.WriteAddressData(0x10, 0x03) // octave ch0 (low nibble), ch1 (high nibble)
	// Enable tone for channel 0 in the tone mixer.
	s.WriteAddressData(0x14, 0x01) // tone mixer: ch0 on
	// Global enable (register 28 bit 0 = sound enabled, bit 1 = sync off).
	s.WriteAddressData(0x1c, 0x01)

	const numSamples = 4410
	buf := make([]int16, 2*numSamples)
	s.GenerateMany(buf, numSamples)

	var maxV, minV int16
	var sumSq float64
	var nonzero int
	for _, v := range buf {
		if v > maxV {
			maxV = v
		}
		if v < minV {
			minV = v
		}
		if v != 0 {
			nonzero++
		}
		sumSq += float64(v) * float64(v)
	}
	rms := math.Sqrt(sumSq / float64(len(buf)))

	t.Logf("max=%d min=%d rms=%.1f nonzero=%d/%d", maxV, minV, rms, nonzero, len(buf))

	if maxV <= 0 {
		t.Errorf("expected positive excursions, got max=%d", maxV)
	}
	if minV >= 0 {
		t.Errorf("expected negative excursions, got min=%d", minV)
	}
	if rms < 100 {
		t.Errorf("expected non-trivial RMS, got %.1f", rms)
	}
}

// TestSilenceWhenDisabled verifies that without a global enable the output is
// silent, confirming the enable path is actually exercised by the tone test.
func TestSilenceWhenDisabled(t *testing.T) {
	s := New(0, 44100)
	s.WriteAddressData(0x00, 0xff)
	s.WriteAddressData(0x08, 0x80)
	s.WriteAddressData(0x10, 0x03)
	s.WriteAddressData(0x14, 0x01)
	// deliberately do NOT enable register 28

	const numSamples = 1000
	buf := make([]int16, 2*numSamples)
	s.GenerateMany(buf, numSamples)

	for i, v := range buf {
		if v != 0 {
			t.Fatalf("expected silence when global enable is off, got buf[%d]=%d", i, v)
		}
	}
}
