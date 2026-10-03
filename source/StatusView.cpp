#include "multiflexitui/TV.h"
#include "multiflexitui/StatusView.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/i18n.h"

#include <algorithm>
#include <cctype>
#include <cstring>
#include <set>

namespace multiflexitui {

namespace {

constexpr uint32_t kServiceRunning = 0x3DDC84; // green
constexpr uint32_t kServiceStopped = 0xE89B3C; // orange
constexpr uint32_t kServiceFailed = 0xF44747;  // red

const std::set<std::string> &serviceKeys() {
    static const std::set<std::string> keys = {"executor", "scheduler", "housekeeper"};
    return keys;
}

std::string toLower(std::string s) {
    std::transform(s.begin(), s.end(), s.begin(),
                   [](unsigned char c) { return static_cast<char>(std::tolower(c)); });
    return s;
}

// Map systemd-ish status text to a foreground color (0 = default).
uint32_t colorForServiceState(const std::string &raw) {
    const std::string s = toLower(raw);
    if (s.empty()) {
        return 0;
    }
    if (s == "failed" || s.find("failed") != std::string::npos) {
        return kServiceFailed;
    }
    if (s == "active" || s == "running" || s == "activating") {
        return kServiceRunning;
    }
    // inactive, dead, stopped, unknown, …
    return kServiceStopped;
}

} // namespace

void StatusListViewer::setRows(std::vector<std::string> rows, std::vector<uint32_t> rowFg) {
    rowFg_ = std::move(rowFg);
    SimpleListViewer::setRows(std::move(rows));
}

void StatusListViewer::draw() {
    // Custom draw so each service row can keep its state color; focused row
    // uses the same tint on a selected background instead of reverse video.
    const ushort hOffset = hScrollBar != nullptr ? static_cast<ushort>(hScrollBar->value) : 0;
    const TColorAttr normalBg = TColorAttr(TColor(TColorRGB(0xF0F0F0)), TColor(TColorRGB(kFieldBg)));
    const TColorAttr selectedBg = TColorAttr(TColor(TColorRGB(0xF0F0F0)), TColor(TColorRGB(0x264F78)));

    for (short y = 0; y < size.y; ++y) {
        const short item = static_cast<short>(topItem + y);
        TDrawBuffer b;
        TColorAttr color = normalBg;

        if (item < range) {
            uint32_t fg = 0;
            if (item >= 0 && static_cast<std::size_t>(item) < rowFg_.size()) {
                fg = rowFg_[static_cast<std::size_t>(item)];
            }
            if (item == focused) {
                color = selectedBg;
                if (fg != 0) {
                    color = TColorAttr(TColor(TColorRGB(fg)), TColor(TColorRGB(0x264F78)));
                }
            } else if (fg != 0) {
                color = TColorAttr(TColor(TColorRGB(fg)), TColor(TColorRGB(kFieldBg)));
            }

            char text[256];
            getText(text, item, 255);
            text[255] = '\0';
            b.moveChar(0, ' ', color, static_cast<ushort>(size.x));
            // Focus marker column 0, text from column 1 (matches TListViewer).
            if (item == focused) {
                b.moveChar(0, 0x10 /* ► */, color, 1);
            }
            b.moveStr(1, text, color, static_cast<ushort>(size.x > 1 ? size.x - 1 : 0), hOffset);
        } else {
            b.moveChar(0, ' ', normalBg, static_cast<ushort>(size.x));
        }
        writeLine(0, y, size.x, 1, b);
    }
}

StatusView::StatusView(CliClient &client)
    : TWindowInit(&TDialog::initFrame),
      TDialog(TRect(6, 2, 74, 22), _("MultiFlexi Status")),
      client_(client) {
    options |= ofCentered;
    makeMaximizable(*this);

    TScrollBar *bar = standardScrollBar(sbVertical | sbHandleKeyboard);
    list_ = new StatusListViewer(TRect(2, 2, 66, 16), bar);
    growFill(list_);
    insert(list_);

    TView *refresh = new AppButton(TRect(2, 17, 14, 19), _("~R~efresh"), cmEntityRefresh, bfNormal);
    stickBottom(refresh);
    insert(refresh);

    startBtn_ = new AppButton(TRect(16, 17, 28, 19), _("~S~tart"), cmServiceStart, bfNormal);
    stickBottom(startBtn_);
    insert(startBtn_);

    stopBtn_ = new AppButton(TRect(30, 17, 42, 19), _("Sto~p~"), cmServiceStop, bfNormal);
    stickBottom(stopBtn_);
    insert(stopBtn_);

    TView *close = new AppButton(TRect(50, 17, 64, 19), _("Close"), cmCancel, bfDefault);
    stickCorner(close);
    insert(close);

    reload();
    selectNext(False);
}

void StatusView::reload() {
    auto r = client_.status();
    std::vector<std::string> lines;
    std::vector<uint32_t> colors;
    keys_.clear();
    if (!r.ok) {
        lines.push_back("Failed to load status:");
        lines.push_back(r.errorMessage);
        colors.push_back(kServiceFailed);
        colors.push_back(kServiceFailed);
        keys_.push_back("");
        keys_.push_back("");
    } else if (r.data.is_object()) {
        for (auto it = r.data.begin(); it != r.data.end(); ++it) {
            std::string key = it.key();
            keys_.push_back(key);
            std::string label = key;
            while (label.size() < 20) {
                label.push_back(' ');
            }
            const std::string value = jsonToString(it.value());
            lines.push_back(label + value);
            uint32_t fg = 0;
            if (serviceKeys().count(key) > 0) {
                fg = colorForServiceState(value);
            }
            colors.push_back(fg);
        }
    } else {
        lines.push_back(r.data.dump(2));
        colors.push_back(0);
        keys_.push_back("");
    }
    if (list_ != nullptr) {
        list_->setRows(std::move(lines), std::move(colors));
        if (!keys_.empty()) {
            list_->focusItem(0);
        }
    }
    updateServiceButtons();
}

std::string StatusView::focusedKey() const {
    if (list_ == nullptr || keys_.empty()) {
        return {};
    }
    const short f = list_->focused;
    if (f < 0 || static_cast<std::size_t>(f) >= keys_.size()) {
        return {};
    }
    return keys_[static_cast<std::size_t>(f)];
}

bool StatusView::focusedIsService() const {
    return serviceKeys().count(focusedKey()) > 0;
}

void StatusView::updateServiceButtons() {
    const Boolean disable = focusedIsService() ? False : True;
    if (startBtn_ != nullptr) {
        startBtn_->setState(sfDisabled, disable);
        startBtn_->drawView();
    }
    if (stopBtn_ != nullptr) {
        stopBtn_->setState(sfDisabled, disable);
        stopBtn_->drawView();
    }
}

void StatusView::controlService(const std::string &action) {
    const std::string key = focusedKey();
    if (serviceKeys().count(key) == 0) {
        return;
    }
    auto r = client_.runJson({"service:" + action, "--name=" + key});
    if (!r.ok) {
        messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
    } else {
        std::string msg = action + " " + key;
        if (r.data.is_object() && r.data.contains("unit_status")) {
            msg += " → " + jsonToString(r.data["unit_status"]);
        }
        messageBox(msg.c_str(), mfInformation | mfOKButton);
    }
    reload();
}

void StatusView::handleEvent(TEvent &event) {
    TDialog::handleEvent(event);
    if (event.what == evCommand) {
        switch (event.message.command) {
        case cmEntityRefresh:
            reload();
            clearEvent(event);
            break;
        case cmServiceStart:
            if (focusedIsService()) {
                controlService("start");
            }
            clearEvent(event);
            break;
        case cmServiceStop:
            if (focusedIsService()) {
                controlService("stop");
            }
            clearEvent(event);
            break;
        default:
            break;
        }
    }
    updateServiceButtons();
}

TColorAttr StatusView::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TDialog::mapColor(index);
}

} // namespace multiflexitui
