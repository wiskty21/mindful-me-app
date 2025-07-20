# MindfulMe - 瞑想タイマーアプリ

iOS App Store申請対応の瞑想・マインドフルネスタイマーアプリです。

## 🧘‍♀️ 特徴

- **シンプルで美しいUI**: iOS Human Interface Guidelinesに準拠
- **PWA対応**: App Store申請可能な品質
- **オフライン動作**: ネット接続なしでも基本機能が利用可能
- **データ統計**: セッション履歴と詳細な統計表示
- **呼吸ガイド**: 初心者向けの呼吸リズムガイド
- **iOS最適化**: Safe Area、タッチ操作、バイブレーション対応

## 🛠 技術構成

### フロントエンド
- **Vue.js 3** + Composition API
- **TypeScript** - 型安全性
- **Vite** - 高速ビルドツール
- **Tailwind CSS** - スタイリング
- **PWA Plugin** - オフライン対応
- **Pinia** - 状態管理

### バックエンド
- **Go 1.21** + Gin Framework
- **MongoDB Atlas** - データベース
- **JWT認証** - セキュアな認証
- **CORS対応** - フロントエンド連携

## 🚀 セットアップ

### 前提条件
- Node.js 18+ 
- Go 1.21+
- MongoDB Atlas アカウント

### 1. リポジトリクローン
```bash
git clone [repository-url]
cd iOS-release-test
```

### 2. バックエンド起動
```bash
cd backend

# 環境変数設定
cp .env.example .env
# .envファイルを編集してMongoDB URI等を設定

# 依存関係インストール
go mod tidy

# サーバー起動
go run main.go
```

### 3. フロントエンド起動
```bash
cd frontend

# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev
```

### 4. アクセス
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080

## 📱 PWAとして使用

### iOS Safari
1. SafariでWebアプリにアクセス
2. 共有ボタン → "ホーム画面に追加"
3. ネイティブアプリのように使用可能

### 本格的なiOSアプリ化
App Store申請用の資料は `app-store-preparation.md` を参照してください。

## 🏗 ビルドとデプロイ

### フロントエンド本番ビルド
```bash
cd frontend
npm run build
```

### バックエンド本番ビルド
```bash
cd backend
go build -o bin/mindfulme-server main.go
```

## 📊 API エンドポイント

### 認証
- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン  
- `POST /api/v1/auth/refresh` - トークンリフレッシュ
- `POST /api/v1/auth/logout` - ログアウト

### セッション
- `GET /api/v1/sessions` - セッション履歴取得
- `POST /api/v1/sessions` - セッション記録
- `PUT /api/v1/sessions/:id` - セッション更新
- `GET /api/v1/sessions/stats` - 統計情報
- `GET /api/v1/sessions/streak` - 連続日数

### プロフィール
- `GET /api/v1/profile` - プロフィール取得
- `PUT /api/v1/profile` - プロフィール更新

## 🎨 UI/UX 設計

### カラーパレット
- **Primary**: #4F46E5 (Indigo)
- **Meditation**: #667EEA (Blue-purple gradient)
- **Background**: #F0F4F8 (Light blue-gray)
- **Text**: #1F2937 (Dark gray)

### 主要画面
1. **ログイン/登録** - 美しいグラデーション背景
2. **ホーム** - 統計サマリーとクイックスタート
3. **タイマー** - 円形プログレス、呼吸ガイド
4. **履歴** - セッション一覧、統計グラフ
5. **プロフィール** - 設定管理、データエクスポート

### iOS最適化
- Safe Area対応
- Touch-friendly UI (44pt最小タップ領域)
- ハプティックフィードバック
- ダークモード対応準備
- アクセシビリティ考慮

## 📱 App Store申請準備

### 必要な作業
1. **Apple Developer アカウント** (年間99ドル)
2. **アプリアイコン** - 1024x1024px + 各サイズ
3. **スクリーンショット** - iPhone/iPad各サイズ
4. **プライバシーポリシー** ✅ 作成済み
5. **利用規約** ✅ 作成済み

### App Store情報
- **アプリ名**: MindfulMe - 瞑想タイマー
- **カテゴリ**: ヘルスケア & フィットネス
- **価格**: 無料（フリーミアムモデル検討）
- **対象年齢**: 4+ (制限なし)

詳細は `app-store-preparation.md` を参照してください。

## 🔒 セキュリティ

- HTTPS通信強制
- JWT認証によるAPIアクセス制御
- パスワードのArgon2ハッシュ化
- 入力値検証とサニタイゼーション
- CORS適切な設定

## 📈 今後の改善予定

### 機能追加
- Apple Watch連携
- ヘルスケアアプリ統合
- Siri Shortcuts対応
- Apple Sign In
- プッシュ通知

### プレミアム機能
- ガイド付き瞑想音声
- 詳細分析とインサイト
- カスタムテーマ
- iCloudバックアップ
- 広告非表示

## 🐛 トラブルシューティング

### バックエンド起動エラー
```bash
# MongoDB接続確認
go run main.go
# エラーメッセージで.envのMONGODB_URIを確認
```

### フロントエンドビルドエラー
```bash
# 依存関係再インストール
rm -rf node_modules package-lock.json
npm install
```

### PWA動作しない
- HTTPS必須（本番環境）
- Service Worker登録確認
- manifest.jsonの設定確認

## 📞 サポート

- **バグ報告**: [GitHub Issues]
- **機能要望**: [GitHub Discussions]  
- **メール**: support@mindfulme.app

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。

---

**MindfulMe開発チーム**  
最終更新: 2025-07-20