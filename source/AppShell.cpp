#include "multiflexitui/TV.h"
#include "multiflexitui/AppShell.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/i18n.h"
#include "multiflexitui/StatusView.h"
#include "multiflexitui/AboutView.h"
#include "multiflexitui/EntityListView.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/Wizards.h"
#include "multiflexitui/JobStreamView.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"

#include <algorithm>

namespace multiflexitui {

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

TMenuBar *MultiFlexiApp::initMenuBar(TRect r) {
    r.b.y = r.a.y + 1;

    TSubMenu &multi = *new TSubMenu("~M~ultiFlexi", kbAltM) +
                      *new TMenuItem("~S~tatus", cmShowStatus, kbAltS) +
                      newLine() +
                      *new TMenuItem("E~x~it", cmQuit, kbAltX, hcNoContext, "Alt-X");

    TSubMenu &companies = *new TSubMenu("~C~ompanies", kbAltC) +
                          *new TMenuItem("~L~ist", cmOpenCompanies, kbNoKey);

    TSubMenu &apps = *new TSubMenu("~A~pplications", kbAltA) +
                     *new TMenuItem("~L~ist", cmOpenApplications, kbNoKey) +
                     *new TMenuItem("~C~onfig fields", cmOpenConfFields, kbNoKey) +
                     *new TMenuItem("Company ~a~pps", cmOpenCompanyApps, kbNoKey) +
                     newLine() +
                     *new TMenuItem("Activation ~W~izard", cmOpenActivationWizard, kbNoKey);

    TSubMenu &runtemplates = *new TSubMenu("~R~unTemplates", kbAltR) +
                             *new TMenuItem("~L~ist", cmOpenRunTemplates, kbNoKey) +
                             *new TMenuItem("~S~tale", cmOpenStaleRunTemplates, kbNoKey);

    TSubMenu &creds = *new TSubMenu("Credentia~l~s", kbAltL) +
                      *new TMenuItem("~C~redentials", cmOpenCredentials, kbNoKey) +
                      *new TMenuItem("~T~ypes", cmOpenCredTypes, kbNoKey) +
                      *new TMenuItem("~P~rototypes", cmOpenCrPrototypes, kbNoKey) +
                      newLine() +
                      *new TMenuItem("Credential ~W~izard", cmOpenCredentialWizard, kbNoKey);

    TSubMenu &jobs = *new TSubMenu("~J~obs", kbAltJ) +
                     *new TMenuItem("~J~obs", cmOpenJobs, kbNoKey) +
                     *new TMenuItem("~T~asks", cmOpenTasks, kbNoKey) +
                     *new TMenuItem("~Q~ueue", cmOpenQueue, kbNoKey) +
                     *new TMenuItem("Queue ~o~verview", cmOpenQueueOverview, kbNoKey) +
                     *new TMenuItem("Job st~a~tus", cmOpenJobStatus, kbNoKey) +
                     *new TMenuItem("Task s~t~atus", cmOpenTaskStatus, kbNoKey) +
                     newLine() +
                     *new TMenuItem("Live ~s~tream…", cmOpenJobStream, kbNoKey) +
                     *new TMenuItem("Arti~f~acts", cmOpenArtifacts, kbNoKey);

    TSubMenu &events = *new TSubMenu("~E~vents", kbAltE) +
                       *new TMenuItem("Event ~s~ources", cmOpenEventSources, kbNoKey) +
                       *new TMenuItem("Event ~r~ules", cmOpenEventRules, kbNoKey);

    TSubMenu &admin = *new TSubMenu("Ad~m~in", kbAltD) +
                      *new TMenuItem("~U~sers", cmOpenUsers, kbNoKey) +
                      *new TMenuItem("~T~okens", cmOpenTokens, kbNoKey) +
                      *new TMenuItem("Deletion ~r~equests", cmOpenUserErasure, kbNoKey) +
                      newLine() +
                      *new TMenuItem("~E~ncryption", cmOpenEncryption, kbNoKey) +
                      *new TMenuItem("~P~rune…", cmOpenPrune, kbNoKey) +
                      *new TMenuItem("Telemetry ~t~est", cmOpenTelemetryTest, kbNoKey) +
                      *new TMenuItem("~I~mport/Export app JSON…", cmOpenImportExport, kbNoKey);

    TSubMenu &window = *new TSubMenu("~W~indow", kbAltW) +
                       *new TMenuItem("~T~ile", cmTile, kbNoKey) +
                       *new TMenuItem("C~a~scade", cmCascade, kbNoKey) +
                       *new TMenuItem("~M~inimize all", cmMinimizeAll, kbNoKey) +
                       *new TMenuItem("~R~estore", cmRestoreWindows, kbNoKey) +
                       *new TMenuItem("Close a~l~l", cmCloseAll, kbNoKey);

    TSubMenu &help = *new TSubMenu("~H~elp", kbAltH) +
                     *new TMenuItem("~A~bout", cmShowAbout, kbF1) +
                     newLine() +
                     *new TMenuItem("Language: ~S~ystem", cmLangSystem, kbNoKey) +
                     *new TMenuItem("Language: ~E~nglish", cmLangEnglish, kbNoKey) +
                     *new TMenuItem("Language: ~C~zech", cmLangCzech, kbNoKey);

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

void MultiFlexiApp::openAdminJsonResult(const std::vector<std::string> &args, const std::string &title) {
    auto r = client_.runJson(args);
    if (!r.ok) {
        messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
        return;
    }
    executeDialog(new TextViewerDialog(title, r.data.dump(2)));
}

void MultiFlexiApp::openPruneDialog() {
    char buf[16] = "1000";
    if (inputBox(_("Prune"), _("Keep latest N jobs/logs"), buf, sizeof(buf) - 1) != cmOK) {
        return;
    }
    openAdminJsonResult({"prune", "--keep=" + std::string(buf)}, "prune");
}

void MultiFlexiApp::openImportExportDialog() {
    char path[256] = "/tmp/app.json";
    if (inputBox(_("App JSON"), _("Path for import (existing) or export target"), path, sizeof(path) - 1) != cmOK) {
        return;
    }
    if (messageBox(_("Import from this path? (No = export selected later via Apps Export action)"),
                   mfYesNoCancel | mfConfirmation) == cmYes) {
        openAdminJsonResult({"application:import-json", "--file=" + std::string(path)}, "import-json");
    } else if (messageBox(_("Validate JSON at path instead?"), mfYesNoCancel | mfConfirmation) == cmYes) {
        openAdminJsonResult({"application:validate-json", "--file=" + std::string(path)}, "validate-json");
    }
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
        // No web QR for MultiFlexi CLI commands; ignore status-line clicks.
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
    if (streamTick_) {
        streamTick_();
    }
}

} // namespace multiflexitui
