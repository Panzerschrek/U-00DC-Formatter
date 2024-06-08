package main

import (
	"strings"
)

// Convert parsed LexTree into string representation.
// Use passed options.
func PrintLexTreeNodes(nodes LexTreeNodeList, depth uint, options *FormattingOptions) string {
	var prev_was_newline bool = true
	builder := strings.Builder{}
	PrintLexTreeNodes_r(nodes, options, &builder, depth, &prev_was_newline, false)

	if !prev_was_newline {
		builder.WriteString(options.line_end_sequence)
	}

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
	depth uint,
	prev_was_newline *bool,
	force_single_line bool) {

	if len(nodes) > 1 &&
		!force_single_line &&
		!CanWriteInSingleLine(nodes, depth, options) {
		// Recursively split and print this list, adding newlines in split points.

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
		subrange_index := 0
		for i := 0; i < len(nodes)-1; i++ {

			if GetLineSplitLexemPriority(&nodes[i].lexem) == max_priority {

				if !*prev_was_newline {
					out.WriteString(options.line_end_sequence)
					*prev_was_newline = true
				}

				subrange_depth := depth
				if subrange_index > 0 {
					subrange_depth++
				}

				PrintLexTreeNodes_r(nodes[last_i:i+1], options, out, subrange_depth, prev_was_newline, false)
				last_i = i + 1
				subrange_index++
			}
		}

		// Process last segment specially.
		if last_i != 0 {

			if !*prev_was_newline {
				out.WriteString(options.line_end_sequence)
				*prev_was_newline = true
			}

			subrange_depth := depth
			if subrange_index > 0 {
				subrange_depth++
			}

			PrintLexTreeNodes_r(nodes[last_i:], options, out, subrange_depth, prev_was_newline, false)
			return
		}
	}

	for i, node := range nodes {

		if *prev_was_newline && !force_single_line {
			for i := uint(0); i < depth; i++ {
				out.WriteString(options.indentation_sequence)
			}
		}

		if node.sub_elements == nil {

			if !*prev_was_newline && i > 0 && WhitespaceIsNeeded(&nodes[i-1].lexem, &node.lexem) {
				out.WriteString(" ")
			}

			out.WriteString(node.lexem.text)
			*prev_was_newline = false

		} else {

			out.WriteString(node.lexem.text)

			if len(node.sub_elements) > 0 {
				out.WriteString(" ")
			}
			*prev_was_newline = false

			PrintLexTreeNodes_r(
				node.sub_elements,
				options,
				out,
				depth,
				prev_was_newline,
				force_single_line)

			if len(node.sub_elements) > 0 {
				out.WriteString(" ")
			}

			out.WriteString(node.trailing_lexem.text)
			*prev_was_newline = false
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

func CanWriteInSingleLine(nodes LexTreeNodeList, depth uint, options *FormattingOptions) bool {

	// Write all in single line and count length.

	builder := strings.Builder{}
	prev_was_newline := false
	PrintLexTreeNodes_r(nodes, options, &builder, depth, &prev_was_newline, true)

	return CalculateLineWidth(builder.String(), options) <= options.max_line_width
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
