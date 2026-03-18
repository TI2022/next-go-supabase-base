## MVP 雛形アプリ仕様（ログイン＋ホーム画面）

このプロジェクトでは、将来の拡張性を意識しつつ、まずは **ログイン機能＋ホーム画面のみ**の MVP を作成してください。  
バックエンド・フロントエンドともに、他機能へ拡張しやすい構成と責務分離を重視します。

---

## 1. 全体構成（MVP）

- リポジトリはモノレポ構成の土台だけ用意し、実装サービスは最小限にとどめる：

```text
repo-root/
├── .env                      # ローカル環境変数
├── compose.yml               # DB や app-service を起動するために必要なら用意
├── web/
│   └── client/               # Next.js フロント
├── app-service/              # Go バックエンド（1 サービスのみ）
├── db/
│   └── migrations/           # DB マイグレーション
└── docs/
    ├── REPO_ARCHITECTURE.md  # アーキテクチャ方針（別途）
    └── MVP_LOGIN_HOME_SPEC.md（本ファイル）
```

- MVP で実装する機能は **2 つだけ**:
  1. ログイン（認証）機能
  2. ログイン後に表示されるホーム画面（簡単な「ようこそ」＋プロフィール or ダミーデータ一覧）

---

## 2. 認証・ログインの要件

- **認証方式**: メール＋パスワード、または OAuth（例: Google）のどちらか一方でよい（どちらでも実装しやすい方を採用）。
- **ユーザー情報の保存**
  - DB（Postgres / Supabase）に保存し、以下のような `users` テーブルを想定：
    - `id`（UUID or bigint）
    - `email`（unique）
    - `name`
    - `created_at` など

- **フロントエンドの要件**
  - `/login` ページでログインフォームを提供する。
  - ログイン成功時は `/`（ホーム画面）へリダイレクトする。
  - 認証済みでないユーザーが `/` にアクセスしたときは `/login` に飛ばす。

- **バックエンドの要件**
  - 認証 API エンドポイントを 1 つ提供する（例: `POST /auth/login`）。
  - セッションは JWT ベースでも HTTP-only Cookie ベースでもよいが、将来の拡張性を考慮してサーバ側で検証しやすい形にする。
  - ログイン API は成功時にトークン or セッションクッキーを返し、以降は `Authorization` ヘッダ or Cookie 経由で認証状態を維持する。

---

## 3. バックエンド（app-service）の要件

### 3.1 言語・構成

- **言語**: Go
- **ディレクトリ構造**（MVP でも将来のサービス分割に備えて層分けする）:

```text
app-service/
├── cmd/
│   └── main.go
└── internal/
    ├── domain/
    │   ├── entity/            # User などドメインモデル
    │   └── repository/        # UserRepository interface
    ├── application/
    │   └── usecase/           # LoginUsecase, GetCurrentUserUsecase など
    ├── adapter/
    │   ├── persistence_adapter/
    │   └── external_adapter/  # MVP では空でもよい
    ├── infrastructure/
    │   ├── persistence/       # UserRepository の実装（DB アクセス）
    │   └── external/          # MVP では空でもよい
    └── interfaces/handler/    # HTTP ハンドラ（/auth/login, /me 等）
```

- **依存方向のルール（重要）**
  - `handler` → `usecase` → `adapter`（interface） → `infrastructure`
  - `domain` は外部に依存しない。

### 3.2 必須エンドポイント

- `POST /auth/login`
  - **入力**: `email` + `password`（または OAuth コード）
  - **処理**: 認証に成功したらトークン or セッションを発行し、クライアントに返す

- `GET /me`
  - **入力**: 認証済みトークン / セッション
  - **出力**: ログイン中ユーザーのプロフィール（`id`, `email`, `name` など）

- `GET /health`
  - 単純なヘルスチェック用（200 OK を返すだけ）

---

## 4. フロントエンド（web/client）の要件

### 4.1 技術スタック

- Next.js (App Router), TypeScript
- UI ライブラリは自由（shadcn/ui や Chakra UI など）だが、デザインはシンプルで構わない。可読性優先。

### 4.2 コンポーネント設計（最小形）

- **Page**
  - `/login/page.tsx`: ログインフォームページ
  - `/page.tsx`: ホーム画面（「ようこそ, {ユーザー名}」＋簡単なカードリストなど）

- **Template / Container**
  - `features/auth/components/LoginTemplate.tsx`
  - `features/home/components/HomeTemplate.tsx`

- **Presentational**
  - `features/auth/components/LoginForm.tsx`
  - `features/home/components/WelcomeCard.tsx` など

### 4.3 状態管理方針（MVP）

- **認証状態**
  - ログイン後に返されるトークン or Cookie をもとに、`GET /me` を叩いてユーザー情報を取得する。
  - `useCurrentUser()` のようなカスタムフックで、現在ユーザーを取得するようにする。

- **ページ表示**
  - `/page.tsx` では `useCurrentUser()` を使い、未ログインなら `/login` にリダイレクト、ログイン済みならホーム UI を表示。

---

## 5. データベース

- **テーブル例**: `users`
  - `id`（UUID or bigint）
  - `email`（unique）
  - `password_hash`（パスワード方式の場合）
  - `name`
  - `created_at`, `updated_at`

- マイグレーションファイルを `db/migrations/` に作成し、スキーマ変更は必ずマイグレーション経由で行う。

---

## 6. 開発・起動手順（MVP）

1. DB 起動（Postgres or Supabase）
2. マイグレーション適用
3. `app-service` 起動（例: `go run ./cmd/main.go`）
4. `web/client` 起動（`npm run dev`）
5. ブラウザで `http://localhost:3000/login` にアクセスし、ログイン → ホーム画面表示まで確認できること

---

## 7. 実装上の注意（AI への特記事項）

- バックエンドのビジネスロジックは必ず usecase 層に書き、handler や infrastructure に直接書かないこと。
- フロントエンドでは、API 呼び出しは `features/auth/api`, `features/home/api` などにまとめ、コンポーネント内で `fetch` や `axios` を直接多用しないこと。
- 認証処理は MVP なので堅牢すぎる実装は不要だが、平文パスワードの保存はしない、など最低限の安全性は守ること。
- 将来、マイクロサービスや BFF を追加できるよう、ディレクトリ構成と依存方向だけは崩さないこと。

---

この `.md` を新規リポジトリの `docs/` に置き、「この仕様どおりに MVP を作って」と AI に渡せば、ログイン＋ホームだけの雛形を、拡張しやすい形で実装させやすくなります。
