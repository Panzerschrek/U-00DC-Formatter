package main

import (
	"strings"
)

// Convert line-by-line representation into text representation, split too ling lines if necessary.
func PrintLines(lines []LogicalLine, options *FormattingOptions) string {

	text_builder := strings.Builder{}

	for _, line := range lines {

		line_builder := strings.Builder{}

		for i := uint(0); i < line.indentation; i++ {
			line_builder.WriteString(options.indentation_sequence)
		}

		line_width := CountIndentationsSize(line.indentation, options)

		for i, lexem := range line.lexems {
			if i > 0 && WhitespaceIsNeeded(&line.lexems[i-1], &lexem) {
				line_builder.WriteString(" ")
				line_width++
			}
			line_builder.WriteString(lexem.text)
			line_width += uint(len(lexem.text))
		}

		line_builder.WriteString(options.line_end_sequence)

		if line_width <= options.max_line_width {
			// Fine - line width does not exeed the limit.
			text_builder.WriteString(line_builder.String())
		} else {
			// Try to split this line.
			// Build lex_tree again, but only for this line.
			lex_tree, err := BuildLexTree(line.lexems)
			_ = err // TODO - handle it
			line_splitted := PrintAndSplitLexTree(lex_tree, line.indentation, options)
			text_builder.WriteString(line_splitted)
		}
	}

	return text_builder.String()
}

func WhitespaceIsNeeded(l *Lexem, r *Lexem) bool {
	switch r.t {
	case LexemTypeNone:

	case LexemTypeLineComment:
		return true

	case LexemTypeIdentifier:
		if l.t == LexemTypeDot || l.t == LexemTypeScope || l.t == LexemTypeTemplateBracketRight {
			return false
		}
		if l.t == LexemTypeAnd {
			if r.text == "mut" || r.text == "imut" || r.text == "constexpr" {
				// Allow "&mut", "&imut", "&constexpr".
				return false
			}
		}
		return true

	case LexemTypeMacroIdentifier:
		return true

	case LexemTypeMacroUniqueIdentifier:
		return true

	case LexemTypeString:
		return true

	case LexemTypeNumber:
		return true

	case LexemTypeLiteralSuffix:

	case LexemTypeBracketLeft:
		if l.t == LexemTypeBracketLeft {
			return true
		}
		return false

	case LexemTypeBracketRight:
		if l.t == LexemTypeBracketLeft {
			return false
		}
		return true

	case LexemTypeSquareBracketLeft:
		return false

	case LexemTypeSquareBracketRight:
		return false

	case LexemTypeBraceLeft:
		return false

	case LexemTypeBraceRight:
		return false

	case LexemTypeTemplateBracketLeft:
		return false

	case LexemTypeTemplateBracketRight:
		return true

	case LexemTypeMacroBracketLeft:
		return true

	case LexemTypeMacroBracketRight:
		return true

	case LexemTypeScope:
		return false

	case LexemTypeComma:
		return false

	case LexemTypeDot:
		if l.t == LexemTypeComma {
			// Add spaces in struct named initializer before ".".
			// But in member access use no space before ".".
			return true
		}
		return false

	case LexemTypeColon:
		return true

	case LexemTypeSemicolon:
		return false

	case LexemTypeQuestion:
		return true

	case LexemTypeAssignment:
		return true

	case LexemTypePlus:
		return true

	case LexemTypeMinus:
		return true // TODO - detect unary minus

	case LexemTypeStar:
		return true

	case LexemTypeSlash:
		return true

	case LexemTypePercent:
		return true

	case LexemTypeAnd:
		// TODO - check cases with &mut
		// TODO - check case wih auto&
		// TODO - check case like var i32& x
		return true

	case LexemTypeOr:
		return true

	case LexemTypeXor:
		return true

	case LexemTypeTilda:
		return false

	case LexemTypeNot:
		return false

	case LexemTypeApostrophe:
		return true

	case LexemTypeAt:
		return true

	case LexemTypeIncrement:
		return false

	case LexemTypeDecrement:
		return false

	case LexemTypeCompareLess:
		return true

	case LexemTypeCompareGreater:
		return true

	case LexemTypeCompareEqual:
		return true

	case LexemTypeCompareNotEqual:
		return true

	case LexemTypeCompareLessOrEqual:
		return true

	case LexemTypeCompareGreaterOrEqual:
		return true

	case LexemTypeCompareOrder:
		return true

	case LexemTypeConjunction:
		return true

	case LexemTypeDisjunction:
		return true

	case LexemTypeAssignAdd:
		return true

	case LexemTypeAssignSub:
		return true

	case LexemTypeAssignMul:
		return true

	case LexemTypeAssignDiv:
		return true

	case LexemTypeAssignRem:
		return true

	case LexemTypeAssignAnd:
		return true

	case LexemTypeAssignOr:
		return true

	case LexemTypeAssignXor:
		return true

	case LexemTypeShiftLeft:
		return true

	case LexemTypeShiftRight:
		return true

	case LexemTypeAssignShiftLeft:
		return true

	case LexemTypeAssignShiftRight:
		return true

	case LexemTypeRightArrow:
		return true

	case LexemTypePointerTypeMark:
		return true

	case LexemTypeReferenceToPointer:
		return true

	case LexemTypePointerToReference:
		return true

	case LexemTypeEllipsis:
		return true

	case LexemTypeEndOfFile:
		return false
	}

	return true
}

