package main

import (
	"strings"
)

// Convert parsed LexTree into string representation.
// Use passed options.
func PrintLexTreeNodes(nodes LexTreeNodeList, options *FormattingOptions) string {
	var prev_was_newline bool = false
	builder := strings.Builder{}
	PrintLexTreeNodes_r(nodes, options, &builder, 0, &prev_was_newline, false)
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
	depth int,
	prev_was_newline *bool,
	force_single_line bool) {

	// In force single line mode replace all newlines with spaces.
	newline_char := "\n"
	if force_single_line {
		newline_char = " "
	}

	if len(nodes) > 1 &&
		!force_single_line &&
		!HasNaturalNewlines(nodes) &&
		!CanWriteInSingleLine(nodes, options) {
		// Recursively split and print this list, adding newlines in split points.

		// Search for the most important lexem type to use it as splitter.
		max_priority := 0
		for _, node := range nodes {
			priority := GetLineSplitLexemPriority(&node.lexem)
			if priority > max_priority {
				max_priority = priority
			}
		}

		// Split this lexems list into parts, using maximum priority lexem type.
		// Add newline after each part.
		last_i := 0
		for i := 0; i < len(nodes)-1; i++ {
			if GetLineSplitLexemPriority(&nodes[i].lexem) == max_priority {

				PrintLexTreeNodes_r(nodes[last_i:i+1], options, out, depth, prev_was_newline, false)
				last_i = i + 1

				if !*prev_was_newline {
					out.WriteString("\n")
					*prev_was_newline = true
				}
			}
		}

		// Process last segment specially.
		if last_i != 0 {
			PrintLexTreeNodes_r(nodes[last_i:], options, out, depth, prev_was_newline, false)
			return
		}
	}

	for i, node := range nodes {

		if node.lexem.t != LexemTypeSemicolon && i > 0 && nodes[i-1].trailing_lexem.t == LexemTypeBraceRight {
			// Add extra empty line after "}", except it is "else".
			// This ensures that global things like classes or functions are always separated by an empty line.
			if node.lexem.text != "else" {
				out.WriteString(newline_char)
				*prev_was_newline = true
			}
		}

		if node.sub_elements == nil {

			if *prev_was_newline && !force_single_line {
				for i := 0; i < depth; i++ {
					out.WriteString("\t")
				}
			}

			if node.lexem.t == LexemTypeSemicolon {

				// Add newline after ";".
				// TODO - do this only if it is necessary.
				out.WriteString(node.lexem.text)
				out.WriteString(newline_char)
				*prev_was_newline = true

			} else if node.lexem.t == LexemTypeLineComment {

				// Always add newline after line comment.
				out.WriteString(node.lexem.text)
				out.WriteString(newline_char)
				*prev_was_newline = true

			} else {

				if !*prev_was_newline && i > 0 && WhitespaceIsNeeded(&nodes[i-1].lexem, &node.lexem) {
					out.WriteString(" ")
				}

				out.WriteString(node.lexem.text)
				*prev_was_newline = false
			}

			if !*prev_was_newline && i > 0 && nodes[i-1].lexem.text == "import" {
				// Add newlines after imports.
				out.WriteString(newline_char)
				*prev_was_newline = true
			}

		} else {

			subelements_contain_natural_newlines := !force_single_line && HasNaturalNewlines(node.sub_elements)

			if subelements_contain_natural_newlines {

				if !*prev_was_newline {
					out.WriteString(newline_char)
				}

				if !force_single_line {
					for i := 0; i < depth; i++ {
						out.WriteString("\t")
					}
				}

				out.WriteString(node.lexem.text)
				out.WriteString(newline_char)
				*prev_was_newline = true

			} else {
				out.WriteString(node.lexem.text)

				if len(node.sub_elements) > 0 {
					out.WriteString(" ")
				}
			}

			// Somewhat hacky namespaces detection.
			// Assuming "{" follows directly after something like "namespace SomeName".
			// TODO - skip also comments, newlines, etc. in this check.
			is_namespace := i >= 2 && nodes[i-1].lexem.t == LexemTypeIdentifier && nodes[i-2].lexem.text == "namespace"

			// Hacky template declaration detection.
			is_template_declaration := node.lexem.t == LexemTypeTemplateBracketLeft && i >= 1 && nodes[i-1].lexem.text == "template"
			_ = is_template_declaration // TODO - use it

			// For namespaces avoid adding extra intendation.
			// TODO - make this behavior configurabe.
			sub_elements_depth := depth + 1
			if is_namespace {
				sub_elements_depth--
			}

			PrintLexTreeNodes_r(
				node.sub_elements,
				options,
				out,
				sub_elements_depth,
				prev_was_newline,
				force_single_line)

			if subelements_contain_natural_newlines {

				if !*prev_was_newline {
					out.WriteString(newline_char)
				}

				if !force_single_line {
					for i := 0; i < depth; i++ {
						out.WriteString("\t")
					}
				}

				out.WriteString(node.trailing_lexem.text)
				out.WriteString(newline_char)
				*prev_was_newline = true

			} else {

				if len(node.sub_elements) > 0 {
					out.WriteString(" ")
				}

				out.WriteString(node.trailing_lexem.text)
			}
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

func CanWriteInSingleLine(nodes LexTreeNodeList, options *FormattingOptions) bool {

	// Write all in single line and count length.

	builder := strings.Builder{}
	prev_was_newline := false
	depth := 0 // TODO - pass it
	PrintLexTreeNodes_r(nodes, options, &builder, depth, &prev_was_newline, true)

	s := builder.String()

	// Evaluate result string length.
	// Treat tabs specially.
	len := uint(0)
	for _, c := range s {
		if c == '\t' {
			len += options.tab_size
		} else {
			len++
		}
	}

	return len <= options.max_line_width
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
