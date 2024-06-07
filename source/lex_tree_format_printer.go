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

	for i, node := range nodes {

		if node.lexem.t != LexemTypeSemicolon && i > 0 && nodes[i-1].trailing_lexem.t == LexemTypeBraceRight {
			// Add extra empty line after "}", except it is "else".
			// This ensures that global things like classes or functions are always separated by an empty line.
			if node.lexem.text != "else" {
				out.WriteString(newline_char)
				*prev_was_newline = true
			}
		}

		if *prev_was_newline && !force_single_line {
			for i := 0; i < depth; i++ {
				out.WriteString("\t")
			}
		}

		if node.sub_elements == nil {

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

			// For now add newlines only before/after {}
			if node.lexem.t == LexemTypeBraceLeft {

				out.WriteString(newline_char)

				if !force_single_line {
					for i := 0; i < depth; i++ {
						out.WriteString("\t")
					}
				}

				out.WriteString(node.lexem.text)
				out.WriteString(newline_char)
				*prev_was_newline = true

			} else if node.lexem.t == LexemTypeBracketLeft {

				out.WriteString(node.lexem.text)

				// Add spaces only in case of non-empty subelements inside ()
				if len(node.sub_elements) > 0 {
					out.WriteString(" ")
				}

			} else {
				out.WriteString(node.lexem.text)
			}

			// Somewhat hacky namespaces detection.
			// Assuming "{" follows directly after something like "namespace SomeName".
			// TODO - skip also comments, newlines, etc. in this check.
			is_namespace := i >= 2 && nodes[i-1].lexem.t == LexemTypeIdentifier && nodes[i-2].lexem.text == "namespace"

			// Hacky template declaration detection.
			is_template_declaration := node.lexem.t == LexemTypeTemplateBracketLeft && i >= 1 && nodes[i-1].lexem.text == "template"

			// For namespaces avoid adding extra intendation.
			// TODO - make this behavior configurabe.
			sub_elements_depth := depth + 1
			if is_namespace {
				sub_elements_depth -= 1
			}

			// Check if can write subelements in single line.
			sub_elements_force_single_line := force_single_line
			if !sub_elements_force_single_line {
				if CanWriteInSingleLine(&node, options) {
					sub_elements_force_single_line = true
				}
			}

			PrintLexTreeNodes_r(
				node.sub_elements,
				options,
				out,
				sub_elements_depth,
				prev_was_newline,
				sub_elements_force_single_line)

			// For now add newlines only before/after {}
			if node.trailing_lexem.t == LexemTypeBraceRight {

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

			} else if node.trailing_lexem.t == LexemTypeBracketRight {

				// Add spaces only in case of non-empty subelements inside ()
				if len(node.sub_elements) > 0 {
					out.WriteString(" ")
				}
				out.WriteString(node.trailing_lexem.text)

			} else if node.trailing_lexem.t == LexemTypeTemplateBracketRight {

				out.WriteString(node.trailing_lexem.text)

				if is_template_declaration {
					out.WriteString(newline_char)
					*prev_was_newline = true
				}

			} else {

				out.WriteString(node.trailing_lexem.text)
			}
		}
	}
}

func CanWriteInSingleLine(node *LexTreeNode, options *FormattingOptions) bool {

	// Fast check - if this node requires newline.
	if node.lexem.t == LexemTypeLineComment || node.lexem.t == LexemTypeSemicolon {
		return false
	}

	// Recursively check contents.
	for _, sub_node := range node.sub_elements {
		if !CanWriteInSingleLine(&sub_node, options) {
			return false
		}
	}

	// More detail check - write all in single line and count length.
	builder := strings.Builder{}
	prev_was_newline := false
	depth := 0 // TODO - pass it
	PrintLexTreeNodes_r(node.sub_elements, options, &builder, depth, &prev_was_newline, true)

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
