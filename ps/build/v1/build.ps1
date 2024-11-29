Write-Output Start building...

fyne-cross windows -arch=amd64 -debug -image "fyneio/fyne-cross-images:1.1.3-windows" -developer "Zihuan Song" -name "Huan档案_x86_64.exe" ./src/cmd/v1

IF ($?) {
    Write-Output Success
} ELSE {
    Write-Output Fail
}