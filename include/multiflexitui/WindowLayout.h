#pragma once

#include "multiflexitui/TV.h"

namespace multiflexitui {

// TDialog clears the zoom and grow flags. Put them back so the frame shows
// the maximize control and the window can fill the desktop.
inline void makeMaximizable(TWindow &window) {
    window.flags |= wfGrow | wfZoom;
}

inline void growWide(TView *view) {
    if (view != nullptr) {
        view->growMode |= gfGrowHiX;
    }
}

inline void growFill(TView *view) {
    if (view != nullptr) {
        view->growMode |= gfGrowHiX | gfGrowHiY;
    }
}

// Keep the view against the bottom edge when the window grows.
inline void stickBottom(TView *view) {
    if (view != nullptr) {
        view->growMode |= gfGrowLoY | gfGrowHiY;
    }
}

inline void stickBottomWide(TView *view) {
    if (view != nullptr) {
        view->growMode |= gfGrowHiX | gfGrowLoY | gfGrowHiY;
    }
}

inline void stickRight(TView *view) {
    if (view != nullptr) {
        view->growMode |= gfGrowLoX | gfGrowHiX;
    }
}

inline void stickCorner(TView *view) {
    if (view != nullptr) {
        view->growMode |= gfGrowLoX | gfGrowHiX | gfGrowLoY | gfGrowHiY;
    }
}

} // namespace multiflexitui
