# MindfulMe 本番デプロイメントガイド

## 推奨構成（初期リリース）

### 1. Vercel（フロントエンド）

```bash
# Vercelインストール
npm i -g vercel

# フロントエンドディレクトリで実行
cd frontend
vercel

# 設定例
# Framework: Vue.js
# Build Command: npm run build
# Output Directory: dist
```

### 2. Railway（バックエンド）

```bash
# Railway CLIインストール
npm i -g @railway/cli

# バックエンドディレクトリで実行
cd backend
railway login
railway link
railway up

# 環境変数設定（Railway ダッシュボード）
# MONGODB_URI=your_mongodb_uri
# JWT_SECRET=your_jwt_secret
# PORT=8080
```

### 3. MongoDB Atlas（データベース）

1. [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)でアカウント作成
2. M0無料クラスター作成
3. 接続文字列取得
4. IPホワイトリスト設定（0.0.0.0/0 for 開発）

## デプロイ手順

### ステップ1: MongoDB Atlas設定
```bash
# 1. クラスター作成後、接続文字列取得
# mongodb+srv://username:password@cluster.mongodb.net/mindfulme

# 2. バックエンドの.envファイルに設定
MONGODB_URI=your_connection_string
```

### ステップ2: バックエンドデプロイ
```bash
cd backend

# Railway にデプロイ
railway up

# デプロイ後のURL確認
railway open
```

### ステップ3: フロントエンドデプロイ
```bash
cd frontend

# API URLを環境変数に設定
echo "VITE_API_URL=https://your-backend.railway.app" > .env.production

# Vercelにデプロイ
vercel --prod
```

## コスト最適化のポイント

1. **初期は無料枠を活用**
   - Vercel: 100GB帯域幅/月
   - Railway: $5無料クレジット
   - MongoDB Atlas: 512MB無料

2. **スケーリング時の移行先**
   - フロント: Cloudflare Pages（無制限）
   - バック: Google Cloud Run（従量課金）
   - DB: MongoDB Atlas M10（専用クラスター）

## 監視・運用

```bash
# ヘルスチェックURL
curl https://your-backend.railway.app/health

# ログ確認
railway logs

# メトリクス
# Vercel Dashboard: Analytics
# Railway Dashboard: Metrics
# MongoDB Atlas: Performance Advisor
```

## セキュリティチェックリスト

- [ ] 本番環境の環境変数設定
- [ ] CORS設定（特定ドメインのみ許可）
- [ ] HTTPS強制
- [ ] レート制限実装
- [ ] MongoDB IPホワイトリスト設定

## トラブルシューティング

### CORS エラー
```go
// backend/middleware/cors.go
AllowOrigins: []string{"https://your-app.vercel.app"}
```

### 接続タイムアウト
- MongoDB Atlas: ネットワークアクセス確認
- Railway: 環境変数確認

### 502 Bad Gateway
- バックエンドのヘルスチェック確認
- ポート設定確認（Railway は PORT 環境変数使用）