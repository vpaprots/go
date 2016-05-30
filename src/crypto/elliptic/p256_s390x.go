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
		
		xyz [12]uint64 //delete
	}
)

var (
	p256            p256Curve
	p256Precomputed *[37][64 * 8]uint64
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
	resx := make([]byte, 32)
	resy := make([]byte, 32)
	resz := make([]byte, 32)
	x := fromBig(x1)
	y := fromBig(y1)
	z := fromBig(z1)
	p256PointDoubleAsm(resx, resy, resz, x, y, z)
	return new(big.Int).SetBytes(resx),new(big.Int).SetBytes(resy), new(big.Int).SetBytes(resz)  
	
}

func (curve p256Curve) TestAdd(x1, y1, z1, x2, y2, z2 *big.Int) (x3,y3,z3 *big.Int) {
	resx := make([]byte, 32)
	resy := make([]byte, 32)
	resz := make([]byte, 32)
	xx1 := fromBig(x1)
	yy1 := fromBig(y1)
	zz1 := fromBig(z1)
	xx2 := fromBig(x2)
	yy2 := fromBig(y2)
	zz2 := fromBig(z2)
	p256PointAddAsm(resx, resy, resz, xx1, yy1, zz1, xx2, yy2, zz2)
	return new(big.Int).SetBytes(resx),new(big.Int).SetBytes(resy), new(big.Int).SetBytes(resz)  
	
}

func (curve p256Curve) TestMul(k *big.Int) *big.Int {
	res := make([]byte, 32)
	x := fromBig(k)
	p256Inverse(res, x)
	return new(big.Int).SetBytes(res)
	
}

