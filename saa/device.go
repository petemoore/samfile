// Ported to Go from Dave Hooper's SAASound (https://github.com/stripwax/SAASound),
// BSD licence — see SAASound-LICENCE.txt.
//
// device.go: wires the sub-units together (6 freq, 2 noise, 2 env, 6 amp) and
// handles register decode and the per-sample tick/output accumulation.
// Ported from src/SAADevice.{cpp,h}.

package saa

type saaDevice struct {
	currentReg    int
	outputEnabled bool
	sync          bool
	oversample    uint32

	noise [2]*saaNoise
	env   [2]*saaEnv
	osc   [6]*saaFreq
	amp   [6]*saaAmp
}

func newSaaDevice() *saaDevice {
	d := &saaDevice{}

	// Noise generators seeded with 0xffffffff (matches CSAADevice ctor).
	d.noise[0] = newSaaNoise(0xffffffff)
	d.noise[1] = newSaaNoise(0xffffffff)

	d.env[0] = newSaaEnv()
	d.env[1] = newSaaEnv()

	// Oscillators wired to noise/env generators exactly as in CSAADevice:
	//   Osc0 -> Noise0;  Osc1 -> Env0;  Osc2 -> (none)
	//   Osc3 -> Noise1;  Osc4 -> Env1;  Osc5 -> (none)
	d.osc[0] = newSaaFreq(d.noise[0], nil)
	d.osc[1] = newSaaFreq(nil, d.env[0])
	d.osc[2] = newSaaFreq(nil, nil)
	d.osc[3] = newSaaFreq(d.noise[1], nil)
	d.osc[4] = newSaaFreq(nil, d.env[1])
	d.osc[5] = newSaaFreq(nil, nil)

	// Amp stages wired to oscillator, noise and (for ch 2,5) env generators.
	d.amp[0] = newSaaAmp(d.osc[0], d.noise[0], nil)
	d.amp[1] = newSaaAmp(d.osc[1], d.noise[0], nil)
	d.amp[2] = newSaaAmp(d.osc[2], d.noise[0], d.env[0])
	d.amp[3] = newSaaAmp(d.osc[3], d.noise[1], nil)
	d.amp[4] = newSaaAmp(d.osc[4], d.noise[1], nil)
	d.amp[5] = newSaaAmp(d.osc[5], d.noise[1], d.env[1])

	d.setClockRate(externalClkHz)
	d.setOversample(defaultOversample)
	return d
}

func (d *saaDevice) setClockRate(clockRate int) {
	for i := 0; i < 6; i++ {
		d.osc[i].setClockRate(clockRate)
	}
	d.noise[0].setClockRate(clockRate)
	d.noise[1].setClockRate(clockRate)
}

func (d *saaDevice) setSampleRate(sampleRate int) {
	for i := 0; i < 6; i++ {
		d.osc[i].setSampleRate(sampleRate)
	}
	d.noise[0].setSampleRate(sampleRate)
	d.noise[1].setSampleRate(sampleRate)
}

func (d *saaDevice) setOversample(oversample uint32) {
	if oversample != d.oversample {
		d.oversample = oversample
		for i := 0; i < 6; i++ {
			d.osc[i].setOversample(oversample)
		}
		d.noise[0].setOversample(oversample)
		d.noise[1].setOversample(oversample)
	}
}

