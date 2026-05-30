#!/bin/bash
# Render all ten FRED E-Tracker tunes to WAV. See FINDINGS.md.
set -euo pipefail
cd "$(dirname "$0")"
ASSETS="${1:-$HOME/git/sam-jukebox/assets}"
go build ./cmd/render/
mkdir -p out
for m in m01 m02 m03 m04 m05 m06 m07 m08 m09 m10; do
    ./render -v -mod "$ASSETS/$m" -out "out/$m.wav"
done
echo "Done -> spike-music/out/*.wav"
