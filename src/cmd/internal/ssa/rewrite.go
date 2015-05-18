// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

import "log"

func applyRewrite(f *Func, r func(*Value) bool) {
	// repeat rewrites until we find no more rewrites
	var curv *Value
	defer func() {
		if curv != nil {
			log.Printf("panic during rewrite of %s\n", curv.LongString())
			// TODO(khr): print source location also
		}
	}()
	for {
		change := false
		for _, b := range f.Blocks {
			for _, v := range b.Values {
				// elide any copies generated during rewriting
				for i, a := range v.Args {
					if a.Op != OpCopy {
						continue
					}
					for a.Op == OpCopy {
						a = a.Args[0]
					}
					v.Args[i] = a
				}

				// apply rewrite function
				curv = v
				if r(v) {
					change = true
				}
			}
		}
		if !change {
			curv = nil
			return
		}
	}
}

// Common functions called from rewriting rules

func is64BitInt(t Type) bool {
	return t.Size() == 8 && t.IsInteger()
}

func is32BitInt(t Type) bool {
	return t.Size() == 4 && t.IsInteger()
}

func isPtr(t Type) bool {
	return t.IsPtr()
}

func isSigned(t Type) bool {
	return t.IsSigned()
}

func typeSize(t Type) int64 {
	return t.Size()
}

// addOff adds two offset aux values.  Each should be an int64.  Fails if wraparound happens.
func addOff(a, b interface{}) interface{} {
	x := a.(int64)
	y := b.(int64)
	z := x + y
	// x and y have same sign and z has a different sign => overflow
	if x^y >= 0 && x^z < 0 {
		log.Panicf("offset overflow %d %d\n", x, y)
	}
	return z
}

func inBounds(idx, len int64) bool {
	return idx >= 0 && idx < len
}