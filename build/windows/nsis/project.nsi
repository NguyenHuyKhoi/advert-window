Unicode true
RequestExecutionLevel user ; Ép buộc quyền User để ghi vào HKCU không cần Admin

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

; Chỉ cài vào ProductName, bỏ CompanyName
InstallDir "$LOCALAPPDATA\${INFO_PRODUCTNAME}"
ShowInstDetails show

Function .onInit
    !insertmacro wails.checkArchitecture
FunctionEnd

Section
    SetShellContext current

    ; 🛠️ FIX 1: Thêm dấu ngoặc kép bao quanh taskkill để tránh lỗi lệnh
    nsExec::ExecToLog 'taskkill /F /IM "advert.exe" /T'
    Sleep 1000

    SetOutPath $INSTDIR
    !insertmacro wails.files

    ; Shortcut
    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"
    CreateShortcut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"

    ; 🔑 FIX 2: AUTO START 
    ; Ghi trực tiếp vào HKCU (Current User). Dấu ngoặc kép bao quanh path là bắt buộc.
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Run" \
      "ForlifeMediaPlayer" '"$INSTDIR\advert.exe"'

    ; 🚀 FIX 3: Khởi động lại App
    Exec '"$INSTDIR\advert.exe"'

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    SetShellContext current

    ; Xóa Registry khi gỡ
    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "ForlifeMediaPlayer"

    RMDir /r $INSTDIR
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.deleteUninstaller
SectionEnd