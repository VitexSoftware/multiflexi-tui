#pragma once

#include "multiflexitui/TV.h"

#include <string>
#include <vector>

namespace multiflexitui {

// A read-only TListViewer backed by a plain vector of already-formatted
// strings, one per row. Used for every tabular/list display in this app
// (companies, evidences, record grids, record detail field lines) instead of
// a separate TListBox/TDataCollection per screen.
class SimpleListViewer : public TListViewer {
public:
    // `hScrollBar` is optional: when given, rows wider than the view become
    // horizontally scrollable (Left/Right arrow keys, or dragging the bar)
    // instead of being silently truncated - needed by RecordListBox once
    // more columns are selected than fit the dialog's width.
    SimpleListViewer(const TRect &bounds, TScrollBar *vScrollBar, TScrollBar *hScrollBar = nullptr) noexcept;

    void setRows(std::vector<std::string> rows);
    const std::string &rowAt(short item) const;
    std::size_t rowCount() const { return rows_.size(); }

    void getText(char *dest, short item, short maxLen) override;
    void draw() override;
    void changeBounds(const TRect &bounds) override;
    void handleEvent(TEvent &event) override;
    TColorAttr mapColor(uchar index) override;

protected:
    std::vector<std::string> rows_;

private:
    // Recomputes hScrollBar's range from the widest current row vs. the
    // view's visible width; called whenever either changes.
    void updateHScrollRange();
};

} // namespace multiflexitui
