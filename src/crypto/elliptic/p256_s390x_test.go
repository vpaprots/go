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

func TestP256Mul(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	pp256, _ := P256().(p256Curve)
	pp256.TestP256Mul();
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


func TestInverse(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	x, _ := new(big.Int).SetString("15792089210356248762697446949407573530086143415290314195533631308867097853951", 10)
	xInvExp, _ := new(big.Int).SetString("78239946472340125005789637834181368510016866213607708536507216111758801279690", 10)
	
	pp256, _ := P256().(p256Curve)
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
	
	pp256, _ := P256().(p256Curve)
	z, _ := new(big.Int).SetString("1", 10)
	x, y, z := pp256.TestDouble(pp256.Gx, pp256.Gy, z)
	
	xExp, _ := new(big.Int).SetString("3ac62a6b166498176eb40d98a27587d8d7c97ac5f11ca1d8d851202b22c3f5c8", 16)
	yExp, _ := new(big.Int).SetString("234aabe92636af27ea8edcd2392f97839c5a74b7ddea27bce94c2d270fb65157", 16)
	zExp, _ := new(big.Int).SetString("57401fa8db8e8e1118a40621ce27d6842bc1e1cef6138faabaf37b85a2a774ea", 16)
	
	if (x.Cmp(xExp)!=0 || z.Cmp(xExp)!=0 || z.Cmp(zExp)!=0) {
		fmt.Printf("EXPECTED: %s\nEXPECTED: %s\nEXPECTED: %s\n", x.Text(16), y.Text(16), z.Text(16),)
		fmt.Printf("ACTUAL:   %s\nACTUAL:   %s\nACTUAL:   %s\n", xExp.Text(16), yExp.Text(16), zExp.Text(16),)
		t.Fail()
	}
}

func TestAdd(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	pp256, _ := P256().(p256Curve)
	z2, _ := new(big.Int).SetString("1", 10) 
	x, y, z := pp256.TestDouble(pp256.Gx, pp256.Gy, z2)
	x, y, z = pp256.TestAdd(pp256.Gx, pp256.Gy, z2, x, y, z)
	
	xExp, _ := new(big.Int).SetString("f59160063c80b047a18194491dc75dc6085cfc92ba9cb84da8d61c3d420ebb84", 16)
	yExp, _ := new(big.Int).SetString("2c5f563ddc8a96fa250ab7d97a624d84206972dadd2c1548c29213deba3b5e11", 16)
	zExp, _ := new(big.Int).SetString("20d428c1b39721225bda48531f6603eb0c14dfa13af43f4b8a3415d6d13a8cc0", 16)
	
		if (x.Cmp(xExp)!=0 || z.Cmp(xExp)!=0 || z.Cmp(zExp)!=0) {
		fmt.Printf("EXPECTED: %s\nEXPECTED: %s\nEXPECTED: %s\n", x.Text(16), y.Text(16), z.Text(16),)
		fmt.Printf("ACTUAL:   %s\nACTUAL:   %s\nACTUAL:   %s\n", xExp.Text(16), yExp.Text(16), zExp.Text(16),)
		t.Fail()
	}
}