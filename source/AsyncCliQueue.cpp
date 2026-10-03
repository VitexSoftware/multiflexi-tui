#include "multiflexitui/AsyncCliQueue.h"
#include "multiflexitui/AppShell.h"
#include "multiflexitui/TV.h"

#include <utility>

namespace multiflexitui {

void AsyncCliQueue::enqueue(CliClient &client, std::vector<std::string> args, Callback callback) {
    Job job;
    job.argv = client.buildArgv(args);
    job.lastCommand.clear();
    for (std::size_t i = 0; i < job.argv.size(); ++i) {
        if (i != 0) {
            job.lastCommand.push_back(' ');
        }
        job.lastCommand += job.argv[i];
    }
    job.callback = std::move(callback);
    pending_.push_back(std::move(job));
    if (!current_) {
        startNext();
    }
}

void AsyncCliQueue::startNext() {
    if (current_ || pending_.empty()) {
        return;
    }
    currentJob_ = std::move(pending_.front());
    pending_.pop_front();
    current_ = ProcessHandle::start(currentJob_.argv);
    if (current_ && current_->done()) {
        const std::string binary = currentJob_.argv.empty() ? std::string() : currentJob_.argv.front();
        CliClient::Result r = CliClient::parseProcessResult(current_->result(), binary, currentJob_.argv);
        r.lastCommand = currentJob_.lastCommand;
        auto cb = std::move(currentJob_.callback);
        current_.reset();
        currentJob_ = Job{};
        if (cb) {
            cb(std::move(r));
        }
        startNext();
    }
}

void AsyncCliQueue::pump() {
    if (!current_) {
        startNext();
        return;
    }
    if (!current_->pump()) {
        return;
    }
    const std::string binary = currentJob_.argv.empty() ? std::string() : currentJob_.argv.front();
    CliClient::Result r = CliClient::parseProcessResult(current_->result(), binary, currentJob_.argv);
    r.lastCommand = currentJob_.lastCommand;
    auto cb = std::move(currentJob_.callback);
    current_.reset();
    currentJob_ = Job{};
    if (cb) {
        cb(std::move(r));
    }
    startNext();
}

bool enqueueCliJson(CliClient &client, std::vector<std::string> args, AsyncCliQueue::Callback callback) {
    if (auto *app = dynamic_cast<MultiFlexiApp *>(TProgram::application)) {
        app->cliQueue().enqueue(client, std::move(args), std::move(callback));
        app->setBusy(true);
        return true;
    }
    auto r = client.runJson(args);
    if (callback) {
        callback(std::move(r));
    }
    return false;
}

} // namespace multiflexitui
