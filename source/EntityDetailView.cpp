#include "multiflexitui/TV.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/EntityPicker.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

#include <algorithm>
#include <cstring>
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

// Password-style input: stores plaintext but draws asterisks.
class PasswordInputLine : public TInputLine {
public:
    PasswordInputLine(const TRect &bounds, int limit) noexcept : TInputLine(bounds, limit) {
    }

    void draw() override {
        if (data == nullptr) {
            TInputLine::draw();
            return;
        }
        const std::string real(data);
        for (std::size_t i = 0; i < real.size(); ++i) {
            data[i] = '*';
        }
        data[real.size()] = '\0';
        TInputLine::draw();
        std::memcpy(data, real.data(), real.size());
        data[real.size()] = '\0';
    }
};

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
    addBtn(_("Act~i~ons"), cmEntityActionsMenu);
    TView *close = new AppButton(TRect(58, 16, 72, 18), _("Close"), cmCancel, bfDefault);
    stickCorner(close);
    insert(close);

    reload();
    selectNext(False);
}

void EntityDetailView::reload() {
    if (def_.supportsGet) {
        const int id = jsonToInt(row_, def_.idKey);
        if (id > 0) {
            auto r = client_.get(def_.cliEntity, id);
            if (r.ok) {
                row_ = r.data;
            }
        }
    }

    std::vector<std::string> lines;
    if (row_.is_object()) {
        for (auto it = row_.begin(); it != row_.end(); ++it) {
            if (it.key() == "stdout" || it.key() == "stderr" || it.key() == "env") {
                const std::string preview = jsonToString(it.value());
                lines.push_back(padLabel(it.key()) +
                                (preview.size() > 40 ? preview.substr(0, 40) + "..." : preview));
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

void EntityDetailView::openActionsMenu() {
    std::vector<std::string> labels;
    for (const auto &a : def_.rowActions) {
        labels.push_back(a.label);
    }
    if (labels.empty()) {
        messageBox(_("No actions available"), mfInformation | mfOKButton);
        return;
    }
    const int idx = pickActionIndex(labels, _("Actions"));
    if (idx >= 0) {
        runAction(static_cast<std::size_t>(idx));
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
        } else if (event.message.command == cmEntityActionsMenu) {
            openActionsMenu();
            clearEvent(event);
        } else if (event.message.command >= cmEntityAction &&
                   event.message.command < cmEntityAction + 50) {
            runAction(static_cast<std::size_t>(event.message.command - cmEntityAction));
            clearEvent(event);
        }
    }
    if (event.what == evKeyDown && event.keyDown.keyCode == kbF2) {
        openActionsMenu();
        clearEvent(event);
    }
}

TColorAttr EntityDetailView::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

EntityForm::EntityForm(CliClient &client, const EntityDef &def, bool createMode, nlohmann::json existing)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(6, 1, 74, 23), createMode ? (def.name + " — New").c_str() : (def.name + " — Edit").c_str()),
      client_(client),
      def_(def),
      createMode_(createMode),
      existing_(std::move(existing)) {
    options |= ofCentered;
    fields_ = createMode_ ? def_.createFields : def_.editFields;

    for (const auto &f : fields_) {
        std::string value;
        if (!createMode_ && existing_.is_object() && existing_.contains(f.key)) {
            value = jsonToString(existing_[f.key]);
        } else if (createMode_ && !f.placeholder.empty()) {
            value = f.placeholder;
        }
        values_[f.key] = value;
    }

    TView *save = new AppButton(TRect(2, 20, 14, 22), _("~S~ave"), cmEntityFormOk, bfDefault);
    stickBottom(save);
    insert(save);
    TView *more = new AppButton(TRect(16, 20, 30, 22), _("~M~ore"), cmEntityFormMore, bfNormal);
    stickBottom(more);
    insert(more);
    TView *prev = new AppButton(TRect(32, 20, 46, 22), _("~P~rev"), cmEntityFormPrevFields, bfNormal);
    stickBottom(prev);
    insert(prev);
    TView *cancel = new AppButton(TRect(52, 20, 66, 22), _("Cancel"), cmCancel, bfNormal);
    stickCorner(cancel);
    insert(cancel);

    rebuildPage();
    selectNext(False);
}

EntityForm::~EntityForm() = default;

void EntityForm::clearInputs() {
    for (TView *v : pageViews_) {
        if (v != nullptr) {
            remove(v);
            TObject::destroy(v);
        }
    }
    pageViews_.clear();
    inputs_.clear();
}

void EntityForm::rebuildPage() {
    // Persist current page values first.
    for (std::size_t i = 0; i < inputs_.size(); ++i) {
        const std::size_t fieldIndex = page_ * kFieldsPerPage + i;
        if (fieldIndex >= fields_.size()) {
            break;
        }
        char buf[256] = {};
        inputs_[i]->getData(buf);
        values_[fields_[fieldIndex].key] = buf;
    }
    clearInputs();

    const std::size_t start = page_ * kFieldsPerPage;
    if (start >= fields_.size() && page_ > 0) {
        --page_;
    }
    const std::size_t begin = page_ * kFieldsPerPage;
    int y = 2;
    for (std::size_t i = begin; i < fields_.size() && i < begin + kFieldsPerPage; ++i) {
        const auto &f = fields_[i];
        auto *label = new TLabel(TRect(2, y, 24, y + 1), f.label.c_str(), nullptr);
        insert(label);
        pageViews_.push_back(label);

        TInputLine *input = nullptr;
        if (f.secret) {
            input = new PasswordInputLine(TRect(25, y, 54, y + 1), 255);
        } else {
            input = new TInputLine(TRect(25, y, 54, y + 1), 255);
        }
        const std::string &value = values_[f.key];
        if (!value.empty()) {
            char buf[256];
            std::strncpy(buf, value.c_str(), sizeof(buf) - 1);
            buf[sizeof(buf) - 1] = '\0';
            input->setData(buf);
        }
        insert(input);
        pageViews_.push_back(input);
        inputs_.push_back(input);

        if (!f.relationEntity.empty()) {
            auto *pick = new AppButton(TRect(55, y, 66, y + 2), _("~F2~"),
                                       static_cast<ushort>(cmEntityFormPick + static_cast<ushort>(i - begin)),
                                       bfNormal);
            insert(pick);
            pageViews_.push_back(pick);
        }
        ++y;
    }
}

void EntityForm::pickRelation(std::size_t visibleIndex) {
    const std::size_t fieldIndex = page_ * kFieldsPerPage + visibleIndex;
    if (fieldIndex >= fields_.size() || visibleIndex >= inputs_.size()) {
        return;
    }
    const auto &f = fields_[fieldIndex];
    if (f.relationEntity.empty()) {
        return;
    }
    const int id = pickEntityId(client_, f.relationEntity, f.label + " — pick");
    if (id <= 0) {
        return;
    }
    const std::string text = std::to_string(id);
    values_[f.key] = text;
    char buf[256];
    std::strncpy(buf, text.c_str(), sizeof(buf) - 1);
    buf[sizeof(buf) - 1] = '\0';
    inputs_[visibleIndex]->setData(buf);
    inputs_[visibleIndex]->drawView();
}

void EntityForm::submit() {
    // Flush current page.
    for (std::size_t i = 0; i < inputs_.size(); ++i) {
        const std::size_t fieldIndex = page_ * kFieldsPerPage + i;
        if (fieldIndex >= fields_.size()) {
            break;
        }
        char buf[256] = {};
        inputs_[i]->getData(buf);
        values_[fields_[fieldIndex].key] = buf;
    }

    for (const auto &f : fields_) {
        if (f.required && values_[f.key].empty()) {
            messageBox((f.label + " is required").c_str(), mfError | mfOKButton);
            return;
        }
    }

    CliClient::Result r;
    if (createMode_) {
        r = client_.create(def_.cliEntity, buildCliArgs(fields_, values_));
    } else {
        const int id = jsonToInt(existing_, def_.idKey);
        r = client_.update(def_.cliEntity, buildCliArgs(fields_, values_, id));
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
    if (event.what == evKeyDown && event.keyDown.keyCode == kbF2) {
        // Pick relation for focused input if possible.
        for (std::size_t i = 0; i < inputs_.size(); ++i) {
            if (inputs_[i]->state & sfFocused) {
                pickRelation(i);
                clearEvent(event);
                return;
            }
        }
    }
    if (event.what == evCommand) {
        if (event.message.command == cmEntityFormOk) {
            submit();
            clearEvent(event);
        } else if (event.message.command == cmEntityFormMore) {
            const std::size_t maxPage = fields_.empty() ? 0 : (fields_.size() - 1) / kFieldsPerPage;
            if (page_ < maxPage) {
                ++page_;
                rebuildPage();
                drawView();
            }
            clearEvent(event);
        } else if (event.message.command == cmEntityFormPrevFields) {
            if (page_ > 0) {
                --page_;
                rebuildPage();
                drawView();
            }
            clearEvent(event);
        } else if (event.message.command >= cmEntityFormPick &&
                   event.message.command < cmEntityFormPick + 20) {
            pickRelation(static_cast<std::size_t>(event.message.command - cmEntityFormPick));
            clearEvent(event);
        }
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

KeyValueDialog::KeyValueDialog(const std::string &title, const nlohmann::json &data)
    : TWindowInit(&TDialog::initFrame), TDialog(TRect(6, 2, 74, 22), title.c_str()) {
    options |= ofCentered;
    makeMaximizable(*this);

    TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
    auto *viewer = new SimpleListViewer(TRect(2, 2, 66, 16), bar);
    growFill(viewer);

    std::vector<std::string> lines;
    if (data.is_object()) {
        for (auto it = data.begin(); it != data.end(); ++it) {
            std::string key = it.key();
            while (key.size() < 22) {
                key.push_back(' ');
            }
            lines.push_back(key + jsonToString(it.value()));
        }
    } else {
        lines.push_back(data.dump(2));
    }
    if (lines.empty()) {
        lines.push_back("(empty)");
    }
    viewer->setRows(std::move(lines));
    insert(viewer);

    TView *close = new AppButton(TRect(52, 17, 66, 19), _("Close"), cmCancel, bfDefault);
    stickCorner(close);
    insert(close);
    selectNext(False);
}

TColorAttr KeyValueDialog::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

} // namespace multiflexitui
