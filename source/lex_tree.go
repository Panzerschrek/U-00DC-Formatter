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

func ParseLexTree_r(lexems *[]Lexem, end_lexem_type LexemType) LexTreeNodeList {
	result := make([]LexTreeNode, 0)

	for len(*lexems) > 0 {
		lexem := &(*lexems)[0]

		if lexem.t == end_lexem_type {
			break
		} else if lexem.t == LexemTypeBraceLeft {

			*lexems = (*lexems)[1:]

			node := LexTreeNode{text: lexem.text}
			node.sub_elements = ParseLexTree_r(lexems, LexemTypeBraceRight)

			if len(*lexems) > 0 {
				trailing_lexem := &(*lexems)[0] // TODO - check it is valid
				*lexems = (*lexems)[1:]
				node.trailing_text = trailing_lexem.text
			}

			result = append(result, node)

		} else if lexem.t == LexemTypeBracketLeft {
			*lexems = (*lexems)[1:]

			node := LexTreeNode{text: lexem.text}
			node.sub_elements = ParseLexTree_r(lexems, LexemTypeBracketRight)

			if len(*lexems) > 0 {
				trailing_lexem := &(*lexems)[0] // TODO - check it is valid
				*lexems = (*lexems)[1:]
				node.trailing_text = trailing_lexem.text
			}

			result = append(result, node)

		} else if lexem.t == LexemTypeSquareBracketLeft {
			*lexems = (*lexems)[1:]

			node := LexTreeNode{text: lexem.text}
			node.sub_elements = ParseLexTree_r(lexems, LexemTypeSquareBracketRight)

			if len(*lexems) > 0 {
				trailing_lexem := &(*lexems)[0] // TODO - check it is valid
				*lexems = (*lexems)[1:]
				node.trailing_text = trailing_lexem.text
			}

			result = append(result, node)

		} else {
			*lexems = (*lexems)[1:]
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
