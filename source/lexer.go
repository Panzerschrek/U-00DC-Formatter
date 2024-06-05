package main

type Lexem struct {
	t    LexemType
	text string
}

type LexemType byte

const (
	LexemTypeNone LexemType = iota

	LexemTypeComment

	LexemTypeIdentifier
	LexemTypeMacroIdentifier
	LexemTypeMacroUniqueIdentifier
	LexemTypeString
	LexemTypeNumber

	LexemTypeLiteralSuffix // For strings, numbers

	LexemTypeBracketLeft  // (
	LexemTypeBracketRight // )

	LexemTypeSquareBracketLeft  // [
	LexemTypeSquareBracketRight // ]

	LexemTypeBraceLeft  // {
	LexemTypeBraceRight // }

	LexemTypeTemplateBracketLeft  // </
	LexemTypeTemplateBracketRight // />

	LexemTypeMacroBracketLeft  // <?
	LexemTypeMacroBracketRight // ?>

	LexemTypeScope // ::

	LexemTypeComma     // ,
	LexemTypeDot       // .
	LexemTypeColon     // :
	LexemTypeSemicolon // ;
	LexemTypeQuestion  // ?

	LexemTypeAssignment // =
	LexemTypePlus       // +
	LexemTypeMinus      // -
	LexemTypeStar       // *
	LexemTypeSlash      // /
	LexemTypePercent    // %

	LexemTypeAnd   // &
	LexemTypeOr    // |
	LexemTypeXor   // ^
	LexemTypeTilda // ~
	LexemTypeNot   // !

	LexemTypeApostrophe // '
	LexemTypeAt         // @

	LexemTypeIncrement // ++
	LexemTypeDecrement // --

	LexemTypeCompareLess           // <
	LexemTypeCompareGreater        // >
	LexemTypeCompareEqual          // ==
	LexemTypeCompareNotEqual       // !=
	LexemTypeCompareLessOrEqual    // <=
	LexemTypeCompareGreaterOrEqual // >=
	LexemTypeCompareOrder          // <=>

	LexemTypeConjunction // &&
	LexemTypeDisjunction // ||

	LexemTypeAssignAdd // +=
	LexemTypeAssignSub // -=
	LexemTypeAssignMul // *=
	LexemTypeAssignDiv // /=
	LexemTypeAssignRem // %=
	LexemTypeAssignAnd // &=
	LexemTypeAssignOr  // |=
	LexemTypeAssignXor // ^=

	LexemTypeShiftLeft  // <<
	LexemTypeShiftRight // >>

	LexemTypeAssignShiftLeft  // <<=
	LexemTypeAssignShiftRight // >>=

	LexemTypeRightArrow // ->

	LexemTypePointerTypeMark    // $
	LexemTypeReferenceToPointer // $<
	LexemTypePointerToReference // $>

	LexemTypeEllipsis // ...

	// Special kind of lexems, that can be created only manually (and not parsed).
	LexemTypeCompletionIdentifier
	LexemTypeCompletionScope
	LexemTypeCompletionDot
	LexemTypeSignatureHelpBracketLeft
	LexemTypeSignatureHelpComma
)

func splitProgramIntoLexems(s string) []Lexem {
	result := make([]Lexem, 0)

	for len(s) > 0 {

		// TODO - process identifiers, numbers, comments, etc.
		
		// Process fixed lexems.
		if len(s) >= 3 { // Fixed lexems of length 3.
			lexem_type := TextToLexem3(s[0:3])
			if lexem_type != LexemTypeNone {
				result = append(result, Lexem{text: s[0:3], t: lexem_type})
				s = s[3:]
				goto end_fixed_lexems_search
			}
		}
		if len(s) >= 2 {  // Fixed lexems of length 2.
			lexem_type := TextToLexem3(s[0:2])
			if lexem_type != LexemTypeNone {
				result = append(result, Lexem{text: s[0:2], t: lexem_type})
				s = s[2:]
				goto end_fixed_lexems_search
			}
		}
		if len(s) >= 1 {  // Fixed lexems of length 1.
			lexem_type := TextToLexem3(s[0:1])
			if lexem_type != LexemTypeNone {
				result = append(result, Lexem{text: s[0:1], t: lexem_type})
				s = s[1:]
				goto end_fixed_lexems_search
			}
		}

		// None of the fixed lexems
		s = s[1:]

	end_fixed_lexems_search:
	}

	return result
}

func TextToLexem1(s string) LexemType {
	switch s {
	case "(":
		return LexemTypeBracketLeft
	case ")":
		return LexemTypeBracketRight
	case "[":
		return LexemTypeSquareBracketLeft
	case "]":
		return LexemTypeSquareBracketRight
	case "{":
		return LexemTypeBraceLeft
	case "}":
		return LexemTypeBraceRight

	case ",":
		return LexemTypeComma
	case ".":
		return LexemTypeDot
	case ":":
		return LexemTypeColon
	case ";":
		return LexemTypeSemicolon
	case "?":
		return LexemTypeQuestion

	case "=":
		return LexemTypeAssignment
	case "+":
		return LexemTypePlus
	case "-":
		return LexemTypeMinus
	case "*":
		return LexemTypeStar
	case "/":
		return LexemTypeSlash
	case "%":
		return LexemTypePercent

	case "<":
		return LexemTypeCompareLess
	case ">":
		return LexemTypeCompareGreater

	case "&":
		return LexemTypeAnd
	case "|":
		return LexemTypeOr
	case "^":
		return LexemTypeXor
	case "~":
		return LexemTypeTilda
	case "!":
		return LexemTypeNot

	case "'":
		return LexemTypeApostrophe
	case "@":
		return LexemTypeAt

	case "$":
		return LexemTypePointerTypeMark
	}

	return LexemTypeNone
}

func TextToLexem2(s string) LexemType {
	switch s {
	case "</":
		return LexemTypeTemplateBracketLeft
	case "/>":
		return LexemTypeTemplateBracketRight

	case "<?":
		return LexemTypeMacroBracketLeft
	case "?>":
		return LexemTypeMacroBracketRight

	case "::":
		return LexemTypeScope

	case "++":
		return LexemTypeIncrement
	case "--":
		return LexemTypeDecrement

	case "==":
		return LexemTypeCompareEqual
	case "!=":
		return LexemTypeCompareNotEqual
	case "<=":
		return LexemTypeCompareLessOrEqual
	case ">=":
		return LexemTypeCompareGreaterOrEqual

	case "&&":
		return LexemTypeConjunction
	case "||":
		return LexemTypeDisjunction

	case "+=":
		return LexemTypeAssignAdd
	case "-=":
		return LexemTypeAssignSub
	case "*=":
		return LexemTypeAssignMul
	case "/=":
		return LexemTypeAssignDiv
	case "%=":
		return LexemTypeAssignRem
	case "&=":
		return LexemTypeAssignAnd
	case "|=":
		return LexemTypeAssignOr
	case "^=":
		return LexemTypeAssignXor

	case "<<":
		return LexemTypeShiftLeft
	case ">>":
		return LexemTypeShiftRight

	case "->":
		return LexemTypeRightArrow

	case "$<":
		return LexemTypeReferenceToPointer
	case "$>":
		return LexemTypePointerToReference
	}

	return LexemTypeNone
}

func TextToLexem3(s string) LexemType {
	switch s {
	case "<=>":
		return LexemTypeCompareOrder
	case "<<=":
		return LexemTypeAssignShiftLeft
	case ">>=":
		return LexemTypeAssignShiftRight
	case "...":
		return LexemTypeEllipsis
	}

	return LexemTypeNone
}
