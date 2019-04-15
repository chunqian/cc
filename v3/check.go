// Copyright 2019 The CC Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc // import "modernc.org/cc/v3"

import (
	"math/bits"
)

type mode = int

const (
	// [2], 6.6 Constant expressions, 6
	//
	// An integer constant expression shall have integer type and shall
	// only have operands that are integer constants, enumeration
	// constants, character constants, sizeof expressions whose results are
	// integer constants, _Alignof expressions, and floating constants that
	// are the immediate operands of casts. Cast operators in an integer
	// constant expression shall only convert arithmetic types to integer
	// types, except as part of an operand to the sizeof or _Alignof
	// operator.
	mIntConstExpr = 1 << iota

	mIntConstExprFloat   // As mIntConstExpr plus accept floating point constants.
	mIntConstExprAnyCast // As mIntConstExpr plus accept any cast.
)

func (n *TranslationUnit) check(ctx *context) {
	for ; n != nil; n = n.TranslationUnit {
		n.ExternalDeclaration.check(ctx)
	}
}

func (n *ExternalDeclaration) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case ExternalDeclarationFuncDef: // FunctionDefinition
		n.FunctionDefinition.check(ctx)
	case ExternalDeclarationDecl: // Declaration
		n.Declaration.check(ctx)
	case ExternalDeclarationAsm: // AsmFunctionDefinition
		n.AsmFunctionDefinition.check(ctx)
	case ExternalDeclarationAsmStmt: // AsmStatement
		n.AsmStatement.check(ctx)
	case ExternalDeclarationEmpty: // ';'
		// nop
	case ExternalDeclarationPragma: // PragmaSTDC
		n.PragmaSTDC.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *PragmaSTDC) check(ctx *context) {
	// nop
}

func (n *AsmFunctionDefinition) check(ctx *context) {
	if n == nil {
		return
	}

	typ := n.DeclarationSpecifiers.check(ctx)
	n.Declarator.check(ctx, typ)
	n.AsmStatement.check(ctx)
}

func (n *AsmStatement) check(ctx *context) {
	if n == nil {
		return
	}

	n.Asm.check(ctx)
	n.AttributeSpecifierList.check(ctx)
}

func (n *Declaration) check(ctx *context) {
	if n == nil {
		return
	}

	typ := n.DeclarationSpecifiers.check(ctx)
	n.InitDeclaratorList.check(ctx, typ)
}

func (n *InitDeclaratorList) check(ctx *context, typ Type) {
	for ; n != nil; n = n.InitDeclaratorList {
		n.AttributeSpecifierList.check(ctx)
		n.InitDeclarator.check(ctx, typ)
	}
}

func (n *InitDeclarator) check(ctx *context, typ Type) {
	if n == nil {
		return
	}

	switch n.Case {
	case InitDeclaratorDecl: // Declarator AttributeSpecifierList
		n.Declarator.check(ctx, typ)
		n.AttributeSpecifierList.check(ctx)
	case InitDeclaratorInit: // Declarator AttributeSpecifierList '=' Initializer
		n.Declarator.check(ctx, typ)
		n.AttributeSpecifierList.check(ctx)
		n.Initializer.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *Initializer) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case InitializerExpr: // AssignmentExpression
		n.AssignmentExpression.check(ctx)
	case InitializerInitList: // '{' InitializerList ',' '}'
		n.InitializerList.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *InitializerList) check(ctx *context) {
	for ; n != nil; n = n.InitializerList {
		n.Designation.check(ctx)
		n.Initializer.check(ctx)
	}
}

func (n *Designation) check(ctx *context) {
	if n == nil {
		return
	}

	n.DesignatorList.check(ctx)
}

func (n *DesignatorList) check(ctx *context) {
	for ; n != nil; n = n.DesignatorList {
		n.Designator.check(ctx)
	}
}

func (n *Designator) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case DesignatorIndex: // '[' ConstantExpression ']'
		n.ConstantExpression.check(ctx, ctx.mode|mIntConstExpr)
		//TODO
	case DesignatorField: // '.' IDENTIFIER
		//TODO
	case DesignatorField2: // IDENTIFIER ':'
		//TODO
	default:
		panic("internal error") //TODOOK
	}
}

