#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/SimpleListViewer.h"

#include <chrono>
#include <string>

namespace multiflexitui {

class JobStreamView : public TWindow {
public:
    JobStreamView(CliClient &client, int jobId);
    ~JobStreamView() override;

    void handleEvent(TEvent &event) override;
    void draw() override;
    TColorAttr mapColor(uchar index) override;

    void tick();
    void refresh(bool force = false);

private:
    void applyJobResult(const CliClient::Result &r, bool force);

    CliClient &client_;
    int jobId_;
    SimpleListViewer *out_ = nullptr;
    TStaticText *meta_ = nullptr;
    std::string stdout_;
    std::string stderr_;
    std::string metaText_;
    bool follow_ = true;
    bool finished_ = false;
    bool pending_ = false;
    std::chrono::steady_clock::time_point lastPoll_{};
};

void openJobStreamPicker(TProgram *app, CliClient &client);
void openJobStreamForId(TProgram *app, CliClient &client, int jobId);

} // namespace multiflexitui
