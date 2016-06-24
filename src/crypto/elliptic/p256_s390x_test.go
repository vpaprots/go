// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package elliptic

import (
	"testing"
	"fmt"
	"math/big"
	"crypto/rand"
	"bytes"
)

/*
X=Z1; Y=Z1; MUL;T-   // T1 = Z1²      T1
X=T ; Y-  ; MUL;T2=T // T2 = T1*Z1    T1   T2
X-  ; Y=X2; MUL;T1=T // T1 = T1*X2    T1   T2
X=T2; Y=Y2; MUL;T-   // T2 = T2*Y2    T1   T2
SUB(T2<T-Y1)         // T2 = T2-Y1    T1   T2
SUB(Y<T1-X1)         // T1 = T1-X1    T1   T2
X=Z1; Y- ;  MUL;Z3:=T// Z3 = Z1*T1         T2
X=Y;  Y- ;  MUL;X=T  // T3 = T1*T1         T2
X- ;  Y- ;  MUL;T4=T // T4 = T3*T1         T2        T4
X- ;  Y=X1; MUL;T3=T // T3 = T3*X1         T2   T3   T4
ADD(T1<T+T)          // T1 = T3+T3    T1   T2   T3   T4
X=T2; Y=T2; MUL;T-   // X3 = T2*T2    T1   T2   T3   T4
SUB(T<T-T1)          // X3 = X3-T1    T1   T2   T3   T4
SUB(T<T-T4) X3:=T    // X3 = X3-T4         T2   T3   T4
SUB(X<T3-T)          // T3 = T3-X3         T2   T3   T4
X- ;  Y- ;  MUL;T3=T // T3 = T3*T2         T2   T3   T4
X=T4; Y=Y1; MUL;T-   // T4 = T4*Y1              T3   T4
SUB(T<T3-T) Y3:=T    // Y3 = T3-T4              T3   T4
*/
// Point add with P2 being affine point
// If sign == 1 -> P2 = -P2
// If sel == 0 -> P3 = P1
// if zero == 0 -> P3 = P2
func p256PointAddAffineAsmBig(P3, P1, P2 *p256Point, sign, sel, zero int) {
	/**
	 * Three operand formula:
	 * Source: 2004 Hankerson–Menezes–Vanstone, page 91. 
	 * T1 = Z1²
     * T2 = T1*Z1
     * T1 = T1*X2
     * T2 = T2*Y2
     * T1 = T1-X1
     * T2 = T2-Y1
     * Z3 = Z1*T1
     * T3 = T1²
     * T4 = T3*T1
     * T3 = T3*X1
     * T1 = 2*T3
     * X3 = T2²
     * X3 = X3-T1
     * X3 = X3-T4
     * T3 = T3-X3
     * T3 = T3*T2
     * T4 = T4*Y1
     * Y3 = T3-T4
	 */
	
	X1 := P1.x[:]
	Y1 := P1.y[:]
	Z1 := P1.z[:]
	X2 := P2.x[:]
	Y2 := P2.y[:]
	X3 := make([]byte, 32) //P3.x[:]
	Y3 := make([]byte, 32) //P3.y[:]
	Z3 := make([]byte, 32) //P3.z[:]
	
	T1 := make([]byte, 32)
	T2 := make([]byte, 32)
	T3 := make([]byte, 32)
	T4 := make([]byte, 32)
	
	if (sign == 1) {
		Y2 = fromBig(new(big.Int).Mod(new(big.Int).Sub(p256Params.P, new(big.Int).SetBytes(Y2)), p256Params.P)) // Y2  = P-Y2
	}

	p256MulAsm(T1, Z1, Z1) // T1 = Z1²
	//fmt.Printf(" --T1 = Z1²  : %s\n", new(big.Int).SetBytes(T1).Text(16))
	p256MulAsm(T2, T1, Z1) // T2 = T1*Z1
	//fmt.Printf(" --T2 = T1*Z1: %s\n", new(big.Int).SetBytes(T2).Text(16))
	p256MulAsm(T1, T1, X2) // T1 = T1*X2
	//fmt.Printf(" --T1 = T1*X2: %s\n", new(big.Int).SetBytes(T1).Text(16))
	p256MulAsm(T2, T2, Y2) // T2 = T2*Y2
	//fmt.Printf(" --T2 = T2*Y2: %s\n", new(big.Int).SetBytes(T2).Text(16))
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T1), new(big.Int).SetBytes(X1)), p256Params.P))) // T1 = T1-X1
	//fmt.Printf(" --T1 = T1-X1: %s\n", new(big.Int).SetBytes(T1).Text(16))
	copy(T2, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T2), new(big.Int).SetBytes(Y1)), p256Params.P))) // T2 = T2-Y1
	//fmt.Printf(" --T2 = T2-Y1: %s\n", new(big.Int).SetBytes(T2).Text(16))
	p256MulAsm(Z3, Z1, T1) // Z3 = Z1*T1
	p256MulAsm(T3, T1, T1) // T3 = T1²
	p256MulAsm(T4, T3, T1) // T4 = T3*T1
	p256MulAsm(T3, T3, X1) // T3 = T3*X1
	//fmt.Printf(" --T3 = T3*X1: %s\n", new(big.Int).SetBytes(T3).Text(16))
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(T3), new(big.Int).SetBytes(T3)), p256Params.P))) // T1 = 2*T3
	//fmt.Printf(" --T1 = 2*T3:  %s\n", new(big.Int).SetBytes(T1).Text(16))
	p256MulAsm(X3, T2, T2) // X3 = T2²
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T1)), p256Params.P))) // X3 = X3-T1
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T4)), p256Params.P))) // X3 = X3-T4
	copy(T3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T3), new(big.Int).SetBytes(X3)), p256Params.P))) // T3 = T3-X3
	p256MulAsm(T3, T3, T2) // T3 = T3*T2
	p256MulAsm(T4, T4, Y1) // T4 = T4*Y1
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T3), new(big.Int).SetBytes(T4)), p256Params.P))) // Y3 = T3-T4
	
	if (sel == 0) {
		copy(P3.x[:], X1)
		copy(P3.y[:], Y1)
		copy(P3.z[:], Z1)
	}

	if (zero == 0) {
		copy(P3.x[:], X2)
		copy(P3.y[:], Y2)
		copy(P3.z[:], []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})  //(p256.z*2^256)%p
	}
	if (zero != 0 && sel != 0) {
		copy(P3.x[:], X3)
		copy(P3.y[:], Y3)
		copy(P3.z[:], Z3)
	}
}

