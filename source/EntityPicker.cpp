#include "multiflexitui/TV.h"
#include "multiflexitui/EntityPicker.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

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

std::string rowLabel(const nlohmann::json &item) {
    const int id = jsonToInt(item, "id");
    std::string name;
    if (item.is_object()) {
        if (item.contains("name")) {
            name = jsonToString(item["name"]);
        } else if (item.contains("login")) {
            name = jsonToString(item["login"]);
        } else if (item.contains("code")) {
            name = jsonToString(item["code"]);
        } else if (item.contains("filename")) {
            name = jsonToString(item["filename"]);
        }
    }
    if (name.empty()) {
        return std::to_string(id);
    }
    return std::to_string(id) + "  " + name;
}

class ActionPickerDialog : public TDialog {
public:
    ActionPickerDialog(const std::vector<std::string> &labels, const std::string &title)
        : TWindowInit(&TDialog::initFrame),
          TDialog(TRect(12, 4, 68, 20), title.c_str()),
          labels_(labels) {
        options |= ofCentered;
        TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
        list_ = new SimpleListViewer(TRect(2, 2, 54, 12), bar);
        growFill(list_);
        insert(list_);
        std::vector<std::string> rows;
        rows.push_back(_("Choose action"));
        for (const auto &l : labels_) {
            rows.push_back(l);
        }
        list_->setRows(std::move(rows));
        TView *ok = new AppButton(TRect(2, 13, 16, 15), _("~O~K"), cmOK, bfDefault);
        stickBottom(ok);
        insert(ok);
        TView *cancel = new AppButton(TRect(40, 13, 54, 15), _("Cancel"), cmCancel, bfNormal);
        stickCorner(cancel);
        insert(cancel);
        selectNext(False);
    }

    int selectedIndex() const { return selected_; }

    void handleEvent(TEvent &event) override {
        TDialog::handleEvent(event);
        auto accept = [&]() {
            if (list_ == nullptr) {
                return;
            }
            const short f = list_->focused;
            if (f <= 0 || static_cast<std::size_t>(f - 1) >= labels_.size()) {
                messageBox(_("Select an item first"), mfError | mfOKButton);
                return;
            }
            selected_ = static_cast<int>(f - 1);
            endModal(cmOK);
        };
        if (event.what == evKeyDown && event.keyDown.keyCode == kbEnter) {
            accept();
            clearEvent(event);
        } else if (event.what == evCommand && event.message.command == cmOK) {
            accept();
            clearEvent(event);
        } else if (event.what == evBroadcast && event.message.command == cmListItemSelected) {
            accept();
            clearEvent(event);
        }
    }

    TColorAttr mapColor(uchar index) override {
        TColorAttr color;
        return windowColor(index, color) ? color : TDialog::mapColor(index);
    }

private:
    std::vector<std::string> labels_;
    SimpleListViewer *list_ = nullptr;
    int selected_ = -1;
};

} // namespace

EntityPickerDialog::EntityPickerDialog(CliClient &client, const std::string &cliEntity,
                                       const std::string &title, int companyId)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(10, 3, 70, 21), title.c_str()),
      client_(client),
      cliEntity_(cliEntity),
      companyId_(companyId) {
    options |= ofCentered;

    TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
    list_ = new SimpleListViewer(TRect(2, 2, 58, 14), bar);
    growFill(list_);
    insert(list_);

    TView *ok = new AppButton(TRect(2, 15, 16, 17), _("~O~K"), cmOK, bfDefault);
    stickBottom(ok);
    insert(ok);
    TView *cancel = new AppButton(TRect(44, 15, 58, 17), _("Cancel"), cmCancel, bfNormal);
    stickCorner(cancel);
    insert(cancel);

    reload();
    selectNext(False);
}

void EntityPickerDialog::reload() {
    ListOptions opts;
    opts.limit = 200;
    opts.offset = 0;
    if (companyId_ > 0) {
        opts.companyId = companyId_;
    }
    auto r = client_.list(cliEntity_, opts);
    rows_.clear();
    std::vector<std::string> lines;
    lines.push_back(_("Select an item"));
    if (!r.ok) {
        lines.push_back(r.errorMessage);
    } else {
        for (const auto &item : asArray(r.data)) {
            rows_.push_back(item);
            lines.push_back(rowLabel(item));
        }
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(lines));
    }
}

void EntityPickerDialog::acceptSelection() {
    if (list_ == nullptr) {
        return;
    }
    const short f = list_->focused;
    if (f <= 0 || static_cast<std::size_t>(f - 1) >= rows_.size()) {
        messageBox(_("Select an item first"), mfError | mfOKButton);
        return;
    }
    selectedRow_ = rows_[static_cast<std::size_t>(f - 1)];
    selectedId_ = jsonToInt(selectedRow_, "id");
    endModal(cmOK);
}

void EntityPickerDialog::handleEvent(TEvent &event) {
    TDialog::handleEvent(event);
    if (event.what == evKeyDown && event.keyDown.keyCode == kbEnter) {
        acceptSelection();
        clearEvent(event);
        return;
    }
    if (event.what == evCommand && event.message.command == cmOK) {
        acceptSelection();
        clearEvent(event);
        return;
    }
    if (event.what == evBroadcast && event.message.command == cmListItemSelected) {
        acceptSelection();
        clearEvent(event);
    }
}

TColorAttr EntityPickerDialog::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

int pickEntityId(CliClient &client, const std::string &cliEntity, const std::string &title, int companyId) {
    auto *dlg = new EntityPickerDialog(client, cliEntity, title, companyId);
    const ushort code = TProgram::deskTop->execView(dlg);
    const int id = (code == cmOK) ? dlg->selectedId() : 0;
    TObject::destroy(dlg);
    return id;
}

int pickActionIndex(const std::vector<std::string> &labels, const std::string &title) {
    if (labels.empty()) {
        return -1;
    }
    auto *dlg = new ActionPickerDialog(labels, title);
    const ushort code = TProgram::deskTop->execView(dlg);
    const int idx = (code == cmOK) ? dlg->selectedIndex() : -1;
    TObject::destroy(dlg);
    return idx;
}

} // namespace multiflexitui
