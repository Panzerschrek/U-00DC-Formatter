package main

import (
	"strings"
)

func PrintLexTreeNodes(nodes LexTreeNodeList, options *FormattingOptions) string {
	var prev_was_newline bool = false
	builder := strings.Builder{}
	PrintLexTreeNodes_r(nodes, options, &builder, 0, &prev_was_newline, true)
	return builder.String()
}

func PrintLexTreeNodes_r(
	nodes LexTreeNodeList,
	options *FormattingOptions,
	out *strings.Builder,
	depth int,
	prev_was_newline *bool,
	semicolon_is_newline bool) {

	for i, node := range nodes {

		if node.lexem.t != LexemTypeSemicolon && i > 0 && nodes[i-1].trailing_lexem.t == LexemTypeBraceRight {
			// Add extra empty line after "}", except it is "else".
			// This ensures that global things like classes or functions are always separated by an empty line.
			if node.lexem.text != "else" {
				out.WriteString("\n")
				*prev_was_newline = true
			}
		}

		if *prev_was_newline {
			for i := 0; i < depth; i++ {
				out.WriteString("\t")
			}
		}

		if node.sub_elements == nil {
			if node.lexem.t == LexemTypeSemicolon {

				// Add newline after ";", if necessary
				if semicolon_is_newline {

					out.WriteString(node.lexem.text)
					out.WriteString("\n")
					*prev_was_newline = true

				} else {

					out.WriteString(node.lexem.text)
					*prev_was_newline = false
				}

			} else if node.lexem.t == LexemTypeLineComment {

				out.WriteString(node.lexem.text)
				out.WriteString("\n")
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
				out.WriteString("\n")
				*prev_was_newline = true
			}

		} else {

			// For now add newlines only before/after {}
			if node.lexem.t == LexemTypeBraceLeft {

				out.WriteString("\n")

				for i := 0; i < depth; i++ {
					out.WriteString("\t")
				}

				out.WriteString(node.lexem.text)
				out.WriteString("\n")
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

			// Insert unconditional newlines after semicolon only in blocks, not in (), [], <//>, etc.
			// This prevents making "for" operator ugly.
			subelements_semicolon_is_newline := node.lexem.t == LexemTypeBraceLeft

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

			PrintLexTreeNodes_r(
				node.sub_elements,
				options,
				out,
				sub_elements_depth,
				prev_was_newline,
				subelements_semicolon_is_newline)

			// For now add newlines only before/after {}
			if node.trailing_lexem.t == LexemTypeBraceRight {

				if !*prev_was_newline {
					out.WriteString("\n")
				}

				for i := 0; i < depth; i++ {
					out.WriteString("\t")
				}

				out.WriteString(node.trailing_lexem.text)
				out.WriteString("\n")
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
					out.WriteString("\n")
					*prev_was_newline = true
				}

			} else {

				out.WriteString(node.trailing_lexem.text)
			}
		}
	}
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
