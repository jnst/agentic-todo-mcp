# MCP API 仕様書

agentic-todo-mcp の Model Context Protocol (MCP) ツール詳細仕様

## 1. 概要

### 1.1 サーバー情報
- **名前**: agentic-todo-mcp
- **バージョン**: 0.0.1
- **説明**: AIエージェント用Markdownベースタスク管理システム

### 1.2 提供ツール
- **タスク管理**: 6ツール (create_task, update_task, delete_task, reorder_task, list_tasks, search_tasks)
- **ADR管理**: 3ツール (create_adr, update_adr_status, list_adrs)
- **コンテキスト管理**: 3ツール (update_context, get_context, search_contexts)

## 2. タスク管理ツール

### 2.1 create_task

新規main-taskを作成し、対応するcontextファイルを生成します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "description": "タスクのタイトル",
      "maxLength": 100,
      "minLength": 1
    },
    "category": {
      "type": "string",
      "description": "タスクのカテゴリ（セクション名）",
      "maxLength": 50
    },
    "description": {
      "type": "string",
      "description": "タスクの詳細説明",
      "maxLength": 500
    },
    "subtasks": {
      "type": "array",
      "description": "サブタスクリスト",
      "items": {
        "type": "string",
        "maxLength": 100
      },
      "maxItems": 20
    }
  },
  "required": ["title"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$",
      "description": "生成されたタスクID"
    },
    "file_path": {
      "type": "string",
      "description": "作成されたcontextファイルパス"
    },
    "created_at": {
      "type": "string",
      "format": "date-time",
      "description": "作成日時"
    }
  }
}
```

#### エラーケース
- `TASK_LIMIT_EXCEEDED`: タスク数上限（999個）に達した場合
- `INVALID_CATEGORY`: カテゴリ名が無効な場合
- `FILE_WRITE_ERROR`: ファイル書き込みエラー

### 2.2 update_task

既存taskの部分更新を行います。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$",
      "description": "更新対象のタスクID"
    },
    "title": {
      "type": "string",
      "description": "新しいタイトル",
      "maxLength": 100
    },
    "status": {
      "type": "string",
      "enum": ["todo", "in_progress", "done"],
      "description": "新しいステータス"
    },
    "category": {
      "type": "string",
      "description": "新しいカテゴリ",
      "maxLength": 50
    },
    "subtasks": {
      "type": "array",
      "description": "サブタスクリスト（全置換）",
      "items": {
        "type": "object",
        "properties": {
          "title": {
            "type": "string",
            "maxLength": 100
          },
          "status": {
            "type": "string",
            "enum": ["todo", "in_progress", "done"]
          }
        },
        "required": ["title"]
      }
    }
  },
  "required": ["task_id"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$"
    },
    "updated_fields": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "更新されたフィールド名のリスト"
    },
    "updated_at": {
      "type": "string",
      "format": "date-time"
    }
  }
}
```

#### エラーケース
- `TASK_NOT_FOUND`: 指定されたタスクIDが存在しない場合
- `INVALID_STATUS`: ステータスが無効な場合
- `FILE_WRITE_ERROR`: ファイル書き込みエラー

### 2.3 delete_task

main-taskを削除し、対応するcontextファイルも同時に削除します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$",
      "description": "削除対象のタスクID"
    }
  },
  "required": ["task_id"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$"
    },
    "deleted_files": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "削除されたファイルのリスト"
    },
    "deleted_at": {
      "type": "string",
      "format": "date-time"
    }
  }
}
```

#### エラーケース
- `TASK_NOT_FOUND`: 指定されたタスクIDが存在しない場合
- `FILE_WRITE_ERROR`: ファイル書き込みエラー

### 2.4 reorder_task

タスクの位置を変更し、優先度を調整します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$",
      "description": "移動対象のタスクID"
    },
    "position": {
      "type": "string",
      "enum": ["first", "last", "before", "after"],
      "description": "移動先の位置"
    },
    "reference_task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$",
      "description": "参照タスクID（position が before/after の場合必須）"
    }
  },
  "required": ["task_id", "position"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$"
    },
    "old_position": {
      "type": "integer",
      "description": "元の位置"
    },
    "new_position": {
      "type": "integer",
      "description": "新しい位置"
    },
    "updated_at": {
      "type": "string",
      "format": "date-time"
    }
  }
}
```

