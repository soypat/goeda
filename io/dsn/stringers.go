// Code generated by "stringer -type=Token -output=stringers.go -trimprefix=Tok"; DO NOT EDIT.

package kicad

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokILLEGAL-0]
	_ = x[TokEOF-1]
	_ = x[TokLPAREN-2]
	_ = x[TokRPAREN-3]
	_ = x[TokIDENT-4]
	_ = x[TokINTEGER-5]
	_ = x[TokSTRING-6]
	_ = x[TokFLOAT-7]
}

const _Token_name = "ILLEGALEOFLPARENRPARENIDENTINTEGERSTRINGFLOAT"

var _Token_index = [...]uint8{0, 7, 10, 16, 22, 27, 34, 40, 45}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}