#!/bin/bash -eu
cd "$(dirname "${0}")"/..
VALID_FORMAT='samfile .* revision.*'
TEMP_SAMFILE_HELP="$(mktemp -t samfile-help-text.XXXXXXXXXX)"
TEMP_SAMFILE_README="$(mktemp -t samfile-readme.XXXXXXXXXX)"
TEMP_SAMFILE_BINARY="$(mktemp -t samfile.XXXXXXXXXX)"
go build -ldflags "-X main.revision=$(git rev-parse HEAD) -X main.version=$(git tag -l 'v*.*.*' --points-at HEAD | sed -n '1s/^v//p')" -o "${TEMP_SAMFILE_BINARY}" github.com/winfreddy88/samfile/cmd/samfile
"${TEMP_SAMFILE_BINARY}" --help > "${TEMP_SAMFILE_HELP}"
echo '```' >> "${TEMP_SAMFILE_HELP}"
sed -e "
   /^${VALID_FORMAT}/,/^\`\`\`\$/!b
   //!d
   /^${VALID_FORMAT}/d;r ${TEMP_SAMFILE_HELP}
   /^\`\`\`\$/d
" README.md > "${TEMP_SAMFILE_README}"
cat "${TEMP_SAMFILE_README}" > README.md
rm "${TEMP_SAMFILE_BINARY}"
rm "${TEMP_SAMFILE_README}"
rm "${TEMP_SAMFILE_HELP}"
