# DOS family `9bc0fb4b949109e8` — samdos2

Self-contained materialisation of one DOS family from the SAM
Coupé corpus. The family is the equivalence class of slot-0 DOS
bodies clustered at 1.5% byte-diff (see
`docs/dos-families.md` for the full table).

## Identity

- **Family-head SHA16:** `9bc0fb4b949109e8`
- **Variants in family:** 126
- **Disks in family:** 291
- **Body length(s):** 10000
- **Load address(es):** 0x008000 (p3), 0x008030 (p3), 0x00e0be (p4), 0x038009 (p15), 0x078009 (p31), 0x080009 (p33)
- **Execution address(es):** 0x000000

## Files in this directory

- `body.bin` — exact slot-0 body of the family-head SHA
  (`9bc0fb4b949109e8`). Header-decoded, so byte 0 is the first byte the
  ROM would copy to the body's load address.
- `body.hex` — xxd-style hex dump of `body.bin` (big families only).
- `variants/*.md` — byte-diff summary of each variant against the
  head, including the differing byte ranges in hex so the variant
  body can be reconstructed from head + diff. Only written for
  families with at least 5 disks; extract any variant body with
  `tools/audit/extract_dos.py <sha>`.
- `src/` — commented original assembly source (only present when
  upstream source is known to assemble to a binary in this family).

## Source binding

