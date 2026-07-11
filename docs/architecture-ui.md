# Arquitetura — `ob ui`

TUI fullscreen para o openbench, reutilizando toda a lógica de negócio existente em `internal/*`.

## Objetivo

Oferecer uma experiência integrada no terminal (estilo lazygit/k9s) sem substituir a CLI atual. Comandos como `ob commit`, `ob docker up` e `ob pr` continuam funcionando para scripts, CI e agentes.

```
┌─ OpenBench UI ─────────────────────────────────────────────────────┐
│ OPENBENCH v0.1.x     owner/repo · feat/my-branch · docker:ok       │
├────────────────────────────────────────────────────────────────────┤
│ Environment                                                        │
│  Docker running · compose.yaml · app running :8080→80              │
├────────────────────────────────────────────────────────────────────┤
│ Git Graph · Changed files · Recent commits · AI Engine             │
├────────────────────────────────────────────────────────────────────┤
│ Next: [U] up [D] down [c] commit [p] push [P] pr [L] dlogs [q] quit│
└────────────────────────────────────────────────────────────────────┘
```

## Princípios

| Princípio | Descrição |
|-----------|-----------|
| **Lógica separada da apresentação** | `internal/app`, `internal/git`, `internal/docker`, `internal/ai`, `internal/pr` não importam Bubble Tea |
| **CLI intacta** | `internal/ui` (ANSI/wizard) permanece para comandos one-shot |
| **Snapshot read-only** | Dashboard carrega `WorkspaceSnapshot` sem efeitos colaterais |
| **Ações via app layer** | Commit/push/PR/Docker chamam `app.Run*` existentes |
| **Progresso desacoplado** | Interface `Progress` permite spinner na TUI e texto na CLI |

## Camadas

```
cmd/ob/main.go
       │
       ├─► ob (overview CLI) ──► app.RunOverview()
       ├─► ob commit/pr/docker/... ──► app.Run*()
       └─► ob ui              ──► tui.Run()
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    ▼                     ▼                     ▼
             app.LoadWorkspace      app.RunCommit         internal/git
             Snapshot()             app.RunDockerUp       internal/docker
                    │                     │               internal/ai
                    ▼                     ▼
             tui/views/dashboard    tui/views/action
             (read-only)            (modal + Progress)
```

### Pacotes

```
internal/
├── app/
│   ├── workspace.go      # WorkspaceSnapshot — dados do dashboard
│   ├── progress.go       # interface Progress (CLI + TUI)
│   ├── runner.go         # RunCommit, RunPush, RunPR (inalterados)
│   └── suggestions.go    # buildNextSteps (regras de negócio)
├── ui/                   # CLI ANSI (Session, Wizard) — legado
└── tui/                  # Bubble Tea
    ├── run.go            # entry point: tui.Run()
    ├── app.go            # root Model (roteamento de telas)
    ├── state.go          # Screen enum, AppState global
    ├── keys.go           # keymap centralizado
    ├── styles.go         # lipgloss (tema Gitai)
    ├── progress.go       # Progress → status bar / modal
    ├── components/
    │   ├── statusbar.go
    │   ├── branchlist.go
    │   ├── filelist.go
    │   └── help.go
    └── views/
        ├── dashboard.go  # tela principal
        ├── diff.go       # viewport com git diff
        ├── commit.go     # preview + confirmação
        ├── pr.go         # preview PR + draft toggle
        └── report.go     # uso/custo (fase 2)
```

## Modelo de dados

### `WorkspaceSnapshot` (`internal/app/workspace.go`)

Agrega tudo que o dashboard precisa em uma única leitura:

```go
type WorkspaceSnapshot struct {
    Overview  *git.Overview
    OpenPR    *pr.PRView   // nil se gh ausente ou sem PR
    Config    *config.Config
    ConfigErr error
    NextSteps []NextStep
    HasGH     bool
}
```

Carregamento:

1. Validar repositório git
2. `repo.Overview(baseBranch)` — já existe
3. `pr.ViewCurrent()` — opcional, best-effort
4. `config.Load()` — pode falhar (usuário não configurou)
5. `buildNextSteps()` — regras já em `suggestions.go`

Refresh: tecla `R` ou timer opcional (30s) dispara novo `LoadWorkspaceSnapshot`.

### `NextStep` (exportado de `suggestions.go`)

```go
type NextStep struct {
    Command string
    Note    string
    Plain   bool
    Action  ActionID // novo: mapeia tecla → ação na TUI
}
```

## Máquina de estados (Bubble Tea)

```
                    ┌─────────────┐
         Init ─────►│  Dashboard  │◄──── refresh (R)
                    └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         ▼                 ▼                 ▼
   ┌──────────┐     ┌──────────┐      ┌──────────┐
   │   Diff   │     │  Commit  │      │    PR    │
   │ (viewport)│     │ (modal)  │      │ (modal)  │
   └──────────┘     └────┬─────┘      └────┬─────┘
                         │                  │
                         ▼                  ▼
                   app.RunCommit      app.RunPR
                   (goroutine)       (goroutine)
                         │                  │
                         ▼                  ▼
                   ┌──────────┐
                   │  Result  │ ──► volta ao Dashboard
                   └──────────┘
```

