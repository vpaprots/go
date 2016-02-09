// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s390x

// This file contains utility functions for use when
// assembling vector instructions.

// vop returns the opcode, element size and condition
// setting for the given (possibly extended) mnemonic.
func vop(as int16) (opcode, es, cs uint32) {
	switch as {
	default:
		return 0, 0, 0
	case AVA:
		return OP_VA, 0, 0
	case AVAB:
		return OP_VA, 0, 0
	case AVAH:
		return OP_VA, 1, 0
	case AVAF:
		return OP_VA, 2, 0
	case AVAG:
		return OP_VA, 3, 0
	case AVAQ:
		return OP_VA, 4, 0
	case AVACC:
		return OP_VACC, 0, 0
	case AVACCB:
		return OP_VACC, 0, 0
	case AVACCH:
		return OP_VACC, 1, 0
	case AVACCF:
		return OP_VACC, 2, 0
	case AVACCG:
		return OP_VACC, 3, 0
	case AVACCQ:
		return OP_VACC, 4, 0
	case AVAC:
		return OP_VAC, 0, 0
	case AVACQ:
		return OP_VAC, 4, 0
	case AVACCC:
		return OP_VACCC, 0, 0
	case AVACCCQ:
		return OP_VACCC, 4, 0
	case AVN:
		return OP_VN, 0, 0
	case AVNC:
		return OP_VNC, 0, 0
	case AVAVG:
		return OP_VAVG, 0, 0
	case AVAVGB:
		return OP_VAVG, 0, 0
	case AVAVGH:
		return OP_VAVG, 1, 0
	case AVAVGF:
		return OP_VAVG, 2, 0
	case AVAVGG:
		return OP_VAVG, 3, 0
	case AVAVGL:
		return OP_VAVGL, 0, 0
	case AVAVGLB:
		return OP_VAVGL, 0, 0
	case AVAVGLH:
		return OP_VAVGL, 1, 0
	case AVAVGLF:
		return OP_VAVGL, 2, 0
	case AVAVGLG:
		return OP_VAVGL, 3, 0
	case AVCKSM:
		return OP_VCKSM, 0, 0
	case AVCEQ:
		return OP_VCEQ, 0, 0
	case AVCEQB:
		return OP_VCEQ, 0, 0
	case AVCEQH:
		return OP_VCEQ, 1, 0
	case AVCEQF:
		return OP_VCEQ, 2, 0
	case AVCEQG:
		return OP_VCEQ, 3, 0
	case AVCEQBS:
		return OP_VCEQ, 0, 1
	case AVCEQHS:
		return OP_VCEQ, 1, 1
	case AVCEQFS:
		return OP_VCEQ, 2, 1
	case AVCEQGS:
		return OP_VCEQ, 3, 1
	case AVCH:
		return OP_VCH, 0, 0
	case AVCHB:
		return OP_VCH, 0, 0
	case AVCHH:
		return OP_VCH, 1, 0
	case AVCHF:
		return OP_VCH, 2, 0
	case AVCHG:
		return OP_VCH, 3, 0
	case AVCHBS:
		return OP_VCH, 0, 1
	case AVCHHS:
		return OP_VCH, 1, 1
	case AVCHFS:
		return OP_VCH, 2, 1
	case AVCHGS:
		return OP_VCH, 3, 1
	case AVCHL:
		return OP_VCHL, 0, 0
	case AVCHLB:
		return OP_VCHL, 0, 0
	case AVCHLH:
		return OP_VCHL, 1, 0
	case AVCHLF:
		return OP_VCHL, 2, 0
	case AVCHLG:
		return OP_VCHL, 3, 0
	case AVCHLBS:
		return OP_VCHL, 0, 1
	case AVCHLHS:
		return OP_VCHL, 1, 1
	case AVCHLFS:
		return OP_VCHL, 2, 1
	case AVCHLGS:
		return OP_VCHL, 3, 1
	case AVCLZ:
		return OP_VCLZ, 0, 0
	case AVCLZB:
		return OP_VCLZ, 0, 0
	case AVCLZH:
		return OP_VCLZ, 1, 0
	case AVCLZF:
		return OP_VCLZ, 2, 0
	case AVCLZG:
		return OP_VCLZ, 3, 0
	case AVCTZ:
		return OP_VCTZ, 0, 0
	case AVCTZB:
		return OP_VCTZ, 0, 0
	case AVCTZH:
		return OP_VCTZ, 1, 0
	case AVCTZF:
		return OP_VCTZ, 2, 0
	case AVCTZG:
		return OP_VCTZ, 3, 0
	case AVEC:
		return OP_VEC, 0, 0
	case AVECB:
		return OP_VEC, 0, 0
	case AVECH:
		return OP_VEC, 1, 0
	case AVECF:
		return OP_VEC, 2, 0
	case AVECG:
		return OP_VEC, 3, 0
	case AVECL:
		return OP_VECL, 0, 0
	case AVECLB:
		return OP_VECL, 0, 0
	case AVECLH:
		return OP_VECL, 1, 0
	case AVECLF:
		return OP_VECL, 2, 0
	case AVECLG:
		return OP_VECL, 3, 0
	case AVERIM:
		return OP_VERIM, 0, 0
	case AVERIMB:
		return OP_VERIM, 0, 0
	case AVERIMH:
		return OP_VERIM, 1, 0
	case AVERIMF:
		return OP_VERIM, 2, 0
	case AVERIMG:
		return OP_VERIM, 3, 0
	case AVERLL:
		return OP_VERLL, 0, 0
	case AVERLLB:
		return OP_VERLL, 0, 0
	case AVERLLH:
		return OP_VERLL, 1, 0
	case AVERLLF:
		return OP_VERLL, 2, 0
	case AVERLLG:
		return OP_VERLL, 3, 0
	case AVERLLV:
		return OP_VERLLV, 0, 0
	case AVERLLVB:
		return OP_VERLLV, 0, 0
	case AVERLLVH:
		return OP_VERLLV, 1, 0
	case AVERLLVF:
		return OP_VERLLV, 2, 0
	case AVERLLVG:
		return OP_VERLLV, 3, 0
	case AVESLV:
		return OP_VESLV, 0, 0
	case AVESLVB:
		return OP_VESLV, 0, 0
	case AVESLVH:
		return OP_VESLV, 1, 0
	case AVESLVF:
		return OP_VESLV, 2, 0
	case AVESLVG:
		return OP_VESLV, 3, 0
	case AVESL:
		return OP_VESL, 0, 0
	case AVESLB:
		return OP_VESL, 0, 0
	case AVESLH:
		return OP_VESL, 1, 0
	case AVESLF:
		return OP_VESL, 2, 0
	case AVESLG:
		return OP_VESL, 3, 0
	case AVESRA:
		return OP_VESRA, 0, 0
	case AVESRAB:
		return OP_VESRA, 0, 0
	case AVESRAH:
		return OP_VESRA, 1, 0
	case AVESRAF:
		return OP_VESRA, 2, 0
	case AVESRAG:
		return OP_VESRA, 3, 0
	case AVESRAV:
		return OP_VESRAV, 0, 0
	case AVESRAVB:
		return OP_VESRAV, 0, 0
	case AVESRAVH:
		return OP_VESRAV, 1, 0
	case AVESRAVF:
		return OP_VESRAV, 2, 0
	case AVESRAVG:
		return OP_VESRAV, 3, 0
	case AVESRL:
		return OP_VESRL, 0, 0
	case AVESRLB:
		return OP_VESRL, 0, 0
	case AVESRLH:
		return OP_VESRL, 1, 0
	case AVESRLF:
		return OP_VESRL, 2, 0
	case AVESRLG:
		return OP_VESRL, 3, 0
	case AVESRLV:
		return OP_VESRLV, 0, 0
	case AVESRLVB:
		return OP_VESRLV, 0, 0
	case AVESRLVH:
		return OP_VESRLV, 1, 0
	case AVESRLVF:
		return OP_VESRLV, 2, 0
	case AVESRLVG:
		return OP_VESRLV, 3, 0
	case AVX:
		return OP_VX, 0, 0
	case AVFAE:
		return OP_VFAE, 0, 0
	case AVFAEB:
		return OP_VFAE, 0, 0
	case AVFAEH:
		return OP_VFAE, 1, 0
	case AVFAEF:
		return OP_VFAE, 2, 0
	case AVFAEBS:
		return OP_VFAE, 0, 1
	case AVFAEHS:
		return OP_VFAE, 1, 1
	case AVFAEFS:
		return OP_VFAE, 2, 1
	case AVFAEZB:
		return OP_VFAE, 0, 2
	case AVFAEZH:
		return OP_VFAE, 1, 2
	case AVFAEZF:
		return OP_VFAE, 2, 2
	case AVFAEZBS:
		return OP_VFAE, 0, 3
	case AVFAEZHS:
		return OP_VFAE, 1, 3
	case AVFAEZFS:
		return OP_VFAE, 2, 3
	case AVFEE:
		return OP_VFEE, 0, 0
	case AVFEEB:
		return OP_VFEE, 0, 0
	case AVFEEH:
		return OP_VFEE, 1, 0
	case AVFEEF:
		return OP_VFEE, 2, 0
	case AVFEEBS:
		return OP_VFEE, 0, 1
	case AVFEEHS:
		return OP_VFEE, 1, 1
	case AVFEEFS:
		return OP_VFEE, 2, 1
	case AVFEEZB:
		return OP_VFEE, 0, 2
	case AVFEEZH:
		return OP_VFEE, 1, 2
	case AVFEEZF:
		return OP_VFEE, 2, 2
	case AVFEEZBS:
		return OP_VFEE, 0, 3
	case AVFEEZHS:
		return OP_VFEE, 1, 3
	case AVFEEZFS:
		return OP_VFEE, 2, 3
	case AVFENE:
		return OP_VFENE, 0, 0
	case AVFENEB:
		return OP_VFENE, 0, 0
	case AVFENEH:
		return OP_VFENE, 1, 0
	case AVFENEF:
		return OP_VFENE, 2, 0
	case AVFENEBS:
		return OP_VFENE, 0, 1
	case AVFENEHS:
		return OP_VFENE, 1, 1
	case AVFENEFS:
		return OP_VFENE, 2, 1
	case AVFENEZB:
		return OP_VFENE, 0, 2
	case AVFENEZH:
		return OP_VFENE, 1, 2
	case AVFENEZF:
		return OP_VFENE, 2, 2
	case AVFENEZBS:
		return OP_VFENE, 0, 3
	case AVFENEZHS:
		return OP_VFENE, 1, 3
	case AVFENEZFS:
		return OP_VFENE, 2, 3
	case AVFA:
		return OP_VFA, 0, 0
	case AVFADB:
		return OP_VFA, 3, 0
	case AWFADB:
		return OP_VFA, 3, 0
	case AWFK:
		return OP_WFK, 0, 0
	case AWFKDB:
		return OP_WFK, 3, 0
	case AVFCE:
		return OP_VFCE, 0, 0
	case AVFCEDB:
		return OP_VFCE, 3, 0
	case AVFCEDBS:
		return OP_VFCE, 3, 1
	case AWFCEDB:
		return OP_VFCE, 3, 0
	case AWFCEDBS:
		return OP_VFCE, 3, 1
	case AVFCH:
		return OP_VFCH, 0, 0
	case AVFCHDB:
		return OP_VFCH, 3, 0
	case AVFCHDBS:
		return OP_VFCH, 3, 1
	case AWFCHDB:
		return OP_VFCH, 3, 0
	case AWFCHDBS:
		return OP_VFCH, 3, 1
	case AVFCHE:
		return OP_VFCHE, 0, 0
	case AVFCHEDB:
		return OP_VFCHE, 3, 0
	case AVFCHEDBS:
		return OP_VFCHE, 3, 1
	case AWFCHEDB:
		return OP_VFCHE, 3, 0
	case AWFCHEDBS:
		return OP_VFCHE, 3, 1
	case AWFC:
		return OP_WFC, 0, 0
	case AWFCDB:
		return OP_WFC, 3, 0
	case AVCDG:
		return OP_VCDG, 0, 0
	case AVCDGB:
		return OP_VCDG, 3, 0
	case AWCDGB:
		return OP_VCDG, 3, 0
	case AVCDLG:
		return OP_VCDLG, 0, 0
	case AVCDLGB:
		return OP_VCDLG, 3, 0
	case AWCDLGB:
		return OP_VCDLG, 3, 0
	case AVCGD:
		return OP_VCGD, 0, 0
	case AVCGDB:
		return OP_VCGD, 3, 0
	case AWCGDB:
		return OP_VCGD, 3, 0
	case AVCLGD:
		return OP_VCLGD, 0, 0
	case AVCLGDB:
		return OP_VCLGD, 3, 0
	case AWCLGDB:
		return OP_VCLGD, 3, 0
	case AVFD:
		return OP_VFD, 0, 0
	case AVFDDB:
		return OP_VFD, 3, 0
	case AWFDDB:
		return OP_VFD, 3, 0
	case AVLDE:
		return OP_VLDE, 0, 0
	case AVLDEB:
		return OP_VLDE, 2, 0
	case AWLDEB:
		return OP_VLDE, 2, 0
	case AVLED:
		return OP_VLED, 0, 0
	case AVLEDB:
		return OP_VLED, 3, 0
	case AWLEDB:
		return OP_VLED, 3, 0
	case AVFM:
		return OP_VFM, 0, 0
	case AVFMDB:
		return OP_VFM, 3, 0
	case AWFMDB:
		return OP_VFM, 3, 0
	case AVFMA:
		return OP_VFMA, 0, 0
	case AVFMADB:
		return OP_VFMA, 3, 0
	case AWFMADB:
		return OP_VFMA, 3, 0
	case AVFMS:
		return OP_VFMS, 0, 0
	case AVFMSDB:
		return OP_VFMS, 3, 0
	case AWFMSDB:
		return OP_VFMS, 3, 0
	case AVFPSO:
		return OP_VFPSO, 0, 0
	case AVFPSODB:
		return OP_VFPSO, 3, 0
	case AWFPSODB:
		return OP_VFPSO, 3, 0
	case AVFLCDB:
		return OP_VFPSO, 3, 0
	case AWFLCDB:
		return OP_VFPSO, 3, 0
	case AVFLNDB:
		return OP_VFPSO, 3, 1
	case AWFLNDB:
		return OP_VFPSO, 3, 1
	case AVFLPDB:
		return OP_VFPSO, 3, 2
	case AWFLPDB:
		return OP_VFPSO, 3, 2
	case AVFSQ:
		return OP_VFSQ, 0, 0
	case AVFSQDB:
		return OP_VFSQ, 3, 0
	case AWFSQDB:
		return OP_VFSQ, 3, 0
	case AVFS:
		return OP_VFS, 0, 0
	case AVFSDB:
		return OP_VFS, 3, 0
	case AWFSDB:
		return OP_VFS, 3, 0
	case AVFTCI:
		return OP_VFTCI, 0, 0
	case AVFTCIDB:
		return OP_VFTCI, 3, 0
	case AWFTCIDB:
		return OP_VFTCI, 3, 0
	case AVGFM:
		return OP_VGFM, 0, 0
	case AVGFMB:
		return OP_VGFM, 0, 0
	case AVGFMH:
		return OP_VGFM, 1, 0
	case AVGFMF:
		return OP_VGFM, 2, 0
	case AVGFMG:
		return OP_VGFM, 3, 0
	case AVGFMA:
		return OP_VGFMA, 0, 0
	case AVGFMAB:
		return OP_VGFMA, 0, 0
	case AVGFMAH:
		return OP_VGFMA, 1, 0
	case AVGFMAF:
		return OP_VGFMA, 2, 0
	case AVGFMAG:
		return OP_VGFMA, 3, 0
	case AVGEF:
		return OP_VGEF, 0, 0
	case AVGEG:
		return OP_VGEG, 0, 0
	case AVGBM:
		return OP_VGBM, 0, 0
	case AVZERO:
		return OP_VGBM, 0, 0
	case AVONE:
		return OP_VGBM, 0, 0
	case AVGM:
		return OP_VGM, 0, 0
	case AVGMB:
		return OP_VGM, 0, 0
	case AVGMH:
		return OP_VGM, 1, 0
	case AVGMF:
		return OP_VGM, 2, 0
	case AVGMG:
		return OP_VGM, 3, 0
	case AVISTR:
		return OP_VISTR, 0, 0
	case AVISTRB:
		return OP_VISTR, 0, 0
	case AVISTRH:
		return OP_VISTR, 1, 0
	case AVISTRF:
		return OP_VISTR, 2, 0
	case AVISTRBS:
		return OP_VISTR, 0, 1
	case AVISTRHS:
		return OP_VISTR, 1, 1
	case AVISTRFS:
		return OP_VISTR, 2, 1
	case AVL:
		return OP_VL, 0, 0
	case AVLR:
		return OP_VLR, 0, 0
	case AVLREP:
		return OP_VLREP, 0, 0
	case AVLREPB:
		return OP_VLREP, 0, 0
	case AVLREPH:
		return OP_VLREP, 1, 0
	case AVLREPF:
		return OP_VLREP, 2, 0
	case AVLREPG:
		return OP_VLREP, 3, 0
	case AVLC:
		return OP_VLC, 0, 0
	case AVLCB:
		return OP_VLC, 0, 0
	case AVLCH:
		return OP_VLC, 1, 0
	case AVLCF:
		return OP_VLC, 2, 0
	case AVLCG:
		return OP_VLC, 3, 0
	case AVLEH:
		return OP_VLEH, 0, 0
	case AVLEF:
		return OP_VLEF, 0, 0
	case AVLEG:
		return OP_VLEG, 0, 0
	case AVLEB:
		return OP_VLEB, 0, 0
	case AVLEIH:
		return OP_VLEIH, 0, 0
	case AVLEIF:
		return OP_VLEIF, 0, 0
	case AVLEIG:
		return OP_VLEIG, 0, 0
	case AVLEIB:
		return OP_VLEIB, 0, 0
	case AVFI:
		return OP_VFI, 0, 0
	case AVFIDB:
		return OP_VFI, 3, 0
	case AWFIDB:
		return OP_VFI, 3, 0
	case AVLGV:
		return OP_VLGV, 0, 0
	case AVLGVB:
		return OP_VLGV, 0, 0
	case AVLGVH:
		return OP_VLGV, 1, 0
	case AVLGVF:
		return OP_VLGV, 2, 0
	case AVLGVG:
		return OP_VLGV, 3, 0
	case AVLLEZ:
		return OP_VLLEZ, 0, 0
	case AVLLEZB:
		return OP_VLLEZ, 0, 0
	case AVLLEZH:
		return OP_VLLEZ, 1, 0
	case AVLLEZF:
		return OP_VLLEZ, 2, 0
	case AVLLEZG:
		return OP_VLLEZ, 3, 0
	case AVLM:
		return OP_VLM, 0, 0
	case AVLP:
		return OP_VLP, 0, 0
	case AVLPB:
		return OP_VLP, 0, 0
	case AVLPH:
		return OP_VLP, 1, 0
	case AVLPF:
		return OP_VLP, 2, 0
	case AVLPG:
		return OP_VLP, 3, 0
	case AVLBB:
		return OP_VLBB, 0, 0
	case AVLVG:
		return OP_VLVG, 0, 0
	case AVLVGB:
		return OP_VLVG, 0, 0
	case AVLVGH:
		return OP_VLVG, 1, 0
	case AVLVGF:
		return OP_VLVG, 2, 0
	case AVLVGG:
		return OP_VLVG, 3, 0
	case AVLVGP:
		return OP_VLVGP, 0, 0
	case AVLL:
		return OP_VLL, 0, 0
	case AVMX:
		return OP_VMX, 0, 0
	case AVMXB:
		return OP_VMX, 0, 0
	case AVMXH:
		return OP_VMX, 1, 0
	case AVMXF:
		return OP_VMX, 2, 0
	case AVMXG:
		return OP_VMX, 3, 0
	case AVMXL:
		return OP_VMXL, 0, 0
	case AVMXLB:
		return OP_VMXL, 0, 0
	case AVMXLH:
		return OP_VMXL, 1, 0
	case AVMXLF:
		return OP_VMXL, 2, 0
	case AVMXLG:
		return OP_VMXL, 3, 0
	case AVMRH:
		return OP_VMRH, 0, 0
	case AVMRHB:
		return OP_VMRH, 0, 0
	case AVMRHH:
		return OP_VMRH, 1, 0
	case AVMRHF:
		return OP_VMRH, 2, 0
	case AVMRHG:
		return OP_VMRH, 3, 0
	case AVMRL:
		return OP_VMRL, 0, 0
	case AVMRLB:
		return OP_VMRL, 0, 0
	case AVMRLH:
		return OP_VMRL, 1, 0
	case AVMRLF:
		return OP_VMRL, 2, 0
	case AVMRLG:
		return OP_VMRL, 3, 0
	case AVMN:
		return OP_VMN, 0, 0
	case AVMNB:
		return OP_VMN, 0, 0
	case AVMNH:
		return OP_VMN, 1, 0
	case AVMNF:
		return OP_VMN, 2, 0
	case AVMNG:
		return OP_VMN, 3, 0
	case AVMNL:
		return OP_VMNL, 0, 0
	case AVMNLB:
		return OP_VMNL, 0, 0
	case AVMNLH:
		return OP_VMNL, 1, 0
	case AVMNLF:
		return OP_VMNL, 2, 0
	case AVMNLG:
		return OP_VMNL, 3, 0
	case AVMAE:
		return OP_VMAE, 0, 0
	case AVMAEB:
		return OP_VMAE, 0, 0
	case AVMAEH:
		return OP_VMAE, 1, 0
	case AVMAEF:
		return OP_VMAE, 2, 0
	case AVMAH:
		return OP_VMAH, 0, 0
	case AVMAHB:
		return OP_VMAH, 0, 0
	case AVMAHH:
		return OP_VMAH, 1, 0
	case AVMAHF:
		return OP_VMAH, 2, 0
	case AVMALE:
		return OP_VMALE, 0, 0
	case AVMALEB:
		return OP_VMALE, 0, 0
	case AVMALEH:
		return OP_VMALE, 1, 0
	case AVMALEF:
		return OP_VMALE, 2, 0
	case AVMALH:
		return OP_VMALH, 0, 0
	case AVMALHB:
		return OP_VMALH, 0, 0
	case AVMALHH:
		return OP_VMALH, 1, 0
	case AVMALHF:
		return OP_VMALH, 2, 0
	case AVMALO:
		return OP_VMALO, 0, 0
	case AVMALOB:
		return OP_VMALO, 0, 0
	case AVMALOH:
		return OP_VMALO, 1, 0
	case AVMALOF:
		return OP_VMALO, 2, 0
	case AVMAL:
		return OP_VMAL, 0, 0
	case AVMALB:
		return OP_VMAL, 0, 0
	case AVMALHW:
		return OP_VMAL, 1, 0
	case AVMALF:
		return OP_VMAL, 2, 0
	case AVMAO:
		return OP_VMAO, 0, 0
	case AVMAOB:
		return OP_VMAO, 0, 0
	case AVMAOH:
		return OP_VMAO, 1, 0
	case AVMAOF:
		return OP_VMAO, 2, 0
	case AVME:
		return OP_VME, 0, 0
	case AVMEB:
		return OP_VME, 0, 0
	case AVMEH:
		return OP_VME, 1, 0
	case AVMEF:
		return OP_VME, 2, 0
	case AVMH:
		return OP_VMH, 0, 0
	case AVMHB:
		return OP_VMH, 0, 0
	case AVMHH:
		return OP_VMH, 1, 0
	case AVMHF:
		return OP_VMH, 2, 0
	case AVMLE:
		return OP_VMLE, 0, 0
	case AVMLEB:
		return OP_VMLE, 0, 0
	case AVMLEH:
		return OP_VMLE, 1, 0
	case AVMLEF:
		return OP_VMLE, 2, 0
	case AVMLH:
		return OP_VMLH, 0, 0
	case AVMLHB:
		return OP_VMLH, 0, 0
	case AVMLHH:
		return OP_VMLH, 1, 0
	case AVMLHF:
		return OP_VMLH, 2, 0
	case AVMLO:
		return OP_VMLO, 0, 0
	case AVMLOB:
		return OP_VMLO, 0, 0
	case AVMLOH:
		return OP_VMLO, 1, 0
	case AVMLOF:
		return OP_VMLO, 2, 0
	case AVML:
		return OP_VML, 0, 0
	case AVMLB:
		return OP_VML, 0, 0
	case AVMLHW:
		return OP_VML, 1, 0
	case AVMLF:
		return OP_VML, 2, 0
	case AVMO:
		return OP_VMO, 0, 0
	case AVMOB:
		return OP_VMO, 0, 0
	case AVMOH:
		return OP_VMO, 1, 0
	case AVMOF:
		return OP_VMO, 2, 0
	case AVNO:
		return OP_VNO, 0, 0
	case AVNOT:
		return OP_VNO, 0, 0
	case AVO:
		return OP_VO, 0, 0
	case AVPK:
		return OP_VPK, 0, 0
	case AVPKH:
		return OP_VPK, 1, 0
	case AVPKF:
		return OP_VPK, 2, 0
	case AVPKG:
		return OP_VPK, 3, 0
	case AVPKLS:
		return OP_VPKLS, 0, 0
	case AVPKLSH:
		return OP_VPKLS, 1, 0
	case AVPKLSF:
		return OP_VPKLS, 2, 0
	case AVPKLSG:
		return OP_VPKLS, 3, 0
	case AVPKLSHS:
		return OP_VPKLS, 1, 1
	case AVPKLSFS:
		return OP_VPKLS, 2, 1
	case AVPKLSGS:
		return OP_VPKLS, 3, 1
	case AVPKS:
		return OP_VPKS, 0, 0
	case AVPKSH:
		return OP_VPKS, 1, 0
	case AVPKSF:
		return OP_VPKS, 2, 0
	case AVPKSG:
		return OP_VPKS, 3, 0
	case AVPKSHS:
		return OP_VPKS, 1, 1
	case AVPKSFS:
		return OP_VPKS, 2, 1
	case AVPKSGS:
		return OP_VPKS, 3, 1
	case AVPERM:
		return OP_VPERM, 0, 0
	case AVPDI:
		return OP_VPDI, 0, 0
	case AVPOPCT:
		return OP_VPOPCT, 0, 0
	case AVREP:
		return OP_VREP, 0, 0
	case AVREPB:
		return OP_VREP, 0, 0
	case AVREPH:
		return OP_VREP, 1, 0
	case AVREPF:
		return OP_VREP, 2, 0
	case AVREPG:
		return OP_VREP, 3, 0
	case AVREPI:
		return OP_VREPI, 0, 0
	case AVREPIB:
		return OP_VREPI, 0, 0
	case AVREPIH:
		return OP_VREPI, 1, 0
	case AVREPIF:
		return OP_VREPI, 2, 0
	case AVREPIG:
		return OP_VREPI, 3, 0
	case AVSCEF:
		return OP_VSCEF, 0, 0
	case AVSCEG:
		return OP_VSCEG, 0, 0
	case AVSEL:
		return OP_VSEL, 0, 0
	case AVSL:
		return OP_VSL, 0, 0
	case AVSLB:
		return OP_VSLB, 0, 0
	case AVSLDB:
		return OP_VSLDB, 0, 0
	case AVSRA:
		return OP_VSRA, 0, 0
	case AVSRAB:
		return OP_VSRAB, 0, 0
	case AVSRL:
		return OP_VSRL, 0, 0
	case AVSRLB:
		return OP_VSRLB, 0, 0
	case AVSEG:
		return OP_VSEG, 0, 0
	case AVSEGB:
		return OP_VSEG, 0, 0
	case AVSEGH:
		return OP_VSEG, 1, 0
	case AVSEGF:
		return OP_VSEG, 2, 0
	case AVST:
		return OP_VST, 0, 0
	case AVSTEH:
		return OP_VSTEH, 0, 0
	case AVSTEF:
		return OP_VSTEF, 0, 0
	case AVSTEG:
		return OP_VSTEG, 0, 0
	case AVSTEB:
		return OP_VSTEB, 0, 0
	case AVSTM:
		return OP_VSTM, 0, 0
	case AVSTL:
		return OP_VSTL, 0, 0
	case AVSTRC:
		return OP_VSTRC, 0, 0
	case AVSTRCB:
		return OP_VSTRC, 0, 0
	case AVSTRCH:
		return OP_VSTRC, 1, 0
	case AVSTRCF:
		return OP_VSTRC, 2, 0
	case AVSTRCBS:
		return OP_VSTRC, 0, 1
	case AVSTRCHS:
		return OP_VSTRC, 1, 1
	case AVSTRCFS:
		return OP_VSTRC, 2, 1
	case AVSTRCZB:
		return OP_VSTRC, 0, 2
	case AVSTRCZH:
		return OP_VSTRC, 1, 2
	case AVSTRCZF:
		return OP_VSTRC, 2, 2
	case AVSTRCZBS:
		return OP_VSTRC, 0, 3
	case AVSTRCZHS:
		return OP_VSTRC, 1, 3
	case AVSTRCZFS:
		return OP_VSTRC, 2, 3
	case AVS:
		return OP_VS, 0, 0
	case AVSB:
		return OP_VS, 0, 0
	case AVSH:
		return OP_VS, 1, 0
	case AVSF:
		return OP_VS, 2, 0
	case AVSG:
		return OP_VS, 3, 0
	case AVSQ:
		return OP_VS, 4, 0
	case AVSCBI:
		return OP_VSCBI, 0, 0
	case AVSCBIB:
		return OP_VSCBI, 0, 0
	case AVSCBIH:
		return OP_VSCBI, 1, 0
	case AVSCBIF:
		return OP_VSCBI, 2, 0
	case AVSCBIG:
		return OP_VSCBI, 3, 0
	case AVSCBIQ:
		return OP_VSCBI, 4, 0
	case AVSBCBI:
		return OP_VSBCBI, 0, 0
	case AVSBCBIQ:
		return OP_VSBCBI, 4, 0
	case AVSBI:
		return OP_VSBI, 0, 0
	case AVSBIQ:
		return OP_VSBI, 4, 0
	case AVSUMG:
		return OP_VSUMG, 0, 0
	case AVSUMGH:
		return OP_VSUMG, 1, 0
	case AVSUMGF:
		return OP_VSUMG, 2, 0
	case AVSUMQ:
		return OP_VSUMQ, 0, 0
	case AVSUMQF:
		return OP_VSUMQ, 1, 0
	case AVSUMQG:
		return OP_VSUMQ, 2, 0
	case AVSUM:
		return OP_VSUM, 0, 0
	case AVSUMB:
		return OP_VSUM, 0, 0
	case AVSUMH:
		return OP_VSUM, 1, 0
	case AVTM:
		return OP_VTM, 0, 0
	case AVUPH:
		return OP_VUPH, 0, 0
	case AVUPHB:
		return OP_VUPH, 0, 0
	case AVUPHH:
		return OP_VUPH, 1, 0
	case AVUPHF:
		return OP_VUPH, 2, 0
	case AVUPLH:
		return OP_VUPLH, 0, 0
	case AVUPLHB:
		return OP_VUPLH, 0, 0
	case AVUPLHH:
		return OP_VUPLH, 1, 0
	case AVUPLHF:
		return OP_VUPLH, 2, 0
	case AVUPLL:
		return OP_VUPLL, 0, 0
	case AVUPLLB:
		return OP_VUPLL, 0, 0
	case AVUPLLH:
		return OP_VUPLL, 1, 0
	case AVUPLLF:
		return OP_VUPLL, 2, 0
	case AVUPL:
		return OP_VUPL, 0, 0
	case AVUPLB:
		return OP_VUPL, 0, 0
	case AVUPLHW:
		return OP_VUPL, 1, 0
	case AVUPLF:
		return OP_VUPL, 2, 0
	}
}

// singleElementMask returns the single element mask bits required for the
// given instruction.
func singleElementMask(as int16) uint32 {
	switch as {
	case AWFADB,
		AWFK,
		AWFKDB,
		AWFCEDB,
		AWFCEDBS,
		AWFCHDB,
		AWFCHDBS,
		AWFCHEDB,
		AWFCHEDBS,
		AWFC,
		AWFCDB,
		AWCDGB,
		AWCDLGB,
		AWCGDB,
		AWCLGDB,
		AWFDDB,
		AWLDEB,
		AWLEDB,
		AWFMDB,
		AWFMADB,
		AWFMSDB,
		AWFPSODB,
		AWFLCDB,
		AWFLNDB,
		AWFLPDB,
		AWFSQDB,
		AWFSDB,
		AWFTCIDB,
		AWFIDB:
		return 8
	}
	return 0
}
