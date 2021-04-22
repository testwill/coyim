package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

func (v *roomView) onModifyBanList() {
	bl := v.newRoomBanListView()
	bl.show()
}

const (
	roomBanListAccountIndex int = iota
	roomBanListAffiliationIndex
	roomBanListReasonIndex
)

type roomBanListView struct {
	roomView *roomView

	dialog             gtki.Window        `gtk-widget:"ban-list-window"`
	addEntryButton     gtki.Button        `gtk-widget:"ban-list-add-entry-button"`
	removeEntryButton  gtki.Button        `gtk-widget:"ban-list-remove-entry-button"`
	list               gtki.TreeView      `gtk-widget:"ban-list-treeview"`
	listSelection      gtki.TreeSelection `gtk-widget:"ban-list-treeview-selection"`
	listContent        gtki.Overlay       `gtk-widget:"ban-list-content-overlay"`
	listView           gtki.Overlay       `gtk-widget:"ban-list-overlay-view"`
	listLoadingView    gtki.Box           `gtk-widget:"ban-list-loading-view"`
	noEntriesView      gtki.Box           `gtk-widget:"ban-list-no-entries-view"`
	noEntriesErrorView gtki.Box           `gtk-widget:"ban-list-error-view"`
	applyButton        gtki.Button        `gtk-widget:"ban-list-apply-changes-button"`

	listModel       gtki.ListStore
	originalBanList []*muc.RoomBanListItem
	cancelChannel   chan bool

	log coylog.Logger
}

func (v *roomView) newRoomBanListView() *roomBanListView {
	bl := &roomBanListView{
		roomView: v,
		log:      v.log.WithField("where", "roomBanListView"),
	}

	bl.initBuilder()
	bl.initDefaults()
	bl.initBanListModel()

	return bl
}

func (bl *roomBanListView) initBuilder() {
	builder := newBuilder("MUCRoomBannedUsersDialog")
	panicOnDevError(builder.bindObjects(bl))

	builder.ConnectSignals(map[string]interface{}{
		"on_no_items_try_again_clicked": bl.requestBanListAgain,
		"on_error_try_again_clicked":    bl.requestBanListAgain,
		"on_jid_edited":                 bl.onUserJidEdited,
		"on_reason_edited":              bl.onReasonEdited,
		"on_add_item":                   bl.onAddNewItem,
		"on_selection_changed":          bl.onSelectionChanged,
		"on_remove_item":                bl.onRemoveItem,
		"on_cancel_clicked":             bl.onCancel,
	})
}

func (bl *roomBanListView) initDefaults() {
	bl.dialog.SetTransientFor(bl.roomView.mainWindow())
	bl.disableButtonsAndInteractions()
}

func (bl *roomBanListView) initBanListModel() {
	model, _ := g.gtk.ListStoreNew(
		// the user's jid
		glibi.TYPE_STRING,
		// the user's affiliation
		glibi.TYPE_STRING,
		// the reason
		glibi.TYPE_STRING,
	)

	bl.listModel = model
	bl.list.SetModel(bl.listModel)
}

// addListItem MUST be called from the UI thread
func (bl *roomBanListView) addListItem(itm *muc.RoomBanListItem) gtki.TreeIter {
	iter := bl.listModel.Append()

	jid := ""
	if itm.Jid != nil {
		jid = fmt.Sprintf("%s", itm.Jid)
	}

	affiliation := ""
	if itm.Affiliation != nil {
		affiliation = affiliationDisplayName(itm.Affiliation)
	}

	bl.listModel.SetValue(iter, roomBanListAccountIndex, jid)
	bl.listModel.SetValue(iter, roomBanListAffiliationIndex, affiliation)
	bl.listModel.SetValue(iter, roomBanListReasonIndex, itm.Reason)

	return iter
}

// show MUST be called from the UI thread
func (bl *roomBanListView) show() {
	bl.refreshBanList()
	bl.dialog.Show()
}

// refreshBanList MUST be called from the UI thread
func (bl *roomBanListView) refreshBanList() {
	bl.listModel.Clear()
	bl.disableButtonsAndInteractions()
	bl.showLoadingAndListViews()

	go bl.requestBanList()
}

