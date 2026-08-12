# PRAMAAN Claude Code Operating Policy

## Workspace
- Repository root: ~/projects/pramaan
- This is a Cosmos SDK blockchain repository.
- Treat consensus code, state schemas, protobuf definitions, genesis, migrations, validator configuration, keys, and Docker configuration as high-risk.

## Default permissions
You may:
- Read files.
- Search code.
- Run read-only discovery commands.
- Run Go tests, build checks, static analysis, and security scans.
- Create audit markdown reports only.

You must ask before:
- Editing any existing source file.
- Editing protobuf definitions or generated protobuf files.
- Editing app wiring, genesis, migrations, keeper state schemas, or authorization logic.
- Running commands that alter Git state.
- Installing packages or changing dependencies.
- Starting, resetting, initializing, or deleting chain data.
- Running Docker commands that create, remove, rebuild, stop, or delete containers, images, volumes, or networks.
- Running commands involving private keys, mnemonics, keyrings, validator keys, or production credentials.

## Never run without explicit approval
- rm -rf
- git reset
- git clean
- git checkout .
- git restore .
- git rebase
- git push
- git commit
- ignite chain serve --reset-once
- pramaand init
- pramaand testnet
- docker compose down -v
- docker system prune
- commands that expose, print, export, or copy private keys or mnemonics

## Audit procedure
1. Inspect Git status before all work.
2. Do not modify code during audit phases.
3. Record exact commands executed and their summarized output in audit reports.
4. Classify findings as:
   - Confirmed static defect
   - Confirmed runtime defect
   - Likely runtime risk
   - Runtime verification required
5. Do not claim mainnet readiness until build, test, static-analysis, and adversarial transaction tests are completed.
