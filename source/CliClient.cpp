#include "multiflexitui/CliClient.h"

#include <sys/stat.h>
#include <sys/wait.h>
#include <fcntl.h>
#include <poll.h>
#include <unistd.h>

#include <chrono>
#include <cstdlib>
#include <ctime>
#include <filesystem>
#include <sstream>
#include <thread>
#include <utility>

namespace multiflexitui {

namespace {

std::string cliLogPath() {
    const char *xdg = std::getenv("XDG_STATE_HOME");
    std::string dir;
    if (xdg != nullptr && xdg[0] != '\0') {
        dir = xdg;
    } else {
        const char *home = std::getenv("HOME");
        if (home == nullptr || home[0] == '\0') {
            return {};
        }
        dir = std::string(home) + "/.local/state";
    }
    return dir + "/multiflexi-tui/multiflexi-tui.log";
}

std::string localTimestamp() {
    const std::time_t now = std::chrono::system_clock::to_time_t(std::chrono::system_clock::now());
    std::tm local {};
    localtime_r(&now, &local);
    char buf[32];
    std::strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", &local);
    return buf;
}

std::string joinArgv(const std::vector<std::string> &argv) {
    std::ostringstream out;
    for (std::size_t i = 0; i < argv.size(); ++i) {
        if (i != 0) {
            out << ' ';
        }
        out << argv[i];
    }
    return out.str();
}

std::string capturedText(const std::string &text) {
    if (text.empty()) {
        return "(empty)\n";
    }
    constexpr std::size_t kMax = 8192;
    std::string body = text.size() > kMax ? text.substr(0, kMax) + "\n... truncated\n" : text;
    if (body.back() != '\n') {
        body.push_back('\n');
    }
    return body;
}

std::string appendCliLog(const std::string &error, const std::vector<std::string> &argv,
                         const ProcessResult &pr) {
    const std::string path = cliLogPath();
    if (path.empty()) {
        return {};
    }
    const std::filesystem::path file(path);
    std::error_code ec;
    std::filesystem::create_directories(file.parent_path(), ec);
    if (ec) {
        return {};
    }
    ::chmod(file.parent_path().c_str(), 0700);
    const int fd = ::open(path.c_str(), O_WRONLY | O_CREAT | O_APPEND, 0600);
    if (fd < 0) {
        return {};
    }
    const std::string entry = localTimestamp() + " " + error + "\ncommand: " + joinArgv(argv) +
                              "\nexit: " + std::to_string(pr.exitCode) + "\nstdout:\n" +
                              capturedText(pr.stdOut) + "stderr:\n" + capturedText(pr.stdErr) + "\n";
    const ssize_t written = ::write(fd, entry.data(), entry.size());
    ::close(fd);
    return written == static_cast<ssize_t>(entry.size()) ? path : std::string();
}

void readAll(int fd, std::string &out) {
    char buf[4096];
    ssize_t n;
    while ((n = read(fd, buf, sizeof(buf))) > 0) {
        out.append(buf, static_cast<std::size_t>(n));
    }
}

void setCloexecNonblock(int fd) {
    int flags = fcntl(fd, F_GETFL, 0);
    if (flags >= 0) {
        fcntl(fd, F_SETFL, flags | O_NONBLOCK);
    }
    flags = fcntl(fd, F_GETFD, 0);
    if (flags >= 0) {
        fcntl(fd, F_SETFD, flags | FD_CLOEXEC);
    }
}

} // namespace

ProcessResult ProcessRunner::run(const std::vector<std::string> &argv,
                                 const std::map<std::string, std::string> &extraEnv) {
    ProcessResult result;
    if (argv.empty()) {
        result.spawnFailed = true;
        result.stdErr = "empty argv";
        return result;
    }

    int outPipe[2];
    int errPipe[2];
    if (pipe(outPipe) != 0) {
        result.spawnFailed = true;
        result.stdErr = "pipe() failed";
        return result;
    }
    if (pipe(errPipe) != 0) {
        close(outPipe[0]);
        close(outPipe[1]);
        result.spawnFailed = true;
        result.stdErr = "pipe() failed";
        return result;
    }

    pid_t pid = fork();
    if (pid < 0) {
        close(outPipe[0]);
        close(outPipe[1]);
        close(errPipe[0]);
        close(errPipe[1]);
        result.spawnFailed = true;
        result.stdErr = "fork() failed";
        return result;
    }

    if (pid == 0) {
        dup2(outPipe[1], STDOUT_FILENO);
        dup2(errPipe[1], STDERR_FILENO);
        close(outPipe[0]);
        close(outPipe[1]);
        close(errPipe[0]);
        close(errPipe[1]);
        for (const auto &item : extraEnv) {
            ::setenv(item.first.c_str(), item.second.c_str(), 1);
        }
        std::vector<char *> cargv;
        cargv.reserve(argv.size() + 1);
        for (const auto &arg : argv) {
            cargv.push_back(const_cast<char *>(arg.c_str()));
        }
        cargv.push_back(nullptr);
        execvp(cargv[0], cargv.data());
        _exit(127);
    }

    close(outPipe[1]);
    close(errPipe[1]);
    std::string outBuf;
    std::string errBuf;
    std::thread outThread([&]() { readAll(outPipe[0], outBuf); });
    std::thread errThread([&]() { readAll(errPipe[0], errBuf); });
    outThread.join();
    errThread.join();
    close(outPipe[0]);
    close(errPipe[0]);

    int status = 0;
    waitpid(pid, &status, 0);
    result.stdOut = std::move(outBuf);
    result.stdErr = std::move(errBuf);
    if (WIFEXITED(status)) {
        result.exitCode = WEXITSTATUS(status);
        if (result.exitCode == 127 && result.stdOut.empty()) {
            result.spawnFailed = true;
        }
    } else {
        result.exitCode = -1;
        result.spawnFailed = true;
    }
    return result;
}

ProcessHandle::ProcessHandle(ProcessHandle &&other) noexcept {
    *this = std::move(other);
}

ProcessHandle &ProcessHandle::operator=(ProcessHandle &&other) noexcept {
    if (this == &other) {
        return *this;
    }
    closeFds();
    if (pid_ > 0 && !done_) {
        // Best-effort reaping on move-from leave; prefer draining via pump.
        int status = 0;
        waitpid(pid_, &status, WNOHANG);
    }
    pid_ = other.pid_;
    outFd_ = other.outFd_;
    errFd_ = other.errFd_;
    outEof_ = other.outEof_;
    errEof_ = other.errEof_;
    done_ = other.done_;
    result_ = std::move(other.result_);
    other.pid_ = -1;
    other.outFd_ = -1;
    other.errFd_ = -1;
    other.outEof_ = true;
    other.errEof_ = true;
    other.done_ = true;
    return *this;
}

ProcessHandle::~ProcessHandle() {
    closeFds();
    if (pid_ > 0 && !done_) {
        int status = 0;
        waitpid(pid_, &status, 0);
        done_ = true;
    }
}

void ProcessHandle::closeFds() {
    if (outFd_ >= 0) {
        close(outFd_);
        outFd_ = -1;
    }
    if (errFd_ >= 0) {
        close(errFd_);
        errFd_ = -1;
    }
}

void ProcessHandle::drainFd(int &fd, std::string &buf, bool &eof) {
    if (fd < 0 || eof) {
        return;
    }
    char tmp[4096];
    for (;;) {
        const ssize_t n = read(fd, tmp, sizeof(tmp));
        if (n > 0) {
            buf.append(tmp, static_cast<std::size_t>(n));
            continue;
        }
        if (n == 0) {
            eof = true;
            close(fd);
            fd = -1;
            return;
        }
        if (errno == EINTR) {
            continue;
        }
        if (errno == EAGAIN || errno == EWOULDBLOCK) {
            return;
        }
        eof = true;
        close(fd);
        fd = -1;
        return;
    }
}

ProcessHandle ProcessHandle::start(const std::vector<std::string> &argv,
                                   const std::map<std::string, std::string> &extraEnv) {
    ProcessHandle h;
    if (argv.empty()) {
        h.done_ = true;
        h.result_.spawnFailed = true;
        h.result_.stdErr = "empty argv";
        return h;
    }

    int outPipe[2];
    int errPipe[2];
    if (pipe(outPipe) != 0 || pipe(errPipe) != 0) {
        if (outPipe[0] >= 0) {
            close(outPipe[0]);
            close(outPipe[1]);
        }
        h.done_ = true;
        h.result_.spawnFailed = true;
        h.result_.stdErr = "pipe() failed";
        return h;
    }

    pid_t pid = fork();
    if (pid < 0) {
        close(outPipe[0]);
        close(outPipe[1]);
        close(errPipe[0]);
        close(errPipe[1]);
        h.done_ = true;
        h.result_.spawnFailed = true;
        h.result_.stdErr = "fork() failed";
        return h;
    }

    if (pid == 0) {
        dup2(outPipe[1], STDOUT_FILENO);
        dup2(errPipe[1], STDERR_FILENO);
        close(outPipe[0]);
        close(outPipe[1]);
        close(errPipe[0]);
        close(errPipe[1]);
        for (const auto &item : extraEnv) {
            ::setenv(item.first.c_str(), item.second.c_str(), 1);
        }
        std::vector<char *> cargv;
        cargv.reserve(argv.size() + 1);
        for (const auto &arg : argv) {
            cargv.push_back(const_cast<char *>(arg.c_str()));
        }
        cargv.push_back(nullptr);
        execvp(cargv[0], cargv.data());
        _exit(127);
    }

    close(outPipe[1]);
    close(errPipe[1]);
    setCloexecNonblock(outPipe[0]);
    setCloexecNonblock(errPipe[0]);
    h.pid_ = pid;
    h.outFd_ = outPipe[0];
    h.errFd_ = errPipe[0];
    h.outEof_ = false;
    h.errEof_ = false;
    h.done_ = false;
    return h;
}

bool ProcessHandle::pump() {
    if (done_) {
        return true;
    }

    drainFd(outFd_, result_.stdOut, outEof_);
    drainFd(errFd_, result_.stdErr, errEof_);

    int status = 0;
    const pid_t waited = waitpid(pid_, &status, WNOHANG);
    if (waited == 0) {
        // Still running; also try a short poll to avoid busy spinning callers.
        pollfd pfds[2];
        int nfds = 0;
        if (outFd_ >= 0) {
            pfds[nfds++] = pollfd{outFd_, POLLIN, 0};
        }
        if (errFd_ >= 0) {
            pfds[nfds++] = pollfd{errFd_, POLLIN, 0};
        }
        if (nfds > 0) {
            poll(pfds, nfds, 0);
            drainFd(outFd_, result_.stdOut, outEof_);
            drainFd(errFd_, result_.stdErr, errEof_);
        }
        return false;
    }

    // Child exited — drain remaining output.
    drainFd(outFd_, result_.stdOut, outEof_);
    drainFd(errFd_, result_.stdErr, errEof_);
    closeFds();
    outEof_ = true;
    errEof_ = true;

    if (waited < 0) {
        result_.exitCode = -1;
        result_.spawnFailed = true;
    } else if (WIFEXITED(status)) {
        result_.exitCode = WEXITSTATUS(status);
        if (result_.exitCode == 127 && result_.stdOut.empty()) {
            result_.spawnFailed = true;
        }
    } else {
        result_.exitCode = -1;
        result_.spawnFailed = true;
    }
    done_ = true;
    pid_ = -1;
    return true;
}

CliClient::CliClient(std::string binaryPath, std::string envFilePath)
    : binaryPath_(std::move(binaryPath)), envFilePath_(std::move(envFilePath)) {
}

void CliClient::setBinaryPath(std::string path) {
    binaryPath_ = std::move(path);
}

void CliClient::setEnvFile(std::string path) {
    envFilePath_ = std::move(path);
}

void CliClient::setRequestObserver(std::function<void(const std::string &)> observer) {
    requestObserver_ = std::move(observer);
}

std::vector<std::string> CliClient::buildArgv(const std::vector<std::string> &args) const {
    std::vector<std::string> argv;
    argv.push_back(binaryPath_);
    if (!envFilePath_.empty()) {
        argv.push_back("--envfile=" + envFilePath_);
    }
    for (const auto &a : args) {
        argv.push_back(a);
    }
    argv.push_back("--format=json");
    return argv;
}

CliClient::Result CliClient::parseProcessResult(const ProcessResult &pr, const std::string &binaryPath,
                                                const std::vector<std::string> &argv) {
    Result result;
    result.lastCommand = joinArgv(argv);
    result.exitCode = pr.exitCode;

    auto noteFailure = [&]() {
        const std::string logPath = appendCliLog(result.errorMessage, argv, pr);
        if (!logPath.empty()) {
            result.errorMessage += "\nLogged to " + logPath;
        }
    };

    if (pr.spawnFailed) {
        result.ok = false;
        result.errorMessage = binaryPath + " not found on PATH (or at the configured --cli path)";
        noteFailure();
        return result;
    }

    const std::string trimmed = [&]() {
        std::string s = pr.stdOut;
        while (!s.empty() && (s.back() == '\n' || s.back() == '\r' || s.back() == ' ')) {
            s.pop_back();
        }
        return s;
    }();

    if (trimmed.empty()) {
        result.data = nlohmann::json::object();
    } else {
        try {
            result.data = nlohmann::json::parse(trimmed);
        } catch (const nlohmann::json::parse_error &e) {
            result.ok = false;
            result.errorMessage = "Invalid JSON from " + binaryPath + ": " + e.what();
            if (!pr.stdErr.empty()) {
                result.errorMessage += " (stderr: " + pr.stdErr + ")";
            }
            noteFailure();
            return result;
        }
    }

    if (pr.exitCode != 0) {
        result.ok = false;
        if (result.data.is_object() && result.data.contains("message")) {
            result.errorMessage = result.data.value("message", std::string("Unknown error"));
        } else if (result.data.is_object() && result.data.contains("error")) {
            result.errorMessage = result.data.value("error", std::string("Unknown error"));
        } else if (!pr.stdErr.empty()) {
            result.errorMessage = pr.stdErr;
        } else {
            result.errorMessage = binaryPath + " exited with code " + std::to_string(pr.exitCode);
        }
        noteFailure();
        return result;
    }

    result.ok = true;
    return result;
}

CliClient::Result CliClient::runJson(const std::vector<std::string> &args) const {
    return runJsonImpl(args, nullptr);
}

CliClient::Result CliClient::runJsonWithEnv(const std::vector<std::string> &args,
                                            const std::map<std::string, std::string> &extraEnv) const {
    return runJsonImpl(args, &extraEnv);
}

CliClient::Result CliClient::runJsonImpl(const std::vector<std::string> &args,
                                         const std::map<std::string, std::string> *overrideEnv) const {
    std::vector<std::string> argv;
    argv.push_back(binaryPath_);
    if ((overrideEnv == nullptr) && !envFilePath_.empty()) {
        argv.push_back("--envfile=" + envFilePath_);
    }
    for (const auto &a : args) {
        argv.push_back(a);
    }
    argv.push_back("--format=json");

    lastCommand_ = joinArgv(argv);
    if (requestObserver_) {
        requestObserver_(lastCommand_);
    }

    const std::map<std::string, std::string> empty;
    const std::map<std::string, std::string> &env = overrideEnv != nullptr ? *overrideEnv : empty;
    ProcessResult pr = ProcessRunner::run(argv, env);
    Result result = parseProcessResult(pr, binaryPath_, argv);
    lastCommand_ = result.lastCommand;
    return result;
}

CliClient::Result CliClient::list(const std::string &entity, int limit, int offset) const {
    ListOptions opts;
    opts.limit = limit;
    opts.offset = offset;
    return list(entity, opts);
}

CliClient::Result CliClient::list(const std::string &entity, const ListOptions &opts) const {
    std::vector<std::string> args{entity + ":list", "--order=D",
                                  "--limit=" + std::to_string(opts.limit),
                                  "--offset=" + std::to_string(opts.offset)};
    if (opts.companyId > 0) {
        args.push_back("--company_id=" + std::to_string(opts.companyId));
    }
    if (!opts.filter.empty()) {
        if (entity == "job") {
            args.push_back("--status=" + opts.filter);
        } else if (entity == "task") {
            args.push_back("--state=" + opts.filter);
        } else if (entity == "artifact") {
            args.push_back("--job_id=" + opts.filter);
        } else if (entity == "conffield") {
            args.push_back("--app_id=" + opts.filter);
        } else {
            args.push_back("--filter=" + opts.filter);
        }
    }
    for (const auto &extra : opts.extraArgs) {
        args.push_back(extra);
    }
    return runJson(args);
}

CliClient::Result CliClient::get(const std::string &entity, int id) const {
    return runJson({entity + ":get", "--id=" + std::to_string(id)});
}

CliClient::Result CliClient::create(const std::string &entity, const std::vector<std::string> &args) const {
    std::vector<std::string> full{entity + ":create"};
    full.insert(full.end(), args.begin(), args.end());
    return runJson(full);
}

CliClient::Result CliClient::update(const std::string &entity, const std::vector<std::string> &args) const {
    std::vector<std::string> full{entity + ":update"};
    full.insert(full.end(), args.begin(), args.end());
    return runJson(full);
}

CliClient::Result CliClient::remove(const std::string &entity, const std::string &deleteAction, int id) const {
    return runJson({entity + ":" + deleteAction, "--id=" + std::to_string(id)});
}

CliClient::Result CliClient::status() const {
    return runJson({"status"});
}

} // namespace multiflexitui
