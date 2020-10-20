package session

import (
	"time"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (m *mucManager) receiveClientMessage(stanza *xmppData.ClientMessage) {
	m.log.WithField("stanza", stanza).Debug("handleMUCReceivedClientMessage()")

	if hasSubject(stanza) {
		m.handleSubjectReceived(stanza)
	}

	switch {
	case isDelayedMessage(stanza):
		m.handleMessageReceived(stanza, m.receiveDelayedMessage)
	case isLiveMessage(stanza):
		m.handleMessageReceived(stanza, m.liveMessageReceived)
	case hasMucUserExtension(stanza):
		m.handleMucUserExtension(stanza)
	}
}

func (m *mucManager) receiveDelayedMessage(roomID jid.Bare, nickname, message string, timestamp time.Time) {
	dh, ok := m.discussionHistory[roomID]
	if !ok {
		dh = data.NewDiscussionHistory()
		m.discussionHistory[roomID] = dh
	}

	dh.AddMessage(nickname, message, timestamp)
}

func (m *mucManager) handleSubjectReceived(stanza *xmppData.ClientMessage) {
	l := m.log.WithFields(log.Fields{
		"from": stanza.From,
		"who":  "handleSubjectReceived",
	})

	roomID, ok := jid.TryParseBare(stanza.From)
	if !ok {
		l.Error("Error trying to get the room ID from stanza")
		return
	}

	dh := m.getDiscussionHistory(roomID)
	if dh != nil {
		m.discussionHistoryReceived(roomID, dh)
	}

	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.WithField("room", roomID).Error("Error trying to read the subject of room")
		return
	}

	s := getSubjectFromStanza(stanza)
	updated := room.UpdateSubject(s)
	if updated {
		m.subjectUpdated(roomID, getNicknameFromStanza(stanza), s)
		return
	}

	m.subjectReceived(roomID, s)
}

func (m *mucManager) handleMessageReceived(stanza *xmppData.ClientMessage, h func(jid.Bare, string, string, time.Time)) {
	roomID, nickname := retrieveRoomIDAndNickname(stanza.From)
	h(roomID, nickname, stanza.Body, retrieveMessageTime(stanza))
}

func (m *mucManager) handleMucUserExtension(stanza *xmppData.ClientMessage) {
	for _, status := range stanza.MucUserExtension.Status {
		switch status.Code {
		// TODO: Implements the publishing of events for states 170, 171, 172, and 173
		case "104":
			m.handleMucRoomConfigurationChanged(stanza)
			m.log.WithField("status code:", status.Code).Info("Room changes published")
		default:
			m.log.WithField("status code:", status.Code).Warn("Unknown status code received")
		}
	}
}

func (m *mucManager) handleMucRoomConfigurationChanged(stanza *xmppData.ClientMessage) {
	l := m.log.WithFields(log.Fields{
		"stanza": stanza,
		"who":    "handleMucUserExtension",
	})

	roomID, ok := jid.TryParseBare(stanza.From)
	if !ok {
		l.Error("Error trying to get the room ID from stanza")
		return
	}

	roomInfo := make(chan *muc.RoomListing)
	go m.getRoom(roomID, roomInfo)

	ri := <-roomInfo
	m.roomConfigurationChanged(roomID, ri)
}

func bodyHasContent(stanza *xmppData.ClientMessage) bool {
	return stanza.Body != ""
}

func isDelayedMessage(stanza *xmppData.ClientMessage) bool {
	return stanza.Delay != nil
}

func isLiveMessage(stanza *xmppData.ClientMessage) bool {
	return bodyHasContent(stanza) && !isDelayedMessage(stanza)
}

func hasSubject(stanza *xmppData.ClientMessage) bool {
	return stanza.Subject != nil
}

func hasMucUserExtension(stanza *xmppData.ClientMessage) bool {
	return stanza.MucUserExtension != nil
}

func getNicknameFromStanza(stanza *xmppData.ClientMessage) string {
	from, ok := jid.TryParseFull(stanza.From)
	if ok {
		return from.Resource().String()
	}

	return ""
}

func getSubjectFromStanza(stanza *xmppData.ClientMessage) string {
	if hasSubject(stanza) {
		return stanza.Subject.Text
	}

	return ""
}
