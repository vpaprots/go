// Derived from Inferno utils/6c/txt.c
// http://code.google.com/p/inferno-os/source/browse/utils/6c/txt.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package s390x

import (
	"cmd/compile/internal/big"
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/s390x"
	"fmt"
)

var resvd = []int{
	s390x.REGZERO,
	s390x.REGSP, // reserved for SP
	s390x.REGG,
	s390x.REGTMP, // REGTMP
	s390x.REGTMP2,
	s390x.FREGCVI,
	s390x.FREGZERO,
	s390x.FREGONE,
	s390x.FREGTWO,
//	s390x.FREGTMP,
}

/*
 * generate
 *	as $c, n
 */
func ginscon(as int, c int64, n2 *gc.Node) {
	var n1 gc.Node

	gc.Nodconst(&n1, gc.Types[gc.TINT64], c)

	if as != s390x.AMOVD && (c < -s390x.BIG || c > s390x.BIG) || n2.Op != gc.OREGISTER || as == s390x.AMULLD {
		// cannot have more than 16-bit of immediate in ADD, etc.
		// instead, MOV into register first.
		var ntmp gc.Node
		gc.Regalloc(&ntmp, gc.Types[gc.TINT64], nil)

		rawgins(s390x.AMOVD, &n1, &ntmp)
		rawgins(as, &ntmp, n2)
		gc.Regfree(&ntmp)
		return
	}

	rawgins(as, &n1, n2)
}

/*
 * generate
 *	as n, $c (CMP/CMPU)
 */
func ginscon2(as int, n2 *gc.Node, c int64) {
	var n1 gc.Node

	gc.Nodconst(&n1, gc.Types[gc.TINT64], c)

	switch as {
	default:
		gc.Fatal("ginscon2")

	case s390x.ACMP:
		if -s390x.BIG <= c && c <= s390x.BIG {
			rawgins(as, n2, &n1)
			return
		}

	case s390x.ACMPU:
		if 0 <= c && c <= 2*s390x.BIG {
			rawgins(as, n2, &n1)
			return
		}
	}

	// MOV n1 into register first
	var ntmp gc.Node
	gc.Regalloc(&ntmp, gc.Types[gc.TINT64], nil)

	rawgins(s390x.AMOVD, &n1, &ntmp)
	rawgins(as, n2, &ntmp)
	gc.Regfree(&ntmp)
}

func ginscmp(op int, t *gc.Type, n1, n2 *gc.Node, likely int) *obj.Prog {
	if gc.Isint[t.Etype] && n1.Op == gc.OLITERAL && n2.Op != gc.OLITERAL {
		// Reverse comparison to place constant last.
		op = gc.Brrev(op)
		n1, n2 = n2, n1
	}

	var r1, r2, g1, g2 gc.Node
	gc.Regalloc(&r1, t, n1)
	gc.Regalloc(&g1, n1.Type, &r1)
	gc.Cgen(n1, &g1)
	gmove(&g1, &r1)
	if gc.Isint[t.Etype] && gc.Isconst(n2, gc.CTINT) {
		ginscon2(optoas(gc.OCMP, t), &r1, n2.Int())
	} else {
		gc.Regalloc(&r2, t, n2)
		gc.Regalloc(&g2, n1.Type, &r2)
		gc.Cgen(n2, &g2)
		gmove(&g2, &r2)
		rawgins(optoas(gc.OCMP, t), &r1, &r2)
		gc.Regfree(&g2)
		gc.Regfree(&r2)
	}
	gc.Regfree(&g1)
	gc.Regfree(&r1)
	return gc.Gbranch(optoas(op, t), nil, likely)
}

// set up nodes representing 2^63
var (
	bigi         gc.Node
	bigf         gc.Node
	bignodes_did bool
)

func bignodes() {
	if bignodes_did {
		return
	}
	bignodes_did = true

	var i big.Int
	i.SetInt64(1)
	i.Lsh(&i, 63)

	gc.Nodconst(&bigi, gc.Types[gc.TUINT64], 0)
	bigi.SetBigInt(&i)

	bigi.Convconst(&bigf, gc.Types[gc.TFLOAT64])
}

/*
 * generate move:
 *	t = f
 * hard part is conversions.
 */
