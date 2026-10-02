#pragma once

#include "multiflexitui/CliClient.h"

#include <functional>
#include <map>
#include <string>
#include <vector>

#include <nlohmann/json.hpp>

namespace multiflexitui {

struct ColumnDef {
    std::string key;   // JSON object key
    std::string title; // column header
    int width = 12;
};

struct FieldDef {
    std::string key;         // JSON key / CLI option name without --
    std::string label;       // form label
    std::string placeholder; // empty-field hint
    bool required = false;
    bool secret = false;     // password-style (still plain TInputLine)
    bool readOnly = false;   // shown in detail only
};

struct EntityAction {
    std::string label;
    char hotkey = 0; // unused in TV menus but kept for docs parity
    // Returns true if the list should refresh afterwards.
    std::function<bool(CliClient &client, const nlohmann::json &row)> run;
};

struct EntityDef {
    std::string name;         // display name
    std::string cliEntity;    // e.g. "company", "run-template"
    std::string deleteAction; // "delete" or "remove"; empty = no delete
    std::string idKey = "id";
    int pageSize = 40;
    bool canCreate = false;
    bool canEdit = false;
    std::vector<ColumnDef> columns;
    std::vector<FieldDef> createFields;
    std::vector<FieldDef> editFields;
    std::vector<EntityAction> rowActions;
    std::vector<EntityAction> listActions;
};

const std::vector<EntityDef> &allEntities();
const EntityDef *findEntity(const std::string &cliEntity);
void registerBuiltinEntities();

// Helpers shared by views.
std::string jsonToString(const nlohmann::json &v);
int jsonToInt(const nlohmann::json &obj, const std::string &key, int fallback = 0);
std::string formatRow(const EntityDef &def, const nlohmann::json &obj);
std::string formatHeader(const EntityDef &def);
std::vector<std::string> buildCliArgs(const std::vector<FieldDef> &fields,
                                      const std::map<std::string, std::string> &values,
                                      int id = -1);

} // namespace multiflexitui
