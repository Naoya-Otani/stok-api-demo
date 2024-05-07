
# Stok Server Demo

このリポジトリは、Stokのバックエンドサーバー（デモ）を定義しています。



## 環境変数

実行するには、.envファイルに以下の環境変数を追加する必要があります。

`DB_USER`

`DB_PASSWORD`

`DB_NAME`

`DB_HOST`

`DB_PORT`

## ローカルで起動する

プロジェクトをクローンする

```bash
  git clone https://github.com/Naoya-Otani/stok-api-demo.git
```

ディレクトリを移動する

```bash
  cd stok-api-demo
```

依存関係のインストールする

```bash
  go mod tidy
```

MySQLを起動する (Homebrewを使用)

```bash
  brew services start mysql
```

MySQLの初期設定

```bash
  mysql_secure_installation
```

MySQLへのログイン

```bash
  mysql -u root -p
```

スキーマ設定等は [`init.sql`](https://github.com/Naoya-Otani/stok-api-demo/blob/main/init.sql)を確認してください

サーバーを起動する

```bash
  go run .
```
## API

#### 全件取得

```http
  GET /products
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `Content-Type` | `application/json` | None |

#### 商品を追加

```http
  POST /products
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `product_name`      | `string` | **必須** 商品名 |
| `brand_id`      | `int` | **必須** ブランドのID |
| `image_paths`      | `[]string` | **オプショナル** 商品画像へのパス |

## デモ動画

https://github.com/Naoya-Otani/stok-api-demo/assets/102457026/826e262d-662f-4f69-a79f-ff8517e3b4d2

## 作業時間

| 項目 | 時間 |
| --- | --- |
| 環境構築 | 10m |
| スキーマ設計 | 90m |
| アプリケーション実装 | 90m |
| リファクタリング | 60m |


## リファレンス

 - [Stok サーバーサイドエンジニア技術課題](https://franky-inc.notion.site/Stok-e7c86b932e364e0f838f4091437d1490)
 - [サンプルデータ](https://docs.google.com/spreadsheets/d/1g2LyTAW3BDACn8Btge_A0yi6LNSPD2c9p25VLbJO38Q/edit?usp=sharing)

