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

func TestP256Mul(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	P256()
	
	exp := make([]byte, 32)
	res := make([]byte, 32)
	x1, _ := new(big.Int).SetString("a007c8559316f82de3d5d9f28b8ffcdf5949bd551f7a1348b8acc00860e058", 16)
	x2, _ := new(big.Int).SetString("66e12d94f3d956202845b2392b6bec594699799c49bd6fa683244c95be79eea2", 16)
	Rinv, _ := new(big.Int).SetString("fffffffe00000003fffffffd0000000200000001fffffffe0000000300000000", 16)
	
	copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256.P)))
	p256Mul(res, fromBig(x1), fromBig(x2))
	
	if (bytes.Compare(exp,res)!=0) {
		fmt.Printf("-EXPECTED %s\n", new(big.Int).SetBytes(exp).Text(16))
		fmt.Printf("-FOUND    %s\n", new(big.Int).SetBytes(res).Text(16))
		fmt.Printf("-TEST in1 %s\n", new(big.Int).SetBytes(fromBig(x1)).Text(16))
		fmt.Printf("-TEST in2 %s\n", new(big.Int).SetBytes(fromBig(x2)).Text(16))
		t.Fail()
	}
}

func TestStressP256Mul(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	
	P256()
	for i:=0; i<1000000; i++ {
		
		exp := make([]byte, 32)
		res := make([]byte, 32)
		x1, _ := rand.Int(rand.Reader, p256.P)
		x2, _ := rand.Int(rand.Reader, p256.P)
		Rinv, _ := new(big.Int).SetString("fffffffe00000003fffffffd0000000200000001fffffffe0000000300000000", 16)
		
		copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256.P)))
		p256Mul(res, fromBig(x1), fromBig(x2))
		
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
	
	copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256.N)))
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
		x1, _ := rand.Int(rand.Reader, p256.N)
		x2, _ := rand.Int(rand.Reader, p256.N)
		Rinv, _ := new(big.Int).SetString("60d066334905c1e907f8b6041e607725badef3e243566fafce1bc8f79c197c79", 16)
		
		copy(exp, fromBig(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(x1, x2), Rinv), p256.N)))
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
	pp256, _ := P256().(p256Curve)
	for i:=0; i<10000; i++ {
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

func BenchmarkInverse(b *testing.B) {
    pp256, _ := P256().(p256Curve)
    x, _ := rand.Int(rand.Reader, pp256.N)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pp256.Inverse(x)
    }
}

func BenchmarkP256Mul(b *testing.B) {
    P256()
    //x, _ := rand.Int(rand.Reader, pp256.N)
    in := make([]byte, 32)
    //copy(in, x.Bytes())
    in[0] = 20
    in[2] = 42
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
		p256Mul(in, in, in)
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
	
	pp256, _ := P256().(p256Curve)
	z2, _ := new(big.Int).SetString("1", 10) 
	x, y, z := pp256.TestDouble(pp256.Gx, pp256.Gy, z2)
	x, y, z = pp256.TestAdd(pp256.Gx, pp256.Gy, z2, x, y, z)
	
	xExp, _ := new(big.Int).SetString("f59160063c80b047a18194491dc75dc6085cfc92ba9cb84da8d61c3d420ebb84", 16)
	yExp, _ := new(big.Int).SetString("2c5f563ddc8a96fa250ab7d97a624d84206972dadd2c1548c29213deba3b5e11", 16)
	zExp, _ := new(big.Int).SetString("20d428c1b39721225bda48531f6603eb0c14dfa13af43f4b8a3415d6d13a8cc0", 16)
	
	if (x.Cmp(xExp)!=0 || y.Cmp(yExp)!=0 || z.Cmp(zExp)!=0) {
		fmt.Printf("EXPECTED: %s\nEXPECTED: %s\nEXPECTED: %s\n", x.Text(16), y.Text(16), z.Text(16),)
		fmt.Printf("ACTUAL:   %s\nACTUAL:   %s\nACTUAL:   %s\n", xExp.Text(16), yExp.Text(16), zExp.Text(16),)
		t.Fail()
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


func TestInitTable(t *testing.T) {
	if testing.Short() {
        t.SkipNow()
    }
	P256()
	precomputeOnce.Do(initTable)
	
	exp1 := new(p256Point)
	
	xExp, _ := new(big.Int).SetString("14af860fcd26d2b48e525f1a46a5122924ae1c304ad63f99ab41b43a43228d83", 16)
	yExp, _ := new(big.Int).SetString("82ceb1dd8a37b527d3e21fcee6a9d694f51865adeb78795ed6baef613f714aa1", 16)
	copy(exp1.x[:], fromBig(xExp))
	copy(exp1.y[:], fromBig(yExp))
	
	if (ComparePoint(&p256Precomputed[36][63], exp1)!=0) {
		PrintPoint("exp", exp1)
		PrintPoint("res", &p256Precomputed[36][63])	
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
