Unicode true
RequestExecutionLevel user ; Cần thiết để ghi vào HKCU không cần quyền Admin

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

; Cài đặt trực tiếp vào Local AppData
InstallDir "$LOCALAPPDATA\${INFO_PRODUCTNAME}"
ShowInstDetails show

Function .onInit
    !insertmacro wails.checkArchitecture
FunctionEnd

Section
    SetShellContext current

    ; Dừng app đang chạy để ghi đè file
    nsExec::ExecToLog 'taskkill /F /IM "advert.exe" /T'
    Sleep 1000

    SetOutPath $INSTDIR
    !insertmacro wails.files

    ; Shortcut
    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"
    CreateShortcut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"

    ; 🔑 AUTO START - ĐÃ ĐỔI KEY THÀNH ForlifeAdvert
    ; Dùng HKCU trực tiếp, đường dẫn bọc trong dấu ngoặc kép '"..."'
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Run" \
      "ForlifeAdvert" '"$INSTDIR\advert.exe"'

    ; 🚀 Chạy app ngay sau khi cài
    Exec '"$INSTDIR\advert.exe"'

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    SetShellContext current

    ; Xóa đúng Key ForlifeAdvert khi gỡ
    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "ForlifeAdvert"

    RMDir /r $INSTDIR
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.deleteUninstaller
SectionEnd