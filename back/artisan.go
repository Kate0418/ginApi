package main

import (
    "fmt"
    "os"
    "time"
    "path/filepath"
    "strings"
    "back/database/migrations"
    "io/ioutil"
    "regexp"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("コマンドを指定してください")
        return
    }

    switch os.Args[1] {
    case "migrate":
        migrations.Migrate()
        fmt.Println("マイグレーションが完了しました")
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

func createModel(modelName string) {
    // モデル名の先頭を大文字にする
    modelName = strings.Title(modelName)

    // モデルファイルの作成
    modelDir := "api/models"
    if err := os.MkdirAll(modelDir, 0755); err != nil {
        fmt.Printf("モデルディレクトリの作成に失敗しました: %v\n", err)
        return
    }

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
}
`, modelName)

    modelPath := filepath.Join(modelDir, strings.ToLower(modelName)+".go")
    if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
        fmt.Printf("モデルファイルの作成に失敗しました: %v\n", err)
        return
    }

    // マイグレーションファイルの作成
    timestamp := time.Now().Format("20060102_150405")
    migrationName := fmt.Sprintf("%s_%s", timestamp, strings.ToLower(modelName))

    migrationContent := fmt.Sprintf(`package migrations

import (
    "back/api/models"
    "back/database"
)

func Create%s() {
    db := database.Gorm()
    db.AutoMigrate(&models.%s{})
}
`, modelName, modelName)

    migrationDir := "database/migrations"
    if err := os.MkdirAll(migrationDir, 0755); err != nil {
        fmt.Printf("マイグレーションディレクトリの作成に失敗しました: %v\n", err)
        return
    }

    migrationPath := filepath.Join(migrationDir, migrationName+".go")
    if err := os.WriteFile(migrationPath, []byte(migrationContent), 0644); err != nil {
        fmt.Printf("マイグレーションファイルの作成に失敗しました: %v\n", err)
        return
    }

    // 00000000_000000_migrate.goファイルの更新
    baseMigratePath := filepath.Join("database/migrations", "00000000_000000_migrate.go")

    // 既存のファイルを読み込む
    content, err := ioutil.ReadFile(baseMigratePath)
    if err != nil {
        fmt.Printf("00000000_000000_migrate.goファイルの読み込みに失敗しました: %v\n", err)
        return
    }

    // 既存の関数呼び出しを保持しながら新しい関数を追加
    contentStr := string(content)

    // Migrate関数の中身を探す
    re := regexp.MustCompile(`func Migrate\(\) {([\s\S]*?)}`)
    matches := re.FindStringSubmatch(contentStr)

    var newContent string
    if len(matches) > 1 {
        // 既存の関数本体を取得
        functionBody := matches[1]
        // 新しい関数呼び出しを追加
        newFunctionBody := functionBody
        if !strings.Contains(newFunctionBody, fmt.Sprintf("Create%s()", modelName)) {
            if strings.TrimSpace(newFunctionBody) == "" {
                newFunctionBody = fmt.Sprintf("\n    Create%s()\n", modelName)
            } else {
                newFunctionBody += fmt.Sprintf("    Create%s()\n", modelName)
            }
        }
        newContent = fmt.Sprintf(`package migrations

func Migrate() {%s}
`, newFunctionBody)
    } else {
        // Migrate関数が見つからない場合は新規作成
        newContent = fmt.Sprintf(`package migrations

func Migrate() {
    Create%s()
}
`, modelName)
    }

    if err := os.WriteFile(baseMigratePath, []byte(newContent), 0644); err != nil {
        fmt.Printf("00000000_000000_migrate.goファイルの更新に失敗しました: %v\n", err)
        return
    }

    fmt.Printf("%sモデルとマイグレーションファイルが正常に作成されました\n", modelName)
}