#### エラーケース
- `TASK_NOT_FOUND`: 指定されたタスクIDが存在しない場合
- `REFERENCE_TASK_NOT_FOUND`: 参照タスクIDが存在しない場合
- `INVALID_POSITION`: 位置指定が無効な場合
- `FILE_WRITE_ERROR`: ファイル書き込みエラー

### 2.5 list_tasks

タスク一覧を取得します。フィルタリング機能付き。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "status": {
      "type": "string",
      "enum": ["todo", "in_progress", "done"],
      "description": "ステータスフィルタ"
    },
    "category": {
      "type": "string",
      "description": "カテゴリフィルタ"
    },
    "limit": {
      "type": "integer",
      "minimum": 1,
      "maximum": 100,
      "default": 50,
      "description": "取得件数上限"
    }
  }
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "tasks": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "task_id": {
            "type": "string",
            "pattern": "^T[0-9]{3}$"
          },
          "title": {
            "type": "string"
          },
          "status": {
            "type": "string",
            "enum": ["todo", "in_progress", "done"]
          },
          "category": {
            "type": "string"
          },
          "subtasks_count": {
            "type": "integer",
            "description": "サブタスク数"
          },
          "priority": {
            "type": "integer",
            "description": "優先度（位置ベース、小さいほど高優先度）"
          }
        }
      }
    },
    "total_count": {
      "type": "integer",
      "description": "総件数"
    }
  }
}
```

### 2.6 search_tasks

全文検索によるタスク検索を行います。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "検索クエリ",
      "minLength": 1,
      "maxLength": 200
    },
    "search_in": {
      "type": "array",
      "items": {
        "type": "string",
        "enum": ["title", "content", "context"]
      },
      "default": ["title", "content"],
      "description": "検索対象"
    },
    "limit": {
      "type": "integer",
      "minimum": 1,
      "maximum": 50,
      "default": 20
    }
  },
  "required": ["query"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "results": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "task_id": {
            "type": "string",
            "pattern": "^T[0-9]{3}$"
          },
          "title": {
            "type": "string"
          },
          "status": {
            "type": "string",
            "enum": ["todo", "in_progress", "done"]
          },
          "category": {
            "type": "string"
          },
          "match_score": {
            "type": "number",
            "minimum": 0,
            "maximum": 1,
            "description": "関連度スコア"
          },
          "matched_content": {
            "type": "string",
            "description": "マッチしたコンテンツの抜粋"
          }
        }
      }
    },
    "total_matches": {
      "type": "integer"
    }
  }
}
```

## 3. ADR管理ツール

### 3.1 create_adr

新規ADRを作成します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "description": "ADRのタイトル",
      "maxLength": 100,
      "minLength": 1
    },
    "context": {
      "type": "string",
      "description": "決定の背景・文脈",
      "maxLength": 2000
    },
    "decision": {
      "type": "string",
      "description": "決定内容",
      "maxLength": 2000
    },
    "rationale": {
      "type": "string",
      "description": "決定理由",
      "maxLength": 2000
    },
    "consequences": {
      "type": "string",
      "description": "決定の影響・結果",
      "maxLength": 2000
    },
    "status": {
      "type": "string",
      "enum": ["Proposed", "Accepted", "Deprecated"],
      "default": "Proposed",
      "description": "ADRステータス"
    }
  },
  "required": ["title", "context", "decision", "rationale"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "adr_id": {
      "type": "string",
      "pattern": "^adr-[0-9]{3}-.*$",
      "description": "生成されたADR ID"
    },
    "file_path": {
      "type": "string",
      "description": "作成されたADRファイルパス"
    },
    "created_at": {
      "type": "string",
      "format": "date-time"
    }
  }
}
```

### 3.2 update_adr_status

ADRのステータスを更新します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "adr_number": {
      "type": "integer",
      "minimum": 1,
      "maximum": 999,
      "description": "更新対象のADR番号"
    },
    "status": {
      "type": "string",
      "enum": ["Proposed", "Accepted", "Deprecated"],
      "description": "新しいステータス"
    },
    "reason": {
      "type": "string",
      "description": "ステータス変更の理由",
      "maxLength": 500
    }
  },
  "required": ["adr_number", "status"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "adr_number": {
      "type": "integer",
      "minimum": 1,
      "maximum": 999
    },
    "old_status": {
      "type": "string"
    },
    "new_status": {
      "type": "string"
    },
    "updated_at": {
      "type": "string",
      "format": "date-time"
    }
  }
}
```

