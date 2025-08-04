import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/app.dart';
import 'core/config/app_config.dart';
import 'core/services/storage_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // 初始化存储服务
  await StorageService.init();
  
  // 初始化应用配置
  await AppConfig.init();
  
  runApp(
    const ProviderScope(
      child: SocialApp(),
    ),
  );
}