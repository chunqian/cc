// Copyright 2019 The CC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc // import "modernc.org/cc/v3"

import (
	"fmt"
	"math"
	"math/big"
)

var (
	_ Value = (*Float128Value)(nil)
	_ Value = (*InitializerValue)(nil)
	_ Value = Complex128Value(0)
	_ Value = Complex256Value{}
	_ Value = Complex64Value(0)
	_ Value = Float32Value(0)
	_ Value = Float64Value(0)
	_ Value = Int64Value(0)
	_ Value = StringValue(0)
	_ Value = Uint64Value(0)
	_ Value = WideStringValue(0)

	_ Operand = (*funcDesignator)(nil)
	_ Operand = (*lvalue)(nil)
	_ Operand = (*operand)(nil)
	_ Operand = noOperand

	noOperand = &operand{typ: noType}
)

type Operand interface {
	Declarator() *Declarator
	IsLValue() bool
	IsNonZero() bool
	IsZero() bool
	Offset() uintptr // Valid only for non nil Declarator() value
	Type() Type
	Value() Value
	convertFromInt(*context, Node, Type) Operand
	convertTo(*context, Node, Type) Operand
	convertToInt(*context, Node, Type) Operand
	integerPromotion(*context, Node) Operand
	isConst() bool
	normalize(*context, Node) Operand
}

type Value interface {
	add(b Value) Value
	and(b Value) Value
	cpl() Value
	div(b Value) Value
	eq(b Value) Value
	ge(b Value) Value
	gt(b Value) Value
	isNonZero() bool
	isZero() bool
	le(b Value) Value
	lsh(b Value) Value
	lt(b Value) Value
	mod(b Value) Value
	mul(b Value) Value
	neg() Value
	neq(b Value) Value
	or(b Value) Value
	rsh(b Value) Value
	sub(b Value) Value
	xor(b Value) Value
}

type WideStringValue StringID

func (v WideStringValue) add(b Value) Value { panic(internalError()) }
func (v WideStringValue) and(b Value) Value { panic(internalError()) }
func (v WideStringValue) cpl() Value        { panic(internalError()) }
func (v WideStringValue) div(b Value) Value { panic(internalError()) }
func (v WideStringValue) eq(b Value) Value  { return boolValue(v == b.(WideStringValue)) }
func (v WideStringValue) isNonZero() bool   { return true }
func (v WideStringValue) isZero() bool      { return false }
func (v WideStringValue) lsh(b Value) Value { panic(internalError()) }
func (v WideStringValue) mod(b Value) Value { panic(internalError()) }
func (v WideStringValue) mul(b Value) Value { panic(internalError()) }
func (v WideStringValue) neg() Value        { panic(internalError()) }
func (v WideStringValue) neq(b Value) Value { return boolValue(v != b.(WideStringValue)) }
func (v WideStringValue) or(b Value) Value  { panic(internalError()) }
func (v WideStringValue) rsh(b Value) Value { panic(internalError()) }
func (v WideStringValue) sub(b Value) Value { panic(internalError()) }
func (v WideStringValue) xor(b Value) Value { panic(internalError()) }

func (v WideStringValue) le(b Value) Value {
	return boolValue(StringID(v).String() <= StringID(b.(WideStringValue)).String())
}

func (v WideStringValue) ge(b Value) Value {
	return boolValue(StringID(v).String() >= StringID(b.(WideStringValue)).String())
}

func (v WideStringValue) gt(b Value) Value {
	return boolValue(StringID(v).String() > StringID(b.(WideStringValue)).String())
}

func (v WideStringValue) lt(b Value) Value {
	return boolValue(StringID(v).String() < StringID(b.(WideStringValue)).String())
}

type StringValue StringID

func (v StringValue) add(b Value) Value { panic(internalError()) }
func (v StringValue) and(b Value) Value { panic(internalError()) }
func (v StringValue) cpl() Value        { panic(internalError()) }
func (v StringValue) div(b Value) Value { panic(internalError()) }
func (v StringValue) eq(b Value) Value  { return boolValue(v == b.(StringValue)) }
func (v StringValue) isNonZero() bool   { return true }
func (v StringValue) isZero() bool      { return false }
func (v StringValue) lsh(b Value) Value { panic(internalError()) }
func (v StringValue) mod(b Value) Value { panic(internalError()) }
func (v StringValue) mul(b Value) Value { panic(internalError()) }
func (v StringValue) neg() Value        { panic(internalError()) }
func (v StringValue) neq(b Value) Value { return boolValue(v != b.(StringValue)) }
func (v StringValue) or(b Value) Value  { panic(internalError()) }
func (v StringValue) rsh(b Value) Value { panic(internalError()) }
func (v StringValue) sub(b Value) Value { panic(internalError()) }
func (v StringValue) xor(b Value) Value { panic(internalError()) }

func (v StringValue) le(b Value) Value {
	return boolValue(StringID(v).String() <= StringID(b.(StringValue)).String())
}

func (v StringValue) ge(b Value) Value {
	return boolValue(StringID(v).String() >= StringID(b.(StringValue)).String())
}

func (v StringValue) gt(b Value) Value {
	return boolValue(StringID(v).String() > StringID(b.(StringValue)).String())
}

func (v StringValue) lt(b Value) Value {
	return boolValue(StringID(v).String() < StringID(b.(StringValue)).String())
}

type Int64Value int64

