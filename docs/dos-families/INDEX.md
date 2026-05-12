# DOS families — per-family directory tree

One subdirectory per family in [`docs/dos-families.md`].
Each contains the family-head body, a hex dump, READMEs
describing the variants and sample disks, and (when
available) the commented original assembly source.

## Summary

- **Families:** 45
- **Total disks across families:** 554
- **Families with upstream source attached:** 2

Big families (≥ 5 disks) materialise all member-variant
bodies under `variants/`. Smaller families only emit the
family-head body to keep the tree compact; if you need a
non-head variant body, run `tools/audit/extract_dos.py`
against any disk that contains it.

## Index

| Rank | Family | Variants | Disks | Length(s) | Source |
|---:|---|---:|---:|---|:-:|
| 1 | [`9bc0fb4b949109e8`](9bc0fb4b949109e8-samdos2/) (samdos2) | 126 | 291 | 10000 | yes |
| 2 | [`a69d4732a3274ede`](a69d4732a3274ede/) | 4 | 113 | 8078 | — |
| 3 | [`13f6279c4d62e8be`](13f6279c4d62e8be/) | 3 | 31 | 15750 | — |
| 4 | [`78bc2964b7516db9`](78bc2964b7516db9/) | 1 | 31 | 8077 | — |
| 5 | [`20e1c593dfd98cca`](20e1c593dfd98cca/) | 3 | 24 | 15700 | — |
| 6 | [`254ae17a87efb171`](254ae17a87efb171/) | 1 | 6 | 8077 | — |
| 7 | [`152b811ed65b651d`](152b811ed65b651d-masterdos-v2.3/) (masterdos-v2.3) | 2 | 6 | 15750 | yes |
| 8 | [`1b0b65f8a9545787`](1b0b65f8a9545787/) | 1 | 5 | 10157 | — |
| 9 | [`e7ead976f53c6003`](e7ead976f53c6003/) | 1 | 4 | 8192 | — |
| 10 | [`91e0f98622d2a6b0`](91e0f98622d2a6b0/) | 1 | 3 | 32631 | — |
| 11 | [`21106301f8545821`](21106301f8545821/) | 1 | 3 | 8077 | — |
| 12 | [`c76f8e68b0d0301b`](c76f8e68b0d0301b/) | 1 | 2 | 15800 | — |
| 13 | [`470699700014483a`](470699700014483a/) | 2 | 2 | 9999 | — |
| 14 | [`0f4d767f9db34845`](0f4d767f9db34845/) | 1 | 2 | 8077 | — |
| 15 | [`78843b6a4b894771`](78843b6a4b894771/) | 1 | 1 | 10191 | — |
| 16 | [`7166b6af2054107e`](7166b6af2054107e/) | 1 | 1 | 14000 | — |
| 17 | [`521478fd84761030`](521478fd84761030/) | 1 | 1 | 15800 | — |
| 18 | [`f0047a502d0d54d9`](f0047a502d0d54d9/) | 1 | 1 | 501 | — |
| 19 | [`39f8558204cb3981`](39f8558204cb3981/) | 1 | 1 | 10000 | — |
| 20 | [`587fa1d449e85ef3`](587fa1d449e85ef3/) | 1 | 1 | 67976 | — |
| 21 | [`16b08ca76ac9bf6c`](16b08ca76ac9bf6c/) | 1 | 1 | 9000 | — |
| 22 | [`571793c2f6a53f92`](571793c2f6a53f92/) | 1 | 1 | 10000 | — |
| 23 | [`50edb1b9a5308f85`](50edb1b9a5308f85/) | 1 | 1 | 10000 | — |
| 24 | [`6a2f65a44273122f`](6a2f65a44273122f/) | 1 | 1 | 8077 | — |
| 25 | [`3d31391ff91d110b`](3d31391ff91d110b/) | 1 | 1 | 8192 | — |
| 26 | [`f9e25435a04c5542`](f9e25435a04c5542/) | 1 | 1 | 8077 | — |
| 27 | [`25b3b8c3de323fc8`](25b3b8c3de323fc8/) | 1 | 1 | 32631 | — |
| 28 | [`c98ea212d3f15722`](c98ea212d3f15722/) | 1 | 1 | 10157 | — |
| 29 | [`5a9d78bd06d11350`](5a9d78bd06d11350/) | 1 | 1 | 36957 | — |
| 30 | [`dc5bc13f03508224`](dc5bc13f03508224/) | 1 | 1 | 32631 | — |
| 31 | [`c3202ec6d71daf64`](c3202ec6d71daf64/) | 1 | 1 | 107317 | — |
| 32 | [`f450085de2d9c53a`](f450085de2d9c53a/) | 1 | 1 | 154354 | — |
| 33 | [`487854350502cf42`](487854350502cf42/) | 1 | 1 | 14000 | — |
| 34 | [`16d35cdb1c766e7f`](16d35cdb1c766e7f/) | 1 | 1 | 8976 | — |
| 35 | [`24727a275424024e`](24727a275424024e/) | 1 | 1 | 15700 | — |
| 36 | [`6e4c75fbba87c8ee`](6e4c75fbba87c8ee/) | 1 | 1 | 15800 | — |
| 37 | [`08160038384ce831`](08160038384ce831/) | 1 | 1 | 32631 | — |
| 38 | [`1ae0eda46245dfa8`](1ae0eda46245dfa8/) | 1 | 1 | 15700 | — |
| 39 | [`68b90ca31c5f14e8`](68b90ca31c5f14e8/) | 1 | 1 | 73567 | — |
| 40 | [`4dc74e1fc51f82bf`](4dc74e1fc51f82bf/) | 1 | 1 | 8100 | — |
| 41 | [`b3e1f498510fc710`](b3e1f498510fc710/) | 1 | 1 | 9792 | — |
| 42 | [`c1dc81fc3674eed2`](c1dc81fc3674eed2/) | 1 | 1 | 10000 | — |
| 43 | [`a3a7f8bf24d650ef`](a3a7f8bf24d650ef/) | 1 | 1 | 15700 | — |
| 44 | [`bec2a8d41401e03d`](bec2a8d41401e03d/) | 1 | 1 | 10000 | — |
| 45 | [`88a75d769da6a53f`](88a75d769da6a53f/) | 1 | 1 | 11044 | — |

Regenerate this tree with `tools/audit/build_family_tree.py`.