func gmove(f *gc.Node, t *gc.Node) {
	if gc.Debug['M'] != 0 {
		fmt.Printf("gmove %v -> %v\n", gc.Nconv(f, obj.FmtLong), gc.Nconv(t, obj.FmtLong))
	}

	ft := int(gc.Simsimtype(f.Type))
	tt := int(gc.Simsimtype(t.Type))
	cvt := (*gc.Type)(t.Type)

	if gc.Iscomplex[ft] || gc.Iscomplex[tt] {
		gc.Complexmove(f, t)
		return
	}

	// cannot have two memory operands
	var a int
	if gc.Ismem(f) && gc.Ismem(t) {
		goto hard
	}

	// convert constant to desired type
	if f.Op == gc.OLITERAL {
		var con gc.Node
		switch tt {
		default:
			f.Convconst(&con, t.Type)

		case gc.TINT32,
			gc.TINT16,
			gc.TINT8:
			var con gc.Node
			f.Convconst(&con, gc.Types[gc.TINT64])
			var r1 gc.Node
			gc.Regalloc(&r1, con.Type, t)
			gins(s390x.AMOVD, &con, &r1)
			gmove(&r1, t)
			gc.Regfree(&r1)
			return

		case gc.TUINT32,
			gc.TUINT16,
			gc.TUINT8:
			var con gc.Node
			f.Convconst(&con, gc.Types[gc.TUINT64])
			var r1 gc.Node
			gc.Regalloc(&r1, con.Type, t)
			gins(s390x.AMOVD, &con, &r1)
			gmove(&r1, t)
			gc.Regfree(&r1)
			return
		}

		f = &con
		ft = tt // so big switch will choose a simple mov

		// constants can't move directly to memory.
		if gc.Ismem(t) {
			goto hard
		}
	}

	// a float-to-int or int-to-float conversion requires the source operand in a register
	if gc.Ismem(f) && ((gc.Isfloat[ft] && gc.Isint[tt]) || (gc.Isint[ft] && gc.Isfloat[tt])) {
		cvt = (*gc.Type)(f.Type)
		goto hard
	}

	// a float32-to-float64 or float64-to-float32 conversion requires the source operand in a register
	if gc.Ismem(f) && gc.Isfloat[ft] && gc.Isfloat[tt] && (ft != tt) {
		cvt = (*gc.Type)(f.Type)
		goto hard
	}

	// float constants come from memory.
	//if(isfloat[tt])
	//	goto hard;

	// 64-bit immediates are also from memory.
	//if(isint[tt])
	//	goto hard;
	//// 64-bit immediates are really 32-bit sign-extended
	//// unless moving into a register.
	//if(isint[tt]) {
	//	if(mpcmpfixfix(con.val.u.xval, minintval[TINT32]) < 0)
	//		goto hard;
	//	if(mpcmpfixfix(con.val.u.xval, maxintval[TINT32]) > 0)
	//		goto hard;
	//}

	// value -> value copy, only one memory operand.
	// figure out the instruction to use.
	// break out of switch for one-instruction gins.
	// goto rdst for "destination must be register".
	// goto hard for "convert to cvt type first".
	// otherwise handle and return.

	switch uint32(ft)<<16 | uint32(tt) {
	default:
		gc.Fatal("gmove %v -> %v", gc.Tconv(f.Type, obj.FmtLong), gc.Tconv(t.Type, obj.FmtLong))

		/*
		 * integer copy and truncate
		 */
	case gc.TINT8<<16 | gc.TINT8, // same size
		gc.TUINT8<<16 | gc.TINT8,
		gc.TINT16<<16 | gc.TINT8,
		// truncate
		gc.TUINT16<<16 | gc.TINT8,
		gc.TINT32<<16 | gc.TINT8,
		gc.TUINT32<<16 | gc.TINT8,
		gc.TINT64<<16 | gc.TINT8,
		gc.TUINT64<<16 | gc.TINT8:
		a = s390x.AMOVB

	case gc.TINT8<<16 | gc.TUINT8, // same size
		gc.TUINT8<<16 | gc.TUINT8,
		gc.TINT16<<16 | gc.TUINT8,
		// truncate
		gc.TUINT16<<16 | gc.TUINT8,
		gc.TINT32<<16 | gc.TUINT8,
		gc.TUINT32<<16 | gc.TUINT8,
		gc.TINT64<<16 | gc.TUINT8,
		gc.TUINT64<<16 | gc.TUINT8:
		a = s390x.AMOVBZ

	case gc.TINT16<<16 | gc.TINT16, // same size
		gc.TUINT16<<16 | gc.TINT16,
		gc.TINT32<<16 | gc.TINT16,
		// truncate
		gc.TUINT32<<16 | gc.TINT16,
		gc.TINT64<<16 | gc.TINT16,
		gc.TUINT64<<16 | gc.TINT16:
		a = s390x.AMOVH

	case gc.TINT16<<16 | gc.TUINT16, // same size
		gc.TUINT16<<16 | gc.TUINT16,
		gc.TINT32<<16 | gc.TUINT16,
		// truncate
		gc.TUINT32<<16 | gc.TUINT16,
		gc.TINT64<<16 | gc.TUINT16,
		gc.TUINT64<<16 | gc.TUINT16:
		a = s390x.AMOVHZ

	case gc.TINT32<<16 | gc.TINT32, // same size
		gc.TUINT32<<16 | gc.TINT32,
		gc.TINT64<<16 | gc.TINT32,
		// truncate
		gc.TUINT64<<16 | gc.TINT32:
		a = s390x.AMOVW

	case gc.TINT32<<16 | gc.TUINT32, // same size
		gc.TUINT32<<16 | gc.TUINT32,
		gc.TINT64<<16 | gc.TUINT32,
		gc.TUINT64<<16 | gc.TUINT32:
		a = s390x.AMOVWZ

	case gc.TINT64<<16 | gc.TINT64, // same size
		gc.TINT64<<16 | gc.TUINT64,
		gc.TUINT64<<16 | gc.TINT64,
		gc.TUINT64<<16 | gc.TUINT64:
		a = s390x.AMOVD

		/*
		 * integer up-conversions
		 */
	case gc.TINT8<<16 | gc.TINT16, // sign extend int8
		gc.TINT8<<16 | gc.TUINT16,
		gc.TINT8<<16 | gc.TINT32,
		gc.TINT8<<16 | gc.TUINT32,
		gc.TINT8<<16 | gc.TINT64,
		gc.TINT8<<16 | gc.TUINT64:
		a = s390x.AMOVB

		goto rdst

	case gc.TUINT8<<16 | gc.TINT16, // zero extend uint8
		gc.TUINT8<<16 | gc.TUINT16,
		gc.TUINT8<<16 | gc.TINT32,
		gc.TUINT8<<16 | gc.TUINT32,
		gc.TUINT8<<16 | gc.TINT64,
		gc.TUINT8<<16 | gc.TUINT64:
		a = s390x.AMOVBZ

		goto rdst

	case gc.TINT16<<16 | gc.TINT32, // sign extend int16
		gc.TINT16<<16 | gc.TUINT32,
		gc.TINT16<<16 | gc.TINT64,
		gc.TINT16<<16 | gc.TUINT64:
		a = s390x.AMOVH

		goto rdst

	case gc.TUINT16<<16 | gc.TINT32, // zero extend uint16
		gc.TUINT16<<16 | gc.TUINT32,
		gc.TUINT16<<16 | gc.TINT64,
		gc.TUINT16<<16 | gc.TUINT64:
		a = s390x.AMOVHZ

		goto rdst

	case gc.TINT32<<16 | gc.TINT64, // sign extend int32
		gc.TINT32<<16 | gc.TUINT64:
		a = s390x.AMOVW

		goto rdst

	case gc.TUINT32<<16 | gc.TINT64, // zero extend uint32
		gc.TUINT32<<16 | gc.TUINT64:
		a = s390x.AMOVWZ

		goto rdst

		/*
		 * float to integer
		 */
	case gc.TFLOAT32<<16 | gc.TUINT8,
		gc.TFLOAT32<<16 | gc.TUINT16:
		cvt = gc.Types[gc.TUINT32]

		goto hard

	case gc.TFLOAT32<<16 | gc.TUINT32:
		a = s390x.ACLFEBR

		goto rdst

	case gc.TFLOAT32<<16 | gc.TUINT64:
		a = s390x.ACLGEBR

		goto rdst

	case gc.TFLOAT64<<16 | gc.TUINT8,
		gc.TFLOAT64<<16 | gc.TUINT16:
		cvt = gc.Types[gc.TUINT32]

		goto hard

	case gc.TFLOAT64<<16 | gc.TUINT32:
		a = s390x.ACLFDBR

		goto rdst

	case gc.TFLOAT64<<16 | gc.TUINT64:
		a = s390x.ACLGDBR

		goto rdst

	case gc.TFLOAT32<<16 | gc.TINT8,
		gc.TFLOAT32<<16 | gc.TINT16:
		cvt = gc.Types[gc.TINT32]

		goto hard

	case gc.TFLOAT32<<16 | gc.TINT32:
		a = s390x.ACFEBRA

		goto rdst

	case gc.TFLOAT32<<16 | gc.TINT64:
		a = s390x.ACGEBRA

		goto rdst

	case gc.TFLOAT64<<16 | gc.TINT8,
		gc.TFLOAT64<<16 | gc.TINT16:
		cvt = gc.Types[gc.TINT32]

		goto hard

	case gc.TFLOAT64<<16 | gc.TINT32:
		a = s390x.ACFDBRA

		goto rdst

	case gc.TFLOAT64<<16 | gc.TINT64:
		a = s390x.ACGDBRA

		goto rdst

		/*
		 * integer to float
		 */
	case gc.TUINT8<<16 | gc.TFLOAT32,
		gc.TUINT16<<16 | gc.TFLOAT32:
		cvt = gc.Types[gc.TUINT32]

		goto hard

	case gc.TUINT32<<16 | gc.TFLOAT32:
		a = s390x.ACELFBR

		goto rdst

	case gc.TUINT64<<16 | gc.TFLOAT32:
		a = s390x.ACELGBR

		goto rdst

	case gc.TUINT8<<16 | gc.TFLOAT64,
		gc.TUINT16<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TUINT32]

		goto hard

	case gc.TUINT32<<16 | gc.TFLOAT64:
		a = s390x.ACDLFBR

		goto rdst

	case gc.TUINT64<<16 | gc.TFLOAT64:
		a = s390x.ACDLGBR

		goto rdst

	case gc.TINT8<<16 | gc.TFLOAT32,
		gc.TINT16<<16 | gc.TFLOAT32:
		cvt = gc.Types[gc.TINT32]

		goto hard

	case gc.TINT32<<16 | gc.TFLOAT32:
		a = s390x.ACEFBRA

		goto rdst

	case gc.TINT64<<16 | gc.TFLOAT32:
		a = s390x.ACEGBRA

		goto rdst

	case gc.TINT8<<16 | gc.TFLOAT64,
		gc.TINT16<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TINT32]

		goto hard

	case gc.TINT32<<16 | gc.TFLOAT64:
		a = s390x.ACDFBRA

		goto rdst

	case gc.TINT64<<16 | gc.TFLOAT64:
		a = s390x.ACDGBRA

		goto rdst

		/*
		 * float to float
		 */
	case gc.TFLOAT32<<16 | gc.TFLOAT32:
		a = s390x.AFMOVS

	case gc.TFLOAT64<<16 | gc.TFLOAT64:
		a = s390x.AFMOVD

	case gc.TFLOAT32<<16 | gc.TFLOAT64:
		a = s390x.ALDEBR
		goto rdst

	case gc.TFLOAT64<<16 | gc.TFLOAT32:
		a = s390x.AFRSP
		goto rdst
	}

	gins(a, f, t)
	return

	// requires register destination