// Functions implemented in p256_asm_s390x.s
// Montgomery multiplication modulo P256
func p256Mul(res, in1, in2 []byte){
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
func p256FromMont(res, in []uint64){
	
}

// iff cond == 1  val <- -val
func p256NegCond(val []uint64, cond int){
	
}

// if cond == 0 res <- b; else res <- a
func p256MovCond(res, a, b []uint64, cond int){
	
}

// Endianess swap
func p256BigToLittle(res []uint64, in []byte){
	
}
func p256LittleToBig(res []byte, in []uint64) {
	
}

// Constant time table access
func p256Select(point, table []uint64, idx int) {
	
}
func p256SelectBase(point, table []uint64, idx int) {
	
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

func (curve p256Curve) TestOrdMul() {
	res := make([]byte, 32)
	x1, _ := new(big.Int).SetString("a007c8559316f82de3d5d9f28b8ffcdf5949bd551f7a1348b8acc00860e058", 16)
	x2, _ := new(big.Int).SetString("66e12d94f3d956202845b2392b6bec594699799c49bd6fa683244c95be79eea2", 16)
	Rinv, _ := new(big.Int).SetString("60d066334905c1e907f8b6041e607725badef3e243566fafce1bc8f79c197c79", 16)
	temp := new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv)
	copy(res, fromBig(new(big.Int).Mod(temp, p256.N)))
	
	t := make([]byte, 32)
	p256OrdMul(t, fromBig(x1), fromBig(x2))
	fmt.Printf("-TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
	fmt.Printf("-TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
	fmt.Printf("-EXPECTED %s\n", new(big.Int).SetBytes(res).Text(16))
	fmt.Printf("-FOUND    %s\n", new(big.Int).SetBytes(t).Text(16))
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

// Point add with in2 being affine point
// If sign == 1 -> in2 = -in2
// If sel == 0 -> res = in1
// if zero == 0 -> res = in2
func p256PointAddAffineAsm(res, in1, in2 []uint64, sign, sel, zero int) {
	
}

// Point add
func p256PointAddAsm(X3, Y3, Z3, X1, Y1, Z1, X2, Y2, Z2 []byte) {
	/*
	 * https://choucroutage.com/Papers/SideChannelAttacks/ctrsa-2011-brown.pdf "Software Implementation of the NIST Elliptic Curves Over Prime Fields"
	 *
	 * A = X₁×Z₂²
	 * B = Y₁×Z₂³
	 * C = X₂×Z₁²-A
	 * D = Y₂×Z₁³-B
	 * X₃ = D² - 2A×C² - C³
	 * Y₃ = D×(A×C² - X₃) - B×C³
	 * Z₃ = Z₁×Z₂×C
	 *
 	 * Three-operand formula (adopted): http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-1998-cmo-2
	 * Temp storage: T1,T2,U1,H,Z3=X3=Y3,S1,R
	 * 
	 * T1 = Z1*Z1
	 * T2 = Z2*Z2
	 * U1 = X1*T2
	 * H  = X2*T1
	 * H  = H-U1
	 * Z3 = Z1*Z2
	 * Z3 = Z3*H << store-out Z3 result reg
	 * 
	 * S1 = Z2*T2
	 * S1 = Y1*S1
	 * R  = Z1*T1
	 * R  = Y2*R
	 * R  = R-S1
	 *
	 * T1 = H*H
	 * T2 = H*T1
	 * U1 = U1*T1
	 * 
	 * X3 = R*R
	 * X3 = X3-T2
	 * T1 = 2*U1
	 * X3 = X3-T1 << store-out X3 result reg
	 * 
	 * T2 = S1*T2
	 * Y3 = U1-X3
	 * Y3 = R*Y3
	 * Y3 = Y3-T2 << store-out X3 result reg
	 */
	
	// Note: This test code was not meant to be pretty! It is written in this convoluted fashion to help debug the real assembly code
	T1 := make([]byte, 32)
	T2 := make([]byte, 32)
	U1 := make([]byte, 32)
	S1 := make([]byte, 32)
	H := make([]byte, 32)
	R := make([]byte, 32)
	
	p256Mul(T1, Z1, Z1) // T1 = Z1*Z1
	p256Mul(T2, Z2, Z2) // T2 = Z2*Z2
	p256Mul(U1, X1, T2) // U1 = X1*T2
	p256Mul( H, X2, T1) // H  = X2*T1
	copy(H, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(H), new(big.Int).SetBytes(U1)), p256.P))) // H  = H-U1
	p256Mul(Z3, Z1, Z2) // Z3 = Z1*Z2
	p256Mul(Z3, Z3,  H) // Z3 = Z3*H << store-out Z3 result reg
	
	p256Mul(S1, Z2, T2) // S1 = Z2*T2
	p256Mul(S1, Y1, S1) // S1 = Y1*S1
	p256Mul( R, Z1, T1) // R  = Z1*T1
	p256Mul( R, Y2,  R) // R  = Y2*R
	copy( R, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(R), new(big.Int).SetBytes(S1)), p256.P))) // R  = R-S1
	
	p256Mul(T1,  H,  H) // T1 = H*H
	p256Mul(T2,  H, T1) // T2 = H*T1
	p256Mul(U1, U1, T1) // U1 = U1*T1
	
	p256Mul(X3,  R,  R) // X3 = R*R
	//fmt.Printf(" --TEST R^2: %s\n", new(big.Int).SetBytes(X3).Text(16),)
	//fmt.Printf(" --TEST  T2: %s\n", new(big.Int).SetBytes(T2).Text(16),)
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T2)), p256.P))) // X3 = X3-T2
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(U1), new(big.Int).SetBytes(U1)), p256.P))) // T1 = 2*U1
	//fmt.Printf(" --TEST  T1: %s\n", new(big.Int).SetBytes(T1).Text(16),)
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T1)), p256.P))) // X3 = X3-T1 << store-out X3 result reg
	//fmt.Printf(" --TEST  X3: %s\n", new(big.Int).SetBytes(X3).Text(16),)
	
	p256Mul(T2, S1, T2) // T2 = S1*T2
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(U1), new(big.Int).SetBytes(X3)), p256.P))) // Y3 = U1-X3
	p256Mul(Y3, R, Y3) // Y3 = R*Y3
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(Y3), new(big.Int).SetBytes(T2)), p256.P))) // Y3 = Y3-T2 << store-out X3 result reg
}

