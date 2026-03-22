@echo off
REM QuizMe Admin User Seeding Script for Windows
REM Color codes support (requires Windows 10+)

setlocal enabledelayedexpansion

REM Default values
set USERNAME=admin
set EMAIL=admin@quizme.com
set PASSWORD=admin123
set FULLNAME=Administrator
set UPDATE_FLAG=

REM Parse arguments
:parse_args
if "%1"=="" goto start_script
if "%1"=="-u" (
  set USERNAME=%2
  shift
  shift
  goto parse_args
)
if "%1"=="--username" (
  set USERNAME=%2
  shift
  shift
  goto parse_args
)
if "%1"=="-e" (
  set EMAIL=%2
  shift
  shift
  goto parse_args
)
if "%1"=="--email" (
  set EMAIL=%2
  shift
  shift
  goto parse_args
)
if "%1"=="-p" (
  set PASSWORD=%2
  shift
  shift
  goto parse_args
)
if "%1"=="--password" (
  set PASSWORD=%2
  shift
  shift
  goto parse_args
)
if "%1"=="-n" (
  set FULLNAME=%2
  shift
  shift
  goto parse_args
)
if "%1"=="--fullname" (
  set FULLNAME=%2
  shift
  shift
  goto parse_args
)
if "%1"=="--update" (
  set UPDATE_FLAG=-update
  shift
  goto parse_args
)
if "%1"=="-h" (
  goto show_help
)
if "%1"=="--help" (
  goto show_help
)
shift
goto parse_args

:show_help
echo.
echo QuizMe Admin User Seeding Script
echo.
echo Usage: create-admin.bat [options]
echo.
echo Options:
echo   -u, --username USERNAME    Admin username (default: admin)
echo   -e, --email EMAIL          Admin email (default: admin@quizme.com)
echo   -p, --password PASSWORD    Admin password (default: admin123)
echo   -n, --fullname FULLNAME    Admin full name (default: Administrator)
echo   --update                   Update existing admin user if exists
echo   -h, --help                 Show this help message
echo.
echo Examples:
echo   create-admin.bat
echo   create-admin.bat -u superadmin -e super@quizme.com -p Pass123
echo   create-admin.bat -u admin -p NewPassword123 --update
echo.
exit /b 0

:start_script
cls
echo.
echo ========================================
echo   QuizMe Admin User Seeding
echo ========================================
echo.

echo Configuration:
echo   Username:   %USERNAME%
echo   Email:      %EMAIL%
echo   Password:   **** (hidden)
echo   Full Name:  %FULLNAME%
if defined UPDATE_FLAG (
  echo   Mode:       UPDATE (if exists)
) else (
  echo   Mode:       CREATE (new)
)
echo.

REM Get the directory where this script is located
cd /d "%~dp0.." || (
  echo Error: Failed to change directory
  exit /b 1
)

echo Running seed command...
echo.

REM Build and run command
set CMD=go run cmd/seed/main.go -username=%USERNAME% -email=%EMAIL% -password=%PASSWORD% -fullname=%FULLNAME%
if defined UPDATE_FLAG (
  set CMD=!CMD! %UPDATE_FLAG%
)

%CMD%

if %errorlevel% equ 0 (
  echo.
  echo Admin user setup completed successfully!
  echo.
  echo You can now log in with:
  echo   Username/Email: %USERNAME% (or %EMAIL%)
  echo   Password:       ****
  echo.
) else (
  echo.
  echo Admin user setup failed (Exit code: %errorlevel%)
  exit /b %errorlevel%
)

endlocal
