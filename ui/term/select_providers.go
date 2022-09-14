package term

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
)

func showProviders(providers []string, selectProviders []string) []string {
	var res []string
	for _, provider := range providers {
		flag := false
		for _, selectProvider := range selectProviders {
			if selectProvider == provider {
				flag = true
				continue
			}
		}
		if flag {
			res = append(res, "[*] "+provider)
		} else {
			res = append(res, "[ ] "+provider)
		}
	}
	return res

}

func SelectProviders(providers []string) []string {

	selectProviders := []string{}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Rows = showProviders(providers, selectProviders)
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.Title = "[Use arrows to move, Space to select, Enter to complete the selection]"
	l.BorderLeft = false
	l.BorderRight = false
	l.BorderTop = false
	l.BorderBottom = false
	l.SelectedRowStyle = ui.NewStyle(ui.ColorRed)
	l.SetRect(0, 0, 200, 10)

	ui.Render(l)

	previousKey := ""
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "j", "<Down>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollDown()
		case "k", "<Up>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollUp()
		case "<C-d>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollHalfPageDown()
		case "<C-u>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollHalfPageUp()
		case "<C-f>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollPageDown()
		case "<C-b>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollPageUp()
		case "g":
			if len(l.Rows) == 0 {
				continue
			}
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Enter>":
			return selectProviders
		case "<Space>":
			if len(l.Rows) == 0 {
				continue
			}
			flag := -1
			for i, provider := range selectProviders {
				if providers[l.SelectedRow] == provider {
					flag = i
					break
				}
			}
			if flag > -1 {
				selectProviders = append(selectProviders[:flag], selectProviders[flag+1:]...)
			} else {
				selectProviders = append(selectProviders, providers[l.SelectedRow])
			}
			l.Rows = showProviders(providers, selectProviders)
		case "<Home>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollTop()
		case "G", "<End>":
			if len(l.Rows) == 0 {
				continue
			}
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(l)

	}
}
