package definitions

func init() {
	add(`AccountDetails`, &defAccountDetails{})
}

type defAccountDetails struct{}

func (*defAccountDetails) String() string {
	return `<interface>
  <object class="GtkListStore" id="proxies-model">
    <columns>
      <!-- proxy -->
      <column type="gchararray"/>
      <!-- real proxy data -->
      <column type="gchararray"/>
    </columns>
  </object>

  <object class="GtkListStore" id="pins-model">
    <columns>
      <!-- Subject -->
      <column type="gchararray"/>
      <!-- Issuer -->
      <column type="gchararray"/>
      <!-- Fingerprint -->
      <column type="gchararray"/>
    </columns>
  </object>

  <object class="GtkDialog" id="AccountDetails">
    <property name="title" translatable="yes">Account Details</property>
    <signal name="close" handler="on_cancel_signal" />
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="margin">10</property>

        <child>
          <object class="GtkBox" id="notification-area">
            <property name="visible">true</property>
            <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
          </object>
        </child>

        <child>
          <object class="GtkNotebook" id="notebook1">
            <property name="visible">True</property>
            <property name="show-border">False</property>
            <property name="page">0</property>
            <property name="margin-bottom">10</property>
            <child>
              <object class="GtkGrid" id="grid">
                <property name="margin-top">15</property>
                <property name="margin-bottom">10</property>
                <property name="margin-start">10</property>
                <property name="margin-end">10</property>
                <property name="row-spacing">12</property>
                <property name="column-spacing">6</property>
                <child>
                  <object class="GtkLabel" id="AccountMessageLabel">
                    <property name="label" translatable="yes">Your account&#xA;(example: kim42@jabber.otr.im)</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="account">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkLabel" id="PasswordLabel">
                    <property name="label" translatable="yes">Password</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="password">
                    <property name="visibility">false</property>
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>

               <child>
                  <object class="GtkLabel" id="DisplayNameLabel">
                    <property name="label" translatable="yes">Account name</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">2</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="displayName">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="showOtherSettings">
                    <property name="label" translatable="yes">Display all settings</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="otherSettings">
                    <signal name="toggled" handler="on_toggle_other_settings" />
                  </object>

                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>
              </object>
            </child>

            <child type="tab">
              <object class="GtkLabel" id="label-tab1">
                <property name="label" translatable="yes">Account</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">0</property>
                <property name="tab-fill">False</property>
              </packing>
            </child>

            <child>
              <object class="GtkGrid" id="otherOptionsGrid">
                <property name="margin-top">15</property>
                <property name="margin-bottom">10</property>
                <property name="margin-start">10</property>
                <property name="margin-end">10</property>
                <property name="row-spacing">12</property>
                <property name="column-spacing">6</property>
                <child>
                  <object class="GtkLabel" id="serverLabel">
                    <property name="label" translatable="yes">Server (leave empty for default)</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="server">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkLabel" id="portLabel">
                    <property name="label" translatable="yes">Port (leave empty for default)</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="port">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="pinningPolicyInstructions">
                    <property name="label" translatable="yes">The pinning policy governs whether we will consider saving information about certificates we have seen before, and how we will react in these cases. The none policy turns off this behavior. Deny will try to use the pins we already have, but will never connect to a server that has another certificate. The add policy will always add new pins - by itself, this is not so useful, but if you turn it off or change it to deny, you will have a list of pins that you can curate later. The add first and ask after policy will pin the first certificate we ever see, and then let the user know if we encounter other certifacates. The add first and deny policy does the same thing except it doesn't ever ask after the first certificate. Finally, the ask policy will always ask what to do when seeing certificates we haven't added.</property>
                    <property name="visible">true</property>
                    <property name="wrap">true</property>
                    <property name="max-width-chars">50</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">2</property>
                    <property name="width">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="pinningPolicyLabel">
                    <property name="label" translatable="yes">Pinning policy</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkComboBoxText" id="pinningPolicyValue">
                    <items>
                      <item translatable="yes" id="none">None</item>
                      <item translatable="yes" id="deny">Deny</item>
                      <item translatable="yes" id="add">Always add</item>
                      <item translatable="yes" id="add-first-ask-rest">Add the first, ask for the rest</item>
                      <item translatable="yes" id="add-first-deny-rest">Add the first, deny the rest</item>
                      <item translatable="yes" id="ask">Always ask</item>
                    </items>
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkPaned" id="hpanedPins">
                    <property name="visible">True</property>
                    <property name="can-focus">True</property>
                    <property name="position">175</property>
                    <child>
                      <object class="GtkScrolledWindow" id="scrolledwindowPins">
                        <property name="hscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
                        <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
                        <property name="width-request">170</property>
                        <property name="height-request">150</property>
                        <property name="margin">5</property>
                        <property name="visible">True</property>
                        <property name="hexpand">True</property>
                        <property name="vexpand">True</property>
                        <property name="can-focus">True</property>
                        <property name="shadow-type">in</property>
                        <child>
                          <object class="GtkTreeView" id="pins-view">
                            <property name="model">pins-model</property>
                            <property name="visible">True</property>
                            <property name="can-focus">True</property>
                            <property name="headers-visible">True</property>
                            <property name="show-expanders">False</property>
                            <property name="reorderable">True</property>
                            <child internal-child="selection">
                              <object class="GtkTreeSelection" id="pins-selection">
                                <property name="mode">GTK_SELECTION_SINGLE</property>
                              </object>
                            </child>
                            <child>
                              <object class="GtkTreeViewColumn" id="pins-subject-column">
                                <property name="title">Subject</property>
                                <child>
                                  <object class="GtkCellRendererText" id="pins-subject-column-rendered"/>
                                  <attributes>
                                    <attribute name="text">0</attribute>
                                  </attributes>
                                </child>
                              </object>
                            </child>
                            <child>
                              <object class="GtkTreeViewColumn" id="pins-issuer-column">
                                <property name="title">Issuer</property>
                                <child>
                                  <object class="GtkCellRendererText" id="pins-issuer-column-rendered"/>
                                  <attributes>
                                    <attribute name="text">1</attribute>
                                  </attributes>
                                </child>
                              </object>
                            </child>
                            <child>
                              <object class="GtkTreeViewColumn" id="pins-fpr-column">
                                <property name="title">Fingerprint</property>
                                <child>
                                  <object class="GtkCellRendererText" id="pins-fpr-column-rendered"/>
                                  <attributes>
                                    <attribute name="text">2</attribute>
                                  </attributes>
                                </child>
                              </object>
                            </child>
                          </object>
                        </child>
                      </object>
                      <packing>
                        <property name="resize">True</property>
                        <property name="shrink">False</property>
                      </packing>
                    </child>
                    <child>
                      <object class="GtkBox" id="vbox-pinbuttons">
                        <property name="margin">5</property>
                        <property name="visible">True</property>
                        <property name="can-focus">False</property>
                        <property name="orientation">vertical</property>
                        <property name="spacing">6</property>
                        <child>
                          <object class="GtkButton" id="remove_pin_button">
                            <property name="label" translatable="yes">_Remove</property>
                            <property name="visible">True</property>
                            <property name="can-focus">True</property>
                            <property name="receives-default">True</property>
                            <property name="use_underline">True</property>
                            <signal name="clicked" handler="on_remove_pin_signal" />
                          </object>
                          <packing>
                            <property name="expand">False</property>
                            <property name="fill">True</property>
                            <property name="position">0</property>
                          </packing>
                        </child>
                      </object>
                      <packing>
                        <property name="resize">False</property>
                        <property name="shrink">False</property>
                      </packing>
                    </child>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">4</property>
                    <property name="width">2</property>
                  </packing>
                </child>
              </object>
            </child>
            <child type="tab">
              <object class="GtkLabel" id="label-tab2">
                <property name="label" translatable="yes">Server</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">1</property>
                <property name="tab-fill">False</property>
              </packing>
            </child>
            <child>
              <object class="GtkBox" id="vbox1">
                <property name="margin">5</property>
                <property name="visible">True</property>
                <property name="can-focus">False</property>
                <property name="orientation">vertical</property>
                <property name="spacing">6</property>
                <child>
                  <object class="GtkPaned" id="hpaned1">
                    <property name="visible">True</property>
                    <property name="can-focus">True</property>
                    <property name="position">175</property>
                    <child>
                      <object class="GtkScrolledWindow" id="scrolledwindow1">
                        <property name="hscrollbar-policy">GTK_POLICY_NEVER</property>
                        <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
                        <property name="width-request">170</property>
                        <property name="height-request">230</property>
                        <property name="margin">5</property>
                        <property name="visible">True</property>
                        <property name="hexpand">True</property>
                        <property name="vexpand">True</property>
                        <property name="can-focus">True</property>
                        <property name="shadow-type">in</property>
                        <child>
                          <object class="GtkTreeView" id="proxies-view">
                            <property name="model">proxies-model</property>
                            <property name="visible">True</property>
                            <property name="can-focus">True</property>
                            <property name="headers-visible">False</property>
                            <property name="show-expanders">False</property>
                            <property name="reorderable">True</property>
                            <signal name="row-activated" handler="on_edit_activate_proxy_signal" />
                            <child internal-child="selection">
                              <object class="GtkTreeSelection" id="selection">
                                <property name="mode">GTK_SELECTION_SINGLE</property>
                              </object>
                            </child>
                            <child>
                              <object class="GtkTreeViewColumn" id="proxy-name-column">
                                <property name="title">proxy-name</property>
                                <child>
                                  <object class="GtkCellRendererText" id="proxy-name-column-rendered"/>
                                  <attributes>
                                    <attribute name="text">0</attribute>
                                  </attributes>
                                </child>
                              </object>
                            </child>
                          </object>
                        </child>
                      </object>
                      <packing>
                        <property name="resize">True</property>
                        <property name="shrink">False</property>
                      </packing>
                    </child>
                    <child>
                      <object class="GtkBox" id="vbox3">
                        <property name="margin">5</property>
                        <property name="visible">True</property>
                        <property name="can-focus">False</property>
                        <property name="orientation">vertical</property>
                        <property name="spacing">6</property>
                        <child>
                          <object class="GtkButton" id="add_button">
                            <property name="label" translatable="yes">_Add...</property>
                            <property name="visible">True</property>
                            <property name="can-focus">True</property>
                            <property name="receives-default">True</property>
                            <property name="use_underline">True</property>
                            <signal name="clicked" handler="on_add_proxy_signal" />
                          </object>
                          <packing>
                            <property name="expand">False</property>
                            <property name="fill">True</property>
                            <property name="position">0</property>
                          </packing>
                        </child>
                        <child>
                          <object class="GtkButton" id="remove_button">
                            <property name="label" translatable="yes">_Remove</property>
                            <property name="visible">True</property>
                            <property name="can-focus">True</property>
                            <property name="receives-default">True</property>
                            <property name="use_underline">True</property>
                            <signal name="clicked" handler="on_remove_proxy_signal" />
                          </object>
                          <packing>
                            <property name="expand">False</property>
                            <property name="fill">True</property>
                            <property name="position">1</property>
                          </packing>
                        </child>
                        <child>
                          <object class="GtkButton" id="edit_button">
                            <property name="label" translatable="yes">_Edit...</property>
                            <property name="visible">True</property>
                            <property name="can-focus">True</property>
                            <property name="receives-default">True</property>
                            <property name="use-underline">True</property>
                            <signal name="clicked" handler="on_edit_proxy_signal" />
                          </object>
                          <packing>
                            <property name="expand">False</property>
                            <property name="fill">True</property>
                            <property name="position">2</property>
                          </packing>
                        </child>
                      </object>
                      <packing>
                        <property name="resize">False</property>
                        <property name="shrink">False</property>
                      </packing>
                    </child>
                  </object>
                  <packing>
                    <property name="expand">True</property>
                    <property name="fill">True</property>
                    <property name="position">0</property>
                  </packing>
                </child>
              </object>
            </child>
            <child type="tab">
              <object class="GtkLabel" id="label-tab3">
                <property name="label" translatable="yes">Proxies</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">2</property>
                <property name="tab-fill">False</property>
              </packing>
            </child>

            <child>
              <object class="GtkBox" id="encryptionOptionsBox">
                <property name="border-width">10</property>
                <property name="homogeneous">false</property>
                <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
                <child>
                  <object class="GtkLabel" id="fingerprintsMessage">
                    <property name="justify">GTK_JUSTIFY_LEFT</property>
                    <property name="halign">GTK_ALIGN_START</property>
                  </object>
                </child>
                <child>
                  <object class="GtkGrid" id="encryptionGrid">
                    <property name="margin-top">15</property>
                    <property name="margin-bottom">10</property>
                    <property name="margin-start">10</property>
                    <property name="margin-end">10</property>
                    <property name="row-spacing">12</property>
                    <property name="column-spacing">6</property>
                    <child>
                      <object class="GtkLabel" id="encryptionImportInstructions">
                        <property name="label" translatable="yes">The below buttons allow you to import private keys and fingerprints. Both of them should be in the Pidgin/libotr format. If you import private keys, your existing private keys will be deleted, since currently there is no way to choose which key to use for encrypted chat.

There are several applications that use the libotr format - Pidgin, Adium and Tor Messenger are most well known ones. Depending on your platform, these files can be found in several different places. Refer to the documentation for the application in question to find out where the files are located for your platform. The filenames to look for are "otr.fingerprints" and "otr.private_key".</property>
                        <property name="visible">true</property>
                        <property name="wrap">true</property>
                        <property name="max-width-chars">50</property>
                      </object>
                      <packing>
                        <property name="left-attach">0</property>
                        <property name="top-attach">0</property>
                        <property name="width">2</property>
                      </packing>
                    </child>
                    <child>
                      <object class="GtkButton" id="import_key_button">
                        <property name="label" translatable="yes">Import Private _Keys...</property>
                        <property name="visible">True</property>
                        <property name="can-focus">True</property>
                        <property name="receives-default">True</property>
                        <property name="use_underline">True</property>
                        <signal name="clicked" handler="on_import_key_signal" />
                      </object>
                      <packing>
                        <property name="left-attach">0</property>
                        <property name="top-attach">1</property>
                      </packing>
                    </child>
                    <child>
                      <object class="GtkButton" id="import_fpr_button">
                        <property name="label" translatable="yes">Import _Fingerprints...</property>
                        <property name="visible">True</property>
                        <property name="can-focus">True</property>
                        <property name="receives-default">True</property>
                        <property name="use_underline">True</property>
                        <signal name="clicked" handler="on_import_fpr_signal" />
                      </object>
                      <packing>
                        <property name="left-attach">1</property>
                        <property name="top-attach">1</property>
                      </packing>
                    </child>
                    <child>
                      <object class="GtkLabel" id="encryptionExportInstructions">
                        <property name="label" translatable="yes">The below buttons allow you to export private keys and fingerprints. Be careful with the files that come out of this process, since they contain potentially sensitive data. The export will be in Pidgin/libotr format.</property>
                        <property name="visible">true</property>
                        <property name="wrap">true</property>
                        <property name="max-width-chars">50</property>
                      </object>
                      <packing>
                        <property name="left-attach">0</property>
                        <property name="top-attach">2</property>
                        <property name="width">2</property>
                      </packing>
                    </child>
                    <child>
                      <object class="GtkButton" id="export_key_button">
                        <property name="label" translatable="yes">Export Private Keys...</property>
                        <property name="visible">True</property>
                        <property name="can-focus">True</property>
                        <property name="receives-default">True</property>
                        <signal name="clicked" handler="on_export_key_signal" />
                      </object>
                      <packing>
                        <property name="left-attach">0</property>
                        <property name="top-attach">3</property>
                      </packing>
                    </child>
                    <child>
                      <object class="GtkButton" id="export_fpr_button">
                        <property name="label" translatable="yes">Export Fingerprints...</property>
                        <property name="visible">True</property>
                        <property name="can-focus">True</property>
                        <property name="receives-default">True</property>
                        <signal name="clicked" handler="on_export_fpr_signal" />
                      </object>
                      <packing>
                        <property name="left-attach">1</property>
                        <property name="top-attach">3</property>
                      </packing>
                    </child>

                  </object>
                </child>
              </object>
            </child>
            <child type="tab">
              <object class="GtkLabel" id="label-tab4">
                <property name="label" translatable="yes">Encryption</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">3</property>
                <property name="tab-fill">False</property>
              </packing>
            </child>
          </object>
        </child>

        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="cancel">
                <property name="label" translatable="yes">Cancel</property>
                <signal name="clicked" handler="on_cancel_signal"/>
              </object>
            </child>
            <child>
              <object class="GtkButton" id="save">
                <property name="label" translatable="yes">Save</property>
                <property name="can-default">true</property>
                <signal name="clicked" handler="on_save_signal"/>
              </object>
            </child>
          </object>
        </child>

      </object>
    </child>
  </object>
</interface>
`
}
