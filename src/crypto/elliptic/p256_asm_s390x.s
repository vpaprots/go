// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


#include "textflag.h"

#define res_ptr R1
#define x_ptr R2
#define y_ptr R3
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define M0    V4
#define M1    V5
#define T0    V6
#define T1    V7
#define T2    V8
#define YDIG  V9

#define ADD1  V16
#define ADD1H V17
#define ADD2  V18
#define ADD2H V19
#define RED1  V20
#define RED1H V21
#define RED2  V22
#define RED2H V23
#define CAR1  V24
#define CAR1M V25

#define MK0   V30
#define K0    V31

DATA p256ordK0<>+0x00(SB)/4, $0xee00bc4f
DATA p256ord<>+0x00(SB)/8, $0xffffffff00000000
DATA p256ord<>+0x08(SB)/8, $0xffffffffffffffff
DATA p256ord<>+0x10(SB)/8, $0xbce6faada7179e84
DATA p256ord<>+0x18(SB)/8, $0xf3b9cac2fc632551
DATA p256<>+0x00(SB)/8, $0xffffffff00000001
DATA p256<>+0x08(SB)/8, $0x0000000000000000
DATA p256<>+0x10(SB)/8, $0x00000000ffffffff
DATA p256<>+0x18(SB)/8, $0xffffffffffffffff
DATA p256<>+0x20(SB)/8, $0x0c0d0e0f1c1d1e1f
DATA p256<>+0x28(SB)/8, $0x0c0d0e0f1c1d1e1f
DATA p256<>+0x30(SB)/8, $0x0000000010111213
DATA p256<>+0x38(SB)/8, $0x1415161700000000
GLOBL p256const0<>(SB), 8, $8
GLOBL p256const1<>(SB), 8, $8
GLOBL p256ordK0<>(SB), 8, $4
GLOBL p256ord<>(SB), 8, $32
GLOBL p256<>(SB), 8, $64

/* ---------------------------------------*/
// func p256OrdMul(res, in1, in2 []byte)
TEXT ·p256OrdMul(SB),NOSPLIT,$0
	MOVD res+0(FP), res_ptr
	MOVD in1+24(FP), x_ptr
	MOVD in2+48(FP), y_ptr

	VZERO T2
	MOVD	$p256ordK0<>+0x00(SB), R4
	//VLEF    $3, 0(R4), K0
	WORD	$0xE7F40000
	BYTE    $0x38
	BYTE    $0x03 //
	MOVD	$p256ord<>+0x00(SB), R4
	VL	16(R4), M0
	VL	0(R4),  M1

	VL (1*16)(x_ptr), X0
	VL (0*16)(x_ptr), X1
	VL (1*16)(y_ptr), Y0
	VL (0*16)(y_ptr), Y1

	//---------------------------------------------------------------------------/
	VREPF $3,  Y0,   YDIG
	VMLF  X0,  YDIG, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMLF  X1,YDIG, ADD2
	VMLHF X0,YDIG, ADD1H
	VMLHF X1,YDIG, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
    /**
   	 * ---+--------+--------+
	 *  T2|   T1   |   T0   |
	 * ---+--------+--------+
	 *           *(add)*
	 *    +--------+--------+
	 *    |   X1   |   X0   |
	 *    +--------+--------+
	 *           *(mul)*
 	 *    +--------+--------+
	 *    |  YDIG  |  YDIG  |
	 *    +--------+--------+
 	 *           *(add)*
	 *    +--------+--------+
	 *    |   M1   |   M0   |
	 *    +--------+--------+
	 *           *(mul)*
 	 *    +--------+--------+
	 *    |   MK0  |   MK0  |
	 *    +--------+--------+
	 *
	 *   ---------------------
   	 *
   	 *    +--------+--------+
	 *    |  ADD2  |  ADD1  |
	 *    +--------+--------+
   	 *  +--------+--------+
	 *  | ADD2H  | ADD1H  |
	 *  +--------+--------+
  	 *    +--------+--------+
	 *    |  RED2  |  RED1  |
	 *    +--------+--------+
   	 *  +--------+--------+
	 *  | RED2H  | RED1H  |
	 *  +--------+--------+
	 */
	VREPF $2,  Y0,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $1,  Y0,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $0,  Y0,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $3,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $2,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $1,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------
	VREPF $0,  Y1,   YDIG
	VMALF X0,  YDIG, T0, ADD1
	VMLF  ADD1,K0,   MK0
	VREPF $3,  MK0,  MK0

	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VMALF  M0,MK0, ADD1, RED1
	VMALHF M0,MK0, ADD1, RED1H
	VMALF  M1,MK0, ADD2, RED2
	VMALHF M1,MK0, ADD2, RED2H

	VSLDB  $12, RED2,RED1, RED1
	VSLDB  $12, T2,  RED2, RED2

	VACCQ  RED1, ADD1H,CAR1
	VAQ    RED1, ADD1H,T0
	VACCQ  RED1H,T0,   CAR1M
	VAQ    RED1H,T0,   T0

	//<< ready for next MK0

	VACQ    RED2, ADD2H,CAR1,  T1
	VACCCQ  RED2, ADD2H,CAR1,  CAR1
	VACCCQ  RED2H,T1, CAR1M, T2
	VACQ    RED2H,T1, CAR1M, T1
	VAQ     CAR1, T2,T2

	//---------------------------------------------------

	VZERO   RED1
	VSCBIQ  M0,T0, CAR1
	VSQ     M0,T0, ADD1
	VSBCBIQ T1,M1,CAR1, CAR1M
	VSBIQ   T1,M1,CAR1, ADD2
	VSBIQ   T2,RED1,CAR1M, T2

	//what output to use, ADD2||ADD1 or T1||T0?
	VSEL    T0,ADD1,T2, T0
	VSEL    T1,ADD2,T2, T1

	VST T0, (1*16)(res_ptr)
	VST T1, (0*16)(res_ptr)
	RET

