@echo off
REM Reseed database - Clears all data and reloads sample data
REM WARNING: This will delete ALL existing data!

echo.
echo ⚠️  WARNING: This will DELETE ALL DATA in the database!
echo.
set /p CONFIRM="Are you sure? Type 'yes' to continue: "

if not "%CONFIRM%"=="yes" (
    echo Cancelled.
    exit /b 0
)

echo.
echo Clearing and reseeding database...

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

REM Drop and recreate database
echo Dropping database %DB_NAME%...
dropdb -U %DB_USER% %DB_NAME% 2>nul
createdb -U %DB_USER% %DB_NAME%

REM Apply schema
echo Applying schema...
psql -U %DB_USER% -d %DB_NAME% -f ..\migrations\001_initial_schema.up.sql

REM Apply seed data
echo Loading sample data...
psql -U %DB_USER% -d %DB_NAME% -f ..\migrations\002_seed_data.up.sql

echo.
echo ✅ Database reseeded successfully!
echo.
echo Sample users (username = password):
echo   - rgodwin, csamad, eknutson, jminihan (admins)
echo   - test, asd1, david (users)
echo.
echo 🎉 Ready to go!
