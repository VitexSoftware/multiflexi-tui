#include "multiflexitui/TV.h"
#include "multiflexitui/Wizards.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

#include <cstring>
#include <sstream>

namespace multiflexitui {

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
        for (const auto &item : asJsonArray(r.data)) {
            companies_.push_back(item);
        }
    }
}

void ActivationWizard::loadApplications() {
    apps_.clear();
    auto r = client_.list("application", 200, 0);
    if (r.ok) {
        for (const auto &item : asJsonArray(r.data)) {
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
        title << _("Assign company-app (Next runs assign)");
        rows[0] = title.str();
        rows.push_back("company_id=" + std::to_string(companyId_));
        rows.push_back("app_id=" + std::to_string(appId_));
        break;
    case 4:
        title << _("Optional: schedule? Type interval (n/h/d/w) or leave blank");
        rows[0] = title.str();
        rows.push_back("run-template id=" + std::to_string(runTemplateId_));
        rows.push_back(_("Enter interval in the input field, then Next to schedule (or leave empty to skip)"));
        break;
    case 5:
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

bool ActivationWizard::advanceFrom(int step) {
    if (step == 1) {
        const short f = list_->focused;
        if (f > 0 && static_cast<std::size_t>(f - 1) < companies_.size()) {
            companyId_ = jsonToInt(companies_[static_cast<std::size_t>(f - 1)], "id");
        }
        if (companyId_ <= 0) {
            messageBox(_("Select a company first"), mfError | mfOKButton);
            return false;
        }
    } else if (step == 2) {
        const short f = list_->focused;
        if (f > 0 && static_cast<std::size_t>(f - 1) < apps_.size()) {
            appId_ = jsonToInt(apps_[static_cast<std::size_t>(f - 1)], "id");
        }
        if (appId_ <= 0) {
            messageBox(_("Select an application first"), mfError | mfOKButton);
            return false;
        }
    } else if (step == 3) {
        auto r = client_.runJson({"company-app:assign", "--company_id=" + std::to_string(companyId_),
                                  "--app_id=" + std::to_string(appId_)});
        if (!r.ok) {
            messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
            return false;
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
    } else if (step == 4) {
        char buf[64] = {};
        extra_->getData(buf);
        interval_ = buf;
        if (!interval_.empty() && runTemplateId_ > 0) {
            // Persist interval on the run-template then schedule.
            auto u = client_.runJson({"run-template:update", "--id=" + std::to_string(runTemplateId_),
                                      "--interv=" + interval_, "--active=1"});
            if (!u.ok) {
                messageBox(u.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            auto r = client_.runJson({"run-template:schedule", "--id=" + std::to_string(runTemplateId_)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            summary_ += "; scheduled interval=" + interval_;
        }
    }
    return true;
}

void ActivationWizard::next() {
    if (!advanceFrom(step_)) {
        return;
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
        const int before = step_;
        if (!advanceFrom(step_)) {
            renderStep();
            return;
        }
        ++step_;
        if (step_ == before) {
            break;
        }
    }
    renderStep();
    if (step_ >= kSteps) {
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
        for (const auto &item : asJsonArray(cr.data)) {
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
                for (const auto &item : asJsonArray(r.data)) {
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
        title << _("Enter field values (Next prompts for each)");
        rows.push_back(title.str());
        rows.push_back("credential id=" + std::to_string(credentialId_));
        rows.push_back("prototype=" + prototypeCode_);
        rows.push_back(_("Press Next to enter secret/field values"));
        break;
    case 6:
        title << _("Done");
        rows.push_back(title.str());
        rows.push_back(_("Credential created with field values."));
        break;
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(rows));
    }
}

bool CredentialWizard::collectAndSaveFields() {
    if (credentialId_ <= 0) {
        return false;
    }
    // Resolve prototype fields from type → prototype code.
    auto typeR = client_.get("credential-type", typeId_);
    if (typeR.ok && typeR.data.is_object()) {
        prototypeCode_ = jsonToString(typeR.data.contains("prototype") ? typeR.data["prototype"]
                                                                       : nlohmann::json());
    }
    fieldDefs_.clear();
    if (!prototypeCode_.empty()) {
        auto pr = client_.runJson({"credential-prototype:get", "--code=" + prototypeCode_});
        if (pr.ok && pr.data.is_object() && pr.data.contains("fields") && pr.data["fields"].is_array()) {
            for (const auto &f : pr.data["fields"]) {
                fieldDefs_.push_back(f);
            }
        }
    }
    for (const auto &f : fieldDefs_) {
        const std::string keyword = jsonToString(f.contains("keyword") ? f["keyword"] : nlohmann::json());
        const std::string name = jsonToString(f.contains("name") ? f["name"] : nlohmann::json());
        if (keyword.empty()) {
            continue;
        }
        char buf[256] = {};
        const std::string hint = jsonToString(f.contains("default_value") ? f["default_value"] : nlohmann::json());
        if (!hint.empty()) {
            std::strncpy(buf, hint.c_str(), sizeof(buf) - 1);
        }
        const std::string label = name.empty() ? keyword : (name + " (" + keyword + ")");
        if (inputBox(_("Credential field"), label.c_str(), buf, sizeof(buf) - 1) != cmOK) {
            continue;
        }
        if (buf[0] == '\0') {
            continue;
        }
        auto r = client_.runJson({"credential:update", "--id=" + std::to_string(credentialId_),
                                  "--field=" + keyword + ":" + std::string(buf)});
        if (!r.ok) {
            messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
            return false;
        }
    }
    return true;
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
        credentialId_ = jsonToInt(r.data, "id");
        if (credentialId_ <= 0 && r.data.is_object() && r.data.contains("credential")) {
            credentialId_ = jsonToInt(r.data["credential"], "id");
        }
    } else if (step_ == 5) {
        if (!collectAndSaveFields()) {
            // still allow finishing without all fields
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
    while (step_ < kSteps) {
        const int before = step_;
        next();
        if (step_ == before) {
            return;
        }
    }
    if (step_ >= kSteps) {
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