// requestBanListAgain MUST be called from the UI thread
func (bl *roomBanListView) requestBanListAgain() {
	bl.showLoadingAndListViews()
	bl.refreshBanList()
}

// showLoadingAndListViews MUST be called from the UI thread
func (bl *roomBanListView) showLoadingAndListViews() {
	bl.noEntriesView.Hide()
	bl.noEntriesErrorView.Hide()

	bl.listLoadingView.Show()
	bl.listView.Show()
}

// hideLoadingAndListViews MUST be called from the UI thread
func (bl *roomBanListView) hideLoadingAndListViews() {
	bl.listLoadingView.Hide()
	bl.listView.Hide()
}

// hasItems MUST be called from the UI thread
func (bl *roomBanListView) hasItems() bool {
	_, ok := bl.listModel.GetIterFirst()
	return ok
}

// requestBanList MUST NOT be called from the UI thread
func (bl *roomBanListView) requestBanList() {
	bl.cancelChannel = make(chan bool)

	blc, ec := bl.roomView.account.session.GetRoomBanList(bl.roomView.roomID())

	go func() {
		defer func() {
			bl.cancelChannel = nil
		}()

		select {
		case items := <-blc:
			bl.onRequestFinish(items)
		case err := <-ec:
			bl.onRequestError(err)
		case <-bl.cancelChannel:
		}
	}()
}

// onRequestFinish MUST NOT be called from the UI thread
func (bl *roomBanListView) onRequestFinish(items []*muc.RoomBanListItem) {
	if len(items) > 0 {
		doInUIThread(func() {
			for _, itm := range items {
				_ = bl.addListItem(itm)
			}
		})
	} else {
		doInUIThread(func() {
			bl.hideLoadingAndListViews()
			bl.noEntriesView.Show()
		})
	}

	bl.addEntryButton.SetSensitive(true)
	bl.originalBanList = items
}

// onRequestError MUST NOT be called from the UI thread
func (bl *roomBanListView) onRequestError(err error) {
	bl.roomView.log.WithError(err).Error("Something happened when requesting the banned users list")

	doInUIThread(func() {
		bl.hideLoadingAndListViews()
		bl.noEntriesErrorView.Show()
	})
}

// onUserJidEdited MUST be called from the UI thread
func (bl *roomBanListView) onUserJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	iter, err := bl.listModel.GetIterFromString(path)
	if err != nil {
		bl.log.WithFields(log.Fields{
			"path":     path,
			"newValue": newValue,
		}).WithError(err).Error("Can't get the iter to update the jid of the banned user")
		return
	}

	if newValue != "" && !jid.Parse(newValue).Valid() {
		bl.log.WithFields(log.Fields{
			"path":        path,
			"newJidValue": newValue,
		}).Error("Can't update the jid of the banned user to an invalid value")
		return
	}

	if err = bl.listModel.SetValue(iter, roomBanListAccountIndex, newValue); err != nil {
		bl.log.WithFields(log.Fields{
			"path":        path,
			"newJidValue": newValue,
		}).WithError(err).Error("Can't set the new value for the jid of the banned user")
	}

	bl.enableApplyIfConditionsAreMet()
}

// onReasonEdited MUST be called from the UI thread
func (bl *roomBanListView) onReasonEdited(_ gtki.CellRendererText, path string, newValue string) {
	iter, err := bl.listModel.GetIterFromString(path)
	if err != nil {
		bl.log.WithFields(log.Fields{
			"path":           path,
			"newReasonValue": newValue,
		}).WithError(err).Error("Can't get the iter to update the reason for the banned user")
		return
	}

	if err = bl.listModel.SetValue(iter, roomBanListReasonIndex, newValue); err != nil {
		bl.log.WithFields(log.Fields{
			"path":           path,
			"newReasonValue": newValue,
		}).WithError(err).Error("Can't set the new value for the reason of the banned user")
	}

	bl.enableApplyIfConditionsAreMet()
}

// onAddNewItem MUST be called from the UI thread
func (bl *roomBanListView) onAddNewItem() {
	bl.listView.Show()
	bl.noEntriesView.Hide()

	iter := bl.addListItem(&muc.RoomBanListItem{
		Affiliation: affiliationFromKnowString(data.AffiliationOutcast),
	})

	bl.unselectSelectedRows()
	bl.listSelection.SelectIter(iter)

	bl.enableApplyIfConditionsAreMet()
}

