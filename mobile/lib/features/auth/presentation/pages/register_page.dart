import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/providers/auth_provider.dart';
import '../../../../core/utils/error_handler.dart';
import '../providers/auth_form_provider.dart';

class RegisterPage extends ConsumerStatefulWidget {
  const RegisterPage({Key? key}) : super(key: key);

  @override
  ConsumerState<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends ConsumerState<RegisterPage> {
  final _formKey = GlobalKey<FormState>();
  final _phoneController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  final _nicknameController = TextEditingController();
  final _ageController = TextEditingController();
  
  String _selectedGender = '男';
  bool _obscurePassword = true;
  bool _obscureConfirmPassword = true;
  bool _agreeToTerms = false;
  String? _lastErrorMessage; // 记录上次显示的错误消息

  @override
  void initState() {
    super.initState();
    
    // 监听认证表单状态变化
    ref.listenManual(authFormProvider, (previous, next) {
      if (next.error != null && next.error != _lastErrorMessage) {
        _lastErrorMessage = next.error;
        ErrorHandler.showErrorSnackBar(context, next.error!);
        // 清除错误状态，避免重复显示
        Future.delayed(const Duration(milliseconds: 100), () {
          ref.read(authFormProvider.notifier).clearError();
        });
      }
    });
  }

  @override
  void dispose() {
    _phoneController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    _nicknameController.dispose();
    _ageController.dispose();
    super.dispose();
  }

  void _handleRegister() async {
    if (_formKey.currentState!.validate() && _agreeToTerms) {
      // 清除之前的错误记录，允许显示新的错误
      _lastErrorMessage = null;
      
      try {
        await ref.read(authFormProvider.notifier).register(
          phone: _phoneController.text.trim(),
          email: _emailController.text.trim().isEmpty ? null : _emailController.text.trim(),
          password: _passwordController.text,
          nickname: _nicknameController.text.trim(),
          age: int.parse(_ageController.text),
          gender: _selectedGender,
        );
        
        if (mounted) {
          ErrorHandler.showSuccessSnackBar(context, '注册成功！');
          context.go('/home');
        }
      } catch (e) {
        // 错误已经在监听器中处理，这里不需要额外处理
      }
    }
  }



  @override
  Widget build(BuildContext context) {
    final authFormState = ref.watch(authFormProvider);
    final theme = Theme.of(context);
    
    return Scaffold(
      appBar: AppBar(
        title: const Text('注册账户'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.go('/login'),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const SizedBox(height: 20),
                
                // 标题
                Text(
                  '创建新账户',
                  style: theme.textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: Colors.black87,
                  ),
                  textAlign: TextAlign.center,
                ),
                
                const SizedBox(height: 8),
                
                Text(
                  '填写以下信息完成注册',
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: Colors.grey[600],
                  ),
                  textAlign: TextAlign.center,
                ),
                
                const SizedBox(height: 32),
                
                // 手机号
                TextFormField(
                  controller: _phoneController,
                  keyboardType: TextInputType.phone,
                  decoration: const InputDecoration(
                    labelText: '手机号',
                    hintText: '请输入手机号',
                    prefixIcon: Icon(Icons.phone),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '请输入手机号';
                    }
                    if (!RegExp(r'^1[3-9]\d{9}$').hasMatch(value)) {
                      return '请输入正确的手机号';
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 16),
                
                // 邮箱（可选）
                TextFormField(
                  controller: _emailController,
                  keyboardType: TextInputType.emailAddress,
                  decoration: const InputDecoration(
                    labelText: '邮箱（可选）',
                    hintText: '请输入邮箱地址',
                    prefixIcon: Icon(Icons.email),
                  ),
                  validator: (value) {
                    if (value != null && value.isNotEmpty) {
                      if (!RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(value)) {
                        return '请输入正确的邮箱格式';
                      }
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 16),
                
                // 昵称
                TextFormField(
                  controller: _nicknameController,
                  decoration: const InputDecoration(
                    labelText: '昵称',
                    hintText: '请输入昵称',
                    prefixIcon: Icon(Icons.person),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '请输入昵称';
                    }
                    if (value.length < 2 || value.length > 20) {
                      return '昵称长度应在2-20个字符之间';
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 16),
                
                // 年龄和性别
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _ageController,
                        keyboardType: TextInputType.number,
                        decoration: const InputDecoration(
                          labelText: '年龄',
                          hintText: '请输入年龄',
                          prefixIcon: Icon(Icons.cake),
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return '请输入年龄';
                          }
                          final age = int.tryParse(value);
                          if (age == null || age < 18 || age > 100) {
                            return '年龄应在18-100之间';
                          }
                          return null;
                        },
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: DropdownButtonFormField<String>(
                        value: _selectedGender,
                        decoration: const InputDecoration(
                          labelText: '性别',
                          prefixIcon: Icon(Icons.wc),
                        ),
                        items: ['男', '女', '其他'].map((String value) {
                          return DropdownMenuItem<String>(
                            value: value,
                            child: Text(value),
                          );
                        }).toList(),
                        onChanged: (String? newValue) {
                          setState(() {
                            _selectedGender = newValue!;
                          });
                        },
                      ),
                    ),
                  ],
                ),
                
                const SizedBox(height: 16),
                
                // 密码
                TextFormField(
                  controller: _passwordController,
                  obscureText: _obscurePassword,
                  decoration: InputDecoration(
                    labelText: '密码',
                    hintText: '请输入密码',
                    prefixIcon: const Icon(Icons.lock),
                    suffixIcon: IconButton(
                      icon: Icon(
                        _obscurePassword ? Icons.visibility : Icons.visibility_off,
                      ),
                      onPressed: () {
                        setState(() {
                          _obscurePassword = !_obscurePassword;
                        });
                      },
                    ),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '请输入密码';
                    }
                    if (value.length < 6) {
                      return '密码至少6位';
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 16),
                
                // 确认密码
                TextFormField(
                  controller: _confirmPasswordController,
                  obscureText: _obscureConfirmPassword,
                  decoration: InputDecoration(
                    labelText: '确认密码',
                    hintText: '请再次输入密码',
                    prefixIcon: const Icon(Icons.lock_outline),
                    suffixIcon: IconButton(
                      icon: Icon(
                        _obscureConfirmPassword ? Icons.visibility : Icons.visibility_off,
                      ),
                      onPressed: () {
                        setState(() {
                          _obscureConfirmPassword = !_obscureConfirmPassword;
                        });
                      },
                    ),
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return '请确认密码';
                    }
                    if (value != _passwordController.text) {
                      return '两次输入的密码不一致';
                    }
                    return null;
                  },
                ),
                
                const SizedBox(height: 24),
                
                // 同意条款
                Row(
                  children: [
                    Checkbox(
                      value: _agreeToTerms,
                      onChanged: (value) {
                        setState(() {
                          _agreeToTerms = value ?? false;
                        });
                      },
                      activeColor: theme.primaryColor,
                    ),
                    Expanded(
                      child: GestureDetector(
                        onTap: () {
                          setState(() {
                            _agreeToTerms = !_agreeToTerms;
                          });
                        },
                        child: RichText(
                          text: TextSpan(
                            style: TextStyle(color: Colors.grey[600], fontSize: 14),
                            children: [
                              const TextSpan(text: '我已阅读并同意'),
                              TextSpan(
                                text: '《用户协议》',
                                style: TextStyle(color: theme.primaryColor),
                              ),
                              const TextSpan(text: '和'),
                              TextSpan(
                                text: '《隐私政策》',
                                style: TextStyle(color: theme.primaryColor),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                
                const SizedBox(height: 32),
                
                // 注册按钮
                SizedBox(
                  height: 50,
                  child: ElevatedButton(
                    onPressed: (authFormState.isLoading || !_agreeToTerms) ? null : _handleRegister,
                    child: authFormState.isLoading
                        ? const CircularProgressIndicator(color: Colors.white)
                        : const Text('注册'),
                  ),
                ),
                
                const SizedBox(height: 16),
                
                // 登录链接
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Text('已有账户？'),
                    TextButton(
                      onPressed: () {
                        context.go('/login');
                      },
                      child: const Text('立即登录'),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}