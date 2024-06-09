package main

import (
	"strings"
)

// Convert parsed LexTree into string representation.
// Use passed options.
func PrintLexTreeNodes(nodes LexTreeNodeList, indentation uint, options *FormattingOptions) string {
	builder := strings.Builder{}

	for i := uint(0); i < indentation; i++ {
		builder.WriteString(options.indentation_sequence)
	}

	current_line_width := CountIndentationsSize(indentation, options)

	PrintLexTreeNodes_r(nodes, options, &builder, indentation, &current_line_width)

	builder.WriteString(options.line_end_sequence)

	return builder.String()
}

// Main formatting routine.
// Since it operates with simple tree-like structure and not proper syntax tree, parsed based on proper language grammatic,
// it uses some heuristics to detect common patterns.
// Such heuristics may occasionally fail.
func PrintLexTreeNodes_r(
	nodes LexTreeNodeList,
	options *FormattingOptions,
	out *strings.Builder,
	indentation uint,
	current_line_width *uint) {

	text_with_further_split := ""
	further_split_current_line_width := *current_line_width

	{
		builder := strings.Builder{}

		for i, node := range nodes {

			if node.sub_elements == nil {

				if i > 0 && WhitespaceIsNeeded(&nodes[i-1].lexem, &node.lexem) {
					builder.WriteString(" ")
					further_split_current_line_width++
				}

				builder.WriteString(node.lexem.text)
				further_split_current_line_width += uint(len(node.lexem.text))

			} else {

				builder.WriteString(node.lexem.text)
				further_split_current_line_width += uint(len(node.lexem.text))

				if len(node.sub_elements) > 0 {
					builder.WriteString(" ")
					further_split_current_line_width++
				}

				PrintLexTreeNodes_r(
					node.sub_elements, options, &builder, indentation, &further_split_current_line_width)

				if len(node.sub_elements) > 0 {
					builder.WriteString(" ")
					further_split_current_line_width++
				}

				builder.WriteString(node.trailing_lexem.text)
				further_split_current_line_width += uint(len(node.trailing_lexem.text))
			}
		}

		text_with_further_split = builder.String()
	}

	text_with_current_split := ""
	current_split_current_line_width := *current_line_width

	// Perform split only if further split fails to achieve target line width.
	if len(nodes) > 1 {
		// Recursively split and print this list, adding newlines in split points.
		builder := strings.Builder{}

		// Search for the most important lexem type to use it as splitter.
		// Ignore last node, because splitting at last node has no sense.
		max_priority := 0
		for _, node := range nodes[:len(nodes)-1] {
			priority := GetLineSplitLexemPriority(&node.lexem)
			if priority > max_priority {
				max_priority = priority
			}
		}

		// Split this lexems list into parts, using maximum priority lexem type.
		// Add newline after each part.
		last_i := 0
		next_indentation := indentation
		for i := 0; i < len(nodes)-1; i++ {

			if GetLineSplitLexemPriority(&nodes[i].lexem) == max_priority {

				PrintLexTreeNodes_r(
					nodes[last_i:i+1], options, &builder, next_indentation, &current_split_current_line_width)
				last_i = i + 1

				next_indentation = indentation + 1

				builder.WriteString(options.line_end_sequence)
				for i := uint(0); i < next_indentation; i++ {
					builder.WriteString(options.indentation_sequence)
				}
				current_split_current_line_width = CountIndentationsSize(next_indentation, options)
			}
		}

		// Process last segment specially.
		PrintLexTreeNodes_r(
			nodes[last_i:], options, &builder, next_indentation, &current_split_current_line_width)

		text_with_current_split = builder.String()
	}

	if len(text_with_current_split) == 0 {
		// Can't split at this level - use recursive split result.
		out.WriteString(text_with_further_split)
		*current_line_width = further_split_current_line_width
	} else if CountNewlines(text_with_current_split) <= CountNewlines(text_with_further_split) {
		// Split at this level gives less or equal lines compared to splits at further levels.
		out.WriteString(text_with_current_split)
		*current_line_width = current_split_current_line_width
	} else {

		// Count max line width for further split.
		max_line_width := uint(0)
		{
			w := *current_line_width

			for _, c := range text_with_further_split {
				if c == '\n' { // TODO - use newline sequence from options.
					max_line_width = max(max_line_width, w)
					w = 0
				} else if c == '\t' {
					w += options.tab_size
				} else {
					w++
				}
			}

			max_line_width = max(max_line_width, w)
		}

		// Split result at this level if further result exceeedes the limit.
		if max_line_width > options.max_line_width {
			out.WriteString(text_with_current_split)
			*current_line_width = current_split_current_line_width
		} else {
			out.WriteString(text_with_further_split)
			*current_line_width = further_split_current_line_width
		}
	}
}

func HasNaturalNewlines(nodes LexTreeNodeList) bool {

	for _, node := range nodes {
		if node.lexem.t == LexemTypeLineComment || node.lexem.t == LexemTypeSemicolon {
			return true
		}

		if HasNaturalNewlines(node.sub_elements) {
			return true
		}
	}

	return false
}

func CountNewlines(s string) uint {
	// TODO - use newline sequence from options.
	count := uint(0)
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}

	return count
}

func CountIndentationsSize(indentation uint, options *FormattingOptions) uint {
	count := uint(0)
	for _, c := range options.indentation_sequence {
		if c == '\t' {
			count += options.tab_size
		} else {
			count++
		}
	}

	return count * indentation
}

// More priority - more likely to split.
func GetLineSplitLexemPriority(l *Lexem) int {
	switch l.t {

	case LexemTypeLineComment:
		return 200

	case LexemTypeSemicolon:
		return 100

	case LexemTypeComma:
		return 99

	case LexemTypeAssignment:
		return 90

	case LexemTypeColon,
		LexemTypeQuestion:
		return 82

	// Use here binary operator priorities.

	case LexemTypeDisjunction:
		return 80
	case LexemTypeConjunction:
		return 79

	case LexemTypeOr:
		return 78
	case LexemTypeXor:
		return 77
	case LexemTypeAnd:
		return 76

	case LexemTypeCompareEqual:
	case LexemTypeCompareNotEqual:
		return 75

	case LexemTypeCompareLess:
	case LexemTypeCompareLessOrEqual:
	case LexemTypeCompareGreater:
	case LexemTypeCompareGreaterOrEqual:
		return 74

	case LexemTypeCompareOrder:
		return 73

	case LexemTypeShiftLeft,
		LexemTypeShiftRight:
		return 72

	case LexemTypePlus:
		return 71
	case LexemTypeMinus: // TODO - what about unary minus?
		return 70

	case LexemTypeStar,
		LexemTypeSlash,
		LexemTypePercent:
		return 69

	case LexemTypeDot:
		return 40

	case LexemTypeBraceLeft:
		return 30

	case LexemTypeBracketLeft,
		LexemTypeSquareBracketLeft,
		LexemTypeTemplateBracketLeft,
		LexemTypeMacroBracketLeft:
		return 20

	case LexemTypeIdentifier:
		return 10

		// TODO - add other lexems
	}

	return 1
}