func (v Int64Value) add(b Value) Value { return v + b.(Int64Value) }
func (v Int64Value) and(b Value) Value { return v & b.(Int64Value) }
func (v Int64Value) cpl() Value        { return ^v }
func (v Int64Value) eq(b Value) Value  { return boolValue(v == b.(Int64Value)) }
func (v Int64Value) ge(b Value) Value  { return boolValue(v >= b.(Int64Value)) }
func (v Int64Value) gt(b Value) Value  { return boolValue(v > b.(Int64Value)) }
func (v Int64Value) isNonZero() bool   { return v != 0 }
func (v Int64Value) isZero() bool      { return v == 0 }
func (v Int64Value) le(b Value) Value  { return boolValue(v <= b.(Int64Value)) }
func (v Int64Value) lt(b Value) Value  { return boolValue(v < b.(Int64Value)) }
func (v Int64Value) mul(b Value) Value { return v * b.(Int64Value) }
func (v Int64Value) neg() Value        { return -v }
func (v Int64Value) neq(b Value) Value { return boolValue(v != b.(Int64Value)) }
func (v Int64Value) or(b Value) Value  { return v | b.(Int64Value) }
func (v Int64Value) sub(b Value) Value { return v - b.(Int64Value) }
func (v Int64Value) xor(b Value) Value { return v ^ b.(Int64Value) }

func (v Int64Value) div(b Value) Value {
	if b.isZero() {
		return nil
	}

	return v / b.(Int64Value)
}

func (v Int64Value) lsh(b Value) Value {
	switch y := b.(type) {
	case Int64Value:
		return v << uint64(y)
	case Uint64Value:
		return v << y
	default:
		panic(internalError())
	}
}

func (v Int64Value) rsh(b Value) Value {
	switch y := b.(type) {
	case Int64Value:
		return v >> uint64(y)
	case Uint64Value:
		return v >> y
	default:
		panic(internalError())
	}
}

func (v Int64Value) mod(b Value) Value {
	if b.isZero() {
		return nil
	}

	return v % b.(Int64Value)
}

type Uint64Value uint64

func (v Uint64Value) add(b Value) Value { return v + b.(Uint64Value) }
func (v Uint64Value) and(b Value) Value { return v & b.(Uint64Value) }
func (v Uint64Value) cpl() Value        { return ^v }
func (v Uint64Value) eq(b Value) Value  { return boolValue(v == b.(Uint64Value)) }
func (v Uint64Value) ge(b Value) Value  { return boolValue(v >= b.(Uint64Value)) }
func (v Uint64Value) gt(b Value) Value  { return boolValue(v > b.(Uint64Value)) }
func (v Uint64Value) isNonZero() bool   { return v != 0 }
func (v Uint64Value) isZero() bool      { return v == 0 }
func (v Uint64Value) le(b Value) Value  { return boolValue(v <= b.(Uint64Value)) }
func (v Uint64Value) lt(b Value) Value  { return boolValue(v < b.(Uint64Value)) }
func (v Uint64Value) mul(b Value) Value { return v * b.(Uint64Value) }
func (v Uint64Value) neg() Value        { return -v }
func (v Uint64Value) neq(b Value) Value { return boolValue(v != b.(Uint64Value)) }
func (v Uint64Value) or(b Value) Value  { return v | b.(Uint64Value) }
func (v Uint64Value) sub(b Value) Value { return v - b.(Uint64Value) }
func (v Uint64Value) xor(b Value) Value { return v ^ b.(Uint64Value) }

func (v Uint64Value) div(b Value) Value {
	if b.isZero() {
		return nil
	}

	return v / b.(Uint64Value)
}

func (v Uint64Value) lsh(b Value) Value {
	switch y := b.(type) {
	case Int64Value:
		return v << uint64(y)
	case Uint64Value:
		return v << y
	default:
		panic(internalError())
	}
}

func (v Uint64Value) rsh(b Value) Value {
	switch y := b.(type) {
	case Int64Value:
		return v >> uint64(y)
	case Uint64Value:
		return v >> y
	default:
		panic(internalError())
	}
}

func (v Uint64Value) mod(b Value) Value {
	if b.isZero() {
		return nil
	}

	return v % b.(Uint64Value)
}

type Float32Value float32

func (v Float32Value) add(b Value) Value { return v + b.(Float32Value) }
func (v Float32Value) and(b Value) Value { panic(internalError()) }
func (v Float32Value) cpl() Value        { panic(internalError()) }
func (v Float32Value) div(b Value) Value { return v / b.(Float32Value) }
func (v Float32Value) eq(b Value) Value  { return boolValue(v == b.(Float32Value)) }
func (v Float32Value) ge(b Value) Value  { return boolValue(v >= b.(Float32Value)) }
func (v Float32Value) gt(b Value) Value  { return boolValue(v > b.(Float32Value)) }
func (v Float32Value) isNonZero() bool   { return v != 0 }
func (v Float32Value) isZero() bool      { return v == 0 }
func (v Float32Value) le(b Value) Value  { return boolValue(v <= b.(Float32Value)) }
func (v Float32Value) lsh(b Value) Value { panic(internalError()) }
func (v Float32Value) lt(b Value) Value  { return boolValue(v < b.(Float32Value)) }
func (v Float32Value) mod(b Value) Value { panic(internalError()) }
func (v Float32Value) mul(b Value) Value { return v * b.(Float32Value) }
func (v Float32Value) neg() Value        { return -v }
func (v Float32Value) neq(b Value) Value { return boolValue(v != b.(Float32Value)) }
func (v Float32Value) or(b Value) Value  { panic(internalError()) }
func (v Float32Value) rsh(b Value) Value { panic(internalError()) }
func (v Float32Value) sub(b Value) Value { return v - b.(Float32Value) }
func (v Float32Value) xor(b Value) Value { panic(internalError()) }

type Float64Value float64

