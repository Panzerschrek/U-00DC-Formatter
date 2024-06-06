package main

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
