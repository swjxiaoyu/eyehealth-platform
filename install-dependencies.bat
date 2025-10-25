@echo off
chcp 65001 >nul
echo ========================================
echo   Eye Health Platform - Install Dependencies
echo ========================================
echo.

:: Set color
color 0B

:: Install backend dependencies
echo [1/3] Installing backend dependencies...
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
echo [2/3] Installing frontend dependencies...
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

:: Install AI service dependencies
echo.
echo [3/3] Installing AI service dependencies...
cd ai-service
echo Installing AI service dependencies...
pip install -r requirements.txt
if %errorlevel% neq 0 (
    echo WARNING: AI service dependencies installation failed!
    echo You can still run the platform without AI service.
) else (
    echo SUCCESS: AI service dependencies installed successfully!
)
cd ..

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