func (v Float64Value) add(b Value) Value { return v + b.(Float64Value) }
func (v Float64Value) and(b Value) Value { panic(internalError()) }
func (v Float64Value) cpl() Value        { panic(internalError()) }
func (v Float64Value) div(b Value) Value { return v / b.(Float64Value) }
func (v Float64Value) eq(b Value) Value  { return boolValue(v == b.(Float64Value)) }
func (v Float64Value) ge(b Value) Value  { return boolValue(v >= b.(Float64Value)) }
func (v Float64Value) gt(b Value) Value  { return boolValue(v > b.(Float64Value)) }
func (v Float64Value) isNonZero() bool   { return v != 0 }
func (v Float64Value) isZero() bool      { return v == 0 }
func (v Float64Value) le(b Value) Value  { return boolValue(v <= b.(Float64Value)) }
func (v Float64Value) lsh(b Value) Value { panic(internalError()) }
func (v Float64Value) lt(b Value) Value  { return boolValue(v < b.(Float64Value)) }
func (v Float64Value) mod(b Value) Value { panic(internalError()) }
func (v Float64Value) mul(b Value) Value { return v * b.(Float64Value) }
func (v Float64Value) neg() Value        { return -v }
func (v Float64Value) neq(b Value) Value { return boolValue(v != b.(Float64Value)) }
func (v Float64Value) or(b Value) Value  { panic(internalError()) }
func (v Float64Value) rsh(b Value) Value { panic(internalError()) }
func (v Float64Value) sub(b Value) Value { return v - b.(Float64Value) }
func (v Float64Value) xor(b Value) Value { panic(internalError()) }

type Float128Value struct {
	n   *big.Float
	nan bool
}

func (v *Float128Value) add(b Value) Value { return v.safe(b, func(x, y *big.Float) { x.Add(x, y) }) }
func (v *Float128Value) and(b Value) Value { panic(internalError()) }
func (v *Float128Value) cpl() Value        { panic(internalError()) }
func (v *Float128Value) div(b Value) Value { return v.safe(b, func(x, y *big.Float) { x.Quo(x, y) }) }
func (v *Float128Value) eq(b Value) Value  { panic(internalError()) }
func (v *Float128Value) ge(b Value) Value  { panic(internalError()) }
func (v *Float128Value) gt(b Value) Value  { return boolValue(v.cmp(b, -1, 0)) }
func (v *Float128Value) isNonZero() bool   { panic(internalError()) }
func (v *Float128Value) isZero() bool      { panic(internalError()) }
func (v *Float128Value) le(b Value) Value  { panic(internalError()) }
func (v *Float128Value) lsh(b Value) Value { panic(internalError()) }
func (v *Float128Value) lt(b Value) Value  { panic(internalError()) }
func (v *Float128Value) mod(b Value) Value { panic(internalError()) }
func (v *Float128Value) mul(b Value) Value { return v.safe(b, func(x, y *big.Float) { x.Mul(x, y) }) }
func (v *Float128Value) neg() Value        { return v.safe(nil, func(x, y *big.Float) { x.Neg(x) }) }
func (v *Float128Value) neq(b Value) Value { panic(internalError()) }
func (v *Float128Value) or(b Value) Value  { panic(internalError()) }
func (v *Float128Value) rsh(b Value) Value { panic(internalError()) }
func (v *Float128Value) sub(b Value) Value { return v.safe(b, func(x, y *big.Float) { x.Sub(x, y) }) }
func (v *Float128Value) xor(b Value) Value { panic(internalError()) }

func (v *Float128Value) cmp(b Value, accept ...int) bool {
	w := b.(*Float128Value)
	if v.nan || w.nan {
		return false
	}

	x := v.n.Cmp(w.n)
	for _, v := range accept {
		if v == x {
			return true
		}
	}
	return false
}

func (v *Float128Value) safe(b Value, f func(*big.Float, *big.Float)) Value {
	var w *Float128Value
	if b != nil {
		w = b.(*Float128Value)
	}
	if v.nan || w != nil && w.nan {
		return &Float128Value{nan: true}
	}

	r := &Float128Value{}

	defer func() {
		switch x := recover().(type) {
		case big.ErrNaN:
			r.n = nil
			r.nan = true
		case nil:
			// ok
		default:
			panic(x)
		}
	}()

	r.n = big.NewFloat(0).Set(v.n)
	var wn *big.Float
	if w != nil {
		wn = w.n
	}
	f(r.n, wn)
	return r
}

type Complex64Value complex64

func (v Complex64Value) add(b Value) Value { return v + b.(Complex64Value) }
func (v Complex64Value) and(b Value) Value { panic(internalError()) }
func (v Complex64Value) cpl() Value        { panic(internalError()) }
func (v Complex64Value) div(b Value) Value { return v / b.(Complex64Value) }
func (v Complex64Value) eq(b Value) Value  { return boolValue(v == b.(Complex64Value)) }
func (v Complex64Value) ge(b Value) Value  { panic(internalError()) }
func (v Complex64Value) gt(b Value) Value  { panic(internalError()) }
func (v Complex64Value) isNonZero() bool   { return v != 0 }
func (v Complex64Value) isZero() bool      { return v == 0 }
func (v Complex64Value) le(b Value) Value  { panic(internalError()) }
func (v Complex64Value) lsh(b Value) Value { panic(internalError()) }
func (v Complex64Value) lt(b Value) Value  { panic(internalError()) }
func (v Complex64Value) mod(b Value) Value { panic(internalError()) }
func (v Complex64Value) mul(b Value) Value { return v * b.(Complex64Value) }
func (v Complex64Value) neg() Value        { return -v }
func (v Complex64Value) neq(b Value) Value { return boolValue(v != b.(Complex64Value)) }
func (v Complex64Value) or(b Value) Value  { panic(internalError()) }
func (v Complex64Value) rsh(b Value) Value { panic(internalError()) }
func (v Complex64Value) sub(b Value) Value { return v - b.(Complex64Value) }
func (v Complex64Value) xor(b Value) Value { panic(internalError()) }

