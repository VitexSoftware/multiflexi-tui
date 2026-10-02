#include "multiflexitui/TV.h"
#include "multiflexitui/EntityListView.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

#include <sstream>

namespace multiflexitui {

EntityListView::EntityListView(CliClient &client, const EntityDef &def)
    : TWindowInit(&TWindow::initFrame),
      TWindow(TRect(1, 1, 78, 23), def.name.c_str(), wnNoNumber),
      client_(client),
      def_(def) {
    flags |= wfGrow | wfZoom | wfClose;
    options |= ofTileable;

    status_ = new TStaticText(TRect(2, 1, 76, 2), "");
    growWide(status_);
    insert(status_);

    TScrollBar *vBar = standardScrollBar(sbVertical | sbHandleKeyboard);
    TScrollBar *hBar = standardScrollBar(sbHorizontal | sbHandleKeyboard);
    list_ = new SimpleListViewer(TRect(1, 2, 76, 18), vBar, hBar);
    growFill(list_);
    insert(list_);

    int x = 2;
    auto addBtn = [&](const char *title, ushort cmd) {
        TView *b = new AppButton(TRect(x, 19, x + 12, 21), title, cmd, bfNormal);
        stickBottom(b);
        insert(b);
        x += 13;
    };
    addBtn(_("~R~efresh"), cmEntityRefresh);
    if (def_.canCreate) {
        addBtn(_("~N~ew"), cmEntityNew);
    }
    if (def_.canEdit) {
        addBtn(_("~E~dit"), cmEntityEdit);
    }
    if (!def_.deleteAction.empty()) {
        addBtn(_("~D~elete"), cmEntityDelete);
    }
    addBtn(_("~P~rev"), cmEntityPrevPage);
    addBtn(_("Nex~t~"), cmEntityNextPage);

    // List-level action buttons (first few)
    for (std::size_t i = 0; i < def_.listActions.size() && i < 3; ++i) {
        const std::string label = "~" + std::to_string(i + 1) + "~ " + def_.listActions[i].label;
        TView *b = new AppButton(TRect(x, 19, x + 16, 21), label.c_str(),
                                 static_cast<ushort>(cmEntityAction + 50 + i), bfNormal);
        stickBottom(b);
        insert(b);
        x += 17;
    }

    selectNext(False);
    refresh();
}

void EntityListView::setStatus(const std::string &text) {
    statusText_ = text;
    // TStaticText has no setter for text after construction in classic TV;
    // redraw via draw() reading statusText_.
    drawView();
}

void EntityListView::draw() {
    TWindow::draw();
    if (status_ != nullptr) {
        TDrawBuffer b;
        const TColorAttr color = mapColor(6);
        b.moveChar(0, ' ', color, static_cast<ushort>(size.x));
        const std::string line = statusText_.empty() ? def_.name : statusText_;
        b.moveStr(2, line, color, static_cast<ushort>(size.x > 4 ? size.x - 4 : 0));
        writeLine(0, 1, size.x, 1, b);
    }
}

void EntityListView::refresh() {
    auto result = client_.list(def_.cliEntity, def_.pageSize, offset_);
    rows_.clear();
    if (!result.ok) {
        setStatus("Error: " + result.errorMessage);
        if (list_ != nullptr) {
            list_->setRows({formatHeader(def_), "(failed to load)"});
        }
        return;
    }

    nlohmann::json arr = result.data;
    if (arr.is_object() && arr.contains("data") && arr["data"].is_array()) {
        arr = arr["data"];
    }
    if (!arr.is_array()) {
        // Some list commands return a bare object; wrap it.
        if (arr.is_object()) {
            arr = nlohmann::json::array({arr});
        } else {
            arr = nlohmann::json::array();
        }
    }

    std::vector<std::string> lines;
    lines.push_back(formatHeader(def_));
    for (const auto &item : arr) {
        rows_.push_back(item);
        lines.push_back(formatRow(def_, item));
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(lines));
    }
    std::ostringstream st;
    st << def_.name << "  offset=" << offset_ << "  rows=" << rows_.size()
       << "  [" << client_.lastCommand() << "]";
    setStatus(st.str());
}

nlohmann::json EntityListView::selectedRow() const {
    if (list_ == nullptr) {
        return {};
    }
    const short focused = list_->focused;
    // Row 0 is header.
    if (focused <= 0) {
        return {};
    }
    const std::size_t idx = static_cast<std::size_t>(focused - 1);
    if (idx >= rows_.size()) {
        return {};
    }
    return rows_[idx];
}

