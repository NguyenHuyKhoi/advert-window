Unicode true
RequestExecutionLevel user

!include "wails_tools.nsh"
!include "MUI.nsh"

ManifestDPIAware true

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe"

InstallDir "$LOCALAPPDATA\${INFO_PRODUCTNAME}"
ShowInstDetails show

Function .onInit
    !insertmacro wails.checkArchitecture
FunctionEnd

Section
    SetShellContext current

    ; Kill app cũ
    nsExec::ExecToLog 'taskkill /F /IM "advert.exe" /T'
    Sleep 1000

    SetOutPath $INSTDIR
    !insertmacro wails.files

    ; Shortcut thường
    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"
    CreateShortcut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"

    ; ✅ AUTO START — CÁCH CHẮC CHẮN NHẤT
    CreateShortcut "$SMSTARTUP\ForlifeAdvert.lnk" "$INSTDIR\advert.exe"

    ; 🚀 AUTO RUN NGAY SAU KHI CÀI
    Exec '"$INSTDIR\advert.exe"'

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    SetShellContext current

    Delete "$SMSTARTUP\ForlifeAdvert.lnk"

    RMDir /r $INSTDIR
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.deleteUninstaller
SectionEnd