#undef res_ptr
#undef x_ptr
#undef y_ptr
#undef X0
#undef X1
#undef Y0
#undef Y1
#undef M0
#undef M1
#undef T0
#undef T1
#undef T2
#undef YDIG

#undef ADD1
#undef ADD1H
#undef ADD2
#undef ADD2H
#undef RED1
#undef RED1H
#undef RED2
#undef RED2H
#undef CAR1
#undef CAR1M

#undef MK0
#undef K0

/* ---------------------------------------*/
// func p256Mul(res, in1, in2 []byte)
#define res_ptr R1
#define x_ptr R2
#define y_ptr R3
#define X0    V0
#define X1    V1
#define Y0    V2
#define Y1    V3
#define M0    V4
#define M1    V5
#define T0    V6
#define T1    V7
#define T2    V8
#define YDIG  V9

#define ADD1  V16
#define ADD11 V17
#define ADD1H V18
#define ADD2  V19
#define ADD2H V20
#define RED1  V21
#define RED2  V22
#define CAR1  V24
#define CAR2  V25

#define SEL1  V29
#define SEL2  V30
#define ZER   V31
  /**
    *        ---+--------+--------+
    *         T2|   T1   |   T0   |
    *        ---+--------+--------+
    *           +--------+------+-+
    *           |  ADD2  |  ADD1+ |
    *           +--------+------+-+
    *         +--------+--------+
    *         | ADD2H  | ADD1H  |
    *         +--------+--------+
    *      ========================
    *      ---+--------+--------+
    *       T2|   T1   |   T0   |
    *      ---+--------+--------+
    *         +--------+------+-+
    *         |  ADD2  |  ADD1+ |
    *         +--------+------+-+
    *       +--------+--------+
    *       | ADD2H  | ADD1H  |
    *       +--------+--------+
    *     =======================
    *    ---+--------+--------+
    *     T2|   T1   |   T0   |
    *    ---+--------+--------+
    *       +--------+--------+
    *       |  RED2  |  RED1  |
    *       +--------+--------+
    *       +--------+
    *       |  RED1  |
    *       +--------+
    */
TEXT ·p256Mul(SB),NOSPLIT,$0
	MOVD res+0(FP), res_ptr
	MOVD in1+24(FP), x_ptr
	MOVD in2+48(FP), y_ptr

	VZERO T2
	VZERO ZER
	MOVD $p256<>+0x00(SB), R4
	VL	16(R4), M0
	VL	0(R4),  M1
	VL	48(R4), SEL2
	VL	32(R4), SEL1

	VL (1*16)(x_ptr), X0
	VL (0*16)(x_ptr), X1
	VL (1*16)(y_ptr), Y0
	VL (0*16)(y_ptr), Y1

	//---------------------------------------------------------------------------/

	VREPF $3,  Y0, YDIG

	VMLF  X0,YDIG, ADD1
	VMLF  X1,YDIG, ADD2
	VMLHF X0,YDIG, ADD1H
	VMLHF X1,YDIG, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VREPF $2,  Y0,   YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED2, RED1, RED2          // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VREPF $1,  Y0, YDIG

	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VREPF $0,  Y0,   YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED2, RED1, RED2         // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VREPF $2,  Y1, YDIG

	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VREPF $2,  Y1,   YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED2, RED1, RED2          // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VREPF $1,  Y1, YDIG

	VMALF  X0,YDIG, T0, ADD1
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD1, T0
	VSLDB  $12, T2,  ADD2, T1

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	//---------------------------------------------------

	VREPF $0,  Y1,   YDIG

	VMALF  X0,YDIG, T0, ADD11
	VMALF  X1,YDIG, T1, ADD2
	VMALHF X0,YDIG, T0, ADD1H
	VMALHF X1,YDIG, T1, ADD2H

	VSLDB  $12, ADD2,ADD11, T0
	VSLDB  $12, T2,  ADD2,  T1

	VPERM ADD11,ADD1, SEL1, RED2   // d1 d0 d1 d0
	VPERM ZER,  RED2, SEL2, RED1   // 0  d1 d0  0
	VSQ   RED2, RED1, RED2         // Guaranteed not to underflow

	VACCQ  T0, ADD1H,CAR1
	VAQ    T0, ADD1H,T0
	VACCCQ T1, ADD2H,CAR1, T2
	VACQ   T1, ADD2H,CAR1, T1

	VACCQ  T0, RED1,CAR1
	VAQ    T0, RED1,T0
	VACCCQ T1, RED2,CAR1, CAR2
	VACQ   T1, RED2,CAR1, T1
	VAQ    T2, CAR2,T2

	//---------------------------------------------------

	VSCBIQ  M0,T0, CAR1
	VSQ     M0,T0, ADD1
	VSBCBIQ T1,M1, CAR1, CAR2
	VSBIQ   T1,M1, CAR1, ADD2
	VSBIQ   T2,ZER,CAR2, T2

	//what output to use, ADD2||ADD1 or T1||T0?
	VSEL    T0,ADD1,T2, T0
	VSEL    T1,ADD2,T2, T1

	VST T0, (1*16)(res_ptr)
	VST T1, (0*16)(res_ptr)
	RET

