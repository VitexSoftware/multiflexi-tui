#pragma once

#include "multiflexitui/CliClient.h"

#include <deque>
#include <functional>
#include <optional>
#include <string>
#include <vector>

namespace multiflexitui {

// One-at-a-time async CLI queue; pump() from MultiFlexiApp::idle().
class AsyncCliQueue {
public:
    using Callback = std::function<void(CliClient::Result)>;

    void enqueue(CliClient &client, std::vector<std::string> args, Callback callback);
    void pump();
    bool busy() const { return current_.has_value() || !pending_.empty(); }
    std::size_t pendingCount() const { return pending_.size() + (current_ ? 1u : 0u); }

private:
    struct Job {
        std::vector<std::string> argv;
        std::string lastCommand;
        Callback callback;
    };

    void startNext();

    std::deque<Job> pending_;
    std::optional<ProcessHandle> current_;
    Job currentJob_;
};

// Helper: enqueue via the running MultiFlexiApp (no-op / sync fallback if unavailable).
bool enqueueCliJson(CliClient &client, std::vector<std::string> args,
                    AsyncCliQueue::Callback callback);

} // namespace multiflexitui
