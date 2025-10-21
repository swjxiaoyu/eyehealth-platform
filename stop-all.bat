@echo off
echo ========================================
echo    去中心化眼健康管理平台 - 停止服务
echo ========================================
echo.

echo [1/3] 查找并停止前端服务 (端口 3000)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr "3000"') do (
    echo 停止进程 %%a
    taskkill /PID %%a /F >nul 2>&1
)

echo [2/3] 查找并停止后端服务 (端口 3001)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr "3001"') do (
    echo 停止进程 %%a
    taskkill /PID %%a /F >nul 2>&1
)

echo [3/3] 停止所有Node.js进程...
taskkill /IM node.exe /F >nul 2>&1

echo.
echo ========================================
echo           所有服务已停止！
echo ========================================
echo.
pause