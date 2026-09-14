# APOD API Migration Notes (EN/JA)

Last updated: 2026-09-14

This note summarizes the APOD feed/API specification used by this repository after migration to the science.nasa.gov APOD Basic endpoints.

Reference:
- https://schlotterer.notion.site/APOD-Feed-And-API-User-Guide-39697d8747c38015a53edfdde76d4f5e

## English

### 1) Primary JSON endpoint

- Base list endpoint:
  - https://science.nasa.gov/wp-json/wp/v2/apod-basic
- Pagination parameters:
  - page
  - per_page (guide says capped at 25)

### 2) Date handling

- Legacy APOD date code format is used in route/filter examples:
  - YYMMDD (example: 140918)
- Guide examples for range filtering:
  - date_from=140901
  - date_to=140930
- JSON response field date is standard format:
  - YYYY-MM-DD

### 3) Single APOD routes

- Single APOD JSON by legacy date code:
  - /wp-json/wp/v2/apod-basic/{YYMMDD}
- Raw HTML route for copy/paste workflows:
  - /wp-json/wp/v2/apod-basic/{YYMMDD}/html

### 4) Important JSON fields

- date
- post_id
- title
- permalink
- media_type
- explanation
- credit
- copyright
- alt
- url
- hdurl
- basic_html
- basic_html_url

### 5) RSS and other related feeds

- APOD Basic RSS:
  - https://science.nasa.gov/feed/apod-basic/
- Future APOD JSON:
  - https://science.nasa.gov/wp-json/wp/v2/future-apods/

### 6) Repository implementation notes

- APOD fetch target is science.nasa.gov APOD Basic endpoint.
- Request flow uses pagination and then client-side filtering for date/range behavior.
- SQLite schema includes permalink column in apod_data.
- Migration is executed for existing DB files as well (not only for new DB files).
- Credit output in social message generation is currently disabled intentionally:
  - new API credit may contain HTML
  - direct output can exceed SNS length constraints
  - TODO marker is left in code for future HTML-to-text + length control support

### 7) Operational cautions

- External API behavior can change without notice.
- Re-verify accepted date filters and formats if date lookup behavior changes.
- Prefer robust fallbacks for media URLs (hdurl/url availability may differ by post).

## 日本語

### 1) 主な JSON エンドポイント

- ベースとなる一覧エンドポイント:
  - https://science.nasa.gov/wp-json/wp/v2/apod-basic
- ページングパラメータ:
  - page
  - per_page（ガイド上は最大 25）

### 2) 日付の扱い

- ルート/フィルタ例では APOD の legacy 日付コードを使用:
  - YYMMDD（例: 140918）
- 期間フィルタ例:
  - date_from=140901
  - date_to=140930
- JSON レスポンスの date フィールドは標準形式:
  - YYYY-MM-DD

### 3) 単一 APOD 取得ルート

- legacy 日付コード指定の単一 JSON:
  - /wp-json/wp/v2/apod-basic/{YYMMDD}
- コピペ向けの生 HTML ルート:
  - /wp-json/wp/v2/apod-basic/{YYMMDD}/html

### 4) 主要 JSON フィールド

- date
- post_id
- title
- permalink
- media_type
- explanation
- credit
- copyright
- alt
- url
- hdurl
- basic_html
- basic_html_url

### 5) RSS など関連フィード

- APOD Basic RSS:
  - https://science.nasa.gov/feed/apod-basic/
- Future APOD JSON:
  - https://science.nasa.gov/wp-json/wp/v2/future-apods/

### 6) このリポジトリでの実装メモ

- APOD 取得先は science.nasa.gov の APOD Basic エンドポイント。
- リクエストはページングで取得し、日付/期間の意味付けはクライアント側でフィルタ。
- SQLite の apod_data に permalink カラムを追加。
- 既存 DB ファイルに対しても migration を実行する。
- SNS 投稿文の credit 出力は現在意図的に無効化:
  - 新 API の credit が HTML を含む場合がある
  - そのまま出すと文字数制限を超えやすい
  - 将来対応用にコードに TODO マーカーを残している

### 7) 運用上の注意

- 外部 API 仕様は予告なく変わる可能性がある。
- 日付指定の挙動が変わった場合は、受理される日付フォーマットを再確認する。
- 投稿によって hdurl/url の有無が異なるため、メディア URL はフォールバック設計を推奨。
