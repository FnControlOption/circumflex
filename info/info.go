package info

import (
	"strings"

	"clx/constants/nerdfonts"
	. "github.com/logrusorgru/aurora/v3"

	"clx/constants/margins"
	"clx/constants/style"
	"clx/keymaps"
	"github.com/charmbracelet/lipgloss"

	text "github.com/MichaelMure/go-term-text"
)

func GetText(screenWidth int, enableNerdFonts bool) string {
	keys := new(keymaps.List)
	keys.Init()

	cmd := lipgloss.NewStyle().Background(lipgloss.Color("237")).Italic(true)

	keys.AddHeader(cmd.Foreground(style.GetMagenta()).Render(" circumflex "))
	keys.AddSeparator()
	keys.AddKeymap("Read comment section", "Enter")
	keys.AddKeymap("Read article in Reader Mode", "Space")
	keys.AddSeparator()
	keys.AddKeymap("Refresh", "r")
	keys.AddKeymap("Change category", "Tab")
	keys.AddSeparator()
	keys.AddKeymap("Open story link in browser", "o")
	keys.AddKeymap("Open comments in browser", "c")
	keys.AddSeparator()
	keys.AddKeymap("Add to favorites", "f")
	keys.AddKeymap("Remove from favorites", "x")
	keys.AddSeparator()
	keys.AddKeymap("Bring up this screen", "i, ?")
	keys.AddKeymap("Quit to prompt", "q")
	keys.AddSeparator()

	keys.AddHeader(cmd.Foreground(style.GetYellow()).Render(" less "))
	keys.AddSeparator()
	keys.AddKeymap("Down one half-window", "d")
	keys.AddKeymap("Up one half-window", "u")
	keys.AddSeparator()
	keys.AddKeymap("Jump to next top-level comment", "n")
	keys.AddKeymap("Jump to previous top-level comment", "N")
	keys.AddSeparator()
	keys.AddKeymap("Help screen", "h")
	keys.AddKeymap("Return to circumflex", "q")
	keys.AddSeparator()

	keys.AddHeader(cmd.Foreground(style.GetBlue()).Render(" legend "))
	keys.AddSeparator()
	keys.AddKeymap("Original Poster", Red(getOP(enableNerdFonts)).String())
	keys.AddKeymap("Parent Poster", Magenta(getPP(enableNerdFonts)).String())
	keys.AddKeymap("Moderator", Green(getMod(enableNerdFonts)).String())

	keymapsWidth := 80
	listOfKeymaps := keys.Print(keymapsWidth)
	listOfKeymapsCentered := alignCenter(listOfKeymaps, screenWidth, keymapsWidth)

	return listOfKeymapsCentered
}

func alignCenter(input string, screenWidth int, keymapsWidth int) string {
	padding := screenWidth/2 - keymapsWidth/2 - margins.MainViewLeftMargin

	if padding < 0 {
		return input
	}

	padToCenterAlign := strings.Repeat(" ", padding)
	output, _ := text.WrapWithPad(input, screenWidth, padToCenterAlign)

	return output
}

func getOP(enableNerdFonts bool) string {
	if enableNerdFonts {
		return nerdfonts.Author
	}

	return "OP"
}

func getPP(enableNerdFonts bool) string {
	if enableNerdFonts {
		return nerdfonts.Author
	}

	return "PP"
}

func getMod(enableNerdFonts bool) string {
	if enableNerdFonts {
		return nerdfonts.Author
	}

	return "mod"
}
