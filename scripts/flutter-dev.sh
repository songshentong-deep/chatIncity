#!/bin/bash

# Flutter开发服务管理脚本

COMMAND=${1:-"start"}
DEVICE=${2:-"chrome"}

case $COMMAND in
    "start")
        echo "🚀 启动Flutter开发服务..."
        
        # 检查是否在项目根目录
        if [ ! -d "mobile" ]; then
            echo "❌ 请在项目根目录运行此脚本"
            exit 1
        fi
        
        # 进入移动端目录
        cd mobile
        
        # 检查Flutter是否安装
        if ! command -v flutter &> /dev/null; then
            echo "❌ Flutter未安装，请先安装Flutter"
            exit 1
        fi
        
        # 清理缓存
        echo "🧹 清理Flutter缓存..."
        flutter clean
        
        # 获取依赖
        echo "📦 获取Flutter依赖..."
        flutter pub get
        
        # 检查可用设备
        echo "📱 检查可用设备..."
        flutter devices
        
        # 启动开发服务
        echo "🔄 启动Flutter应用 (设备: $DEVICE)..."
        if [ "$DEVICE" = "chrome" ]; then
            flutter run -d chrome --web-port=3001
        else
            flutter run -d $DEVICE
        fi
        ;;
        
    "stop")
        echo "🛑 停止Flutter开发服务..."
        
        # 停止Flutter进程
        FLUTTER_PIDS=$(ps aux | grep "flutter.*run" | grep -v grep | awk '{print $2}')
        
        if [ -n "$FLUTTER_PIDS" ]; then
            echo "找到Flutter进程: $FLUTTER_PIDS"
            kill $FLUTTER_PIDS 2>/dev/null
            sleep 2
            
            # 强制停止仍在运行的进程
            REMAINING=$(ps aux | grep "flutter.*run" | grep -v grep | awk '{print $2}')
            if [ -n "$REMAINING" ]; then
                echo "强制停止剩余Flutter进程: $REMAINING"
                kill -9 $REMAINING 2>/dev/null
            fi
            echo "✅ Flutter服务已停止"
        else
            echo "ℹ️  没有找到运行中的Flutter服务"
        fi
        
        # 停止Chrome进程（如果是Web开发）
        CHROME_FLUTTER=$(ps aux | grep "flutter_tools_chrome_device" | grep -v grep | awk '{print $2}')
        if [ -n "$CHROME_FLUTTER" ]; then
            echo "停止Flutter Chrome进程..."
            kill $CHROME_FLUTTER 2>/dev/null
        fi
        ;;
        
    "restart")
        echo "🔄 重启Flutter开发服务..."
        $0 stop
        sleep 2
        $0 start $DEVICE
        ;;
        
    "build")
        echo "🔨 构建Flutter应用..."
        cd mobile
        
        BUILD_TYPE=${2:-"web"}
        case $BUILD_TYPE in
            "web")
                echo "构建Web版本..."
                flutter build web
                ;;
            "apk")
                echo "构建Android APK..."
                flutter build apk
                ;;
            "ios")
                echo "构建iOS版本..."
                flutter build ios
                ;;
            *)
                echo "未知构建类型: $BUILD_TYPE"
                echo "支持的类型: web, apk, ios"
                exit 1
                ;;
        esac
        ;;
        
    "test")
        echo "🧪 运行Flutter测试..."
        cd mobile
        flutter test
        ;;
        
    "analyze")
        echo "🔍 分析Flutter代码..."
        cd mobile
        flutter analyze
        ;;
        
    "format")
        echo "💅 格式化Flutter代码..."
        cd mobile
        flutter format .
        ;;
        
    "devices")
        echo "📱 显示可用设备..."
        cd mobile
        flutter devices
        ;;
        
    "debug")
        echo "🐛 Flutter调试工具..."
        cd mobile
        
        DEBUG_TYPE=${2:-"info"}
        case $DEBUG_TYPE in
            "info")
                echo "📊 Flutter环境信息:"
                flutter --version
                echo ""
                echo "📱 可用设备:"
                flutter devices
                echo ""
                echo "🔍 项目依赖状态:"
                flutter pub deps
                ;;
            "doctor")
                echo "🩺 Flutter环境诊断:"
                flutter doctor -v
                ;;
            "logs")
                echo "📋 查看Flutter日志..."
                echo "提示: 在Flutter运行时，可以按以下键进行调试:"
                echo "  r - 热重载"
                echo "  R - 热重启"
                echo "  h - 显示帮助"
                echo "  d - 分离调试器"
                echo "  c - 清除屏幕"
                echo "  q - 退出"
                ;;
            "clean")
                echo "🧹 清理Flutter缓存和构建文件..."
                flutter clean
                echo "删除 .dart_tool 目d录..."
                rm -rf .dart_tool
                echo "删除 build 目录..."
                rm -rf build
                echo "重新获取依赖..."
                flutter pub get
                echo "✅ 清理完成"
                ;;
            "inspector")
                echo "🔍 启动Flutter Inspector..."
                flutter inspector
                ;;
            "performance")
                echo "⚡ 性能分析模式启动..."
                flutter run --profile -d ${3:-chrome}
                ;;
            "verbose")
                echo "📝 详细日志模式启动..."
                flutter run --verbose -d ${3:-chrome}
                ;;
            *)
                echo "未知调试类型: $DEBUG_TYPE"
                echo "支持的调试类型:"
                echo "  info        - 显示环境信息"
                echo "  doctor      - 环境诊断"
                echo "  logs        - 日志说明"
                echo "  clean       - 深度清理"
                echo "  inspector   - 启动Inspector"
                echo "  performance - 性能分析模式"
                echo "  verbose     - 详细日志模式"
                exit 1
                ;;
        esac
        ;;
        
    "logs")
        echo "📋 查看Flutter运行日志..."
        cd mobile
        
        LOG_TYPE=${2:-"all"}
        case $LOG_TYPE in
            "all")
                echo "显示所有Flutter进程日志..."
                ps aux | grep flutter | grep -v grep
                ;;
            "crash")
                echo "查看崩溃日志..."
                if [ -d "build/app/outputs/flutter-apk" ]; then
                    find build -name "*.log" -exec echo "=== {} ===" \; -exec cat {} \;
                fi
                ;;
            "web")
                echo "Web开发服务器日志..."
                echo "请在浏览器开发者工具中查看控制台日志"
                ;;
            *)
                echo "未知日志类型: $LOG_TYPE"
                echo "支持的日志类型: all, crash, web"
                ;;
        esac
        ;;
        
    "profile")
        echo "📊 Flutter性能分析..."
        cd mobile
        
        PROFILE_TYPE=${2:-"cpu"}
        case $PROFILE_TYPE in
            "cpu")
                echo "启动CPU性能分析..."
                flutter run --profile --trace-startup -d ${3:-chrome}
                ;;
            "memory")
                echo "启动内存分析..."
                flutter run --profile --enable-software-rendering -d ${3:-chrome}
                ;;
            "size")
                echo "分析应用大小..."
                flutter build apk --analyze-size
                ;;
            *)
                echo "未知分析类型: $PROFILE_TYPE"
                echo "支持的分析类型: cpu, memory, size"
                ;;
        esac
        ;;
        
    "fix")
        echo "🔧 Flutter问题修复工具..."
        cd mobile
        
        FIX_TYPE=${2:-"common"}
        case $FIX_TYPE in
            "common")
                echo "执行常见问题修复..."
                echo "1. 清理缓存..."
                flutter clean
                echo "2. 重新获取依赖..."
                flutter pub get
                echo "3. 修复权限问题..."
                if [ -d "android" ]; then
                    chmod +x android/gradlew
                fi
                echo "4. 检查环境..."
                flutter doctor
                ;;
            "gradle")
                echo "修复Gradle问题..."
                if [ -d "android" ]; then
                    cd android
                    ./gradlew clean
                    cd ..
                fi
                flutter clean
                flutter pub get
                ;;
            "ios")
                echo "修复iOS问题..."
                if [ -d "ios" ]; then
                    cd ios
                    rm -rf Pods
                    rm Podfile.lock
                    pod install
                    cd ..
                fi
                ;;
            "web")
                echo "修复Web问题..."
                flutter clean
                flutter pub get
                flutter build web --release
                ;;
            *)
                echo "未知修复类型: $FIX_TYPE"
                echo "支持的修复类型: common, gradle, ios, web"
                ;;
        esac
        ;;
        
    "help")
        echo "Flutter开发服务管理脚本"
        echo ""
        echo "用法: $0 <command> [options]"
        echo ""
        echo "基本命令:"
        echo "  start [device]     - 启动Flutter开发服务 (默认: chrome)"
        echo "  stop              - 停止Flutter开发服务"
        echo "  restart [device]  - 重启Flutter开发服务"
        echo "  build [type]      - 构建Flutter应用 (web/apk/ios)"
        echo "  test              - 运行测试"
        echo "  analyze           - 代码分析"
        echo "  format            - 代码格式化"
        echo "  devices           - 显示可用设备"
        echo ""
        echo "调试命令:"
        echo "  debug [type]      - 调试工具 (info/doctor/logs/clean/inspector/performance/verbose)"
        echo "  logs [type]       - 查看日志 (all/crash/web)"
        echo "  profile [type]    - 性能分析 (cpu/memory/size)"
        echo "  fix [type]        - 问题修复 (common/gradle/ios/web)"
        echo "  diagnose          - 全面环境诊断"
        echo "  quick [action]    - 快速操作 (reload/restart/clean-start/upgrade)"
        echo ""
        echo "示例:"
        echo "  $0 start chrome           # 在Chrome中启动"
        echo "  $0 debug doctor           # 环境诊断"
        echo "  $0 profile cpu            # CPU性能分析"
        echo "  $0 fix common             # 常见问题修复"
        echo "  $0 debug clean            # 深度清理"
        echo "  $0 logs crash             # 查看崩溃日志"
        echo "  $0 diagnose               # 全面环境诊断"
        echo "  $0 quick clean-start      # 清理并启动"
        ;;
        
    "diagnose")
        echo "🩺 Flutter开发环境全面诊断..."
        cd mobile
        
        echo "=== 1. Flutter环境检查 ==="
        flutter --version
        echo ""
        
        echo "=== 2. 设备连接检查 ==="
        flutter devices
        echo ""
        
        echo "=== 3. 依赖状态检查 ==="
        if [ -f "pubspec.yaml" ]; then
            echo "✅ pubspec.yaml 存在"
            if [ -f "pubspec.lock" ]; then
                echo "✅ pubspec.lock 存在"
            else
                echo "⚠️  pubspec.lock 不存在，需要运行 flutter pub get"
            fi
        else
            echo "❌ pubspec.yaml 不存在"
        fi
        echo ""
        
        echo "=== 4. 构建文件检查 ==="
        if [ -d ".dart_tool" ]; then
            echo "✅ .dart_tool 目录存在"
        else
            echo "⚠️  .dart_tool 目录不存在"
        fi
        
        if [ -d "build" ]; then
            echo "✅ build 目录存在"
            echo "构建文件大小: $(du -sh build 2>/dev/null | cut -f1)"
        else
            echo "ℹ️  build 目录不存在（正常，首次运行会创建）"
        fi
        echo ""
        
        echo "=== 5. 平台特定检查 ==="
        if [ -d "android" ]; then
            echo "✅ Android 平台配置存在"
            if [ -f "android/gradlew" ]; then
                if [ -x "android/gradlew" ]; then
                    echo "✅ Gradle wrapper 可执行"
                else
                    echo "⚠️  Gradle wrapper 不可执行，需要修复权限"
                fi
            fi
        fi
        
        if [ -d "ios" ]; then
            echo "✅ iOS 平台配置存在"
            if [ -f "ios/Podfile" ]; then
                echo "✅ Podfile 存在"
            fi
        fi
        
        if [ -d "web" ]; then
            echo "✅ Web 平台配置存在"
        fi
        echo ""
        
        echo "=== 6. 常见问题检查 ==="
        # 检查端口占用
        if lsof -i :3000 >/dev/null 2>&1; then
            echo "⚠️  端口 3000 被占用"
        else
            echo "✅ 端口 3000 可用"
        fi
        
        # 检查网络连接
        if ping -c 1 pub.dev >/dev/null 2>&1; then
            echo "✅ 网络连接正常 (可访问 pub.dev)"
        else
            echo "⚠️  网络连接问题，无法访问 pub.dev"
        fi
        echo ""
        
        echo "=== 7. 建议操作 ==="
        echo "如果发现问题，可以尝试以下命令："
        echo "  $0 fix common     # 修复常见问题"
        echo "  $0 debug clean    # 深度清理"
        echo "  $0 debug doctor   # 详细环境诊断"
        ;;
        
    "quick")
        echo "⚡ Flutter快速操作..."
        cd mobile
        
        QUICK_ACTION=${2:-"reload"}
        case $QUICK_ACTION in
            "reload")
                echo "🔄 热重载 (如果Flutter正在运行，按 'r' 键)"
                echo "提示: 在Flutter运行的终端中按 'r' 进行热重载"
                ;;
            "restart")
                echo "🔄 热重启 (如果Flutter正在运行，按 'R' 键)"
                echo "提示: 在Flutter运行的终端中按 'R' 进行热重启"
                ;;
            "clean-start")
                echo "🧹 清理并启动..."
                flutter clean
                flutter pub get
                flutter run -d ${3:-chrome}
                ;;
            "upgrade")
                echo "⬆️  升级Flutter和依赖..."
                flutter upgrade
                flutter pub upgrade
                ;;
            *)
                echo "未知快速操作: $QUICK_ACTION"
                echo "支持的操作: reload, restart, clean-start, upgrade"
                ;;
        esac
        ;;
        
    *)
        echo "❌ 未知命令: $COMMAND"
        echo "运行 '$0 help' 查看帮助信息"
        exit 1
        ;;
esac