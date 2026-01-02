# Communication Rules

- チャットは日本語で行ってください。
- コード内のコメントは英語で記述してください。
- **Comments must strictly follow Godoc (for Go) and TSDoc (for TypeScript) standards.**

# Project Structure

- **Backend**: `backend/` (Go + Echo)
- **Frontend**: `frontend/` (Angular + Bun)

# UI/UX Guidelines

- **Component Library**: `@angular/material` + `@angular/cdk`
- **Animations**: Use `@angular/animations` extensively for rich interactions
  - Page transitions, list stagger, dialog enter/leave
  - Micro-interactions on buttons, cards, inputs
- **Design Principles**:
  - **Lightweight**: Lazy-load modules, tree-shake unused Material components
  - **Rich & Modern**: Glassmorphism, subtle shadows, smooth transitions
  - **Accessible**: Follow Material Design accessibility guidelines
- **Theme**: Custom dark/light theme with CSS variables
- **Icons**: Material Symbols (variable font)

# Commands

## Backend

- `cd backend`
- **Run Server**: `go run ./cmd/sabakan`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`

## Frontend

- `cd frontend`
- **Start Dev Server**: `bun run start` (Runs `ng serve`)
- **Build**: `bun run build` (Runs `ng build`)
- **Watch Build**: `bun run watch`
- **Storybook**: `bun run storybook`
- **Build Storybook**: `bun run build-storybook`

# Code Style

- **TypeScript**:
  - Validated by `@tsconfig/strictest` and `oxlint`.
  - Use `ESNext` features.
  - Strict null checks and type safety are enforced.
- **Go**:
  - Follow standard Go formatting (`go fmt`).
  - Use idiomatic Go (check `effective_go`).
- **General**:
  - Keep functions small and focused.
  - Document public APIs.
  - Use **CRLF** for line endings.

# Workflow

- **Planning**: ALWAYS create/update an `implementation_plan.md` before starting complex tasks.
- **TDD**: Adopt a strict Test-Driven Development (TDD) workflow. Write tests *before* implementation.
- **Incremental Changes**: Make small, verifiable changes. Use formatting commands often.
- **Verification**: Verify changes by running the build or tests after every significant step.
- **Checklists**: Use checklists in planning documents to track progress.
- **Storybook**: Always implement frontend components using Storybook.

# Container Support

- **Primary**: Podman (Initial focus).
- **Secondary**: Docker (Future support planned).
- **Base Images**: Use Alpine or Debian-slim for minimal image sizes.

### Development Container (`Containerfile.dev`)

開発・検証用コンテナ。Podman-in-Podmanでコンテナ管理機能のテストが可能。

```bash
# ビルド
podman build -f Containerfile.dev -t sabakan:dev .

# 実行（開発モード）
podman run -it --rm \
  -v ./backend:/app/backend:Z \
  -v ./frontend:/app/frontend:Z \
  -p 1323:1323 -p 4200:4200 -p 6006:6006 \
  --privileged \
  sabakan:dev
```

### Production Container (`Containerfile`)

本番用マルチステージビルド。Go + Angular の最小イメージ。

```bash
# ビルド
podman build -f Containerfile -t sabakan:latest .

# 実行
podman run -d --name sabakan \
  -v /run/podman/podman.sock:/run/podman/podman.sock:Z \
  -v sabakan-data:/data:Z \
  -p 1323:1323 \
  sabakan:latest
```

### Compose (`compose.yml`)

```bash
# 開発（デフォルト）
podman compose up

# 本番
podman compose --profile prod up -d
```

# Configuration

- **Format**: TOML
- **System Config**: `backend/config.toml` (server, database, logging settings)
- **Game Config**: `backend/games/<game-id>/config.toml` (per-game container and mod settings)
- **Examples**: `config.example.toml`, `games/game.example.toml`

# Internationalization (i18n)

- **Supported Locales**: English (en), Japanese (ja)
- **Translation Files**: `frontend/src/locale/messages.*.xlf`
- **Extract Messages**: `ng extract-i18n`

# CI/CD

- **GitHub Actions**: `.github/workflows/docs.yml`
- **GitHub Pages**: Storybook, TypeDoc, GoDoc

# Authentication

- **Initial Admin**: `admin/admin` (change on first login)
- **OAuth Providers**: Google, Discord
- **API Tokens**: Configurable via settings UI and environment variables
