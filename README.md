# MultiFlexi TUI
![Packaging: deb](https://img.shields.io/badge/packaging-.deb-red?logo=debian&logoColor=white)

A terminal user interface for [multiflexi-cli](https://github.com/VitexSoftware/multiflexi-cli), built with [Turbo Vision](https://github.com/magiblot/tvision) (tvision) — the same stack as [abraflexi-tui](https://github.com/VitexSoftware/abraflexi-tui).

## Features

- **Status dashboard** from `multiflexi-cli status`
- **Entity CRUD** for companies, applications, config fields, run templates, jobs, tasks, credentials, tokens, users, artifacts, credential types/prototypes, company-apps, queue, event sources/rules, and GDPR deletion requests
- **Special actions**: schedule / clone run templates, job stdout & stderr, token generate, user roles, queue fix/truncate, prototype sync, event-source test, app JSON export, …
- **Activation Wizard** and **Credential Wizard** (TUI approximations of the web wizards)
- **Live job streaming** via polled `job:get` (stdout/stderr follow-tail)
- **Admin helpers**: encryption status, prune, job/task/queue metrics, telemetry test, app JSON import/validate
- **Czech** UI strings via gettext (`po/cs.po`)

## Prerequisites

- C++17 toolchain, CMake ≥ 3.16
- `libtvision-dev`, `nlohmann-json3-dev`, `libncurses-dev`, `gettext`
- `multiflexi-cli` on `PATH` (runtime)

## Build

```bash
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build -j
sudo cmake --install build
```

## Usage

```bash
multiflexi-tui
multiflexi-tui --cli=/usr/bin/multiflexi-cli --envfile=/etc/multiflexi/multiflexi.env
```

Environment overrides: `MULTIFLEXI_TUI_CLI`, `MULTIFLEXI_TUI_LOCALEDIR`.

### Keyboard (global)

| Key | Action |
|-----|--------|
| `Alt+M` | MultiFlexi menu (Status, Exit) |
| `Alt+S` | Status |
| `Alt+C` / `Alt+A` / `Alt+R` / … | Entity menus |
| `F1` | About |
| `Alt+X` | Quit |
| `F6` / `Shift+F6` | Next / previous window |
| `F10` | Menu |

Entity lists: Enter opens detail, `n`/`Ins` create, `e` edit, `d` delete, `r` refresh.

## Debian package

```bash
dpkg-buildpackage -b -us -uc
sudo apt install ../multiflexi-tui_*.deb
```

## Architecture

All I/O goes through `CliClient` (`fork`+`execvp` of `multiflexi-cli … --format=json`). UI widgets follow abraflexi-tui patterns (`AppShell`, `SimpleListViewer`, `WindowColors`, …). Entity screens are driven by `EntityRegistry` definitions in `source/EntityRegistry.cpp`.

## License

MIT — Vitex Software
