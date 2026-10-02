#pragma once

#include <functional>
#include <map>
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

class ProcessRunner {
public:
    static ProcessResult run(const std::vector<std::string> &argv,
                             const std::map<std::string, std::string> &extraEnv = {});
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

    // Convenience helpers matching the former Go client.
    Result list(const std::string &entity, int limit, int offset) const;
    Result get(const std::string &entity, int id) const;
    Result create(const std::string &entity, const std::vector<std::string> &args) const;
    Result update(const std::string &entity, const std::vector<std::string> &args) const;
    Result remove(const std::string &entity, const std::string &deleteAction, int id) const;
    Result status() const;

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
