# DOS catalog

Empirical survey of the SAM Coupé corpus at `~/sam-corpus/disks/`,
fingerprinting each disk's DOS by following the actual ROM ↔ DOS
load contract (rom-disasm:20473-20598, the BOOTEX routine).

## ROM bootability gate

Before recording any fingerprint, the survey applies ROM BOOTEX's
own bootability check (BTCK loop at rom-disasm:20586-20596):
T4S1 bytes 256..259 are XORed with the literal `"BOOT"` and
AND-masked with 0x5F. If any of the four bytes mismatches, ROM
prints `NO DOS` and refuses to boot — the disk's T4S1 and slot 0
are irrelevant for DOS classification, so the survey skips them.

## ROM-contract fingerprint (bootable disks only)

ROM's BOOT path reads **exactly track-4-sector-1 (512 bytes)** into
`0x8000..0x81FF` (rom-disasm:20550, `LD HL,0x8000` + `INI` loop at
RSA5), validates the 4-byte signature, then jumps to `0x8009`.
**Nothing else** is part of the ROM contract — the bootstrap at
0x8009 is per-DOS code that chooses where to find the rest of
itself.

Each disk's strict ROM-contract fingerprint is therefore:

- **Length:** 512 bytes (constant)
- **Load address:** 0x8000 (constant)
- **Execution address:** 0x8009 (constant)
- **Content SHA-256:** varies per disk (see table)

- Total disks scanned: **800**
- Disks ROM would refuse (`NO DOS`): **237**
- Disks ROM would boot: **563**
- Unique T4S1 contents (among bootable): **184**
- Disks with malformed images: **0**

### Unique T4S1 contents (most-used first)

