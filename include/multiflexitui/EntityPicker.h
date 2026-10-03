#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/SimpleListViewer.h"

#include <string>
#include <vector>

#include <nlohmann/json.hpp>

namespace multiflexitui {

// Modal picker: list entity rows, Enter selects; returns chosen id via selectedId().
class EntityPickerDialog : public TDialog {
public:
    EntityPickerDialog(CliClient &client, const std::string &cliEntity, const std::string &title,
                       int companyId = 0);

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

    int selectedId() const { return selectedId_; }
    nlohmann::json selectedRow() const { return selectedRow_; }

private:
    void reload();
    void acceptSelection();

    CliClient &client_;
    std::string cliEntity_;
    int companyId_ = 0;
    SimpleListViewer *list_ = nullptr;
    std::vector<nlohmann::json> rows_;
    int selectedId_ = 0;
    nlohmann::json selectedRow_;
};

// Returns selected id, or 0 if cancelled.
int pickEntityId(CliClient &client, const std::string &cliEntity, const std::string &title,
                 int companyId = 0);

// Modal list of action labels; returns chosen index or -1.
int pickActionIndex(const std::vector<std::string> &labels, const std::string &title);

} // namespace multiflexitui