rdst:
	{
		var r1 gc.Node
		gc.Regalloc(&r1, t.Type, t)

		gins(a, f, &r1)
		gmove(&r1, t)
		gc.Regfree(&r1)
		return
	}

	// requires register intermediate
hard:
	var r1 gc.Node
	gc.Regalloc(&r1, cvt, t)

	gmove(f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return
}

func intLiteral(n *gc.Node) (x int64, ok bool) {
	switch {
	case n == nil:
		return
	case gc.Isconst(n, gc.CTINT):
		return n.Int(), true
	case gc.Isconst(n, gc.CTBOOL):
		return int64(obj.Bool2int(n.Bool())), true
	}
	return
}

// gins is called by the front end.
// It synthesizes some multiple-instruction sequences
// so the front end can stay simpler.
func gins(as int, f, t *gc.Node) *obj.Prog {
	if as >= obj.A_ARCHSPECIFIC {
		if x, ok := intLiteral(f); ok {
			ginscon(as, x, t)
			return nil // caller must not use
		}
	}
	if as == s390x.ACMP || as == s390x.ACMPU {
		if x, ok := intLiteral(t); ok {
			ginscon2(as, f, x)
			return nil // caller must not use
		}
	}
	return rawgins(as, f, t)
}

/*
 * generate one instruction:
 *	as f, t
 */
