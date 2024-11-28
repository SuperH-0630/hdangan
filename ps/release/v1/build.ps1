Write-Output Start building release...
IF (Test-Path ./target/release/Huan档案.exe) {
    Remove-Item ./target/release/Huan档案.exe
}
fyne package --os windows --src ./src/cmd/v1 --release
Move-Item ./src/cmd/v1/Huan档案.exe ./target/release/Huan档案.exe
Write-Output Finish