type Complex128Value complex128

func (v Complex128Value) add(b Value) Value { return v + b.(Complex128Value) }
func (v Complex128Value) and(b Value) Value { panic(internalError()) }
func (v Complex128Value) cpl() Value        { panic(internalError()) }
func (v Complex128Value) div(b Value) Value { return v / b.(Complex128Value) }
func (v Complex128Value) eq(b Value) Value  { return boolValue(v == b.(Complex128Value)) }
func (v Complex128Value) ge(b Value) Value  { panic(internalError()) }
func (v Complex128Value) gt(b Value) Value  { panic(internalError()) }
func (v Complex128Value) isNonZero() bool   { return v != 0 }
func (v Complex128Value) isZero() bool      { return v == 0 }
func (v Complex128Value) le(b Value) Value  { panic(internalError()) }
func (v Complex128Value) lsh(b Value) Value { panic(internalError()) }
func (v Complex128Value) lt(b Value) Value  { panic(internalError()) }
func (v Complex128Value) mod(b Value) Value { panic(internalError()) }
func (v Complex128Value) mul(b Value) Value { return v * b.(Complex128Value) }
func (v Complex128Value) neg() Value        { return -v }
func (v Complex128Value) neq(b Value) Value { return boolValue(v != b.(Complex128Value)) }
func (v Complex128Value) or(b Value) Value  { panic(internalError()) }
func (v Complex128Value) rsh(b Value) Value { panic(internalError()) }
func (v Complex128Value) sub(b Value) Value { return v - b.(Complex128Value) }
func (v Complex128Value) xor(b Value) Value { panic(internalError()) }

type Complex256Value struct {
	re, im *Float128Value
}

func (v Complex256Value) add(b Value) Value {
	w := b.(Complex256Value)
	return Complex256Value{v.re.add(w.re).(*Float128Value), v.im.add(w.im).(*Float128Value)}
}

func (v Complex256Value) and(b Value) Value { panic(internalError()) }
func (v Complex256Value) cpl() Value        { panic(internalError()) }
func (v Complex256Value) div(b Value) Value { panic(internalError()) }
func (v Complex256Value) eq(b Value) Value  { panic(internalError()) }
func (v Complex256Value) ge(b Value) Value  { panic(internalError()) }
func (v Complex256Value) gt(b Value) Value  { panic(internalError()) }
func (v Complex256Value) isNonZero() bool   { panic(internalError()) }
func (v Complex256Value) isZero() bool      { panic(internalError()) }
func (v Complex256Value) le(b Value) Value  { panic(internalError()) }
func (v Complex256Value) lsh(b Value) Value { panic(internalError()) }
func (v Complex256Value) lt(b Value) Value  { panic(internalError()) }
func (v Complex256Value) mod(b Value) Value { panic(internalError()) }
func (v Complex256Value) mul(b Value) Value { panic(internalError()) }
func (v Complex256Value) neg() Value        { panic(internalError()) }
func (v Complex256Value) neq(b Value) Value { panic(internalError()) }
func (v Complex256Value) or(b Value) Value  { panic(internalError()) }
func (v Complex256Value) rsh(b Value) Value { panic(internalError()) }
func (v Complex256Value) sub(b Value) Value { panic(internalError()) }
func (v Complex256Value) xor(b Value) Value { panic(internalError()) }

type lvalue struct {
	Operand
	declarator *Declarator
}

func (o *lvalue) Declarator() *Declarator { return o.declarator }
func (o *lvalue) IsLValue() bool          { return true }

func (o *lvalue) isConst() bool {
	if o.Value() != nil {
		return true
	}

	d := o.Declarator()
	return d != nil && (d.Linkage != None || d.IsStatic())
}

func (o *lvalue) convertTo(ctx *context, n Node, to Type) (r Operand) {
	return &lvalue{Operand: o.Operand.convertTo(ctx, n, to), declarator: o.declarator}
}

type funcDesignator struct {
	Operand
	declarator *Declarator
}

func (o *funcDesignator) Declarator() *Declarator { return o.declarator }
func (o *funcDesignator) IsLValue() bool          { return false }
func (o *funcDesignator) isConst() bool           { return true }

func (o *funcDesignator) convertTo(ctx *context, n Node, to Type) (r Operand) {
	return &lvalue{Operand: o.Operand.convertTo(ctx, n, to), declarator: o.declarator}
}

type operand struct {
	typ    Type
	value  Value
	offset uintptr
}

func (o *operand) Declarator() *Declarator { return nil }
func (o *operand) Offset() uintptr         { return o.offset }
func (o *operand) IsLValue() bool          { return false }
func (o *operand) IsNonZero() bool         { return o.value != nil && o.value.isNonZero() }
func (o *operand) IsZero() bool            { return o.value != nil && o.value.isZero() }
func (o *operand) Type() Type              { return o.typ }
func (o *operand) Value() Value            { return o.value }

func (o *operand) isConst() bool {
	if o.Value() != nil {
		return true
	}

	d := o.Declarator()
	return d != nil && (d.Linkage != None || d.IsStatic())
}