- **Binary-of-record SHA16:** `3cca541beb3f9fe9` (full SHA `3cca541beb3f9fe93402a770997945b2be852e69f278d2b176ba0bbc4fbb6077`)
- **Upstream source:** `/Users/pmoore/git/samdos` (src/*.s)
- **Reference binary in upstream:** `/Users/pmoore/git/samdos/res/samdos2.reference.bin`

- **Upstream archive:** https://ftp.nvg.ntnu.no/pub/sam-coupe/sources/SamDos2InCometFormatMasterv1.2.zip

### Notes from upstream README

> Source from https://github.com/stefandrissen/samdos (Stefan
> Drissen). HEAD is samdos2, assembled byte-identical to
> variants/3cca541beb3f9fe9.bin. The git history of that repo
> carries the five upstream comp1..comp5 versions as separate
> commits.
>
> The upstream archive `SamDos2InCometFormatMasterv1.2.zip`
> contains a SAM .dsk image (also fetched into
> `upstream/SamDos2InCometFormatMasterv1.2.dsk` if available)
> with all 28 source files (a1.s..h2.s, comp1.s..comp5.s,
> gm1.s, gm2.s, ldit1.s..ldit3.s, ld1.s).

## Variants in this family

| Variant SHA16 | Disks | Length | Load | Exec | Within-fam diff vs head |
|---|---:|---:|---|---|---:|
| `3cca541beb3f9fe9` ←source-of-record | 30 | 10000 | 0x078009 (p31) | — | 0.640% |
| `9bc0fb4b949109e8` | 19 | 10000 | 0x038009 (p15) | — | head |
| `a242cc7b48a65541` | 16 | 10000 | 0x080009 (p33) | — | 0.390% |
| `fb0b95fdf1091299` | 15 | 10000 | 0x078009 (p31) | — | 0.200% |
| `739a6586e599ce10` | 13 | 10000 | 0x078009 (p31) | — | 0.550% |
| `2a3e618918a5520b` | 11 | 10000 | 0x038009 (p15) | — | 0.160% |
| `31771043e0729972` | 10 | 10000 | 0x00e0be (p4) | — | 1.090% |
| `4c639e1325d9940f` | 10 | 10000 | 0x078009 (p31) | — | 0.710% |
| `1293e6ab6739adb3` | 8 | 10000 | 0x078009 (p31) | — | 0.490% |
| `14e4d256a8f0f0ec` | 6 | 10000 | 0x038009 (p15) | — | 0.150% |
| `d46ac18e935ff556` | 6 | 10000 | 0x078009 (p31) | — | 0.180% |
| `5f76f756fd1e6bee` | 5 | 10000 | 0x078009 (p31) | — | 0.200% |
| `3273b9e8ec9b1ecd` | 4 | 10000 | 0x038009 (p15) | — | 0.190% |
| `c753cf32520d0060` | 3 | 10000 | 0x078009 (p31) | — | 0.130% |
| `3a4be5a0b749094b` | 3 | 10000 | 0x078009 (p31) | — | 1.130% |
| `d4d0e56b1342da0f` | 3 | 10000 | 0x038009 (p15) | — | 0.040% |
| `01a3ab44786698d7` | 3 | 10000 | 0x038009 (p15) | — | 0.310% |
| `61f76732d61326d7` | 3 | 10000 | 0x078009 (p31) | — | 0.680% |
| `401e21508fe023a5` | 3 | 10000 | 0x038009 (p15) | — | 0.680% |
| `cf4c3b28ab2d0144` | 3 | 10000 | 0x038009 (p15) | — | 0.160% |
| `59210e67e270d468` | 2 | 10000 | 0x078009 (p31) | — | 0.550% |
| `cf248342d24596e3` | 2 | 10000 | 0x038009 (p15) | — | 0.770% |
| `b112c7d5d4340027` | 2 | 10000 | 0x038009 (p15) | — | 0.700% |
| `0415ede5ca82e665` | 2 | 10000 | 0x038009 (p15) | — | 0.550% |
| `58a02b8abd11d55b` | 2 | 10000 | 0x008030 (p3) | — | 0.180% |
| `df1ea02fd28e8997` | 2 | 10000 | 0x038009 (p15) | — | 0.460% |
| `6c53f7bb9dd34172` | 2 | 10000 | 0x038009 (p15) | — | 0.440% |
| `f2199b85e766e967` | 2 | 10000 | 0x078009 (p31) | — | 0.690% |
| `440f66f1aab2a9f3` | 2 | 10000 | 0x078009 (p31) | — | 0.670% |
| `6f0c8686c2eafafe` | 2 | 10000 | 0x078009 (p31) | — | 0.630% |
| `c3c34de6552f864a` | 2 | 10000 | 0x078009 (p31) | — | 0.650% |
| `daeae9e0daf098a9` | 1 | 10000 | 0x078009 (p31) | — | 0.620% |
| `7089ca659242dde6` | 1 | 10000 | 0x078009 (p31) | — | 0.540% |
| `15ecc4085019ca35` | 1 | 10000 | 0x078009 (p31) | — | 0.270% |
| `557b1a4bc1e3e723` | 1 | 10000 | 0x078009 (p31) | — | 0.270% |
| `80b973784d34c635` | 1 | 10000 | 0x078009 (p31) | — | 0.420% |
| `46644ee5c8ceb328` | 1 | 10000 | 0x078009 (p31) | — | 0.200% |
| `17f417b9a7ca6444` | 1 | 10000 | 0x078009 (p31) | — | 0.190% |
| `a7a17dd5b77a5555` | 1 | 10000 | 0x078009 (p31) | — | 0.330% |
| `c85f4d0eb8be8c3a` | 1 | 10000 | 0x078009 (p31) | — | 1.170% |
| `a49f2fd382925282` | 1 | 10000 | 0x078009 (p31) | — | 0.390% |
| `05017e20003819a7` | 1 | 10000 | 0x078009 (p31) | — | 0.560% |
| `1b853b911475c903` | 1 | 10000 | 0x038009 (p15) | — | 0.090% |
| `d852f669d7278df8` | 1 | 10000 | 0x078009 (p31) | — | 0.740% |
| `eb1f8c0158856a40` | 1 | 10000 | 0x038009 (p15) | — | 0.140% |
| `a57d9139b31938c8` | 1 | 10000 | 0x038009 (p15) | — | 0.090% |
| `f6f3c9f45cab0b92` | 1 | 10000 | 0x038009 (p15) | — | 0.550% |
| `801951340a516cac` | 1 | 10000 | 0x038009 (p15) | — | 0.360% |
| `5a93056404830c42` | 1 | 10000 | 0x038009 (p15) | — | 0.220% |
| `69c86a8b6360d876` | 1 | 10000 | 0x078009 (p31) | — | 0.200% |
| `5c389f3cfa578eaa` | 1 | 10000 | 0x038009 (p15) | — | 0.620% |
| `064224e791163227` | 1 | 10000 | 0x078009 (p31) | — | 0.610% |
| `1b6ab6536f4e23e2` | 1 | 10000 | 0x038009 (p15) | — | 0.100% |
| `5c55e456b89a177c` | 1 | 10000 | 0x078009 (p31) | — | 0.660% |
| `d481994aba88a063` | 1 | 10000 | 0x078009 (p31) | — | 0.270% |
| `0b232c8d426c1534` | 1 | 10000 | 0x078009 (p31) | — | 0.540% |
| `ff9990ed35284292` | 1 | 10000 | 0x078009 (p31) | — | 0.550% |
| `4c400685ae82a2ab` | 1 | 10000 | 0x078009 (p31) | — | 0.490% |
| `78093dc9050f1fe4` | 1 | 10000 | 0x078009 (p31) | — | 0.500% |
| `ba47e62e42fa6c83` | 1 | 10000 | 0x078009 (p31) | — | 0.200% |
| `a041404837bb1f94` | 1 | 10000 | 0x078009 (p31) | — | 0.390% |
| `891bc46720d223f0` | 1 | 10000 | 0x078009 (p31) | — | 0.170% |
| `8ac2f0dc49b4ffe6` | 1 | 10000 | 0x078009 (p31) | — | 0.380% |
| `98d52753a0ca13cf` | 1 | 10000 | 0x038009 (p15) | — | 0.690% |
| `8a81ca8c8271f5ab` | 1 | 10000 | 0x038009 (p15) | — | 0.380% |
| `b366bf0c8841fbc3` | 1 | 10000 | 0x078009 (p31) | — | 0.550% |
| `080e2ab3c97bc23c` | 1 | 10000 | 0x078009 (p31) | — | 0.440% |
| `6a6feb5c4a7b1d91` | 1 | 10000 | 0x038009 (p15) | — | 0.270% |
| `abfdf5499cae5497` | 1 | 10000 | 0x078009 (p31) | — | 0.640% |
| `298b3e7150d3b3be` | 1 | 10000 | 0x078009 (p31) | — | 0.120% |
| `3be4384702b9054a` | 1 | 10000 | 0x078009 (p31) | — | 0.640% |
| `49aa9a49b0964aba` | 1 | 10000 | 0x078009 (p31) | — | 0.620% |
| `864506f419c442e3` | 1 | 10000 | 0x078009 (p31) | — | 0.610% |
| `96e56fea822f33ab` | 1 | 10000 | 0x078009 (p31) | — | 0.920% |
| `7587b034dcbcb1b7` | 1 | 10000 | 0x078009 (p31) | — | 0.710% |
| `6f82050a2a0a43d3` | 1 | 10000 | 0x078009 (p31) | — | 0.680% |
| `b4fe3be82ac7a7b0` | 1 | 10000 | 0x078009 (p31) | — | 0.190% |
| `222cf6e7c0f515d8` | 1 | 10000 | 0x078009 (p31) | — | 0.700% |
| `82e592ca1b61476a` | 1 | 10000 | 0x078009 (p31) | — | 0.580% |
| `6fc1cdddbab20aa1` | 1 | 10000 | 0x078009 (p31) | — | 0.650% |
| `167956b902b73b6a` | 1 | 10000 | 0x078009 (p31) | — | 0.590% |
| `4eb0176725953059` | 1 | 10000 | 0x078009 (p31) | — | 0.410% |
| `703faf63d1cd4532` | 1 | 10000 | 0x078009 (p31) | — | 0.530% |
| `1e2e616c034ca8b5` | 1 | 10000 | 0x078009 (p31) | — | 0.660% |
| `522789453c5843a2` | 1 | 10000 | 0x078009 (p31) | — | 0.590% |
| `1d6f023e8820be49` | 1 | 10000 | 0x078009 (p31) | — | 0.700% |
| `da0c1ba197cd8010` | 1 | 10000 | 0x078009 (p31) | — | 0.710% |
| `bf76f513f5ff657c` | 1 | 10000 | 0x078009 (p31) | — | 0.580% |
| `81669e1f07422cb0` | 1 | 10000 | 0x078009 (p31) | — | 0.520% |
| `fa82a7e1450378c3` | 1 | 10000 | 0x078009 (p31) | — | 0.580% |
| `1aac98174ed8a75c` | 1 | 10000 | 0x078009 (p31) | — | 0.500% |
| `6100659ce98702b2` | 1 | 10000 | 0x078009 (p31) | — | 0.300% |
| `8e148d8cc88349bc` | 1 | 10000 | 0x078009 (p31) | — | 0.410% |
| `f22c4afc0cb24491` | 1 | 10000 | 0x078009 (p31) | — | 0.540% |
| `9cd1813018d2687c` | 1 | 10000 | 0x078009 (p31) | — | 0.460% |
| `8544bcebf58dfc8f` | 1 | 10000 | 0x078009 (p31) | — | 0.530% |
| `792d52b0c2ef5e61` | 1 | 10000 | 0x078009 (p31) | — | 0.550% |
| `173055c1d00f401e` | 1 | 10000 | 0x078009 (p31) | — | 0.540% |
| `d6152674405d34eb` | 1 | 10000 | 0x078009 (p31) | — | 0.460% |
| `41be122a791eae30` | 1 | 10000 | 0x078009 (p31) | — | 0.600% |
| `87f1a1f97596fa74` | 1 | 10000 | 0x078009 (p31) | — | 1.020% |
| `fd07f97df85dc17a` | 1 | 10000 | 0x078009 (p31) | — | 0.760% |
| `e39224f96381c658` | 1 | 10000 | 0x078009 (p31) | — | 0.540% |
| `e7312179003c9fe2` | 1 | 10000 | 0x078009 (p31) | — | 0.690% |
| `03f24fabd55fc7d5` | 1 | 10000 | 0x078009 (p31) | — | 0.610% |
| `da9251cfcaa28773` | 1 | 10000 | 0x078009 (p31) | — | 0.460% |
| `98785744cd87ccc2` | 1 | 10000 | 0x078009 (p31) | — | 0.560% |
| `f017586eb9b3272d` | 1 | 10000 | 0x078009 (p31) | — | 0.590% |
| `612826e7b8ebc753` | 1 | 10000 | 0x078009 (p31) | — | 0.520% |
| `999499e4b78c2e0c` | 1 | 10000 | 0x078009 (p31) | — | 0.610% |
| `13334f4bf5b5fa2c` | 1 | 10000 | 0x078009 (p31) | — | 0.430% |
| `3ca200dbb6588f8b` | 1 | 10000 | 0x078009 (p31) | — | 0.320% |
| `902ff6d466ae5e14` | 1 | 10000 | 0x078009 (p31) | — | 0.750% |
| `d031665e45db3c86` | 1 | 10000 | 0x078009 (p31) | — | 0.240% |
| `b9b42cf1534fc63d` | 1 | 10000 | 0x078009 (p31) | — | 0.150% |
| `2a375be311585fba` | 1 | 10000 | 0x078009 (p31) | — | 0.280% |
| `3663b9b5c46b3394` | 1 | 10000 | 0x078009 (p31) | — | 0.330% |
| `4b9daee4e076df6d` | 1 | 10000 | 0x078009 (p31) | — | 0.320% |
| `f5b9bb560f6a0d3b` | 1 | 10000 | 0x038009 (p15) | — | 0.590% |
| `e88ff2a2e71e8516` | 1 | 10000 | 0x038009 (p15) | — | 0.420% |
| `214159b9cffa2359` | 1 | 10000 | 0x038009 (p15) | — | 0.170% |
| `5f7c9e75fbb60438` | 1 | 10000 | 0x038009 (p15) | — | 0.550% |
| `7b713de1bc50a02f` | 1 | 10000 | 0x038009 (p15) | — | 0.180% |
| `dedf6e1bee3729d4` | 1 | 10000 | 0x078009 (p31) | — | 0.420% |
| `774e500863e77368` | 1 | 10000 | 0x078009 (p31) | — | 0.290% |
| `a99e1f2160ecbd59` | 1 | 10000 | 0x008000 (p3) | 0x000000 | 0.660% |

## Sample disks

- 18 Rated Poker for 512k (19xx) (Supplement Software)
- 32 Colour Demo by Gordon Wallis (1992) (PD)
- Adventures of Captain Comic_ The (19xx) (Lars)
- Aliens vs Predator Demo by Gordon Wallis (1991) (PD)
- Allan Stevens - 50 Programs to Play and Write (19xx)
- Allan Stevens - Capricorn Software Disk 1 (1994)
- Allan Stevens - Capricorn Software Disk 2 (1994)
- Allan Stevens - Capricorn Software Disk 3 Unfinished (1994)
- Allan Stevens - Colour Cycle (19xx)
- Allan Stevens - Home Utilities (1994)

... and 281 more.
