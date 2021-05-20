package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryFieldContainer struct {
	fields []*roomConfigSummaryField

	widget  gtki.Box     `gtk-widget:"room-config-field-box"`
	content gtki.ListBox `gtk-widget:"room-config-fields-content"`
}

func newRoomConfigSummaryFieldContainer(f []*roomConfigSummaryField) hasRoomConfigFormField {
	field := &roomConfigSummaryFieldContainer{
		fields: f,
	}

	field.initBuilder()
	field.initDefaults()

	return field
}

func (fc *roomConfigSummaryFieldContainer) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryFieldContainer")
	panicOnDevError(builder.bindObjects(fc))
}

func (fc *roomConfigSummaryFieldContainer) initDefaults() {
	fc.content.Add(fc.fields[0].widget)
	for _, f := range fc.fields[1:] {
		fc.content.Add(createSeparator(gtki.HorizontalOrientation))
		fc.content.Add(f.widget)
	}
}

func (fc *roomConfigSummaryFieldContainer) fieldWidget() gtki.Widget {
	return fc.widget
}

// refreshContent MUST NOT be called from the UI thread
func (fc *roomConfigSummaryFieldContainer) refreshContent() {}

// collectFieldValue MUST be called from the UI thread
func (fc *roomConfigSummaryFieldContainer) collectFieldValue() {
	for _, f := range fc.fields {
		f.collectFieldValue()
	}
}

// isValid implements the hasRoomConfigFormField interface
func (fc *roomConfigSummaryFieldContainer) isValid() bool {
	return true
}

// showValidationErrors implements the hasRoomConfigFormField interface
func (fc *roomConfigSummaryFieldContainer) showValidationErrors() {}