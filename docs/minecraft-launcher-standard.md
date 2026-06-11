# 执行摘要  
本文深入研究 Minecraft Java 版的启动流程和命令行规则，涵盖官方启动器的命令行选项和环境变量，以及游戏目录、资源文件、版本 JSON、身份验证令牌等机制。我们比较了不同启动场景下的差异，重点说明通过 **Authlib-Injector** 注入认证、**离线模式**（无认证服务器）、**原版干净客户端**、**安装 Forge 模组** 和 **安装 Fabric 模组** 时的具体要求。文中将详细列出各场景下所需和可选的 JVM 参数与启动器参数，给出 Bash 示例命令，并讨论常见陷阱和安全注意事项。最后提供对比表格和流程图（Mermaid）梳理启动器→认证→JVM→游戏的执行顺序。所有示例使用 Linux Bash 语法，如有 Windows 差异会特别说明。  

## 官方启动器命令行选项与环境变量  
Minecraft 官方启动器本身也支持一些命令行选项。新启动器（版本 2.x）支持如 `--demo`（演示模式）和 `--workDir <路径>`（指定 `.minecraft` 目录位置）等参数。例如，`--workDir .` 可让启动器使用当前目录作为游戏目录。启动器还支持 `--clean`（清理游戏和运行时目录）等选项。不过，这些是启动器本身的参数，不影响游戏进程的 JVM 启动参数。  

真正启动游戏时，启动器会构造完整的 `java` 调用：指定 `-cp` 类路径、`-Djava.library.path` 指向本地库目录、主类（一般是 `net.minecraft.client.main.Main`）以及游戏参数。例如，官方命令可能类似：  

```bash
java -Djava.library.path=/path/to/natives \
     -cp "/path/to/libs/*" \
     net.minecraft.client.main.Main \
     --username 玩家名 --version 1.20.4 --gameDir "/home/user/.minecraft" \
     --assetsDir "/home/user/.minecraft/assets" --assetIndex 1.20.4 \
     --uuid <UUID> --accessToken <TOKEN> --userType mojang --versionType release
```  

其中 `-Djava.library.path` 指向解压后的本地文件（natives），类路径包括所有库 Jar。主类 `Main` 位于 `com.mojang:minecraft` 包中。游戏参数（如 `--username`、`--uuid`、`--accessToken` 等）由版本 JSON 和认证返回的会话数据提供。需要注意，启动器在 Windows 下会将本地库解压到 `%TEMP%` 目录，需通过 `-Djava.library.path` 指定；在 Linux 下一般指定 `natives` 文件夹。环境变量方面，**系统级的 `_JAVA_OPTIONS` 可能影响启动**（如强制内存设置），应确保其不会干扰 Minecraft。  

## 1) Authlib-Injector 注入认证  
**解释：** Authlib-Injector 是一个代理库，允许自建登录服务器（Yggdrasil API）替代 Mojang 官方服务器。通过将其作为 Java Agent 加载，游戏进程中请求验证时会重定向到指定地址。这对于使用私有账号库或测试服务器非常有用。注入时需要指定代理 Jar 路径和目标服务器 URL，例如 `-javaagent:/path/authlib-injector.jar=https://my.auth.server`。  

**参数要求：** JVM 参数中必须包含 `-javaagent`：  
- `-javaagent:/绝对/路径/authlib-injector.jar=<认证服务器 URL>`。例如 `-javaagent:/home/user/libs/authlib-injector.jar=http://localhost:8080`。  
- 其他 JVM 参数可照常设置（内存、GC 等）。  
游戏程序参数基本同在线模式：必须提供 `--username`、`--uuid`、`--accessToken` 等，只是这些值应来自自定义认证服务器接口。格式与官方相同（UUID 不带短横线，accessToken 通常为 20+ 位十六进制字符串）。可以使用 `--userType mojang` 或其他未迁移账户类型。  

**示例命令（Linux Bash）：**  
```bash
java -javaagent:/opt/auth/authlib-injector.jar=https://auth.example.com \
     -Djava.library.path=./natives \
     -cp "./libraries/*:versions/1.20.4/1.20.4.jar" \
     net.minecraft.client.main.Main \
     --username OfflineUser --version 1.20.4 \
     --gameDir "/home/user/.minecraft" --assetsDir "/home/user/.minecraft/assets" \
     --assetIndex 1.20.4 --uuid 0123456789abcdef0123456789abcdef \
     --accessToken AB12CD34EF56GH78IJ90KL12MN34OP56 \
     --userType mojang
```  
请将 `--accessToken` 和 `--uuid` 替换为从自有认证服务器获取的有效值。  

**环境变量：** 无特殊环境变量，仅需确保系统能找到 Java 可执行文件。可以通过 `JAVA_HOME` 或修改 PATH 来指定 Java。  

