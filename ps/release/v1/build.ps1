Write-Output Start building release...

fyne-cross windows -arch=amd64 -image="fyneio/fyne-cross-images:1.1.3-windows" -developer="Zihuan Song" -dir . -name="Huan档案_x86_64.exe" ./src/cmd/v1

IF ($?) {
    Write-Output Success
} ELSE {
    Write-Output Fail
}