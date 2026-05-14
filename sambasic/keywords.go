package sambasic

var keywordTable = [...]string{
	"PI",          // 0x3B  (index 0)
	"RND",         // 0x3C
	"POINT",       // 0x3D
	"FREE",        // 0x3E
	"LENGTH",      // 0x3F
	"ITEM",        // 0x40
	"ATTR",        // 0x41
	"FN",          // 0x42
	"BIN",         // 0x43
	"XMOUSE",      // 0x44
	"YHOUSE",      // 0x45
	"XPEN",        // 0x46
	"YPEN",        // 0x47
	"RAMTOP",      // 0x48
	"",            // 0x49  reserved
	"INSTR",       // 0x4A
	"INKEY$",      // 0x4B
	"SCREEN$",     // 0x4C
	"MEM$",        // 0x4D
	"",            // 0x4E  reserved
	"PATH$",       // 0x4F
	"STRING$",     // 0x50
	"",            // 0x51  reserved
	"",            // 0x52  reserved
	"SIN",         // 0x53
	"COS",         // 0x54
	"TAN",         // 0x55
	"ASN",         // 0x56
	"ACS",         // 0x57
	"ATN",         // 0x58
	"LN",          // 0x59
	"EXP",         // 0x5A
	"ABS",         // 0x5B
	"SGN",         // 0x5C
	"SQR",         // 0x5D
	"INT",         // 0x5E
	"USR",         // 0x5F
	"IN",          // 0x60
	"PEEK",        // 0x61
	"LPEEK",       // 0x62
	"DVAR",        // 0x63
	"SVAR",        // 0x64
	"BUTTON",      // 0x65
	"EOF",         // 0x66
	"PTR",         // 0x67
	"",            // 0x68  reserved
	"UDG",         // 0x69
	"",            // 0x6A  reserved
	"LEN",         // 0x6B
	"CODE",        // 0x6C
	"VAL$",        // 0x6D
	"VAL",         // 0x6E
	"TRUNC$",      // 0x6F
	"CHR$",        // 0x70
	"STR$",        // 0x71
	"BIN$",        // 0x72
	"HEX$",        // 0x73
	"USR$",        // 0x74
	"",            // 0x75  reserved
	"NOT",         // 0x76
	"",            // 0x77  reserved
	"",            // 0x78  reserved
	"",            // 0x79  reserved
	"MOD",         // 0x7A
	"DIV",         // 0x7B
	"BOR",         // 0x7C
	"",            // 0x7D  reserved
	"BAND",        // 0x7E
	"OR",          // 0x7F
	"AND",         // 0x80
	"<>",          // 0x81
	"<=",          // 0x82
	">=",          // 0x83
	"",            // 0x84  reserved
	"USING",       // 0x85  (index 74, first SingleByteKeyword)
	"WRITE",       // 0x86
	"AT",          // 0x87
	"TAB",         // 0x88
	"OFF",         // 0x89
	"WHILE",       // 0x8A
	"UNTIL",       // 0x8B
	"LINE",        // 0x8C
	"THEN",        // 0x8D
	"TO",          // 0x8E
	"STEP",        // 0x8F
	"DIR",         // 0x90
	"FORMAT",      // 0x91
	"ERASE",       // 0x92
	"MOVE",        // 0x93
	"SAVE",        // 0x94
	"LOAD",        // 0x95
	"MERGE",       // 0x96
	"VERIFY",      // 0x97
	"OPEN",        // 0x98
	"CLOSE",       // 0x99
	"CIRCLE",      // 0x9A
	"PLOT",        // 0x9B
	"LET",         // 0x9C
	"BLITZ",       // 0x9D
	"BORDER",      // 0x9E
	"CLS",         // 0x9F
	"PALETTE",     // 0xA0
	"PEN",         // 0xA1
	"PAPER",       // 0xA2
	"FLASH",       // 0xA3
	"BRIGHT",      // 0xA4
	"INVERSE",     // 0xA5
	"OVER",        // 0xA6
	"FATPIX",      // 0xA7
	"CSIZE",       // 0xA8
	"BLOCKS",      // 0xA9
	"MODE",        // 0xAA
	"GRAB",        // 0xAB
	"PUT",         // 0xAC
	"BEEP",        // 0xAD
	"SOUND",       // 0xAE
	"NEW",         // 0xAF
	"RUN",         // 0xB0
	"STOP",        // 0xB1
	"CONTINUE",    // 0xB2
	"CLEAR",       // 0xB3
	"GO TO",       // 0xB4
	"GO SUB",      // 0xB5
	"RETURN",      // 0xB6
	"REM",         // 0xB7
	"READ",        // 0xB8
	"DATA",        // 0xB9
	"RESTORE",     // 0xBA
	"PRINT",       // 0xBB
	"LPRINT",      // 0xBC
	"LIST",        // 0xBD
	"LLIST",       // 0xBE
	"DUMP",        // 0xBF
	"FOR",         // 0xC0
	"NEXT",        // 0xC1
	"PAUSE",       // 0xC2
	"DRAW",        // 0xC3
	"DEFAULT",     // 0xC4
	"DIM",         // 0xC5
	"INPUT",       // 0xC6
	"RANDOMIZE",   // 0xC7
	"DEF FN",      // 0xC8
	"DEF KEYCODE", // 0xC9
	"DEF PROC",    // 0xCA
	"END PROC",    // 0xCB
	"RENUM",       // 0xCC
	"DELETE",      // 0xCD
	"REF",         // 0xCE
	"COPY",        // 0xCF
	"",            // 0xD0  reserved
	"KEYIN",       // 0xD1
	"LOCAL",       // 0xD2
	"LOOP IF",     // 0xD3
	"DO",          // 0xD4
	"LOOP",        // 0xD5
	"EXIT IF",     // 0xD6
	"IF",          // 0xD7  long IF
	"IF",          // 0xD8  short IF
	"ELSE",        // 0xD9  long ELSE
	"ELSE",        // 0xDA  short ELSE
	"END IF",      // 0xDB
	"KEY",         // 0xDC
	"ON ERROR",    // 0xDD
	"ON",          // 0xDE
	"GET",         // 0xDF
	"OUT",         // 0xE0
	"POKE",        // 0xE1
	"DPOKE",       // 0xE2
	"RENAME",      // 0xE3
	"CALL",        // 0xE4
	"ROLL",        // 0xE5
	"SCROLL",      // 0xE6
	"SCREEN",      // 0xE7
	"DISPLAY",     // 0xE8
	"BOOT",        // 0xE9
	"LABEL",       // 0xEA
	"FILL",        // 0xEB
	"WINDOW",      // 0xEC
	"AUTO",        // 0xED
	"POP",         // 0xEE
	"RECORD",      // 0xEF
	"DEVICE",      // 0xF0
	"PROTECT",     // 0xF1
	"HIDE",        // 0xF2
	"ZAP",         // 0xF3
	"POW",         // 0xF4
	"BOOM",        // 0xF5
	"ZOOM",        // 0xF6
	"",            // 0xF7  reserved
	"",            // 0xF8  reserved
	"",            // 0xF9  reserved
	"",            // 0xFA  reserved
	"",            // 0xFB  reserved
	"",            // 0xFC  reserved
	"",            // 0xFD  reserved
	"",            // 0xFE  reserved
	"INK",         // 0xFF  INK→PEN shim (grammar §3.3); finalise rewrites to 0xA1
}