func rawgins(as int, f *gc.Node, t *gc.Node) *obj.Prog {
	// TODO(austin): Add self-move test like in 6g (but be careful
	// of truncation moves)

	p := gc.Prog(as)
	gc.Naddr(&p.From, f)
	gc.Naddr(&p.To, t)

	switch as {
	// Bad things the front end has done to us. Crash to find call stack.
	case s390x.AAND, s390x.AMULLD:
		if p.From.Type == obj.TYPE_CONST {
			gc.Debug['h'] = 1
			gc.Fatal("bad inst: %v", p)
		}
	case s390x.ACMP, s390x.ACMPU:
		if p.From.Type == obj.TYPE_MEM || p.To.Type == obj.TYPE_MEM {
			gc.Debug['h'] = 1
			gc.Fatal("bad inst: %v", p)
		}
	}

	if gc.Debug['g'] != 0 {
		fmt.Printf("%v\n", p)
	}

	w := int32(0)
	switch as {
	case s390x.AMOVB,
		s390x.AMOVBU,
		s390x.AMOVBZ,
		s390x.AMOVBZU:
		w = 1

	case s390x.AMOVH,
		s390x.AMOVHU,
		s390x.AMOVHZ,
		s390x.AMOVHZU:
		w = 2

	case s390x.AMOVW,
		s390x.AMOVWU,
		s390x.AMOVWZ,
		s390x.AMOVWZU:
		w = 4

	case s390x.AMOVD,
		s390x.AMOVDU:
		if p.From.Type == obj.TYPE_CONST || p.From.Type == obj.TYPE_ADDR {
			break
		}
		w = 8
	}

	if w != 0 && ((f != nil && p.From.Width < int64(w)) || (t != nil && p.To.Type != obj.TYPE_REG && p.To.Width > int64(w))) {
		gc.Dump("f", f)
		gc.Dump("t", t)
		gc.Fatal("bad width: %v (%d, %d)\n", p, p.From.Width, p.To.Width)
	}

	return p
}

