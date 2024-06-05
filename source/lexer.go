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

	result = append(result, Lexem{text: "TODO", t: LexemTypeEllipsis})

	return result
}
