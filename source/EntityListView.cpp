#include "multiflexitui/TV.h"
#include "multiflexitui/EntityListView.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/EntityPicker.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/AppShell.h"
#include "multiflexitui/AsyncCliQueue.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

#include <algorithm>
#include <cctype>
#include <sstream>

namespace multiflexitui {

namespace {

std::string toLower(std::string s) {
    std::transform(s.begin(), s.end(), s.begin(), [](unsigned char c) { return static_cast<char>(std::tolower(c)); });
    return s;
}

bool clientSideMatch(const EntityDef &def, const nlohmann::json &row, const std::string &filter) {
    if (filter.empty()) {
        return true;
    }
    const std::string needle = toLower(filter);
    const std::string hay = toLower(formatRow(def, row));
    return hay.find(needle) != std::string::npos;
}

} // namespace

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

    insert(new TLabel(TRect(2, 2, 10, 3), _("Filter"), nullptr));
    filterInput_ = new TInputLine(TRect(11, 2, 50, 3), 80);
    growWide(filterInput_);
    insert(filterInput_);
    TView *apply = new AppButton(TRect(52, 2, 64, 4), _("~A~pply"), cmEntityApplyFilter, bfNormal);
    stickRight(apply);
    insert(apply);

    TScrollBar *vBar = standardScrollBar(sbVertical | sbHandleKeyboard);
    TScrollBar *hBar = standardScrollBar(sbHorizontal | sbHandleKeyboard);
    list_ = new SimpleListViewer(TRect(1, 4, 76, 18), vBar, hBar);
    growFill(list_);
    insert(list_);

