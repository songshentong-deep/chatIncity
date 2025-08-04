import 'package:flutter/material.dart';
import 'package:dio/dio.dart';

/// 错误处理工具类
class ErrorHandler {
  static final Map<String, DateTime> _lastShownErrors = {};
  static const Duration _errorCooldown = Duration(seconds: 2);
  
  /// 清除特定操作的错误记录，允许重新显示错误
  static void clearErrorRecord(String operationType) {
    _lastShownErrors.remove(operationType);
  }

  /// 根据错误类型显示相应的弹窗，避免重复显示
  static void showErrorDialog(BuildContext context, dynamic error) {
    final errorMessage = _extractErrorMessage(error);
    final now = DateTime.now();
    
    // 检查是否在冷却期内显示过相同错误
    if (_lastShownErrors.containsKey(errorMessage)) {
      final lastShown = _lastShownErrors[errorMessage]!;
      if (now.difference(lastShown) < _errorCooldown) {
        return; // 跳过重复显示
      }
    }
    
    _lastShownErrors[errorMessage] = now;
    
    String title = '错误';
    String message = errorMessage;
    
    if (error is DioException) {
      final errorInfo = _handleDioError(error);
      title = errorInfo['title'] ?? '错误';
      message = errorInfo['message'] ?? '发生未知错误，请稍后重试';
    }
    
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Row(
            children: [
              Icon(Icons.error_outline, color: Colors.red, size: 24),
              SizedBox(width: 8),
              Text(title, style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            ],
          ),
          content: Text(message, style: TextStyle(fontSize: 16)),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: Text('确定', style: TextStyle(color: Colors.blue)),
            ),
          ],
        );
      },
    );
  }
  
  /// 显示成功提示，避免重复显示
  static void showSuccessSnackBar(BuildContext context, String message) {
    final now = DateTime.now();
    
    // 检查是否在冷却期内显示过相同消息
    if (_lastShownErrors.containsKey(message)) {
      final lastShown = _lastShownErrors[message]!;
      if (now.difference(lastShown) < _errorCooldown) {
        return;
      }
    }
    
    _lastShownErrors[message] = now;
    
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            Icon(Icons.check_circle, color: Colors.white, size: 20),
            SizedBox(width: 8),
            Text(message, style: TextStyle(color: Colors.white)),
          ],
        ),
        backgroundColor: Colors.green,
        duration: Duration(seconds: 2),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
        margin: const EdgeInsets.all(16),
      ),
    );
  }
  
  /// 显示警告提示，避免重复显示
  static void showWarningSnackBar(BuildContext context, String message) {
    final now = DateTime.now();
    
    if (_lastShownErrors.containsKey(message)) {
      final lastShown = _lastShownErrors[message]!;
      if (now.difference(lastShown) < _errorCooldown) {
        return;
      }
    }
    
    _lastShownErrors[message] = now;
    
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            Icon(Icons.warning, color: Colors.white, size: 20),
            SizedBox(width: 8),
            Text(message, style: TextStyle(color: Colors.white)),
          ],
        ),
        backgroundColor: Colors.orange,
        duration: Duration(seconds: 3),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
        margin: const EdgeInsets.all(16),
      ),
    );
  }

  /// 显示错误提示，避免重复显示
  /// [operationType] 操作类型，用于区分不同的操作（如 'login', 'register'）
  static void showErrorSnackBar(BuildContext context, String message, {String? operationType}) {
    final now = DateTime.now();
    final key = operationType ?? message;
    
    if (_lastShownErrors.containsKey(key)) {
      final lastShown = _lastShownErrors[key]!;
      if (now.difference(lastShown) < _errorCooldown) {
        return;
      }
    }
    
    _lastShownErrors[key] = now;
    
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            const Icon(Icons.error_outline, color: Colors.white),
            const SizedBox(width: 8),
            Expanded(child: Text(message)),
          ],
        ),
        backgroundColor: Colors.red[600],
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
        margin: const EdgeInsets.all(16),
        duration: const Duration(seconds: 4),
      ),
    );
  }
  
  /// 处理 Dio 网络错误
  static Map<String, String> _handleDioError(DioException error) {
    String title = '网络错误';
    String message = '网络连接失败，请检查网络设置';
    
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
        title = '连接超时';
        message = '连接服务器超时，请检查网络连接';
        break;
      case DioExceptionType.sendTimeout:
        title = '发送超时';
        message = '发送数据超时，请重试';
        break;
      case DioExceptionType.receiveTimeout:
        title = '接收超时';
        message = '接收数据超时，请重试';
        break;
      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode;
        final responseData = error.response?.data;
        
        switch (statusCode) {
          case 400:
            title = '请求错误';
            message = _extractErrorMessage(responseData) ?? '请求参数错误，请检查输入信息';
            break;
          case 401:
            title = '认证失败';
            message = _extractErrorMessage(responseData) ?? '用户名或密码错误';
            break;
          case 403:
            title = '权限不足';
            message = _extractErrorMessage(responseData) ?? '您没有权限执行此操作';
            break;
          case 404:
            title = '服务未找到';
            message = '请求的服务不存在';
            break;
          case 409:
            title = '数据冲突';
            message = _extractErrorMessage(responseData) ?? '用户已存在或数据冲突';
            break;
          case 422:
            title = '数据验证失败';
            message = _extractErrorMessage(responseData) ?? '输入数据格式不正确';
            break;
          case 500:
            title = '服务器错误';
            message = '服务器内部错误，请稍后重试';
            break;
          case 502:
            title = '网关错误';
            message = '服务器网关错误，请稍后重试';
            break;
          case 503:
            title = '服务不可用';
            message = '服务暂时不可用，请稍后重试';
            break;
          default:
            title = '请求失败';
            message = _extractErrorMessage(responseData) ?? '请求失败，状态码: $statusCode';
        }
        break;
      case DioExceptionType.cancel:
        title = '请求取消';
        message = '请求已被取消';
        break;
      case DioExceptionType.connectionError:
        title = '连接错误';
        message = '无法连接到服务器，请检查网络连接';
        break;
      default:
        title = '网络错误';
        message = '网络请求失败，请稍后重试';
    }
    
    return {'title': title, 'message': message};
  }
  
  /// 从响应数据中提取错误消息
  static String _extractErrorMessage(dynamic responseData) {
    if (responseData == null) return '发生未知错误，请稍后重试';
    
    try {
      if (responseData is Map<String, dynamic>) {
        // 尝试从不同的字段中提取错误消息
        return responseData['message'] ?? 
               responseData['error'] ?? 
               responseData['msg'] ?? 
               responseData['detail'] ??
               '发生未知错误，请稍后重试';
      } else if (responseData is String) {
        return responseData;
      }
    } catch (e) {
      // 解析失败，返回默认消息
    }
    
    return '发生未知错误，请稍后重试';
  }
  
  /// 显示加载对话框
  static void showLoadingDialog(BuildContext context, {String message = '加载中...'}) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
          content: Row(
            children: [
              CircularProgressIndicator(),
              SizedBox(width: 16),
              Text(message),
            ],
          ),
        );
      },
    );
  }
  
  /// 隐藏加载对话框
  static void hideLoadingDialog(BuildContext context) {
    Navigator.of(context).pop();
  }

  /// 清理过期的错误记录
  static void cleanupExpiredErrors() {
    final now = DateTime.now();
    _lastShownErrors.removeWhere((key, value) => 
        now.difference(value) > _errorCooldown);
  }
}

/// 自定义异常类
class AppException implements Exception {
  final String message;
  final int? code;
  
  AppException(this.message, {this.code});
  
  @override
  String toString() => message;
}

/// 认证异常
class AuthException extends AppException {
  AuthException(String message, {int? code}) : super(message, code: code);
}

/// 网络异常
class NetworkException extends AppException {
  NetworkException(String message, {int? code}) : super(message, code: code);
}

/// 验证异常
class ValidationException extends AppException {
  ValidationException(String message, {int? code}) : super(message, code: code);
}