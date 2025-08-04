import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:dio/dio.dart';
import '../../../../core/config/app_config.dart';
import '../../../../core/providers/auth_provider.dart';
import '../../data/models/user_model.dart';

// 认证表单状态
class AuthFormState {
  final bool isLoading;
  final String? error;

  const AuthFormState({
    this.isLoading = false,
    this.error,
  });

  AuthFormState copyWith({
    bool? isLoading,
    String? error,
  }) {
    return AuthFormState(
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
    );
  }
}

// 认证表单管理器
class AuthFormNotifier extends StateNotifier<AuthFormState> {
  final Ref ref;
  late final Dio _dio;

  AuthFormNotifier(this.ref) : super(const AuthFormState()) {
    _dio = Dio(BaseOptions(
      baseUrl: AppConfig.baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
    ));
  }

  // 登录
  Future<void> login({
    required String phone,
    required String password,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _dio.post('/auth/login', data: {
        'phone': phone,
        'password': password,
      });

      if (response.statusCode == 200) {
        final data = response.data['data'];
        final token = data['token'] as String;
        final userJson = data['user'] as Map<String, dynamic>;
        final user = UserModel.fromJson(userJson);

        await ref.read(authStateProvider.notifier).login(token, user);
        state = state.copyWith(isLoading: false);
      } else {
        throw Exception('登录失败');
      }
    } on DioException catch (e) {
      final errorMessage = _handleLoginError(e);
      state = state.copyWith(isLoading: false, error: errorMessage);
      throw Exception(errorMessage);
    } catch (e) {
      final errorMessage = '登录失败: ${e.toString()}';
      state = state.copyWith(isLoading: false, error: errorMessage);
      throw Exception(errorMessage);
    }
  }

  // 处理登录错误
  String _handleLoginError(DioException e) {
    // 网络连接错误
    if (e.type == DioExceptionType.connectionTimeout) {
      return '连接超时，请检查网络连接';
    }
    if (e.type == DioExceptionType.receiveTimeout) {
      return '请求超时，请重试';
    }
    if (e.type == DioExceptionType.connectionError) {
      return '网络连接失败，请检查网络设置';
    }

    // HTTP状态码错误
    final statusCode = e.response?.statusCode;
    final responseData = e.response?.data;
    
    // 尝试从响应中获取错误码和消息
    String? errorCode;
    String? errorMessage;
    
    if (responseData is Map<String, dynamic>) {
      errorCode = responseData['code']?.toString();
      errorMessage = responseData['message'] as String?;
    }

    // 根据HTTP状态码处理
    switch (statusCode) {
      case 400:
        return _handleBadRequestError(errorCode, errorMessage);
      case 401:
        return _handleUnauthorizedError(errorCode, errorMessage);
      case 403:
        return '账户被禁用，请联系客服';
      case 404:
        return '服务不可用，请稍后重试';
      case 429:
        return '登录尝试过于频繁，请稍后再试';
      case 500:
      case 502:
      case 503:
        return '服务器繁忙，请稍后重试';
      default:
        return errorMessage ?? '登录失败，请重试';
    }
  }

  // 处理400错误（请求参数错误）
  String _handleBadRequestError(String? errorCode, String? errorMessage) {
    switch (errorCode) {
      case 'INVALID_PHONE':
        return '手机号格式不正确';
      case 'INVALID_PASSWORD':
        return '密码格式不正确';
      case 'MISSING_PARAMS':
        return '请填写完整的登录信息';
      case 'PHONE_NOT_REGISTERED':
        return '该手机号尚未注册';
      default:
        return errorMessage ?? '请求参数错误';
    }
  }

  // 处理401错误（认证失败）
  String _handleUnauthorizedError(String? errorCode, String? errorMessage) {
    switch (errorCode) {
      case 'WRONG_PASSWORD':
        return '密码错误，请重新输入';
      case 'USER_NOT_FOUND':
        return '用户不存在，请检查手机号';
      case 'ACCOUNT_LOCKED':
        return '账户已被锁定，请联系客服';
      case 'ACCOUNT_SUSPENDED':
        return '账户已被暂停，请联系客服';
      case 'PASSWORD_EXPIRED':
        return '密码已过期，请重置密码';
      case 'TOO_MANY_ATTEMPTS':
        return '登录失败次数过多，请稍后再试';
      default:
        return errorMessage ?? '用户名或密码错误';
    }
  }

