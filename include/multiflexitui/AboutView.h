#pragma once

#include "multiflexitui/TV.h"
#include "multiflexitui/SimpleListViewer.h"

#include <string>
#include <vector>

namespace multiflexitui {

// Version and homepages of the libraries and utilities this program uses.
class AboutView : public TDialog {
public:
    AboutView();

    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

private:
    void openSelected();

    SimpleListViewer *links_;
    std::vector<std::string> urls_;
};

} // namespace multiflexitui
