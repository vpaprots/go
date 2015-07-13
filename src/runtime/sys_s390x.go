// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build s390x

package runtime

import "unsafe"

// adjust Gobuf as if it executed a call to fn with context ctxt
// and then did an immediate Gosave.
func gostartcall(buf *gobuf, fn, ctxt unsafe.Pointer) {
	if buf.lr != 0 {
		throw("invalid use of gostartcall")
	}
	buf.lr = buf.pc
	buf.pc = uintptr(fn)
	buf.ctxt = ctxt
}

// Called to rewind context saved during morestack back to beginning of function.
// To help us, the linker emits a jmp back to the beginning right after the
// call to morestack. We just have to decode and apply that jump.
func rewindmorestack(buf *gobuf) {
	var inst uint64
	if buf.pc&1 == 0 && buf.pc != 0 {
		inst = *(*uint64)(unsafe.Pointer(buf.pc))
		inst = inst >> 16 //get rid of junk bytes
		//print("runtime: rewind pc=", hex(buf.pc), " to pc=", hex(inst), "\n");
		if inst>>32 == 0xc0f4 {
			inst = inst & 0xFFFFFFFF
			//print("runtime: rewind inst1 = ",hex(inst),"\n")
			inst = (inst << 1) & 0xFFFFFFFF
			//print("runtime: rewind inst2 = ",hex(inst),"\n")
			buf.pc += uintptr(int32(inst))
			//print("runtime: rewind pc = ",hex(buf.pc),"\n")
			return
		}
	}
	print("runtime: pc=", hex(buf.pc), " ", hex(inst), "\n")
	throw("runtime: misuse of rewindmorestack")
}