### 3.3 list_adrs

ADR一覧を取得します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "status": {
      "type": "string",
      "enum": ["Proposed", "Accepted", "Deprecated"],
      "description": "ステータスフィルタ"
    },
    "limit": {
      "type": "integer",
      "minimum": 1,
      "maximum": 100,
      "default": 50
    }
  }
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "adrs": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "adr_number": {
            "type": "integer",
            "minimum": 1,
            "maximum": 999
          },
          "title": {
            "type": "string"
          },
          "status": {
            "type": "string",
            "enum": ["Proposed", "Accepted", "Deprecated"]
          },
          "created_at": {
            "type": "string",
            "format": "date-time"
          },
          "updated_at": {
            "type": "string",
            "format": "date-time"
          }
        }
      }
    },
    "total_count": {
      "type": "integer"
    }
  }
}
```

## 4. コンテキスト管理ツール

### 4.1 update_context

main-taskのコンテキスト情報を追加・更新します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$",
      "description": "対象のタスクID"
    },
    "content": {
      "type": "string",
      "description": "コンテキスト内容",
      "maxLength": 10000
    },
    "append": {
      "type": "boolean",
      "default": false,
      "description": "既存内容に追記するかどうか"
    },
    "section": {
      "type": "string",
      "description": "追記するセクション名",
      "maxLength": 50
    }
  },
  "required": ["task_id", "content"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$"
    },
    "file_path": {
      "type": "string",
      "description": "更新されたcontextファイルパス"
    },
    "updated_at": {
      "type": "string",
      "format": "date-time"
    }
  }
}
```

### 4.2 get_context

main-taskのコンテキスト情報を取得します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$",
      "description": "対象のタスクID"
    }
  },
  "required": ["task_id"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "string",
      "pattern": "^T[0-9]{3}$"
    },
    "content": {
      "type": "string",
      "description": "コンテキスト内容"
    },
    "file_path": {
      "type": "string",
      "description": "contextファイルパス"
    },
    "updated_at": {
      "type": "string",
      "format": "date-time"
    }
  }
}
```

### 4.3 search_contexts

全コンテキストファイルを検索します。

#### 入力スキーマ
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "検索クエリ",
      "minLength": 1,
      "maxLength": 200
    },
    "limit": {
      "type": "integer",
      "minimum": 1,
      "maximum": 50,
      "default": 20
    }
  },
  "required": ["query"]
}
```

#### 出力スキーマ
```json
{
  "type": "object",
  "properties": {
    "results": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "task_id": {
            "type": "string",
            "pattern": "^T[0-9]{3}$"
          },
          "match_score": {
            "type": "number",
            "minimum": 0,
            "maximum": 1
          },
          "matched_content": {
            "type": "string",
            "description": "マッチしたコンテンツの抜粋"
          },
          "file_path": {
            "type": "string"
          }
        }
      }
    },
    "total_matches": {
      "type": "integer"
    }
  }
}
```

## 5. MCPリソース

### 5.1 ファイルリソース