// Point double
//http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian.html#doubling-dbl-2007-bl
//http://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
//http://www.hyperelliptic.org/EFD/g1p/auto-shortw-projective-3.html
func p256PointDoubleAsmBig(P3, P1 *p256Point) {
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
	X1 := P1.x[:]
	Y1 := P1.y[:]
	Z1 := P1.z[:]
	X3 := P3.x[:]
	Y3 := P3.y[:]
	Z3 := P3.z[:]
	
	T1 := make([]byte, 32)
	T2 := make([]byte, 32)
	T3 := make([]byte, 32)
	
	p256MulAsm(T1, Z1, Z1) //T1 = Z1²
	//fmt.Printf(" --T1 = Z1²  : %s\n", new(big.Int).SetBytes(T1).Text(16))
	copy(T2, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X1), new(big.Int).SetBytes(T1)), p256Params.P))) //T2 = X1-T1
	//fmt.Printf(" --T2 = X1-T1: %s\n", new(big.Int).SetBytes(T2).Text(16))
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(X1), new(big.Int).SetBytes(T1)), p256Params.P))) //T1 = X1+T1
	//fmt.Printf(" --T1 = X1+T1: %s\n", new(big.Int).SetBytes(T1).Text(16))
	p256MulAsm(T2, T2, T1) //T2 = T2*T1
	//fmt.Printf(" --T2 = T2*T1: %s\n", new(big.Int).SetBytes(T2).Text(16))
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(T2), new(big.Int).SetBytes(T2)), p256Params.P))) //T2 = 3*T2
	copy(T2, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(T1), new(big.Int).SetBytes(T2)), p256Params.P)))
	//fmt.Printf(" --T2 = 3*T2 : %s\n", new(big.Int).SetBytes(T2).Text(16))
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(Y1), new(big.Int).SetBytes(Y1)), p256Params.P))) // Y3 = 2*Y1
	//fmt.Printf(" --Y3 = 2*Y1 : %s\n", new(big.Int).SetBytes(Y3).Text(16))
	p256MulAsm(Z3, Y3, Z1) // Z3 = Y3*Z1
	//fmt.Printf(" --Z3 = Y3*Z1: %s\n", new(big.Int).SetBytes(Z3).Text(16))
	p256MulAsm(Y3, Y3, Y3) // Y3 = Y3²
	//fmt.Printf(" --Y3 = Y3²  : %s\n", new(big.Int).SetBytes(Y3).Text(16))
	p256MulAsm(T3, Y3, X1) // T3 = Y3*X1
	//fmt.Printf(" --T3 = Y3*X1: %s\n", new(big.Int).SetBytes(T3).Text(16))
	p256MulAsm(Y3, Y3, Y3) // Y3 = Y3²
	//fmt.Printf(" --Y3 = Y3²  : %s\n", new(big.Int).SetBytes(Y3).Text(16))
	if (1 == Y3[31]&0x01) { // Y3 = half*Y3
		copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Rsh(new(big.Int).Add(new(big.Int).SetBytes(Y3), p256Params.P), 1), p256Params.P)))
	} else {
		copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Rsh(new(big.Int).SetBytes(Y3), 1), p256Params.P)))
	}
	//fmt.Printf(" --Y3 = hal*Y3: %s\n", new(big.Int).SetBytes(Y3).Text(16))
	p256MulAsm(X3, T2, T2) // X3 = T2²
	//fmt.Printf(" --X3 = T2²   : %s\n", new(big.Int).SetBytes(X3).Text(16))
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(T3), new(big.Int).SetBytes(T3)), p256Params.P))) // T1 = 2*T3
	//fmt.Printf(" --T1 = 2*T3 : %s\n", new(big.Int).SetBytes(T1).Text(16))
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T1)), p256Params.P))) // X3 = X3-T1
	//fmt.Printf(" --X3 = X3-T1: %s\n", new(big.Int).SetBytes(X3).Text(16))
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T3), new(big.Int).SetBytes(X3)), p256Params.P))) // T1 = T3-X3
	//fmt.Printf(" --T1 = T3-X3: %s\n", new(big.Int).SetBytes(T1).Text(16))
	p256MulAsm(T1, T1, T2) // T1 = T1*T2
	//fmt.Printf(" --T1 = T1*T2: %s\n", new(big.Int).SetBytes(T1).Text(16))
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(T1), new(big.Int).SetBytes(Y3)), p256Params.P))) // Y3 = T1-Y3
	//fmt.Printf(" --Y3 = T1-Y3: %s\n", new(big.Int).SetBytes(Y3).Text(16))
}

