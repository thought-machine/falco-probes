#!/bin/bash
set -Eeuo pipefail

clear_remainder="\033[0K"

util::info() {
    printf "$(ansi::resetColor)$(ansi::magentaIntense)💡 %s$(ansi::resetColor)\n" "$@"
}

util::infor() {
    printf "$(ansi::resetColor)$(ansi::magentaIntense)💡 %s$(ansi::resetColor)\n" "$@"
}

util::rinfor() {
    printf "\r$(ansi::resetColor)$(ansi::magentaIntense)💡 %s$(ansi::resetColor)${clear_remainder}" "$@"
}

util::warn() {
  printf "$(ansi::resetColor)$(ansi::yellowIntense)⚠️  %s$(ansi::resetColor)\n" "$@"
}

util::error() {
  printf "$(ansi::resetColor)$(ansi::bold)$(ansi::redIntense)❌ %s$(ansi::resetColor)\n" "$@"
}

util::rerror() {
  printf "\r$(ansi::resetColor)$(ansi::bold)$(ansi::redIntense)❌ %s$(ansi::resetColor)${clear_remainder}\n" "$@"
}

util::success() {
  printf "$(ansi::resetColor)$(ansi::greenIntense)✅ %s$(ansi::resetColor)\n" "$@"
}

util::rsuccess() {
  printf "\r$(ansi::resetColor)$(ansi::greenIntense)✅ %s$(ansi::resetColor)${clear_remainder}\n" "$@"
}

util::retry() {
  "${@}" || sleep 1; "${@}" || sleep 5; "${@}"
}

util::prompt() {
  prompt=$(printf "$(ansi::bold)$()❔ %s [y/N]$(ansi::resetColor)\n" "$@")
  read -p "${prompt}" yn
  case $yn in
      [Yy]* ) ;;
      * ) util::error "Did not receive happy input, exiting."; exit 1;;
  esac
}

util::prompt_skip() {
  prompt=$(printf "$(ansi::bold)$()❔ %s [y/N]$(ansi::resetColor)\n" "$@")
  read -p "${prompt}" yn
  case $yn in
      [Yy]* ) return 0;;
      * ) util::warn "Did not receive happy input, skipping."; return 1;;
  esac
}