#### 5.1.1 インデックスファイル
- **URI**: `file://.todo/index.md`
- **説明**: プロジェクト全体のナビゲーション
- **アクセス**: 読み取り専用

#### 5.1.2 タスクファイル
- **URI**: `file://.todo/task.md`
- **説明**: 全タスク管理ファイル
- **アクセス**: 読み取り専用

#### 5.1.3 ADRファイル
- **URI**: `file://.todo/adr/{adr-id}.md`
- **説明**: 個別ADRファイル
- **アクセス**: 読み取り専用

#### 5.1.4 コンテキストファイル
- **URI**: `file://.todo/context/{task-id}.md`
- **説明**: 個別タスクコンテキスト
- **アクセス**: 読み取り専用

## 6. エラーハンドリング

### 6.1 エラーコード一覧

#### ファイル操作エラー
- `FILE_NOT_FOUND`: ファイルが見つからない
- `FILE_READ_ERROR`: ファイル読み取りエラー
- `FILE_WRITE_ERROR`: ファイル書き込みエラー
- `PERMISSION_DENIED`: ファイルアクセス権限エラー

#### バリデーションエラー
- `INVALID_TASK_ID`: タスクIDの形式が無効
- `INVALID_ADR_NUMBER`: ADR番号が無効
- `INVALID_STATUS`: ステータスが無効
- `INVALID_CATEGORY`: カテゴリ名が無効
- `CONTENT_TOO_LONG`: コンテンツが最大長を超過

#### ビジネスロジックエラー
- `TASK_NOT_FOUND`: タスクが見つからない
- `ADR_NOT_FOUND`: ADRが見つからない
- `REFERENCE_TASK_NOT_FOUND`: 参照タスクIDが存在しない
- `INVALID_POSITION`: 位置指定が無効
- `TASK_LIMIT_EXCEEDED`: タスク数上限に達した
- `ADR_LIMIT_EXCEEDED`: ADR数上限に達した

### 6.2 エラーレスポンス形式

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ",
    "details": {
      "field": "エラー対象フィールド",
      "value": "エラー値"
    }
  }
}
```

## 7. パフォーマンス要件

### 7.1 応答時間
- **通常操作**: 100ms以内
- **検索操作**: 500ms以内
- **大量データ処理**: 2秒以内

### 7.2 スループット
- **並行セッション**: 最大10セッション
- **ファイル処理**: 10,000ファイルまで効率的に処理
- **メモリ使用量**: 100MB以下

### 7.3 最適化
- **インメモリインデックス**: 高速検索のためのキャッシュ
- **差分更新**: 変更部分のみの更新
- **バッチ処理**: 複数操作の一括処理対応

## 8. 仕様の明確化

### 8.1 ステータス表現の対応

Markdownファイル内の表現とAPI内の表現の対応関係：

| Markdownチェックボックス | API内ステータス | 説明 |
|:---|:---|:---|
| `[ ]` | `"todo"` | 未着手 |
| `[-]` | `"in_progress"` | 作業中 |
| `[x]` | `"done"` | 完了 |

### 8.2 ADR番号管理

- ADRの識別は連番（1-999）で行う
- ファイル名は `adr-{number:03d}-{title}.md` 形式
- APIでは `adr_number` のみで特定し、タイトル変更に影響されない

### 8.3 タスクの位置管理

- 位置は task.md 内のカテゴリ内での順序を示す
- 上位ほど高優先度（position番号は小さい）
- reorder_task で位置変更可能

### 8.4 総合評価

**素晴らしい点:**
- 明確なコンセプト: AIエージェントのコンテキスト記憶補助
- 人間とAIの協調: Markdownフォーマットの採用
- ADRの統合: 背景情報の管理
- 詳細なAPI仕様: JSON Schemaによる厳密な定義

**改善された点:**
- delete_task の追加: 同期削除機能の実装
- reorder_task の追加: 位置ベース優先度管理
- ADR番号による識別: 不変IDでの管理
- ステータス対応の明確化: Markdown↔API変換ルール