// Point double
//http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian.html#doubling-dbl-2007-bl
//http://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
//http://www.hyperelliptic.org/EFD/g1p/auto-shortw-projective-3.html
func p256PointDoubleAsm(X3, Y3, Z3, X1, Y1, Z1 []byte) {
	/*
	 * http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2004-hmv	
	 * Cost: 4M + 4S + 1*half + 5add + 2*2 + 1*3.
	 * Source: 2004 Hankerson–Menezes–Vanstone, page 91.
	 * 	A  = 3(X₁-Z₁²)×(X₁+Z₁²)
	 * 	B  = 2Y₁
	 * 	Z₃ = B×Z₁
	 * 	C  = B²
	 * 	D  = C×X₁
	 * 	X₃ = A²-2D
	 * 	Y₃ = (D-X₃)×A-C²/2
	 * 
	 * Three-operand formula:
	 *       T1 = Z1²
	 *       T2 = X1-T1
	 *       T1 = X1+T1
	 *       T2 = T2*T1
	 *       T2 = 3*T2
	 *       Y3 = 2*Y1
	 *       Z3 = Y3*Z1
	 *       Y3 = Y3²
	 *       T3 = Y3*X1
	 *       Y3 = Y3²
	 *       Y3 = half*Y3
	 *       X3 = T2²
	 *       T1 = 2*T3
	 *       X3 = X3-T1
	 *       T1 = T3-X3
	 *       T1 = T1*T2
	 *       Y3 = T1-Y3
	 */
	
	// Note: This test code was not meant to be pretty! It is written in this convoluted fashion to help debug the real assembly code
	T1 := make([]byte, 32)
	T2 := make([]byte, 32)
	T3 := make([]byte, 32)
	
	p256Mul(T1, Z1, Z1) //T1 = Z1²
	copy(T2, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X1), new(big.Int).SetBytes(T1)), p256.P))) //T2 = X1-T1
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(X1), new(big.Int).SetBytes(T1)), p256.P))) //T1 = X1+T1
	p256Mul(T2, T2, T1) //T2 = T2*T1
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(T2), new(big.Int).SetBytes(T2)), p256.P))) //T2 = 3*T2
	copy(T2, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(T1), new(big.Int).SetBytes(T2)), p256.P)))
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(Y1), new(big.Int).SetBytes(Y1)), p256.P))) // Y3 = 2*Y1
	p256Mul(Z3, Y3, Z1) // Z3 = Y3*Z1
	p256Mul(Y3, Y3, Y3) // Y3 = Y3²
	p256Mul(T3, Y3, X1) // T3 = Y3*X1
	p256Mul(Y3, Y3, Y3) // Y3 = Y3²
	if (1 == Y3[31]&0x01) { // Y3 = half*Y3
		copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Rsh(new(big.Int).Add(new(big.Int).SetBytes(Y3), p256.P), 1), p256.P)))
	} else {
		copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Rsh(new(big.Int).SetBytes(Y3), 1), p256.P)))
	}
	p256Mul(X3, T2, T2) // X3 = T2²
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(T3), new(big.Int).SetBytes(T3)), p256.P))) // T1 = 2*T3
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T1)), p256.P))) // X3 = X3-T1
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T3), new(big.Int).SetBytes(X3)), p256.P))) // T1 = T3-X3
	p256Mul(T1, T1, T2) // T1 = T1*T2
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T1), new(big.Int).SetBytes(Y3)), p256.P))) // Y3 = T1-Y3
}

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

// p256GetScalar endian-swaps the big-endian scalar value from in and writes it
// to out. If the scalar is equal or greater than the order of the group, it's
// reduced modulo that order.
func p256GetScalar(out []uint64, in []byte) {
	n := new(big.Int).SetBytes(in)

	if n.Cmp(p256.N) >= 0 {
		n.Mod(n, p256.N)
	}
	fromBigDel(out, n)
}

// p256Mul operates in a Montgomery domain with R = 2^256 mod p, where p is the
// underlying field of the curve. (See initP256 for the value.) Thus rr here is
// R×R mod p. See comment in Inverse about how this is used.
var rr = []byte{  0x00, 0x00, 0x00, 0x04, 0xff, 0xff, 0xff, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
	              0xff, 0xff, 0xff, 0xfb, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}

