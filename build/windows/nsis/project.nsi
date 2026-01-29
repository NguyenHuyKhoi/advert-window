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

; 🔥 FIX: Bỏ INFO_COMPANYNAME, cài thẳng vào thư mục App
InstallDir "$LOCALAPPDATA\${INFO_PRODUCTNAME}"

ShowInstDetails show

Function .onInit
    !insertmacro wails.checkArchitecture
FunctionEnd

Section
    !insertmacro wails.setShellContext

    ; 🔥 Kill running app (advert.exe) để không bị lock file
    nsExec::ExecToLog 'taskkill /IM advert.exe /F'

    Sleep 800

    SetOutPath $INSTDIR
    !insertmacro wails.files

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\advert.exe"

    ; 🔑 FIX AUTO-START: 
    ; 1. Ghi vào HKCU (Current User)
    ; 2. Sử dụng đúng tên file advert.exe
    ; 3. Bọc dấu ngoặc kép '"..."' để xử lý khoảng trắng trong đường dẫn Windows
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Run" \
      "ForlifeMediaPlayer" '"$INSTDIR\advert.exe"'

    ; 🚀 Relaunch ngay lập tức sau khi cài/update
    Exec '"$INSTDIR\advert.exe"'

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    !insertmacro wails.setShellContext

    ; Xóa Registry khởi động khi gỡ app
    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "ForlifeMediaPlayer"

    RMDir /r "$AppData\${INFO_PRODUCTNAME}"
    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols

    !insertmacro wails.deleteUninstaller
SectionEnd