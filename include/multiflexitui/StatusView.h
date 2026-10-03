#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/SimpleListViewer.h"
#include "multiflexitui/AppButton.h"

#include <cstdint>
#include <string>
#include <vector>

namespace multiflexitui {

// List viewer that can tint individual rows (used for service state colors).
class StatusListViewer : public SimpleListViewer {
public:
    using SimpleListViewer::SimpleListViewer;

    void setRows(std::vector<std::string> rows, std::vector<uint32_t> rowFg);
    void draw() override;

private:
    std::vector<uint32_t> rowFg_; // 0 = default list color
};

class StatusView : public TDialog {
public:
    explicit StatusView(CliClient &client);

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

private:
    void reload();
    void updateServiceButtons();
    std::string focusedKey() const;
    bool focusedIsService() const;
    void controlService(const std::string &action);

    CliClient &client_;
    StatusListViewer *list_ = nullptr;
    AppButton *startBtn_ = nullptr;
    AppButton *stopBtn_ = nullptr;
    std::vector<std::string> keys_;
};

} // namespace multiflexitui
