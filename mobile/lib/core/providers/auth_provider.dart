import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../services/storage_service.dart';
import '../../features/auth/data/models/user_model.dart';

// 认证状态
class AuthState {
  final bool isAuthenticated;
  final UserModel? user;
  final String? token;
  final bool isLoading;
  final String? error;

  const AuthState({
    this.isAuthenticated = false,
    this.user,
    this.token,
    this.isLoading = false,
    this.error,
  });

  AuthState copyWith({
    bool? isAuthenticated,
    UserModel? user,
    String? token,
    bool? isLoading,
    String? error,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      user: user ?? this.user,
      token: token ?? this.token,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
    );
  }
}

// 认证状态管理器
class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier() : super(const AuthState()) {
    _loadAuthState();
  }

  // 加载认证状态
  Future<void> _loadAuthState() async {
    final token = await StorageService.getToken();
    final userJson = await StorageService.getUserInfo();
    
    if (token != null && userJson != null) {
      final user = UserModel.fromJson(userJson);
      state = state.copyWith(
        isAuthenticated: true,
        user: user,
        token: token,
      );
    }
  }

  // 登录
  Future<void> login(String token, UserModel user) async {
    state = state.copyWith(isLoading: true, error: null);
    
    try {
      await StorageService.saveToken(token);
      await StorageService.saveUserInfo(user.toJson());
      
      state = state.copyWith(
        isAuthenticated: true,
        user: user,
        token: token,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  // 登出
  Future<void> logout() async {
    await StorageService.clearAll();
    state = const AuthState();
  }

  // 更新用户信息
  Future<void> updateUser(UserModel user) async {
    await StorageService.saveUserInfo(user.toJson());
    state = state.copyWith(user: user);
  }

  // 清除错误
  void clearError() {
    state = state.copyWith(error: null);
  }
}

// Provider
final authStateProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier();
});