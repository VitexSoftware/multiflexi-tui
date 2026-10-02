#include "multiflexitui/TV.h"
#include "multiflexitui/Wizards.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

#include <sstream>

namespace multiflexitui {

namespace {

nlohmann::json asArray(const nlohmann::json &data) {
    if (data.is_array()) {
        return data;
    }
    if (data.is_object() && data.contains("data") && data["data"].is_array()) {
        return data["data"];
    }
    return nlohmann::json::array();
}

} // namespace

ActivationWizard::ActivationWizard(CliClient &client)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(6, 2, 74, 22), _("Activation Wizard")),
      client_(client) {
    options |= ofCentered;
    makeMaximizable(*this);

    prompt_ = new TStaticText(TRect(2, 1, 66, 2), "");
    insert(prompt_);

    TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
    list_ = new SimpleListViewer(TRect(2, 3, 66, 14), bar);
    growFill(list_);
    insert(list_);

    extra_ = new TInputLine(TRect(2, 15, 40, 16), 64);
    stickBottomWide(extra_);
    insert(extra_);

    TView *prev = new AppButton(TRect(2, 17, 14, 19), _("~B~ack"), cmWizardPrev, bfNormal);
    stickBottom(prev);
    insert(prev);
    TView *next = new AppButton(TRect(16, 17, 28, 19), _("~N~ext"), cmWizardNext, bfDefault);
    stickBottom(next);
    insert(next);
    TView *fin = new AppButton(TRect(30, 17, 44, 19), _("~F~inish"), cmWizardFinish, bfNormal);
    stickBottom(fin);
    insert(fin);
    TView *cancel = new AppButton(TRect(52, 17, 66, 19), _("Cancel"), cmCancel, bfNormal);
    stickCorner(cancel);
    insert(cancel);

    loadCompanies();
    loadApplications();
    renderStep();
    selectNext(False);
}

void ActivationWizard::loadCompanies() {
    companies_.clear();
    auto r = client_.list("company", 200, 0);
    if (r.ok) {
        for (const auto &item : asArray(r.data)) {
            companies_.push_back(item);
        }
    }
}

void ActivationWizard::loadApplications() {
    apps_.clear();
    auto r = client_.list("application", 200, 0);
    if (r.ok) {
        for (const auto &item : asArray(r.data)) {
            apps_.push_back(item);
        }
    }
}

void ActivationWizard::renderStep() {
    std::ostringstream title;
    title << _("Step") << " " << step_ << "/" << kSteps << ": ";
    std::vector<std::string> rows;
    rows.push_back(title.str());

    switch (step_) {
    case 1:
        title << _("Pick a company");
        rows[0] = title.str();
        for (const auto &c : companies_) {
            rows.push_back(std::to_string(jsonToInt(c, "id")) + "  " + jsonToString(c["name"]));
        }
        break;
    case 2:
        title << _("Pick an application");
        rows[0] = title.str();
        for (const auto &a : apps_) {
            rows.push_back(std::to_string(jsonToInt(a, "id")) + "  " + jsonToString(a["name"]));
        }
        break;
    case 3:
        title << _("Assign company-app (Enter Next to run)");
        rows[0] = title.str();
        rows.push_back("company_id=" + std::to_string(companyId_));
        rows.push_back("app_id=" + std::to_string(appId_));
        break;
    case 4:
        title << _("Optional: schedule interval (n/h/d/w) in input");
        rows[0] = title.str();
        rows.push_back("runtemplate will be created by assign");
        break;
    case 5:
        title << _("Optional: enable & schedule now?");
        rows[0] = title.str();
        rows.push_back("Type 'yes' in the input to schedule after assign");
        break;
    case 6:
        title << _("Confirm");
        rows[0] = title.str();
        rows.push_back("Company " + std::to_string(companyId_) + " + App " + std::to_string(appId_));
        break;
    case 7:
        title << _("Summary");
        rows[0] = title.str();
        rows.push_back(summary_.empty() ? "(not finished yet)" : summary_);
        break;
    default:
        break;
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(rows));
    }
}