**类路径组装：** 类路径与正常启动相同，包含所有版本 JSON 定义的库。Authlib-Injector 本身无需加入类路径，只需作为 Java Agent 加载。  

**注意事项与安全：** Authlib-Injector 会改变客户端与官方认证服务器的交互，因此要保证**认证服务可信且加密**。访问 token 仍然是敏感信息，不应在日志或命令行中泄露。使用 HTTP 时注意风险，推荐 HTTPS。由于游戏内外的 UUID 是从客户端提交的，使用伪造的 UUID 可能导致与服务端数据不匹配。  

## 2) 离线模式启动  
**解释：** 离线模式下，游戏不连接 Mojang 认证服务器，玩家可使用任意用户名本地进入游戏。此模式常用于无网络环境或测试，但多人游戏（非离线模式服务器）中无法区分真实玩家身份。离线启动依旧需要提供伪造的会话信息给游戏进程。  

**参数要求：** 与正常模式类似，但生成“离线”令牌：  
- `--username`：离线用户名（可随意，但与游戏内角色对应）。  
- `--uuid`：根据用户名生成的离线 UUID（可通过 `UUID.nameUUIDFromBytes("OfflinePlayer:"+用户名)` 算出，或使用任意固定值）。  
- `--accessToken`：随意字符串（例如多次启动可用固定随机值）。PolyMC 等启动器会生成随机令牌。  
- `--userType legacy`（旧版账户）或 `--userType offline`。官方示例多用 `legacy`。  
- 其他参数与版本相关：`--version`、`--gameDir`、`--assetsDir`、`--assetIndex` 等与普通启动相同。  
例如：  
```bash
java -Djava.library.path=./natives \
     -cp "./libraries/*:versions/1.20.4/1.20.4.jar" \
     net.minecraft.client.main.Main \
     --username MyOfflineName --version 1.20.4 \
     --gameDir "/home/user/.minecraft" --assetsDir "/home/user/.minecraft/assets" \
     --assetIndex 1.20.4 --uuid 9f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c \
     --accessToken OFFLINETOKEN1234567890 --userType legacy
```  

**环境变量：** 同正常模式，无需额外变量。  

**类路径组装：** 与原版相同。使用版本 JSON 中的库和游戏 Jar 即可。  

**注意事项与安全：** 离线模式无身份验证，任何人都可用相同用户名（及离线 UUID）登录，因此只能用于信任的单机环境。不要在公网服务器上启用离线模式，否则容易被冒名顶替。离线令牌和 UUID 不涉及敏感信息，但仍建议不要在命令行历史中泄露习惯用户名和 UUID 信息。  

