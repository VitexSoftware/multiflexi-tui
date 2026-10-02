#include "multiflexitui/WindowColors.h"

namespace multiflexitui {

bool windowColor(uchar index, TColorAttr &out) {
    static const uint32_t fg = 0xF0F0F0;
    static const uint32_t mutedFg = 0x8FA0B8;
    static const uint32_t bg = kWindowBg;      // neutral dark blue, matches SimpleListViewer/AppButton's face.
    static const uint32_t fieldBg = kFieldBg;  // lighter "this is editable" tone - inputs only.
    static const uint32_t accentBg = 0x264F78; // focused input / selected scrollbar control.
    static const uint32_t shadowFg = 0x000000; // the actual drop-shadow ink.

    switch (index) {
    case 1: // Frame passive
        out = TColorAttr(TColor(TColorRGB(mutedFg)), TColor(TColorRGB(bg)));
        return true;
    case 2: // Frame active
    case 3: // Frame icon
        // Same `bg` as plain StaticText, not a separate lighter tone: a
        // distinct frame background left stray cells showing that tone
        // wherever a redraw exposed a not-yet-repainted gap (reported: a
        // patch of "wrong" color in the 1-cell gaps between toolbar
        // buttons after an overlapping dialog closed - now also fixed at
        // the source with an explicit backdrop, see RecordListView.cpp).
        // The active/passive frame is still readable from its box-drawing
        // glyphs and brighter foreground alone. `fieldBg` stays reserved
        // for input lines, below, where a lighter tone is the point: it is
        // the "click here to type" affordance, not decoration.
        out = TColorAttr(TColor(TColorRGB(fg)), TColor(TColorRGB(bg)));
        return true;
    case 4: // ScrollBar page area
        out = TColorAttr(TColor(TColorRGB(mutedFg)), TColor(TColorRGB(bg)));
        return true;
    case 5: // ScrollBar controls
        out = TColorAttr(TColor(TColorRGB(fg)), TColor(TColorRGB(accentBg)));
        return true;
    case 6: // StaticText
        out = TColorAttr(TColor(TColorRGB(fg)), TColor(TColorRGB(bg)));
        return true;
    case 15: // Button shadow (AppButton falls through to here on purpose).
        // TButton::drawState() (tbutton.cpp) paints this single color pair
        // in two very different roles: (a) the actual drop-shadow - solid
        // CP437 half/full-block glyphs (0xDC/0xDB/0xDF) along the button's
        // right edge and, offset one cell right, its bottom row, where the
        // glyph's ink makes the *foreground* what the eye reads; and
        // (b) a blank filler cell forced onto the button's own left edge on
        // every row (and the first 2 cells of the bottom row), meant to be
        // invisible - a space glyph only ever shows its *background*, so
        // for that filler to disappear into the dialog rather than read as
        // a matching shadow bar down the left side too (reported: a
        // "two-sided" shadow), the background here must be the dialog's own
        // background, not the shadow ink color.
        out = TColorAttr(TColor(TColorRGB(shadowFg)), TColor(TColorRGB(bg)));
        return true;
    case 19: // InputLine normal text
        // Deliberately lighter than the surrounding dialog (`fieldBg`, not
        // `bg`): an input line that's the same color as plain label text
        // gives no visual cue that it's a "click here and type" field at
        // all, even when empty and unfocused - reported as a real
        // usability loss, not just cosmetics, after this briefly matched
        // `bg` like every other background in the dialog.
        out = TColorAttr(TColor(TColorRGB(fg)), TColor(TColorRGB(fieldBg)));
        return true;
    case 20: // InputLine selected text (focused)
        out = TColorAttr(TColor(TColorRGB(fg)), TColor(TColorRGB(accentBg)));
        return true;
    default:
        return false;
    }
}

} // namespace multiflexitui
