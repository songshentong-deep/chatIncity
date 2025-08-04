# 前端错误处理和弹窗反馈系统

## 🎯 功能概述

已实现完整的前端错误处理和弹窗反馈系统，包括：

- ✅ 统一的错误处理机制
- ✅ 根据后端状态码显示不同的错误弹窗
- ✅ 网络错误、认证错误、验证错误的分类处理
- ✅ 成功提示和警告提示
- ✅ 加载状态管理
- ✅ Token 自动管理和过期处理

## 📁 文件结构

```
mobile/lib/
├── core/
│   ├── utils/
│   │   └── error_handler.dart          # 错误处理工具类
│   └── services/
│       └── http_service.dart           # HTTP 服务类
└── features/
    └── auth/
        ├── data/
        │   └── auth_service.dart       # 认证服务
        ├── models/
        │   └── auth_models.dart        # 数据模型
        └── presentation/
            ├── login_page.dart         # 登录页面
            └── register_page.dart      # 注册页面
```

## 🔧 核心功能

### 1. 错误处理工具类 (ErrorHandler)

#### 显示错误弹窗
```dart
// 自动根据错误类型显示相应弹窗
ErrorHandler.showErrorDialog(context, error);
```

#### 显示成功提示
```dart
ErrorHandler.showSuccessSnackBar(context, '登录成功！');
```

#### 显示警告提示
```dart
ErrorHandler.showWarningSnackBar(context, '功能暂未开放');
```

#### 显示/隐藏加载对话框
```dart
ErrorHandler.showLoadingDialog(context, message: '登录中...');
ErrorHandler.hideLoadingDialog(context);
```

### 2. HTTP 服务类 (HttpService)

#### 自动错误处理
- 自动添加认证头
- 自动处理 Token 过期
- 统一的错误分类和处理
- 请求/响应日志记录

#### 使用示例
```dart
final httpService = HttpService();

try {
  final response = await httpService.post('/auth/login', data: loginData);
  // 处理成功响应
} catch (e) {
  // 错误会被自动分类为 AuthException、NetworkException 等
  ErrorHandler.showErrorDialog(context, e);
}
```

### 3. 认证服务 (AuthService)

#### 登录
```dart
final authService = AuthService();

try {
  final response = await authService.login(loginRequest);
  if (response.isSuccess) {
    // 登录成功，Token 已自动保存
  }
} catch (e) {
  // 显示错误弹窗
  ErrorHandler.showErrorDialog(context, e);
}
```

#### 注册
```dart
try {
  final response = await authService.register(registerRequest);
  if (response.isSuccess) {
    // 注册成功
  }
} catch (e) {
  ErrorHandler.showErrorDialog(context, e);
}
```

## 🎨 错误弹窗样式

### 错误对话框
- 红色错误图标
- 清晰的标题和消息
- 蓝色确定按钮
- 不可点击背景关闭

### 成功提示
- 绿色背景的 SnackBar
- 白色勾选图标
- 2秒自动消失
- 浮动样式

### 警告提示
- 橙色背景的 SnackBar
- 白色警告图标
- 3秒自动消失

## 📊 状态码处理

### 400 - 请求错误
- 标题：请求错误
- 消息：显示后端返回的具体错误信息
- 用于：参数验证失败等

### 401 - 认证失败
- 标题：认证失败
- 消息：用户名或密码错误
- 自动清除本地 Token

### 403 - 权限不足
- 标题：权限不足
- 消息：您没有权限执行此操作

### 404 - 服务未找到
- 标题：服务未找到
- 消息：请求的服务不存在

### 409 - 数据冲突
- 标题：数据冲突
- 消息：用户已存在或数据冲突
- 用于：注册时用户已存在等

### 422 - 数据验证失败
- 标题：数据验证失败
- 消息：输入数据格式不正确

### 500/502/503 - 服务器错误
- 标题：服务器错误
- 消息：服务器内部错误，请稍后重试

### 网络错误
- 连接超时
- 发送/接收超时
- 连接错误
- 请求取消

## 🔄 使用流程

### 1. 初始化 HTTP 服务
```dart
void main() {
  HttpService().init();
  runApp(MyApp());
}
```

### 2. 在页面中使用
```dart
class LoginPage extends StatefulWidget {
  // ...
  
  Future<void> _performLogin() async {
    setState(() {
      _isLoading = true;
    });
    
    try {
      final response = await _authService.login(loginRequest);
      
      if (response.isSuccess) {
        ErrorHandler.showSuccessSnackBar(context, '登录成功！');
        Navigator.pushReplacementNamed(context, '/home');
      } else {
        ErrorHandler.showErrorDialog(context, response.message);
      }
    } catch (e) {
      ErrorHandler.showErrorDialog(context, e);
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }
}
```

## 🎯 自定义错误类型

### AppException
- 基础异常类
- 包含错误消息和状态码

### AuthException
- 认证相关异常
- 继承自 AppException

### NetworkException
- 网络相关异常
- 继承自 AppException

### ValidationException
- 验证相关异常
- 继承自 AppException

## 📱 页面集成示例

### 登录页面特性
- 手机号格式验证
- 密码长度验证
- 加载状态显示
- 错误弹窗反馈
- 成功提示

### 注册页面特性
- 完整的表单验证
- 年龄和性别选择
- 密码确认验证
- 邮箱格式验证（可选）
- 统一的错误处理

## 🔧 配置要求

### pubspec.yaml 依赖
```yaml
dependencies:
  flutter:
    sdk: flutter
  dio: ^5.4.0
  shared_preferences: ^2.2.2
  # 其他依赖...
```

### 应用配置
确保在 `mobile/lib/core/config/app_config.dart` 中配置正确的 API 地址：
```dart
class AppConfig {
  static const String apiBaseUrl = 'http://localhost:8001/api/v1';
}
```

## 🎉 使用效果

### 成功场景
- 登录成功：绿色 SnackBar 提示 "登录成功！"
- 注册成功：绿色 SnackBar 提示 "注册成功！"
- 自动跳转到主页面

### 失败场景
- 网络错误：显示连接失败弹窗
- 用户名密码错误：显示认证失败弹窗
- 用户已存在：显示数据冲突弹窗
- 参数验证失败：显示具体的验证错误信息

### 加载状态
- 按钮显示加载动画
- 禁用重复提交
- 可选的全屏加载对话框

---

## 📞 技术支持

这个错误处理系统提供了完整的用户反馈机制，确保用户在任何情况下都能收到清晰的提示信息。所有的错误都会被适当分类和处理，提供良好的用户体验。