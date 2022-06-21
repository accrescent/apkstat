package apkstat

func arabParents() map[uint32]uint32 {
	return map[uint32]uint32{
		0x61724145: 0x61729420, // ar-AE -> ar-015
		0x6172445A: 0x61729420, // ar-DZ -> ar-015
		0x61724548: 0x61729420, // ar-EH -> ar-015
		0x61724C59: 0x61729420, // ar-LY -> ar-015
		0x61724D41: 0x61729420, // ar-MA -> ar-015
		0x6172544E: 0x61729420, // ar-TN -> ar-015
	}
}

func devaParents() map[uint32]uint32 {
	return map[uint32]uint32{
		0x68690000: 0x656E494E, // hi-Latn -> en-IN
	}
}

func hantParents() map[uint32]uint32 {
	return map[uint32]uint32{
		0x7A684D4F: 0x7A68484B, // zh-Hant-MO -> zh-Hant-HK
	}
}

func latnParents() map[uint32]uint32 {
	return map[uint32]uint32{
		0x656E80A1: 0x656E8400, // en-150 -> en-001
		0x656E4147: 0x656E8400, // en-AG -> en-001
		0x656E4149: 0x656E8400, // en-AI -> en-001
		0x656E4154: 0x656E80A1, // en-AT -> en-150
		0x656E4155: 0x656E8400, // en-AU -> en-001
		0x656E4242: 0x656E8400, // en-BB -> en-001
		0x656E4245: 0x656E80A1, // en-BE -> en-150
		0x656E424D: 0x656E8400, // en-BM -> en-001
		0x656E4253: 0x656E8400, // en-BS -> en-001
		0x656E4257: 0x656E8400, // en-BW -> en-001
		0x656E425A: 0x656E8400, // en-BZ -> en-001
		0x656E4343: 0x656E8400, // en-CC -> en-001
		0x656E4348: 0x656E80A1, // en-CH -> en-150
		0x656E434B: 0x656E8400, // en-CK -> en-001
		0x656E434D: 0x656E8400, // en-CM -> en-001
		0x656E4358: 0x656E8400, // en-CX -> en-001
		0x656E4359: 0x656E8400, // en-CY -> en-001
		0x656E4445: 0x656E80A1, // en-DE -> en-150
		0x656E4447: 0x656E8400, // en-DG -> en-001
		0x656E444B: 0x656E80A1, // en-DK -> en-150
		0x656E444D: 0x656E8400, // en-DM -> en-001
		0x656E4552: 0x656E8400, // en-ER -> en-001
		0x656E4649: 0x656E80A1, // en-FI -> en-150
		0x656E464A: 0x656E8400, // en-FJ -> en-001
		0x656E464B: 0x656E8400, // en-FK -> en-001
		0x656E464D: 0x656E8400, // en-FM -> en-001
		0x656E4742: 0x656E8400, // en-GB -> en-001
		0x656E4744: 0x656E8400, // en-GD -> en-001
		0x656E4747: 0x656E8400, // en-GG -> en-001
		0x656E4748: 0x656E8400, // en-GH -> en-001
		0x656E4749: 0x656E8400, // en-GI -> en-001
		0x656E474D: 0x656E8400, // en-GM -> en-001
		0x656E4759: 0x656E8400, // en-GY -> en-001
		0x656E484B: 0x656E8400, // en-HK -> en-001
		0x656E4945: 0x656E8400, // en-IE -> en-001
		0x656E494C: 0x656E8400, // en-IL -> en-001
		0x656E494D: 0x656E8400, // en-IM -> en-001
		0x656E494E: 0x656E8400, // en-IN -> en-001
		0x656E494F: 0x656E8400, // en-IO -> en-001
		0x656E4A45: 0x656E8400, // en-JE -> en-001
		0x656E4A4D: 0x656E8400, // en-JM -> en-001
		0x656E4B45: 0x656E8400, // en-KE -> en-001
		0x656E4B49: 0x656E8400, // en-KI -> en-001
		0x656E4B4E: 0x656E8400, // en-KN -> en-001
		0x656E4B59: 0x656E8400, // en-KY -> en-001
		0x656E4C43: 0x656E8400, // en-LC -> en-001
		0x656E4C52: 0x656E8400, // en-LR -> en-001
		0x656E4C53: 0x656E8400, // en-LS -> en-001
		0x656E4D47: 0x656E8400, // en-MG -> en-001
		0x656E4D4F: 0x656E8400, // en-MO -> en-001
		0x656E4D53: 0x656E8400, // en-MS -> en-001
		0x656E4D54: 0x656E8400, // en-MT -> en-001
		0x656E4D55: 0x656E8400, // en-MU -> en-001
		0x656E4D56: 0x656E8400, // en-MV -> en-001
		0x656E4D57: 0x656E8400, // en-MW -> en-001
		0x656E4D59: 0x656E8400, // en-MY -> en-001
		0x656E4E41: 0x656E8400, // en-NA -> en-001
		0x656E4E46: 0x656E8400, // en-NF -> en-001
		0x656E4E47: 0x656E8400, // en-NG -> en-001
		0x656E4E4C: 0x656E80A1, // en-NL -> en-150
		0x656E4E52: 0x656E8400, // en-NR -> en-001
		0x656E4E55: 0x656E8400, // en-NU -> en-001
		0x656E4E5A: 0x656E8400, // en-NZ -> en-001
		0x656E5047: 0x656E8400, // en-PG -> en-001
		0x656E504B: 0x656E8400, // en-PK -> en-001
		0x656E504E: 0x656E8400, // en-PN -> en-001
		0x656E5057: 0x656E8400, // en-PW -> en-001
		0x656E5257: 0x656E8400, // en-RW -> en-001
		0x656E5342: 0x656E8400, // en-SB -> en-001
		0x656E5343: 0x656E8400, // en-SC -> en-001
		0x656E5344: 0x656E8400, // en-SD -> en-001
		0x656E5345: 0x656E80A1, // en-SE -> en-150
		0x656E5347: 0x656E8400, // en-SG -> en-001
		0x656E5348: 0x656E8400, // en-SH -> en-001
		0x656E5349: 0x656E80A1, // en-SI -> en-150
		0x656E534C: 0x656E8400, // en-SL -> en-001
		0x656E5353: 0x656E8400, // en-SS -> en-001
		0x656E5358: 0x656E8400, // en-SX -> en-001
		0x656E535A: 0x656E8400, // en-SZ -> en-001
		0x656E5443: 0x656E8400, // en-TC -> en-001
		0x656E544B: 0x656E8400, // en-TK -> en-001
		0x656E544F: 0x656E8400, // en-TO -> en-001
		0x656E5454: 0x656E8400, // en-TT -> en-001
		0x656E5456: 0x656E8400, // en-TV -> en-001
		0x656E545A: 0x656E8400, // en-TZ -> en-001
		0x656E5547: 0x656E8400, // en-UG -> en-001
		0x656E5643: 0x656E8400, // en-VC -> en-001
		0x656E5647: 0x656E8400, // en-VG -> en-001
		0x656E5655: 0x656E8400, // en-VU -> en-001
		0x656E5753: 0x656E8400, // en-WS -> en-001
		0x656E5A41: 0x656E8400, // en-ZA -> en-001
		0x656E5A4D: 0x656E8400, // en-ZM -> en-001
		0x656E5A57: 0x656E8400, // en-ZW -> en-001
		0x65734152: 0x6573A424, // es-AR -> es-419
		0x6573424F: 0x6573A424, // es-BO -> es-419
		0x65734252: 0x6573A424, // es-BR -> es-419
		0x6573425A: 0x6573A424, // es-BZ -> es-419
		0x6573434C: 0x6573A424, // es-CL -> es-419
		0x6573434F: 0x6573A424, // es-CO -> es-419
		0x65734352: 0x6573A424, // es-CR -> es-419
		0x65734355: 0x6573A424, // es-CU -> es-419
		0x6573444F: 0x6573A424, // es-DO -> es-419
		0x65734543: 0x6573A424, // es-EC -> es-419
		0x65734754: 0x6573A424, // es-GT -> es-419
		0x6573484E: 0x6573A424, // es-HN -> es-419
		0x65734D58: 0x6573A424, // es-MX -> es-419
		0x65734E49: 0x6573A424, // es-NI -> es-419
		0x65735041: 0x6573A424, // es-PA -> es-419
		0x65735045: 0x6573A424, // es-PE -> es-419
		0x65735052: 0x6573A424, // es-PR -> es-419
		0x65735059: 0x6573A424, // es-PY -> es-419
		0x65735356: 0x6573A424, // es-SV -> es-419
		0x65735553: 0x6573A424, // es-US -> es-419
		0x65735559: 0x6573A424, // es-UY -> es-419
		0x65735645: 0x6573A424, // es-VE -> es-419
		0x6E620000: 0x6E6F0000, // nb -> no
		0x6E6E0000: 0x6E6F0000, // nn -> no
		0x7074414F: 0x70745054, // pt-AO -> pt-PT
		0x70744348: 0x70745054, // pt-CH -> pt-PT
		0x70744356: 0x70745054, // pt-CV -> pt-PT
		0x70744751: 0x70745054, // pt-GQ -> pt-PT
		0x70744757: 0x70745054, // pt-GW -> pt-PT
		0x70744C55: 0x70745054, // pt-LU -> pt-PT
		0x70744D4F: 0x70745054, // pt-MO -> pt-PT
		0x70744D5A: 0x70745054, // pt-MZ -> pt-PT
		0x70745354: 0x70745054, // pt-ST -> pt-PT
		0x7074544C: 0x70745054, // pt-TL -> pt-PT
	}
}

func bParents() map[uint32]uint32 {
	return map[uint32]uint32{
		0x61725842: 0x61729420, // ar-XB -> ar-015
	}
}

type scriptParent struct {
	script [4]uint8
	map_   map[uint32]uint32
}

func scriptParents() []scriptParent {
	return []scriptParent{
		{[4]uint8{'A', 'r', 'a', 'b'}, arabParents()},
		{[4]uint8{'D', 'e', 'v', 'a'}, devaParents()},
		{[4]uint8{'H', 'a', 'n', 't'}, hantParents()},
		{[4]uint8{'L', 'a', 't', 'n'}, latnParents()},
		{[4]uint8{'~', '~', '~', 'B'}, bParents()},
	}
}