func (n *AssignmentExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case AssignmentExpressionCond: // ConditionalExpression
		n.ConditionalExpression.check(ctx)
	case AssignmentExpressionAssign: // UnaryExpression '=' AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionMul: // UnaryExpression "*=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionDiv: // UnaryExpression "/=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionMod: // UnaryExpression "%=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionAdd: // UnaryExpression "+=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionSub: // UnaryExpression "-=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionLsh: // UnaryExpression "<<=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionRsh: // UnaryExpression ">>=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionAnd: // UnaryExpression "&=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionXor: // UnaryExpression "^=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	case AssignmentExpressionOr: // UnaryExpression "|=" AssignmentExpression
		n.UnaryExpression.check(ctx)
		n.AssignmentExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *UnaryExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case UnaryExpressionPostfix: // PostfixExpression
		n.PostfixExpression.check(ctx)
	case UnaryExpressionInc: // "++" UnaryExpression
		n.UnaryExpression.check(ctx)
	case UnaryExpressionDec: // "--" UnaryExpression
		n.UnaryExpression.check(ctx)
	case UnaryExpressionAddrof: // '&' CastExpression
		ctx.not(n, mIntConstExpr)
		n.CastExpression.check(ctx)
	case UnaryExpressionDeref: // '*' CastExpression
		ctx.not(n, mIntConstExpr)
		n.CastExpression.check(ctx)
	case UnaryExpressionPlus: // '+' CastExpression
		n.CastExpression.check(ctx)
	case UnaryExpressionMinus: // '-' CastExpression
		n.CastExpression.check(ctx)
	case UnaryExpressionCpl: // '~' CastExpression
		n.CastExpression.check(ctx)
	case UnaryExpressionNot: // '!' CastExpression
		n.CastExpression.check(ctx)
	case UnaryExpressionSizeofExpr: // "sizeof" UnaryExpression
		ctx.push(ctx.mode &^ mIntConstExpr)
		n.UnaryExpression.check(ctx)
		ctx.pop()
		//TODO
	case UnaryExpressionSizeofType: // "sizeof" '(' TypeName ')'
		ctx.push(ctx.mode)
		if ctx.mode&mIntConstExpr != 0 {
			ctx.mode |= mIntConstExprAnyCast
		}
		n.TypeName.check(ctx)
		ctx.pop()
		//TODO
	case UnaryExpressionLabelAddr: // "&&" IDENTIFIER
		ctx.not(n, mIntConstExpr)
		//TODO
	case UnaryExpressionAlignofExpr: // "_Alignof" UnaryExpression
		ctx.push(ctx.mode &^ mIntConstExpr)
		n.UnaryExpression.check(ctx)
		ctx.pop()
	case UnaryExpressionAlignofType: // "_Alignof" '(' TypeName ')'
		ctx.push(ctx.mode)
		if ctx.mode&mIntConstExpr != 0 {
			ctx.mode |= mIntConstExprAnyCast
		}
		n.TypeName.check(ctx)
		ctx.pop()
	case UnaryExpressionImag: // "__imag__" UnaryExpression
		ctx.not(n, mIntConstExpr)
		n.UnaryExpression.check(ctx)
	case UnaryExpressionReal: // "__real__" UnaryExpression
		ctx.not(n, mIntConstExpr)
		n.UnaryExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *TypeName) check(ctx *context) {
	if n == nil {
		return
	}

	typ := n.SpecifierQualifierList.check(ctx)
	n.typ = n.AbstractDeclarator.check(ctx, typ)
}

func (n *AbstractDeclarator) check(ctx *context, typ Type) Type {
	if n == nil {
		return nil
	}

	switch n.Case {
	case AbstractDeclaratorPtr: // Pointer
		n.Pointer.check(ctx)
		return nil //TODO
	case AbstractDeclaratorDecl: // Pointer DirectAbstractDeclarator
		n.Pointer.check(ctx)
		n.DirectAbstractDeclarator.check(ctx)
		return nil //TODO
	default:
		panic("internal error") //TODOOK
	}
}

func (n *DirectAbstractDeclarator) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case DirectAbstractDeclaratorDecl: // '(' AbstractDeclarator ')'
		n.AbstractDeclarator.check(ctx, nil)
	case DirectAbstractDeclaratorArr: // DirectAbstractDeclarator '[' TypeQualifiers AssignmentExpression ']'
		n.DirectAbstractDeclarator.check(ctx)
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		n.AssignmentExpression.check(ctx)
	case DirectAbstractDeclaratorStaticArr: // DirectAbstractDeclarator '[' "static" TypeQualifiers AssignmentExpression ']'
		n.DirectAbstractDeclarator.check(ctx)
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		n.AssignmentExpression.check(ctx)
	case DirectAbstractDeclaratorArrStatic: // DirectAbstractDeclarator '[' TypeQualifiers "static" AssignmentExpression ']'
		n.DirectAbstractDeclarator.check(ctx)
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		n.AssignmentExpression.check(ctx)
	case DirectAbstractDeclaratorArrStar: // DirectAbstractDeclarator '[' '*' ']'
		n.DirectAbstractDeclarator.check(ctx)
	case DirectAbstractDeclaratorFunc: // DirectAbstractDeclarator '(' ParameterTypeList ')'
		n.DirectAbstractDeclarator.check(ctx)
		n.ParameterTypeList.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ParameterTypeList) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case ParameterTypeListList: // ParameterList
		n.ParameterList.check(ctx)
	case ParameterTypeListVar: // ParameterList ',' "..."
		n.ParameterList.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ParameterList) check(ctx *context) {
	for ; n != nil; n = n.ParameterList {
		n.ParameterDeclaration.check(ctx)
	}
}

