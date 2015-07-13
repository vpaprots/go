// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file encapsulates some of the odd characteristics of the
// s390x instruction set, to minimize its interaction
// with the core of the assembler.

package arch

import "cmd/internal/obj/s390x"

func jumpS390x(word string) bool {
	switch word {
	case "BC", "BCL", "BEQ", "BGE", "BGT", "BL", "BLE", "BLT", "BNE", "BR", "BVC", "BVS", "CALL", "JMP":
		return true
	}
	return false
}

// Iss390xRLD reports whether the op (as defined by an s390x.A* constant) is
// one of the RLD-like instructions that require special handling.
// The FMADD-like instructions behave similarly.
func IsS390xRLD(op int) bool {
	switch op {
	case s390x.ARLDC, s390x.ARLDCCC, s390x.ARLDCL, s390x.ARLDCLCC,
		s390x.ARLDCR, s390x.ARLDCRCC, s390x.ARLDMI, s390x.ARLDMICC,
		s390x.ARLWMI, s390x.ARLWMICC, s390x.ARLWNM, s390x.ARLWNMCC:
		return true
	case s390x.AFMADD, s390x.AFMADDCC, s390x.AFMADDS, s390x.AFMADDSCC,
		s390x.AFMSUB, s390x.AFMSUBCC, s390x.AFMSUBS, s390x.AFMSUBSCC,
		s390x.AFNMADD, s390x.AFNMADDCC, s390x.AFNMADDS, s390x.AFNMADDSCC,
		s390x.AFNMSUB, s390x.AFNMSUBCC, s390x.AFNMSUBS, s390x.AFNMSUBSCC:
		return true
	}
	return false
}

// Iss390xCMP reports whether the op (as defined by an s390x.A* constant) is
// one of the CMP instructions that require special handling.
func IsS390xCMP(op int) bool {
	switch op {
	case s390x.ACMP, s390x.ACMPU, s390x.ACMPW, s390x.ACMPWU:
		return true
	}
	return false
}

// Iss390xNEG reports whether the op (as defined by an s390x.A* constant) is
// one of the NEG-like instructions that require special handling.
func IsS390xNEG(op int) bool {
	switch op {
	case s390x.AADDMECC, s390x.AADDMEVCC, s390x.AADDMEV, s390x.AADDME,
		s390x.AADDZECC, s390x.AADDZEVCC, s390x.AADDZEV, s390x.AADDZE,
		s390x.ACNTLZDCC, s390x.ACNTLZD, s390x.ACNTLZWCC, s390x.ACNTLZW,
		s390x.AEXTSBCC, s390x.AEXTSB, s390x.AEXTSHCC, s390x.AEXTSH,
		s390x.AEXTSWCC, s390x.AEXTSW, s390x.ANEGCC, s390x.ANEGVCC,
		s390x.ANEGV, s390x.ANEG,
		s390x.ASUBMECC, s390x.ASUBMEVCC, s390x.ASUBMEV,
		s390x.ASUBME, s390x.ASUBZECC, s390x.ASUBZEVCC, s390x.ASUBZEV,
		s390x.ASUBZE:
		return true
	}
	return false
}

// IsS390xStorageAndStorage reports whether the op (as defined by an s390x.A* constant) refers
// to an storage-and-storage format instruction such as mvc, clc, xc, oc or nc.
func IsS390xStorageAndStorage(op int) bool {
	switch op {
	case s390x.AMVC, s390x.ACLC, s390x.AXC, s390x.AOC, s390x.ANC:
		return true
	}
	return false
}

func s390xRegisterNumber(name string, n int16) (int16, bool) {
	switch name {
	case "AR":
		if 0 <= n && n <= 15 {
			return s390x.REG_AR0 + n, true
		}
	case "F":
		if 0 <= n && n <= 15 {
			return s390x.REG_F0 + n, true
		}
	case "R":
		if 0 <= n && n <= 15 {
			return s390x.REG_R0 + n, true
		}
	}
	return 0, false
}
