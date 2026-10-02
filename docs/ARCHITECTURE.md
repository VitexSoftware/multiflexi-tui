# multiflexi-tui architecture (C++ / Turbo Vision)

## Overview

`multiflexi-tui` is a Turbo Vision application that shells out to `multiflexi-cli`
with `--format=json`. It does not talk to the database or REST API directly.

## Layers

1. **CliClient / ProcessRunner** — `fork` + `execvp`, never a shell; parses JSON with nlohmann::json.
2. **EntityRegistry** — describes each CLI entity (columns, create/edit fields, row/list actions).
3. **EntityListView / EntityDetailView / EntityForm** — generic CRUD windows.
4. **Wizards** — ActivationWizard, CredentialWizard multi-step dialogs.
5. **JobStreamView** — polls `job:get` from `AppShell::idle()` for live stdout/stderr.
6. **AppShell** — menus aligned with the web MainMenu groups; window tile/cascade.

## Packaging

CMake + `libtvision-dev` + `nlohmann-json3-dev` + gettext. Debian package depends on `multiflexi-cli`.
