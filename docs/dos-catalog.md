# DOS catalog

Empirical survey of the SAM Coupé corpus at `~/sam-corpus/disks/`.
For each disk, slot 0's file body — the loaded boot code, i.e.
the actual on-disk DOS — is read by walking its sector chain,
then SHA-256-hashed after stripping trailing 0x00 / 0xFF padding.
Two bodies with the same SHA are the same DOS; identifying which
DOS each SHA corresponds to is a manual step (disassembly /
source comparison), not something this survey tries to guess.

The slot-0 filename and the dir-entry file-type byte are dir
metadata, not part of the ROM ↔ DOS contract, so the survey
deliberately ignores them.

- Total disks scanned: **800**
- Disks with slot-0 DOS: **785**
- Disks with no slot-0 file (non-bootable archives): **15**
- Unique DOS bodies (by SHA-256, truncated to 16 hex chars): **410**
- Disks with malformed images: **0**

## Unique DOSes (most-used first)

| SHA-256 (16) | Disks | Body length (trimmed) | Body length (raw) |
|---|---:|---:|---:|
| `e22f10fffeba727f` | 107 | 8086 | 8160 |
| `827e541b8ad6c557` | 31 | 8086 | 8160 |
| `894f0bb155a4609b` | 30 | 8086 | 10200 |
| `0378541b547f6810` | 26 | 15735 | 15810 |
| `9067648326300f51` | 19 | 8086 | 10200 |
| `c4bd6d9e5f923af1` | 16 | 15690 | 15810 |
| `9dc893bae55e96fa` | 15 | 8086 | 10200 |
| `6bb51e2a43a187a1` | 13 | 8086 | 10200 |
| `64de9ac9b06b63ae` | 11 | 8086 | 10200 |
| `a5f6beac5230a201` | 10 | 10009 | 10200 |
| `116810d8f735510f` | 10 | 8086 | 10200 |
| `db155a1280ef3371` | 8 | 8086 | 10200 |
| `a4f908aa24edc985` | 8 | 8086 | 10200 |
| `1debdf0d8623b795` | 8 | 8086 | 10200 |
| `82a088a02982f94a` | 6 | 8086 | 8160 |
| `ca9afddd778a6a68` | 6 | 8086 | 10200 |
| `7c2a0cf35499b180` | 6 | 8086 | 10200 |
| `405f91b39e7b67cb` | 6 | 15690 | 15810 |
| `dd35dac3cced33c3` | 5 | 8086 | 10200 |
| `3c15c69dd62e2fc2` | 5 | 10165 | 10200 |
| `ee8b4ea2889eb3da` | 4 | 15735 | 15810 |
| `7b9b424bd6c0ca28` | 4 | 15740 | 15810 |
| `fa0fb0c9b7cd5f06` | 4 | 8086 | 8670 |
| `6782e00fc79e4303` | 4 | 8086 | 8160 |
| `54921a4c8710e66c` | 3 | 8086 | 10200 |
| `75b9e0bd60246fe1` | 3 | 8086 | 10200 |
| `5210e3d8ab5fd350` | 3 | 8164 | 10200 |
| `b268a4d7765c9c40` | 3 | 32299 | 32640 |
| `1b50e1b88428669a` | 3 | 8086 | 8160 |
| `e533d061fe37ce96` | 3 | 8086 | 10200 |
| `278a76211b2ed7a0` | 3 | 8086 | 10200 |
| `ffd57cff2a579224` | 3 | 8086 | 10200 |
| `ab91340c16daa333` | 3 | 8086 | 10200 |
| `05785df78ff3380d` | 3 | 8086 | 10200 |
| `abfaf0a2bfa7d3d1` | 3 | 1885 | 2040 |
| `2855102edf7e6ac8` | 2 | 509 | 1020 |
| `40abc664d94e05f8` | 2 | 8086 | 10200 |
| `101fc03c260bbab2` | 2 | 15808 | 15810 |
| `a9e83db7df3d7190` | 2 | 8086 | 10200 |
| `e281c4d7e3ba417d` | 2 | 8086 | 10200 |
| `f2e895e9a3552af7` | 2 | 8086 | 10200 |
| `2035bc16ee6f5935` | 2 | 8086 | 10200 |
| `7306966382dd7018` | 2 | 8086 | 10200 |
| `e56581e23cfa5d4a` | 2 | 510 | 510 |
| `f0e7082ce038ef26` | 2 | 25216 | 25500 |
| `d490eb86605e6b00` | 2 | 4080 | 4590 |
| `d51c0517b874dbf3` | 2 | 8086 | 10200 |
| `a3e23e57f466f3b2` | 2 | 8086 | 10200 |
| `60e33a9cf7c1adf5` | 2 | 8086 | 10200 |
| `40e62f1147e20a7c` | 2 | 8086 | 10200 |
| `83e1e3e42e830963` | 1 | 8086 | 10200 |
| `3f89b596def1d871` | 1 | 29599 | 30090 |
| `14067c6c4a35bc50` | 1 | 8086 | 10200 |
| `398d4ada9fb846e1` | 1 | 8086 | 10200 |
| `778efea49992e636` | 1 | 8086 | 10200 |
| `cc34cbb6d414bd4c` | 1 | 8086 | 10200 |
| `50910d0af28ac480` | 1 | 8086 | 10200 |
| `aa916895758bdf5d` | 1 | 2040 | 2550 |
| `da98ffd3b080f183` | 1 | 37229 | 37740 |
| `4eb81b995a2b3dbb` | 1 | 3362 | 3570 |
| `7a3202aad5e83e7a` | 1 | 8086 | 10200 |
| `bfc78d159d74383b` | 1 | 10200 | 10200 |
| `1a734ae773ad4552` | 1 | 13946 | 14280 |
| `301f82da7d6512f8` | 1 | 15800 | 15810 |
| `d456ed4b59e85664` | 1 | 742 | 1020 |
| `bd600bbf2daeefd1` | 1 | 510 | 510 |
| `333688b3dc517634` | 1 | 69 | 510 |
| `11d10e1743f5b5ba` | 1 | 504 | 510 |
| `05287b50f4014713` | 1 | 544 | 1020 |
| `311a4c303bbf0d68` | 1 | 458761 | 459000 |
| `05a0ee9cc5a17bdc` | 1 | 10009 | 10200 |
| `fbb9a22ca375f411` | 1 | 8086 | 10200 |
| `a10beb8927c96d2d` | 1 | 8086 | 10200 |
| `35680c24c92d8dd5` | 1 | 32640 | 32640 |
| `3eef01a939260ea4` | 1 | 8839 | 9180 |
| `9e43ea60dc10e9a7` | 1 | 12240 | 12750 |
| `7ecfa14f3664ed32` | 1 | 1539 | 2550 |
| `dd196204f7ae1e4e` | 1 | 2040 | 2550 |
| `b76dee6e3dc76b9d` | 1 | 8164 | 10200 |
| `77240f27c0603b8b` | 1 | 14133 | 14280 |
| `73fd9b77d3cb276b` | 1 | 5072 | 5100 |
| `c8314fac02c9af88` | 1 | 9 | 510 |
| `7d753d88494ae303` | 1 | 15740 | 15810 |
| `7aeb10851c891575` | 1 | 1020 | 1530 |
| `9d8cf00a8420e23d` | 1 | 988 | 1020 |
| `00af5392e0117f1f` | 1 | 147179 | 196860 |
| `448eef87513ef076` | 1 | 982 | 1020 |
| `bde6352f3e8e2e76` | 1 | 53550 | 54060 |
| `f15caed0430cf05a` | 1 | 8086 | 10200 |
| `fbd44d06d74d8ad6` | 1 | 1163 | 2040 |
| `7f49052f9526d7f5` | 1 | 735 | 1530 |
| `60eaf12f77ed2b5a` | 1 | 8086 | 10200 |
| `5f6e2e5879d51362` | 1 | 8086 | 10200 |
| `0e17d92bc897da37` | 1 | 14593 | 14790 |
| `45c03aaa2a266199` | 1 | 8086 | 10200 |
| `6d39630ab8d52950` | 1 | 53550 | 54060 |
| `20d494c977a3e620` | 1 | 2040 | 2550 |
| `98f9037bd9a44032` | 1 | 252 | 1020 |
| `da188f0c83ce6fd3` | 1 | 61878 | 62220 |
| `8f52ba40668ff4fb` | 1 | 78633 | 79050 |
| `0f8e8d0fb73b6157` | 1 | 3458 | 3570 |
| `fc676f5887cce004` | 1 | 67984 | 68340 |
| `6553c93d96602632` | 1 | 8086 | 9180 |
| `650b69947333c3c0` | 1 | 8086 | 10200 |
| `fef61954fc3b5118` | 1 | 8086 | 10200 |
| `047e0bdfc326984f` | 1 | 8086 | 8160 |
| `8ec687be626ffbb5` | 1 | 8086 | 8670 |
| `f325f3b860914620` | 1 | 8086 | 8160 |
| `e5e54d21ddfa7d8d` | 1 | 1020 | 1530 |
| `e76efb3b2370bfa9` | 1 | 44889 | 45390 |
| `2a52aa9e744db79a` | 1 | 11685 | 11730 |
| `246eca4b3eebc9d7` | 1 | 35964 | 36210 |
| `f7931d26f72ddbd7` | 1 | 30726 | 31110 |
| `299c19906774e76a` | 1 | 1404 | 3060 |
| `669dc5a92fa7bb6f` | 1 | 30268 | 30600 |
| `1e1b18a2a85fcd43` | 1 | 38347 | 38760 |
| `c9766802ca7b06e6` | 1 | 32753 | 33150 |
| `81971381b3eaec31` | 1 | 1313 | 3060 |
| `0d0b68e1e39e6c63` | 1 | 39942 | 40290 |
| `7fe9597073473650` | 1 | 35617 | 35700 |
| `6c9a4e93010ab464` | 1 | 19709 | 19890 |
| `ab7cf7256e4893b7` | 1 | 38944 | 39270 |
| `be7b7432d9835ab5` | 1 | 43769 | 43860 |
| `a59d962e903612aa` | 1 | 33349 | 33660 |
| `554882c83eb3a8a9` | 1 | 14751 | 14790 |
| `1fb754db9879be3e` | 1 | 11016 | 11220 |
| `cbdafe9e6f299457` | 1 | 40586 | 40800 |
| `eac4f8901d74db7b` | 1 | 9254 | 9690 |
| `94a408d570950ad0` | 1 | 34183 | 34680 |
| `37e3f9cf15beeea3` | 1 | 29074 | 29580 |
| `f12bb94f0d8b0336` | 1 | 41664 | 41820 |
| `a5ecc53e40004b5c` | 1 | 38271 | 38760 |
| `9d5eb06334725b0b` | 1 | 40197 | 40290 |
| `5dc1fad85888e0ce` | 1 | 35380 | 35700 |
| `541fe3a0d1abe785` | 1 | 31834 | 32130 |
| `ddac9f696e53e32e` | 1 | 36399 | 36720 |
| `33067aa63e0fa8f5` | 1 | 34301 | 34680 |
| `aac521ba31eebefc` | 1 | 39127 | 39270 |
| `ee5308a0adec1226` | 1 | 33490 | 33660 |
| `e7edec13f68563c8` | 1 | 36121 | 36210 |
| `6b06689f0ac032c1` | 1 | 36857 | 37230 |
| `66d8bdc0952ebcbf` | 1 | 36157 | 36210 |
| `e37ead4408db3079` | 1 | 32440 | 32640 |
| `ef07b5b782cb288f` | 1 | 4284 | 4590 |
| `2cd54ff421c04de6` | 1 | 5098 | 5100 |
| `40640e4536a6ab90` | 1 | 18360 | 18870 |
| `bdf1186a14b3cf67` | 1 | 81929 | 82110 |
| `3af88087f7e157ac` | 1 | 10165 | 10200 |
| `fd236cd923702a76` | 1 | 8086 | 10200 |
| `b1e9f36039aa7a9d` | 1 | 8086 | 10200 |
| `1c44a9d91c551e8e` | 1 | 8086 | 8160 |
| `8b85368ae38ea5fa` | 1 | 8086 | 10200 |
| `cd23461de7d045f4` | 1 | 36965 | 37230 |
| `3f35cec8e2bde547` | 1 | 5978 | 6120 |
| `ed16f57220106aa1` | 1 | 8086 | 10200 |
| `860ea8b1ad9be97f` | 1 | 8086 | 10200 |
| `98bea751220b489c` | 1 | 8086 | 10200 |
| `99e568983101d966` | 1 | 777 | 1020 |
| `a1c9d79ba2fcf38f` | 1 | 8086 | 10200 |
| `fdd0f7e9212b5be0` | 1 | 8086 | 10200 |
| `5512ee5f28dc1706` | 1 | 4069 | 4590 |
| `e3b0c44298fc1c14` | 1 | 0 | 510 |
| `1c39d04b5a8a4af6` | 1 | 10378 | 10710 |
| `8ab3a6307a919231` | 1 | 8086 | 10200 |
| `eb4399b247e0434a` | 1 | 15740 | 15810 |
| `afe10942f94ba6d4` | 1 | 2040 | 2550 |
| `3f3883bc03bed501` | 1 | 2040 | 2550 |
| `a692ae9d23a57656` | 1 | 1530 | 2040 |
| `2f62bf26ab644326` | 1 | 8086 | 10200 |
| `61d175867f73fb2a` | 1 | 1019 | 1530 |
| `367ca9df53acc4bf` | 1 | 1530 | 2040 |
| `764da7e9d8ac88a6` | 1 | 36073 | 36720 |
| `f7cd1ed21b71ee69` | 1 | 1998 | 2550 |
| `f7bc69664a87883b` | 1 | 8086 | 10200 |
| `5a1781d443c35839` | 1 | 8086 | 10200 |
| `ee8ff35bde01fe49` | 1 | 1530 | 2040 |
| `4be07efe19b7cb0f` | 1 | 8086 | 10200 |
| `78ef60fcf8a1e715` | 1 | 53550 | 54060 |
| `a3f470d48e7559f4` | 1 | 5859 | 6120 |
| `b4bc39e6c8cc560e` | 1 | 2195 | 2550 |
| `3c3b3f7046455122` | 1 | 2039 | 2550 |
| `54a6d2b12e0c17b3` | 1 | 53550 | 54060 |
| `998be3b72a986452` | 1 | 1530 | 2040 |
| `3fbf90ed0c97796e` | 1 | 18360 | 18870 |
| `132f4f07e0d85ed7` | 1 | 12907 | 13260 |
| `2e68bd68c8b20da5` | 1 | 4053 | 4590 |
| `46f89fa54b2f33f7` | 1 | 2040 | 2550 |
| `d9d113eb76c51f63` | 1 | 13200 | 13260 |
| `2f96f4b6cf96f7e9` | 1 | 8086 | 10200 |
| `ec06160c88987e5d` | 1 | 32440 | 32640 |
| `55c64d9534412f7c` | 1 | 107325 | 107610 |
| `99015aeaa1d2f763` | 1 | 154362 | 154530 |
| `2dfdf40f1a72973d` | 1 | 5646 | 6120 |
| `bacacd3e207ce221` | 1 | 8738 | 9180 |
| `90a45ecc5c12aad7` | 1 | 1020 | 1530 |
| `cecd5cdf87336d25` | 1 | 2040 | 2550 |
| `6bf5445940a8a76d` | 1 | 9702 | 10200 |
| `225b2304bdbedad2` | 1 | 8086 | 10200 |
| `2341c9493b5dbeeb` | 1 | 8086 | 10200 |
| `01f788a493ca03bf` | 1 | 14009 | 14280 |
| `3992aaec02aab425` | 1 | 8086 | 10200 |
| `050325c2c98552c8` | 1 | 16393 | 16830 |
| `16b08b5a6164003d` | 1 | 8086 | 10200 |
| `4eab21a727463c3b` | 1 | 15735 | 15810 |
| `6f083938079ca1ec` | 1 | 8086 | 10200 |
| `5f3550d1ca8989f0` | 1 | 8086 | 10200 |
| `5793ce5635564711` | 1 | 8086 | 10200 |
| `7b8cbc6085ac4be3` | 1 | 18441 | 18870 |
| `09a809b2c3f6a881` | 1 | 24625 | 24990 |
| `a7a83fd615a00232` | 1 | 15690 | 15810 |
| `541ffa1c199769e4` | 1 | 8086 | 9180 |
| `7588eb488ea3963d` | 1 | 8086 | 10200 |
| `623d94e58d851687` | 1 | 24625 | 24990 |
| `56cbdab3f9d1f6fe` | 1 | 164250 | 164730 |
| `6e355ad8d6d322e4` | 1 | 510 | 510 |
| `9f7f247d0f209596` | 1 | 1530 | 1530 |
| `c3415c19c2d25f9b` | 1 | 15693 | 15810 |
| `e9382a64e59ffd93` | 1 | 18659 | 18870 |
| `09510607cff5eedf` | 1 | 4735 | 5100 |
| `227f1272dc95b5be` | 1 | 24990 | 25500 |
| `df40a1eaf134685b` | 1 | 14361 | 14790 |
| `b2a05cc13487f8df` | 1 | 38473 | 38760 |
| `9709bdf03809e0ae` | 1 | 16827 | 17340 |
| `a49f846d7e922479` | 1 | 566 | 1020 |
| `410c7abb0ab561b1` | 1 | 2040 | 2550 |
| `364df841599ec69e` | 1 | 48213 | 48450 |
| `d010d2dd7fc4b45c` | 1 | 15868 | 16320 |
| `d4fc5229fa786b26` | 1 | 8086 | 10200 |
| `d6db2024f3b55251` | 1 | 2040 | 2550 |
| `397794db3e9032a4` | 1 | 2039 | 2550 |
| `f2222766ea225e48` | 1 | 2040 | 2550 |
| `edfac77a1256a479` | 1 | 1020 | 1530 |
| `54788c4e1d5a3d30` | 1 | 12239 | 12750 |
| `8213d749e526c8c7` | 1 | 11726 | 12750 |
| `ddc9534d9f14fe4c` | 1 | 15800 | 15810 |
| `f71ecf6155c4c639` | 1 | 32300 | 32640 |
| `e1cb1393429f6606` | 1 | 8086 | 10200 |
| `cad0f3208a632a36` | 1 | 8086 | 10200 |
| `bfd0de70f0212791` | 1 | 2331 | 2550 |
| `6b3819e43403182a` | 1 | 1018 | 1530 |
| `110e242c901d54f5` | 1 | 1530 | 2040 |
| `f3a35085259dbe41` | 1 | 8086 | 10200 |
| `2590f1c411796e2d` | 1 | 46 | 510 |
| `90517122bf2481c5` | 1 | 8086 | 10200 |
| `e76eeeea7692fb07` | 1 | 37230 | 37740 |
| `b0ebdee710caeb5b` | 1 | 3674 | 4080 |
| `821df6540285fd9e` | 1 | 44370 | 44880 |
| `0c34a31ade2849d0` | 1 | 2040 | 2550 |
| `ff4eb979d939874d` | 1 | 1017 | 1530 |
| `3a2e943bac812b12` | 1 | 3603 | 4080 |
| `c43e949b399d2fce` | 1 | 4825 | 5100 |
| `d695ae0e9211769f` | 1 | 45545 | 45900 |
| `7d5073f9468416e2` | 1 | 12240 | 12750 |
| `133c99f187347302` | 1 | 11729 | 12750 |
| `d07853609211e15e` | 1 | 1530 | 2040 |
| `098c0f672c7fbe3f` | 1 | 4080 | 4590 |
| `0449b0f60d31acc7` | 1 | 1530 | 2040 |
| `ea84c2e9ec5c1973` | 1 | 1020 | 1530 |
| `d9feb1ddc0e17771` | 1 | 8086 | 10200 |
| `59c63c8de8e3edca` | 1 | 2040 | 2550 |
| `4ef8faa15ddf815c` | 1 | 8086 | 10200 |
| `a09dfe6fd32bddcc` | 1 | 8086 | 8160 |
| `10c694d747b12705` | 1 | 17358 | 17850 |
| `8e9b529fa294ea21` | 1 | 8086 | 10200 |
| `b3c3ee71cc95d3a5` | 1 | 4080 | 4590 |
| `5ae7ecf1649b8972` | 1 | 3183 | 3570 |
| `200278d35e1e7d4b` | 1 | 46 | 510 |
| `6c7cb5197d3034af` | 1 | 8086 | 10200 |
| `3d231309c58b9c7b` | 1 | 8086 | 10200 |
| `2ab28418f3673bef` | 1 | 8086 | 10200 |
| `08317b88daef6576` | 1 | 8086 | 10200 |
| `bb17425f57b905f7` | 1 | 8086 | 10200 |
| `041b2d1f1dc791cd` | 1 | 165 | 510 |
| `42c388a53f73ca44` | 1 | 8086 | 10200 |
| `8fb01afeca36dbcd` | 1 | 665 | 1020 |
| `9f32255832a86dd5` | 1 | 34435 | 34680 |
| `153b4c1b750833e1` | 1 | 2040 | 2550 |
| `59a2a3377c71df6d` | 1 | 54118 | 54570 |
| `78345b4235aae4dd` | 1 | 4735 | 5100 |
| `8969e9a30b1e5863` | 1 | 53550 | 54060 |
| `3db3af4d4a245344` | 1 | 18360 | 18870 |
| `f9d0ed730f76838b` | 1 | 15695 | 15810 |
| `4fff977d8013377f` | 1 | 73575 | 73950 |
| `e9dd512c58f870b4` | 1 | 8086 | 10200 |
| `1e009491ead9a1e9` | 1 | 8086 | 8160 |
| `23fd8a30002c7978` | 1 | 25735 | 26010 |
| `111177e603253ae9` | 1 | 255 | 510 |
| `1e6670db1775b582` | 1 | 8086 | 10200 |
| `e850f05bd65f5a12` | 1 | 8086 | 10200 |
| `632db00d73f2bcdb` | 1 | 8086 | 10200 |
| `b3e093d53d3e8428` | 1 | 8086 | 10200 |
| `74b0fcff076393f9` | 1 | 8086 | 10200 |
| `9813a2b223cd3390` | 1 | 8086 | 10200 |
| `9548a3aa5d7098f4` | 1 | 8086 | 10200 |
| `4e0b7599aa852830` | 1 | 8086 | 10200 |
| `7722a00047e84140` | 1 | 8086 | 10200 |
| `f70d22ec2642fb03` | 1 | 8086 | 10200 |
| `a465e0b79fb55157` | 1 | 8086 | 10200 |
| `999a129332b2916b` | 1 | 8086 | 10200 |
| `cf06e3c38e44072d` | 1 | 8086 | 10200 |
| `8ba62ef306bf80d5` | 1 | 8086 | 10200 |
| `b33a4a2558f9125d` | 1 | 8086 | 10200 |
| `03182e8ae606f607` | 1 | 8086 | 10200 |
| `d09a2f4f45505f97` | 1 | 8086 | 10200 |
| `6c8ba91597ed8c3c` | 1 | 8086 | 10200 |
| `b5888b0cd6e79794` | 1 | 8086 | 10200 |
| `9d6f3291cdc4f938` | 1 | 8086 | 10200 |
| `06f000fe436a5a67` | 1 | 8086 | 10200 |
| `d80749fad2c5ec95` | 1 | 8086 | 10200 |
| `0920b0f727f516cf` | 1 | 8086 | 10200 |
| `adbfcbfe401f73ff` | 1 | 8086 | 10200 |
| `eef4d45e68ad2765` | 1 | 8086 | 10200 |
| `ae4c23bb6344861d` | 1 | 8086 | 10200 |
| `3893517bfe602537` | 1 | 8086 | 10200 |
| `d64a50beb45e4415` | 1 | 8086 | 10200 |
| `e064fcc49fe12839` | 1 | 8086 | 10200 |
| `d0d935496574bcb4` | 1 | 8086 | 10200 |
| `30facf34e3ae8416` | 1 | 8086 | 10200 |
| `85fe103f9df025d2` | 1 | 8086 | 10200 |
| `92514ebc590f5d1d` | 1 | 8086 | 10200 |
| `0b8b4046fa7f5d3a` | 1 | 8086 | 10200 |
| `6d56e93be9d356ba` | 1 | 8086 | 10200 |
| `d36b77540f6f7237` | 1 | 8086 | 10200 |
| `2a8bb90a7cce3e6c` | 1 | 12240 | 12750 |
| `e76cba1ae4eee4da` | 1 | 12237 | 12750 |
| `a020a583ec5c5a57` | 1 | 2009 | 2550 |
| `9424150ac355dcce` | 1 | 2040 | 2550 |
| `30be24dbcfcae6a3` | 1 | 10698 | 10710 |
| `aff6bbf5bcc5ad46` | 1 | 4080 | 4590 |
| `fde198d34bd1af0f` | 1 | 5289 | 5610 |
| `71c9b9601b80bf06` | 1 | 2330 | 2550 |
| `a9ff3d1e798d8cfc` | 1 | 4590 | 5100 |
| `9919ef4aac9605c2` | 1 | 9801 | 10200 |
| `1d4d39218c139c85` | 1 | 8086 | 10200 |
| `93ae595946cbfc07` | 1 | 8086 | 10200 |
| `1a0266804ed4c785` | 1 | 8086 | 10200 |
| `cffb14eb8520690f` | 1 | 8086 | 10200 |
| `d01f9b12d501b19e` | 1 | 24350 | 24480 |
| `912af1d501dd8de0` | 1 | 30235 | 30600 |
| `23d231cec05f234e` | 1 | 30296 | 30600 |
| `f067bee72231d802` | 1 | 30296 | 30600 |
| `263f716215abf909` | 1 | 36865 | 37230 |
| `f7d93f36612a7685` | 1 | 3570 | 4590 |
| `eface0f0b955164b` | 1 | 1085 | 2040 |
| `22b364ea5beabbed` | 1 | 132 | 510 |
| `ba53f23f387c40f0` | 1 | 2040 | 2550 |
| `4371857cb647d5c3` | 1 | 3761 | 4080 |
| `1b08f26b7abce60b` | 1 | 1527 | 2040 |
| `2de501366d3e37f8` | 1 | 2040 | 2550 |
| `4ffe30c8b8bcd872` | 1 | 1020 | 2040 |
| `f8f3afcc1330ebfe` | 1 | 2039 | 2550 |
| `be9f4198d9851f5f` | 1 | 510 | 1530 |
| `b05585902262df5b` | 1 | 60 | 510 |
| `08c2c5aa041946d5` | 1 | 1039 | 1530 |
| `08cc627ae4ea1d14` | 1 | 8093 | 8160 |
| `9ed339a39a563d41` | 1 | 1530 | 2550 |
| `fbd3854b912a449e` | 1 | 39330 | 39780 |
| `721551247cd4f9ca` | 1 | 2570 | 3060 |
| `1670a7d17b0132c1` | 1 | 2040 | 2550 |
| `6adf935424f52438` | 1 | 8086 | 10200 |
| `2ea0d9173dddd47f` | 1 | 7909 | 8160 |
| `51b1c68aeddaf053` | 1 | 9181 | 10200 |
| `d1b70a1445b0372f` | 1 | 4543 | 4590 |
| `ca4f5f99c00d0f7b` | 1 | 51662 | 52020 |
| `376164478e5c066b` | 1 | 6992 | 7140 |
| `bbd08e667a089486` | 1 | 169 | 510 |
| `17578664befe3147` | 1 | 30618 | 31110 |
| `0c67c69e9f94e664` | 1 | 16644 | 16830 |
| `a851cb71b60b313b` | 1 | 11154 | 11220 |
| `ef3fda2cfc8f57d9` | 1 | 29790 | 30090 |
| `3a3eb37d99089ed4` | 1 | 22184 | 22440 |
| `4b4590a60fbcbfc5` | 1 | 5185 | 5610 |
| `845ff405604811e5` | 1 | 53513 | 54060 |
| `c8a6c7c087a2dc4f` | 1 | 227 | 510 |
| `4bdb9f49e49834a2` | 1 | 74142 | 74460 |
| `6589bc32b3c2f92e` | 1 | 74141 | 74460 |
| `9756db858428a533` | 1 | 74142 | 74460 |
| `71660009a2b74579` | 1 | 74141 | 74460 |
| `d36ed4a4ddfb1341` | 1 | 18360 | 18870 |
| `7cf97cba81af3631` | 1 | 12240 | 12750 |
| `2a701ee27032f6f5` | 1 | 15690 | 15810 |
| `b4aca8aff4efdffc` | 1 | 8086 | 10200 |
| `03cd50b5e50c5bb9` | 1 | 8018 | 8160 |
| `5fb3e7460965c75c` | 1 | 12240 | 12750 |
| `425bce275aec6a70` | 1 | 333 | 510 |
| `4051f72e2146bfe9` | 1 | 24621 | 24990 |
| `8e7932233ceecf38` | 1 | 1425 | 1530 |
| `4b3a6b20f3184ede` | 1 | 2461 | 2550 |
| `eda9270d197b7064` | 1 | 20910 | 21420 |
| `9acde11c04c4dfa8` | 1 | 12230 | 12750 |
| `2d94fc4c29675c4d` | 1 | 3700 | 4080 |
| `651ceca0418ff661` | 1 | 4080 | 4590 |
| `204baca3dad8428c` | 1 | 122910 | 123420 |
| `c6dd90280914b377` | 1 | 15692 | 15810 |
| `69a211a4de1747c7` | 1 | 521 | 1020 |
| `ed45d7b6dea661f3` | 1 | 37173 | 37740 |
| `d75d94c539bfd8f2` | 1 | 8086 | 10200 |
| `32ef028957538851` | 1 | 8086 | 10200 |
| `e47caa45f77e0f3b` | 1 | 8086 | 10200 |
| `15a8e15fd680dbf3` | 1 | 4161 | 4590 |
| `ed635e2282aec119` | 1 | 20198 | 38250 |
| `90c5cca7472ca621` | 1 | 631 | 1020 |
| `0df25ed6fe3f0156` | 1 | 9319 | 10200 |
| `473bdd91856e0a3d` | 1 | 8086 | 10200 |
| `d373434c5da09096` | 1 | 8086 | 10200 |
| `6e691c7af719521f` | 1 | 4639 | 5100 |
| `923d7c33e6239ab1` | 1 | 8086 | 10200 |
| `fdc9c5194995c4d9` | 1 | 67438 | 67830 |
| `1e9cec8cf567b429` | 1 | 11053 | 11220 |
| `e59feafa029fbe18` | 1 | 8086 | 10200 |

