#include "multiflexitui/CliClient.h"

#include <iostream>

int main() {
    multiflexitui::CliClient client;
    auto r = client.status();
    if (!r.ok) {
        std::cerr << "status failed: " << r.errorMessage << "\n";
        // Still exit 0 if CLI missing in CI without multiflexi — treat spawn failure specially
        return r.exitCode == 127 || r.errorMessage.find("not found") != std::string::npos ? 0 : 1;
    }
    std::cout << "status ok keys=" << (r.data.is_object() ? r.data.size() : 0) << "\n";
    return 0;
}