/*
 * return Axxx for Oxxx on type t.
 */
func optoas(op int, t *gc.Type) int {
	if t == nil {
		gc.Fatal("optoas: t is nil")
	}

	a := int(obj.AXXX)
	switch uint32(op)<<16 | uint32(gc.Simtype[t.Etype]) {
	default:
		gc.Fatal("optoas: no entry for op=%v type=%v", gc.Oconv(int(op), 0), t)

	case gc.OEQ<<16 | gc.TBOOL,
		gc.OEQ<<16 | gc.TINT8,
		gc.OEQ<<16 | gc.TUINT8,
		gc.OEQ<<16 | gc.TINT16,
		gc.OEQ<<16 | gc.TUINT16,
		gc.OEQ<<16 | gc.TINT32,
		gc.OEQ<<16 | gc.TUINT32,
		gc.OEQ<<16 | gc.TINT64,
		gc.OEQ<<16 | gc.TUINT64,
		gc.OEQ<<16 | gc.TPTR32,
		gc.OEQ<<16 | gc.TPTR64,
		gc.OEQ<<16 | gc.TFLOAT32,
		gc.OEQ<<16 | gc.TFLOAT64:
		a = s390x.ABEQ

	case gc.ONE<<16 | gc.TBOOL,
		gc.ONE<<16 | gc.TINT8,
		gc.ONE<<16 | gc.TUINT8,
		gc.ONE<<16 | gc.TINT16,
		gc.ONE<<16 | gc.TUINT16,
		gc.ONE<<16 | gc.TINT32,
		gc.ONE<<16 | gc.TUINT32,
		gc.ONE<<16 | gc.TINT64,
		gc.ONE<<16 | gc.TUINT64,
		gc.ONE<<16 | gc.TPTR32,
		gc.ONE<<16 | gc.TPTR64,
		gc.ONE<<16 | gc.TFLOAT32,
		gc.ONE<<16 | gc.TFLOAT64:
		a = s390x.ABNE

	case gc.OLT<<16 | gc.TINT8, // ACMP
		gc.OLT<<16 | gc.TINT16,
		gc.OLT<<16 | gc.TINT32,
		gc.OLT<<16 | gc.TINT64,
		gc.OLT<<16 | gc.TUINT8,
		// ACMPU
		gc.OLT<<16 | gc.TUINT16,
		gc.OLT<<16 | gc.TUINT32,
		gc.OLT<<16 | gc.TUINT64,
		gc.OLT<<16 | gc.TFLOAT32,
		// AFCMPU
		gc.OLT<<16 | gc.TFLOAT64:
		a = s390x.ABLT

	case gc.OLE<<16 | gc.TINT8, // ACMP
		gc.OLE<<16 | gc.TINT16,
		gc.OLE<<16 | gc.TINT32,
		gc.OLE<<16 | gc.TINT64,
		gc.OLE<<16 | gc.TUINT8,
		// ACMPU
		gc.OLE<<16 | gc.TUINT16,
		gc.OLE<<16 | gc.TUINT32,
		gc.OLE<<16 | gc.TUINT64,
		gc.OLE<<16 | gc.TFLOAT32, // TODO(WGO)
		gc.OLE<<16 | gc.TFLOAT64: // TODO(WGO)
		// on Power, No OLE for floats, because it mishandles NaN.
		// on Power, Front end must reverse comparison or use OLT and OEQ together.
		a = s390x.ABLE

	case gc.OGT<<16 | gc.TINT8,
		gc.OGT<<16 | gc.TINT16,
		gc.OGT<<16 | gc.TINT32,
		gc.OGT<<16 | gc.TINT64,
		gc.OGT<<16 | gc.TUINT8,
		gc.OGT<<16 | gc.TUINT16,
		gc.OGT<<16 | gc.TUINT32,
		gc.OGT<<16 | gc.TUINT64,
		gc.OGT<<16 | gc.TFLOAT32,
		gc.OGT<<16 | gc.TFLOAT64:
		a = s390x.ABGT

	case gc.OGE<<16 | gc.TINT8,
		gc.OGE<<16 | gc.TINT16,
		gc.OGE<<16 | gc.TINT32,
		gc.OGE<<16 | gc.TINT64,
		gc.OGE<<16 | gc.TUINT8,
		gc.OGE<<16 | gc.TUINT16,
		gc.OGE<<16 | gc.TUINT32,
		gc.OGE<<16 | gc.TUINT64,
		gc.OGE<<16 | gc.TFLOAT32, // TODO(WGO)
		gc.OGE<<16 | gc.TFLOAT64: // TODO(WGO)
		// No OGE for floats, because it mishandles NaN.
		// Front end must reverse comparison or use OLT and OEQ together.
		a = s390x.ABGE

	case gc.OCMP<<16 | gc.TBOOL,
		gc.OCMP<<16 | gc.TINT8,
		gc.OCMP<<16 | gc.TINT16,
		gc.OCMP<<16 | gc.TINT32,
		gc.OCMP<<16 | gc.TPTR32,
		gc.OCMP<<16 | gc.TINT64:
		a = s390x.ACMP

	case gc.OCMP<<16 | gc.TUINT8,
		gc.OCMP<<16 | gc.TUINT16,
		gc.OCMP<<16 | gc.TUINT32,
		gc.OCMP<<16 | gc.TUINT64,
		gc.OCMP<<16 | gc.TPTR64:
		a = s390x.ACMPU

	case gc.OCMP<<16 | gc.TFLOAT32:
		a = s390x.ACEBR

	case gc.OCMP<<16 | gc.TFLOAT64:
		a = s390x.AFCMPU

	case gc.OAS<<16 | gc.TBOOL,
		gc.OAS<<16 | gc.TINT8:
		a = s390x.AMOVB

	case gc.OAS<<16 | gc.TUINT8:
		a = s390x.AMOVBZ

	case gc.OAS<<16 | gc.TINT16:
		a = s390x.AMOVH

	case gc.OAS<<16 | gc.TUINT16:
		a = s390x.AMOVHZ

	case gc.OAS<<16 | gc.TINT32:
		a = s390x.AMOVW

	case gc.OAS<<16 | gc.TUINT32,
		gc.OAS<<16 | gc.TPTR32:
		a = s390x.AMOVWZ

	case gc.OAS<<16 | gc.TINT64,
		gc.OAS<<16 | gc.TUINT64,
		gc.OAS<<16 | gc.TPTR64:
		a = s390x.AMOVD

	case gc.OAS<<16 | gc.TFLOAT32:
		a = s390x.AFMOVS

	case gc.OAS<<16 | gc.TFLOAT64:
		a = s390x.AFMOVD

	case gc.OADD<<16 | gc.TINT8,
		gc.OADD<<16 | gc.TUINT8,
		gc.OADD<<16 | gc.TINT16,
		gc.OADD<<16 | gc.TUINT16,
		gc.OADD<<16 | gc.TINT32,
		gc.OADD<<16 | gc.TUINT32,
		gc.OADD<<16 | gc.TPTR32,
		gc.OADD<<16 | gc.TINT64,
		gc.OADD<<16 | gc.TUINT64,
		gc.OADD<<16 | gc.TPTR64:
		a = s390x.AADD

	case gc.OADD<<16 | gc.TFLOAT32:
		a = s390x.AFADDS

	case gc.OADD<<16 | gc.TFLOAT64:
		a = s390x.AFADD

	case gc.OSUB<<16 | gc.TINT8,
		gc.OSUB<<16 | gc.TUINT8,
		gc.OSUB<<16 | gc.TINT16,
		gc.OSUB<<16 | gc.TUINT16,
		gc.OSUB<<16 | gc.TINT32,
		gc.OSUB<<16 | gc.TUINT32,
		gc.OSUB<<16 | gc.TPTR32,
		gc.OSUB<<16 | gc.TINT64,
		gc.OSUB<<16 | gc.TUINT64,
		gc.OSUB<<16 | gc.TPTR64:
		a = s390x.ASUB

	case gc.OSUB<<16 | gc.TFLOAT32:
		a = s390x.AFSUBS

	case gc.OSUB<<16 | gc.TFLOAT64:
		a = s390x.AFSUB

	case gc.OMINUS<<16 | gc.TINT8,
		gc.OMINUS<<16 | gc.TUINT8,
		gc.OMINUS<<16 | gc.TINT16,
		gc.OMINUS<<16 | gc.TUINT16,
		gc.OMINUS<<16 | gc.TINT32,
		gc.OMINUS<<16 | gc.TUINT32,
		gc.OMINUS<<16 | gc.TPTR32,
		gc.OMINUS<<16 | gc.TINT64,
		gc.OMINUS<<16 | gc.TUINT64,
		gc.OMINUS<<16 | gc.TPTR64:
		a = s390x.ANEG

	case gc.OAND<<16 | gc.TINT8,
		gc.OAND<<16 | gc.TUINT8,
		gc.OAND<<16 | gc.TINT16,
		gc.OAND<<16 | gc.TUINT16,
		gc.OAND<<16 | gc.TINT32,
		gc.OAND<<16 | gc.TUINT32,
		gc.OAND<<16 | gc.TPTR32,
		gc.OAND<<16 | gc.TINT64,
		gc.OAND<<16 | gc.TUINT64,
		gc.OAND<<16 | gc.TPTR64:
		a = s390x.AAND

	case gc.OOR<<16 | gc.TINT8,
		gc.OOR<<16 | gc.TUINT8,
		gc.OOR<<16 | gc.TINT16,
		gc.OOR<<16 | gc.TUINT16,
		gc.OOR<<16 | gc.TINT32,
		gc.OOR<<16 | gc.TUINT32,
		gc.OOR<<16 | gc.TPTR32,
		gc.OOR<<16 | gc.TINT64,
		gc.OOR<<16 | gc.TUINT64,
		gc.OOR<<16 | gc.TPTR64:
		a = s390x.AOR

	case gc.OXOR<<16 | gc.TINT8,
		gc.OXOR<<16 | gc.TUINT8,
		gc.OXOR<<16 | gc.TINT16,
		gc.OXOR<<16 | gc.TUINT16,
		gc.OXOR<<16 | gc.TINT32,
		gc.OXOR<<16 | gc.TUINT32,
		gc.OXOR<<16 | gc.TPTR32,
		gc.OXOR<<16 | gc.TINT64,
		gc.OXOR<<16 | gc.TUINT64,
		gc.OXOR<<16 | gc.TPTR64:
		a = s390x.AXOR

		// TODO(minux): handle rotates
	//case CASE(OLROT, TINT8):
	//case CASE(OLROT, TUINT8):
	//case CASE(OLROT, TINT16):
	//case CASE(OLROT, TUINT16):
	//case CASE(OLROT, TINT32):
	//case CASE(OLROT, TUINT32):
	//case CASE(OLROT, TPTR32):
	//case CASE(OLROT, TINT64):
	//case CASE(OLROT, TUINT64):
	//case CASE(OLROT, TPTR64):
	//	a = 0//???; RLDC?
	//	break;

	case gc.OLSH<<16 | gc.TINT8,
		gc.OLSH<<16 | gc.TUINT8,
		gc.OLSH<<16 | gc.TINT16,
		gc.OLSH<<16 | gc.TUINT16,
		gc.OLSH<<16 | gc.TINT32,
		gc.OLSH<<16 | gc.TUINT32,
		gc.OLSH<<16 | gc.TPTR32,
		gc.OLSH<<16 | gc.TINT64,
		gc.OLSH<<16 | gc.TUINT64,
		gc.OLSH<<16 | gc.TPTR64:
		a = s390x.ASLD

	case gc.ORSH<<16 | gc.TUINT8,
		gc.ORSH<<16 | gc.TUINT16,
		gc.ORSH<<16 | gc.TUINT32,
		gc.ORSH<<16 | gc.TPTR32,
		gc.ORSH<<16 | gc.TUINT64,
		gc.ORSH<<16 | gc.TPTR64:
		a = s390x.ASRD

	case gc.ORSH<<16 | gc.TINT8,
		gc.ORSH<<16 | gc.TINT16,
		gc.ORSH<<16 | gc.TINT32,
		gc.ORSH<<16 | gc.TINT64:
		a = s390x.ASRAD

		// TODO(minux): handle rotates
	//case CASE(ORROTC, TINT8):
	//case CASE(ORROTC, TUINT8):
	//case CASE(ORROTC, TINT16):
	//case CASE(ORROTC, TUINT16):
	//case CASE(ORROTC, TINT32):
	//case CASE(ORROTC, TUINT32):
	//case CASE(ORROTC, TINT64):
	//case CASE(ORROTC, TUINT64):
	//	a = 0//??? RLDC??
	//	break;

	case gc.OHMUL<<16 | gc.TINT64:
		a = s390x.AMULHD

	case gc.OHMUL<<16 | gc.TUINT64,
		gc.OHMUL<<16 | gc.TPTR64:
		a = s390x.AMULHDU

	case gc.OMUL<<16 | gc.TINT8,
		gc.OMUL<<16 | gc.TINT16,
		gc.OMUL<<16 | gc.TINT32,
		gc.OMUL<<16 | gc.TINT64:
		a = s390x.AMULLD

	case gc.OMUL<<16 | gc.TUINT8,
		gc.OMUL<<16 | gc.TUINT16,
		gc.OMUL<<16 | gc.TUINT32,
		gc.OMUL<<16 | gc.TPTR32,
		// don't use word multiply, the high 32-bit are undefined.
		gc.OMUL<<16 | gc.TUINT64,
		gc.OMUL<<16 | gc.TPTR64:
		// for 64-bit multiplies, signedness doesn't matter.
		a = s390x.AMULLD

	case gc.OMUL<<16 | gc.TFLOAT32:
		a = s390x.AFMULS

	case gc.OMUL<<16 | gc.TFLOAT64:
		a = s390x.AFMUL

	case gc.ODIV<<16 | gc.TINT8,
		gc.ODIV<<16 | gc.TINT16,
		gc.ODIV<<16 | gc.TINT32,
		gc.ODIV<<16 | gc.TINT64:
		a = s390x.ADIVD

	case gc.ODIV<<16 | gc.TUINT8,
		gc.ODIV<<16 | gc.TUINT16,
		gc.ODIV<<16 | gc.TUINT32,
		gc.ODIV<<16 | gc.TPTR32,
		gc.ODIV<<16 | gc.TUINT64,
		gc.ODIV<<16 | gc.TPTR64:
		a = s390x.ADIVDU

	case gc.ODIV<<16 | gc.TFLOAT32:
		a = s390x.AFDIVS

	case gc.ODIV<<16 | gc.TFLOAT64:
		a = s390x.AFDIV

	case gc.OSQRT<<16 | gc.TFLOAT64:
		a = s390x.AFSQRT
	}

	return a
}

const (
	ODynam   = 1 << 0
	OAddable = 1 << 1
)

func xgen(n *gc.Node, a *gc.Node, o int) bool {
	// TODO(minux)

	return -1 != 0 /*TypeKind(100016)*/
}

func sudoclean() {
	return
}

/*
 * generate code to compute address of n,
 * a reference to a (perhaps nested) field inside
 * an array or struct.
 * return 0 on failure, 1 on success.
 * on success, leaves usable address in a.
 *
 * caller is responsible for calling sudoclean
 * after successful sudoaddable,
 * to release the register used for a.
 */
func sudoaddable(as int, n *gc.Node, a *obj.Addr) bool {
	// TODO(minux)

	*a = obj.Addr{}
	return false
}
