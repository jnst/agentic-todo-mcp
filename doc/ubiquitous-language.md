# Ubiquitous Language

agentic-todo-mcp プロジェクトで使用される共通言語の定義

## Core Concepts

### agentic-todo-mcp

- プロジェクト名、開発コード名

### task
- main-task, sub-task の2階層
- todo/in-progress/done の3状態を持つ
- .todo/task.mdに記述する
- taskは3つの状態を持つ
  - [ ] 未着手
  - [-] 作業中
  - [x] 完了

### main-task (メインタスク)
- 一意のtask-id（`T001`, `T002`, `T010`）を持つ
- main-taskはcontextと1:1で紐付く

### task-id
- task-id は各 main-task の末尾に、半角スペース+ハッシュタグ形式で記述 
  - 例
    - Reactコンポーネントを検討 #T001
    - 要件定義を作成 #T002

### sub-task (サブタスク)
- main-task の子タスク
- task-id は持たない
- context は持たない

### context （コンテキスト）
- main-task に対応するメタ情報
- タスクIDと1:1で対応するMarkdownファイル
- 会議録、決定履歴、知識蓄積を含む
- .todo/cotext 配下に配置
  - 例
    - .todo/cotext/T001.md
    - .todo/cotext/T002.md

### category
- taskを分類するための概念
- 任意の名称を使用できる

### ADR（Architecture Decision Record）
- 設計判断の記録
- main-task 完了時に必要に応じて追加する
- 決定の背景、理由、影響を記録
- Status（Accepted/Deprecated）を管理

## File Structure Terms

### task.md
- 全タスクを管理するメインファイル
- セクション分割による状態管理
- シンプルなMarkdownチェックリスト形式
- 上位ほど優先度が高い

### context/{task-id}.md
- 個別タスクのコンテキスト情報
- task-id と同名のファイル
- 詳細な背景情報や作業履歴

### adr/adr-{number}-{title}.md
- 連番付きのADRファイル
- 設計判断の記録と追跡
- 影響分析の基礎データ
- ファイル名の命名規則
  - .todo/adr/adr-001-core-knowledge.md
  - .todo/adr/adr-002-techstack.md

### .todo/index.md
- プロジェクト全体のインデックス
- ナビゲーションとメタ情報
