package main

import (
	"strings"
	"unicode/utf8"
)

type Lexem struct {
	t    LexemType
	text string
}

type LexemType byte

const (
	LexemTypeNone LexemType = iota

	LexemTypeLineComment

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

	LexemTypeEndOfFile
)

func SplitProgramIntoLexems(s string) []Lexem {
	result := make([]Lexem, 0)

	for len(s) > 0 {

		c, c_size := utf8.DecodeRuneInString(s)

		// TODO - parse multiline comments

		if IsWhitespace(c) {

			s = s[c_size:] // Skip whitespaces.

		} else if c == '/' && len(s) > c_size && s[1] == '/' {

			// Line comment
			comment := Lexem{t: LexemTypeLineComment}
			s_before := s

			for len(s) > 0 {
				c, c_size := utf8.DecodeRuneInString(s)

				if IsNewline(c) {
					break
				}
				s = s[c_size:]
			}

			comment.text = string(s_before[:len(s_before)-len(s)])
			result = append(result, comment)

		} else if IsIdentifierStartChar(c) {

			result = append(result, ParseIdentifier(&s))

		} else if IsNumberStartChar(c) {

			result = append(result, ParseNumber(&s))

		} else if c == '"' {

			result = append(result, ParseString(&s))

		} else {
			// Process fixed lexems.

			if len(s) >= 3 { // Fixed lexems of length 3.
				lexem_type := TextToLexem3(s[0:3])
				if lexem_type != LexemTypeNone {
					result = append(result, Lexem{text: s[0:3], t: lexem_type})
					s = s[3:]
					continue
				}
			}
			if len(s) >= 2 { // Fixed lexems of length 2.
				lexem_type := TextToLexem2(s[0:2])
				if lexem_type != LexemTypeNone {
					result = append(result, Lexem{text: s[0:2], t: lexem_type})
					s = s[2:]
					continue
				}
			}
			if len(s) >= 1 { // Fixed lexems of length 1.
				lexem_type := TextToLexem1(s[0:1])
				if lexem_type != LexemTypeNone {
					result = append(result, Lexem{text: s[0:1], t: lexem_type})
					s = s[1:]
					continue
				}
			}

			// None of the fixed lexems.
			// TODO - generate error
			s = s[1:]
		}
	}

	result = append(result, Lexem{t: LexemTypeEndOfFile})

	return result
}

func IsWhitespace(c rune) bool {
	return c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t' || c == '\v' || c <= 0x1F || c == 0x7F
}

func IsNewline(c rune) bool {
	// See https://en.wikipedia.org/wiki/Newline#Unicode.
	return c == '\n' || // line feed
		c == '\r' || // carriage return
		c == '\f' || // form feed
		c == '\v' || // vertical tab
		c == 0x0085 || // Next line
		c == 0x2028 || // line separator
		c == 0x2029 // paragraph separator
}

func IsIdentifierStartChar(c rune) bool {
	// HACK - manually define allowed "letters".
	// TODO - use something, like symbol category from unicode.
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= 0x0400 && c <= 0x04FF) || // Cyrillic
		(c >= 0x0500 && c <= 0x0527) || // Extended cyrillic
		(c >= 0x00C0 && c <= 0x00D6) || // Additional latin symbols
		(c >= 0x00D8 && c <= 0x00F6) || // Additional latin symbols
		(c >= 0x00F8 && c <= 0x00FF) || // Additional latin symbols
		(c >= 0x0100 && c <= 0x017F) || // Extended latin part A
		(c >= 0x0180 && c <= 0x024F) // Extended latin part B
}

func IsIdentifierChar(c rune) bool {
	return IsIdentifierStartChar(c) || IsNumberStartChar(c) || c == '_'
}

func IsNumberStartChar(c rune) bool {
	return c >= '0' && c <= '9'
}

func ParseIdentifier(s *string) Lexem {

	s_initial := *s

	for len(*s) > 0 {
		c, c_size := utf8.DecodeRuneInString(*s)
		if !IsIdentifierChar(c) {
			break
		}

		*s = (*s)[c_size:]
	}

	return Lexem{t: LexemTypeIdentifier, text: string(s_initial[:len(s_initial)-len(*s)])}
}

func ParseNumber(s *string) Lexem {

	s_initial := *s

	var number_func = func(c rune) int {
		if c >= '0' && c <= '9' {
			return int(c) - '0'
		}

		return -1
	}

	base := 10

	if strings.HasPrefix(*s, "0b") {

		base = 2

		number_func = func(c rune) int {
			if c >= '0' && c <= '1' {
				return int(c) - '0'
			}
			return -1
		}

		*s = (*s)[2:]

	} else if strings.HasPrefix(*s, "0o") {

		base = 8

		number_func = func(c rune) int {
			if c >= '0' && c <= '7' {
				return int(c) - '0'
			}
			return -1
		}

		*s = (*s)[2:]

	} else if strings.HasPrefix(*s, "0x") {

		base = 16

		number_func = func(c rune) int {
			if c >= '0' && c <= '9' {
				return int(c) - '0'
			}
			if c >= 'a' && c <= 'f' {
				return int(c) - 'a' + 10
			}
			if c >= 'A' && c <= 'F' {
				return int(c) - 'a' + 10
			}
			return -1
		}

		*s = (*s)[2:]
	}

	// Parse integer part.
	for len(*s) > 0 {
		c, c_size := utf8.DecodeRuneInString(*s)
		if number_func(c) < 0 {
			break
		}
		*s = (*s)[c_size:]
	}

	if strings.HasPrefix(*s, ".") {
		// Parse fractional part.
		*s = (*s)[1:]

		for len(*s) > 0 {
			c, c_size := utf8.DecodeRuneInString(*s)
			if number_func(c) < 0 {
				break
			}
			*s = (*s)[c_size:]
		}
	}

	if base == 10 && strings.HasPrefix(*s, "e") {
		// Parse exponent
		*s = (*s)[1:]

		// Optional negative sign.
		if strings.HasPrefix(*s, "-") {
			*s = (*s)[1:]
		}

		// Exponent digits
		for len(*s) > 0 {
			c, c_size := utf8.DecodeRuneInString(*s)
			if number_func(c) < 0 {
				break
			}
			*s = (*s)[c_size:]
		}
	}

	c, _ := utf8.DecodeRuneInString(*s)
	if IsIdentifierStartChar(c) {
		// Type suffix.
		ParseIdentifier(s)
	}

	return Lexem{t: LexemTypeNumber, text: string(s_initial[:len(s_initial)-len(*s)])}
}

func ParseString(s *string) Lexem {

	s_initial := *s

	*s = (*s)[1:] // Skip initial "

	for len(*s) > 0 {
		c, c_size := utf8.DecodeRuneInString(*s)
		if c == '\\' {
			// TODO - check if escape sequence is correct.
			*s = (*s)[2:]
			continue
		} else if c == '"' {
			*s = (*s)[1:]
			break
		} else {
			*s = (*s)[c_size:]
		}
	}

	return Lexem{t: LexemTypeString, text: string(s_initial[:len(s_initial)-len(*s)])}
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
