#!/bin/bash

# MindfulMe - デプロイスクリプト
# iOS App Store申請準備用

set -e

echo "🧘‍♀️ MindfulMe デプロイスクリプト開始"
echo "========================================"

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ログ関数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 環境確認
check_requirements() {
    log_info "環境要件を確認中..."
    
    if ! command -v node &> /dev/null; then
        log_error "Node.js がインストールされていません"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        log_error "Go がインストールされていません"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        log_error "npm がインストールされていません"
        exit 1
    fi
    
    NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
    if [ "$NODE_VERSION" -lt 18 ]; then
        log_error "Node.js 18以上が必要です (現在: $(node --version))"
        exit 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2 | cut -d'.' -f1,2)
    if [ "$(echo "$GO_VERSION >= 1.21" | bc)" -eq 0 ]; then
        log_error "Go 1.21以上が必要です (現在: $(go version))"
        exit 1
    fi
    
    log_success "環境要件OK"
}

# バックエンドビルド
build_backend() {
    log_info "バックエンドをビルド中..."
    
    cd backend
    
    # 依存関係確認
    if [ ! -f "go.mod" ]; then
        log_error "go.mod が見つかりません"
        exit 1
    fi
    
    # .envファイル確認
    if [ ! -f ".env" ]; then
        log_warning ".env ファイルが見つかりません"
        log_info ".env.example をコピーして設定してください"
        cp .env.example .env
    fi
    
    # 依存関係インストール
    go mod tidy
    
    # テスト実行
    log_info "テスト実行中..."
    if go test ./... -v; then
        log_success "テスト成功"
    else
        log_warning "テストが失敗しましたが、ビルドを継続します"
    fi
    
    # ビルド
    mkdir -p bin
    
    # 本番用ビルド
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/mindfulme-server-linux main.go
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/mindfulme-server-darwin main.go
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/mindfulme-server-windows.exe main.go
    
    cd ..
    log_success "バックエンドビルド完了"
}

