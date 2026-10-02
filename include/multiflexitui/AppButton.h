#pragma once

#include "multiflexitui/TV.h"

namespace multiflexitui {

// A drop-in TButton replacement with fixed truecolor RGB instead of the
// classic 16-color BIOS palette entries. TButton::draw() resolves its
// colors through mapColor(1..8) (normal/default/selected/disabled/shortcut
// x3/shadow, see tvision's cpButton), which walks up the owning dialog's
// and application's palette chain. Removing this override was tried
// (letting stock TButton colors resolve through tvision's own canonical
// RGB tables) on the theory that a truecolor-capable terminal
// (COLORTERM=truecolor/24bit) would no longer need it - but tvision's own
// stock "Button normal" color for this app's blue window palette IS a
// bright red regardless of truecolor, reproducing the exact red-button
// problem this class exists to avoid. So the override stays; see
// WindowColors.h for the matching dialog-background color `faceBg` uses.
class AppButton : public TButton {
public:
    using TButton::TButton;

    TColorAttr mapColor(uchar index) override;
};

} // namespace multiflexitui
