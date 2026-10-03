#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/SimpleListViewer.h"

#include <vector>

#include <nlohmann/json.hpp>

namespace multiflexitui {

class EntityListView : public TWindow {
public:
    EntityListView(CliClient &client, const EntityDef &def);
    ~EntityListView() override;

    void handleEvent(TEvent &event) override;
    void draw() override;
    TColorAttr mapColor(uchar index) override;

    void refresh();
    void openDetail();
    void openCreate();
    void openEdit();
    void doDelete();
    void runRowAction(std::size_t index);
    void runListAction(std::size_t index);
    void openActionsMenu();

private:
    void setStatus(const std::string &text);
    void applyFilterFromInput();
    nlohmann::json selectedRow() const;
    void applyListResult(const CliClient::Result &result);

    CliClient &client_;
    EntityDef def_;
    SimpleListViewer *list_ = nullptr;
    TInputLine *filterInput_ = nullptr;
    TStaticText *status_ = nullptr;
    std::vector<nlohmann::json> rows_;
    int offset_ = 0;
    std::string filter_;
    std::string statusText_;
    bool loading_ = false;
};

void openEntityList(TProgram *app, CliClient &client, const std::string &cliEntity);

} // namespace multiflexitui