## Sample disks per DOS

Sample list of corpus disks for each unique DOS (up to 5).
Use any of these to extract the body for disassembly:

```bash
# Extract slot 0's body for disassembly:
python3 ~/git/samfile/tools/audit/extract_dos.py <hash-prefix>
```

### `e22f10fffeba727f` (107 disks, body=8086 bytes)

- Blast Turbo_ by James R Curry (1995) (PD)
- COMMIX V2.00 by S. Grodkowski (1995) (PD)
- COMMIX V2.01 by S. Grodkowski (1995) (PD)
- COMMIX V2.02 by S. Grodkowski (1995) (PD)
- Easydisc V4.9 (1995) (Saturn Software)

### `827e541b8ad6c557` (31 disks, body=8086 bytes)

- Sam Adventure Club Issue 01 (Nov 1991)
- Sam Adventure Club Issue 02 (Jan 1992)
- Sam Adventure Club Issue 03 (Mar 1992) _a1_
- Sam Adventure Club Issue 03 (Mar 1992)
- Sam Adventure Club Issue 04 (May 1992) _a1_

### `894f0bb155a4609b` (30 disks, body=8086 bytes)

- 32 Colour Demo by Gordon Wallis (1992) (PD)
- Allan Stevens - Home Utilities (1994)
- Allan Stevens - Home Utilities - Seven Pack (1994)
- Comms Loader (19xx)
- F-16 Combat Pilot Demo (1991) (PD)