void EntityListView::openDetail() {
    auto row = selectedRow();
    if (row.is_null() || row.empty()) {
        return;
    }
    const int id = jsonToInt(row, def_.idKey);
    if (id > 0) {
        auto r = client_.get(def_.cliEntity, id);
        if (r.ok) {
            row = r.data;
        }
    }
    owner->execView(new EntityDetailView(client_, def_, std::move(row)));
    refresh();
}

void EntityListView::openCreate() {
    if (owner->execView(new EntityForm(client_, def_, true)) == cmOK) {
        refresh();
    }
}

void EntityListView::openEdit() {
    auto row = selectedRow();
    if (row.is_null() || row.empty()) {
        return;
    }
    if (owner->execView(new EntityForm(client_, def_, false, row)) == cmOK) {
        refresh();
    }
}

void EntityListView::doDelete() {
    auto row = selectedRow();
    if (row.is_null() || row.empty() || def_.deleteAction.empty()) {
        return;
    }
    const int id = jsonToInt(row, def_.idKey);
    if (messageBox(("Delete " + def_.name + " #" + std::to_string(id) + "?").c_str(),
                   mfYesNoCancel | mfConfirmation) != cmYes) {
        return;
    }
    auto r = client_.remove(def_.cliEntity, def_.deleteAction, id);
    if (!r.ok) {
        messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
        return;
    }
    refresh();
}

void EntityListView::runRowAction(std::size_t index) {
    if (index >= def_.rowActions.size()) {
        return;
    }
    auto row = selectedRow();
    if (row.is_null() || row.empty()) {
        return;
    }
    if (def_.rowActions[index].run && def_.rowActions[index].run(client_, row)) {
        refresh();
    }
}

void EntityListView::runListAction(std::size_t index) {
    if (index >= def_.listActions.size()) {
        return;
    }
    nlohmann::json empty = nlohmann::json::object();
    if (def_.listActions[index].run && def_.listActions[index].run(client_, empty)) {
        refresh();
    }
}

void EntityListView::handleEvent(TEvent &event) {
    TWindow::handleEvent(event);

    if (event.what == evKeyDown) {
        switch (event.keyDown.keyCode) {
        case kbEnter:
            openDetail();
            clearEvent(event);
            return;
        case kbIns:
            if (def_.canCreate) {
                openCreate();
                clearEvent(event);
            }
            return;
        default:
            if (event.keyDown.charScan.charCode == 'r') {
                refresh();
                clearEvent(event);
            } else if (event.keyDown.charScan.charCode == 'n' && def_.canCreate) {
                openCreate();
                clearEvent(event);
            } else if (event.keyDown.charScan.charCode == 'e' && def_.canEdit) {
                openEdit();
                clearEvent(event);
            } else if (event.keyDown.charScan.charCode == 'd' && !def_.deleteAction.empty()) {
                doDelete();
                clearEvent(event);
            }
            break;
        }
    }

    if (event.what == evCommand) {
        switch (event.message.command) {
        case cmEntityRefresh:
            refresh();
            clearEvent(event);
            break;
        case cmEntityNew:
            openCreate();
            clearEvent(event);
            break;
        case cmEntityEdit:
            openEdit();
            clearEvent(event);
            break;
        case cmEntityDelete:
            doDelete();
            clearEvent(event);
            break;
        case cmEntityPrevPage:
            offset_ = std::max(0, offset_ - def_.pageSize);
            refresh();
            clearEvent(event);
            break;
        case cmEntityNextPage:
            offset_ += def_.pageSize;
            refresh();
            clearEvent(event);
            break;
        default:
            if (event.message.command >= cmEntityAction + 50 &&
                event.message.command < cmEntityAction + 50 + 10) {
                runListAction(static_cast<std::size_t>(event.message.command - (cmEntityAction + 50)));
                clearEvent(event);
            } else if (event.message.command >= cmEntityAction &&
                       event.message.command < cmEntityAction + 50) {
                runRowAction(static_cast<std::size_t>(event.message.command - cmEntityAction));
                clearEvent(event);
            }
            break;
        }
    }

    if (event.what == evBroadcast && event.message.command == cmListItemSelected) {
        openDetail();
        clearEvent(event);
    }
}

TColorAttr EntityListView::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TWindow::mapColor(index);
}

void openEntityList(TProgram *app, CliClient &client, const std::string &cliEntity) {
    const EntityDef *def = findEntity(cliEntity);
    if (def == nullptr) {
        messageBox(("Unknown entity: " + cliEntity).c_str(), mfError | mfOKButton);
        return;
    }
    TWindow *w = new EntityListView(client, *def);
    app->deskTop->insert(w);
}

} // namespace multiflexitui
