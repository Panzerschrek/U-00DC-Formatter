package main

// TODO - use better name?
type LogicalLine = struct {
	indentation uint
	lexems      []Lexem
}

// Convert lex tree into line by line representation.
func SplitLexTreeIntoLines(nodes LexTreeNodeList) []LogicalLine {

	result := make([]LogicalLine, 0)
	AddNewLine(&result, 0)

	prev_was_newline := false
	SplitLexTreeIntoLines_r(nodes, 0, &result, &prev_was_newline)

	return result
}

func SplitLexTreeIntoLines_r(nodes LexTreeNodeList, indentation uint, out *[]LogicalLine, prev_was_newline *bool) {
	for i, node := range nodes {

		if node.lexem.t != LexemTypeSemicolon && i > 0 && nodes[i-1].trailing_lexem.t == LexemTypeBraceRight {
			// Add extra empty line after "}", except it is "else", ".".
			// This ensures that global things like classes or functions are always separated by an empty line.
			if !(node.lexem.t == LexemTypeDot || node.lexem.text == "else") {
				AddNewLine(out, indentation)
				*prev_was_newline = true
			}
		}

		if node.sub_elements == nil {

			AppendToLastLine(out, node.lexem)

			if node.lexem.t == LexemTypeSemicolon {

				// Add newline after ";".
				// TODO - do this only if it is necessary (allow ";" in single-line "for" operator).
				AddNewLine(out, indentation)
				*prev_was_newline = true

			} else if node.lexem.t == LexemTypeLineComment {

				// Always add newline after line comment.
				AddNewLine(out, indentation)
				*prev_was_newline = true

			} else {
				*prev_was_newline = false
			}

			if !*prev_was_newline && i > 0 && nodes[i-1].lexem.text == "import" {
				// Add newlines after imports.
				AddNewLine(out, indentation)
				*prev_was_newline = true
			}

		} else {

			subelements_contain_natural_newlines := HasNaturalNewlines(node.sub_elements)

			if subelements_contain_natural_newlines {

				if !*prev_was_newline {
					AddNewLine(out, indentation)
				}

				AppendToLastLine(out, node.lexem)
				AddNewLine(out, indentation+1)
				*prev_was_newline = true

			} else {
				AppendToLastLine(out, node.lexem)
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
			sub_elements_indentation := indentation + 1
			if is_namespace {
				sub_elements_indentation--
			}

			SplitLexTreeIntoLines_r(
				node.sub_elements,
				sub_elements_indentation,
				out,
				prev_was_newline)

			if subelements_contain_natural_newlines {

				if !*prev_was_newline {
					AddNewLine(out, indentation)
				} else {
					(*out)[len(*out)-1].indentation = indentation
				}

				AppendToLastLine(out, node.trailing_lexem)
				AddNewLine(out, indentation)
				*prev_was_newline = true

			} else {

				AppendToLastLine(out, node.trailing_lexem)
			}
		}
	}
}

func AppendToLastLine(lines *[]LogicalLine, lexem Lexem) {
	line := &(*lines)[len(*lines)-1]
	line.lexems = append(line.lexems, lexem)
}

func AddNewLine(lines *[]LogicalLine, indentation uint) {
	*lines = append(*lines, LogicalLine{indentation: indentation, lexems: make([]Lexem, 0)})
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
