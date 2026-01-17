# Commit Record

一個使用 AI 自動生成每日工作總結的 Git 提交記錄分析工具。

## 專案簡介

Commit Record 是一個基於 Go 語言開發的命令列工具，能夠自動分析 Git 專案的提交記錄，並透過 OpenAI API 生成繁體中文的每日工作總結。支援單一專案或多專案的批次處理，幫助開發者快速回顧和記錄每日的工作成果。

## 功能特色

- **自動化提交分析**：自動讀取 Git 提交歷史並分析變更內容
- **AI 驅動總結**：使用 OpenAI GPT 模型生成結構化的工作總結
- **多專案支援**：可同時處理多個專案的提交記錄
- **雙模式運行**：
  - `--bench`：批次處理模式，一次 API 調用處理所有專案
  - `--single`：單一處理模式，逐個專案分別處理
- **自動持久化**：將生成的總結自動保存為 Markdown 檔案
- **繁體中文輸出**：生成的總結內容以繁體中文呈現

## 環境需求

- Go 1.16 或更高版本
- Git
- OpenAI API Key

## 安裝

### 1. Clone 專案

```bash
git clone <repository-url>
cd commit-record
```

### 2. 安裝依賴

```bash
cd src
go mod download
```

### 3. 編譯

```bash
go build -o commit-record main.go
```

或使用 Makefile（如果有提供）：

```bash
make build
```

## 配置

### 設定 OpenAI API Key

必須設定 `OPENAI_API_KEY` 環境變數：

```bash
export OPENAI_API_KEY="your-openai-api-key"
```

建議將此設定加入到你的 shell 配置檔案中（如 `.bashrc` 或 `.zshrc`）：

```bash
echo 'export OPENAI_API_KEY="your-openai-api-key"' >> ~/.zshrc
source ~/.zshrc
```

## 使用方法

### 基本語法

```bash
commit-record [--bench|--single] <project-path-1> [<project-path-2> ...]
```

### 參數說明

- `--bench`：批次處理模式（推薦用於多專案）
  - 將所有專案的提交記錄合併後，一次性調用 OpenAI API
  - 適合需要整合多個專案工作成果的情況
  - 節省 API 調用次數

- `--single`：單一處理模式
  - 逐個專案分別調用 OpenAI API
  - 適合需要獨立分析每個專案的情況
  - 生成更詳細的個別專案總結

- `<project-path>`：Git 專案的路徑（可以是相對或絕對路徑）

**注意**：
- 如果未指定模式標誌，預設使用 `--bench` 批次處理模式
- 不能同時使用 `--bench` 和 `--single` 標誌

### 使用範例

#### 1. 處理單一專案（使用預設批次模式）

```bash
./commit-record /path/to/my-project
```

#### 2. 批次處理多個專案

```bash
./commit-record --bench /path/to/project1 /path/to/project2 /path/to/project3
```

#### 3. 單一模式處理多個專案

```bash
./commit-record --single /path/to/project1 /path/to/project2
```

#### 4. 使用相對路徑

```bash
./commit-record --bench ../frontend ../backend
```

#### 5. 處理當前目錄

```bash
./commit-record --single .
```

## 輸出範例

工具會自動生成包含以下內容的 Markdown 檔案，儲存於 `~/work_conclusion/YYYY-MM-DD/專案名稱.md`：

```markdown
# 每日工作總結 - 2026-01-17

## 專案：my-awesome-project

### 主要工作項目

1. **功能開發**
   - 實作使用者認證模組
   - 新增登入/登出功能

2. **問題修復**
   - 修復資料庫連線逾時問題
   - 解決前端頁面載入緩慢的效能問題

3. **優化改進**
   - 重構 API 層級架構
   - 提升測試覆蓋率至 85%

### 技術亮點

- 採用 JWT 實作安全的身分驗證機制
- 優化資料庫查詢，效能提升 40%

### 明日計畫

- 完成使用者權限管理功能
- 進行整合測試
```

**檔案位置範例：**
```
~/work_conclusion/
└── 2026-01-17/
    ├── my-awesome-project.md
    ├── frontend-app.md
    └── backend-api.md
```

## 專案結構

```
commit-record/
├── src/
│   ├── main.go                    # 主程式入口
│   ├── internal/
│   │   ├── git/                   # Git 操作相關
│   │   ├── proxies/               # OpenAI API 代理
│   │   ├── repositories/          # 資料持久化
│   │   └── services/              # 業務邏輯服務
│   ├── go.mod
│   └── go.sum
└── README.md
```

## 工作原理

1. **讀取提交記錄**：從指定的 Git 專案中讀取今日的所有提交
2. **提取變更資訊**：分析每個提交的訊息和檔案變更
3. **AI 分析**：將提交資訊傳送給 OpenAI API 進行智能分析
4. **生成總結**：接收 AI 生成的結構化工作總結
5. **保存結果**：將總結保存為 Markdown 格式的檔案

## 常見問題

### Q: 如何更改使用的 OpenAI 模型？

A: 編輯 `src/internal/proxies/open_ai_proxy.go` 檔案，修改 `Model` 參數（目前使用 `openai.GPT4oMini`）。

### Q: 生成的總結檔案保存在哪裡？

A: 檔案會保存在使用者家目錄的 `~/work_conclusion/YYYY-MM-DD/` 目錄下，每個專案一個獨立的 Markdown 檔案，檔名為專案名稱。例如：`~/work_conclusion/2026-01-17/my-project.md`。

### Q: 可以自訂總結的格式嗎？

A: 可以，修改 `src/internal/proxies/open_ai_proxy.go` 中的 `GetConclusion` 和 `GetBatchConclusion` 方法內的 prompt 即可調整輸出格式。

### Q: 支援哪些語言？

A: 目前預設使用繁體中文，可以修改 prompt 來支援其他語言。

## 授權

[請根據實際情況添加授權資訊]

## 貢獻

歡迎提交 Issue 和 Pull Request！

## 聯絡方式

[請根據實際情況添加聯絡資訊]
