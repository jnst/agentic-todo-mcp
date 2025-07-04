# TODO: agentic-todo-mcp 実装タスク（TDD実践）

## 実装状況
- プロジェクト状況: 設計・計画フェーズ完了、TDD実装開始
- 開発スタイル: t-wadaのTDD（Test-Driven Development）を実践

## TDD実践原則
1. **Red**: 失敗するテストを書く
2. **Green**: テストを通す最小限のコードを書く
3. **Refactor**: コードを改善する

## 優先度：高（High Priority） - 開発環境基盤構築

### 1. Formatter/Linter環境セットアップ（最優先）
- [x] gofmt の設定と実行確認
- [x] goimports の導入（import文の自動整理）
- [x] golangci-lint の導入（包括的なlinter）
- [x] .golangci.yml の設定ファイル作成
- [x] VSCode/IDE設定の整備（保存時自動実行）
- [x] Makefileの作成（開発コマンドの統一）
- [x] go.mod dependencies の追加（go-cmp, gomock, MCP SDK）

### 2. CI/CD基盤構築
- [ ] GitHub Actions の設定
- [ ] formatter/linter の自動実行
- [ ] Pull Request での自動チェック
- [ ] コミット前フックの設定（pre-commit）
- [ ] 品質ゲートの設定（linter通過必須）

### 3. テスト環境セットアップ
- [ ] Go testing フレームワークの準備
- [ ] go-cmp ライブラリの導入（深い比較、差分表示）
- [ ] uber/gomock の導入（モック生成・管理）
- [ ] テストディレクトリ構造の設計
- [ ] テストカバレッジの計測設定

### 4. MCP Server基盤のTDD実装
- [ ] mcp.Server の初期化をTDDで実装する
  - [ ] Server作成のテストケース作成
  - [ ] 最小限のServer初期化実装
  - [ ] Server設定のリファクタリング
- [ ] Transport層をTDDで実装する
  - [ ] StdioTransport のテストケース作成
  - [ ] Transport接続の実装
  - [ ] エラーハンドリングの実装

### 5. コアデータモデルのTDD実装
- [ ] Task構造体をTDDで実装する
  - [ ] Task構造体のテストケース作成
  - [ ] 最小限のTask構造体実装
  - [ ] Task構造体のリファクタリング
- [ ] ADR構造体をTDDで実装する
  - [ ] ADR構造体のテストケース作成
  - [ ] 最小限のADR構造体実装
  - [ ] ADR構造体のリファクタリング
- [ ] Context構造体をTDDで実装する
  - [ ] Context構造体のテストケース作成
  - [ ] 最小限のContext構造体実装
  - [ ] Context構造体のリファクタリング
- [ ] バリデーション機能をTDDで実装する

## 優先度：中（Medium Priority） - コア機能のTDD実装

### 6. ファイル操作層のTDD実装
- [ ] MarkdownパーサーをTDDで実装する
  - [ ] パーサーのテストケース作成（様々な入力パターン）
  - [ ] 最小限のパーサー実装
  - [ ] パーサーのリファクタリング
- [ ] ファイルI/OをTDDで実装する
  - [ ] ファイルI/Oのテストケース作成（gomockでモック使用）
  - [ ] 最小限のファイルI/O実装
  - [ ] ファイルI/Oのリファクタリング
- [ ] データ整合性チェックをTDDで実装する

### 7. MCPツール実装のTDD実装
- [ ] タスク管理MCPツールをTDDで実装する（6ツール）
  - [ ] create_task MCPツールのテスト・実装
    - [ ] mcp.NewServerTool でのツール定義
    - [ ] ハンドラー関数の実装（CallToolParams → CallToolResult）
    - [ ] 入力パラメータのバリデーション
  - [ ] update_task MCPツールのテスト・実装
  - [ ] delete_task MCPツールのテスト・実装
  - [ ] reorder_task MCPツールのテスト・実装
  - [ ] list_tasks MCPツールのテスト・実装
  - [ ] search_tasks MCPツールのテスト・実装
- [ ] ADR管理MCPツールをTDDで実装する（3ツール）
  - [ ] create_adr MCPツールのテスト・実装
  - [ ] update_adr_status MCPツールのテスト・実装
  - [ ] list_adrs MCPツールのテスト・実装
- [ ] コンテキスト管理MCPツールをTDDで実装する（3ツール）
  - [ ] update_context MCPツールのテスト・実装
  - [ ] get_context MCPツールのテスト・実装
  - [ ] search_contexts MCPツールのテスト・実装

