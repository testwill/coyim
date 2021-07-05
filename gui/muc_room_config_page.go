package gui

import (
	"github.com/coyim/coyim/i18n"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

var roomConfigPagesFields = map[mucRoomConfigPageID][]muc.RoomConfigFieldType{
	roomConfigInformationPageIndex: {
		muc.RoomConfigFieldName,
		muc.RoomConfigFieldDescription,
		muc.RoomConfigFieldLanguage,
		muc.RoomConfigFieldIsPublic,
		muc.RoomConfigFieldIsPersistent,
	},
	roomConfigAccessPageIndex: {
		muc.RoomConfigFieldPassword,
		muc.RoomConfigFieldIsMembersOnly,
		muc.RoomConfigFieldAllowInvites,
	},
	roomConfigPermissionsPageIndex: {
		muc.RoomConfigFieldWhoIs,
		muc.RoomConfigFieldIsModerated,
		muc.RoomConfigFieldCanChangeSubject,
		muc.RoomConfigFieldAllowPrivateMessages,
		muc.RoomConfigFieldPresenceBroadcast,
	},
	roomConfigPositionsPageIndex: {
		muc.RoomConfigFieldOwners,
		muc.RoomConfigFieldAdmins,
		muc.RoomConfigFieldMembers,
	},
	roomConfigOthersPageIndex: {
		muc.RoomConfigFieldMaxOccupantsNumber,
		muc.RoomConfigFieldMaxHistoryFetch,
		muc.RoomConfigFieldEnableLogging,
		muc.RoomConfigFieldAllowVisitorNickchange,
		muc.RoomConfigFieldAllowVoiceRequest,
		muc.RoomConfigFieldAllowSubscription,
		muc.RoomConfigFieldMembersByDefault,
		muc.RoomConfigFieldAllowVisitorStatus,
		muc.RoomConfigAllowPrivateMessagesFromVisitors,
		muc.RoomConfigPublicList,
	},
}

var roomConfigAdvancedFields = []muc.RoomConfigFieldType{
	muc.RoomConfigFieldAllowQueryUsers,
	muc.RoomConfigFieldPubsub,
	muc.RoomConfigFieldVoiceRequestMinInteval,
}

type roomConfigPage struct {
	u      *gtkUI
	form   *muc.RoomConfigForm
	fields []hasRoomConfigFormField

	title               string
	pageID              mucRoomConfigPageID
	roomConfigComponent *mucRoomConfigComponent

	page                gtki.Overlay     `gtk-widget:"room-config-page-overlay"`
	header              gtki.Label       `gtk-widget:"room-config-page-header-label"`
	content             gtki.Box         `gtk-widget:"room-config-page-content"`
	notificationsArea   gtki.Box         `gtk-widget:"notifications-box"`
	autojoinContent     gtki.Box         `gtk-widget:"room-config-autojoin-content"`
	autojoinCheckButton gtki.CheckButton `gtk-widget:"room-config-autojoin"`

	notifications  *notificationsComponent
	loadingOverlay *loadingOverlayComponent
	doAfterRefresh *callbacksSet

	log coylog.Logger
}

func (c *mucRoomConfigComponent) newConfigPage(pageID mucRoomConfigPageID) *roomConfigPage {
	p := &roomConfigPage{
		u:                   c.u,
		roomConfigComponent: c,
		title:               configPageDisplayTitle(pageID),
		pageID:              pageID,
		loadingOverlay:      c.u.newLoadingOverlayComponent(),
		doAfterRefresh:      newCallbacksSet(),
		form:                c.form,
		log: c.log.WithFields(log.Fields{
			"page": pageID,
		}),
	}

	p.initBuilder()
	p.initDefaults()
	mucStyles.setRoomConfigPageStyle(p.content)

	return p
}

func (p *roomConfigPage) initBuilder() {
	builder := newBuilder("MUCRoomConfigPage")
	panicOnDevError(builder.bindObjects(p))
	builder.ConnectSignals(map[string]interface{}{
		"on_autojoin_toggled": func() {
			p.roomConfigComponent.updateAutoJoin(p.autojoinCheckButton.GetActive())
		},
	})

	p.notifications = p.u.newNotificationsComponent()
	p.loadingOverlay = p.u.newLoadingOverlayComponent()
	p.notificationsArea.Add(p.notifications.contentBox())
}

func (p *roomConfigPage) initDefaults() {
	p.initIntroPage()
	switch p.pageID {
	case roomConfigSummaryPageIndex:
		p.initSummary()
		return
	case roomConfigPositionsPageIndex:
		p.initOccupants()
		return
	case roomConfigOthersPageIndex:
		p.initKnownFields()
		p.initUnknownFields()
		p.initAdvancedOptionsFields()
		return
	}
	p.initKnownFields()
}

func (p *roomConfigPage) initIntroPage() {
	intro := configPageDisplayIntro(p.pageID)
	if intro == "" {
		p.header.SetVisible(false)
		return
	}
	p.header.SetText(intro)
}

func (p *roomConfigPage) initKnownFields() {
	if knownFields, ok := roomConfigPagesFields[p.pageID]; ok {
		booleanFields := []*roomConfigFormFieldBoolean{}
		for _, kf := range knownFields {
			if knownField, ok := p.form.GetKnownField(kf); ok {
				field, err := roomConfigFormFieldFactory(kf, roomConfigFieldsTexts[kf], knownField.ValueType())
				if err != nil {
					p.log.WithError(err).Error("Room configuration form field not supported")
					continue
				}
				if f, ok := field.(*roomConfigFormFieldBoolean); ok {
					booleanFields = append(booleanFields, f)
					continue
				}
				p.addField(field)
			}
		}
		if len(booleanFields) > 0 {
			p.addField(newRoomConfigFormFieldBooleanContainer(booleanFields))
		}
	}
}

func (p *roomConfigPage) initUnknownFields() {
	booleanFields := []*roomConfigFormFieldBoolean{}
	for _, ff := range p.form.GetUnknownFields() {
		field, err := roomConfigFormUnknownFieldFactory(newRoomConfigFieldTextInfo(ff.Label, ff.Description), ff.ValueType())
		if err != nil {
			p.log.WithError(err).Error("Room configuration form field not supported")
			continue
		}
		if f, ok := field.(*roomConfigFormFieldBoolean); ok {
			booleanFields = append(booleanFields, f)
			continue
		}
		p.addField(field)
	}
	if len(booleanFields) > 0 {
		p.addField(newRoomConfigFormFieldBooleanContainer(booleanFields))
	}
}

func (p *roomConfigPage) initAdvancedOptionsFields() {
	booleanFields := []*roomConfigFormFieldBoolean{}
	advancedFields := []hasRoomConfigFormField{}
	for _, aff := range roomConfigAdvancedFields {
		if knownField, ok := p.form.GetKnownField(aff); ok {
			field, err := roomConfigFormFieldFactory(aff, roomConfigFieldsTexts[aff], knownField.ValueType())
			if err != nil {
				p.log.WithError(err).Error("Room configuration form field not supported")
				continue
			}
			if f, ok := field.(*roomConfigFormFieldBoolean); ok {
				booleanFields = append(booleanFields, f)
				continue
			}
			advancedFields = append(advancedFields, field)
		}
	}
	if len(booleanFields) > 0 {
		advancedFields = append(advancedFields, newRoomConfigFormFieldBooleanContainer(booleanFields))
	}

	if len(advancedFields) > 0 {
		p.addField(newRoomConfigFormFieldAdvancedOptionsContainer(advancedFields))
	}
}

func (p *roomConfigPage) initSummary() {
	p.initSummaryFields(roomConfigInformationPageIndex)
	p.initSummaryFields(roomConfigAccessPageIndex)
	p.initSummaryFields(roomConfigPermissionsPageIndex)
	p.initSummaryFields(roomConfigPositionsPageIndex)
	p.initSummaryFields(roomConfigOthersPageIndex)
	p.autojoinCheckButton.SetActive(p.roomConfigComponent.autoJoin)
	p.autojoinContent.Show()
}

func (p *roomConfigPage) initSummaryFields(pageID mucRoomConfigPageID) {
	p.addField(newRoomConfigFormFieldLinkButton(pageID, p.roomConfigComponent.setCurrentPage))
	if pageID == roomConfigPositionsPageIndex {
		p.initOccupantsSummaryFields()
		return
	}

	fields := []hasRoomConfigFormField{}
	for _, kf := range roomConfigPagesFields[pageID] {
		if knownField, ok := p.form.GetKnownField(kf); ok {
			fields = append(fields, newRoomConfigSummaryField(kf, roomConfigFieldsTexts[kf], knownField.ValueType()))
		}
	}

	if pageID == roomConfigOthersPageIndex {
		fields = append(fields, p.otherPageSummaryFields()...)
	}

	p.addField(newRoomConfigSummaryFieldContainer(fields))
}

func (p *roomConfigPage) otherPageSummaryFields() []hasRoomConfigFormField {
	fields := []hasRoomConfigFormField{}

	for _, ff := range p.form.GetUnknownFields() {
		fields = append(fields, newRoomConfigSummaryField(muc.RoomConfigFieldUnexpected, newRoomConfigFieldTextInfo(ff.Label, ff.Description), ff.ValueType()))
	}

	advancedFields := []hasRoomConfigFormField{}
	for _, aff := range roomConfigAdvancedFields {
		if knownField, ok := p.form.GetKnownField(aff); ok {
			advancedFields = append(advancedFields, newRoomConfigSummaryField(aff, roomConfigFieldsTexts[aff], knownField.ValueType()))
		}
	}

	if len(advancedFields) > 0 {
		fields = append(fields, newAdvancedOptionSummaryField(advancedFields))
	}

	return fields
}

func (p *roomConfigPage) initOccupantsSummaryFields() {
	fields := []hasRoomConfigFormField{
		newRoomConfigSummaryOccupantField(&data.OwnerAffiliation{}, p.form.GetOccupantsByAffiliation),
		newRoomConfigSummaryOccupantField(&data.AdminAffiliation{}, p.form.GetOccupantsByAffiliation),
		newRoomConfigSummaryOccupantField(&data.OutcastAffiliation{}, p.form.GetOccupantsByAffiliation),
	}
	p.addField(newRoomConfigSummaryFieldContainer(fields))
}

func (p *roomConfigPage) initOccupants() {
	p.addField(newRoomConfigPositions(&data.OwnerAffiliation{}, p.form.GetOccupantsByAffiliation, p.form.UpdateRoomOccupantsByAffiliation))
	p.content.Add(createSeparator(gtki.HorizontalOrientation))
	p.addField(newRoomConfigPositions(&data.AdminAffiliation{}, p.form.GetOccupantsByAffiliation, p.form.UpdateRoomOccupantsByAffiliation))
	p.content.Add(createSeparator(gtki.HorizontalOrientation))
	p.addField(newRoomConfigPositions(&data.OutcastAffiliation{}, p.form.GetOccupantsByAffiliation, p.form.UpdateRoomOccupantsByAffiliation))
}

func (p *roomConfigPage) addField(field hasRoomConfigFormField) {
	p.fields = append(p.fields, field)
	p.content.Add(field.fieldWidget())
	p.doAfterRefresh.add(field.refreshContent)
}

// isValid MUST be called from the UI thread
func (p *roomConfigPage) isValid() bool {
	isValid := true
	for _, f := range p.fields {
		if !f.isValid() {
			f.showValidationErrors()
			isValid = false
		}
	}
	return isValid
}

func (p *roomConfigPage) updateFieldValues() {
	for _, f := range p.fields {
		f.updateFieldValue()
	}
}

// refresh MUST be called from the UI thread
func (p *roomConfigPage) refresh() {
	p.page.ShowAll()
	p.hideLoadingOverlay()
	p.clearErrors()
	p.doAfterRefresh.invokeAll()
}

// clearErrors MUST be called from the ui thread
func (p *roomConfigPage) clearErrors() {
	p.notifications.clearErrors()
}

// notifyError MUST be called from the ui thread
func (p *roomConfigPage) notifyError(m string) {
	p.notifications.notifyOnError(m)
}

// onConfigurationApply MUST be called from the ui thread
func (p *roomConfigPage) onConfigurationApply() {
	p.showLoadingOverlay(i18n.Local("Saving room configuration"))
}

// onConfigurationApplyError MUST be called from the ui thread
func (p *roomConfigPage) onConfigurationApplyError() {
	p.hideLoadingOverlay()
}

// showLoadingOverlay MUST be called from the ui thread
func (p *roomConfigPage) showLoadingOverlay(m string) {
	p.loadingOverlay.setSolid()
	p.loadingOverlay.showWithMessage(m)
}

// hideLoadingOverlay MUST be called from the ui thread
func (p *roomConfigPage) hideLoadingOverlay() {
	p.loadingOverlay.hide()
}