| SHA-256 (16) | Disks | Length | Load | Exec |
|---|---:|---:|---:|---:|
| `9b020cca69c03dcb` | 107 | 512 | 0x8000 | 0x8009 |
| `8cbf7dc090f81166` | 31 | 512 | 0x8000 | 0x8009 |
| `1bd8f6e214482d77` | 30 | 512 | 0x8000 | 0x8009 |
| `96f8319f17594aac` | 26 | 512 | 0x8000 | 0x8009 |
| `983d8d27cd374a17` | 23 | 512 | 0x8000 | 0x8009 |
| `ef32fb0142103c74` | 19 | 512 | 0x8000 | 0x8009 |
| `18d4a709b74df92d` | 15 | 512 | 0x8000 | 0x8009 |
| `842fb716a9e6379a` | 13 | 512 | 0x8000 | 0x8009 |
| `2486448f3e41ac4d` | 11 | 512 | 0x8000 | 0x8009 |
| `2fa94f89c30c6bcd` | 10 | 512 | 0x8000 | 0x8009 |
| `e01b39fb0c72780c` | 10 | 512 | 0x8000 | 0x8009 |
| `81a555670453a6dd` | 8 | 512 | 0x8000 | 0x8009 |
| `d7581afadfc8f8e4` | 8 | 512 | 0x8000 | 0x8009 |
| `87ce8e53bc14cd97` | 8 | 512 | 0x8000 | 0x8009 |
| `5b5597160b802bd8` | 6 | 512 | 0x8000 | 0x8009 |
| `4df12b39b4b7e183` | 6 | 512 | 0x8000 | 0x8009 |
| `6d2dad9c4ad72337` | 6 | 512 | 0x8000 | 0x8009 |
| `c240ccbd57f17e97` | 5 | 512 | 0x8000 | 0x8009 |
| `4b6b0a49ccdb1ee5` | 5 | 512 | 0x8000 | 0x8009 |
| `f6a6d2816bdac9db` | 5 | 512 | 0x8000 | 0x8009 |
| `e6faa50f76b3a75d` | 4 | 512 | 0x8000 | 0x8009 |
| `50b66a27b761f244` | 4 | 512 | 0x8000 | 0x8009 |
| `85d16f4eee2ba67e` | 4 | 512 | 0x8000 | 0x8009 |
| `b7ee4bb6f8c60ee8` | 3 | 512 | 0x8000 | 0x8009 |
| `ccb5768474312da9` | 3 | 512 | 0x8000 | 0x8009 |
| `e9a26ac37608bbf8` | 3 | 512 | 0x8000 | 0x8009 |
| `4511c4ed8b49a32b` | 3 | 512 | 0x8000 | 0x8009 |
| `c416aefbf00942b8` | 3 | 512 | 0x8000 | 0x8009 |
| `1b68f942e88f84f8` | 3 | 512 | 0x8000 | 0x8009 |
| `3b818ee9f809e4e4` | 3 | 512 | 0x8000 | 0x8009 |
| `55b1d9130dddc27f` | 3 | 512 | 0x8000 | 0x8009 |
| `ee9c7fe2dcfbafce` | 3 | 512 | 0x8000 | 0x8009 |
| `766d9dfb3b136556` | 3 | 512 | 0x8000 | 0x8009 |
| `deef05cf23285e89` | 2 | 512 | 0x8000 | 0x8009 |
| `29295a7f5c78068c` | 2 | 512 | 0x8000 | 0x8009 |
| `f6e9d43dca417bd5` | 2 | 512 | 0x8000 | 0x8009 |
| `eb6088cc0f6fccf7` | 2 | 512 | 0x8000 | 0x8009 |
| `a3c1ac35602ddd9d` | 2 | 512 | 0x8000 | 0x8009 |
| `86259913bc664d5c` | 2 | 512 | 0x8000 | 0x8009 |
| `9c8c8fb21ea6d771` | 2 | 512 | 0x8000 | 0x8009 |
| `09bb9af3ccbb6569` | 2 | 512 | 0x8000 | 0x8009 |
| `62ce35c94d4786ac` | 2 | 512 | 0x8000 | 0x8009 |
| `8ead55d69f64b27d` | 2 | 512 | 0x8000 | 0x8009 |
| `e930e2a2e5c851d8` | 2 | 512 | 0x8000 | 0x8009 |
| `4bddd4ce7743bad1` | 2 | 512 | 0x8000 | 0x8009 |
| `cd07a902828f24a9` | 2 | 512 | 0x8000 | 0x8009 |
| `db9487711eb9cbc0` | 2 | 512 | 0x8000 | 0x8009 |
| `ac5e85b75ab8aa63` | 2 | 512 | 0x8000 | 0x8009 |
| `2872a6d98436bab2` | 2 | 512 | 0x8000 | 0x8009 |
| `d2b64f81ab3d53e0` | 2 | 512 | 0x8000 | 0x8009 |
| `3418754d24020111` | 2 | 512 | 0x8000 | 0x8009 |
| `8680c87220ffd529` | 1 | 512 | 0x8000 | 0x8009 |
| `82c1ce414c995032` | 1 | 512 | 0x8000 | 0x8009 |
| `4b843655eabfb41f` | 1 | 512 | 0x8000 | 0x8009 |
| `1747832f2c385750` | 1 | 512 | 0x8000 | 0x8009 |
| `392a6f9ef594e6a3` | 1 | 512 | 0x8000 | 0x8009 |
| `1201a50c15e00147` | 1 | 512 | 0x8000 | 0x8009 |
| `5232f93ac580d4fe` | 1 | 512 | 0x8000 | 0x8009 |
| `5c0acba653cea138` | 1 | 512 | 0x8000 | 0x8009 |
| `8e0b10ab0b8831fc` | 1 | 512 | 0x8000 | 0x8009 |
| `1fe179b85df586c1` | 1 | 512 | 0x8000 | 0x8009 |
| `eb029b91accf8306` | 1 | 512 | 0x8000 | 0x8009 |
| `8b505a8add8077b8` | 1 | 512 | 0x8000 | 0x8009 |
| `4216a6eb511bb68f` | 1 | 512 | 0x8000 | 0x8009 |
| `0d702aa49de7e970` | 1 | 512 | 0x8000 | 0x8009 |
| `dc9b46f03af77b21` | 1 | 512 | 0x8000 | 0x8009 |
| `9b89e60179e6fbef` | 1 | 512 | 0x8000 | 0x8009 |
| `e19e91b7bf613135` | 1 | 512 | 0x8000 | 0x8009 |
| `497cf3561859ac5e` | 1 | 512 | 0x8000 | 0x8009 |
| `7a10b158f2a55a35` | 1 | 512 | 0x8000 | 0x8009 |
| `b79e46b21ee52d86` | 1 | 512 | 0x8000 | 0x8009 |
| `609acaf30bb85615` | 1 | 512 | 0x8000 | 0x8009 |
| `c31ab37a7681a16f` | 1 | 512 | 0x8000 | 0x8009 |
| `dd231ec7c6bd9604` | 1 | 512 | 0x8000 | 0x8009 |
| `c10caee3f9705914` | 1 | 512 | 0x8000 | 0x8009 |
| `44e5906feef11c02` | 1 | 512 | 0x8000 | 0x8009 |
| `9d504f31c57ecb24` | 1 | 512 | 0x8000 | 0x8009 |
| `fc9699be5a610098` | 1 | 512 | 0x8000 | 0x8009 |
| `007443365e7a65b7` | 1 | 512 | 0x8000 | 0x8009 |
| `fc6861899ced7a76` | 1 | 512 | 0x8000 | 0x8009 |
| `16287a10580ec570` | 1 | 512 | 0x8000 | 0x8009 |
| `fdb19d33613d456a` | 1 | 512 | 0x8000 | 0x8009 |
| `a166783458773a64` | 1 | 512 | 0x8000 | 0x8009 |
| `4688e18a7817bd59` | 1 | 512 | 0x8000 | 0x8009 |
| `c01d0af829a6f9f0` | 1 | 512 | 0x8000 | 0x8009 |
| `875529ddd84c5456` | 1 | 512 | 0x8000 | 0x8009 |
| `081824d021819d38` | 1 | 512 | 0x8000 | 0x8009 |
| `bb597c830bc92f93` | 1 | 512 | 0x8000 | 0x8009 |
| `a02484f7fdf17a9b` | 1 | 512 | 0x8000 | 0x8009 |
| `3e8fb17eed5c1089` | 1 | 512 | 0x8000 | 0x8009 |
| `9861f5bea0977727` | 1 | 512 | 0x8000 | 0x8009 |
| `d7b40324bcd5bdcd` | 1 | 512 | 0x8000 | 0x8009 |
| `a5f6ae9617f6be74` | 1 | 512 | 0x8000 | 0x8009 |
| `af646b4a1ceab99b` | 1 | 512 | 0x8000 | 0x8009 |
| `ceaeb89ba7914fd0` | 1 | 512 | 0x8000 | 0x8009 |
| `f060a0bced13e9ee` | 1 | 512 | 0x8000 | 0x8009 |
| `0ddeea1bde03ae81` | 1 | 512 | 0x8000 | 0x8009 |
| `848b730cb1f0d63d` | 1 | 512 | 0x8000 | 0x8009 |
| `391dd30b22a1e03f` | 1 | 512 | 0x8000 | 0x8009 |
| `8507ae08cb66654e` | 1 | 512 | 0x8000 | 0x8009 |
| `466277c1f8bfddd4` | 1 | 512 | 0x8000 | 0x8009 |
| `4990e2095bc6ac78` | 1 | 512 | 0x8000 | 0x8009 |
| `0dbcc86ffc7fdc71` | 1 | 512 | 0x8000 | 0x8009 |
| `68d68da587d6ab9d` | 1 | 512 | 0x8000 | 0x8009 |
| `be3a13574ef0f492` | 1 | 512 | 0x8000 | 0x8009 |
| `70b7c0403e3d6d75` | 1 | 512 | 0x8000 | 0x8009 |
| `de065797cbf44f45` | 1 | 512 | 0x8000 | 0x8009 |
| `d86321bd57c032fc` | 1 | 512 | 0x8000 | 0x8009 |
| `cd74791ef40672d6` | 1 | 512 | 0x8000 | 0x8009 |
| `40b7d159ca1f98b8` | 1 | 512 | 0x8000 | 0x8009 |
| `802b66f9c2de2a68` | 1 | 512 | 0x8000 | 0x8009 |
| `c608b04bf2f58ad3` | 1 | 512 | 0x8000 | 0x8009 |
| `45cbfa4660d5bda6` | 1 | 512 | 0x8000 | 0x8009 |
| `cf6ac48f09b3bb80` | 1 | 512 | 0x8000 | 0x8009 |
| `9e54788e1478dddd` | 1 | 512 | 0x8000 | 0x8009 |
| `59fc1d05800ae746` | 1 | 512 | 0x8000 | 0x8009 |
| `4661389efb98548b` | 1 | 512 | 0x8000 | 0x8009 |
| `074f0dde54899bbe` | 1 | 512 | 0x8000 | 0x8009 |
| `5cc4f71bf2f0dc66` | 1 | 512 | 0x8000 | 0x8009 |
| `80b7e5122cca472b` | 1 | 512 | 0x8000 | 0x8009 |
| `f6315b991e75b954` | 1 | 512 | 0x8000 | 0x8009 |
| `635835804cbf3c3d` | 1 | 512 | 0x8000 | 0x8009 |
| `e50e223cb17d4a0b` | 1 | 512 | 0x8000 | 0x8009 |
| `2e5a09de810292e1` | 1 | 512 | 0x8000 | 0x8009 |
| `d7fc904aa0622fde` | 1 | 512 | 0x8000 | 0x8009 |
| `332479ede4de2b0e` | 1 | 512 | 0x8000 | 0x8009 |
| `a63800a0044fa750` | 1 | 512 | 0x8000 | 0x8009 |
| `7e6641021bda7c21` | 1 | 512 | 0x8000 | 0x8009 |
| `eaa74f313a230bd5` | 1 | 512 | 0x8000 | 0x8009 |
| `97a5dba04d1e87ee` | 1 | 512 | 0x8000 | 0x8009 |
| `eb6d45fd22a60d6a` | 1 | 512 | 0x8000 | 0x8009 |
| `8ff769559974bb83` | 1 | 512 | 0x8000 | 0x8009 |
| `eb4ebacd71013015` | 1 | 512 | 0x8000 | 0x8009 |
| `7666ce4478663e80` | 1 | 512 | 0x8000 | 0x8009 |
| `daa2037181cecd6d` | 1 | 512 | 0x8000 | 0x8009 |
| `558c377038944c47` | 1 | 512 | 0x8000 | 0x8009 |
| `2d958a3033ae1fd1` | 1 | 512 | 0x8000 | 0x8009 |
| `e7ee00c820ea2ee5` | 1 | 512 | 0x8000 | 0x8009 |
| `a52ec1720ccd875f` | 1 | 512 | 0x8000 | 0x8009 |
| `5e19aff8608563c6` | 1 | 512 | 0x8000 | 0x8009 |
| `f636bce222a7e102` | 1 | 512 | 0x8000 | 0x8009 |
| `73507e38fdafae07` | 1 | 512 | 0x8000 | 0x8009 |
| `2b1858715d257e5a` | 1 | 512 | 0x8000 | 0x8009 |
| `afc9144c626311df` | 1 | 512 | 0x8000 | 0x8009 |
| `9f8ae5b3a3645d58` | 1 | 512 | 0x8000 | 0x8009 |
| `f4ae74ae7b9a677e` | 1 | 512 | 0x8000 | 0x8009 |
| `909d6f52f7267002` | 1 | 512 | 0x8000 | 0x8009 |
| `5a6f3adb185855d6` | 1 | 512 | 0x8000 | 0x8009 |
| `752512b1e76ca555` | 1 | 512 | 0x8000 | 0x8009 |
| `d169e8fb265b0cbc` | 1 | 512 | 0x8000 | 0x8009 |
| `13d214689c741686` | 1 | 512 | 0x8000 | 0x8009 |
| `37b138dc8823dd7e` | 1 | 512 | 0x8000 | 0x8009 |
| `39595ccad70a52b0` | 1 | 512 | 0x8000 | 0x8009 |
| `fa03bc0cf4ec1976` | 1 | 512 | 0x8000 | 0x8009 |
| `101b835ed1854201` | 1 | 512 | 0x8000 | 0x8009 |
| `3402b5dfd32c0b70` | 1 | 512 | 0x8000 | 0x8009 |
| `b07e204bc67e81ed` | 1 | 512 | 0x8000 | 0x8009 |
| `a3a21a733c1b7984` | 1 | 512 | 0x8000 | 0x8009 |
| `62ddd85e1c700d02` | 1 | 512 | 0x8000 | 0x8009 |
| `ed494e30032687f8` | 1 | 512 | 0x8000 | 0x8009 |
| `ad7b6b8436ed30c6` | 1 | 512 | 0x8000 | 0x8009 |
| `509c153360d1bb74` | 1 | 512 | 0x8000 | 0x8009 |
| `e740834e48361532` | 1 | 512 | 0x8000 | 0x8009 |
| `679d9b9ae42a6d99` | 1 | 512 | 0x8000 | 0x8009 |
| `b22c4c13a9aafdd1` | 1 | 512 | 0x8000 | 0x8009 |
| `d1a00eca89311d3a` | 1 | 512 | 0x8000 | 0x8009 |
| `32434baea90291fe` | 1 | 512 | 0x8000 | 0x8009 |
| `8fb23dceaf632798` | 1 | 512 | 0x8000 | 0x8009 |
| `63614da27dd3c44c` | 1 | 512 | 0x8000 | 0x8009 |
| `8a5769615c77d097` | 1 | 512 | 0x8000 | 0x8009 |
| `22a0d1caac693dcb` | 1 | 512 | 0x8000 | 0x8009 |
| `862e0d9bd3bf2c65` | 1 | 512 | 0x8000 | 0x8009 |
| `8c5493a5074c4c27` | 1 | 512 | 0x8000 | 0x8009 |
| `cf666ea3103bbc3a` | 1 | 512 | 0x8000 | 0x8009 |
| `31fc8f03692898f8` | 1 | 512 | 0x8000 | 0x8009 |
| `6a02ba1b7c14577c` | 1 | 512 | 0x8000 | 0x8009 |
| `6da48ff761e1394a` | 1 | 512 | 0x8000 | 0x8009 |
| `7caca696e909b466` | 1 | 512 | 0x8000 | 0x8009 |
| `0f83766f267f4a3c` | 1 | 512 | 0x8000 | 0x8009 |
| `08decb1ee6ae3703` | 1 | 512 | 0x8000 | 0x8009 |
| `22874ebd1077b63f` | 1 | 512 | 0x8000 | 0x8009 |
| `643aff0a77b2fd74` | 1 | 512 | 0x8000 | 0x8009 |
| `c09ad9634703d939` | 1 | 512 | 0x8000 | 0x8009 |
| `e564b0c37367397a` | 1 | 512 | 0x8000 | 0x8009 |

