# multiflexi-tui architecture (C++ / Turbo Vision)

## Overview

`multiflexi-tui` is a Turbo Vision application that shells out to `multiflexi-cli`
with `--format=json`. It does not talk to the database or REST API directly.

## Layers

1. **CliClient / ProcessRunner / ProcessHandle** — `fork` + `execvp`, never a shell; parses JSON with nlohmann::json. Sync runner for modal actions; non-blocking `ProcessHandle` for the async queue.
2. **AsyncCliQueue** — one-at-a-time CLI jobs pumped from `AppShell::idle()`; used by entity list refresh and job streams.
3. **EntityRegistry** — describes each CLI entity (columns, create/edit fields, row/list actions, `supportsGet`, `supportsCompanyScope`).
4. **EntityListView / EntityDetailView / EntityForm / EntityPicker** — generic CRUD windows, filter, actions menu, relation pickers.
5. **Wizards** — ActivationWizard, CredentialWizard multi-step dialogs.
6. **JobStreamView** — polls `job:get` via the async queue; multiple stream tick callbacks.
7. **AppShell** — menus, active company badge, Busy indicator, window tile/cascade.

## Packaging

CMake + `libtvision-dev` + `nlohmann-json3-dev` + gettext. Debian package depends on `multiflexi-cli`.