func (n *ParameterDeclaration) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case ParameterDeclarationDecl: // DeclarationSpecifiers Declarator AttributeSpecifierList
		typ := n.DeclarationSpecifiers.check(ctx)
		n.Declarator.check(ctx, typ)
		n.AttributeSpecifierList.check(ctx)
	case ParameterDeclarationAbstract: // DeclarationSpecifiers AbstractDeclarator
		typ := n.DeclarationSpecifiers.check(ctx)
		n.AbstractDeclarator.check(ctx, typ)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *Pointer) check(ctx *context) (t Type) {
	t = noType
	if n == nil {
		return t
	}

	switch n.Case {
	case PointerTypeQual: // '*' TypeQualifiers
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
	case PointerPtr: // '*' TypeQualifiers Pointer
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		n.Pointer.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
	if n.typeQualifiers != nil {
		t = n.typeQualifiers.check(ctx, (*DeclarationSpecifiers)(nil))
	}
	return t
}

func (n *TypeQualifiers) check(ctx *context, typ **baseType) {
	for ; n != nil; n = n.TypeQualifiers {
		switch n.Case {
		case TypeQualifiersTypeQual: // TypeQualifier
			if *typ == nil {
				*typ = &baseType{}
			}
			n.TypeQualifier.check(ctx, *typ)
		case TypeQualifiersAttribute: // AttributeSpecifier
			n.AttributeSpecifier.check(ctx)
		default:
			panic("internal error") //TODOOK
		}
	}
}

func (n *TypeQualifier) check(ctx *context, typ *baseType) {
	if n == nil {
		return
	}

	switch n.Case {
	case TypeQualifierConst: // "const"
		typ.flags |= btfConst
	case TypeQualifierRestrict: // "restrict"
		typ.flags |= btfRestrict
	case TypeQualifierVolatile: // "volatile"
		typ.flags |= btfVolatile
	case TypeQualifierAtomic: // "_Atomic"
		typ.flags |= btfAtomic
	default:
		panic("internal error") //TODOOK
	}
}

func (n *SpecifierQualifierList) check(ctx *context) Type {
	n0 := n
	typ := &baseType{}
	for ; n != nil; n = n.SpecifierQualifierList {
		switch n.Case {
		case SpecifierQualifierListTypeSpec: // TypeSpecifier SpecifierQualifierList
			n.TypeSpecifier.check(ctx)
		case SpecifierQualifierListTypeQual: // TypeQualifier SpecifierQualifierList
			n.TypeQualifier.check(ctx, typ)
		case SpecifierQualifierListAlignSpec: // AlignmentSpecifier SpecifierQualifierList
			n.AlignmentSpecifier.check(ctx)
		case SpecifierQualifierListAttribute: // AttributeSpecifier SpecifierQualifierList
			n.AttributeSpecifier.check(ctx)
		default:
			panic("internal error") //TODOOK
		}
	}
	return typ.check(ctx, n0)
}

func (n *TypeSpecifier) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case
		TypeSpecifierVoid,       // "void"
		TypeSpecifierChar,       // "char"
		TypeSpecifierShort,      // "short"
		TypeSpecifierInt,        // "int"
		TypeSpecifierInt128,     // "__int128"
		TypeSpecifierLong,       // "long"
		TypeSpecifierFloat,      // "float"
		TypeSpecifierFloat16,    // "__fp16"
		TypeSpecifierDecimal32,  // "_Decimal32"
		TypeSpecifierDecimal64,  // "_Decimal64"
		TypeSpecifierDecimal128, // "_Decimal128"
		TypeSpecifierFloat32,    // "_Float32"
		TypeSpecifierFloat32x,   // "_Float32x"
		TypeSpecifierFloat64,    // "_Float64"
		TypeSpecifierFloat64x,   // "_Float64x"
		TypeSpecifierFloat128,   // "_Float128"
		TypeSpecifierFloat80,    // "__float80"
		TypeSpecifierDouble,     // "double"
		TypeSpecifierSigned,     // "signed"
		TypeSpecifierUnsigned,   // "unsigned"
		TypeSpecifierBool,       // "_Bool"
		TypeSpecifierComplex:    // "_Complex"
		// nop
	case TypeSpecifierStruct: // StructOrUnionSpecifier
		n.StructOrUnionSpecifier.check(ctx)
	case TypeSpecifierEnum: // EnumSpecifier
		n.EnumSpecifier.check(ctx)
	case TypeSpecifierTypedefName: // TYPEDEFNAME
		// nop
	case TypeSpecifierTypeofExpr: // "typeof" '(' Expression ')'
		n.Expression.check(ctx)
	case TypeSpecifierTypeofType: // "typeof" '(' TypeName ')'
		n.TypeName.check(ctx)
	case TypeSpecifierAtomic: // AtomicTypeSpecifier
		n.AtomicTypeSpecifier.check(ctx)
	case
		TypeSpecifierFract, // "_Fract"
		TypeSpecifierSat,   // "_Sat"
		TypeSpecifierAccum: // "_Accum"
		// nop
	default:
		panic("TODO") //TODOOK
	}
}

