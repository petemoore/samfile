;W 004,00000,027,16383

               DUMP 1,&0009        ;SAVE "MD23"CODE 32777,15750
               ORG  &4009


;          ANDREW J.A. WRIGHT
;
;          COPYRIGHT 1990
;
;          MASTER DOS


;--------------
Part_A1:

FS:            EQU  &4000

COMM:          EQU  224
TRCK:          EQU  225
SECT:          EQU  226
DTRQ:          EQU  227
DTRQ2:         EQU  243

HCHRWR:        EQU  167
HCHRRD:        EQU  168

GND:           EQU  &4009

VAR2:          EQU  &5A00

MRND:          EQU  &A5
MIN:           EQU  191
MOUT:          EQU  223

DEVL:          EQU  &5A06
DEVN:          EQU  &5A07

PFLAGT:        EQU  &5A50
INVERT:        EQU  &5A54
WINDRHS:       EQU  &5A56
SPOSNL:        EQU  &5A6E
CUSCRNP:       EQU  &5A78
CSTAT:         EQU  &5A7B
SAVARSP:       EQU  &5A81
SAVARS:        EQU  &5A82
NUMENDP:       EQU  &5A84
NUMEND:        EQU  &5A85
NVARSP:        EQU  &5A87
NVARS:         EQU  &5A88
ELINP:         EQU  &5A93
ELINE:         EQU  &5A94
CHADP:         EQU  &5A96
CHADD:         EQU  &5A97
; KCURP:       EQU  &5A99          ;NotUsed
; KCUR:        EQU  &5A9A          ;NotUsed
PROGP:         EQU  &5A9F
PROG:          EQU  &5AA0
XPTRP:         EQU  &5AA2
XPTR:          EQU  &5AA3
OPSTORE:       EQU  &5AB5
INQUFG:        EQU  &5ABA
CURCMD:        EQU  &5B74
DOSFLG:        EQU  &5BC2
DOSCNT:        EQU  &5BC3
FLAGS:         EQU  &5C3B
TVFLAG:        EQU  &5C3C
ERRSP:         EQU  &5C3D
BORDCR:        EQU  &5C4B          ;BORDCOL
CHANS:         EQU  &5C4F
CURCHL:        EQU  &5C51
DOSSTK:        EQU  &5C59
FRAMES:        EQU  &5C78

FOWIA:         EQU  &0004          ;1*
HLJUMP:        EQU  &0005          ;1*
IYJUMP:        EQU  &0006          ;7*
; NRWRITE:     EQU  &000D          ;NotUsed
; NRREAD:      EQU  &00AC          ;NotUsed
WKROOM:        EQU  &0109          ;1*
JMKRBIG:       EQU  &010C          ;3*
STREAM:        EQU  &0112          ;6*
; POMSG:       EQU  &0115          ;NotUsed
EXPNUM:        EQU  &0118          ;4*
EXPSTR:        EQU  &011B          ;1*
EXPEXP:        EQU  &011E          ;1*
GETINT:        EQU  &0121          ;2*
GETSTR:        EQU  &0124          ;1*
STKSTR:        EQU  &0127          ;1*
JCLSBL:        EQU  &014E          ;1*
CLSLOW:        EQU  &0151          ;1*
JRECLAIM:      EQU  &0163          ;4*
; KBFLSH:      EQU  &0166          ;NotUsed
RDKEY:         EQU  &0169          ;1*
BEEPR:         EQU  &016F          ;1*
JSTRS:         EQU  &017E          ;1*
JGTTOK:        EQU  &018A          ;1*

SELURPG:       EQU  &3FDF      ;!   ;4*
; INCURPGDE:   EQU  &3FEB      ;!	;NotUsed
CHKHLR:        EQU  &3FEF      ;!   ;1*
INCURPAGE:     EQU  &3FF2      ;!   ;2*

FTADD:         EQU  &A280          ;8*      (SCR in section C)
UNSTLEN:       EQU  &3F8C      ;!  ;1*

CALC:          EQU  &EF
DIVN:          EQU  &05
SWOP:          EQU  &06
IDIV:          EQU  &09
MOD:           EQU  &08
STO4:          EQU  &D5
RCL4:          EQU  &DD
FIVELIT:       EQU  &27
EXIT:          EQU  &33
EXIT2:         EQU  &34
DUP:           EQU  &25

L1303:         EQU  &5A69          ;1*
CMDV:          EQU  &5AF4          ;1*
OVERF:         EQU  &5BB9          ;2*
COMAD:         EQU  &5BDA          ;1*
; INSBF:       EQU  &4F00          ;NotUsed

CHBTLO:        EQU  11
CHBTHI:        EQU  12
CHREC:         EQU  13
CHNAME:        EQU  14
CHFLAG:        EQU  24
CHDRIV:        EQU  25
RECFLG:        EQU  67
RECNUM:        EQU  68
RCLNLO:        EQU  69
RCLNHI:        EQU  70

DRES:          EQU  %00001001
SEEK:          EQU  %00011011
STPIN:         EQU  %01011011
STPOUT:        EQU  %01111011

DRSEC:         EQU  %10000000
DWSEC:         EQU  %10100010
RADD:          EQU  %11000000
RTRK:          EQU  %11100000
WTRK:          EQU  %11110010      ;*

RESP:          EQU  233
PPORT:         EQU  232
ULA:           EQU  254

MDRV:          EQU  11
MFLG:          EQU  12
RPTL:          EQU  13
RPTH:          EQU  14
BUFL:          EQU  15
BUFH:          EQU  16
NSRL:          EQU  17
NSRH:          EQU  18
FFSA:          EQU  19
NAME:          EQU  20
RTYP:          EQU  230            ;30
CNTH:          EQU  30
CNTL:          EQU  31
FTRK:          EQU  32
FSCT:          EQU  33
FSAM:          EQU  34

DIRT:          EQU  &FA ;DISP TO TAG VALUE IN DIR FILE DIR ENTRY
WTYP:          EQU  230
PTRL:          EQU  231
LENL:          EQU  233
;RDRAM:        EQU  39
RDRAM:         EQU  275

WRRAM:         EQU  275
NTRK:          EQU  785
NSECT:         EQU  786

NSAM:          EQU  GND-9
SAM:           EQU  NSAM+15

DCHAN:         EQU  GND+&3C00-9
SVBC:          EQU  DCHAN
SVDE:          EQU  DCHAN+2
RFDH:          EQU  DCHAN+4
SVHL:          EQU  DCHAN+5
SVIX:          EQU  DCHAN+7
REG1:          EQU  DCHAN+9
DRIVE:         EQU  DCHAN+11
FLAG3:         EQU  DCHAN+12
RPT:           EQU  DCHAN+13
BUF:           EQU  DCHAN+15
NSR:           EQU  DCHAN+17
FSA:           EQU  DCHAN+19

DRAM:          EQU  FSA+256

;PATH NAMES FOR DRIVES 1 AND 2

MPL:           EQU  38             ;MAX PATH LEN

PTH1:          EQU  DRAM+512
PTH2:          EQU  PTH1+MPL       ;TO ABOUT 7F60H

STR:           EQU  &7F90          ;REGISTER STORE AND STACK
                                   ;USES 7F68-7F91H (COPY)
                                   ;OR STACK USE+22 REGS (SNAP)

SYSP:          EQU  &7FF0          ;SYNTAX STACK

DFT:           EQU  21             ;DIRECTORY FILE TYPE
RDLIM:         EQU  8              ;ALLOW RAM DISCS 3-7
ALLOCT:        EQU  &5100
MRPRT:         EQU  128            ;MEGA RAM PORT
XTRA:          EQU  &5896          ;XTRA HOOKS ETC. IN SYS PAGE


;--------------
Part_B1:

               ORG  GND

Fix_L4009_42:  ;4.2 Chg
;*PTHRD:       EQU  &40D0

Fix_L4009_43:  ;4.3 Chg
PTHRD:         EQU  &40D2          ;TEMP RAM DISC PATH NAME

;FixEnd
           ;THIS AREA (256 BYTES) WILL BE OVER-WRITTEN BY NSAM

START:         LD   HL,&8000+511
               LD   E,(HL)
               DEC  HL
               LD   D,(HL)         ;T/S OF REST OF FILE

DOS:           XOR  A
               LD   (START+FS),A
               LD   (START+1+FS),HL ;USE AS VARIABLE

               LD   A,E
               OUT  (SECT),A

DOS2:          IN   A,(COMM)
               BIT  0,A
               JR   NZ,DOS2

               IN   A,(TRCK)
               CP   D
               JR   Z,DOS4

               LD   A,STPOUT
               JR   NC,DOS3
               LD   A,STPIN

DOS3:          OUT  (COMM),A
               LD   B,20
DEL1:          DJNZ DEL1
               JR   DOS2

DOS4:          DI
               LD   A,DRSEC
               OUT  (COMM),A
               LD   B,20
DEL2:          DJNZ DEL2

               LD   HL,(START+1+FS)
               LD   BC,DTRQ
               JR   DOS6

DOS5:          INI

DOS6:          IN   A,(COMM)
               BIT  1,A
               JR   NZ,DOS5
               RRCA
               JR   C,DOS6

;CHECK DISC ERR COUNT

               AND  &0D
               JR   Z,DOS8

               LD   A,(START+FS)
               INC  A
               LD   (START+FS),A
               PUSH AF
               AND  2
               JR   Z,DOS7

               LD   A,DRES
               OUT  (COMM),A
               LD   B,20

DEL3:          DJNZ DEL3

DOS7:          POP  AF
               CP   10
               JR   C,DOS2

               RST  8
               DEFB 19

DOS8:          DEC  HL
               LD   E,(HL)
               DEC  HL
               LD   D,(HL)
               LD   A,D
               OR   E
               JR   NZ,DOS

               LD   HL,SVHDR+FS
               LD   B,SNME-SVHDR

SVERL:         LD   (HL),A
               INC  HL
               DJNZ SVERL          ;CLEAR SYS VARS FOR
                                   ;CONSISTENT RESULT
               LD   A,(&5CB4)      ;LAST PAGE
               LD   (PORT2+FS),A
               DEC  A
               LD   (SNPRT2+FS),A
               IN   A,(251)
               AND  &1F
               LD   (DOSFLG),A     ;DOSFLG = DOS PAGE

               LD   H,&51          ;ALLOCT/256
               LD   L,A
               LD   (HL),&60       ;ALLOCT MARKED "DISK USE"

               LD   HL,&0144       ;DEVICE: "D", 1
               LD   (&5A06),HL

               XOR  A
               LD   (SAMCNT+FS),A
               LD   (&5BC3),A      ;"DOS NOT IN CONTROL"

               LD   BC,DTRQ2

TDLP:          OUT  (C),B
               LD   A,20

TDDL:          DEC  A
               JR   NZ,TDDL        ;WAIT 50 USEC

               IN   A,(C)
               CP   B
               LD   A,0            ;"NO SECOND DRIVE"
               JR   NZ,TD2         ;JR IF NO SECOND DISC CHIP

               DJNZ TDLP

               LD   A,(TRAKS2+FS)
               AND  A
               JR   NZ,TD3         ;JR IF TRACKS ALREADY SET
                                   ;FOR DISC 2
               LD   A,128+80

TD2:           LD   (TRAKS2+FS),A
                                   ;    IF "BOOT"  / IF "BOOT 1"
TD3:           POP  HL             ;RETURN TO ROM1 / NEXT STAT
               POP  DE             ;NEXT STAT      / ERR HANDLER
               POP  BC             ;ERR HANDLER    /      ?
               PUSH BC
               PUSH DE
               PUSH HL
               BIT  7,H
               JR   Z,TD4          ;JR IF HL=NEXT STAT (BOOT 1)

               EX   DE,HL
               LD   D,B
               LD   E,C
L40CC: ;*
TD4:           LD   (NEXTST+FS),HL ;KEEP A RECORD TO USE IN SOME
               LD   (L1303),DE     ; SYNTAX HANDLING
               CALL MRINIT+FS      ;MEGA RAM INIT

               LD   HL,SERDT+FS
               LD   DE,&4BA0
               LD   BC,DTEND-SERDT
               LDIR              ;COPY CODE TO SYS PAGE AT 4BA0H

               LD   HL,&4BB0
               LD   (&5ADE),HL     ;PRTOKV

               LD   HL,XTRA+MTV-PFV
               LD   (&5AFA),HL     ;MTOKV

               LD   HL,&4BB0+EVV-PVECT
               LD   (&5AF6),HL     ;EVALUV

               LD   HL,&4BB0+SLVP-PVECT
               LD   (CMDV),HL
               JP   INIP3+FS       ;IN HOOKS.DOS

L40FC:         DEFB 0,0,0,0                                 ;***
Pad_4100:      DEFS &4100-$

FFHL:          DEFM "BO"
FFDE:          DEFB "O","T"+&80

ENTSP:         DEFW 0
SNPRT0:        DEFB &1F
SNPRT1:        DEFB 2
SNPRT2:        DEFB 0
SNPSVA:        DEFB 0

;CLEARED FROM HERE ON:

SVHDR:         DEFW 0
CCHAD:         DEFW 0
CNT:           DEFW 0

DSC:           DEFB 0
DCT:           DEFB 0
DST:           DEFB 0
               DEFS 7
NBOT:          DEFB 0
RCMR:          DEFB 0
COUNT:         DEFB 0
SVA:           DEFB 0
SVC:           DEFB 0
NRFLG:         DEFB 0     ;WAS SAMCNT
RMSE:          DEFB 0
SMSE:          DEFB 0

SVDPT:         DEFW 0
SVTRS:         DEFW 0
SVBUF:         DEFW 0
SVCNT:         DEFW 0
HLDI:          DEFW 0
PTRSCR:        DEFW 0

PORT1:         DEFB 0
PORT2:         DEFB 0
PORT3:         DEFB 0
SVCST:         DEFB 0

TSTR1:         DEFB 0
OSTR1:         DEFB 0
CSTR1:         DEFB 0
HSTR1:         DEFB 0

DSTR1:         DEFB 0
FSTR1:         DEFB 0
SSTR1:         DEFB 0
LSTR1:         DEFB 0
NSTR1:         DEFB 0
               DEFS 14
HD001:         DEFB 0
HD0B1:         DEFW 0
HD0D1:         DEFW 0
HD0F1:         DEFW 0
PGES1:         DEFB 0
PAGE1:         DEFB 0

DSTR2:         DEFB 0
FSTR2:         DEFB 0
SSTR2:         DEFB 0
LSTR2:         DEFB 0
NSTR2:         DEFB 0
               DEFS 14
HD002:         DEFB 0
HD0B2:         DEFW 0
HD0D2:         DEFW 0
HD0F2:         DEFW 0
PGES2:         DEFB 0
PAGE2:         DEFB 0

NSTR3:         DEFB 0
               DEFS 14

UIFA:          DEFB 0
               DEFS 47

DIFA:          DEFB 0
               DEFS 47

HKA:           DEFB 0
HKHL:          DEFW 0
HKDE:          DEFW 0
HKBC:          DEFW 0

;END OF CLEARED SECTION

SNME:          DEFB &13
               DEFM "SNAP          "
               DEFB &13
SNLEN:         DEFW 49152
SNADD:         DEFW 16384
               DEFW 0
               DEFW &FFFF

FSLOT:         DEFW 0  ;SECT/TRK OF FREE DIR SLOT, OR 00?? IF NO
FSLTE:         DEFB 0  ;0 OR 1 FOR DIR ENTRY IN SECTOR

;MAIN PROGRAM ENTRY POINT

L41FF:         DEFB 0                                       ;***
Pad_4200:      DEFS &4200-$

               JP   HOOK
               JP   SYNTAX
               JP   NMI

TEMPB1:        DEFB 0
DTKSX:         DEFB 0     ;USED BY SNDFL AS DTKS
HKSP:          DEFW 0
RDAT:          DEFB 0
SAMDR:         DEFB 0

L420F:         DEFB 0                                       ;***
Pad_4210:      DEFS &4210-$

               DEFW ERRTBL+FS

TEMPW1:        DEFW 0
TEMPW2:        DEFW 0
TEMPW3:        DEFW 0
TEMPW4:        DEFW 0
TCNT:          DEFW 0     ;DIR'S 'TOTAL FILES ON DISC' COUNTER
FCNT:          DEFW 0     ;DIR'S 'FILES IN CURRENT DIR' COUNTER
NEXTST:        DEFW 0

Pad_4220:      DEFS &4220-$

DVAR:          EQU  $     ;4220H

Fix_L4220_42:  ;4.2 Chg
;*RBCC:        DEFB 0  ;7  ;0 (v2.2=0 / v2.3=7 ?)
;*TRAKS1:      DEFB 128+80 ;1
;*TRAKS2:      DEFB 0      ;2
;*STPRAT:      DEFB 0      ;3
;*STPRT2:      DEFB 0      ;4
;*CHDIR:       DEFB " "    ;5  DIR SPACE CHAR
;*NSTAT:       DEFB 1      ;6
;*VERS:        DEFB 42 ;43 ;7  VERSION 2.2 =([DVAR7]-20)/10
;*DCOLS:       DEFB 0      ;8  AUTO-SET FOR DIR COLUMN NUMBER
;*SRTFG:       DEFB 0  ;1  ;9  SORTED DIR "ON"
;*DELIM:       DEFB &0D    ;10 POINT'S DELIMITER
;*FNSEP:       DEFB 0  ;".";11 SEPARATOR IN FILE NAMES
;*                         ;   - USED BY PFNAME
;*RTSYM:       DEFM "\/"   ;12 (2) ROOT SYMBOLS
;*SKEW:        DEFB &FF    ;14 FOR SKEW 1. 0FEH FOR SKEW 2
;*ODEF:        DEFB 1      ;15 DEFAULT DRIVE
;*DTKS:        DEFB 4      ;16 NUM. OF DIR TRACKS (4 OR MORE)
;*CDIRT:       DEFW 0      ;17 (2) CURRENT DIR CODE (TEMP)
;*                         ;   AND EXX VERSION
;*DTFLG:       DEFB 0      ;19 NZ IF DATES TO BE SHOWN IN DIR
;*SAMCNT:      DEFB 0      ;20
;*MAXT:        DEFB 0      ;21 MAX TAG VALUE FOR SUB DIRS
;*             DEFW SAMHK  ;22 (2) ADDR OF HOOKS
;*MSFLG:       DEFB 1      ;24 INVERT CHARS >127 MOVED TO SCREEN
;*                         ;   USE 1 TO PRINT ALL EXCEPT 0FFH

Fix_L4220_43:  ;4.3 Chg
RBCC:          DEFB 7  ;0
TRAKS1:        DEFB 128+80 ;1 (v2.2=0 / v2.3=7 ?)
TRAKS2:        DEFB 0      ;2
STPRAT:        DEFB 0      ;3
STPRT2:        DEFB 0      ;4
CHDIR:         DEFB " "    ;5  DIR SPACE CHAR
NSTAT:         DEFB 1      ;6
VERS:          DEFB 43 ;42 ;7  VERSION 2.3 =([DVAR7]-20)/10
DCOLS:         DEFB 0      ;8  AUTO-SET FOR DIR COLUMN NUMBER
SRTFG:         DEFB 1  ;0  ;9  SORTED DIR "ON"
DELIM:         DEFB &0D    ;10 POINT'S DELIMITER
FNSEP:         DEFB ".";0  ;11 SEPARATOR IN FILE NAMES
                           ;   - USED BY PFNAME
RTSYM:         DEFM "\/"   ;12 (2) ROOT SYMBOLS
SKEW:          DEFB &FF    ;14 FOR SKEW 1. 0FEH FOR SKEW 2
ODEF:          DEFB 1      ;15 DEFAULT DRIVE
DTKS:          DEFB 4      ;16 NUM. OF DIR TRACKS (4 OR MORE)
CDIRT:         DEFW 0      ;17 (2) CURRENT DIR CODE (TEMP)
                           ;   AND EXX VERSION
DTFLG:         DEFB 0      ;19 NZ IF DATES TO BE SHOWN IN DIR
SAMCNT:        DEFB 0      ;20
MAXT:          DEFB 0      ;21 MAX TAG VALUE FOR SUB DIRS
               DEFW SAMHK  ;22 (2) ADDR OF HOOKS
MSFLG:         DEFB 0  ;1  ;24 INVERT CHARS >127 MOVED TO SCREEN
                           ;   USE 1 TO PRINT ALL EXCEPT 0FFH

;FixEnd

MSUPC:         DEFB "."    ;25 CHAR TO USE FOR CHARS <32 OR =255
NMIKP:         DEFB 4      ;26 PAGE AT 8000H IF NMI & "1" OR "5"
NMIKA:         DEFW &0004  ;27 (2) ADDR CALLED IF DITTO.
                           ;   E=PORT (BITS 4-0)
DWAI:          DEFB 0      ;29 NUMBER OF .25 SECS BEFORE SAVE,-1
EXTADD:        CALL CMR    ;
ONERR:         DEFW 0      ;33 (2)
               RET         ;
EAPG:          DEFB 0      ;36 PAGE USED IF ONERR IS ABOVE 8000H
MSINC:         DEFW &0200  ;37 (2) MULTI-SECTOR INCREMENT

;TABLE OF TRKS/DRIVE FOR EACH RAM DISC
RDDT:          DEFB 0,0,0,0,0     ;(5) 39-43 RAM DISCS 3,4,5,6,7
                                  ;    START AT 0 TKS

;TABLE OF FIRST PAGE FOR EACH RAM DISC
FIPT:          DEFB 0,0,0,0,0     ;(5) 44-48

;TABLE OF DIRECTORY NUMBER FOR EACH DRIVE
CDIT:          DEFB 0,0,0,0,0,0,0 ;(7) 49-55

;TABLE OF PATH LENGTH FOR EACH DRIVE
PLT:           DEFB 2,2,2,2,2,2,2 ;(7) 56-62

;TABLE OF RANDOM WORDS FOR EACH DRIVE
CRWT:          DEFW 0,0,0,0,0,0,0 ;(14) 63-76

SAMRN:         DEFW 0      ;(2) 77 SAM RND NO. (DISC RND NO.
                           ;       WHEN SAM CREATED)
TDVAR:         DEFW 0      ;(2) 79

DATDT:         DEFM "00/00/00"     ;(9) 81-88 DD/MM/YY
               DEFB &0D            ;    89             CR

               DEFB 31,1,12,1,99,0 ;(6) 90-95 DD/MM/YY
                                   ;          HI/LOW LIMITS
;???
TIMDT:         DEFM "00:00:00"     ;(9) 96-103 HH:MM:SS
               DEFB &0D            ;    104             CR

               DEFB 23,0,59,0,59,0 ;(6) 105-110 HH/MM/SS
                                   ;            HI/LOW LIMITS

;TABLE OF DRIVES EACH DRIVE PRETENDS TO BE
DRPT:          DEFB 1,2,3,4,5,6,7  ;(7) 111-117

;TABLE OF BITS, 1 PER 16K PAGE IN A MEGA RAM=64 BITS =8 BYTES
;4 POSSIBLE MEGA RAMS=32 BYTES

MRTAB:         DEFS &20            ;(32) 118-149

CKPT:          DEFB &EF            ;150 CLOCK PORT
BEEPT:         DEFW &0085          ;151 (2) BEEP TIME
XXPTR:         DEFW 0              ;(2) 153 XPTR STORE

               DEFB SYNTAX-CTAB/3  ;153 NO. OF CMDS IN CTAB

;CMD VALUE AND ADDR TABLE
L42BC: ;*
CTAB:          DEFB &86        ;WRITE
               DEFW WRITE

               DEFB &90        ;DIR
               DEFW DIR

               DEFB &91        ;FORMAT
               DEFW WFOD

               DEFB &92        ;ERASE
               DEFW ERAZ

               DEFB &93        ;MOVE
               DEFW MOVE

               DEFB &95        ;LOAD
               DEFW LOAD

               DEFB &98        ;OPEN
               DEFW OPEN

               DEFB &99        ;CLOSE
               DEFW CLOSE

               DEFB &B3        ;CLEAR
               DEFW CLEAR

               DEFB &B8        ;READ
               DEFW READ

               DEFB &CF        ;COPY
               DEFW COPY

               DEFB &E3        ;RENAME
               DEFW RENAM

               DEFB &E4        ;CALL
               DEFW CALL_Label

               DEFB &F1        ;PROTECT
               DEFW PROT

               DEFB &F2        ;HIDE
               DEFW HIDE

               DEFB &F7        ;BACKUP
               DEFW BACKUP

               DEFB &F8        ;TIME
               DEFW TIME

               DEFB &F9        ;DATE
               DEFW DATE

               DEFB 0          ;LAST ENTRY FOR NOT FOUND
               DEFW CNF

L42F5: ;*
SYNTAX:        DI
               CP   53         ;NO DOS
               JR   Z,SYNT1

               CP   29         ;NOT UNDERSTOOD
               JR   NZ,ST3HP

SYNT1:         LD   (SVCST),A  ;ERROR NUMBER
               LD   HL,NRFLG   ;NO RECURSE FLAG
               INC  (HL)
               DEC  (HL)
               LD   (HL),1     ;"NO RECURSE NOW"

ST3HP:         JP   NZ,SYNT3   ;JP IF E.G. DOS CALLED EXPT1NUM
                               ; AND GOT VAL "#!"
                               ;GIVING NONSENSE ERROR - DO NOT
                               ; RECURSIVELY CALL DOS!
               CALL SETSTK
               LD   BC,ENDS
               PUSH BC
               CALL ZFSP       ;ZERO FLAG3 AND HKSP
               CALL RESREG
               CALL GTIXD

;GET CHADD AND SAVE IT

               CALL NRRDD
               DEFW CHADD

               LD   (CCHAD),BC

;GET START OF STATEMENT

               CALL NRRDD
               DEFW CSTAT

               CALL NRWRD
               DEFW CHADD

               CALL GCHR
               CP   ":"
               CALL Z,GTNC    ;SKIP ANY COLON, CSTAT POINTS TO
                              ;":" IF E.G. "ON X: DIR: PRINT..."

               EX   DE,HL
               LD   HL,CTAB-1
               LD   B,(HL)     ;NO. OF CMDS
               DEC  HL

LCMDL:         INC  HL
               INC  HL
               CP   (HL)
               INC  HL
               JR   Z,CMDF

               DJNZ LCMDL    ;IF NOT FOUND, END ON "CNF" ROUTINE

CMDF:          LD   C,(HL)
               INC  HL
               LD   B,(HL)
               EX   DE,HL      ;HL=CHAD
               PUSH BC
               RET

;NO MATCH IN TABLE

CNF:           INC  A          ;POINT IS FF 3D
               JR   NZ,CKESV

               CALL GTNC
               CP   &3D
               JP   Z,POINTC

;CHECK EXTERNAL SYNTAX VECTOR

CKESV:         POP  BC         ;JUNK ENDS
               EX   (SP),HL    ;CHAD TO STACK, GET ENTSP
               LD   (ENTSP),HL ;ORIG
               LD   BC,(CCHAD)
               CALL NRWRD
               DEFW CHADD

               POP  BC         ;CHAD PTR TO CMD OR AFTER FF
               LD   HL,(ONERR)
               LD   A,H
               OR   L
               LD   A,(SVCST)
               JR   Z,SYNT3

               BIT  7,H
               JP   Z,EXTADD

               IN   A,(251)
               LD   D,A        ;ORIG
               LD   A,(EAPG)
               OUT  (251),A
               LD   A,(SVCST)
               JP   (HL)
L437F: ;*
SYNT3:         LD   E,0        ;NO ACTION
               LD   HL,NRFLG   ;RECURSE OK
               LD   (HL),E
               RET               ;A=ERROR NO.

;SAMDOS HOOK CODE ROUTINE

HOOK:          DI
               ADD  A,A        ;JUNK BIT 7, MAKE WORD OFFSET
               LD   (SVCST),A  ;HOOK CODE OFFSET
               EX   AF,AF'
               LD   (HKA),A
               CALL SETSTK
               LD   (SVHDR),IX
               EXX
               LD   (HKHL),HL
               LD   (HKDE),DE
               LD   (HKBC),BC
               EXX
               CALL ZFSP
               CALL NRWR       ;A=0
               DEFW XPTR+1     ;NO ERROR IN CASE PRINT
               LD   HL,(SVCST) ;H=JUNK
               LD   A,(HKA)
               LD   DE,SAMHK
               CALL INDJP      ;INDEXED JP USING L

;RETURN FROM HOOK CODE O.K

               CALL BCR
               EXX               ;PASS HL,DE,BC OUT VIA ALT REGS
               XOR  A          ;"NO ERROR"
               LD   E,A        ;"NO ACTION"
               JP   RENT       ;RESTORE ORIG ENTSP
L43C0: ;*
SETSTK:        POP  IY         ;RET ADDR
               CALL NRRDD
               DEFW DOSSTK     ;READ - SET TO 7FF0H BY BASIC

               LD   H,B        ; AT EACH CMD
               LD   L,C
               POP  AF
               POP  DE
               POP  BC     ;GET ALL 3 VALUES FROM STACK AT 8000H
               LD   SP,HL      ;NEW STACK SO THAT
                               ; DOS->CALL ROM->HOOKS CAN AVOID
                               ; OVERWRITING MAIN DOS STACK
               PUSH BC
               PUSH DE
               PUSH AF
               LD   HL,(ENTSP)
               PUSH HL         ;OLD VALUE RESTORED ON EXIT
               LD   (ENTSP),SP
               JP   (IY)

;ZERO FLAG3, HKSP

ZFSP:          XOR  A
               LD   (FLAG3),A
               LD   H,A
               LD   L,A
               LD   (HKSP),HL
               RET

;RESET ALL REGISTERS

RESREG:        LD   HL,TSTR1
               LD   B,UIFA-TSTR1
               LD   A,&FF

RESR1:         LD   (HL),A
               INC  HL
               DJNZ RESR1

               LD   (RDAT),A
               RET

;COMMAND CODE TABLE
L43F3: ;*
SAMHK:         DEFW INIT    ;1;128
               DEFW HGTHD     ;129
               DEFW HLOAD     ;130
               DEFW HVERY     ;131
               DEFW HSAVE   ;2;132
               DEFW SKSAFE    ;133
               DEFW HOPEN     ;134
               DEFW HCLOS     ;135
               DEFW HAUTO   ;3;136
               DEFW HSKTD     ;137 SEEK TRACK D
               DEFW HDUMMY ;S ;138 FORMAT TRACK UNDER HEAD,
                              ;    USING DE AS 1ST T/S
               DEFW HVAR      ;139
               DEFW HEOF    ;4;140
               DEFW HPTR      ;141
               DEFW HPATH     ;142
               DEFW HLDPG     ;143 LIKE 130 BUT A=PAGE
               DEFW HVEPG   ;5;144 LIKE 131 BUT A=PAGE
               DEFW HSDIR     ;145 SET DIR. A=PAGE, DE=START,
                              ;    BC=LEN OF NAME
               DEFW ROFSM     ;146
               DEFW HOFLE     ;147
               DEFW SBYT    ;6;148
               DEFW HWSAD     ;149
               DEFW HKSB      ;150 SAVE ADE FROM HL
               DEFW HDBOP     ;151 O/P BC FROM DE TO FILE (IX)
               DEFW SCFSM   ;7;152
               DEFW HORDER    ;153 A=LEN TO SORT ON, BC=STR. LEN
                              ;    DE=NO., HL=START
               DEFW HDUMMY ;S ;154
               DEFW HDUMMY ;S ;155
               DEFW HDUMMY ;S ;156
               DEFW HDUMMY ;S ;157
               DEFW HGFLE     ;158
               DEFW LBYT      ;159
               DEFW HRSAD     ;160
               DEFW HLDBK     ;161
               DEFW HFRSAD    ;162 FAR READ IX SECTORS
                              ;    FROM DSC A, TRK D, SCT E
                              ;    TO PAGE C, OFFSET HL
               DEFW HFWSAD    ;163 FAR WRITE IX SECTORS
                              ;    TO DSC A, TRK D, SCT E
                              ;    FROM PAGE C, OFFSET HL
               DEFW REST      ;164
               DEFW PCAT      ;165
               DEFW HERAZ     ;166
               DEFW MCHWR     ;167
               DEFW MCHRD     ;168
               DEFW HPTV      ;169 PRINT TOKEN "A"
               DEFW HPFF      ;170 POST FF
               DEFW HGTTK     ;171 GET TOKEN
               DEFW HKLEN     ;172 EVALUATOR PATCH
               DEFW HSLMV     ;173 SAVE/LOAD ETC PATCH
               DEFW RCPTCH    ;174 RUN/CLEAR PATCH


;--------------
Part_C11:


COMMP:         PUSH AF
               LD   A,(DSC)
               LD   C,A
               POP  AF
               RET

TRCKP:         CALL COMMP
               INC  C
               RET

COMMR:         PUSH BC
               CALL COMMP
               IN   A,(C)
               POP  BC
               RET

;RETURN Z IF PGES1/DE IS ZERO
L4465: ;*
CKDE:          LD   A,D
               OR   E
               RET  NZ

               LD   A,(PGES1)
               AND  A
               RET  Z

               DEC  A
               LD   (PGES1),A
               LD   DE,16384
               JR   CKDE

;PRECOMPENSATION CALCULATOR

PRECMX:        LD   C,DWSEC

PRECMP:        CALL TSTD       ;GET TRACKS ON DISC
               RRA
               AND  &3F        ;** MASK TOP BIT (SIDEDNES),
               LD   B,A        ;HALVE TRK/SIDE USUALLY 40D
               LD   A,D
               AND  &7F
               SUB  B         ;E.G. TRACK 0-40=CY, SO JR LEAVING
                              ;BIT 1 SET  (PRECOMP=DISABLED)
                              ;E.G. TRACK 70-40=NC, NO JR,
                              ;RES 1=ENABLE
               JR   C,SADC     ;JR IF CURRENT TRK IN D IS AN
                               ; INNER TRACK
               RES  1,C        ;ELSE TURN *ON* PRECMP

;SEND A DISC COMMAND

SADC:          CALL BUSY
SDCX:          LD   A,C
               CALL COMMP
               OUT  (C),A
               LD   B,20
SDC1:          DJNZ SDC1
               RET

DWAIT:         CALL TIRD
               RET  NC

               LD   C,&C0      ;READ ADDRESS CMD CODE
               CALL SADC     ;SPIN UP DISC IN CASE IT IS OFF
                             ;(OR STEP IN/OUT WILL START DRIVE,
                             ;THEN WSAD WILL WRITE TOO SOON
                             ;'COS "DRIVE RUNNING") WAIT TILL
                             ;FINISHED CMD (ALTERS SECT REG)
;TEST FOR CHIP BUSY

BUSY:          CALL COMMR
               BIT  0,A
               RET  Z
               CALL BRKTST
               JR   BUSY

;WRITE SECTOR IF ALTERED

WRIF:          CALL GTNSR      ;GET DE=CUR SECTOR

WRIF2:         BIT  3,(IX+MFLG)
               RET  Z          ;RET IF SECTOR HAS NOT BEEN
                               ;WRITTEN TO
NWSAD:         CALL DWAIT

;WRITE SECTOR AT DE

WSAD:          CALL TIRDXDCT
               JP   NC,RDWSCT

               DI
L44BB: ;*
WSA1:          CALL CTAS
               CALL PRECMX     ;WRITE SECTOR CMD
               LD   A,(DSC)
               LD   (WSA3+1),A ;SELF-MOD STATUS PORT
               ADD  A,3
               LD   C,A        ;DATA PORT
               CALL GTBUF
               CALL WSA3
               RLCA
               BIT  6,A
               JP   NZ,REP23   ;ERROR IF WRITE-PROTECTED

               CALL CDEC
               JR   WSA1

WSA2:          OUTI
WSA3:          IN   A,(0)
               RRCA
               RET  NC         ;RET IF DONE
               RRCA
               JR   C,WSA2     ;JR IF DISK READY FOR BYTE
               JR   WSA3

RSADSV:        LD   DE,(SVDE)

;READ SECTOR AT DE

RSAD:          CALL TIRDXDCT
               JP   NC,RDRSCT

RSA1:          CALL RSSR
               CALL RDDATA
               RLCA
               CALL CDEC
               JR   RSA1

RDDATA:        LD   A,(DSC)    ;STATUS PORT
               LD   (RSA3+1),A
               ADD  A,3
               LD   (RSA2+1),A
               JR   RSA3

RSA2:          IN   A,(0)
               LD   (HL),A
               INC  HL

RSA3:          IN   A,(0)
               RRCA
               RET  NC         ;RET IF READ SECTOR CMD FINISHED

               RRA
               JR   NC,RSA3    ;JR IF NO BYTE IS READY TO READ
               JR   RSA2

;OPTIONAL RSAD OR NRSAD (IF SAM WANTED)

ORSAD:         BIT  5,(IX+4)
               JR   Z,RSAD

;READ SECTOR AT DE

NRSAD:         CALL TIRDXDCT
               JP   NC,NRDRSCT
L4522: ;*
NRSA1:         CALL RSSR
               EXX
               LD   L,&FF
               LD   D,NSAM/256
               EXX
               PUSH DE
               CALL NRDDATA
               POP  DE
               RLCA
               CALL CDEC
               JR   NRSA1

NRDDATA:       LD   A,(DSC)        ;STATUS PORT
               LD   (TGT1+1),A
               LD   (NRSA3+1),A
               ADD  A,3
               LD   (NRSA2+1),A
               LD   B,2            ;MASK FOR BIT 1
               JR   NRSA3

NRS22:         LD   H,D      ;NSAM MSB
               AND  A        ;Z IF DIR ENTRY IS ERASED OR UNUSED
               JR   NZ,NRS25

               LD   H,A      ;DUMP DATA TO ROM

NRS25:         EXX
TGT1:          IN   A,(0)
               AND  B
               JR   Z,NRSA3  ;JR IF NO BYTE IS READY TO READ

NRSA2:         IN   A,(0)
               LD   (HL),A
               INC  HL
               EXX
               INC  L
               JR   Z,NRS22  ;JR IF START OF DIR ENTRY

               OR   (HL)
               LD   (HL),A
               EXX

NRSA3:         IN   A,(0)
               RRCA
               RET  NC       ;RET IF READ SECTOR CMD FINISHED
               RRA
               JR   NC,NRSA3 ;JR IF NO BYTE IS READY TO READ
               JR   NRSA2

;SEARCH VERSION OF RSAD
;ENTRY: DE=T/S, (TEMPW1)=DELIM COUNT, (DELIM)=DELIM
;EXIT: TEMPW3 HOLDS NEW DELIM COUNT, TEMPW2=LOCN
L4567: ;*
SRSAD:         CALL TIRDXDCT
               JP   NC,RDSSAD

SRSA1:         CALL RSSR
               LD   A,(DSC)        ;STATUS PORT
               LD   (SRSA3+1),A
               ADD  A,3
               LD   (SRSA2+1),A
               PUSH DE
               LD   DE,(TEMPW1)    ;DELIM COUNT
               LD   A,(DELIM)
               LD   B,A
               CALL SRSA3
               LD   (TEMPW3),DE    ;NEW COUNT
               POP  DE
               RLCA
               CALL CDEC
               JR   SRSA1

SRSA2:         IN   A,(0)
               LD   (HL),A
               INC  HL
               CP   B
               JR   Z,SRSA4

SRSA3:         IN   A,(0)
               RRCA
               RET  NC         ;RET IF READ SECTOR CMD FINISHED

               RRA
               JR   NC,SRSA3   ;JR IF NO BYTE IS READY TO READ
               JR   SRSA2

SRSA4:         DEC  DE
               LD   A,D
               OR   E
               JR   NZ,SRSA3

               LD   (TEMPW2),HL
               JR   SRSA3

;CHECK DISC ERR COUNT

CDEC:          AND  &1C
               JR   NZ,CDE1    ;JR IF AN ERROR WAS DETECTED

               CALL CLRRPT
               POP  HL         ;JUNK RET ADDR
               JP   GTBUF

CDE1:          PUSH AF
               LD   A,(DCT)
               INC  A
               LD   (DCT),A
               CP   10
               JP   NC,REP4

               POP  AF
               BIT  4,A
               JR   NZ,CTSL    ;JR IF RECORD NOT FOUND

               CALL INSTP
               CALL OUTSTP
               CALL OUTSTP
               JP   INSTP

;CONFIRM TRACK/SECTOR LOCATION

CTSL:          LD   C,RADD
               CALL SADC
               LD   HL,DST
               CALL RDDATA
               AND  &1C
               JR   NZ,CTS1    ;JR IF ERROR

               CALL TRCKP
               LD   A,(DST)
               OUT  (C),A
               RET

CTS1:          LD   A,(DCT)
               INC  A
               LD   (DCT),A
               CP   8
               JP   NC,REP5

               AND  2
               JR   Z,CTS2

               PUSH DE
               CALL REST
               POP  DE
               JR   CTSL

CTS2:          CALL INSTP
               JR   CTSL

;AHL=START, CDE=LEN

HLDPG:         OUT  (251),A

;LOAD FILE ALREADY OPENED BY HGTHD. HL=START, CDE=LEN, PAGED IN.

HLOAD:         LD   BC,&4BB0+HLDP-PVECT
               CALL NETPA
               CALL DSCHD
               CALL LDBLK

;SEEK LAST TRACK IN DIRECTORY, OR TRACK 3 IF LAST TRACK IS 4
;(AVOID 1ST FILE!)
;USED TO LEAVE HEAD SOMEWHERE MORE SAFE IN CASE OF RESET
L4617: ;*
SKSAFE:        CALL TIRD
               RET  NC

               LD   A,(DTKS)
               DEC  A
               LD   D,A
               LD   E,1        ;PROB NOT NEEDED
               CP   4
               JR   NZ,SEEKD

               DEC  D
               JR   SEEKD

;CONFIRM TRACK AND SEEK

CTAS:          LD   A,D
               OR   E
   ;         * JR   NZ,CTA1

  ;          * CALL BITF2
               JP   Z,REP27    ;END OF FILE

;            * LD   SP,(ENTSP)
 ;           * XOR  A
  ;          * LD   E,A
   ;         * RET

CTA1:          CALL BCC        ;BORDER COLOUR CHANGE
               EXX

HSKTD:         EXX

;SEEK TRACK D. USED BY FORMAT, SKSAFE

SEEKD:         CALL SELD
               INC  A
               INC  A
               LD   C,A        ;SECTOR PORT
               OUT  (C),E

CTA2:          LD   A,D
               AND  &7F
               LD   B,A
               CALL BUSY
               CALL TRCKP
               IN   A,(C)
               CP   B
               RET  Z

;CHECK IF PAGE OVER C000H EVERY TIME TRACK ALTERS

               PUSH AF
               CALL BITF6
               JR   Z,CTA3     ;JR IF NOT LOAD/SAVE BLOCK

CTA25:         LD   A,(SVHL+1)
               CP   &C0
               JR   C,CTA3

               RES  6,A
               LD   (SVHL+1),A
               IN   A,(251)
               INC  A
               LD   (PORT1),A
               OUT  (251),A

CTA3:          POP  AF
               CALL NC,OUTSTP
               CALL C,INSTP
               JR   CTA2

OUTSTP:        LD   C,STPOUT
               JR   STEP

INSTP:         LD   C,STPIN

STEP:          CALL SADC

;STEP DELAY ROUTINE

STPDEL:        PUSH HL
               LD   HL,STPRAT
               LD   A,(DSC)
               BIT  4,A
               JR   Z,STPD1

               INC  HL
STPD1:         LD   A,(HL)
               POP  HL
               AND  A

STPD2:         RET  Z

STPDX:         PUSH AF
               LD   BC,150

STPD4:         DEC  BC
               LD   A,B
               OR   C
               JR   NZ,STPD4

               POP  AF
               DEC  A
               JR   STPD2

;RESTORE DISC DRIVE

RESTX:
REST:          LD   DE,&0001
               CALL TIRD
               RET  NC

               CALL TFIHO      ;TEST FOR INDEX HOLE
               JP   Z,REP6     ;"CHECK DISC IN DRIVE" IF NO HOLE

;TEST FOR TRACK 00

RSLP4:         CALL COMMR
               BIT  2,A
               JP   NZ,BUSY

;STEP OUT ONE TRACK

               CALL OUTSTP
               JR   RSLP4

;TEST FOR INDEX HOLE
;EXIT: NZ=OK, Z=TIMED OUT, NO HOLE

TFIHO:         CALL SELD
               LD   C,&D0
               CALL SDCX       ;RESET DISC CHIP
               LD   B,0

TFIHL:         DJNZ TFIHL
               LD   H,&80      ;LOOPS BEFORE TIME-OUT

TFI2:          CALL COMMR
               BIT  1,A
               JR   Z,TFI3

               DEC  HL
               LD   A,H
               OR   L
               JR   NZ,TFI2
               RET

TFI3:          CALL COMMR
               BIT  1,A
               RET  NZ

               DEC  HL
               LD   A,H
               OR   L
               JR   NZ,TFI3
               RET

;FORMAT...TO... SR

FTOSR:         CALL TIRD
               JP   C,EXDAT

REP10:         CALL DERR
               DEFB 91         ;INVALID DEVICE IF RAM DISC

HFWCD:         LD   A,(HKA)    ;DRIVE
               CALL CODN       ;CONVERT
               JR   CKDRX

;CHECK VALID SPECIFIER DISC

CKDISC:        LD   A,(LSTR1)
               CP   "D"
               JR   NZ,REP10   ;"INVALID DEVICE"

;CHECK DRIVE NUMBER

CKDRV:         LD   A,(DSTR1)

CKDRX:         CP   1
               JR   Z,CKDV1

               CP   2
               JR   Z,CKDV0

               DEC  A
               CP   RDLIM-1
               JP   NC,REP22

               INC  A          ;ALLOWS 3 TO RDLIM-1
               JR   CKDV1

CKDV0:         LD   A,(RBCC+2)
               CP   0
               JP   Z,REP22    ;"NO SUCH DRIVE"

               LD   A,2

CKDV1:         LD   (DRIVE),A
               RET

;SERIAL SET DRIVE, RETURN CURRENT ONE IN C

SSDRV:         LD   A,(DRIVE)
               LD   C,A
               LD   A,(IX+MDRV)

SSDRV2:        LD   (DRIVE),A

;SELECT DISC AND SIDE - GET PORT BASE ADDR IN A AND DSC VAR
;(224/228/240/244 FOR DISC 1 SIDE 1/2, DISC 2 SIDE 1/2)
L4718: ;*
SELD:          CALL TIRD         ;TEST IF RAM DISK
               RET  NC           ;RET IF IT IS

               CP   2
               LD   B,%11100000  ;224
               JR   NZ,SEL1

               LD   B,%11110000  ;240

SEL1:          LD   A,D
               AND  &80
               JR   Z,SEL2

               LD   A,%00000100

SEL2:          OR   B
               LD   (DSC),A
               RET

;CONVERT T/S IN D/E INTO NUMBER IN BC. A=C

CONM:          PUSH HL
               LD   H,0
               LD   L,D          ;HL=TRK
               LD   B,H
               LD   C,L
               ADD  HL,HL        ;*2
               ADD  HL,HL        ;*4
               ADD  HL,BC        ;*5
               ADD  HL,HL        ;*10
               LD   C,E
               DEC  C            ;0-9
               LD   A,D
               CP   4            ;NC IF ON TRACK 4 OR MORE
               ADC  HL,BC        ;HL=SECTOR (1 OR MORE,
                                 ; SUB 1 IF T4 OR MORE)
               ADD  HL,HL        ;FILE NUMBER DIV 2+1
               DEC  HL           ;FILE NUMBER DIV 2
               LD   C,(IX+RPTH)  ;DIR ENTRY 0/1
               ADD  HL,BC        ;FILE NUMBER
               LD   B,H
               LD   C,L
               LD   A,C
               POP  HL
               RET

;TEST FOR BUFFER FULL

TFBF:          CALL GRPNT
               LD   A,C
               CP   254
               RET  NZ
               LD   A,B
               DEC  A
               RET

;LOAD DATA BLOCK FROM DISK TO HL.

HLDBK:         LD   (PGES1),A
               EXX

;PGES1/DE=COUNT

LDBLK:         CALL SETF6
               JR   LBLOK

LDB1:          LD   A,(HL)
               CALL INCRPT
               LD   HL,(SVHL)
               LD   (HL),A
               INC  HL
               DEC  DE

LBLOK:         LD   (SVHL),HL
               CALL CKDE
               RET  Z

LDB2:          CALL TFBF
               JR   NZ,LDB1

               LD   (SVDE),DE
               LD   D,(HL)
               INC  HL
               LD   E,(HL)
               DI

LDB3:          CALL CCNT
               JP   C,LDB8

               INC  HL
               LD   (SVDE),HL
               CALL TIRD
               JR   C,LDB35

               CALL RDLB         ;RAM DISC
               JR   LDB3         ;DE=T/S, SVHL=OK

LDB35:         XOR  A
               LD   (DCT),A
               CALL SVNSR

LDB4:          CALL CTAS
               LD   C,DRSEC
               CALL SADC
               EXX
               LD   DE,2
               CALL GTBUF
               LD   A,(DSC)
               LD   (LDB6+1),A   ;STATUS PORT
               ADD  A,3
               LD   (LDB5+1),A   ;DATA PORT
               EXX
               LD   DE,510
               LD   HL,(SVHL)
               JR   LDB6

LDB5:          IN   A,(0)
               LD   (HL),A
               INC  HL
               DEC  DE
               LD   A,D
               OR   E
               JR   Z,LDB65

LDB6:          IN   A,(0)
               BIT  1,A
               JR   NZ,LDB5

               RRCA
               JR   C,LDB6

               AND  &0D
               JR   Z,LDB7

               CALL GTNSR
               CALL CDE1
               JR   LDB4

LDB65:         EXX
               JR   LDB6

LDB7:          LD   (SVHL),HL
               CALL GTBUF
               LD   D,(HL)
               INC  HL
               LD   E,(HL)
               JP   LDB3

LDB8: ;      * CALL LDINT
               CALL RSAD
               LD   DE,(SVDE)
               JP   LDB2

;CALCULATE COUNT

CCNT:          LD   HL,(SVDE)    ;BYTES LEFT IN THIS 16K BLOCK
               LD   BC,510
               SCF
               SBC  HL,BC
               RET  NC           ;RET IF 511 OR MORE LEFT

               LD   A,(PGES1)
               AND  A
               SCF                 ;SIGNAL "LAST SECTOR NOW"
               RET  Z            ;RET IF IT'S TRUE

               DEC  A
               LD   (PGES1),A    ;DEC PAGES TO DO
               LD   HL,(SVDE)
               LD   BC,16384
               ADD  HL,BC
               LD   (SVDE),HL    ;ADJUST REMAINDER
               JR   CCNT

;GET SCREEN MEMORY AND POINTER

GETSCR:        IN   A,(251)
               LD   (PORT1),A
               LD   A,(PORT2)
               OUT  (251),A
               LD   HL,(PTRSCR)
               RET

;PUT SCREEN MEMORY AND POINTER

PUTSCR:        LD   (PTRSCR),HL
               LD   A,(PORT1)
               OUT  (251),A
               RET

;REAL HOOK

HKSB:          LD   (PGES1),A
               EXX

;SAVE DATA BLOCK ON DISC. LEAVE LAST SECTOR IN BUFFER,
;PADDED WITH SPACES, BUT WITH RPT PTING PAST LAST REAL DATUM
;SAVE PGES1/DE BYTES FROM HL
L482D: ;*
HSVBK:         CALL DDEL   ;DELAY 0.25 SEC IN CASE OFSM TOO FAST
HSVB2:         CALL SBLOK
               RET  C          ;RET IF SINGLE SECTOR IN BUFFER,
                               ; PADDED
               CALL SVNSR      ;SAVE DE TO NEXT T/S VARIABLES
               LD   DE,(SVDE)  ;BYTES LEFT TO DO
               LD   HL,(SVHL)
               JR   HSVB2

;DELAY, SAVE BLOCK

DSVBL:         CALL DDEL     ;DELAY ABOUT 0.25 SEC IF NOT RAMDIC

;SAVE DATA BLOCK ON DISC
;SAVE PGES1/DE BYTES FROM HL, PADDING LAST SECTOR WITH SPACES

SVBLK:         CALL SBLOK
               JP   C,WSAD   ;WRITE SINGLE SECTOR

               JP   SVBF     ;OR LAST SECTOR OF BLOCK (HL=0)
L484C: ;*
SVB1:          LD   (HL),D
               CALL INCRPT   ;INC IX+RPTL/H
               LD   HL,(SVHL)
               INC  HL
               POP  DE
               DEC  DE       ;DEC BYTES-TO-SAVE

;TEST FOR ZERO BLOCK COUNT

SBLOK:         CALL SETF6
               CALL CKDE
               JP   Z,SVBSI  ;JP IF ALL BYTES NOW IN BUFFER
                             ;- PAD BUFFER AND SAVE IT
;SAVE CHAR IN REGISTER D

               PUSH DE
               LD   D,(HL)
L4861: ;*
Fix_L4861_4x:  ;4.2 & 4.3 Add    v2.2     v2.3    Src2.3      ;*
               LD   (SVHL),HL  ;L7C05   =       =             ;*
               CALL TFBF       ;L474C   =       =             ;*
               JR   NZ,SVB1    ;L484C   =       =             ;*
               POP  DE                                        ;*
               LD   (SVDE),DE  ;L7C02   =       =             ;*
               CALL FNFS       ;L496A   =       =             ;*
               LD   (HL),D                                    ;*
               INC  HL                                        ;*
               LD   (HL),E                                    ;*
               EX   DE,HL                                     ;*
               CALL SWPNSR     ;L6CAB / L6CAE / L6CA2_SWPNSR ?;*
               CALL TIRD       ;L5FDB = L5FDB / L5FCF_TIRD ?  ;*
               JP   NC,RDSB    ;L746E / L7473 / L7467_RDSB ?  ;*
;FixEnd                                                       ;*

               PUSH DE
               CALL GETSCR   ;SCREEN PAGED IN
               LD   HL,FTADD
;            * LD   (HL),D
               LD   BC,0     ;SECT COUNT=0
               JR   SVB2A
L488A: ;*
SVB2:          PUSH HL       ;SCRN PTR
               CALL CCNT     ;CHECK BYTES REMAINING, SUB 511
               EX   DE,HL
               POP  HL       ;SCRN PTR
               JR   C,SVB3   ;JR IF 510 OR LESS BYTES TO GO
                             ;- START SAVING SECTORS
               INC  DE       ;CORRECT FOR SUB 511 ->SUB 510 NOW
               LD   (SVDE),DE
 ;           * LD   D,(HL)
               CALL FFNS     ;FIND NEXT FREE SECT IN SAM
               LD   (HL),D
               INC  HL
               LD   (HL),E   ;STORE IN SCREEN TABLE
               INC  HL
  ;          * LD   (HL),D
               LD   BC,(SVCNT)
               INC  BC       ;INC SECTOR COUNT

SVB2A:         LD   (SVCNT),BC
               JR   SVB2

SVB3:          XOR  A
               LD   (HL),A
               LD   D,H
               LD   E,L
               INC  DE
               LD   BC,&01FF ;ALLOW UP TO 512 BYTES OF ZEROS
               LDIR          ; FOR LAST SECTOR

               LD   HL,FTADD
               CALL PUTSCR   ;SET PTR VAR TO START OF LIST
                             ;OF T/S, SWITCH SCREEN OFF
               POP  DE
               CALL WSAD     ;WRITE FIRST SECTOR, SETUP DSC
               DI
               XOR  A
               LD   (DCT),A  ;ZERO ERRORS
               CALL GTNSR    ;GET T/S IN D/E
               EXX
               LD   HL,(PTRSCR)  ;SOURCE OF LAST 2 BYTES
               EXX                 ; (OR MORE, IF LAST SECT)

;LOOP HERE FOR EACH SECTOR IN SVCNT

SVB4:          LD   HL,(SVCNT)
               LD   A,H
               OR   L
               RET  Z          ;RET IF LAST TO DO NEXT. NC

               BIT  7,H
               RET  NZ         ;RET IF NEGATIVE (FINAL SECTOR
L48D3: ;*
SVBF:          DEC  HL
               LD   (SVCNT),HL
               PUSH DE         ;T/S
               CALL CTAS       ;CONFIRM TRACK
               CALL PRECMX     ;SEND "WRITE SECTOR" CMD
               LD   A,(PORT2)  ;SCREEN PORT
               EX   AF,AF'     ;STORED IN A' FOR SPEED
               LD   A,(DSC)    ;STATUS PORT
               LD   (SVB6+1),A ;SELF-MOD CODE
               ADD  A,3
               LD   C,A        ;DATA PORT
               EXX
               LD   C,A        ;C'=DATA PORT TOO
               LD   DE,&0300   ;DE IS A LARGE NUMBER SO
               EXX
               LD   DE,&02FE   ;510 BYTE COUNTER
               LD   HL,(SVCNT)
               INC  HL
               LD   A,H
               OR   L
               JR   NZ,SVB45   ;JR IF NOT LAST SECT

               LD   DE,(SVDE)  ;REMAINING BYTES
               LD   A,E
               AND  A
               JR   Z,SVB45    ;NO INC IF E=0 SO COUNTER OK

               INC  D

SVB45:         LD   HL,(SVHL)  ;SRC
               JR   SVB6

SVB5:          OUTI
               DEC  E        ;LSB OF BYTE COUNT
               JR   Z,SVBL2
L490F: ;*
SVB6:          IN   A,(0)    ;SELF-MOD STATUS PORT
               BIT  1,A
               JR   NZ,SVB5  ;JR IF DISK READY FOR BYTE
               RRCA
               JR   C,SVB6   ;JR IF BUSY - DATA STILL TO BE SENT
               EXX             ;HL=MAIN PTR
               POP  DE       ;T/S
               AND  &0D
               JR   NZ,SVB7  ;JR IF ERROR

;SCREEN IS STILL IN

               LD   (SVHL),HL ;UPDATE SRC PTR (IN 510 BYTE STEPS

               EXX
               DEC  HL
               LD   E,(HL)
               DEC  HL
               LD   D,(HL)    ;D/E=BYTES JUST SAVED (NEXT T/S)
               INC  HL
               INC  HL
               LD   A,(PORT1)
               OUT  (251),A   ;NON-SCREEN
               PUSH DE
               EXX

               POP  DE        ;T/S
               JP   SVB4      ;LOOP FOR ALL SECTORS

;ENTERED WHEN AN ERROR OCCURRED DURING SAVE

SVB7:          LD   A,(PORT1)
               OUT  (251),A   ;NON-SCREEN
               CALL CDE1      ;CHECK/SET T/S
               EXX
               DEC  HL
               DEC  HL
               EXX
               JP   SVB4      ;TRY AGAIN - MAX 10 ERRORS IN
                              ; ENTIRE BLOCK
;ENTERED WHEN E COUNTS TO 0

SVBL2:         DEC  D         ;DEC MSB OF COUNT (PRE-INCED)
               JR   NZ,SVB6

               EXX              ;SWITCH TO HL' - PTS TO T/S OR 0
               EX   AF,AF'
               OUT  (251),A   ;SCREEN ON ???
               JR   SVB6

;PAD A SINGLE SECTOR FOR A SMALL BLOCK

SVBSI:         CALL GRPNT
SVBS1:         LD   A,B
               SUB  2
               OR   C
               JR   Z,SVBS3

SVBS2:         LD   (HL),0   ;PAD WITH ZEROS
               INC  HL       ;NO CHANGE OF RPT
               INC  BC
               JR   SVBS1

SVBS3:         CALL GTNSR    ;GET T/S
               SCF             ;"SINGLE SECTOR"
               RET


;--------------
Part_C12:


;FAST FIND NEXT SCT IF FNFS USED BEFORE

FFNS:          LD   DE,(FFDE)
               PUSH HL
               PUSH BC
               LD   HL,(FFHL)
               JR   FNS2

;FIND NEXT FREE SECTOR

FNFS:          LD   DE,&0401   ;TRK 4, SCT 1
               PUSH HL
               PUSH BC
               LD   HL,SAM-1

FNS1:          INC  HL

FNS2:          LD   A,(HL)
               INC  A
               JR   NZ,FNS3    ;JR IF NOT FF (=ALL USED)

               LD   A,E
               ADD  A,8
               LD   E,A
               SUB  11
               INC  A          ;EQUIV. OF SUB 10, BUT CY IF <=0
               JR   C,FNS1     ;JR IF SECTOR WAS OK BEFORE SUB

               LD   E,A
               CALL FNS5       ;NEXT TRACK
               JR   FNS1

FNS3:          LD   B,1        ;MASK FOR BIT 0
               LD   (FFDE),DE
               LD   (FFHL),HL

FNS4:          LD   A,(HL)
               AND  B
               JR   Z,FNS6     ;JR IF FND FREE

               CALL ISECT
               CALL Z,FNS5     ;NEXT TRK IF NEEDED
               RLC  B
               JR   FNS4

;CALLED BY POINT SR "FITS"

FNS5:          INC  D
               CALL TSTD       ;GET TRK LIMIT
               LD   C,A
               LD   A,(DTKS)
               SUB  4
               JR   NC,FNS55   ;JR IF NORMAL DISC
               NEG
               ADD  A,C        ;FIDDLE LIMIT TO ALLOW
                               ;E.G TRK 42 ON 40 TRK, 1 DTK
               LD   C,A        ;RAM DISC SO TKS 4-42 OK.

FNS55:         LD   A,C        ;LIMIT
               CP   D          ;CP LIMIT, TRK
               JP   Z,REP24    ;NOT ENOUGH SPACE

               CALL TIRD
               JR   C,FNS56    ;JR IF NOT RAM DISC

               LD   A,&50      ;SIDE 1 ON RAM DISC=80 TRK
               DEFB &FE        ;"JR+1"

FNS56:         LD   A,C
               AND  &7F
               CP   D
               RET  NZ

               LD   D,&80
               RET

;FOUND FREE SECT

FNS6:          LD   A,(HL)
               OR   B
               LD   (HL),A     ;SET BIT IN SAM
               LD   C,L        ;L=DISP FROM START OF SAM,+0FH
               LD   A,B
               LD   B,0
               PUSH IX
               ADD  IX,BC
               OR   (IX+FSAM-&0F)
               LD   (IX+FSAM-&0F),A
               POP  IX
               CALL ICNT       ;INC SECTS USED COUNT
               POP  BC
               POP  HL
               RET

;TEST TRACKS ON DISC

TSTD:          PUSH HL
               LD   HL,TRAKS1
               CALL TIRD
               JR   C,TSD0

               CALL RTSTD      ;RAMDISC
               POP  HL
               RET

TSD0:          DEC  A
               JR   Z,TSD1

               INC  HL
TSD1:          LD   A,(HL)
               POP  HL
               RET

;PRINT FILE NAME

PFNME:         LD   B,1
               CALL GRPNTB

;PRINT FILE NAME FROM HL. USES HL,A,B
;ADVANCES HL BY 10

PFNM0:         PUSH HL
               LD   BC,&0A0A

PFNM1:         LD   A,(FNSEP) ;FILE NAME SEPARATOR - USUALLY "."
               CP   (HL)
               LD   A,(HL)
               JR   NZ,PFNM12

               LD   A,B
               CP   5
               JR   NC,PFNM15  ;PRINT PAD SPACE IF "." AND CHARS
                               ;STILL TO PRINT >=5
               LD   A,(FNSEP)

PFNM12:        INC  HL
               CP   &20
               JR   NZ,PFNM2

PFNM15:        LD   A,(CHDIR)

PFNM2:         CALL PNT
               DJNZ PFNM1

               POP  HL
               ADD  HL,BC      ;BC=10
               RET

;FILE DIR HAND/ROUT.

FDHR:          DI
               PUSH AF
               CALL GTIXD
               POP  AF
               LD   (DCHAN+4),A
               XOR  A
               LD   (FSLOT),A  ;NO FREE SLOT
               LD   (MAXT),A   ;MAX TAG=0
               CALL REST
               CALL ORSAD
               CALL SDTKS      ;SET DIR TKS, CHECK RAND NO
               PUSH HL
               BIT  2,(IX+4)
               CALL NZ,PDIRH   ;PRINT DIR HDR NOW THAT PATH$
                               ;CHECKED, IF COMPLEX DIR WANTED
               LD   A,(DTKS)
               SUB  5
               JR   C,FDH05    ;JR IF SAM OK (4 DIR TRAC

               INC  A
               LD   B,A        ;B=EXTRA DTKS (1-35)
               ADD  A,A
               ADD  A,A
               ADD  A,B        ;A=5-175 (EXTRA DIR SECTS/2)
               LD   C,A
               RRA
               RRA               ;EXTRA DIR SECTS/8=BYTES
               AND  &3F        ;A=BYTES TO MARK FF IN SAM (1-43)
               LD   B,A
               LD   HL,SAM
               LD   A,(HL)
               OR   &FE
               LD   (HL),A     ;T4,S1 NOT RESERVED - BUT KEEP
               DEC  B          ; CURRENT STATUS!
               JR   Z,FDH03

FDH02:         INC  HL
               LD   (HL),&FF
               DJNZ FDH02

FDH03:         LD   A,C
               ADD  A,A
               AND  &07        ;EXTRA BITS (0,2,4 OR 6)
               JR   Z,FDH05

               LD   B,A
               XOR  A

FDH04:         SCF
               RLA
               DJNZ FDH04

               INC  HL
               OR   (HL)
               LD   (HL),A

FDH05:         POP  HL
               JR   FDH25

FDH1:          CALL ORSAD
FDH2:          CALL POINT

FDH25:         LD   A,(HL)
               AND  A
               JP   Z,FDHF     ;JP IF FREE SLOT

;TEST FOR 'P' NUMBER

               LD   A,(DCHAN+4)
               RRA
               JR   NC,FDH3    ;JR IF NOT E.G. LOAD 3

               CALL CONM       ;CONVERT T/S TO P NUMBER IN BC
               LD   A,(FSTR1)  ;LSB
               CP   C

FDHDH:         JP   NZ,FDHd

               LD   A,(TEMPW4+1) ;MSB
               CP   B
               JR   NZ,FDHDH
               RET

;GET TRACK AND SECTOR COUNT

FDH3:          AND  &03
               JP   Z,FDH9     ;JP IF NEITHER TYPE OF DIRECTORY

               LD   A,(HL)     ;TYPE
               EX   AF,AF'
               LD   B,H
               LD   C,L
               INC  H
               DEC  HL
               DEC  HL
               LD   A,(HL)     ;DIRECTORY TAG - BYTE FE
               LD   HL,11
               ADD  HL,BC

               LD   B,(HL)
               INC  HL
               LD   C,(HL)
               LD   (SVBC),BC ;SECTORS USED BY FILE
               LD   HL,(CNT)
               ADD  HL,BC
               LD   (CNT),HL  ;TOTAL USED SECTORS SO FAR ON DISC
               LD   HL,(TCNT)
               INC  HL
               LD   (TCNT),HL ;TOTAL FILES SO FAR ON DISC
               LD   B,A
               LD   A,(CDIRT) ;CUR DIR
               CP   &FF
               JR   Z,FDH35   ;JR IF "ALL"

               CP   B
               JP   NZ,FDHd   ;JP IF NOT FILE FROM CURRENT DIR

FDH35:         LD   HL,(FCNT)
               INC  HL
               LD   (FCNT),HL ;TOTAL FILES SO FAR IN TH

;TEST IF WE SHOULD PRINT NAME

               EX   AF,AF'
               RLA               ;BIT 7 OF FILE TYPE=1 IF "HIDDE
               JR   C,FDHH     ;JR IF HIDDEN

               LD   A,(NSTR1+2) ;FF IF NO NAME ASKED FOR
               INC  A
               JR   Z,FDH4

               CALL CKNAM
               JR   NZ,FDHH   ;JR IF NOT A DIR OF THIS NAME TYPE

FDH4:          BIT  1,(IX+4)
               JR   Z,FDH5     ;JR IF COMPLEX DIR,
                               ;CONTINUE IF SIMPLE

;COPY FILE NAME FROM SECT BUFFER TO SPECIAL BUFFER OF NAMES

               PUSH DE
               LD   A,(DCHAN+RPTH) ;GET DIR ENTRY (0 OR 1)
               LD   B,A
               LD   C,1
               LD   HL,(DCHAN+BUFL);GET SECT BUFFER ADDR IN HL
               ADD  HL,BC          ;PT TO FILE NAME IN DIR ENTRY
               LD   DE,(PTRSCR)
               LD   BC,10
               LDIR
               LD   (PTRSCR),DE
               POP  DE

FDHH:          JP   FDHd

FDH5:          CALL POINT
               LD   A,(HL)
               BIT  6,A
               JR   Z,FDH4A  ;JR IF FILE NOT PROTECTED - PRT NUM
                             ;ELSE PRINT " * "
               CALL SPC
               LD   A,"*"
               CALL PNT
               CALL SPC
               JR   FDH4D    ;JR TO PRINT FILE NAME

;PRINT FILE NUMBER AS 2 OR 3 DIGITS E.G. "12 " OR "123"

FDH4A:         CALL CONM     ;GET FILE NUMBER IN BC
               PUSH DE
               LD   HL,99
               AND  A
               SBC  HL,BC
               LD   H,B
               LD   L,C
               LD   A,&20
               JR   C,FDH4B  ;JR IF >99

               CALL PNUM2
               CALL SPC
               AND  A        ;NC

FDH4B:         CALL C,PNUM3
               POP  DE

FDH4D:         CALL PFNME    ;PRINT FILE NAME FOR DIR 1/2

;PRINT SECTOR COUNT AS E.G. " 20 "

               PUSH DE
               CALL POIDFT
               JR   Z,FDH85  ;JR IF DIR

               PUSH AF
               LD   A,&20
               LD   BC,(SVBC) ;SECTORS USED
               LD   HL,999
               AND  A
               SBC  HL,BC    ;CY IF >999 SECTS USED
               LD   H,B
               LD   L,C
               JR   C,FDH84

               CALL PNUM3
               CALL SPC
               AND  A

FDH84:         CALL C,PNUM4  ;USE 4 DIGITS, BUTTING ONTO TYPE,
                             ;IF NEEDED E.G. 1234OPENTYPE

;PRINT TYPE OF FILE E.G. "BASIC ","C "

               POP  AF

FDH85:         CALL PNTYP
               POP  DE
               JR   FDHd

;TEST FOR SPECIFIC FILE NAME

FDH9:          LD   A,(DCHAN+4)
               AND  &18
               JR   Z,FDHd   ;JR IF NEITHER BIT 3 NOR 4 ARE SET

               LD   A,(HL)
               AND  &1F
               CP   DFT      ;DIR
               JR   NZ,FDH92

               LD   BC,DIRT
               ADD  HL,BC    ;PT TO DIR TAG
               LD   A,(MAXT) ;CURRENT MAX
               CP   (HL)     ;CP TAG OF THIS DIR
               JR   NC,FDH91 ;JR IF NO NEW MAX TAG

               LD   A,(HL)
               LD   (MAXT),A
               AND  A

FDH91:         SBC  HL,BC    ;PT TO DIR TAG FOR FILE

FDH92:         INC  H
               DEC  HL
               DEC  HL
               LD   A,(CDIRT)
               CP   (HL)
               JR   Z,FDH95  ;JR IF RIGHT DIR

               INC  A
               JR   NZ,FDHd  ;JR UNLESS 0FFH "ANY" DIR

FDH95:         CALL CKNAM
               RET  Z

;CALCULATE NEXT DIRECTORY ENTRY

FDHd:          LD   A,(DCHAN+RPTH)
               DEC  A
               JR   Z,FDHe   ;JR IF WE HAVE JUST DONE SECOND DIR
                             ;ENTRY IN SECT
               CALL CLRRPT
               INC  (IX+RPTH)  ;NEXT ENTRY
               JP   FDH2

FDHe:          CALL ISECT    ;NEXT SECTOR
               JP   NZ,FDH1  ;JP IF NOT LAST SECT IN TRACK

               INC  D        ;INC TRACK
               LD   A,D
               CP   4
               JR   NZ,FDHE2

               INC  E        ;SECT=2 IF TRACK=4 (SKIP DOS)

FDHE2:         LD   A,(DTKS)
               CP   D
               JP   NZ,FDH1

               AND  A        ;NZ
               RET

;FOUND FREE DIRECTORY SPACE - SEE IF THAT IS WHAT IS WANTED

FDHF:          CALL FSLSR
               LD   A,(IX+4)
               CPL
               BIT  6,A
               RET  Z        ;RET IF GOT WHAT WE WANTED

               INC  HL
               LD   A,(HL)
               AND  A        ;**
               JR   NZ,FDHd  ;JR IF DIR ENTRY HAS BEEN USED AT
                             ;SOME TIME (NAME START NOT CHR$ 0)
               INC  A        ;NZ
               RET             ;USED PART OF DIR DEALT WITH

FSLSR:         LD   A,(FSLOT)
               AND  A
               RET  NZ         ;JR IF WE HAVE ALREADY NOTED A
                               ;"FIRST FREE SLOT"
               LD   (FSLOT),DE
               LD   A,(DCHAN+RPTH)
               LD   (FSLTE),A  ;RECORD 0 OR 1 FOR WHICH ENTRY
               RET               ; IN SECTOR IS FREE

;CHECK FILE NAME IN DIR
;EXIT: Z IF MATCHED

CKNAM:         PUSH DE
               CALL POINT
               LD   B,11
               LD   DE,NSTR1
               BIT  3,(IX+4)
               JR   Z,CKNM2  ;JR IF TYPE IRREL - SKIP TYPE

CKNM1:         LD   A,(DE)
               CP   "*"
               JR   Z,CKNM5  ;MATCH ANY LEN IF "*"

               CP   "?"
               JR   Z,CKNM2  ;MATCH ANY CHAR IF "?"

               XOR  (HL)
               AND  &DF
               JR   NZ,CKNM4

CKNM2:         INC  DE
               INC  HL
               DJNZ CKNM1

CKNM3:         XOR  A        ;MATCH=Z
CKNM4:         POP  DE
               RET

CKNM5:         INC  DE
               INC  HL
               LD   A,(DE)
               CP   "."
               JR   NZ,CKNM3 ;E.G. "*" OR "A*" MATCHES ANYTHING
                             ;FROM "*" ON, BUT E.G. "*.Z" LOOKS
                             ;FOR MATCHING "." IN FILE NAME
CKNM6:         LD   A,(HL)
               CP   "."
               JR   Z,CKNM2  ;IF MATCHING "." FOUND,
                             ; CHECK SUFFIX VS REQUESTED
               INC  HL
               DJNZ CKNM6

               POP  DE
               RET             ;NO MATCH

;CLEAR SAM

CLSAM:         XOR  A
               LD   HL,SAM
               LD   B,195

CLSML:         LD   (HL),A
               INC  HL
               DJNZ CLSML
               RET

;RND/SERIAL OPEN FILE SECTOR ADDRESS MAP - INC SAMCNT
;OF LONG-TERM OPEN FILES (CHANS BUFFERS)

ROFSM:         CALL OFSM
               LD   A,(SAMCNT)
               INC  A
               LD   (SAMCNT),A
               DEC  A
               JR   NZ,SOFS2   ;JR IF NOT FIRST FILE

               CALL GRWA       ;GET CUR. RAND WORD ADDR IN HL
               LD   C,(HL)
               INC  HL
               LD   B,(HL)
               LD   (SAMRN),BC ;RND NO. OF DISC USED FOR SAM
               LD   A,(DRIVE)
               LD   (SAMDR),A

SOFS2:         XOR  A          ;NC - OK
               RET

GOFSM:         CALL GTIXD

;OPEN FILE SECTOR ADDRESS MAP

OFSM:          PUSH IX
               LD   A,(SAMCNT)
               AND  A
               CALL Z,CLSAM  ;NO CLEAR IF OPEN-TYPE FILE(S) HAS
                             ;SAM IN USE
               LD   A,&30  ;BIT 5=CREATE SAM,BIT 4=TEST FOR NAME
                           ;(BETTER THAN G+DOS BECAUSE EVEN IF
                           ;DISC CHANGED SECTORS USED IN EITHER
                           ;DISC ARE SET BITS IN SAM. G+DOS
                           ;USED LD A,10H - NO ORING OF BITS)
               CALL FDHR
               JR   NZ,OFM4  ;JR IF NAME NOT FOUND

;FILE NAME ALREADY USED

OFM25:         PUSH DE
               CALL BITF4
               JP   NZ,REP28 ;"FILE NAME USED" IF "OPEN DIR"

               CALL POIDFT
               JP   Z,REP28  ;OR IF DIR FILE

               CALL NRRD
               DEFW OVERF

               AND  A
               JR   Z,OFM3   ;NO "OVERWRITE? Y/N" IF SAVE OVER
                             ;WANTED
               PUSH HL
               CALL SKSAFE   ;IN CASE "N"
               POP  HL
               BIT  6,(HL)
               JP   NZ,REP36 ;"PROTECTED FILE"

               CALL PMO5     ;"OVERWRITE"
               CALL FNM7K    ;FILE NAME, "Y/N", KEY
               JR   Z,OFM29  ;JR IF "Y"

               POP  DE
               POP  IX
               SCF             ;SIGNAL ERROR
               RET             ;AND ABORT

OFM29:         CALL DWAIT    ;IN CASE STOPPED

OFM3:          CALL DDEL
               CALL POINT
               LD   BC,&000F
               LD   (HL),B   ;MARK ENTRY IN BUFFER AS "DELETED"
               ADD  HL,BC    ;POINT TO SAM OF FILE IN BUFFER
               LD   DE,SAM   ;POINT TO BAM FOR DISC SO FAR
               LD   B,195

DBAML:         LD   A,(DE)
               XOR  (HL)     ;REVERSE BITS ORED IN BY FDHR
               LD   (DE),A   ;UPDATE BAM TO TAKE ACCOUNT OF
               INC  E        ; ERASURE
               INC  HL
               DJNZ DBAML

               POP  DE
               CALL FSLSR
               CALL WSAD     ;WRITE DIR SECT TO DISC
               LD   IX,DCHAN
               CALL FDH1     ;COMPLETE BAM BY SCANNING REST OF
                             ;DIRECTORY
               JR   Z,OFM25  ;JR IF SECOND OR OTHER COPY...
                             ;SHOLD NEVER BE!

;CLEAR AREA FOR NEW DIR ENTRY IN BUFFER

OFM4:          POP  IX
               PUSH IX
               LD   B,0

OFM5:          LD   (IX+FFSA),0
               INC  IX
               DJNZ OFM5

               POP  IX
               PUSH IX
               LD   HL,NSTR1
               LD   B,11
               CALL OFM6
               POP  IX
               PUSH IX
               LD   BC,220
               ADD  IX,BC
               LD   HL,UIFA+15
               LD   B,48-15
               CALL OFM6
               POP  IX

               CALL BITF4
               CALL Z,FNFS     ;AVOID FNFS IF EXISTING OPENTYPE
                               ;FILE, OR IF "OPEN DIR"
               CALL SVNSR
               LD   (IX+FTRK),D ;FIRST TRACK
               LD   (IX+FSCT),E ;AND SECTOR OF FILE IN DIR
               XOR  A           ;NC=OK
               JR   CLRRPT

;GET BUFFER ADDRESS

GTIXD:         LD   IX,DCHAN
               LD   HL,DRAM
               LD   (BUF),HL

;CLEAR RAM POINTER

CLRRPT:        LD   (IX+RPTL),0
               LD   (IX+RPTH),0
               RET

OFM6:          LD   A,(HL)
               LD   (IX+FFSA),A
               INC  HL
               INC  IX
               DJNZ OFM6
               RET

BEEP:          PUSH HL
               PUSH DE
               PUSH BC
               PUSH IX
               LD   HL,&036A
               LD   DE,(BEEPT)
               CALL CMR
               DEFW BEEPR

               POP  IX
               POP  BC
               POP  DE
               POP  HL
               RET

;SPECIAL CFSM USED WHEN SEVERAL BLOCKS MAY BE WRITTEN TO SAME
;FILE AND LAST SECTOR NEEDS SAVING BECAUSE HSVBL USED.

SCFSM:         CALL GTNSR
               CALL WSAD     ;LAST SECTOR

;CLOSE FILE SECTOR MAP

CFSM:          PUSH IX
               LD   DE,(FSLOT)
               INC  E
               DEC  E
               JR   Z,CLOSX  ;DIR PROBABLY FULL - BUT DO A FULL
                             ;CHECK IN CASE WE ARE CLOSING AN
                             ;"OUT" FILE (FSLOT HAS BEEN SET TO
                             ;ZERO BECAUSE IT MIGHT HAVE BEEN
                             ;OBSOLETE IF SAVE USED SINCE OPEN)
               CALL RSAD     ;READ DIR SECTOR WITH FREE SLOTS
               LD   A,(FSLTE)
               LD   (DCHAN+RPTH),A ;PT RPT TO FREE SLOT
               JR   NCF25

;JPED TO FROM CLOSP

CLOSX:         LD   A,&40    ;FIND FREE SLOT
               CALL FDHR
               JP   NZ,REP25 ;DIRECTORY FULL

;UPDATE DIRECTORY
;JPED TO FROM CLOSP

NCF25:         CALL POINT
               LD   (SVIX),IX
               POP  IX
               PUSH IX

               LD   A,E
               DEC  A
               OR   D
               OR   B        ;OFFSET FROM 'POINT'
               LD   B,0      ;B=COUNT FOR 256 BYTES
               JR   NZ,CFM3  ;JR IF NOT T0,S1,ENTRY1

               LD   BC,&0D20A ;C=0AH
               CALL CFMC     ;COPY BYTES 0-D1 FROM BUFFER
               ADD  HL,BC    ;SKIP DISC NAME - KEEP AS READ FROM
               ADD  IX,BC    ;(BC=000AH)
               LD   B,&20
               CALL CFMC     ;COPY BYTES DC-FB FROM BUFFER, ONLY
                             ;COPY UP TO FBH ON 1ST DIR ENTRY
                             ;KEEP RAND WORD, DIR TAG AND 'EXTRA
                             ;DIR TRACKS' BYTE AS READ FROM DISC
               LD   A,(IX+FFSA+2) ;READ DIR TAG
               INC  HL
               INC  HL
               LD   (HL),A   ;USE DIR TAG FROM BUFFER, NOT DISC
               CP   A        ;Z

CFM3:          CALL NZ,CFMC
               LD   IX,(SVIX)
               CALL DATSET   ;COPY DATE/TIME TO +F2 TO +FB
               CALL POINT
               INC  H
               DEC  HL
               DEC  HL
               LD   A,(CDIRT)
               LD   (HL),A   ;TAG FILE WITH DIRECTORY CODE
               CALL WSAD
               CALL SKSAFE
               POP  IX
               RET

CFMC:          LD   A,(IX+FFSA)
               LD   (HL),A
               INC  IX
               INC  HL
               DJNZ CFMC
               RET

;GET A FILE FROM DISC

GTFLE:         LD   A,(FSTR1)

;TEST FOR FILE NUMBER

               INC  A
               JR   Z,GTFL3  ;JR IF NOT 'LOAD N' (FSTR1=FF)

               LD   A,1      ;"NUMBERED FILE"
               CALL FDHR     ;FIND FILE
               JP   NZ,REP26 ;ERROR IF DOESNT EXIST

GTFLX:         CALL GTFSR    ;FIDDLE ZX TYPE CODES
               CP   21
               JP   NC,REP13 ;ERROR IF DIR OR OPENTYPE

               SUB  &12
               ADC  A,0
               JP   Z,REP13  ;ERROR IF TYPE 11H OR 12H

               LD   DE,NSTR1
               LD   BC,11
               LDIR
               LD   B,4
               XOR  A
               CALL LCNTB

               CALL B211

               CALL POINT
               LD   DE,UIFA
               LD   BC,11
               LDIR

               LD   B,15
               CALL LCNTA

               LD   B,48-26
               LD   A,&FF
               CALL LCNTB
               JR   GTFL4

GTFL3:         LD   A,&10      ;NAMED FILE, IGNORE TYPE
               CALL FDHR
               JP   NZ,REP26

GTFL4:         CALL GTFSR
               LD   A,(NSTR1)
               LD   C,A
               SUB  20
               ADC  A,0
               JR   NZ,GTFL5A  ;JR IF DESIRED TYPE <> CODE OR
                               ;SCREEN$
               LD   A,(HL)
               SUB  20
               ADC  A,0
               JR   NZ,GTFL5A  ;JR IF LOADED TYPE <>CODE/SCREEN$

               LD   (HL),C     ;LOADED TYPE=DESIRED TYPE

GTFL5A:        LD   A,C
               CP   (HL)
               JP   NZ,REP13   ;WRONG FILE TYPE IF NO MATCH

               CALL BITF7

;CALLED BY COPY - Z SET (F7 USED FOR OTHER THINGS)

GTFL2:         PUSH AF         ;"FLAG 7" STATUS - NZ IF ZX
               LD   DE,HD002
               CALL B211       ;211-219 ZX FILE DATA TO HD002

               DEC  HL
               LD   DE,STR-30
               LD   C,42
               LDIR            ;210-251 WITH SNP REGS AT STR-20

               CALL NMMOV
               POP  AF
               JR   Z,GTFL7    ;JR IF NOT ZX FILE

               LD   B,11
               CALL LCNTA

               LD   B,48-26
               LD   A,&FF
               CALL LCNTB

               LD   HL,(HD0B2) ;LEN
               XOR  A
               CALL PAGEFORM
               RES  7,H
               LD   (DIFA+34),A
               LD   (PGES2),A
               LD   (DIFA+35),HL
               LD   (HD0B2),HL

               LD   HL,(HD0D2) ;ADDR
               XOR  A
               CALL PAGEFORM
               DEC  A
               AND  &1F
               LD   (DIFA+31),A
               LD   (PAGE2),A
               LD   (DIFA+32),HL
               LD   (HD0D2),HL
               JR   GTFL8

GTFL7:         LD   B,220
               CALL GRPNTB
               LD   BC,48-15
               LDIR

GTFL8:         LD   B,13
               CALL GRPNTB
               LD   D,(HL)
               INC  HL
               LD   E,(HL)
               LD   (SVDE),DE
               RET

NMMOV:         CALL POINT
               LD   DE,DIFA
               LD   BC,11
               LDIR            ;TYPE/NAME
               LD   B,4

;LOAD (DE) WITH COUNT B

LCNTA:         LD   A,&20

LCNTB:         LD   (DE),A
               INC  DE
               DJNZ LCNTB
               RET

GTFSR:         CALL POIDFT
               CP   &07
               JR   NZ,GTFS1
               LD   A,&04

GTFS1:         CP   &10
               JR   NC,GTFS2   ;JR IF SAM FILE

               DEC  A
               OR   &10
               CALL SETF7      ;"ZX FILE"

GTFS2:         LD   (HL),A
               RET

;READ SECTOR SR

RSSR:          CALL CTAS
               DI
               LD   C,DRSEC
               CALL SADC
GTBUF:         LD   L,(IX+BUFL)
               LD   H,(IX+BUFH)
               RET

;FIND A FILE ROUTINE
FINDC:         LD   A,&10
               CALL FDHR

;GET RAM AND POINTER TO BUFFER

POINT:         LD   B,0

GRPNTB:        LD   (IX+RPTL),B

GRPNT:         LD   L,(IX+BUFL)
               LD   H,(IX+BUFH)
               LD   B,(IX+RPTH)
               LD   C,(IX+RPTL)
               ADD  HL,BC
               RET

;GET THE NEXT TRACK/SECTOR

GTNSR:         LD   D,(IX+NSRH)
               LD   E,(IX+NSRL)
               RET

POIDFT:        CALL POINT
               LD   A,(HL)
               AND  &1F
               CP   DFT        ;DIR
               RET

B211:          LD   B,211
               CALL GRPNTB
               LD   BC,9
               PUSH HL
               LDIR
               POP  HL
               RET


;--------------
Part_D1:


;CHECK IS IT END OF LINE

CIEL:          CALL GCHR
               CP   &0D        ;CR
               RET  Z

               CP   &3A        ;":"
               RET

;CHECK FOR SYNTAX ONLY

CFSO:          LD   (SVA),A
               CALL NRRD
               DEFW FLAGS

               AND  &80
               LD   A,(SVA)
               RET

;CHECK FOR END OF SYNTAX

CEOS:          CALL CIEL
               JR   NZ,REP0HC

ABORT:         CALL CFSO
               RET  NZ

;END OF STATEMENT

ENDS:          LD   E,0        ;NO ACTION
               DEFB &21        ;"JR+2"

ENDSX:         LD   E,1

END1:          XOR  A
               CALL NRWR
               DEFW XPTR+1     ;NO ERROR

               CALL BCR
               LD   (NRFLG),A  ;RECURSE OK
               LD   SP,(ENTSP)

;FROM END OF HOOKS

RENT:          POP  HL
               LD   (ENTSP),HL ;OLD VALUE
               RET

;TEST FOR BREAK ROUTINE

BRKTST:        LD   A,&F7
               IN   A,(249)
               AND  &20
               RET  NZ

REP3:          CALL DERR
               DEFB 84

;INSIST SEPARATOR C

ISEPX:         CALL GTNC

ISEP:          CP   C

REP0HC:        JP   NZ,REP0    ;ERROR OR SKIP CHAR

;GET THE NEXT CHAR.

GTNC:          CALL CMR
               DEFW &0020

               RET

;GET THE CHAR UNDER POINTER

GCHR:          CALL CMR
               DEFW &0018

               RET

;CY IF A LETTER OR DIGIT

ALPHANUM:      CALL ALPHA
               RET  C

;CY IF A DIGIT

NUMERIC:       CP   "9"+1      ;NC IF TOO HIGH
               RET  NC

               CP   "0"
               CCF
               RET

;READ A DOUBLE ROM WORD IN BC

NRRDD:         EX   (SP),HL
               PUSH DE
               CALL GTHL
               PUSH DE
               CALL RDBC
               JR   PPXR

;READ A ROM SYSTEM VARIABLE (PARAM WORD) TO A REG

NRRD:          EX   (SP),HL
               PUSH DE
               CALL GTHL
               PUSH DE
               CALL RDA
               JR   PPXR

;WRITE A DOUBLE WORD FROM A ROM SYS VAR (PARAM WORD) TO BC

NRWRD:         EX   (SP),HL
               PUSH DE
               CALL GTHL
               PUSH DE
               CALL WRTBC
               JR   PPXR

;WRITE ROM SYSTEM VARIABLE

NRWR:          EX   (SP),HL
               PUSH DE
               CALL GTHL
               PUSH DE

;WRITE A TO (HL) IN SYS PAGE

               LD   E,A
               IN   A,(251)
               LD   D,A
               XOR  A
               OUT  (251),A
               SET  7,H
               RES  6,H
               LD   (HL),E
               LD   A,D
               OUT  (251),A
               LD   A,E

PPXR:          POP  HL
               POP  DE
               EX   (SP),HL
               RET

;PLACE "NEXTST" ADDR - CALLED BY "OPEN"
;IN CASE NEXTSTAT NOT ON STACK, PLACE IT. TRANSFORMS:
;NEXT STAT/ERRSP TO FOWIA/NEXTSTAT/ERRSP
;   *                *             (*=SP)
;ERRSP TO FOWIA/NEXTSTAT/ERRSP
;  *               *             (*=SP)

PLNS:          LD   HL,(ENTSP)
               INC  HL
               INC  HL
               INC  HL
               INC  HL       ;MAIN ROM STACK PTR ON DOS STACK
               DEC  (HL)
               DEC  (HL)     ;DEC STORED SP LSB
                        ;(ADD ONE ADDR TO ENSURE INTS DON'T HIT)
               CALL NRRDD
               DEFW ERRSP    ;READ ERRSP INTO BC

               LD   H,B
               LD   L,C
               LD   BC,(NEXTST)
               CALL DWRBC    ;DEC HL TO PT TO ADDR BELOW ERRSP,
               DEC  HL       ; WRITE BC
               LD   BC,FOWIA ;NULL

DWRBC:         DEC  HL
               DEC  HL

;WRITE BC TO (HL) IN SYS PAGE

WRTBC:         IN   A,(251)
               EX   AF,AF'
               XOR  A
               OUT  (251),A
               SET  7,H
               RES  6,H
               LD   (HL),C
               INC  HL
               LD   (HL),B
               JR   BCRWC

;READ BC FROM HL IN SYS PAGE

RDBC:          IN   A,(251)
               EX   AF,AF'
               XOR  A
               OUT  (251),A
               SET  7,H
               RES  6,H
               LD   C,(HL)
               INC  HL
               LD   B,(HL)
               JR   BCRWC

;REPLACE CALL CMR:DW NRREAD - FASTER

RDA:           IN   A,(251)
               EX   AF,AF'
               XOR  A
               OUT  (251),A
               SET  7,H
               RES  6,H
               LD   A,(HL)

BCRWC:         EX   AF,AF'
               OUT  (251),A
               EX   AF,AF'
               RET

GTHL:          LD   E,(HL)
               INC  HL
               LD   D,(HL)
               INC  HL
               EX   DE,HL
               RET

;GET THE HALF FLAG3

HLFG:          EX   (SP),HL
               PUSH HL       ;ORIG HL ON STACK, THEN RET ADDR
               LD   HL,PHLR
               EX   (SP),HL
               PUSH HL       ;ORIG HL, PHLR, RET ADDR
               LD   HL,FLAG3
               RET

PHLR:          POP  HL
               RET

;SET FLAG3 SUBROUTINE

SETF0:         CALL HLFG
               SET  0,(HL)
               RET

SETF1:         CALL HLFG
               SET  1,(HL)
               RET

SETF2:         CALL HLFG
               SET  2,(HL)
               RET

SETF3:         CALL HLFG
               SET  3,(HL)
               RET

SETF4:         CALL HLFG
               SET  4,(HL)
               RET

SETF5:         CALL HLFG
               SET  5,(HL)
               RET

SETF6:         CALL HLFG
               SET  6,(HL)
               RET

SETF7:         CALL HLFG
               SET  7,(HL)
               RET

;BIT TEST OF FLAG3 ROUTS.

BITF0:         CALL HLFG
               BIT  0,(HL)
               RET

BITF1:         CALL HLFG
               BIT  1,(HL)
               RET

BITF2:         CALL HLFG
               BIT  2,(HL)
               RET

BITF3:         CALL HLFG
               BIT  3,(HL)
               RET

BITF4:         CALL HLFG
               BIT  4,(HL)
               RET

BITF5:         CALL HLFG
               BIT  5,(HL)
               RET

BITF6:         CALL HLFG
               BIT  6,(HL)
               RET

BITF7:         CALL HLFG
               BIT  7,(HL)
               RET

;BORDER COLOUR CHANGE

BCC:           LD   A,(RBCC)
               AND  &0F
               RET  Z

               AND  7
               AND  E
               OUT  (ULA),A
               RET

;BORDER COLOUR RESTORE

BCR:           PUSH AF
               CALL NRRD
               DEFW BORDCR         ;BORDCOL

               OUT  (ULA),A
               POP  AF
               RET

;ERROR REPORT MESSAGES

REP4:          LD   A,85
               DEFB &21

REP5:          LD   A,86
               DEFB &21

REP6:          LD   A,87
               DEFB &21

REP12:         LD   A,93
               DEFB &21

REP13:         LD   A,94
               DEFB &21

REP18:         LD   A,99
               DEFB &21

REP19:         LD   A,100
               DEFB &21

REP20:         LD   A,101
               DEFB &21

REP22:         LD   A,103
               DEFB &21

REP23:         LD   A,104
               DEFB &21

REP24:         LD   A,105
               DEFB &21

REP25:         LD   A,106
               DEFB &21

REP27:         LD   A,22
               DEFB &21      ;EOF

REP28:         LD   A,109
               DEFB &21

REP30:         LD   A,111
               DEFB &21

REP31:         LD   A,112
               DEFB &21

REP32:         LD   A,113
               DEFB &21

REP33:         LD   A,114
               DEFB &21

REP35:         LD   A,115
               DEFB &21

REP36:         LD   A,116

               LD   (DERV),A
               CALL DERR
DERV:          DEFB 0

;CONVERT NUMBER IN A TO DIGITS IN ABC

CONR:          PUSH DE
               LD   H,0
               LD   L,A
               LD   DE,100
               LD   C," "
               CALL CONR1
               PUSH AF       ;HUNDREDS DIGIT
               LD   DE,10
               CALL CONR1
               PUSH AF       ;TENS DIGIT
               LD   A,L
               ADD  A,&30
               LD   B,A
               POP  AF
               LD   C,A
               POP  AF
               POP  DE
               RET

CONR1:         XOR  A
CONR2:         SBC  HL,DE
               JR   C,CONR3

               INC  A
               JR   CONR2

CONR3:         ADD  HL,DE
               AND  A
               JR   NZ,CONR4

               LD   A,C      ;SPACE OR ZERO
               RET

CONR4:         LD   C,"0"    ;**
               ADD  A,C
               RET

;SAMDOS ERROR PRINT
;DE HOLDS TRACK AND SECTOR

DERR:          CALL BCR

               LD   HL,(HKSP)
               LD   A,H
               OR   L
               JR   Z,DERR1

               LD   SP,HL
               RET

DERR1:         LD   A,D
               CALL CONR
               LD   (PRTRK),A
               LD   (FMTRK),A
               LD   (PRTRK+1),BC
               LD   (FMTRK+1),BC
               LD   A,E
               CALL CONR
               LD   (PRSEC),BC
               CALL NRRDD
               DEFW CHADD

               CALL NRWRD
               DEFW XPTR

               CALL NRRD
               DEFW CHADP

               CALL NRWR
               DEFW XPTRP

               XOR  A
               LD   (FLAG3),A
               LD   (NRFLG),A  ;RECURSE OK
               LD   E,A
               POP  HL
               LD   A,(HL)     ;ERR NO.
               LD   SP,&7FFA   ;OK WHETHER CMD OR HOOK
               RET

INVALID:       EQU  0
ERROR:         EQU  7
TREAM:         EQU  8
NO:            EQU  11
SNOTS:         EQU  17
SNAME:         EQU  18
TOOMANY:       EQU  20
TATEMENT:      EQU  21
FILE:          EQU  23

ERRTBL:        DEFB " "+&80            ;0

               DEFB " "+&80            ;1

               DEFB " "+&80            ;2

               DEFM "Escape requeste"  ;3
               DEFB "d"+&80

               DEFM "TRK-"             ;4
PRTRK:         DEFB &20
               DEFB &20
               DEFB &20
               DEFM ",SCT-"
PRSEC:         DEFB &20
               DEFB &20
               DEFM ",Erro"
               DEFB "r"+&80

               DEFM "Format TRK-"      ;5
FMTRK:         DEFB &20
               DEFB &20
               DEFB &20
               DEFM " los"
               DEFB "t"+&80

               DEFM "Check disk in "   ;6
               DEFM "driv"
               DEFB "e"+&80

               DEFB " "+&80            ;7

               DEFB " "+&80            ;8

               DEFB " "+&80            ;9

               DEFB INVALID            ;10
               DEFM "devic"
               DEFB "e"+&80

               DEFB " "+&80            ;11

               DEFM "Verify faile"     ;12
               DEFB "d"+&80

               DEFM "Wrong "           ;13
               DEFB FILE
               DEFM " typ"
               DEFB "e"+&80

               DEFB " "+&80            ;14

               DEFB " "+&80            ;15

               DEFB " "+&80            ;16

               DEFB " "+&80            ;17

               DEFM "Reading "         ;18
               DEFM "a write "
               DEFB FILE+&80

               DEFM "Writing "         ;19
               DEFM "a read "
               DEFB FILE+&80

               DEFB NO                 ;20
               DEFM "AUTO* "
               DEFB FILE+&80

               DEFB " "+&80            ;21

               DEFB NO                 ;22
               DEFM "such driv"
               DEFB "e"+&80

               DEFM "Disk is write "   ;23
               DEFM "protecte"
               DEFB "d"+&80

               DEFM "Disk ful"         ;24
               DEFB "l"+&80

               DEFM "Directory ful"    ;25
               DEFB "l"+&80

               DEFM "File"             ;26
               DEFB SNOTS
               DEFM "foun"
               DEFB "d"+&80

               DEFB " "+&80            ;27

               DEFM "File"             ;28
               DEFB SNAME
               DEFM " use"
               DEFB "d"+&80

               DEFB " "+&80            ;29

               DEFB "S"                ;30
               DEFB TREAM
               DEFM " use"
               DEFB "d"+&80

               DEFM "Channel use"      ;31
               DEFB "d"+&80

               DEFM "Directory"        ;32
               DEFB SNOTS
               DEFM "foun"
               DEFB "d"+&80

               DEFM "Directory"        ;33
               DEFB SNOTS
               DEFM "empt"
               DEFB "y"+&80

               DEFB NO                 ;34
               DEFM "pages fre"
               DEFB "e"+&80

               DEFM "PROTECTED "       ;35
               DEFB FILE+&80

;NON MASK INT ROUT.
NMI:           LD   (STR),SP
               LD   SP,STR
               LD   A,I
               PUSH AF
               PUSH HL
               PUSH BC
               PUSH DE
               EX   AF,AF'
               EXX
               PUSH AF
               PUSH HL
               PUSH BC
               PUSH DE
               PUSH IX
               PUSH IY

;PUSH RETURN ADDRESS ON STACK

               LD   HL,SNAP7
               PUSH HL
               LD   (HKSP),SP

;TEST FOR SNAPSHOT TYPE

SNAP29:        LD   A,4
               OUT  (251),A    ;ZX RAM AT 8000H
               IM   1

SNAP3:         LD   BC,&F7FE
               IN   E,(C)      ;BITS 0-4=DIGITS 1-5
               BIT  0,E
               JR   Z,SNAP31   ;JR IF "1"

               BIT  4,E
               JR   NZ,SNAP32  ;JR IF NOT "5"

SNAP31:        LD   HL,(NMIKA)
               LD   A,(NMIKP)  ;NORMALLY 4
               OUT  (251),A
               CALL HLJUMP     ;JUST RET USUALLY (0004H)
               JR   SNAP29

SNAP32:        BIT  1,E        ;RET IF "2"
               RET  Z

               BIT  2,E
               JR   NZ,SNAP3A  ;JR IF NOT "3"

               LD   A,&14      ;SCREEN$ TYPE
               LD   D,&1B      ;SCREEN LENGTH MSB
               JR   SNAP4

SNAP3A:        BIT  3,E
               JR   NZ,SNAP3B  ;JR IF NOT "4"
               LD   A,&05      ;48K SNAPSHOT TYPE
               LD   D,&C0      ;LEN MSB
               JR   SNAP4

SNAP3B:        INC  A
               AND  7
               OUT  (C),A      ;BORDER
               LD   B,&FE
               IN   E,(C)
               BIT  2,E
               JR   NZ,SNAP3   ;LOOP IF NOT "X"

SNAP3C:        IN   E,(C)
               BIT  2,E
               JR   Z,SNAP3C   ;LOOP TILL NOT "X"
               CALL &005F      ;DELAY BC
               LD   A,(SNPRT0)
               OUT  (251),A
               LD   A,(SNPRT2)
               OUT  (252),A
               EI
               JP   ENDS

;SAVE VARIABLES OF FILE

SNAP4:         LD   (SNME),A
               LD   HL,&8000   ;ZX RAM STARTS AT 8000H
               LD   E,L        ;ZERO E
               LD   (SNLEN),DE
               LD   (SNADD),HL

;TEST FOR DIRECTORY SPACE

               LD   IX,DCHAN
               LD   B,&FE
               IN   A,(C)
               RRA
               LD   A,1
               JR   C,SNAP4A   ;JR IF NOT "SHIFT"

               INC  A          ;DRIVE 2

SNAP4A:        CALL CKDRX
               LD   A,&40
               CALL FDHR
               JR   NZ,SNAP3           ;LOOP IF NO SPACE

;FORM SNAPSHOT FILE NAME

               LD   A,D
               AND  7
               JR   Z,SNAP5

               ADD  A,&30
               LD   (SNME+5),A

SNAP5:         LD   L,E
               SLA  L
               DEC  L
               LD   A,(IX+RPTH)
               ADD  A,L
               ADD  A,&40
               LD   (SNME+6),A

;TRANSFER NAME TO FILE AREA

               LD   HL,SNME
               LD   DE,NSTR1
               LD   BC,24
               LDIR

;OPEN A FILE

               XOR  A
               LD   (FLAG3),A
               LD   (PGES1),A
               CALL OFSM

;SAVE REGISTERS IN DIRECTORY

               LD   HL,STR-20
               LD   DE,FSA+220
               LD   BC,22
               LD   A,(NSTR1)
               CP   5
               JR   Z,SNAP6  ;JR IF 48K SNAP - NO HDR. COPY REGS

               XOR  A
               LD   (DE),A     ;FLAGS
               INC  DE
               LD   (DE),A     ;SCREEN MODE
               CALL SVHD
               LD   HL,SNPTAB  ;SCREEN$ START, LEN
               LD   DE,FSA+236
               LD   BC,7

SNAP6:         LDIR
               LD   HL,(SNADD)
               LD   DE,(SNLEN)
               CALL DSVBL
               JP   CFSM

SNPTAB:        DEFB &6E,&00,&80  ;START (IF 256K MACHINE)
               DEFB &00,&00,&1B  ;LEN
               DEFB &FF

;RETURN ADDRESS OF SNAPSHOT

SNAP7:         DI
               LD   A,3
               OUT  (251),A    ;SPECTRUM "ROM" AT 8000H
               LD   HL,0
               LD   (HKSP),HL
               LD   SP,STR-20
               POP  IY
               POP  IX
               POP  DE
               POP  BC
               POP  HL
               POP  AF
               EX   AF,AF'
               EXX
               LD   HL,NMI
               LD   (&B8F6),HL
               LD   HL,SNPRT0
               LD   DE,&B8F8
               LD   BC,3
               LDIR

  ;          * LD   A,(SNPRT0)
   ;         * LD   (&B8F8),A
   ;         * LD   A,(SNPRT1)
   ;         * LD   (&B8F9),A
   ;         * LD   A,(SNPRT2)
   ;         * LD   (&B8FA),A  ;STORE INFO IN EMULATOR ROM AREA

               POP  DE
               POP  BC
               POP  HL
               POP  AF
               LD   I,A
               AND  A
               JR   Z,SNAP8

               CP   &3F
               JR   Z,SNAP8

               IM   2
SNAP8:         LD   SP,(STR)
               JP   &B900


;--------------
Part_E1:


;DISC FORMAT ROUTINE

DFMT:          DI
               CALL CKDRV
               CALL SELD
               LD   B,10

DFMTA:         PUSH BC
               CALL INSTP
               POP  BC
               DJNZ DFMTA

               CALL REST
               DI
               LD   C,DRSEC
               CALL SADC
               CALL GTBUF
               CALL RDDATA
               AND  &0D
               JR   NZ,DFMTB   ;JR IF ERROR ON READ T0/S1

               CALL PMO6       ;"FORMAT "
  ;          * LD   HL,DRAM+&D2
  ;          * CALL PNDN2      ;PRINT DISC NAME
               CALL PM7K       ;"Y/N", KEY
               RET  NZ

DFMTB:         CALL GETSCR
               CALL PMOA       ;PRINT "FORMAT DISK AT TRACK "
               LD   DE,&0001   ;T0/S1
               CALL DMT1       ;PREPARE TRACK DATA
               LD   HL,FTADD+&0176 ;END OF T0/S1 ENTRY 1
               CALL FESET
               JR   FMT1A

FMT1:          CALL DMT1       ;PREPARE TRACK DATA

FMT1A:         CALL SCTRK      ;DISPLAY TRACK NUMBER
               CALL FMTSR
               BIT  5,A
               JP   NZ,REP23   ;JP IF WRITE-PROTECTED

               INC  D
               CALL TSTD
               CP   D
               JR   Z,FMT7     ;JR IF ALL TRACKS DONE

               AND  &7F
               CP   D
               JR   Z,FMT6     ;JR IF TIME FOR SIDE 2

               CALL INSTP
               LD   A,(SKEW)   ;&FF GIVES SKEW 1, &FE: 2, &00: 0
               DEC  E          ;E=0-9
               ADD  A,E
               JR   C,FMT5     ;JR IF E.G. FF+02=01
                               ;           FF+01=00
                               ;CONT IF    FF+00=FF
                               ;JR IF E.G. FE+02=00
                               ;CONT IF    FE+01=FF
                               ;CONT IF    FE+00=FE

               ADD  A,10       ;FF->9, FE->8
FMT5:          LD   E,A
               INC  E          ;1-10
               JR   FMT1

FMT6:          CALL REST
               LD   D,&80
               CALL SELD
               JR   FMT1

FMT7:          CALL REST       ;DE=0001
               LD   A,(DSTR2)
               CP   &FF
               JR   Z,FMT11A ;JR IF SIMPLE FORMAT, NOT DISK COPY

  ;          * LD   HL,(BUF)
   ;         * LD   (SVBUF),HL

FMT8:          CALL PMOB
               CALL SCTRK
               LD   A,(DSTR2)
               CALL CKDRX
               LD   HL,FTADD

FMT9:          LD   (BUF),HL
               CALL RSAD
               INC  H
               INC  H
               CALL ISECT
               JR   NZ,FMT9    ;READ 10 SECTORS TO BUFFER

               LD   A,(DSTR1)
               CALL CKDRX
               LD   HL,FTADD
FMT10:         LD   (BUF),HL
               CALL WSAD
               INC  H
               INC  H
               CALL ISECT
               JR   NZ,FMT10   ;WRITE 10 SECTORS FROM BUFFER

               CALL ITRCK
               JR   NZ,FMT8

   ;         * LD   HL,(SVBUF)
    ;        * LD   (BUF),HL
               CALL GTIXD
               JR   FMT12

FMT11:         CALL RSAD
               CALL ISECT
               JR   NZ,FMT11

               CALL INSTP
               JR   FMT11B

FMT11A:        CALL PMOC

FMT11B:        CALL SCTRK
               CALL ITRCK
               JR   NZ,FMT11

FMT12:         CALL CLSL
               LD   HL,FTADD
               CALL PUTSCR
               JP   REST

FMTSR:         DI
               CALL COMMP      ;GET C=STAT PORT
               LD   A,C
               LD   (WSA3+1),A ;SELF-MOD STATUS PORT
               INC  C
               INC  C
               INC  C          ;DATA PORT
               PUSH BC
               LD   C,WTRK
               CALL PRECMP
               POP  BC
               LD   HL,FTADD
               JP   WSA3       ;KEEP DELAY BETWEEN PRECMP
                               ; AND WSA3 SMALL
;PRINT TRACK ON SCREEN

SCTRK:         PUSH DE
               LD   A,21       ;COLUMN 21
               CALL NRWR
               DEFW SPOSNL

               LD   L,D
               LD   H,0
               LD   A,&20
               CALL PNUM3
               DI
               POP  DE
               RET

;INCREMENT TRACK. Z IF ALL DONE

ITRCK:         INC  D
               CALL TSTD
               CP   D
               RET  Z          ;RET IF FINISHED ALL TRKS

               LD   H,A
               CALL TIRD
               LD   A,H
               JR   C,ITRK2    ;JR IF NOT RAM DISC

               LD   A,&50

ITRK2:         AND  &7F
               CP   D
               RET  NZ         ;RET IF NOT FINISHED SIDE 1

               CALL REST
               LD   D,&7F
               INC  D          ;NZ
               RET

;DOUBLE DENSITY FORMAT

DMT1:          LD   HL,FTADD
               LD   BC,&3C4E   ;60 (3CH) TRACK HDR BYTES O
               CALL WFM
               LD   B,10       ;10 SECTORS

DMT2:          PUSH BC
               LD   BC,&0C00
               CALL WFM
               LD   BC,&03F5
               CALL WFM
               LD   C,&FE
               CALL WFMIB
               LD   C,D
               RES  7,C
               CALL WFMIB      ;TRACK
               LD   A,D
               AND  &80
               RLCA
               LD   C,A
               CALL WFMIB      ;SIDE
               LD   C,E
               CALL ISECT
               CALL WFMIB      ;SECTOR
               LD   C,2
               CALL WFMIB      ;512-BYTE SECTS
               LD   C,&F7
               CALL WFMIB
               LD   BC,&164E
               CALL WFM
               LD   BC,&0C00
               CALL WFM
               LD   BC,&03F5
               CALL WFM
               LD   C,&FB
               CALL WFMIB
               LD   BC,0
               CALL WFM
               CALL WFM        ;512 BYTES FOR DATA
               LD   C,&F7
               CALL WFMIB
               LD   BC,&1B4E   ;3 BYTES EXTRA FOR MORE TIME
               CALL WFM        ; BETWEEN SECTORS
               POP  BC
               DJNZ DMT2
               LD   C,&4E
               CALL WFM
               CALL WFM
               DEC  B          ;TRACK TAIL OF 768 BYTES
                               ; (ONLY ABOUT 200 USED)
;WRITE FORMAT IN MEMORY

WFMIB:         INC  B

WFM:           LD   (HL),C
               INC  HL
               DJNZ WFM
               RET

;PLACE DTKS-4 AT (HL), AND RND WORD AT (HL-3), AND NAME

FESET:         LD   A,(DTKS)
               SUB  4
               LD   (HL),A     ;MARK DIR ENTRY 1 (T0,S1)
                               ; WITH TKS/DIR-4
;CALLED BY RENAME

FESE2:         PUSH HL
               DEC  HL
               DEC  HL
               LD   A,R
               LD   (HL),A
               DEC  HL
               CALL NRRD
               DEFW FRAMES     ;FRAMES LOW

               LD   (HL),A     ;RND NO USES R REG AND FRAMES
               LD   BC,-42
               ADD  HL,BC      ;PT TO DISC NAME DEST (BYTE D2H)
               PUSH DE
               LD   DE,NSTR1+1
               EX   DE,HL
               LD   A,(HL)
               AND  &DF
               CP   "D"
               JR   NZ,FESE3

               CALL C11SP
               JR   NZ,FESE3

               LD   (HL),"*"   ;ALTER "D      " TO "*       " SO
                          ;FORMAT "D" DOES NOT USE DISC NAME "D"
FESE3:         LD   BC,10
               LDIR            ;COPY NAME TO TRACK BUFFER
               POP  DE
               POP  HL
               RET

;INCREMENT SECTOR ROUTINE

ISECT:         INC  E
               LD   A,E
               CP   11
               RET  NZ
               LD   E,1
               RET

;PRINT TYPE OF FILE

PNTYP:         AND  &1F
               PUSH AF
               CP   22
               JR   C,PNTY1

               LD   A,13

PNTY1:         LD   B,A
               LD   HL,DRTAB
               CALL PTV2
               POP  AF         ;TYPE
               CP   16
               JR   NZ,PNTY3

               LD   B,242      ;BASIC PROGRAM
               CALL GRPNTB
               LD   A,(HL)
               AND  &C0
               JR   NZ,PNTY5

               INC  HL
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               EX   DE,HL
               CALL PNUM5
               JR   PNTY5

PNTY3:         CP   19
               JR   NZ,PNTY4

               LD   B,236      ;CODE
               CALL GRPNTB
               CALL GTVAL
               INC  C
               EX   DE,HL
               PUSH DE
               LD   A,&20
               CALL PNUM6
               LD   A,","
               CALL PNT
               POP  HL
               CALL GTVAL
               EX   DE,HL
               XOR  A
               CALL PNUM6

PNTY4:         CP   4
               JR   NZ,PNTY5

               LD   B,215              ;ZX CODE
               CALL GRPNTB
               LD   D,(HL)
               DEC  HL
               LD   E,(HL)
               EX   DE,HL
               PUSH DE
               CALL PNUM5
               LD   A,","
               CALL PNT
               POP  HL
               DEC  HL
               LD   D,(HL)
               DEC  HL
               LD   E,(HL)
               EX   DE,HL
               XOR  A
               CALL PNUM5X

PNTY5:         LD   A,(DTFLG)
               AND  A
               JR   Z,PNTY6

               LD   A,&17
               CALL PNT        ;TAB
               LD   A,35
               CALL PNT
               CALL PNT
               CALL PNDAT      ;PRINT DATE/TIME

PNTY6:         JP   PNCR

;GET NUMBER FROM HEADER

GTVAL:         LD   A,(HL)
               AND  &1F
               LD   C,A
               INC  HL
               LD   E,(HL)
               INC  HL
               LD   A,(HL)
               AND  &7F
               LD   D,A
               INC  HL
               RET

DRTAB:         DEFB &A0
               DEFB ZXS
               DEFM "BASI"
               DEFB "C"+&80            ;1
               DEFB ZXS,"D",ARRAY+&80  ;2
               DEFB ZXS,"$",ARRAY+&80  ;3
               DEFB ZXS+&80            ;4
               DEFB ZXS
               DEFM "SNP 48"
               DEFB "K"+&80            ;5
               DEFM "MD.FIL"
               DEFB "E"+&80            ;6
               DEFB ZXS,SCREENS+&80    ;7
               DEFM "SPECIA"
               DEFB "L"+&80            ;8
               DEFB ZXS
               DEFM "SNP 128"
               DEFB "K"+&80            ;9
               DEFM "OPENTYP"
               DEFB "E"+&80            ;10
               DEFM "N/A EXECUT"
               DEFB "E"+&80            ;11
               DEFB WHAT+&80           ;12
               DEFB WHAT+&80           ;13
               DEFB WHAT+&80           ;14
               DEFB WHAT+&80           ;15
               DEFM "BASI"
               DEFB "C"+&80            ;16
               DEFB "D",ARRAY+&80      ;17
               DEFB "$",ARRAY+&80      ;18
               DEFB "C"+&80            ;19
               DEFB SCREENS+&80        ;20
               DEFM "    DI"
               DEFB "R"+&80            ;21

;PRINT NUMBER IN HL

PNUM6:         LD   (SVA),A
               XOR  A
               LD   DE,0
               RR   C
               RR   D
               RR   C
               RR   D
               LD   A,D
               ADD  A,H
               LD   H,A
               LD   A,C
               ADC  A,E
               LD   B,A
               LD   DE,34464
               LD   C,1        ;65536
               LD   A,(SVA)
               CALL PNM2
               JR   PNUM5Y

PNUM5:         LD   A,&20

PNUM5X:        LD   B,0
PNUM5Y:        LD   C,0
               LD   DE,10000
               CALL PNM2
PNUM4:         LD   DE,1000
               CALL PNM1
PNUM3:         LD   DE,100
               CALL PNM1
PNUM2:         LD   DE,10
               CALL PNM1
PNUM1:         LD   A,L

;PRINT DIGIT

PNTDI:         ADD  A,&30
               JR   PNT

PNM1:          LD   BC,0
PNM2:          PUSH AF
               LD   A,B
               LD   B,0
               AND  A
PNM3:          SBC  HL,DE
               SBC  A,C
               JR   C,PNM4

               INC  B
               JR   PNM3
PNM4:          ADD  HL,DE
               ADC  A,C
               LD   C,A
               LD   A,B
               LD   B,C
               AND  A
               JR   NZ,PNM5

               POP  DE
               ADD  A,D
               RET  Z
               JR   PNT

PNM5:          CALL PNTDI
               POP  DE
               LD   A,&30
               RET

;SEND A SPACE CHARACTER

SPC:           LD   A,&20

;OUTPUT A CHAR TO CURRENT CHAN

PNT:           PUSH AF
               CALL CMR
               DEFW &0010

               POP  AF
               RET

;PRINT TEXT MESSAGE

PTM:           POP  HL

PTM2:          LD   A,(HL)
               AND  &7F
               CP   &0D
               JR   NC,PTM4

               PUSH HL
               PUSH DE
               AND  A
               JR   NZ,PTM3    ;1-12 ARE COMPRESSION CODES

               CALL CLSL       ;0=CLSL
               CP   A          ;Z

PTM3:          LD   B,A
               LD   HL,MCPT
               CALL NZ,PTV2    ;PNT MSG B FROM LIST AT HL

               POP  DE
               POP  HL
               SCF

PTM4:          CALL NC,PNT
               BIT  7,(HL)
               RET  NZ

               INC  HL
               JR   PTM2

;0=CLEAR LOWER SCREEN

SDISKS:        EQU  1
PAK:           EQU  2
ATTK:          EQU  3
YENO:          EQU  4
ENTER:         EQU  5
SDOS:          EQU  6
SFREE:         EQU  7
WHAT:          EQU  8
ARRAY:         EQU  9
ZXS:           EQU  10
SCREENS:       EQU  11

MCPT:          DEFB &A0
               DEFM " dis"
               DEFB "k"+&80
               DEFM "press a ke"
               DEFB "y"+&80
               DEFM "at trac"
               DEFB "k"+&80
               DEFM " (y/n"
               DEFB ")"+&80
               DEFM "Inser"
               DEFB "t"+&80
               DEFM " DO"
               DEFB "S"+&80
               DEFM " Fre"
               DEFB "e"+&80
               DEFM "WHAT"
               DEFB "?"+&80
               DEFM ".ARRA"
               DEFB "Y"+&80
               DEFB "Z","X"+&80
               DEFM "SCREEN"
               DEFB "$"+&80

;SCREEN ROUTINES

PMO3:          CALL PTM
               DEFB &0D
               DEFM "Number of"
               DEFB SFREE
               DEFM "K-Bytes ="
               DEFB &A0

PMO5:          CALL PTM
               DEFB 0
               DEFM "OVERWRITE "
               DEFB 34+&80

PMO6:          CALL PTM
               DEFB 0
               DEFM "FORMAT "
               DEFB 34+&80

PMO7:          CALL PTM
               DEFB &22,YENO+&80

PMO9:          CALL PTM
               DEFB 0,ENTER
               DEFM "source"
               DEFB SDISKS
               DEFB PAK+&80

PMOA:          CALL PTM
               DEFB 0
               DEFM "Format"
               DEFB SDISKS,ATTK+&80

PMOB:          CALL PTM
               DEFB 0
               DEFM "  Copy"
               DEFB SDISKS,ATTK+&80

PMOC:          CALL PTM
               DEFB 0
               DEFM "Verify"
               DEFB SDISKS,ATTK+&80

PMOD:          CALL PTM
               DEFB 0,ENTER
               DEFM "target"
               DEFB SDISKS,PAK+&80

PMOE:          CALL PTM
               DEFM " Fil"
               DEFB "e"+&80

PMOF:          CALL PTM
               DEFB SFREE
               DEFM "Slo"
               DEFB "t"+&80

PMOG:          CALL PTM
               DEFB 0
               DEFM "LOADING"
               DEFB " "+&80

PMOH:          CALL PTM
               DEFB 0
               DEFM " SAVING"
               DEFB " "+&80

PMOSD:         CALL PTM
               DEFB "S","A","M",SDOS
               DEFM "  "
               DEFB " "+&80

PMYNAE:        CALL PTM
               DEFB 34
               DEFM " (y/n/a/e"
               DEFB ")"+&80

;PRINT DISC NAME

PNDNM:         LD   HL,DNAME

;FROM FORMAT

PNDN2:         LD   A,(HL)   ;00 IF OLD DOS, OR FF IF NEV USED
                             ;IT FOR IBU OR FIRST NAME CHAR
               RES  7,(HL)   ;SO NEV CAN USE BIT 7 AS IBU MARKER
               INC  A
               CP   2
               JR   C,PMOSD  ;PRINT "   SAM DOS " IF OLD DOS

               LD   A,(HL)
               CP   "*"
               JR   NZ,PMO8  ;PRINT DISC NAME IF THERE IS ONE,
                             ;PRINT "MASTER DOS " IF NO NAME
PMOMD:         CALL PTM
               DEFM "MASTER"
               DEFB SDOS+&80

PMO8:          CALL PTM
DNAME:       ;*DEFS 10        ;- ;(Alloc 10 bytes) OVER-WRITTEN
               DEFW 0,0,0,0,0 ;+ ;(Clear 10 bytes) OVER-WRITTEN
               DEFB " "+&80

PMOOF:         CALL PTM
               DEFM "OPEN Fil"
               DEFB "e"+&80

MSGUN:         DEFM "UN "    ;SPACE, BACKSPACE CANCELS LEADING
               DEFB 8        ; SPACE ON KWDS

OHNM:          CALL CLSL
               XOR  A
               CALL WIQF     ;TOKENS ON
               CALL NRRD
               DEFW CURCMD

               CP   241
               JR   C,OHNM2  ;JR IF ERASE OR COPY

               PUSH AF
               CALL BITF1
               LD   HL,MSGUN
               LD   B,4
               CALL NZ,PBFHL ;PRINT "UN" IF PROTECT OR HIDE OFF

               POP  AF

OHNM2:         CALL PNT      ;PRINT PROTECT, HIDE, ERASE OR COPY
               LD   A,&22
               CALL PNT      ;NOW E.G. ERASE "
               LD   A,1

WIQF:          CALL NRWR
               DEFW INQUFG   ;TOKENS OFF

               RET

FNMAE:         CALL PFNME    ;PRINT FILE NAME
               CALL PMYNAE   ;Y/N/A/E
               CALL CYES1
               RET  Z        ;RET IF "Y"

               CP   "E"
               JP   Z,ENDS   ;END

               CP   "A"
               RET  NZ

               LD   HL,FLAG3
               RES  7,(HL)   ;"?" OPTION OFF - DO ALL
               CP   A        ;Z - "Y"
               RET

FNM7K:         CALL PFNME    ;PRINT FILE NAME

PM7K:          CALL PMO7     ;"Y/N"

;COMPARE FOR Y or N

CYES:          CALL BEEP

CYES1:         CALL RDKY
               JR   NC,CYES1

               AND  &DF
               CP   "Y"
               PUSH AF

CYES2:         CALL RDKY
               JR   C,CYES2

               CALL CLSL
               POP  AF
               RET

TSPCE1:        CALL BITF5
               RET  Z

               CALL PMOD
               JR   TSPC1

TSPCE2:        CALL BITF5
               RET  Z

               CALL PMO9

TSPC1:         CALL RDKY
               JR   NC,TSPC1

TSPC2:         CALL RDKY
               JR   C,TSPC2

               CALL CLSL

RDKY:          PUSH IX
               CALL CMR
               DEFW RDKEY

               POP  IX
               RET


;--------------
Part_F11:


;ZX ENTRY INTO RAM SPACE

CALL_Label:    LD   C,&AA    ;MODE
               CALL ISEPX
               CALL EVNUM
               CALL CEOS

               LD   A,C
               CP   2
               JP   NC,IOOR

               LD   A,4
               OUT  (252),A  ;SCREEN IN PAGE 4, SPECTRUM MODE
               DEC  C
               JP   Z,SNAP7  ;JP IF "CALL MODE 1"

;CALL MODE 0
               DEC  A
               OUT  (251),A  ;PAGE 3 (SPECTRUM "ROM") AT 8000H
               DI
               JP   &B914    ;JP TO "ROM"

FLTOFL:        CALL EVNAM
               LD   C,&8E    ;TO
               CALL ISEP
               CALL EVNAM2
               CALL GCHR
               JR   NMQU1

;EVAL NAME WITH "?" OPTION

NMQU:          CALL EVNAM

NMQU1:         CP   "?"
               JR   NZ,NMQU2

               CALL GTNC     ;SKIP "?"
               CALL SETF7

NMQU2:         JP   CEOS

;COPY/BACKUP SR

COBUS:         CALL EVFINS
               CALL EXDATX   ;EXDAT AND EXX STRINGS AT 4F00/4F80
               CALL EVFINS   ; IN SYS PAGE
               CALL EXDAT
               LD   A,(DSTR1)
               LD   B,A
               LD   A,(DSTR2)
               CP   B
               RET  NZ       ;RET IF 2 DISKS IN USE

               CP   3
               JP   C,SETF5  ;SET FLAG 5 IF SINGLE-DISC COPY
                             ; WITH 1 OR 2
               RET

;DISC COPY ROUTINE

COPY:          CALL GTNC
               CP   &A6      ;OVER
               JR   NZ,COPYB

               CALL GTNC
               XOR  A        ;0="SAVE OVER"

COPYB:         CALL NRWR
               DEFW OVERF

               CALL FLTOFL
               CALL COBUS    ;COPY/BACKUP SR
               CALL N2TN3    ;COPY NSTR2 TO NSTR3 TO ACT AS
                             ;SAVED NAME TEMPLATE
               CALL BUDT
               LDIR          ;INIT BUFFER

COPY11:        CALL CKDRV
               CALL BITF1
               LD   HL,(TEMPW4)
               LD   A,(TEMPB1)
               JR   NZ,COYP3 ;JR IF STILL PRT OF LAST FILE TO DO

               CALL REFBUF   ;REFRESH BUFFER
               LD   HL,FLAG3
               RES  3,(HL)   ;"FIRST PASS - OFSM NEEDED"
               CALL SNDFX    ;FIND FILE - ABORT IF NOT FOUND
               JR   C,COPY11 ;JR IF DIR TYPE

               CALL OHASR
               JR   NZ,COPY11  ;JR IF "?" AND "N"

               LD   B,11
               CALL GRPNTB
               LD   B,(HL)
               INC  HL
               LD   C,(HL)
               PUSH BC       ;NO. OF SECTORS IN FILE
               XOR  A        ;Z SO "NOT ZX FILE" TO GTFL2
               CALL GTFL2
               CALL RSADSV
               POP  BC
               LD   HL,510
               CALL M510     ;GET AHL=HL+BC*510 (MAX FILE SIZE)
               CALL PAGEFORM
               RES  7,H      ;A=PAGES, HL=LEN MOD 16K

;ENTRY HERE IF 2ND OR LATER PASSES

COYP3:         PUSH AF
               PUSH HL
               CALL PMOG     ;"LOADING "
               LD   HL,DIFA+1
               CALL PFNM0
               CALL FFPG     ;GET IN B NO. OF PAGES IN BIGGEST
                             ;FREE BLOCK PAGE IN E
               INC  B
               DEC  B
               JP   Z,REP35  ;ERROR IF NO FREE PAGES

               LD   HL,FLAG3
               RES  1,(HL)     ;"CLOSE AFTER THIS"
               LD   L,E        ;PAGE
               POP  DE         ;LEN MOD 16K
               POP  AF         ;REQUIRED PAGES, MINUS 1
               CP   B          ;FREE PAGES
               JR   C,FCP1     ;JR IF WE CAN COMPLETE THE COPY
                               ;IN THIS PASS
               SUB  B
               LD   (TEMPB1),A  ;PAGES LEFT TO DO
               LD   (TEMPW4),DE ;LEN MOD 16K LEFT
               CALL SETF1       ;"NO CLOSE"
               LD   A,B
               LD   DE,0        ;LEN MOD 16K=0

FCP1:          LD   (TEMPW1),DE ;LEN MOD 16K TO DO (FOR SVBLK)
               LD   H,A
               LD   (TEMPW2),HL ;PGES1/PAGE OF BUFFER
               CALL GCOP
               CALL LDBLK
               CALL TSPCE1
               CALL BSWOP
               CALL TRX
               CALL CKDRV
               CALL BITF3
               JR   NZ,CYSV1   ;JR IF NOT FIRST PASS

               CALL OFSM
               JR   C,COPY3    ;JR IF "OVERWRITE?" AND N

               CALL SETF3      ;"NO OFSM NEXT TIME"

CYSV1:         CALL PMOH       ;"SAVING "
               LD   HL,NSTR1+1
               CALL PFNM0
               CALL DWAIT      ;SPIN DISK IF STOPPED
               CALL GCOP
               LD   DE,(TEMPW1) ;LEN
               CALL HSVBK
               CALL BITF1
               JR   NZ,COPY3   ;NO CLOSE IF PART STILL TO DO

               LD   HL,STR-30
               LD   DE,FSA+210
               LD   BC,42
               LDIR
               CALL SCFSM

COPY3:         CALL BSWOP
               CALL TSPCE2     ;MSG IF SINGLE DRIVE
               JP   COPY11

GCOP:          LD   HL,(TEMPW2)
               LD   A,H
               LD   (PGES1),A
               LD   A,L
               OUT  (251),A
               LD   HL,&8000
               RET

;FIND FREE PAGES
;EXIT: SYS PAGE AT 8000H, B=SIZE OF MAX FREE BLOCK IN PAGES,
;      E=PAGE NUMBER

FFPG:          XOR  A
               OUT  (251),A
               LD   HL,ALLOCT+FS+31
               LD   B,0      ;BIGGEST BLOCK SO FAR=0 PAGES

FFM:           INC  L

FFPL:          DEC  L
               RET  Z

               LD   A,(HL)
               AND  A
               JR   NZ,FFPL  ;LOOP, LOOKING FOR A FREE PAGE

               LD   C,0      ;CURRENT BLOCK=0 PAGES
CFPL:          DEC  L
               LD   A,(HL)
               INC  C
               AND  A
               JR   Z,CFPL   ;LOOP WHILE PAGES FREE, COUNTING

               LD   A,B
               CP   C
               JR   NC,FFM   ;JR IF CURRENT BLOCK NO BIGGER
                             ; THAN MAX
               LD   B,C      ;NEW MAX SIZE
               LD   E,L
               INC  E        ;POINT TO FIRST PAGE IN FREE BLOCK
                             ;=NEW START PAGE
               JR   FFM

BUDT:          CALL GETSCR
               LD   HL,RPT
               LD   DE,&B800 ;ALLOWS 1580H BYTES FOR SECTOR LIST
                             ; AT A280H
               LD   BC,&0306 ;LEN OF RPT, BUF, NSR, FSA, DRAM
               RET

CLSL:          PUSH IX       ;USED BY MODE 1/2 CLS
               CALL CMR
               DEFW CLSLOW

               POP  IX
               RET

;CALL UP DIRECTORY

DIR:           CALL FDFSR      ;CALL GTIXD, GTDEF, NULL NAME
               LD   A,2
               LD   (SSTR1),A  ;STREAM 2

               XOR  A
               LD   (DTFLG),A  ;"PRINT DATES" OFF
               CALL GTNC
               CP   "="
               JP   Z,STDIR

               CP   "#"
               JR   NZ,CAT0

               CALL EVSRM
               CALL SEPARX

CAT0:          CP   249        ;"DATE" TOKEN
               JR   NZ,CAT1

               LD   (DTFLG),A  ;"PRINT DATES" ON (NZ)
               CALL GTNC       ;SKIP "DATE"
               CALL CIEL
               JR   Z,CAT2     ;"DIR DATE" GIVES DETAILED,
                               ; DATED DIR
CAT1:          CALL ALLSR
               CALL CIEL
               JR   Z,CAT1a    ;JUST "DIR" GIVES SIMPLE DIR

               CALL EVEXP      ;EVAL DRIVE, OR FILE NAME
               JR   NC,CAT11   ;JR IF IT WAS A NAME

               CALL SEPAR      ;,/;/"
               JR   NZ,CAT12

               LD   A,(DSTR1)
               PUSH AF         ;Z
               CALL EVNAM      ;EVAL NAME - ,/;/" FOLLOWED DRIVE
               CALL NZ,EVFINS  ; NUMBER
               POP  AF         ;Z
               LD   (DSTR1),A  ;KEEP DRIVE NO SET BY NUMBER

CAT11:         CALL NZ,EVFINS  ;PROCESS NAME IF RUN TIME
               CALL GCHR

CAT12:         CP   "!"
               JR   NZ,CAT2    ;JR IF E.G. DIR 1,"NAME" OR DIR 1
                               ; - DETAILED
               CALL GTNC       ;SKIP "!" - E.G. DIR 1,"NAME"!
               CALL ALLSR      ; OR DIR 1!-SIMPLE

CAT1a:         CALL CEOS
               LD   A,2        ;SIMPLE DIR
               JR   PCAT

CAT2:          CALL ALLSR
               CALL CEOS
               LD   A,1        ;WINDOW
               CALL CMR
               DEFW JCLSBL     ;MAIN CLS

               LD   A,4        ;SINGLE COLUMN DETAILED DIR

;SETUP FOR DIRECTORY
;NSTR1+1,*
;SSTR1  ,2
;A =2 FOR SIMPLE, 4 FOR DETAILED DIR

PCAT:          PUSH AF
               CALL CKDRV
               LD   A,(SSTR1)
               CALL CMR
               DEFW STREAM

               CALL ZDVS       ;ZERO VARS
               POP  AF         ;2 IF SIMPLE
               CP   2
               JR   NZ,PCAT2

               CALL DITOB      ;DIR TO BUFFER
               CALL PDIRH      ;PRINT HEADER
               LD   HL,(PTRSCR)
               LD   DE,&A000
               AND  A
               SBC  HL,DE      ;HL=TEXT LEN (1SL
L5A79: ;*
Fix_L5A79_4x:  ;4.2 & 4.3 Add    v2.2    v2.3    Src2.3       ;*
               JR   Z,PCN3     ;L5AA9 ? L5AA9 ? L5AA9_PCN3 ?  ;*
               LD   BC,&000A                                  ;*
               LD   D,B                                       ;*
               LD   E,B                                       ;*
L5A80:         INC  DE                                        ;*
               SBC  HL,BC                                     ;*
               JR   NZ,L5A80   ;L5A80 = L5A80 = L5A80         ;*
;FixEnd                                                       ;*

               LD   HL,&A000   ;HL=START, DE=FILES
               LD   A,(SRTFG)
               AND  A
               CALL NZ,ORDER

               CALL COLB     ;NAMES/LINE IN B. COULD BE 1,2,3...
               LD   C,B

PCNML:         PUSH BC       ;B=NAMES/LINE COUNTER, C=RELOAD
               CALL PFNM0    ;PRINT NAME, ADVANCE HL
               POP  BC
               DEC  DE
               LD   A,D
               OR   E
               JR   Z,PCN3   ;EXIT IF ALL DONE

               LD   A,&20
               DJNZ PCN2     ;JR AND USE SPACE UNLESS LAST
                             ; NAME ON LINE
               LD   A,&0D
               LD   B,C      ;RELOAD COUNTER

PCN2:          CALL PNT      ;CR OR SPACE
               JR   PCNML
L5AA9: ;*
PCN3:          CALL PUTSCR
               CALL PNCR
               JR   PCN4

PCAT2:         CALL FDHR     ;DO COMPLEX DIRECTORY

PCN4:          CALL PMO3     ;"Number of Free K-Bytes = "
               CALL STATS    ;GET HL=FREE SECTS, DE=FREE SLOTS
               PUSH DE
               SRL  H
               RR   L        ;HALVE FREE SECTS
               XOR  A
               CALL PNUM4    ;FREE K
               CALL PNCR     ;CARRIAGE RETURN
               LD   HL,(FCNT)
               PUSH HL
               XOR  A
               CALL PNUM4    ;PRINT NO. OF FILES IN CURRENT DIR
               CALL PMOE     ;PRINT " File"
               POP  HL
               CALL PLUR
               LD   A,","
               CALL PNT
               CALL SPC
               POP  HL
               PUSH HL
               XOR  A
               CALL PNUM4    ;PRINT SLOTS FREE
               CALL PMOF     ;PRINT " Free Slot"
               POP  HL
               CALL PLUR

PNCR:          LD   A,&0D

PTHP:          JP   PNT

PLUR:          DEC  HL
               LD   A,H
               OR   L
               RET  Z

               LD   A,"s"
               JR   PTHP     ;MAKE IT PLURAL UNLESS ONLY 1 FILE

;PRINT DIRECTORY HEADER

PDIRH:         CALL PNCR
               CALL PNDNM    ;PRINT DISC NAME
               CALL GPATD    ;GET BC=LEN, HL=START OF PATH DATA
               LD   B,C
               CALL PBFHL    ;PRINT B FROM HL
               CALL PNCR
               JR   PNCR

ZDVS:          LD   HL,0
               LD   (CNT),HL   ;ZERO "SECTORS USED" COUNTER
               LD   (TCNT),HL  ;ZERO "FILES ON DISC" COUNTER
               LD   (FCNT),HL  ;ZERO "FILES IN DIR" COUNTER
               RET

DITOB:         CALL GETSCR
               LD   HL,&A000
               LD   (PTRSCR),HL ;INIT BUFFER PTR
               LD   A,2
               JP   FDHR        ;SIMPLE DIR TO SCREEN BUFFER

;GET HL=FREE SECTS, DE=FREE SLOTS

STATS:         CALL TSTD     ;TRACKS TO A
               LD   C,A
               LD   A,(DTKS)
               LD   B,A
               CALL TIRD     ;CY IF NOT RAM DISC
               LD   A,C
               BIT  7,A
               JR   Z,TRK0   ;JR IF SINGLE SIDED

               JR   C,TRKM1  ;JR IF NORMAL DISK

               SUB  48       ;E.G 130->82, 168->120
               DEFB &FE      ;"JR+1"

TRKM1:         ADD  A,A      ;*2 FOR DOUBLE-SIDED NORMAL.
                             ; 208->160, 168->80
TRK0:          SUB  B        ;NUMBER OF NON-DIRECTORY TRACKS
               LD   HL,0
               LD   B,10     ;10 SECTS/TRK
               LD   D,H
               LD   E,A

TRKL:          ADD  HL,DE    ;CALC SECTORS FOR DATA
               DJNZ TRKL

               LD   A,(DTKS)
               CP   5
               JR   C,TRK1

               INC  HL       ;ALLOW FOR UNUSED T4/S1 GIVING
                             ; EXTRA SPACE IF DTKS 5 OR MORE
TRK1:          LD   DE,(CNT)
               AND  A
               SBC  HL,DE
               PUSH HL       ;FREE SECTS
               LD   A,(DTKS)
               LD   C,A
               ADD  A,A      ;*2
               ADD  A,A      ;*4
               ADD  A,C      ;*5 (5-195)
               LD   L,A
               LD   H,0
               ADD  HL,HL    ;*10
               LD   A,C
               CP   5
               JR   C,PCT2   ;JR IF 4-TRACK DIR

               DEC  HL       ;ALLOW FOR 2 UNUSABLE T4/S1 ENTRIES

PCT2:          ADD  HL,HL     ;*20
               LD   DE,(TCNT) ;FILES ON DISC
               AND  A
               SBC  HL,DE     ;ALL SLOTS-FILES=FREE SLOTS
               EX   DE,HL
               POP  HL
               RET

;ENTRY: D=STRINGS TO SORT, HL=FIRST STRING

ORDER:         LD   BC,10    ;NAME LEN
               LD   A,C
               DEFB &FE      ;"JR+1"

;A=LEN TO SORT ON, BC=STRING LEN, DE=NO. OF STRINGS, HL=START

HORDER:        EXX

               DEC  A
               LD   (SLNTG+1),A ;LEN TO MATCH ON AFTER 1ST
               PUSH DE          ;FILES TO SORT
               PUSH HL          ;START

MAINSORTLP:    PUSH DE          ;FILES
               PUSH DE
               EXX
               POP  BC          ;ELEMENTS TO SORT COUNTER
               INC  C
               DEC  C
               JR   Z,NBBMP     ;EG COUNT 0100H LEFT ALONE

               INC  B      ;COUNT 0106H->0206H SO MSB/LSB CNT OK

NBBMP:         EXX
               PUSH HL          ;START OF CURRENT EL.

;ENTRY: BC"=STRINGS TO EXAMINE
;       HL=TOP OF LIST (CURRENT MAX), BC=STR. LEN

LDMAX:         LD   D,H
               LD   E,L

FETFIRST:      LD   A,(DE)

ALORDLP:       ADD  HL,BC
               EXX
               DEC  C
               JR   Z,SRTMB

SRTME:         EXX

FIRSTCP:       CP   (HL)
               JR   C,ALORDLP  ;IF (HL) LOWER THAN MAX,
                               ; NO NEW MAX - LOOP
               JR   NZ,LDMAX

               PUSH HL
               PUSH DE       ;SAVE STARTS OF BOTH STRINGS

SLNTG:         LD   B,0      ;SELF-MOD. 9 CHARS LEFT TO MATCH ON
                             ; USUALLY
STRH1:         INC  DE       ;MAX PTR - IN COMMON MEMORY
               INC  HL
               LD   A,(DE)
               CP   (HL)     ;CMP EQUIVALENT LETTERS IN EACH
               JR   NZ,STRH4 ;EXIT IF MISMATCH FOUND

STRH3:         DJNZ STRH1    ;LOOP, MATCHING UP TO B' CHARS

STRH4:         POP  DE
               POP  HL       ;GET BACK STARTS OF BOTH STRINGS
               LD   B,0
               JR   C,FETFIRST ;EITHER PICK UP OLD MAX 1. LETTER

               JR   LDMAX    ;OR RELOAD DE WITH NEW MAX PTR,
                             ; AS NEEDED
SRTMB:         DJNZ SRTME    ;DEC MSB
               EXX

;GOT MAX STRING

               POP  HL       ;HL PTS CURRENT EL., DE TO MAX
               LD   B,C

SORTMOVELP:    LD   A,(DE)  ;SWOP TOP OF LIST (HL) WITH MAX (DE)
               EX   AF,AF'
               LD   A,(HL)
               LD   (DE),A
               EX   AF,AF'
               LD   (HL),A
               INC  HL
               INC  DE
               DJNZ SORTMOVELP

               POP  DE       ;ELEMENTS TO SORT COUNTER
               DEC  DE
               LD   A,D
               OR   E
               JR   NZ,MAINSORTLP

               POP  HL
               POP  DE
               RET

;GET IN B NUMBER OF COLUMNS FOR DIRECTORY

COLB:          LD   A,(DCOLS)
               LD   B,A
               AND  A
               RET  NZ       ;RET IF FIXED NUMBER ELSE
                             ;DETERMINE ACCORDING TO WINDW WIDTH
               CALL NRRDD
               DEFW WINDRHS  ;WINDRHS/LHS IN C/B

               LD   A,C
               SUB  B
               INC  A      ;WINDOW WIDTH (NORMALLY 32, 64 OR 85)
               LD   B,1
               SUB  11
               RET  C

               INC  A
               DEC  B

COLBL:         INC  B        ;COUNT NAMES THAT WILL FIT ACROSS
                             ; CURRENT WINDOW
               SUB  11       ;NAME LEN, INC SPACE
               JR   NC,COLBL
               RET

;OVER OPTION

OVERO:         CALL GTNC
               CP   &A6        ;OVER
               RET  NZ

SF1S:          CALL SETF1
               JR   GTNCH

;RENAME/ERASE DIR OPTION

REDIR:         CALL GTNC

REDI2:         CP   144        ;"DIR"
               RET  NZ

               CALL SETF4
               JR   GTNCH      ;SKIP "DIR"

;ALLOW "?" TO MEAN "ALL FILES"

ALLSR:         SUB  "?"
               RET  NZ

               DEC  A
               LD   (CDIRT),A  ;CDIRT=FF
GTNCH:         JP   GTNC       ;SKIP "?"

;ERASE A FILE

ERAZ:          CALL OVERO
               CALL REDI2
               CALL NMQU
               CALL EVFINS

ERAZ3:         CALL SNDFX
               JR   Z,ERAZ33   ;JR IF FLAG=ERASE FILE, NOT DIR

               JR   NC,ERAZ3   ;SKIP ERASE IF NOT DIR TYPE

               AND  A          ;NC

ERAZ33:        JR   C,ERAZ3    ;SKIP RENAME IF DIR TYPE

               LD   A,(HL)
               CALL BITF1
               JR   NZ,ERAZ45  ;JR IF "ERASE OVER"

               BIT  6,A
               JR   NZ,ERASP   ;JR IF PROTECTED

ERAZ45:        CALL OHASR
               JR   NZ,ERAZ3   ;JR IF "?" AND "N"

               CALL BITF4
               JR   Z,ERAZ46   ;JR IF ERASE FILE, NOT DIR

               LD   BC,DIRT
               ADD  HL,BC
               LD   A,(HL)
               LD   (CDIRT+1),A ;SET TEMP DIR TO SUBJECT OF
                                ;ERASURE (this is EXX version,
                                ;moved to cdirt by EXDAT)
               CALL EXDAT      ;SAVE CURRENT FILE DETAILS
               LD   L,"*"
               LD   (NSTR1+1),HL ;H CANNOT BE "." (2EH) SO NAME
                                 ; IS "*"+NOT "."
               LD   A,&10    ;FIND FILE, IGNORE TYPE
               CALL FDHR     ;FIND ANY FILE IN THIS DIRECTORY
               JP   Z,REP33  ;"DIRECTORY NOT EMPTY" IF A FILE
                             ; FOUND
               CALL EXDAT    ;ORIG CDIRT AND NAME BACK
               CALL PTSVT    ;LOAD AND POINT TO DIR FILE ENTRY
                             ; AGAIN
ERAZ46:        LD   (HL),0
               CALL WSAD
               JR   ERAZ3

ERASP:         CALL BEEP
               CALL SETF3    ;"TRIED TO ERASE PROTECTED FILE"
               JR   ERAZ3

;OPTIONAL HANDLER FOR ERASE, PROTECT AND HIDE

OHAND:         CALL OHASR
               RET  NZ       ;RET IF "?" OPTION AND "N"

               LD   (HL),A   ;ERASED/PROTECTED/HIDDEN
               JR   NWSADH

OHASR:         CALL SETF0    ;"GOT ONE"
               CALL BITF7
               RET  Z        ;RET IF NOT "?" OPTION

               PUSH AF
               PUSH DE
               PUSH HL
               CALL OHNM     ;CMD NAME
               CALL FNMAE    ;FILE NAME, Y/N/A/E, KEY.
               POP  HL       ; Z IF YES OR ALL
               POP  DE
               POP  BC
               LD   A,B
               RET             ;Z IF YES

;Z IF DIRECTORY OK

CKDIR:         INC  H
               DEC  HL
               DEC  HL
               LD   A,(CDIRT)
               CP   &FF
               RET  Z        ;RET IF "MATCH ANY"

               CP   (HL)
               RET

;RENAME A FILE

RENAM:         CALL REDIR    ;SKIP ANY "DIR", SET FLAG4
               CP   142      ;"TO"
               JR   NZ,RENM1

               CALL EVNAMX
               CALL CEOS

;RENAME DISK
               CALL EVFINS
               LD   DE,&0001
               CALL RSAD
               CALL CLRRPT
               LD   B,255
               CALL GRPNTB
               CALL FESE2      ;COPY NAME, NEW RND NO.

NWSADH:        JP   WSAD

RENM1:         CALL FLTOFL
               CALL COBUS
               CALL N2TN3      ;COPY 2ND NAME TO TEMPLATE AREA
               LD   A,(DSTR1)  ;FIRST FILE DRIVE
               CALL CKDRX      ;SO E.G. RENAME "D2:ASD"
                               ; TO "FROG" OK
RENM2:         CALL REFBUF     ;REFRESH BUFFER
               CALL SNDFX      ;FIND 1ST NAME FILE
                               ; - ABORT IF NOT FOUND
               JR   Z,RENM3    ;JR IF FLAG= RENAME FILE, NOT DIR

               JR   NC,RENM2   ;SKIP RENAME IF NOT DIR TYPE

               AND  A          ;NC

RENM3:         JR   C,RENM2    ;SKIP RENAME IF DIR TYPE

               CALL OHASR
               JR   NZ,RENM2   ;JR IF "?" AND "N"

               PUSH HL
               CALL EXDAT      ;FIRST FILE NAME TO NSTR2
               POP  HL
               CALL TRX0       ;COPY LOADED NAME TO NSTR1,
               CALL FINDC      ; USE TEMPLATE ON IT
               JP   Z,REP28    ;"FILE NAME USED"

               CALL PTSVT      ;PT TO FIRST FILE ENTRY AGAIN
               PUSH DE         ;T/S
               PUSH HL
               INC  H
               DEC  HL
               DEC  HL         ;SUB DIR TAG
               LD   A,(CDIRT)  ;SUB DIR VALUE FOR 2ND NAME
               LD   (HL),A     ;SET SUB DIR TAG
               POP  DE
               INC  DE         ;NAME POSN
               LD   HL,NSTR1+1
               LD   BC,10
               LDIR
               POP  DE         ;T/S
               CALL WSAD
               CALL SETF0      ;"DONE ONE"
               CALL EXDAT
               JR   RENM2

REFBUF:        CALL BITF2
               RET  Z          ;RET IF SNDFL NOT USED YET

RFB2:          LD   DE,(SVTRS)
               JP   RSAD       ;ENSURE DRAM CONTAINS DIR ENTRIES
                               ;IN CASE 2ND ONE WANTED
;PT TO SVTRS ENTRY

PTSVT:         CALL RFB2
               LD   A,(SVDPT)
               LD   (IX+RPTH),A
               JP   POINT

;COPY NAME FROM 2 TO 3 TO ACT AS NAME TEMPLATE

N2TN3:         LD   HL,NSTR2+1
               LD   DE,NSTR3+1
               LD   BC,14
               LDIR
               RET

;PROTECT ROUTINE

PROT:          LD   A,&40
               DEFB &21        ;"JR+2"

;HIDE FILE ROUTINE

HIDE:          LD   A,&80
SFBT:          LD   (HSTR1),A
               CALL GTNC
               CP   &89        ;OFF
               CALL Z,SF1S     ;SET F1, SKIP

               CALL NMQU
               CALL EVFINS

SFB2:          CALL SNDFX
               LD   A,(HSTR1)
               LD   C,A
               CPL
               LD   B,A
               LD   A,(HL)
               AND  B
               CALL BITF1
               JR   NZ,SFB4

               OR   C
SFB4:          CALL OHAND
               JR   SFB2

;FROM ERASE, HIDE, PROTECT

SNDFX:         CALL SNDFL
               JR   C,SNDTC    ;JR IF GOT ONE

               POP  HL         ;JUNK RETURN
               CALL BITF0
               RET  NZ         ;OK IF DONE AT LEAST 1 FILE

               CALL BITF3      ;SET BY ERASE. ALWAYS 0 FROM COPY
               JP   NZ,REP36   ;"PROTECTED FILE"

REP26:         CALL DERR       ;"FILE NOT FOUND"
               DEFB 107

;GET INFO FOR ERASE/RENAME
;RET: CY IF DIR TYPE, Z IF FLAG4="DO FILES", NZ IF IT ="DO DIRS"

SNDTC:         CALL POINT
               LD   A,(HL)
               AND  &1F        ;TYPE
               SUB  DFT
               CP   1          ;CY IF DIR TYPE
               JP   BITF4


;--------------
Part_F12:


;SELECT A FILE IN THE DIRECTORY, KEEPING PLACE IN DIRECTORY
;SO LATER PARTS CAN BE LOOKED AT LATER.
;DRAM IS UNCHANGED AND DOES NOT NEED RE-READING UNLESS NEXT
;SECTOR NEEDED

SNDFL:         CALL BITF2
               JR   NZ,SNDF3

               CALL SETF2
               CALL GTIXD
               XOR  A
               LD   (IX+4),A     ;FDHR FLAGS=0 FOR CHKNM
               LD   (IX+RPTH),A  ;1ST ENTRY
               CALL REST
               CALL RSAD         ;T0/S1
               CALL SDTKS        ;SET DIR TKS, CHECK RAND NO
               JR   SNDF15

SNDF1:         CALL RSAD

SNDF15:        LD   (SVTRS),DE

SNDF2:         CALL POINT
               LD   A,(IX+RPTH)
               LD   (SVDPT),A
               LD   A,(HL)
               AND  A
               JR   Z,SNDF25   ;JR IF FREE SLOT

               CALL CKDIR
               CALL Z,CKNAM
               JR   NZ,SNDF3   ;JR IF NOT MATCHED

               SCF
               RET

SNDF25:        INC  HL
               OR   (HL)
               RET  Z          ;RET IF ALL USED DIRECTORY

SNDF3:         LD   A,(SVDPT)
               LD   (IX+RPTH),A
               DEC  A
               JR   Z,SNDF4    ;JR IF PTR WAS TO 2ND ENTRY

               CALL CLRRPT
               INC  (IX+RPTH)
               JR   SNDF2      ;DEAL WITH 2ND ENTRY

SNDF4:         LD   DE,(SVTRS)
               CALL ISECT
               JR   NZ,SNDF1

               INC  D
               LD   HL,DTKS
               LD   A,D
               CP   4
               JR   NZ,SNDF5

               INC  E          ;SKIP T4,S1

SNDF5:         CP   (HL)
               JR   C,SNDF1
               RET

;SEPARATOR REPORT ROUTINE

SEPARX:        CALL SEPAR
               RET  Z

REP0:          CALL DERR
               DEFB 29

;THE SEPARATOR SUBROUTINE

SEPAR:         CP   ","
               JR   Z,SEPA1

               CP   ";"
               JR   Z,SEPA1

               CP   &22        ;"""
               RET

SEPA1:         CALL GTNC
  ;          * LD   (SVA),A
  ;          * XOR  A
  ;          * LD   A,(SVA)
               CP   A
               RET

;EVALUATE PARAMETERS IN SYNTAX FOR READ AT/WRITE AT D,T,S,ADDR

EVPRM:         LD   C,&87      ;AT
               CALL ISEPX

;GET DRIVE NUMBER
               CALL EVDNM

;GET TRACK NUMBER
               CALL SEPARX
               CALL EVNUM
               PUSH BC         ;C=TRK IF RUN TIME

;GET SECTOR NUMBER
               CALL SEPARX
               CALL EVNUM
               POP  DE         ;E=TRK
               LD   B,E        ;BC=T/S
               PUSH BC

;GET ADDRESS
               CALL SEPARX
               CALL EVADDR     ;Z IF SYNTAX
               PUSH AF
               PUSH HL         ;AHL=ADDR
               CALL CIEL
               LD   BC,1       ;DEFAULT=1 SECT
               JR   Z,EVPR5    ;JR IF NO "SECTORS" PARAM

               CALL EVNUMX     ;NUMBER OF SECTORS
               JR   Z,EVPR5
               DEC  HL
               LD   A,H
               CP   4
               JP   NC,IOOR    ;ALLOW 0000-03FFH IN DECED,
                               ;      0001-0400H IN ORIG,
                               ;BC= 1-1024 SECTORS
EVPR5:         POP  HL
               POP  AF         ;AHL=ADDR, NZ IF RUN
               CALL NZ,SCASD   ;SET/CHECK ADDR IN AHL,BC=SECTORS

               POP  HL
               LD   (HKDE),HL  ;T/S
               CALL CEOS

               LD   A,(DSTR1)
               LD   (HKA),A
               RET

;SAVE HEADER INFORMATION

SVHD:          LD   HL,HD001
               LD   DE,FSA+211
               LD   B,9

SVHD1:         LD   A,(HL)
               LD   (DE),A
               CALL SBYT
               INC  HL
               INC  DE
               DJNZ SVHD1
               RET

;REMOVE FLOATING POINT NUMBERS

REMFP:         LD   A,(HL)
               CP   &0E
               JR   NZ,REMP1

               LD   BC,6
               CALL CMR
               DEFW JRECLAIM

REMP1:         LD   A,(HL)
               INC  HL
               CP   &0D
               JR   NZ,REMFP

               RET

;LOAD COMMAND (FAILED NORMAL SYNTAX)

LOAD:          CALL GTNC       ;CHAR AFTER LOAD
               CALL CFSO
               CALL Z,REMFP    ;REMOVE FP FORMS FROM HL ON
                               ; IF SYNTAX TIME
               CALL EVNUM
               JP   Z,CEOS

               LD   (TEMPW4),HL ;FILE NUM CAN BE 2 BYTES
               LD   A,L
               LD   (FSTR1),A
               CALL GTFLE

AUTOX:         LD   A,(DIFA)
               CP   &14
               JR   NZ,DLVM1 ;JR IF NOT SCREEN$ OR FORMER TYPE 5

               CALL BITF7
               JR   Z,DLVM1    ;JR IF SCREEN$

;48K SNAPSHOT IS FOUND

               CALL RSADSV
               LD   BC,&0300+250
               LD   HL,SNPRT0

RDPRTL:        IN   A,(C)
               LD   (HL),A
               INC  HL
               INC  C
               DJNZ RDPRTL

     ;       * IN   A,(250)
     ;       * LD   (SNPRT0),A
     ;       * IN   A,(251)
     ;       * LD   (SNPRT1),A
     ;       * IN   A,(252)
     ;       * LD   (SNPRT2),A

               LD   A,4
               OUT  (251),A    ;ZX "SCREEN" AT 8000H
               OUT  (252),A    ;DISPLAY PAGE 4, SCREEN MODE 1
               LD   HL,&8000
               LD   DE,&4000
               LD   A,2
               LD   (PGES1),A
               CALL LDBLK      ;LOAD 48K TO ZX IMAGE
               CALL SKSAFE
               JP   SNAP7

DLVM1:         CP   &10
               JR   NZ,DLVM2   ;JR IF NOT BASIC

               CALL BITF7
               JP   NZ,REP13   ;"WRONG FILE TYPE" IF ZX BA

               CALL NRRDD
               DEFW PROG

               PUSH BC
               POP  HL
               CALL NRRD
               DEFW PROGP

               LD   (UIFA+31),A
               LD   (UIFA+32),HL
               EX   DE,HL
               LD   C,A
               PUSH BC
               CALL NRRDD
               DEFW ELINE

               PUSH BC
               POP  HL
               CALL NRRD
               DEFW ELINP

               POP  BC
               DEC  HL
               BIT  7,H
               JR   NZ,LAB2

               DEC  A

LAB2:          PUSH BC
               PUSH DE
               CALL AHLN
               PUSH AF
               EX   DE,HL
               LD   A,C
               CALL AHLN
               EX   DE,HL
               LD   C,A
               POP  AF
               AND  A
               SBC  HL,DE
               SBC  A,C
               POP  DE
               POP  BC
               CALL PAGEFORM
               LD   (UIFA+34),A
               LD   (UIFA+35),HL
               XOR  A
               LD   (UIFA+15),A

DLVM2:         LD   IX,&4B00
               CALL TXINF
               CALL TXHED
               LD   HL,SYSP-1
               SET  6,(HL)   ;ENSURE ROM1 ON ON EXIT FROM DOS BY
                             ;ALTERING PORT 251 STATUS ON STACK
               JP   ENDSX    ;NOW JUMP TO ROM1 TO LOAD IS OK

;GET 19-BIT NUMBER

AHLN:          CALL AHLNX
               AND  &07
               RET

;GET 20-BIT NUMBER

AHLNX:         RLC  H
               RLC  H
               RRA
               RR   H
               RRA
               RR   H
               AND  &0F
               RET

IOOR:          CALL DERR
               DEFB 30

;WRITE FORMAT ON DISC

WFOD:          CALL FDFSR
               LD   A,4
               LD   (DTKS),A   ;DEFAULT TRACKS/DIR
               XOR  A
               LD   (TEMPW1),A ;DEFAULT TKS/DISC
               CALL GTNC
               CP   &8E        ;TO
               JR   Z,WFOD1    ;JR IF 'FORMAT TO "NAME"'

               CALL CIEL
               JR   Z,WFOD2X

               CALL EVNAM
               CP   ","
               JR   NZ,WFOD0   ;JR IF 'FORMAT "NAME" TO "NAME"

               CALL EVNUMX
               JR   Z,WFOD01

               LD   A,B
               AND  A
               JR   NZ,IOOR

               LD   A,C
               LD   (DTKS),A

WFOD01:        CALL GCHR
               CP   ","
               JR   NZ,WFOD2

               CALL EVNUMX     ;TOTAL TKS OF RAM DISC
               JR   Z,WFOD2

               LD   A,B
               AND  A
               JR   NZ,IOOR

               LD   A,C
               LD   (TEMPW1),A ;TKS/DISC
               AND  A
               JR   NZ,WFOD2

WIOOR:         JP   IOOR

WFOD0:         CP   &8E        ;TO
               JR   NZ,WFOD2

               CALL EVNAM2X    ;SECOND NAME OF TWO
               CALL CEOS
               CALL EVFINS     ;FIRST
               JR   WFODTC

WFOD1:         CALL EVNAM2X    ;SECOND (ONLY) NAME
               CALL CEOS

WFODTC:        CALL FTOSR      ;EXDAT AND CHECK DRIVE
               CALL EVFINS     ;SECOND NAME
               CALL FTOSR

WFOD3:         JP   DFMT

WFOD2X:        CALL CEOS
               JR   WFOD3X

WFOD2:         CALL CEOS
               CALL EVFINS     ;FIRST NAME

WFOD3X:        LD   A,(DTKS)
               LD   C,A
               CP   40
               JR   NC,WIOOR

               LD   B,4        ;LOW DTK LIMIT
               LD   A,(DSTR1)
               CP   3
               JR   C,WFOD00

               LD   B,1        ;ALLOW 1 OR MORE DIR TKS
               LD   A,C
               AND  A
               JR   Z,WFODB    ;JR IF E.G. FORMAT"D3:",0
                               ; - ERASE RAM DISK
               LD   A,(TEMPW1)
               AND  A
               JR   Z,WIOOR    ;E.G. FORMAT "D3:",5,0 ILLEGAL
                              ;AS IS FORMAT "D3:",5
               JR   WFOD07

WFOD00:        LD   A,(TEMPW1)
               AND  A
               JR   NZ,WIOOR   ;E.G. FORMAT "D1:",5,10 ILLEGAL.
                               ;ONLY FORMAT "D1",4 ETC O.K.
WFOD07:        LD   A,C
               CP   B
               JR   C,WIOOR ;ALLOW 4-39 DIR TRACKS(80-780 FILES)
                            ;IF REAL DISC OR 1-39 IF RAMDISC
               LD   A,(TEMPW1)
               AND  A
               JR   Z,WFODB    ;JR IF REAL DISC FORMAT

               LD   C,A
               LD   B,159      ;LIMIT
               LD   A,(DTKS)
               SUB  C
               JR   NC,WIOOR   ;TOTAL TKS MUST BE >DTKS

               ADD  A,C
               SUB  4
               JR   NC,WFOD02  ;JR IF DIR IS NORMAL OR LARGE

               ADD  A,B
               LD   B,A    ;LOWER LIMIT - 3 DTK=158,2=157,1=156

WFOD02:        LD   A,C
               SUB  2
               CP   B

WIORH:         JR   NC,WIOOR   ;DTKS=1, TOT TKS 2-157 OK
                               ;DTKS=2, TOT TKS 3-158
                               ;DTKS=3, TOT TKS 4-159
                               ;DTKS=4, TOT TKS 5-160
                               ;DTKS=5, TOT TKS 6-160

WFODB:         LD   A,(DSTR1)
               CP   3
               JR   C,WFOD3

               JP   FORMRD     ;FORMAT RAMDISC

TIRDXDCT:      XOR  A
               LD   (DCT),A

;TEST IF RAM DISC - NC IF SO

TIRD:          LD   A,(DRIVE)
               CP   3
               RET

;EVALUATE DRIVE NO. OR FILE NAME (USED BY DIR)
;EXIT: CY IF DRIVE

EVEXP:         CP   "*"
               JR   Z,EVDN1

               CALL CMR
               DEFW EXPEXP

               JR   Z,EVEXP2   ;JR IF STRING

               CALL EVNU2      ;UNSTACK NUMBER IF RUNNING
               SCF               ;"DRIVE"
               JR   EVDN3      ;RET, OR SET DRIVE

EVEXP2:        CALL EVST2      ;UNSTACK STRING IF RUNNING
               SCF
               CCF               ;NC="NAME"
               JR   EVNMR      ;RET, OR EVNAME

;EVALUATE DRIVE NUMBER - WRITE AT *,T,S,A ALLOWED

EVDNM:         CALL GCHR
               CP   "*"
               JR   NZ,EVDN2

EVDN1:         CALL GTNC
               PUSH AF
               CALL GTDD       ;DEFAULT NUMBER
               LD   A,(DRIVE)
               LD   C,A
               POP  AF         ;CHAR AFTER "*"
               CALL CFSO
               SCF               ;"NUMBER"
               JR   EVDN3

EVDN2:         CALL EVNUM

EVDN3:         RET  Z
               PUSH AF
               LD   A,C
               CALL CODN       ;CONVERT DRIVE NUMBER IN A TO A
               CALL DRSET      ;SET DRIVE, CDIRT
               POP  AF
               RET

CODN:          DEC  A
               CP   7
               INC  A
               JR   NC,CODN2   ;JR IF NOT 1-7 - LEAVE ALONE

               PUSH BC
               PUSH HL
               LD   C,A
               LD   B,0
               LD   HL,DRPT-1  ;DRIVE NUMBER PRETEND TABLE
               ADD  HL,BC
               LD   A,(HL)
               POP  HL
               POP  BC

CODN2:         LD   (DSTR1),A
               RET

HEVNAM:        LD   BC,(HKBC)
               LD   DE,(HKDE)
               SET  7,D
               RES  6,D        ;BUFFER ADDR ADJUST 4F10->8F10H
               XOR  A
               LD   B,A
               LD   (SVC),A    ;PAGE 0
               JR   EVNM2

REP8:          CALL DERR
               DEFB 18         ;INVALID FILE NAME

EVNAMX:        CALL GTNC

;EVALUATE FILE NAME
;EXIT: Z/NZ FOR SYN/RUN, A=CHAR AFTER

EVNAM:         CALL EVSTR

EVNMR:         RET  Z          ;RET IF SYNTAX TIME

EVNM2:         PUSH AF
               LD   A,C
               OR   B
               JR   Z,REP8     ;"INVALID FILE NAME" IF LEN 0

               LD   HL,79
               SBC  HL,BC
               JR   C,REP8     ;ERROR IF TOO LONG

               IN   A,(251)
               PUSH AF
               PUSH BC         ;REAL LEN
               LD   HL,14
               SBC  HL,BC
               JR   NC,EVNM0

               LD   BC,14

EVNM0:         LD   HL,NSTR1
               LD   A,15

EVNM1:         LD   (HL),&20
               INC  HL
               DEC  A
               JR   NZ,EVNM1

               LD   HL,NSTR1+1
               EX   DE,HL
               LD   A,(SVC)
               AND  &1F
               OUT  (251),A
               PUSH HL
               LDIR
               POP  HL
               POP  BC         ;REAL LEN
               LD   A,C
               CALL NRWR
               DEFW &4F60      ;STORE REAL LEN

               LD   DE,&4F10
               CALL CMR
               DEFW &008F      ;LDIR TO BUFFER IN SYS PAGE

               POP  AF
               OUT  (251),A
               POP  AF
               RET

EVNAM2X:       CALL GTNC

;EVALUATE SECOND FILE NAME
;EXIT: Z/NZ FOR SYNTAX/RUN. A CORUPT

EVNAM2:        CALL EXDATX
               CALL EVNAM

;SWOP LONG NAMES IN SYS PAGE BUFFER AS WELL AS OTHER DATA

EXDATX:        CALL EXDAT
               IN   A,(251)
               PUSH AF
               PUSH BC
               PUSH DE
               PUSH HL
               XOR  A
               OUT  (251),A
               LD   BC,&0058
               LD   HL,&8F10
               LD   DE,&8F68
               CALL EXDT1
               POP  HL
               POP  DE
               POP  BC
               POP  AF
               OUT  (251),A
               RET

EXDAT:         PUSH AF
               PUSH BC
               PUSH DE
               PUSH HL
               LD   HL,(CDIRT)
               LD   A,H
               LD   H,L
               LD   L,A
               LD   (CDIRT),HL ;SWOP MAIN AND EXX DIRECTORIES
               LD   BC,28
               LD   DE,DSTR1
               LD   HL,DSTR2
               CALL EXDT1
               POP  HL
               POP  DE
               POP  BC
               POP  AF
               RET

BSWOP:         CALL EXDAT
               LD   A,(DTKSX)
               EX   AF,AF'
               LD   A,(DTKS)
               LD   (DTKSX),A
               EX   AF,AF'
               LD   (DTKS),A   ;SWOP DTKS/DTKSX
               CALL BUDT       ;GET HL=RPT, DE=PAGED IN BUFFER,
                               ; BC=0306H
;SWOP BC AT HL/DE

EXDT1:         LD   A,(DE)
               EX   AF,AF'
               LD   A,(HL)
               EX   AF,AF'     ;* -
               LD   (HL),A     ;* LD (DE),A
               EX   AF,AF'     ;* EX AF,AF'
               LD   (DE),A     ;* LD (HL),A
               INC  DE
               INC  HL
               DEC  BC
               LD   A,B
               OR   C
               JR   NZ,EXDT1
               RET

;EVALUATE STRING EXPRESSION

EVSTR:         CALL CMR
               DEFW EXPSTR

EVST2:         CALL CFSO
               RET  Z

               PUSH AF
               CALL CMR
               DEFW GETSTR

               LD   (SVC),A
               POP  AF
               RET

;EVALUATE STREAM INFORMATION

EVSRM:         CALL GTNC
EVSRMX:        CALL EVNUM
               RET  Z

               PUSH AF
               LD   A,C
               CP   17
               JR   NC,EVSRE

               LD   (SSTR1),A
               POP  AF
               RET

EVSRE:         CALL DERR
               DEFB 21         ;"INVALID STREAM NUMBER"

EVNUMX:        CALL GTNC

;EVALUATE NUMBER ROUTINE

EVNUM:         CALL CMR
               DEFW EXPNUM

EVNU2:         CALL CFSO
               RET  Z

CGTINT:        PUSH AF
               CALL CMR
               DEFW GETINT

               POP  AF
               RET

;EVALUATE BIG NUMBER ROUTINE

EVBNUM:        CALL CMR
               DEFW EXPNUM

               CALL CFSO
               RET  Z

               PUSH AF
               IN   A,(251)
               PUSH AF
               IN   A,(250)
               INC  A
               CALL SELURPG    ;DOS AT 8000H TOO
               CALL CMR
               DEFW FPDT+FS

               POP  AF
               OUT  (251),A
               CALL CGTINT
               PUSH HL         ;X DIV 64K
               CALL CGTINT     ;BC=X MOD 64K
               POP  DE         ;DE=X DIV 64K
               POP  AF
               RET

FPDT:          DEFB CALC           ;X
               DEFB DUP            ;X,X
               DEFB FIVELIT
               DEFB &91,0,0,0,0    ;X,X,64K
               DEFB STO4           ;X,X,64K
               DEFB MOD            ;X,X MOD 64K
               DEFB SWOP           ;X MOD 64K,X
               DEFB RCL4           ;X MOD 64K,X,64K
               DEFB IDIV           ;X MOD 64K, X DIV 64K
               DEFB EXIT2

;EVALUATE ADDRESS ROUTINE

EVADDR:        CALL CMR
               DEFW EXPNUM

               CALL CFSO
               RET  Z

               CALL CMR
               DEFW UNSTLEN    ;AHL=PAGES/ ADDR MOD 16K

               SET  7,H        ;HL=OFFSET
               DEC  A
               LD   C,0
               INC  C          ;NZ
               RET

;CHECK FOR ALPHA CHAR

ALPHA:         CP   &41
               CCF
               RET  NC

               CP   &5B
               RET  C

               CP   &61
               CCF
               RET  NC

               CP   &7B
               RET

;TRANSFER FILE NAMES IN COPY
;E.G. COPY "*" TO "*" (NSTR3) WITH "FROG" (NSTR1) LOADED
;      GIVES "FROG" IN NSTR1
;     COPY "*" TO "?X.BIN" WITH "AB.BIN" LOADED GIVES "AX.BIN"

TRX:           LD   HL,DIFA

;CALLED BY RENAME

TRX0:          LD   DE,NSTR1
               LD   BC,15
               LDIR            ;NSTR1=DATA FROM DISC
               LD   DE,NSTR1+1
               LD   HL,NSTR3+1
               LD   B,10

TRX1:          LD   A,(HL)
               CP   "*"
               JR   Z,TRX3

               CP   "?"
               JR   Z,TRX2     ;LEAVE CHARS IN TGT OPPOSITE "?"
                               ; ALONE
               LD   (DE),A     ;TGT=SRC UNLESS SRC="?"

TRX2:          INC  HL
               INC  DE
               DJNZ TRX1
               RET

TRX3:          INC  HL
               LD   A,(HL)
               CP   "."
               RET  NZ         ;RET IF SRC="*???" UNLESS "*."

TRX4:          LD   A,(DE)
               CP   "."
               JR   Z,TRX2

               INC  DE
               DJNZ TRX4       ;LOOK FOR "." IN TGT, THEN MATCH
               RET               ; EXTENSIONS

GDIFA:         CALL NMMOV
               LD   B,220
               CALL GRPNTB
               LD   BC,33
               LDIR

               LD   HL,DIFA
               LD   DE,UIFA
               LD   C,48
               LDIR

               LD   B,13
               CALL GRPNTB
               LD   D,(HL)
               INC  HL
               LD   E,(HL)
               RET


;--------------
Part_G1:


;HOOK CODE ROUTINES

;INPUT A HEADER FROM IX

RXHED:         CALL RXSS
               JR   NZ,REP10H
               RET

RXSS:          CALL RXHSR
               CALL HCONR
               LD   A,(LSTR1)
               CP   "D"
               RET

;INPUT A HEADER FROM IX, ALLOW DEVICE D/T/N.
;RETURNS NC IF D, CY IF T/N

RXHED2:        CALL RXSS
               RET  Z

               CP   "T"
               JR   Z,EVFL75

               CP   "N"

REP10H:        JP   NZ,REP10   ;INVALID DEVICE

EVFL75:        CALL NRWR
               DEFW &5BB7      ;TEMP DEVICE

               LD   A,(DSTR1)
               CP   8
               JR   NC,EVFL8   ;JR IF NO SENSIBLE TAPE SPEED
                               ; SPECIFIED
               LD   A,112    ;DEFAULT TAPE SPEED (IRREL FOR NET)

EVFL8:         CALL NRWR
               DEFW &5BB8      ;TEMP SPEED

               LD   HL,NSTR1+1
               LD   A,(HL)
               SUB  &20
               JR   NZ,EVFL8A

               DEC  A          ;A=FF
               LD   (UIFA+1),A ;NULL NAME
               JR   EVFL8B

EVFL8A:        LD   DE,UIFA+1
               LD   BC,14
               LDIR

;OUTPUT UIFA TO ROMUIFA

TXINF:
EVFL8B:        LD   BC,0       ;COPY TO ROM HDR
               LD   DE,UIFA
               PUSH IX
               JR   TXRM       ;ENDS WITH SCF

;CONVERT NEW HDR TO OLD

HCONR:         CALL RESREG
               LD   HL,UIFA+1
               CALL EVFILE
               LD   A,(&7FFF)  ;ENTRY LRPORT VALUE ON STACK
               BIT  6,A
               CALL NZ,GTLNM   ;IF ROM 1 USED HGTHD THEN USE
                               ; LONG NAME
               LD   A,(UIFA)
               LD   (NSTR1),A
               LD   (HD001),A

               LD   A,(UIFA+31)
               LD   (PAGE1),A

               LD   HL,(UIFA+32)
               LD   (HD0D1),HL

               LD   A,(UIFA+34)
               AND  &1F
               LD   (PGES1),A

               LD   HL,(UIFA+35)
               RES  7,H
               LD   (UIFA+35),HL
               LD   (HD0B1),HL
               RET

;HOOK TO OPEN A FILE FOR LOAD/VERIFY - IX=HDR

HGTHD:         PUSH IX
               CALL RXHED2
               JR   C,HGTH2    ;JR IF TAPE OR NET

               CALL CKDRV
               CALL GTFLE
               POP  IX

;OUTPUT DIFA TO ROMDIFA (DIFA TO IX+80)

TXHED:         PUSH IX
               POP  HL
               PUSH HL
               LD   BC,&0024
               ADD  HL,BC      ;IX+24H
               CALL RDA        ;MSB OF LEN MOD 16K (HDR) TO A
               LD   HL,DIFA+&24
               XOR  (HL)
               AND  &80        ;TAKE BIT 7 FROM HDR
                               ; (SOMETIMES HI, SOMETIMES LO)
               XOR  (HL)       ;VERIFY DEPENDS ON EQUALIT
               LD   (HL),A
               LD   DE,DIFA
               LD   BC,80

;OUTPUT A HEADER

TXRM:          POP  HL         ;WAS IX
               ADD  HL,BC      ;NC. DEST IN HL=IX OR IX+80
                               ;SRC IN DE=DIFA OR UIFA
               JR   RXTXC

RXHSR:         LD   DE,UIFA
               PUSH IX
               POP  HL         ;SRC=HL, DEST=UIFA
               SCF

RXTXC:         BIT  7,H
               IN   A,(251)
               LD   BC,251
               JR   NZ,RXH1    ;JR IF NO PAGING REQUIRED

               SET  7,H
               RES  6,H
               OUT  (C),B      ;SYSTEM PAGE

RXH1:          JR   C,RXH2

               EX   DE,HL

RXH2:          LD   C,48
               LDIR
               OUT  (251),A
               SCF               ;FOR TXINF
               RET

HGTH2:         POP  HL         ;JUNK
               LD   HL,(ENTSP)
               INC  HL
               INC  HL
               INC  HL
               INC  HL
               INC  (HL)       ;STORED SP LSB (EVEN)
               INC  (HL)   ;DISCARD ONE ADDR - RET TO LOAD FILE
               LD   E,3        ;LOAD/MERGE VERIFY ENTIRE FILE
               JP   END1       ; FROM T/N

DSCHD:         CALL GTIXD
               CALL LDHDX

               LD   HL,(HKHL)
               LD   (HD0D1),HL ;START

               LD   A,(HKBC)   ;CDE WAS LEN
               LD   C,A
               LD   (PGES1),A

               LD   DE,(HKDE)
               RES  7,D
               LD   (HD0B1),DE
               RET

NETPA:         LD   A,(LSTR1)
               CP   "N"
               RET  NZ

               POP  HL         ;JUNK RET ADDR
               CALL OWSTK
               EXX               ;START, LEN TO HL, CDE FOR
               RET

HVEPG:         OUT  (251),A

;VERIFY FILE ALREADY OPENED BY HGTHD. HL=START,CDE=LEN,PAGED IN.

HVERY:         LD   BC,&4BB0+HVEP-PVECT
               CALL NETPA
               CALL DSCHD
               LD   (IX+RPTL),9

HVER1:         LD   A,D
               OR   E
               JR   NZ,HVER2

               LD   A,C
               AND  A
               JP   Z,SKSAFE

               DEC  C
               LD   DE,16384

HVER2:         CALL LBYT
               CP   (HL)
               JP   NZ,REP12   ;VERIFY FAILED

               DEC  DE
               INC  HL
               LD   A,H
               CP   &C0
               JR   C,HVER1

               CALL INCURPAGE
               JR   HVER1

HSAVE:         CALL RXHED2     ;GET HEADER, ALLOW DEVICES D/T/N
               JR   C,HSAVE2   ;JR IF T/N

               CALL CKDRV
               IN   A,(251)
               LD   (PORT3),A
               LD   A,(UIFA+31)
               CALL SELURPG
               CALL GOFSM
               JR   C,HSAVE1   ;JR IF "OVERWRITE?" AND N

               CALL SVHD
               LD   HL,(HD0D1)
               LD   DE,(HD0B1)
               CALL DSVBL
               CALL CFSM
HSAVE1:        LD   A,(PORT3)
               OUT  (251),A
               RET

HSAVE2:        LD   E,2        ;SAVE ENTIRE FILE TO TAPE OR NET
               JP   END1

DDEL:          CALL TIRD
               RET  NC         ;RET IF RAM DISC

               LD   A,(DWAI)
               INC  A

DDLP:          PUSH AF
               XOR  A
               CALL STPDX      ;DELAY ABOUT 0.25 SEC
               POP  AF
               DEC  A
               JR   NZ,DDLP
               RET

HVAR:          CALL CGTINT     ;GET DVAR PARAM
               LD   HL,DVAR
               ADD  HL,BC      ;ADD DVAR BASE
               IN   A,(250)
               AND  &1F
               ADD  A,2        ;A=DOS PAGE+1
               ADD  HL,HL
               ADD  HL,HL      ;LEFT JUSTIFY 16K OFFSET -
                               ; JUNK UPPER BITS
               LD   B,&96      ;INIT EXPONENT

HVAR1:         DEC  B
               ADD  HL,HL
               RLA
               BIT  7,A
               JR   Z,HVAR1

               RES  7,A        ;RES SGN BIT
               LD   E,A
               LD   A,B        ;EXPONENT

HVAR2:         LD   D,H
               LD   C,L
               LD   B,0
               JP   HSTKS

FNLN2:         LD   A,2
               DEFB &21        ;"JR+2"

HPTR:          LD   A,1
               DEFB &FE        ;"JR+1"

HEOF:          XOR  A
               PUSH IX
               LD   C,251
               IN   B,(C)
               PUSH BC
               PUSH AF         ;0 IF EOF, 1 IF PTR
               CALL PESR       ;AHL=PTR
               POP  DE
               DEC  D
               JR   Z,EPCOM    ;JR IF PTR

               DEC  D
               JR   NZ,HEOF2

               CALL GLEN       ;GET FILE LEN IN AHL
               JR   EPCOM

HEOF2:         CALL CPPTR
               LD   HL,0
               LD   A,H
               JR   NZ,EPCOM

               INC  HL         ;HL=1, EOF=TRUE

EPCOM:         POP  BC
               OUT  (C),B
               POP  IX

;STACK AHL AS A NUMBER

STKFP:         LD   B,A
               OR   H
               OR   L
               JR   Z,STKHL    ;ZERO IS A SPECIAL CASE

               LD   A,B
               LD   B,&98      ;INIT EXPONENT
               JR   HVAR1

;STKHL ON FPCS. USED BY EOF

STKHL:         LD   A,H
               LD   H,L
               LD   L,A
               XOR  A
               LD   E,A
               JR   HVAR2

AUTNAM:        DEFB 1
               DEFB &FF
               DEFB &FF
               DEFB "D"
               DEFB &10
               DEFM "AUTO*     "
               DEFM "    "
               DEFB 0
               DEFW &FFFF
               DEFW &FFFF
               DEFW &FFFF
               DEFW &FFFF

;LOOK FOR AN AUTO FILE

HAUTO:         CALL AUINSR
               JP   NZ,REP20   ;ERROR IF NOT FOUND

AUINC:         CALL GTFLX
               JP   AUTOX

AUINSR:        LD   A,&95      ;LOAD TOK
               CALL NRWR
               DEFW CURCMD

               LD   HL,AUTNAM
               LD   DE,DSTR1
               LD   BC,28
               LDIR
               CALL GTDEF
               CALL CKDRV
               LD   A,&10      ;"LOOK FOR NAME"
               JP   FDHR

INIT:          CALL AUINSR
               RET  NZ         ;RET IF NOT FOUND

               JR   AUINC

;HOOK OPEN FILE

HOFLE:         CALL RXHED
               CALL CKDRV
               CALL GOFSM
               JP   NC,SVHD    ;JP IF NOT "OVERWRITE?"+N

               RET

HGFLE:         CALL RXHED
               CALL GTFLE

LDHDX:         CALL RSADSV

;LOAD HEADER INFORMATION

LDHD:          LD   B,9
LDHD1:         CALL LBYT
               DJNZ LDHD1
               RET

HERAZ:         CALL RXHED
               CALL CKDRV
               CALL FINDC
               JP   NZ,REP26

               LD   (HL),0
               JP   WSAD

;WRITE AT A TRACK AND SECTOR

WRITE:         CALL EVPRM
               JR   HFWSAD

;READ AT A TRACK AND SECTOR

READ:          CALL EVPRM
               JR   HFRSAD

;HOOK READ 1 SECTOR AT DE

HRSAD:         CALL CALS

;HOOK FAR READ MANY SECTORS AT DE

HFRSAD:        LD   IY,RSAD
               JR   FRWSR

;HOOK WRITE 1 SECTOR AT DE

HWSAD:         CALL CALS

;HOOK FAR WRITE MANY SECTORS AT DE

HFWSAD:        CALL HFWCD      ;CONVERT HKA AND CHECK DRIVE
               CALL SELD
               CALL DWAIT
               LD   IY,WSAD

FRWSR:         IN   A,(251)
               PUSH AF
               XOR  A
               LD   (RDAT),A  ;SO RAM DISC KEEPS TRACK UNFIDDLED
               CALL HFWCD
               CALL GTIXD
               LD   HL,(HKHL)  ;ADDR
               LD   BC,(SVHDR) ;SECTORS
               LD   A,(HKBC)   ;ADDR PAGE
               AND  &1F
               OUT  (251),A
               LD   DE,(HKDE)  ;T/S

FRWL:          PUSH BC
               CALL CHKHL
               LD   (BUF),HL
               CALL IYJUMP
               LD   BC,(MSINC) ;USUALLY 0200H
               ADD  HL,BC
               CALL ISECT
               JR   NZ,FRW2

               PUSH HL
               CALL ITRCK
               POP  HL

FRW2:          POP  BC
               DEC  BC
               LD   A,B
               OR   C
               JR   NZ,FRWL

               POP  AF
               OUT  (251),A
               JP   GTIXD      ;NORMAL BUF AGAIN

;CALCULATE ADDRESS SECTION

CALS:          LD   HL,(HKHL)
               XOR  A
               CALL PAGEFORM
               DEC  A          ;PAGE 0-2 FOR SECTS B/C/D
               LD   BC,1       ;1 SECT TO DO

;CALLED BY EVPRM. AHL=ADDR, BC=SECTS

SCASD:         LD   (HKBC),A
               INC  A
               JP   Z,IOOR     ;0-3FFFH ILLEGAL ADDR

               LD   (HKHL),HL  ;OFFSET
               LD   (SVHDR),BC ;SECTORS TO DO

HDUMMY:        RET               ;DUMMY HOOK (WAS S:)

;CALLED BY DIR, DIR$

FDFSR:         LD   A,"*"
               LD   (NSTR1+1),A ;NULL NAME
               LD   (NSTR1+2),A ;ENSURE NOT "." IF DIR$

;GET DEFAULT LETTER AND NUMBER

GTDEF:         CALL NRRD
               DEFW DEVL

               LD   (LSTR1),A

;GET DEFAULT NUMBER

GTDD:          CALL NRRD
               DEFW DEVN

               AND  A
               JR   Z,GTDF2

               CP   RDLIM+1
               JR   C,GTDF3    ;JR IF LEGAL DRIVE
                               ; - ELSE PROB. T SPEED

GTDF2:         LD   A,(ODEF)   ;USE "OTHER" DEFAULT

GTDF3:         PUSH HL
               LD   (ODEF),A
               CALL CODN       ;CONVERT DRV NUM
               CALL DRSET
               POP  HL
               RET

;EVALUATE FILE INFORMATION AT (HL). 14 CHARS IF HL=NSTR1+1

EVFILE:        CALL GTDEF
               LD   (SVHL),HL
               LD   A,(HL)
               CP   &FF
               JR   NZ,EVFL0   ;JR UNLESS NULL NAME

               LD   A,"T"      ;LOAD "" BECOMES LOAD "T:"
               LD   (HL),A
               INC  HL
               LD   (HL),":"
               DEC  HL

EVFL0:         AND  &DF

;CHECK FOR FIRST DIGIT

EVFL1:         LD   C,A
               INC  HL
               LD   A,(HL)
               CP   ":"
               JR   Z,EVFL3

               SUB  &30
               CP   10
               JR   NC,EVFL4

;CHECK FOR SECOND DIGIT

               LD   D,A        ;SAVE FIRST DIGIT
               LD   A,C        ;FIRST CHAR
               CP   "D"
               INC  HL
               LD   A,(HL)     ;CHAR AFTER DIGIT
               JR   NZ,EVFL12  ;JR IF NOT "D1XXXXX"

               CP   ":"
               JR   NZ,EVFL11  ;JR IF NOT "D1:XXXX"

               CALL C11SP
               JR   NZ,EVFL2   ;JR IF NOT "D1:       "

               JR   EVFL14

EVFL11:        CP   " "
               JR   NZ,EVFL12  ;JR IF NOT E.G. "D2 XXXXX"

               CALL C11SP
               JR   NZ,EVFL4   ;JR IF NOT "D1      "

EVFL14:        LD   (HL),":"
               INC  HL
               LD   (HL),"*"   ;E.G. "D1" OR "D1:" BECOME "D1:*"
               DEC  HL         ;PT TO ":"
               JR   EVFL2

EVFL12:        SUB  &30
               CP   10
               JR   NC,EVFL4   ;JR IF 2ND DIGIT NOT FOUND

;EVALUATE 2 DIGIT NUMBER
               LD   E,A
               LD   A,D
               ADD  A,A
               ADD  A,A
               ADD  A,D
               ADD  A,A
               ADD  A,E
               LD   D,A
               INC  HL

;CHECK FOR ":"
               LD   A,(HL)
               CP   ":"
               JR   NZ,EVFL4

EVFL2:         LD   A,D
               CALL CODN

EVFL3:         LD   A,C
               LD   (LSTR1),A  ;LETTER
               INC  HL         ;SKIP ":"
               JR   EVFL5

EVFL4:         LD   HL,(SVHL)

;FILE NAME START

EVFL5:         LD   DE,(SVHL)
               AND  A
               SBC  HL,DE
               LD   (TEMPW3),HL ;STORE NO. OF CHARS TRIMMED OFF
               ADD  HL,DE       ; FRONT (0-4)
               LD   BC,10
               LD   DE,NSTR1+1
               LDIR            ;LEFT JUSTIFY NAME (OR COPY TO
               LD   B,4        ; NSTR1+1)
               CALL LCNTA      ;PAD WITH 4 SPACES...NOT NEEDED..
               LD   A,(DSTR1)

DRSET:         LD   (DRIVE),A
               CALL GCDIA      ;GET CDIRP ADDR
               LD   A,(HL)
               LD   (CDIRT),A
               RET

;CHECK FOR 11 SPACES

C11SP:         LD   B,11
               PUSH HL

C11LP:         INC  HL
               LD   A,(HL)
               CP   " "
               JR   NZ,C11E

               DJNZ C11LP

C11E:          POP  HL
               RET


;--------------
Part_MOVE:


;MOVE COMMAND

MOVE:          CALL EVMOV
               CP   142        ;"TO"
               JP   NZ,REP0    ;"NONSENSE"

               CALL EXDAT
               CALL EVMOV
               CALL EXDAT
               CALL CEOS
               CALL SETF2      ;"MOVE"
               XOR  A
               OUT  (251),A
               LD   A,MIN
               LD   (FSTR1),A
               CALL OPMOV      ;OPEN "IN" TEMP CHANNEL - GET
                               ; NSTR1=ADDR
               JP   C,SNOP     ;IF OPMOV RETURNS C WHEN CREATING
                               ; "IN", IT WAS DUE TO A NON-OPEN
                               ; STREAM BEING USED
               LD   A,(FSTR1)  ;DRIVE NO. FOR 2ND CHANNEL IS
               LD   (NSTR2),A  ; SAME AS FIRST??
               CALL EXDAT
               LD   A,MOUT
               LD   (FSTR1),A
               LD   IX,DCHAN
               LD   HL,FLAG3
               RES  4,(HL)     ;SO OFSM NORMAL (PREVIOUS OPMOV
                               ; MIGHT HAVE SET IT)
               CALL OPMOV      ;OPEN "OUT" TEMP CHANNEL
               JR   NC,MOVA    ;JR IF OK, IF "OVERWRITE?" & "N"
                               ; OR STREAM NOT OPEN THEN CONT.
               PUSH AF         ;Z IF STRM NOT OPEN
               LD   IX,(NSTR2) ;ADDR OF FIRST CHANNEL (IN NSTR1
               LD   BC,FS      ; AFTER EXDAT)
               ADD  IX,BC
               LD   A,(IX+4)
               CP   "D"+&80
               JR   NZ,MVNRC   ;ONLY RECLAIM FIRST CHANNEL IF IT
                               ; IS TEMP "D"
               CALL DECSAM
               CALL RCLMX

MVNRC:         POP  AF
               JP   Z,SNOP     ;STREAM NOT OPEN ERROR

               RET

;RECLAIM A CHANNEL PTED TO BY IX IN SECT C

RCLMX:         LD   C,(IX+9)
               LD   B,(IX+10)
               PUSH BC
               PUSH IX
               POP  HL
               CALL CMR
               DEFW JRECLAIM

               POP  BC
               RET

MOVA:          CALL EXDAT      ;FIRST CHAN ADDR IN NSTR1 AGAIN
               CALL TOSCQ
               JP   NZ,MOVE1   ;JR IF NOT MOVE D TO S/T/

               XOR  A
               LD   (TVFLAG+FS),A  ;NOT AUTO-LIST
               INC  A
               LD   (INQUFG+FS),A  ;IN QUOTES SO NO KEYWORDS
               LD   A,D
               AND  &1F
               CP   &0A
               JR   Z,MOVE1    ;JR IF OPENTYPE

               PUSH AF
               LD   B,9
               CALL MOVJ       ;JUNK HEADER
               POP  AF
               CP   &10
               JR   NZ,MOVCD   ;IF NOT PROGRAM, SUPPRESS
                               ; UNPRINTABLES
               XOR  A
               LD   (INQUFG+FS),A  ;NOT IN QUOTES - KEYWORDS ON

MOVLN:         CALL MOVRC      ;LINE NO. MSB
               CP   &FF
               JR   Z,MEOF     ;END NOW IF END OF PROG
                               ; - IGNORE VARS
               PUSH AF
               CALL MOVRC      ;LINE NO. LSB
               LD   HL,(NSTR2)
               LD   (CURCHL+FS),HL
               POP  HL
               LD   L,A        ;HL=LINE NO
               CALL PNUM5
               LD   B,2
               CALL MOVJ       ;JUNK LINE LEN DATA

MVSLP:         CALL MOVRC
               CP   &0E
               CALL Z,MOVJ6    ;JUNK 5-BYTE FORMS

               PUSH AF
               CALL MOVWC
               POP  AF
               CP   &0D
               JR   NZ,MVSLP

               JR   MOVLN

MOVJ6:         LD   B,6

MOVJ:          PUSH BC
               CALL MOVRC
               POP  BC
               DJNZ MOVJ
               RET

;MOVE NON-OPENTYPE, NON-PROGRAM TO S/K/T

MOVCD:         CALL MOVRC      ;READ CHAR
               JR   NC,MEOF    ;JR IF EOF

               BIT  7,A
               JR   Z,MCD1

               LD   C,A
               LD   A,(MSFLG)
               AND  A
               LD   A,C
               JR   Z,MCD0     ;JR IF >127 TO BE INVERTED

               CP   &FF
               JR   NZ,MCD2
               JR   MCD15

MCD0:          LD   HL,INVERT+FS
               LD   (HL),&FF   ;INVERSE

MCD1:          AND  &7F
               CP   &20
               JR   NC,MCD2

MCD15:         LD   A,(MSUPC)  ;USUALLY "."

MCD2:          CALL MOVWC      ;WRITE PRINTABLE CHAR
               XOR  A
               LD   (INVERT+FS),A
               JR   MOVCD

MOVE1:         CALL MOVRC      ;READ CHAR
               JR   NC,MEOF    ;JR IF EOF

               CALL MOVWC      ;WRITE CHAR
               JR   MOVE1

MEOF:          XOR  A
               LD   (FLAG3),A
               CALL EXDAT
               CALL CLMOV
               CALL EXDAT
               CALL CLMOV
               JP   CLTEMP

;DISC TO SCREEN? Z IF SO (D TO S/K/P). EXIT WITH D=SRC FILE TYPE

TOSCQ:         LD   HL,(NSTR1)
               LD   BC,FS+4
               ADD  HL,BC      ;SYS PAGE IS AT 8000H. PT TO CHAN
               LD   A,(HL)     ; LETTER
               CP   "D"+&80
               RET  NZ         ;RET IF DISC NOT SRC

               LD   BC,FFSA-4
               ADD  HL,BC
               LD   D,(HL)     ;SRC FILE TYPE
               LD   HL,(NSTR2)
               LD   BC,FS+4
               ADD  HL,BC      ;PT TO CHAN LETTER
               LD   A,(HL)
               CP   "S"
               RET  Z

               CP   "P"
               RET  Z

               CP   "K"
               RET

;MOVE - READ CHAR
;EXIT: CY IF GOT CHAR IN A, NC IF EOF

MOVRC:         LD   HL,(NSTR1)
               LD   (CURCHL+FS),HL

MOVRC2:        LD   HL,(CURCHL+FS)
               LD   DE,FS+2
               ADD  HL,DE      ;SYS PAGE IS AT 8000H
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               EX   DE,HL      ;HL=INPUT ADDRESS
               LD   A,H
               CP   &4B
               JR   Z,DOSIP    ;JR IF DOS

               LD   (MTARG),HL
               CALL CMR
MTARG:         DEFW 0

               JR   GIPC

DOSIP:         CALL MCHRD

GIPC:          RET  C          ;RET IF GOT CHAR

               JR   Z,MOVRC2   ;LOOP UNLESS EOF

               RET

;MOVE - WRITE CHAR IN A

MOVWC:         LD   HL,(NSTR2)
               LD   (CURCHL+FS),HL
               LD   DE,FS+1
               ADD  HL,DE
               LD   D,A
               LD   A,(HL)
               CP   &4B
               LD   A,D
               JP   Z,MCHWR

               JP   PNT

;EVALUATE A MOVE SYNTAX

EVMOV:         CALL GTNC
               CP   "#"
               JP   Z,EVSRM

EVSYN:         CALL EVNAM
               RET  Z          ;RET IF SYNTAX TIME

EVSY2:         PUSH AF
               CALL EVFINS
               POP  AF
               RET

HEVSY:         CALL HEVNAM
               JR   EVSY2

;OPEN A MOVE CHANNEL. RETURN NSTR1=CHANNEL ADDR IN CHANNELS
;(SECT B) CY IF ERROR, (AND Z IF STREAM NOT OPEN), NC IF OK

OPMOV:         LD   A,(SSTR1)
               INC  A
               JR   Z,OPMV1    ;JR IF NOT A STREAM

               DEC  A
               CALL STRMD
               LD   D,A
               LD   A,B
               OR   C
               SCF
               RET  Z          ;RET IF STREAM NOT OPEN, WITH CY

               LD   A,D
               CALL CMR
               DEFW STREAM     ;SET STREAM

               LD   IX,(CURCHL+FS)
               JR   OPMV2

OPMV1:         LD   A,(LSTR1)
               AND  &DF
               CP   "D"
               JP   NZ,REP0

               CALL CKDRV
               CALL OPEND
               LD   A,1
               INC  A        ;NZ = ERROR NOT DUE STREAM NOT OPEN
               RET  C        ;RET IF ABORTED (OVERWRITE?+N)

               LD   A,(NSTR1)
               LD   (FSTR1),A
               LD   BC,-FS
               ADD  IX,BC

OPMV2:         LD   (NSTR1),IX
               AND  A          ;NC=OK
               RET

;CLOSE A MOVE CHANNEL

CLMOV:         LD   A,(SSTR1)
               INC  A
               RET  NZ         ;RET IF THERE WAS A STREAM
                               ; - DON'T CLOSE IT
               LD   IX,(NSTR1)
               LD   BC,FS
               ADD  IX,BC

;ENTRY: IX POINTS TO CHANNEL IN SECT C

DELD:          PUSH IX
               POP  HL
               LD   DE,-FS
               ADD  HL,DE      ;CORRECT TO SECT B, LIKE CHANS
               LD   DE,(CHANS+FS)
               OR   A
               SBC  HL,DE
               INC  HL
               LD   (SVTRS),HL ;CHANNEL DISP
               JP   CLRCHD

;RECLAIM TEMPORARY CHANNELS

CLTEMP:        LD   IX,(CHANS+FS)
               LD   DE,6*5+FS
               ADD  IX,DE

CLTM1:         LD   A,(IX+0)
               CP   &0D
               RET  Z

               LD   A,(IX+4)
               CP   "D"+&80
               JR   NZ,CLTM2

               CALL DELD
               JR   CLTEMP

CLTM2:         CALL BITF1
               JR   Z,CLTM3    ;JR IF NOT CLEAR #

               CALL RCLMX
               JR   CLTEMP

CLTM3:         LD   E,(IX+9)
               LD   D,(IX+10)
               ADD  IX,DE
               JR   CLTM1

BACKUP:        CALL GTNC
               CALL FLTOFL     ;EVAL FILE TO FILE
               CALL COBUS      ;EVFILES, SET SINGLE-DISC FLAG
               CALL CKDRV      ; AS NEEDED
               LD   HL,&8000
               LD   (HKHL),HL  ;SRC/DEST
               CALL CLSAM      ;CLEAR SAM
               LD   A,&20
               CALL FDHR       ;CREATE SAM
               LD   B,195      ;BYTES IN SAM
               LD   HL,SAM+194

SBKSL:         LD   A,(HL)
               AND  A
               JR   NZ,BKU2  ;EXIT WHEN 1ST USED MAP BYTE FOUND

               DEC  HL
               DJNZ SBKSL

BKU2:          LD   L,B
               LD   H,0
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,HL      ;HL=SECTORS USED
               LD   BC,40
               ADD  HL,BC      ;ALLOW FOR DIRECTORY TKS 0-3
               INC  B          ;NZ
               LD   DE,&0001
               LD   (HKDE),DE  ;START AT T0/S1

BKUL:          PUSH DE         ;T/S TO WRITE TO
               PUSH HL         ;SECTS LEFT TO DO
               CALL Z,TSPCE2   ;PROMPT FOR SRC DISC IF 1-DRIVE
                               ; BACKUP (Z IF ENTRY FROM LOOP)
               CALL FFPG
               LD   A,B
               AND  A
               JP   Z,REP35    ;ERROR IF NO PAGES FREE

               LD   (HKBC),DE  ;HKBC=PAGE
               LD   L,A
               LD   H,0
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,HL
               EX   DE,HL      ;DE=SECTORS FREE AT 32/PAGE
               POP  HL         ;LEFT TO DO
               SBC  HL,DE      ;SUB LEFT TO DO,FREE
               PUSH HL         ;NEW LEFT TO DO, UNLESS LAST PASS
               JR   C,BKU3     ;JR IF CAN FINISH IN THIS PASS
               JR   Z,BKU3

               EX   DE,HL      ;HL=FREE
               JR   BKU4

BKU3:          ADD  HL,DE      ;HL=LEFT TO DO
               CALL SETF1      ;"LAST PASS"

BKU4:          LD   (SVHDR),HL ;IX=SECTS TO DO
               LD   A,(DSTR1)
               LD   (HKA),A    ;DRIVE
               CALL HFRSAD     ;READ SECTS
               POP  HL         ;LEFT TO DO
               EX   (SP),HL
               LD   (HKDE),HL  ;T/S TO WRITE TO NEXT
               PUSH DE         ;T/S REACHED ON SRC
               CALL TSPCE1     ;PROMPT FOR DEST DISC IF 1-DRIVE
               CALL EXDAT      ; BACKUP
               LD   A,(DSTR1)
               LD   (HKA),A    ;DRIVE
               CALL BITF6
               JR   NZ,BKU5    ;JR IF NOT FIRST PASS

               LD   A,(HKBC)
               OUT  (251),A
               LD   HL,&80FF
               CALL FESE2      ;SET RND WORD AND NAME
               CALL SETF6      ;"NOT FIRST PASS"

BKU5:          CALL HFWSAD     ;WRITE SECTS
               POP  HL
               LD   (HKDE),HL  ;T/S TO READ FROM NEXT IF MORE
                               ;PASSES, DE=T/S TO WRITE TO NEXT
                               ;(DIFFERENT, IF RAMDISK)
               CALL EXDAT
               POP  HL         ;LEFT TO DO OR JUNK
               CALL BITF1
               JR   Z,BKUL     ;LOOP UNTIL LAST PASS

               RET


;--------------
Part_SER2:


;HKHL=ADDR OF MSB OF PTR IN STRMS (IN SECTION B)
;HKC=FILE NAME LEN
;HKDE=FILE NAME START (USUALLY IN BUFFER)
;CHAD PTS TO 0D/: OR IN OR OUT

HOPEN:         CALL GCHR
               CALL OPSR
               CALL HEVSY
               LD   HL,(HKHL)
               DEC  HL
               LD   BC,&5C16
               AND  A
               SBC  HL,BC
               LD   A,L
               SRL  A
               LD   (SSTR1),A
               JR   OPEX

OPSR:

Fix_L686C_42:  ;4.2 Cut

Fix_L686C_43:  ;4.3 Add
               LD   (FSTR1),A

;FixEnd
               CP   &0D
               RET  Z

               CP   ":"
               RET  Z

               CP   &FF
               JR   NZ,OPSR0

               CALL GTNC
               LD   C,MIN      ;USE ZX "IN" CODE
               CP   &60        ;FF 60 IS "IN" CODE
               JR   Z,OPSR1

               LD   C,MRND     ;ZX "RND" CODE
               CP   &3C        ;FF 3C IS "RND" CODE
               JR   Z,OPSR1

OPSR0:         LD   C,MOUT     ;ZX OUT CODE
               CP   &E0        ;OUT TOKEN
               JP   NZ,REP0

OPSR1:         LD   A,C
               LD   (FSTR1),A
               JP   GTNC

;OPEN # COMMAND SYNTAX ROUTINE

OPEN:          CALL CFSO
               CALL Z,REMFP  ;REMOVE ANY FP FORMS IN SYNTAX TIME

               CALL GTNC
               CP   "#"
               JP   NZ,OPNDIR

               CALL EVSRM      ;SKIP, EVAL STREAM
               CALL SEPARX     ;INSIST ,/; (SKIPPED) OR QUOTE

               CALL EVSYN      ;NAME
               CALL OPSR       ;DEAL WITH CR/:/IN/OUT
               CALL PLNS       ;PLACE NEXT STAT ADDR
               CALL CEOS

OPEX:          LD   A,(SSTR1)
               CALL STRMD      ;BC=CURRENT DISP IN STRMS
               CP   4
               JR   C,OPEN25   ;DON'T WORRY IF STREAMS 0-3
                               ; ALREADY OPEN
               LD   HL,5*5+1   ;DISP OF LAST STANDARD CHANNEL
               AND  A
               SBC  HL,BC
               JP   C,REP30    ;"STREAM USED" IF STREAM OPEN TO
                               ; NON-STANDARD NOW
OPEN25:        LD   A,(LSTR1)
               AND  &DF
               CP   "D"
               JP   NZ,REP0    ;"NONSENSE"

               CALL CKDRV
               LD   A,10
               LD   (NSTR1),A  ;FILE TYPE

;OPEN A STREAM TO 'D' CHANNEL

OPDST:         LD   A,(SSTR1)
               ADD  A,A
               LD   HL,&5C16+FS
               LD   E,A
               LD   D,0
               ADD  HL,DE
               PUSH HL         ;ADDR OF PTR IN STRMS
               CALL OPEND
               POP  DE
               RET  C          ;RET IF ERROR

   ;         * BIT  0,(IX+MFLG)
    ;        * JR   Z,OPDST1   ;JR IF OPEN IN OR RND

    ;        * CALL COMMP
    ;        * IN   A,(C)
    ;        * BIT  6,A
    ;        * JR   Z,OPDST1   ;NEED TO WRITE TO GET THIS BIT TO
                               ; WORK!
    ;        * CALL RCLMX
     ;       * JP REP23

OPDST1:        RES  7,(IX+4)   ;"D" NOT "D"+80H - PERM CHANNEL
               EX   DE,HL
               LD   (HL),E
               INC  HL
               LD   (HL),D     ;STRM PTR
               RET

;OPEN A  'D' DISC CHANNEL

OPEND:         XOR  A
               OUT  (251),A
               LD   IX,(CHANS+FS)
               LD   DE,6*5+FS
               ADD  IX,DE      ;SKIP 6 STANDARD CHANNELS

OPND1:         LD   A,(IX+0)
               CP   &0D
               JR   Z,OPND4    ;JR IF CHANS TERMINATOR FOUND
                               ; - NEW CHANNEL
;FOUND AN OPEN CHANNEL

               LD   A,(IX+4)
               AND  &5F
               CP   "D"
               JR   NZ,OPND3

               LD   A,(DSTR1)
               CP   (IX+MDRV)
               JR   NZ,OPND3

;CHECK NAME OF CHANNEL

               PUSH IX
               POP  HL
               LD   DE,NAME
               ADD  HL,DE
               EX   DE,HL
               LD   HL,NSTR1+1
               LD   B,10

OPND2:         LD   A,(DE)
               XOR  (HL)
               AND  &DF
               JR   NZ,OPND3

               INC  HL
               INC  DE
               DJNZ OPND2

               JP   REP31      ;"CHANNEL USED"

;GET THE LENGTH OF CHANNEL

OPND3:         LD   E,(IX+9)
               LD   D,(IX+10)
               ADD  IX,DE
               JR   OPND1

;IT IS A NEW CHANNEL - NOW TEST DIRECTORY FOR FILENAME

OPND4:         PUSH IX
               LD   A,&10
               CALL FDHR       ;Z IF FOUND
               POP  HL         ;ADDR OF CHANS TERMINATOR
               PUSH HL         ; - LOCN FOR NEW CHAN
               LD   A,(FSTR1)
               JP   NZ,OPND45  ;JR IF NOT FOUND

               PUSH AF
               CALL CRMCH      ;CREATE CHANNEL
               CALL POINT      ;HL=DIR ENTRY  (USES IX=DCHAN)
               POP  AF

               POP  IX         ;CHANS PTR
               PUSH HL         ;DIR ENTRY PTR
               LD   C,4+0      ;BITS 1 AND 0 SHOW READ, BIT 2
               CP   MIN        ; SHOWS EXISTS
               JR   Z,OPND44

               CP   MOUT
               JP   Z,REP19    ;WRITING A READ FILE IF EXIST
                               ; AND OUT
               CP   MRND
               JR   NZ,OPND44  ;JR IF DEFAULT - IN

               LD   C,4+2      ;BITS 1 AND 0 SHOW RND, BIT 2
                               ; SHOWS EXISTS
OPND44:        LD   A,C
               CALL RAMST
               LD   DE,FFSA
               LD   BC,&0100
               PUSH IX
               POP  HL         ;CHANNEL
               ADD  HL,DE
               EX   DE,HL      ;DE=TYPE/NAME ETC DEST IN CHANNEL
               POP  HL         ;DIR ENTRY
               PUSH HL         ;DIR ENTRY
               PUSH DE         ;TYPE IN CHANNEL
               LDIR    ;COPY 256 BYTES FROM DIR ENTRY TO CHANNE
                       ; INCLUDES TYPE, NAME, ????, BAM, LENGTH
               POP  DE         ;TYPE IN CHANNEL
               POP  HL         ;DIR ENTRY
               LD   BC,11
               LD   A,(HL)
               LD   (NSTR1),A
               LDIR

               INC  HL
               INC  HL
               LD   D,(HL)
               INC  HL
               LD   E,(HL)
               PUSH DE          ;FIRST T/S
               XOR  A
               LD   (IX+CNTL),A
               LD   (IX+CNTH),A ;FILE'S SECTOR=0 TO START WITH
               CALL RSADM    ;INC SECTOR, LOAD AND MARK 1ST SECT
               CALL OPND8    ;WITH T/S. GET STREAM OFFSET IN HL
               POP  DE
               PUSH HL
               LD   (IX+FTRK),D ;TRK
               LD   (IX+FSCT),E ;SECT

               PUSH IX
               POP  HL
               LD   BC,FFSA+&F2
               ADD  HL,BC       ;PT TO EXECUTION ADDR OF FFFFF
                                ; UNLESS G+DOS
               LD   BC,9        ;HDR LEN
               LD   A,(IX+FFSA) ;TYPE
               BIT  6,A         ;PROTECT BIT
               JR   Z,OPNDUP

               LD   (IX+MFLG),4 ;ENSURE PROTECTED FILES ARE
                                ; READ-ONLY
OPNDUP:        AND  &1F
               CP   &0A
               JR   NZ,NOTF    ;JR IF NOT OPENTYPE

               LD   A,(HL)
               INC  A
               JR   NZ,OOTF2   ;JR IF G+DOS OPEN-TYPE FILE

               LD   C,A        ;C=0. NO HDR IN OPENTYPE FILE
               DEC  A          ;A=FF SO NO JR

NOTF:          CP   &10
               JR   C,OOTF     ;JR IF ZX TYPE

               DEC  HL       ;PT TO LEN DATA FOR ALL FILE TYPES
               LD   D,(HL)
               DEC  HL
               LD   E,(HL)
               DEC  HL
               LD   A,(HL)
               EX   DE,HL
               CALL AHLNX      ;GET 20-BIT LEN
               ADD  HL,BC      ;HDR LEN OR ZERO
               ADC  A,0
               LD   B,A
               JR   GT19B

OOTF:          LD   HL,&C000
               LD   B,L
               CP   5
               JR   Z,GT19B    ;JR IF 48K SNAP

OOTF2:         CALL PTLEN
               DEC  HL
               LD   D,(HL)
               DEC  HL
               LD   E,(HL)
               DEC  HL
               DEC  HL
               LD   B,(HL)
               EX   DE,HL      ;BHL=LEN IN ZX FORM

GT19B:         CALL D510       ;GET HL=LEN MOD 510, HL'=DIV 510
               EX   DE,HL
               CALL PTLEN      ;NC
               LD   (HL),E
               INC  HL
               LD   (HL),D     ;LEN MOD 510
               INC  HL
               EXX
               PUSH HL         ;LEN DIV 510 (LEN CNT)
               EXX
               POP  DE
               LD   (HL),E
               INC  HL
               LD   (HL),D     ;LEN DIV 510
               XOR  A
               LD   (NSTR1+1),A ;NO NAME SO NO "OVERWRITE?"
               PUSH IX
               LD   IX,DCHAN   ;ENSURE CLEARING HITS DCHAN,
                               ; NOT CHANNEL
               CALL SETF4      ;"NO SECTOR NEEDED"
               CALL ROFSM      ;SET UP SAM, CLEAR, INC SAMCNT
               POP  IX

;ENSURE EXTENDING RND FILE NEVER USES EARLIER SECTORS. LOOK BACU
;BACKWARDS THRU SAM, FORCE TO FFS AFTER LAST USED SECTOR SEEN

               LD   HL,SAM+195

CSL1:          DEC  L
               LD   A,(HL)
               AND  A          ;(NC)
               JR   Z,CSL1     ;LOOK FOR LAST BYTE WITH ANY BIT
                               ; (SECT) USED
               INC  L

CSL2:          DEC  L
               LD   (HL),&FF
               JR   NZ,CSL2

               POP  HL         ;STREAM OFFSET
               AND  A          ;NC=OK (NOT NEEDED...)
               RET

;FILE NOT FOUND
;OPEN A NEW "RND" FILE OR AN "OUT" FILE
L6A1C: ;* &6A1F
OPND45:        CP   MIN
               JP   Z,REP26    ;"FILE NOT FOUND" IF "IN" WITH
                               ; NEW FILE
               PUSH AF
               CALL CRMCH
               POP  AF
               POP  IX

               CP   MRND
               LD   A,2
               JR   Z,OPND6    ;JR IF RND (BITS 1-0 = 10)
                               ; BIT 2=0 (NEW FILE)
               DEC  A          ;BITS 1-0 SHOW OUT (01). OUT IS
                               ; DEFAULT.
OPND6:         CALL OPND7      ;SETS FTRK/FSCT, CNT=1
                               ;EXIT WITH HL=STRM PTR
               JP   NC,SETLEN

               RET               ;C=ABORTED

OPND7:         CALL RAMST      ;SET MFLG ETC
               CALL DDEL     ;DELAY IF NOT RAMDISC IN CASE 1ST
                             ;SCT WRITTEN TOO SOON

               CALL ROFSM   ;CREATE SAM FOR DRIVE, INIT AREA AT
                            ;IX, INC SAMCNT. AS SECTORS ARE
                            ;WRITTEN, SAM ACCUMULATES USED SECTS
                            ;FOR DRIVE, FSAM DITTO FOR FILE. BUT
                            ;IF 2 OUT FILES OPEN TO DIFFERENT
                            ;DRIVES, 1 SAM=DISASTER! EVEN 2 OUT
                            ;FILES ON 1 DRIVE MUST ENSURE 2ND
                            ;DOES NOT CLEAR SAM OR BITS
                            ;REFLECTING FSAM ARE LOST!

               JR   NC,OPND75  ;JR IF NOT ABORT DUE "OVERWRITE"
                               ; AND "N"
               LD   BC,787
               PUSH IX
               POP  HL
               CALL CMR
               DEFW JRECLAIM

               SCF               ;"OPEND ABORTED"
               RET

OPND75:        CALL BITF2
               JR   Z,OPND8    ;JR IF NOT MOVE

               PUSH IX
               POP  HL
               LD   BC,239
               ADD  HL,BC
               EX   DE,HL      ;DEST=BYTE 220 IN FFSA (DEST DISC
                               ; CHANNEL)
               LD   HL,(NSTR2)
               LD   B,FS/256
               ADD  HL,BC      ;DITTO FOR SRC (IF DISC CHANNEL)
               LD   BC,36
               LDIR            ;COPY 220-255

;CALCULATE STREAM OFFSET

OPND8:         PUSH IX
               POP  HL
               LD   DE,FS
               AND  A
               SBC  HL,DE
               LD   DE,(CHANS+FS)
               SBC  HL,DE
               INC  HL
               RET               ;NC - OK

;CREATE A 'D'+80H CHANNEL AT HL

CRMCH:         PUSH HL
               XOR  A
               LD   BC,787
               CALL CMR
               DEFW JMKRBIG    ;OPEN SPACE FOR CHANNEL

               POP  DE
               LD   HL,MTBLS
               LD   BC,11
               LDIR            ;O/P,I/P,D+80H,0,0,0,0,LEN
               EX   DE,HL
               LD   BC,787-11

;CLEAR THE NEW CHANNEL AREA

CRMC1:         LD   (HL),0
               INC  HL
               DEC  BC
               LD   A,B
               OR   C
               JR   NZ,CRMC1

               RET

;DISC 'D' CHANNEL DATA

MTBLS:         DEFW &4BA0
               DEFW &4BA9
               DEFB "D"+&80
               DEFW 0
               DEFW 0
               DEFW 787        ;CHAN LEN (IX+9/10)

;HOOK 151 - DISC BLOCK O/P BC FROM DE. FASTER THAN REPEATED USE
;OF RST 08H FOR SINGLE CHARS

HDBOP:         LD   DE,(HKDE)  ;SRC
               LD   HL,(HKBC)  ;COUNT
               BIT  7,D
               JR   NZ,DBOL    ;JR IF SRC IN SECT C

               SET  7,D
               RES  6,D        ;ELSE ADJUST SYS PAGE SRC
               XOR  A          ; (E.G. STR$) TO SECT C
               OUT  (251),A    ;SYS PAGE IN SECT C
                               ; (0013H RESETS ORIG PG AT END)
DBOL:          LD   A,H
               OR   L
               RET  Z

               DEC  HL         ;DEC CHAR COUNT
               PUSH DE
               PUSH HL
               LD   A,(DE)
               CALL MCHWR
               POP  HL
               POP  DE
               INC  DE
               BIT  6,D
               CALL NZ,&3FEB   ;INCPAGDE
               JR   DBOL

;CLOSE# STREAMS COMMAND

CLOSE:         LD   C,"*"
               CALL ISEPX      ;INSIST "*"
               CALL CIEL
               JR   Z,CLOS1  ;JR IF CR/0D - CLOSE * IS CLOSE ALL

               CALL EVSRMX     ;EVAL STREAM NO.
               CALL CEOS
               LD   A,(SSTR1)
               JR   CLSRM

CLOS1:         CALL CEOS
               JR   CLRS1

;CLEAR# STREAMS COMMAND

CLEAR:         LD   C,"#"
               CALL ISEPX
               CALL CEOS
               CALL SETF1      ;"CLEAR#"

CLRS1:         XOR  A

CLRS2:         PUSH AF
               CALL CLSRM
               POP  AF
               INC  A
               CP   16
               JR   C,CLRS2

               CALL CLTEMP
               XOR  A
               LD   (SAMCNT),A
               LD   (FLAG3),A
               RET

HCLOS:         LD   HL,(HKDE)
               LD   BC,&5C16
               AND  A
               SBC  HL,BC
               LD   A,L
               SRL  A
               LD   (SSTR1),A

;CLEAR STREAM AND CLOSE CHANNEL

CLSRM:         CALL STRMD
               LD   A,C
               OR   B
               RET  Z          ;RET IF CLOSED ALREADY

;CLOSE THE STREAM

               LD   (SVTRS),BC
               PUSH HL         ;PTR TO STRMS
               LD   HL,(CHANS+FS)
               DEC  HL
               ADD  HL,BC
               SET  7,H
               RES  6,H        ;PT TO CHANNEL SWITCHED IN SECT C
               EX   (SP),HL    ;HL=PTR TO STRMS

               LD   BC,0
               LD   DE,-&9C1E
               EX   DE,HL
               ADD  HL,DE      ;ADD -9C1EH, STRMS PTR
               JR   C,CLOSE1   ;JR IF STREAM >3

               LD   BC,TABLE+8
               ADD  HL,BC
               LD   C,(HL)
               INC  HL
               LD   B,(HL)

CLOSE1:        EX   DE,HL
               SET  7,H
               RES  6,H        ;ADDR IN STREAMS
               LD   (HL),C
               INC  HL
               LD   (HL),B
               POP  IX         ;CHANNEL START
               LD   A,B
               OR   C
               RET  NZ         ;RET IF STREAM 0-3 JUST CLOSED

;TEST FOR DISC 'D' CHANNEL

               LD   A,(IX+4)
               AND  &5F
               CP   "D"
               RET  NZ

CLRCHD:        LD   A,(IX+MFLG)
               AND  &03
               JR   NZ,CLRC2   ;JR IF ITS NOT AN "IN" FILE

               CALL DECSAM   ;(DECSAM DONE BY SDCM FOR RND/OUT)
               JR   RCLAIM

DECSAM:        LD   A,(SAMCNT)
               AND  A
               RET  Z          ;NO DEC BELOW ZERO!

               DEC  A
               LD   (SAMCNT),A
               RET

CLRC2:         CALL BITF1
               CALL Z,SDCM     ;CALL IF NOT CLEAR #

;RECLAIM THE CHANNEL

RCLAIM:        CALL RCLMX

;CLOSE AND UPDATE STREAM DISP

               XOR  A
               LD   HL,&5C16+FS

RCLM1:         LD   (SVHL),HL
               LD   E,(HL)
               INC  HL
               LD   D,(HL)     ;DISP OF A STREAM
               LD   HL,(SVTRS)
               AND  A
               SBC  HL,DE
               JR   NC,RCLM4   ;JR IF NO NEED TO ALTER IT

               EX   DE,HL
               AND  A
               SBC  HL,BC
               EX   DE,HL
               LD   HL,(SVHL)
               LD   (HL),E
               INC  HL
               LD   (HL),D     ;REDUCED DISP REPLACED
               DEC  DE
               LD   HL,(CHANS+FS)
               ADD  HL,DE
               LD   DE,FS+4
               ADD  HL,DE
               LD   D,A
               LD   A,(HL)
               AND  &5F
               CP   "D"
               LD   A,D
               JR   NZ,RCLM4

               LD   DE,BUFL-4
               ADD  HL,DE
               LD   E,(HL)
               INC  HL
               LD   D,(HL)     ;OLD MBUFF
               EX   DE,HL
               SBC  HL,BC
               EX   DE,HL
               LD   (HL),D
               DEC  HL
               LD   (HL),E     ;ADJUSTED MBUFF

RCLM4:         LD   HL,(SVHL)
               INC  HL
               INC  HL
               INC  A
               CP   16
               JR   C,RCLM1    ;LOOP FOR STREAMS 0-15

               RET

SDCM:          CALL SSDRV
               XOR  A
               LD   (FSLOT),A  ;ENSURE NO USE OF FSLOT
                               ;DE=T/S
               CALL WRIF       ;WRITE CURRENT SECTOR IF IT HAS
               CALL DECSAM     ; BEEN ALTERED
               BIT  5,(IX+MFLG)
               RET  Z          ;RET IF FILE WAS NOT ALTERED - NO
                               ; NEED TO ALTER DIRECTORY ENTRY
               CALL GLEN       ;GET FILE LEN IN AHL, SECTS IN BC
               LD   D,A
               LD   A,(IX+FFSA)
               AND  &1F
               CP   &10
               LD   A,D
               JR   C,SDCM2    ;JR IF OPENTYPE AND ZX FILES

               LD   DE,-9
               ADD  HL,DE
               ADC  A,&FF   ;AHL=AHL-9 (HDR NOT INCLUDED IN LEN)

SDCM2:         LD   (IX+CNTL),C ;OVER-WRITE PTR
               LD   (IX+CNTH),B
               EX   DE,HL
               CALL PTLEN
               DEC  HL
               LD   (HL),D     ;OLD FORMAT MID LEN - SO G+DOS
                               ;CAN READ (UNLESS FIRST DIR
                               ;ENTRY - NAME USES AREA)
               DEC  HL
               LD   (HL),E     ;LOW
               DEC  HL
               DEC  HL
               LD   (HL),A     ;HI - BYTE D2
               EX   DE,HL
               CALL PAGEFORM
               EX   DE,HL
               LD   BC,&001D
               ADD  HL,BC      ;PT TO BYTE EF
               LD   (HL),A
               INC  HL
               LD   (HL),E
               INC  HL
               LD   (HL),D     ;NEW FORM OF LEN
               INC  HL
               LD   (HL),&FF   ;"EXEC ADDR" OF FFXXXX
               PUSH IX
               BIT  2,(IX+MFLG)
               JP   Z,CLOSX    ;JR IF FILE IS A NEW ONE - MAKE
                               ;NEW DIR ENTRY ELSE DEALING WITH
                               ;AN ALTERED EXISTING FILE
               PUSH IX
               POP  HL
               LD   BC,FFSA    ;19
               ADD  HL,BC
               LD   DE,NSTR1
               LD   C,11
               LDIR            ;COPY NAME TO NSTR1

               LD   A,8
               CALL FDHR       ;LOOK FOR FILE ENTRY
               JP   NZ,REP26   ;ERROR IF NOT FOUND

               JP   NCF25      ;UPDATE DIRECTORY - OVERWRITING
                               ; OLD ENTRY

TABLE:         DEFB 1,0,1,0,6,0,16,0

;HOOK ROUTINE TO READ BYTE FROM DISC. USED BY "D" CHANNEL

MCHRD:         CALL MCHIN
               RET  C          ;RET IF GOT CHAR

               RET  Z          ;RET IF NOT EOF

               CALL BITF2
               JP   Z,REP27    ;"EOF" IF NOT A MOVE CMD

               OR   1          ;NC,NZ - EOF
               RET

;DISC 'D' CHANNEL INPUT SUBROUTINE

MCHIN:         IN   A,(251)
               PUSH AF
               XOR  A
               OUT  (251),A
               PUSH IX             ;KEEP FPC HAPPY DURING INKEY$
               LD   HL,TVFLAG+FS
               RES  3,(HL)         ;"NO NEED TO COPY LINE TO LS"
               LD   IX,(CURCHL+FS) ; - NEEDED?
               LD   BC,FS
               ADD  IX,BC
               BIT  0,(IX+MFLG)
               JP   NZ,REP18     ;"READING A WRITE FILE"
                                 ; IF TYPE="OUT"
               CALL CPPTR
               JR   NZ,MCHN1     ;JR IF NOT EOF

               ADD  A,&0D        ;NC+NZ=EOF (WHEN PTR=LEN)
               JR   MCHN2

MCHN1:         CALL LBYT
               SCF                 ;"GOT KEY"

MCHN2:         PUSH AF
               POP  BC           ;AF IN BC FOR TRANSMISSION OUT
               POP  IX
               EX   AF,AF'
               POP  AF
               OUT  (251),A
               EX   AF,AF'
               RET

;HOOK ROUTINE TO WRITE BYTE IN A TO DISC. USED BY "D" CHANNEL

MCHWR:         LD   D,A
               IN   A,(251)
               PUSH AF
               XOR  A
               OUT  (251),A
               LD   IX,(CURCHL+FS)
               LD   BC,FS
               ADD  IX,BC
               LD   A,(IX+MFLG)
               AND  &03
               JP   Z,REP19      ;"WRITING A READ FILE" IF "IN"

               LD   A,(CURCMD+FS)
               CP   &C6          ;VALUE FOR "INPUT"
               JR   Z,MCHW2      ;NO WRITE IF SO

               CALL CPPTR        ;Z IF PTR=LEN
               PUSH AF
               LD   A,D
               CALL NSBYT        ;SAVE BYTE
               SET  3,(IX+MFLG)  ;"SECTOR WRITTEN TO"
               SET  5,(IX+MFLG)  ;"FILE WRITTEN TO"
               POP  AF
               CALL Z,SETLEN     ;COPY PTR TO LEN IF WRITING TO

               LD   A,(&5C4B+FS) ;RESTORE BORDER COLOUR - QUICK
               OUT  (ULA),A

MCHW2:         POP  AF
               OUT  (251),A
               RET

SBCSR:         CALL FNFS
               LD   (HL),D       ;T
               INC  HL
               LD   (HL),E       ;S COMPLETE BUFFER OF 512
               EX   DE,HL

;SWAP THE NEXT TRACK/SECTOR

SWPNSR:        CALL GTNSR
               LD   (IX+NSRH),H
               LD   (IX+NSRL),L
               RET

SBYT:          PUSH BC
               PUSH HL
               PUSH AF
               CALL TFBF     ;HL=ADDR OF WRITE POINT
                             ;BC=RPT (DISP FROM BUFFER START TO
                             ;WRITE POINT) Z IF BC=510
               JR   NZ,SBT1

SBT2:          PUSH DE
               CALL SBCSR    ;PLACE NEXT T/S IN IX+,
                             ; GET DE=CURRENT T/S
               CALL WSAD     ;EXITS WITH HL POINTING TO BUFFER
                             ; START, RPT RESET
               JR   SBT0

;NEW SAVE BYTE TO DISC - SERIAL FILES

NSBYT:         PUSH BC
               PUSH HL
               PUSH AF
               CALL TFBF
               JR   NZ,SBT1  ;JR IF BUFFER NOT FULL

               PUSH DE
               PUSH HL
               CALL CPPTR    ;CP PTR WITH FILE LEN
               POP  HL
               JR   NZ,SBT3  ;JR IF WE ARE NOT AT FILE END

               CALL SBCSR
               CALL SSDRV    ;SELECT DRIVE
               PUSH BC       ;PREV DRIVE
               CALL NWSAD    ;EXIT WITH HL AT BUFFER START
               POP  BC
               LD   A,C
               CALL SSDRV2   ;PREV
               PUSH HL
               LD   D,H
               LD   E,L
               INC  DE
               LD   BC,&01FF
               XOR  A        ;Z
               LD   (HL),A
               LDIR          ;BLANK NEW SECTOR
               POP  HL

SBT3:          CALL NZ,RDWRSR

SBT0:          POP  DE
SBT1:          POP  AF
               LD   (HL),A
               POP  HL
               POP  BC

;INCREMENT RAM POINTER

INCRPT:        INC  (IX+RPTL)
               RET  NZ

               INC  (IX+RPTH)
               RET

SETLEN:        PUSH HL
               CALL PTLEN        ;NC
               LD   A,(IX+RPTL)
               LD   (HL),A
               INC  HL
               LD   A,(IX+RPTH)
               LD   (HL),A
               INC  HL
               LD   A,(IX+CNTL)
               LD   (HL),A
               INC  HL
               LD   A,(IX+CNTH)
               LD   (HL),A
               POP  HL
    ;        * SET  4,(IX+MFLG)  ;"FILE HAS BEEN EXTENDED"
               RET

;CP PTR WITH LEN. Z IF MATCH (AND HL=LEN MSB OF 4)

CPPTR:         CALL PTLEN
               LD   A,(IX+RPTL)
               CP   (HL)         ;HL PTS TO LEN LOW (LIKE RPT)
               RET  NZ

               INC  HL
               LD   A,(IX+RPTH)
               CP   (HL)
               RET  NZ

               INC  HL
               LD   A,(IX+CNTL)
               CP   (HL)
               RET  NZ

               INC  HL
               LD   A,(IX+CNTH)
               SUB  (HL)
               RET                 ;A=0 IF EOF, Z FLAG

;LOAD BYTE FROM DISC

LBYT:          PUSH BC
               PUSH DE
               PUSH HL
               CALL TFBF
               CALL Z,RDWRSR ;CALL IF BUFFER FULL.
                             ;WRITE THIS SECTOR IF ALTERED,
                             ;READ NEXT
LBT1:          LD   A,(HL)
               POP  HL
               POP  DE
               POP  BC
               JR   INCRPT

;PTR/EOF SR
;EXIT: AHL=PTR VALUE, DE=CHANNEL START

PESR:          CALL CMR
               DEFW GETINT

               INC  H
               DEC  H
               JR   NZ,INVST

               CP   16
               JR   NC,INVST

PESR2:         CALL STRMD
               LD   A,B
               OR   C
               JR   Z,SNOP     ;ERROR IF STREAM NOT OPEN

               DEC  BC
               LD   HL,(CHANS+FS)
               ADD  HL,BC
               LD   BC,FS
               ADD  HL,BC      ;CHANNEL START
               PUSH HL
               POP  IX
               LD   A,(IX+4)
               CP   "D"
               JP   NZ,REP10   ;"INVALID DEVICE"

FPTR:          LD   C,(IX+CNTL)
               LD   B,(IX+CNTH) ;PTR IN 510'S
               LD   L,(IX+RPTL)
               LD   H,(IX+RPTH) ;PTR MOD 510
               JP   M510

INVST:         CALL DERR
               DEFB 21          ;"INVALID STREAM NUMBER"

SNOP:          CALL DERR
               DEFB 47          ;"STREAM IS NOT OPEN"

;STREAM DISPLACEMENT
;ENTRY: A=STREAM NO.
;EXIT: A UNCHANGE, HL=ADDR IN STREAMS, BC=PTR VALUE FROM STREAMS

STRMD:         PUSH AF
               ADD  A,A
               LD   C,A
               XOR  A
               OUT  (251),A
               LD   B,A
               LD   HL,&5C16+FS ;STREAM ZERO
               ADD  HL,BC
               LD   C,(HL)
               INC  HL
               LD   B,(HL)      ;BC=CURRENT DISP IN STRMS
               DEC  HL
               POP  AF
               RET

;SET IX+MFLG, IX+BUFL/H USING A AND BC
;(SET TYPE=IN (0), OUT (1) OR RND (2), SET BUFFER LOCN=IX+BC)

RAMST:         PUSH HL
               LD   (IX+MFLG),A
               PUSH IX
               POP  HL
               LD   BC,WRRAM
               ADD  HL,BC
               LD   (IX+BUFL),L
               LD   (IX+BUFH),H
               LD   A,(DSTR1)
               LD   (IX+MDRV),A
               POP  HL
               RET

;POINT #S,X OR POINT #S,OVER X

POINTC:        LD   C,"#"
               CALL ISEPX
               CALL EVSRMX     ;EVAL STREAM
               CALL SEPARX     ;,/
               CP   &A6        ;OVER
               JP   Z,PTREC

               CALL EVBNUM     ;BC=X MOD 64K, DE=X DIV 64K
               CALL CEOS
               PUSH BC
               PUSH DE
               LD   A,(SSTR1)
               CALL PESR2    ;PT IX TO CHANNEL,CHECK OPEN "D" IN
               CALL GLEN
               POP  BC         ;X DIV 64K

Fix_L6DCE_42:  ;4.2 Cut

Fix_L6DD1_43:  ;4.3 Add
               LD   B,C

;FixEnd
               POP  DE
               SBC  HL,DE
               SBC  A,B
               JR   C,REP27H

               EX   DE,HL
               CALL D510
               PUSH HL
               EXX
               LD   (IX+CNTL),L
               LD   (IX+CNTH),H
               EXX

;FIND T/S OF DE'TH SECTOR IN A FILE, USING FSAM

;ENTRY: IX PTS TO FILE AREA, HL'=SECT (1=FIRST SECTOR)
;EXIT: D=TRACK, E=SECT, NZ, OR Z=NOT FOUND (NOT THAT MANY
;SECTORS IN FSAM). USES AF,BC,DE,HL, HL', D'

FITS:          PUSH IX
               POP  HL
               LD   BC,FSAM
               ADD  HL,BC      ;PT HL TO FSAM
               PUSH HL
               LD   B,195

FTSL:          LD   A,(HL)
               INC  HL
               AND  A
               JR   NZ,FTS2    ;JR IF MAP BYTE HAS ANY USED
                               ; SECTORS IN IT
FTDB:          DJNZ FTSL

REP27H:        JP   REP27      ;EOF IF RAN OUT OF MAP BITS

FTS2:          LD   E,8        ;8 BITS TO EXAMINE

FTSBL:         RRA
               JR   NC,FTS3    ;JR IF SECTOR NOT USED

               EXX
               DEC  HL         ;DEC 'DESIRED SECTOR' COUNTER
               LD   D,A
               LD   A,H
               OR   L
               LD   A,D
               EXX
               JR   Z,FTS4     ;JR IF GOT THE ONE WE WANTED

FTS3:          DEC  E
               JR   NZ,FTSBL   ;LOOP FOR ALL BITS

               JR   FTDB       ;NEXT BYTE IF ALL BITS DONE

FTS4:          POP  BC         ;FSAM START
               SCF
               SBC  HL,BC      ;HL=MAP BYTE
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,HL      ;MAP SECTOR (DIV 8)
               LD   D,0
               LD   A,8
               SUB  E
               LD   E,A
               ADD  HL,DE      ;D=0, E=BIT (7-0 FOR BIT 7-0 I
               LD   D,&03      ;FIRST DATA TRACK=4

FTSSL:         CALL FNS5       ;INC D, CHECK FOR END OF SIDE
               LD   BC,-10
               ADD  HL,BC
               JP   C,FTSSL

               SBC  HL,BC
               LD   H,D
               INC  L

               CALL GTNSR
               EX   DE,HL
               AND  A
               SBC  HL,DE
               JR   Z,PTLX5

               PUSH DE
               CALL WRIF
               POP  DE
               CALL RSADM2

PTLX5:         POP  DE

SETRPT:

Fix_L6E38_42:  ;4.2 Chg
;*             JP   L6EBC

Fix_L6E3C_43:  ;4.3 Chg
               LD   (IX+RPTL),E
               LD   (IX+RPTH),D
               RET

;FixEnd

PTREC:         CALL EVNUMX
               CALL CEOS
               PUSH BC
               LD   A,(SSTR1)
               CALL PESR2      ;AHL=CURRENT PTR
               CALL GRPNT      ;HL=BUFF PTR, BC=DISP IN BUFFER
               EX   DE,HL      ; (0-509)
               LD   HL,510
               AND  A
               SBC  HL,BC
               LD   B,H
               LD   C,L        ;BC=BYTES LEFT IN BUFFER
               EX   DE,HL      ;HL=BUFF PTR
               POP  DE         ;ITEMS TO SKIP
               LD   A,D
               OR   E
               RET  Z

PTRCSL:        LD   A,B
               OR   C
               JR   Z,PTRSL    ;IF SECTOR SEARCHED, ENTER SECT
                               ; LOOP (NEVER ON FIRST PASS)
               LD   A,(DELIM)
               CPIR
               JR   NZ,PTRSL   ;JR IF DELIM NOT FOUND - SEARCH
                               ; NEXT SECT
               DEC  DE
               LD   A,D
               OR   E
               JR   NZ,PTRCSL  ;LOOP UNTIL SKIPPED RIGHT NUMBER
                               ;OF DELIMS IF BC=0, USE NEXT SECT
               LD   HL,510
               JR   PTRCE

PTRSL:         LD   (TEMPW1),DE  ;DELIM COUNT
               XOR  A
               LD   (TEMPW2+1),A ;"NOT FOUND"
               CALL WRIF
               CALL GTNSC
               CALL SRSAD   ;READ NEXT SECTOR, LOOKING FOR DELIM
               CALL MSECT        ;MARK SECT WITH CUR T/S
               CALL ICNT
               LD   DE,(TEMPW3)  ;COUNTER
               PUSH HL           ;BUFFER START
               PUSH IX
               POP  HL           ;CHANNEL START
               LD   BC,NTRK
               ADD  HL,BC        ;PT TO NTRK
               LD   BC,(TEMPW2) ;RECORD PTR:
                             ;IF 00XX, NO MATCH
                             ;IF NTRK OR LESS, REAL MATCH
                             ;IF NSECT, SPURIOUS MATCH ON TRK
                             ;IF NSECT+1, SPURIOUST MATCH ON SCT
               LD   A,B
               AND  A
               JR   Z,PTRC3    ;JR IF DE NOT COUNTED DOWN YET,
                               ; NO ADDR IN BC
               SBC  HL,BC
               ADD  HL,BC
               JR   NC,PTRCOK  ;JR IF OK

PTRC3:         POP  AF         ;JUNK BUFFER START
               LD   A,(DELIM)
               CP   (HL)
               JR   NZ,PTRC4

               INC  DE         ;CORRECT COUNTER FOR SPURIOUS
                               ; MATCH ON TRACK
PTRC4:         INC  HL
               CP   (HL)
               JR   NZ,PTRSL   ;JR IF SECTOR WASN'T A SPURIOUS
                               ; MATCH
               INC  DE         ;CORRECT COUNTER
               JR   PTRSL

PTRCOK:        POP  BC          ;BUFFER START
               LD   HL,(TEMPW2) ;PTR WITHIN BUFFER

PTRCE:         SBC  HL,BC   ;RECORD PTR-BUFFER START (0001-01FE)
                            ;(WHEN ENTRY AT PRCE, HL=510,
                            ;BC=BYTES LEFT)
               EX   DE,HL
               LD   HL,510
               AND  A
               SBC  HL,DE

Fix_L6EBA_42:  ;4.2 Chg
;*             JR   Z,PTRC5
;*L6EBC:       LD   (IX+&0D),E
;*             LD   (IX+&0E),D
;*             RET

Fix_L6EC2_43:  ;4.3 Chg
               JP   NZ,SETRPT  ;JP IF RECORD IN THIS SECTOR

;FixEnd

PTRC5:         CALL GTNSC
               JR   RSADM

;READ NEXT T/S FROM IX+NTRK

GTNSC:         PUSH IX
               POP  HL
               LD   DE,NTRK
               ADD  HL,DE
               LD   D,(HL)
               INC  HL
               LD   E,(HL)     ;DE=NEXT T/S
               LD   A,D
               OR   E
               JP   Z,REP27    ;EOF IF LAST SECTOR

               RET

;CALLED WHEN PTR HAS REACHED BUFFER END ON READ OR WRITE (NOT
;AT FILE END) WRITES CURRENT SECTOR IF IT HAS BEEN WRITTEN TO,
;BEFORE READING NEXT ONE

RDWRSR:        LD   D,(HL)     ;TRACK
               INC  HL
               LD   E,(HL)     ;SECT FROM END OF CURRENT BUFFER
               PUSH DE
               EX   DE,HL
               CALL SWPNSR     ;SET NEW CUR. T/S FROM HL,
                               ; GET DE=OLD CUR. T/S
               CALL WRIF2      ;IF SECTOR HAS BEEN WRITTEN TO,
               POP  DE         ; SAVE IT

RSADM:         CALL ICNT

RSADM2:        CALL SSDRV      ;SELECT DRIVE
               PUSH BC         ;PREV DRIVE
               CALL RSAD
               POP  BC
               LD   A,C
               CALL SSDRV2     ;PREV

MSECT:         RES  3,(IX+MFLG) ;"NOT WRITTEN TO YET"
                                ;MARK NSRL/H WITH DE

;SAVE THE NEXT TRACK/SECTOR

SVNSR:         LD   (IX+NSRH),D
               LD   (IX+NSRL),E
               RET

;INC CNT (CURRENT SECTOR NUMBER)

ICNT:          INC  (IX+CNTL)
               RET  NZ

               INC  (IX+CNTH)
               RET

;ENTRY: BHL=24-BIT NUMBER
;EXIT: HL=NUMBER MOD 510, HL'=NUMBER DIV 510

D510:          XOR  A
               LD   D,A        ;D MUST BE ZERO LATER
               EXX
               LD   L,A
               LD   H,A        ;SECTOR COUNT INITED
               EXX
               LD   A,B
               LD   BC,510     ;DBC=510

D510L:         EXX
               INC  HL
               EXX
               SBC  HL,BC
               SBC  A,D
               JR   NC,D510L

               ADD  HL,BC      ;HL=DISP IN SECTOR
               RET               ;HL=SECTOR

;GET FILE LEN IN AHL. NC FROM M510

GLEN:          CALL PTLEN      ;HL PTS TO LEN LOW (LIKE RPT)
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               INC  HL
               LD   C,(HL)
               INC  HL
               LD   B,(HL)     ;BC=LEN DIV 510
               EX   DE,HL      ;HL=LEN MOD 510

;LET AHL=HL+BC*510

M510:          PUSH BC         ;NEEDED BY 'GLEN'
               XOR  A
               LD   DE,510
               JR   PTAD2

PTADL:         EX   AF,AF'
               ADD  HL,DE
               ADC  A,0

PTAD2:         EX   AF,AF'
               DEC  BC
               LD   A,B
               OR   C
               JR   NZ,PTADL   ;LOOP UNTIL A'HL=PTR

               EX   AF,AF'     ;NC
               POP  BC
               RET

PTLEN:         PUSH IX
               POP  HL
               LD   BC,LENL
               ADD  HL,BC      ;HL PTS TO LEN LOW (LIKE RPT)
               RET


;--------------
Part_TIME:


DATE:          LD   HL,DATDT
               DEFB &FD        ;"JR+3"

TIME:          LD   HL,TIMDT

TIMDC:         LD   (TDVAR),HL
               CALL RDCLK      ;UPDATE BUFFERS
               CALL GTNC
               CALL CIEL
               JR   Z,TIMPR    ;JR IF PRINTING OF DATA,
                               ; NOT SETTING DATA
               CALL EVSTR      ;DE=START, BC=LEN
               CALL CEOS

               LD   A,(SVC)
               CALL SELURPG    ;PAGE STRING IN
               LD   HL,(TDVAR)
               LD   B,6

TIML:          LD   A,C
               AND  A
               LD   A,"0"
               JR   Z,TIM1     ;PAD WITH ZEROS ONCE STRING LEN=0

               LD   A,(DE)
               INC  DE
               DEC  C
               CALL NUMERIC
               JR   NC,TIML    ;LOOP IF NOT A DIGIT

TIM1:          LD   (HL),A

TIM2:          INC  HL
               LD   A,(HL)
               CALL NUMERIC
               JR   NC,TIM2    ;LOOP IF NOT A DIGIT

               DJNZ TIML       ;LOOP PLACING 6 NUMBERS FROM
                               ; STRING, OR PAD
               CALL DTVCK      ;CHECK VALUES
               JP   C,IOOR

;WRITE DATE AND TIME TO CLOCK HARDWARE FROM BUFFERS

WRCLK:         LD   IY,SCTWO
               JR   CKHW

TIMPR:         CALL CEOS

               LD   A,2
               CALL CMR
               DEFW STREAM

               LD   B,9
               LD   HL,(TDVAR)

;PRINT B FROM HL - USED TO PRINT DATE/TIME/PATH NAME

PBFHL:         LD   A,(HL)
               INC  HL
               CALL PNT
               DJNZ PBFHL

               RET

;READ CLOCK HARDWARE INFO INTO DATE AND TIME BUFFERS

RDCLK:         LD   IY,PLTWO

CKHW:          LD   A,(CKPT)   ;CLOCK PORT
               AND  A
               RET  Z          ;ABORT IF NO CLOCK!

               DI              ;CRASHES IF NOT...
               PUSH HL
               LD   C,A
               LD   B,&D0      ;CONTROL REGISTER
               LD   HL,2000

CKLP:          LD   A,1
               OUT  (C),A      ;HOLD ON
               IN   A,(C)
               BIT  1,A
               JR   Z,CKLX     ;EXIT IF NOT BUSY

               XOR  A
               OUT  (C),A      ;HOLD OFF
               DEC  HL
               CP   H
               JR   NZ,CKLP

               POP  HL         ;GIVE UP - NO CLOCK??
               RET

CKLX:          LD   B,&50      ;HOURS-H REGISTER
               LD   HL,TIMDT
               CALL IYJUMP     ;HANDLE HOURS
               CALL IYJUMP     ;HANDLE MINS
               CALL IYJUMP     ;HANDLE SECS
               LD   HL,DATDT
               LD   B,&70      ;DAY-H
               CALL IYJUMP     ;DAY
               LD   B,&90      ;MTH-H
               CALL IYJUMP     ;MONTH
               LD   B,&B0      ;YR-H
               CALL IYJUMP     ;YEAR
               XOR  A
               LD   B,&D0
               OUT  (C),A      ;HOLD OFF
               POP  HL
               RET

;ENTRY: HL PTS TO DEST, B=CLK REGISTER, C=PORT
;EXIT: HL INCED BY 3, REGISTER DECED BY 2, 2  DIGITS PLACED

PLTWO:         CALL PLONE
               DEC  HL

PLONE:         IN   A,(C)      ;READ DIGIT
               AND  &0F
               ADD  A,&30
               CP   &3A
               JR   NC,SCPLC   ;ENSURE HARDWARE ERRORS DO NOT
                               ; PLACE A NON-DIGIT
               LD   (HL),A
               JR   SCPLC

;ENTRY: HL PTS TO SRC, B=CLK REGISTER, C=PORT
;EXIT: HL INCED BY 3, REGISTER DECED BY 2, 2  DIGITS SENT TO CLK

SCTWO:         CALL SCONE
               DEC  HL

SCONE:         LD   A,(HL)     ;READ DIGIT
               SUB  &30
               OUT  (C),A

SCPLC:         INC  HL
               INC  HL
               LD   A,B
               SUB  &10
               LD   B,A        ;NEXT PORT
               RET

;DATE/TIME SR TO CHECK VALUES
;NC IF OK, CY IF >LIMIT

DTVCK:         LD   DE,(TDVAR)
               LD   HL,9
               ADD  HL,DE  ;DE PTS TO DATA TO CHECK, HL TO LIMIT
               LD   B,3    ;VALUES TO CHECK

VCLP:          CALL GDTV       ;GET VALUE IN A
               LD   C,A
               LD   A,(HL)
               INC  HL
               CP   C
               RET  C          ;RET IF VALUE>HI LIMIT

               LD   A,C
               CP   (HL)
               RET  C          ;RET IF VALUE<LOW LIMIT

               INC  HL
               DJNZ VCLP

               RET

;CALLED BY CFSM TO DATE STAMP FILE DIR ENTRY

DATSET:        CALL NRRD
               DEFW CURCMD

               CP   &CF        ;COPY
               RET  Z

               PUSH DE
               CALL RDCLK
               CALL POINT
               LD   BC,&00F5
               ADD  HL,BC      ;POSN F5 IN DIR ENTRY IS
                               ; DDMMYYHHSS (BCD)
               LD   DE,DATDT   ;SRC
               LD   B,3        ;3 PAIRS OF DIGITS

SDTL1:         CALL GDTV       ;GET VALUE IN A
               AND  A
               JR   Z,SDTL3    ;ABORT IF DY, MTH OR YEAR=0 (2000
                               ;(2000? PITY) INITIAL VALUE IS
                               ;00/00/00 SO NO STAMPING
               LD   (HL),A
               INC  HL
               DJNZ SDTL1

               LD   DE,TIMDT
               LD   B,2        ;(HRS/MINS)

SDTL2:         CALL GDTV       ;GET VALUE IN A
               LD   (HL),A
               INC  HL
               DJNZ SDTL2

SDTL3:         POP  DE
               RET

;GET DATE/TIME VALUE IN A FROM (DE) AND (DE+1)
;ADVANCES DE BY 3, USES C

GDTV:          LD   A,(DE)     ;MS DIGIT
               INC  DE
               SUB  &30
               LD   C,A
               ADD  A,A
               ADD  A,A
               ADD  A,C
               ADD  A,A        ;*10

               LD   C,A
               LD   A,(DE)     ;LS DIGIT
               INC  DE
               INC  DE         ;SKIP SEPARATOR
               SUB  &30
               ADD  A,C
               RET

;PRINT DATE/TIME AS DD/MM/YY HH:MM

PNDAT:         LD   B,&F5
               CALL GRPNTB     ;PT HL TO DATE
               LD   A,(HL)
               INC  A        ;DOS 2.0 FILES HAVE FF ->00
                             ;G+DOS FILES HAVE 00 ->01
                             ;DOS 3.0 FILES HAVE 01 OR MORE ->02
               CP   2
               RET  C          ;RET IF NO DATE

               LD   C,"/"
               CALL PSIX       ;DATE, WITH "/" SEPARATOR
               LD   C,":"
               LD   A," "
               JR   PFOUR      ;TIME, WITH ":" SEPARATOR

;PRINT SIX DIGITS FROM 3 BYTES AT (HL) AS:
; NN (SEPARATOR) NN (SEPARATOR) NN
;ENTRY: C=SEPARATOR

PSIX:          CALL PTWO
               LD   A,C

PFOUR:         CALL PNT
               CALL PTWO
               LD   A,C
               CALL PNT

PTWO:          PUSH HL
               PUSH DE
               PUSH BC
               LD   L,(HL)
               LD   H,0
               LD   A,"0"
               CALL PNUM2
               POP  BC
               POP  DE
               POP  HL
               INC  HL
               RET


;--------------
Part_SUBD:


;OPEN DIRECTORY - JUMPED TO FROM "OPEN" WHEN NEXT CHAR<>"#"

OPNDIR:        LD   C,144      ;DIR
               CALL ISEP
               CALL EVSYN      ;NAME
               CALL PLNS       ;PLACE NEXT STAT ADDR
               CALL CEOS
               LD   A,DFT      ;"DIRECTORY"
               LD   (NSTR1),A
               CALL SETF4      ;"NO SECTOR NEEDED"
               CALL GOFSM
               LD   A,(MAXT)   ;SET BY FDHR
               INC  A          ;TAG VALUE FOR THIS DIRECTORY'S
               CP   &FF        ; FILES
               JP   Z,REP25    ;ERROR IF TOO MANY DIR FILES

               LD   (FSA+DIRT),A ;COPIED TO SECTOR BUFFER
               JP   CFSM

;SET CURRENT DIRECTORY
; - JUMPED TO FROM "DIR" CMD WHEN "=" FOLLOWS IT

STDIR:         CALL GTNC       ;SKIP "="
               CP   &5E        ;UPARROW
               JR   Z,STDUQ

               CALL RTCK
               JR   NZ,STDQ

STDUQ:         EX   DE,HL      ;START OF "STRING" TO DE
               CALL GTNC
               CALL CEOS

               LD   BC,1       ;LEN
               IN   A,(251)
               DEFB &FE        ;"JR+1"

;HOOK SET DIR. A=PAGE, BC=LEN, DE=START, OF NAME

HSDIR:         EXX
               LD   (SVC),A
               CALL EVNM2      ;EVAL AS NAME, AVOIDING STRING
               JR   STDQC      ; FETCH

STDQ:          CALL EVNAM
               CALL CEOS

STDQC:         LD   HL,NSTR1+1
               CALL EVFILE
               CALL CKDISC
               CALL GPLA       ;GET PATH LEN ADDR IN HL
               LD   A,(HL)
               LD   (TEMPW1),A ;SAVE PATH LEN TO ALLOW UP ARROW
                               ; TRIM/NO TRIM
               LD   HL,NSTR1+1
               LD   A,(HL)
               CALL RTCK       ;CHECK IF ROOT
               JR   NZ,STDR0

               CALL C11SP
               JR   NZ,STDR0   ;JR IF NOT ROOT BY ITSELF
                               ;ELSE SET ROOT (NO DISC ACCESS)

;SET ROOT=PERM DIR. CALLED IF RAND WORD CHANGES, USED IF DIR=/

SETRT:         LD   BC,2       ;PATH LEN FOR "1:" OR "2:"
               JP   STDR2      ;SELECT ROOT (B=0)

STDR0:         CP   94         ;UP-ARROW
               JR   NZ,STDR1

               CALL C11SP
               JR   NZ,STDR1   ;JR IF NOT UP-ARROW BY ITSELF

               CALL UPAR2  ;CHANGE CDIRT TO LEVEL OF PARENT DIR
                           ;IF THERE IS A PARENT DIR.
                           ;SAVE TRIMED PATH LEN IN TEMPW1
               CALL GPLA
               LD   A,(TEMPW1)
               LD   (HL),A     ;TRIM LENGTH
               CALL GCDIA      ;GET CDIRP ADDR
               LD   A,(CDIRT)
               LD   (HL),A
               RET

STDR1:         IN   A,(251)
               PUSH AF
               CALL PRPTH    ;GET B=LEN, HL=START, OF PAGED-IN
                             ; PATH STRING
               PUSH BC       ;LEN IN B
               PUSH HL
               CALL DTREE ;DESCEND PATH TILL LAST FILE IN STRING
               CALL FNDIR    ;FIND DIRECTORY

               CALL GPLA
               POP  DE
               POP  BC       ;B=LEN

SDFCL:         LD   A,(DE)
               CALL RTCK
               LD   A,2  ;PATH LEN IS 2 IF DIR NAME STARTS WITH
                         ;ROOT. LATER NAMES ADD TO THIS STUMP,
                         ;E.G. DIR="\GAMES\SCAB" SETS PATH$ TO
                         ;"1:\GAMES\SCRAB" WHATEVER CURRENT DIR
                         ;IS. OTHERWISE E.G. IF PATH$="1:\GAMES"
                         ;AND DIR="SCRAB" THEN PATH$ BECOMES
                         ;"1:\GAMES\SCRAB"

               JR   Z,SKIPF    ;JR IF ROOT

               LD   A,(DE)
               CP   94         ;UP-ARROW
               JR   NZ,STDNR

               LD   A,(TEMPW1) ;LEN LEFT BY UPAR SR

SKIPF:         LD   (HL),A   ;SET LENGTH
               INC  DE       ;SKIP LEADING ROOT OR ARROW SYMBOL.
                             ;ROOT SYM IS ADDED AUTOMATICALLY AS
                             ;NAME IS ADDED TO PATH
               DJNZ SDFCL      ;DEC LEN

               CALL DERR
               DEFB 18

STDNR:         PUSH BC
               CALL GPATD      ;GET HL=PATH START, BC=LEN
               ADD  HL,BC
               EX   DE,HL      ;DE=END OF PATH DATA, HL=NAME
               POP  AF
               LD   B,A        ;B=NAME CHARS TO COPY
               INC  B          ;PLUS 1 FOR DIVIDER
               LD   A,(RTSYM)  ;DIVIDER SYMBOL
               JR   SDN2

SDNL:          LD   A,(HL)
               INC  HL

SDN2:          EX   AF,AF'
               INC  C
               LD   A,C
               CP   MPL+1      ;MAX LEN+1
               JR   NC,SDTL    ;JR IF NO ROOM LEFT (WOULD BE
                               ; 33RD CHAR NEXT)
               EX   AF,AF'
               CALL RTCK
               JR   NZ,SDN3

               LD   A,(RTSYM)  ;ALWAYS USE FIRST ALTERNATIVE
                               ; DIVIDER IN PATH$
SDN3:          LD   (DE),A
               INC  DE
               DJNZ SDNL       ;COPY NAME

               INC  C

SDTL:          DEC  DE
               DEC  C
               LD   A,(DE)
               CP   " "
               JR   Z,SDTL     ;TRIM SPACES
                               ;C=NEW PATH LEN
               PUSH BC
               CALL POINT      ;START OF DIR ENTRY
               POP  BC
               LD   DE,DIRT    ;DISP TO TAG VALUE FOR FILES IN
               ADD  HL,DE      ; THIS DIR
               LD   B,(HL)
               POP  AF
               OUT  (251),A

STDR2:         CALL GCDIA      ;GET CDIRP ADDR
               LD   (HL),B     ;SET CDIRP
               LD   A,B
               LD   (CDIRT),A  ;AND CDIRT
               CALL GPLA       ;GET PATH LEN ADDR
               LD   (HL),C

OSRDPN:        CALL TIRD
               JP   NC,SRDPN   ;STORE RD PATH NAME IF RAM DISC

               RET

;ROOT CHECK - Z IF ROOT SYMBOL

RTCK:          PUSH HL
               LD   HL,(RTSYM) ;USUALLY "\" AND "/"
               CP   H
               JR   Z,RTCK2

               CP   L

RTCK2:         POP  HL
               RET

;FIND DIRECTORY. EXIT WITH HL=POINT

FNDIR:         CALL SNDFL
               JR   NC,REP32H  ;ERROR IF DIRECTORY NOT FOUND

               LD   HL,FLAG3
               RES  2,(HL)   ;"NOT INITIALISED" SO LATER USE OF
                             ;FNDFL AND SNDFL STARTS FROM T0/S1
               CALL POIDFT
               RET  Z          ;RET IF DIR TYPE

REP32H:        JP   REP32      ;ERROR IF NOT RIGHT TYPE

EVFINS:        LD   HL,NSTR1+1
               CALL EVFILE
               CALL CKDISC

;GET LAST NAME

GTLNM:         IN   A,(251)
               PUSH AF
               CALL PRPTH      ;GET B=LEN, HL=START, OF PAGED-IN
                               ; PATH STRING
               CALL DTREE      ;DESCEND PATH TILL LAST FILE
               POP  AF         ; IN NSTR1+1
               OUT  (251),A
               RET

;PREPARE PATH

PRPTH:         XOR  A
               OUT  (251),A      ;PAGE IN STRING BUFFER

               LD   A,(&4F60+FS) ;REAL STRING LENGTH
               LD   HL,(TEMPW3)  ;CHARS TRIMMED OFF FRONT
               SUB  L            ; E.G. "D2:"
               JR   C,PRP2     ;JR IF E.G. "" OR "D2" ALTERED
                               ; TO "T:" OR "D2:"
                               ; SO TRIMMED CHARS > ORIG LEN
               JR   Z,PRP2     ;JR IF ALL CHARS TRIMMED
                               ; (E.G. "D2:")
               LD   B,A        ;CORRECTED LEN
               LD   DE,&8F10   ;SYS PAGE BUFFER
               ADD  HL,DE      ;CORRECTED START
               RET               ;LEN IN B, START IN HL

PRP2:          LD   HL,PRPD    ;USE "  " IF TAPE
               LD   B,1
               LD   A,(LSTR1)
               CP   "T"
               RET  Z

               INC  HL         ;ELSE "*"
               RET

PRPD:          DEFB " "
               DEFB "*"

;DESCEND PATH TREE
;ENTRY: HL PTS TO PATH STRING, B=LEN (<>)
;ACTION: PATH IS BROKEN INTO DIRECTORY NAMES AND EACH DIRECTORY
;IS SELECTED UNTIL A FINAL NAME IS RETURNED IN NSTR1+1

DTREE:         CALL UPARW      ;CHECK IF FIRST CHAR IN NAME IS
                               ;UP-ARROW, HANDLE IF SO
               CALL RTCK
               JR   NZ,DTREL

               INC  HL         ;SKIP "\"
               DEC  B
               JR   Z,IFNE     ;E.G. LOAD "\" IS AN ERROR

               XOR  A
               LD   (CDIRT),A  ;ROOT IS TEMP DIRECTORY

DTREL:         CALL EFLNM      ;GET FILE NAME
               RET  C          ;RET IF LAST NAME

               PUSH HL
               PUSH BC
               CALL FNDIR      ;FIND DIRECTORY FILE
               LD   BC,DIRT    ;DISP TO TAG VALUE FOR FILES IN
               ADD  HL,BC      ; THIS DIR
               LD   A,(HL)
               LD   (CDIRT),A
               POP  BC
               POP  HL
               JR   DTREL

;EXTRACT FILE NAME FROM STRING
;
;ENTRY: HL POINTS TO POSN IN SRC STRING, B=SRC LEN LEFT
;EXIT: C=FILE NAME LEN,        NSTR1+1 HOLDS NAME
;      IF CY, NAME IS THE LAST ONE
;      HL POINTS TO NEXT NAME START (ROOT SYMBOL IS SKIPPED),
;      B IS UPDATED

;E.G. "UTILS\PDS\TEST" BROKEN DOWN TO: UTILS, PDS, TEST

EFLNM:         LD   DE,NSTR1+1
               PUSH DE
               LD   A," "
               LD   C,10

EFCL:          LD   (DE),A
               INC  DE
               DEC  C
               JR   NZ,EFCL  ;INITIAL CLEAR OF NAME. C BECOMES 0

               POP  DE

EFNLP:         LD   A,(HL)
               INC  HL
               DEC  B
               CALL RTCK
               RET  Z          ;RET IF ROOT. NC

               INC  B
               INC  C          ;NAME LEN
               EX   AF,AF'
               LD   A,C
               CP   11
               JR   Z,IFNE

               EX   AF,AF'
               LD   (DE),A
               INC  DE
               DJNZ EFNLP

               SCF               ;"LAST NAME"
               RET

IFNE:          CALL DERR
               DEFB 18         ;"INVALID FILE NAME"

;UP-ARROW MEANS "GO UP ONE LEVEL" - SEARCH PATH$ FOR LAST ROOT
;ENTRY: HL POINTS TO STRING, B=LEN                       ;SYMBOL
;EXIT: HL POINTS PAST UP-ARROW, B IS LESS

UPARW:         LD   A,(HL)
               CP   94
               RET  NZ

               INC  HL         ;SKIP LEADING UP-ARROW
               DEC  B
               JR   Z,IFNE     ;E.G. LOAD "UP-ARROW" IS AN ERROR

UPAR2:         PUSH BC
               PUSH HL
               CALL GPATD      ;GET HL=PATH START, BC=LEN
               ADD  HL,BC
               LD   A,(RTSYM)

FLRSL:         DEC  HL
               DEC  C
               JR   Z,STDUF    ;JR IF NO ROOT FOUND - NO ACTION

               CP   (HL)
               JR   NZ,FLRSL

               LD   A,C
               LD   (TEMPW1),A ;NEW LEN STORED FOR USE BY "DIR="
               PUSH HL
               CALL GPLA       ;PT HL TO LEN
               LD   A,(HL)     ;OLD
               POP  HL
               INC  HL         ;PT TO START OF TRIMMED DIR NAME
               LD   DE,NSTR1+1
               SUB  C
               DEC  A          ;LEN OF TRIMMED NAME
               CP   11
               JR   C,STDTN    ;SHOULD BE JR ALWAYS....

               LD   A,10

STDTN:         LD   C,A
               LD   B,0
               LDIR
               LD   C,A
               LD   A,10
               SUB  C
               JR   Z,STDPP    ;JR IF NO PADS NEEDED

               EX   DE,HL

SPDL:          LD   (HL),&20
               INC  HL
               DEC  A
               JR   NZ,SPDL

STDPP:         LD   HL,CDIRT
               LD   A,(HL)
               LD   (HL),&FF   ;"ANY"

STFPL:         PUSH AF
               CALL FNDIR      ;FIND PARENT DIRECTORY
               CALL SETF2      ;"RESTART FROM CURRENT T/S"
               LD   BC,DIRT
               ADD  HL,BC
               POP  AF         ;ORIG CDIRT
               CP   (HL)
               JR   NZ,STFPL   ;JR IF NAME CORRECT BUT NOT REAL
                               ;PARENT ACCORDING TO TAG
               LD   BC,254-DIRT
               ADD  HL,BC
               LD   A,(HL)
               LD   (CDIRT),A  ;MAKE CURRENT LEVEL SAME
               LD   HL,FLAG3
               RES  2,(HL)     ;START FROM T0/S1
               CALL OSRDPN           ;STORE RD PATH NAME

STDUF:         POP  HL
               POP  BC
               LD   A,(HL)
               RET

;GET BC=PATH LEN, HL=PATH, FOR DRIVE

GPATD:         CALL GPLA
               LD   B,0
               LD   C,(HL)     ;BC=PATH LEN
               CALL TIRD       ;A=DRIVE, NC IF RAM DISC
               JR   NC,GPATD2

               LD   HL,PTH1
               DEC  A
               RET  Z

               LD   HL,PTH2
               RET

GPATD2:        PUSH BC
               PUSH DE
               CALL MRDPN  ;CALLED WITH NC - FETCH NAME TO PTHRD
               LD   HL,PTHRD
               POP  DE
               POP  BC
               RET

;RETURN TRACKS/DISC FOR CURRENT RAM DISC IN A

RTSTD:         LD   HL,RDDT-3  ;DRIVE 3 IS FIRST ENTRY
               CALL GPLA2      ;GET ADDR
               LD   A,(HL)     ;T/DISK
               RET

;GET FIRST PAGE ADDR FOR RAM DISC

GFPA:          LD   HL,FIPT-3
               JR   GPLA2

;GET RND WORD ADDR

GRWA:          LD   HL,CRWT-2
               LD   A,(DRIVE)
               ADD  A,A
               JR   GPLA3

;GET CDIRP ADDR

GCDIA:         LD   HL,CDIT-1
               JR   GPLA2

;GET PATH LEN ADDR

GPLA:          LD   HL,PLT-1

GPLA2:         LD   A,(DRIVE)

GPLA3:         ADD  A,L
               LD   L,A
               RET  NC

               INC  H
               RET

;SET DIR TKS FROM SECTOR JUST READ. COMPARE RAND WORD ON DISC
;WITH CURRENT, ALTER DIR TO ROOT IF A NEW DISC, UPDATE CURRENT
;RAND WORD IN MEMORY

SDTKS:         PUSH DE
               CALL POINT
               PUSH HL
               INC  H
               DEC  HL
               LD   A,(HL)
               ADD  A,4
               LD   (DTKS),A
               DEC  HL
               DEC  HL
               LD   C,(HL)
               DEC  HL
               LD   B,(HL)     ;BC=RND WORD FROM DISC
               CALL GRWA       ;GET CUR. RAND WORD ADDR IN HL
               LD   A,C
               CP   (HL)
               JR   NZ,SDTK2   ;JR IF NEW DISC

               INC  HL
               LD   A,B
               CP   (HL)       ;NZ IF NEW DISC
               DEC  HL

SDTK2:         LD   (HL),C
               INC  HL
               LD   (HL),B
               POP  HL         ;POINT VALUE
               PUSH HL
               PUSH BC
               PUSH HL
               CALL NZ,SETRT   ;SET ROOT IF NEW DISC

               POP  HL         ;POINT VALUE
               LD   BC,&00D2
               ADD  HL,BC      ;DISC NAME ON DISC
               LD   DE,DNAME
               LD   C,10
               LDIR            ;COPY NAME TO MSG BUFFER
               POP  BC         ;CURRENT RND NO.

               LD   A,(SAMDR)
               LD   HL,DRIVE
               CP   (HL)
               JR   NZ,SDTK4
               LD   A,(SAMCNT)
               AND  A
               JR   Z,SDTK4  ;JR IF NO OPEN FILES TO WARN ABOUT

               LD   HL,(SAMRN) ;SAM RND NO.
               SBC  HL,BC
               JR   Z,SDTK4    ;JR IF SAME DISC AS WHEN

               CALL CLSL
               CALL PMOOF      ;"OPEN file"
               CALL BEEP
               LD   A,(SSTR1)
               CALL CMR
               DEFW STREAM

SDTK4:         POP  HL
               POP  DE
               RET

;--------------
Part_RAMD:


;RAM DISC WRITE SECTOR AT DE. D=TRACK, E=SECTOR.

RDWSCT:        CALL GTBUF    ;GET SRC ADDR IN HL. MAY BE 8000H+
                             ;IF WRITE AT OR SAVE BLOCK
               LD   A,DRAM/256
               CP   H
               JR   Z,RDW2   ;JR IF SRC=DRAM (FIRST SECTORS,
                             ; SERIAL FILES)
               PUSH DE       ;T/S
               LD   DE,DRAM
               LD   BC,&0200
               PUSH DE
               LDIR
               POP  HL     ;SRC=DRAM SO PAGING OF DEST POSSIBLE
                           ;(used by WRITE AT)
               POP  DE       ;T/S

RDW2:          PUSH DE       ;T/S
               PUSH HL       ;SRC
               CALL RDADR    ;GET RD ADDR IN HL
               EX   DE,HL    ;DEST IN RD INTO DE
               POP  HL       ;SRC IN RAM
               LD   BC,&0200
               LDIR

RDW3:          OUT  (251),A

RDW4:          POP  DE       ;T/S
               CALL CLRRPT
               JP   GTBUF

;RAM DISC SEARCH SECTOR AT DE
;ENTRY: DE=T/S, (TEMPW1)=DELIM COUNT, (DELIM)=DELIM
;EXIT: TEMPW3 HOLDS NEW DELIM COUNT, TEMPW2=LOCN
;COULD BE FASTER IF SECTOR ONLY READ TO BUFFER WHEN MATCHED AND
;COUNT=0: ELSE JUST READ LAST 2 BYTES??

RDSSAD:        CALL RDRSCT     ;HL=DATA START
               PUSH DE
               PUSH HL
               LD   DE,(TEMPW1)
               LD   BC,&0002

RDS1:          LD   A,(DELIM)

RDSLP:         CP   (HL)
               INC  HL
               JR   Z,RDS2

RDS12:         DJNZ RDSLP

               DEC  C
               JR   NZ,RDSLP

               JR   RDS3

RDS2:          DEC  DE
               LD   A,D
               OR   E

Fix_L7394_42:  ;4.2 Cut
;*             JR   NZ,RDS1

Fix_L7396_43:  ;4.3 Add
               LD   A,(DELIM)
               JR   NZ,RDS12

;FixEnd

               LD   (TEMPW2),HL  ;LOCN OF MATCH AND COUNT=0

RDS3:          LD   (TEMPW3),DE  ;NEW COUNT
               POP  HL
               POP  DE
               RET

;RAM DISC PREPARE SAM FROM DIR ENTRY READ

NRDRSCT:       CALL RDRSCT
               PUSH DE
               PUSH HL
               EX   DE,HL
               LD   HL,NSAM
               LD   B,2

NRDROL:        LD   A,(DE)
               INC  D
               AND  A
               JR   Z,NRDR2    ;JR IF AN ERASED ENTRY

               DEC  D

NRDRL:         LD   A,(DE)
               OR   (HL)
               LD   (HL),A
               INC  DE
               INC  L
               JR   NZ,NRDRL

NRDR2:         DJNZ NRDROL

               POP  HL
               POP  DE
               RET

;RAM DISC READ SECTOR AT DE. D=TRACK, E=SECTOR.

RDRSCT:        CALL GTBUF
               LD   BC,&0200

RDRS2:         PUSH DE         ;T/S
               PUSH HL         ;MAIN MEM PTR
               PUSH BC
               CALL RDADR      ;GET SRC ADDR IN HL
               LD   BC,&0200
               LD   DE,DRAM
               LDIR            ;COPY RAMDISC TO DRAM

               OUT  (251),A    ;ORIG URPORT
               POP  BC
               POP  HL
               LD   A,DRAM/256
               CP   H
               JR   Z,RDW4     ;RET IF DEST=DRAM - DONE IT

               EX   DE,HL      ;DE=BUF
               LD   HL,DRAM
               LDIR

               JR   RDW4

;RAM DISC ADDRESS - GET HL=ADDR, PAGE SELECTED, FROM D/E (T/S)
;GET A=CUR. URPORT. BC, DE USED

RDADR:         LD   A,(RDAT)
               AND  A
               JR   Z,RDAD2

               CALL SELFP
               LD   HL,(&82FF)
               OUT  (251),A
               LD   A,L
               ADD  A,3
               CP   D          ;CP DTKS-1,TRK
               JR   NC,RDAD2 ;JR IF TRK IN DIRECTORY - NO FIDDLE
                             ; (ALWAYS JR ON TRACK 0 - WHEN DTK
                             ; IS OBSOLETE!)
               SUB  3
               JR   NC,RDAD2 ;LEAVE TRACK ALONE IF 4 DIR TKS OR
                             ; MORE. A=-1 IF DTKS=3, -2 IF 2,
                             ; -3 IF 1 DTK
               ADD  A,D
               LD   D,A    ;E.G. TRACK 4 ACCESS BECOMES TRACK 1
                           ;SO IF DTKS=1, DISK=40 TRK, TRK 0=0
                           ;TRKS 4-42 BECOME 1-39
RDAD2:         CALL RTSTD    ;GET A=TRACKS/DRIVE
               AND  A
               JP   Z,REP22  ;"NO SUCH DRIVE" IF TKS=0

               LD   C,A
               LD   A,D
               CP   C
REP4H:         JP   NC,REP4  ;TRK/SCT ERROR IF TRACK >=(LIMIT)

               AND  &7F
               CP   80
               JR   NC,REP4H ;ERROR IF E.G. TRK 81

               BIT  7,D
               JR   Z,RDAD3  ;JR IF "SIDE 1"

               LD   A,D
               SUB  48       ;128->80
               LD   D,A

RDAD3:         CALL CPFTS    ;GET PAGEFORM OF DISC ADDR
               PUSH HL
               CALL PTRD2    ;PT TO ENTRY "A" IN LIST
               LD   L,(HL)   ;GET PAGE VALUE
               OUT  (251),A
               LD   A,L
               POP  HL
               JR   SELRDP

;SEL FIRST PAGE

SELFP:         CALL GFPA
               LD   A,(HL)   ;A=FIRST PAGE

;SELECT RAM DISC PAGE. ENTRY: A=PAGE. 00-1F=INTERNAL RAM, 20-FFL
; 20-FFH=EXTERNAL. EXIT: A=ORIG URPORT VALUE

SELRDP:        DI
               PUSH DE         ;SAVE THROUGHOUT
               LD   D,A        ;40-7F=MEGA RAM 1,
                               ;80-BF=MEGA RAM 2,
                               ;C0-FF=MEGA RAM 3,
                               ;20-3F=MEGA RAM 0 SECOND HALF
               CP   &20        ;??? SUB &20 ???
               JR   C,SRDP2    ;JR IF MAIN RAM

               OUT  (MRPRT),A  ;SELECT MEGA RAM
               LD   D,&80      ;XMEM BIT HIGH

SRDP2:         IN   A,(251)
               PUSH AF
               LD   A,D
               OUT  (251),A
               POP  AF         ;ORIG PAGE
               POP  DE         ;ORIG
               RET

;CALCULATE PAGE FORM FROM T/S IN DE

CPFTS:         LD   L,D
               LD   H,0
               LD   B,H
               LD   C,L
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,BC
               ADD  HL,HL      ;TRACK*10
               LD   C,E
               DEC  C
               LD   A,C
               CP   10
               JP   NC,REP27   ;"END OF FILE" IF SECTOR NOT 1-10

               ADD  HL,BC      ;SECTOR NUMBER 0 TO DISC LIMIT
               PUSH HL
               LD   C,31
               XOR  A          ;NC
               DEC  A          ;A=0FFH

DIV31L:        INC  A
               SBC  HL,BC
               JP   NC,DIV31L

               POP  HL         ;SECT NO 0 TO LIMIT-1
               LD   C,A
               ADD  HL,BC      ;SECT=SECT+INT(SECT/31) TO AVOID
                               ; EVERY 32ND SECT
               INC  HL         ;AVOID SECT 0,32,64,96 ETC
               ADD  HL,HL      ;HL=SECT NO.*2
               LD   A,H
               LD   H,L
               LD   L,B      ;AHL=20-BIT DISPLACEMENT (512*SECT)
                             ;NOW GET A=PG DISP, HL=OFFSET

;TRANSFORM 20-BIT NUMBER IN AHL TO PAGE, ADDR (8000-BFFF)

PAGEFORM:      RL   H
               RLA
               RL   H
               RLA             ;NC. PAGE NOW OK
               RR   H
               SCF
               RR   H        ;ADDR NOW OK IN 8000-BFFF FORM
               RET

;RAM DISC SAVE BLOCK

RDSB:          CALL WSAD       ;SAVE FIRST BLOCK
               LD   BC,0       ;SECT COUNT=0
               JR   RDSB2

RDSBL: ;     * PUSH DE         ;PREV TRK
               CALL CCNT       ;CHECK BYTES REMAINING, SUB 511
               EX   DE,HL
               LD   HL,(SVHL)  ;SRC
               JR   C,RDSB3    ;JR IF 510 OR LESS BYTES TO GO
                               ; - PREPARE LAST SCT
               INC  DE        ;CORRECT FOR SUB 511 ->SUB 510 NOW
               LD   (SVDE),DE
               LD   DE,DRAM    ;TEMP STORE
               LD   BC,510
               LDIR            ;COPY DATA TO TEMP BUFFER IN
                               ; COMMON MEMORY
               CALL CHKHL
               LD   (SVHL),HL
        ;    * POP DE          ;PREV TRK
               CALL FFNS       ;FAST FIND NEXT FREE SECT IN SAM
               LD   L,D
               LD   H,E
               LD   (DRAM+510),HL

               EX   DE,HL
               CALL SWPNSR     ;GET T/S IN D/E
               CALL WSAD
               LD   BC,(SVCNT)
               INC  BC         ;INC SECTOR COUNT

RDSB2:         LD   (SVCNT),BC
               JR   RDSBL

RDSB3:  ;    * POP DE          ;PREV TRK
               LD   BC,(SVDE)
               LD   (IX+RPTH),B
               LD   (IX+RPTL),C
               LD   DE,DRAM
               PUSH BC
               LDIR            ;LAST PART-BUFFER COPIED

               EX   DE,HL      ;BUFFER DEST IN HL
               POP  BC         ;PTR
               JP   SVBS1

;FORMAT RAM DISC. A=DSTR1

FORMRD:        CP   RDLIM
               JP   NC,REP22   ;"NO SUCH DRIVE"

               DI
               LD   (DRIVE),A
               XOR  A
               OUT  (251),A    ;SYS PAGE AT 8000H
               CALL RTSTD      ;GET TKS/DRIVE
               AND  A
               JR   Z,FRMRD2   ;JR IF NEW RAMDISC

               CALL PTRD       ;PT TO TABLE START - 50 BYTES OF
                               ;PAGES USED+00H IN FIRST PAGE
               LD   DE,DRAM
               PUSH DE
               LD   BC,52
               LDIR            ;COPY TO BUFFER
               POP  HL
               OUT  (251),A    ;ORIG
               LD   D,ALLOCT/256+&40

FRMRDL:        LD   A,(HL)
               INC  HL
               LD   E,A
               AND  A          ;TERMINATOR IS 0
               JR   Z,FRMRD2 ;EXIT IF ALL USED PAGES DE-ALLOCTED

               AND  &E0
               JR   NZ,FRMRD1  ;JR IF MEGA RAM PAGE

               LD   (DE),A     ;FREE PAGE IN ALLOCT
               JR   FRMRDL

FRMRD1:        LD   A,E
               CALL RMRBIT     ;FREE BIT IN MEGA RAM MAP
               JR   FRMRDL

FRMRD2:        LD   HL,DRAM    ;PT TO TABLE START - 50 BYTES OF
               LD   B,52       ; PAGES USED
               XOR  A

RDCLL:         LD   (HL),A
               INC  HL
               DJNZ RDCLL      ;CLEAR TABLE

               LD   A,(TEMPW1) ;TKS/DISC
               AND  A
               JP   Z,FTKPD    ;JUST ERASE ANY EXISTING RAM DISC
                               ; IF TKS=0
               LD   D,A
               DEC  D
               LD   E,10       ;LAST SECTOR ON DISC
               CALL CPFTS      ;GET PAGE FORM OF T/S ADDR
               INC  A
               LD   C,A        ;PAGES NEEDED
               CALL CNTFP      ;COUNT FREE PAGES IN A AND B
               CP   C
               JP   C,OOMERR

               LD   A,C
               PUSH AF         ;PAGES NEEDED
               LD   B,0
               LD   HL,DRAM
               ADD  HL,BC      ;LAST REQUIRED POSN IN TABLE,+1
               POP  BC         ;B=PAGES NEEDED
               LD   DE,ALLOCT+FS+32

FRMRDEL:       CALL SFMRP  ;SEARCH FOR FREE MR PAGE & RESERVE IT
               JR   Z,FRMRD6   ;JR IF PAGE FREE (A=PAGE)

               DEC  E          ;INTERNAL RAM PAGE NO./ALLOCT PTR
               XOR  A
               OUT  (251),A    ;SYS PAGE IN AT 8000H
               LD   A,(DE)
               AND  A
               JR   NZ,FRMRDEL ;JR IF PAGE NOT FREE

FRMRD5:        LD   A,(DRIVE)
               OR   &D0
               LD   (DE),A     ;RESERVE IN ALLOCT
               LD   A,E

FRMRD6:        DEC  HL         ;NEXT TABLE POSN
               LD   (HL),A     ;PAGE NUMBER TO RAM DISC'S LIST
               CALL SELRDP
               EXX
               LD   HL,RDCODE
               LD   DE,&8002
               LD   BC,RDCE-RDCODE+1
               LDIR            ;CREATE MOVER CODE IN UNUSED
                               ; "SECTOR" OF EACH PAGE
               LD   HL,&8020
               LD   B,H        ;B=128

FTCCL:         LD   (HL),&ED
               INC  HL
               LD   (HL),&A0
               INC  HL
               DJNZ FTCCL      ;CREATE 128 LDIs

               LD   (HL),&3D   ;DEC A
               INC  HL
               LD   (HL),&C2   ;JP NZ
               INC  HL
               LD   (HL),&20
               INC  HL
               LD   (HL),&80   ;8020H
               INC  HL
               LD   (HL),&C9   ;RET
               LD   B,&3E      ;MAX NUMBER OF POSSIBLE DIR
                               ; ENTRIES IN PAGE
               LD   HL,&8200   ;FIRST POSSIBLE DIR ENTRY

FZDL:          LD   (HL),C
               INC  L
               LD   (HL),C     ;ENSURE IT STARTS 00/00
               DEC  L          ; SO "NEVER USED"
               INC  H          ;NEXT ENTRY
               DJNZ FZDL

               EXX
               DJNZ FRMRDEL

               LD   HL,&82FF
               CALL FESET      ;SET DTKS AND RND WORD AND NAME
               LD   HL,DRAM
               LD   DE,&8125   ;AFTER MULTI-LDI CODE
               LD   BC,52
               LD   A,(HL)
               LDIR            ;COPY PAGE TABLE FROM TEMP BUFFE
                               ; TO 1ST PAGE
               PUSH AF
               CALL GFPA       ;FIRST PAGE ADDR
               POP  AF
               LD   (HL),A     ;SET FIRST PAGE VAR

FTKPD:         CALL RTSTD      ;GET ADDR OF TKS/DISC
               LD   A,(TEMPW1) ;DESIRED TKS/DISC
               CP   80
               JR   C,NRDTF

               ADD  A,48       ;E.G. 80->128, 160->208

NRDTF:         LD   (HL),A     ;NZ SO DISC USABLE NOW, OR ZERO
               XOR  A          ; SO NON-EXISTENT
               OUT  (251),A    ;ENSURE XMEM BIT LOW
               RET

OOMERR:        CALL DERR
               DEFB 1

PTRD:          XOR  A

PTRD2:         LD   C,A
               CALL SELFP      ;SEL FIRST PAGE, EXIT WITH A=ORIG
               LD   HL,&8125   ;START OF FULL TABLE
               LD   B,0
               ADD  HL,BC      ;PT TO REQUIRED ENTRY
               RET

;STORE RAM DISC PATH NAME

SRDPN:         SCF

;MOVE RAM DISC PATH NAME
;ENTRY: (DRIVE) SET, NC IF FETCH FROM RD, CY IF STORE

MRDPN:         DI
               EX   AF,AF'
               CALL RTSTD      ;GET A=TRACKS/DRIVE
               AND  A
               JP   Z,REP22    ;"NO SUCH DRIVE" IF TKS=0

               CALL SELFP      ;SEL FIRST PAGE, EXIT WITH A=ORIG
               LD   HL,&8160
               LD   DE,PTHRD+2
               EX   AF,AF'
               JR   NC,MRDPN2

               EX   DE,HL

MRDPN2:        LD   BC,MPL-2
               LDIR            ;COPY PATH NAME FROM RAM DISC
                               ; TO/FROM TEMP BUFFER
               LD   A,(DRIVE)
               ADD  A,&30
               LD   L,A
               LD   H,":"
               LD   (PTHRD),HL
               EX   AF,AF'
               OUT  (251),A
               RET

;RAM DISC LOAD BLOCK SR

RDLB:          LD   HL,(SVHL)  ;DEST
               CALL SDCHK2
               JR   NC,RDLB2

               LD   BC,510
               CALL RDRS2
               LD   A,(DRAM+510)
               LD   B,A
               LD   A,(DRAM+511)
               LD   C,A
               JR   RDLB3

RDLB2:         CALL FRDRD2  ;LOAD 510 BYTES TO DEST FROM T/S D/E

RDLB3:         LD   HL,(SVHL)
               LD   DE,510
               ADD  HL,DE
               CALL CHKHL
               LD   (SVHL),HL
               LD   D,B
               LD   E,C        ;NEXT T/S
               JP   CLRRPT

CHKHL:         BIT  6,H
               RET  Z

               CALL INCURPAGE
               LD   (PORT1),A
               RET

FRDRD2:        PUSH DE         ;T/S
               PUSH HL         ;MAIN MEMORY PTR - DEST
               CALL RDADR      ;PT HL TO SRC IN RAMD, PAGE IN
               POP  DE         ;DEST
               PUSH AF         ;ORIG PAGE
               LD   C,A        ;MAIN MEM PAGE
               SCF             ;"510"
               EX   AF,AF'
               DEC  C          ;TO BE PAGED IN AT 4000H
               IN   A,(250)
               LD   B,A
               XOR  C
               AND  &E0
               XOR  C          ;A=VALUE FOR PORT 250 TO PAGE
               CALL &8002      ; IN DEST AT 4000H
               LD   B,(HL)
               INC  HL
               LD   C,(HL)     ;NEXT T/S
               POP  AF
               OUT  (251),A
               POP  DE         ;CURRENT T/S
               RET

SDCHK:         CALL GTBUF    ;PT HL TO BUFF (EITHER DRAM
                             ; OR 8000-BFFF)
SDCHK2:        LD   A,H
               CP   &80
               RET  C        ;RET IF HL IN DRAM
               PUSH HL
               LD   BC,&41FF
               ADD  HL,BC
               POP  HL
               RET  C        ;CY IF HL WILL CROSS PAGE BOUNDARY
               RES  7,H      ; - USE BUFF
               SET  6,H
               RET             ;ADDR IN 4000-7FFF AREA NOW

;AT 8002H

RDCODE:        LD   (&8000),SP
               LD   SP,&8200
               OUT  (250),A    ;PAGE IN DEST IN SECTION B
               PUSH BC         ;ORIG PORT 250 VALUE
               EX   AF,AF'
               LD   A,4        ;4 MOVES OF 128 BYTES
               JR   C,RDC1     ;JR IF 510 BYTES WANTED

               CALL &8020      ;MOVER - DO 128*4=512
               AND  A          ;NC

RDC1:          CALL C,&8024    ;MOVE 1FEH BYTES
               POP  AF
               OUT  (250),A
               LD   SP,(&8000)
RDCE:          RET

;SEARCH FOR FREE MR PAGE
;EXIT: IF Z, A=PAGE NUMBER, PAGE RESERVED

SFMRP:         XOR  A

SFMRL:         DEC  A
               CALL TMRBIT
               JR   Z,SMRBIT   ;JR IF FREE

               CP   &20
               JR   NZ,SFMRL

               AND  A          ;NZ
               RET

;TEST MEGA RAM BIT (PAGE). RETURN Z IF PAGE "A" FREE, NZ IF IN
;USE OR NON-EXISTENT. TABLE SET UP AT BOOT TIME, ALL 0FFHS IF NO
;MEGARAMS. ENTRY:A=00-FF

TMRBIT:        PUSH HL
               PUSH AF
               CALL MRADDR     ;GET ADDR, MASK
               AND  (HL)       ;Z IF BIT LOW (FREE)
               POP  HL
               LD   A,H
               POP  HL
               RET

;SET MEGA RAM BIT A. ALL REGS SAVED

SMRBIT:        PUSH AF
               PUSH HL
               CALL MRADDR
               OR   (HL)       ;NZ
               JR   SRMRC

;RESET MEGA RAM BIT A. ALL REGS SAVED

RMRBIT:        PUSH AF
               PUSH HL
               CALL MRADDR
               CPL               ;SET MASK BIT LOW
               AND  (HL)

SRMRC:         LD   (HL),A
               POP  HL
               POP  AF
               RET

;GET MEGA RAM TABLE ADDR IN HL, MASK IN A (BIT SET HI)
;ENTRY: A=BIT
;EXIT: HL AND A SET, OTHER REGS SAVED

MRADDR:        PUSH BC
               LD   B,A
               RRA
               RRA
               RRA
               AND  &1F      ;DIV BY 8
               LD   C,A
               LD   A,B
               AND  &07
               INC  A
               LD   B,A      ;B=1 TO 8 FOR BIT 7-0
               LD   A,1

MRBL:          RRCA
               DJNZ MRBL     ;A=80H IF B WAS 1, 40H IF B WAS 2..

               LD   HL,MRTAB
               ADD  HL,BC    ;PT TO BYTE (B=0)
               POP  BC
               RET

CNTFP:         CALL CFMRP    ;GET FREE MEGA RAM PAGES IN E (0-2)
               LD   B,E
               LD   HL,ALLOCT+FS+32

FRMRAL:        CALL RDA      ;READ A FROM SYS PAGE
               AND  A
               JR   NZ,FRMRD4

               INC  B        ;INC FREE PG CNT

FRMRD4:        DEC  L
               JR   NZ,FRMRAL

               LD   A,B
               RET

;COUNT FREE MEGA RAM PAGES - RESULT IN DE. BC SAVED

CFMRP:         LD   HL,MRTAB+4 ;AVOID USE OF MEGA RAM 0 1ST HALF
               LD   E,28       ;BYTES TO LOOK AT

;CALLED BY BOOT

CFMI:          PUSH BC
               XOR  A
               LD   D,A

CFPEL:         LD   C,(HL)
               INC  HL
               LD   B,8

CFPBL:         RR   C
               JR   C,CFP3

               INC  A        ;INC A BY NUMBER OF LOW (FREE) BITS
               JR   Z,CFPF   ;JR IF HIT 256 WHEN CALL FROM BOOT

CFP3:          DJNZ CFPBL    ;DO 8 BITS

               DEC  E
               JR   NZ,CFPEL ;DO 32 OR 28 BYTES

               DEC  D

CFPF:          INC  D
               LD   E,A
               POP  BC
               RET


;--------------
Part_HOOKS:


;HOOK 169 - PRINT TOKEN A

HPTV:          EXX               ;GET HL=PASSED IN XPTR
               CP   &FF
               JR   Z,HPFN

               SUB  246-FUNO
               PUSH AF
               PUSH HL         ;XPTR
               CALL NRRD
               DEFW FLAGS

               RRA
               CALL NC,SPC

               POP  BC         ;XPTR
               CALL NRWRD
               DEFW XPTR

               POP  BC

;ENTRY HERE AVOIDS LEADING SPACE. B=CHAR

PTV1:          LD   HL,KWMT
PTV2:          CALL PTVS
               JP   SPC

PTVS:          LD   A,(HL)
               INC  HL
               RLA
               JR   NC,PTVS

               DJNZ PTVS
               JP   PTM2       ;PRINT MSG FROM HL

HPFN:          LD   (XXPTR),HL ;SAVE XPTR TILL LATER - AVOID
                               ; PRINT KILLING ERROR MARKER
               LD   BC,251
               IN   E,(C)
               OUT  (C),B
               LD   B,E
               LD   HL,(CURCHL+FS)
               SET  7,H
               RES  6,H
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               LD   (OPSTORE+FS),DE
               LD   DE,XTRA    ;XTRA VECTORS - POST FF HOOK
               LD   (HL),D
               DEC  HL
               LD   (HL),E
               OUT  (C),B
               RET

;HOOK 170 - POST FF PRINT

HPFF:          PUSH AF         ;CHAR AFTER FF
               LD   BC,(XXPTR)
               CALL NRWRD
               DEFW XPTR

               POP  AF
               CP   &30
               JR   C,HPF3

               CP   &30+FUNO
               CCF
               JR   C,HPF3     ;JR IF NOT NEW FN
                               ;ELSE DEAL WITH 30H+
               SUB  &2F
               LD   BC,251
               IN   E,(C)
               OUT  (C),B
               LD   B,E
               LD   HL,(CURCHL+FS)
               SET  7,H
               RES  6,H
               LD   DE,(OPSTORE+FS)
               LD   (HL),E
               INC  HL
               LD   (HL),D     ;RESTORE NORMAL O/P ADDR
               OUT  (C),B
               LD   B,A
               LD   HL,KWMT
               CALL PTVS   ;PRINT NEW FN NAME - NO LEADING SPACE
               AND  A          ;NC=DONE IT

HPF3:          PUSH AF
               POP  BC
               RET

FUNO:          EQU  7
CMNO:          EQU  3

KWMT:          DEFB &A0
               DEFM "TIME"
               DEFB "$"+&80     ;30H $
               DEFM "DATE"
               DEFB "$"+&80     ;31H $
               DEFM "INP"
               DEFB "$"+&80     ;32H $
               DEFM "DIR"
               DEFB "$"+&80     ;33H $
               DEFM "FSTA"
               DEFB "T"+&80     ;34H N
               DEFM "DSTA"
               DEFB "T"+&80     ;35H N
               DEFM "FPAGE"
               DEFB "S"+&80     ;36H N
      ;      * DEFM "SCRA"
      ;      * DEFB "D"+&80     ;37H N
               DEFM "BACKU"
               DEFB "P"+&80     ;247
               DEFM "TIM"
               DEFB "E"+&80     ;248
               DEFM "DAT"
               DEFB "E"+&80     ;249
    ;        * DERM "ALTE"
    ;        * DEFB "R"+&80     ;250
    ;        * DEFM "SOR"
    ;        * DEFB "T"+&80     ;251

;HOOK 171 -  GET TOKEN

HGTTK:         LD   C,251
               IN   B,(C)
               PUSH BC
               XOR  A
               OUT  (251),A
               LD   HL,GTDT
               LD   DE,&8F00
               LD   BC,GTDTE-GTDT
               LDIR            ;EXTRA CODE TO BUFFER

               IN   A,(250)
               INC  A
               CALL SELURPG    ;PAGE KEYWORD LIST IN AT 8000H
               EXX
   ;         * LD   DE,(HKDE)
               LD   HL,KWMT-1+FS
               LD   A,FUNO+CMNO+1 ;ITEMS IN LIST+1
               CALL CMR
               DEFW JGTTOK

               POP  BC
               OUT  (C),B
               JR   Z,HGT2     ;RET IF NO MATCH

               CP   FUNO+1
               JR   NC,HGT1    ;JR IF CMD

               EX   DE,HL
               AND  A
               SBC  HL,DE
               ADD  A,&2F      ;FNS=30H+
               SCF
               JR   HGT2

HGT1:          ADD  A,&BB-FUNO ;247+

HGT2:          PUSH AF
               POP  BC
               RET

;COPIED TO 4F00H TO DEAL WITH FNS

GTDT:          POP  IY
               LD   BC,17
               ADD  IY,BC      ;PT TO 17 BYTES FURTHER ON
               POP  DE         ; IN TOKENISE SR
               ADD  HL,DE
               EX   DE,HL
               LD   (HL),&FF   ;FN LEADER PLACED IN BASIC LINE
               JP   (IY)

GTDTE:

;HOOK 172 - LENGTH FUNCTION PATCH

HKLEN:         EXX               ;HL=RET ADDR TO SCANNING
               INC  HL
               INC  HL
               INC  HL
               LD   C,(HL)     ;DISP TO "IMMED CODES"
               LD   B,0        ; (PART OF "JR")
               ADD  HL,BC
               PUSH HL         ;IMMED CODES
               LD   C,17
               ADD  HL,BC      ;PT TO NUMCONT
               LD   C,(HL)
               INC  HL
               LD   B,(HL)     ;BC=NUMCONT
      ;      * LD   A,(HKA)
               CP   &34-&1A
               JR   NC,HEVV2   ;JR IF NUMERIC RESULT

               LD   BC,6
               ADD  HL,BC
               LD   C,(HL)
               INC  HL
               LD   B,(HL)     ;BC=STRCONT

HEVV2:         PUSH AF
               CALL OWSTK   ;OVER-WRITE ADDR AFTER HOOK CODE ON
                            ;MAIN ROM STACK WITH NUMCONT OR
                            ;STRCONT ACCORDING TO TYPE OF RESULT

               CALL GTIXD
               POP  AF
               POP  HL         ;IMMED CODES
               CP   &3F-&1A    ;LENGTH
               JR   Z,FNLENG

               SUB  &30-&1A
               JP   C,REP0

               CP   FUNO
               JP   NC,REP0

               ADD  A,A
               LD   L,A
               LD   DE,DFNT

;USED BY MAIN HOOK ROUTINE. DE=TABLE, L=ENTRY

INDJP:         LD   H,0
               ADD  HL,DE
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               EX   DE,HL
               JP   (HL)

FNLENG:        LD   BC,10
               ADD  HL,BC
               LD   C,(HL)
               INC  HL
               LD   B,(HL)     ;IMFNTAB
               LD   HL,8
               ADD  HL,BC
               LD   E,(HL)
               INC  HL
               LD   D,(HL)     ;IMLENGTH
               IN   A,(250)
               OR   &40
               OUT  (250),A    ;ROM1 ON
               IN   A,(251)
               PUSH AF
               XOR  A
               OUT  (251),A    ;SYS PAGE AT 8000H
               EX   DE,HL
               LD   DE,&8D80
               LD   BC,&0078
               LDIR

               LD   HL,PDATA
               LD   C,15
               LDIR

               LD   A,&0B
               LD   (&8D80+&6C),A
               LD   HL,&4D80+&36
               LD   (&8D80+&24),HL
               POP  AF
               OUT  (251),A    ;ORIG
               CALL NRRDD
               DEFW CHADD

               INC  BC
               LD   A,(BC)
               CP   "#"
               JR   NZ,NRML

               CALL GTNC       ;PT TO "#"
               CALL GTNC       ;SKIP "#"
               CALL CMR
               DEFW EXPNUM

               CALL FABORT

;FNLENGTH
               JP   FNLN2

NRML:          CALL CMR
               DEFW &4D80  ;CALL MODIFIED "LENGTH" CODE IN CDBUF

               RET

;CODE TO PATCH "LENGTH" FOR LENGTH(0,A$) ERROR ON PAGE BOUNDARY

PDATA:         BIT  6,H
               JR   Z,PDAT2

               RES  6,H
               LD   A,(&5123)
               INC  A
               LD   (&5123),A

PDAT2:         JR   $-&56

DFNT:          DEFW FNTIME
               DEFW FNDATE
               DEFW INPST      ;INP$
               DEFW FNDIRS
               DEFW FSTAT
               DEFW DSTAT
               DEFW FPAGES
     ;       * DEFW SCRAD

FNDIRS:        CALL FDFSR      ;GET "ANY" NAME, GTDEF
               CALL GTNC
               CP   "("
               JR   Z,FNDI2

               CALL FABORT
               JR   FNDI3

FNDI2:         CALL EVNAMX
               LD   C,")"
               CALL ISEP
               CALL FABORT
               CALL EVFINS

FNDI3:         CALL DITOB      ;DIR TO BUFFER
               LD   HL,(PTRSCR)
               LD   DE,&A000
               AND  A
               SBC  HL,DE      ;HL=TEXT LEN (10 PER NAME)
               LD   B,H
               LD   C,L
               IN   A,(251)
               CALL HSTKS
               JP   PUTSCR

;DSTAT(D,X) - STATUS OF DISC D. IF D=*, USE DEFAULT DRIVE

;X=1 GIVES 0 IF WRITE PROTECTED OR NO FREE SLOTS, ELSE GIVES
;X=2 GIVES 1 IF WRITE PROTECTED                     ;FREE SPACE
;X=3 GIVES FREE SPACE (SECTS*510-9)
;X=4 GIVES FREE SLOTS
;X=5 GIVES TOTAL FILES
;X=6 GIVES FILES IN CURRENT DIR
;X=7 GIVES DTKS
;X=8 GIVES CURRENT DRIVE

DSTAT:         CALL SIBKS      ;INSIST "("
               CALL EVDNM      ;DRIVE NUMBER TO DSTR1
               CALL CNB        ;SKIP COMMA, GET X IN BC.
                               ; NO RET IF SYNTAX
               LD   A,C        ;PARAM
               CP   8
               LD   A,(DSTR1)  ;DRIVE
               JR   Z,NSTKAH   ;STACK CURRENT DRIVE IF PARAM=8

               SUB  2
               JR   NZ,DST0

               LD   A,(DVAR+2)
               AND  A
               JR   NZ,DST0    ;JR IF DRIVE 2 EXISTS

               CALL HOC0       ;GET AEHL=-1

HVAR2H:        JP   HVAR2      ;STACK-1

DST0:          CALL CKDRV
               LD   B,A
               PUSH BC
               CALL HOCHK
               POP  BC
               JR   Z,HVAR2H   ;IF NO DISC, STACK -1

               LD   A,B        ;DRIVE
               DEC  C
               JR   NZ,DST1    ;JR UNLESS SPECIAL WP/FREE
                               ; SLOT/SPACE
               CALL WPCHK
               LD   A,0
               JR   NZ,STKA    ;ZERO SPACE FREE IF WRITE PROT

               JR   DST15

DST1:          DEC  C
               JR   Z,DST3     ;JR TO CHECK WRITE PROTECT

               LD   A,C
               CP   6
               LD   A,(DRIVE)

NSTKAH:        JR   Z,STKA
               LD   A,C
               JP   NC,IOOR

DST15:         PUSH AF
               IN   A,(251)
               PUSH AF
               CALL ZDVS       ;ZERO VARS
               CALL DITOB      ;DIR TO BUFFER
               CALL STATS      ;GET HL=FREE SECTS, DE=FREE SLOTS
               ADD  HL,HL      ;FREE 256-BYTES
               POP  AF
               OUT  (251),A
               POP  AF
               EX   DE,HL
               AND  A
               JR   NZ,DST17
               LD   A,H
               OR   L

SHLHP:         JP   Z,STKHL    ;NO SPACE IF NO SLOTS
               JR   DST18

DST17:         CP   2
               JR   Z,SHLHP    ;STACK FREE SLOTS

               DEC  A
               JR   NZ,DST2

DST18:         LD   A,D
               OR   E
               JR   Z,STKA     ;JR IF NO FREE SECTORS

               LD   L,0
               LD   A,D
               LD   H,E        ;AHL=SECTS*512, DE=SECTS*2
               AND  A
               SBC  HL,DE
               SBC  A,0        ;AHL=SECTS*510
               LD   DE,9
               SBC  HL,DE
               SBC  A,0        ;ALLOW FOR HEADER
               JP   STKFP

DST2:          LD   HL,(TCNT)
               SUB  2
               JR   Z,DSSTK

               LD   HL,(FCNT)
               DEC  A
               JR   Z,DSSTK

               LD   A,(DTKS)
               JR   STKA

DST3:          CALL WPCHK

STKA:          LD   L,A
               LD   H,0

DSSTK:         JP   STKHL

;FREE PAGES FN

FPAGES:        CALL GTNC
               CALL FABORT
               CALL CNTFP
               JR   STKA

;ENTRY: A=DRIVE

WPCHK:         CP   3
               LD   A,0
               JR   NC,WPC2    ;JR IF RAM DISC

               CALL SELD
               ADD  A,2
               LD   C,A
               OUT  (C),A      ;IMPOSSIBLE SECTOR
               CALL PRECMX
               CALL BUSY
               RLCA
               RLCA             ;BIT 6 (WRITE PROT) TO BIT 0

WPC2:          AND  &01
               RET

SIBKS:         LD   C,"("
               JP   ISEPX

;INP$(#STREAM,NO. OF CHARS) FUNCTION

INPST:         CALL SIBKS      ;INSIST "("
               LD   C,"#"
               CALL ISEP
               CALL NNB        ;GET CHARS IN BC, STREAM IN DE.
                               ;NO RET IF SYNTAX TIME
               LD   HL,16
               AND  A
               SBC  HL,DE
               JP   C,INVST    ;LIMIT STREAM TO 0-16

               DEC  BC         ;Z->FFFF
               LD   A,B
               CP   &40
               JP   NC,IOOR

               PUSH BC
               CALL NRRDD
               DEFW CURCHL

               POP  HL         ;CHARS
               PUSH BC         ;CURCHL
               IN   A,(251)
               PUSH AF
               INC  HL         ;HL=1-4000H
               PUSH HL
               LD   A,E
               CALL CMR
               DEFW STREAM

               POP  BC
               CALL CMR
               DEFW WKROOM

               PUSH BC
               PUSH DE

INPSL:         PUSH BC
               PUSH DE
               IN   A,(251)
               PUSH AF
               XOR  A
               OUT  (251),A    ;SYS PAGE
               CALL MOVRC2
               JP   NC,REP27   ;EOF

               LD   D,A
               POP  AF
               OUT  (251),A    ;WKROOM
               LD   A,D
               POP  DE
               POP  BC
               LD   (DE),A
               INC  DE
               DEC  BC
               LD   A,B
               OR   C
               JR   NZ,INPSL

               POP  DE         ;START
               POP  BC         ;LEN
               IN   A,(251)
               EX   AF,AF'
               POP  AF
               OUT  (251),A    ;ORIG
               EX   AF,AF'
               CALL HSTKS      ;STACK ADEBC
               POP  BC         ;ORIG CURCHL
               CALL NRWRD
               DEFW CURCHL

               RET

;EVAL NUMBER, COMMA, NUMBER, BRACKET. ABORT IF SYNTAX, ELSE
;RETURN 1ST NUM IN DE, 2ND IN BC

NNB:           CALL EVNUM      ;GET NUM IN HL, Z IF SYNTAX

CNB:           PUSH HL
               LD   C,","
               CALL ISEP
               CALL EVNUM
               PUSH HL
               LD   C,")"
               CALL ISEP
               POP  BC         ;2ND
               POP  DE         ;1ST

FABORT:        CALL CFSO
               RET  NZ

               POP  HL
               RET

;FILE STATUS FUNCTION
;FSTAT("NAME",1)=FILE NUMBER
;FSTAT("NAME",2)=FILE LENGTH
;FSTAT("NAME",3)=FILE TYPE
;FSTAT("NAME",4)=FILE TYPE AND PROTECT/HIDE BITS

FSTAT:         CALL GTDEF
               CALL SIBKS      ;INSIST ON "("
               CALL EVNAM
               CALL CNB        ;COMMA,NUMBER, ")", ABORT

               PUSH BC         ;NUMBER
               DEC  BC         ;LEGAL VALUES NOW 0-3
               LD   HL,3
               AND  A
               SBC  HL,BC
               JP   C,IOOR

               CALL EVFINS
               CALL HOCHK
               POP  BC
               JP   Z,HVAR2    ;STACK -1 IF NO HOLE OR
                               ; NON-FORMATED RAMDISC
               PUSH BC
               CALL FINDC      ;LOOK FOR NAMED FILE, POINT
               POP  BC
               LD   A,0
               JR   NZ,STKAH   ;JR IF NOT FOUND - STACK 0

               DEC  C
               JR   Z,FST2

               LD   A,(HL)     ;TYPE
               DEC  C
               JR   Z,FST3

               DEC  C
               JR   NZ,STKAH   ;JR IF PARAM WAS 4

               AND  &1F

STKAH:         JP   STKA

FST2:          CALL CONM   ;CONVERT T/S IN D/E TO A NUMBER IN BC
               LD   H,B
               LD   L,C

STHLH:         JP   STKHL

FST3:          LD   DE,0
               AND  &1F
               CP   5
               LD   A,3
               JR   Z,FST4     ;JR IF 48K SNAP, ADE=48K

               LD   BC,&00EF
               ADD  HL,BC      ;POINT TO LEN DATA IN BUFFER
               LD   A,(HL)     ; (DIR ENTRY)
               INC  HL
               LD   E,(HL)
               INC  HL
               LD   D,(HL)

FST4:          EX   DE,HL      ;AHL=PAGE FORM

STK20B:        CALL AHLNX      ;GET 20-BIT NUMBER
               JP   STKFP

;SCRAD       * CALL GTNC
 ;           * CALL FABORT
  ;          * CALL NRRD
   ;         * DEFW CUSCRNP
    ;        * AND  &1F
     ;       * INC  A
      ;      * LD   HL,&8000
       ;     * JR   STK20B

HOCHK:         LD   A,(DRIVE)
               CP   3
               JR   C,HOC1

               CALL RTSTD
               AND  A
               RET  NZ         ;RET IF RAM DISC FORMATTED

HOC0:          LD   H,A
               LD   L,A
               JR   HOC2

HOC1:          LD   D,A        ;SIDE 1 (TRACK <80)
               CALL TFIHO      ;TEST FOR INDEX HOLE. Z IF NONE,
                               ; A=0, HL=0
               RET  NZ         ;RET IF HOLE FOUND

HOC2:          DEC  HL         ;VALUE
               LD   E,H        ;SGN BYTE. A=0 AND EHL=FFFFFF
               RET             ;Z                   ;IF NO HOLE

FNTIME:        LD   HL,TIMDT   ;TIME START
               DEFB &FD        ;"JR+3"

FNDATE:        LD   HL,DATDT   ;DATE START

;TIME/DATE FN COMMON

TDFNC:         PUSH HL
               CALL GTNC
               POP  HL         ;STR START
               CALL FABORT
               CALL RDCLK
               LD   BC,8       ;STR LEN
               JR   HPTH2

;HOOK CALL FOR PATH$

HPATH:         CALL GTDEF
               CALL GPATD

HPTH2:         SET  7,H
               RES  6,H
               EX   DE,HL
               IN   A,(250)
               INC  A          ;DOS PAGE

HSTKS:         CALL CMR
               DEFW STKSTR

               RET

;HOOK 173 - VIA CMDV - SAVE/LOAD/MERGE/VERIFY

HSLMV:         PUSH AF
               CALL NRRDD
               DEFW COMAD      ;CMD ADDRT TO BC

        ;    * LD   A,(HKA)
               POP  AF
               CALL NRWR
               DEFW CURCMD     ;WRITE CURCMD

               SUB  &90
               ADD  A,A
               LD   L,A
               LD   H,0
               ADD  HL,BC
               LD   C,250
               IN   B,(C)
               SET  6,B
               OUT  (C),B      ;ROM1 ON
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               IN   A,(251)
               PUSH AF
               XOR  A
               OUT  (251),A
               PUSH DE         ;SLVM ADDR

               LD   HL,MGPT
               LD   DE,&9000-MGPE+MGPT ;DEST IN BUFFER FOR MERGE
               LD   BC,MGPE-MGPT       ; PATCH CODE
               LDIR                    ;DE ENDS AT 9000H

               POP  HL
               LD   C,&58
               LDIR            ;DE ADV TO 8158H, SRC TO SLVM+58

               PUSH HL
               LD   HL,SVDT
               LD   C,MGPT-SVDT
               LDIR

               POP  HL
               LD   (&9000+&59),HL ;PATCH JP TO SLMV+58H

               LD   HL,(&9000+&25) ;READ SBFSR ADDR
               LD   (&9000+&5C),HL ;PATCH CALL SBFSR REPLACEMENT
               LD   HL,&5000+&5B
               LD   (&9000+&25),HL ;PATCH CALL SBFSR
               POP  AF
               OUT  (251),A            ;ORIG
             ;-LD   BC,&5000-MGPE+MGPT ;START OF MERGE PATCH IN
               LD   BC,&5000+MGPT-MGPE ;START OF MERGE PATCH IN

;OVER WRITE MAIN ROM STACK - FIRST ADDR WITH BC.

OWSTK:  ;    * LD   HL,(ENTSP)
       ;     * INC  HL
        ;    * INC  HL
        ;    * LD   E,(HL)
        ;    * INC  HL
        ;    * LD   D,(HL)
        ;    * EX   DE,HL      ;HL=MAIN ROM STACK PTR READ
                               ; FROM DOS STACK
               LD   HL,(&7FFC)
               JP   WRTBC

SVDT:          JP   0
               CALL 0          ;SBFSR
               LD   A,C
               LD   DE,&4F60
               LD   (DE),A
               DEC  DE         ;4F5F
               LD   HL,&4F4F
               LD   BC,&50   ;FOR 4F00-4F4F COPIED TO 4F10-4F5F
               LDDR           ;COPY TO SAFETY IN CASE
                             ;E.G. "S" CODE UDG " "
               INC  DE         ;4F10H
               LD   C,A
               CP   15
               RET  C

               LD   C,14
               RET

;CORRECT MERGE BUG (IF TOO MANY NUMS MERGE FOR GAP SIZE BETWEEN
;NUMEND AND SAVARS)

MGPT:          LD   A,(CURCMD)
               CP   &96        ;MERGE
               JR   NZ,MGP2

               LD   HL,(SAVARS)
               LD   DE,(NUMEND)
               LD   A,(NUMENDP)
               LD   C,A
               LD   A,(SAVARSP)
               CP   C
               JR   Z,MGP1     ;JR IF SAME PAGE

               RES  7,D
               SET  6,D        ;SUB 4000H
               AND  A

MGP1:          SBC  HL,DE      ;GET GAP SIZE
               LD   A,H
               CP   6
               JR   NC,MGP2    ;JR IF LOTS OF SPACE TO MERGE
                               ; NUMS (1.5K OR MORE)
               IN   A,(251)
               PUSH AF
               LD   A,C
               OUT  (251),A     ;PAGE IN NUMENDP
               XOR  A
               LD   BC,&0610    ;ABC=SPACE
               LD   HL,(NUMEND) ;MAKE ROOM HERE
               CALL JMKRBIG
               POP  AF
               OUT  (251),A

MGP2:          RST  &20
               LD   BC,&5FFA
               OUT  (C),B       ;ROM1 ON (CMR TURNS OFF ROM1)

MGPE:

;HOOK 174 - RUN/CLEAR PATCH

RCPTCH:        CALL FABORT
               IN   A,(251)
               PUSH AF
               CALL NRRDD
               DEFW NVARS

               LD   A,B
               CP   &BC
               JR   C,RCP3
                             ;NVARS IS TOO CLOSE TO END OF
                             ; PAGE - EXPAND
               CALL FPGE     ;POST PROG GAP TO PUT NVARS IN
               XOR  A        ; NEXT PAGE
               LD   BC,&0400
               CALL CMR
               DEFW JMKRBIG  ;NVARS NOW AT LEAST C000H->8000H
                             ; IN NEXT PAGE
               JR   RCP5

RCP3:          CP   &85
               JR   C,RCP5 ;RET IF NVARS TOO CLOSE TO PAGE START
                           ; FOR ANY POST-PROG GAP TO BE CLOSED.
               PUSH BC       ;NVARS
               CALL FPGE
               EX   DE,HL    ;DE=PROG END
               CALL NRRD
               DEFW NVARSP

               POP  HL
               LD   C,A      ;CHL=NVARS
               IN   A,(251)  ;PROG END P
               XOR  C
               AND  &1F
               JR   Z,RCP4   ;JR IF NORMAL SUBTRACT OK
                             ; - PAGES MATCH
               SET  6,H      ;ADD 4000H TO NVARS

RCP4:          SCF             ;NORMAL GAP=1, SO SCF FOR NORMAL=
               SBC  HL,DE
               JR   Z,RCP5   ;RET IF NO GAP TO CLOSE

               LD   B,H
               LD   C,L
               EX   DE,HL    ;RECLAIM AT PROG END
               CALL CMR
               DEFW JRECLAIM

               LD   (HL),&FF ;ENSURE PROG TERMINATED

RCP5:          POP  AF
               OUT  (251),A
               RET

FPGE:          CALL NRRDD
               DEFW PROG     ;GAP AFTER PROGRAM TO MOVE NVARS
                             ; TO NEXT PAGE
               PUSH BC
               CALL NRRD
               DEFW PROGP

               OUT  (251),A
               POP  HL

FPEL:          LD   D,(HL)
               INC  D
               RET  Z        ;RET IF TERMINATOR

               INC  HL
               INC  HL
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               INC  HL
               ADD  HL,DE
               CALL CHKHLR
               JR   FPEL

;CALL MAIN ROM - REWRITTEN TO ALLOW HOOKS TO BE USED BY ROUTINES
; CALLED BY CMR

CMR:           EXX
               POP  HL
               LD   E,(HL)
               INC  HL
               LD   D,(HL)
               INC  HL
               PUSH HL
               LD   HL,(&7FFC) ;MAIN ROM SP ON ENTRY TO DOS
               PUSH HL
               LD   HL,CMRF
               PUSH HL
               LD   C,A
               IN   A,(251)
               LD   H,A
               IN   A,(250)
               LD   B,A
               INC  A
               AND  &1F
               OUT  (251),A    ;DOS AT 8000H ALSO
               JP   CMR2+FS

CMR2:          LD   A,&1F
               DI
               OUT  (250),A       ;SYS PAGE AT 4000H
               LD   IY,(DOSSTK)
               LD   A,H           ;ORIG URPORT
               LD   HL,0
               ADD  HL,SP
               LD   (DOSSTK),HL   ;SO CALLS TO DOS DO NOT HIT
                                  ; CURRENT STACK
               LD   SP,(&7FFC+FS) ;NEW STACK JUST BELOW MAIN ROM
               EI                    ; STACK
               PUSH HL            ;ORIG SP
               PUSH BC            ;B=ORIG LRPORT
               PUSH IY            ;ORIG DOSSTK
               LD   HL,XTRA+L58A4-PFV ;RET ADDR TO PAGING SR
               PUSH HL                ; AFTER DEFKEYS
               PUSH DE            ;ROUTINE ADDR TO CALL
               JP   XTRA+CMR3-PFV ;PAGE IN ORIG URPORT,
                                  ; SET A REG, RET
CMRF:          POP  IY
               LD   (&7FFC),IY    ;ORIG IN CASE HOOK USE CHANGED
               RET                  ;IT AVOID STACK CREEP

;JUMPED TO AT BOOT TIME (INIP3+FS)

INIP3:         LD   A,&5B
               LD   BC,&D65B   ;LOOK FOR 5BD65B OF LD DE,(UMSG)
               LD   HL,&3D00   ;AT 3DB1 IN ROM 30
               CALL FTHREE+FS
               LD   DE,-4
               ADD  HL,DE      ;POSTFF ADDR
               PUSH HL

               LD   HL,PFV+FS
               LD   DE,XTRA
               LD   BC,LDEND-PFV
               LDIR            ;DATA TO AFTER DKLIM

               POP  HL
               LD   (XTRA+PFTRG+1-PFV),HL

               LD   HL,&3A31     ;:/1
               LD   (PTH1+FS),HL ;"1:"
               INC  L
               LD   (PTH2+FS),HL ;INIT PATH FOR DRIVE 2
               LD   A,(CKPT+FS)
               LD   C,A
               LD   B,&F0      ;CONTROL REG
               LD   A,&05      ;0101 NOT TEST/24H/RUN /REST
               OUT  (C),A      ;REST HI
               OUT  (C),A      ;WE CAN SET 24H NOW
               LD   A,&04      ;0100 NOT TEST/24H/RUN /NOT REST
               OUT  (C),A
               LD   HL,SYSP
               LD   (DOSSTK),HL
               LD   A,1
               LD   (DRIVE+FS),A

               LD   A,&CF
               LD   BC,&82C9
               LD   HL,&E200   ;E2B6 IN ROM 30
               CALL FTHREE+FS
               LD   DE,4
               ADD  HL,DE
               LD   (&4BB0+HLVTG-PVECT+1),HL ;LOAD/VERIFY HOOK
               RET                             ; PATCH

FTHREE:        LD   DE,0

F3OL:          PUSH BC         ;CHARS 2 AND 3
               LD   B,A

F3IL:          LD   A,D
               LD   D,E
               INC  HL
               LD   E,(HL)
               CP   B
               JR   NZ,F3IL

               POP  BC
               EX   DE,HL
               SBC  HL,BC
               EX   DE,HL
               JR   NZ,F3OL

               DEC  HL
               DEC  HL
               RET

;MEGA RAM INIT - USED AT BOOT TIME WHEN DOS AT 8000H.
; CALL MRINIT+FS

MRINIT:        DI
               LD   A,(DOSFLG)
               DEC  A
               OUT  (250),A     ;DOS PAGED IN AT 4000H AS WELL
               LD   (TEMPW1),SP ; AS 8000H (BOOT)
               LD   SP,&8000
               CALL MRIP2       ;CALL SR IN SECT B. EXIT WITH
               LD   HL,(TEMPW1) ; DE=MR PAGES
               LD   BC,&5FFA
               OUT  (C),B       ;NORMAL PAGE AT 4000H (ROM1 ON)
               LD   SP,HL
               LD   A,D
               OR   E
               RET  Z           ;RET IF NO FREE PAGES

               EX   DE,HL
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,HL
               ADD  HL,HL          ;*16=K (16-4096)
               LD   (TEMPW1+FS),HL

               RST  &28
               DEFB &2A            ;LKADDRW
               DEFW TEMPW1+FS
               DEFB EXIT

               LD   A,2
               CALL STREAM
               CALL JSTRS
               CALL &0013
               LD   DE,XMMSG+FS
               LD   BC,MRIP2-XMMSG
               JP   &0013          ;PRINT BC FROM DE

XMMSG:         DEFM "K External Memory"
               DEFB &0D

MRIP2:         IN   A,(251)
               PUSH AF
               OR   &80        ;EXTERNAL MEM BIT ON
               OUT  (251),A
               LD   C,0        ;START AT PAGE 0

MRIL:          LD   A,C
               OUT  (MRPRT),A  ;SELECT MEGA RAM (IF PRESENT)
               LD   HL,&8000
               XOR  A
               LD   (HL),A
               CP   (HL)
               JR   NZ,MRI2    ;JR IF CANNOT BE SET TO ZERO

               INC  A
               LD   (HL),A
               CP   (HL)
               JR   NZ,MRI2    ;JR IF CANNOT BE SET TO 0FFH

               LD   A,C
               CALL RMRBIT     ;RES MEGA RAM BIT - FREE
               LD   HL,0
               ADD  HL,SP

               EXX
               LD   SP,&C000
               LD   HL,0
               LD   BC,&0004

CMRL:          PUSH HL
               PUSH HL
               PUSH HL
               PUSH HL
               PUSH HL
               PUSH HL
               PUSH HL
               PUSH HL         ;CLEAR 16 BYTES
               DJNZ CMRL       ;DO IT 256 TIMES=4K

               DEC  C
               JR   NZ,CMRL    ;DO 4K FOUR TIMES

               EXX
               LD   SP,HL
               JR   MRI3

MRI2:          LD   A,C
               CALL SMRBIT     ;SET MEGA RAM BIT A - NOT PRESENT
L7CF8: ;* &7CFD
MRI3:          INC  C
               JR   NZ,MRIL    ;DO PAGES 00-FF

               POP  AF
               OUT  (251),A    ;EXTERNAL MEM BIT OFF

               LD   HL,MRTAB
               LD   E,32       ;BYTES TO LOOK AT
               JP   CFMI       ;COUNT FREE PAGES
L7D06: ;* &7D0B
;--------------
;Copy to &4BA0:
;SERIAL I/O SRS FOR OPENTYPE FILES

SERDT:         DEFB &40        ;"BLOCK MOVE SUPPORTED"
               JR   SERD2
;???
               RST  &08
               DEFB 151        ;O/P BC FROM DE
               RET

SERD2:         RST  &08
               DEFB HCHRWR
               RET
;???
               RST  &08
               DEFB HCHRRD
               EXX
               PUSH BC
               POP  AF
               RET
L4BAF: ;* &7D1A
               DEFB &00        ;???
;4BB0H
;PRINT TOKENS

PVECT:         CP   247
               RET  C

               POP  HL
               LD   HL,(XPTR)
               RST  &08
               DEFB 169        ;HPTV - DOS PRINT TOKEN
               RET

;EVALUATOR PATCH

EVV:           CP   &3F-&1A
               JR   Z,EVV2     ;RET UNLESS "LENGTH"

               CP   &3B-&1A
               RET  NC         ;DEAL WITH <3BH (PI)

EVV2:          POP  HL
               RST  &08
               DEFB 172        ;HKLEN - NO RETURN

;CMDV PATCH FOR SLVM

SLVP:          LD   HL,SYSP    ;7FF0H
               LD   (DOSSTK),HL
               CP   &94        ;SAVE
               RET  C          ;RET IF LESS THAN "SAVE"

               CP   &98
               JR   NC,RCP     ;JR IF NOT SLMV

               POP  HL
               RST  &08
               DEFB 173        ;HSLMV - "RET" TO MODIFIED CODE
                               ; IN BUFFER
RCP:           CP   &B3        ;CLEAR
               JR   Z,RCP2

               CP   &B0        ;RUN
               RET  NZ

RCP2:          PUSH AF
               RST  &08
               DEFB 174        ;RUN/CLEAR PATCH
               POP  AF
               RET

;NET PATCHES

HVEP:          AND  A          ;NC
               DEFB &3E        ;"JR+1"

HLDP:          SCF
               EXX
               LD   A,&FF
HLVTG:         JP   0          ;ROM LOAD/VERIFY

DTEND:

L7D4F: ;* &7D54
;--------------
;Copy to &5896:
;O/P ADDR USED WHEN CHAR AFTER FF IS PRINTED AT XTRA

PFV:           RST  &08
               DEFB 170      ;HPFF

               EXX
               PUSH BC
               POP  AF
PFTRG:         JP   C,0      ;JP POSTFF IF DOS DIDN'T HANDLE IT

               RET             ;JUST RET IF DONE IT

CMR3:          OUT  (251),A  ;ORIG URPORT
               LD   A,C      ;ORIG A
               EXX
               RET

L58A4:         EX   AF,AF'
               POP  IY
               LD   (DOSSTK),IY
               POP  AF         ;ORIG LRPORT
               POP  IY         ;ORIG SP
               DI
               LD   SP,IY
               JP   &0286      ;OUT (250),A: EI: EX AF,AF': RET

;MATCH TOKENS

MTV:           RST  &08
               DEFB 171

               EXX
               PUSH BC
               POP  AF
               RET  Z

               RET  NC         ;NEEDED?

               JP   &4F00      ;EXTENSION TO DEAL WITH FNS

LDEND:

Fix_L7D77_42:  ;4.2 Chg
;*           ;*DEFS 24                 ;- 24 NoClear
;*             DEFW 0,0,0,0,0,0,0,0    ;+ 24 Null
;*             DEFW 0,0,0,0            ;+

Fix_L7D7C_43:  ;4.3 Chg
             ;*DEFS 19                 ;- 19 NoClear
               DEFW 0,0,0,0,0,0,0,0    ;+ 19 Null
               DEFB 0,0,0              ;+

;FixEnd
;--------------

END:
;---------------------------------------------------------------
Len:           EQU  END-START
