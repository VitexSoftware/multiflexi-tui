#pragma once

#include <functional>
#include <map>
#include <optional>
#include <string>
#include <vector>

#include <nlohmann/json.hpp>

namespace multiflexitui {

struct ProcessResult {
    int exitCode = -1;
    std::string stdOut;
    std::string stdErr;
    bool spawnFailed = false;
};

// Blocking one-shot runner (kept for modal actions / tests).
class ProcessRunner {
public:
    static ProcessResult run(const std::vector<std::string> &argv,
                             const std::map<std::string, std::string> &extraEnv = {});
};

// Non-blocking child process: fork+execvp, drain pipes via pump().
class ProcessHandle {
public:
    ProcessHandle() = default;
    ProcessHandle(const ProcessHandle &) = delete;
    ProcessHandle &operator=(const ProcessHandle &) = delete;
    ProcessHandle(ProcessHandle &&other) noexcept;
    ProcessHandle &operator=(ProcessHandle &&other) noexcept;
    ~ProcessHandle();

    static ProcessHandle start(const std::vector<std::string> &argv,
                               const std::map<std::string, std::string> &extraEnv = {});

    // Returns true when the child has exited and pipes are drained.
    bool pump();
    bool done() const { return done_; }
    bool valid() const { return pid_ > 0 || done_; }
    ProcessResult result() const { return result_; }

private:
    void closeFds();
    void drainFd(int &fd, std::string &buf, bool &eof);

    pid_t pid_ = -1;
    int outFd_ = -1;
    int errFd_ = -1;
    bool outEof_ = true;
    bool errEof_ = true;
    bool done_ = false;
    ProcessResult result_;
};

struct ListOptions {
    int limit = 40;
    int offset = 0;
    int companyId = 0;
    std::string filter;                 // entity-specific (e.g. job status, task state)
    std::vector<std::string> extraArgs; // raw --key=value extras
};

// Thin JSON-returning wrapper around multiflexi-cli. Always appends --format=json.
class CliClient {
public:
    explicit CliClient(std::string binaryPath = "multiflexi-cli", std::string envFilePath = "");

    struct Result {
        bool ok = false;
        nlohmann::json data;
        std::string errorMessage;
        int exitCode = -1;
        std::string lastCommand;
    };

    Result runJson(const std::vector<std::string> &args) const;
    Result runJsonWithEnv(const std::vector<std::string> &args,
                          const std::map<std::string, std::string> &extraEnv) const;

    Result list(const std::string &entity, int limit, int offset) const;
    Result list(const std::string &entity, const ListOptions &opts) const;
    Result get(const std::string &entity, int id) const;
    Result create(const std::string &entity, const std::vector<std::string> &args) const;
    Result update(const std::string &entity, const std::vector<std::string> &args) const;
    Result remove(const std::string &entity, const std::string &deleteAction, int id) const;
    Result status() const;

    // Build argv (binary + envfile + args + --format=json) without executing.
    std::vector<std::string> buildArgv(const std::vector<std::string> &args) const;
    static Result parseProcessResult(const ProcessResult &pr, const std::string &binaryPath,
                                     const std::vector<std::string> &argv);

    void setBinaryPath(std::string path);
    void setEnvFile(std::string path);
    void setRequestObserver(std::function<void(const std::string &)> observer);

    const std::string &binaryPath() const { return binaryPath_; }
    const std::string &envFile() const { return envFilePath_; }
    const std::string &lastCommand() const { return lastCommand_; }

private:
    Result runJsonImpl(const std::vector<std::string> &args,
                       const std::map<std::string, std::string> *overrideEnv) const;

    std::string binaryPath_;
    std::string envFilePath_;
    mutable std::string lastCommand_;
    mutable std::function<void(const std::string &)> requestObserver_;
};

} // namespace multiflexitui