## Bootstrap-convention fingerprint (slot 0's file body)

Every known SAM DOS's bootstrap chooses to fetch the rest of
its code from slot 0's file (it doesn't have to — the ROM
contract is silent on this — but they all do). Decoding
slot 0's 9-byte file-header gives the precise byte count,
load address, and execution address declared by the file.
This is the *bootstrap-convention* fingerprint, not part of
the ROM contract; the actual bootstrap code at 0x8009 might
ignore these values and load anything from anywhere.

- Bootable disks with a usable slot-0 file: **554**
- Bootable disks with no usable slot-0 file: **9**
- Unique slot-0 bodies: **179**

### Unique slot-0 bodies (most-used first)

| SHA-256 (16) | Disks | Length | Load | Exec |
|---|---:|---:|---:|---:|
| `a69d4732a3274ede` | 107 | 8078 | 0x38009 | (none) |
| `78bc2964b7516db9` | 31 | 8077 | 0x78009 | (none) |
| `3cca541beb3f9fe9` | 30 | 10000 | 0x78009 | (none) |
| `13f6279c4d62e8be` | 26 | 15750 | 0x10000 | (none) |
| `9bc0fb4b949109e8` | 19 | 10000 | 0x38009 | (none) |
| `20e1c593dfd98cca` | 17 | 15700 | 0x10000 | (none) |
| `a242cc7b48a65541` | 16 | 10000 | 0x80009 | (none) |
| `fb0b95fdf1091299` | 15 | 10000 | 0x78009 | (none) |
| `739a6586e599ce10` | 13 | 10000 | 0x78009 | (none) |
| `2a3e618918a5520b` | 11 | 10000 | 0x38009 | (none) |
| `31771043e0729972` | 10 | 10000 | 0xe0be | (none) |
| `4c639e1325d9940f` | 10 | 10000 | 0x78009 | (none) |
| `1293e6ab6739adb3` | 8 | 10000 | 0x78009 | (none) |
| `254ae17a87efb171` | 6 | 8077 | 0x78009 | (none) |
| `14e4d256a8f0f0ec` | 6 | 10000 | 0x38009 | (none) |
| `d46ac18e935ff556` | 6 | 10000 | 0x78009 | (none) |
| `7e456418b0e18330` | 6 | 15700 | 0x10000 | (none) |
| `5f76f756fd1e6bee` | 5 | 10000 | 0x78009 | (none) |
| `1b0b65f8a9545787` | 5 | 10157 | 0x8009 | (none) |
| `152b811ed65b651d` | 5 | 15750 | 0x8000 | (none) |
| `2c189da5491097d9` | 4 | 15750 | 0x10000 | (none) |
| `3273b9e8ec9b1ecd` | 4 | 10000 | 0x38009 | (none) |
| `e7ead976f53c6003` | 4 | 8192 | 0x78009 | (none) |
| `fd6fa2869b471a1e` | 4 | 8078 | 0x78009 | (none) |
| `c753cf32520d0060` | 3 | 10000 | 0x78009 | (none) |
| `3a4be5a0b749094b` | 3 | 10000 | 0x78009 | (none) |
| `91e0f98622d2a6b0` | 3 | 32631 | 0x10000 | (none) |
| `21106301f8545821` | 3 | 8077 | 0x78009 | (none) |
| `d4d0e56b1342da0f` | 3 | 10000 | 0x38009 | (none) |
| `01a3ab44786698d7` | 3 | 10000 | 0x38009 | (none) |
| `61f76732d61326d7` | 3 | 10000 | 0x78009 | (none) |
| `401e21508fe023a5` | 3 | 10000 | 0x38009 | (none) |
| `cf4c3b28ab2d0144` | 3 | 10000 | 0x38009 | (none) |
| `59210e67e270d468` | 2 | 10000 | 0x78009 | (none) |
| `c76f8e68b0d0301b` | 2 | 15800 | 0x8009 | (none) |
| `cf248342d24596e3` | 2 | 10000 | 0x38009 | (none) |
| `b112c7d5d4340027` | 2 | 10000 | 0x38009 | (none) |
| `0415ede5ca82e665` | 2 | 10000 | 0x38009 | (none) |
| `58a02b8abd11d55b` | 2 | 10000 | 0x8030 | (none) |
| `df1ea02fd28e8997` | 2 | 10000 | 0x38009 | (none) |
| `6c53f7bb9dd34172` | 2 | 10000 | 0x38009 | (none) |
| `0f4d767f9db34845` | 2 | 8077 | 0x78009 | (none) |
| `f2199b85e766e967` | 2 | 10000 | 0x78009 | (none) |
| `440f66f1aab2a9f3` | 2 | 10000 | 0x78009 | (none) |
| `6f0c8686c2eafafe` | 2 | 10000 | 0x78009 | (none) |
| `c3c34de6552f864a` | 2 | 10000 | 0x78009 | (none) |
| `daeae9e0daf098a9` | 1 | 10000 | 0x78009 | (none) |
| `7089ca659242dde6` | 1 | 10000 | 0x78009 | (none) |
| `15ecc4085019ca35` | 1 | 10000 | 0x78009 | (none) |
| `557b1a4bc1e3e723` | 1 | 10000 | 0x78009 | (none) |
| `80b973784d34c635` | 1 | 10000 | 0x78009 | (none) |
| `46644ee5c8ceb328` | 1 | 10000 | 0x78009 | (none) |
| `78843b6a4b894771` | 1 | 10191 | 0x8009 | (none) |
| `7166b6af2054107e` | 1 | 14000 | 0x8009 | (none) |
| `521478fd84761030` | 1 | 15800 | 0x8009 | (none) |
| `f0047a502d0d54d9` | 1 | 501 | 0xe0be | (none) |
| `39f8558204cb3981` | 1 | 10000 | 0x8009 | (none) |
| `17f417b9a7ca6444` | 1 | 10000 | 0x78009 | (none) |
| `a7a17dd5b77a5555` | 1 | 10000 | 0x78009 | (none) |
| `c85f4d0eb8be8c3a` | 1 | 10000 | 0x78009 | (none) |
| `a49f2fd382925282` | 1 | 10000 | 0x78009 | (none) |
| `05017e20003819a7` | 1 | 10000 | 0x78009 | (none) |
| `1b853b911475c903` | 1 | 10000 | 0x38009 | (none) |
| `d852f669d7278df8` | 1 | 10000 | 0x78009 | (none) |
| `587fa1d449e85ef3` | 1 | 67976 | 0x8000 | 0x0000 |
| `16b08ca76ac9bf6c` | 1 | 9000 | 0x78009 | (none) |
| `571793c2f6a53f92` | 1 | 10000 | 0x78009 | (none) |
| `50edb1b9a5308f85` | 1 | 10000 | 0x78009 | (none) |
| `6a2f65a44273122f` | 1 | 8077 | 0x78009 | (none) |
| `3d31391ff91d110b` | 1 | 8192 | 0x78009 | (none) |
| `f9e25435a04c5542` | 1 | 8077 | 0x78009 | (none) |
| `25b3b8c3de323fc8` | 1 | 32631 | 0x10000 | (none) |
| `c98ea212d3f15722` | 1 | 10157 | 0x8009 | (none) |
| `eb1f8c0158856a40` | 1 | 10000 | 0x38009 | (none) |
| `9a3a718414aa71d0` | 1 | 8078 | 0x38009 | (none) |
| `a57d9139b31938c8` | 1 | 10000 | 0x38009 | (none) |
| `5a9d78bd06d11350` | 1 | 36957 | 0x8000 | 0x0000 |
| `f6f3c9f45cab0b92` | 1 | 10000 | 0x38009 | (none) |
| `801951340a516cac` | 1 | 10000 | 0x38009 | (none) |
| `5a93056404830c42` | 1 | 10000 | 0x38009 | (none) |
| `69c86a8b6360d876` | 1 | 10000 | 0x78009 | (none) |
| `5c389f3cfa578eaa` | 1 | 10000 | 0x38009 | (none) |
| `de2280bba1c6ba10` | 1 | 15750 | 0x1c009 | (none) |
| `064224e791163227` | 1 | 10000 | 0x78009 | (none) |
| `1b6ab6536f4e23e2` | 1 | 10000 | 0x38009 | (none) |
| `5c55e456b89a177c` | 1 | 10000 | 0x78009 | (none) |
| `d481994aba88a063` | 1 | 10000 | 0x78009 | (none) |
| `0b232c8d426c1534` | 1 | 10000 | 0x78009 | (none) |
| `dc5bc13f03508224` | 1 | 32631 | 0x10000 | (none) |
| `c3202ec6d71daf64` | 1 | 107317 | 0x8000 | 0x0000 |
| `f450085de2d9c53a` | 1 | 154354 | 0x8000 | 0x0000 |
| `ff9990ed35284292` | 1 | 10000 | 0x78009 | (none) |
| `4c400685ae82a2ab` | 1 | 10000 | 0x78009 | (none) |
| `487854350502cf42` | 1 | 14000 | 0x8009 | (none) |
| `0730ab3eb08701f4` | 1 | 9999 | 0x78009 | (none) |
| `470699700014483a` | 1 | 9999 | 0x78009 | (none) |
| `9057be073b6d042a` | 1 | 15750 | 0x10000 | (none) |
| `78093dc9050f1fe4` | 1 | 10000 | 0x78009 | (none) |
| `ba47e62e42fa6c83` | 1 | 10000 | 0x78009 | (none) |
| `a041404837bb1f94` | 1 | 10000 | 0x78009 | (none) |
| `e80ef68803505adf` | 1 | 15700 | 0x10000 | (none) |
| `16d35cdb1c766e7f` | 1 | 8976 | 0x78009 | (none) |
| `891bc46720d223f0` | 1 | 10000 | 0x78009 | (none) |
| `24727a275424024e` | 1 | 15700 | 0x10000 | (none) |
| `8ac2f0dc49b4ffe6` | 1 | 10000 | 0x78009 | (none) |
| `6e4c75fbba87c8ee` | 1 | 15800 | 0x8009 | (none) |
| `08160038384ce831` | 1 | 32631 | 0x10000 | (none) |
| `98d52753a0ca13cf` | 1 | 10000 | 0x38009 | (none) |
| `8a81ca8c8271f5ab` | 1 | 10000 | 0x38009 | (none) |
| `b366bf0c8841fbc3` | 1 | 10000 | 0x78009 | (none) |
| `080e2ab3c97bc23c` | 1 | 10000 | 0x78009 | (none) |
| `6a6feb5c4a7b1d91` | 1 | 10000 | 0x38009 | (none) |
| `abfdf5499cae5497` | 1 | 10000 | 0x78009 | (none) |
| `8a73776828631b6b` | 1 | 8078 | 0x38009 | (none) |
| `298b3e7150d3b3be` | 1 | 10000 | 0x78009 | (none) |
| `3be4384702b9054a` | 1 | 10000 | 0x78009 | (none) |
| `49aa9a49b0964aba` | 1 | 10000 | 0x78009 | (none) |
| `864506f419c442e3` | 1 | 10000 | 0x78009 | (none) |
| `96e56fea822f33ab` | 1 | 10000 | 0x78009 | (none) |
| `7587b034dcbcb1b7` | 1 | 10000 | 0x78009 | (none) |
| `6f82050a2a0a43d3` | 1 | 10000 | 0x78009 | (none) |
| `1ae0eda46245dfa8` | 1 | 15700 | 0x10000 | (none) |
| `68b90ca31c5f14e8` | 1 | 73567 | 0x8000 | 0x0000 |
| `b4fe3be82ac7a7b0` | 1 | 10000 | 0x78009 | (none) |
| `4dc74e1fc51f82bf` | 1 | 8100 | 0x7530 | (none) |
| `222cf6e7c0f515d8` | 1 | 10000 | 0x78009 | (none) |
| `82e592ca1b61476a` | 1 | 10000 | 0x78009 | (none) |
| `6fc1cdddbab20aa1` | 1 | 10000 | 0x78009 | (none) |
| `167956b902b73b6a` | 1 | 10000 | 0x78009 | (none) |
| `4eb0176725953059` | 1 | 10000 | 0x78009 | (none) |
| `703faf63d1cd4532` | 1 | 10000 | 0x78009 | (none) |
| `1e2e616c034ca8b5` | 1 | 10000 | 0x78009 | (none) |
| `522789453c5843a2` | 1 | 10000 | 0x78009 | (none) |
| `1d6f023e8820be49` | 1 | 10000 | 0x78009 | (none) |
| `da0c1ba197cd8010` | 1 | 10000 | 0x78009 | (none) |
| `bf76f513f5ff657c` | 1 | 10000 | 0x78009 | (none) |
| `81669e1f07422cb0` | 1 | 10000 | 0x78009 | (none) |
| `fa82a7e1450378c3` | 1 | 10000 | 0x78009 | (none) |
| `1aac98174ed8a75c` | 1 | 10000 | 0x78009 | (none) |
| `6100659ce98702b2` | 1 | 10000 | 0x78009 | (none) |
| `8e148d8cc88349bc` | 1 | 10000 | 0x78009 | (none) |
| `f22c4afc0cb24491` | 1 | 10000 | 0x78009 | (none) |
| `9cd1813018d2687c` | 1 | 10000 | 0x78009 | (none) |
| `8544bcebf58dfc8f` | 1 | 10000 | 0x78009 | (none) |
| `792d52b0c2ef5e61` | 1 | 10000 | 0x78009 | (none) |
| `173055c1d00f401e` | 1 | 10000 | 0x78009 | (none) |
| `d6152674405d34eb` | 1 | 10000 | 0x78009 | (none) |
| `41be122a791eae30` | 1 | 10000 | 0x78009 | (none) |
| `87f1a1f97596fa74` | 1 | 10000 | 0x78009 | (none) |
| `fd07f97df85dc17a` | 1 | 10000 | 0x78009 | (none) |
| `e39224f96381c658` | 1 | 10000 | 0x78009 | (none) |
| `e7312179003c9fe2` | 1 | 10000 | 0x78009 | (none) |
| `03f24fabd55fc7d5` | 1 | 10000 | 0x78009 | (none) |
| `da9251cfcaa28773` | 1 | 10000 | 0x78009 | (none) |
| `98785744cd87ccc2` | 1 | 10000 | 0x78009 | (none) |
| `f017586eb9b3272d` | 1 | 10000 | 0x78009 | (none) |
| `612826e7b8ebc753` | 1 | 10000 | 0x78009 | (none) |
| `999499e4b78c2e0c` | 1 | 10000 | 0x78009 | (none) |
| `13334f4bf5b5fa2c` | 1 | 10000 | 0x78009 | (none) |
| `3ca200dbb6588f8b` | 1 | 10000 | 0x78009 | (none) |
| `902ff6d466ae5e14` | 1 | 10000 | 0x78009 | (none) |
| `b3e1f498510fc710` | 1 | 9792 | 0x17d00 | (none) |
| `d031665e45db3c86` | 1 | 10000 | 0x78009 | (none) |
| `b9b42cf1534fc63d` | 1 | 10000 | 0x78009 | (none) |
| `2a375be311585fba` | 1 | 10000 | 0x78009 | (none) |
| `3663b9b5c46b3394` | 1 | 10000 | 0x78009 | (none) |
| `4b9daee4e076df6d` | 1 | 10000 | 0x78009 | (none) |
| `c1dc81fc3674eed2` | 1 | 10000 | 0x38009 | (none) |
| `f5b9bb560f6a0d3b` | 1 | 10000 | 0x38009 | (none) |
| `a3a7f8bf24d650ef` | 1 | 15700 | 0x10000 | (none) |
| `e88ff2a2e71e8516` | 1 | 10000 | 0x38009 | (none) |
| `214159b9cffa2359` | 1 | 10000 | 0x38009 | (none) |
| `5f7c9e75fbb60438` | 1 | 10000 | 0x38009 | (none) |
| `bec2a8d41401e03d` | 1 | 10000 | 0x38009 | (none) |
| `7b713de1bc50a02f` | 1 | 10000 | 0x38009 | (none) |
| `dedf6e1bee3729d4` | 1 | 10000 | 0x78009 | (none) |
| `774e500863e77368` | 1 | 10000 | 0x78009 | (none) |
| `88a75d769da6a53f` | 1 | 11044 | 0x8009 | 0x0000 |
| `a99e1f2160ecbd59` | 1 | 10000 | 0x8000 | 0x0000 |

