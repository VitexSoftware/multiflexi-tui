#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/CliClient.h"
#include "multiflexitui/SimpleListViewer.h"

#include <string>
#include <vector>

#include <nlohmann/json.hpp>

namespace multiflexitui {

// Multi-step activation wizard approximating the web ActivationWizard.
class ActivationWizard : public TDialog {
public:
    explicit ActivationWizard(CliClient &client);

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

private:
    void renderStep();
    void next();
    void prev();
    void finish();
    void loadCompanies();
    void loadApplications();

    CliClient &client_;
    int step_ = 1;
    static constexpr int kSteps = 7;
    SimpleListViewer *list_ = nullptr;
    TStaticText *prompt_ = nullptr;
    TInputLine *extra_ = nullptr;
    std::vector<nlohmann::json> companies_;
    std::vector<nlohmann::json> apps_;
    int companyId_ = 0;
    int appId_ = 0;
    int runTemplateId_ = 0;
    std::string summary_;
};

class CredentialWizard : public TDialog {
public:
    explicit CredentialWizard(CliClient &client);

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

private:
    void renderStep();
    void next();
    void prev();
    void finish();

    CliClient &client_;
    int step_ = 1;
    static constexpr int kSteps = 5;
    SimpleListViewer *list_ = nullptr;
    TStaticText *prompt_ = nullptr;
    TInputLine *nameInput_ = nullptr;
    std::vector<nlohmann::json> companies_;
    std::vector<nlohmann::json> types_;
    int companyId_ = 0;
    int typeId_ = 0;
    std::string credName_;
};

} // namespace multiflexitui
