import 'package:flutter/foundation.dart';

class AppConfig {
  static late AppConfig _instance;
  static AppConfig get instance => _instance;
  
  // API配置
  final String _baseUrl;
  final String _wsUrl;
  
  // 静态访问器
  static String get baseUrl => _instance._baseUrl;
  static String get wsBaseUrl => _instance._wsUrl;
  static String get apiBaseUrl => _instance._baseUrl;
  
  // 应用配置
  late final String appName;
  late final String version;
  late final bool isDebug;
  
  // 第三方服务配置
  late final String mapApiKey;
  late final String pushAppKey;
  
  AppConfig._({
    required String baseUrl,
    required String wsUrl,
    required this.appName,
    required this.version,
    required this.isDebug,
    required this.mapApiKey,
    required this.pushAppKey,
  }) : _baseUrl = baseUrl,
       _wsUrl = wsUrl;
  
  static Future<void> init() async {
    _instance = AppConfig._(
      baseUrl: _getBaseUrl(),
      wsUrl: _getWsUrl(),
      appName: '社交App',
      version: '1.0.0',
      isDebug: kDebugMode,
      mapApiKey: const String.fromEnvironment('MAP_API_KEY', defaultValue: ''),
      pushAppKey: const String.fromEnvironment('PUSH_APP_KEY', defaultValue: ''),
    );
  }
  
  static String _getBaseUrl() {
    if (kDebugMode) {
      // 开发环境 - 用户服务运行在8001端口
      return 'http://localhost:8001/api/v1';
    } else {
      // 生产环境
      return 'https://api.socialapp.com/api/v1';
    }
  }
  
  static String _getWsUrl() {
    if (kDebugMode) {
      return 'ws://localhost:8004';
    } else {
      return 'wss://ws.socialapp.com';
    }
  }
}

// API端点常量
class ApiEndpoints {
  // 认证相关
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String refreshToken = '/auth/refresh';
  static const String logout = '/auth/logout';
  
  // 用户相关
  static const String profile = '/user/profile';
  static const String updateProfile = '/user/profile';
  static const String uploadAvatar = '/user/avatar';
  
  // 匹配相关
  static const String recommendations = '/match/recommendations';
  static const String like = '/match/like';
  static const String pass = '/match/pass';
  static const String matches = '/match/matches';
  
  // 聊天相关
  static const String chatRooms = '/chat/rooms';
  static const String messages = '/chat/messages';
  static const String sendMessage = '/chat/send';
  
  // 位置相关
  static const String updateLocation = '/location/update';
  static const String nearbyUsers = '/location/nearby';
}

// 应用常量
class AppConstants {
  // 缓存键
  static const String userTokenKey = 'user_token';
  static const String userInfoKey = 'user_info';
  static const String settingsKey = 'app_settings';
  
  // 默认值
  static const int defaultPageSize = 20;
  static const double defaultSearchRadius = 10.0; // 10公里
  static const int maxInterestTags = 10;
  static const int minAge = 18;
  static const int maxAge = 60;
}