// [0]6.3.1.8
//
// Many operators that expect operands of arithmetic type cause conversions and
// yield result types in a similar way. The purpose is to determine a common
// real type for the operands and result. For the specified operands, each
// operand is converted, without change of type domain, to a type whose
// corresponding real type is the common real type. Unless explicitly stated
// otherwise, the common real type is also the corresponding real type of the
// result, whose type domain is the type domain of the operands if they are the
// same, and complex otherwise. This pattern is called the usual arithmetic
// conversions:
func usualArithmeticConversions(ctx *context, n Node, a, b Operand) (Operand, Operand) {
	if a.Type().Kind() == Invalid || b.Type().Kind() == Invalid {
		return noOperand, noOperand
	}

	if !a.Type().IsArithmeticType() || !b.Type().IsArithmeticType() {
		panic(internalError())
	}

	if a.Type() == nil || b.Type() == nil {
		return a, b
	}

	a = a.normalize(ctx, n)
	b = b.normalize(ctx, n)
	if a == noOperand || b == noOperand {
		return noOperand, noOperand
	}

	at := a.Type()
	bt := b.Type()
	cplx := at.IsComplexType() || bt.IsComplexType()

	// First, if the corresponding real type of either operand is long
	// double, the other operand is converted, without change of type
	// domain, to a type whose corresponding real type is long double.
	if at.Kind() == ComplexLongDouble || bt.Kind() == ComplexLongDouble || at.Kind() == LongDouble || bt.Kind() == LongDouble {
		switch {
		case cplx:
			return a.convertTo(ctx, n, ctx.cfg.ABI.Type(ComplexLongDouble)), b.convertTo(ctx, n, ctx.cfg.ABI.Type(ComplexLongDouble))
		default:
			return a.convertTo(ctx, n, ctx.cfg.ABI.Type(LongDouble)), b.convertTo(ctx, n, ctx.cfg.ABI.Type(LongDouble))
		}
	}

	// Otherwise, if the corresponding real type of either operand is
	// double, the other operand is converted, without change of type
	// domain, to a type whose corresponding real type is double.
	if at.Kind() == ComplexDouble || bt.Kind() == ComplexDouble || at.Kind() == Double || bt.Kind() == Double {
		switch {
		case cplx:
			return a.convertTo(ctx, n, ctx.cfg.ABI.Type(ComplexDouble)), b.convertTo(ctx, n, ctx.cfg.ABI.Type(ComplexDouble))
		default:
			return a.convertTo(ctx, n, ctx.cfg.ABI.Type(Double)), b.convertTo(ctx, n, ctx.cfg.ABI.Type(Double))
		}
	}

	// Otherwise, if the corresponding real type of either operand is
	// float, the other operand is converted, without change of type
	// domain, to a type whose corresponding real type is float.
	if at.Kind() == ComplexFloat || bt.Kind() == ComplexFloat || at.Kind() == Float || bt.Kind() == Float {
		switch {
		case cplx:
			return a.convertTo(ctx, n, ctx.cfg.ABI.Type(ComplexFloat)), b.convertTo(ctx, n, ctx.cfg.ABI.Type(ComplexFloat))
		default:
			return a.convertTo(ctx, n, ctx.cfg.ABI.Type(Float)), b.convertTo(ctx, n, ctx.cfg.ABI.Type(Float))
		}
	}

	if cplx {
		panic(internalErrorf("TODO %v, %v", at, bt))
	}

	if !a.Type().IsIntegerType() || !b.Type().IsIntegerType() {
		panic(internalError())
	}

	// Otherwise, the integer promotions are performed on both operands.
	a = a.integerPromotion(ctx, n)
	b = b.integerPromotion(ctx, n)
	at = a.Type()
	bt = b.Type()

	// Then the following rules are applied to the promoted operands:

	// If both operands have the same type, then no further conversion is
	// needed.
	if at.Kind() == bt.Kind() {
		return a, b
	}

	// Otherwise, if both operands have signed integer types or both have
	// unsigned integer types, the operand with the type of lesser integer
	// conversion rank is converted to the type of the operand with greater
	// rank.
	abi := ctx.cfg.ABI
	if abi.isSignedInteger(at.Kind()) == abi.isSignedInteger(bt.Kind()) {
		t := a.Type()
		if intConvRank[bt.Kind()] > intConvRank[at.Kind()] {
			t = b.Type()
		}
		return a.convertTo(ctx, n, t), b.convertTo(ctx, n, t)

	}

	// Otherwise, if the operand that has unsigned integer type has rank
	// greater or equal to the rank of the type of the other operand, then
	// the operand with signed integer type is converted to the type of the
	// operand with unsigned integer type.
	switch {
	case a.Type().IsSignedType(): // b is unsigned
		if intConvRank[bt.Kind()] >= intConvRank[a.Type().Kind()] {
			return a.convertTo(ctx, n, b.Type()), b
		}
	case b.Type().IsSignedType(): // a is unsigned
		if intConvRank[at.Kind()] >= intConvRank[b.Type().Kind()] {
			return a, b.convertTo(ctx, n, a.Type())
		}
	default:
		panic(fmt.Errorf("TODO %v %v", a, b))
	}

	// Otherwise, if the type of the operand with signed integer type can
	// represent all of the values of the type of the operand with unsigned
	// integer type, then the operand with unsigned integer type is
	// converted to the type of the operand with signed integer type.
	var signed Type
	switch {
	case abi.isSignedInteger(at.Kind()): // b is unsigned
		signed = a.Type()
		if at.Size() > bt.Size() {
			return a, b.convertTo(ctx, n, a.Type())
		}
	case abi.isSignedInteger(bt.Kind()): // a is unsigned
		signed = b.Type()
		if bt.Size() > at.Size() {
			return a.convertTo(ctx, n, b.Type()), b
		}

	}

	// Otherwise, both operands are converted to the unsigned integer type
	// corresponding to the type of the operand with signed integer type.
	var typ Type
	switch signed.Kind() {
	case Int:
		//TODO if a.IsEnumConst || b.IsEnumConst {
		//TODO 	return a, b
		//TODO }

		typ = abi.Type(UInt)
	case Long:
		typ = abi.Type(ULong)
	case LongLong:
		typ = abi.Type(ULongLong)
	default:
		panic(internalError())
	}
	return a.convertTo(ctx, n, typ), b.convertTo(ctx, n, typ)
}