## Sample disks per slot-0 body

Extract any of these for disassembly:
```bash
python3 ~/git/samfile/tools/audit/extract_dos.py <hash-prefix>
```

### `a69d4732a3274ede` (107 disks, length=8078 bytes, load=0x38009, exec=(none))

- Blast Turbo_ by James R Curry (1995) (PD)
- COMMIX V2.00 by S. Grodkowski (1995) (PD)
- COMMIX V2.01 by S. Grodkowski (1995) (PD)
- COMMIX V2.02 by S. Grodkowski (1995) (PD)
- Easydisc V4.9 (1995) (Saturn Software)

### `78bc2964b7516db9` (31 disks, length=8077 bytes, load=0x78009, exec=(none))

- Sam Adventure Club Issue 01 (Nov 1991)
- Sam Adventure Club Issue 02 (Jan 1992)
- Sam Adventure Club Issue 03 (Mar 1992) _a1_
- Sam Adventure Club Issue 03 (Mar 1992)
- Sam Adventure Club Issue 04 (May 1992) _a1_

### `3cca541beb3f9fe9` (30 disks, length=10000 bytes, load=0x78009, exec=(none))

- 32 Colour Demo by Gordon Wallis (1992) (PD)
- Allan Stevens - Home Utilities (1994)
- Allan Stevens - Home Utilities - Seven Pack (1994)
- Comms Loader (19xx)
- F-16 Combat Pilot Demo (1991) (PD)

