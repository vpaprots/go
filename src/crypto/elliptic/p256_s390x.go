// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build s390x

package elliptic

import (
	"math/big"
	"sync"
	"fmt"
	"bytes"
)

type (
	p256Curve struct {
		*CurveParams
	}

	p256Point struct {
		x [32]byte
		y [32]byte
		z [32]byte
	}
)

var (
	p256            p256Curve
	p256Precomputed *[37][64]p256Point
	precomputeOnce  sync.Once
)

func initP256() {
	// See FIPS 186-3, section D.2.3
	p256.CurveParams = &CurveParams{Name: "P-256"}
	p256.P, _ = new(big.Int).SetString("115792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	p256.N, _ = new(big.Int).SetString("115792089210356248762697446949407573529996955224135760342422259061068512044369", 10)
	p256.B, _ = new(big.Int).SetString("5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b", 16)
	p256.Gx, _ = new(big.Int).SetString("6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296", 16)
	p256.Gy, _ = new(big.Int).SetString("4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5", 16)
	p256.BitSize = 256
}

func (curve p256Curve) Params() *CurveParams {
	return curve.CurveParams
}

func (curve p256Curve) TestDouble(x1, y1, z1 *big.Int) (x3,y3,z3 *big.Int) {
	res  := new(p256Point)
	in   := new(p256Point)
	copy(in.x[:], fromBig(x1))
	copy(in.y[:], fromBig(y1))
	copy(in.z[:], fromBig(z1))
	p256PointDoubleAsm(res, in)
	return new(big.Int).SetBytes(res.x[:]),new(big.Int).SetBytes(res.y[:]), new(big.Int).SetBytes(res.z[:])  
}

func (curve p256Curve) TestInv(k *big.Int) *big.Int {
	res := make([]byte, 32)
	x := fromBig(k)
	p256Inverse(res, x)
	return new(big.Int).SetBytes(res)
}

// Functions implemented in p256_asm_s390x.s
// Montgomery multiplication modulo P256
func p256Mul(res, in1, in2 []byte)

func p256MulBig(res, in1, in2 []byte){
	x1 := new(big.Int).SetBytes(in1)
	x2 := new(big.Int).SetBytes(in2)
	Rinv, _ := new(big.Int).SetString("fffffffe00000003fffffffd0000000200000001fffffffe0000000300000000", 16) //minv(2^256,p)
	temp := new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv)
	//t := make([]byte, 32)
	//p256OrdMul(t, in1, in2)
	
	copy(res, fromBig(new(big.Int).Mod(temp, p256.P)))

	/*if (!bytes.Equal(t,res)) {
		fmt.Printf("TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
		fmt.Printf("TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
		fmt.Printf("EXPECTED %s\n", new(big.Int).SetBytes(res).Text(16))
		fmt.Printf("FOUND    %s\n", new(big.Int).SetBytes(t).Text(16))
	} /*else {
		fmt.Printf("+\n")
	}*/
}

// Montgomery square modulo P256
func p256Sqr(res, in []byte){
	p256Mul(res, in, in)
}

// Montgomery multiplication by 1
func p256FromMont(res, in []byte){
	x1 := new(big.Int).SetBytes(in)
	Rinv, _ := new(big.Int).SetString("fffffffe00000003fffffffd0000000200000001fffffffe0000000300000000", 16) //minv(2^256,p)
	temp := new(big.Int).Mul(x1, Rinv)
	
	copy(res, fromBig(new(big.Int).Mod(temp, p256.P)))
}

// iff cond == 1  val <- -val
func p256NegCond(val *p256Point, cond int){
	if (cond == 1) {
		copy(val.y[:], fromBig(new(big.Int).Mod(new(big.Int).Sub(p256.P, new(big.Int).SetBytes(val.y[:])), p256.P)))
	}
}

// if cond == 0 res <- b; else res <- a
func p256MovCond(res, a, b *p256Point, cond int){
	if (cond == 0) {
		copy(res.x[:], b.x[:])
		copy(res.y[:], b.y[:])
		copy(res.z[:], b.z[:])
	} else {
		copy(res.x[:], a.x[:])
		copy(res.y[:], a.y[:])
		copy(res.z[:], a.z[:])
	}
}