func (d *saaDevice) writeData(data byte) {
	switch d.currentReg {
	// Amplitude data (==> Amp)
	case 0:
		d.amp[0].setAmpLevel(data)
	case 1:
		d.amp[1].setAmpLevel(data)
	case 2:
		d.amp[2].setAmpLevel(data)
	case 3:
		d.amp[3].setAmpLevel(data)
	case 4:
		d.amp[4].setAmpLevel(data)
	case 5:
		d.amp[5].setAmpLevel(data)

	// Freq data (==> Osc)
	case 8:
		d.osc[0].setFreqOffset(data)
	case 9:
		d.osc[1].setFreqOffset(data)
	case 10:
		d.osc[2].setFreqOffset(data)
	case 11:
		d.osc[3].setFreqOffset(data)
	case 12:
		d.osc[4].setFreqOffset(data)
	case 13:
		d.osc[5].setFreqOffset(data)

	// Freq octave data (==> Osc) for channels 0,1
	case 16:
		d.osc[0].setFreqOctave(data & 0x07)
		d.osc[1].setFreqOctave((data >> 4) & 0x07)
	// Freq octave data (==> Osc) for channels 2,3
	case 17:
		d.osc[2].setFreqOctave(data & 0x07)
		d.osc[3].setFreqOctave((data >> 4) & 0x07)
	// Freq octave data (==> Osc) for channels 4,5
	case 18:
		d.osc[4].setFreqOctave(data & 0x07)
		d.osc[5].setFreqOctave((data >> 4) & 0x07)

	// Tone mixer control (==> Amp)
	case 20:
		d.amp[0].setToneMixer(data & 0x01)
		d.amp[1].setToneMixer(data & 0x02)
		d.amp[2].setToneMixer(data & 0x04)
		d.amp[3].setToneMixer(data & 0x08)
		d.amp[4].setToneMixer(data & 0x10)
		d.amp[5].setToneMixer(data & 0x20)

	// Noise mixer control (==> Amp)
	case 21:
		d.amp[0].setNoiseMixer(data & 0x01)
		d.amp[1].setNoiseMixer(data & 0x02)
		d.amp[2].setNoiseMixer(data & 0x04)
		d.amp[3].setNoiseMixer(data & 0x08)
		d.amp[4].setNoiseMixer(data & 0x10)
		d.amp[5].setNoiseMixer(data & 0x20)

	// Noise frequency/source control (==> Noise)
	case 22:
		d.noise[0].setSource(int(data & 0x03))
		d.noise[1].setSource(int((data >> 4) & 0x03))

	// Envelope control data (==> Env)
	case 24:
		d.env[0].setEnvControl(int(data))
	case 25:
		d.env[1].setEnvControl(int(data))

	// Global enable and reset (sync) controls
	case 28:
		// Reset (sync) bit
		bSync := (data & 0x02) != 0
		if bSync != d.sync {
			for i := 0; i < 6; i++ {
				d.osc[i].doSync(bSync)
			}
			d.noise[0].doSync(bSync)
			d.noise[1].doSync(bSync)
			for i := 0; i < 6; i++ {
				d.amp[i].doSync(bSync)
			}
			d.sync = bSync
		}

		// Global mute bit
		bOutputEnabled := (data & 0x01) != 0
		if bOutputEnabled != d.outputEnabled {
			for i := 0; i < 6; i++ {
				d.amp[i].setMute(!bOutputEnabled)
			}
			d.outputEnabled = bOutputEnabled
		}

	default:
		// register not used within the SAA-1099 architecture; ignore.
	}
}

func (d *saaDevice) writeAddress(reg byte) {
	d.currentReg = int(reg) & 31
	if d.currentReg == 24 {
		d.env[0].externalClock()
	} else if d.currentReg == 25 {
		d.env[1].externalClock()
	}
}

// tickAndOutputStereo advances the chip by one output sample (looping over the
// oversample factor) and returns the accumulated left/right mixed amplitudes.
func (d *saaDevice) tickAndOutputStereo(bitmask byte) (leftMixed, rightMixed uint32) {
	var accumLeft, accumRight uint32
	for i := 1 << d.oversample; i > 0; i-- {
		d.noise[0].tick()
		d.noise[1].tick()
		for c := 0; c < 6; c++ {
			tl, tr := d.amp[c].tickAndOutputStereo()
			if bitmask&(1<<uint(c)) != 0 {
				accumLeft += tl
				accumRight += tr
			}
		}
	}
	return accumLeft, accumRight
}
