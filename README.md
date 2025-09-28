# bubbletea

A Go project.

## Getting Started

### Prerequisites

- Go 1.24 or later

### Usage

```bash
go run main.go
```

## 概要

- 画面を4分割し、各ペインにテーブル（ダミーデータ）を表示します。
- ウィンドウサイズの変更に追従してレイアウトが自動調整されます。
- 終了は `q`、`Esc`、`Ctrl+C` で行います。

## 操作方法

- フォーカス移動: `Tab`（次のペイン） / `Shift+Tab`（前のペイン）
- スクロール（フォーカス中のペインに適用）:
  - `↑/↓` 行単位で移動
  - `PageUp/PageDown` ページ単位で移動
  - `Home/End` 先頭/末尾へ移動
- 終了: `q` / `Esc` / `Ctrl+C`

ヒント: 端末サイズが小さい場合は列がトリムされます。十分な横幅・高さでの実行を推奨します。

## License

This project is licensed under the MIT License.
