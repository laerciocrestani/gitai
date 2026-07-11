#!/usr/bin/env bash
# Wrapper — prefira install.sh na raiz do repositório.
#
#   curl -fsSL https://raw.githubusercontent.com/laerciocrestani/openbench/main/install.sh | bash
#   ./install.sh
#   ob update

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

usage() {
  cat <<'EOF'
openbench setup

Prefira:
  ./install.sh                 Instalação completa (Go + ob + PATH + config)
  ./uninstall.sh               Remove openbench, config e PATH do instalador
  ob config                    Wizard de configuração
  ob update                    git pull + reinstala

Este script ainda aceita:
  ./scripts/setup.sh install   → ./install.sh
  ./scripts/setup.sh uninstall → ./uninstall.sh
  ./scripts/setup.sh config
  ./scripts/setup.sh update
EOF
}

run_ob() {
  if command -v ob >/dev/null 2>&1; then
    ob "$@"
    return
  fi
  if command -v openbench >/dev/null 2>&1; then
    openbench "$@"
    return
  fi
  if [[ -x "${HOME}/go/bin/openbench" ]]; then
    "${HOME}/go/bin/openbench" "$@"
    return
  fi
  if ! command -v go >/dev/null 2>&1; then
    echo "✗ Go não encontrado. Rode: ./install.sh" >&2
    exit 1
  fi
  (cd "$REPO_ROOT" && go run ./cmd/ob "$@")
}

main() {
  local cmd="${1:-help}"
  case "$cmd" in
    install)   exec "$REPO_ROOT/install.sh" "${@:2}" ;;
    uninstall) exec "$REPO_ROOT/uninstall.sh" "${@:2}" ;;
    config)    run_ob config ;;
    update)    run_ob update ;;
    help|-h|--help) usage ;;
    *) echo "✗ Comando desconhecido: $cmd" >&2; usage >&2; exit 1 ;;
  esac
}

main "$@"
