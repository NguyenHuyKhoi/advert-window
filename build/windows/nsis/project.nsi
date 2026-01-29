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
    nsExec::ExecToLog 'taskkill /IM advert.exe /F'
    Sleep 800

    SetOutPath $INSTDIR
    !insertmacro wails.files

    ; ✅ AUTO START — GIỜ SẼ GHI ĐÚNG USER HKCU
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Run" \
      "ForlifeMediaPlayer" '"$INSTDIR\advert.exe"'

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"
    CreateShortcut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"

    Exec '"$INSTDIR\advert.exe"'

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    SetShellContext current

    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "ForlifeMediaPlayer"

    RMDir /r $INSTDIR
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.deleteUninstaller
SectionEnd
