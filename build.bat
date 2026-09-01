@echo off
REM Build script for Windows - Derleme betiği

setlocal enabledelayedexpansion

set PROJECT_ROOT=%~dp0
set SRC_DIR=%PROJECT_ROOT%src
set BUILD_DIR=%PROJECT_ROOT%build
set CONFIG_DIR=%PROJECT_ROOT%config

echo.
echo 🔨 Randevu Tracker Derleniyor...
echo ================================
echo Kaynak: %SRC_DIR%
echo Build: %BUILD_DIR%
echo.

if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

cd /d "%SRC_DIR%"

echo 📦 Tüm Go dosyaları derleniyor...
go build -o "%BUILD_DIR%\randevu_tracker.exe" *.go

if %errorlevel% equ 0 (
    echo.
    echo ✅ Derleme başarılı!
    echo 📍 Binary konumu: %BUILD_DIR%\randevu_tracker.exe
    dir "%BUILD_DIR%\randevu_tracker.exe"
    echo.
    echo 💡 Başlatmak için:
    echo    %BUILD_DIR%\randevu_tracker.exe
) else (
    echo.
    echo ❌ Derleme başarısız!
    exit /b 1
)
