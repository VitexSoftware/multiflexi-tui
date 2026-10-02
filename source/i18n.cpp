#include "multiflexitui/i18n.h"

#include <cstdlib>
#include <cstring>
#include <initializer_list>

#include <climits>
#include <sys/stat.h>
#include <unistd.h>

#if defined(__GLIBC__)
extern "C" int _nl_msg_cat_cntr;
#endif

namespace multiflexitui {

namespace {

bool isDirectory(const std::string &path) {
    struct stat st{};
    return ::stat(path.c_str(), &st) == 0 && S_ISDIR(st.st_mode);
}

std::string executableDir() {
    char buf[PATH_MAX];
    const ssize_t len = ::readlink("/proc/self/exe", buf, sizeof(buf) - 1);
    if (len <= 0) {
        return {};
    }
    buf[len] = '\0';
    const std::string exe(buf);
    const std::size_t slash = exe.find_last_of('/');
    return slash == std::string::npos ? std::string() : exe.substr(0, slash);
}

std::string localeDir() {
    if (const char *override_ = std::getenv("MULTIFLEXI_TUI_LOCALEDIR")) {
        return override_;
    }
    const std::string exeDir = executableDir();
    if (!exeDir.empty() && isDirectory(exeDir + "/locale")) {
        return exeDir + "/locale";
    }
#ifndef MULTIFLEXI_TUI_LOCALEDIR
#define MULTIFLEXI_TUI_LOCALEDIR "/usr/share/locale"
#endif
    return MULTIFLEXI_TUI_LOCALEDIR;
}

bool trySetLocale(std::initializer_list<const char *> candidates) {
    for (const char *candidate : candidates) {
        if (setlocale(LC_ALL, candidate) != nullptr) {
            return true;
        }
    }
    return false;
}

} // namespace

void initI18n() {
    setlocale(LC_ALL, "");
    const char *current = setlocale(LC_ALL, nullptr);
    if (current == nullptr || std::strcmp(current, "C") == 0 || std::strcmp(current, "POSIX") == 0) {
        setlocale(LC_ALL, "C.UTF-8");
    }
    bindtextdomain(MULTIFLEXI_TUI_GETTEXT_DOMAIN, localeDir().c_str());
    bind_textdomain_codeset(MULTIFLEXI_TUI_GETTEXT_DOMAIN, "UTF-8");
    textdomain(MULTIFLEXI_TUI_GETTEXT_DOMAIN);
}

void setLanguage(const std::string &lang) {
    if (lang.empty()) {
        unsetenv("LANGUAGE");
        setlocale(LC_ALL, "");
    } else {
        setenv("LANGUAGE", lang.c_str(), 1);
        if (lang == "cs") {
            trySetLocale({"cs_CZ.UTF-8", "cs_CZ.utf8", "cs_CZ"});
        } else if (lang == "en") {
            trySetLocale({"en_US.UTF-8", "en_US.utf8", "en_GB.UTF-8", "en_GB.utf8"});
        }
    }
#if defined(__GLIBC__)
    ++_nl_msg_cat_cntr;
#endif
}

} // namespace multiflexitui
