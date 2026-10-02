#pragma once

namespace multiflexitui {

constexpr unsigned short cmShowStatus = 1000;
constexpr unsigned short cmShowAbout = 1001;
constexpr unsigned short cmShowWebQr = 1002; // reserved (status-line click)

constexpr unsigned short cmOpenCompanies = 1100;
constexpr unsigned short cmOpenApplications = 1101;
constexpr unsigned short cmOpenRunTemplates = 1102;
constexpr unsigned short cmOpenJobs = 1103;
constexpr unsigned short cmOpenTasks = 1104;
constexpr unsigned short cmOpenCredentials = 1105;
constexpr unsigned short cmOpenTokens = 1106;
constexpr unsigned short cmOpenUsers = 1107;
constexpr unsigned short cmOpenArtifacts = 1108;
constexpr unsigned short cmOpenCredTypes = 1109;
constexpr unsigned short cmOpenCrPrototypes = 1110;
constexpr unsigned short cmOpenCompanyApps = 1111;
constexpr unsigned short cmOpenQueue = 1112;
constexpr unsigned short cmOpenEventSources = 1113;
constexpr unsigned short cmOpenEventRules = 1114;
constexpr unsigned short cmOpenConfFields = 1115;

constexpr unsigned short cmOpenActivationWizard = 1120;
constexpr unsigned short cmOpenCredentialWizard = 1121;
constexpr unsigned short cmOpenJobStream = 1122;
constexpr unsigned short cmOpenUserErasure = 1123;
constexpr unsigned short cmOpenEncryption = 1124;
constexpr unsigned short cmOpenPrune = 1125;
constexpr unsigned short cmOpenJobStatus = 1126;
constexpr unsigned short cmOpenTaskStatus = 1127;
constexpr unsigned short cmOpenQueueOverview = 1128;
constexpr unsigned short cmOpenTelemetryTest = 1129;
constexpr unsigned short cmOpenStaleRunTemplates = 1130;
constexpr unsigned short cmOpenImportExport = 1131;

constexpr unsigned short cmEntityRefresh = 1200;
constexpr unsigned short cmEntityNew = 1201;
constexpr unsigned short cmEntityEdit = 1202;
constexpr unsigned short cmEntityDelete = 1203;
constexpr unsigned short cmEntityDetail = 1204;
constexpr unsigned short cmEntityPrevPage = 1205;
constexpr unsigned short cmEntityNextPage = 1206;
constexpr unsigned short cmEntityAction = 1207; // infoPtr = action index
constexpr unsigned short cmEntitySave = 1208;
constexpr unsigned short cmEntityFormOk = 1209;

constexpr unsigned short cmMinimizeAll = 200;
constexpr unsigned short cmRestoreWindows = 201;
constexpr unsigned short cmCloseAll = 202;

constexpr unsigned short cmLangSystem = 1300;
constexpr unsigned short cmLangEnglish = 1301;
constexpr unsigned short cmLangCzech = 1302;

constexpr unsigned short cmWizardNext = 1400;
constexpr unsigned short cmWizardPrev = 1401;
constexpr unsigned short cmWizardFinish = 1402;

constexpr unsigned short cmStreamRefresh = 1410;
constexpr unsigned short cmStreamToggleFollow = 1411;

} // namespace multiflexitui