### 8. 検索機能のTDD実装
- [ ] 検索アルゴリズムをTDDで実装する
  - [ ] 全文検索のテストケース作成
  - [ ] 関連度スコア算出のテスト・実装
  - [ ] インメモリインデックスのテスト・実装

## 優先度：低（Low Priority） - 統合・最適化のTDD実装

### 9. MCP統合のTDD実装
- [ ] MCP Serverの統合をTDDで実装する
  - [ ] server.AddTools() での全ツール登録テスト
  - [ ] server.Run() での stdio通信テスト
  - [ ] エンドツーエンドの統合テスト作成
  - [ ] エラーハンドリングの網羅的テスト
- [ ] エラーハンドリングをTDDで実装する
  - [ ] エラーケースの網羅的テスト作成
  - [ ] 統一的なエラー処理の実装
- [ ] パフォーマンステストを実装する
  - [ ] 応答時間要件のテスト（100ms以内）
  - [ ] 並行処理のテスト
  - [ ] メモリ使用量のテスト

### 10. 継続的リファクタリング
- [ ] 各Red-Green-Refactorサイクルでコード改善
- [ ] テストカバレッジの向上
- [ ] コードの可読性・保守性の向上

## 技術要件
- Go 1.24.3
- github.com/modelcontextprotocol/go-sdk v0.1.0
- github.com/google/go-cmp （テスト比較）
- github.com/uber-go/mock （モック生成）
- 応答時間: 通常操作 < 100ms, 検索操作 < 500ms
- サポートファイル数: 最大10,000ファイル

## 開発ツール選択方針

### Formatter/Linter（最優先）
- **gofmt**: Go標準フォーマッター
- **goimports**: import文の自動整理
- **golangci-lint**: 包括的なlinter（deadcode, unused, misspell等）
- **デグレーション防止**: AI開発で特に重要、必須の整備

### テストツール
- **go-cmp**: 深い比較・差分表示のメインツール
- **uber/gomock**: モック生成・管理（必要な場合のみ）
- **標準testing**: Go標準のテストフレームワーク
- **testify等のライブラリ**: できるだけ使わない方針

### MCP SDK関連の重要ポイント
- **mcp.NewServer()**: サーバー初期化（名前、バージョン、オプション）
- **mcp.NewServerTool()**: ツール定義（名前、説明、ハンドラー、入力スキーマ）
- **server.AddTools()**: サーバーにツールを登録
- **server.Run()**: stdio上でサーバー実行
- **CallToolParams/CallToolResult**: ツール呼び出しの入出力型
- **mcp.StdioTransport**: 標準入出力での通信

## ディレクトリ構造計画
```
agentic-todo-mcp/
├── cmd/
│   └── server/
│       └── main.go           # MCPサーバーエントリーポイント
├── internal/
│   ├── config/              # 設定管理
│   ├── models/              # データモデル (Task, ADR, Context)
│   ├── storage/             # ファイル操作・永続化
│   ├── parser/              # Markdownパーサー
│   ├── search/              # 検索・インデックス
│   └── mcp/                 # MCPツール実装
├── pkg/
│   └── types/               # 公開型定義
├── .todo/                   # 管理対象ディレクトリ
│   ├── task.md
│   ├── index.md
│   ├── context/
│   └── adr/
├── tests/                   # テストファイル
└── docs/                    # ドキュメント
```

## 参考資料
- @doc/requirements.md: 詳細技術要件
- @doc/mcp-spec.md: MCP API仕様
- @doc/ubiquitous-language.md: 用語定義
- @CLAUDE.md: プロジェクト指示書

## MCP Go SDKの基本実装パターン

### サーバー初期化
```go
server := mcp.NewServer("agentic-todo-mcp", "v0.1.0", nil)
```

### ツール定義
```go
type CreateTaskParams struct {
    Title       string   `json:"title"`
    Category    string   `json:"category,omitempty"`
    Description string   `json:"description,omitempty"`
    Subtasks    []string `json:"subtasks,omitempty"`
}

func CreateTaskHandler(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateTaskParams]) (*mcp.CallToolResultFor[any], error) {
    // 実装
    return &mcp.CallToolResultFor[any]{
        Content: []mcp.Content{&mcp.TextContent{Text: "Task created"}},
    }, nil
}

server.AddTools(
    mcp.NewServerTool("create_task", "Create new main-task", CreateTaskHandler, mcp.Input(
        mcp.Property("title", mcp.Description("Task title")),
        mcp.Property("category", mcp.Description("Task category")),
    )),
)
```

### サーバー実行
```go
if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
    log.Fatal(err)
}
```