### `0378541b547f6810` (26 disks, body=15735 bytes)

- E-Tracker Program Disk V1.2 (19xx) (FRED Publishing)
- Pics from the Net 06 (19xx) (PD)
- Pics from the Net 08 (19xx) (PD)
- Pics from the Net 09 (19xx) (PD)
- Pics from the Net 10 (19xx) (PD)

### `9067648326300f51` (19 disks, body=8086 bytes)

- Curse of the Serpent_s Eye_ The (1994) (Dream World) _a1_
- Curse of the Serpent_s Eye_ The (1994) (Dream World) _a2_
- Curse of the Serpent_s Eye_ The (1994) (Dream World) _a3_
- Curse of the Serpent_s Eye_ The (1994) (Dream World)
- FRED Magazine Issue 02 (1990)

### `c4bd6d9e5f923af1` (16 disks, body=15690 bytes)

- Spectrum 128 Music Disk 2 (19xx) (PD)
- Spectrum Games Compilation 02 (1992)
- Spectrum Games Compilation 03 (1992)
- Spectrum Games Compilation 04 (1992)
- Spectrum Games Compilation 05 (1992)

### `9dc893bae55e96fa` (15 disks, body=8086 bytes)

- Metempsychosis Demo - Christine (19xx)
- Metempsychosis Demo - Highlander (19xx)
- Metempsychosis Unreleased Demo - Internal_highlander (19xx)
- Metempsychosis Unreleased Demo - Internal_joy_pen (19xx)
- Metempsychosis Unreleased Demo - Mega_mix (19xx)

### `6bb51e2a43a187a1` (13 disks, body=8086 bytes)

- Cheats By Paul Crompton (1994) (PD)
- Driver Icons (1995) (Saturn Software)
- Sam Adventure Club Library Disk 3 SamScratch V3.2 (19xx)
- Sam Supplement Magazine Issue 06 (Mar 1991)
- Sam Supplement Magazine Issue 20 (May 1992)

### `64de9ac9b06b63ae` (11 disks, body=8086 bytes)

- Allan Stevens - Capricorn Software Disk 1 (1994)
- Allan Stevens Compilation - Gallery 1 (19xx)
- Allan Stevens Compilation - Gallery 2 (19xx)
- Allan Stevens Compilation - Gallery 3 (19xx)
- Allan Stevens Compilation - Games Disk 1 (19xx)

### `a5f6beac5230a201` (10 disks, body=10009 bytes)

- Banzai - The Games Compliation by Dan Doore (1995) (PD) _a1_
- Banzai - The Games Compliation by Dan Doore (1995) (PD)
- Banzai Babes 1 - Cindy _ Claudia by Dan Doore (1994) (PD)
- Banzai Babes 2 - Claudia - Elle - Kate by Dan Doore (1994) (PD)
- Banzai Pictures II - The Atari Job by Dan Doore (1994) (PD)

### `116810d8f735510f` (10 disks, body=8086 bytes)

