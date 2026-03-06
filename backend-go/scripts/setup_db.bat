@echo off
REM Database setup script for Noddit (Windows)
REM Run this to create the database and apply migrations

echo Setting up database...

REM Set default values
set DB_USER=postgres
set DB_NAME=userdb

REM Load .env if exists
if exist ..\.env (
    for /f "tokens=1,2 delims==" %%a in (..\.env) do (
        if "%%a"=="DB_USER" set DB_USER=%%b
        if "%%a"=="DB_NAME" set DB_NAME=%%b
    )
)

echo Database: %DB_NAME%
echo User: %DB_USER%

REM Create database
createdb -U %DB_USER% %DB_NAME% 2>nul || echo Database %DB_NAME% already exists

REM Apply migrations
echo Applying migrations...
psql -U %DB_USER% -d %DB_NAME% -f ..\migrations\001_initial_schema.up.sql

echo.
echo ✅ Database setup complete!
echo.
echo To load sample data from the old database, run:
echo   psql -U %DB_USER% -d %DB_NAME% -f ..\..\backend\database\dbexport.pgsql
