import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/auth/presentation/pages/register_page.dart';
import '../../features/home/presentation/pages/home_page.dart';
import '../../features/match/presentation/pages/match_page.dart';
import '../../features/chat/presentation/pages/chat_list_page.dart';
import '../../features/chat/presentation/pages/chat_page.dart';
import '../../features/profile/presentation/pages/profile_page.dart';
import '../../features/profile/presentation/pages/edit_profile_page.dart';
import '../providers/auth_provider.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);
  
  return GoRouter(
    initialLocation: '/login',
    redirect: (context, state) {
      final isLoggedIn = authState.isAuthenticated;
      final currentPath = state.uri.path;
      final isLoggingIn = currentPath == '/login' || currentPath == '/register';
      
      if (!isLoggedIn && !isLoggingIn) {
        return '/login';
      }
      if (isLoggedIn && isLoggingIn) {
        return '/home';
      }
      return null;
    },
    routes: [
      // 认证相关路由
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterPage(),
      ),
      
      // 主应用路由
      ShellRoute(
        builder: (context, state, child) => MainShell(child: child),
        routes: [
          GoRoute(
            path: '/home',
            builder: (context, state) => const HomePage(),
          ),
          GoRoute(
            path: '/match',
            builder: (context, state) => const MatchPage(),
          ),
          GoRoute(
            path: '/chat',
            builder: (context, state) => const ChatListPage(),
            routes: [
              GoRoute(
                path: ':chatId',
                builder: (context, state) {
                  final chatId = state.pathParameters['chatId']!;
                  final chatName = state.uri.queryParameters['name'];
                  return ChatPage(
                    chatId: chatId,
                    chatName: chatName,
                  );
                },
              ),
            ],
          ),
          GoRoute(
            path: '/profile',
            builder: (context, state) => const ProfilePage(),
            routes: [
              GoRoute(
                path: 'edit',
                builder: (context, state) => const EditProfilePage(),
              ),
            ],
          ),
        ],
      ),
    ],
  );
});

// 主应用Shell，包含底部导航栏
class MainShell extends StatelessWidget {
  final Widget child;
  
  const MainShell({Key? key, required this.child}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: child,
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        selectedItemColor: Theme.of(context).primaryColor,
        unselectedItemColor: Colors.grey,
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            label: '首页',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.favorite),
            label: '匹配',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.chat),
            label: '聊天',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person),
            label: '我的',
          ),
        ],
        onTap: (index) {
          switch (index) {
            case 0:
              context.go('/home');
              break;
            case 1:
              context.go('/match');
              break;
            case 2:
              context.go('/chat');
              break;
            case 3:
              context.go('/profile');
              break;
          }
        },
      ),
    );
  }
}