### `13f6279c4d62e8be` (26 disks, length=15750 bytes, load=0x10000, exec=(none))

- E-Tracker Program Disk V1.2 (19xx) (FRED Publishing)
- Pics from the Net 06 (19xx) (PD)
- Pics from the Net 08 (19xx) (PD)
- Pics from the Net 09 (19xx) (PD)
- Pics from the Net 10 (19xx) (PD)

### `9bc0fb4b949109e8` (19 disks, length=10000 bytes, load=0x38009, exec=(none))

- Curse of the Serpent_s Eye_ The (1994) (Dream World) _a1_
- Curse of the Serpent_s Eye_ The (1994) (Dream World) _a2_
- Curse of the Serpent_s Eye_ The (1994) (Dream World) _a3_
- Curse of the Serpent_s Eye_ The (1994) (Dream World)
- FRED Magazine Issue 02 (1990)

### `20e1c593dfd98cca` (17 disks, length=15700 bytes, load=0x10000, exec=(none))

- Spectrum 128 Music Disk 2 (19xx) (PD)
- Spectrum Games Compilation 02 (1992)
- Spectrum Games Compilation 03 (1992)
- Spectrum Games Compilation 04 (1992)
- Spectrum Games Compilation 05 (1992)

### `a242cc7b48a65541` (16 disks, length=10000 bytes, load=0x80009, exec=(none))

- Adventures of Captain Comic_ The (19xx) (Lars)
- Astroball Demo by Balor Knight (1992) (PD)
- Defender - Persona and Digital Reality (1998) (Chris Pile)
- Dino Sourcerer (1993) (Softdisk Inc.)
- Ice Chicken Demo by ESI (1995) (PD)

### `fb0b95fdf1091299` (15 disks, length=10000 bytes, load=0x78009, exec=(none))

- Metempsychosis Demo - Christine (19xx)
- Metempsychosis Demo - Highlander (19xx)
- Metempsychosis Unreleased Demo - Internal_highlander (19xx)
- Metempsychosis Unreleased Demo - Internal_joy_pen (19xx)
- Metempsychosis Unreleased Demo - Mega_mix (19xx)

### `739a6586e599ce10` (13 disks, length=10000 bytes, load=0x78009, exec=(none))

- Cheats By Paul Crompton (1994) (PD)
- Driver Icons (1995) (Saturn Software)
- Sam Adventure Club Library Disk 3 SamScratch V3.2 (19xx)
- Sam Supplement Magazine Issue 06 (Mar 1991)
- Sam Supplement Magazine Issue 20 (May 1992)

### `2a3e618918a5520b` (11 disks, length=10000 bytes, load=0x38009, exec=(none))

- Allan Stevens - Capricorn Software Disk 1 (1994)
- Allan Stevens Compilation - Gallery 1 (19xx)
- Allan Stevens Compilation - Gallery 2 (19xx)
- Allan Stevens Compilation - Gallery 3 (19xx)
- Allan Stevens Compilation - Games Disk 1 (19xx)

### `31771043e0729972` (10 disks, length=10000 bytes, load=0xe0be, exec=(none))

- Banzai - The Games Compliation by Dan Doore (1995) (PD) _a1_
- Banzai - The Games Compliation by Dan Doore (1995) (PD)
- Banzai Babes 1 - Cindy _ Claudia by Dan Doore (1994) (PD)
- Banzai Babes 2 - Claudia - Elle - Kate by Dan Doore (1994) (PD)
- Banzai Pictures II - The Atari Job by Dan Doore (1994) (PD)

### `4c639e1325d9940f` (10 disks, length=10000 bytes, load=0x78009, exec=(none))

- Bats _n_ Balls Demo by Lord Insanity (1992) (PD)
- Craft Demo by ESI (1992) (PD)
- Hexagonia Demo by Fuxoft (1992) (PD)
- Lettis Demo by Daniel Cannon (1992) (PD)
- Sam News Disk Issue 1 (Jan 1992) (SAM Computers LTD)

### `1293e6ab6739adb3` (8 disks, length=10000 bytes, load=0x78009, exec=(none))

- FRED Magazine Issue 35 (1993) _a1_
- FRED Magazine Issue 35 (1993)
- FRED Magazine Issue 36 (1993)
- FRED Magazine Issue 37 (1993)
- FRED Magazine Issue 38 (1993) _a1_

### `254ae17a87efb171` (6 disks, length=8077 bytes, load=0x78009, exec=(none))

- Best of ENCELADUS_ Birthday Pack Edition (19xx) (Relion)
- ENCELADUS Magazine Issue 05 (Jun 1991) (Relion Software)
- ENCELADUS Magazine Issue 06 (Aug 1991) (Relion Software) _a1_
- ENCELADUS Magazine Issue 06 (Aug 1991) (Relion Software)
- ENCELADUS Magazine Issue 07 (Oct 1991) (Relion Software)

### `14e4d256a8f0f0ec` (6 disks, length=10000 bytes, load=0x38009, exec=(none))

- FRED Magazine Issue 15 (1991)
- FRED Magazine Issue 20 (1992) _a1_
- FRED Magazine Issue 20 (1992)
- FRED Magazine Issue 21 (1992) _a1_
- FRED Magazine Issue 21 (1992)

### `d46ac18e935ff556` (6 disks, length=10000 bytes, load=0x78009, exec=(none))

