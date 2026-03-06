@echo off
REM Seed database with sample data (Windows)
REM Run this AFTER setup_db.bat to populate with example posts

echo Seeding database with sample data...

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

REM Apply seed migration
psql -U %DB_USER% -d %DB_NAME% -f ..\migrations\002_seed_data.up.sql

echo.
echo ✅ Database seeded successfully!
echo.
echo Sample users created (username = password):
echo   - rgodwin (super_admin)
echo   - csamad (super_admin)
echo   - eknutson (super_admin)
echo   - jminihan (super_admin)
echo   - test (user)
echo   - asd1 (user)
echo   - david (user)
echo.
echo Sample communities:
echo   - Cats
echo   - Dogs
echo   - Harold
echo   - Gardening
echo   - star_wars
echo.
echo 🎉 Ready to explore!
