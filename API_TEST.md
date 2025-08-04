# API 接口测试文档

## 🚀 服务状态

### ✅ 已启动的服务
- **用户服务** (端口 8001) - 登录注册接口 ✅
- **认证服务** (端口 8002) - Token验证接口 ✅
- **其他微服务** (端口 8003-8008) - 各功能服务 ✅

### 🗄️ 数据库状态
- **PostgreSQL** (端口 5432) ✅
- **MongoDB** (端口 27017) ✅  
- **Redis** (端口 6379) ✅
- **数据库表已迁移** ✅
- **兴趣标签数据已初始化** ✅

## 📱 前端接口调用指南

### 🔑 认证接口 (用户服务 - 端口 8001)

#### 1. 用户注册
```bash
POST http://localhost:8001/api/v1/auth/register
Content-Type: application/json

{
  "phone": "13800138000",
  "email": "test@example.com", 
  "password": "123456",
  "nickname": "测试用户",
  "age": 25,
  "gender": "male"  // male, female, other
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user": {
      "id": "71f2aee5-c967-4e53-8065-00cf24fc2c00",
      "phone": "13800138000",
      "email": "test@example.com",
      "nickname": "测试用户",
      "avatar": "",
      "age": 25,
      "gender": "male"
    }
  }
}
```

#### 2. 用户登录
```bash
POST http://localhost:8001/api/v1/auth/login
Content-Type: application/json

{
  "phone": "13800138000",    // 手机号或邮箱二选一
  "email": "test@example.com", // 手机号或邮箱二选一
  "password": "123456"
}
```

**响应格式同注册接口**

#### 3. 刷新Token
```bash
POST http://localhost:8001/api/v1/auth/refresh
Authorization: Bearer <your_token>
```

### 👤 用户信息接口 (需要认证)

#### 1. 获取用户资料
```bash
GET http://localhost:8001/api/v1/user/profile
Authorization: Bearer <your_token>
```

#### 2. 更新用户资料
```bash
PUT http://localhost:8001/api/v1/user/profile
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "bio": "个人简介",
  "location": "北京市"
}
```

#### 3. 获取兴趣标签列表
```bash
GET http://localhost:8001/api/v1/interests
```

## 🧪 测试命令

### 快速测试脚本
```bash
# 测试注册
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "phone":"13900139000",
    "email":"user2@example.com",
    "password":"123456",
    "nickname":"用户2",
    "age":28,
    "gender":"female"
  }'

# 测试登录
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone":"13900139000",
    "password":"123456"
  }'

# 测试获取兴趣标签
curl http://localhost:8001/api/v1/interests
```

## 🔧 前端配置建议

### 1. API Base URL 配置
```javascript
// 开发环境
const API_BASE_URL = 'http://localhost:8001/api/v1';

// 生产环境
const API_BASE_URL = 'https://your-domain.com/api/v1';
```

### 2. 请求拦截器配置
```javascript
// 添加认证头
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器处理错误
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Token过期，跳转到登录页
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### 3. Flutter HTTP 配置
```dart
class ApiService {
  static const String baseUrl = 'http://localhost:8001/api/v1';
  
  static Future<Map<String, String>> _getHeaders() async {
    final token = await SharedPreferences.getInstance()
        .then((prefs) => prefs.getString('token'));
    
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }
  
  static Future<http.Response> post(String endpoint, Map<String, dynamic> data) {
    return http.post(
      Uri.parse('$baseUrl$endpoint'),
      headers: await _getHeaders(),
      body: jsonEncode(data),
    );
  }
}
```

## 🚨 常见问题解决

### 1. CORS 跨域问题
如果前端遇到跨域问题，后端已配置CORS中间件，支持：
- 所有来源 (`*`)
- 常用HTTP方法 (GET, POST, PUT, DELETE, OPTIONS)
- 认证头 (Authorization)

### 2. Token 过期处理
- Token 有效期：24小时
- 使用 `/api/v1/auth/refresh` 刷新Token
- 401状态码表示需要重新登录

### 3. 参数验证错误
注册接口必需参数：
- `phone` (必需)
- `password` (必需，最少6位)
- `nickname` (必需)
- `age` (必需，18-100岁)
- `gender` (必需，male/female/other)

## 📊 数据库查看

### 查看注册用户
```bash
docker exec social-app-postgres psql -U postgres -d social_app -c "SELECT id, phone, email, nickname, age, gender, created_at FROM users;"
```

### 查看兴趣标签
```bash
docker exec social-app-postgres psql -U postgres -d social_app -c "SELECT id, name, category, icon FROM interest_tags LIMIT 10;"
```

---

## ✅ 总结

**登录注册接口已经可以正常使用！**

- ✅ 用户服务 (8001端口) 正常运行
- ✅ 数据库连接正常
- ✅ 数据表已创建并初始化
- ✅ 注册接口测试通过
- ✅ 登录接口测试通过
- ✅ CORS已配置，支持前端调用

**前端可以直接调用 `http://localhost:8001/api/v1/auth/` 下的接口进行登录注册功能开发。**