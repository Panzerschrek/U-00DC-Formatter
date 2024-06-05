package main

import (
	"fmt"
)

type LexTreeNodeList = []LexTreeNode

type LexTreeNode struct {
	text          string // Only for simple nodes
	sub_elements  LexTreeNodeList
	trailing_text string
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

			node := LexTreeNode{text: lexem.text, sub_elements: ParseLexTree_r(lexems, LexemTypeBraceRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeBraceRight) {
				panic("non-matching }")
			}

			node.trailing_text = (*lexems)[0].text
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeBracketLeft {

			node := LexTreeNode{text: lexem.text, sub_elements: ParseLexTree_r(lexems, LexemTypeBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeBracketRight) {
				panic("non-matching )")
			}

			node.trailing_text = (*lexems)[0].text
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeSquareBracketLeft {

			node := LexTreeNode{text: lexem.text, sub_elements: ParseLexTree_r(lexems, LexemTypeSquareBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeSquareBracketRight) {
				panic("non-matching ]")
			}

			node.trailing_text = (*lexems)[0].text
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeTemplateBracketLeft {

			node := LexTreeNode{text: lexem.text, sub_elements: ParseLexTree_r(lexems, LexemTypeTemplateBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeTemplateBracketRight) {
				panic("non-matching />")
			}

			node.trailing_text = (*lexems)[0].text
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeMacroBracketLeft {

			node := LexTreeNode{text: lexem.text, sub_elements: ParseLexTree_r(lexems, LexemTypeMacroBracketRight)}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeMacroBracketRight) {
				panic("non-matching ?>")
			}

			node.trailing_text = (*lexems)[0].text
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else {
			result = append(result, LexTreeNode{text: lexem.text})
		}
	}

	return result
}

func PrintLexTreeNodes(nodes LexTreeNodeList) {
	var prev_was_newline bool = false
	PrintLexTreeNodes_r(nodes, 0, &prev_was_newline)
}

func PrintLexTreeNodes_r(nodes LexTreeNodeList, depth int, prev_was_newline *bool) {

	for i, node := range nodes {

		if *prev_was_newline {
			for i := 0; i < depth; i++ {
				fmt.Print("\t")
			}
		}

		if node.sub_elements == nil {
			if node.text == ";" {
				// Always add newline after ";"
				fmt.Print(node.text, "\n")
				*prev_was_newline = true
			} else {
				fmt.Print(node.text, " ")
				*prev_was_newline = false
			}

			if !*prev_was_newline && i > 0 && nodes[i-1].text == "import" {
				// Add newlines after imports.
				fmt.Print("\n")
				*prev_was_newline = true
			}

			// TODO - add newlines after line comments.

		} else {

			// For now add newlines only before/after {}
			if node.text == "{" {
				fmt.Print("\n")
				for i := 0; i < depth; i++ {
					fmt.Print("\t")
				}
				fmt.Print(node.text)
				fmt.Print("\n")
				*prev_was_newline = true
			} else {
				fmt.Print(node.text, " ")
			}

			PrintLexTreeNodes_r(node.sub_elements, depth+1, prev_was_newline)

			// For now add newlines only before/after {}
			if node.trailing_text == "}" {
				if !*prev_was_newline {
					fmt.Print("\n")
				}

				for i := 0; i < depth; i++ {
					fmt.Print("\t")
				}
				fmt.Print(node.trailing_text)
				fmt.Print("\n")
				*prev_was_newline = true
			} else {
				fmt.Print(node.trailing_text, " ")
			}
		}
	}
}
