#pragma once

#include "multiflexitui/TV.h"
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

private:
    void openStatus();
    void openAbout();
    void openEntity(const std::string &cliEntity);
    void openActivationWizard();
    void openCredentialWizard();
    void openJobStream();
    void openAdminJsonResult(const std::vector<std::string> &args, const std::string &title);
    void openPruneDialog();
    void openImportExportDialog();
    void reloadMenuAndStatusLine();
    void minimizeAll();
    void restoreWindows();
    void closeAllWindows();

    struct MinimizedWindow {
        TWindow *window;
        TRect bounds;
    };

    CliClient client_;
    AppStatusLine *statusLine_ = nullptr;
    bool pendingStatus_ = true;
    bool pendingMenuReload_ = false;
    std::string language_;
    std::vector<MinimizedWindow> minimized_;

    // Job stream polling hook (set by JobStreamView).
public:
    std::function<void()> streamTick_;
};

} // namespace multiflexitui
