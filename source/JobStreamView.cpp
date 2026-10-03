#include "multiflexitui/TV.h"
#include "multiflexitui/JobStreamView.h"
#include "multiflexitui/AppButton.h"
#include "multiflexitui/AppShell.h"
#include "multiflexitui/AsyncCliQueue.h"
#include "multiflexitui/Commands.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/WindowColors.h"
#include "multiflexitui/WindowLayout.h"
#include "multiflexitui/i18n.h"

#include <sstream>

namespace multiflexitui {

JobStreamView::~JobStreamView() {
    if (auto *mf = dynamic_cast<MultiFlexiApp *>(TProgram::application)) {
        mf->removeStreamTick(this);
    }
}

JobStreamView::JobStreamView(CliClient &client, int jobId)
    : TWindowInit(&TWindow::initFrame),
      TWindow(TRect(2, 1, 78, 23), ("Job stream #" + std::to_string(jobId)).c_str(), wnNoNumber),
      client_(client),
      jobId_(jobId) {
    flags |= wfGrow | wfZoom | wfClose;
    options |= ofTileable;

    meta_ = new TStaticText(TRect(2, 1, 76, 2), "");
    growWide(meta_);
    insert(meta_);

    TScrollBar *vBar = standardScrollBar(sbVertical | sbHandleKeyboard);
    TScrollBar *hBar = standardScrollBar(sbHorizontal | sbHandleKeyboard);
    out_ = new SimpleListViewer(TRect(1, 2, 76, 18), vBar, hBar);
    growFill(out_);
    insert(out_);

    TView *ref = new AppButton(TRect(2, 19, 16, 21), _("~R~efresh"), cmStreamRefresh, bfNormal);
    stickBottom(ref);
    insert(ref);
    TView *fol = new AppButton(TRect(18, 19, 36, 21), _("~F~ollow"), cmStreamToggleFollow, bfNormal);
    stickBottom(fol);
    insert(fol);

    refresh(true);
    selectNext(False);
}

void JobStreamView::draw() {
    TWindow::draw();
    TDrawBuffer b;
    const TColorAttr color = mapColor(6);
    b.moveChar(0, ' ', color, static_cast<ushort>(size.x));
    b.moveStr(2, metaText_, color, static_cast<ushort>(size.x > 4 ? size.x - 4 : 0));
    writeLine(0, 1, size.x, 1, b);
}

void JobStreamView::applyJobResult(const CliClient::Result &r, bool /*force*/) {
    pending_ = false;
    if (!r.ok) {
        metaText_ = "Error: " + r.errorMessage;
        drawView();
        return;
    }

    stdout_ = jsonToString(r.data.contains("stdout") ? r.data["stdout"] : nlohmann::json());
    stderr_ = jsonToString(r.data.contains("stderr") ? r.data["stderr"] : nlohmann::json());
    const std::string begin = jsonToString(r.data.contains("begin") ? r.data["begin"] : nlohmann::json());
    const std::string end = jsonToString(r.data.contains("end") ? r.data["end"] : nlohmann::json());
    const std::string exitcode =
        jsonToString(r.data.contains("exitcode") ? r.data["exitcode"] : nlohmann::json());
    const int pid = jsonToInt(r.data, "pid");

    finished_ = !end.empty() || (!exitcode.empty() && exitcode != "null");
    std::ostringstream meta;
    meta << "job=" << jobId_ << " pid=" << pid << " begin=" << begin << " end=" << end
         << " exit=" << exitcode << (follow_ ? " [follow]" : " [paused]")
         << (finished_ ? " DONE" : " LIVE");
    metaText_ = meta.str();

    std::vector<std::string> lines;
    lines.push_back("--- stdout ---");
    {
        std::istringstream in(stdout_);
        std::string line;
        while (std::getline(in, line)) {
            lines.push_back(line);
        }
    }
    lines.push_back("--- stderr ---");
    {
        std::istringstream in(stderr_);
        std::string line;
        while (std::getline(in, line)) {
            lines.push_back(line);
        }
    }
    if (out_ != nullptr) {
        const short prev = out_->focused;
        out_->setRows(std::move(lines));
        if (follow_ && out_->rowCount() > 0) {
            out_->focusItem(static_cast<short>(out_->rowCount() - 1));
        } else if (prev >= 0) {
            out_->focusItem(prev);
        }
    }
    drawView();
    lastPoll_ = std::chrono::steady_clock::now();
}

void JobStreamView::refresh(bool force) {
    if (pending_) {
        return;
    }
    pending_ = true;
    auto *self = this;
    enqueueCliJson(client_, {"job:get", "--id=" + std::to_string(jobId_)},
                   [self, force](CliClient::Result r) {
                       if (auto *app = dynamic_cast<MultiFlexiApp *>(TProgram::application)) {
                           app->setBusy(app->cliQueue().busy());
                       }
                       bool alive = false;
                       if (TProgram::deskTop != nullptr) {
                           TView *p = TProgram::deskTop->first();
                           while (p != nullptr) {
                               if (p == self) {
                                   alive = true;
                                   break;
                               }
                               p = p->nextView();
                           }
                       }
                       if (alive) {
                           self->applyJobResult(r, force);
                       }
                   });
}

void JobStreamView::tick() {
    if (finished_ || !follow_ || pending_) {
        return;
    }
    const auto now = std::chrono::steady_clock::now();
    if (now - lastPoll_ < std::chrono::seconds(1)) {
        return;
    }
    refresh(false);
}

void JobStreamView::handleEvent(TEvent &event) {
    TWindow::handleEvent(event);
    if (event.what == evCommand) {
        if (event.message.command == cmStreamRefresh) {
            refresh(true);
            clearEvent(event);
        } else if (event.message.command == cmStreamToggleFollow) {
            follow_ = !follow_;
            metaText_ += follow_ ? " [follow on]" : " [follow off]";
            drawView();
            clearEvent(event);
        }
    }
}

TColorAttr JobStreamView::mapColor(uchar index) {
    TColorAttr color;
    return windowColor(index, color) ? color : TWindow::mapColor(index);
}

void openJobStreamForId(TProgram *app, CliClient &client, int jobId) {
    auto *view = new JobStreamView(client, jobId);
    app->deskTop->insert(view);
    if (auto *mf = dynamic_cast<MultiFlexiApp *>(app)) {
        mf->addStreamTick(view, [view]() { view->tick(); });
    }
}

void openJobStreamPicker(TProgram *app, CliClient &client) {
    char buf[32] = {};
    if (inputBox(_("Job ID"), _("Open live job stream"), buf, sizeof(buf) - 1) != cmOK) {
        return;
    }
    int id = 0;
    try {
        id = std::stoi(buf);
    } catch (...) {
        messageBox(_("Invalid job id"), mfError | mfOKButton);
        return;
    }
    openJobStreamForId(app, client, id);
}

} // namespace multiflexitui