void ActivationWizard::next() {
    if (step_ == 1) {
        const short f = list_->focused;
        if (f > 0 && static_cast<std::size_t>(f - 1) < companies_.size()) {
            companyId_ = jsonToInt(companies_[static_cast<std::size_t>(f - 1)], "id");
        }
        if (companyId_ <= 0) {
            messageBox(_("Select a company first"), mfError | mfOKButton);
            return;
        }
    } else if (step_ == 2) {
        const short f = list_->focused;
        if (f > 0 && static_cast<std::size_t>(f - 1) < apps_.size()) {
            appId_ = jsonToInt(apps_[static_cast<std::size_t>(f - 1)], "id");
        }
        if (appId_ <= 0) {
            messageBox(_("Select an application first"), mfError | mfOKButton);
            return;
        }
    } else if (step_ == 3) {
        auto r = client_.runJson({"company-app:assign", "--company_id=" + std::to_string(companyId_),
                                  "--app_id=" + std::to_string(appId_)});
        if (!r.ok) {
            messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
            return;
        }
        if (r.data.is_object() && r.data.contains("runtemplate_id")) {
            runTemplateId_ = jsonToInt(r.data, "runtemplate_id");
        } else if (r.data.is_object() && r.data.contains("id")) {
            runTemplateId_ = jsonToInt(r.data, "id");
        }
        summary_ = "Assigned company " + std::to_string(companyId_) + " app " + std::to_string(appId_);
        if (runTemplateId_ > 0) {
            summary_ += " (run-template " + std::to_string(runTemplateId_) + ")";
        }
    } else if (step_ == 5) {
        char buf[64] = {};
        extra_->getData(buf);
        std::string ans = buf;
        if ((ans == "yes" || ans == "y" || ans == "Y") && runTemplateId_ > 0) {
            auto r = client_.runJson({"run-template:schedule", "--id=" + std::to_string(runTemplateId_)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
            } else {
                summary_ += "; scheduled";
            }
        }
    }

    if (step_ < kSteps) {
        ++step_;
        renderStep();
    }
}

void ActivationWizard::prev() {
    if (step_ > 1) {
        --step_;
        renderStep();
    }
}

void ActivationWizard::finish() {
    while (step_ < kSteps) {
        next();
        if (step_ >= kSteps) {
            break;
        }
        // Avoid infinite loop if next() refused to advance
        if (step_ < kSteps) {
            // force advance only after successful critical steps
            break;
        }
    }
    renderStep();
    if (step_ == kSteps) {
        messageBox(summary_.empty() ? _("Wizard finished") : summary_.c_str(), mfInformation | mfOKButton);
        endModal(cmOK);
    }
}

void ActivationWizard::handleEvent(TEvent &event) {
    TDialog::handleEvent(event);
    if (event.what == evCommand) {
        switch (event.message.command) {
        case cmWizardNext:
            next();
            clearEvent(event);
            break;
        case cmWizardPrev:
            prev();
            clearEvent(event);
            break;
        case cmWizardFinish:
            finish();
            clearEvent(event);
            break;
        default:
            break;
        }
    }
}

TColorAttr ActivationWizard::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

// ---- Credential Wizard ----

CredentialWizard::CredentialWizard(CliClient &client)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(6, 2, 74, 22), _("Credential Wizard")),
      client_(client) {
    options |= ofCentered;
    makeMaximizable(*this);

    prompt_ = new TStaticText(TRect(2, 1, 66, 2), "");
    insert(prompt_);

    TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
    list_ = new SimpleListViewer(TRect(2, 3, 66, 13), bar);
    growFill(list_);
    insert(list_);

    insert(new TLabel(TRect(2, 14, 18, 15), _("Name"), nullptr));
    nameInput_ = new TInputLine(TRect(19, 14, 50, 15), 64);
    stickBottomWide(nameInput_);
    insert(nameInput_);

    TView *prev = new AppButton(TRect(2, 17, 14, 19), _("~B~ack"), cmWizardPrev, bfNormal);
    stickBottom(prev);
    insert(prev);
    TView *next = new AppButton(TRect(16, 17, 28, 19), _("~N~ext"), cmWizardNext, bfDefault);
    stickBottom(next);
    insert(next);
    TView *fin = new AppButton(TRect(30, 17, 44, 19), _("~F~inish"), cmWizardFinish, bfNormal);
    stickBottom(fin);
    insert(fin);
    TView *cancel = new AppButton(TRect(52, 17, 66, 19), _("Cancel"), cmCancel, bfNormal);
    stickCorner(cancel);
    insert(cancel);

    auto cr = client_.list("company", 200, 0);
    if (cr.ok) {
        for (const auto &item : asArray(cr.data)) {
            companies_.push_back(item);
        }
    }
    renderStep();
    selectNext(False);
}

