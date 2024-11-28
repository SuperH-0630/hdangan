Write-Output Start building...
fyne package --os windows --src ./src/cmd/v1
Move-Item ./src/cmd/v1/Huan档案.exe ./target/build/Huan档案.exe
Write-Output Finish