func (n *AtomicTypeSpecifier) check(ctx *context) {
	if n == nil {
		return
	}

	n.TypeName.check(ctx)
}

func (n *EnumSpecifier) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case EnumSpecifierDef: // "enum" AttributeSpecifierList IDENTIFIER '{' EnumeratorList ',' '}'
		n.AttributeSpecifierList.check(ctx)
		n.EnumeratorList.check(ctx)
	case EnumSpecifierTag: // "enum" AttributeSpecifierList IDENTIFIER
		n.AttributeSpecifierList.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *EnumeratorList) check(ctx *context) {
	for ; n != nil; n = n.EnumeratorList {
		n.Enumerator.check(ctx)
	}
}

func (n *Enumerator) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case EnumeratorIdent: // IDENTIFIER AttributeSpecifierList
		n.AttributeSpecifierList.check(ctx)
	case EnumeratorExpr: // IDENTIFIER AttributeSpecifierList '=' ConstantExpression
		n.AttributeSpecifierList.check(ctx)
		n.ConstantExpression.check(ctx, ctx.mode|mIntConstExpr)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ConstantExpression) check(ctx *context, mode mode) {
	if n == nil {
		return
	}

	ctx.push(mode)
	n.ConditionalExpression.check(ctx)
	ctx.pop()
}

func (n *StructOrUnionSpecifier) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case StructOrUnionSpecifierDef: // StructOrUnion AttributeSpecifierList IDENTIFIER '{' StructDeclarationList '}'
		n.StructOrUnion.check(ctx)
		n.AttributeSpecifierList.check(ctx)
		n.StructDeclarationList.check(ctx)
	case StructOrUnionSpecifierTag: // StructOrUnion AttributeSpecifierList IDENTIFIER
		n.StructOrUnion.check(ctx)
		n.AttributeSpecifierList.check(ctx)
	default:
		panic("interanal error") //TODOOK
	}
}

func (n *StructDeclarationList) check(ctx *context) (s []*field) {
	for ; n != nil; n = n.StructDeclarationList {
		s = append(s, n.StructDeclaration.check(ctx)...)
	}
	return s
}

func (n *StructDeclaration) check(ctx *context) (s []*field) {
	if n == nil {
		return nil
	}

	typ := n.SpecifierQualifierList.check(ctx)
	return n.StructDeclaratorList.check(ctx, typ)
}

func (n *StructDeclaratorList) check(ctx *context, typ Type) (s []*field) {
	for ; n != nil; n = n.StructDeclaratorList {
		s = append(s, n.StructDeclarator.check(ctx, typ))
	}
	return s
}

func (n *StructDeclarator) check(ctx *context, typ Type) *field {
	if n == nil {
		return nil
	}

	sf := &field{
		typ: typ,
	}
	switch n.Case {
	case StructDeclaratorDecl: // Declarator
		n.Declarator.check(ctx, typ)
		sf.name = n.Declarator.Name()
	case StructDeclaratorBitField: // Declarator ':' ConstantExpression AttributeSpecifierList
		n.Declarator.check(ctx, typ)
		sf.name = n.Declarator.Name()
		n.ConstantExpression.check(ctx, ctx.mode|mIntConstExpr)
		n.AttributeSpecifierList.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
	return sf
}

func (n *StructOrUnion) check(ctx *context) (isStruct bool) {
	if n == nil {
		return false
	}

	switch n.Case {
	case StructOrUnionStruct: // "struct"
		return true
	case StructOrUnionUnion: // "union"
		return false
	default:
		panic("internal error") //TODOOK
	}
}

func (n *CastExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case CastExpressionUnary: // UnaryExpression
		n.UnaryExpression.check(ctx)
	case CastExpressionCast: // '(' TypeName ')' CastExpression
		n.TypeName.check(ctx)
		ctx.push(ctx.mode)
		if m := ctx.mode; m&mIntConstExpr != 0 && m&mIntConstExprAnyCast == 0 {
			if t := n.TypeName.Type(); t != nil && t.Kind() != Int {
				ctx.mode &^= mIntConstExpr
			}
			ctx.mode |= mIntConstExprFloat
		}
		n.CastExpression.check(ctx)
		ctx.pop()
	default:
		panic("internal error") //TODOOK
	}
}

