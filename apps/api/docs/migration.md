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

## 基本的な流れ

```
1. Goのモデルを編集
   ↓
2. スキーマSQLを生成（go run cmd/schema/main.go）
   ↓
3. マイグレーションファイルを生成（atlas migrate diff）
   ↓
4. マイグレーションファイルをレビュー
   ↓
5. マイグレーションを適用（atlas migrate apply）
   ↓
6. 変更をGitにコミット（atlas.sumも含む）
```

## マイグレーション手順

### 1. モデルの変更

例: Todoモデルに新しいカラムを追加

```go
// apps/api/internal/model/todo.model.go
type Todo struct {
    ID          uint      `json:"id" bun:",pk,autoincrement"`
    Title       string    `json:"title" bun:",notnull"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed" bun:",default:false"`
    Priority    int       `json:"priority" bun:",default:0"` // 新規追加
    CreatedAt   time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt   time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}
```

### 2. スキーマSQLの生成

```bash
# ローカル環境
make generate-schema
```

これにより `apps/api/schema.sql` が更新されます。

### 3. マイグレーションファイルの生成

```bash
# ローカル環境
make migrate-diff NAME=add_priority_to_todos
```

`apps/api/migrations/YYYYMMDDHHMMSS_add_priority_to_todos.sql` が生成されます。

### 4. マイグレーションファイルのレビュー

生成されたSQLファイルを確認:

```sql
-- apps/api/migrations/20250125120000_add_priority_to_todos.sql
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

### 7. Gitへのコミット

マイグレーションファイルと atlas.sum をコミット:

```bash
git add apps/api/migrations/
git add apps/api/atlas.sum
git commit -m "feat: Add priority column to todos table"
```

## よく使うコマンド

### スキーマ生成

```bash
make generate-schema
```

Goのモデルから `schema.sql` を生成します。

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

### 特定バージョンへのマイグレーション

```bash
docker exec -i go_todo_server \
  atlas migrate apply \
  --env local \
  --config file://atlas.hcl \
  --to-version 20250125120000
```

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
# PostgreSQLに接続
docker exec -it go_todo_db psql -U user -d go_todo_db

# atlas_dev データベースを作成
CREATE DATABASE atlas_dev;
\q
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
  1. `atlas_dev` に現在のスキーマ（`schema.sql`）を適用
  2. 実際のDB（`go_todo_db`）と `atlas_dev` を比較
  3. 差分をマイグレーションファイルとして生成
- データは保存されず、毎回クリーンな状態で使用される
- 本番DBとは完全に分離されているため安全

### atlas.sum について

- マイグレーションファイルの**整合性を保証**するためのハッシュファイル
- 各マイグレーションファイルのSHA256ハッシュ値を記録
- **必ずGitで管理する**（コミット必須）
- マイグレーションファイルを手動編集した場合は `make migrate-hash` で再計算

### ベストプラクティス

1. **マイグレーションファイルは編集しない**
   - 一度生成したら基本的に編集しない
   - どうしても必要な場合は `make migrate-hash` を実行

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

## 参考リンク

- [Atlas CLI Documentation](https://atlasgo.io/docs)
- [Bun ORM Documentation](https://bun.uptrace.dev/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
