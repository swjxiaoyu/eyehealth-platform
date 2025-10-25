@echo off
chcp 65001 >nul
echo ========================================
echo   Eye Health Platform - Install Dependencies
echo ========================================
echo.

:: Set color
color 0B

:: Check Node.js
echo [1/4] Checking Node.js installation...
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Node.js is not installed!
    echo.
    echo Please install Node.js first:
    echo 1. Go to https://nodejs.org/
    echo 2. Download and install the LTS version
    echo 3. Restart your computer after installation
    echo 4. Run this script again
    echo.
    pause
    exit /b 1
)

node --version
echo SUCCESS: Node.js is installed

:: Check npm
echo.
echo [2/4] Checking npm installation...
npm --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: npm is not available!
    pause
    exit /b 1
)

npm --version
echo SUCCESS: npm is available

:: Install backend dependencies
echo.
echo [3/4] Installing backend dependencies...
cd backend
echo Installing backend dependencies...
npm install
if %errorlevel% neq 0 (
    echo ERROR: Backend dependencies installation failed!
    echo Please check your internet connection and try again.
    pause
    exit /b 1
)
echo SUCCESS: Backend dependencies installed successfully!
cd ..

:: Install frontend dependencies
echo.
echo [4/4] Installing frontend dependencies...
cd frontend
echo Installing frontend dependencies...
npm install
if %errorlevel% neq 0 (
    echo ERROR: Frontend dependencies installation failed!
    echo Please check your internet connection and try again.
    pause
    exit /b 1
)
echo SUCCESS: Frontend dependencies installed successfully!
cd ..

:: Check Python for AI service
echo.
echo [Optional] Checking Python for AI service...
python --version >nul 2>&1
if %errorlevel% equ 0 (
    python --version
    echo Installing AI service dependencies...
    cd ai-service
    pip install -r requirements.txt
    if %errorlevel% neq 0 (
        echo WARNING: AI service dependencies installation failed!
        echo You can still run the platform without AI service.
    ) else (
        echo SUCCESS: AI service dependencies installed successfully!
    )
    cd ..
) else (
    echo WARNING: Python is not installed, AI service will be skipped.
    echo To enable AI service, please install Python from https://python.org/
)

echo.
echo ========================================
echo        Installation Complete!
echo ========================================
echo.
echo All dependencies have been installed successfully!
echo.
echo Next steps:
echo 1. Run 'start-all.bat' to start all services
echo 2. Or run 'start-dev.bat' for development mode
echo.
echo Service addresses:
echo   * Frontend: http://localhost:3000
echo   * Backend:  http://localhost:3001
echo   * AI Service: http://localhost:8000 (if Python is installed)
echo.
pause