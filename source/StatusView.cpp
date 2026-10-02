#include "multiflexitui/TV.h"
#include "multiflexitui/StatusView.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/i18n.h"

namespace multiflexitui {

StatusView::StatusView(CliClient &client)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(8, 2, 72, 22), _("MultiFlexi Status")),
      client_(client) {
    options |= ofCentered;
    makeMaximizable(*this);

    TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
    list_ = new SimpleListViewer(TRect(2, 2, 62, 16), bar);
    growFill(list_);
    insert(list_);

    TView *refresh = new AppButton(TRect(2, 17, 16, 19), _("~R~efresh"), cmEntityRefresh, bfNormal);
    stickBottom(refresh);
    insert(refresh);
    TView *close = new AppButton(TRect(48, 17, 62, 19), _("Close"), cmCancel, bfDefault);
    stickCorner(close);
    insert(close);

    reload();
    selectNext(False);
}

void StatusView::reload() {
    auto r = client_.status();
    std::vector<std::string> lines;
    if (!r.ok) {
        lines.push_back("Failed to load status:");
        lines.push_back(r.errorMessage);
    } else if (r.data.is_object()) {
        for (auto it = r.data.begin(); it != r.data.end(); ++it) {
            std::string key = it.key();
            while (key.size() < 20) {
                key.push_back(' ');
            }
            lines.push_back(key + jsonToString(it.value()));
        }
    } else {
        lines.push_back(r.data.dump(2));
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(lines));
    }
}

void StatusView::handleEvent(TEvent &event) {
    TDialog::handleEvent(event);
    if (event.what == evCommand && event.message.command == cmEntityRefresh) {
        reload();
        clearEvent(event);
    }
}

TColorAttr StatusView::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

} // namespace multiflexitui
