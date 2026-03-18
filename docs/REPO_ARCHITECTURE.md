# リポジトリ構成ガイド（モノレポ）

このドキュメントは、**責務の分離・変更容易性・拡張性**を前提にしたリポジトリ構成の意図とルールを定義する。  
他のプロジェクトに同じ構成を適用する際や、AI に「この構成で作って／直して」と指示するときの参照として使う。

---

## 1. 全体の意図

- **モノレポ**: フロント・BFF・マイクロサービス・バッチ・DB 定義を 1 リポジトリで管理し、変更の一貫性とデプロイ単位の切り分けを両立する。
- **レイヤー分離**: 各サービス内で「ドメイン → ユースケース → アダプタ → インフラ」を分け、**ビジネスロジックが DB・外部 API・プロトコルに依存しない**ようにする。
- **単一の API 境界**: クライアントは BFF（GraphQL）だけを叩き、BFF がマイクロサービス（gRPC）に振り分ける。クライアントが gRPC を直接知る必要はない。
- **環境の一元化**: ローカルはルートの `.env` と Supabase + Docker Compose で再現し、各サービスは必要最小限の環境変数だけを持つ。

---

## 2. ディレクトリ構成と責務

```
（ルート）
├── .env                    # ローカル用の環境変数（Compose・開発時の参照元）
├── compose.yml             # BFF + マイクロサービスの起動・ネットワーク
├── Makefile                # よく使う操作のエイリアス（supabase, compose, コンテナ ssh）
├── web/                    # クライアントアプリ群
│   ├── client-v2/          # メイン Web フロント（Next.js, GraphQL, Supabase Auth）
│   ├── client/             # 別版 or 旧版クライアント
│   └── admin/              # 管理画面
├── bff/                    # Backend for Frontend
│   └── apollo-gateway/     # GraphQL API（NestJS, Apollo）, gRPC クライアント
├── micro-service/          # ドメイン別 Go マイクロサービス（gRPC サーバ）
│   ├── content-service/
│   ├── bookmark-service/
│   ├── my-feed-service/
│   ├── favorite-service/
│   └── user-service/
├── batch-service/          # バッチ・クローラー・シード（Go）
├── supabase/               # DB スキーマ（マイグレーション）とローカル設定
└── docs/                   # 本ドキュメントなど
```