const (
	USING       SingleByteKeyword = 0x85
	WRITE       SingleByteKeyword = 0x86
	AT          SingleByteKeyword = 0x87
	TAB         SingleByteKeyword = 0x88
	OFF         SingleByteKeyword = 0x89
	WHILE       SingleByteKeyword = 0x8A
	UNTIL       SingleByteKeyword = 0x8B
	LINE        SingleByteKeyword = 0x8C
	THEN        SingleByteKeyword = 0x8D
	TO          SingleByteKeyword = 0x8E
	STEP        SingleByteKeyword = 0x8F
	DIR         SingleByteKeyword = 0x90
	FORMAT      SingleByteKeyword = 0x91
	ERASE       SingleByteKeyword = 0x92
	MOVE        SingleByteKeyword = 0x93
	SAVE        SingleByteKeyword = 0x94
	LOAD        SingleByteKeyword = 0x95
	MERGE       SingleByteKeyword = 0x96
	VERIFY      SingleByteKeyword = 0x97
	OPEN        SingleByteKeyword = 0x98
	CLOSE       SingleByteKeyword = 0x99
	CIRCLE      SingleByteKeyword = 0x9A
	PLOT        SingleByteKeyword = 0x9B
	LET         SingleByteKeyword = 0x9C
	BLITZ       SingleByteKeyword = 0x9D
	BORDER      SingleByteKeyword = 0x9E
	CLS         SingleByteKeyword = 0x9F
	PALETTE     SingleByteKeyword = 0xA0
	PEN         SingleByteKeyword = 0xA1
	PAPER       SingleByteKeyword = 0xA2
	FLASH       SingleByteKeyword = 0xA3
	BRIGHT      SingleByteKeyword = 0xA4
	INVERSE     SingleByteKeyword = 0xA5
	OVER        SingleByteKeyword = 0xA6
	FATPIX      SingleByteKeyword = 0xA7
	CSIZE       SingleByteKeyword = 0xA8
	BLOCKS      SingleByteKeyword = 0xA9
	MODE        SingleByteKeyword = 0xAA
	GRAB        SingleByteKeyword = 0xAB
	PUT         SingleByteKeyword = 0xAC
	BEEP        SingleByteKeyword = 0xAD
	SOUND       SingleByteKeyword = 0xAE
	NEW         SingleByteKeyword = 0xAF
	RUN         SingleByteKeyword = 0xB0
	STOP        SingleByteKeyword = 0xB1
	CONTINUE    SingleByteKeyword = 0xB2
	CLEAR       SingleByteKeyword = 0xB3
	GO_TO       SingleByteKeyword = 0xB4
	GO_SUB      SingleByteKeyword = 0xB5
	RETURN      SingleByteKeyword = 0xB6
	REM         SingleByteKeyword = 0xB7
	READ        SingleByteKeyword = 0xB8
	DATA        SingleByteKeyword = 0xB9
	RESTORE     SingleByteKeyword = 0xBA
	PRINT       SingleByteKeyword = 0xBB
	LPRINT      SingleByteKeyword = 0xBC
	LIST        SingleByteKeyword = 0xBD
	LLIST       SingleByteKeyword = 0xBE
	DUMP        SingleByteKeyword = 0xBF
	FOR         SingleByteKeyword = 0xC0
	NEXT        SingleByteKeyword = 0xC1
	PAUSE       SingleByteKeyword = 0xC2
	DRAW        SingleByteKeyword = 0xC3
	DEFAULT     SingleByteKeyword = 0xC4
	DIM         SingleByteKeyword = 0xC5
	INPUT       SingleByteKeyword = 0xC6
	RANDOMIZE   SingleByteKeyword = 0xC7
	DEF_FN      SingleByteKeyword = 0xC8
	DEF_KEYCODE SingleByteKeyword = 0xC9
	DEF_PROC    SingleByteKeyword = 0xCA
	END_PROC    SingleByteKeyword = 0xCB
	RENUM       SingleByteKeyword = 0xCC
	DELETE      SingleByteKeyword = 0xCD
	REF         SingleByteKeyword = 0xCE
	COPY        SingleByteKeyword = 0xCF
	// 0xD0 reserved
	KEYIN      SingleByteKeyword = 0xD1
	LOCAL      SingleByteKeyword = 0xD2
	LOOP_IF    SingleByteKeyword = 0xD3
	DO         SingleByteKeyword = 0xD4
	LOOP       SingleByteKeyword = 0xD5
	EXIT_IF    SingleByteKeyword = 0xD6
	IF_LONG    SingleByteKeyword = 0xD7
	IF_SHORT   SingleByteKeyword = 0xD8
	ELSE_LONG  SingleByteKeyword = 0xD9
	ELSE_SHORT SingleByteKeyword = 0xDA
	END_IF     SingleByteKeyword = 0xDB
	KEY        SingleByteKeyword = 0xDC
	ON_ERROR   SingleByteKeyword = 0xDD
	ON         SingleByteKeyword = 0xDE
	GET        SingleByteKeyword = 0xDF
	OUT        SingleByteKeyword = 0xE0
	POKE       SingleByteKeyword = 0xE1
	DPOKE      SingleByteKeyword = 0xE2
	RENAME     SingleByteKeyword = 0xE3
	CALL       SingleByteKeyword = 0xE4
	ROLL       SingleByteKeyword = 0xE5
	SCROLL     SingleByteKeyword = 0xE6
	SCREEN     SingleByteKeyword = 0xE7
	DISPLAY    SingleByteKeyword = 0xE8
	BOOT       SingleByteKeyword = 0xE9
	LABEL      SingleByteKeyword = 0xEA
	FILL       SingleByteKeyword = 0xEB
	WINDOW     SingleByteKeyword = 0xEC
	AUTO_KW    SingleByteKeyword = 0xED
	POP        SingleByteKeyword = 0xEE
	RECORD     SingleByteKeyword = 0xEF
	DEVICE     SingleByteKeyword = 0xF0
	PROTECT    SingleByteKeyword = 0xF1
	HIDE       SingleByteKeyword = 0xF2
	ZAP        SingleByteKeyword = 0xF3
	POW        SingleByteKeyword = 0xF4
	BOOM       SingleByteKeyword = 0xF5
	ZOOM       SingleByteKeyword = 0xF6
	// INK is the editor's compat shim for Spectrum users: GETTOKEN matches
	// "INK" at table slot 0xFF, then TOKMAIN rewrites the token to PEN
	// (0xA1) before storing. See grammar spec §3.3. The 0xFF byte never
	// appears in stored program text — finalise() applies the rewrite.
	INK SingleByteKeyword = 0xFF
)

