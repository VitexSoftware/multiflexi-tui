#include "multiflexitui/TV.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

#include <algorithm>
#include <cstring>
#include <map>
#include <sstream>

namespace multiflexitui {

namespace {

std::string padLabel(const std::string &s) {
    std::string out = s;
    while (out.size() < 22) {
        out.push_back(' ');
    }
    return out;
}

} // namespace

EntityDetailView::EntityDetailView(CliClient &client, const EntityDef &def, nlohmann::json row)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(4, 2, 76, 22), (def.name + " detail").c_str()),
      client_(client),
      def_(def),
      row_(std::move(row)) {
    options |= ofCentered;
    makeMaximizable(*this);

    TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
    fields_ = new SimpleListViewer(TRect(2, 2, 72, 15), bar);
    growFill(fields_);
    insert(fields_);

    int x = 2;
    auto addBtn = [&](const char *title, ushort cmd) {
        TView *b = new AppButton(TRect(x, 16, x + 14, 18), title, cmd, bfNormal);
        stickBottom(b);
        insert(b);
        x += 15;
    };
    addBtn(_("~R~efresh"), cmEntityRefresh);
    if (def_.canEdit) {
        addBtn(_("~E~dit"), cmEntityEdit);
    }
    for (std::size_t i = 0; i < def_.rowActions.size() && i < 4; ++i) {
        std::string label = def_.rowActions[i].label;
        if (label.size() > 10) {
            label = label.substr(0, 10);
        }
        TView *b = new AppButton(TRect(x, 16, x + 14, 18), label.c_str(),
                                 static_cast<ushort>(cmEntityAction + i), bfNormal);
        stickBottom(b);
        insert(b);
        x += 15;
    }
    TView *close = new AppButton(TRect(58, 16, 72, 18), _("Close"), cmCancel, bfDefault);
    stickCorner(close);
    insert(close);

    reload();
    selectNext(False);
}

void EntityDetailView::reload() {
    const int id = jsonToInt(row_, def_.idKey);
    if (id > 0) {
        auto r = client_.get(def_.cliEntity, id);
        if (r.ok) {
            row_ = r.data;
        }
    }

    std::vector<std::string> lines;
    if (row_.is_object()) {
        for (auto it = row_.begin(); it != row_.end(); ++it) {
            if (it.key() == "stdout" || it.key() == "stderr" || it.key() == "env") {
                const std::string preview = jsonToString(it.value());
                lines.push_back(padLabel(it.key()) + (preview.size() > 40 ? preview.substr(0, 40) + "..." : preview));
            } else {
                lines.push_back(padLabel(it.key()) + jsonToString(it.value()));
            }
        }
    } else {
        lines.push_back(row_.dump(2));
    }
    if (fields_ != nullptr) {
        fields_->setRows(std::move(lines));
    }
}

void EntityDetailView::runAction(std::size_t index) {
    if (index >= def_.rowActions.size()) {
        return;
    }
    if (def_.rowActions[index].run) {
        def_.rowActions[index].run(client_, row_);
        reload();
    }
}

void EntityDetailView::handleEvent(TEvent &event) {
    TDialog::handleEvent(event);
    if (event.what == evCommand) {
        if (event.message.command == cmEntityRefresh) {
            reload();
            clearEvent(event);
        } else if (event.message.command == cmEntityEdit && def_.canEdit) {
            if (owner->execView(new EntityForm(client_, def_, false, row_)) == cmOK) {
                reload();
            }
            clearEvent(event);
        } else if (event.message.command >= cmEntityAction &&
                   event.message.command < cmEntityAction + 50) {
            runAction(static_cast<std::size_t>(event.message.command - cmEntityAction));
            clearEvent(event);
        }
    }
}

TColorAttr EntityDetailView::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

EntityForm::EntityForm(CliClient &client, const EntityDef &def, bool createMode, nlohmann::json existing)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(8, 2, 72, 22), createMode ? (def.name + " — New").c_str() : (def.name + " — Edit").c_str()),
      client_(client),
      def_(def),
      createMode_(createMode),
      existing_(std::move(existing)) {
    options |= ofCentered;
    fields_ = createMode_ ? def_.createFields : def_.editFields;

    int y = 2;
    for (const auto &f : fields_) {
        if (y >= 16) {
            break;
        }
        insert(new TLabel(TRect(2, y, 24, y + 1), f.label.c_str(), nullptr));
        TInputLine *input = new TInputLine(TRect(25, y, 60, y + 1), 255);
        std::string value;
        if (!createMode_ && existing_.is_object() && existing_.contains(f.key)) {
            value = jsonToString(existing_[f.key]);
        } else if (createMode_ && !f.placeholder.empty()) {
            // placeholder doubles as optional default for create forms
            value = f.placeholder;
        }
        if (!value.empty()) {
            char buf[256];
            std::strncpy(buf, value.c_str(), sizeof(buf) - 1);
            buf[sizeof(buf) - 1] = '\0';
            input->setData(buf);
        }
        insert(input);
        inputs_.push_back(input);
        ++y;
    }

    TView *save = new AppButton(TRect(2, 18, 16, 20), _("~S~ave"), cmEntityFormOk, bfDefault);
    stickBottom(save);
    insert(save);
    TView *cancel = new AppButton(TRect(48, 18, 62, 20), _("Cancel"), cmCancel, bfNormal);
    stickCorner(cancel);
    insert(cancel);
    selectNext(False);
}

void EntityForm::submit() {
    std::map<std::string, std::string> values;
    for (std::size_t i = 0; i < fields_.size() && i < inputs_.size(); ++i) {
        char buf[256] = {};
        inputs_[i]->getData(buf);
        values[fields_[i].key] = buf;
        if (fields_[i].required && values[fields_[i].key].empty()) {
            messageBox((fields_[i].label + " is required").c_str(), mfError | mfOKButton);
            return;
        }
    }

    CliClient::Result r;
    if (createMode_) {
        r = client_.create(def_.cliEntity, buildCliArgs(fields_, values));
    } else {
        const int id = jsonToInt(existing_, def_.idKey);
        r = client_.update(def_.cliEntity, buildCliArgs(fields_, values, id));
    }
    if (!r.ok) {
        messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
        return;
    }
    saved_ = true;
    endModal(cmOK);
}

void EntityForm::handleEvent(TEvent &event) {
    TDialog::handleEvent(event);
    if (event.what == evCommand && event.message.command == cmEntityFormOk) {
        submit();
        clearEvent(event);
    }
}

TColorAttr EntityForm::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

TextViewerDialog::TextViewerDialog(const std::string &title, const std::string &body)
    : TWindowInit(&TDialog::initFrame), TDialog(TRect(2, 1, 78, 23), title.c_str()) {
    options |= ofCentered;
    makeMaximizable(*this);

    TScrollBar *vBar = standardScrollBar(sbVertical | sbHandleKeyboard);
    TScrollBar *hBar = standardScrollBar(sbHorizontal | sbHandleKeyboard);
    auto *viewer = new SimpleListViewer(TRect(1, 1, 75, 19), vBar, hBar);
    growFill(viewer);

    std::vector<std::string> lines;
    std::istringstream in(body);
    std::string line;
    while (std::getline(in, line)) {
        lines.push_back(line);
    }
    if (lines.empty()) {
        lines.push_back("(empty)");
    }
    viewer->setRows(std::move(lines));
    insert(viewer);

    TView *close = new AppButton(TRect(60, 20, 74, 22), _("Close"), cmCancel, bfDefault);
    stickCorner(close);
    insert(close);
    selectNext(False);
}

TColorAttr TextViewerDialog::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

} // namespace multiflexitui
