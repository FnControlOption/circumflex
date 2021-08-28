package parser

import (
	"clx/markdown"
	"errors"
	"regexp"
	"strings"
)

func Parse(text string) []*markdown.Block {
	var blocks []*markdown.Block

	enDash := "–"
	emDash := "—"
	normalDash := "-"

	// en- and em-dashes are occasionally used or list items.
	// converting them to normal dashes lets us parse more list items.
	text = strings.ReplaceAll(text, enDash, normalDash)
	text = strings.ReplaceAll(text, emDash, normalDash)

	text = strings.ReplaceAll(text, markdown.BoldStart, "")
	text = strings.ReplaceAll(text, markdown.BoldStop, "")

	lines := strings.Split(text+"\n", "\n")
	temp := new(tempBuffer)

	isInsideQuote := false
	isInsideCode := false
	isInsideText := false
	isInsideList := false
	isInsideTable := false

	for _, line := range lines {
		lineWithoutLeadingWhitespace := strings.TrimLeft(line, " ")

		if isInsideCode {
			if strings.HasPrefix(lineWithoutLeadingWhitespace, "```") {
				isInsideCode = false

				appendedBlocks, err := appendNonEmptyBuffer(temp, blocks)
				if err == nil {
					blocks = appendedBlocks
				}

				temp.reset()

				continue
			}

			temp.append("\n" + line)

			continue
		}

		if line == "" {
			appendedBlocks, err := appendNonEmptyBuffer(temp, blocks)
			if err == nil {
				blocks = appendedBlocks
			}

			temp.reset()

			isInsideQuote = false
			isInsideText = false
			isInsideList = false
			isInsideTable = false

			continue
		}

		if isInsideTable {
			temp.append("\n" + line)

			continue
		}

		if isInsideText {
			temp.append(" " + line)

			continue
		}

		if isInsideList {
			temp.append("\n" + line)

			continue
		}

		if isInsideQuote {
			line = strings.TrimPrefix(line, ">")
			line = strings.TrimPrefix(line, " ")

			temp.append("\n" + line)

			continue
		}

		switch {
		case strings.HasPrefix(lineWithoutLeadingWhitespace, `![`):
			temp.kind = markdown.Image
			temp.text = line

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "> "):
			temp.kind = markdown.Quote
			temp.text = strings.TrimPrefix(line, "> ")

			isInsideQuote = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "```"):
			temp.kind = markdown.Code
			temp.text = ""

			isInsideCode = true

		case isListItem(lineWithoutLeadingWhitespace):
			if isSameTypeAsPreviousItem(markdown.List, blocks) {
				lastItem := len(blocks) - 1

				temp.kind = markdown.List
				temp.text = blocks[lastItem].Text + "\n" + line

				blocks = RemoveIndex(blocks, lastItem)
				isInsideList = true

				continue
			}

			temp.kind = markdown.List
			temp.text = line

			isInsideList = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "|"):
			if isSameTypeAsPreviousItem(markdown.Table, blocks) {
				lastItem := len(blocks) - 1

				temp.kind = markdown.Table
				temp.text = blocks[lastItem].Text + "\n" + line

				blocks = RemoveIndex(blocks, lastItem)
				isInsideTable = true

				continue
			}

			temp.kind = markdown.Table
			temp.text = line

			isInsideTable = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "* * *"):
			temp.kind = markdown.Divider
			temp.text = line

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "# "):
			temp.kind = markdown.H1
			temp.text = line

			isInsideText = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "## "):
			temp.kind = markdown.H2
			temp.text = line

			isInsideText = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "### "):
			temp.kind = markdown.H3
			temp.text = line

			isInsideText = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "#### "):
			temp.kind = markdown.H4
			temp.text = line

			isInsideText = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "##### "):
			temp.kind = markdown.H5
			temp.text = line

			isInsideText = true

		case strings.HasPrefix(lineWithoutLeadingWhitespace, "###### "):
			temp.kind = markdown.H6
			temp.text = line

			isInsideText = true

		default:
			temp.kind = markdown.Text
			temp.text = line

			isInsideText = true
		}
	}

	return blocks
}

func RemoveIndex(s []*markdown.Block, index int) []*markdown.Block {
	return append(s[:index], s[index+1:]...)
}

func isListItem(text string) bool {
	if text == "" {
		return false
	}

	exp := regexp.MustCompile(`^\s*(-|\d+\. )`)
	listToken := exp.FindString(text)

	return listToken != ""
}

func isSameTypeAsPreviousItem(itemType int, blocks []*markdown.Block) bool {
	if len(blocks) == 0 {
		return false
	}

	previousItem := len(blocks) - 1

	return blocks[previousItem].Kind == itemType
}

func appendNonEmptyBuffer(temp *tempBuffer, blocks []*markdown.Block) ([]*markdown.Block, error) {
	if temp.kind == markdown.Text && temp.text == "" {
		return nil, errors.New("buffer is empty")
	}

	b := markdown.Block{
		Kind: temp.kind,
		Text: temp.text,
	}

	return append(blocks, &b), nil
}

type tempBuffer struct {
	kind int
	text string
}

func (b *tempBuffer) reset() {
	b.kind = 0
	b.text = ""
}

func (b *tempBuffer) append(text string) {
	b.text += text
}
