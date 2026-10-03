#include "multiflexitui/TV.h"
#include "multiflexitui/AppStatusLine.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/i18n.h"

namespace multiflexitui {

AppStatusLine::AppStatusLine(const TRect &bounds, TStatusDef &aDefs) noexcept : TStatusLine(bounds, aDefs) {
}

void AppStatusLine::setCurrentUrl(const std::string &url) {
    if (url == currentUrl_) {
        return;
    }

    currentUrl_ = url;
    drawView();
}

void AppStatusLine::setCompany(const std::string &company) {
    if (company == company_) {
        return;
    }

    company_ = company;
    drawView();
}

void AppStatusLine::setQueryMode(const std::string &mode) {
    if (mode == queryMode_) {
        return;
    }

    queryMode_ = mode;
    drawView();
}

void AppStatusLine::draw() {
    TDrawBuffer b;
    const TAttrPair cNormal = getColor(0x0301);
    const TAttrPair cNormDisabled = getColor(0x0202);
    // Accent for the active-company badge so it reads as context, not a hotkey.
    const TColorAttr cCompany(TColor(TColorRGB(0xF2C23E)), TColor(TColorRGB(kWindowBg)));
    b.moveChar(0, ' ', cNormal, static_cast<ushort>(size.x));
    urlStart_ = -1;
    urlEnd_ = -1;

    int used = 0;

    for (TStatusItem *item = items; item != nullptr; item = item->next) {
        if (item->text == nullptr) {
            continue;
        }

        const int length = cstrlen(item->text);

        if (used + length < size.x) {
            const TAttrPair color = commandEnabled(item->command) ? cNormal : cNormDisabled;
            b.moveChar(static_cast<ushort>(used), ' ', color, 1);
            b.moveCStr(static_cast<ushort>(used + 1), item->text, color);
            b.moveChar(static_cast<ushort>(used + length + 1), ' ', color, 1);
        }

        used += length + 2;
    }

    // Right-side context: company (always when set) + optional Busy badge.
    std::string right;
    if (!company_.empty()) {
        right = std::string(_("Co:")) + " " + company_;
    }
    if (!queryMode_.empty()) {
        if (!right.empty()) {
            right += "  ";
        }
        right += queryMode_;
    }

    int rightStart = size.x;
    if (!right.empty() && used + 2 < size.x) {
        int avail = size.x - used - 1;
        if (static_cast<int>(right.size()) > avail) {
            // Keep the Co: prefix; ellipsize the company name from the left.
            if (avail > 4) {
                right = "..." + right.substr(right.size() - static_cast<std::size_t>(avail - 3));
            } else {
                right.resize(static_cast<std::size_t>(avail));
            }
        }
        rightStart = size.x - static_cast<int>(right.size());
        if (rightStart < used + 1) {
            rightStart = used + 1;
        }
        const TColorAttr color = company_.empty() ? TColorAttr(cNormal) : cCompany;
        b.moveStr(static_cast<ushort>(rightStart), right, color);
    }

    // CLI command fills the gap between hotkeys and the right badge.
    if (!currentUrl_.empty() && used + 1 < rightStart) {
        const int avail = rightStart - used - 2;
        if (avail > 0) {
            std::string shown = currentUrl_;
            if (static_cast<int>(shown.size()) > avail) {
                if (avail <= 3) {
                    shown = shown.substr(shown.size() - static_cast<std::size_t>(avail));
                } else {
                    shown = "..." + shown.substr(shown.size() - static_cast<std::size_t>(avail - 3));
                }
            }

            int start = rightStart - 1 - static_cast<int>(shown.size());
            if (start < used + 1) {
                start = used + 1;
            }

            TColorAttr link(cNormal);
            link.setStyle(static_cast<ushort>(link.getStyle() | slUnderline));
            b.moveStr(static_cast<ushort>(start), shown, link);
            urlStart_ = start;
            urlEnd_ = start + static_cast<int>(shown.size());
        }
    }

    writeLine(0, 0, static_cast<ushort>(size.x), 1, b);
}

void AppStatusLine::handleEvent(TEvent &event) {
    if (event.what == evMouseDown && urlStart_ >= 0) {
        const TPoint mouse = makeLocal(event.mouse.where);

        if (mouse.y == 0 && mouse.x >= urlStart_ && mouse.x < urlEnd_) {
            message(owner, evCommand, cmShowWebQr, nullptr);
            clearEvent(event);
            return;
        }
    }

    TStatusLine::handleEvent(event);
}

} // namespace multiflexitui