- Bats _n_ Balls Demo by Lord Insanity (1992) (PD)
- Craft Demo by ESI (1992) (PD)
- Hexagonia Demo by Fuxoft (1992) (PD)
- Lettis Demo by Daniel Cannon (1992) (PD)
- Sam News Disk Issue 1 (Jan 1992) (SAM Computers LTD)

### `db155a1280ef3371` (8 disks, body=8086 bytes)

- Adventures of Captain Comic_ The (19xx) (Lars)
- Defender - Persona and Digital Reality (1998) (Chris Pile)
- Dino Sourcerer (1993) (Softdisk Inc.)
- Morse Code Tutor (1999) (R.J. Wilkins) (PD)
- Outwrite V2 (1992) (Rj Will Kinson)

### `a4f908aa24edc985` (8 disks, body=8086 bytes)

- Astroball Demo by Balor Knight (1992) (PD)
- Ice Chicken Demo by ESI (1995) (PD)
- Manic Miner (1992) (Revelation Software)
- SCPDSA Demo Disk 1 (1992) (PD)
- Splat Demo by Colin Jordan (1991) (PD)

### `1debdf0d8623b795` (8 disks, body=8086 bytes)

- FRED Magazine Issue 35 (1993) _a1_
- FRED Magazine Issue 35 (1993)
- FRED Magazine Issue 36 (1993)
- FRED Magazine Issue 37 (1993)
- FRED Magazine Issue 38 (1993) _a1_

### `82a088a02982f94a` (6 disks, body=8086 bytes)

- Best of ENCELADUS_ Birthday Pack Edition (19xx) (Relion)
- ENCELADUS Magazine Issue 05 (Jun 1991) (Relion Software)
- ENCELADUS Magazine Issue 06 (Aug 1991) (Relion Software) _a1_
- ENCELADUS Magazine Issue 06 (Aug 1991) (Relion Software)
- ENCELADUS Magazine Issue 07 (Oct 1991) (Relion Software)

### `ca9afddd778a6a68` (6 disks, body=8086 bytes)

- FRED Magazine Issue 15 (1991)
- FRED Magazine Issue 20 (1992) _a1_
- FRED Magazine Issue 20 (1992)
- FRED Magazine Issue 21 (1992) _a1_
- FRED Magazine Issue 21 (1992)

### `7c2a0cf35499b180` (6 disks, body=8086 bytes)

- Metempsychosis Unreleased Demo - CD_demo (19xx)
- Metempsychosis Unreleased Demo - Internal_advert (19xx)
- Metempsychosis Unreleased Demo - Internal_small_demos (19xx)
- Metempsychosis Unreleased Demo - Jukebox (19xx)
- Metempsychosis Unreleased Demo - RTJ_pdm3 (19xx)

### `405f91b39e7b67cb` (6 disks, body=15690 bytes)

- Mouse Flash 1.1 for MDOS (19xx)
- Outwrite V2.0 (1992) (Chezron Software)
- PAX Disk 1 (1996) (Glenco)
- SC Filer V2.0 (1991) (Steve_s Software)
- Sam Adventure System Test Disk (1992) (Axxent Software)

### `dd35dac3cced33c3` (5 disks, body=8086 bytes)

- Allan Stevens - 50 Programs to Play and Write (19xx)
- Allan Stevens - Capricorn Software Disk 2 (1994)
- Allan Stevens - Colour Cycle (19xx)
- Allan Stevens - Spectrum Screens (19xx)
- Allan Stevens - World Finder (1994)

### `3c15c69dd62e2fc2` (5 disks, body=10165 bytes)

- Blitz Magazine Issue 2 (1997) (Persona)
- Blitz Magazine Issue 4A (1997) (Persona)
- COMET to ASCII by Simon Cooke (1995)
- FRED Magazine Issue 82 (1997)
- Fashoom_ (1997) (Sad Snail Productions)

### `ee8b4ea2889eb3da` (4 disks, body=15735 bytes)

- Amiga MODS Disk (19xx)
- ETrackerv1.2
- Interlaced RGB Viewer Pics (19xx) (PD)
- Metempsychosis Sample Disk 7 (19xx)

### `7b9b424bd6c0ca28` (4 disks, body=15740 bytes)

- FRED Magazine - Morkography (1992)
- Flight of Fantasy and Occult Connection Adventures (19xx)
- MasterDOS V2.1 (19xx)
- Recover-E (1995)

### `fa0fb0c9b7cd5f06` (4 disks, body=8086 bytes)

- Mono Clipart Samples V1.0 (Nov 1995) (Steve_s Software)
- SC Compressor 2 (1991) (Steve_s Software)
- SC Monitor Pro 1.2_ TurboMon 1.0 (1992) (Steve_s Software)
- SC PD 3 by Steve_s Software (1992) (PD)

### `6782e00fc79e4303` (4 disks, body=8086 bytes)

- Screens Viewer Disk 1 (19xx)
- Screens Viewer Disk 2 (19xx)
- Screens Viewer Disk 3 (19xx)
- Screens Viewer Disk 4 (19xx)

### `54921a4c8710e66c` (3 disks, body=8086 bytes)

- Allan Stevens Compilation - Games Disk 3 (19xx)
- FRED Magazine Issue 02 (1990) _a1_
- Neil Holmes_ Boing_ Graphics (1992) (Noesis Software)

### `75b9e0bd60246fe1` (3 disks, body=8086 bytes)

- Arcadia Disk Magazine _1 (1991)
- Arcadia Disk Magazine _2 (1991)
- Arcadia Disk Magazine _3 (1991)

### `5210e3d8ab5fd350` (3 disks, body=8164 bytes)

- Blinky Samples Disk 1 (1997) (Edwin Blink)
- Blinky Samples Disk 3 (1997) (Edwin Blink)
- Blitz 6 Menu by Edwin Blink (1997) (PD)

### `b268a4d7765c9c40` (3 disks, body=32299 bytes)

- DS12 Duff Capers Music Demo (2003) (PD)
- DS7 RainBow Scroller (2002) (PD)
- Duff Capers v0.51 by Tobermory (2003) (PD)

### `1b50e1b88428669a` (3 disks, body=8086 bytes)

- ENCELADUS Magazine Issue 09 (Feb 1992) (Relion Software)
- ENCELADUS Magazine Issue 11 (Jun 1992) (Relion Software)
- ENCELADUS Magazine Issue 12 (Oct 1992) (Relion Software)

### `e533d061fe37ce96` (3 disks, body=8086 bytes)

- FRED Magazine Issue 01 (1990) _a1_
- FRED Magazine Issue 01 (1990)
- FRED Magazine Issue 01 _ 02 (1990)

### `278a76211b2ed7a0` (3 disks, body=8086 bytes)

- FRED Magazine Issue 12 (1991) _a1_
- FRED Magazine Issue 12 (1991) _a2_
- FRED Magazine Issue 12 (1991)

### `ffd57cff2a579224` (3 disks, body=8086 bytes)

- Golden Sword of Bhakhor_ The (1997) (Persona)
- Sam Supplement Magazine Issue 45 (Jun 1994)
- Sam Supplement Magazine Issue 46 (Jul 1994)

### `ab91340c16daa333` (3 disks, body=8086 bytes)

- Sam Paper Magazine Issue 6 (19xx)
- Sam Paper Magazine Issue 8 (19xx)
- Sam Paper Magazine Issue 9 (19xx)

### `05785df78ff3380d` (3 disks, body=8086 bytes)

- Visually 2 (19xx) (Zenith Graphics)
- Visually 3 (19xx) (Zenith Graphics)
- Visually 5 (19xx) (Zenith Graphics)

### `abfaf0a2bfa7d3d1` (3 disks, body=1885 bytes)

- Zeddy - ZX81 Emu and Programs 1 (19xx) (PD)
- Zeddy - ZX81 Emu and Programs 2 (19xx) (PD)
- Zeddy - ZX81 Emu and Programs 3 (19xx) (PD)

### `2855102edf7e6ac8` (2 disks, body=509 bytes)

- All Star Belch by Tobermory (2001) (PD)
- Blancmange Burps by Tobermory (2001) (PD)

### `40abc664d94e05f8` (2 disks, body=8086 bytes)

- Allan Stevens Compilation - Spectrum Disk 1 (19xx)
- Sam CD2 Utility (1990) (Kobrahsoft)

### `101fc03c260bbab2` (2 disks, body=15808 bytes)

- B-DOS V1.7N (1999) (Martijn Groen _ Edwin Blink) (PD)
- Blancmange Burps 2 - Smell the Glove_ by Tobermory (2001) (PD)

### `a9e83db7df3d7190` (2 disks, body=8086 bytes)

- Banzai - The Demos _ Utils by Dan Doore (1994) (PD)
- Banzai - The Double by Dan Doore (1994) (PD)

### `e281c4d7e3ba417d` (2 disks, body=8086 bytes)

- Best of KAPSA I-V_ The (1993) (KAPSA) _a1_
- Best of KAPSA I-V_ The (1993) (KAPSA)

### `f2e895e9a3552af7` (2 disks, body=8086 bytes)

- Edition 3 (1991) (Zenith Graphics) (PD)
- Zenith Edition 3 (19xx) (Zenith Graphics)

### `2035bc16ee6f5935` (2 disks, body=8086 bytes)

- Fastline Public Domain Library Disk 13 (19xx)
- Silly Demo 1 by Lord Insanity (1990) (PD)

### `7306966382dd7018` (2 disks, body=8086 bytes)

- Impatience - Triltex Viking (1991) (Fred Publishing)
- Occult Adventure by David Munden (1993) (PD)

### `e56581e23cfa5d4a` (2 disks, body=510 bytes)

- Mike AJ Demo Disk 1 (1991)
- Mike AJ Disc 7 (19xx)

### `f0e7082ce038ef26` (2 disks, body=25216 bytes)

- Sam Adventure Club Issue 09b (Mar 1993) _a1_
- Sam Adventure Club Issue 09b (Mar 1993)

### `d490eb86605e6b00` (2 disks, body=4080 bytes)

- Sam Coupe Demo Collection Disk 10 (1992) (PD) _b1_
- Sam Coupe Demo Collection Disk 10 (1992) (PD)

### `d51c0517b874dbf3` (2 disks, body=8086 bytes)

- Sam Supplement Magazine - The Best Of (19xx)
- Sam Supplement Magazine Issue 12A (Sep 1991)

### `a3e23e57f466f3b2` (2 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 33 (Jun 1993)
- Sam Supplement Magazine Issue 34 (Jul 1993)

### `60e33a9cf7c1adf5` (2 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 37 (Oct 1993)
- Sam Supplement Magazine Issue 38 (Nov 1993)

### `40e62f1147e20a7c` (2 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 50 (Nov 1994) _b1_
- Sam Supplement Magazine Issue 50 (Nov 1994)

### `83e1e3e42e830963` (1 disks, body=8086 bytes)

- 18 Rated Poker for 512k (19xx) (Supplement Software)

### `3f89b596def1d871` (1 disks, body=29599 bytes)

- AMRAD Amateur Radio Logbook (1994) (Spencer)

### `14067c6c4a35bc50` (1 disks, body=8086 bytes)

- Aliens vs Predator Demo by Gordon Wallis (1991) (PD)

### `398d4ada9fb846e1` (1 disks, body=8086 bytes)

- Allan Stevens - Capricorn Software Disk 3 Unfinished (1994)

### `778efea49992e636` (1 disks, body=8086 bytes)

- Allan Stevens - Learn With Timmy Under 6s (1994)

### `cc34cbb6d414bd4c` (1 disks, body=8086 bytes)

- Allan Stevens - Learning Games (19xx)

### `50910d0af28ac480` (1 disks, body=8086 bytes)

- Allan Stevens Compilation - Games Disk 4 (19xx)

### `aa916895758bdf5d` (1 disks, body=2040 bytes)

- Alternative Module Player _ Tunes V1.2 (1994) (Aley Keprt)

