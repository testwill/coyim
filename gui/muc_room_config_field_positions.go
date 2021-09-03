package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

const positionsListJidColumnIndex = 0

type roomConfigPositionsOptions struct {
	affiliation            data.Affiliation
	occupantList           muc.RoomOccupantItemList
	setOccupantList        func(muc.RoomOccupantItemList) // setOccupantList WILL be called from the UI thread
	setRemovedOccupantList func(muc.RoomOccupantItemList) // setRemovedOccupantList WILL be called from the UI thread
	displayErrors          func()                         // displayErrors WILL be called from the UI thread
	parentWindow           gtki.Window
}

type roomConfigFieldPositions struct {
	builder                   *builder
	affiliation               data.Affiliation
	originalOccupantsList     muc.RoomOccupantItemList
	setOccupantList           func(occupants muc.RoomOccupantItemList)
	updateRemovedOccupantList func(occupantsToRemove muc.RoomOccupantItemList)
	showErrorNotification     func()

	content               gtki.Box      `gtk-widget:"room-config-positions-content"`
	header                gtki.Label    `gtk-widget:"room-config-position-header"`
	description           gtki.Label    `gtk-widget:"room-config-position-description"`
	positionsListContent  gtki.Box      `gtk-widget:"room-config-positions-list-content"`
	positionsList         gtki.TreeView `gtk-widget:"room-config-positions-list"`
	positionsAddButton    gtki.Button   `gtk-widget:"room-config-position-add"`
	positionsRemoveButton gtki.Button   `gtk-widget:"room-config-position-remove"`
	positionsRemoveLabel  gtki.Label    `gtk-widget:"room-config-position-remove-label"`

	positionsListController *mucRoomConfigListController
	onListChanged           *callbacksSet
}

func newRoomConfigFieldPositions(options roomConfigPositionsOptions) *roomConfigFieldPositions {
	rcp := &roomConfigFieldPositions{
		affiliation:               options.affiliation,
		originalOccupantsList:     options.occupantList,
		setOccupantList:           options.setOccupantList,
		updateRemovedOccupantList: options.setRemovedOccupantList,
		showErrorNotification:     options.displayErrors,
		onListChanged:             newCallbacksSet(),
	}

	rcp.loadUIDefinition()
	rcp.initDefaults()

	return rcp
}

func (p *roomConfigFieldPositions) setUIBuilder(b *builder) {
	p.builder = b
}

func (p *roomConfigFieldPositions) connectUISignals(b *builder) {
	b.ConnectSignals(map[string]interface{}{
		"on_jid_edited": p.onOccupantJidEdited,
	})
}

func (p *roomConfigFieldPositions) loadUIDefinition() {
	buildUserInterface("MUCRoomConfigFieldPositions", p, p.setUIBuilder, p.connectUISignals)
}

func (p *roomConfigFieldPositions) initDefaults() {
	p.initPositionLabels()
	mucStyles.setHelpTextStyle(p.content)
}

func (p *roomConfigFieldPositions) initPositionLabels() {
	p.header.SetText(getFieldTextByAffiliation(p.affiliation).headerLabel)
	p.description.SetText(getFieldTextByAffiliation(p.affiliation).descriptionLabel)
}

// addItemsToListController MUST be called from the UI thread
func (p *roomConfigFieldPositions) addItemsToListController() {
	jids := []string{}
	for _, o := range p.originalOccupantsList {
		jids = append(jids, o.Jid.String())
	}
	p.positionsListController.listComponent.addListItems(jids)
}

func (p *roomConfigFieldPositions) refreshContentLists() {
	p.positionsListContent.SetVisible(p.positionsListController.hasItems())
}

func (p *roomConfigFieldPositions) onOccupantJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	p.updateOccupantListCellForString(p.positionsListController, positionsListJidColumnIndex, path, newValue)
	p.onListChanged.invokeAll()
}

func (p *roomConfigFieldPositions) updateOccupantListCellForString(controller *mucRoomConfigListController, column int, path string, newValue string) {
	if controller.updateCellForString(column, path, newValue) {
		log.WithFields(log.Fields{
			"path":        path,
			"newText":     newValue,
			"affiliation": p.affiliation.Name(),
		}).Debug("The occupant's jid can't be updated")
	}
}

func (p *roomConfigFieldPositions) updateFieldValue() {
	p.refreshOccupantLists(p.currentOccupantList())
}

func (p *roomConfigFieldPositions) refreshOccupantLists(currentList muc.RoomOccupantItemList) {
	occupantsList := muc.RoomOccupantItemList{}
	for _, oi := range currentList {
		oi.MustBeUpdated = p.isNewOccupant(oi)
		occupantsList = append(occupantsList, oi)
	}
	p.setOccupantList(occupantsList)

	deletedOccupantsList := muc.RoomOccupantItemList{}
	for _, oi := range p.originalOccupantsList {
		if !currentList.IncludesJid(oi.Jid) {
			oi.ChangeAffiliationToNone()
			oi.MustBeUpdated = true
			deletedOccupantsList = append(deletedOccupantsList, oi)
		}
	}
	p.updateRemovedOccupantList(deletedOccupantsList)
}

func (p *roomConfigFieldPositions) isNewOccupant(o *muc.RoomOccupantItem) bool {
	return !p.originalOccupantsList.IncludesJid(o.Jid)
}

func (p *roomConfigFieldPositions) currentOccupantList() muc.RoomOccupantItemList {
	positionsList := []*muc.RoomOccupantItem{}
	for _, item := range p.positionsListController.listItems() {
		positionsList = append(positionsList, &muc.RoomOccupantItem{
			Jid:         jid.Parse(item),
			Affiliation: p.affiliation,
		})
	}
	return positionsList
}

func (p *roomConfigFieldPositions) showValidationErrors() {
	if p.showErrorNotification != nil {
		p.showErrorNotification()
	}
}

func (p *roomConfigFieldPositions) fieldWidget() gtki.Widget {
	return p.content
}

func (p *roomConfigFieldPositions) refreshContent() {
	p.refreshContentLists()
}

func (p *roomConfigFieldPositions) isValid() bool {
	return !(p.affiliation.IsOwner() && len(p.originalOccupantsList) != 0 && len(p.currentOccupantList()) == 0)
}

func (p *roomConfigFieldPositions) hasListChanged() bool {
	ol := append(muc.RoomOccupantItemList{}, p.originalOccupantsList...)
	cl := append(muc.RoomOccupantItemList{}, p.currentOccupantList()...)

	if len(ol) != len(cl) {
		return true
	}

	for _, i := range cl {
		if !ol.IncludesJid(i.Jid) {
			return true
		}
	}

	return false
}

// fieldKey implements the hasRoomConfigFormField interface
func (p *roomConfigFieldPositions) fieldKey() muc.RoomConfigFieldType {
	return muc.RoomConfigFieldUnexpected
}
