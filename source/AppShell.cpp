#include "multiflexitui/TV.h"
#include "multiflexitui/AppShell.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/i18n.h"
#include "multiflexitui/StatusView.h"
#include "multiflexitui/AboutView.h"
#include "multiflexitui/EntityListView.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/EntityPicker.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/Wizards.h"
#include "multiflexitui/JobStreamView.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"

#include <algorithm>
#include <cstring>

namespace multiflexitui {

namespace {

class ImportExportDialog : public TDialog {
public:
    explicit ImportExportDialog(CliClient &client)
        : TWindowInit(&TDialog::initFrame),
          TDialog(TRect(8, 5, 72, 16), _("App JSON Import/Export")),
          client_(client) {
        options |= ofCentered;
        insert(new TLabel(TRect(2, 2, 10, 3), _("Path"), nullptr));
        path_ = new TInputLine(TRect(11, 2, 62, 3), 255);
        const std::string def = defaultExportPath("app.json");
        char buf[256];
        std::strncpy(buf, def.c_str(), sizeof(buf) - 1);
        buf[sizeof(buf) - 1] = '\0';
        path_->setData(buf);
        insert(path_);

        insert(new AppButton(TRect(2, 5, 16, 7), _("~I~mport"), cmImportJson, bfNormal));
        insert(new AppButton(TRect(18, 5, 34, 7), _("~V~alidate"), cmValidateJson, bfNormal));
        insert(new AppButton(TRect(36, 5, 52, 7), _("E~x~port…"), cmExportJsonHint, bfNormal));
        insert(new AppButton(TRect(48, 9, 62, 11), _("Close"), cmCancel, bfDefault));
        selectNext(False);
    }

