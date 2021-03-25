package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAssistantSidebar struct {
	assistant *roomConfigAssistant
	content   gtki.Box `gtk-widget:"assistant-options"`
}

func (rc *roomConfigAssistant) newRoomConfigAssistantSidebar() *roomConfigAssistantSidebar {
	sb := &roomConfigAssistantSidebar{
		assistant: rc,
	}

	sb.initBuilder()

	return sb
}

func (sb *roomConfigAssistantSidebar) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistantSidebar")
	panicOnDevError(b.bindObjects(sb))

	b.ConnectSignals(map[string]interface{}{
		"row_selected": sb.rowSelected,
	})
}

func (sb *roomConfigAssistantSidebar) rowSelected(l gtki.ListBox, r gtki.ListBoxRow) {
	// TODO: update assistant page
}

func (sb *roomConfigAssistantSidebar) getContent() gtki.Box {
	return sb.content
}
