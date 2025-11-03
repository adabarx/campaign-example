@echo off
setlocal

if "%1"=="" goto help
if "%1"=="install" goto install
if "%1"=="generate" goto generate
if "%1"=="build" goto build
if "%1"=="run" goto run
if "%1"=="dev" goto dev
if "%1"=="test" goto test
if "%1"=="clean" goto clean
if "%1"=="help" goto help

echo Unknown command: %1
goto help

:install
echo Installing dependencies...
go mod download
go install github.com/a-h/templ/cmd/templ@latest
goto end

:generate
echo Generating templ code and static HTML...
templ generate
go run main.go build.go --generate
goto end

:build
call :generate
echo Building binary...
go build -o campaign.exe
goto end

:run
call :generate
echo Running server...
go run main.go build.go
goto end

:dev
call :run
goto end

:test
echo Manual test:
echo 1. go run main.go --generate
echo 2. curl http://localhost:3000
echo 3. curl -X POST http://localhost:3000/api/donations -d "name=John&email=john@example.com&amount=50" -H "Content-Type: application/x-www-form-urlencoded"
echo 4. curl http://localhost:3000/api/stats
echo 5. curl http://localhost:3000/api/recent-donors
goto end

:clean
echo Cleaning build artifacts...
if exist public\ rmdir /s /q public
if exist campaign.exe del /q campaign.exe
if exist templates\*_templ.go del /q templates\*_templ.go
goto end

:help
echo Available commands:
echo   build.bat install   - Install dependencies
echo   build.bat generate  - Generate templ code and static HTML
echo   build.bat build     - Build binary
echo   build.bat run       - Run server with static generation
echo   build.bat dev       - Alias for run
echo   build.bat clean     - Clean build artifacts
echo   build.bat test      - Show manual test commands
goto end

:end
endlocal
