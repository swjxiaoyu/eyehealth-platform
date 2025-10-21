# 去中心化眼健康管理平台 - 快速启动指南

## 🚀 一键启动

### 方法一：使用批处理脚本（推荐）

1. **启动所有服务**
   ```bash
   start-all.bat
   ```
   - 自动启动前端和后端服务
   - 自动打开浏览器访问 http://localhost:3000
   - 服务运行在独立窗口中

2. **停止所有服务**
   ```bash
   stop-all.bat
   ```
   - 自动停止所有相关进程
   - 释放端口 3000 和 3001

### 方法二：手动启动

1. **启动前端**
   ```bash
   cd frontend
   npm run dev
   ```

2. **启动后端**
   ```bash
   cd backend
   npm run start:dev
   ```

## 📋 服务地址

- **前端应用**: http://localhost:3000
- **后端API**: http://localhost:3001
- **API文档**: http://localhost:3001/api/docs

## 🔧 环境要求

- Node.js v18+
- PostgreSQL 数据库
- npm 包管理器

## 📝 功能测试

### 用户注册测试
1. 访问 http://localhost:3000/auth/register
2. 填写注册信息：
   - 姓名：必填
   - 邮箱：必填，格式正确
   - 出生日期：必填
   - 密码：至少6位
   - 确认密码：与密码一致

### API测试
```bash
# 测试注册API
node test-api-new.js

# 测试其他API
curl http://localhost:3001/api/v1/health
```

## 🐛 常见问题

### 1. 端口被占用
```bash
# 查看端口占用
netstat -ano | findstr "3000"
netstat -ano | findstr "3001"

# 停止占用进程
taskkill /PID <进程ID> /F
```

### 2. 数据库连接失败
- 确保PostgreSQL服务正在运行
- 检查数据库配置：`backend/.env`
- 数据库密码：`swj21bsss`

### 3. 前端构建错误
```bash
cd frontend
npm install
npm run build
```

### 4. 后端启动失败
```bash
cd backend
npm install
npm run build
npm run start:dev
```

## 📊 项目结构

```
eyehealth-platform/
├── frontend/          # Next.js 前端应用
├── backend/           # NestJS 后端API
├── ai-service/        # Python AI服务
├── infrastructure/    # Docker和K8s配置
├── start-all.bat      # 一键启动脚本
├── stop-all.bat       # 一键停止脚本
└── README.md          # 项目说明
```

## 🎯 下一步

1. ✅ 前后端服务正常运行
2. ✅ 用户注册功能正常
3. 🔄 完善其他功能模块
4. 🔄 集成AI服务
5. 🔄 部署到生产环境

## 📞 技术支持

如遇到问题，请检查：
1. 控制台错误信息
2. 网络连接状态
3. 服务端口占用情况
4. 数据库连接状态