func (n *PostfixExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case PostfixExpressionPrimary: // PrimaryExpression
		n.PrimaryExpression.check(ctx)
	case PostfixExpressionIndex: // PostfixExpression '[' Expression ']'
		n.PostfixExpression.check(ctx)
		n.Expression.check(ctx)
	case PostfixExpressionCall: // PostfixExpression '(' ArgumentExpressionList ')'
		n.PostfixExpression.check(ctx)
		n.ArgumentExpressionList.check(ctx)
	case PostfixExpressionSelect: // PostfixExpression '.' IDENTIFIER
		n.PostfixExpression.check(ctx)
		//TODO
	case PostfixExpressionPSelect: // PostfixExpression "->" IDENTIFIER
		n.PostfixExpression.check(ctx)
		//TODO
	case PostfixExpressionInc: // PostfixExpression "++"
		n.PostfixExpression.check(ctx)
	case PostfixExpressionDec: // PostfixExpression "--"
		n.PostfixExpression.check(ctx)
	case PostfixExpressionComplit: // '(' TypeName ')' '{' InitializerList ',' '}'
		n.TypeName.check(ctx)
		n.InitializerList.check(ctx)
	case PostfixExpressionTypeCmp: // "__builtin_types_compatible_p" '(' TypeName ',' TypeName ')'
		n.TypeName.check(ctx)
		n.TypeName2.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ArgumentExpressionList) check(ctx *context) {
	for ; n != nil; n = n.ArgumentExpressionList {
		n.AssignmentExpression.check(ctx)
	}
}

func (n *Expression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case ExpressionAssign: // AssignmentExpression
		n.AssignmentExpression.check(ctx)
	case ExpressionComma: // Expression ',' AssignmentExpression
		n.Expression.check(ctx)
		n.AssignmentExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *PrimaryExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case PrimaryExpressionIdent: // IDENTIFIER
		ctx.not(n, mIntConstExpr)
		//TODO
	case PrimaryExpressionInt: // INTCONST
		//TODO
	case PrimaryExpressionFloat: // FLOATCONST
		if ctx.mode&mIntConstExpr != 0 && ctx.mode&mIntConstExprFloat == 0 {
			ctx.errNode(n, "invalid integer constant expression")
		}
		//TODO
	case PrimaryExpressionEnum: // ENUMCONST
		//TODO
	case PrimaryExpressionChar: // CHARCONST
		//TODO
	case PrimaryExpressionLChar: // LONGCHARCONST
		//TODO
	case PrimaryExpressionString: // STRINGLITERAL
		ctx.not(n, mIntConstExpr)
		//TODO
	case PrimaryExpressionLString: // LONGSTRINGLITERAL
		ctx.not(n, mIntConstExpr)
		//TODO
	case PrimaryExpressionExpr: // '(' Expression ')'
		n.Expression.check(ctx)
	case PrimaryExpressionStmt: // '(' CompoundStatement ')'
		ctx.not(n, mIntConstExpr)
		n.CompoundStatement.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ConditionalExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case ConditionalExpressionLOr: // LogicalOrExpression
		n.LogicalOrExpression.check(ctx)
	case ConditionalExpressionCond: // LogicalOrExpression '?' Expression ':' ConditionalExpression
		n.LogicalOrExpression.check(ctx)
		n.Expression.check(ctx)
		n.ConditionalExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *LogicalOrExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case LogicalOrExpressionLAnd: // LogicalAndExpression
		n.LogicalAndExpression.check(ctx)
	case LogicalOrExpressionLOr: // LogicalOrExpression "||" LogicalAndExpression
		n.LogicalOrExpression.check(ctx)
		n.LogicalAndExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *LogicalAndExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case LogicalAndExpressionOr: // InclusiveOrExpression
		n.InclusiveOrExpression.check(ctx)
	case LogicalAndExpressionLAnd: // LogicalAndExpression "&&" InclusiveOrExpression
		n.LogicalAndExpression.check(ctx)
		n.InclusiveOrExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *InclusiveOrExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case InclusiveOrExpressionXor: // ExclusiveOrExpression
		n.ExclusiveOrExpression.check(ctx)
	case InclusiveOrExpressionOr: // InclusiveOrExpression '|' ExclusiveOrExpression
		n.InclusiveOrExpression.check(ctx)
		n.ExclusiveOrExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ExclusiveOrExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case ExclusiveOrExpressionAnd: // AndExpression
		n.AndExpression.check(ctx)
	case ExclusiveOrExpressionXor: // ExclusiveOrExpression '^' AndExpression
		n.ExclusiveOrExpression.check(ctx)
		n.AndExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *AndExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case AndExpressionEq: // EqualityExpression
		n.EqualityExpression.check(ctx)
	case AndExpressionAnd: // AndExpression '&' EqualityExpression
		n.AndExpression.check(ctx)
		n.EqualityExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *EqualityExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case EqualityExpressionRel: // RelationalExpression
		n.RelationalExpression.check(ctx)
	case EqualityExpressionEq: // EqualityExpression "==" RelationalExpression
		n.EqualityExpression.check(ctx)
		n.RelationalExpression.check(ctx)
	case EqualityExpressionNeq: // EqualityExpression "!=" RelationalExpression
		n.EqualityExpression.check(ctx)
		n.RelationalExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *RelationalExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case RelationalExpressionShift: // ShiftExpression
		n.ShiftExpression.check(ctx)
	case RelationalExpressionLt: // RelationalExpression '<' ShiftExpression
		n.RelationalExpression.check(ctx)
		n.ShiftExpression.check(ctx)
	case RelationalExpressionGt: // RelationalExpression '>' ShiftExpression
		n.RelationalExpression.check(ctx)
		n.ShiftExpression.check(ctx)
	case RelationalExpressionLeq: // RelationalExpression "<=" ShiftExpression
		n.RelationalExpression.check(ctx)
		n.ShiftExpression.check(ctx)
	case RelationalExpressionGeq: // RelationalExpression ">=" ShiftExpression
		n.RelationalExpression.check(ctx)
		n.ShiftExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ShiftExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case ShiftExpressionAdd: // AdditiveExpression
		n.AdditiveExpression.check(ctx)
	case ShiftExpressionLsh: // ShiftExpression "<<" AdditiveExpression
		n.ShiftExpression.check(ctx)
		n.AdditiveExpression.check(ctx)
	case ShiftExpressionRsh: // ShiftExpression ">>" AdditiveExpression
		n.ShiftExpression.check(ctx)
		n.AdditiveExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *AdditiveExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case AdditiveExpressionMul: // MultiplicativeExpression
		n.MultiplicativeExpression.check(ctx)
	case AdditiveExpressionAdd: // AdditiveExpression '+' MultiplicativeExpression
		n.AdditiveExpression.check(ctx)
		n.MultiplicativeExpression.check(ctx)
	case AdditiveExpressionSub: // AdditiveExpression '-' MultiplicativeExpression
		n.AdditiveExpression.check(ctx)
		n.MultiplicativeExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *MultiplicativeExpression) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case MultiplicativeExpressionCast: // CastExpression
		n.CastExpression.check(ctx)
	case MultiplicativeExpressionMul: // MultiplicativeExpression '*' CastExpression
		n.MultiplicativeExpression.check(ctx)
		n.CastExpression.check(ctx)
	case MultiplicativeExpressionDiv: // MultiplicativeExpression '/' CastExpression
		n.MultiplicativeExpression.check(ctx)
		n.CastExpression.check(ctx)
	case MultiplicativeExpressionMod: // MultiplicativeExpression '%' CastExpression
		n.MultiplicativeExpression.check(ctx)
		n.CastExpression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *Declarator) check(ctx *context, typ Type) Type {
	if n == nil {
		return noType
	}

	if n.Pointer != nil {
		typ = &pointerType{
			typ:            typ,
			typeQualifiers: n.Pointer.check(ctx),
			align:          byte(ctx.cfg.ABI.alignKind(ctx, n.Pointer, Ptr)),
			fieldAlign:     byte(ctx.cfg.ABI.fieldAlignKind(ctx, n.Pointer, Ptr)),
		}
	}
	n.typ = n.DirectDeclarator.check(ctx, typ)
	n.AttributeSpecifierList.check(ctx)
	return n.typ
}