- Metempsychosis Unreleased Demo - CD_demo (19xx)
- Metempsychosis Unreleased Demo - Internal_advert (19xx)
- Metempsychosis Unreleased Demo - Internal_small_demos (19xx)
- Metempsychosis Unreleased Demo - Jukebox (19xx)
- Metempsychosis Unreleased Demo - RTJ_pdm3 (19xx)

### `7e456418b0e18330` (6 disks, length=15700 bytes, load=0x10000, exec=(none))

- Mouse Flash 1.1 for MDOS (19xx)
- Outwrite V2.0 (1992) (Chezron Software)
- PAX Disk 1 (1996) (Glenco)
- SC Filer V2.0 (1991) (Steve_s Software)
- Sam Adventure System Test Disk (1992) (Axxent Software)

### `5f76f756fd1e6bee` (5 disks, length=10000 bytes, load=0x78009, exec=(none))

- Allan Stevens - 50 Programs to Play and Write (19xx)
- Allan Stevens - Capricorn Software Disk 2 (1994)
- Allan Stevens - Colour Cycle (19xx)
- Allan Stevens - Spectrum Screens (19xx)
- Allan Stevens - World Finder (1994)

### `1b0b65f8a9545787` (5 disks, length=10157 bytes, load=0x8009, exec=(none))

- Blitz Magazine Issue 2 (1997) (Persona)
- Blitz Magazine Issue 4A (1997) (Persona)
- COMET to ASCII by Simon Cooke (1995)
- FRED Magazine Issue 82 (1997)
- Fashoom_ (1997) (Sad Snail Productions)

### `152b811ed65b651d` (5 disks, length=15750 bytes, load=0x8000, exec=(none))

- CometAssembler1.8EdwinBlink
- FRED Magazine - Morkography (1992)
- Flight of Fantasy and Occult Connection Adventures (19xx)
- MasterDOS V2.1 (19xx)
- Recover-E (1995)

### `2c189da5491097d9` (4 disks, length=15750 bytes, load=0x10000, exec=(none))

- Amiga MODS Disk (19xx)
- ETrackerv1.2
- Interlaced RGB Viewer Pics (19xx) (PD)
- Metempsychosis Sample Disk 7 (19xx)

### `3273b9e8ec9b1ecd` (4 disks, length=10000 bytes, load=0x38009, exec=(none))

- Arcadia Disk Magazine _1 (1991)
- Arcadia Disk Magazine _2 (1991)
- Arcadia Disk Magazine _3 (1991)
- Arcadia Disk Magazine _5 (1991)

### `e7ead976f53c6003` (4 disks, length=8192 bytes, load=0x78009, exec=(none))

- Mono Clipart Samples V1.0 (Nov 1995) (Steve_s Software)
- SC Compressor 2 (1991) (Steve_s Software)
- SC Monitor Pro 1.2_ TurboMon 1.0 (1992) (Steve_s Software)
- SC PD 3 by Steve_s Software (1992) (PD)

### `fd6fa2869b471a1e` (4 disks, length=8078 bytes, load=0x78009, exec=(none))

- Screens Viewer Disk 1 (19xx)
- Screens Viewer Disk 2 (19xx)
- Screens Viewer Disk 3 (19xx)
- Screens Viewer Disk 4 (19xx)

### `c753cf32520d0060` (3 disks, length=10000 bytes, load=0x78009, exec=(none))

- Allan Stevens Compilation - Games Disk 3 (19xx)
- FRED Magazine Issue 02 (1990) _a1_
- Neil Holmes_ Boing_ Graphics (1992) (Noesis Software)

### `3a4be5a0b749094b` (3 disks, length=10000 bytes, load=0x78009, exec=(none))

- Blinky Samples Disk 1 (1997) (Edwin Blink)
- Blinky Samples Disk 3 (1997) (Edwin Blink)
- Blitz 6 Menu by Edwin Blink (1997) (PD)

### `91e0f98622d2a6b0` (3 disks, length=32631 bytes, load=0x10000, exec=(none))

- DS12 Duff Capers Music Demo (2003) (PD)
- DS7 RainBow Scroller (2002) (PD)
- Duff Capers v0.51 by Tobermory (2003) (PD)

### `21106301f8545821` (3 disks, length=8077 bytes, load=0x78009, exec=(none))

- ENCELADUS Magazine Issue 09 (Feb 1992) (Relion Software)
- ENCELADUS Magazine Issue 11 (Jun 1992) (Relion Software)
- ENCELADUS Magazine Issue 12 (Oct 1992) (Relion Software)

### `d4d0e56b1342da0f` (3 disks, length=10000 bytes, load=0x38009, exec=(none))

- FRED Magazine Issue 01 (1990) _a1_
- FRED Magazine Issue 01 (1990)
- FRED Magazine Issue 01 _ 02 (1990)

### `01a3ab44786698d7` (3 disks, length=10000 bytes, load=0x38009, exec=(none))

- FRED Magazine Issue 12 (1991) _a1_
- FRED Magazine Issue 12 (1991) _a2_
- FRED Magazine Issue 12 (1991)

### `61f76732d61326d7` (3 disks, length=10000 bytes, load=0x78009, exec=(none))

- Golden Sword of Bhakhor_ The (1997) (Persona)
- Sam Supplement Magazine Issue 45 (Jun 1994)
- Sam Supplement Magazine Issue 46 (Jul 1994)

### `401e21508fe023a5` (3 disks, length=10000 bytes, load=0x38009, exec=(none))

- Sam Paper Magazine Issue 6 (19xx)
- Sam Paper Magazine Issue 8 (19xx)
- Sam Paper Magazine Issue 9 (19xx)

### `cf4c3b28ab2d0144` (3 disks, length=10000 bytes, load=0x38009, exec=(none))

- Visually 2 (19xx) (Zenith Graphics)
- Visually 3 (19xx) (Zenith Graphics)
- Visually 5 (19xx) (Zenith Graphics)

### `59210e67e270d468` (2 disks, length=10000 bytes, load=0x78009, exec=(none))

- Allan Stevens Compilation - Spectrum Disk 1 (19xx)
- Sam CD2 Utility (1990) (Kobrahsoft)

### `c76f8e68b0d0301b` (2 disks, length=15800 bytes, load=0x8009, exec=(none))

- B-DOS V1.7N (1999) (Martijn Groen _ Edwin Blink) (PD)
- Blancmange Burps 2 - Smell the Glove_ by Tobermory (2001) (PD)

### `cf248342d24596e3` (2 disks, length=10000 bytes, load=0x38009, exec=(none))

- Banzai - The Demos _ Utils by Dan Doore (1994) (PD)
- Banzai - The Double by Dan Doore (1994) (PD)

### `b112c7d5d4340027` (2 disks, length=10000 bytes, load=0x38009, exec=(none))

- Best of KAPSA I-V_ The (1993) (KAPSA) _a1_
- Best of KAPSA I-V_ The (1993) (KAPSA)

### `0415ede5ca82e665` (2 disks, length=10000 bytes, load=0x38009, exec=(none))

- Edition 3 (1991) (Zenith Graphics) (PD)
- Zenith Edition 3 (19xx) (Zenith Graphics)

### `58a02b8abd11d55b` (2 disks, length=10000 bytes, load=0x8030, exec=(none))

- Entropy Demo (1992) (PD)
- Goblin Mountain Adventure (1993) (Sam PD Sware Lib.) (PD)

### `df1ea02fd28e8997` (2 disks, length=10000 bytes, load=0x38009, exec=(none))

- Fastline Public Domain Library Disk 13 (19xx)
- Silly Demo 1 by Lord Insanity (1990) (PD)

### `6c53f7bb9dd34172` (2 disks, length=10000 bytes, load=0x38009, exec=(none))

- Impatience - Triltex Viking (1991) (Fred Publishing)
- Occult Adventure by David Munden (1993) (PD)

### `0f4d767f9db34845` (2 disks, length=8077 bytes, load=0x78009, exec=(none))

- Sam Adventure Club Issue 09b (Mar 1993) _a1_
- Sam Adventure Club Issue 09b (Mar 1993)

### `f2199b85e766e967` (2 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine - The Best Of (19xx)
- Sam Supplement Magazine Issue 12A (Sep 1991)

### `440f66f1aab2a9f3` (2 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 33 (Jun 1993)
- Sam Supplement Magazine Issue 34 (Jul 1993)

### `6f0c8686c2eafafe` (2 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 37 (Oct 1993)
- Sam Supplement Magazine Issue 38 (Nov 1993)

