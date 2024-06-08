package main

type FormattingOptions struct {
	indentation_sequence string
	line_end_sequence    string
	tab_size             uint
	max_line_width       uint
}

func GetDefaultFormattingOptions() FormattingOptions {
	return FormattingOptions{
		indentation_sequence: "\t",
		line_end_sequence:    "\n",
		tab_size:             4,
		max_line_width:       120}
}
