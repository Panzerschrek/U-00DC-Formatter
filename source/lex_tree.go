package main

import (
	"errors"
)

// Do not perform proper syntax analysis.
// It's not possible due to complication with macros.
// Build simple tree structure instead - where lexems inside paired symbols ( (), [], {}, <//>, <??>) are grouped together.

type LexTreeNodeList = []LexTreeNode

type LexTreeNode struct {
	lexem          Lexem
	sub_elements   LexTreeNodeList
	trailing_lexem Lexem
}

func BuildLexTree(lexems []Lexem) (LexTreeNodeList, error) {
	return ParseLexTree_r(&lexems, LexemTypeEndOfFile)
}

// Parse until specified end lexem.
func ParseLexTree_r(lexems *[]Lexem, end_lexem_type LexemType) (LexTreeNodeList, error) {
	result := make([]LexTreeNode, 0)

	for len(*lexems) > 0 {
		lexem := &(*lexems)[0]

		if lexem.t == end_lexem_type {
			break
		}

		*lexems = (*lexems)[1:]

		if lexem.t == LexemTypeBraceLeft {

			sub_elements, err := ParseLexTree_r(lexems, LexemTypeBraceRight)
			if err != nil {
				return nil, err
			}

			node := LexTreeNode{lexem: *lexem, sub_elements: sub_elements}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeBraceRight) {
				return nil, errors.New("non-matching }")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeBracketLeft {

			sub_elements, err := ParseLexTree_r(lexems, LexemTypeBracketRight)
			if err != nil {
				return nil, err
			}

			node := LexTreeNode{lexem: *lexem, sub_elements: sub_elements}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeBracketRight) {
				return nil, errors.New("non-matching )")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeSquareBracketLeft {

			sub_elements, err := ParseLexTree_r(lexems, LexemTypeSquareBracketRight)
			if err != nil {
				return nil, err
			}

			node := LexTreeNode{lexem: *lexem, sub_elements: sub_elements}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeSquareBracketRight) {
				return nil, errors.New("non-matching ]")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeTemplateBracketLeft {

			sub_elements, err := ParseLexTree_r(lexems, LexemTypeTemplateBracketRight)
			if err != nil {
				return nil, err
			}

			node := LexTreeNode{lexem: *lexem, sub_elements: sub_elements}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeTemplateBracketRight) {
				return nil, errors.New("non-matching />")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else if lexem.t == LexemTypeMacroBracketLeft {

			sub_elements, err := ParseLexTree_r(lexems, LexemTypeMacroBracketRight)
			if err != nil {
				return nil, err
			}

			node := LexTreeNode{lexem: *lexem, sub_elements: sub_elements}

			if !(len(*lexems) > 0 && (*lexems)[0].t == LexemTypeMacroBracketRight) {
				return nil, errors.New("non-matching ?>")
			}

			node.trailing_lexem = (*lexems)[0]
			*lexems = (*lexems)[1:]

			result = append(result, node)

		} else {
			result = append(result, LexTreeNode{lexem: *lexem})
		}
	}

	return result, nil
}
