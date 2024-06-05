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
	for _, node := range nodes {
		PrintLexTreeNode_r(0, &node)
	}
}

func PrintLexTree(node *LexTreeNode) {
	PrintLexTreeNode_r(0, node)
}

func PrintLexTreeNode_r(depth int, node *LexTreeNode) {

	if len(node.text) > 0 {
		for i := 0; i < depth; i++ {
			fmt.Print("\t")
		}
		fmt.Println(node.text)
	}

	for _, sub_node := range node.sub_elements {
		PrintLexTreeNode_r(depth+1, &sub_node)
	}

	if len(node.trailing_text) > 0 {
		for i := 0; i < depth; i++ {
			fmt.Print("\t")
		}
		fmt.Println(node.trailing_text)
	}
}