### `c3c34de6552f864a` (2 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 50 (Nov 1994) _b1_
- Sam Supplement Magazine Issue 50 (Nov 1994)

### `daeae9e0daf098a9` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- 18 Rated Poker for 512k (19xx) (Supplement Software)

### `7089ca659242dde6` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Aliens vs Predator Demo by Gordon Wallis (1991) (PD)

### `15ecc4085019ca35` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Allan Stevens - Capricorn Software Disk 3 Unfinished (1994)

### `557b1a4bc1e3e723` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Allan Stevens - Learn With Timmy Under 6s (1994)

### `80b973784d34c635` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Allan Stevens - Learning Games (19xx)

### `46644ee5c8ceb328` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Allan Stevens Compilation - Games Disk 4 (19xx)

### `78843b6a4b894771` (1 disks, length=10191 bytes, load=0x8009, exec=(none))

- B-DOS V1.5A (1997) (Martijn Groen _ Edwin Blink) (PD)

### `7166b6af2054107e` (1 disks, length=14000 bytes, load=0x8009, exec=(none))

- B-DOS V1.7D (1999) (Martijn Groen _ Edwin Blink) (PD)

### `521478fd84761030` (1 disks, length=15800 bytes, load=0x8009, exec=(none))

- B-DOS V1.7J (1999) (Martijn Groen _ Edwin Blink) (PD)

### `f0047a502d0d54d9` (1 disks, length=501 bytes, load=0xe0be, exec=(none))

- Banzai Pictures I by Dan Doore (1994) (PD)

### `39f8558204cb3981` (1 disks, length=10000 bytes, load=0x8009, exec=(none))

- Blinky Samples Disk 4 (1997) (Edwin Blink)

### `17f417b9a7ca6444` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Blitz Magazine Issue 4B - Sound Machine (1997) (Persona) _b1_

### `a7a17dd5b77a5555` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Blokker (19xx) (Stephen McGreal)

### `c85f4d0eb8be8c3a` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Bunnik 2K MOD Slide (2000) (Edwin Blink)

### `a49f2fd382925282` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Defender Compilation (19xx)

### `05017e20003819a7` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Demo Disk (1994) (PD)

### `1b853b911475c903` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Diaz Demo 2 (19xx)

### `d852f669d7278df8` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Domino Box (1992) (Supplement Software)

### `587fa1d449e85ef3` (1 disks, length=67976 bytes, load=0x8000, exec=0x0000)

- E-Tunes Player (19xx) (Andrew Collier)

### `16b08ca76ac9bf6c` (1 disks, length=9000 bytes, load=0x78009, exec=(none))

- ENCELADUS - Complete Guide to SAMBASIC Parts 1-7 (1994) (Relion)

### `571793c2f6a53f92` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- ENCELADUS Magazine Issue 01 (Oct 1990) (Relion Software)

### `50edb1b9a5308f85` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- ENCELADUS Magazine Issue 02 (Dec 1990) (Relion Software)

### `6a2f65a44273122f` (1 disks, length=8077 bytes, load=0x78009, exec=(none))

- ENCELADUS Magazine Issue 03 (Feb 1991) (Relion Software)

### `3d31391ff91d110b` (1 disks, length=8192 bytes, load=0x78009, exec=(none))

- ENCELADUS Magazine Issue 04 (Apr 1991) (Relion Software)

### `f9e25435a04c5542` (1 disks, length=8077 bytes, load=0x78009, exec=(none))

- ENCELADUS Magazine Issue 10 (Apr 1992) (Relion Software)

### `25b3b8c3de323fc8` (1 disks, length=32631 bytes, load=0x10000, exec=(none))

- EXPLOSION - ZX SPECTRUM 48 Emulator _ COMMANDER (1996)

### `c98ea212d3f15722` (1 disks, length=10157 bytes, load=0x8009, exec=(none))

- Entropy Demo (1992) (PD) _a2_

### `eb1f8c0158856a40` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- FRED Magazine Issue 28 (1992)

### `9a3a718414aa71d0` (1 disks, length=8078 bytes, load=0x38009, exec=(none))

- FRED Magazine Issue 33 (1993)

### `a57d9139b31938c8` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- FRED Magazine Issue 34 (1993)

### `5a9d78bd06d11350` (1 disks, length=36957 bytes, load=0x8000, exec=0x0000)

- FRED Magazine Issue 65 Menu (1995)

### `f6f3c9f45cab0b92` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Fastline Public Domain Library Disk 10 (19xx)

### `801951340a516cac` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Fastline Public Domain Library Disk 12 (19xx)

### `5a93056404830c42` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Fastline Public Domain Library Disk 16 (19xx)

### `69c86a8b6360d876` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Football League Manager (1994) (Key Software-FRED Publishing)

### `5c389f3cfa578eaa` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Fredatives 4 (1992)

### `de2280bba1c6ba10` (1 disks, length=15750 bytes, load=0x1c009, exec=(none))

- H-DOS V2.12 HD Loader V2.0 (1996)

### `064224e791163227` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Images SAM Software - Roboblob _ Give a Dog a Bone (19xx)

### `1b6ab6536f4e23e2` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Integrated Logic_s PD Disk (1990) (PD)

### `5c55e456b89a177c` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Jigsaw Disk - Jigsaw Creator (1992) (Colony Software)

### `d481994aba88a063` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- KEdisk V2.04 (19xx)

### `0b232c8d426c1534` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Lovehearts (19xx) (Supplement Software)

### `dc5bc13f03508224` (1 disks, length=32631 bytes, load=0x10000, exec=(none))

- MDOS _ MBASIC for Formatting Discs in 2 Drives (19xx)

### `c3202ec6d71daf64` (1 disks, length=107317 bytes, load=0x8000, exec=0x0000)

- MNEMOtech Demo 1 (19xx) (PD)

### `f450085de2d9c53a` (1 disks, length=154354 bytes, load=0x8000, exec=0x0000)

- MNEMOtech Demo 2 (19xx) (PD)

### `ff9990ed35284292` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Mega Demo 6 by Supplement Software (1991) (PD)

### `4c400685ae82a2ab` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Mega Text Demo III by Masters of Magic (19xx) (PD)

### `487854350502cf42` (1 disks, length=14000 bytes, load=0x8009, exec=(none))

- Megaboot V2.3 (Atom HD Interface) (1999) (M.Groen)

### `0730ab3eb08701f4` (1 disks, length=9999 bytes, load=0x78009, exec=(none))

- Megadisk 1 - Puzzles (19xx)

### `470699700014483a` (1 disks, length=9999 bytes, load=0x78009, exec=(none))

- Metempsychosis Unreleased Demo - Demos4metemp (19xx)

### `9057be073b6d042a` (1 disks, length=15750 bytes, load=0x10000, exec=(none))

- Metempsychosis Unreleased Demo - Internal_digi_utils (19xx)

### `78093dc9050f1fe4` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Metempsychosis Unreleased Demo - Kinetik (19xx)

### `ba47e62e42fa6c83` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Metempsychosis Unreleased Demo - RTJ_pdm1 (19xx)

### `a041404837bb1f94` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Metempsychosis Unreleased Demo - RTJ_pdm2 (19xx)

### `e80ef68803505adf` (1 disks, length=15700 bytes, load=0x10000, exec=(none))

- Metempsychosis Unreleased Demo - Wizard (19xx)

### `16d35cdb1c766e7f` (1 disks, length=8976 bytes, load=0x78009, exec=(none))

- Metempsychosis pdm12 (19xx)

### `891bc46720d223f0` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Metempsychosis slide_1 (19xx)

### `24727a275424024e` (1 disks, length=15700 bytes, load=0x10000, exec=(none))

- Mike AJ Disc 6-Edwin (19xx)

### `8ac2f0dc49b4ffe6` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Network Sigma Issue 6 (Feb-Mar 1996) (Saturn Software)

### `6e4c75fbba87c8ee` (1 disks, length=15800 bytes, load=0x8009, exec=(none))

- Open 3D V082 by Tobermory (2001) (PD)

### `08160038384ce831` (1 disks, length=32631 bytes, load=0x10000, exec=(none))

- Ore Warz II (1990) (William McGugan)

### `98d52753a0ca13cf` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Outlet - Sam _ Spectrum Mag Issue 33 (May 1990)

### `8a81ca8c8271f5ab` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Outlet - Sam _ Spectrum Mag Issue 56 (April 1992)

### `b366bf0c8841fbc3` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- PD_90 The Best of 1990 (1990) (PD)

