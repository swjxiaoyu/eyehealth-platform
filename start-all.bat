@echo off
chcp 65001 >nul
echo ========================================
echo   Eye Health Platform - Start All Services
echo ========================================
echo.

:: Set color
color 0A

:: Check Node.js
echo [1/7] Checking Node.js installation...
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Node.js is not installed, please install Node.js first
    echo Download from: https://nodejs.org/
    pause
    exit /b 1
)
echo SUCCESS: Node.js is installed

:: Check npm
echo Checking npm installation...
npm --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: npm is not available
    pause
    exit /b 1
)
echo SUCCESS: npm is available

:: Install backend dependencies
echo.
echo [2/7] Installing backend dependencies...
cd backend
if not exist "node_modules" (
    echo Installing backend dependencies...
    npm install
    if %errorlevel% neq 0 (
        echo ERROR: Backend dependencies installation failed
        pause
        exit /b 1
    )
    echo SUCCESS: Backend dependencies installed
) else (
    echo SUCCESS: Backend dependencies already installed
)
cd ..

:: Install frontend dependencies
echo.
echo [3/7] Installing frontend dependencies...
cd frontend
if not exist "node_modules" (
    echo Installing frontend dependencies...
    npm install
    if %errorlevel% neq 0 (
        echo ERROR: Frontend dependencies installation failed
        pause
        exit /b 1
    )
    echo SUCCESS: Frontend dependencies installed
) else (
    echo SUCCESS: Frontend dependencies already installed
)
cd ..

:: Check Docker status
echo.
echo [4/7] Checking Docker service status...
docker version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Docker is not running, please start Docker Desktop first
    pause
    exit /b 1
)
echo SUCCESS: Docker service is running

:: Start MinIO service
echo.
echo [5/7] Starting MinIO storage service...
docker start eyehealth-minio >nul 2>&1
if %errorlevel% neq 0 (
    echo WARNING: MinIO container does not exist, creating...
    docker run -d --name eyehealth-minio -p 9000:9000 -p 9001:9001 -e MINIO_ROOT_USER=admin -e MINIO_ROOT_PASSWORD=password123 minio/minio server /data --console-address ":9001" >nul 2>&1
    if %errorlevel% neq 0 (
        echo ERROR: MinIO startup failed
        pause
        exit /b 1
    )
    echo SUCCESS: MinIO container created and started
) else (
    echo SUCCESS: MinIO service started
)

:: Wait for MinIO to be ready
echo    Waiting for MinIO service to be ready...
timeout /t 3 /nobreak >nul

:: Start backend service
echo.
echo [6/7] Starting backend service...
cd backend
start "Backend Service" cmd /k "npm run start:dev"
cd ..
echo SUCCESS: Backend service starting...

:: Wait for backend to start
echo    Waiting for backend service to be ready...
timeout /t 8 /nobreak >nul

:: Start frontend service
echo.
echo [7/7] Starting frontend service...
cd frontend
start "Frontend Service" cmd /k "npm run dev"
cd ..
echo SUCCESS: Frontend service starting...

:: Wait for frontend to start
echo    Waiting for frontend service to be ready...
timeout /t 5 /nobreak >nul

:: Start AI service (optional)
echo.
echo [Optional] Starting AI service...
cd ai-service
if exist "main.py" (
    python --version >nul 2>&1
    if %errorlevel% equ 0 (
        if not exist "venv" (
            echo Creating Python virtual environment...
            python -m venv venv
        )
        echo Installing AI service dependencies...
        pip install -r requirements.txt >nul 2>&1
        start "AI Service" cmd /k "python main.py"
        echo SUCCESS: AI service starting...
    ) else (
        echo WARNING: Python is not installed, skipping AI service
    )
) else (
    echo WARNING: AI service file does not exist, skipping
)
cd ..

:: Show service status
echo.
echo ========================================
echo           Services Started Successfully
echo ========================================
echo.
echo Service Status:
echo   * MinIO Storage: http://localhost:9000 (admin/password123)
echo   * MinIO Console: http://localhost:9001
echo   * Backend API:   http://localhost:3001
echo   * Frontend App:  http://localhost:3000
echo   * AI Service:    http://localhost:8000 (if started)
echo.
echo Management Commands:
echo   * Stop all services: stop-all.bat
echo   * View Docker containers: docker ps
echo.
echo Quick Access:
echo   * Open frontend app: http://localhost:3000
echo   * API documentation: http://localhost:3001/api
echo.

:: Ask if open browser
set /p choice="Open browser to access frontend app? (y/n): "
if /i "%choice%"=="y" (
    start http://localhost:3000
    echo SUCCESS: Browser opened
)

echo.
echo All services started successfully!
echo    Press any key to exit...
pause >nul