// Constant time table access
func p256Select(point *p256Point, table []p256Point, idx int) {
	if idx==0 {
		copy(point.x[:], make([]byte, 32))
		copy(point.y[:], make([]byte, 32))
		copy(point.z[:], make([]byte, 32))
	} else {
		copy(point.x[:], table[idx-1].x[:])
		copy(point.y[:], table[idx-1].y[:])
		copy(point.z[:], table[idx-1].z[:])
	}
}

func p256SelectBase(point *p256Point, table []p256Point, idx int) {
	if idx==0 {
		copy(point.x[:], table[0].z[:])
		copy(point.y[:], table[0].z[:])
	} else {
		copy(point.x[:], table[idx-1].x[:])
		copy(point.y[:], table[idx-1].y[:])
	}
}

// Montgomery multiplication modulo Ord(G)
func p256OrdMul(res, in1, in2 []byte)

func p256OrdMulBig(res, in1, in2 []byte) {
	x1 := new(big.Int).SetBytes(in1)
	x2 := new(big.Int).SetBytes(in2)
	Rinv, _ := new(big.Int).SetString("60d066334905c1e907f8b6041e607725badef3e243566fafce1bc8f79c197c79", 16) //minv(2^256,n)
	temp := new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv)
	t := make([]byte, 32)
	p256OrdMul(t, in1, in2)
	
	copy(res, fromBig(new(big.Int).Mod(temp, p256.N)))

	if (!bytes.Equal(t,res)) {
		fmt.Printf("TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
		fmt.Printf("TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
		fmt.Printf("EXPECTED %s\n", new(big.Int).SetBytes(res).Text(16))
		fmt.Printf("FOUND    %s\n", new(big.Int).SetBytes(t).Text(16))
	} /*else {
		fmt.Printf("+\n")
	}*/
}

/*
x=0xfef7163fe956730df28c999458d9c038a17b9500f51bd2f803cabdf9818bc253
RR=0xa3ff46f14bce132cd59447e8378fe08999ca0e402b77090c7215405740ffd73b
r0=x*mydigit(RR,0)
r1=x*mydigit(RR,1)
r2=x*mydigit(RR,2)
r3=x*mydigit(RR,3)
r4=x*mydigit(RR,4)
r5=x*mydigit(RR,5)
r6=x*mydigit(RR,6)
r7=x*mydigit(RR,7)
t0=(r0+mydigit(r0*k0,0)*n)>>32
t1=(t0+r1+mydigit((t0+r1)*k0,0)*n)>>32
t2=(t1+r2+mydigit((t1+r2)*k0,0)*n)>>32
t3=(t2+r3+mydigit((t2+r3)*k0,0)*n)>>32
t4=(t3+r4+mydigit((t3+r4)*k0,0)*n)>>32
t5=(t4+r5+mydigit((t4+r5)*k0,0)*n)>>32
t6=(t5+r6+mydigit((t5+r6)*k0,0)*n)>>32
t7=(t6+r7+mydigit((t6+r7)*k0,0)*n)>>32
t0
t1
t2
t3
t4
t5
t6
t7

x2=0x3ac62a6b166498176eb40d98a27587d8d7c97ac5f11ca1d8d851202b22c3f5c8
y2=0x234aabe92636af27ea8edcd2392f97839c5a74b7ddea27bce94c2d270fb65157
z2=0x57401fa8db8e8e1118a40621ce27d6842bc1e1cef6138faabaf37b85a2a774ea
x1=0x6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296
y1=0x4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5
z1=1

Rinv=minv(2^256,p)
A=Rinv^2*x1*z2^2
B=Rinv^3*y1*z2^3
C=Rinv^2*x2*z1^2-A
D=Rinv^3*y2*z1^3-B
A=A%p
B=B%p
C=C%p
D=D%p
x3 = D^2*Rinv - 2*A×C^2*Rinv^2 - C^3*Rinv^2
Y₃ = D×(A×C² - X₃) - B×C³
Z₃ = Z₁×Z2×C
(z1*z2*C)%p

define mydigit(x,n) = (x>>(32*n))%2^32
x=0xa007c8559316f82de3d5d9f28b8ffcdf5949bd551f7a1348b8acc00860e058
RR=0x66e12d94f3d956202845b2392b6bec594699799c49bd6fa683244c95be79eea2
r0=x*mydigit(RR,0)
r1=x*mydigit(RR,1)
r2=x*mydigit(RR,2)
r3=x*mydigit(RR,3)
r4=x*mydigit(RR,4)
r5=x*mydigit(RR,5)
r6=x*mydigit(RR,6)
r7=x*mydigit(RR,7)
t0=(r0+mydigit(r0*k0,0)*p)>>32
t1=(t0+r1+mydigit((t0+r1)*k0,0)*p)>>32
t2=(t1+r2+mydigit((t1+r2)*k0,0)*p)>>32
t3=(t2+r3+mydigit((t2+r3)*k0,0)*p)>>32
t4=(t3+r4+mydigit((t3+r4)*k0,0)*p)>>32
t5=(t4+r5+mydigit((t4+r5)*k0,0)*p)>>32
t6=(t5+r6+mydigit((t5+r6)*k0,0)*p)>>32
t7=(t6+r7+mydigit((t6+r7)*k0,0)*p)>>32
t0
t1
t2
t3
t4
t5
t6
t7


define mydigit(x,n) = (x>>(32*n))%2^32
define mydigit2(x,n) = (x>>(64*n))%2^64
x=0xa007c8559316f82de3d5d9f28b8ffcdf5949bd551f7a1348b8acc00860e058
RR=0x66e12d94f3d956202845b2392b6bec594699799c49bd6fa683244c95be79eea2
r0=x*mydigit(RR,0)
r1=x*mydigit(RR,1)
r2=x*mydigit(RR,2)
r3=x*mydigit(RR,3)
r4=x*mydigit(RR,4)
r5=x*mydigit(RR,5)
r6=x*mydigit(RR,6)
r7=x*mydigit(RR,7)

t0=r0+(r1<<32)
red0=mydigit2(t0,0)*p+mydigit2(t0,0)
t1=(t0>>64)+r2+(r3<<32)+(red0>>64)
red1=mydigit2(t1,0)*p+mydigit2(t1,0)
t2=(t1>>64)+r4+(r5<<32)+(red1>>64)
red2=mydigit2(t2,0)*p+mydigit2(t2,0)
t3=(t2>>64)+r6+(r7<<32)+(red2>>64)
red3=mydigit2(t3,0)*p+mydigit2(t3,0)


t0=(r0+mydigit(r0*k0,0)*p)>>32
t1=(t0+r1+mydigit((t0+r1)*k0,0)*p)>>32
t2=(t1+r2+mydigit((t1+r2)*k0,0)*p)>>32
t3=(t2+r3+mydigit((t2+r3)*k0,0)*p)>>32
t4=(t3+r4+mydigit((t3+r4)*k0,0)*p)>>32
t5=(t4+r5+mydigit((t4+r5)*k0,0)*p)>>32
t6=(t5+r6+mydigit((t5+r6)*k0,0)*p)>>32
t7=(t6+r7+mydigit((t6+r7)*k0,0)*p)>>32
t0
t1
t2
t3
t4
t5
t6
t7

*/
// Montgomery square modulo Ord(G), repeated n times
func p256OrdSqr(res, in []byte, n int) {
	copy(res, in)
	for i := 0; i < n; i += 1 {
		p256OrdMul(res, res, res)
	}
}

func p256OrdSqrBig(res, in []byte, n int) {
	x := new(big.Int).SetBytes(in)
	Rinv, _ := new(big.Int).SetString("60d066334905c1e907f8b6041e607725badef3e243566fafce1bc8f79c197c79", 16)
	for i := 0; i < n; i += 1 {
		x = new(big.Int).Mul(new(big.Int).Mul(x, x), Rinv)
	    x = new(big.Int).Mod(x, p256.N)	
	}
	copy(res, fromBig(x))	
}

// Point add with P2 being affine point
// If sign == 1 -> P2 = -P2
// If sel == 0 -> P3 = P1
// if zero == 0 -> P3 = P2
func p256PointAddAffineAsm(P3, P1, P2 *p256Point, sign, sel, zero int)

// Point add
func p256PointAddAsm(P3, P1, P2 *p256Point)
func p256PointDoubleAsm(P3, P1 *p256Point)

func (curve p256Curve) Inverse(k *big.Int) *big.Int {
	if k.Cmp(p256.N) >= 0 {
		// This should never happen.
		reducedK := new(big.Int).Mod(k, p256.N)
		k = reducedK
	}

	// table will store precomputed powers of x. The 32 bytes at index
	// i store x^(i+1).
	var table [15][32]byte

	x := fromBig(k)
	// This code operates in the Montgomery domain where R = 2^256 mod n
	// and n is the order of the scalar field. (See initP256 for the
	// value.) Elements in the Montgomery domain take the form a×R and
	// multiplication of x and y in the calculates (x × y × R^-1) mod n. RR
	// is R×R mod n thus the Montgomery multiplication x and RR gives x×R,
	// i.e. converts x into the Montgomery domain. Stored in BigEndian form
	RR := []byte{0x66, 0xe1, 0x2d, 0x94, 0xf3, 0xd9, 0x56, 0x20, 0x28, 0x45, 0xb2, 0x39, 0x2b, 0x6b, 0xec, 0x59, 
		         0x46, 0x99, 0x79, 0x9c, 0x49, 0xbd, 0x6f, 0xa6, 0x83, 0x24, 0x4c, 0x95, 0xbe, 0x79, 0xee, 0xa2}
	
	p256OrdMul(table[0][:], x, RR)

	// Prepare the table, no need in constant time access, because the
	// power is not a secret. (Entry 0 is never used.)
	for i := 2; i < 16; i += 2 {
		p256OrdSqr(table[i-1][:], table[(i/2)-1][:], 1)
		p256OrdMul(table[i][:], table[i-1][:], table[0][:])
	}

	copy(x, table[14][:])  // f

	p256OrdSqr(x[0:32], x[0:32], 4)
	p256OrdMul(x[0:32], x[0:32], table[14][:]) // ff
	t := make([]byte, 32)
	copy(t, x)

	p256OrdSqr(x, x, 8)
	p256OrdMul(x, x, t) // ffff
	copy(t, x)

	p256OrdSqr(x, x, 16)
	p256OrdMul(x, x, t) // ffffffff
	copy(t, x)

	p256OrdSqr(x, x, 64) // ffffffff0000000000000000
	p256OrdMul(x, x, t)  // ffffffff00000000ffffffff
	p256OrdSqr(x, x, 32) // ffffffff00000000ffffffff00000000
	p256OrdMul(x, x, t)  // ffffffff00000000ffffffffffffffff

	// Remaining 32 windows
	expLo := [32]byte{0xb, 0xc, 0xe, 0x6, 0xf, 0xa, 0xa, 0xd, 0xa, 0x7, 0x1, 0x7, 0x9, 0xe, 0x8, 0x4, 
		              0xf, 0x3, 0xb, 0x9, 0xc, 0xa, 0xc, 0x2, 0xf, 0xc, 0x6, 0x3, 0x2, 0x5, 0x4, 0xf}
	for i := 0; i < 32; i++ {
		p256OrdSqr(x, x, 4)
		p256OrdMul(x, x, table[expLo[i]-1][:])
	}

	// Multiplying by one in the Montgomery domain converts a Montgomery
	// value out of the domain.
	one := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		          0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	p256OrdMul(x, x, one)

	return new(big.Int).SetBytes(x)
}

func (curve p256Curve) InverseBig(k *big.Int) *big.Int {
	if k.Cmp(p256.N) >= 0 {
		// This should never happen.
		reducedK := new(big.Int).Mod(k, p256.N)
		k = reducedK
	}

	// table will store precomputed powers of x. The 32 bytes at index
	// i store x^(i+1).
	var table [15][32]byte

	x := fromBig(k)
	// This code operates in the Montgomery domain where R = 2^256 mod n
	// and n is the order of the scalar field. (See initP256 for the
	// value.) Elements in the Montgomery domain take the form a×R and
	// multiplication of x and y in the calculates (x × y × R^-1) mod n. RR
	// is R×R mod n thus the Montgomery multiplication x and RR gives x×R,
	// i.e. converts x into the Montgomery domain. Stored in BigEndian form
	RR := []byte{0x66, 0xe1, 0x2d, 0x94, 0xf3, 0xd9, 0x56, 0x20, 0x28, 0x45, 0xb2, 0x39, 0x2b, 0x6b, 0xec, 0x59, 
		         0x46, 0x99, 0x79, 0x9c, 0x49, 0xbd, 0x6f, 0xa6, 0x83, 0x24, 0x4c, 0x95, 0xbe, 0x79, 0xee, 0xa2}
	
	p256OrdMulBig(table[0][:], x, RR)

	// Prepare the table, no need in constant time access, because the
	// power is not a secret. (Entry 0 is never used.)
	for i := 2; i < 16; i += 2 {
		p256OrdSqrBig(table[i-1][:], table[(i/2)-1][:], 1)
		p256OrdMulBig(table[i][:], table[i-1][:], table[0][:])
	}

	copy(x, table[14][:])  // f

	p256OrdSqrBig(x[0:32], x[0:32], 4)
	p256OrdMulBig(x[0:32], x[0:32], table[14][:]) // ff
	t := make([]byte, 32)
	copy(t, x)

	p256OrdSqrBig(x, x, 8)
	p256OrdMulBig(x, x, t) // ffff
	copy(t, x)

	p256OrdSqrBig(x, x, 16)
	p256OrdMulBig(x, x, t) // ffffffff
	copy(t, x)

	p256OrdSqrBig(x, x, 64) // ffffffff0000000000000000
	p256OrdMulBig(x, x, t)  // ffffffff00000000ffffffff
	p256OrdSqrBig(x, x, 32) // ffffffff00000000ffffffff00000000
	p256OrdMulBig(x, x, t)  // ffffffff00000000ffffffffffffffff

	// Remaining 32 windows
	expLo := [32]byte{0xb, 0xc, 0xe, 0x6, 0xf, 0xa, 0xa, 0xd, 0xa, 0x7, 0x1, 0x7, 0x9, 0xe, 0x8, 0x4, 
		              0xf, 0x3, 0xb, 0x9, 0xc, 0xa, 0xc, 0x2, 0xf, 0xc, 0x6, 0x3, 0x2, 0x5, 0x4, 0xf}
	for i := 0; i < 32; i++ {
		p256OrdSqrBig(x, x, 4)
		p256OrdMulBig(x, x, table[expLo[i]-1][:])
	}

	// Multiplying by one in the Montgomery domain converts a Montgomery
	// value out of the domain.
	one := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		          0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	p256OrdMulBig(x, x, one)

	return new(big.Int).SetBytes(x)
}


// fromBig converts a *big.Int into a format used by this code.
func fromBig(big *big.Int) []byte {
	// This could be done a lot more efficiently...
	res := big.Bytes()
	if (32 == len(res)) {
		return res
	}
	t := make([]byte, 32)
	offset := 32-len(res)
	for i:=len(res)-1; i>=0; i-- {
		t[i+offset] = res[i]
	}
	return t
}

// fromBig converts a *big.Int into a format used by this code.
func fromBigDel(out []uint64, big *big.Int) {
	for i := range out {
		out[i] = 0
	}

	for i, v := range big.Bits() {
		out[i] = uint64(v)
	}
}

// p256GetScalar makes sure byte array will have 32 byte elements, If the scalar
// is equal or greater than the order of the group, it's reduced modulo that order.
func p256GetScalar(in []byte) ([]byte) {
	n := new(big.Int).SetBytes(in)

	if n.Cmp(p256.N) >= 0 {
		n.Mod(n, p256.N)
	}
	return fromBig(n)
}

// p256Mul operates in a Montgomery domain with R = 2^256 mod p, where p is the
// underlying field of the curve. (See initP256 for the value.) Thus rr here is
// R×R mod p. See comment in Inverse about how this is used.
var rr = []byte{  0x00, 0x00, 0x00, 0x04, 0xff, 0xff, 0xff, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
	              0xff, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}

// (This is one, in the Montgomery domain.)
var one = []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			        0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

func maybeReduceModP(in *big.Int) *big.Int {
	if in.Cmp(p256.P) < 0 {
		return in
	}
	return new(big.Int).Mod(in, p256.P)
}

func (curve p256Curve) CombinedMult(bigX, bigY *big.Int, baseScalar, scalar []byte) (x, y *big.Int) {
	var r1, r2 p256Point
	r1.p256BaseMult(p256GetScalar(baseScalar))

	copy(r2.x[:], fromBig(maybeReduceModP(bigX)))
	copy(r2.y[:], fromBig(maybeReduceModP(bigY)))
	copy(r2.z[:], one)
	p256Mul(r2.x[:], r2.x[:], rr[:])
	p256Mul(r2.y[:], r2.y[:], rr[:])

	r2.p256ScalarMult(p256GetScalar(scalar))
	p256PointAddAsm(&r1, &r1, &r2)
	return r1.p256PointToAffine()
}

func (curve p256Curve) ScalarBaseMult(scalar []byte) (x, y *big.Int) {
	var r p256Point
	r.p256BaseMult(p256GetScalar(scalar))
	return r.p256PointToAffine()
}

func (curve p256Curve) ScalarMult(bigX, bigY *big.Int, scalar []byte) (x, y *big.Int) {
	var r p256Point
	copy(r.x[:], fromBig(maybeReduceModP(bigX)))
	copy(r.y[:], fromBig(maybeReduceModP(bigY)))
	//PrintPoint("r", &r)
	copy(r.z[:], one)
	p256Mul(r.x[:], r.x[:], rr[:])
	p256Mul(r.y[:], r.y[:], rr[:])
	//PrintPoint("r", &r)
	r.p256ScalarMult(p256GetScalar(scalar))
	return r.p256PointToAffine()
}

func (p *p256Point) p256PointToAffine() (x, y *big.Int) {
	zInv := make([]byte, 32)
	zInvSq := make([]byte, 32)
	p256Inverse(zInv, p.z[:])
	//fmt.Printf("zInv %s\n", new(big.Int).SetBytes(zInv).Text(16))
	p256Sqr(zInvSq, zInv)
	p256Mul(zInv, zInv, zInvSq)

	p256Mul(zInvSq, p.x[:], zInvSq)
	p256Mul(zInv, p.y[:], zInv)

	p256FromMont(zInvSq, zInvSq)
	p256FromMont(zInv, zInv)

	return new(big.Int).SetBytes(zInvSq), new(big.Int).SetBytes(zInv)
}

// p256Inverse sets out to in^-1 mod p.
func p256Inverse(out, in []byte) {
	var stack [6 * 32]byte
	p2 := stack[32*0 : 32*0+32]
	p4 := stack[32*1 : 32*1+32]
	p8 := stack[32*2 : 32*2+32]
	p16 := stack[32*3 : 32*3+32]
	p32 := stack[32*4 : 32*4+32]

	p256Sqr(out, in)
	p256Mul(p2, out, in) // 3*p

	p256Sqr(out, p2)
	p256Sqr(out, out)
	p256Mul(p4, out, p2) // f*p

	p256Sqr(out, p4)
	p256Sqr(out, out)
	p256Sqr(out, out)
	p256Sqr(out, out)
	p256Mul(p8, out, p4) // ff*p

	p256Sqr(out, p8)

	for i := 0; i < 7; i++ {
		p256Sqr(out, out)
	}
	p256Mul(p16, out, p8) // ffff*p

	p256Sqr(out, p16)
	for i := 0; i < 15; i++ {
		p256Sqr(out, out)
	}
	p256Mul(p32, out, p16) // ffffffff*p

	p256Sqr(out, p32)

	for i := 0; i < 31; i++ {
		p256Sqr(out, out)
	}
	p256Mul(out, out, in)

	for i := 0; i < 32*4; i++ {
		p256Sqr(out, out)
	}
	p256Mul(out, out, p32)

	for i := 0; i < 32; i++ {
		p256Sqr(out, out)
	}
	p256Mul(out, out, p32)

	for i := 0; i < 16; i++ {
		p256Sqr(out, out)
	}
	p256Mul(out, out, p16)

	for i := 0; i < 8; i++ {
		p256Sqr(out, out)
	}
	p256Mul(out, out, p8)

	p256Sqr(out, out)
	p256Sqr(out, out)
	p256Sqr(out, out)
	p256Sqr(out, out)
	p256Mul(out, out, p4)

	p256Sqr(out, out)
	p256Sqr(out, out)
	p256Mul(out, out, p2)

	p256Sqr(out, out)
	p256Sqr(out, out)
	p256Mul(out, out, in)
}

func boothW5(in uint) (int, int) {
	var s uint = ^((in >> 5) - 1)
	var d uint = (1 << 6) - in - 1
	d = (d & s) | (in & (^s))
	d = (d >> 1) + (d & 1)
	return int(d), int(s & 1)
}

func boothW7(in uint) (int, int) {
	var s uint = ^((in >> 7) - 1)
	var d uint = (1 << 8) - in - 1
	d = (d & s) | (in & (^s))
	d = (d >> 1) + (d & 1)
	return int(d), int(s & 1)
}

func initTable() {
	p256Precomputed = new([37][64]p256Point) //z coordinate not used
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}

	t1 := new(p256Point)
	t2 := new(p256Point)
	*t2 = basePoint

	zInv := make([]byte, 32)
	zInvSq := make([]byte, 32)
	for j := 0; j < 64; j++ {
		*t1 = *t2
		for i := 0; i < 37; i++ {
			// The window size is 7 so we need to double 7 times.
			if i != 0 {
				for k := 0; k < 7; k++ {
					p256PointDoubleAsm(t1, t1)
				}
			}
			// Convert the point to affine form. (Its values are
			// still in Montgomery form however.)
			p256Inverse(zInv, t1.z[:])
			p256Sqr(zInvSq, zInv)
			p256Mul(zInv, zInv, zInvSq)

			p256Mul(t1.x[:], t1.x[:], zInvSq)
			p256Mul(t1.y[:], t1.y[:], zInv)

			copy(t1.z[:], basePoint.z[:])
			// Update the table entry
			copy(p256Precomputed[i][j].x[:], t1.x[:])
			copy(p256Precomputed[i][j].y[:], t1.y[:])
		}
		if j == 0 {
			p256PointDoubleAsm(t2, &basePoint)
		} else {
			p256PointAddAsm(t2, t2, &basePoint)
		}
	}
}