// [0]6.3.1.1-2
//
// If an int can represent all values of the original type, the value is
// converted to an int; otherwise, it is converted to an unsigned int. These
// are called the integer promotions. All other types are unchanged by the
// integer promotions.
func (o *operand) integerPromotion(ctx *context, n Node) Operand {
	t := o.Type()
	if t2 := integerPromotion(ctx, t); t2.Kind() != t.Kind() {
		return o.convertTo(ctx, n, t2)
	}

	return o
}

// [0]6.3.1.1-2
//
// If an int can represent all values of the original type, the value is
// converted to an int; otherwise, it is converted to an unsigned int. These
// are called the integer promotions. All other types are unchanged by the
// integer promotions.
func integerPromotion(ctx *context, t Type) Type {
	// github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf-sign-2.c
	//
	// This test checks promotion of bitfields.  Bitfields
	// should be promoted very much like chars and shorts:
	//
	// Bitfields (signed or unsigned) should be promoted to
	// signed int if their value will fit in a signed int,
	// otherwise to an unsigned int if their value will fit
	// in an unsigned int, otherwise we don't promote them
	// (ANSI/ISO does not specify the behavior of bitfields
	// larger than an unsigned int).
	if t.IsBitFieldType() {
		f := t.BitField()
		intBits := int(ctx.cfg.ABI.Types[Int].Size) * 8
		switch {
		case t.IsSignedType():
			if f.BitFieldWidth() < intBits-1 {
				return ctx.cfg.ABI.Type(Int)
			}
		default:
			if f.BitFieldWidth() < intBits {
				return ctx.cfg.ABI.Type(Int)
			}
		}
		return t
	}

	switch t.Kind() {
	case Invalid:
		return t
	case Char, SChar, UChar, Short, UShort:
		return ctx.cfg.ABI.Type(Int)
	default:
		return t
	}
}

func (o *operand) convertTo(ctx *context, n Node, to Type) Operand {
	if o.Type().Kind() == Invalid {
		return o
	}

	v := o.Value()
	r := &operand{typ: to, offset: o.offset, value: v}
	if v == nil {
		return r
	}

	if o.Type().Kind() == to.Kind() {
		return r.normalize(ctx, n)
	}

	if o.Type().IsIntegerType() {
		return o.convertFromInt(ctx, n, to)
	}

	if to.IsIntegerType() {
		return o.convertToInt(ctx, n, to)
	}

	switch o.Type().Kind() {
	case Array:
		switch to.Kind() {
		case Ptr:
			return r
		default:
			panic("TODO653")
		}
	case ComplexFloat:
		v := v.(Complex64Value)
		switch to.Kind() {
		case ComplexDouble:
			r.value = Complex128Value(v)
		case Float:
			r.value = Float32Value(real(v))
		case Double:
			r.value = Float64Value(real(v))
		case ComplexLongDouble:
			panic("TODO663")
		default:
			panic("TODO665")
		}
	case ComplexDouble:
		v := v.(Complex128Value)
		switch to.Kind() {
		case ComplexFloat:
			r.value = Complex64Value(v)
		case ComplexLongDouble:
			panic("TODO673")
		case Float:
			r.value = Float32Value(real(v))
		case Double:
			r.value = Float64Value(real(v))
		default:
			panic("TODO681")
		}
	case Float:
		v := v.(Float32Value)
		switch to.Kind() {
		case ComplexFloat:
			r.value = Complex64Value(complex(v, 0))
		case ComplexDouble:
			r.value = Complex128Value(complex(v, 0))
		case Double:
			r.value = Float64Value(v)
		case ComplexLongDouble:
			panic("TODO693")
		case LongDouble:
			r.value = &Float128Value{n: big.NewFloat(float64(v))}
		default:
			panic(fmt.Sprintf("TODO695 %s", to.Kind()))
		}
	case Double:
		v := v.(Float64Value)
		switch to.Kind() {
		case ComplexFloat:
			r.value = Complex64Value(complex(v, 0))
		case ComplexDouble:
			r.value = Complex128Value(complex(v, 0))
		case LongDouble:
			f := float64(v)
			switch {
			case math.IsNaN(f):
				r.value = &Float128Value{nan: true}
			default:
				r.value = &Float128Value{n: big.NewFloat(f)}
			}
		case Float:
			r.value = Float32Value(v)
		case ComplexLongDouble:
			panic("TODO709")
		default:
			panic("TODO711")
		}
	case LongDouble:
		v := v.(*Float128Value)
		switch to.Kind() {
		case Double:
			if v.nan {
				r.value = Float64Value(math.NaN())
				break
			}

			d, _ := v.n.Float64()
			r.value = Float64Value(d)
		case ComplexLongDouble:
			if v.nan {
				r.value = Complex256Value{v, &Float128Value{nan: true}}
				break
			}

			r.value = Complex256Value{v, &Float128Value{n: big.NewFloat(0)}}
		default:
			panic(fmt.Sprintf("TODO813 %v", to.Kind()))
		}
	default:
		panic(internalErrorf("%v: %v -> %v %v", n.Position(), o.Type(), to, to.Kind()))
	}
	return r.normalize(ctx, n)
}

type signedSaturationLimit struct {
	fmin, fmax float64
	min, max   int64
}

type unsignedSaturationLimit struct {
	fmax float64
	max  uint64
}