### `da98ffd3b080f183` (1 disks, body=37229 bytes)

- Andy Monk_s Music (19xx) (PD)

### `4eb81b995a2b3dbb` (1 disks, body=3362 bytes)

- Arcadia Disk Magazine _3b (1991)

### `7a3202aad5e83e7a` (1 disks, body=8086 bytes)

- Arcadia Disk Magazine _5 (1991)

### `bfc78d159d74383b` (1 disks, body=10200 bytes)

- B-DOS V1.5A (1997) (Martijn Groen _ Edwin Blink) (PD)

### `1a734ae773ad4552` (1 disks, body=13946 bytes)

- B-DOS V1.7D (1999) (Martijn Groen _ Edwin Blink) (PD)

### `301f82da7d6512f8` (1 disks, body=15800 bytes)

- B-DOS V1.7J (1999) (Martijn Groen _ Edwin Blink) (PD)

### `d456ed4b59e85664` (1 disks, body=742 bytes)

- Banzai Babes 1 - Cindy _ Claudia by Dan Doore (1994) (PD) _a1_

### `bd600bbf2daeefd1` (1 disks, body=510 bytes)

- Banzai Pictures I by Dan Doore (1994) (PD)

### `333688b3dc517634` (1 disks, body=69 bytes)

- Basic Scrolly Pokers by Chris White (1992)

### `11d10e1743f5b5ba` (1 disks, body=504 bytes)

- Bats _n_ Balls by David Gommerman (1992) (Revelation)

### `05287b50f4014713` (1 disks, body=544 bytes)

- Beta DOS 1.0 For the Plus D (1990) (Betasoft)

### `311a4c303bbf0d68` (1 disks, body=458761 bytes)

- Blinky Samples Disk 2 (1997) (Edwin Blink)

### `05a0ee9cc5a17bdc` (1 disks, body=10009 bytes)

- Blinky Samples Disk 4 (1997) (Edwin Blink)

### `fbb9a22ca375f411` (1 disks, body=8086 bytes)

- Blitz Magazine Issue 4B - Sound Machine (1997) (Persona) _b1_

### `a10beb8927c96d2d` (1 disks, body=8086 bytes)

- Blokker (19xx) (Stephen McGreal)

### `35680c24c92d8dd5` (1 disks, body=32640 bytes)

- Blondie and Dagwood_ Arkanoid_ Prince of the Yolk Folk (199x)

### `3eef01a939260ea4` (1 disks, body=8839 bytes)

- Blue Disk Show_ The (19xx)

### `9e43ea60dc10e9a7` (1 disks, body=12240 bytes)

- Boing _ Sphera (1992) (Noesis Software)

### `7ecfa14f3664ed32` (1 disks, body=1539 bytes)

- Bombed Out_Nuclear Waste_Magic Caves_Blockade (1990) (Enigma)

### `dd196204f7ae1e4e` (1 disks, body=2040 bytes)

- Bowin and the Count Dracula (1991) (Lucosoft and Revelation)

### `b76dee6e3dc76b9d` (1 disks, body=8164 bytes)

- Bunnik 2K MOD Slide (2000) (Edwin Blink)

### `77240f27c0603b8b` (1 disks, body=14133 bytes)

- COMET Opcodes by Tobermory (2001) (PD)

### `73fd9b77d3cb276b` (1 disks, body=5072 bytes)

- COMET Z80 Assembler V1.8 (1992) (Revelation)

### `c8314fac02c9af88` (1 disks, body=9 bytes)

- Colour Clipart Samples V1.0 (Nov 1995) (Steve_s Software)

### `7d753d88494ae303` (1 disks, body=15740 bytes)

- CometAssembler1.8EdwinBlink

### `7aeb10851c891575` (1 disks, body=1020 bytes)

- Compressor Utilities (1993) (ESI) (PD)

### `9d8cf00a8420e23d` (1 disks, body=988 bytes)

- Contact Demos by Chris White (1991) (PD)

### `00af5392e0117f1f` (1 disks, body=147179 bytes)

- DS4 80 Pixel Demo by Tobermory _ Andrew Collier (2001) (PD)

### `448eef87513ef076` (1 disks, body=982 bytes)

- DS5 SIN Wave Distort Toy by Tobermory _ Steve Taylor (2001) (PD)

### `bde6352f3e8e2e76` (1 disks, body=53550 bytes)

- Daton MasterBASIC Demos (1991) (PD)

### `f15caed0430cf05a` (1 disks, body=8086 bytes)

- Defender Compilation (19xx)

### `fbd44d06d74d8ad6` (1 disks, body=1163 bytes)

- Demo Collection 1 (19xx) (PD)

### `7f49052f9526d7f5` (1 disks, body=735 bytes)

- Demo Collection 2 (19xx) (PD)

### `60eaf12f77ed2b5a` (1 disks, body=8086 bytes)

- Demo Disk (1994) (PD)

### `5f6e2e5879d51362` (1 disks, body=8086 bytes)

- Diaz Demo 2 (19xx)

### `0e17d92bc897da37` (1 disks, body=14593 bytes)

- Dissasembler (1993) (Chris White)

### `45c03aaa2a266199` (1 disks, body=8086 bytes)

- Domino Box (1992) (Supplement Software)

### `6d39630ab8d52950` (1 disks, body=53550 bytes)

- Dyzonium_ WaterWorks_ WOP Gamma_ Bugulators (19xx)

### `20d494c977a3e620` (1 disks, body=2040 bytes)

- E-Copier V2.0 (1993) (Chris White)

### `98f9037bd9a44032` (1 disks, body=252 bytes)

- E-Demo by Simon Cooke (1992) (PD)

### `da188f0c83ce6fd3` (1 disks, body=61878 bytes)

- E-Mag Demo (1993) (PD)

### `8f52ba40668ff4fb` (1 disks, body=78633 bytes)

- E-Tracker Program Disk (19xx) (FRED Publishing)

### `0f8e8d0fb73b6157` (1 disks, body=3458 bytes)

- E-Tracker Tunes (19xx) (ESI)

### `fc676f5887cce004` (1 disks, body=67984 bytes)

- E-Tunes Player (19xx) (Andrew Collier)

### `6553c93d96602632` (1 disks, body=8086 bytes)

- ENCELADUS - Complete Guide to SAMBASIC Parts 1-7 (1994) (Relion)

### `650b69947333c3c0` (1 disks, body=8086 bytes)

- ENCELADUS Magazine Issue 01 (Oct 1990) (Relion Software)

### `fef61954fc3b5118` (1 disks, body=8086 bytes)

- ENCELADUS Magazine Issue 02 (Dec 1990) (Relion Software)

### `047e0bdfc326984f` (1 disks, body=8086 bytes)

- ENCELADUS Magazine Issue 03 (Feb 1991) (Relion Software)

### `8ec687be626ffbb5` (1 disks, body=8086 bytes)

- ENCELADUS Magazine Issue 04 (Apr 1991) (Relion Software)

### `f325f3b860914620` (1 disks, body=8086 bytes)

- ENCELADUS Magazine Issue 10 (Apr 1992) (Relion Software)

### `e5e54d21ddfa7d8d` (1 disks, body=1020 bytes)

- ESI Demos (19xx) (PD)

### `e76efb3b2370bfa9` (1 disks, body=44889 bytes)

- EXPLOSION - 0 - GAMES FOR EMULATOR EXPLOSION (1996)

### `2a52aa9e744db79a` (1 disks, body=11685 bytes)

- EXPLOSION - 1 - GAMES FOR EMULATOR EXPLOSION (1996)

### `246eca4b3eebc9d7` (1 disks, body=35964 bytes)

- EXPLOSION - 2 - GAMES FOR EMULATOR EXPLOSION (1996)

### `f7931d26f72ddbd7` (1 disks, body=30726 bytes)

- EXPLOSION - 3 - GAMES FOR EMULATOR EXPLOSION (1996)

### `299c19906774e76a` (1 disks, body=1404 bytes)

- EXPLOSION - 4 - GAMES FOR EMULATOR EXPLOSION (1996)

### `669dc5a92fa7bb6f` (1 disks, body=30268 bytes)

- EXPLOSION - 5 - GAMES FOR EMULATOR EXPLOSION (1996)

### `1e1b18a2a85fcd43` (1 disks, body=38347 bytes)

- EXPLOSION - 6 - GAMES FOR EMULATOR EXPLOSION (1996)

### `c9766802ca7b06e6` (1 disks, body=32753 bytes)

- EXPLOSION - 7 - GAMES FOR EMULATOR EXPLOSION (1996)

### `81971381b3eaec31` (1 disks, body=1313 bytes)

- EXPLOSION - 8 - GAMES FOR EMULATOR EXPLOSION (1996)

### `0d0b68e1e39e6c63` (1 disks, body=39942 bytes)

- EXPLOSION - A - GAMES FOR EMULATOR EXPLOSION (1996)

### `7fe9597073473650` (1 disks, body=35617 bytes)

- EXPLOSION - B - GAMES FOR EMULATOR EXPLOSION (1996)

### `6c9a4e93010ab464` (1 disks, body=19709 bytes)

- EXPLOSION - C - GAMES FOR EMULATOR EXPLOSION (1996)

### `ab7cf7256e4893b7` (1 disks, body=38944 bytes)

- EXPLOSION - D - GAMES FOR EMULATOR EXPLOSION (1996)

### `be7b7432d9835ab5` (1 disks, body=43769 bytes)

- EXPLOSION - E - GAMES FOR EMULATOR EXPLOSION (1996)

### `a59d962e903612aa` (1 disks, body=33349 bytes)

- EXPLOSION - F - GAMES FOR EMULATOR EXPLOSION (1996)

### `554882c83eb3a8a9` (1 disks, body=14751 bytes)

- EXPLOSION - G - GAMES FOR EMULATOR EXPLOSION (1996)

### `1fb754db9879be3e` (1 disks, body=11016 bytes)

- EXPLOSION - H - GAMES FOR EMULATOR EXPLOSION (1996)

### `cbdafe9e6f299457` (1 disks, body=40586 bytes)

- EXPLOSION - I - GAMES FOR EMULATOR EXPLOSION (1996)

### `eac4f8901d74db7b` (1 disks, body=9254 bytes)

- EXPLOSION - J - GAMES FOR EMULATOR EXPLOSION (1996)

### `94a408d570950ad0` (1 disks, body=34183 bytes)

- EXPLOSION - K - GAMES FOR EMULATOR EXPLOSION (1996)

### `37e3f9cf15beeea3` (1 disks, body=29074 bytes)

- EXPLOSION - L - GAMES FOR EMULATOR EXPLOSION (1996)

### `f12bb94f0d8b0336` (1 disks, body=41664 bytes)

- EXPLOSION - M - GAMES FOR EMULATOR EXPLOSION (1996)

### `a5ecc53e40004b5c` (1 disks, body=38271 bytes)

- EXPLOSION - N - GAMES FOR EMULATOR EXPLOSION (1996)

### `9d5eb06334725b0b` (1 disks, body=40197 bytes)

- EXPLOSION - O - GAMES FOR EMULATOR EXPLOSION (1996)

### `5dc1fad85888e0ce` (1 disks, body=35380 bytes)

- EXPLOSION - P - GAMES FOR EMULATOR EXPLOSION (1996)

### `541fe3a0d1abe785` (1 disks, body=31834 bytes)

- EXPLOSION - Q - GAMES FOR EMULATOR EXPLOSION (1996)

### `ddac9f696e53e32e` (1 disks, body=36399 bytes)

- EXPLOSION - R - GAMES FOR EMULATOR EXPLOSION (1996)

### `33067aa63e0fa8f5` (1 disks, body=34301 bytes)

- EXPLOSION - S - GAMES FOR EMULATOR EXPLOSION (1996)