| レイヤー | 役割 | 変更が及ぶ範囲 |
|----------|------|----------------|
| **web/** | UI・UX。BFF の GraphQL と Supabase Auth のみ利用。 | 画面・ルーティング・フロント状態 |
| **bff/** | 認証・認可・集約。GraphQL スキーマとリゾルバ、gRPC 呼び出し。 | API 契約・レスポンス形・他サービス呼び出し |
| **micro-service/** | ドメインごとの永続化・他サービス連携。gRPC で提供。 | 当該ドメインのロジック・テーブル・他サービス連携 |
| **batch-service/** | データ投入・クロール・集計。DB 直接 or 内部 API。 | シード・クロール仕様・新ジョブ追加 |
| **supabase/** | スキーマの真実の源。全サービスが参照する DB 定義。 | テーブル・カラム・RLS・関数 |

---

## 3. 責務の分離（詳細）

### 3.1 クライアント（web/client-v2 等）

- **やること**: GraphQL クエリ/ミューテーション、Supabase Auth（ログイン・セッション）、UI。
- **やらないこと**: 他サービスへの直アクセス、DB 直叩き、gRPC の知識。
- **変更容易性**: BFF のスキーマが変わったら型・クエリを合わせて変更。認証だけ Supabase に依存。

### 3.2 BFF（bff/apollo-gateway）

- **やること**:
  - GraphQL スキーマの定義（`schema/<ドメイン>/`）。
  - リゾルバで gRPC クライアントを呼び出し、レスポンスを GraphQL に変換。
  - JWT 検証（Supabase）と CORS など API 境界の設定。
- **やらないこと**: ビジネスロジックの実装、DB 直接アクセス。
- **変更容易性**: 新フィールド・新リソースは「スキーマ追加 → リゾルバ → 既存 gRPC または新 gRPC 呼び出し」。

構成の目安:

- `src/app/<ドメイン>/` … リゾルバ・サービス（例: content, bookmark, favorite）。
- `src/app/grpc/` … 各マイクロサービス向け gRPC クライアント。
- `src/schema/<ドメイン>/` … そのドメインの GraphQL（schema.graphql, query.graphql, mutation.graphql）。

### 3.3 マイクロサービス（micro-service/*-service）

各サービスは **Clean Architecture / ヘキサゴナル** に近い層分けをしている。

- **cmd/main.go**: 起動のみ。設定読み込み・DI・サーバ起動。
- **internal/domain/**: エンティティ・リポジトリインターフェース・ドメイン型。**インフラに依存しない**。
- **internal/application/usecase/**: ユースケース。リポジトリ interface 経由で永続化・他サービス呼び出し。
- **internal/adapter/persistence_adapter/**: ユースケースから使う「永続化の入口」。repository を呼ぶ。
- **internal/adapter/external_adapter/**: 他マイクロサービス呼び出し（gRPC クライアント）のラッパー。
- **internal/infrastructure/persistence/**: 実際の DB アクセス（sqlboiler 等）。
- **internal/infrastructure/external/**: 他サービスへの gRPC 呼び出し実装。
- **internal/interfacess/handler/**: gRPC ハンドラ。リクエストを受け、ユースケースを呼ぶ。

依存の向き:  
`handler` → `usecase` → `adapter`（interface）  
実装は `infrastructure` にあり、`adapter` がそれを使う。

- **変更容易性**: DB スキーマ変更は entity（sqlboiler 再生成）と repository 実装。他サービス連携変更は external_adapter + infrastructure/external。ビジネスルールは usecase のみ。

### 3.4 バッチ（batch-service）

- **cmd/<ジョブ名>/**: エントリポイントとそのジョブ専用ユースケース。
  - 例: `cmd/migrate-seed`, `cmd/trend-article-crawler`, `cmd/article-company-crawler`。
- **共有**: `database/`（DB 接続）, `entity/`（sqlboiler）, `domain/`, `infrastructure/`（API・RSS・Supabase 等）。
- **責務**: 「いつ・何を実行するか」は cmd、「どう取得・どう保存するか」は usecase と infrastructure。
- **変更容易性**: 新ジョブは `cmd/<新ジョブ>/main.go` と usecase を追加。既存ジョブの挙動変更は該当 usecase と infrastructure。

### 3.5 Supabase（supabase/）

- **役割**: マイグレーション（`supabase/migrations/`）がスキーマの唯一の定義。BFF・マイクロサービス・バッチはすべてこの DB を参照。
- **変更容易性**: テーブル追加・変更はマイグレーション追加 → `supabase db reset` または `migrate up`。他レイヤーは entity/型を合わせて更新。

---

## 4. 変更容易性のチェックリスト

| 変更内容 | 触る場所 | 注意点 |
|----------|----------|--------|
| 新テーブル・カラム | supabase/migrations, 各サービスの entity（sqlboiler 再生成） | マイグレーション順序・既存データ |
| 新 GraphQL フィールド・型 | bff schema + リゾルバ、必要なら gRPC | クライアントのクエリ・型 |
| 新 gRPC メソッド | protocol-buffers リポジトリ、該当サービスの handler + usecase | BFF の gRPC クライアント呼び出し追加 |
| 新マイクロサービス | micro-service/<新サービス>, compose.yml, BFF にモジュール・gRPC クライアント | ネットワーク・環境変数 |
| 新バッチジョブ | batch-service/cmd/<新ジョブ>, 必要なら infrastructure | .env（ローカルは host: 127.0.0.1, port: 54322 等） |
| 認証・認可の変更 | BFF（ガード・コンテキスト）, Supabase Auth 設定 | クライアントのトークン送り方 |

---

## 5. 拡張性のルール

- **新クライアントを足す**: `web/<新クライアント>/` を追加。BFF の GraphQL と Supabase のみ利用。必要なら CORS に Origin 追加。
- **新マイクロサービスを足す**:  
  - `micro-service/<新サービス>/` を content-service 等の構造に倣って作成。  
  - `compose.yml` にサービス追加。  
  - BFF に `schema/<ドメイン>/`, `app/<ドメイン>/`, `app/grpc/grpc-<新サービス>-client.service.ts` を追加。
- **新バッチを足す**: `batch-service/cmd/<新ジョブ>/main.go` と usecase。DB は既存 `database.Init()` と entity を流用可能。
- **新ドメイン（新リソース）を BFF に足す**: `schema/<ドメイン>/` を作り、リゾルバで既存 or 新規 gRPC を呼ぶ。

---

## 6. AI 向け指示テンプレート

別プロジェクトで同じ構成を採用するとき、AI に次のように指示できる。

```
- リポジトリ構成は docs/REPO_ARCHITECTURE.md に従う（本ドキュメントで定義したモノレポ構成）。
- 責務: web は GraphQL と認証のみ。BFF は GraphQL と gRPC 集約。マイクロサービスはドメインごとに Go で gRPC。バッチは cmd ごとにエントリポイントを分ける。
- 変更容易性: スキーマ変更は supabase migrations。API 契約は BFF の schema と gRPC。ビジネスロジックは usecase にのみ書く。
- 拡張: 新サービスは micro-service/<名>/ を既存サービスと同じ internal 構成で。新ジョブは batch-service/cmd/<名>/。新 API は BFF の schema とリゾルバから。
- 環境: ルート .env を Compose が参照。各サービスは必要なら自ディレクトリに .env（ローカル用）。batch をホストから実行するときは DB は 127.0.0.1:54322 などホストから見えるアドレスにする。
```

---

## 7. 関連ファイル一覧（参照用）

| ファイル | 役割 |
|----------|------|
| `compose.yml` | BFF と各マイクロサービスの定義・ネットワーク・環境変数 |
| `Makefile` | supabase-start/stop, dcu/dcd/dcb, コンテナ ssh |
| `.env` | ローカル用の BFF_PORT, POSTGRES_*, *_CONTAINER_NAME 等 |
| `supabase/config.toml` | ローカル Supabase のポート・オプション |
| `bff/apollo-gateway/src/app/grpc/*.ts` | 各マイクロサービスへの gRPC 呼び出し |
| `micro-service/*/internal/` | ドメイン・usecase・adapter・infrastructure・handler |
| `batch-service/cmd/*/main.go` | 各バッチのエントリポイント |
| `batch-service/.env` | ローカル実行時は POSTGRES_HOST=127.0.0.1, POSTGRES_PORT=54322 を推奨 |

---

このドキュメントをプロジェクトルートの `docs/REPO_ARCHITECTURE.md` として配置し、README や CONTRIBUTING からリンクしておくと、人間・AI の両方が「どこをどう変えるか」を揃えやすい。
