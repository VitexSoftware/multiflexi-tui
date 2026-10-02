#include "multiflexitui/AppButton.h"
#include "multiflexitui/WindowColors.h"

namespace multiflexitui {

namespace {

// Matches cpButton's 8 local indices: 1=normal, 2=default, 3=selected
// (pressed/focused), 4=disabled, 5-7=shortcut letter in each of those three
// states, 8=shadow. Index 8 is deliberately NOT handled here: TButton::draw()
// paints it as a flat fill (see tbutton.cpp, cShadow = getColor(8)), which
// falls through to TButton::mapColor(8) and from there to the owning
// window's own mapColor() (see WindowColors.cpp, index 15) - a fixed, solid
// black foreground on `kWindowBg`, so it blends into the dialog everywhere
// except the actual drop-shadow glyphs (see that file's comment on index 15
// for why the two need different backgrounds).
//
// `faceBg` is `kWindowBg` - the same background WindowColors.cpp uses for
// the dialog itself (StaticText/InputLine/frame interior): a distinct
// accent tone here was tried, but then the thin dialog-background strip
// between adjacent buttons read as a *third*, mismatched color next to the
// button face and the shadow, instead of the classic TV look where the
// button face is flush with its dialog and only the shadow sets it apart.
TColorAttr colorFor(uchar index) {
    constexpr uint32_t faceBg = kWindowBg;

    switch (index) {
    case 1: // normal
    case 2: // default (primary action)
        return TColorAttr(TColor(TColorRGB(0xFFFFFF)), TColor(TColorRGB(faceBg)));
    case 3: // selected/pressed
        return TColorAttr(TColor(TColorRGB(0x000000)), TColor(TColorRGB(0xFFFFFF)));
    case 4: // disabled
        return TColorAttr(TColor(TColorRGB(0x5A6B80)), TColor(TColorRGB(faceBg)));
    case 5: // shortcut letter, normal
    case 6: // shortcut letter, default
        return TColorAttr(TColor(TColorRGB(0xF2C23E)), TColor(TColorRGB(faceBg)));
    case 7: // shortcut letter, selected
        return TColorAttr(TColor(TColorRGB(0xB34700)), TColor(TColorRGB(0xFFFFFF)));
    default:
        return TColorAttr(TColor(TColorRGB(0xFFFFFF)), TColor(TColorRGB(faceBg)));
    }
}

} // namespace

TColorAttr AppButton::mapColor(uchar index) {
    if (index >= 1 && index <= 7) {
        return colorFor(index);
    }

    return TButton::mapColor(index);
}

} // namespace multiflexitui