const (
	PI_2B      TwoByteKeyword = 0x3B
	RND_2B     TwoByteKeyword = 0x3C
	POINT_2B   TwoByteKeyword = 0x3D
	FREE_2B    TwoByteKeyword = 0x3E
	LENGTH_2B  TwoByteKeyword = 0x3F
	ITEM_2B    TwoByteKeyword = 0x40
	ATTR_2B    TwoByteKeyword = 0x41
	FN_2B      TwoByteKeyword = 0x42
	BIN_2B     TwoByteKeyword = 0x43
	XMOUSE_2B  TwoByteKeyword = 0x44
	YHOUSE_2B  TwoByteKeyword = 0x45
	XPEN_2B    TwoByteKeyword = 0x46
	YPEN_2B    TwoByteKeyword = 0x47
	RAMTOP_2B  TwoByteKeyword = 0x48
	INSTR_2B   TwoByteKeyword = 0x4A
	INKEY_2B   TwoByteKeyword = 0x4B
	SCREEN_2B  TwoByteKeyword = 0x4C
	MEM_2B     TwoByteKeyword = 0x4D
	PATH_2B    TwoByteKeyword = 0x4F
	STRING_2B  TwoByteKeyword = 0x50
	SIN_2B     TwoByteKeyword = 0x53
	COS_2B     TwoByteKeyword = 0x54
	TAN_2B     TwoByteKeyword = 0x55
	ASN_2B     TwoByteKeyword = 0x56
	ACS_2B     TwoByteKeyword = 0x57
	ATN_2B     TwoByteKeyword = 0x58
	LN_2B      TwoByteKeyword = 0x59
	EXP_2B     TwoByteKeyword = 0x5A
	ABS_2B     TwoByteKeyword = 0x5B
	SGN_2B     TwoByteKeyword = 0x5C
	SQR_2B     TwoByteKeyword = 0x5D
	INT_2B     TwoByteKeyword = 0x5E
	USR_2B     TwoByteKeyword = 0x5F
	IN_2B      TwoByteKeyword = 0x60
	PEEK_2B    TwoByteKeyword = 0x61
	LPEEK_2B   TwoByteKeyword = 0x62
	DVAR_2B    TwoByteKeyword = 0x63
	SVAR_2B    TwoByteKeyword = 0x64
	BUTTON_2B  TwoByteKeyword = 0x65
	EOF_2B     TwoByteKeyword = 0x66
	PTR_2B     TwoByteKeyword = 0x67
	UDG        TwoByteKeyword = 0x69
	LEN_2B     TwoByteKeyword = 0x6B
	CODE       TwoByteKeyword = 0x6C
	VAL_STR    TwoByteKeyword = 0x6D
	VAL_2B     TwoByteKeyword = 0x6E
	TRUNC_STR  TwoByteKeyword = 0x6F
	CHR_STR    TwoByteKeyword = 0x70
	STR_STR    TwoByteKeyword = 0x71
	BIN_STR    TwoByteKeyword = 0x72
	HEX_STR    TwoByteKeyword = 0x73
	USR_STR    TwoByteKeyword = 0x74
	NOT_2B     TwoByteKeyword = 0x76
	MOD_2B     TwoByteKeyword = 0x7A
	DIV_2B     TwoByteKeyword = 0x7B
	BOR_2B     TwoByteKeyword = 0x7C
	BAND_2B    TwoByteKeyword = 0x7E
	OR_2B      TwoByteKeyword = 0x7F
	AND_2B     TwoByteKeyword = 0x80
	NOT_EQUAL  TwoByteKeyword = 0x81
	LESS_EQUAL TwoByteKeyword = 0x82
	GR_EQUAL   TwoByteKeyword = 0x83
)

func KeywordName(tokenByte byte, extended bool) (string, bool) {
	if tokenByte < 0x3B {
		return "", false
	}
	idx := int(tokenByte - 0x3B)
	if idx >= len(keywordTable) {
		return "", false
	}
	if !extended && tokenByte < 0x85 {
		return "", false
	}
	name := keywordTable[idx]
	if name == "" {
		return "", false
	}
	return name, true
}
