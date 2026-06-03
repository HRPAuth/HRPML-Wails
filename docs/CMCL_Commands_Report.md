# CMCL Commands Analysis Report

**Tool:** Console Minecraft Launcher (CMCL)  
**Version:** v2.2.2  
**Analysis Date:** 2026-05-18

---

## Overview

**CMCL (Console Minecraft Launcher)** is a command-line interface tool for managing Minecraft Java Edition. It provides comprehensive functionality for version management, account handling, mod installation, and launcher configuration.

---

## 1. Basic Commands

| Command | Short Option | Description | Example |
|---------|--------------|-------------|---------|
| `cmcl [<version>]` | - | Start Minecraft with selected/specified version | `cmcl 1.19` |
| `--help` | `-h` | Display help documentation | `cmcl -h` |
| `--list[=<game dir>]` | `-l` | List all game versions | `cmcl -l` |
| `--print[=<version>]` | `-p` | Print startup command | `cmcl -p 1.19` |
| `--select=<version>` | `-s` | Select default version | `cmcl -s 1.18.2` |
| `--about` | `-a` | Show about information | `cmcl -a` |
| `--check-for-updates` | `-c` | Check for CMCL updates | `cmcl -c` |

---

## 2. Version Management (`install` function)

### 2.1 Install Options

| Option | Description | Example |
|--------|-------------|---------|
| `-n, --name` | Custom storage name | `cmcl install 1.18.2 -n my1.18.2` |
| `-s, --select` | Select after installation | `cmcl install 1.18.2 -s` |
| `--fabric[=<version>]` | Install Fabric | `cmcl install 1.19 --fabric=0.14.12` |
| `--forge[=<version>]` | Install Forge | `cmcl install 1.19 --forge` |
| `--liteloader[=<version>]` | Install LiteLoader | `cmcl install 1.12 --liteloader=1.12-SNAPSHOT` |
| `--optifine[=<version>]` | Install OptiFine | `cmcl install 1.19 --optifine` |
| `--quilt[=<version>]` | Install Quilt | `cmcl install 1.19 --quilt` |
| `-t, --thread` | Download thread count | `cmcl install 1.19 --thread=96` |
| `--no-assets` | Skip assets download | `cmcl install 1.19 --no-assets` |
| `--no-libraries` | Skip libraries download | `cmcl install 1.19 --no-libraries` |

### 2.2 Show Installable Versions

```bash
cmcl install --show=<version type> [-t <time range>]
```

**Version Types:**
- `a/all` - All versions
- `r/release` - Release versions
- `s/snapshot` - Snapshot versions
- `oa/oldAlpha` - Old Alpha versions
- `ob/oldBeta` - Old Beta versions

**Example:** `cmcl install --show=s --time=2020-05-09/2021-10-23`

---

## 3. Version Operations (`version` function)

| Option | Description | Example |
|--------|-------------|---------|
| `--info` | View version information | `cmcl version 1.19 --info` |
| `-d, --delete` | Delete version | `cmcl version -d` |
| `--rename=<name>` | Rename version | `cmcl version --rename=myversion` |
| `--complete[=assets|libraries|natives]` | Complete missing files | `cmcl version 1.19 --complete` |
| `--config=<name> [<content>]` | Set version-specific config | `cmcl version 1.19 --config maxMemory 2048` |
| `--fabric/forge/liteloader/optifine/quilt` | Install mod loader | `cmcl version 1.19 --fabric --api` |
| `--isolate` | Set working directory isolation | `cmcl version 1.19 --isolate` |
| `--unset-isolate` | Unset isolation | `cmcl version 1.19 --unset-isolate` |
| `-p, --print-command` | Print startup command | `cmcl version 1.19 -p` |
| `--export-script=<file>` | Export launch script (bat/sh) | `cmcl version --export-script=start.sh` |
| `--export-script-ps=<file>` | Export PowerShell script | `cmcl version --export-script-ps=start.ps1` |

---

## 4. Account Management (`account` function)

### 4.1 Account Operations

| Option | Description | Example |
|--------|-------------|---------|
| `-s, --select=<order>` | Select account | `cmcl account -s 3` |
| `-l, --list` | List all accounts | `cmcl account --list` |
| `-d, --delete=<order>` | Delete account | `cmcl account --delete=4` |
| `-r, --refresh` | Refresh selected account | `cmcl account --refresh` |
| `--cape[=<path>]` | Set/unset custom cape | `cmcl account --cape=cape.png` |
| `--download-skin=<path>` | Download skin file | `cmcl account --download-skin=skin.png` |
| `--skin[=steve|alex|<path>]` | Set skin | `cmcl account --skin=steve` |

### 4.2 Login Methods

| Method | Description | Required Parameters |
|--------|-------------|-------------------|
| `offline` | Offline account | `-n, --name=<player name>` |
| `microsoft` | Microsoft account | None |
| `authlib` | Authlib-injector | `--address=<server address>` |
| `nide8auth` | Nide8Auth | `--serverId=<server ID>` |

