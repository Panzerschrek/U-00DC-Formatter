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

	text_by_lines := SplitLexTreeIntoLines(lex_tree)
	text_formatted := PrintLines(text_by_lines, &options)
	fmt.Print(text_formatted)
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
