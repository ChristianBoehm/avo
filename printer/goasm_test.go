package printer_test

import (
	"testing"

	"github.com/mmcloughlin/avo/attr"
	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/printer"
	"github.com/mmcloughlin/avo/reg"
)

func TestBasic(t *testing.T) {
	ctx := build.NewContext()
	ctx.Function("add")
	ctx.SignatureExpr("func(x, y uint64) uint64")
	x := ctx.Load(ctx.Param("x"), reg.RAX)
	y := ctx.Load(ctx.Param("y"), reg.R9)
	ctx.ADDQ(x, y)
	ctx.Store(y, ctx.ReturnIndex(0))
	ctx.RET()

	AssertPrintsLines(t, ctx, printer.NewGoAsm, []string{
		"// Code generated by avo. DO NOT EDIT.",
		"",
		"// func add(x uint64, y uint64) uint64",
		"TEXT ·add(SB), $0-24",
		"\tMOVQ x(FP), AX",
		"\tMOVQ y+8(FP), R9",
		"\tADDQ AX, R9",
		"\tMOVQ R9, ret+16(FP)",
		"\tRET",
		"",
	})
}

func TestTextDecl(t *testing.T) {
	ctx := build.NewContext()

	ctx.Function("noargs")
	ctx.SignatureExpr("func()")
	ctx.AllocLocal(16)
	ctx.RET()

	ctx.Function("withargs")
	ctx.SignatureExpr("func(x, y uint64) uint64")
	ctx.RET()

	ctx.Function("withattr")
	ctx.SignatureExpr("func()")
	ctx.Attributes(attr.NOSPLIT | attr.TLSBSS)
	ctx.RET()

	AssertPrintsLines(t, ctx, printer.NewGoAsm, []string{
		"// Code generated by avo. DO NOT EDIT.",
		"",
		"// func noargs()",
		"TEXT ·noargs(SB), $16", // expect only the frame size
		"\tRET",
		"",
		"// func withargs(x uint64, y uint64) uint64",
		"TEXT ·withargs(SB), $0-24", // expect both frame size and argument size
		"\tRET",
		"",
		"// func withattr()",
		"TEXT ·withattr(SB), NOSPLIT|TLSBSS, $0", // expect to see attributes
		"\tRET",
		"",
	})
}

func TestConstraints(t *testing.T) {
	ctx := build.NewContext()
	ctx.ConstraintExpr("linux,386 darwin,!cgo")
	ctx.ConstraintExpr("!noasm")

	AssertPrintsLines(t, ctx, printer.NewGoAsm, []string{
		"// Code generated by avo. DO NOT EDIT.",
		"",
		"// +build linux,386 darwin,!cgo",
		"// +build !noasm",
		"",
	})
}