**Examples:**
```bash
cmcl account --login=offline --name=Alexander
cmcl account --login=microsoft
cmcl account --login=authlib --address=127.0.0.1
```

---

## 5. Configuration Management (`config` function)

| Option | Description | Example |
|--------|-------------|---------|
| `<config name> [<content>]` | Get/set configuration | `cmcl config javaPath` / `cmcl config maxMemory 2048` |
| `-a, --all` | Output all config | `cmcl config -a` |
| `--getRaw[=<indent>]` | Output raw config | `cmcl config --getRaw` |
| `-d, --delete=<name>` | Delete config item | `cmcl config -d javaPath` |
| `-c, --clear` | Clear all config | `cmcl config --clear` |
| `-v, --view` | View settable configs | `cmcl config -v` |

---

## 6. JVM Arguments (`jvmArgs` function)

| Option | Description | Example |
|--------|-------------|---------|
| `-p, --print[=<indent>] [-v=<version>]` | Print arguments | `cmcl jvmArgs -p2 -v1.19` |
| `-a, --add=<content> [-v=<version>]` | Add argument | `cmcl jvmArgs --add="-Dfile.encoding=UTF-8"` |
| `-d, --delete=<order> [-v=<version>]` | Delete argument | `cmcl jvmArgs --delete=2` |

---

## 7. Game Arguments (`gameArgs` function)

| Option | Description | Example |
|--------|-------------|---------|
| `-p, --print[=<indent>] [-v=<version>]` | Print arguments | `cmcl gameArgs --print` |
| `-a, --add=<name> [<content>] [-v=<version>]` | Add argument | `cmcl gameArgs -a width 256` |
| `-d, --delete=<name> [-v=<version>]` | Delete argument | `cmcl gameArgs --delete=demo` |

---

## 8. Mod Management (`mod` function)

**Supported Sources:** CurseForge (`cf`), Modrinth (`mr`)

| Option | Description | Example |
|--------|-------------|---------|
| `--url=<url>` | Download mod from URL | `cmcl mod --url=https://...` |
| `--install --name=<name>` | Search and install by name | `cmcl mod --install -nMinis --limit=30` |
| `--install --id=<id>` | Install by ID | `cmcl mod --install --source=curseforge --id=297344` |
| `--info --name=<name>` | Search and display info | `cmcl mod --info -n Sodium` |
| `--info --id=<id>` | Display mod info by ID | `cmcl mod --info --source=mr --id=GBeCx05I` |

---

## 9. Modpack Management (`modpack` function)

**Supported Sources:** CurseForge (`cf`), Modrinth (`mr`)

| Option | Description | Example |
|--------|-------------|---------|
| `--url=<url>` | Download from URL | `cmcl modpack --url=https://...` |
| `--file=<path>` | Install local modpack | `cmcl modpack --file=modpack.zip` |
| `--install --name=<name>` | Search and install | `cmcl modpack --install -nRLCraft --storage="New Game"` |
| `--install --id=<id>` | Install by ID | `cmcl modpack --install --source=curseforge --id=285109` |
| `--info --name=<name>` | Search and display info | `cmcl modpack --info -n "Sugar Optimization"` |

**Additional Options:**
- `-k, --keep-file` - Keep downloaded file
- `--storage=<name>` - Custom storage name

---

## 10. Simplified Commands (`simplify` function)

| Option | Description | Example |
|--------|-------------|---------|
| `-p, --print` | View all simplified commands | `cmcl simplify -p` |
| `-s, --set=<cmd> "<original>"` | Create shortcut | `cmcl simplify -s ds2 "config downloadSource 2"` |
| `-d, --delete=<cmd>` | Delete shortcut | `cmcl simplify -d ds2` |

---

## Function Summary

| Function | Description |
|----------|-------------|
| `install` | Install Minecraft versions with mod loaders |
| `version` | Manage installed versions |
| `account` | Handle player accounts and authentication |
| `config` | Configure launcher settings |
| `simplify` | Create command shortcuts |
| `jvmArgs` | Custom JVM startup arguments |
| `gameArgs` | Custom game arguments |
| `mod` | Search and install mods from CurseForge/Modrinth |
| `modpack` | Search and install modpacks from CurseForge/Modrinth |

---

## Conclusion

CMCL provides a comprehensive command-line interface for Minecraft management with the following key capabilities:

1. **Version Management**: Install, delete, rename, and manage Minecraft versions with support for multiple mod loaders (Fabric, Forge, LiteLoader, OptiFine, Quilt)
2. **Account Management**: Support for offline, Microsoft, Authlib-injector, and Nide8Auth accounts with skin/cape customization
3. **Configuration**: Fine-grained control over launcher settings and per-version configurations
4. **Customization**: JVM and game argument management for performance tuning
5. **Mod/Modpack Management**: Integrated search and installation from CurseForge and Modrinth
6. **Quality of Life**: Simplified commands for frequently used operations

The tool follows standard CLI conventions with both short (`-h`) and long (`--help`) options, and supports flexible input formats for options (space-separated, equals sign, or direct attachment for short options).