func maybeReduceModP(in *big.Int) *big.Int {
	if in.Cmp(p256.P) < 0 {
		return in
	}
	return new(big.Int).Mod(in, p256.P)
}

/*func (curve p256Curve) CombinedMult(bigX, bigY *big.Int, baseScalar, scalar []byte) (x, y *big.Int) {
	scalarReversed := make([]uint64, 4)
	var r1, r2 p256Point
	p256GetScalar(scalarReversed, baseScalar)
	r1.p256BaseMult(scalarReversed)

	p256GetScalar(scalarReversed, scalar)
	fromBigDel(r2.xyz[0:4], maybeReduceModP(bigX))
	fromBigDel(r2.xyz[4:8], maybeReduceModP(bigY))
	p256Mul(r2.xyz[0:4], r2.xyz[0:4], rr[:])
	p256Mul(r2.xyz[4:8], r2.xyz[4:8], rr[:])

	// This sets r2's Z value to 1, in the Montgomery domain.
	r2.xyz[8] = 0x0000000000000001
	r2.xyz[9] = 0xffffffff00000000
	r2.xyz[10] = 0xffffffffffffffff
	r2.xyz[11] = 0x00000000fffffffe

	r2.p256ScalarMult(scalarReversed)
	p256PointAddAsm(r1.xyz[:], r1.xyz[:], r2.xyz[:])
	return r1.p256PointToAffine()
}

func (curve p256Curve) ScalarBaseMult(scalar []byte) (x, y *big.Int) {
	scalarReversed := make([]uint64, 4)
	p256GetScalar(scalarReversed, scalar)

	var r p256Point
	r.p256BaseMult(scalarReversed)
	return r.p256PointToAffine()
}

func (curve p256Curve) ScalarMult(bigX, bigY *big.Int, scalar []byte) (x, y *big.Int) {
	scalarReversed := make([]uint64, 4)
	p256GetScalar(scalarReversed, scalar)

	var r p256Point
	fromBigDel(r.xyz[0:4], maybeReduceModP(bigX))
	fromBigDel(r.xyz[4:8], maybeReduceModP(bigY))
	p256Mul(r.xyz[0:4], r.xyz[0:4], rr[:])
	p256Mul(r.xyz[4:8], r.xyz[4:8], rr[:])
	// This sets r2's Z value to 1, in the Montgomery domain.
	r.xyz[8] = 0x0000000000000001
	r.xyz[9] = 0xffffffff00000000
	r.xyz[10] = 0xffffffffffffffff
	r.xyz[11] = 0x00000000fffffffe

	r.p256ScalarMult(scalarReversed)
	return r.p256PointToAffine()
}

func (p *p256Point) p256PointToAffine() (x, y *big.Int) {
	zInv := make([]uint64, 4)
	zInvSq := make([]uint64, 4)
	p256Inverse(zInv, p.xyz[8:12])
	p256Sqr(zInvSq, zInv)
	p256Mul(zInv, zInv, zInvSq)

	p256Mul(zInvSq, p.xyz[0:4], zInvSq)
	p256Mul(zInv, p.xyz[4:8], zInv)

	p256FromMont(zInvSq, zInvSq)
	p256FromMont(zInv, zInv)

	xOut := make([]byte, 32)
	yOut := make([]byte, 32)
	p256LittleToBig(xOut, zInvSq)
	p256LittleToBig(yOut, zInv)

	return new(big.Int).SetBytes(xOut), new(big.Int).SetBytes(yOut)
}*/

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
	
	fmt.Printf("-TEST in  %s\n", new(big.Int).SetBytes(in).Text(16))
	fmt.Printf("-TEST out %s\n", new(big.Int).SetBytes(out).Text(16))
	fmt.Printf("-TEST p2  %s\n", new(big.Int).SetBytes(p2).Text(16))
	fmt.Printf("-TEST p4  %s\n", new(big.Int).SetBytes(p4).Text(16))
	fmt.Printf("-TEST p8  %s\n", new(big.Int).SetBytes(p8).Text(16))
	fmt.Printf("-TEST p16 %s\n", new(big.Int).SetBytes(p16).Text(16))
	fmt.Printf("-TEST p32 %s\n", new(big.Int).SetBytes(p32).Text(16))
}

