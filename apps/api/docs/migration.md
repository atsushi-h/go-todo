# データベースマイグレーション手順

このドキュメントでは、Atlasを使用したデータベースマイグレーションの手順を説明します。

## 目次

- [環境構成](#環境構成)
- [基本的な流れ](#基本的な流れ)
- [マイグレーション手順](#マイグレーション手順)
- [よく使うコマンド](#よく使うコマンド)
- [トラブルシューティング](#トラブルシューティング)
- [補足情報](#補足情報)

## 環境構成

プロジェクトでは以下の環境を想定しています:

### local環境
- ローカル開発環境（Dockerコンテナ）
- Atlas CLIをDockerコンテナ内で実行
- DB接続先: `go_todo_db:5432`（Dockerネットワーク内）

### dev環境
- デプロイされた開発環境（AWS RDS等）
- CI/CD（GitHub Actions等）から実行
- DB接続先: 開発環境のRDSエンドポイント

### prod環境
- 本番環境（AWS RDS等）
- CI/CDから実行（手動トリガー推奨）
- DB接続先: 本番環境のRDSエンドポイント

## ディレクトリ構成

```
apps/api/
├── db/
│   ├── schema.sql          # スキーマ定義（真実の源）
│   ├── query/              # SQLクエリ定義
│   │   ├── todo.sql
│   │   └── user.sql
│   ├── migrations/         # マイグレーションファイル
│   │   ├── 20251125113905_init.sql
│   │   └── atlas.sum
│   └── sqlc/               # sqlc生成コード（自動生成）
│       ├── db.go
│       ├── models.go
│       ├── todo.sql.go
│       └── user.sql.go
├── sqlc.yaml               # sqlc設定
└── atlas.hcl               # Atlas設定
```

## 基本的な流れ

```
1. db/schema.sql を編集
   ↓
2. マイグレーションファイルを生成（atlas migrate diff）
   ↓
3. マイグレーションファイルをレビュー
   ↓
4. マイグレーションを適用（atlas migrate apply）
   ↓
5. sqlcコードを再生成（sqlc generate）
   ↓
6. 変更をGitにコミット（atlas.sumも含む）
```

または、一括コマンドを使用:

```bash
make schema-update NAME=<migration_name>
```

## マイグレーション手順

### 1. スキーマの変更

例: Todoテーブルに新しいカラムを追加

```sql
-- apps/api/db/schema.sql
CREATE TABLE todos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    priority INTEGER NOT NULL DEFAULT 0,  -- 新規追加
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 2. クエリの追加・更新（必要な場合）

```sql
-- apps/api/db/query/todo.sql

-- name: ListTodosByUserOrderByPriority :many
SELECT * FROM todos
WHERE user_id = $1
ORDER BY priority DESC, created_at DESC;
```

### 3. マイグレーションファイルの生成

```bash
# ローカル環境
make migrate-diff NAME=add_priority_to_todos
```

`apps/api/db/migrations/YYYYMMDDHHMMSS_add_priority_to_todos.sql` が生成されます。

### 4. マイグレーションファイルのレビュー

生成されたSQLファイルを確認:

```sql
-- apps/api/db/migrations/20250125120000_add_priority_to_todos.sql
-- Add column "priority" to table: "todos"
ALTER TABLE "todos" ADD COLUMN "priority" integer NOT NULL DEFAULT 0;
```

必要に応じて手動で編集可能です。

### 5. マイグレーションの状態確認

```bash
# ローカル環境
make migrate-status
```

出力例:
```
Migration Status: PENDING
  Current Version: 20250125100000
  Next Version:    20250125120000
  Total Pending:   1
```

### 6. マイグレーションの適用

```bash
# ローカル環境
make migrate-apply
```

成功すると以下のような出力:
```
Migrating to version 20250125120000 (1 migration)
  -> 20250125120000_add_priority_to_todos.sql ................ ok (15ms)
  -------------------------------
  -> 0s
```

### 7. sqlcコードの再生成

```bash
make sqlc-generate
```

これにより `db/sqlc/` 配下のコードが更新されます。

### 8. Gitへのコミット

マイグレーションファイルと atlas.sum をコミット:

```bash
git add apps/api/db/
git commit -m "feat: Add priority column to todos table"
```

## よく使うコマンド

### マイグレーションファイル生成

```bash
make migrate-diff NAME=<migration_name>
```

現在のスキーマとDBの差分からマイグレーションファイルを生成します。

### マイグレーション状態確認

```bash
make migrate-status
```

適用済み/未適用のマイグレーションを確認します。

### マイグレーション適用

```bash
make migrate-apply
```

未適用のマイグレーションを全て適用します。

### sqlcコード生成

```bash
make sqlc-generate
```

`db/query/` 内のSQLから Go コードを生成します。

### スキーマ変更の一括処理

```bash
make schema-update NAME=<migration_name>
```

マイグレーション生成とsqlcコード生成を一括実行します。

### 全コード生成

```bash
make generate
```

OpenAPI + sqlc のコードを一括生成します。

### ハッシュ値の再計算

```bash
make migrate-hash
```

マイグレーションファイルを手動編集した場合に実行します。

## トラブルシューティング

### エラー: "atlas_dev database does not exist"

**原因**: Atlas作業用の一時データベースが作成されていない

**解決方法**:
```bash
make create-atlas-dev-db
```

または手動で:
```bash
docker exec -it go_todo_db psql -U user -d go_todo_db -c "CREATE DATABASE atlas_dev;"
```

### エラー: "migration checksum mismatch"

**原因**: マイグレーションファイルを直接編集したが、`atlas.sum` が更新されていない

**解決方法**:
```bash
make migrate-hash
```

### エラー: "connection refused"

**原因**:
- Atlas CLIの実行場所とDB接続先のホスト名が一致していない
- DBコンテナが起動していない

**解決方法**:
```bash
# DBコンテナの状態確認
docker ps | grep go_todo_db

# Dockerコンテナ内から実行する場合は go_todo_db を使用
# ホストマシンから実行する場合は localhost を使用
```

### マイグレーションが中途半端に失敗した場合

```bash
# 現在の状態を確認
make migrate-status

# 必要に応じて手動でSQLを実行してデータを修正
docker exec -it go_todo_db psql -U user -d go_todo_db

# atlas_migrationsテーブルを確認
SELECT * FROM atlas_migrations;

# 必要に応じてバージョンをロールバック（慎重に！）
DELETE FROM atlas_migrations WHERE version = '20250125120000';
```

## 補足情報

### atlas_dev データベースとは

- Atlasが差分計算のために使用する**一時的な作業用データベース**
- マイグレーション実行時に以下の処理で使用される:
  1. `atlas_dev` に現在のスキーマ（`db/schema.sql`）を適用
  2. 実際のDB（`go_todo_db`）と `atlas_dev` を比較
  3. 差分をマイグレーションファイルとして生成
- データは保存されず、毎回クリーンな状態で使用される
- 本番DBとは完全に分離されているため安全

### atlas.sum について

- マイグレーションファイルの**整合性を保証**するためのハッシュファイル
- 各マイグレーションファイルのSHA256ハッシュ値を記録
- **必ずGitで管理する**（コミット必須）
- マイグレーションファイルを手動編集した場合は `make migrate-hash` で再計算

### sqlc について

- SQLクエリから型安全なGoコードを自動生成するツール
- `db/query/*.sql` に定義したクエリが `db/sqlc/*.go` に生成される
- 生成されたコードは直接編集しない（再生成で上書きされる）

### ベストプラクティス

1. **db/schema.sql を真実の源として管理**
   - スキーマ変更は必ず `db/schema.sql` から行う
   - マイグレーションファイルは差分として生成される

2. **小さく頻繁にマイグレーション**
   - 1つのマイグレーションで1つの変更
   - 複雑な変更は複数のマイグレーションに分割

3. **本番適用前にdev環境でテスト**
   - dev環境で必ずテストしてから本番適用
   - ロールバック手順も事前に確認

4. **破壊的変更に注意**
   - カラム削除やテーブル削除は慎重に
   - 可能であれば段階的に実施:
     1. 新カラム追加 → アプリ移行
     2. 旧カラム削除

5. **atlas.sum は必ずコミット**
   - マイグレーションファイルと同時にコミット
   - CIでチェックサム検証を行う

6. **sqlc生成コードは編集しない**
   - クエリを変更したい場合は `db/query/*.sql` を編集
   - `make sqlc-generate` で再生成

## 参考リンク

- [Atlas CLI Documentation](https://atlasgo.io/docs)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [pgx Documentation](https://github.com/jackc/pgx)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
