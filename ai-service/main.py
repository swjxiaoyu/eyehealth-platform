from fastapi import FastAPI, HTTPException, UploadFile, File, Form
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn
import os
from typing import List, Dict, Any, Optional
import logging
from datetime import datetime
import json

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# 创建FastAPI应用
app = FastAPI(
    title="眼健康AI推荐服务",
    description="基于多模态数据的眼健康产品推荐AI服务",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# 配置CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://localhost:3001"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 全局变量存储模型（实际应用中应该使用模型管理服务）
models = {}
embeddings_cache = {}

@app.on_event("startup")
async def startup_event():
    """应用启动时初始化"""
    logger.info("AI服务启动中...")
    
    # 初始化模型
    await initialize_models()
    
    logger.info("AI服务启动完成")

@app.on_event("shutdown")
async def shutdown_event():
    """应用关闭时清理"""
    logger.info("AI服务关闭中...")
    # 清理资源
    logger.info("AI服务已关闭")

async def initialize_models():
    """初始化AI模型"""
    try:
        # 这里应该加载实际的模型
        # 目前使用模拟数据
        models["text_model"] = "text_model_loaded"
        models["vision_model"] = "vision_model_loaded"
        models["multimodal_model"] = "multimodal_model_loaded"
        
        logger.info("模型初始化完成")
    except Exception as e:
        logger.error(f"模型初始化失败: {e}")
        raise

@app.get("/")
async def root():
    """根路径"""
    return {
        "message": "眼健康AI推荐服务",
        "version": "1.0.0",
        "status": "running",
        "timestamp": datetime.now().isoformat()
    }

@app.get("/health")
async def health_check():
    """健康检查"""
    return {
        "status": "healthy",
        "timestamp": datetime.now().isoformat(),
        "models": list(models.keys())
    }

@app.post("/api/v1/recommendation")
async def get_recommendations(
    user_id: str = Form(...),
    context: str = Form(...),
    top_k: int = Form(5)
):
    """获取个性化推荐"""
    try:
        # 解析上下文数据
        context_data = json.loads(context)
        
        # 模拟推荐逻辑
        recommendations = await generate_recommendations(user_id, context_data, top_k)
        
        return {
            "recommendations": recommendations,
            "model_version": "v1.0.0",
            "confidence": 0.85,
            "processing_time": 0.5,
            "timestamp": datetime.now().isoformat()
        }
        
    except Exception as e:
        logger.error(f"推荐生成失败: {e}")
        raise HTTPException(status_code=500, detail="推荐生成失败")

@app.post("/api/v1/analyze")
async def analyze_data(
    file: UploadFile = File(...),
    analysis_type: str = Form("multimodal")
):
    """分析上传的数据"""
    try:
        # 读取文件内容
        content = await file.read()
        
        # 根据文件类型进行分析
        if file.content_type.startswith('image/'):
            result = await analyze_image(content)
        elif file.content_type == 'application/pdf':
            result = await analyze_pdf(content)
        elif file.content_type == 'text/plain':
            result = await analyze_text(content)
        else:
            raise HTTPException(status_code=400, detail="不支持的文件类型")
        
        return {
            "analysis_type": analysis_type,
            "file_name": file.filename,
            "file_type": file.content_type,
            "result": result,
            "timestamp": datetime.now().isoformat()
        }
        
    except Exception as e:
        logger.error(f"数据分析失败: {e}")
        raise HTTPException(status_code=500, detail="数据分析失败")

@app.post("/api/v1/embedding")
async def generate_embedding(
    text: str = Form(...),
    model: str = Form("text")
):
    """生成文本嵌入向量"""
    try:
        # 检查缓存
        if text in embeddings_cache:
            return {
                "embedding": embeddings_cache[text],
                "model": model,
                "cached": True
            }
        
        # 生成嵌入向量（模拟）
        embedding = await generate_text_embedding(text)
        
        # 缓存结果
        embeddings_cache[text] = embedding
        
        return {
            "embedding": embedding,
            "model": model,
            "cached": False,
            "dimension": len(embedding)
        }
        
    except Exception as e:
        logger.error(f"嵌入向量生成失败: {e}")
        raise HTTPException(status_code=500, detail="嵌入向量生成失败")

async def generate_recommendations(user_id: str, context: Dict[str, Any], top_k: int) -> List[Dict[str, Any]]:
    """生成推荐结果（模拟）"""
    # 模拟推荐逻辑
    recommendations = [
        {
            "product_id": "prod_001",
            "product_name": "护眼蓝光眼镜",
            "score": 0.95,
            "reason": "基于您的屏幕使用时间推荐",
            "category": "protective_equipment"
        },
        {
            "product_id": "prod_002", 
            "product_name": "叶黄素补充剂",
            "score": 0.88,
            "reason": "适合您的年龄和用眼习惯",
            "category": "supplements"
        },
        {
            "product_id": "prod_003",
            "product_name": "眼部按摩仪",
            "score": 0.82,
            "reason": "缓解眼部疲劳",
            "category": "equipment"
        }
    ]
    
    return recommendations[:top_k]

async def analyze_image(content: bytes) -> Dict[str, Any]:
    """分析图像数据（模拟）"""
    return {
        "type": "image_analysis",
        "detected_objects": ["eye", "face"],
        "quality_score": 0.85,
        "recommendations": ["建议进行专业眼科检查"]
    }

async def analyze_pdf(content: bytes) -> Dict[str, Any]:
    """分析PDF文档（模拟）"""
    return {
        "type": "document_analysis",
        "extracted_text": "模拟提取的文本内容",
        "confidence": 0.90,
        "key_metrics": {
            "vision_acuity": "20/20",
            "pressure": "正常"
        }
    }

async def analyze_text(content: bytes) -> Dict[str, Any]:
    """分析文本数据（模拟）"""
    text = content.decode('utf-8')
    return {
        "type": "text_analysis",
        "sentiment": "positive",
        "keywords": ["眼睛", "健康", "检查"],
        "summary": "用户关注眼健康问题"
    }

async def generate_text_embedding(text: str) -> List[float]:
    """生成文本嵌入向量（模拟）"""
    # 模拟生成384维的嵌入向量
    import random
    random.seed(hash(text))
    return [random.uniform(-1, 1) for _ in range(384)]

if __name__ == "__main__":
    # 从环境变量获取配置
    host = os.getenv("API_HOST", "0.0.0.0")
    port = int(os.getenv("API_PORT", 8000))
    workers = int(os.getenv("API_WORKERS", 1))
    
    uvicorn.run(
        "main:app",
        host=host,
        port=port,
        workers=workers,
        reload=True if os.getenv("NODE_ENV") == "development" else False
    )