func (n *DirectDeclarator) check(ctx *context, typ Type) Type {
	if n == nil {
		return noType
	}

	switch n.Case {
	case DirectDeclaratorIdent: // IDENTIFIER Asm
		n.Asm.check(ctx)
		return typ
	case DirectDeclaratorDecl: // '(' AttributeSpecifierList Declarator ')'
		n.AttributeSpecifierList.check(ctx)
		return n.Declarator.check(ctx, typ)
	case DirectDeclaratorArr: // DirectDeclarator '[' TypeQualifiers AssignmentExpression ']'
		n.DirectDeclarator.check(ctx, typ)
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		n.AssignmentExpression.check(ctx)
		//TODO type
	case DirectDeclaratorStaticArr: // DirectDeclarator '[' "static" TypeQualifiers AssignmentExpression ']'
		n.DirectDeclarator.check(ctx, typ)
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		n.AssignmentExpression.check(ctx)
		//TODO type
	case DirectDeclaratorArrStatic: // DirectDeclarator '[' TypeQualifiers "static" AssignmentExpression ']'
		n.DirectDeclarator.check(ctx, typ)
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		n.AssignmentExpression.check(ctx)
		//TODO type
	case DirectDeclaratorStar: // DirectDeclarator '[' TypeQualifiers '*' ']'
		n.DirectDeclarator.check(ctx, typ)
		n.TypeQualifiers.check(ctx, &n.typeQualifiers)
		//TODO type
	case DirectDeclaratorFuncParam: // DirectDeclarator '(' ParameterTypeList ')'
		n.DirectDeclarator.check(ctx, typ)
		n.ParameterTypeList.check(ctx)
		//TODO type
	case DirectDeclaratorFuncIdent: // DirectDeclarator '(' IdentifierList ')'
		n.DirectDeclarator.check(ctx, typ)
		n.IdentifierList.check(ctx)
		//TODO type
	default:
		panic("internal error") //TODOOK
	}
	return noType //TODO-
}

func (n *IdentifierList) check(ctx *context) {
	for ; n != nil; n = n.IdentifierList {
		//TODO
	}
}

func (n *Asm) check(ctx *context) {
	if n == nil {
		return
	}

	n.AsmQualifierList.check(ctx)
	n.AsmArgList.check(ctx)
}

