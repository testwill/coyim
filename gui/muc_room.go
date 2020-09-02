package gui

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomView struct {
	builder *builder
	u       *gtkUI

	log              coylog.Logger
	account          *account
	jid              jid.Bare
	onJoin           chan bool
	lastError        error
	lastErrorMessage string
	sync.RWMutex

	window           gtki.Window      `gtk-widget:"roomWindow"`
	boxJoinRoomView  gtki.Box         `gtk-widget:"boxJoinRoomView"`
	nicknameEntry    gtki.Entry       `gtk-widget:"nicknameEntry"`
	passwordCheck    gtki.CheckButton `gtk-widget:"passwordCheck"`
	passwordLabel    gtki.Label       `gtk-widget:"passwordLabel"`
	passwordEntry    gtki.Entry       `gtk-widget:"passwordEntry"`
	roomJoinButton   gtki.Button      `gtk-widget:"roomJoinButton"`
	spinnerJoinView  gtki.Spinner     `gtk-widget:"joinSpinner"`
	notificationArea gtki.Box         `gtk-widget:"boxNotificationArea"`

	errorNotif   *errorNotification
	notification gtki.InfoBar

	boxRoomView        gtki.Box        `gtk-widget:"boxRoomView"`
	roomChatTextBuffer gtki.TextBuffer `gtk-widget:"roomChatTextBuffer"`
	rosterPanel        gtki.Box        `gtk-widget:"panel"`
	panelToggle        gtki.Button     `gtk-widget:"panel-toggle"`
	membersModel       gtki.ListStore  `gtk-widget:"room-members-model"`
	membersView        gtki.TreeView   `gtk-widget:"room-members-tree"`
}

func (v *roomView) clearErrors() {
	v.errorNotif.Hide()
}

func (v *roomView) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}

	v.errorNotif.ShowMessage(err)
}

func (v *roomView) initUIBuilder() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.builder.ConnectSignals(map[string]interface{}{
		"on_show_window":         v.validateInput,
		"on_nickname_changed":    v.validateInput,
		"on_password_changed":    v.validateInput,
		"on_password_checked":    v.onPasswordChecked,
		"on_room_cancel_clicked": v.window.Destroy,
		"on_room_join_clicked":   v.joinRoom,
		"on_close_window":        v.onCloseWindow,
		"on_toggle_roster_panel": v.toggleRosterPanel,
	})
}

func (v *roomView) onPasswordChecked() {
	v.setPasswordSensitiveBasedOnCheck()
	v.validateInput()
}

func (v *roomView) onCloseWindow() {
	_ = v.account.roomManager.LeaveRoom(v.jid)
}

func (v *roomView) initDefaults() {
	v.log = v.account.log.WithField("room", v.jid)
	v.errorNotif = newErrorNotification(v.notificationArea)
	v.setPasswordSensitiveBasedOnCheck()
	v.window.SetTitle(i18n.Localf("Room: [%s]", v.jid))
}

func (v *roomView) setPasswordSensitiveBasedOnCheck() {
	a := v.passwordCheck.GetActive()
	v.passwordLabel.SetSensitive(a)
	v.passwordEntry.SetSensitive(a)
}

func (v *roomView) hasValidNickname() bool {
	nickname, _ := v.nicknameEntry.GetText()
	return len(nickname) > 0
}

func (v *roomView) hasValidPassword() bool {
	cv := v.passwordCheck.GetActive()
	if !cv {
		return true
	}
	password, _ := v.passwordEntry.GetText()
	return len(password) > 0
}

func (v *roomView) validateInput() {
	v.clearErrors()
	sensitiveValue := v.hasValidNickname() && v.hasValidPassword()
	v.roomJoinButton.SetSensitive(sensitiveValue)
}

func (v *roomView) togglePanelView() {
	doInUIThread(func() {
		value := v.boxJoinRoomView.IsVisible()
		v.boxJoinRoomView.SetVisible(!value)
		v.boxRoomView.SetVisible(value)
	})
}

func (v *roomView) toggleRosterPanel() {
	iv := v.rosterPanel.IsVisible()
	v.rosterPanel.SetVisible(!iv)
	if !iv {
		v.panelToggle.SetProperty("label", i18n.Local("Hide panel"))
		return
	}

	v.panelToggle.SetProperty("label", i18n.Local("Show panel"))
}

// startSpinner should be called from UI thread
func (v *roomView) startSpinner() {
	v.spinnerJoinView.Start()
	v.spinnerJoinView.SetVisible(true)
	v.roomJoinButton.SetSensitive(false)
}

// stopSpinner should be called from UI thread
func (v *roomView) stopSpinner() {
	v.spinnerJoinView.Stop()
	v.spinnerJoinView.SetVisible(false)
	if !v.errorNotif.IsVisible() {
		v.roomJoinButton.SetSensitive(true)
	}
}

