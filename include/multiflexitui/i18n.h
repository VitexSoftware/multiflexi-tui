#pragma once

#include <libintl.h>
#include <locale.h>

#include <string>

#ifndef MULTIFLEXI_TUI_GETTEXT_DOMAIN
#define MULTIFLEXI_TUI_GETTEXT_DOMAIN "multiflexi-tui"
#endif

#define _(s) gettext(s)
#define N_(s) (s)

namespace multiflexitui {

void initI18n();
void setLanguage(const std::string &lang);

} // namespace multiflexitui
