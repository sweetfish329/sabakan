# Sabakan Developer Guide

## Communication Rules

- チャットは日本語で行ってください。
- コード内のコメントは英語で記述してください。
- **Comments must strictly follow Godoc (for Go) and TSDoc (for TypeScript) standards.**

## Project Structure

- **Backend**: `backend/` (Go + Echo)
- **Frontend**: `frontend/` (Angular + Bun)

## UI/UX Guidelines

- **Component Library**: `@angular/material` + `@angular/cdk`
- **Animations**: Use CSS Animations / Transitions or View Transitions API.
  - **Deprecated**: `@angular/animations` (`trigger`, `transition`, etc.) are deprecated in Angular v21+ and must not be used.
- **Design Principles**:
  - **Lightweight**: Lazy-load modules, tree-shake unused Material components
  - **Rich & Modern**: Glassmorphism, subtle shadows, smooth transitions, premium feel
  - **Accessible**: Follow Material Design accessibility guidelines
- **Theme**: Custom dark/light theme with CSS variables
- **Icons**: Material Symbols (variable font)
- **Aesthetic Guidelines**:
  - **Glassmorphism**: Use for cards, overlays, and sticky headers. (Blur + Transparency + Thin Border)
  - **Gradients**: Use sparingly. Avoid heavy, high-contrast linear gradients for backgrounds. Use subtle, mesh-like gradients or solid colors with glass effects.
  - **Shadows**: Soft, diffused shadows for depth. Avoid harsh, black shadows.

## Commands

### Backend

- `cd backend`
- **Run Server**: `go run ./cmd/sabakan`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`

### Frontend

- `cd frontend`
- **Start Dev Server**: `bun run start` (Runs `ng serve`)
- **Build**: `bun run build` (Runs `ng build`)
- **Watch Build**: `bun run watch`
- **Storybook**: `bun run storybook`
- **Build Storybook**: `bun run build-storybook`

## Code Style

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

## Workflow

> [!IMPORTANT]
> **YOU MUST FOLLOW THESE RULES WITHOUT EXCEPTION.**

### Planning Phase (REQUIRED for complex tasks)

1. **ALWAYS** create/update `implementation_plan.md` before starting complex tasks
2. Ask clarifying questions before coding
3. Use checklists in planning documents to track progress

### Test-Driven Development (TDD) - MANDATORY

> [!CAUTION]
> **NEVER write implementation code before tests.** This is a strict requirement.

#### Backend (Go) TDD Workflow

```text
1. Write test first → go test ./... (tests FAIL - Red)
2. Write minimal implementation → go test ./... (tests PASS - Green)
3. Refactor → go test ./... (tests still PASS)
4. Commit
```

**Example prompt pattern:**

```text
"Write tests for [feature] first. Do NOT write any implementation code yet.
The tests should verify [expected behavior]. Run the tests and confirm they fail."
```

#### Frontend (Angular) TDD Workflow

```text
1. Create Storybook story first
2. Write unit test (*.spec.ts) → tests FAIL
3. Implement component → tests PASS
4. Verify in Storybook visually
5. Commit
```

### Storybook-First Development - MANDATORY

> [!CAUTION]
> **NEVER create Angular components without Storybook stories.** This is a strict requirement.

#### Required Files for Each Component

```text
component-name/
├── component-name.ts          # Component implementation
├── component-name.html        # Template
├── component-name.scss        # Styles
├── component-name.spec.ts     # Unit tests
└── component-name.stories.ts  # Storybook stories (REQUIRED)
```

#### Storybook Workflow

1. **Create story file FIRST** with visual states (default, loading, error, etc.)
2. Run `bun run storybook` to view component in isolation
3. Implement component to match story expectations
4. Iterate visually until design is correct
5. Commit

### Incremental Changes

- Make small, verifiable changes
- Run `go fmt ./...` (backend) or formatting commands often
- Run tests after every significant step
- Commit frequently with meaningful messages

### Verification Checklist (Before Commit)

- [ ] All tests pass (`go test ./...` or `bun run test`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Storybook stories exist for new components
- [ ] Build succeeds (`go build ./...` or `bun run build`)

## Container Support

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

## Configuration

- **Format**: TOML
- **System Config**: `backend/config.toml` (server, database, logging settings)
- **Game Config**: `backend/games/<game-id>/config.toml` (per-game container and mod settings)
- **Examples**: `config.example.toml`, `games/game.example.toml`

## Internationalization (i18n)

- **Supported Locales**: English (en), Japanese (ja)
- **Translation Files**: `frontend/src/locale/messages.*.xlf`
- **Extract Messages**: `ng extract-i18n`

## CI/CD

- **GitHub Actions**: `.github/workflows/docs.yml`
- **GitHub Pages**: Storybook, TypeDoc, GoDoc

## Authentication

- **Initial Admin**: `admin/admin` (change on first login)
- **OAuth Providers**: Google, Discord
- **API Tokens**: Configurable via settings UI and environment variables
