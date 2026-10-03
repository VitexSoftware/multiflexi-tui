#include "multiflexitui/CliClient.h"

#include <iostream>

int main() {
    multiflexitui::CliClient client;
    auto r = client.status();
    if (!r.ok) {
        std::cerr << "status failed: " << r.errorMessage << "\n";
        return r.exitCode == 127 || r.errorMessage.find("not found") != std::string::npos ? 0 : 1;
    }
    std::cout << "status ok keys=" << (r.data.is_object() ? r.data.size() : 0) << "\n";

    // Non-blocking ProcessHandle drain smoke test (echo).
    auto handle = multiflexitui::ProcessHandle::start({"true"});
    int spins = 0;
    while (!handle.pump() && spins < 1000) {
        ++spins;
    }
    if (!handle.done() || handle.result().exitCode != 0) {
        std::cerr << "ProcessHandle pump failed exit=" << handle.result().exitCode << "\n";
        return 1;
    }
    std::cout << "async process handle ok spins=" << spins << "\n";
    return 0;
}
