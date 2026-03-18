## フロントエンド設計ガイド（AI 向け）

このドキュメントは、**React/Next.js ベースのアプリ**でコンポーネント設計と状態管理方針を統一するためのガイドです。  
新規アプリでも、原則としてこの方針に従って実装してください。

---

## 1. 全体方針

- **責務の分離**: 「見た目を作るコンポーネント」と「データやロジックを扱うコンポーネント」を明確に分ける。
- **変更容易性**: UI の変更は見た目コンポーネントだけを触ればよい、データ取得変更はフックや API レイヤーだけを触ればよい、という構造を目指す。
- **スケーラビリティ**: コンポーネント数や画面数が増えても、命名・配置・責務のルールが崩れないようにする。
- **型安全性**: TypeScript を前提とし、GraphQL や REST の型情報をできる限り自動生成して利用する。

---

## 2. ディレクトリ構成（例）

```text
web/client/
├── src/
│   ├── app/                        # ルーティング（Next.js App Router 前提）
│   │   ├── (public)/...           # 非ログインエリア
│   │   └── (logged-in)/...        # ログイン後エリア
│   ├── features/                  # ドメインごとの UI/ロジック
│   │   ├── articles/
│   │   │   ├── components/        # Presentational Component
│   │   │   ├── hooks/             # 状態管理・データ取得フック
│   │   │   ├── api/               # BFF/バックエンド呼び出し
│   │   │   └── types/             # ドメイン固有の型
│   │   ├── bookmarks/
│   │   └── ...
│   ├── shared/                    # 複数ドメインで使う共通モジュール
│   │   ├── components/            # 共通 UI（ボタン、レイアウト、モーダル等）
│   │   ├── hooks/                 # 共通フック（メディアクエリ、モーダル状態等）
│   │   ├── lib/                   # 汎用ユーティリティ
│   │   └── types/                 # 汎用型
│   ├── graphql/                   # GraphQL クエリ/型（Codegen 前提）
│   └── lib/                       # ルーティング/環境変数/クライアントなど
└── ...
```

---

## 3. コンポーネント設計方針

### 3.1 コンポーネントの種類

1. **Presentational Component（見た目担当）**
   - **役割**
     - HTML / CSS / UI ライブラリ（shadcn/ui 等）を使った「表示」に専念する。
     - データ取得・ビジネスロジック・副作用を持たない。
   - **入力**
     - Props として「すでに用意されたデータ」と「イベントハンドラ（コールバック）」だけを受け取る。
   - **配置**
     - 各 feature の `components/` 配下（例: `features/articles/components/ArticleList.tsx` など）。
   - **例**
     - `ArticleCard`, `BookmarkList`, `Sidebar`, `DialogContent` など。

2. **Container / Template Component（画面の組み立て＋データ受け取り）**
   - **役割**
     - 複数の Presentational Component を組み合わせて「画面」を構成する。
     - データ取得フックやミューテーションフックを呼び出し、その結果を Presentational に渡す。
   - **入力**
     - ルーティングパラメータやグローバル状態から必要な情報を取得する。
   - **配置**
     - `features/<domain>/components/Template/` など。
   - **例**
     - `ArticleListTemplate`, `BookmarkTemplate`, `FavoriteArticleFolderListTemplate` など。

3. **Page Component（ルーティング境界）**
   - **役割**
     - Next.js の `app/` ディレクトリ配下の `page` / `layout` として、URL パラメータやメタ情報を処理し Template に渡すだけにする。
   - **配置**
     - 例: `src/app/(logged-in)/articles/page.tsx` など。
   - **原則**
     - ここにはビジネスロジックを書かない。Template/Container に委譲する。

### 3.2 責務分離のルール

- **Presentational**
  - ビジネス用語を知らなくても理解できるレベルの「UI の見た目」だけにする。
  - API クライアント・GraphQL hook・`useEffect` でのデータ取得を使わない。
- **Container / Template**
  - データ取得・変換・ハンドラ作成などのロジックをまとめる。
  - 画面ごとの「状態管理のハブ」になる。
- **Page**
  - URL 〜 Template の橋渡しのみ。将来、別フレームワークに移行するときの境界としても機能する。

---

## 4. 状態管理方針

### 4.1 原則

1. **サーバーソースのデータは「サーバーキャッシュ or データフェッチ専用フック」で扱う**
   - GraphQL: Apollo Client / urql 等の hooks を利用する。
   - REST: React Query / SWR などを利用する。
   - 画面ローカルの `useState` でサーバーデータを全部抱え込まない。

2. **UI 状態（モーダルの開閉、タブ選択など）はコンポーネントローカル or Context に閉じ込める**
   - ドメインロジックと無関係な UI 状態は、グローバルストアに乗せない。

3. **グローバル状態は最小限に**
   - 認証情報（ログインユーザー）
   - テーマ（ダークモード）や言語設定
   - 画面間で共有する必要のある一部のフィルタ条件などだけをグローバルストア（Context / Redux / Zustand 等）で扱う。

### 4.2 推奨レイヤー

- **データ取得（サーバー状態）**
  - GraphQL の場合:
    - `graphql/` にクエリと型を定義（Codegen で TS 型を生成）。
    - `features/<domain>/hooks/` 内で、クエリごとのカスタムフックを作る。
      - 例: `useArticleListQuery`, `useBookmarkListQuery`。
  - REST の場合:
    - `features/<domain>/api/` にエンドポイントを定義。
    - `features/<domain>/hooks/` に React Query / SWR をラップしたカスタムフックを置く。

- **ローカル UI 状態**
  - モーダル開閉、ドロワー開閉、選択中タブなどは、Template または Presentational コンポーネント内で `useState` / `useReducer` を使って管理する。

- **グローバル状態**
  - 認証: Supabase Auth / NextAuth などの仕組みを使い、Context や専用 provider コンポーネントでラップする。
  - 認証済みユーザー情報へのアクセスは、`useCurrentUser()` のような共通フックから行う。

---

## 5. API／型の扱い方針

- GraphQL を利用する場合は、必ず Codegen で型を生成し、手書きの型を極力減らす。
- フロント側では、**「GraphQL 型 → ドメイン用の軽量な型」**への変換関数を用意しても良い（API 変更に強くするため）。
  - API 呼び出しは、`features/<domain>/api` または hooks に閉じ込める。
  - UI コンポーネント内で `fetch` / `axios` / `client.query` を直接呼ばない。

---

## 6. 再利用性・デザインシステム

- 共通 UI コンポーネント（ボタン、モーダル、カード枠など）は `shared/components/ui/` に集約する。
- ドメイン固有の文言・アイコン・スタイルが入ったものは、各 `features/<domain>/components/` に置く。
- 共通レイアウト（ヘッダー、サイドバー、フッターなど）は `shared/components/layout/` にまとめる。

---

## 7. AI への短い指示（コピペ用）

このフロントエンドでは、`docs/FRONTEND_DESIGN_GUIDE.md` に記載されたルールに従って、

- コンポーネントを「Presentational / Container(Template) / Page」の 3 層に分割し、
- サーバー状態は GraphQL / REST の専用フックで扱い、
- グローバル状態は認証など最小限にとどめてください。  
  各機能は `features/<domain>/` 単位で切り出し、UI 部品は `components/`、データ取得は `hooks/` と `api/` に整理してください。
