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

	PrintLexTreeNodes_r(nodes, options, &builder, indentation, current_line_width)

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
	current_line_width uint) {

	text_with_further_split := ""

	{
		builder := strings.Builder{}

		w := current_line_width

		for i, node := range nodes {

			if node.sub_elements == nil {

				if i > 0 && WhitespaceIsNeeded(&nodes[i-1].lexem, &node.lexem) {
					builder.WriteString(" ")
					w++
				}

				builder.WriteString(node.lexem.text)
				w += uint(len(node.lexem.text))

			} else {

				builder.WriteString(node.lexem.text)
				w += uint(len(node.lexem.text))

				if len(node.sub_elements) > 0 {
					builder.WriteString(" ")
					w++
				}

				PrintLexTreeNodes_r(node.sub_elements, options, &builder, indentation, w) // TODO - we need to update current_line_width recursively

				if len(node.sub_elements) > 0 {
					builder.WriteString(" ")
					w++
				}

				builder.WriteString(node.trailing_lexem.text)
				w += uint(len(node.trailing_lexem.text))
			}
		}

		text_with_further_split = builder.String()
	}

	// Count max line width for further split.
	max_line_width := uint(0)
	{
		w := current_line_width

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

	text_with_current_split := ""

	// Perform split only if further split fails to achieve target line width.
	if max_line_width > options.max_line_width &&
		len(nodes) > 1 {
		// Recursively split and print this list, adding newlines in split points.
		builder := strings.Builder{}

		w := current_line_width

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

				PrintLexTreeNodes_r(nodes[last_i:i+1], options, &builder, next_indentation, w)
				last_i = i + 1

				next_indentation = indentation + 1

				builder.WriteString(options.line_end_sequence)
				for i := uint(0); i < next_indentation; i++ {
					builder.WriteString(options.indentation_sequence)
				}
				w = CountIndentationsSize(next_indentation, options)
			}
		}

		// Process last segment specially.
		PrintLexTreeNodes_r(nodes[last_i:], options, &builder, next_indentation, w)

		text_with_current_split = builder.String()
	}

	if len(text_with_current_split) == 0 {
		// Can't split at this level - use recursive split result.
		out.WriteString(text_with_further_split)
	} else if CountNewlines(text_with_current_split) <= CountNewlines(text_with_further_split) {
		// Split at this level gives less or equal lines compared to splits at further levels.
		out.WriteString(text_with_current_split)
	} else {

		// Split result at this level if further level can't achieve target number of lines.
		if max_line_width > options.max_line_width {
			out.WriteString(text_with_current_split)
		} else {
			out.WriteString(text_with_further_split)
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

	case LexemTypeShiftLeft:
	case LexemTypeShiftRight:
		return 72

	case LexemTypePlus:
		return 71
	case LexemTypeMinus: // TODO - what about unary minus?
		return 70

	case LexemTypeStar:
	case LexemTypeSlash:
	case LexemTypePercent:
		return 69

	case LexemTypeDot:
		return 40

	case LexemTypeBraceLeft:
		return 30

	case LexemTypeBracketLeft:
	case LexemTypeSquareBracketLeft:
	case LexemTypeTemplateBracketLeft:
	case LexemTypeMacroBracketLeft:
		return 20

	case LexemTypeIdentifier:
		return 10

		// TODO - add other lexems
	}

	return 1
}
