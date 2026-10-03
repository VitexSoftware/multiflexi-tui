#include "multiflexitui/TV.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/EntityListView.h"
#include "multiflexitui/EntityPicker.h"
#include "multiflexitui/AppShell.h"
#include "multiflexitui/JobStreamView.h"
#include "multiflexitui/i18n.h"

#include <cstdlib>
#include <cstring>
#include <filesystem>
#include <sstream>

namespace multiflexitui {

namespace {

std::vector<EntityDef> g_entities;

std::string pad(const std::string &s, int width) {
    if (static_cast<int>(s.size()) >= width) {
        return s.substr(0, static_cast<std::size_t>(width));
    }
    return s + std::string(static_cast<std::size_t>(width) - s.size(), ' ');
}

bool showText(CliClient &client, const std::string &title, const std::string &body) {
    (void)client;
    TProgram::application->executeDialog(new TextViewerDialog(title, body));
    return false;
}

bool showJobOutput(CliClient &client, const nlohmann::json &row, const char *field, const char *label) {
    const int id = jsonToInt(row, "id");
    auto r = client.get("job", id);
    if (!r.ok) {
        messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
        return false;
    }
    std::string content = jsonToString(r.data.contains(field) ? r.data[field] : nlohmann::json());
    if (content.empty()) {
        content = "(empty)";
    }
    return showText(client, std::string("Job ") + std::to_string(id) + " — " + label, content);
}

bool promptPath(const std::string &title, const std::string &label, std::string &path) {
    char buf[240];
    std::strncpy(buf, path.c_str(), sizeof(buf) - 1);
    buf[sizeof(buf) - 1] = '\0';
    if (inputBox(title.c_str(), label.c_str(), buf, sizeof(buf) - 1) != cmOK) {
        return false;
    }
    path = buf;
    return !path.empty();
}

FieldDef field(const char *key, const char *label, const char *placeholder = "", bool required = false,
               bool secret = false, const char *relation = "") {
    FieldDef f;
    f.key = key;
    f.label = label;
    f.placeholder = placeholder;
    f.required = required;
    f.secret = secret;
    f.relationEntity = relation;
    return f;
}

} // namespace

nlohmann::json asJsonArray(const nlohmann::json &data) {
    if (data.is_array()) {
        return data;
    }
    if (data.is_object() && data.contains("data") && data["data"].is_array()) {
        return data["data"];
    }
    if (data.is_object()) {
        return nlohmann::json::array({data});
    }
    return nlohmann::json::array();
}

std::string defaultExportPath(const std::string &filename) {
    const char *xdg = std::getenv("XDG_STATE_HOME");
    std::string dir;
    if (xdg != nullptr && xdg[0] != '\0') {
        dir = std::string(xdg) + "/multiflexi";
    } else {
        const char *home = std::getenv("HOME");
        if (home != nullptr && home[0] != '\0') {
            dir = std::string(home) + "/.local/state/multiflexi";
        } else {
            dir = ".";
        }
    }
    std::error_code ec;
    std::filesystem::create_directories(dir, ec);
    return dir + "/" + filename;
}

std::string jsonToString(const nlohmann::json &v) {
    if (v.is_null()) {
        return {};
    }
    if (v.is_string()) {
        return v.get<std::string>();
    }
    if (v.is_boolean()) {
        return v.get<bool>() ? "true" : "false";
    }
    if (v.is_number_integer()) {
        return std::to_string(v.get<long long>());
    }
    if (v.is_number_float()) {
        return std::to_string(v.get<double>());
    }
    return v.dump();
}

int jsonToInt(const nlohmann::json &obj, const std::string &key, int fallback) {
    if (!obj.is_object() || !obj.contains(key) || obj[key].is_null()) {
        return fallback;
    }
    if (obj[key].is_number_integer()) {
        return obj[key].get<int>();
    }
    if (obj[key].is_string()) {
        try {
            return std::stoi(obj[key].get<std::string>());
        } catch (...) {
            return fallback;
        }
    }
    return fallback;
}

std::string formatHeader(const EntityDef &def) {
    std::ostringstream out;
    for (const auto &col : def.columns) {
        out << pad(col.title, col.width) << ' ';
    }
    return out.str();
}

std::string formatRow(const EntityDef &def, const nlohmann::json &obj) {
    std::ostringstream out;
    for (const auto &col : def.columns) {
        std::string val;
        if (obj.is_object() && obj.contains(col.key)) {
            val = jsonToString(obj[col.key]);
        }
        out << pad(val, col.width) << ' ';
    }
    return out.str();
}

std::vector<std::string> buildCliArgs(const std::vector<FieldDef> &fields,
                                      const std::map<std::string, std::string> &values, int id) {
    std::vector<std::string> args;
    if (id >= 0) {
        args.push_back("--id=" + std::to_string(id));
    }
    for (const auto &f : fields) {
        auto it = values.find(f.key);
        if (it == values.end()) {
            continue;
        }
        if (it->second.empty() && !f.required) {
            continue;
        }
        args.push_back("--" + f.key + "=" + it->second);
    }
    return args;
}

const std::vector<EntityDef> &allEntities() {
    return g_entities;
}

const EntityDef *findEntity(const std::string &cliEntity) {
    for (const auto &e : g_entities) {
        if (e.cliEntity == cliEntity) {
            return &e;
        }
    }
    return nullptr;
}

void registerBuiltinEntities() {
    g_entities.clear();

    {
        EntityDef e;
        e.name = _("Companies");
        e.cliEntity = "company";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", _("Name"), 24}, {"slug", "Slug", 14}, {"ic", "IC", 12},
                     {"enabled", "En", 4}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("slug", "Slug", "", true),
            field("ic", "IC"),
            field("email", "Email"),
            field("enabled", _("Enabled"), "1"),
        };
        e.editFields = e.createFields;
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Applications");
        e.cliEntity = "application";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", _("Name"), 22}, {"code", "Code", 12}, {"version", "Ver", 8},
                     {"enabled", "En", 4}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("executable", "Executable", "", true),
            field("code", "Code"),
            field("description", _("Description")),
            field("homepage", "Homepage"),
            field("ociimage", "OCI Image"),
            field("version", "Version"),
            field("enabled", _("Enabled"), "1"),
        };
        e.editFields = e.createFields;
        e.rowActions.push_back({_("Config fields"), 'c', [](CliClient &client, const nlohmann::json &row) {
            const int appId = jsonToInt(row, "id");
            auto r = client.runJson(
                {"conffield:list", "--app_id=" + std::to_string(appId), "--limit=200", "--offset=0"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            std::ostringstream body;
            for (const auto &item : asJsonArray(r.data)) {
                body << jsonToString(item.contains("id") ? item["id"] : nlohmann::json()) << "  "
                     << jsonToString(item.contains("keyname") ? item["keyname"] : nlohmann::json()) << "  "
                     << jsonToString(item.contains("type") ? item["type"] : nlohmann::json()) << "\n";
            }
            showText(client, "Config fields for app " + std::to_string(appId), body.str());
            return false;
        }});
        e.rowActions.push_back({_("Show config"), 's', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"application:show-config", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "application:show-config", r.data.dump(2));
            return false;
        }});
        e.rowActions.push_back({_("Export JSON"), 'x', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            std::string path = defaultExportPath("multiflexi-app-" + std::to_string(id) + ".json");
            if (!promptPath(_("Export JSON"), _("Output path"), path)) {
                return false;
            }
            auto r = client.runJson(
                {"application:export-json", "--id=" + std::to_string(id), "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(("Exported to " + path).c_str(), mfInformation | mfOKButton);
            return false;
        }});
        e.listActions.push_back({_("Remove JSON"), 'R', [](CliClient &client, const nlohmann::json &) {
            std::string path = defaultExportPath("app.json");
            if (!promptPath(_("Remove JSON"), _("Path to app JSON"), path)) {
                return false;
            }
            auto r = client.runJson({"application:remove-json", "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Application removed from JSON definition."), mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Config Fields");
        e.cliEntity = "conffield";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"app_id", "App", 6}, {"keyname", "Key", 18}, {"type", "Type", 10},
                     {"required", "Req", 4}};
        e.createFields = {
            field("app_id", "App ID", "", true, false, "application"),
            field("keyname", "Key name", "", true),
            field("type", "Type", "string", true),
            field("description", _("Description")),
            field("defval", "Default"),
            field("required", "Required", "0"),
        };
        e.editFields = e.createFields;
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("RunTemplates");
        e.cliEntity = "run-template";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.supportsCompanyScope = true;
        e.columns = {{"id", "ID", 6}, {"name", _("Name"), 22}, {"company_id", "Co", 5}, {"app_id", "App", 5},
                     {"interv", "Interval", 10}, {"active", "Act", 4}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("app_id", "App ID", "", true, false, "application"),
            field("company_id", "Company ID", "", true, false, "company"),
            field("interv", "Interval", "n"),
            field("executor", "Executor", "Native"),
            field("active", "Active", "0"),
            field("cron", "Cron"),
        };
        e.editFields = e.createFields;
        e.rowActions.push_back({_("Schedule"), 's', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"run-template:schedule", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(("Scheduled run-template " + std::to_string(id)).c_str(), mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({_("Clone"), 'c', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"run-template:clone", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Clone created (disabled)."), mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({_("List credentials"), 'l', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"run-template:list-credentials", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "Credentials for RT " + std::to_string(id), r.data.dump(2));
            return false;
        }});
        e.rowActions.push_back({_("Assign credential"), 'a', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            const int credId = pickEntityId(client, "credential", _("Pick credential"));
            if (credId <= 0) {
                return false;
            }
            auto r = client.runJson({"run-template:assign-credential", "--id=" + std::to_string(id),
                                     "--credential_id=" + std::to_string(credId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Credential assigned."), mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({_("Unassign credential"), 'u', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            const int credId = pickEntityId(client, "credential", _("Pick credential to unassign"));
            if (credId <= 0) {
                return false;
            }
            auto r = client.runJson({"run-template:unassign-credential", "--id=" + std::to_string(id),
                                     "--credential_id=" + std::to_string(credId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Credential unassigned."), mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Jobs");
        e.cliEntity = "job";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 8}, {"runtemplate_id", "RT", 6}, {"company_id", "Co", 5},
                     {"exitcode", "Exit", 5}, {"begin", "Begin", 20}, {"executor", "Exec", 10}};
        e.createFields = {
            field("runtemplate_id", "RunTemplate ID", "", true, false, "run-template"),
            field("scheduled", "Scheduled", "now", true),
            field("executor", "Executor", "Native"),
            field("schedule_type", "Schedule type", "adhoc"),
        };
        e.editFields = {
            field("exitcode", "Exit code"),
            field("executor", "Executor"),
        };
        e.rowActions.push_back({_("Stdout"), 'o', [](CliClient &c, const nlohmann::json &row) {
            return showJobOutput(c, row, "stdout", "Stdout");
        }});
        e.rowActions.push_back({_("Stderr"), 'e', [](CliClient &c, const nlohmann::json &row) {
            return showJobOutput(c, row, "stderr", "Stderr");
        }});
        e.rowActions.push_back({_("Live stream"), 's', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            openJobStreamForId(TProgram::application, client, id);
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Tasks");
        e.cliEntity = "task";
        e.deleteAction = "";
        e.columns = {{"id", "ID", 8}, {"runtemplate_id", "RT", 6}, {"state", "State", 12},
                     {"window_start", "Start", 20}, {"deadline", "Deadline", 20}};
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Credentials");
        e.cliEntity = "credential";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", _("Name"), 24}, {"company_id", "Co", 5},
                     {"credential_type_id", "Type", 6}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("company_id", "Company ID", "", true, false, "company"),
            field("credential_type_id", "Credential type ID", "", true, false, "credential-type"),
        };
        e.editFields = {field("name", _("Name"), "", true)};
        e.listActions.push_back({_("Encrypt existing"), 'E', [](CliClient &client, const nlohmann::json &) {
            auto r = client.runJson({"credential:encrypt-existing"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Encryption pass completed."), mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Tokens");
        e.cliEntity = "token";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"user", "User", 16}, {"user_id", "UID", 5}, {"start", "Start", 20},
                     {"until", "Until", 20}};
        e.createFields = {
            field("user_id", "User ID", "", true, false, "user"),
            field("until", "Until"),
        };
        e.editFields = {field("until", "Until")};
        e.rowActions.push_back({_("Generate"), 'g', [](CliClient &client, const nlohmann::json &row) {
            const int userId = jsonToInt(row, "user_id");
            auto r = client.runJson({"token:generate", "--user_id=" + std::to_string(userId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "Generated token", r.data.dump(2));
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Users");
        e.cliEntity = "user";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"login", "Login", 16}, {"email", "Email", 24}, {"enabled", "En", 4}};
        e.createFields = {
            field("login", "Login", "", true),
            field("email", "Email", "", true),
            field("firstname", "First name"),
            field("lastname", "Last name"),
            field("password", _("Password"), "", true, true),
            field("enabled", _("Enabled"), "1"),
        };
        e.editFields = {
            field("email", "Email"),
            field("firstname", "First name"),
            field("lastname", "Last name"),
            field("enabled", _("Enabled")),
        };
        e.rowActions.push_back({_("Roles"), 'r', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-role:list", "--user_id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "Roles for user " + std::to_string(id), r.data.dump(2));
            return false;
        }});
        e.rowActions.push_back({_("Make Admin"), 'a', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-role:set", "--user_id=" + std::to_string(id), "--roles=admin"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Admin role set."), mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({_("Assign company"), 'c', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            const int companyId = pickEntityId(client, "company", _("Pick company"));
            if (companyId <= 0) {
                return false;
            }
            auto r = client.runJson({"user-company:assign", "--user_id=" + std::to_string(id),
                                     "--company_id=" + std::to_string(companyId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Assigned."), mfInformation | mfOKButton);
            return false;
        }});
        e.rowActions.push_back({_("Unassign company"), 'u', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            const int companyId = pickEntityId(client, "company", _("Pick company to unassign"));
            if (companyId <= 0) {
                return false;
            }
            auto r = client.runJson({"user-company:unassign", "--user_id=" + std::to_string(id),
                                     "--company_id=" + std::to_string(companyId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Unassigned."), mfInformation | mfOKButton);
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Artifacts");
        e.cliEntity = "artifact";
        e.deleteAction = "";
        e.columns = {{"id", "ID", 6}, {"job_id", "Job", 8}, {"filename", "Filename", 28},
                     {"content_type", "Type", 18}};
        e.rowActions.push_back({_("Save"), 's', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            std::string filename = jsonToString(row.contains("filename") ? row["filename"] : nlohmann::json());
            if (filename.empty()) {
                filename = "artifact-" + std::to_string(id) + ".bin";
            }
            std::string path = defaultExportPath(filename);
            if (!promptPath(_("Save artifact"), _("Output path"), path)) {
                return false;
            }
            auto r = client.runJson({"artifact:save", "--id=" + std::to_string(id), "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(("Saved to " + path).c_str(), mfInformation | mfOKButton);
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Credential Types");
        e.cliEntity = "credential-type";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", _("Name"), 22}, {"prototype", "Prototype", 16},
                     {"company_id", "Co", 5}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("company_id", "Company ID", "", true, false, "company"),
            field("prototype", "Prototype code", "", true),
        };
        e.editFields = {field("name", _("Name"), "", true), field("url", "URL")};
        e.listActions.push_back({_("Import JSON"), 'I', [](CliClient &client, const nlohmann::json &) {
            std::string path = defaultExportPath("credential-type.json");
            if (!promptPath(_("Import credential type"), _("JSON path"), path)) {
                return false;
            }
            auto r = client.runJson({"credential-type:import-json", "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.listActions.push_back({_("Validate JSON"), 'V', [](CliClient &client, const nlohmann::json &) {
            std::string path = defaultExportPath("credential-type.json");
            if (!promptPath(_("Validate credential type"), _("JSON path"), path)) {
                return false;
            }
            auto r = client.runJson({"credential-type:validate-json", "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "credential-type:validate-json", r.data.dump(2));
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Cred Prototypes");
        e.cliEntity = "credential-prototype";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"code", "Code", 16}, {"name", _("Name"), 24}, {"version", "Ver", 8}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("code", "Code", "", true),
            field("description", _("Description")),
        };
        e.editFields = e.createFields;
        e.listActions.push_back({_("Sync all"), 'S', [](CliClient &client, const nlohmann::json &) {
            auto r = client.runJson({"credential-prototype:sync"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Prototypes synchronized."), mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({_("Export JSON"), 'x', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            std::string path = defaultExportPath("prototype-" + std::to_string(id) + ".json");
            if (!promptPath(_("Export prototype"), _("Output path"), path)) {
                return false;
            }
            auto r = client.runJson(
                {"credential-prototype:export-json", "--id=" + std::to_string(id), "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(("Exported to " + path).c_str(), mfInformation | mfOKButton);
            return false;
        }});
        e.listActions.push_back({_("Import JSON"), 'I', [](CliClient &client, const nlohmann::json &) {
            std::string path = defaultExportPath("prototype.json");
            if (!promptPath(_("Import prototype"), _("JSON path"), path)) {
                return false;
            }
            auto r = client.runJson({"credential-prototype:import-json", "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Company Apps");
        e.cliEntity = "company-app";
        e.deleteAction = "";
        e.supportsGet = false;
        e.supportsCompanyScope = true;
        e.columns = {{"id", "ID", 6}, {"company_id", "Co", 5}, {"company_name", _("Company"), 18},
                     {"app_id", "App", 5}, {"app_name", _("Application"), 22}};
        e.listActions.push_back({_("Assign"), 'a', [](CliClient &client, const nlohmann::json &) {
            const int companyId = pickEntityId(client, "company", _("Pick company"));
            if (companyId <= 0) {
                return false;
            }
            const int appId = pickEntityId(client, "application", _("Pick application"));
            if (appId <= 0) {
                return false;
            }
            auto r = client.runJson({"company-app:assign", "--company_id=" + std::to_string(companyId),
                                     "--app_id=" + std::to_string(appId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Assigned."), mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({_("Unassign"), 'u', [](CliClient &client, const nlohmann::json &row) {
            const int companyId = jsonToInt(row, "company_id");
            const int appId = jsonToInt(row, "app_id");
            if (messageBox(_("Unassign this application from the company?"), mfYesNoCancel | mfConfirmation) !=
                cmYes) {
                return false;
            }
            auto r = client.runJson({"company-app:unassign", "--company_id=" + std::to_string(companyId),
                                     "--app_id=" + std::to_string(appId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Queue");
        e.cliEntity = "queue";
        e.deleteAction = "";
        e.supportsGet = false;
        e.columns = {{"id", "ID", 8}, {"runtemplate_id", "RT", 6}, {"scheduled", "Scheduled", 20},
                     {"schedule_type", "Type", 10}, {"executor", "Exec", 10}};
        e.listActions.push_back({_("Fix"), 'f', [](CliClient &client, const nlohmann::json &) {
            auto r = client.runJson({"queue:fix"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "queue:fix", r.data.dump(2));
            return true;
        }});
        e.listActions.push_back({_("Truncate"), 'T', [](CliClient &client, const nlohmann::json &) {
            if (messageBox(_("Truncate the entire job queue?"), mfYesNoCancel | mfConfirmation) != cmYes) {
                return false;
            }
            auto r = client.runJson({"queue:truncate"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Queue truncated."), mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Event Sources");
        e.cliEntity = "event-source";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", _("Name"), 22}, {"adapter", "Adapter", 14}, {"enabled", "En", 4}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("adapter", "Adapter", "", true),
            field("config", "Config JSON", "{}"),
            field("enabled", _("Enabled"), "1"),
        };
        e.editFields = e.createFields;
        e.rowActions.push_back({_("Test"), 't', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"event-source:test", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "event-source:test", r.data.dump(2));
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Event Rules");
        e.cliEntity = "event-rule";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", _("Name"), 22}, {"event_source_id", "Src", 5},
                     {"runtemplate_id", "RT", 5}, {"enabled", "En", 4}};
        e.createFields = {
            field("name", _("Name"), "", true),
            field("event_source_id", "Event source ID", "", true, false, "event-source"),
            field("runtemplate_id", "RunTemplate ID", "", true, false, "run-template"),
            field("enabled", _("Enabled"), "1"),
        };
        e.editFields = e.createFields;
        g_entities.push_back(std::move(e));
    }

    {
        EntityDef e;
        e.name = _("Deletion Requests");
        e.cliEntity = "user-erasure";
        e.deleteAction = "";
        e.supportsGet = false;
        e.columns = {{"id", "ID", 6}, {"user_id", "User", 6}, {"status", "Status", 12},
                     {"created_at", "Created", 20}};
        e.rowActions.push_back({_("Approve"), 'a', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-erasure:approve", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.rowActions.push_back({_("Reject"), 'r', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-erasure:reject", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.rowActions.push_back({_("Process"), 'p', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            if (messageBox(_("Process approved erasure (destructive)?"), mfYesNoCancel | mfConfirmation) !=
                cmYes) {
                return false;
            }
            auto r = client.runJson({"user-erasure:process", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.rowActions.push_back({_("Audit"), 'u', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-erasure:audit", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "Erasure audit", r.data.dump(2));
            return false;
        }});
        e.listActions.push_back({_("Create request"), 'n', [](CliClient &client, const nlohmann::json &) {
            const int userId = pickEntityId(client, "user", _("Pick user"));
            if (userId <= 0) {
                return false;
            }
            auto r = client.runJson({"user-erasure:create", "--user_id=" + std::to_string(userId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.listActions.push_back({_("Cleanup audits"), 'C', [](CliClient &client, const nlohmann::json &) {
            auto r = client.runJson({"user-erasure:cleanup"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(_("Old GDPR audit logs cleaned up."), mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }
}

} // namespace multiflexitui