void CredentialWizard::renderStep() {
    std::vector<std::string> rows;
    std::ostringstream title;
    title << _("Step") << " " << step_ << "/" << kSteps << ": ";
    switch (step_) {
    case 1:
        title << _("Pick company");
        rows.push_back(title.str());
        for (const auto &c : companies_) {
            rows.push_back(std::to_string(jsonToInt(c, "id")) + "  " + jsonToString(c["name"]));
        }
        break;
    case 2:
        title << _("Pick credential type");
        rows.push_back(title.str());
        types_.clear();
        {
            auto r = client_.runJson({"credential-type:list", "--company_id=" + std::to_string(companyId_),
                                      "--limit=200", "--offset=0"});
            if (r.ok) {
                for (const auto &item : asArray(r.data)) {
                    types_.push_back(item);
                    rows.push_back(std::to_string(jsonToInt(item, "id")) + "  " + jsonToString(item["name"]));
                }
            } else {
                rows.push_back(r.errorMessage);
            }
        }
        break;
    case 3:
        title << _("Enter credential name in the Name field");
        rows.push_back(title.str());
        rows.push_back("company=" + std::to_string(companyId_) + " type=" + std::to_string(typeId_));
        break;
    case 4:
        title << _("Confirm create");
        rows.push_back(title.str());
        rows.push_back("name=" + credName_);
        break;
    case 5:
        title << _("Done");
        rows.push_back(title.str());
        rows.push_back(_("Credential created (values can be edited later)."));
        break;
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(rows));
    }
}

void CredentialWizard::next() {
    if (step_ == 1) {
        const short f = list_->focused;
        if (f > 0 && static_cast<std::size_t>(f - 1) < companies_.size()) {
            companyId_ = jsonToInt(companies_[static_cast<std::size_t>(f - 1)], "id");
        }
        if (companyId_ <= 0) {
            messageBox(_("Select a company"), mfError | mfOKButton);
            return;
        }
    } else if (step_ == 2) {
        const short f = list_->focused;
        if (f > 0 && static_cast<std::size_t>(f - 1) < types_.size()) {
            typeId_ = jsonToInt(types_[static_cast<std::size_t>(f - 1)], "id");
        }
        if (typeId_ <= 0) {
            messageBox(_("Select a credential type"), mfError | mfOKButton);
            return;
        }
    } else if (step_ == 3) {
        char buf[64] = {};
        nameInput_->getData(buf);
        credName_ = buf;
        if (credName_.empty()) {
            messageBox(_("Name is required"), mfError | mfOKButton);
            return;
        }
    } else if (step_ == 4) {
        auto r = client_.create("credential",
                                {"--name=" + credName_, "--company_id=" + std::to_string(companyId_),
                                 "--credential_type_id=" + std::to_string(typeId_)});
        if (!r.ok) {
            messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
            return;
        }
    }
    if (step_ < kSteps) {
        ++step_;
        renderStep();
    }
}

void CredentialWizard::prev() {
    if (step_ > 1) {
        --step_;
        renderStep();
    }
}

void CredentialWizard::finish() {
    if (step_ < 4) {
        messageBox(_("Complete the steps first"), mfError | mfOKButton);
        return;
    }
    if (step_ == 4) {
        next();
    }
    if (step_ == 5) {
        endModal(cmOK);
    }
}

void CredentialWizard::handleEvent(TEvent &event) {
    TDialog::handleEvent(event);
    if (event.what == evCommand) {
        switch (event.message.command) {
        case cmWizardNext:
            next();
            clearEvent(event);
            break;
        case cmWizardPrev:
            prev();
            clearEvent(event);
            break;
        case cmWizardFinish:
            finish();
            clearEvent(event);
            break;
        default:
            break;
        }
    }
}

TColorAttr CredentialWizard::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

} // namespace multiflexitui