    void handleEvent(TEvent &event) override {
        TDialog::handleEvent(event);
        if (event.what != evCommand) {
            return;
        }
        char buf[256] = {};
        path_->getData(buf);
        const std::string path = buf;
        if (event.message.command == cmImportJson) {
            auto r = client_.runJson({"application:import-json", "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
            } else {
                messageBox(_("Import completed."), mfInformation | mfOKButton);
            }
            clearEvent(event);
        } else if (event.message.command == cmValidateJson) {
            auto r = client_.runJson({"application:validate-json", "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
            } else {
                TProgram::application->executeDialog(new TextViewerDialog("validate-json", r.data.dump(2)));
            }
            clearEvent(event);
        } else if (event.message.command == cmExportJsonHint) {
            messageBox(_("Use Applications → row Actions → Export JSON to export a selected app."),
                       mfInformation | mfOKButton);
            clearEvent(event);
        }
    }

    TColorAttr mapColor(uchar index) override {
        TColorAttr color;
        return windowColor(index, color) ? color : TDialog::mapColor(index);
    }

private:
    CliClient &client_;
    TInputLine *path_ = nullptr;
};

class PruneDialog : public TDialog {
public:
    PruneDialog()
        : TWindowInit(&TDialog::initFrame), TDialog(TRect(14, 6, 66, 18), _("Prune")) {
        options |= ofCentered;

        insert(new TLabel(TRect(2, 2, 30, 3), _("What to prune"), nullptr));
        // Bit0 = logs, bit1 = jobs — both on by default.
        checks_ = new TCheckBoxes(TRect(3, 3, 30, 5),
                                  new TSItem(_("~L~ogs table"), new TSItem(_("~J~obs table"), nullptr)));
        {
            ushort marks = 0x3; // both checked
            checks_->setData(&marks);
        }
        insert(checks_);

        keep_ = new TInputLine(TRect(30, 6, 42, 7), 12);
        char buf[16] = "1000";
        keep_->setData(buf);
        insert(keep_);
        insert(new TLabel(TRect(2, 6, 28, 7), _("Keep latest N records"), keep_));

        insert(new AppButton(TRect(2, 9, 16, 11), _("~P~rune"), cmOK, bfDefault));
        insert(new AppButton(TRect(34, 9, 48, 11), _("Cancel"), cmCancel, bfNormal));
        selectNext(False);
    }

    ushort selectedMask() const {
        ushort marks = 0;
        if (checks_ != nullptr) {
            checks_->getData(&marks);
        }
        return marks;
    }

    bool pruneLogs() const { return (selectedMask() & 0x1) != 0; }
    bool pruneJobs() const { return (selectedMask() & 0x2) != 0; }

    std::string keepCount() const {
        char buf[16] = {};
        if (keep_ != nullptr) {
            keep_->getData(buf);
        }
        std::string s = buf;
        while (!s.empty() && (s.back() == ' ' || s.back() == '\0')) {
            s.pop_back();
        }
        return s.empty() ? "1000" : s;
    }

    TColorAttr mapColor(uchar index) override {
        TColorAttr color;
        return windowColor(index, color) ? color : TDialog::mapColor(index);
    }

private:
    TCheckBoxes *checks_ = nullptr;
    TInputLine *keep_ = nullptr;
};

} // namespace

MultiFlexiApp::MultiFlexiApp()
    : TProgInit(&MultiFlexiApp::initStatusLine, &MultiFlexiApp::initMenuBar, &MultiFlexiApp::initDeskTop) {
    statusLine_ = dynamic_cast<AppStatusLine *>(statusLine);
}

void MultiFlexiApp::configure(std::string cliBinary, std::string envFile) {
    client_.setBinaryPath(std::move(cliBinary));
    client_.setEnvFile(std::move(envFile));
    client_.setRequestObserver([this](const std::string &cmd) {
        if (statusLine_ != nullptr) {
            statusLine_->setCurrentUrl(cmd);
        }
    });
    registerBuiltinEntities();
}

void MultiFlexiApp::setBusy(bool busy) {
    busy_ = busy;
    updateBusyBadge();
}

void MultiFlexiApp::updateBusyBadge() {
    if (statusLine_ == nullptr) {
        return;
    }
    if (busy_ || cliQueue_.busy()) {
        statusLine_->setQueryMode(_("Busy…"));
    } else {
        statusLine_->setQueryMode("");
    }
}

void MultiFlexiApp::setActiveCompany(int id, const std::string &name) {
    activeCompanyId_ = id;
    activeCompanyName_ = name;
    if (statusLine_ != nullptr) {
        if (id > 0) {
            statusLine_->setCompany(name.empty() ? ("co:" + std::to_string(id)) : name);
        } else {
            statusLine_->setCompany("");
        }
    }
}

void MultiFlexiApp::addStreamTick(void *owner, std::function<void()> tick) {
    removeStreamTick(owner);
    streamTicks_.push_back(StreamTick{owner, std::move(tick)});
}

void MultiFlexiApp::removeStreamTick(void *owner) {
    streamTicks_.erase(std::remove_if(streamTicks_.begin(), streamTicks_.end(),
                                      [owner](const StreamTick &t) { return t.owner == owner; }),
                       streamTicks_.end());
}

TMenuBar *MultiFlexiApp::initMenuBar(TRect r) {
    r.b.y = r.a.y + 1;

    TSubMenu &multi = *new TSubMenu(_("~M~ultiFlexi"), kbAltM) +
                      *new TMenuItem(_("~S~tatus"), cmShowStatus, kbAltS) +
                      *new TMenuItem(_("Set ~c~ompany…"), cmOpenSetCompany, kbNoKey) +
                      newLine() +
                      *new TMenuItem(_("E~x~it"), cmQuit, kbAltX, hcNoContext, "Alt-X");

    TSubMenu &companies = *new TSubMenu(_("~C~ompanies"), kbAltC) +
                          *new TMenuItem(_("~L~ist"), cmOpenCompanies, kbNoKey);

    TSubMenu &apps = *new TSubMenu(_("~A~pplications"), kbAltA) +
                     *new TMenuItem(_("~L~ist"), cmOpenApplications, kbNoKey) +
                     *new TMenuItem(_("~C~onfig fields"), cmOpenConfFields, kbNoKey) +
                     *new TMenuItem(_("Company ~a~pps"), cmOpenCompanyApps, kbNoKey) +
                     newLine() +
                     *new TMenuItem(_("Activation ~W~izard"), cmOpenActivationWizard, kbNoKey);

    TSubMenu &runtemplates = *new TSubMenu(_("~R~unTemplates"), kbAltR) +
                             *new TMenuItem(_("~L~ist"), cmOpenRunTemplates, kbNoKey) +
                             *new TMenuItem(_("~S~tale"), cmOpenStaleRunTemplates, kbNoKey);

    TSubMenu &creds = *new TSubMenu(_("Credentia~l~s"), kbAltL) +
                      *new TMenuItem(_("~C~redentials"), cmOpenCredentials, kbNoKey) +
                      *new TMenuItem(_("~T~ypes"), cmOpenCredTypes, kbNoKey) +
                      *new TMenuItem(_("~P~rototypes"), cmOpenCrPrototypes, kbNoKey) +
                      newLine() +
                      *new TMenuItem(_("Credential ~W~izard"), cmOpenCredentialWizard, kbNoKey);

    TSubMenu &jobs = *new TSubMenu(_("~J~obs"), kbAltJ) +
                     *new TMenuItem(_("~J~obs"), cmOpenJobs, kbNoKey) +
                     *new TMenuItem(_("~T~asks"), cmOpenTasks, kbNoKey) +
                     *new TMenuItem(_("~Q~ueue"), cmOpenQueue, kbNoKey) +
                     *new TMenuItem(_("Queue ~o~verview"), cmOpenQueueOverview, kbNoKey) +
                     *new TMenuItem(_("Job st~a~tus"), cmOpenJobStatus, kbNoKey) +
                     *new TMenuItem(_("Task s~t~atus"), cmOpenTaskStatus, kbNoKey) +
                     newLine() +
                     *new TMenuItem(_("Live ~s~tream…"), cmOpenJobStream, kbNoKey) +
                     *new TMenuItem(_("Arti~f~acts"), cmOpenArtifacts, kbNoKey);

    TSubMenu &events = *new TSubMenu(_("~E~vents"), kbAltE) +
                       *new TMenuItem(_("Event ~s~ources"), cmOpenEventSources, kbNoKey) +
                       *new TMenuItem(_("Event ~r~ules"), cmOpenEventRules, kbNoKey);

    TSubMenu &admin = *new TSubMenu(_("Ad~m~in"), kbAltD) +
                      *new TMenuItem(_("~U~sers"), cmOpenUsers, kbNoKey) +
                      *new TMenuItem(_("~T~okens"), cmOpenTokens, kbNoKey) +
                      *new TMenuItem(_("Deletion ~r~equests"), cmOpenUserErasure, kbNoKey) +
                      newLine() +
                      *new TMenuItem(_("~E~ncryption"), cmOpenEncryption, kbNoKey) +
                      *new TMenuItem(_("~P~rune…"), cmOpenPrune, kbNoKey) +
                      *new TMenuItem(_("Telemetry ~t~est"), cmOpenTelemetryTest, kbNoKey) +
                      *new TMenuItem(_("~I~mport/Export app JSON…"), cmOpenImportExport, kbNoKey);

    TSubMenu &window = *new TSubMenu(_("~W~indow"), kbAltW) +
                       *new TMenuItem(_("~T~ile"), cmTile, kbNoKey) +
                       *new TMenuItem(_("C~a~scade"), cmCascade, kbNoKey) +
                       *new TMenuItem(_("~M~inimize all"), cmMinimizeAll, kbNoKey) +
                       *new TMenuItem(_("~R~estore"), cmRestoreWindows, kbNoKey) +
                       *new TMenuItem(_("Close a~l~l"), cmCloseAll, kbNoKey);

    TSubMenu &help = *new TSubMenu(_("~H~elp"), kbAltH) +
                     *new TMenuItem(_("~A~bout"), cmShowAbout, kbF1) +
                     newLine() +
                     *new TMenuItem(_("Language: ~S~ystem"), cmLangSystem, kbNoKey) +
                     *new TMenuItem(_("Language: ~E~nglish"), cmLangEnglish, kbNoKey) +
                     *new TMenuItem(_("Language: ~C~zech"), cmLangCzech, kbNoKey);

    return new TMenuBar(r, multi + companies + apps + runtemplates + creds + jobs + events + admin + window + help);
}

TStatusLine *MultiFlexiApp::initStatusLine(TRect r) {
    r.a.y = r.b.y - 1;
    return new AppStatusLine(r, *new TStatusDef(0, 0xFFFF) +
                                    *new TStatusItem(_("~Alt-X~ Exit"), kbAltX, cmQuit) +
                                    *new TStatusItem(_("~F1~ About"), kbF1, cmShowAbout) +
                                    *new TStatusItem(_("~Alt-S~ Status"), kbAltS, cmShowStatus) +
                                    *new TStatusItem(nullptr, kbF10, cmMenu));
}

void MultiFlexiApp::openStatus() {
    executeDialog(new StatusView(client_));
}

void MultiFlexiApp::openAbout() {
    executeDialog(new AboutView());
}

void MultiFlexiApp::openEntity(const std::string &cliEntity) {
    openEntityList(this, client_, cliEntity);
}

void MultiFlexiApp::openActivationWizard() {
    executeDialog(new ActivationWizard(client_));
}

void MultiFlexiApp::openCredentialWizard() {
    executeDialog(new CredentialWizard(client_));
}

void MultiFlexiApp::openJobStream() {
    openJobStreamPicker(this, client_);
}

void MultiFlexiApp::openSetCompany() {
    if (activeCompanyId_ > 0) {
        const ushort ans =
            messageBox(_("Clear active company filter? (No = pick another)"), mfYesNoCancel | mfConfirmation);
        if (ans == cmYes) {
            setActiveCompany(0, "");
            return;
        }
        if (ans == cmCancel) {
            return;
        }
    }
    auto *dlg = new EntityPickerDialog(client_, "company", _("Set active company"));
    const ushort code = deskTop->execView(dlg);
    const int id = (code == cmOK) ? dlg->selectedId() : 0;
    std::string name;
    if (id > 0) {
        const auto row = dlg->selectedRow();
        if (row.is_object() && row.contains("name")) {
            name = jsonToString(row["name"]);
        } else if (row.is_object() && row.contains("slug")) {
            name = jsonToString(row["slug"]);
        }
        if (name.empty()) {
            name = std::to_string(id);
        }
    }
    TObject::destroy(dlg);
    if (id <= 0) {
        return;
    }
    setActiveCompany(id, name);
}

void MultiFlexiApp::openAdminJsonResult(const std::vector<std::string> &args, const std::string &title) {
    auto r = client_.runJson(args);
    if (!r.ok) {
        messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
        return;
    }
    if (r.data.is_object()) {
        executeDialog(new KeyValueDialog(title, r.data));
    } else {
        executeDialog(new TextViewerDialog(title, r.data.dump(2)));
    }
}

void MultiFlexiApp::openPruneDialog() {
    auto *dlg = new PruneDialog();
    const ushort code = deskTop->execView(dlg);
    const bool logs = dlg->pruneLogs();
    const bool jobs = dlg->pruneJobs();
    const std::string keep = dlg->keepCount();
    TObject::destroy(dlg);
    if (code != cmOK) {
        return;
    }
    if (!logs && !jobs) {
        messageBox(_("Select at least one of: logs, jobs."), mfError | mfOKButton);
        return;
    }
    std::vector<std::string> args{"prune", "--keep=" + keep};
    if (logs) {
        args.push_back("--logs");
    }
    if (jobs) {
        args.push_back("--jobs");
    }
    openAdminJsonResult(args, "prune");
}

void MultiFlexiApp::openImportExportDialog() {
    executeDialog(new ImportExportDialog(client_));
}

void MultiFlexiApp::reloadMenuAndStatusLine() {
    TRect mb = menuBar->getExtent();
    TMenuBar *newMb = initMenuBar(mb);
    destroy(menuBar);
    menuBar = newMb;
    insert(menuBar);

    TRect sl = statusLine->getExtent();
    TStatusLine *newSl = initStatusLine(sl);
    destroy(statusLine);
    statusLine = newSl;
    statusLine_ = dynamic_cast<AppStatusLine *>(statusLine);
    insert(statusLine);
    if (activeCompanyId_ > 0) {
        setActiveCompany(activeCompanyId_, activeCompanyName_);
    }
    updateBusyBadge();
}

void MultiFlexiApp::minimizeAll() {
    minimized_.clear();
    TView *p = deskTop->first();
    while (p != nullptr) {
        TView *next = p->nextView();
        if (auto *w = dynamic_cast<TWindow *>(p)) {
            if ((w->flags & wfClose) != 0) {
                MinimizedWindow m{w, w->getBounds()};
                minimized_.push_back(m);
                TRect r = w->getBounds();
                r.b.y = r.a.y + 1;
                w->locate(r);
            }
        }
        p = next;
    }
}

void MultiFlexiApp::restoreWindows() {
    for (auto &m : minimized_) {
        if (m.window != nullptr) {
            m.window->locate(m.bounds);
        }
    }
    minimized_.clear();
}

void MultiFlexiApp::closeAllWindows() {
    TView *p = deskTop->first();
    while (p != nullptr) {
        TView *next = p->nextView();
        if (dynamic_cast<TWindow *>(p) != nullptr) {
            message(p, evCommand, cmClose, nullptr);
        }
        p = next;
    }
    minimized_.clear();
}

void MultiFlexiApp::handleEvent(TEvent &event) {
    TApplication::handleEvent(event);

    if (event.what != evCommand) {
        return;
    }

    switch (event.message.command) {
    case cmShowStatus:
        openStatus();
        clearEvent(event);
        break;
    case cmShowAbout:
        openAbout();
        clearEvent(event);
        break;
    case cmShowWebQr:
        clearEvent(event);
        break;
    case cmOpenSetCompany:
        openSetCompany();
        clearEvent(event);
        break;
    case cmOpenCompanies:
        openEntity("company");
        clearEvent(event);
        break;
    case cmOpenApplications:
        openEntity("application");
        clearEvent(event);
        break;
    case cmOpenConfFields:
        openEntity("conffield");
        clearEvent(event);
        break;
    case cmOpenRunTemplates:
        openEntity("run-template");
        clearEvent(event);
        break;
    case cmOpenJobs:
        openEntity("job");
        clearEvent(event);
        break;
    case cmOpenTasks:
        openEntity("task");
        clearEvent(event);
        break;
    case cmOpenCredentials:
        openEntity("credential");
        clearEvent(event);
        break;
    case cmOpenTokens:
        openEntity("token");
        clearEvent(event);
        break;
    case cmOpenUsers:
        openEntity("user");
        clearEvent(event);
        break;
    case cmOpenArtifacts:
        openEntity("artifact");
        clearEvent(event);
        break;
    case cmOpenCredTypes:
        openEntity("credential-type");
        clearEvent(event);
        break;
    case cmOpenCrPrototypes:
        openEntity("credential-prototype");
        clearEvent(event);
        break;
    case cmOpenCompanyApps:
        openEntity("company-app");
        clearEvent(event);
        break;
    case cmOpenQueue:
        openEntity("queue");
        clearEvent(event);
        break;
    case cmOpenEventSources:
        openEntity("event-source");
        clearEvent(event);
        break;
    case cmOpenEventRules:
        openEntity("event-rule");
        clearEvent(event);
        break;
    case cmOpenUserErasure:
        openEntity("user-erasure");
        clearEvent(event);
        break;
    case cmOpenActivationWizard:
        openActivationWizard();
        clearEvent(event);
        break;
    case cmOpenCredentialWizard:
        openCredentialWizard();
        clearEvent(event);
        break;
    case cmOpenJobStream:
        openJobStream();
        clearEvent(event);
        break;
    case cmOpenEncryption:
        openAdminJsonResult({"encryption:status"}, "encryption:status");
        if (messageBox(_("Initialize encryption key if missing?"), mfYesNoCancel | mfConfirmation) == cmYes) {
            openAdminJsonResult({"encryption:init"}, "encryption:init");
        }
        clearEvent(event);
        break;
    case cmOpenPrune:
        openPruneDialog();
        clearEvent(event);
        break;
    case cmOpenJobStatus:
        openAdminJsonResult({"job:status"}, "job:status");
        clearEvent(event);
        break;
    case cmOpenTaskStatus:
        openAdminJsonResult({"task:status"}, "task:status");
        clearEvent(event);
        break;
    case cmOpenQueueOverview:
        openAdminJsonResult({"queue:overview"}, "queue:overview");
        clearEvent(event);
        break;
    case cmOpenTelemetryTest:
        openAdminJsonResult({"telemetry:test"}, "telemetry:test");
        clearEvent(event);
        break;
    case cmOpenStaleRunTemplates:
        openAdminJsonResult({"run-template:stale"}, "run-template:stale");
        clearEvent(event);
        break;
    case cmOpenImportExport:
        openImportExportDialog();
        clearEvent(event);
        break;
    case cmMinimizeAll:
        minimizeAll();
        clearEvent(event);
        break;
    case cmRestoreWindows:
        restoreWindows();
        clearEvent(event);
        break;
    case cmCloseAll:
        closeAllWindows();
        clearEvent(event);
        break;
    case cmLangSystem:
        language_.clear();
        setLanguage("");
        pendingMenuReload_ = true;
        clearEvent(event);
        break;
    case cmLangEnglish:
        language_ = "en";
        setLanguage("en");
        pendingMenuReload_ = true;
        clearEvent(event);
        break;
    case cmLangCzech:
        language_ = "cs";
        setLanguage("cs");
        pendingMenuReload_ = true;
        clearEvent(event);
        break;
    default:
        break;
    }
}

void MultiFlexiApp::idle() {
    TApplication::idle();
    if (pendingStatus_) {
        pendingStatus_ = false;
        openStatus();
    }
    if (pendingMenuReload_) {
        pendingMenuReload_ = false;
        reloadMenuAndStatusLine();
    }
    cliQueue_.pump();
    updateBusyBadge();
    for (auto &t : streamTicks_) {
        if (t.tick) {
            t.tick();
        }
    }
}

} // namespace multiflexitui