func PrintAndSplitLexTree(nodes LexTreeNodeList, indentation uint, options *FormattingOptions) string {
	builder := strings.Builder{}

	for i := uint(0); i < indentation; i++ {
		builder.WriteString(options.indentation_sequence)
	}

	current_line_width := CountIndentationsSize(indentation, options)

	PrintAndSplitLexTree_r(nodes, options, &builder, indentation, &current_line_width)

	builder.WriteString(options.line_end_sequence)

	return builder.String()
}

// Main recursive routine for splitting of lex_tree into multiple lines.
// Has exponentioal complexity, but it should not be a big problem, since single-line lexical trees are pretty small.
func PrintAndSplitLexTree_r(
	nodes LexTreeNodeList,
	options *FormattingOptions,
	out *strings.Builder,
	indentation uint,
	current_line_width *uint) {

	split_results := make([]SplittingResult, 0)

	// Perform split at current level to give it more priority.
	current_level_split_result := SplitNodeListAtCurrentLevel(nodes, options, indentation, *current_line_width)
	if current_level_split_result != nil {
		split_results = append(split_results, *current_level_split_result)
	}

	further_level_split_result := PrintAndSplitNodeListAtFurtherLevels(nodes, options, indentation, *current_line_width)
	split_results = append(split_results, further_level_split_result)

	best_result := ChooseBestSplitResult(split_results, *current_line_width, options)
	out.WriteString(best_result.text)
	*current_line_width = best_result.current_line_width
}

func SplitNodeListAtCurrentLevel(
	nodes LexTreeNodeList,
	options *FormattingOptions,
	indentation uint,
	current_line_width uint) *SplittingResult {

	if len(nodes) <= 1 {
		return nil // Can-t split single node.
	}

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

			PrintAndSplitLexTree_r(nodes[last_i:i+1], options, &builder, next_indentation, &current_line_width)
			last_i = i + 1

			next_indentation = indentation + 1

			builder.WriteString(options.line_end_sequence)
			for i := uint(0); i < next_indentation; i++ {
				builder.WriteString(options.indentation_sequence)
			}
			current_line_width = CountIndentationsSize(next_indentation, options)
		}
	}

	// Process last segment specially.
	PrintAndSplitLexTree_r(nodes[last_i:], options, &builder, next_indentation, &current_line_width)

	return &SplittingResult{current_line_width: current_line_width, text: builder.String()}
}

func PrintAndSplitNodeListAtFurtherLevels(
	nodes LexTreeNodeList,
	options *FormattingOptions,
	indentation uint,
	current_line_width uint) SplittingResult {

	builder := strings.Builder{}

	for i, node := range nodes {

		if node.sub_elements == nil {

			if i > 0 && WhitespaceIsNeeded(&nodes[i-1].lexem, &node.lexem) {
				builder.WriteString(" ")
				current_line_width++
			}

			builder.WriteString(node.lexem.text)
			current_line_width += uint(len(node.lexem.text))

		} else {

			builder.WriteString(node.lexem.text)
			current_line_width += uint(len(node.lexem.text))

			if len(node.sub_elements) > 0 {
				builder.WriteString(" ")
				current_line_width++
			}

			PrintAndSplitLexTree_r(
				node.sub_elements, options, &builder, indentation, &current_line_width)

			if len(node.sub_elements) > 0 {
				builder.WriteString(" ")
				current_line_width++
			}

			builder.WriteString(node.trailing_lexem.text)
			current_line_width += uint(len(node.trailing_lexem.text))
		}
	}

	return SplittingResult{current_line_width: current_line_width, text: builder.String()}
}

type SplittingResult struct {
	current_line_width uint
	text               string
}

// Choose best result based on number of lines.
// If number of lines is equal, choose firs result with such number.
func ChooseBestSplitResult(
	results []SplittingResult, current_line_width uint, options *FormattingOptions) *SplittingResult {

	if len(results) == 0 {
		panic("No splitting results!")
		return nil
	}
	if len(results) == 1 {
		return &results[0]
	}

	type ResultStats struct {
		num_lines                uint
		exceeds_line_width_limit bool
	}

	stats := make([]ResultStats, len(results))

	for i, r := range results {

		num_newlines := uint(0)
		max_line_width := uint(0)
		w := current_line_width

		for _, c := range r.text {
			if c == '\n' { // TODO - use newline sequence from options.
				max_line_width = max(max_line_width, w)
				w = 0
				num_newlines++
			} else if c == '\t' {
				w += options.tab_size
			} else {
				w++
			}
		}

		max_line_width = max(max_line_width, w)

		stats[i].num_lines = num_newlines
		stats[i].exceeds_line_width_limit = max_line_width > options.max_line_width
	}

	num_results_over_line_limit := uint(0)
	for _, s := range stats {
		if s.exceeds_line_width_limit {
			num_results_over_line_limit++
		}
	}

	var res *SplittingResult = nil
	min_num_lines := uint(1024 * 1024)

	if num_results_over_line_limit < uint(len(stats)) {
		for i, r := range results {
			if !stats[i].exceeds_line_width_limit && stats[i].num_lines < min_num_lines {
				min_num_lines = stats[i].num_lines
				res = &r
			}
		}
	} else {
		// Fallback in cases where it isn't possible to split below the limit.
		for i, r := range results {
			if stats[i].num_lines < min_num_lines {
				min_num_lines = stats[i].num_lines
				res = &r
			}
		}
	}

	if res == nil {
		panic("Unexpected missing result!")
	}
	return res
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
