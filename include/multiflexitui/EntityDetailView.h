#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/SimpleListViewer.h"

#include <map>
#include <string>
#include <vector>

#include <nlohmann/json.hpp>

namespace multiflexitui {

class EntityDetailView : public TDialog {
public:
    EntityDetailView(CliClient &client, const EntityDef &def, nlohmann::json row);

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

private:
    void reload();
    void runAction(std::size_t index);
    void openActionsMenu();

    CliClient &client_;
    EntityDef def_;
    nlohmann::json row_;
    SimpleListViewer *fields_ = nullptr;
};

class EntityForm : public TDialog {
public:
    EntityForm(CliClient &client, const EntityDef &def, bool createMode, nlohmann::json existing = {});
    ~EntityForm() override;

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

    bool saved() const { return saved_; }

private:
    void rebuildPage();
    void submit();
    void pickRelation(std::size_t visibleIndex);
    void clearInputs();

    CliClient &client_;
    EntityDef def_;
    bool createMode_;
    nlohmann::json existing_;
    std::vector<FieldDef> fields_;
    std::vector<TInputLine *> inputs_;
    std::vector<TView *> pageViews_;
    std::map<std::string, std::string> values_;
    std::size_t page_ = 0;
    static constexpr std::size_t kFieldsPerPage = 10;
    bool saved_ = false;
};

class TextViewerDialog : public TDialog {
public:
    TextViewerDialog(const std::string &title, const std::string &body);

    TColorAttr mapColor(uchar index) override;
};

// Key/value viewer for flat JSON admin results.
class KeyValueDialog : public TDialog {
public:
    KeyValueDialog(const std::string &title, const nlohmann::json &data);

    TColorAttr mapColor(uchar index) override;
};

} // namespace multiflexitui