func (p *p256Point) p256BaseMult(scalar []byte) {
	precomputeOnce.Do(initTable)
	//fmt.Printf("%s\n", new(big.Int).SetBytes(scalar).Text(16))
	wvalue := (uint(scalar[31]) << 1) & 0xff
	sel, sign := boothW7(uint(wvalue))
	//fmt.Printf("sel %d sign %d wval %d\n", sel, sign, wvalue)
	p256SelectBase(p, p256Precomputed[0][:], sel)
	for i := 0; i < 37; i++ {
		//PrintPoint("p", &p256Precomputed[0][i])
	}
	
	p256NegCond(p, sign)

	copy(p.z[:], one[:])
	//PrintPoint("p", p)
	var t0 p256Point
	copy(t0.z[:], one[:])

	index := uint(6)
	zero := sel

	for i := 1; i < 37; i++ {
		if index < 247 {
			wvalue = ((uint(scalar[31-index/8]) >> (index % 8)) + (uint(scalar[31-index/8-1]) << (8 - (index % 8)))) & 0xff
		} else {
			wvalue = (uint(scalar[31-index/8]) >> (index % 8)) & 0xff
		}
		index += 7
		sel, sign = boothW7(uint(wvalue))
		//fmt.Printf("sel %d sign %d zero %d wval %d\n", sel, sign, zero, wvalue)
		p256SelectBase(&t0, p256Precomputed[i][:], sel)
		//PrintPoint("t0", &t0)
		p256PointAddAffineAsm(p, p, &t0, sign, sel, zero)
		//PrintPoint("p", p)
		zero |= sel
	}
}

