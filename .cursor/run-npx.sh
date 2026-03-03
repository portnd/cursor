#!/bin/bash
# Wrapper so MCP can find npx when Cursor doesn't inherit terminal PATH.
set -e
export PATH="/usr/local/bin:/opt/homebrew/bin:$HOME/.volta/bin:$PATH"

# nvm: source and add default node to PATH
if [ -s "$HOME/.nvm/nvm.sh" ]; then
  export NVM_DIR="$HOME/.nvm"
  . "$NVM_DIR/nvm.sh" 2>/dev/null || true
  # nvm may not set PATH in non-interactive; add default node bin if present
  for nvmdir in "$HOME/.nvm/versions/node"/v*/bin; do
    [ -x "$nvmdir/npx" ] && export PATH="$nvmdir:$PATH" && break
  done
fi
# fnm
if [ -d "$HOME/.local/share/fnm" ]; then
  export PATH="$HOME/.local/share/fnm/aliases/default/bin:$PATH"
  [ -d "$HOME/.local/share/fnm/node-versions" ] && for d in "$HOME/.local/share/fnm/node-versions"/*/installation/bin; do [ -x "$d/npx" ] && export PATH="$d:$PATH" && break; done
fi

# Find npx: try PATH first, then known locations
NPX=""
if command -v npx &>/dev/null; then
  NPX="npx"
elif [ -x "/opt/homebrew/bin/npx" ]; then
  NPX="/opt/homebrew/bin/npx"
elif [ -x "/usr/local/bin/npx" ]; then
  NPX="/usr/local/bin/npx"
else
  for d in "$HOME/.nvm/versions/node"/v*/bin "$HOME/.local/share/fnm/aliases/default/bin" "$HOME/.volta/bin"; do
    [ -x "$d/npx" ] && NPX="$d/npx" && break
  done
fi

if [ -z "$NPX" ]; then
  echo "npx not found. Install Node.js: brew install node  or from https://nodejs.org" >&2
  exit 1
fi
exec "$NPX" -y "@agent-infra/mcp-server-browser@latest" "$@"
