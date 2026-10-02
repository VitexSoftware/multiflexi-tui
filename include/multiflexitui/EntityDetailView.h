#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/SimpleListViewer.h"

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

    CliClient &client_;
    EntityDef def_;
    nlohmann::json row_;
    SimpleListViewer *fields_ = nullptr;
};

class EntityForm : public TDialog {
public:
    // createMode=true uses createFields; false uses editFields and includes --id.
    EntityForm(CliClient &client, const EntityDef &def, bool createMode, nlohmann::json existing = {});

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

    bool saved() const { return saved_; }

private:
    void submit();

    CliClient &client_;
    EntityDef def_;
    bool createMode_;
    nlohmann::json existing_;
    std::vector<TInputLine *> inputs_;
    std::vector<FieldDef> fields_;
    bool saved_ = false;
};

class TextViewerDialog : public TDialog {
public:
    TextViewerDialog(const std::string &title, const std::string &body);

    TColorAttr mapColor(uchar index) override;
};

} // namespace multiflexitui