# フロントエンドビルド
build_frontend() {
    log_info "フロントエンドをビルド中..."
    
    cd frontend
    
    # package.json確認
    if [ ! -f "package.json" ]; then
        log_error "package.json が見つかりません"
        exit 1
    fi
    
    # 依存関係インストール
    log_info "依存関係をインストール中..."
    npm ci
    
    # 型チェック
    if command -v npx &> /dev/null; then
        log_info "TypeScript型チェック中..."
        if npx vue-tsc --noEmit; then
            log_success "型チェック成功"
        else
            log_warning "型チェックで警告がありますが、ビルドを継続します"
        fi
    fi
    
    # リンターチェック
    if npm run lint &> /dev/null; then
        log_success "リンターチェック成功"
    else
        log_warning "リンターで警告がありますが、ビルドを継続します"
    fi
    
    # 本番ビルド
    log_info "本番用ビルド実行中..."
    npm run build
    
    # ビルド結果確認
    if [ -d "dist" ]; then
        log_success "フロントエンドビルド完了"
        log_info "ビルドサイズ:"
        du -sh dist/*
    else
        log_error "ビルドが失敗しました"
        exit 1
    fi
    
    cd ..
}

# PWA確認
check_pwa() {
    log_info "PWA設定を確認中..."
    
    cd frontend/dist
    
    # 必要ファイル確認
    files=("manifest.webmanifest" "sw.js" "index.html")
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            log_success "$file 存在確認"
        else
            log_error "$file が見つかりません"
            exit 1
        fi
    done
    
    # アイコンファイル確認
    if [ -d "icons" ]; then
        icon_count=$(find icons -name "*.png" | wc -l)
        log_info "アイコンファイル数: $icon_count"
    else
        log_warning "iconsディレクトリが見つかりません"
    fi
    
    cd ../..
    log_success "PWA設定確認完了"
}

# iOS App Store 準備確認
check_app_store_readiness() {
    log_info "App Store申請準備状況を確認中..."
    
    # 必要ファイル確認
    files=(
        "app-store-preparation.md"
        "privacy-policy.md" 
        "terms-of-service.md"
        "README.md"
    )
    
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            log_success "$file 存在確認"
        else
            log_warning "$file が見つかりません"
        fi
    done
    
    log_info "App Store申請のために以下が必要です:"
    echo "  🔧 Apple Developer アカウント (年間$99)"
    echo "  🎨 アプリアイコン (1024x1024px + 各サイズ)"
    echo "  📱 スクリーンショット (iPhone/iPad各サイズ)"
    echo "  📄 詳細は app-store-preparation.md を参照"
}

# パフォーマンステスト
performance_test() {
    log_info "パフォーマンステストを実行中..."
    
    # バックエンドパフォーマンス
    cd backend
    if [ -f "bin/mindfulme-server-darwin" ]; then
        log_info "バックエンドサーバーサイズ:"
        ls -lh bin/mindfulme-server-*
    fi
    cd ..
    
    # フロントエンドパフォーマンス
    cd frontend
    if [ -d "dist" ]; then
        log_info "フロントエンドバンドルサイズ:"
        find dist -name "*.js" -o -name "*.css" | head -10 | xargs ls -lh
        
        # 合計サイズ
        total_size=$(du -sh dist | cut -f1)
        log_info "総ビルドサイズ: $total_size"
    fi
    cd ..
}

# セキュリティチェック
security_check() {
    log_info "セキュリティチェック実行中..."
    
    # npm audit (フロントエンド)
    cd frontend
    if npm audit --audit-level=high &> /dev/null; then
        log_success "npm audit - 重大な脆弱性なし"
    else
        log_warning "npm audit - 脆弱性が検出されました"
        log_info "詳細: npm audit"
    fi
    cd ..
    
    # Go security check
    cd backend
    if command -v gosec &> /dev/null; then
        log_info "gosec セキュリティスキャン実行中..."
        if gosec ./... &> /dev/null; then
            log_success "gosec - セキュリティ問題なし"
        else
            log_warning "gosec - 潜在的なセキュリティ問題が検出されました"
        fi
    else
        log_info "gosec がインストールされていません (推奨: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)"
    fi
    cd ..
}

# デプロイ成功メッセージ
deployment_success() {
    log_success "🎉 デプロイ完了!"
    echo ""
    echo "=========================================="
    echo "🧘‍♀️ MindfulMe - 瞑想タイマーアプリ"
    echo "=========================================="
    echo ""
    echo "📁 ビルド成果物:"
    echo "  Backend: backend/bin/"
    echo "  Frontend: frontend/dist/"
    echo ""
    echo "🚀 次のステップ:"
    echo "  1. MongoDB Atlas の設定"
    echo "  2. 本番サーバーへのデプロイ"
    echo "  3. ドメイン設定とHTTPS化"
    echo "  4. App Store申請準備"
    echo ""
    echo "📖 詳細情報:"
    echo "  README.md - セットアップガイド"
    echo "  app-store-preparation.md - App Store申請ガイド"
    echo ""
    echo "✨ 素晴らしい瞑想アプリの完成です！"
}

# メイン実行
main() {
    echo "デプロイ対象を選択してください:"
    echo "1) フルビルド (バックエンド + フロントエンド)"
    echo "2) バックエンドのみ"
    echo "3) フロントエンドのみ"
    echo "4) App Store準備確認のみ"
    read -p "選択 (1-4): " choice
    
    case $choice in
        1)
            check_requirements
            build_backend
            build_frontend
            check_pwa
            performance_test
            security_check
            check_app_store_readiness
            deployment_success
            ;;
        2)
            check_requirements
            build_backend
            ;;
        3)
            check_requirements
            build_frontend
            check_pwa
            ;;
        4)
            check_app_store_readiness
            ;;
        *)
            log_error "無効な選択です"
            exit 1
            ;;
    esac
}

# スクリプト実行
main "$@"