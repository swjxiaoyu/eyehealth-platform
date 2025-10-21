@echo off
echo ========================================
echo    去中心化眼健康管理平台 - 一键启动
echo ========================================
echo.

echo [1/4] 检查服务状态...
netstat -ano | findstr "3000" >nul
if %errorlevel% equ 0 (
    echo 前端服务已在运行 (端口 3000)
) else (
    echo 前端服务未运行
)

netstat -ano | findstr "3001" >nul
if %errorlevel% equ 0 (
    echo 后端服务已在运行 (端口 3001)
) else (
    echo 后端服务未运行
)

echo.
echo [2/4] 启动前端服务...
cd frontend
start "前端服务" cmd /k "npm run dev"
cd ..

echo.
echo [3/4] 等待前端启动...
timeout /t 3 /nobreak >nul

echo [4/4] 启动后端服务...
cd backend
start "后端服务" cmd /k "npm run start:dev"
cd ..

echo.
echo ========================================
echo           启动完成！
echo ========================================
echo.
echo 前端地址: http://localhost:3000
echo 后端地址: http://localhost:3001
echo API文档: http://localhost:3001/api/docs
echo.
echo 按任意键打开浏览器...
pause >nul

echo 正在打开浏览器...
start http://localhost:3000

echo.
echo 提示：
echo - 前端和后端服务将在新窗口中运行
echo - 关闭对应窗口即可停止服务
echo - 如需停止所有服务，请运行 stop-services.bat
echo.
pause