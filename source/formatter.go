package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	file_contents := ReadFile(args[0])

	lexems, err := SplitProgramIntoLexems(file_contents)
	if err != nil {
		panic(err)
	}

	if false {
		for _, lexem := range lexems {
			fmt.Print(lexem.text)
			fmt.Print(" ")
		}
	}
	fmt.Print("\n")

	lex_tree, err := BuildLexTree(lexems)
	if err != nil {
		panic(err)
	}

	options := GetDefaultFormattingOptions()
	text_formatted := PrintLexTreeNodes(lex_tree, &options)
	if false {
		fmt.Print(text_formatted)
	}

	text_by_lines := SplitLexTreeIntoLines(lex_tree)
	for _, line := range text_by_lines {
		for i := uint(0); i < line.indentation; i++ {
			fmt.Print("\t")
		}

		for _, lexem := range line.lexems {
			fmt.Print(lexem.text, " ")
		}

		fmt.Print("\n")
	}
}

func ReadFile(s string) string {
	file, e := os.Open(s)
	if e != nil {
		panic(e)
	}

	stat, e := file.Stat()
	if e != nil {
		panic(e)
	}

	size := stat.Size()

	bytes := make([]byte, size)

	read_size, e := file.Read(bytes)
	if e != nil {
		panic(e)
	}
	// TODO - read in loop
	if int64(read_size) != size {
		panic("Unexpected read size!")
	}

	file.Close()

	return string(bytes)
}
