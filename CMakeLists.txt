cmake_minimum_required(VERSION 3.16)
project(multiflexi-tui LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release)
endif()

find_package(PkgConfig REQUIRED)
pkg_check_modules(TVISION REQUIRED IMPORTED_TARGET tvision)
add_library(tvision::tvision ALIAS PkgConfig::TVISION)

find_package(nlohmann_json REQUIRED)
find_package(Threads REQUIRED)
find_package(Intl REQUIRED)

find_program(MSGFMT_EXECUTABLE msgfmt REQUIRED)
find_program(MSGMERGE_EXECUTABLE msgmerge)
find_program(XGETTEXT_EXECUTABLE xgettext)

set(MULTIFLEXI_TUI_LOCALEDIR "${CMAKE_INSTALL_PREFIX}/share/locale" CACHE PATH
    "Directory .mo catalogs are installed to and loaded from at runtime")
set(MULTIFLEXI_TUI_LANGUAGES cs)
set(MULTIFLEXI_TUI_MO_FILES "")

foreach(lang ${MULTIFLEXI_TUI_LANGUAGES})
    set(po_file "${CMAKE_SOURCE_DIR}/po/${lang}.po")
    set(mo_dir "${CMAKE_BINARY_DIR}/locale/${lang}/LC_MESSAGES")
    set(mo_file "${mo_dir}/multiflexi-tui.mo")
    file(MAKE_DIRECTORY "${mo_dir}")
    add_custom_command(
        OUTPUT "${mo_file}"
        COMMAND "${MSGFMT_EXECUTABLE}" --check -o "${mo_file}" "${po_file}"
        DEPENDS "${po_file}"
        COMMENT "Compiling po/${lang}.po -> ${lang}.mo"
    )
    list(APPEND MULTIFLEXI_TUI_MO_FILES "${mo_file}")
    install(FILES "${mo_file}"
        DESTINATION "${CMAKE_INSTALL_PREFIX}/share/locale/${lang}/LC_MESSAGES"
        RENAME multiflexi-tui.mo)
endforeach()

add_custom_target(translations ALL DEPENDS ${MULTIFLEXI_TUI_MO_FILES})

add_executable(multiflexi-tui
    source/main.cpp
    source/CliClient.cpp
    source/i18n.cpp
    source/SimpleListViewer.cpp
    source/WindowColors.cpp
    source/AppShell.cpp
    source/AppStatusLine.cpp
    source/AppButton.cpp
    source/StatusView.cpp
    source/AboutView.cpp
    source/EntityRegistry.cpp
    source/EntityListView.cpp
    source/EntityDetailView.cpp
    source/Wizards.cpp
    source/JobStreamView.cpp
)

target_include_directories(multiflexi-tui PRIVATE include)
target_link_libraries(multiflexi-tui PRIVATE
    tvision::tvision
    nlohmann_json::nlohmann_json
    Threads::Threads
    ${Intl_LIBRARIES}
)
target_include_directories(multiflexi-tui PRIVATE ${Intl_INCLUDE_DIRS})
target_compile_definitions(multiflexi-tui PRIVATE
    MULTIFLEXI_TUI_LOCALEDIR="${MULTIFLEXI_TUI_LOCALEDIR}"
    MULTIFLEXI_TUI_GETTEXT_DOMAIN="multiflexi-tui"
)
add_dependencies(multiflexi-tui translations)

install(TARGETS multiflexi-tui RUNTIME DESTINATION bin)

add_executable(multiflexi-tui-cliclient-check
    tools/cliclient_check.cpp
    source/CliClient.cpp
)
target_include_directories(multiflexi-tui-cliclient-check PRIVATE include)
target_link_libraries(multiflexi-tui-cliclient-check PRIVATE
    nlohmann_json::nlohmann_json
    Threads::Threads
)
