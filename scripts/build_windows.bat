@echo off
REM =========================
REM Build script for Windows
REM =========================
REM Usage: build_windows.bat

echo 🔹 Building API for Windows...

REM Create bin folder if not exists
if not exist ..\bin mkdir ..\bin

REM Build binary
set GOOS=windows
set GOARCH=amd64
go build -o ..\bin\app-windows.exe ..\cmd\api

echo Build complete: ..\bin\app-windows.exe
pause
