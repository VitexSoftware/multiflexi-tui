#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/SimpleListViewer.h"

namespace multiflexitui {

class StatusView : public TDialog {
public:
    explicit StatusView(CliClient &client);

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

private:
    void reload();

    CliClient &client_;
    SimpleListViewer *list_ = nullptr;
};

} // namespace multiflexitui