/*func (p *p256Point) p256StorePoint(r *[16 * 4 * 3]uint64, index int) {
	copy(r[index*12:], p.xyz[:])
}*/

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

/*func initTable() {
	p256Precomputed = new([37][64 * 8]uint64)

	basePoint := []uint64{
		0x79e730d418a9143c, 0x75ba95fc5fedb601, 0x79fb732b77622510, 0x18905f76a53755c6, //(p256.x1*2^256)%p
		0xddf25357ce95560a, 0x8b4ab8e4ba19e45c, 0xd2e88688dd21f325, 0x8571ff1825885d85,
		0x0000000000000001, 0xffffffff00000000, 0xffffffffffffffff, 0x00000000fffffffe,
	}
	t1 := make([]uint64, 12)
	t2 := make([]uint64, 12)
	copy(t2, basePoint)

	zInv := make([]uint64, 4)
	zInvSq := make([]uint64, 4)
	for j := 0; j < 64; j++ {
		copy(t1, t2)
		for i := 0; i < 37; i++ {
			// The window size is 7 so we need to double 7 times.
			if i != 0 {
				for k := 0; k < 7; k++ {
					p256PointDoubleAsm(t1, t1)
				}
			}
			// Convert the point to affine form. (Its values are
			// still in Montgomery form however.)
			p256Inverse(zInv, t1[8:12])
			p256Sqr(zInvSq, zInv)
			p256Mul(zInv, zInv, zInvSq)

			p256Mul(t1[:4], t1[:4], zInvSq)
			p256Mul(t1[4:8], t1[4:8], zInv)

			copy(t1[8:12], basePoint[8:12])
			// Update the table entry
			copy(p256Precomputed[i][j*8:], t1[:8])
		}
		if j == 0 {
			p256PointDoubleAsm(t2, basePoint)
		} else {
			p256PointAddAsm(t2, t2, basePoint)
		}
	}
}

/*func (p *p256Point) p256BaseMult(scalar []uint64) {
	precomputeOnce.Do(initTable)

	wvalue := (scalar[0] << 1) & 0xff
	sel, sign := boothW7(uint(wvalue))
	p256SelectBase(p.xyz[0:8], p256Precomputed[0][0:], sel)
	p256NegCond(p.xyz[4:8], sign)

	// (This is one, in the Montgomery domain.)
	p.xyz[8] = 0x0000000000000001
	p.xyz[9] = 0xffffffff00000000
	p.xyz[10] = 0xffffffffffffffff
	p.xyz[11] = 0x00000000fffffffe

	var t0 p256Point
	// (This is one, in the Montgomery domain.)
	t0.xyz[8] = 0x0000000000000001
	t0.xyz[9] = 0xffffffff00000000
	t0.xyz[10] = 0xffffffffffffffff
	t0.xyz[11] = 0x00000000fffffffe

	index := uint(6)
	zero := sel

	for i := 1; i < 37; i++ {
		if index < 192 {
			wvalue = ((scalar[index/64] >> (index % 64)) + (scalar[index/64+1] << (64 - (index % 64)))) & 0xff
		} else {
			wvalue = (scalar[index/64] >> (index % 64)) & 0xff
		}
		index += 7
		sel, sign = boothW7(uint(wvalue))
		p256SelectBase(t0.xyz[0:8], p256Precomputed[i][0:], sel)
		p256PointAddAffineAsm(p.xyz[0:12], p.xyz[0:12], t0.xyz[0:8], sign, sel, zero)
		zero |= sel
	}
}

func (p *p256Point) p256ScalarMult(scalar []uint64) {
	// precomp is a table of precomputed points that stores powers of p
	// from p^1 to p^16.
	var precomp [16 * 4 * 3]uint64
	var t0, t1, t2, t3 p256Point

	// Prepare the table
	p.p256StorePoint(&precomp, 0) // 1

	p256PointDoubleAsm(t0.xyz[:], p.xyz[:])
	p256PointDoubleAsm(t1.xyz[:], t0.xyz[:])
	p256PointDoubleAsm(t2.xyz[:], t1.xyz[:])
	p256PointDoubleAsm(t3.xyz[:], t2.xyz[:])
	t0.p256StorePoint(&precomp, 1)  // 2
	t1.p256StorePoint(&precomp, 3)  // 4
	t2.p256StorePoint(&precomp, 7)  // 8
	t3.p256StorePoint(&precomp, 15) // 16

	p256PointAddAsm(t0.xyz[:], t0.xyz[:], p.xyz[:])
	p256PointAddAsm(t1.xyz[:], t1.xyz[:], p.xyz[:])
	p256PointAddAsm(t2.xyz[:], t2.xyz[:], p.xyz[:])
	t0.p256StorePoint(&precomp, 2) // 3
	t1.p256StorePoint(&precomp, 4) // 5
	t2.p256StorePoint(&precomp, 8) // 9

	p256PointDoubleAsm(t0.xyz[:], t0.xyz[:])
	p256PointDoubleAsm(t1.xyz[:], t1.xyz[:])
	t0.p256StorePoint(&precomp, 5) // 6
	t1.p256StorePoint(&precomp, 9) // 10

	p256PointAddAsm(t2.xyz[:], t0.xyz[:], p.xyz[:])
	p256PointAddAsm(t1.xyz[:], t1.xyz[:], p.xyz[:])
	t2.p256StorePoint(&precomp, 6)  // 7
	t1.p256StorePoint(&precomp, 10) // 11

	p256PointDoubleAsm(t0.xyz[:], t0.xyz[:])
	p256PointDoubleAsm(t2.xyz[:], t2.xyz[:])
	t0.p256StorePoint(&precomp, 11) // 12
	t2.p256StorePoint(&precomp, 13) // 14

	p256PointAddAsm(t0.xyz[:], t0.xyz[:], p.xyz[:])
	p256PointAddAsm(t2.xyz[:], t2.xyz[:], p.xyz[:])
	t0.p256StorePoint(&precomp, 12) // 13
	t2.p256StorePoint(&precomp, 14) // 15

	// Start scanning the window from top bit
	index := uint(254)
	var sel, sign int

	wvalue := (scalar[index/64] >> (index % 64)) & 0x3f
	sel, _ = boothW5(uint(wvalue))

	p256Select(p.xyz[0:12], precomp[0:], sel)
	zero := sel

	for index > 4 {
		index -= 5
		p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		p256PointDoubleAsm(p.xyz[:], p.xyz[:])
		p256PointDoubleAsm(p.xyz[:], p.xyz[:])

		if index < 192 {
			wvalue = ((scalar[index/64] >> (index % 64)) + (scalar[index/64+1] << (64 - (index % 64)))) & 0x3f
		} else {
			wvalue = (scalar[index/64] >> (index % 64)) & 0x3f
		}

		sel, sign = boothW5(uint(wvalue))

		p256Select(t0.xyz[0:], precomp[0:], sel)
		p256NegCond(t0.xyz[4:8], sign)
		p256PointAddAsm(t1.xyz[:], p.xyz[:], t0.xyz[:])
		p256MovCond(t1.xyz[0:12], t1.xyz[0:12], p.xyz[0:12], sel)
		p256MovCond(p.xyz[0:12], t1.xyz[0:12], t0.xyz[0:12], zero)
		zero |= sel
	}

	p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	p256PointDoubleAsm(p.xyz[:], p.xyz[:])
	p256PointDoubleAsm(p.xyz[:], p.xyz[:])

	wvalue = (scalar[0] << 1) & 0x3f
	sel, sign = boothW5(uint(wvalue))

	p256Select(t0.xyz[0:], precomp[0:], sel)
	p256NegCond(t0.xyz[4:8], sign)
	p256PointAddAsm(t1.xyz[:], p.xyz[:], t0.xyz[:])
	p256MovCond(t1.xyz[0:12], t1.xyz[0:12], p.xyz[0:12], sel)
	p256MovCond(p.xyz[0:12], t1.xyz[0:12], t0.xyz[0:12], zero)
}
*/