func (n *AsmArgList) check(ctx *context) {
	for ; n != nil; n = n.AsmArgList {
		n.AsmExpressionList.check(ctx)
	}
}

func (n *AsmExpressionList) check(ctx *context) {
	for ; n != nil; n = n.AsmExpressionList {
		n.AsmIndex.check(ctx)
		n.AssignmentExpression.check(ctx)
	}
}

func (n *AsmIndex) check(ctx *context) {
	if n == nil {
		return
	}

	n.Expression.check(ctx)
}

func (n *AsmQualifierList) check(ctx *context) {
	for ; n != nil; n = n.AsmQualifierList {
		n.AsmQualifier.check(ctx)
	}
}

func (n *AsmQualifier) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case AsmQualifierVolatile: // "volatile"
		//TODO
	case AsmQualifierInline: // "inline"
		//TODO
	case AsmQualifierGoto: // "goto"
		//TODO
	default:
		panic("internal error") //TODOOK
	}
}

func (n *AttributeSpecifierList) check(ctx *context) {
	for ; n != nil; n = n.AttributeSpecifierList {
		n.AttributeSpecifier.check(ctx)
	}
}

func (n *AttributeSpecifier) check(ctx *context) {
	if n == nil {
		return
	}

	n.AttributeValueList.check(ctx)
}

func (n *AttributeValueList) check(ctx *context) {
	for ; n != nil; n = n.AttributeValueList {
		n.AttributeValue.check(ctx)
	}
}

func (n *AttributeValue) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case AttributeValueIdent: // IDENTIFIER
		//TODO
	case AttributeValueExpr: // IDENTIFIER '(' ExpressionList ')'
		n.ExpressionList.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ExpressionList) check(ctx *context) {
	for ; n != nil; n = n.ExpressionList {
		n.AssignmentExpression.check(ctx)
	}
}

func (n *DeclarationSpecifiers) check(ctx *context) Type {
	n0 := n
	typ := &baseType{}
	for ; n != nil; n = n.DeclarationSpecifiers {
		switch n.Case {
		case DeclarationSpecifiersStorage: // StorageClassSpecifier DeclarationSpecifiers
			n.StorageClassSpecifier.check(ctx, typ)
		case DeclarationSpecifiersTypeSpec: // TypeSpecifier DeclarationSpecifiers
			n.TypeSpecifier.check(ctx)
		case DeclarationSpecifiersTypeQual: // TypeQualifier DeclarationSpecifiers
			n.TypeQualifier.check(ctx, typ)
		case DeclarationSpecifiersFunc: // FunctionSpecifier DeclarationSpecifiers
			n.FunctionSpecifier.check(ctx, typ)
		case DeclarationSpecifiersAlignSpec: // AlignmentSpecifier DeclarationSpecifiers
			n.AlignmentSpecifier.check(ctx)
		case DeclarationSpecifiersAttribute: // AttributeSpecifier DeclarationSpecifiers
			n.AttributeSpecifier.check(ctx)
		default:
			panic("internal error") //TODOOK
		}
	}
	return typ.check(ctx, n0)
}

func (n *AlignmentSpecifier) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case AlignmentSpecifierAlignasType: // "_Alignas" '(' TypeName ')'
		n.TypeName.check(ctx)
	case AlignmentSpecifierAlignasExpr: // "_Alignas" '(' ConstantExpression ')'
		n.ConstantExpression.check(ctx, ctx.mode|mIntConstExpr)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *FunctionSpecifier) check(ctx *context, typ *baseType) {
	if n == nil {
		return
	}

	switch n.Case {
	case FunctionSpecifierInline: // "inline"
		typ.flags |= btfInline
	case FunctionSpecifierNoreturn: // "_Noreturn"
		typ.flags |= btfNoReturn
	default:
		panic("internal error") //TODOOK
	}
}

func (n *StorageClassSpecifier) check(ctx *context, typ *baseType) {
	if n == nil {
		return
	}

	switch n.Case {
	case StorageClassSpecifierTypedef: // "typedef"
		typ.flags |= btfTypedef
	case StorageClassSpecifierExtern: // "extern"
		typ.flags |= btfExtern
	case StorageClassSpecifierStatic: // "static"
		typ.flags |= btfStatic
	case StorageClassSpecifierAuto: // "auto"
		typ.flags |= btfAuto
	case StorageClassSpecifierRegister: // "register"
		typ.flags |= btfRegister
	case StorageClassSpecifierThreadLocal: // "_Thread_local"
		typ.flags |= btfThreadLocal
	default:
		panic("internal error") //TODOOK
	}
	c := bits.OnesCount(uint(typ.flags & (btfTypedef | btfExtern | btfStatic | btfAuto | btfRegister | btfThreadLocal)))
	if c == 1 {
		return
	}

	// [2], 6.7.1, 2
	if c == 2 && typ.flags&btfThreadLocal != 0 {
		if typ.flags&(btfStatic|btfExtern) != 0 {
			return
		}
	}

	ctx.errNode(n, "at most, one storage-class specifier may be given in the declaration specifiers in a declaration")
}