// onSelectionChanged MUST be called from the UI thread
func (bl *roomBanListView) onSelectionChanged() {
	totalSelected := len(bl.getSeledtedRows())
	bl.removeEntryButton.SetSensitive(totalSelected > 0)
	bl.removeEntryButton.SetTooltipText(i18n.Local("Remove selected item"))
	if totalSelected > 1 {
		bl.removeEntryButton.SetTooltipText(i18n.Local("Remove selected items"))
	}
}

// onRemoveItem MUST be called from the UI thread
func (bl *roomBanListView) onRemoveItem() {
	for _, path := range bl.getSeledtedRows() {
		iter, _ := bl.listModel.GetIter(path)
		bl.listModel.Remove(iter)
	}

	bl.enableApplyIfConditionsAreMet()
}

// unselectSelectedRows MUST be called from the UI thread
func (bl *roomBanListView) unselectSelectedRows() {
	for _, path := range bl.getSeledtedRows() {
		bl.listSelection.UnselectPath(path)
	}
}

// getSeledtedRows MUST be called from the UI thread
func (bl *roomBanListView) getSeledtedRows() []gtki.TreePath {
	return bl.listSelection.GetSelectedRows(bl.listModel)
}

// onCancel MUST be called from the UI thread
func (bl *roomBanListView) onCancel() {
	go bl.cancelActiveRequestListening()
	bl.dialog.Destroy()
}

// disableButtonsAndInteractions MUST be called from the UI thread
func (bl *roomBanListView) disableButtonsAndInteractions() {
	bl.addEntryButton.SetSensitive(false)
	bl.removeEntryButton.SetSensitive(false)
	bl.applyButton.SetSensitive(false)
}

// cancelActiveRequestListening MUST NOT be called from the UI thread
func (bl *roomBanListView) cancelActiveRequestListening() {
	if bl.cancelChannel != nil {
		bl.cancelChannel <- true
	}
}

// enableApplyIfConditionsAreMet MUST be called from the UI thread
func (bl *roomBanListView) enableApplyIfConditionsAreMet() {
	listHasChanged := bl.isTheListUpdated()

	listValuesAreValid := true
	for _, itm := range bl.listFromModel() {
		if itm.Jid.String() == "" {
			listValuesAreValid = false
			break
		}
	}

	bl.applyButton.SetSensitive(listHasChanged && listValuesAreValid)
}

// isTheListUpdated MUST be called from the UI thread
func (bl *roomBanListView) isTheListUpdated() bool {
	currentList := bl.listFromModel()

	if len(currentList) != len(bl.originalBanList) {
		return true
	}

	for idx, itm := range bl.originalBanList {
		currentItm := currentList[idx]
		if currentItm.Jid.String() != itm.Jid.String() ||
			currentItm.Affiliation.IsDifferentFrom(itm.Affiliation) ||
			currentItm.Reason != itm.Reason {
			return true
		}
	}

	return false
}

// listFromModel MUST be called from the UI thread
func (bl *roomBanListView) listFromModel() []*muc.RoomBanListItem {
	list := []*muc.RoomBanListItem{}

	iter, ok := bl.listModel.GetIterFirst()
	for ok {
		account := bl.columnStringValueFromListModelIter(iter, roomBanListAccountIndex)
		affiliation := bl.columnStringValueFromListModelIter(iter, roomBanListAccountIndex)
		reason := bl.columnStringValueFromListModelIter(iter, roomBanListAccountIndex)

		list = append(list, &muc.RoomBanListItem{
			Jid:         jid.Parse(account),
			Affiliation: affiliationFromKnowString(affiliation),
			Reason:      reason,
		})

		ok = bl.listModel.IterNext(iter)
	}

	return list
}

// columnValueFromListModelIter MUST be called from the UI thread
func (bl *roomBanListView) columnStringValueFromListModelIter(iter gtki.TreeIter, column int) string {
	v, _ := bl.listModel.GetValue(iter, column)
	s, _ := v.GetString()

	return s
}

func affiliationFromKnowString(a string) data.Affiliation {
	affiliation, _ := data.AffiliationFromString(a)
	return affiliation
}