func p256PointAddAsmBig(P3, P1, P2 *p256Point) {
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
	 * Z3 = Z3*H << store-out Z3 result reg.. could override Z1, if slices have same backing array
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
	 * Y3 = Y3-T2 << store-out Y3 result reg
	 */
	
	// Note: This test code was not meant to be pretty! It is written in this convoluted fashion to help debug the real assembly code
	X1 := P1.x[:]
	Y1 := P1.y[:]
	Z1 := P1.z[:]
	X2 := P2.x[:]
	Y2 := P2.y[:]
	Z2 := P2.z[:]
	X3 := P3.x[:]
	Y3 := P3.y[:]
	Z3 := P3.z[:]
	
	T1 := make([]byte, 32)
	T2 := make([]byte, 32)
	U1 := make([]byte, 32)
	S1 := make([]byte, 32)
	H := make([]byte, 32)
	R := make([]byte, 32)
	
	p256MulAsm(T1, Z1, Z1) // T1 = Z1*Z1
	//fmt.Printf(" --T1 = Z1*Z1: %s\n", new(big.Int).SetBytes(T1).Text(16))
	p256MulAsm( R, Z1, T1) // R  = Z1*T1
	//fmt.Printf(" --R  = Z1*T1: %s\n", new(big.Int).SetBytes(R).Text(16))
	p256MulAsm( H, X2, T1) // H  = X2*T1
	//fmt.Printf(" --H  = X2*T1: %s\n", new(big.Int).SetBytes(H).Text(16))
	p256MulAsm(T2, Z2, Z2) // T2 = Z2*Z2
	//fmt.Printf(" --T2 = Z2*Z2: %s\n", new(big.Int).SetBytes(T2).Text(16))
	p256MulAsm(S1, Z2, T2) // S1 = Z2*T2
	//fmt.Printf(" --S1 = Z2*T2: %s\n", new(big.Int).SetBytes(S1).Text(16))
	p256MulAsm(U1, X1, T2) // U1 = X1*T2
	//fmt.Printf(" --U1 = X1*T2: %s\n", new(big.Int).SetBytes(U1).Text(16))
	
	copy(H, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(H), new(big.Int).SetBytes(U1)), p256Params.P))) // H  = H-U1
	//fmt.Printf(" --H  = H-U1 : %s\n", new(big.Int).SetBytes(H).Text(16))
	p256MulAsm(Z3, Z1, Z2) // Z3 = Z1*Z2
	//fmt.Printf(" --Z3 = Z1*Z2: %s\n", new(big.Int).SetBytes(Z3).Text(16))
	p256MulAsm(Z3, Z3,  H) // Z3 = Z3*H << store-out Z3 result reg
	//fmt.Printf(" --Z3 = Z3*H : %s\n", new(big.Int).SetBytes(Z3).Text(16))
	
	p256MulAsm(S1, Y1, S1) // S1 = Y1*S1
	//fmt.Printf(" --S1 = Y1*S1: %s\n", new(big.Int).SetBytes(S1).Text(16))
	p256MulAsm( R, Y2,  R) // R  = Y2*R
	//fmt.Printf(" --R  = Y2*R : %s\n", new(big.Int).SetBytes(R).Text(16))
	copy( R, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(R), new(big.Int).SetBytes(S1)), p256Params.P))) // R  = R-S1
	//fmt.Printf(" --R  = R-S1 : %s\n", new(big.Int).SetBytes(R).Text(16))
	
	p256MulAsm(T1,  H,  H) // T1 = H*H
	//fmt.Printf(" --T1 = H*H  : %s\n", new(big.Int).SetBytes(T1).Text(16))
	p256MulAsm(T2,  H, T1) // T2 = H*T1
	//fmt.Printf(" --T2 = H*T1 : %s\n", new(big.Int).SetBytes(T2).Text(16))
	p256MulAsm(U1, U1, T1) // U1 = U1*T1
	//fmt.Printf(" --U1 = U1*T1: %s\n", new(big.Int).SetBytes(U1).Text(16))
	
	p256MulAsm(X3,  R,  R) // X3 = R*R
	//fmt.Printf(" --X3 = R*R  : %s\n", new(big.Int).SetBytes(X3).Text(16))
	//fmt.Printf(" --TEST R^2: %s\n", new(big.Int).SetBytes(X3).Text(16),)
	//fmt.Printf(" --TEST  T2: %s\n", new(big.Int).SetBytes(T2).Text(16),)
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T2)), p256Params.P))) // X3 = X3-T2
	//fmt.Printf(" --X3 = X3-T2: %s\n", new(big.Int).SetBytes(X3).Text(16))
	copy(T1, fromBig(new(big.Int).Mod(new(big.Int).Add(new(big.Int).SetBytes(U1), new(big.Int).SetBytes(U1)), p256Params.P))) // T1 = 2*U1
	//fmt.Printf(" --T1 = 2*U1 : %s\n", new(big.Int).SetBytes(T1).Text(16))
	//fmt.Printf(" --TEST  T1: %s\n", new(big.Int).SetBytes(T1).Text(16),)
	copy(X3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(X3), new(big.Int).SetBytes(T1)), p256Params.P))) // X3 = X3-T1 << store-out X3 result reg
	//fmt.Printf(" --X3 = X3-T1: %s\n", new(big.Int).SetBytes(X3).Text(16))
	//fmt.Printf(" --TEST  X3: %s\n", new(big.Int).SetBytes(X3).Text(16),)
	
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(U1), new(big.Int).SetBytes(X3)), p256Params.P))) // Y3 = U1-X3
	//fmt.Printf(" --Y3 = U1-X3: %s\n", new(big.Int).SetBytes(Y3).Text(16))
	p256MulAsm(Y3, R, Y3) // Y3 = R*Y3
	//fmt.Printf(" --Y3 = R*Y3 : %s\n", new(big.Int).SetBytes(Y3).Text(16))
	p256MulAsm(T2, S1, T2) // T2 = S1*T2
	//fmt.Printf(" --T2 = S1*T2: %s\n", new(big.Int).SetBytes(T2).Text(16))
	copy(Y3, fromBig(new(big.Int).Mod(new(big.Int).Sub(new(big.Int).SetBytes(Y3), new(big.Int).SetBytes(T2)), p256Params.P))) // Y3 = Y3-T2 << store-out X3 result reg
	//fmt.Printf(" --Y3 = Y3-T2: %s\n", new(big.Int).SetBytes(Y3).Text(16))
	
	// X=Z1; Y=Z1; MUL; T-   // T1 = Z1*Z1
	// X-  ; Y=T ; MUL; R=T  // R  = Z1*T1
	// X=X2; Y-  ; MUL; H=T  // H  = X2*T1
	// X=Z2; Y=Z2; MUL; T-   // T2 = Z2*Z2
	// X-  ; Y=T ; MUL; S1=T // S1 = Z2*T2
	// X=X1; Y-  ; MUL; U1=T // U1 = X1*T2
	// SUB(H<H-T)            // H  = H-U1
	// X=Z1; Y=Z2; MUL; T-   // Z3 = Z1*Z2
	// X=T ; Y=H ; MUL; Z3:=T// Z3 = Z3*H << store-out Z3 result reg.. could override Z1, if slices have same backing array
	// X=Y1; Y=S1; MUL; S1=T // S1 = Y1*S1
	// X=Y2; Y=R ; MUL; T-   // R  = Y2*R
	// SUB(R<T-S1)           // R  = R-S1
	// X=H ; Y=H ; MUL; T-   // T1 = H*H
	// X-  ; Y=T ; MUL; T2=T // T2 = H*T1
	// X=U1; Y-  ; MUL; U1=T // U1 = U1*T1
	// X=R ; Y=R ; MUL; T-   // X3 = R*R
	// SUB(T<T-T2)           // X3 = X3-T2
	// ADD(X<U1+U1)          // T1 = 2*U1
	// SUB(T<T-X) X3:=T      // X3 = X3-T1 << store-out X3 result reg
	// SUB(Y<U1-T)           // Y3 = U1-X3
	// X=R ; Y-  ; MUL; U1=T // Y3 = R*Y3
	// X=S1; Y=T2; MUL; T-   // T2 = S1*T2            
	// SUB(T<U1-T); Y3:=T    // Y3 = Y3-T2 << store-out Y3 result reg
}