### `aac521ba31eebefc` (1 disks, body=39127 bytes)

- EXPLOSION - T - GAMES FOR EMULATOR EXPLOSION (1996)

### `ee5308a0adec1226` (1 disks, body=33490 bytes)

- EXPLOSION - U - GAMES FOR EMULATOR EXPLOSION (1996)

### `e7edec13f68563c8` (1 disks, body=36121 bytes)

- EXPLOSION - V - GAMES FOR EMULATOR EXPLOSION (1996)

### `6b06689f0ac032c1` (1 disks, body=36857 bytes)

- EXPLOSION - X - GAMES FOR EMULATOR EXPLOSION (1996)

### `66d8bdc0952ebcbf` (1 disks, body=36157 bytes)

- EXPLOSION - Z - GAMES FOR EMULATOR EXPLOSION (1996)

### `e37ead4408db3079` (1 disks, body=32440 bytes)

- EXPLOSION - ZX SPECTRUM 48 Emulator _ COMMANDER (1996)

### `ef07b5b782cb288f` (1 disks, body=4284 bytes)

- EXPLOSION MELODY DISK - FOR SELECT A CREATE BLOCK (1996)

### `2cd54ff421c04de6` (1 disks, body=5098 bytes)

- Easydisc V5.0 (1995) (Saturn Software)

### `40640e4536a6ab90` (1 disks, body=18360 bytes)

- Edwin Blink_s Samples (1991) (PD)

### `bdf1186a14b3cf67` (1 disks, body=81929 bytes)

- Entropy Demo (1992) (PD) _a1_

### `3af88087f7e157ac` (1 disks, body=10165 bytes)

- Entropy Demo (1992) (PD) _a2_

### `fd236cd923702a76` (1 disks, body=8086 bytes)

- Entropy Demo (1992) (PD)

### `b1e9f36039aa7a9d` (1 disks, body=8086 bytes)

- FRED Magazine Issue 28 (1992)

### `1c44a9d91c551e8e` (1 disks, body=8086 bytes)

- FRED Magazine Issue 33 (1993)

### `8b85368ae38ea5fa` (1 disks, body=8086 bytes)

- FRED Magazine Issue 34 (1993)

### `cd23461de7d045f4` (1 disks, body=36965 bytes)

- FRED Magazine Issue 65 Menu (1995)

### `3f35cec8e2bde547` (1 disks, body=5978 bytes)

- FRED _ SamCo Demo Disk (1990) (PD)

### `ed16f57220106aa1` (1 disks, body=8086 bytes)

- Fastline Public Domain Library Disk 10 (19xx)

### `860ea8b1ad9be97f` (1 disks, body=8086 bytes)

- Fastline Public Domain Library Disk 12 (19xx)

### `98bea751220b489c` (1 disks, body=8086 bytes)

- Fastline Public Domain Library Disk 16 (19xx)

### `99e568983101d966` (1 disks, body=777 bytes)

- Font Loader (1991) (Phantom Software)

### `a1c9d79ba2fcf38f` (1 disks, body=8086 bytes)

- Football League Manager (1994) (Key Software-FRED Publishing)

### `fdd0f7e9212b5be0` (1 disks, body=8086 bytes)

- Fredatives 4 (1992)

### `5512ee5f28dc1706` (1 disks, body=4069 bytes)

- GFX Demo by Doug Holmes (1991) (PD)

### `e3b0c44298fc1c14` (1 disks, body=0 bytes)

- Game Dos - DVar 8.92n3 (19xx)

### `1c39d04b5a8a4af6` (1 disks, body=10378 bytes)

- GamesMaster 1.2 (1992) (Betasoft)

### `8ab3a6307a919231` (1 disks, body=8086 bytes)

- Goblin Mountain Adventure (1993) (Sam PD Sware Lib.) (PD)

### `eb4399b247e0434a` (1 disks, body=15740 bytes)

- H-DOS V2.12 HD Loader V2.0 (1996)

### `afe10942f94ba6d4` (1 disks, body=2040 bytes)

- Hexagonia _b1_ _ Witching Hour (19xx)

### `3f3883bc03bed501` (1 disks, body=2040 bytes)

- Highway Code_ Quizball _ Editor (1991) (Revelation)

### `a692ae9d23a57656` (1 disks, body=1530 bytes)

- ICONMASTER V1.0 by Steve Taylor (1993) (Revelation)

### `2f62bf26ab644326` (1 disks, body=8086 bytes)

- Images SAM Software - Roboblob _ Give a Dog a Bone (19xx)

### `61d175867f73fb2a` (1 disks, body=1019 bytes)

- In Comet Format Master V1.2 by Chris White (1990)

### `367ca9df53acc4bf` (1 disks, body=1530 bytes)

- In Comet Format Work Copy V1.2 by Chris White (1990)

### `764da7e9d8ac88a6` (1 disks, body=36073 bytes)

- Infinity - E-Tracker Crack Menu Demo by Entropy (1993) (PD)

### `f7cd1ed21b71ee69` (1 disks, body=1998 bytes)

- Integrated Logic_s Madonna Strip Show (1990) (PD)

### `f7bc69664a87883b` (1 disks, body=8086 bytes)

- Integrated Logic_s PD Disk (1990) (PD)

### `5a1781d443c35839` (1 disks, body=8086 bytes)

- Jigsaw Disk - Jigsaw Creator (1992) (Colony Software)

### `ee8ff35bde01fe49` (1 disks, body=1530 bytes)

- Juggler Demo (1990) _f1_

### `4be07efe19b7cb0f` (1 disks, body=8086 bytes)

- KEdisk V2.04 (19xx)

### `78ef60fcf8a1e715` (1 disks, body=53550 bytes)

- Kim Wilde - Gary Moore Samples (1990) (PD)

### `a3f470d48e7559f4` (1 disks, body=5859 bytes)

- Lemmings - Assemble Disk 1 (19xx) (Chris White)

### `b4bc39e6c8cc560e` (1 disks, body=2195 bytes)

- Lemmings - DOS Test Side 1 (19xx) (Chris White)

### `3c3b3f7046455122` (1 disks, body=2039 bytes)

- Lemmings - Front End GFX (19xx) (Chris White)

### `54a6d2b12e0c17b3` (1 disks, body=53550 bytes)

- Lemmings - GFX IFF 1 (19xx) (Chris White)

### `998be3b72a986452` (1 disks, body=1530 bytes)

- Lemmings - GFX IFF 2 (19xx) (Chris White)

### `3fbf90ed0c97796e` (1 disks, body=18360 bytes)

- Lemmings - Stuff From DMA (19xx) (Chris White)

### `132f4f07e0d85ed7` (1 disks, body=12907 bytes)

- Lemmings Raw Data Disk 1 (1992) (DMA-Chris White)

### `2e68bd68c8b20da5` (1 disks, body=4053 bytes)

- Lemmings Raw Data Disk 2 (1992) (DMA-Chris White)

### `46f89fa54b2f33f7` (1 disks, body=2040 bytes)

- Lemmings Raw Data Disk 3 (1992) (DMA-Chris White)

### `d9d113eb76c51f63` (1 disks, body=13200 bytes)

- Little Joke - Lords Demo by Lord Insanity (1992) (PD)

### `2f96f4b6cf96f7e9` (1 disks, body=8086 bytes)

- Lovehearts (19xx) (Supplement Software)

### `ec06160c88987e5d` (1 disks, body=32440 bytes)

- MDOS _ MBASIC for Formatting Discs in 2 Drives (19xx)

### `55c64d9534412f7c` (1 disks, body=107325 bytes)

- MNEMOtech Demo 1 (19xx) (PD)

### `99015aeaa1d2f763` (1 disks, body=154362 bytes)

- MNEMOtech Demo 2 (19xx) (PD)

### `2dfdf40f1a72973d` (1 disks, body=5646 bytes)

- Manic Miner_ Splat_ Mr Pac_ Snake Mania_ Craft Compilation (19xx)

### `bacacd3e207ce221` (1 disks, body=8738 bytes)

- Map Print Routines by Chris White (1990)

### `90a45ecc5c12aad7` (1 disks, body=1020 bytes)

- MasterDOS - MasterBASIC BootDisk Creator (19xx)

### `cecd5cdf87336d25` (1 disks, body=2040 bytes)

- MasterDOS - Utility Disk V3.0 (199x)

### `6bf5445940a8a76d` (1 disks, body=9702 bytes)

- MasterDOS File Manager V1.0 (1991)

### `225b2304bdbedad2` (1 disks, body=8086 bytes)

- Mega Demo 6 by Supplement Software (1991) (PD)

### `2341c9493b5dbeeb` (1 disks, body=8086 bytes)

- Mega Text Demo III by Masters of Magic (19xx) (PD)

### `01f788a493ca03bf` (1 disks, body=14009 bytes)

- Megaboot V2.3 (Atom HD Interface) (1999) (M.Groen)

### `3992aaec02aab425` (1 disks, body=8086 bytes)

- Megadisk 1 - Puzzles (19xx)

### `050325c2c98552c8` (1 disks, body=16393 bytes)

- Metempsychosis Unreleased Demo - Adrian_letter2 (19xx)

### `16b08b5a6164003d` (1 disks, body=8086 bytes)

- Metempsychosis Unreleased Demo - Demos4metemp (19xx)

### `4eab21a727463c3b` (1 disks, body=15735 bytes)

- Metempsychosis Unreleased Demo - Internal_digi_utils (19xx)

### `6f083938079ca1ec` (1 disks, body=8086 bytes)

- Metempsychosis Unreleased Demo - Kinetik (19xx)

### `5f3550d1ca8989f0` (1 disks, body=8086 bytes)

- Metempsychosis Unreleased Demo - RTJ_pdm1 (19xx)

### `5793ce5635564711` (1 disks, body=8086 bytes)

- Metempsychosis Unreleased Demo - RTJ_pdm2 (19xx)

### `7b8cbc6085ac4be3` (1 disks, body=18441 bytes)

- Metempsychosis Unreleased Demo - Samples1 (19xx)

### `09a809b2c3f6a881` (1 disks, body=24625 bytes)

- Metempsychosis Unreleased Demo - WIP1 (19xx)

### `a7a83fd615a00232` (1 disks, body=15690 bytes)

- Metempsychosis Unreleased Demo - Wizard (19xx)

### `541ffa1c199769e4` (1 disks, body=8086 bytes)

- Metempsychosis pdm12 (19xx)

### `7588eb488ea3963d` (1 disks, body=8086 bytes)

- Metempsychosis slide_1 (19xx)

### `623d94e58d851687` (1 disks, body=24625 bytes)

- Metempsychosis slide_2 (19xx)

### `56cbdab3f9d1f6fe` (1 disks, body=164250 bytes)

- Metempsychosis term_1b (19xx)

### `6e355ad8d6d322e4` (1 disks, body=510 bytes)

- Mike AJ Demo Disk 1 (1991) _a1_

### `9f7f247d0f209596` (1 disks, body=1530 bytes)

- Mike AJ Demo Disk 7 (1991)

### `c3415c19c2d25f9b` (1 disks, body=15693 bytes)

- Mike AJ Disc 6-Edwin (19xx)

### `e9382a64e59ffd93` (1 disks, body=18659 bytes)

- Mind Games II (199x) (Enigma Variations)

### `09510607cff5eedf` (1 disks, body=4735 bytes)

- Misc. Games 1 (19xx)

### `227f1272dc95b5be` (1 disks, body=24990 bytes)

- Misc. Text Files for Secretary Word Processor (19xx)

### `df40a1eaf134685b` (1 disks, body=14361 bytes)

- Misc. Utilities 1 (19xx)

### `b2a05cc13487f8df` (1 disks, body=38473 bytes)

- Misc. Utilities 2 (19xx)

### `9709bdf03809e0ae` (1 disks, body=16827 bytes)

