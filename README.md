# 去中心化眼健康管理平台

## 项目概述

这是一个基于区块链技术的去中心化眼健康管理平台，提供用户数据管理、产品溯源、多模态AI推荐、远程医疗等功能。

## 技术架构

### 前端
- **Web**: Next.js + React + TypeScript + Tailwind CSS + shadcn/ui
- **Mobile**: React Native + Expo

### 后端
- **API Gateway**: 统一认证、限流、日志
- **微服务**: NestJS + TypeScript
- **AI服务**: Python + FastAPI + PyTorch + BentoML

### 数据存储
- **关系数据库**: PostgreSQL
- **对象存储**: MinIO (S3兼容)
- **向量数据库**: Milvus
- **去中心化存储**: IPFS

### 区块链
- **平台**: Hyperledger Fabric
- **智能合约**: Go/Java/Node.js

### 基础设施
- **容器化**: Docker + Kubernetes
- **CI/CD**: GitHub Actions + ArgoCD
- **监控**: Prometheus + Grafana + ELK
- **密钥管理**: HashiCorp Vault

## 项目结构

```
eyehealth-platform/
├── frontend/                 # Next.js 前端应用
├── mobile/                   # React Native 移动应用
├── backend/                  # NestJS 后端微服务
├── ai-service/              # Python FastAPI AI服务
├── blockchain/              # Hyperledger Fabric 区块链
├── infrastructure/          # Kubernetes 和基础设施配置
├── contracts/               # 智能合约
├── docs/                    # 项目文档
└── scripts/                 # 部署和工具脚本
```

## 核心功能

### MVP 功能列表
1. **用户认证** - DID + 传统登录
2. **文件上传** - 加密存储 + 区块链哈希记录
3. **AI推荐** - 基于多模态数据的个性化产品推荐
4. **产品溯源** - 二维码扫描展示产品全链路信息
5. **订单管理** - 支付集成 + 智能合约自动退款

### 非功能性需求
- **隐私安全**: 数据加密存储，只写区块链哈希
- **可用性**: MVP目标99.5%可用性
- **性能**: 推荐实时响应<300-500ms
- **可扩展性**: 微服务 + Kubernetes支持水平扩展
- **合规性**: 符合医疗数据合规要求（GDPR等）

## 快速开始

### 环境要求
- Node.js 18+
- Python 3.10+
- Docker & Docker Compose
- Kubernetes (可选，用于生产环境)

### 本地开发

1. **启动基础设施服务**
```bash
cd infrastructure
docker-compose up -d
```

2. **启动后端服务**
```bash
cd backend
npm install
npm run start:dev
```

3. **启动AI服务**
```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

4. **启动前端应用**
```bash
cd frontend
npm install
npm run dev
```

## 部署

### 开发环境
```bash
./scripts/deploy-dev.sh
```

### 生产环境
```bash
./scripts/deploy-prod.sh
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目链接: [https://github.com/swjxiaoyu/eyehealth-platform](https://github.com/swjxiaoyu/eyehealth-platform)
- 问题反馈: [Issues](https://github.com/swjxiaoyu/eyehealth-platform/issues)