    int x = 2;
    auto addBtn = [&](const char *title, ushort cmd, int width = 12) {
        TView *b = new AppButton(TRect(x, 19, x + width, 21), title, cmd, bfNormal);
        stickBottom(b);
        insert(b);
        x += width + 1;
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
    addBtn(_("Act~i~ons"), cmEntityActionsMenu, 12);

    selectNext(False);
    refresh();
}

EntityListView::~EntityListView() = default;

void EntityListView::setStatus(const std::string &text) {
    statusText_ = text;
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

void EntityListView::applyFilterFromInput() {
    if (filterInput_ == nullptr) {
        return;
    }
    char buf[128] = {};
    filterInput_->getData(buf);
    filter_ = buf;
    offset_ = 0;
    refresh();
}

void EntityListView::applyListResult(const CliClient::Result &result) {
    loading_ = false;
    rows_.clear();
    if (!result.ok) {
        setStatus("Error: " + result.errorMessage);
        if (list_ != nullptr) {
            list_->setRows({formatHeader(def_), "(failed to load)"});
        }
        return;
    }

    nlohmann::json arr = asJsonArray(result.data);
    std::vector<std::string> lines;
    lines.push_back(formatHeader(def_));

    const bool useClientFilter =
        !filter_.empty() && def_.cliEntity != "job" && def_.cliEntity != "task" && def_.cliEntity != "artifact";

    for (const auto &item : arr) {
        if (useClientFilter && !clientSideMatch(def_, item, filter_)) {
            continue;
        }
        rows_.push_back(item);
        lines.push_back(formatRow(def_, item));
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(lines));
    }
    std::ostringstream st;
    st << def_.name << "  offset=" << offset_ << "  rows=" << rows_.size();
    if (!filter_.empty()) {
        st << "  filter=" << filter_;
    }
    st << "  [" << result.lastCommand << "]";
    setStatus(st.str());
}

void EntityListView::refresh() {
    if (loading_) {
        return;
    }
    loading_ = true;
    setStatus(_("Busy…"));

    ListOptions opts;
    opts.limit = def_.pageSize;
    opts.offset = offset_;
    if (def_.supportsCompanyScope) {
        if (auto *app = dynamic_cast<MultiFlexiApp *>(TProgram::application)) {
            if (app->activeCompanyId() > 0) {
                opts.companyId = app->activeCompanyId();
            }
        }
    }
    if (!filter_.empty() &&
        (def_.cliEntity == "job" || def_.cliEntity == "task" || def_.cliEntity == "artifact")) {
        opts.filter = filter_;
    }

    std::vector<std::string> args{def_.cliEntity + ":list", "--order=D",
                                  "--limit=" + std::to_string(opts.limit),
                                  "--offset=" + std::to_string(opts.offset)};
    if (opts.companyId > 0) {
        args.push_back("--company_id=" + std::to_string(opts.companyId));
    }
    if (!opts.filter.empty()) {
        if (def_.cliEntity == "job") {
            args.push_back("--status=" + opts.filter);
        } else if (def_.cliEntity == "task") {
            args.push_back("--state=" + opts.filter);
        } else if (def_.cliEntity == "artifact") {
            args.push_back("--job_id=" + opts.filter);
        }
    }

    // Capture weak-ish: if window closed before callback, skip update.
    auto *self = this;
    const bool enqueued = enqueueCliJson(client_, args, [self](CliClient::Result r) {
        if (auto *app = dynamic_cast<MultiFlexiApp *>(TProgram::application)) {
            app->setBusy(app->cliQueue().busy());
        }
        // Ensure the view is still on the desktop.
        bool alive = false;
        if (TProgram::deskTop != nullptr) {
            TView *p = TProgram::deskTop->first();
            while (p != nullptr) {
                if (p == self) {
                    alive = true;
                    break;
                }
                p = p->nextView();
            }
        }
        if (alive) {
            self->applyListResult(r);
        }
    });
    if (!enqueued) {
        // sync fallback already invoked callback
    }
}

nlohmann::json EntityListView::selectedRow() const {
    if (list_ == nullptr) {
        return {};
    }
    const short focused = list_->focused;
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
    if (def_.supportsGet) {
        const int id = jsonToInt(row, def_.idKey);
        if (id > 0) {
            auto r = client_.get(def_.cliEntity, id);
            if (r.ok) {
                row = r.data;
            }
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

void EntityListView::openActionsMenu() {
    std::vector<std::string> labels;
    std::vector<std::pair<bool, std::size_t>> map; // isList, index
    for (std::size_t i = 0; i < def_.listActions.size(); ++i) {
        labels.push_back(std::string("[List] ") + def_.listActions[i].label);
        map.push_back({true, i});
    }
    for (std::size_t i = 0; i < def_.rowActions.size(); ++i) {
        labels.push_back(def_.rowActions[i].label);
        map.push_back({false, i});
    }
    if (labels.empty()) {
        messageBox(_("No actions available"), mfInformation | mfOKButton);
        return;
    }
    const int idx = pickActionIndex(labels, _("Actions"));
    if (idx < 0 || static_cast<std::size_t>(idx) >= map.size()) {
        return;
    }
    if (map[static_cast<std::size_t>(idx)].first) {
        runListAction(map[static_cast<std::size_t>(idx)].second);
    } else {
        runRowAction(map[static_cast<std::size_t>(idx)].second);
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
        case kbF2:
            openActionsMenu();
            clearEvent(event);
            return;
        default:
            if (event.keyDown.charScan.charCode == '/') {
                if (filterInput_ != nullptr) {
                    filterInput_->select();
                }
                clearEvent(event);
            } else if (event.keyDown.charScan.charCode == 'r') {
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
            } else if (event.keyDown.charScan.charCode == 'i') {
                openActionsMenu();
                clearEvent(event);
            } else {
                // Row-action hotkeys
                const char ch = event.keyDown.charScan.charCode;
                if (ch != 0) {
                    for (std::size_t i = 0; i < def_.rowActions.size(); ++i) {
                        if (def_.rowActions[i].hotkey != 0 && def_.rowActions[i].hotkey == ch) {
                            runRowAction(i);
                            clearEvent(event);
                            return;
                        }
                    }
                    for (std::size_t i = 0; i < def_.listActions.size(); ++i) {
                        if (def_.listActions[i].hotkey != 0 && def_.listActions[i].hotkey == ch) {
                            runListAction(i);
                            clearEvent(event);
                            return;
                        }
                    }
                }
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
        case cmEntityActionsMenu:
            openActionsMenu();
            clearEvent(event);
            break;
        case cmEntityApplyFilter:
            applyFilterFromInput();
            clearEvent(event);
            break;
        default:
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
