package main

import (
	"fmt"
)

type LexTreeNodeList = []LexTreeNode

type LexTreeNode struct {
	lexem          Lexem
	sub_elements   LexTreeNodeList
	trailing_lexem Lexem
}

func BuildLexTree(lexems []Lexem) LexTreeNodeList {
	return ParseLexTree_r(&lexems, LexemTypeEndOfFile)
}

// Parse until specified end of line.
func ParseLexTree_r(lexems *[]Lexem, end_lexem_type LexemType) LexTreeNodeList {
	result := make([]LexTreeNode, 0)

	for len(*lexems) > 0 {
		lexem := &(*lexems)[0]

		if lexem.t == end_lexem_type {
			break
		}

		*lexems = (*lexems)[1:]

		if lexem.t == LexemTypeBraceLeft {

			node := LexTreeNode{lexem: *lexem, sub_elements: ParseLexTree_r(lexems, LexemTypeBraceRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeBraceRight) {
				panic("non-matching }")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeBracketLeft {

			node := LexTreeNode{lexem: *lexem, sub_elements: ParseLexTree_r(lexems, LexemTypeBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeBracketRight) {
				panic("non-matching )")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeSquareBracketLeft {

			node := LexTreeNode{lexem: *lexem, sub_elements: ParseLexTree_r(lexems, LexemTypeSquareBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeSquareBracketRight) {
				panic("non-matching ]")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeTemplateBracketLeft {

			node := LexTreeNode{lexem: *lexem, sub_elements: ParseLexTree_r(lexems, LexemTypeTemplateBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeTemplateBracketRight) {
				panic("non-matching />")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeMacroBracketLeft {

			node := LexTreeNode{lexem: *lexem, sub_elements: ParseLexTree_r(lexems, LexemTypeMacroBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeMacroBracketRight) {
				panic("non-matching ?>")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else {
			result = append(result, LexTreeNode{lexem: *lexem})
		}
	}

	return result
}

func PrintLexTreeNodes(nodes LexTreeNodeList) {
	var prev_was_newline bool = false
	PrintLexTreeNodes_r(nodes, 0, &prev_was_newline, true)
}

func PrintLexTreeNodes_r(nodes LexTreeNodeList, depth int, prev_was_newline *bool, semicolon_is_newline bool) {

	for i, node := range nodes {

		if *prev_was_newline {
			for i := 0; i < depth; i++ {
				fmt.Print("\t")
			}
		}

		if node.sub_elements == nil {
			if node.lexem.t == LexemTypeSemicolon {

				// Add newline after ";", if necessary
				if semicolon_is_newline {

					fmt.Print(node.lexem.text, "\n")
					*prev_was_newline = true

				} else {

					fmt.Print(node.lexem.text, " ")
					*prev_was_newline = false
				}

			} else if node.lexem.t == LexemTypeLineComment {

				fmt.Print(node.lexem.text, "\n")
				*prev_was_newline = true

			} else {

				// Add spaces between lexems.
				// TODO - make this confugurable - depending on lexem types.
				fmt.Print(node.lexem.text, " ")
				*prev_was_newline = false
			}

			if !*prev_was_newline && i > 0 && nodes[i-1].lexem.text == "import" {
				// Add newlines after imports.
				fmt.Print("\n")
				*prev_was_newline = true
			}

		} else {

			// For now add newlines only before/after {}
			if node.lexem.t == LexemTypeBraceLeft {

				fmt.Print("\n")

				for i := 0; i < depth; i++ {
					fmt.Print("\t")
				}

				fmt.Print(node.lexem.text)
				fmt.Print("\n")
				*prev_was_newline = true

			} else {

				fmt.Print(node.lexem.text, " ")
			}

			// Insert unconditional newlines after semicolon only in blocks, not in (), [], <//>, etc.
			// This prevents making "for" operator ugly.
			subelements_semicolon_is_newline := node.lexem.t == LexemTypeBraceLeft

			// Somewhat hacky namespaces detection.
			// Assuming "{" follows directly after something like "namespace SomeName".
			// TODO - skip also comments, newlines, etc. in this check. 
			is_namespace := i >= 2 && nodes[i - 1].lexem.t == LexemTypeIdentifier && nodes[ i - 2 ].lexem.text == "namespace"

			// For namespaces avoid adding extra intendation.
			// TODO - make this behavior configurabe.
			sub_elements_depth := depth + 1
			if is_namespace {
				sub_elements_depth -= 1
			}

			PrintLexTreeNodes_r(node.sub_elements, sub_elements_depth, prev_was_newline, subelements_semicolon_is_newline)

			// For now add newlines only before/after {}
			if node.trailing_lexem.t == LexemTypeBraceRight {

				if !*prev_was_newline {
					fmt.Print("\n")
				}

				for i := 0; i < depth; i++ {
					fmt.Print("\t")
				}

				fmt.Print(node.trailing_lexem.text)
				fmt.Print("\n")
				*prev_was_newline = true

			} else {

				fmt.Print(node.trailing_lexem.text, " ")
			}
		}
	}
}
