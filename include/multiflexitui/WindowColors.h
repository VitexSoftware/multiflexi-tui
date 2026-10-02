#pragma once

#include "multiflexitui/TV.h"

namespace multiflexitui {

// Shared with AppButton.cpp: a button's face must match this exact
// background so it reads as part of the dialog it sits on (not a
// disconnected block), while its shadow (see windowColor()'s index 15,
// below) stays unambiguously darker than both, for the classic TV raised-
// button 3D look.
constexpr uint32_t kWindowBg = 0x14233A;

// Shared with SimpleListViewer.cpp: the "this is content, not decoration"
// tone for input lines and list views - a shade lighter than `kWindowBg` so
// an empty, unfocused input line or an idle list still visibly reads as
// "click here" / "this is a data area", not just plain dialog background
// (reported as a real usability loss when it briefly matched `kWindowBg`).
constexpr uint32_t kFieldBg = 0x1B3A5C;

// Explicit-RGB replacement for the handful of TWindow/TDialog palette
// indices this app's widgets resolve through by default (see the "Palette
// layout" comment next to TWindow::getPalette/TDialog::getPalette in
// tvision's dialogs.h - both share the same index layout): frame, static
// text, scrollbar, input lines, and the button shadow that AppButton
// deliberately leaves unhandled (see AppButton.h). Same "don't trust what a
// specific palette index happens to render as under an arbitrary terminal
// ANSI-16 theme" reasoning as AppButton/StatusView/SimpleListViewer
// (reported: a theme that recolors most of the app red/orange left plain
// window backgrounds, input lines and the button shadow all showing as
// light red).
//
// Call this first from every TWindow/TDialog subclass's mapColor()
// override; it returns false for indices it doesn't handle (labels,
// clusters, buttons - AppButton already colors those directly), so the
// caller should fall back to its base class implementation in that case.
bool windowColor(uchar index, TColorAttr &out);

} // namespace multiflexitui
