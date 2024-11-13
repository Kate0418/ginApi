package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "back/database"
    "back/database/migrations"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("コマンドを指定してください")
        return
    }

    switch os.Args[1] {
    case "migrate":
        if len(os.Args) < 3 {
            // 引数が不足している場合は全てのモデルをマイグレート
            migrateAll()
        } else {
            // 特定のモデルをマイグレート
            migrateModel(os.Args[2])
        }
    case "model":
        if len(os.Args) < 3 {
            fmt.Println("モデル名を指定してください")
            return
        }
        createModel(os.Args[2])
    default:
        fmt.Println("不明なコマンドです")
    }
}

func migrateAll() {
    db := database.Gorm()
    for modelName, model := range migrations.Models {
        err := db.AutoMigrate(model)
        if err != nil {
            fmt.Printf("%s のマイグレーションに失敗しました: %v\n", modelName, err)
            continue
        }
        fmt.Printf("%s のマイグレーションが完了しました\n", modelName)
    }
    fmt.Println("全てのマイグレーションが完了しました")
}

func migrateModel(modelName string) {
    db := database.Gorm()
    model, exists := migrations.Models[modelName]
    if !exists {
        fmt.Printf("モデル '%s' が見つかりません\n", modelName)
        return
    }

    err := db.AutoMigrate(model)
    if err != nil {
        fmt.Printf("マイグレーションに失敗しました: %v\n", err)
        return
    }
    fmt.Printf("'%s' のマイグレーションが完了しました\n", modelName)
}

func createModel(modelName string) {
    // モデル名の先頭を大文字にする
    modelName = strings.Title(modelName)

    // モデルディレクトリの作成
    modelDir := filepath.Join("api", "models")
    if err := os.MkdirAll(modelDir, 0755); err != nil {
        fmt.Printf("モデルディレクトリの作成に失敗しました: %v\n", err)
        return
    }

    // モデルファイルの内容を生成
    modelContent := fmt.Sprintf(`package models

import (
    "time"
    "gorm.io/gorm"
)

type %s struct {
    ID        uint           `+"`gorm:\"primaryKey\"`"+`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `+"`gorm:\"index\"`"+`
}`, modelName)

    // モデルファイルの作成
    modelPath := filepath.Join(modelDir, strings.ToLower(modelName)+".go")
    if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
        fmt.Printf("モデルファイルの作成に失敗しました: %v\n", err)
        return
    }

    fmt.Printf("%sモデルが正常に作成されました\n", modelName)

    // マイグレーションファイルの更新
    updateMigrationFile(modelName)
}

func updateMigrationFile(modelName string) {
    migrationPath := filepath.Join("database", "migrations", "00000000_000000_migrate.go")

    // 既存のファイルを読み込む
    content, err := os.ReadFile(migrationPath)
    if err != nil {
        fmt.Printf("マイグレーションファイルの読み込みに失敗しました: %v\n", err)
        return
    }

    // ファイルの内容を文字列に変換
    contentStr := string(content)

    // "var Models = map[string]interface{}{" を検索
    marker := "var Models = map[string]interface{}{"
    parts := strings.SplitN(contentStr, marker, 2)
    if len(parts) != 2 {
        fmt.Println("マイグレーションファイルの形式が不正です")
        return
    }

    // 新しいモデルのエントリを作成
    newModelEntry := fmt.Sprintf("\n    \"%s\": &models.%s{},",
        strings.ToLower(modelName),
        modelName)

    // 新しい内容を組み立て
    newContent := parts[0] +
                 marker +
                 newModelEntry +
                 parts[1]

    // ファイルに書き戻す
    err = os.WriteFile(migrationPath, []byte(newContent), 0644)
    if err != nil {
        fmt.Printf("マイグレーションファイルの更新に失敗しました: %v\n", err)
        return
    }

    fmt.Println("マイグレーションファイルを更新しました")
}
