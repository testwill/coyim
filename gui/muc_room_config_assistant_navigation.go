package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

const noRowIndex = -1

type roomConfigAssistantNavigation struct {
	assistant *roomConfigAssistant

	content    gtki.Box     `gtk-widget:"room-config-assistant-navigation-content"`
	navigation gtki.ListBox `gtk-widget:"room-config-assistant-navigation-list"`
}

func (rc *roomConfigAssistant) newRoomConfigAssistantNavigation() *roomConfigAssistantNavigation {
	rcn := &roomConfigAssistantNavigation{
		assistant: rc,
	}

	rcn.initBuilder()

	return rcn
}

func (rcn *roomConfigAssistantNavigation) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistantNavigation")
	panicOnDevError(b.bindObjects(rcn))

	b.ConnectSignals(map[string]interface{}{
		"row_selected": rcn.onRowSelected,
	})
}

// onRowSelected MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) onRowSelected(_ gtki.ListBox, r gtki.ListBoxRow) {
	rcn.assistant.updateAssistantPage(r.GetIndex())
}

// selectOptionByIndex MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) selectOptionByIndex(idx int) {
	row := rcn.navigation.GetRowAtIndex(idx)
	rowIndex := getListBoxRowIndex(row)
	currentRowIndex := getListBoxRowIndex(rcn.navigation.GetSelectedRow())

	if rowIndex != noRowIndex && rowIndex != currentRowIndex {
		rcn.navigation.SelectRow(row)
	}
}

func getListBoxRowIndex(r gtki.ListBoxRow) int {
	if r != nil {
		return r.GetIndex()
	}
	return noRowIndex
}
