-- 社交App数据库初始化脚本

-- 数据库已通过环境变量 POSTGRES_DB 创建，直接使用

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建兴趣标签表
CREATE TABLE IF NOT EXISTS interest_tags (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    category VARCHAR(30) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    color VARCHAR(7),
    usage_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_interest_tags_category ON interest_tags(category);
CREATE INDEX IF NOT EXISTS idx_interest_tags_active ON interest_tags(is_active);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为interest_tags表创建更新时间触发器
CREATE TRIGGER update_interest_tags_updated_at 
    BEFORE UPDATE ON interest_tags 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();