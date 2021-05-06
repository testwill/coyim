package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldListSuite struct{}

var _ = Suite(&MucRoomConfigFieldListSuite{})

func (*MucRoomConfigFieldListSuite) Test_newRoomConfigFieldListValue(c *C) {
	cases := []struct {
		values   []string
		expected string
	}{
		{
			[]string{},
			"",
		},
		{
			[]string{"whatever"},
			"whatever",
		},
		{
			[]string{"bla", "foo"},
			"bla",
		},
		{
			[]string{"foo", "bla"},
			"foo",
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldListValue(mock.values, []string{})
		c.Assert(field.Selected(), Equals, mock.expected)
		c.Assert(field.Value(), DeepEquals, []string{mock.expected})
	}
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldListValue_SetValue(c *C) {
	field := newRoomConfigFieldListValue([]string{"bla", "foo"}, []string{})
	c.Assert(field.Selected(), DeepEquals, "bla")
	c.Assert(field.Value(), DeepEquals, []string{"bla"})

	field.SetValue("foo")
	c.Assert(field.Selected(), DeepEquals, "foo")
	c.Assert(field.Value(), DeepEquals, []string{"foo"})

	field.SetValue("whatever")
	c.Assert(field.Selected(), DeepEquals, "whatever")
	c.Assert(field.Value(), DeepEquals, []string{"whatever"})

	field.SetValue("")
	c.Assert(field.Selected(), DeepEquals, "")
	c.Assert(field.Value(), DeepEquals, []string{""})

	field.SetValue(1000)
	c.Assert(field.Selected(), DeepEquals, "")
	c.Assert(field.Value(), DeepEquals, []string{""})

	field.SetValue("abc")
	c.Assert(field.Selected(), DeepEquals, "abc")
	c.Assert(field.Value(), DeepEquals, []string{"abc"})
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldListValue_Options(c *C) {
	field := newRoomConfigFieldListValue(nil, []string{"bla", "foo"})
	c.Assert(field.Options(), DeepEquals, []string{"bla", "foo"})

	field.SetOptions(nil)
	c.Assert(field.Options(), DeepEquals, []string(nil))

	field.SetOptions([]string{"foo"})
	c.Assert(field.Options(), DeepEquals, []string{"foo"})

	field.SetOptions([]string{"whatever"})
	c.Assert(field.Options(), DeepEquals, []string{"whatever"})
}