var (
	signedSaturationLimits = [...]signedSaturationLimit{
		1: {math.Nextafter(math.MinInt8, 0), math.Nextafter(math.MaxInt8, 0), math.MinInt8, math.MaxInt8},
		2: {math.Nextafter(math.MinInt16, 0), math.Nextafter(math.MaxInt16, 0), math.MinInt16, math.MaxInt16},
		4: {math.Nextafter(math.MinInt32, 0), math.Nextafter(math.MaxInt32, 0), math.MinInt32, math.MaxInt32},
		8: {math.Nextafter(math.MinInt64, 0), math.Nextafter(math.MaxInt64, 0), math.MinInt32, math.MaxInt64},
	}

	unsignedSaturationLimits = [...]unsignedSaturationLimit{
		1: {math.Nextafter(math.MaxUint8, 0), math.MaxUint8},
		2: {math.Nextafter(math.MaxUint16, 0), math.MaxUint16},
		4: {math.Nextafter(math.MaxUint32, 0), math.MaxUint32},
		8: {math.Nextafter(math.MaxUint64, 0), math.MaxUint64},
	}
)

func (o *operand) convertToInt(ctx *context, n Node, to Type) (r Operand) {
	v := o.Value()
	switch o.Type().Kind() {
	case Float:
		v := float64(v.(Float32Value))
		switch {
		case to.IsSignedType():
			limits := &signedSaturationLimits[to.Size()]
			if v > limits.fmax {
				return (&operand{typ: to, value: Int64Value(limits.max)}).normalize(ctx, n)
			}

			if v < limits.fmin {
				return (&operand{typ: to, value: Int64Value(limits.min)}).normalize(ctx, n)
			}

			return (&operand{typ: to, value: Int64Value(v)}).normalize(ctx, n)
		default:
			limits := &unsignedSaturationLimits[to.Size()]
			if v > limits.fmax {
				return (&operand{typ: to, value: Uint64Value(limits.max)}).normalize(ctx, n)
			}

			return (&operand{typ: to, value: Uint64Value(v)}).normalize(ctx, n)
		}
	case Double:
		v := float64(v.(Float64Value))
		switch {
		case to.IsSignedType():
			limits := &signedSaturationLimits[to.Size()]
			if v > limits.fmax {
				return (&operand{typ: to, value: Int64Value(limits.max)}).normalize(ctx, n)
			}

			if v < limits.fmin {
				return (&operand{typ: to, value: Int64Value(limits.min)}).normalize(ctx, n)
			}

			return (&operand{typ: to, value: Int64Value(v)}).normalize(ctx, n)
		default:
			limits := &unsignedSaturationLimits[to.Size()]
			if v > limits.fmax {
				return (&operand{typ: to, value: Uint64Value(limits.max)}).normalize(ctx, n)
			}

			return (&operand{typ: to, value: Uint64Value(v)}).normalize(ctx, n)
		}
	case LongDouble:
		panic("TODO791")
	case Ptr:
		var v uint64
		switch x := o.Value().(type) {
		case Int64Value:
			v = uint64(x)
		case Uint64Value:
			v = uint64(x)
		case *InitializerValue:
			return (&operand{typ: to})
		default:
			panic(internalErrorf("%v: %T", n.Position(), x))
		}
		switch {
		case to.IsSignedType():
			return (&operand{typ: to, value: Int64Value(v)}).normalize(ctx, n)
		default:
			return (&operand{typ: to, value: Uint64Value(v)}).normalize(ctx, n)
		}
	case Array:
		return &operand{typ: to}
	}
	panic("TODO")
}

func (o *operand) convertFromInt(ctx *context, n Node, to Type) (r Operand) {
	var v uint64
	switch x := o.Value().(type) {
	case Int64Value:
		v = uint64(x)
	case Uint64Value:
		v = uint64(x)
	default:
		ctx.errNode(n, "conversion to integer: invalid value")
		return &operand{typ: to}
	}

	if to.IsIntegerType() {
		switch {
		case to.IsSignedType():
			return (&operand{typ: to, value: Int64Value(v)}).normalize(ctx, n)
		default:
			return (&operand{typ: to, value: Uint64Value(v)}).normalize(ctx, n)
		}
	}

	switch to.Kind() {
	case ComplexFloat:
		switch {
		case o.Type().IsSignedType():
			return (&operand{typ: to, value: Complex64Value(complex(float64(int64(v)), 0))}).normalize(ctx, n)
		default:
			return (&operand{typ: to, value: Complex64Value(complex(float64(v), 0))}).normalize(ctx, n)
		}
	case ComplexDouble:
		switch {
		case o.Type().IsSignedType():
			return (&operand{typ: to, value: Complex128Value(complex(float64(int64(v)), 0))}).normalize(ctx, n)
		default:
			return (&operand{typ: to, value: Complex128Value(complex(float64(v), 0))}).normalize(ctx, n)
		}
	case Float:
		switch {
		case o.Type().IsSignedType():
			return (&operand{typ: to, value: Float32Value(float64(int64(v)))}).normalize(ctx, n)
		default:
			return (&operand{typ: to, value: Float32Value(float64(v))}).normalize(ctx, n)
		}
	case ComplexLongDouble:
		panic("TODO896")
	case Double:
		switch {
		case o.Type().IsSignedType():
			return (&operand{typ: to, value: Float64Value(int64(v))}).normalize(ctx, n)
		default:
			return (&operand{typ: to, value: Float64Value(v)}).normalize(ctx, n)
		}
	case LongDouble:
		switch {
		case o.Type().IsSignedType():
			return (&operand{typ: to, value: &Float128Value{n: big.NewFloat(0).SetInt64(int64(v))}}).normalize(ctx, n)
		default:
			return (&operand{typ: to, value: &Float128Value{n: big.NewFloat(0).SetUint64(v)}}).normalize(ctx, n)
		}
	case Ptr:
		return (&operand{typ: to, value: Uint64Value(v)}).normalize(ctx, n)
	case Struct, Union, Void, Int128, UInt128:
		return &operand{typ: to}
	}
	panic(internalErrorf("%q, %q", to, to.Kind()))
}