### `080e2ab3c97bc23c` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Personal Filing System V2.07 (1994) (Hilton Comp. Services)

### `6a6feb5c4a7b1d91` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Prince of Persia Demo (1990) (Revelation-Chris White)

### `abfdf5499cae5497` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Printer Port Music Sample Player (19xx)

### `8a73776828631b6b` (1 disks, length=8078 bytes, load=0x38009, exec=(none))

- Public the Third (19xx)

### `298b3e7150d3b3be` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Robocop - Rolling Demo (19xx) (PD)

### `3be4384702b9054a` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Amateur Programming _ Electronics Issue 1 (Feb 1992)

### `49aa9a49b0964aba` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Amateur Programming _ Electronics Issue 2 (Mar 1992)

### `864506f419c442e3` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Amateur Programming _ Electronics Issue 4 (May 1992)

### `96e56fea822f33ab` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Amateur Programming _ Electronics Issue 5 (Jun 1992)

### `7587b034dcbcb1b7` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Amateur Programming _ Electronics Issue 6 (Aug 1992)

### `6f82050a2a0a43d3` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Cards (1994) (Supplement Software)

### `1ae0eda46245dfa8` (1 disks, length=15700 bytes, load=0x10000, exec=(none))

- Sam D I C E V1.0 for MasterDOS (1991) (Kobrahsoft)

### `68b90ca31c5f14e8` (1 disks, length=73567 bytes, load=0x8000, exec=0x0000)

- Sam Mines (19xx) (PD)

### `b4fe3be82ac7a7b0` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Paint (19xx)

### `4dc74e1fc51f82bf` (1 disks, length=8100 bytes, load=0x7530, exec=(none))

- Sam Prime (19xx)

### `222cf6e7c0f515d8` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 02 (Nov 1990)

### `82e592ca1b61476a` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 03 (Dec 1990)

### `6fc1cdddbab20aa1` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 04 (Jan 1991)

### `167956b902b73b6a` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 05 (Feb 1991)

### `4eb0176725953059` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 07 (Apr 1991)

### `703faf63d1cd4532` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 08 (May 1991)

### `1e2e616c034ca8b5` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 09 (Jun 1991)

### `522789453c5843a2` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 10 (Jul 1991)

### `1d6f023e8820be49` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 11 (Aug 1991)

### `da0c1ba197cd8010` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 12B Freeware (Sep 1991)

### `bf76f513f5ff657c` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 13 (Oct 1991)

### `81669e1f07422cb0` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 14 (Nov 1991)

### `fa82a7e1450378c3` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 15 (Dec 1991)

### `1aac98174ed8a75c` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 16 (Jan 1992)

### `6100659ce98702b2` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 17 (Feb 1992)

### `8e148d8cc88349bc` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 18 (Mar 1992)

### `f22c4afc0cb24491` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 19 (Apr 1992)

### `9cd1813018d2687c` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 22 (Jul 1992)

### `8544bcebf58dfc8f` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 24 (Sep 1992)

### `792d52b0c2ef5e61` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 27 (Dec 1992)

### `173055c1d00f401e` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 28 (Jan 1993)

### `d6152674405d34eb` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 32 (May 1993)

### `41be122a791eae30` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 35 (Aug 1993)

### `87f1a1f97596fa74` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 39 (Dec 1993)

### `fd07f97df85dc17a` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 40 (Jan 1994)

### `e39224f96381c658` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 43 (Apr 1994)

### `e7312179003c9fe2` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 47 (Aug 1994)

### `03f24fabd55fc7d5` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 48 (Sep 1994)

### `da9251cfcaa28773` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 49 (Oct 1994)

### `98785744cd87ccc2` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 51 (Dec 1994)

### `f017586eb9b3272d` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 52 (Jan 1995)

### `612826e7b8ebc753` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 53 (Feb 1995)

### `999499e4b78c2e0c` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Issue 54 (Mar 1995)

### `13334f4bf5b5fa2c` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Supplement Magazine Xmas Disk (Dec 1992)

### `3ca200dbb6588f8b` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam Utils (1993)

### `902ff6d466ae5e14` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sam X (19xx) (Supplement Software)

### `b3e1f498510fc710` (1 disks, length=9792 bytes, load=0x17d00, exec=(none))

- Samsational Complete Guide to SAM PD Software (1992) (SCPDSA)

### `d031665e45db3c86` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sandman_s Shadow 1 (1993) (PD)

### `b9b42cf1534fc63d` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sandman_s Shadow 2 (1993) (PD)

### `2a375be311585fba` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Sandman_s Shadow 3 (1993) (PD)

### `3663b9b5c46b3394` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Satellite_92 (1992)

### `4b9daee4e076df6d` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Spectrum 128 - Myth and Escape From Singe_s Castle (19xx)

### `c1dc81fc3674eed2` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Spectrum Emulator (Sept 04) (1990)

### `f5b9bb560f6a0d3b` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Spectrum Games Compilation 13 (1992)

### `a3a7f8bf24d650ef` (1 disks, length=15700 bytes, load=0x10000, exec=(none))

- TurboMON V1.0 (19xx)

### `e88ff2a2e71e8516` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Visually 1 (19xx) (Zenith Graphics)

### `214159b9cffa2359` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Visually 4 (19xx) (Zenith Graphics)

### `5f7c9e75fbb60438` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Visually 6 (19xx) (Zenith Graphics)

### `bec2a8d41401e03d` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Zenith Edition 1 (19xx) (Zenith Graphics)

### `7b713de1bc50a02f` (1 disks, length=10000 bytes, load=0x38009, exec=(none))

- Zenith Edition 2 (19xx) (Zenith Graphics)

### `dedf6e1bee3729d4` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- Zenith Edition 2-5 (19xx) (Zenith Graphics)

### `774e500863e77368` (1 disks, length=10000 bytes, load=0x78009, exec=(none))

- pete-made

### `88a75d769da6a53f` (1 disks, length=11044 bytes, load=0x8009, exec=0x0000)

- trinity

### `a99e1f2160ecbd59` (1 disks, length=10000 bytes, load=0x8000, exec=0x0000)

- trinload

## Disks ROM would refuse with `NO DOS` (237)

T4S1[256..259] doesn't match `BOOT` (masked 0x5F).
ROM BOOTEX prints `NO DOS` and refuses to load.
These disks have no bootable DOS; skipped from
classification.

Sample (first 20):
- AMRAD Amateur Radio Logbook (1994) (Spencer)
- All Star Belch by Tobermory (2001) (PD)
- Alternative Module Player _ Tunes V1.2 (1994) (Aley Keprt)
- Andy Monk_s Music (19xx) (PD)
- Arcadia Disk Magazine _3b (1991)
- Arnie_s Samples (1991)
- Banzai Babes 1 - Cindy _ Claudia by Dan Doore (1994) (PD) _a1_
- Basic Scrolly Pokers by Chris White (1992)
- Bats _n_ Balls by David Gommerman (1992) (Revelation)
- Beta DOS 1.0 For the Plus D (1990) (Betasoft)
- Blancmange Burps by Tobermory (2001) (PD)
- Blinky Samples Disk 2 (1997) (Edwin Blink)
- Blue Disk Show_ The (19xx)
- Boing _ Sphera (1992) (Noesis Software)
- Bombed Out_Nuclear Waste_Magic Caves_Blockade (1990) (Enigma)
- Bowin and the Count Dracula (1991) (Lucosoft and Revelation)
- COMET Opcodes by Tobermory (2001) (PD)
- COMET Z80 Assembler V1.8 (1992) (Revelation)
- Colour Clipart Samples V1.0 (Nov 1995) (Steve_s Software)
- Compressor Utilities (1993) (ESI) (PD)

## Bootable but no slot-0 file (9)

Disks where ROM would boot (BOOT signature present)
but slot 0 is erased / its chain malformed, so the
slot-0 bootstrap-convention fingerprint doesn't apply.
The T4S1 fingerprint above still classifies them.

- Blondie and Dagwood_ Arkanoid_ Prince of the Yolk Folk (199x)
- Integrated Logic_s Demo Disk and Utils (1990) (PD)
- Lyra 3 Megademo by ESI (1993) (PD) _a2_
- Lyra 3 Megademo by ESI (1993) (PD)
- Mike AJ Demo Disk 1 (1991) _a1_
- Mike AJ Demo Disk 1 (1991)
- Mike AJ Disc 7 (19xx)
- SAMart and Slideshow (19xx) (Sam PD Sware Lib.) (PD)
- newdisk

