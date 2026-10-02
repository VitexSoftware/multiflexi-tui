#pragma once

#include "multiflexitui/TV.h"

#include <string>

namespace multiflexitui {

// TStatusLine plus the URL of the MultiFlexi request currently in flight.
// Hotkey items stay on the left; the URL is right-aligned and, when the
// terminal is narrow, ellipsized from the left so the evidence/id remain
// visible. TStatusLine::drawSelect() is private, so draw() repeats the
// hotkey painting from tstatusl.cpp and then writes the URL into the gap.
class AppStatusLine : public TStatusLine {
public:
    AppStatusLine(const TRect &bounds, TStatusDef &aDefs) noexcept;

    void draw() override;
    void handleEvent(TEvent &event) override;
    void setCurrentUrl(const std::string &url);
    void setCompany(const std::string &company);
    void setQueryMode(const std::string &mode);

    const std::string &currentUrl() const { return currentUrl_; }

private:
    std::string currentUrl_;
    std::string company_;
    std::string queryMode_;
    int urlStart_ = -1;
    int urlEnd_ = -1;
};

} // namespace multiflexitui
