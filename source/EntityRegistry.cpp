#include "multiflexitui/TV.h"
#include "multiflexitui/EntityRegistry.h"
#include "multiflexitui/EntityDetailView.h"
#include "multiflexitui/i18n.h"

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

} // namespace

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

    // --- Companies ---
    {
        EntityDef e;
        e.name = "Companies";
        e.cliEntity = "company";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", "Name", 24}, {"slug", "Slug", 14}, {"ic", "IC", 12}, {"enabled", "En", 4}};
        e.createFields = {
            {"name", "Name", "", true},
            {"slug", "Slug", "", true},
            {"ic", "IC", "", false},
            {"email", "Email", "", false},
            {"enabled", "Enabled", "1", false},
        };
        e.editFields = e.createFields;
        g_entities.push_back(std::move(e));
    }

    // --- Applications ---
    {
        EntityDef e;
        e.name = "Applications";
        e.cliEntity = "application";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", "Name", 22}, {"code", "Code", 12}, {"version", "Ver", 8}, {"enabled", "En", 4}};
        e.createFields = {
            {"name", "Name", "", true},
            {"executable", "Executable", "", true},
            {"code", "Code", "", false},
            {"description", "Description", "", false},
            {"homepage", "Homepage", "", false},
            {"ociimage", "OCI Image", "", false},
            {"version", "Version", "", false},
            {"enabled", "Enabled", "1", false},
        };
        e.editFields = e.createFields;
        e.rowActions.push_back({"Config fields", 'c', [](CliClient &client, const nlohmann::json &row) {
            const int appId = jsonToInt(row, "id");
            auto r = client.runJson({"conffield:list", "--app_id=" + std::to_string(appId), "--limit=200", "--offset=0"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            std::ostringstream body;
            if (r.data.is_array()) {
                for (const auto &item : r.data) {
                    body << jsonToString(item.contains("id") ? item["id"] : nlohmann::json()) << "  "
                         << jsonToString(item.contains("keyname") ? item["keyname"] : nlohmann::json()) << "  "
                         << jsonToString(item.contains("type") ? item["type"] : nlohmann::json()) << "\n";
                }
            } else {
                body << r.data.dump(2);
            }
            showText(client, "Config fields for app " + std::to_string(appId), body.str());
            return false;
        }});
        e.rowActions.push_back({"Export JSON", 'x', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            const std::string path = "/tmp/multiflexi-app-" + std::to_string(id) + ".json";
            auto r = client.runJson({"application:export-json", "--id=" + std::to_string(id), "--file=" + path});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(("Exported to " + path).c_str(), mfInformation | mfOKButton);
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    // --- ConfFields ---
    {
        EntityDef e;
        e.name = "Config Fields";
        e.cliEntity = "conffield";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"app_id", "App", 6}, {"keyname", "Key", 18}, {"type", "Type", 10}, {"required", "Req", 4}};
        e.createFields = {
            {"app_id", "App ID", "", true},
            {"keyname", "Key name", "", true},
            {"type", "Type", "string", true},
            {"description", "Description", "", false},
            {"defval", "Default", "", false},
            {"required", "Required", "0", false},
        };
        e.editFields = e.createFields;
        g_entities.push_back(std::move(e));
    }

    // --- RunTemplates ---
    {
        EntityDef e;
        e.name = "RunTemplates";
        e.cliEntity = "run-template";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", "Name", 22}, {"company_id", "Co", 5}, {"app_id", "App", 5},
                     {"interv", "Interval", 10}, {"active", "Act", 4}};
        e.createFields = {
            {"name", "Name", "", true},
            {"app_id", "App ID", "", true},
            {"company_id", "Company ID", "", true},
            {"interv", "Interval", "n", false},
            {"executor", "Executor", "Native", false},
            {"active", "Active", "0", false},
            {"cron", "Cron", "", false},
        };
        e.editFields = e.createFields;
        e.rowActions.push_back({"Schedule", 's', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"run-template:schedule", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox(("Scheduled run-template " + std::to_string(id)).c_str(), mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({"Clone", 'c', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"run-template:clone", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox("Clone created (disabled).", mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({"List credentials", 'l', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"run-template:list-credentials", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "Credentials for RT " + std::to_string(id), r.data.dump(2));
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    // --- Jobs ---
    {
        EntityDef e;
        e.name = "Jobs";
        e.cliEntity = "job";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 8}, {"runtemplate_id", "RT", 6}, {"company_id", "Co", 5},
                     {"exitcode", "Exit", 5}, {"begin", "Begin", 20}, {"executor", "Exec", 10}};
        e.createFields = {
            {"runtemplate_id", "RunTemplate ID", "", true},
            {"scheduled", "Scheduled", "now", true},
            {"executor", "Executor", "Native", false},
            {"schedule_type", "Schedule type", "adhoc", false},
        };
        e.editFields = {
            {"exitcode", "Exit code", "", false},
            {"executor", "Executor", "", false},
        };
        e.rowActions.push_back({"Stdout", 'o', [](CliClient &c, const nlohmann::json &row) {
            return showJobOutput(c, row, "stdout", "Stdout");
        }});
        e.rowActions.push_back({"Stderr", 'r', [](CliClient &c, const nlohmann::json &row) {
            return showJobOutput(c, row, "stderr", "Stderr");
        }});
        g_entities.push_back(std::move(e));
    }

    // --- Tasks ---
    {
        EntityDef e;
        e.name = "Tasks";
        e.cliEntity = "task";
        e.deleteAction = "";
        e.columns = {{"id", "ID", 8}, {"runtemplate_id", "RT", 6}, {"state", "State", 12},
                     {"window_start", "Start", 20}, {"deadline", "Deadline", 20}};
        g_entities.push_back(std::move(e));
    }

    // --- Credentials ---
    {
        EntityDef e;
        e.name = "Credentials";
        e.cliEntity = "credential";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", "Name", 24}, {"company_id", "Co", 5}, {"credential_type_id", "Type", 6}};
        e.createFields = {
            {"name", "Name", "", true},
            {"company_id", "Company ID", "", true},
            {"credential_type_id", "Credential type ID", "", true},
        };
        e.editFields = {{"name", "Name", "", true}};
        e.listActions.push_back({"Encrypt existing", 'E', [](CliClient &client, const nlohmann::json &) {
            auto r = client.runJson({"credential:encrypt-existing"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox("Encryption pass completed.", mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    // --- Tokens ---
    {
        EntityDef e;
        e.name = "Tokens";
        e.cliEntity = "token";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"user", "User", 16}, {"user_id", "UID", 5}, {"start", "Start", 20}, {"until", "Until", 20}};
        e.createFields = {
            {"user_id", "User ID", "", true},
            {"until", "Until", "", false},
        };
        e.editFields = {{"until", "Until", "", false}};
        e.rowActions.push_back({"Generate", 'g', [](CliClient &client, const nlohmann::json &row) {
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

    // --- Users ---
    {
        EntityDef e;
        e.name = "Users";
        e.cliEntity = "user";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"login", "Login", 16}, {"email", "Email", 24}, {"enabled", "En", 4}};
        e.createFields = {
            {"login", "Login", "", true},
            {"email", "Email", "", true},
            {"firstname", "First name", "", false},
            {"lastname", "Last name", "", false},
            {"password", "Password", "", true},
            {"enabled", "Enabled", "1", false},
        };
        e.editFields = {
            {"email", "Email", "", false},
            {"firstname", "First name", "", false},
            {"lastname", "Last name", "", false},
            {"enabled", "Enabled", "", false},
        };
        e.rowActions.push_back({"Roles", 'r', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-role:list", "--user_id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "Roles for user " + std::to_string(id), r.data.dump(2));
            return false;
        }});
        e.rowActions.push_back({"Make Admin", 'a', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-role:set", "--user_id=" + std::to_string(id), "--roles=admin"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox("Admin role set.", mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({"Assign company", 'c', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            char buf[32] = {};
            if (inputBox("Company ID", "Assign user to company", buf, sizeof(buf) - 1) != cmOK) {
                return false;
            }
            auto r = client.runJson({"user-company:assign", "--user_id=" + std::to_string(id),
                                     "--company_id=" + std::string(buf)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox("Assigned.", mfInformation | mfOKButton);
            return false;
        }});
        g_entities.push_back(std::move(e));
    }

    // --- Artifacts ---
    {
        EntityDef e;
        e.name = "Artifacts";
        e.cliEntity = "artifact";
        e.deleteAction = "";
        e.columns = {{"id", "ID", 6}, {"job_id", "Job", 8}, {"filename", "Filename", 28}, {"content_type", "Type", 18}};
        e.rowActions.push_back({"Save", 's', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            std::string filename = jsonToString(row.contains("filename") ? row["filename"] : nlohmann::json());
            if (filename.empty()) {
                filename = "artifact-" + std::to_string(id) + ".bin";
            }
            const std::string path = "/tmp/" + filename;
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

    // --- CredTypes ---
    {
        EntityDef e;
        e.name = "Credential Types";
        e.cliEntity = "credential-type";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", "Name", 22}, {"prototype", "Prototype", 16}, {"company_id", "Co", 5}};
        e.createFields = {
            {"name", "Name", "", true},
            {"company_id", "Company ID", "", true},
            {"prototype", "Prototype code", "", true},
        };
        e.editFields = {{"name", "Name", "", true}, {"url", "URL", "", false}};
        g_entities.push_back(std::move(e));
    }

    // --- CrPrototypes ---
    {
        EntityDef e;
        e.name = "Cred Prototypes";
        e.cliEntity = "credential-prototype";
        e.deleteAction = "delete";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"code", "Code", 16}, {"name", "Name", 24}, {"version", "Ver", 8}};
        e.createFields = {
            {"name", "Name", "", true},
            {"code", "Code", "", true},
            {"description", "Description", "", false},
        };
        e.editFields = e.createFields;
        e.listActions.push_back({"Sync all", 'S', [](CliClient &client, const nlohmann::json &) {
            auto r = client.runJson({"credential-prototype:sync"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox("Prototypes synchronized.", mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    // --- CompanyApps ---
    {
        EntityDef e;
        e.name = "Company Apps";
        e.cliEntity = "company-app";
        e.deleteAction = "";
        e.columns = {{"id", "ID", 6}, {"company_id", "Co", 5}, {"company_name", "Company", 18},
                     {"app_id", "App", 5}, {"app_name", "Application", 22}};
        e.listActions.push_back({"Assign", 'a', [](CliClient &client, const nlohmann::json &) {
            char company[32] = {};
            char app[32] = {};
            if (inputBox("Company ID", "Assign application", company, sizeof(company) - 1) != cmOK) {
                return false;
            }
            if (inputBox("Application ID", "Assign application", app, sizeof(app) - 1) != cmOK) {
                return false;
            }
            auto r = client.runJson({"company-app:assign", "--company_id=" + std::string(company),
                                     "--app_id=" + std::string(app)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox("Assigned.", mfInformation | mfOKButton);
            return true;
        }});
        e.rowActions.push_back({"Unassign", 'u', [](CliClient &client, const nlohmann::json &row) {
            const int companyId = jsonToInt(row, "company_id");
            const int appId = jsonToInt(row, "app_id");
            if (messageBox("Unassign this application from the company?", mfYesNoCancel | mfConfirmation) != cmYes) {
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

    // --- Queue ---
    {
        EntityDef e;
        e.name = "Queue";
        e.cliEntity = "queue";
        e.deleteAction = "";
        e.columns = {{"id", "ID", 8}, {"runtemplate_id", "RT", 6}, {"scheduled", "Scheduled", 20},
                     {"schedule_type", "Type", 10}, {"executor", "Exec", 10}};
        e.listActions.push_back({"Fix", 'f', [](CliClient &client, const nlohmann::json &) {
            auto r = client.runJson({"queue:fix"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "queue:fix", r.data.dump(2));
            return true;
        }});
        e.listActions.push_back({"Truncate", 'T', [](CliClient &client, const nlohmann::json &) {
            if (messageBox("Truncate the entire job queue?", mfYesNoCancel | mfConfirmation) != cmYes) {
                return false;
            }
            auto r = client.runJson({"queue:truncate"});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            messageBox("Queue truncated.", mfInformation | mfOKButton);
            return true;
        }});
        g_entities.push_back(std::move(e));
    }

    // --- EventSources ---
    {
        EntityDef e;
        e.name = "Event Sources";
        e.cliEntity = "event-source";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", "Name", 22}, {"adapter", "Adapter", 14}, {"enabled", "En", 4}};
        e.createFields = {
            {"name", "Name", "", true},
            {"adapter", "Adapter", "", true},
            {"config", "Config JSON", "{}", false},
            {"enabled", "Enabled", "1", false},
        };
        e.editFields = e.createFields;
        e.rowActions.push_back({"Test", 't', [](CliClient &client, const nlohmann::json &row) {
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

    // --- EventRules ---
    {
        EntityDef e;
        e.name = "Event Rules";
        e.cliEntity = "event-rule";
        e.deleteAction = "remove";
        e.canCreate = true;
        e.canEdit = true;
        e.columns = {{"id", "ID", 6}, {"name", "Name", 22}, {"event_source_id", "Src", 5},
                     {"runtemplate_id", "RT", 5}, {"enabled", "En", 4}};
        e.createFields = {
            {"name", "Name", "", true},
            {"event_source_id", "Event source ID", "", true},
            {"runtemplate_id", "RunTemplate ID", "", true},
            {"enabled", "Enabled", "1", false},
        };
        e.editFields = e.createFields;
        g_entities.push_back(std::move(e));
    }

    // --- User erasure (GDPR) ---
    {
        EntityDef e;
        e.name = "Deletion Requests";
        e.cliEntity = "user-erasure";
        e.deleteAction = "";
        e.columns = {{"id", "ID", 6}, {"user_id", "User", 6}, {"status", "Status", 12}, {"created_at", "Created", 20}};
        e.rowActions.push_back({"Approve", 'a', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-erasure:approve", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.rowActions.push_back({"Reject", 'r', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-erasure:reject", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.rowActions.push_back({"Process", 'p', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            if (messageBox("Process approved erasure (destructive)?", mfYesNoCancel | mfConfirmation) != cmYes) {
                return false;
            }
            auto r = client.runJson({"user-erasure:process", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        e.rowActions.push_back({"Audit", 'u', [](CliClient &client, const nlohmann::json &row) {
            const int id = jsonToInt(row, "id");
            auto r = client.runJson({"user-erasure:audit", "--id=" + std::to_string(id)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            showText(client, "Erasure audit", r.data.dump(2));
            return false;
        }});
        e.listActions.push_back({"Create request", 'n', [](CliClient &client, const nlohmann::json &) {
            char userId[32] = {};
            if (inputBox("User ID", "Create erasure request", userId, sizeof(userId) - 1) != cmOK) {
                return false;
            }
            auto r = client.runJson({"user-erasure:create", "--user_id=" + std::string(userId)});
            if (!r.ok) {
                messageBox(r.errorMessage.c_str(), mfError | mfOKButton);
                return false;
            }
            return true;
        }});
        g_entities.push_back(std::move(e));
    }
}

} // namespace multiflexitui