- Misc. Utilities 3 (19xx)

### `a49f846d7e922479` (1 disks, body=566 bytes)

- Misc. Utilities 4 (19xx)

### `410c7abb0ab561b1` (1 disks, body=2040 bytes)

- Music Disk (19xx) (Unk)

### `364df841599ec69e` (1 disks, body=48213 bytes)

- Music Player 2.1 (19xx)

### `d010d2dd7fc4b45c` (1 disks, body=15868 bytes)

- Neil Holmes_ Boing_ Graphics (1992) (Noesis Software) _a1_

### `d4fc5229fa786b26` (1 disks, body=8086 bytes)

- Network Sigma Issue 6 (Feb-Mar 1996) (Saturn Software)

### `d6db2024f3b55251` (1 disks, body=2040 bytes)

- Night Breed Demo (1990) (PD)

### `397794db3e9032a4` (1 disks, body=2039 bytes)

- Oh No More Lemmings - Data Before Compression (19xx) (Chris White)

### `f2222766ea225e48` (1 disks, body=2040 bytes)

- Oh No More Lemmings - GFX IFF (19xx) (Chris White)

### `edfac77a1256a479` (1 disks, body=1020 bytes)

- Oh No More Lemmings - Make Files (19xx) (Chris White)

### `54788c4e1d5a3d30` (1 disks, body=12239 bytes)

- Oh No More Lemmings - Make Master Side 1_2 (19xx) (Chris White)

### `8213d749e526c8c7` (1 disks, body=11726 bytes)

- Oh No More Lemmings - Raw Data IFFs (19xx) (Chris White)

### `ddc9534d9f14fe4c` (1 disks, body=15800 bytes)

- Open 3D V082 by Tobermory (2001) (PD)

### `f71ecf6155c4c639` (1 disks, body=32300 bytes)

- Ore Warz II (1990) (William McGugan)

### `e1cb1393429f6606` (1 disks, body=8086 bytes)

- Outlet - Sam _ Spectrum Mag Issue 33 (May 1990)

### `cad0f3208a632a36` (1 disks, body=8086 bytes)

- Outlet - Sam _ Spectrum Mag Issue 56 (April 1992)

### `bfd0de70f0212791` (1 disks, body=2331 bytes)

- PAX Disk 2 (1996) (Glenco)

### `6b3819e43403182a` (1 disks, body=1018 bytes)

- PC Suite V2.2 (1991) (Spencer)

### `110e242c901d54f5` (1 disks, body=1530 bytes)

- PDS Boot Disk (1998) (Rob Holman)

### `f3a35085259dbe41` (1 disks, body=8086 bytes)

- PD_90 The Best of 1990 (1990) (PD)

### `2590f1c411796e2d` (1 disks, body=46 bytes)

- Park DBS (19xx)

### `90517122bf2481c5` (1 disks, body=8086 bytes)

- Personal Filing System V2.07 (1994) (Hilton Comp. Services)

### `e76eeeea7692fb07` (1 disks, body=37230 bytes)

- Pickasso_s GFX by Steven Pick (1990) (PD)

### `b0ebdee710caeb5b` (1 disks, body=3674 bytes)

- Pipe Mania_ EFPOTRM_ Klax_ Tetris_ Defenders of Earth (19xx)

### `821df6540285fd9e` (1 disks, body=44370 bytes)

- Porno Disk 1 (19xx) (SAM Psychopaths Incorporated)

### `0c34a31ade2849d0` (1 disks, body=2040 bytes)

- Porno TV (19xx) (PD)

### `ff4eb979d939874d` (1 disks, body=1017 bytes)

- Prince of Persia - Assemble Disk (1990) (Chris White)

### `3a2e943bac812b12` (1 disks, body=3603 bytes)

- Prince of Persia - Assembly _ Pokers (1990) (Chris White)

### `c43e949b399d2fce` (1 disks, body=4825 bytes)

- Prince of Persia - Screens 1 (1990) (Chris White) _a1_

### `d695ae0e9211769f` (1 disks, body=45545 bytes)

- Prince of Persia - Screens 1 (1990) (Chris White)

### `7d5073f9468416e2` (1 disks, body=12240 bytes)

- Prince of Persia - Screens 2 (1990) (Chris White)

### `133c99f187347302` (1 disks, body=11729 bytes)

- Prince of Persia - Show Levels (1990) (Chris White)

### `d07853609211e15e` (1 disks, body=1530 bytes)

- Prince of Persia - Source (1990) (Chris White)

### `098c0f672c7fbe3f` (1 disks, body=4080 bytes)

- Prince of Persia - Trans Disk _ GameDOS (1990) (Chris White)

### `0449b0f60d31acc7` (1 disks, body=1530 bytes)

- Prince of Persia - YS Cover Mount Demo (1990) (Chris White)

### `ea84c2e9ec5c1973` (1 disks, body=1020 bytes)

- Prince of Persia Demo (1990) (Revelation-Chris White) _a1_

### `d9feb1ddc0e17771` (1 disks, body=8086 bytes)

- Prince of Persia Demo (1990) (Revelation-Chris White)

### `59c63c8de8e3edca` (1 disks, body=2040 bytes)

- Print Routines_Tables_Testers by Chris White (19xx)

### `4ef8faa15ddf815c` (1 disks, body=8086 bytes)

- Printer Port Music Sample Player (19xx)

### `a09dfe6fd32bddcc` (1 disks, body=8086 bytes)

- Public the Third (19xx)

### `10c694d747b12705` (1 disks, body=17358 bytes)

- Raytrace Disk 1 (1994) (Colonysoft)

### `8e9b529fa294ea21` (1 disks, body=8086 bytes)

- Robocop - Rolling Demo (19xx) (PD)

### `b3c3ee71cc95d3a5` (1 disks, body=4080 bytes)

- SC DTP (1991) (Steve_s Software)

### `5ae7ecf1649b8972` (1 disks, body=3183 bytes)

- SC PD 1 - Speclone_Compressor by Steve_s Software (1991) (PD)

### `200278d35e1e7d4b` (1 disks, body=46 bytes)

- Sam Adventure System V1.0 (1992) (Axxent Software)

### `6c7cb5197d3034af` (1 disks, body=8086 bytes)

- Sam Amateur Programming _ Electronics Issue 1 (Feb 1992)

### `3d231309c58b9c7b` (1 disks, body=8086 bytes)

- Sam Amateur Programming _ Electronics Issue 2 (Mar 1992)

### `2ab28418f3673bef` (1 disks, body=8086 bytes)

- Sam Amateur Programming _ Electronics Issue 4 (May 1992)

### `08317b88daef6576` (1 disks, body=8086 bytes)

- Sam Amateur Programming _ Electronics Issue 5 (Jun 1992)

### `bb17425f57b905f7` (1 disks, body=8086 bytes)

- Sam Amateur Programming _ Electronics Issue 6 (Aug 1992)

### `041b2d1f1dc791cd` (1 disks, body=165 bytes)

- Sam BASIC Demo Programs by Chris White (1991) (PD)

### `42c388a53f73ca44` (1 disks, body=8086 bytes)

- Sam Cards (1994) (Supplement Software)

### `8fb01afeca36dbcd` (1 disks, body=665 bytes)

- Sam Coupe Demo Collection Disk 01 (1992) (PD)

### `9f32255832a86dd5` (1 disks, body=34435 bytes)

- Sam Coupe Demo Collection Disk 02 (1992) (PD)

### `153b4c1b750833e1` (1 disks, body=2040 bytes)

- Sam Coupe Demo Collection Disk 03 (1992) (PD)

### `59a2a3377c71df6d` (1 disks, body=54118 bytes)

- Sam Coupe Demo Collection Disk 04 (1992) (PD)

### `78345b4235aae4dd` (1 disks, body=4735 bytes)

- Sam Coupe Demo Collection Disk 06 (1992) (PD)

### `8969e9a30b1e5863` (1 disks, body=53550 bytes)

- Sam Coupe Demo Collection Disk 07 (1992) (PD)

### `3db3af4d4a245344` (1 disks, body=18360 bytes)

- Sam Coupe Demo Collection Disk 08 (1992) (PD)

### `f9d0ed730f76838b` (1 disks, body=15695 bytes)

- Sam D I C E V1.0 for MasterDOS (1991) (Kobrahsoft)

### `4fff977d8013377f` (1 disks, body=73575 bytes)

- Sam Mines (19xx) (PD)

### `e9dd512c58f870b4` (1 disks, body=8086 bytes)

- Sam Paint (19xx)

### `1e009491ead9a1e9` (1 disks, body=8086 bytes)

- Sam Prime (19xx)

### `23fd8a30002c7978` (1 disks, body=25735 bytes)

- Sam Public Quarterly Issue 5 (19xx)

### `111177e603253ae9` (1 disks, body=255 bytes)

- Sam Supplement Magazine Issue 01 (Sep 1990)

### `1e6670db1775b582` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 02 (Nov 1990)

### `e850f05bd65f5a12` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 03 (Dec 1990)

### `632db00d73f2bcdb` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 04 (Jan 1991)

### `b3e093d53d3e8428` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 05 (Feb 1991)

### `74b0fcff076393f9` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 07 (Apr 1991)

### `9813a2b223cd3390` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 08 (May 1991)

### `9548a3aa5d7098f4` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 09 (Jun 1991)

### `4e0b7599aa852830` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 10 (Jul 1991)

### `7722a00047e84140` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 11 (Aug 1991)

### `f70d22ec2642fb03` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 12B Freeware (Sep 1991)

### `a465e0b79fb55157` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 13 (Oct 1991)

### `999a129332b2916b` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 14 (Nov 1991)

### `cf06e3c38e44072d` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 15 (Dec 1991)

### `8ba62ef306bf80d5` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 16 (Jan 1992)

### `b33a4a2558f9125d` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 17 (Feb 1992)

### `03182e8ae606f607` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 18 (Mar 1992)

### `d09a2f4f45505f97` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 19 (Apr 1992)

### `6c8ba91597ed8c3c` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 22 (Jul 1992)

### `b5888b0cd6e79794` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 24 (Sep 1992)

### `9d6f3291cdc4f938` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 27 (Dec 1992)

### `06f000fe436a5a67` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 28 (Jan 1993)

### `d80749fad2c5ec95` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 32 (May 1993)

### `0920b0f727f516cf` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 35 (Aug 1993)

### `adbfcbfe401f73ff` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 39 (Dec 1993)

### `eef4d45e68ad2765` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 40 (Jan 1994)

### `ae4c23bb6344861d` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 43 (Apr 1994)

### `3893517bfe602537` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 47 (Aug 1994)

### `d64a50beb45e4415` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 48 (Sep 1994)

### `e064fcc49fe12839` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 49 (Oct 1994)

### `d0d935496574bcb4` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 51 (Dec 1994)

### `30facf34e3ae8416` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 52 (Jan 1995)

### `85fe103f9df025d2` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 53 (Feb 1995)

### `92514ebc590f5d1d` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Issue 54 (Mar 1995)

### `0b8b4046fa7f5d3a` (1 disks, body=8086 bytes)

- Sam Supplement Magazine Xmas Disk (Dec 1992)

### `6d56e93be9d356ba` (1 disks, body=8086 bytes)

- Sam Utils (1993)

### `d36b77540f6f7237` (1 disks, body=8086 bytes)

- Sam X (19xx) (Supplement Software)

### `2a8bb90a7cce3e6c` (1 disks, body=12240 bytes)

- SamCo Birthday Pack Best of Enceledus (1991)

### `e76cba1ae4eee4da` (1 disks, body=12237 bytes)

- SamCo Birthday Pack Demo and Previews (1991)

### `a020a583ec5c5a57` (1 disks, body=2009 bytes)

- SamCo Birthday Pack Games and Utils (1991) (Revelation) _a1_