func (n *FunctionDefinition) check(ctx *context) {
	if n == nil {
		return
	}

	typ := n.DeclarationSpecifiers.check(ctx)
	n.Declarator.check(ctx, typ)
	n.DeclarationList.check(ctx)
	n.CompoundStatement.check(ctx)
}

func (n *CompoundStatement) check(ctx *context) {
	if n == nil {
		return
	}

	n.BlockItemList.check(ctx)
}

func (n *BlockItemList) check(ctx *context) {
	for ; n != nil; n = n.BlockItemList {
		n.BlockItem.check(ctx)
	}
}

func (n *BlockItem) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case BlockItemDecl: // Declaration
		n.Declaration.check(ctx)
	case BlockItemStmt: // Statement
		n.Statement.check(ctx)
	case BlockItemLabel: // LabelDeclaration
		n.LabelDeclaration.check(ctx)
	case BlockItemFuncDef: // DeclarationSpecifiers Declarator CompoundStatement
		typ := n.DeclarationSpecifiers.check(ctx)
		n.Declarator.check(ctx, typ)
		n.CompoundStatement.check(ctx)
	case BlockItemPragma: // PragmaSTDC
		n.PragmaSTDC.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *LabelDeclaration) check(ctx *context) {
	if n == nil {
		return
	}

	n.IdentifierList.check(ctx)
}

func (n *Statement) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case StatementLabeled: // LabeledStatement
		n.LabeledStatement.check(ctx)
	case StatementCompound: // CompoundStatement
		n.CompoundStatement.check(ctx)
	case StatementExpr: // ExpressionStatement
		n.ExpressionStatement.check(ctx)
	case StatementSelection: // SelectionStatement
		n.SelectionStatement.check(ctx)
	case StatementIteration: // IterationStatement
		n.IterationStatement.check(ctx)
	case StatementJump: // JumpStatement
		n.JumpStatement.check(ctx)
	case StatementAsm: // AsmStatement
		n.AsmStatement.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *JumpStatement) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case JumpStatementGoto: // "goto" IDENTIFIER ';'
		//TODO
	case JumpStatementGotoExpr: // "goto" '*' Expression ';'
		n.Expression.check(ctx)
	case JumpStatementContinue: // "continue" ';'
		//TODO
	case JumpStatementBreak: // "break" ';'
		//TODO
	case JumpStatementReturn: // "return" Expression ';'
		n.Expression.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *IterationStatement) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case IterationStatementWhile: // "while" '(' Expression ')' Statement
		n.Expression.check(ctx)
		n.Statement.check(ctx)
	case IterationStatementDo: // "do" Statement "while" '(' Expression ')' ';'
		n.Statement.check(ctx)
		n.Expression.check(ctx)
	case IterationStatementFor: // "for" '(' Expression ';' Expression ';' Expression ')' Statement
		n.Expression.check(ctx)
		n.Expression2.check(ctx)
		n.Expression3.check(ctx)
		n.Statement.check(ctx)
	case IterationStatementForDecl: // "for" '(' Declaration Expression ';' Expression ')' Statement
		n.Declaration.check(ctx)
		n.Expression.check(ctx)
		n.Expression2.check(ctx)
		n.Statement.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *SelectionStatement) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case SelectionStatementIf: // "if" '(' Expression ')' Statement
		n.Expression.check(ctx)
		n.Statement.check(ctx)
	case SelectionStatementIfElse: // "if" '(' Expression ')' Statement "else" Statement
		n.Expression.check(ctx)
		n.Statement.check(ctx)
		n.Statement2.check(ctx)
	case SelectionStatementSwitch: // "switch" '(' Expression ')' Statement
		n.Expression.check(ctx)
		n.Statement.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *ExpressionStatement) check(ctx *context) {
	if n == nil {
		return
	}

	n.Expression.check(ctx)
	n.AttributeSpecifierList.check(ctx)
}

func (n *LabeledStatement) check(ctx *context) {
	if n == nil {
		return
	}

	switch n.Case {
	case LabeledStatementLabel: // IDENTIFIER ':' AttributeSpecifierList Statement
		n.AttributeSpecifierList.check(ctx)
		n.Statement.check(ctx)
	case LabeledStatementCaseLabel: // "case" ConstantExpression ':' Statement
		n.ConstantExpression.check(ctx, ctx.mode|mIntConstExpr)
		n.Statement.check(ctx)
	case LabeledStatementRange: // "case" ConstantExpression "..." ConstantExpression ':' Statement
		n.ConstantExpression.check(ctx, ctx.mode|mIntConstExpr)
		n.ConstantExpression2.check(ctx, ctx.mode|mIntConstExpr)
		n.Statement.check(ctx)
	case LabeledStatementDefault: // "default" ':' Statement
		n.Statement.check(ctx)
	default:
		panic("internal error") //TODOOK
	}
}

func (n *DeclarationList) check(ctx *context) {
	for ; n != nil; n = n.DeclarationList {
		n.Declaration.check(ctx)
	}
}
