#include "multiflexitui/TV.h"
#include "multiflexitui/AppShell.h"
#include "multiflexitui/i18n.h"

#include <cstdlib>
#include <cstring>
#include <iostream>
#include <string>

namespace {

std::string envOr(const char *name, const std::string &fallback) {
    const char *v = std::getenv(name);
    return (v != nullptr) ? std::string(v) : fallback;
}

bool takeFlag(const std::string &arg, const char *prefix, std::string &out) {
    const std::size_t len = std::strlen(prefix);
    if (arg.compare(0, len, prefix) == 0) {
        out = arg.substr(len);
        return true;
    }
    return false;
}

} // namespace

int main(int argc, char **argv) {
    multiflexitui::initI18n();

    std::string cliBinary = envOr("MULTIFLEXI_TUI_CLI", "multiflexi-cli");
    std::string envFile;

    for (int i = 1; i < argc; ++i) {
        std::string arg = argv[i];
        std::string value;
        if (arg == "-h" || arg == "--help") {
            std::cout << "Usage: multiflexi-tui [--cli=PATH] [--envfile=PATH]\n";
            return 0;
        }
        if (takeFlag(arg, "--cli=", value)) {
            cliBinary = value;
        } else if (takeFlag(arg, "--envfile=", value)) {
            envFile = value;
        }
    }

    multiflexitui::MultiFlexiApp app;
    app.configure(cliBinary, envFile);
    app.run();
    return 0;
}
