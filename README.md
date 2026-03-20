## 🙋 About Tatsunori Iijima
I'm a web and mobile developer, working mostly with Ruby on Rails, JavaScript! 

### 🌱 my skills

[![My Skills](https://skillicons.dev/icons?i=ts,js,python,dart,php,nodejs,html,css,sass,tailwind,react,nextjs,vue,nuxtjs,express,nestjs,flutter,docker,aws,gcp,graphql,mysql,postgres,firebase,supabase,prisma,jest,npm,yarn,webpack&perline=10)](https://skillicons.dev)

### 🍪 likes
* Play tennis
* Drive my car

### 📩 contact me
*  [linkdin](https://www.linkedin.com/in/%E8%BE%B0%E5%89%87-%E9%A3%AF%E5%B3%B6-88953a34a/)





## 📈 Status

<img alt="Top Langs" height="150px" src="https://github-readme-stats.vercel.app/api/top-langs/?username=TI2022&layout=compact&count_private=true&show_icons=true&theme=tokyonight" />          <img alt="github stats" height="150px" src="https://github-readme-stats.vercel.app/api?username=TI2022&count_private=true&show_icons=true&show_icons=true&theme=tokyonight" />

[![trophy](https://github-profile-trophy.vercel.app/?username=TI2022&theme=onedark&column=7)](https://github.com/ryo-ma/github-profile-trophy)

<!-- [![Anurag's GitHub stats](https://github-readme-stats.vercel.app/api?username=YukiOnishi1129&theme=onedark)](https://github.com/anuraghazra/github-readme-stats)


[![Top Langs](https://github-readme-stats.vercel.app/api/top-langs/?username=YukiOnishi1129&theme=github_dark&layout=compact
)](https://github.com/anuraghazra/github-readme-stats) -->


<!-- 
<a href="https://app.daily.dev/yuki"><img src="https://api.daily.dev/devcards/v2/IytwLEYk5PX0HyTXp5pEg.png?type=default&r=9ex" width="356" alt="yuki's Dev Card"/></a>
 -->

 <!-- 
 ![](https://github-profile-summary-cards.vercel.app/api/cards/profile-details?username=YukiOnishi1129&theme=2077)
  -->







<!-- ### Hi there 👋 -->

<!--
**YukiOnishi1129/YukiOnishi1129** is a ✨ _special_ ✨ repository because its `README.md` (this file) appears on your GitHub profile.

Here are some ideas to get you started:

- 🔭 I’m currently working on ...
- 🌱 I’m currently learning ...
- 👯 I’m looking to collaborate on ...
- 🤔 I’m looking for help with ...
- 💬 Ask me about ...
- 📫 How to reach me: ...
- 😄 Pronouns: ...
- ⚡ Fun fact: ...
--># next-go-supabase-base
# next-go-supabase-base

新規アプリのモノレポ（フロント + app-service + DB）のリポジトリです。

## MVP 起動手順（Login + Home）

前提:
- Supabase DB コンテナが起動していること
- ルート `.env` が設定済みであること（`POSTGRES_*`）

### 1. DB マイグレーション適用

```bash
set -a
source .env
set +a

psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DATABASE" -f db/migrations/0000_enable_pgcrypto.sql
psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DATABASE" -f db/migrations/0001_create_users.sql
psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DATABASE" -f db/migrations/0002_seed_dummy_user.sql
psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DATABASE" -f db/migrations/0003_update_dummy_user_password_hash.sql
```

### 2. app-service 起動

```bash
cd app-service
set -a
source ../.env
set +a
go run ./cmd/main.go
```

### 3. web/client 起動

別ターミナルで:

```bash
cd web/client
npm run dev
```

### 4. 動作確認

- `http://localhost:3000/login` を開く
- 次でログイン:
  - email: `demo@example.com`
  - password: `plain-text-demo-password`
- ログイン成功後、`/` で Welcome とダミー一覧が表示されることを確認
