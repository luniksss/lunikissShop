@echo off
setlocal enabledelayedexpansion

echo Lunikiss Shop - Migration Creator
echo ================================

if "%1"=="" (
    set /p MIGRATION_NAME=Enter migration name:
) else (
    set MIGRATION_NAME=%1
)

if "!MIGRATION_NAME!"=="" (
    echo Error: Migration name cannot be empty!
    pause
    exit /b 1
)

REM Создаем папку для миграций если не существует
if not exist "migrations" mkdir migrations

REM Генерируем timestamp в формате YYYYMMDDHHMMSS
for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set datetime=%%I
set TIMESTAMP=%datetime:~0,14%

REM Создаем имя файла
set MIGRATION_FILE=migrations\!TIMESTAMP!_!MIGRATION_NAME!.sql

echo Creating migration: !MIGRATION_FILE!

REM Создаем файл миграции с шаблоном
(
echo -- Migration: !MIGRATION_NAME!
echo -- Created: %date% %time%
echo -- Author: %USERNAME%
echo.
echo -- UP Migration
echo.
echo -- DOWN Migration (rollback)
echo.
) > "!MIGRATION_FILE!"

echo.
echo Migration file created successfully: !MIGRATION_FILE!
echo.
echo Edit the file to add your SQL statements.
echo.

REM Открываем файл в редакторе по умолчанию
notepad "!MIGRATION_FILE!"

endlocal