func (v *roomView) joinRoomWithNickname(nickname string) {
	v.log.WithField("nickname", nickname).Debug("joinRoomWithNickname()")

	doInUIThread(v.startSpinner)

	go func() {
		err := v.account.session.JoinRoom(v.jid, nickname)
		if err != nil {
			doInUIThread(func() {
				v.stopSpinner()
				v.log.WithField("nickname", nickname).WithError(err).Error("An error occurred while trying to join the room.")
			})
		}
	}()

	go v.whenJoinRoomFinishes(nickname)
}

func (v *roomView) whenJoinRoomFinishes(nickname string) {
	defer func() {
		doInUIThread(v.stopSpinner)
	}()

	hasJoined, ok := <-v.onJoin
	if !ok {
		doInUIThread(func() {
			v.lastErrorMessage = i18n.Local("An error happened while trying to join the room, please check your connection or try again.")
			v.notifyOnError(v.lastErrorMessage)
		})
		return
	}

	if !hasJoined {
		// TODO: We should do the better for the user, if the room doesn't exists maybe we should
		// allow the user to create the room or tell him something to try as solution
		if v.lastErrorMessage == "" {
			v.lastErrorMessage = i18n.Local("An error happened while trying to join the room, please check your connection or make sure the room exists.")
		}

		v.log.WithFields(log.Fields{
			"nickname": nickname,
			"message":  v.lastErrorMessage,
		}).Error("An error happened while trying to join the room")

		doInUIThread(func() {
			v.notifyOnError(v.lastErrorMessage)
		})

		return
	}

	doInUIThread(func() {
		v.clearErrors()
		v.togglePanelView()
	})
}

func (v *roomView) joinRoom() {
	v.clearErrors()

	v.onJoin = make(chan bool)
	nickname, _ := v.nicknameEntry.GetText()

	go v.joinRoomWithNickname(nickname)
}

func (u *gtkUI) newRoom(a *account, ident jid.Bare) *muc.Room {
	room := muc.NewRoom(ident)

	view := &roomView{
		account: a,
		jid:     ident,
		u:       u,
	}

	view.initUIBuilder()
	view.initDefaults()

	room.Opaque = view

	return room

}

func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare) {
	room, ok := a.roomManager.GetRoom(ident)
	if !ok {
		room = u.newRoom(a, ident)
		a.roomManager.AddRoom(room)
	}

	view := getViewFromRoom(room)

	if !ok {
		view.window.Show()
		return
	}

	view.window.Present()
}

func getViewFromRoom(r *muc.Room) *roomView {
	return r.Opaque.(*roomView)
}

func (v *roomView) addLineToChatText(text string) {
	i := v.roomChatTextBuffer.GetEndIter()

	t := fmt.Sprintf("%s\n", text)
	v.roomChatTextBuffer.Insert(i, t)
}

func (v *roomView) showOccupantLeftRoom(nickname jid.Resource) {
	doInUIThread(func() {
		v.addLineToChatText(i18n.Localf("%s left the room", nickname))
	})
}

func (v *roomView) updateOccupantsInModel(occupants []*muc.Occupant) {
	doInUIThread(func() {
		v.membersModel.Clear()
		for _, o := range occupants {
			iter := v.membersModel.Append()
			_ = v.membersModel.SetValue(iter, 0, v.getIconBaseOnVoice(o.Role).GetPixbuf())
			_ = v.membersModel.SetValue(iter, 1, o.Nick)
			_ = v.membersModel.SetValue(iter, 2, v.getAffiliationForRosterPanel(o.Affiliation))
			_ = v.membersModel.SetValue(iter, 3, getRoleNameForTooltip(o.Role))
		}
	})
}

func getRoleNameForTooltip(r muc.Role) string {
	switch r.Name() {
	case muc.RoleNone:
		return i18n.Local("Role: None")
	case muc.RoleParticipant:
		return i18n.Local("Role: Participant")
	case muc.RoleVisitor:
		return i18n.Local("Role: Visitor")
	case muc.RoleModerator:
		return i18n.Local("Role: Moderator")
	default:
		return ""
	}
}

func (v *roomView) getAffiliationForRosterPanel(a muc.Affiliation) string {
	switch a.Name() {
	case muc.AffiliationAdmin:
		return i18n.Local("Admin")
	case muc.AffiliationOwner:
		return i18n.Local("Owner")
	default:
		return ""
	}
}

func (v *roomView) getIconBaseOnVoice(r muc.Role) Icon {
	if r.HasVoice() {
		return statusIcons["available"]
	}
	return statusIcons["offline"]
}
