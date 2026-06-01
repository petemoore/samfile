//go:build race

package audio

// raceEnabled is true when built with -race (which also enables checkptr).
const raceEnabled = true