## 3) 原版（干净客户端）启动  
**解释：** 原版启动即官方推荐的在线模式启动。启动器先通过 Mojang/Yggdrasil 服务获取 `accessToken`、`clientToken` 以及玩家的 `uuid`。然后构造命令行启动游戏主类 `net.minecraft.client.main.Main`。整个流程中访问 Mojang API 严格按照 [wiki.vg 验证流程](https://wiki.vg/Authentication) 进行，令牌通常为十六进制字符串。  

**参数要求：** 必需的游戏参数包括：  
- `--username <用户名>`：在线账户名。  
- `--uuid <玩家 UUID>`：从 Mojang 返回的角色 UUID（不带短横线）。  
- `--accessToken <访问令牌>`：从 Mojang 登录接口获得的令牌。  
- `--userType mojang`（或微软账户类型 `msa`），`--versionType release`（或版本类型，如 `release`、`snapshot` 等）。  
- `--gameDir <游戏目录>`：游戏主目录（默认为 `~/.minecraft`）。  
- `--assetsDir <资源目录>`：资源文件目录（默认为 `<gameDir>/assets`）。  
- `--assetIndex <资源版本>`：一般与游戏版本号相同。  
- JVM 参数如 `-Xmx`, `-Xms` 可根据需要指定堆内存大小。**必须包含** `-Djava.library.path` 指向本地库（natives）目录。  
示例：  
```bash
java -Xms1G -Xmx2G -Djava.library.path=./natives \
     -cp "./libraries/*:versions/1.20.4/1.20.4.jar" \
     net.minecraft.client.main.Main \
     --username Steve --version 1.20.4 \
     --gameDir "/home/user/.minecraft" \
     --assetsDir "/home/user/.minecraft/assets" \
     --assetIndex 1.20.4 \
     --uuid 123e4567e89b12d3a456426655440000 \
     --accessToken a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6 \
     --userType mojang --versionType release
```  
（其中 `--uuid` 和 `--accessToken` 需要通过官方认证接口获得。）  

**环境变量：** 使用系统中可用的 Java，通常设置 `JAVA_HOME` 或修改 PATH 以使用正确的 Java 版本。Java 版本应符合官方要求（目前约为 Java 17 以上）。  

**类路径组装：** 所有版本 JSON 列出的库都要加入类路径，包括 `com.mojang:authlib`、`com.mojang:realms` 等。可以使用脚本或库自动解析 `~/.minecraft/versions/1.20.4/1.20.4.json` 生成类路径。千万**不要**手动省略任何必要库，否则会缺少类和崩溃。  

**常见问题与注意：** Windows 下类路径分隔符使用 `;`，Linux/macOS 用 `:`。确保 `natives` 目录存在且与 `-Djava.library.path` 路径匹配，否则游戏将无法加载本地库。启动器版本更新时会改变版本目录名（如 `1.20.4`），要及时调整命令。由于用户令牌和 UUID 来自安全 API，绝不要硬编码用户名/密码信息，令牌若泄露可能被他人使用登录。  

## 4) 带 Forge 模组的启动  
**解释：** Forge 是常用模组加载器，它扩展了官方启动流程。1.12 及以前版本使用 LaunchWrapper，需要在命令行中添加 Tweaker 以激活 Forge。新版本 Forge 已迁移到 Mojang 的新的启动系统（Fabric 类似模式），但仍需要指定一些参数。总体而言，启动时仍调用游戏主类，只是会加载 Forge 的核心代码。  

**参数要求（以 Forge 1.12 为例）：**  
- 主类变为 `net.minecraft.launchwrapper.Launch`（由 Forge 提供）。  
- 需要加入 `--tweakClass net.minecraftforge.fml.common.launcher.FMLTweaker` 参数。这告诉 Minecraft 使用 Forge Tweaker 进行启动。  
- 添加 `--versionType Forge` 或其他版本标识。  
- 其余参数与原版相同：`--username`、`--uuid`、`--accessToken`、`--gameDir`、`--assetsDir` 等。  
- 可能还有 Forge 特殊参数（最新版本可能自动添加）：如 `--launchTarget forgeclient`、`--fml.forgeVersion <版本>`、`--fml.mcVersion <MC版本>` 等，具体见版本 JSON 或安装脚本。  
示例（Forge 1.12.2）：  
```bash
java -Djava.library.path=./natives \
     -cp "./libraries/*:forge-1.12.2.jar" \
     net.minecraft.launchwrapper.Launch \
     --username Alex --version 1.12.2 \
     --gameDir "/home/user/.minecraft" \
     --assetsDir "/home/user/.minecraft/assets" \
     --assetIndex 1.12 \
     --uuid 123e4567e89b12d3a456426655440000 \
     --accessToken abcdef0123456789abcdef0123456789 \
     --userType mojang \
     --tweakClass net.minecraftforge.fml.common.launcher.FMLTweaker \
     --versionType Forge
```  
其中 `forge-1.12.2.jar` 是 Forge 安装生成的客户端 Jar（包含修改）。在新版本 Forge/ForgeGradle 逻辑中，主类可能自动设置为 `Launch`，Tweaker 参数由版本 JSON 管理，无需手动指定。  

**环境变量：** 与原版相同，无额外要求。  

**类路径组装：** 类路径需包含 Minecraft 的所有库**加上** Forge 核心库（通常版本目录下的 `forge-<ver>-universal.jar`）。具体库列表请参考对应 `forge-<ver>.json` 文件，它通常会继承并追加额外库。不要遗漏 `net.minecraftforge:launchwrapper` 等依赖。  

**注意事项：** Forge 在命令行中增加了 Tweaker，使得服务器端验证会话时仍使用原生流程。确保使用真实可用的 `accessToken`。不同 Forge 版本 Tweaker 路径可能有变化（详见对应文档或版本 JSON）。与原版启动一样，确保本地库路径正确、采用相同的 `--uuid` 和 `--accessToken` 格式。由于 Forge 本身不提供命令行参数校验，错误使用 Tweaker 名称或遗漏必须的参数会导致游戏无法启动（典型表现为 “Could not find or load main class” 或 crash）。  

## 5) 带 Fabric 模组的启动  
**解释：** Fabric 是另一种轻量级模组加载器。Fabric 加载器安装后，启动时会使用其专用主类（如 **Knot**）替代 Minecraft 原始主类。该主类会初始化 Fabric 环境并再调用 Minecraft 代码。Fabric 启动流程与原版类似，但需要使用 Fabric 的 Loader Jar。  

**参数要求：**  
- 主类变为 Fabric Loader 提供的类：较新版本常用 `net.fabricmc.loader.launch.knot.KnotClient`（Fabric Loader 0.7+）或 `net.fabricmc.loader.launch.FabricClient`（Fabric 老版本）。  
- 其他参数同原版：`--username`、`--uuid`、`--accessToken`、`--gameDir`、`--assetsDir`、`--assetIndex` 等。  
- 无需 `--tweakClass` 参数。Fabric 会自行处理 mod 加载。  
示例（Fabric 1.20.4）：  
```bash
java -Djava.library.path=./natives \
     -cp "./libraries/*:fabric-loader-0.14.23.jar:game.jar" \
     net.fabricmc.loader.launch.knot.KnotClient \
     --username Herobrine --version 1.20.4 \
     --gameDir "/home/user/.minecraft" \
     --assetsDir "/home/user/.minecraft/assets" \
     --assetIndex 1.20.4 \
     --uuid 123e4567e89b12d3a456426655440000 \
     --accessToken deadbeefdeadbeefdeadbeefdeadbeef \
     --userType mojang --versionType release
```  
其中 `fabric-loader-0.14.23.jar` 是 Fabric 安装器生成的 Loader Jar，`game.jar` 是原版的 `1.20.4.jar` 或重命名后的游戏 Jar。KnotClient 主类加载后会查找并注入所有 Fabric 模组。  

**环境变量：** 同原版，无特别要求。  

**类路径组装：** 类路径需要包含游戏所有库和 Fabric Loader Jar。Fabric 的版本 JSON 会指定 Loader 和库的位置，一般包括 `com.mojang:authlib`、Fabric Loader 本身及其依赖。可参考 `versions/1.20.4-fabric.json`。不要遗漏 Fabric Loader Jar。  

**注意事项：** Fabric 启动时不要混用错误版本的 Loader 和游戏主 Jar，否则会 `NoClassDefFoundError`。确保 Loader 版本与游戏版本兼容。Fabric 不依赖 LaunchWrapper Tweaker，因此省去了许多参数，但类路径设置更复杂。访问令牌依旧由 Mojang 管理，Fabric 不更改该流程。  

## 参数对比表  

| 参数/场景      | Authlib-Injector 注入认证 | 离线模式     | 原版在线       | Forge     | Fabric      |
|---------------|-------------------------|------------|--------------|----------|------------|
| `-javaagent`  | `是：-javaagent:/path/authlib-injector.jar=<认证服>` | 无         | 无           | 无       | 无         |
| `--accessToken` | 自定义服务器提供令牌    | 任意离线令牌   | 来自 Mojang 官方 | 来自 Mojang | 来自 Mojang |
| `--uuid`      | 自定义服务器提供 UUID   | 离线 UUID（如从 `OfflinePlayer:` 计算） | 官方返回 UUID  | 官方返回 UUID  | 官方返回 UUID  |
| `--userType`  | `mojang`（或自定义类型）| `legacy`（或 `offline`） | `mojang`（或 `msa`） | `mojang`（或 `msa`） | `mojang`（或 `msa`） |
| `--versionType` | 通常 `release`        | `release`  | `release`    | `Forge` | `release`  |
| `--tweakClass` | 无                    | 无         | 无           | `net.minecraftforge.fml.common.launcher.FMLTweaker` | 无         |
| 主类 (Main)   | `net.minecraft.client.main.Main` | `net.minecraft.client.main.Main` | `net.minecraft.client.main.Main` | `net.minecraft.launchwrapper.Launch` | `net.fabricmc.loader.launch.knot.KnotClient` |
| 类路径        | 原版库 + 注入器 Jar 不在类路径 | 原版库      | 原版库         | 原版库 + Forge 核心 Jar | 原版库 + Fabric Loader Jar |
  
## 启动流程时序图（Mermaid）  

```mermaid
flowchart LR
  A[启动器启动] --> B[认证流程]
  B --> C{在线模式?}
  C -- 是 --> D[调用 Mojang/Yggdrasil<br/>获取 accessToken/UUID]
  C -- 否(离线/注入) --> E[跳过官方认证<br/>生成或重定向令牌]
  D --> F[组装 JVM 参数与类路径]
  E --> F
  F --> G[启动 Java 进程]
  G --> H[执行游戏主类(Main)]
```  

图示说明：启动器先进入认证阶段（在线模式则向服务器请求令牌；离线或 Authlib-Injector 模式跳过官方服务器），获得必要的 `accessToken` 和 `UUID` 后，启动器将版本 JSON 中的库和资源加入类路径，组装 JVM 启动参数。最后启动 `java` 进程，执行主类（如 `net.minecraft.client.main.Main`、`Launch` 或 Fabric 的主类），进入游戏。  

**参考资料：** 官方启动器命令行选项来自 Minecraft Wiki；Authlib-Injector 的用法说明来自其文档；原版和 Forge 启动参数示例来自社区讨论；Yggdrasil 验证流程和令牌格式参考 wiki.vg。