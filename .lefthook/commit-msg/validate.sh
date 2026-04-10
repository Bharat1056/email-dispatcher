#!/usr/bin/env bash

commit_msg_file="$1"
commit_msg=$(head -1 "$commit_msg_file")

# Conventional commit pattern: type: description
pattern="^(feat|fix|chore|docs|style|refactor|test|ci|perf|build): .+"

if ! echo "$commit_msg" | grep -qE "$pattern"; then
    echo ""
    echo "ERROR: Invalid commit message format!"
    echo ""
    echo "  Got: $commit_msg"
    echo ""
    echo "  Expected format: type: description"
    echo "  Allowed types: feat, fix, chore, docs, style, refactor, test, ci, perf, build"
    echo ""
    echo "  Examples:"
    echo "    feat: add user authentication"
    echo "    fix: resolve database connection timeout"
    echo "    chore: update dependencies"
    echo "    docs: update API documentation"
    echo "    style: format code with prettier"
    echo "    refactor: simplify user service logic"
    echo "    test: add unit tests for login module"
    echo "    ci: update GitHub Actions workflow"
    echo "    perf: improve query performance"
    echo "    build: update webpack configuration"
    echo ""
    exit 1
fi
