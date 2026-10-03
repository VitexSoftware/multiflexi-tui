#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/AsyncCliQueue.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/AppStatusLine.h"

#include <functional>
#include <string>
#include <vector>

namespace multiflexitui {

class MultiFlexiApp : public TApplication {
public:
    MultiFlexiApp();

    void configure(std::string cliBinary, std::string envFile);

    static TMenuBar *initMenuBar(TRect r);
    static TStatusLine *initStatusLine(TRect r);

    void handleEvent(TEvent &event) override;
    void idle() override;

    CliClient &client() { return client_; }
    AsyncCliQueue &cliQueue() { return cliQueue_; }

    void setBusy(bool busy);
    void setActiveCompany(int id, const std::string &name);
    int activeCompanyId() const { return activeCompanyId_; }
    const std::string &activeCompanyName() const { return activeCompanyName_; }

    void addStreamTick(void *owner, std::function<void()> tick);
    void removeStreamTick(void *owner);

private:
    void openStatus();
    void openAbout();
    void openEntity(const std::string &cliEntity);
    void openActivationWizard();
    void openCredentialWizard();
    void openJobStream();
    void openSetCompany();
    void openAdminJsonResult(const std::vector<std::string> &args, const std::string &title);
    void openPruneDialog();
    void openImportExportDialog();
    void reloadMenuAndStatusLine();
    void minimizeAll();
    void restoreWindows();
    void closeAllWindows();
    void updateBusyBadge();

    struct MinimizedWindow {
        TWindow *window;
        TRect bounds;
    };

    struct StreamTick {
        void *owner = nullptr;
        std::function<void()> tick;
    };

    CliClient client_;
    AsyncCliQueue cliQueue_;
    AppStatusLine *statusLine_ = nullptr;
    bool pendingStatus_ = true;
    bool pendingMenuReload_ = false;
    bool busy_ = false;
    std::string language_;
    std::vector<MinimizedWindow> minimized_;
    std::vector<StreamTick> streamTicks_;
    int activeCompanyId_ = 0;
    std::string activeCompanyName_;
};

} // namespace multiflexitui
