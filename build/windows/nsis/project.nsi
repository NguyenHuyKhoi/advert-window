Unicode true
!include "wails_tools.nsh"

VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

ManifestDPIAware true

!include "MUI.nsh"

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

; Bỏ CompanyName như yêu cầu
InstallDir "$LOCALAPPDATA\${INFO_PRODUCTNAME}"
ShowInstDetails show

Function .onInit
    !insertmacro wails.checkArchitecture
FunctionEnd

Section
    ; Ép buộc ngữ cảnh là User để có quyền ghi HKCU
    SetShellContext current

    ; 🔥 Kill app cũ
    nsExec::ExecToLog 'taskkill /IM advert.exe /F'
    Sleep 800

    SetOutPath $INSTDIR
    !insertmacro wails.files

    ; 🔑 GHI REGISTRY TRƯỚC KHI CHẠY APP
    ; Dùng HKCU thay vì SHCTX để đảm bảo ghi đúng vào hình bạn chụp
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Run" \
      "ForlifeMediaPlayer" '"$INSTDIR\advert.exe"'

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"

    ; 🚀 Chạy app sau cùng
    Exec '"$INSTDIR\advert.exe"'

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols
    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    SetShellContext current
    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "ForlifeMediaPlayer"

    RMDir /r $INSTDIR
    RMDir /r "$AppData\${INFO_PRODUCTNAME}"

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols
    !insertmacro wails.deleteUninstaller
SectionEnd