func (o *operand) normalize(ctx *context, n Node) (r Operand) {
	if o.Type().IsIntegerType() {
		switch {
		case o.Type().IsSignedType():
			if x, ok := o.value.(Uint64Value); ok {
				o.value = Int64Value(x)
			}
		default:
			if x, ok := o.value.(Int64Value); ok {
				o.value = Uint64Value(x)
			}
		}
		switch x := o.Value().(type) {
		case Int64Value:
			if v := convertInt64(int64(x), o.Type(), ctx); v != int64(x) {
				o.value = Int64Value(v)
			}
		case Uint64Value:
			v := uint64(x)
			switch o.Type().Size() {
			case 1:
				v &= 0xff
			case 2:
				v &= 0xffff
			case 4:
				v &= 0xffffffff
			}
			if v != uint64(x) {
				o.value = Uint64Value(v)
			}
		case *InitializerValue, nil:
			// ok
		default:
			panic(internalErrorf("%T %v", x, x))
		}
		return o
	}

	switch o.Type().Kind() {
	case ComplexFloat:
		switch o.Value().(type) {
		case Complex64Value, nil:
			return o
		default:
			panic(internalError())
		}
	case ComplexDouble:
		switch o.Value().(type) {
		case Complex128Value, nil:
			return o
		default:
			panic(internalError())
		}
	case ComplexLongDouble:
		switch o.Value().(type) {
		case Complex256Value, nil:
			return o
		default:
			panic(fmt.Sprintf("TODO934 %v", o.Type().Kind()))
		}
	case Float:
		switch o.Value().(type) {
		case Float32Value, *InitializerValue, nil:
			return o
		default:
			panic(internalError())
		}
	case Double:
		switch x := o.Value().(type) {
		case Float64Value, *InitializerValue, nil:
			return o
		default:
			panic(internalErrorf("%T %v", x, x))
		}
	case LongDouble:
		switch x := o.Value().(type) {
		case *Float128Value, nil:
			return o
		default:
			panic(internalErrorf("%T %v TODO980 %v", x, x, n.Position()))
		}
	case Ptr:
		switch o.Value().(type) {
		case Int64Value, Uint64Value, *InitializerValue, StringValue, WideStringValue, nil:
			return o
		default:
			panic(internalError())
		}
	case Array, Void, Function, Struct, Union:
		return o
	case ComplexChar, ComplexInt, ComplexLong, ComplexLongLong, ComplexShort, ComplexUInt, ComplexUShort:
		ctx.errNode(n, "unsupported type: %s", o.Type())
		return noOperand
	}
	panic(internalErrorf("%v, %v", o.Type(), o.Type().Kind()))
}

func convertInt64(n int64, t Type, ctx *context) int64 {
	abi := ctx.cfg.ABI
	k := t.Kind()
	if k == Enum {
		//TODO
	}
	signed := abi.isSignedInteger(k)
	switch sz := abi.size(k); sz {
	case 1:
		switch {
		case signed:
			switch {
			case int8(n) < 0:
				return n | ^math.MaxUint8
			default:
				return n & math.MaxUint8
			}
		default:
			return n & math.MaxUint8
		}
	case 2:
		switch {
		case signed:
			switch {
			case int16(n) < 0:
				return n | ^math.MaxUint16
			default:
				return n & math.MaxUint16
			}
		default:
			return n & math.MaxUint16
		}
	case 4:
		switch {
		case signed:
			switch {
			case int32(n) < 0:
				return n | ^math.MaxUint32
			default:
				return n & math.MaxUint32
			}
		default:
			return n & math.MaxUint32
		}
	case 8:
		return n
	default:
		panic(internalError())
	}
}

func boolValue(b bool) Value {
	if b {
		return Int64Value(1)
	}

	return Int64Value(0)
}

type initializer interface {
	List() []*Initializer
	IsConst() bool
}

type InitializerValue struct {
	typ         Type
	initializer initializer
}

func (v *InitializerValue) IsConst() bool        { return v.initializer.IsConst() }
func (v *InitializerValue) List() []*Initializer { return v.initializer.List() }
func (v *InitializerValue) Type() Type           { return v.typ }
func (v *InitializerValue) add(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) and(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) cpl() Value           { panic(internalError()) }
func (v *InitializerValue) div(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) eq(b Value) Value     { panic(internalError()) }
func (v *InitializerValue) ge(b Value) Value     { panic(internalError()) }
func (v *InitializerValue) gt(b Value) Value     { panic(internalError()) }
func (v *InitializerValue) le(b Value) Value     { panic(internalError()) }
func (v *InitializerValue) lsh(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) lt(b Value) Value     { panic(internalError()) }
func (v *InitializerValue) mod(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) mul(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) neg() Value           { panic(internalError()) }
func (v *InitializerValue) neq(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) or(b Value) Value     { panic(internalError()) }
func (v *InitializerValue) rsh(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) sub(b Value) Value    { panic(internalError()) }
func (v *InitializerValue) xor(b Value) Value    { panic(internalError()) }

func (v *InitializerValue) isNonZero() bool {
	for _, v := range v.List() {
		if v.AssignmentExpression.Operand.IsNonZero() {
			return true
		}
	}
	return false
}

func (v *InitializerValue) isZero() bool {
	for _, v := range v.List() {
		if v.AssignmentExpression.Operand.IsNonZero() {
			return false
		}
	}
	return true
}
