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

private:
    void setStatus(const std::string &text);
    nlohmann::json selectedRow() const;

    CliClient &client_;
    EntityDef def_;
    SimpleListViewer *list_ = nullptr;
    TStaticText *status_ = nullptr;
    std::vector<nlohmann::json> rows_;
    int offset_ = 0;
    std::string statusText_;
};

void openEntityList(TProgram *app, CliClient &client, const std::string &cliEntity);

} // namespace multiflexitui