func Testp256MulAsm(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	P256()
	
	exp := make([]byte, 32)
	res := make([]byte, 32)
	x1, _ := new(big.Int).SetString("a007c8559316f82de3d5d9f28b8ffcdf5949bd551f7a1348b8acc00860e058", 16)
	x2, _ := new(big.Int).SetString("66e12d94f3d956202845b2392b6bec594699799c49bd6fa683244c95be79eea2", 16)
	Rinv, _ := new(big.Int).SetString("fffffffe00000003fffffffd0000000200000001fffffffe0000000300000000", 16)
	
	copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256Params.P)))
	p256MulAsm(res, fromBig(x1), fromBig(x2))
	
	if (bytes.Compare(exp,res)!=0) {
		fmt.Printf("-EXPECTED %s\n", new(big.Int).SetBytes(exp).Text(16))
		fmt.Printf("-FOUND    %s\n", new(big.Int).SetBytes(res).Text(16))
		fmt.Printf("-TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
		fmt.Printf("-TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
		t.Fail()
	}
}

func TestStressp256MulAsm(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	P256()
	for i:=0; i<1000000; i++ {
		
		exp := make([]byte, 32)
		res := make([]byte, 32)
		x1, _ := rand.Int(rand.Reader, p256Params.P)
		x2, _ := rand.Int(rand.Reader, p256Params.P)
		Rinv, _ := new(big.Int).SetString("fffffffe00000003fffffffd0000000200000001fffffffe0000000300000000", 16)
		
		copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256Params.P)))
		p256MulAsm(res, fromBig(x1), fromBig(x2))
		
		if (bytes.Compare(exp,res)!=0) {
			fmt.Printf("-EXPECTED %s\n", new(big.Int).SetBytes(exp).Text(16))
			fmt.Printf("-FOUND    %s\n", new(big.Int).SetBytes(res).Text(16))
			fmt.Printf("-TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
			fmt.Printf("-TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
			t.FailNow()
		}
		
		if ( 0 == i%10000) {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")
}

func TestP256OrdMul(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	P256()
	
	exp := make([]byte, 32)
	res := make([]byte, 32)
	x1, _ := new(big.Int).SetString("a007c8559316f82de3d5d9f28b8ffcdf5949bd551f7a1348b8acc00860e058", 16)
	x2, _ := new(big.Int).SetString("66e12d94f3d956202845b2392b6bec594699799c49bd6fa683244c95be79eea2", 16)
	Rinv, _ := new(big.Int).SetString("60d066334905c1e907f8b6041e607725badef3e243566fafce1bc8f79c197c79", 16)
	
	copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256Params.N)))
	p256OrdMul(res, fromBig(x1), fromBig(x2))
	
	if (bytes.Compare(exp,res)!=0) {
		fmt.Printf("-EXPECTED %s\n", new(big.Int).SetBytes(exp).Text(16))
		fmt.Printf("-FOUND    %s\n", new(big.Int).SetBytes(res).Text(16))
		fmt.Printf("-TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
		fmt.Printf("-TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
		t.Fail()
	}
}

func TestStressP256OrdMul(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	P256()
	for i:=0; i<100000; i++ {
		
		exp := make([]byte, 32)
		res := make([]byte, 32)
		x1, _ := rand.Int(rand.Reader, p256Params.N)
		x2, _ := rand.Int(rand.Reader, p256Params.N)
		Rinv, _ := new(big.Int).SetString("60d066334905c1e907f8b6041e607725badef3e243566fafce1bc8f79c197c79", 16)
		
		copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256Params.N)))
		p256OrdMul(res, fromBig(x1), fromBig(x2))
		
		if (bytes.Compare(exp,res)!=0) {
			fmt.Printf("-EXPECTED %s\n", new(big.Int).SetBytes(exp).Text(16))
			fmt.Printf("-FOUND    %s\n", new(big.Int).SetBytes(res).Text(16))
			fmt.Printf("-TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
			fmt.Printf("-TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
			t.FailNow()
		}
		
		if ( 0 == i%1000) {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")
}

func TestStressInverse(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	pp256, _ := P256().(p256CurveFast)
	for i:=0; i<10000; i++ {
		x, _ := rand.Int(rand.Reader, p256Params.N)
		xInv := pp256.Inverse(x)
		xInv2 := pp256.InverseBig(x)
		if (xInv.Cmp(xInv2)!=0) {
			fmt.Printf("EXPECTED: %s\nFOUND:    %s\n", xInv2.String(), xInv.String())
			t.FailNow()
		} 
		if ( 0 == i%100) {
			fmt.Printf(".")
		}
	}
}

func BenchmarkInverse(b *testing.B) {
    pp256, _ := P256().(p256CurveFast)
    x, _ := rand.Int(rand.Reader, p256Params.N)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pp256.Inverse(x)
    }
}

func Benchmarkp256MulAsm(b *testing.B) {
    P256()
    //x, _ := rand.Int(rand.Reader, pp256.N)
    in := make([]byte, 32)
    //copy(in, x.Bytes())
    in[0] = 20
    in[2] = 42
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
		p256MulAsm(in, in, in)
    }
}

