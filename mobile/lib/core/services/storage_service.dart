import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';

class StorageService {
  static SharedPreferences? _prefs;
  
  static const String _tokenKey = 'auth_token';
  static const String _userInfoKey = 'user_info';
  static const String _settingsKey = 'app_settings';

  static Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }

  // Token 相关
  static Future<void> saveToken(String token) async {
    await _prefs?.setString(_tokenKey, token);
  }

  static Future<String?> getToken() async {
    return _prefs?.getString(_tokenKey);
  }

  static Future<void> removeToken() async {
    await _prefs?.remove(_tokenKey);
  }

  // 用户信息相关
  static Future<void> saveUserInfo(Map<String, dynamic> userInfo) async {
    final jsonString = jsonEncode(userInfo);
    await _prefs?.setString(_userInfoKey, jsonString);
  }

  static Future<Map<String, dynamic>?> getUserInfo() async {
    final jsonString = _prefs?.getString(_userInfoKey);
    if (jsonString != null) {
      return jsonDecode(jsonString) as Map<String, dynamic>;
    }
    return null;
  }

  static Future<void> removeUserInfo() async {
    await _prefs?.remove(_userInfoKey);
  }

  // 应用设置相关
  static Future<void> saveSettings(Map<String, dynamic> settings) async {
    final jsonString = jsonEncode(settings);
    await _prefs?.setString(_settingsKey, jsonString);
  }

  static Future<Map<String, dynamic>?> getSettings() async {
    final jsonString = _prefs?.getString(_settingsKey);
    if (jsonString != null) {
      return jsonDecode(jsonString) as Map<String, dynamic>;
    }
    return null;
  }

  // 通用方法
  static Future<void> saveString(String key, String value) async {
    await _prefs?.setString(key, value);
  }

  static Future<String?> getString(String key) async {
    return _prefs?.getString(key);
  }

  static Future<void> saveBool(String key, bool value) async {
    await _prefs?.setBool(key, value);
  }

  static Future<bool?> getBool(String key) async {
    return _prefs?.getBool(key);
  }

  static Future<void> saveInt(String key, int value) async {
    await _prefs?.setInt(key, value);
  }

  static Future<int?> getInt(String key) async {
    return _prefs?.getInt(key);
  }

  // 清除所有数据
  static Future<void> clearAll() async {
    await _prefs?.clear();
  }

  // 清除特定键
  static Future<void> remove(String key) async {
    await _prefs?.remove(key);
  }
}