### Root `Model`

```go
type Model struct {
    screen    Screen
    snapshot  *app.WorkspaceSnapshot
    width     int
    height    int
    loading   bool
    err       error

    dashboard views.DashboardModel
    diff      views.DiffModel
    action    views.ActionModel  // commit/push/pr em progresso
}
```

Comandos assíncronos (padrão Bubble Tea):

```go
type snapshotLoadedMsg struct { snap *app.WorkspaceSnapshot; err error }
type actionDoneMsg     struct { result *app.Result; err error }
```

Ações longas (IA) rodam em `tea.Cmd` com goroutine; a UI mostra spinner na status bar.

## Interface `Progress`

Desacopla `app.Run*` da saída textual:

```go
// internal/app/progress.go
type Progress interface {
    Step(label string, fn func() error) error
    StepQuiet(fn func() error) error
    Detail(msg string)
    Info(msg string)
    Success(msg string)
}
```

| Implementação | Uso |
|---------------|-----|
| `ui.Session` | CLI (`ob commit`, etc.) — já existe, adapter fino |
| `tui.Progress` | Atualiza status bar + log lateral no modal |

Migração incremental: `Options.Progress Progress` substitui `Options.UI *ui.Session` quando preenchido.

## Keymap

| Tecla | Ação | Condição |
|-------|------|----------|
| `q` / `Ctrl+C` | Sair | sempre |
| `r` | Refresh snapshot | dashboard |
| `d` | Ver diff staged/branch | arquivos alterados |
| `c` | Commit com IA | working tree dirty |
| `p` | Push | ahead ou dirty |
| `P` | PR | commits ahead of base |
| `s` | Sync | behind > 0 |
| `o` | Abrir PR no browser | PR aberto |
| `?` | Ajuda | sempre |
| `↑/↓` | Navegar listas | branches / files |
| `Enter` | Selecionar arquivo → diff | file list |

Teclas derivadas de `NextSteps` — cada step com `ActionID` mapeia para handler.

## Dependências

```
github.com/charmbracelet/bubbletea   # framework TUI
github.com/charmbracelet/lipgloss    # estilos
github.com/charmbracelet/bubbles     # list, viewport, spinner, help
```

Sem dependências novas em `internal/app`, `internal/git`, `internal/ai`.

## Entry points

| Comando | Comportamento |
|---------|---------------|
| `ob ui` | Abre TUI (explícito) |
| `ob` | Abre TUI quando `interactive_ui: true` e terminal interativo |
| `OB_NO_UI=1` | Força overview CLI (sobrescreve config) |
| `NO_COLOR=1` | Sem cores — convenção Unix (sobrescreve `ui_color`) |

Detecção de terminal: `term.IsTerminal` + `OB_NO_UI` + `CI` — mesma lógica de `ui.Session`.

## Fases de implementação

### Fase 1 — Dashboard (MVP)

- [x] `WorkspaceSnapshot` + `LoadWorkspaceSnapshot`
- [x] `ob ui` com layout básico
- [x] Branches, files, commits, next steps, status bar
- [x] Refresh (`r`) e quit (`q`)

### Fase 2 — Ações

- [x] Interface `Progress` + adapter TUI
- [x] Modal de commit (preview mensagem IA → confirmar)
- [x] Modal de PR (preview body → draft toggle → criar)
- [x] Diff viewer com `bubbles/viewport`

### Fase 3 — Polish

- [x] Tela de report/usage
- [x] Temas (dark/light via `NO_COLOR`)
- [x] `ob` default → TUI quando TTY interativo
- [x] Testes de keymap e snapshot

## Testes

| Camada | Estratégia |
|--------|------------|
| `app/workspace.go` | Unit test com git repo fixture (como `overview_test.go`) |
| `app/suggestions.go` | Já testado em `suggestions_test.go` |
| `tui/keys.go` | Tabela de tecla → ação |
| `tui/views/*` | `tea.NewProgram(model, tea.WithInput(nil))` — smoke sem terminal |

Evitar testes frágeis de renderização pixel-a-pixel.

## Riscos e mitigações

| Risco | Mitigação |
|-------|-----------|
| Terminal pequeno (< 80×24) | Layout mínimo com scroll; mensagem se muito pequeno |
| IA lenta bloqueando UI | Goroutine + spinner; cancel com `Ctrl+C` no modal |
| Duplicação CLI/TUI | Snapshot e `Run*` compartilhados; zero lógica git na TUI |
| Conflito com wizard `config` | `ob config` permanece CLI; link `?` → abre hint |