func BenchmarkP256Sqr(b *testing.B) {
    P256()
    //x, _ := rand.Int(rand.Reader, pp256.N)
    in := make([]byte, 32)
    //copy(in, x.Bytes())
    in[0] = 20
    in[2] = 42
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
		p256MulAsm(in, in, in)
    }
}

func TestInverse(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	x, _ := new(big.Int).SetString("15792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	xInvExp, _ := new(big.Int).SetString("78239946472340125005789637834181368510016866213607708536507216111758801279690", 10)
	
	pp256, _ := P256().(p256CurveFast)
	xInv := pp256.TestInv(x)
	
	if (xInv.Cmp(xInvExp)!=0) {
		fmt.Printf("EXPECTED: %s\nACTUAL:   %s\n", xInv.Text(16), xInv.Text(16))
		t.Fail()
	}
}

func TestDouble(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	pp256, _ := P256().(p256CurveFast)
	z, _ := new(big.Int).SetString("1", 10)
	x, y, z := pp256.TestDouble(p256Params.Gx, p256Params.Gy, z)
	
	xExp, _ := new(big.Int).SetString("3ac62a6b166498176eb40d98a27587d8d7c97ac5f11ca1d8d851202b22c3f5c8", 16)
	yExp, _ := new(big.Int).SetString("234aabe92636af27ea8edcd2392f97839c5a74b7ddea27bce94c2d270fb65157", 16)
	zExp, _ := new(big.Int).SetString("57401fa8db8e8e1118a40621ce27d6842bc1e1cef6138faabaf37b85a2a774ea", 16)
	
	if (x.Cmp(xExp)!=0 || y.Cmp(yExp)!=0 || z.Cmp(zExp)!=0) {
		fmt.Printf("EXPECTED: %s\nEXPECTED: %s\nEXPECTED: %s\n", x.Text(16), y.Text(16), z.Text(16),)
		fmt.Printf("ACTUAL:   %s\nACTUAL:   %s\nACTUAL:   %s\n", xExp.Text(16), yExp.Text(16), zExp.Text(16),)
		t.Fail()
	}
}

func TestAdd(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	P256()
	
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	p1  := new(p256Point)
	res  := new(p256Point)
	exp1 := new(p256Point)
	
	p256PointDoubleAsm(p1, p2)
	xExp, _ := new(big.Int).SetString("46b25328fadd64b3d27a6dd305abf1aedb90d91e2b4557e0385baba2d08d581", 16)
	yExp, _ := new(big.Int).SetString("8501da36024740683dc01b7ac378aaf0ada17e37f7d757024d3d8a1222392169", 16)
	zExp, _ := new(big.Int).SetString("a544fc8b9b4e66fb3f7e28e7822754d67c47431d67b9c376b5098d22b457a054", 16)
	copy(exp1.x[:], fromBig(xExp))
	copy(exp1.y[:], fromBig(yExp))
	copy(exp1.z[:], fromBig(zExp))
	
	p256PointAddAsmBig(res, p1, p2)  // res = p1 + p2
	p256PointAddAsm(res, p1, p2)  // res = p1 + p2
	if (ComparePoint(res, exp1)!=0) {
		fmt.Printf("[@4] Expected res == p1 + p2\n")
		PrintPoint("exp", exp1)
		PrintPoint("res", res)
		PrintPoint("in1", p1)
		PrintPoint("in2", p2)	
		t.Fail();
	}
}

func ComparePoint(p1, p2 *p256Point) (int) {
	if (bytes.Compare(p1.x[:], p2.x[:])==0 && 
		bytes.Compare(p1.y[:], p2.y[:])==0 && 
		bytes.Compare(p1.z[:], p2.z[:])==0) {
		return 0
	}
	return 1
}

func PrintPoint(msg string, p1 *p256Point) {
	fmt.Printf("INPUT %s.x:    %s\nINPUT %s.y:    %s\nINPUT %s.z:    %s\n\n", msg, new(big.Int).SetBytes(p1.x[:]).Text(16), 
		                                                                      msg, new(big.Int).SetBytes(p1.y[:]).Text(16), 
		                                                                      msg, new(big.Int).SetBytes(p1.z[:]).Text(16))
}

func TestAddAffine(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	P256()
	
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	p1  := new(p256Point)
	res  := new(p256Point)
	exp1 := new(p256Point)
	exp2 := new(p256Point)
	exp3 := new(p256Point)
	
	p256PointDoubleAsm(p1, p2)
	xExp, _ := new(big.Int).SetString("46b25328fadd64b3d27a6dd305abf1aedb90d91e2b4557e0385baba2d08d581", 16)
	yExp, _ := new(big.Int).SetString("8501da36024740683dc01b7ac378aaf0ada17e37f7d757024d3d8a1222392169", 16)
	zExp, _ := new(big.Int).SetString("a544fc8b9b4e66fb3f7e28e7822754d67c47431d67b9c376b5098d22b457a054", 16)
	copy(exp1.x[:], fromBig(xExp))
	copy(exp1.y[:], fromBig(yExp))
	copy(exp1.z[:], fromBig(zExp))
	yExp, _ = new(big.Int).SetString("7a8e00e6da77a27b2d17797722de0cda74b5471c45e61ba3220daca8316aa9f5", 16)
	copy(exp2.x[:], p2.x[:])
	copy(exp2.y[:], fromBig(yExp))
	copy(exp2.z[:], p2.z[:])
	xExp, _ = new(big.Int).SetString("a16ecd2edf99bba8d6ad839e9d4d7ef17a0e7dd416e9f09fea66f67e26465c37", 16)
	yExp, _ = new(big.Int).SetString("2d1bce3c02c56f7d94d97e00136ea2423688b905a687f94ff8b1364a1df174b", 16)
	zExp, _ = new(big.Int).SetString("a544fc8b9b4e66fb3f7e28e7822754d67c47431d67b9c376b5098d22b457a054", 16)
	copy(exp3.x[:], fromBig(xExp))
	copy(exp3.y[:], fromBig(yExp))
	copy(exp3.z[:], fromBig(zExp))
	
	// If sign == 1 -> P2 = -P2
	// If sel == 0 -> P3 = P1
	// if zero == 0 -> P3 = P2
	//func p256PointAddAffineAsm(P3, P1, P2 *p256Point, sign, sel, zero int)
	p256PointAddAffineAsm(res, p1, p2, 0, 0, 0)  // res = p2
	if (ComparePoint(res, p2)!=0) {
		fmt.Printf("[@1] Expected res == in2)\n")
		PrintPoint("in2", p2)
		PrintPoint("res", res)
		t.Fail();
	}
	
	p256PointAddAffineAsm(res, p1, p2, 0, 0, 1)  // res = p1
	if (ComparePoint(res, p1)!=0) {
		fmt.Printf("[@2] Expected res == in1)\n")
		PrintPoint("in1", p1)
		PrintPoint("res", res)
		t.Fail();
	}
	p256PointAddAffineAsm(res, p1, p2, 0, 1, 0)  // res = p2
	if (ComparePoint(res, p2)!=0) {
		fmt.Printf("[@3] Expected res == in2)\n")
		PrintPoint("in2", p2)
		PrintPoint("res", res)
		t.Fail();
	}
	p256PointAddAffineAsm(res, p1, p2, 0, 1, 1)  // res = p1 + p2
	if (ComparePoint(res, exp1)!=0) {
		fmt.Printf("[@4] Expected res == p1 + p2\n")
		PrintPoint("exp", exp1)
		PrintPoint("res", res)
		PrintPoint("in1", p1)
		PrintPoint("in2", p2)	
		t.Fail();
	}
	
	p256PointAddAffineAsm(res, p1, p2, 1, 0, 0)  // res = -p2
	if (ComparePoint(res, exp2)!=0) {
		fmt.Printf("[@5] Expected res == -in2)\n")
		PrintPoint("in2", exp2)
		PrintPoint("res", res)
		t.Fail();
	}
	
	p256PointAddAffineAsm(res, p1, p2, 1, 0, 1)  // res = p1
	if (ComparePoint(res, p1)!=0) {
		fmt.Printf("[@6] Expected res == in1)\n")
		PrintPoint("in1", p1)
		PrintPoint("res", res)
		t.Fail();
	}
		
	p256PointAddAffineAsm(res, p1, p2, 1, 1, 0)  // res = -p2
	if (ComparePoint(res, exp2)!=0) {
		fmt.Printf("[@7] Expected res == -in2)\n")
		PrintPoint("in2", p2)
		PrintPoint("res", exp2)
		t.Fail();
	}
		
	p256PointAddAffineAsm(res, p1, p2, 1, 1, 1)  // res = p1 + (-p2)
	if (ComparePoint(res, exp3)!=0) {
		fmt.Printf("[@8] Expected res == p1 + (-p2)\n")
		PrintPoint("exp", exp3)
		PrintPoint("res", res)	
		t.Fail();
	}
}

func TestAddAffineFine(t *testing.T) {
	t.SkipNow()
	
	P256()
	
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	p1  := new(p256Point)
	res  := new(p256Point)
	exp1 := new(p256Point)
	exp2 := new(p256Point)
	exp3 := new(p256Point)
	
	p256PointDoubleAsm(p1, p2)
	xExp, _ := new(big.Int).SetString("46b25328fadd64b3d27a6dd305abf1aedb90d91e2b4557e0385baba2d08d581", 16)
	yExp, _ := new(big.Int).SetString("8501da36024740683dc01b7ac378aaf0ada17e37f7d757024d3d8a1222392169", 16)
	zExp, _ := new(big.Int).SetString("a544fc8b9b4e66fb3f7e28e7822754d67c47431d67b9c376b5098d22b457a054", 16)
	copy(exp1.x[:], fromBig(xExp))
	copy(exp1.y[:], fromBig(yExp))
	copy(exp1.z[:], fromBig(zExp))
	yExp, _ = new(big.Int).SetString("7a8e00e6da77a27b2d17797722de0cda74b5471c45e61ba3220daca8316aa9f5", 16)
	copy(exp2.x[:], p2.x[:])
	copy(exp2.y[:], fromBig(yExp))
	copy(exp2.z[:], p2.z[:])
	xExp, _ = new(big.Int).SetString("a16ecd2edf99bba8d6ad839e9d4d7ef17a0e7dd416e9f09fea66f67e26465c37", 16)
	yExp, _ = new(big.Int).SetString("2d1bce3c02c56f7d94d97e00136ea2423688b905a687f94ff8b1364a1df174b", 16)
	zExp, _ = new(big.Int).SetString("a544fc8b9b4e66fb3f7e28e7822754d67c47431d67b9c376b5098d22b457a054", 16)
	copy(exp3.x[:], fromBig(xExp))
	copy(exp3.y[:], fromBig(yExp))
	copy(exp3.z[:], fromBig(zExp))
	
	// If sign == 1 -> P2 = -P2
	// If sel == 0 -> P3 = P1
	// if zero == 0 -> P3 = P2
	//func p256PointAddAffineAsm(P3, P1, P2 *p256Point, sign, sel, zero int)

	p256PointAddAffineAsmBig(res, p1, p2, 0, 1, 1)  // res = p1 + p2
	p256PointAddAffineAsm(res, p1, p2, 0, 1, 1)  // res = p1 + p2
	if (ComparePoint(res, exp1)!=0) {
		fmt.Printf("[@4] Expected res == p1 + p2\n")
		PrintPoint("exp", exp1)
		PrintPoint("res", res)
		PrintPoint("in1", p1)
		PrintPoint("in2", p2)	
		t.Fail();
	}
	
}

func TestStressP256AddAffine(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	P256()
	
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	res  := new(p256Point)
	exp1 := new(p256Point)
	
	p256PointDoubleAsm(res, p2)
	
	for i:=0; i<100000; i++ {
		cond, _ := rand.Int(rand.Reader, big.NewInt(3))
		if (cond.Int64() == 2) {  
			p256PointDoubleAsm(res, res)
		}
		
		p256PointAddAffineAsmBig(exp1, res, p2, 0, 1, 1)  // res = p1 + p2
		p256PointAddAffineAsm(res, res, p2, 0, 1, 1)  // res = p1 + p2
		if (ComparePoint(res, exp1)!=0) {
			fmt.Printf("[@4] Expected res == p1 + p2\n")
			PrintPoint("exp", exp1)
			PrintPoint("res", res)
			//PrintPoint("in1", p1)
			//PrintPoint("in2", p2)	
			t.FailNow();
		}
		
		if ( 0 == i%1000) {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")
}


func TestStressP256Add(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	P256()
	
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	res  := new(p256Point)
	exp1 := new(p256Point)
	
	p256PointDoubleAsm(res, p2)
	
	for i:=0; i<100000; i++ {
		cond, _ := rand.Int(rand.Reader, big.NewInt(3))
		if (cond.Int64() == 2) {  
			p256PointDoubleAsm(res, res)
		}
		
		p256PointAddAsmBig(exp1, res, p2)  // res = p1 + p2
		p256PointAddAsm(res, res, p2)  // res = p1 + p2
		if (ComparePoint(res, exp1)!=0) {
			fmt.Printf("[@4] Expected res == p1 + p2\n")
			PrintPoint("exp", exp1)
			PrintPoint("res", res)
			//PrintPoint("in1", p1)
			//PrintPoint("in2", p2)	
			t.FailNow();
		}
		
		if ( 0 == i%1000) {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")
}

func TestStressP256Double(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	P256()
	
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	res  := new(p256Point)
	exp1 := new(p256Point)
	p1  := new(p256Point)
	
	p256PointDoubleAsm(res, p2)
	
	for i:=0; i<100000; i++ {
		*p1 = *p2
		cond, _ := rand.Int(rand.Reader, big.NewInt(3))
		if (cond.Int64() == 2) {  
			p256PointAddAffineAsm(p1, p1, p2, 0, 1, 1)
		}
		
		p256PointDoubleAsmBig(exp1, p1)
		p256PointDoubleAsm(res, p1)
		if (ComparePoint(res, exp1)!=0) {
			fmt.Printf("[@4] Expected res == p1 + p2\n")
			PrintPoint("exp", exp1)
			PrintPoint("res", res)
			PrintPoint("in1", p1)
			//PrintPoint("in2", p2)	
			t.FailNow();
		}
		
		if ( 0 == i%1000) {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n")
}

func TestDoubleFine1(t *testing.T) {
	//t.SkipNow()
	
	P256()

	p1  := new(p256Point)
	p2  := new(p256Point)
	res  := new(p256Point)
	exp1 := new(p256Point)
	
	xExp, _ := new(big.Int).SetString("5de4f991e25e717ad2aad563a7090a86c6b2e33a31355d9bf26ac1b54f2bf829", 16)
	yExp, _ := new(big.Int).SetString("f9f84f1f77bd1fdeda98128dce880027cf48a95eaa13daa3dab967123ccc950d", 16)
	zExp, _ := new(big.Int).SetString("ae3fe314b10bb0aa5d10d11ba43e64b169571c87433c8b9bbe4a6af9d2aac15", 16)
	copy(exp1.x[:], fromBig(xExp))
	copy(exp1.y[:], fromBig(yExp))
	copy(exp1.z[:], fromBig(zExp))
	
	xExp, _ = new(big.Int).SetString("18905f76a53755c679fb732b7762251075ba95fc5fedb60179e730d418a9143c", 16)
	yExp, _ = new(big.Int).SetString("8571ff1825885d85d2e88688dd21f3258b4ab8e4ba19e45cddf25357ce95560a", 16)
	zExp, _ = new(big.Int).SetString("fffffffeffffffffffffffffffffffff000000000000000000000001", 16)
	copy(p1.x[:], fromBig(xExp))
	copy(p1.y[:], fromBig(yExp))
	copy(p1.z[:], fromBig(zExp))
	
	// If sign == 1 -> P2 = -P2
	// If sel == 0 -> P3 = P1
	// if zero == 0 -> P3 = P2
	//func p256PointAddAffineAsm(P3, P1, P2 *p256Point, sign, sel, zero int)

	p256PointDoubleAsmBig(p2, p1)
	p256PointDoubleAsm(res, p1)
	if (ComparePoint(res, exp1)!=0) {
		fmt.Printf("Expected res == p1 + p1\n")
		PrintPoint("exp", exp1)
		PrintPoint("res", res)
		PrintPoint("in1", p1)
		PrintPoint("in2", p2)	
		t.Fail();
	}
	
}

func BenchmarkP256AddAffine(b *testing.B) {
    P256()
    
    basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	res  := new(p256Point)
	p256PointDoubleAsm(res, p2)
	
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
		p256PointAddAffineAsm(res, res, p2, 0, 1, 1)
    }
}

func BenchmarkP256Add(b *testing.B) {
    P256()
    
    basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	res  := new(p256Point)
	p256PointDoubleAsm(res, p2)
	
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
		p256PointAddAsm(res, res, p2)
    }
}

func BenchmarkP256Double(b *testing.B) {
    P256()
    
    basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	res  := new(p256Point)
	
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
		p256PointDoubleAsm(res, p2)
    }
}

func TestDoubleFine2(t *testing.T) {
	//t.SkipNow()
	
	P256()
	
	basePoint := p256Point{ 
		x: [32]byte{0x18, 0x90, 0x5f, 0x76, 0xa5, 0x37, 0x55, 0xc6, 0x79, 0xfb, 0x73, 0x2b, 0x77, 0x62, 0x25, 0x10, 
			0x75, 0xba, 0x95, 0xfc, 0x5f, 0xed, 0xb6, 0x01, 0x79, 0xe7, 0x30, 0xd4, 0x18, 0xa9, 0x14, 0x3c}, //(p256.x*2^256)%p
		y:[32]byte{0x85, 0x71, 0xff, 0x18, 0x25, 0x88, 0x5d, 0x85, 0xd2, 0xe8, 0x86, 0x88, 0xdd, 0x21, 0xf3, 0x25,
			0x8b, 0x4a, 0xb8, 0xe4, 0xba, 0x19, 0xe4, 0x5c, 0xdd, 0xf2, 0x53, 0x57, 0xce, 0x95, 0x56, 0x0a}, //(p256.y*2^256)%p
		z:[32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},  //(p256.z*2^256)%p
	}
	p2 := &basePoint
	p1  := new(p256Point)
	res  := new(p256Point)
	exp1 := new(p256Point)

	p256PointDoubleAsmBig(exp1, p2)	
	p256PointDoubleAsm(res, p2)
	if (ComparePoint(res, exp1)!=0) {
		fmt.Printf("[@4] Expected res == p1 + p1\n")
		PrintPoint("exp", exp1)
		PrintPoint("res", res)
		PrintPoint("in1", p1)
		t.Fail();
	}
	
}

//func TestLoopBooth(t *testing.T) {
//	// (This is one, in the Montgomery domain.)
//	scalar1 := [32]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
//			        0xff, 0xff, 0xff, 0xff, 0x0c, 0x0b, 0x0a, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}
//	
//	wvalue := (uint64(scalar1[31]) << 1) & 0xff
//	sel, sign := boothW7(uint(wvalue))
//	fmt.Printf("+sel %d, sign %d, wvalue %d\n", sel, sign, wvalue)
//
//	scalar2 := make([]uint64, 4)
//	// (This is one, in the Montgomery domain.)
//	scalar2[0] = 0x0807060504030201
//	scalar2[1] = 0xffffffff0c0b0a09
//	scalar2[2] = 0xffffffffffffffff
//	scalar2[3] = 0x00000000fffffffe
//	
//	
//	wvalue2 := (scalar2[0] << 1) & 0xff
//	sel, sign = boothW7(uint(wvalue2))
//	fmt.Printf("-sel %d, sign %d, wvalue %d\n\n", sel, sign, wvalue2)
//	
//	index := uint(6)
//	for i := 1; i < 37; i++ {
//		if index < 247 {
//			wvalue = ((uint64(scalar1[31-index/8]) >> (index % 8)) + (uint64(scalar1[31-index/8-1]) << (8 - (index % 8)))) & 0xff
//		} else {
//			wvalue = (uint64(scalar1[31-index/8]) >> (index % 8)) & 0xff
//		}
//		sel, sign = boothW7(uint(wvalue))
//		fmt.Printf("+[%d] sel %d, sign %d, wvalue %d scalar1[%d]) >> (%d)) + ((scalar1[%d]) << (%d)\n", index, sel, sign, wvalue, 31-index/8-1, index % 8, 31-index/8,8 - (index % 8))
//		
//		
//		if index < 192 {
//			wvalue2 = ((scalar2[index/64] >> (index % 64)) + (scalar2[index/64+1] << (64 - (index % 64)))) & 0xff
//			sel, sign = boothW7(uint(wvalue2))
//			fmt.Printf("-[%d] sel %d, sign %d, wvalue %d   (scalar2[%d] >> (%d)) + (scalar2[%d] << (%d))\n\n", index, sel, sign, wvalue2, index/64, index % 64, index/64+1, 64 - (index % 64))
//		} else {
//			wvalue2 = (scalar2[index/64] >> (index % 64)) & 0xff
//			sel, sign = boothW7(uint(wvalue2))
//			fmt.Printf("-[%d] sel %d, sign %d, wvalue %d   scalar2[%d] >> (%d)\n\n", index, sel, sign, wvalue2, index/64, index % 64)
//		}
//		index += 7
//	}
//}

//func TestBaseMult2(t *testing.T) {
//	if testing.Short() {
//        t.SkipNow()
//    }
//	pp256, _ := P256().(p256Curve)
//	
//	k, _ := new(big.Int).SetString("a544fc8b9b4e66fb3f7e28e7822754d67c47431d67b9c376b5098d22b457a054", 16)
//	x, y := pp256.ScalarBaseMult(k.Bytes())
//	fmt.Printf("x %s\ny %s\n", x.Text(16), y.Text(16))
//}

//func TestMult2(t *testing.T) {
//	if testing.Short() {
//        t.SkipNow()
//    }
//	pp256, _ := P256().(p256Curve)
//	
//	xExp, _ := new(big.Int).SetString("b70e0cbd6bb4bf7f321390b94a03c1d356c21122343280d6115c1d21", 16)
//	yExp, _ := new(big.Int).SetString("bd376388b5f723fb4c22dfe6cd4375a05a07476444d5819985007e34", 16)
//	
//	//k, _ := new(big.Int).SetString("a544fc8b9b4e66fb3f7e28e7822754d67c47431d67b9c376b5098d22b457a054", 16)
//	k, _ := new(big.Int).SetString("2", 16)
//	x, y := pp256.ScalarMult(xExp, yExp, k.Bytes())
//	fmt.Printf("x %s\ny %s\n", x.Text(16), y.Text(16))
//}

func TestInitTable(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	P256()
	
	exp1 := new(p256Point)
	
	xExp, _ := new(big.Int).SetString("14af860fcd26d2b48e525f1a46a5122924ae1c304ad63f99ab41b43a43228d83", 16)
	yExp, _ := new(big.Int).SetString("82ceb1dd8a37b527d3e21fcee6a9d694f51865adeb78795ed6baef613f714aa1", 16)
	copy(exp1.x[:], fromBig(xExp))
	copy(exp1.y[:], fromBig(yExp))
	
	if (ComparePoint(&p256PreFast[36][63], exp1)!=0) {
		PrintPoint("exp", exp1)
		PrintPoint("res", &p256PreFast[36][63])	
		t.Fail();
	}

//	for j := 0; j < 64; j++ {
//		for i := 0; i < 37; i++ {
//			fmt.Printf("INPUT %d|%d.x:    %s\nINPUT %d|%d.y:    %s\n\n", i, j, new(big.Int).SetBytes(p256Precomputed[i][j].x[:]).Text(16), 
//				                                                         i, j, new(big.Int).SetBytes(p256Precomputed[i][j].y[:]).Text(16))
//		}
//		fmt.Printf("=================================================\n")
//	}
}
