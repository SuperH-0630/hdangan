Write-Output Start building...

fyne-cross windows -arch=386 -debug -image="fyneio/fyne-cross-images:1.1.3-windows" -developer="Zihuan Song" -dir . -name="Huan档案_x368_debug.exe" ./src/cmd/v1

IF ($?) {
    Write-Output Success
} ELSE {
    Write-Output Fail
}