### `9424150ac355dcce` (1 disks, body=2040 bytes)

- SamCo News Disk - Sam Juggler (1991)

### `30be24dbcfcae6a3` (1 disks, body=10698 bytes)

- SamCo News Disk 1 (Jan 1992)

### `aff6bbf5bcc5ad46` (1 disks, body=4080 bytes)

- SamCo News Disk 2 (Feb 1992)

### `fde198d34bd1af0f` (1 disks, body=5289 bytes)

- SamCo News Disk 3 (Mar 1992)

### `71c9b9601b80bf06` (1 disks, body=2330 bytes)

- SamCo News Disk 4 (1992)

### `a9ff3d1e798d8cfc` (1 disks, body=4590 bytes)

- SamCo News Disk 5 (1992)

### `9919ef4aac9605c2` (1 disks, body=9801 bytes)

- Samsational Complete Guide to SAM PD Software (1992) (SCPDSA)

### `1d4d39218c139c85` (1 disks, body=8086 bytes)

- Sandman_s Shadow 1 (1993) (PD)

### `93ae595946cbfc07` (1 disks, body=8086 bytes)

- Sandman_s Shadow 2 (1993) (PD)

### `1a0266804ed4c785` (1 disks, body=8086 bytes)

- Sandman_s Shadow 3 (1993) (PD)

### `cffb14eb8520690f` (1 disks, body=8086 bytes)

- Satellite_92 (1992)

### `d01f9b12d501b19e` (1 disks, body=24350 bytes)

- Sci-Fi Screens Disk 1 (1994) (ColonySoft)

### `912af1d501dd8de0` (1 disks, body=30235 bytes)

- Screens Shots 1 (1991)

### `23d231cec05f234e` (1 disks, body=30296 bytes)

- Screens Shots 2 (1991)

### `f067bee72231d802` (1 disks, body=30296 bytes)

- Screens Shots 3 (1991)

### `263f716215abf909` (1 disks, body=36865 bytes)

- Screens Shots 4 (1991)

### `f7d93f36612a7685` (1 disks, body=3570 bytes)

- Screens Shots 5 (1991)

### `eface0f0b955164b` (1 disks, body=1085 bytes)

- Secretary Word Processor_ The (1992) (A.N.Stevens)

### `22b364ea5beabbed` (1 disks, body=132 bytes)

- Silly Demo 1 by Lord Insanity (1990) (PD) _a1_

### `ba53f23f387c40f0` (1 disks, body=2040 bytes)

- Simon Cooke - Fast E-Tracker Player for Lemmings (1992)

### `4371857cb647d5c3` (1 disks, body=3761 bytes)

- Simon Cooke - Fred 18 Freebie (1991)

### `1b08f26b7abce60b` (1 disks, body=1527 bytes)

- Simon Cooke - Graphic Converters (1993)

### `2de501366d3e37f8` (1 disks, body=2040 bytes)

- Simon Cooke - Holidaying At SamTech (1991)

### `4ffe30c8b8bcd872` (1 disks, body=1020 bytes)

- Simon Cooke - Loading Tetris (19xx)

### `f8f3afcc1330ebfe` (1 disks, body=2039 bytes)

- Simon Cooke - Rick Dangerous Stuff (1990)

### `be9f4198d9851f5f` (1 disks, body=510 bytes)

- Simon Cooke - Samples (1991)

### `b05585902262df5b` (1 disks, body=60 bytes)

- Simon Cooke - Stuff (1992)

### `08c2c5aa041946d5` (1 disks, body=1039 bytes)

- Simon Cooke - Work Disk (1992)

### `08cc627ae4ea1d14` (1 disks, body=8093 bytes)

- Small C - Compiler V3.1 (19xx) (Rumsoft)

### `9ed339a39a563d41` (1 disks, body=1530 bytes)

- Sound Machine (1991) (Paul Angel)

### `fbd3854b912a449e` (1 disks, body=39330 bytes)

- Speccy EmulFiles (1991) (Chris White)

### `721551247cd4f9ca` (1 disks, body=2570 bytes)

- Speccy Emulators (1991) (Chris White)

### `1670a7d17b0132c1` (1 disks, body=2040 bytes)

- Speccy Utils (1991) (Chris White)

### `6adf935424f52438` (1 disks, body=8086 bytes)

- Spectrum 128 - Myth and Escape From Singe_s Castle (19xx)

### `2ea0d9173dddd47f` (1 disks, body=7909 bytes)

- Spectrum 128 Music Disk 1 (19xx) (PD)

### `51b1c68aeddaf053` (1 disks, body=9181 bytes)

- Spectrum Emulator (Sept 04) (1990)

### `d1b70a1445b0372f` (1 disks, body=4543 bytes)

- Spectrum Games (128K) Disk 01 (1991)

### `ca4f5f99c00d0f7b` (1 disks, body=51662 bytes)

- Spectrum Games (128K) Disk 02 (1991)

### `376164478e5c066b` (1 disks, body=6992 bytes)

- Spectrum Games (128K) Disk 03 (1991)

### `bbd08e667a089486` (1 disks, body=169 bytes)

- Spectrum Games (128K) Disk 04 (1991)

### `17578664befe3147` (1 disks, body=30618 bytes)

- Spectrum Games (128K) Disk 05 (1991)

### `0c67c69e9f94e664` (1 disks, body=16644 bytes)

- Spectrum Games (128K) Disk 06 (1991)

### `a851cb71b60b313b` (1 disks, body=11154 bytes)

- Spectrum Games (128K) Disk 07 (1991)

### `ef3fda2cfc8f57d9` (1 disks, body=29790 bytes)

- Spectrum Games (128K) Disk 08 (1991)

### `3a3eb37d99089ed4` (1 disks, body=22184 bytes)

- Spectrum Games (128K) Disk 09 (1991)

### `4b4590a60fbcbfc5` (1 disks, body=5185 bytes)

- Spectrum Games (128K) Disk 10 (1991)

### `845ff405604811e5` (1 disks, body=53513 bytes)

- Spectrum Games (128K) Disk 11 (1991)

### `c8a6c7c087a2dc4f` (1 disks, body=227 bytes)

- Spectrum Games (128K) Disk 12 (1991)

### `4bdb9f49e49834a2` (1 disks, body=74142 bytes)

- Spectrum Games (48K) Disk 1 (1991)

### `6589bc32b3c2f92e` (1 disks, body=74141 bytes)

- Spectrum Games (48K) Disk 2 (1991)

### `9756db858428a533` (1 disks, body=74142 bytes)

- Spectrum Games (48K) Disk 3 (1991)

### `71660009a2b74579` (1 disks, body=74141 bytes)

- Spectrum Games (48K) Disk 4 (1991)

### `d36ed4a4ddfb1341` (1 disks, body=18360 bytes)

- Spectrum Games (48K) Disk 5 (1991)

### `7cf97cba81af3631` (1 disks, body=12240 bytes)

- Spectrum Games (48K) Disk 6 (1991)

### `2a701ee27032f6f5` (1 disks, body=15690 bytes)

- Spectrum Games Compilation 12 (1992)

### `b4aca8aff4efdffc` (1 disks, body=8086 bytes)

- Spectrum Games Compilation 13 (1992)

### `03cd50b5e50c5bb9` (1 disks, body=8018 bytes)

- Spectrum _ SAM Computing Issue 3 (19xx)

### `5fb3e7460965c75c` (1 disks, body=12240 bytes)

- Speed King Hacking (1991) (Chris White)

### `425bce275aec6a70` (1 disks, body=333 bytes)

- Star Demo (19xx) (PD)

### `4051f72e2146bfe9` (1 disks, body=24621 bytes)

- Star Wars Slideshow 1 (19xx) (PD)

### `8e7932233ceecf38` (1 disks, body=1425 bytes)

- Steffan Drisson_s Multiple Stuff (1991)

### `4b3a6b20f3184ede` (1 disks, body=2461 bytes)

- Stuart_s Leonardi_s Rotater Stuff (1990)

### `eda9270d197b7064` (1 disks, body=20910 bytes)

- Stuart_s Leonardi_s Vector Stuff (1990)

### `9acde11c04c4dfa8` (1 disks, body=12230 bytes)

- Top Gun 1Mb MDOS Demo Disk 1 (1991) (PD) _a1_

### `2d94fc4c29675c4d` (1 disks, body=3700 bytes)

- Top Gun 1Mb MDOS Demo Disk 1 (1991) (PD)

### `651ceca0418ff661` (1 disks, body=4080 bytes)

- Top Gun 1Mb MDOS Demo Disk 2 (1991) (PD) _a1_

### `204baca3dad8428c` (1 disks, body=122910 bytes)

- Top Gun 1Mb MDOS Demo Disk 2 (1991) (PD)

### `c6dd90280914b377` (1 disks, body=15692 bytes)

- TurboMON V1.0 (19xx)

### `69a211a4de1747c7` (1 disks, body=521 bytes)

- Various E-Tracker Stuff (1991) (PD)

### `ed45d7b6dea661f3` (1 disks, body=37173 bytes)

- Various Unused GFX (1991) (Evolution)

### `d75d94c539bfd8f2` (1 disks, body=8086 bytes)

- Visually 1 (19xx) (Zenith Graphics)

### `32ef028957538851` (1 disks, body=8086 bytes)

- Visually 4 (19xx) (Zenith Graphics)

### `e47caa45f77e0f3b` (1 disks, body=8086 bytes)

- Visually 6 (19xx) (Zenith Graphics)

### `15a8e15fd680dbf3` (1 disks, body=4161 bytes)

- Walker 1 Meg Demo (1991) (PD)

### `ed635e2282aec119` (1 disks, body=20198 bytes)

- Wobbles Painful History_ The (2001) (The Wobbles)

### `90c5cca7472ca621` (1 disks, body=631 bytes)

- ZUB Demo by Simon Cooke (1992) (PD)

### `0df25ed6fe3f0156` (1 disks, body=9319 bytes)

- Zenith Edition 1 (19xx) (Zenith Graphics)

### `473bdd91856e0a3d` (1 disks, body=8086 bytes)

- Zenith Edition 2 (19xx) (Zenith Graphics)

### `d373434c5da09096` (1 disks, body=8086 bytes)

- Zenith Edition 2-5 (19xx) (Zenith Graphics)

### `6e691c7af719521f` (1 disks, body=4639 bytes)

- newdisk

### `923d7c33e6239ab1` (1 disks, body=8086 bytes)

- pete-made

### `fdc9c5194995c4d9` (1 disks, body=67438 bytes)

- test__1

### `1e9cec8cf567b429` (1 disks, body=11053 bytes)

- trinity

### `e59feafa029fbe18` (1 disks, body=8086 bytes)

- trinload

## Disks with no slot-0 file (15)

Non-bootable archive disks. Dir entries were written by
*some* DOS but the DOS isn't on the disk.

Sample (first 20):
- Arnie_s Samples (1991)
- Images SAM Software - Pictures _ Demos (19xx)
- Integrated Logic_s Demo Disk and Utils (1990) (PD)
- Lyra 3 Megademo by ESI (1993) (PD) _a1_
- Lyra 3 Megademo by ESI (1993) (PD) _a2_
- Lyra 3 Megademo by ESI (1993) (PD)
- Mouse Driver V2.0 by Steve Taylor (19xx) (Sam PD Sware Lib.) (PD)
- PC2SAM Disk 6 (1994) (ColonySoft)
- SAMart and Slideshow (19xx) (Sam PD Sware Lib.) (PD)
- Sam Adventure System V1.2 (1992) (Axxent Software)
- SamCo-Tech PC Contracts (1989)
- Samples and Instruments (1991) (PD)
- Sega Graphic Converters (1992) (Chris White)
- Source Code Samples Disk 1 (1990) (Chris White)
- Source Code Samples Disk 2 (1990) (Chris White)