func (p *p256Point) p256ScalarMult(scalar []byte) {
	// precomp is a table of precomputed points that stores powers of p
	// from p^1 to p^16.
	//var precomp [16 * 4 * 3]uint64
	var precomp [16]p256Point
	var t0, t1, t2, t3 p256Point

	// Prepare the table
	//p.p256StorePoint(&precomp, 0) // 1
	*&precomp[0] = *p
	
	p256PointDoubleAsm(&t0, p)
	p256PointDoubleAsm(&t1, &t0)
	p256PointDoubleAsm(&t2, &t1)
	p256PointDoubleAsm(&t3, &t2)
	*&precomp[1] = t0  // 2
	*&precomp[3] = t1  // 4
	*&precomp[7] = t2  // 8
	*&precomp[15]= t3  // 16

	p256PointAddAsm(&t0, &t0, p)
	p256PointAddAsm(&t1, &t1, p)
	p256PointAddAsm(&t2, &t2, p)
	*&precomp[2] = t0  // 3
	*&precomp[4] = t1  // 5
	*&precomp[8] = t2  // 9

	p256PointDoubleAsm(&t0, &t0)
	p256PointDoubleAsm(&t1, &t1)
	*&precomp[5] = t0  // 6
	*&precomp[9] = t1  // 10

	p256PointAddAsm(&t2, &t0, p)
	p256PointAddAsm(&t1, &t1, p)
	*&precomp[6] = t2  // 7
	*&precomp[10]= t1  // 11

	p256PointDoubleAsm(&t0, &t0)
	p256PointDoubleAsm(&t2, &t2)
	*&precomp[11] = t0  // 12
	*&precomp[13] = t2  // 14

	p256PointAddAsm(&t0, &t0, p)
	p256PointAddAsm(&t2, &t2, p)
	*&precomp[12] = t0  // 13
	*&precomp[14] = t2  // 15
	//PrintPoint("t2", &t2)
	// Start scanning the window from top bit
	index := uint(254)
	var sel, sign int

	wvalue := (uint(scalar[31-index/8]) >> (index % 8)) & 0x3f
	sel, _ = boothW5(uint(wvalue))
	//PrintPoint("p", p)
	p256Select(p, precomp[:], sel)
	zero := sel
	//PrintPoint("p", p)
	//fmt.Printf("sel %d sign %d zero %d wval %d\n", sel, sign, zero, wvalue)
	for index > 4 {
		index -= 5
		//PrintPoint("p", p)
		p256PointDoubleAsm(p, p)
		//PrintPoint("p", p)
		p256PointDoubleAsm(p, p)
		p256PointDoubleAsm(p, p)
		p256PointDoubleAsm(p, p)
		p256PointDoubleAsm(p, p)
		//PrintPoint("p", p)
		if index < 247 {
			wvalue = ((uint(scalar[31-index/8]) >> (index % 8)) + (uint(scalar[31-index/8-1]) << (8 - (index % 8)))) & 0x3f
		} else {
			wvalue = (uint(scalar[31-index/8]) >> (index % 8)) & 0x3f
		}

		sel, sign = boothW5(uint(wvalue))

		p256Select(&t0, precomp[:], sel)
		p256NegCond(&t0, sign)
		p256PointAddAsm(&t1, p, &t0)
		p256MovCond(&t1, &t1, p, sel)
		p256MovCond(p, &t1, &t0, zero)
		zero |= sel
	}

	p256PointDoubleAsm(p, p)
	p256PointDoubleAsm(p, p)
	p256PointDoubleAsm(p, p)
	p256PointDoubleAsm(p, p)
	p256PointDoubleAsm(p, p)

	wvalue = (uint(scalar[31]) << 1) & 0x3f
	sel, sign = boothW5(uint(wvalue))

	p256Select(&t0, precomp[:], sel)
	p256NegCond(&t0, sign)
	p256PointAddAsm(&t1, p, &t0)
	p256MovCond(&t1, &t1, p, sel)
	p256MovCond(p, &t1, &t0, zero)
}
