package muc

// RoomConfigFieldListValue contains information of the value of the list single field
type RoomConfigFieldListValue struct {
	value   string
	options []string
}

func newRoomConfigFieldListValue(values, options []string) *RoomConfigFieldListValue {
	return &RoomConfigFieldListValue{formFieldSingleString(values), options}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListValue) Value() []string {
	return []string{v.value}
}

// SetSelected sets the current selected value
func (v *RoomConfigFieldListValue) SetSelected(s string) {
	v.value = s
}

// Selected returns the current selected value
func (v *RoomConfigFieldListValue) Selected() string {
	return v.value
}

// Options returns the available options for the field
func (v *RoomConfigFieldListValue) Options() []string {
	return v.options
}

// SetOptions updates the options for the field
func (v *RoomConfigFieldListValue) SetOptions(options []string) {
	v.options = options
}
