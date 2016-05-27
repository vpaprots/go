// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package elliptic

import (
	"testing"
	"fmt"
	"math/big"
	"crypto/rand"
)

func TestInverse(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	x, _ := new(big.Int).SetString("15792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	pp256, _ := P256().(p256Curve)
	xInv := pp256.Inverse(x)
	fmt.Printf("EXPECTED: %s\n", xInv.String())
}

func TestStressInverse(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	pp256, _ := P256().(p256Curve)
	for i:=0; i<100000; i++ {
		x, _ := rand.Int(rand.Reader, pp256.N)
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
	fmt.Printf("\nSUCCESS!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
}

func TestSpecificInverse(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	pp256, _ := P256().(p256Curve)
	pp256.TestOrdMul();
}

func BenchmarkInverse(b *testing.B) {
    pp256, _ := P256().(p256Curve)
    x, _ := rand.Int(rand.Reader, pp256.N)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pp256.Inverse(x)
    }
}


func TestMul(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	x, _ := new(big.Int).SetString("15792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	pp256, _ := P256().(p256Curve)
	xInv := pp256.TestMul(x)
	fmt.Printf("EXPECTED: %s\n", xInv.String())
}

func TestDouble(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	pp256, _ := P256().(p256Curve)
	z, _ := new(big.Int).SetString("1", 10)
	x, y, z := pp256.TestDouble(pp256.Gx, pp256.Gy, z)
	fmt.Printf("EXPECTED: %s\nEXPECTED: %s\nEXPECTED: %s\n", x.Text(16), y.Text(16), z.Text(16),)
}
