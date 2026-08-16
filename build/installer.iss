
#define AppName "ExNES"
#define AppPublisher "Emptiedfull"
#define AppURL        "https://github.com/Emptiedfull/ExNES"
#define AppExeName    "exnes.exe"

#define StageDir "..\dist\ExNES-windows-x86_64"

#ifndef AppVersion
    #define AppVersion "0.0.0"
#endif

[Setup]
AppId={{7C3F1A62-5E4D-4B18-9A77-2F0B6D3C8E51}}
AppName={#AppName}
AppVersion={#AppVersion}
AppVerName={#AppName} {#AppVersion}
AppPublisher={#AppPublisher}
AppPublisherURL={#AppURL}
AppSupportURL={#AppURL}/issues
AppUpdatesURL={#AppURL}/releases

PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog

DefaultDirName={autopf}\{#AppName}
DefaultGroupName={#AppName}
DisableProgramGroupPage=yes
UninstallDisplayIcon={app}\{#AppExeName}

OutputDir=..\dist
OutputBaseFilename=ExNES-{#AppVersion}-windows-x64-setup
Compression=lzma2/max
SolidCompression=yes
WizardStyle=modern

ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

[Languages]
Name:"English"; MessageFile: "compiler:Default.isl" 

[Tasks]
Name: "desktopicon"; Description: "Create shortcut"; GroupDescription: "Shortcuts:"; Flags: unchecked
Name: "associate"; Description: "Link .nes with ExNES";  GroupDescription: "File associations:"

[Icons]
Name: "{autoprograms}\{#AppName}"; Filename: "{app}\{#AppExeName}"
Name: "{autodesktop}\{#AppName}"; Filename: "{app}\{#AppExeName}"; Tasks: desktopicon

[Files]
Source: "{#StageDir}\{#AppExeName}";DestDir: "{app}"; Flags: ignoreversion
Source: "{#StageDir}\*.dll";DestDir: "app"; Flags: ignoreversion

[Registry]
Root: HKA; Subkey: "Software\Classes\.nes\OpenWithProgids"; ValueType: string; ValueName: "ExNES.rom"; ValueData: ""; Flags:uninsdeletevalue; Tasks: associate
Root: HKA; Subkey: "Software\Classes\ExNES.rom"; ValueType: string; ValueName: ""; ValueData: "NES ROM"; Flags: uninsdeletevalue; Tasks: associate
Root: HKA; Subkey: "Software\Classes\ExNES.rom\DefaultIcon"; \ ValueType: string; ValueName:"";ValueData: "{app}\{#AppExeName},0" Tasks: associate
Root: HKA; Subkey: "Software\Classes\ExNES.rom\shell\open\command"; ValueType: string; ValueName:""; ValueData: """{app}\{#AppExeName}"" ""%1"""; Tasks: associate

[Run]
Filename: "{app}\{#AppExeName}"; Description: "Launch {#AppName}"; Flags: nowait postinstall skipifsilient

