#!/usr/bin/env bash
# Render all ten FRED E-Tracker tunes to looping audio files for a music
# library. See FINDINGS.md.
#
# Each module is decoded to its EXACT single loop (detected via the engine's
# order-list pointer wrap — no intro, loops from the start), the loop body is
# baked LOOPS times, encoded, named by tune title, and tagged (title/album/
# track + loop points). Default format is Apple-AAC .m4a (near-universal);
# override with FMT=flac|wavpack|opus|mp3|wav.
#
# Usage: render-all.sh [-l N] [ASSETS_DIR]
#   -l/--loops N   times to repeat the loop body (default 4)
#   FMT=...        output format (env; default aac)
#   ASSETS_DIR     dir containing m01..m10 (default ~/git/sam-jukebox/assets)
set -euo pipefail
cd "$(dirname "$0")"

LOOPS=4
while [[ $# -gt 0 ]]; do
  case "$1" in
    -l|--loops) LOOPS="$2"; shift 2;;
    *) ASSETS="$1"; shift;;
  esac
done
ASSETS="${ASSETS:-$HOME/git/sam-jukebox/assets}"
FMT="${FMT:-aac}"
OUT=out
go build ./cmd/render/
mkdir -p "$OUT"

case "$FMT" in
  aac)     ext=m4a;  enc=(-c:a aac_at -b:a 128k);;
  flac)    ext=flac; enc=(-c:a flac -compression_level 12);;
  wavpack) ext=wv;   enc=(-c:a wavpack -compression_level 8);;
  opus)    ext=opus; enc=(-c:a libopus -b:a 96k);;
  mp3)     ext=mp3;  enc=(-c:a libmp3lame -b:a 192k);;
  wav)     ext=wav;  enc=(-c:a pcm_s16le);;
  *) echo "unknown FMT=$FMT" >&2; exit 2;;
esac

# module | title | FRED issue | year | original SAM filename
tunes=(
  "m01|C.J's|51|1994|e1"
  "m02|Bubble Bobble|53|1995|e1"
  "m03|Satch|54|1995|e1"
  "m04|Crazy Stuff|55|1995|e1"
  "m05|Short But Bubbley|55|1995|e2"
  "m06|Happy Birthday|56|1995|e1"
  "m07|Happy Alien|56|1995|e2"
  "m08|Willy's World|56|1995|e3"
  "m09|Just The Way|56|1995|e4"
  "m10|M.S.C.D.T.R.H.U.F.R.D|56|1995|e5"
)

track=0
for entry in "${tunes[@]}"; do
  IFS='|' read -r m title issue year orig <<<"$entry"
  track=$((track+1))
  tmp="$(mktemp -t "$m").wav"

  # render prints: "<path>: intro=I + loop=N frames x<loops> -> ..."
  info=$(./render -loops "$LOOPS" -mod "$ASSETS/$m" -out "$tmp")
  intro=$(sed -E 's/.*intro=([0-9]+) .*/\1/' <<<"$info")
  frames=$(sed -E 's/.*loop=([0-9]+) frames.*/\1/' <<<"$info")
  intro_samples=$(( intro  * 44100 / 50 ))   # loop region starts after the intro
  loop_samples=$((  frames * 44100 / 50 ))

  ffmpeg -y -loglevel error -i "$tmp" "${enc[@]}" \
    -metadata title="$title" \
    -metadata album="SAM Coupe E-Tunes (FRED Magazine)" \
    -metadata album_artist="FRED Magazine (SAM Coupe)" \
    -metadata artist="E-Tracker / Andrew Collier player" \
    -metadata track="$track/10" \
    -metadata date="$year" \
    -metadata genre="Chiptune" \
    -metadata comment="FRED issue $issue, original file '$orig'. Decoded from E-Tracker SAA1099 module; ${LOOPS}x loop, loop=$loop_samples samples." \
    -metadata LOOP_START="$intro_samples" \
    -metadata LOOP_LENGTH="$loop_samples" \
    "$OUT/$title.$ext"
  rm -f "$tmp"
  printf "  %-26s -> %s  (loop %.1fs x%d)\n" \
    "$title.$ext" "$OUT/$title.$ext" "$(echo "scale=1;$loop_samples/44100"|bc)" "$LOOPS"
done
echo "Done -> $OUT/*.$ext"