  // 注册
  Future<void> register({
    required String phone,
    String? email,
    required String password,
    required String nickname,
    required int age,
    required String gender,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _dio.post('/auth/register', data: {
        'phone': phone,
        'email': email,
        'password': password,
        'nickname': nickname,
        'age': age,
        'gender': gender,
      });

      if (response.statusCode == 200) {
        final data = response.data['data'];
        final token = data['token'] as String;
        final userJson = data['user'] as Map<String, dynamic>;
        final user = UserModel.fromJson(userJson);

        await ref.read(authStateProvider.notifier).login(token, user);
        state = state.copyWith(isLoading: false);
      } else {
        throw Exception('注册失败');
      }
    } on DioException catch (e) {
      final errorMessage = _handleRegisterError(e);
      state = state.copyWith(isLoading: false, error: errorMessage);
      throw Exception(errorMessage);
    } catch (e) {
      final errorMessage = '注册失败: ${e.toString()}';
      state = state.copyWith(isLoading: false, error: errorMessage);
      throw Exception(errorMessage);
    }
  }

  // 处理注册错误
  String _handleRegisterError(DioException e) {
    // 网络连接错误
    if (e.type == DioExceptionType.connectionTimeout) {
      return '连接超时，请检查网络连接';
    }
    if (e.type == DioExceptionType.receiveTimeout) {
      return '请求超时，请重试';
    }
    if (e.type == DioExceptionType.connectionError) {
      return '网络连接失败，请检查网络设置';
    }

    // HTTP状态码错误
    final statusCode = e.response?.statusCode;
    final responseData = e.response?.data;
    
    // 尝试从响应中获取错误码和消息
    String? errorCode;
    String? errorMessage;
    
    if (responseData is Map<String, dynamic>) {
      errorCode = responseData['code']?.toString();
      errorMessage = responseData['message'] as String?;
    }

    // 根据HTTP状态码处理
    switch (statusCode) {
      case 400:
        return _handleRegisterBadRequestError(errorCode, errorMessage);
      case 409:
        return _handleRegisterConflictError(errorCode, errorMessage);
      case 422:
        return _handleRegisterValidationError(errorCode, errorMessage);
      case 429:
        return '注册请求过于频繁，请稍后再试';
      case 500:
      case 502:
      case 503:
        return '服务器繁忙，请稍后重试';
      default:
        return errorMessage ?? '注册失败，请重试';
    }
  }

  // 处理注册400错误（请求参数错误）
  String _handleRegisterBadRequestError(String? errorCode, String? errorMessage) {
    switch (errorCode) {
      case 'INVALID_PHONE':
        return '手机号格式不正确';
      case 'INVALID_EMAIL':
        return '邮箱格式不正确';
      case 'INVALID_PASSWORD':
        return '密码格式不正确，至少6位字符';
      case 'INVALID_NICKNAME':
        return '昵称格式不正确，2-20个字符';
      case 'INVALID_AGE':
        return '年龄必须在18-100之间';
      case 'INVALID_GENDER':
        return '请选择正确的性别';
      case 'MISSING_PARAMS':
        return '请填写完整的注册信息';
      default:
        return errorMessage ?? '请求参数错误';
    }
  }

  // 处理注册409错误（资源冲突）
  String _handleRegisterConflictError(String? errorCode, String? errorMessage) {
    switch (errorCode) {
      case 'PHONE_EXISTS':
        return '该手机号已被注册，请直接登录';
      case 'EMAIL_EXISTS':
        return '该邮箱已被注册，请更换邮箱';
      case 'NICKNAME_EXISTS':
        return '该昵称已被使用，请更换昵称';
      default:
        return errorMessage ?? '用户信息已存在';
    }
  }

  // 处理注册422错误（数据验证失败）
  String _handleRegisterValidationError(String? errorCode, String? errorMessage) {
    switch (errorCode) {
      case 'WEAK_PASSWORD':
        return '密码强度不够，请包含字母和数字';
      case 'INVALID_PHONE_FORMAT':
        return '请输入正确的手机号格式';
      case 'INVALID_EMAIL_FORMAT':
        return '请输入正确的邮箱格式';
      case 'NICKNAME_TOO_SHORT':
        return '昵称至少需要2个字符';
      case 'NICKNAME_TOO_LONG':
        return '昵称不能超过20个字符';
      case 'AGE_TOO_YOUNG':
        return '用户年龄必须满18岁';
      case 'AGE_TOO_OLD':
        return '用户年龄不能超过100岁';
      default:
        return errorMessage ?? '数据验证失败';
    }
  }

  // 清除错误
  void clearError() {
    state = state.copyWith(error: null);
  }
}

// Provider
final authFormProvider = StateNotifierProvider<AuthFormNotifier, AuthFormState>((ref) {
  return AuthFormNotifier(ref);
});