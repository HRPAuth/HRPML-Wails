# HRPAuth 后端 API 端点分析报告

## 0. 根路径 API (Backend Status)

### 请求入口
- **URL**: `/status`
- **请求方法**: GET
- **路由处理**: 内联处理 (`public/index.php`)

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
无

### 处理操作
1. 返回 Backend 状态信息

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 用途 |
|--------|----------|------|
| 200 | JSON: `{"status": "online", "backend": {"name": "string", "url": "string", "version": "string", "php_version": "string", "server_time": "string"}, "message": "string"}` | 返回后端状态 |

## 0.1 Yggdrasil 元数据 API

### 请求入口
- **URL**: `/`
- **请求方法**: GET
- **路由处理**: `modules/zggdrasilapi/src/meta.php`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
无

### 处理操作
1. 返回 Yggdrasil API 服务器元数据

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值
```json
{
  "meta": {
    "serverName": "string",
    "implementationName": "string",
    "implementationVersion": "string",
    "links": {
      "homepage": "string",
      "register": "string"
    },
    "feature.non_email_login": "boolean",
    "feature.legacy_skin_api": "boolean",
    "feature.no_mojang_namespace": "boolean",
    "feature.enable_mojang_anti_features": "boolean",
    "feature.enable_profile_key": "boolean",
    "feature.username_check": "boolean"
  },
  "skinDomains": ["string"],
  "signaturePublickey": "string"
}
```

## 1. 登录 API

### 请求入口
- **URL**: `/login`
- **请求方法**: POST
- **路由处理**: `controllers/AuthController@login`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱地址 |
| password | string | 是 | 用户密码 |

### 文件系统与数据库操作
- **数据库操作**:
  - 查询用户信息：`SELECT uid, password, totp FROM users WHERE email = ? LIMIT 1`
  - 更新用户 token：`UPDATE users SET remember_token = ? WHERE uid = ?`

### 处理操作
1. 验证请求方法是否为 POST
2. 解析 JSON 请求数据
3. 验证邮箱格式
4. 查询用户信息（包含 totp 字段）
5. 验证密码
6. 生成随机 token
7. 更新用户 token 到数据库
8. 判断用户是否配置了 TOTP（检查 totp 字段是否有内容）
9. 返回登录结果（包含 totp 状态）

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "message": "Login successful", "token": "string", "uid": "number", "totp": "number"}` | - | 登录成功，返回用户 token、uid 和 TOTP 状态 |
| 400 | - | `{"success": false, "message": "Invalid email"}` | 邮箱格式错误 |
| 401 | - | `{"success": false, "message": "Email or password incorrect"}` | 邮箱或密码错误 |
| 405 | - | `{"success": false, "message": "Method Not Allowed"}` | 请求方法错误 |

### 返回字段说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 是否登录成功 |
| message | string | 操作结果消息 |
| token | string | 用户登录令牌，用于后续请求认证 |
| uid | number | 用户唯一标识符 |
| totp | number | TOTP 状态：`1` 表示用户已配置 TOTP，`0` 表示未配置 |

## 2. 注册 API

### 请求入口
- **URL**: `/register`
- **请求方法**: POST
- **路由处理**: `controllers/AuthController@register`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱地址 |
| username | string | 是 | 用户名 |
| password | string | 是 | 用户密码 |

### 文件系统与数据库操作
- **数据库操作**:
  - 检查邮箱是否存在：`SELECT uid FROM users WHERE email = ? LIMIT 1`
  - 创建用户：`INSERT INTO users (email, username, score, password, ip, last_sign_at, register_at, verified, verification_token) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

### 处理操作
1. 验证请求方法是否为 POST
2. 解析 JSON 请求数据
3. 验证邮箱格式
4. 验证用户名长度
5. 验证密码长度
6. 连接数据库
7. 检查邮箱是否已注册
8. 密码哈希处理
9. 插入用户数据
10. 返回注册结果

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "uid": "number", "message": "Register successful"}` | - | 注册成功，返回用户 uid |
| 400 | - | `{"success": false, "message": "Invalid email"}` | 邮箱格式错误 |
| 400 | - | `{"success": false, "message": "Username too short"}` | 用户名长度不足 |
| 400 | - | `{"success": false, "message": "Password too short"}` | 密码长度不足 |
| 409 | - | `{"success": false, "message": "Email already registered"}` | 邮箱已注册 |
| 405 | - | `{"success": false, "message": "Method Not Allowed"}` | 请求方法错误 |
| 500 | - | `{"success": false, "message": "Database error"}` | 数据库错误 |

## 3. 邮件验证 API

### 请求入口
- **URL**: `/email-verification`
- **请求方法**: POST
- **路由处理**: `controllers/EmailVerificationController@handle`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
根据不同的 action 参数，请求参数有所不同：

#### 3.1 发送测试邮件 (`action=send-test-email`)
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| action | string | 是 | 固定值 "send-test-email" |
| to | string | 是 | 收件人邮箱地址 |
| subject | string | 是 | 邮件主题 |
| message | string | 是 | 邮件内容 |

#### 3.2 发送验证码 (`action=send-verification-code`)
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| action | string | 是 | 固定值 "send-verification-code" |
| email | string | 是 | 用户邮箱地址 |

#### 3.3 验证验证码 (`action=verify-code`)
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| action | string | 是 | 固定值 "verify-code" |
| email | string | 是 | 用户邮箱地址 |
| code | string | 是 | 验证码 |

### 文件系统与数据库操作
- **文件系统操作**:
  - 发送邮件：使用 fsockopen 连接 SMTP 服务器
- **数据库操作**:
  - 更新用户验证状态：`UPDATE users SET verified = 1 WHERE email = ?`
- **缓存操作**:
  - 存储验证码：`storeVerificationCode($email, $code)`
  - 获取验证码：`getVerificationCode($email)`
  - 删除验证码：`deleteVerificationCode($email)`

### 处理操作
1. 验证请求方法是否为 POST
2. 解析 JSON 请求数据
3. 根据 action 参数执行不同操作：
   - **send-test-email**: 验证参数，发送测试邮件
   - **send-verification-code**: 验证邮箱，检查是否已发送验证码，生成验证码，存储验证码，发送邮件
   - **verify-code**: 验证邮箱和验证码，删除验证码，更新用户验证状态
4. 返回操作结果

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "message": "Email sent successfully", "data": {"to": "string", "subject": "string"}}` | - | 测试邮件发送成功 |
| 200 | `{"success": true, "message": "Verification code sent successfully"}` | - | 验证码发送成功 |
| 200 | `{"success": true, "message": "Verification successful"}` | - | 验证码验证成功 |
| 400 | - | `{"success": false, "message": "Recipient email cannot be empty"}` | 测试邮件收件人邮箱为空 |
| 400 | - | `{"success": false, "message": "Invalid recipient email format"}` | 测试邮件收件人邮箱格式错误 |
| 400 | - | `{"success": false, "message": "Email subject cannot be empty"}` | 测试邮件主题为空 |
| 400 | - | `{"success": false, "message": "Email content cannot be empty"}` | 测试邮件内容为空 |
| 400 | - | `{"success": false, "message": "Invalid email"}` | 邮箱格式错误 |
| 400 | - | `{"success": false, "message": "Verification code is required"}` | 验证码为空 |
| 400 | - | `{"success": false, "message": "Verification code expired or not found"}` | 验证码过期或不存在 |
| 400 | - | `{"success": false, "message": "Invalid verification code"}` | 验证码错误 |
| 400 | - | `{"success": false, "message": "Invalid action"}` | 无效的 action |
| 404 | - | `{"success": false, "message": "User not found or already verified"}` | 用户不存在或已验证 |
| 405 | - | `{"success": false, "message": "Method Not Allowed"}` | 请求方法错误 |
| 429 | - | `{"success": false, "message": "Verification code already sent, please wait"}` | 验证码已发送，请等待 |
| 500 | - | `{"success": false, "message": "Failed to store verification code"}` | 存储验证码失败 |
| 500 | - | `{"success": false, "message": "Failed to update verification status"}` | 更新验证状态失败 |
| 500 | - | `{"success": false, "message": "错误信息"}` | 邮件发送失败 |

## 4. 获取用户信息 API

### 请求入口
- **URL**: `/user`
- **请求方法**: POST
- **路由处理**: `controllers/UserController@getUser`

### 请求值类型
- **Content-Type**: `application/json` (也支持表单数据或 URL 参数)

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| remember_token | string | 是 | 用户登录 token |
| uid | string | 否 | 用户唯一标识符 |
| email | string | 否 | 用户邮箱地址 |

### 文件系统与数据库操作
- **数据库操作**:
  - 查询用户信息：`SELECT uid, email, username, avatar, verified FROM users WHERE remember_token = ? [AND uid = ?] [AND email = ?]`
  - 回退查询（无 avatar 字段）：`SELECT uid, email, username, verified FROM users WHERE remember_token = ? [AND uid = ?] [AND email = ?]`

### 处理操作
1. 启动会话
2. 连接数据库
3. 从不同来源获取参数：
   - POST 请求的 JSON 数据
   - POST 请求的表单数据
   - URL 参数
4. 验证 token 是否为空
5. 构建查询条件，根据提供的参数添加 uid 和 email 比对
6. 尝试查询用户信息（带 avatar 字段）
7. 如果失败，尝试回退查询（无 avatar 字段）
8. 验证用户是否存在
9. 构建用户数据
10. 返回用户信息

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "message": "获取用户信息成功", "data": {"uid": "number", "email": "string", "username": "string", "avatar": "string|null", "verified": "boolean"}}` | - | 获取用户信息成功 |
| 200 | - | `{"success": false, "message": "未登录或登录已过期"}` | token 为空 |
| 200 | - | `{"success": false, "message": "用户不存在或token无效"}` | 用户不存在或 token 无效 |
| 200 | - | `{"success": false, "message": "服务器错误"}` | 数据库错误 |

## 5. 登出 API

### 请求入口
- **URL**: `/logout`
- **请求方法**: GET
- **路由处理**: `controllers/AuthController@logout`

### 请求值类型
- **Content-Type**: 无特定要求

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| remember_token | string | 否 | 用户登录 token |

### 文件系统与数据库操作
- **数据库操作**:
  - 清除用户 token：`UPDATE users SET remember_token = NULL WHERE remember_token = ?`
- **会话操作**:
  - 清除会话数据：`$_SESSION = []`
  - 销毁会话：`session_destroy()`

### 处理操作
1. 启动会话
2. 从不同来源获取 token：
   - POST 请求的 JSON 数据
   - POST 请求的表单数据
   - URL 参数
3. 如果 token 存在，清除数据库中的 token
4. 清除会话数据
5. 销毁会话
6. 重定向到登录页面

### 返回值类型
- **Content-Type**: 无特定返回值，执行重定向

### 期望的返回值用途
- 执行登出操作并重定向到登录页面

## 6. 测试用户 API（开发用）

### 请求入口
- **URL**: `/test-user`
- **请求方法**: GET
- **路由处理**: 内联处理 (`public/index.php`)

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
无

### 处理操作
1. 加载 UserController
2. 调用 getUser 方法

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | - | - | 测试用端点 |

## 7. TOTP 生成 API

### 请求入口
- **URL**: `/totpgen`
- **请求方法**: GET
- **路由处理**: `controllers/TOTPController@generate`

### 请求值类型
- **Content-Type**: 无特定要求

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| secret | string | 是 | TOTP 密钥 |

### 文件系统与数据库操作
- **无**

### 处理操作
1. 验证 secret 参数是否存在
2. 生成 TOTP 验证码
3. 返回验证码

### 返回值类型
- **Content-Type**: `text/plain`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `字符串类型的 6 位数字验证码` | - | 生成 TOTP 验证码 |
| 400 | - | `Missing secret` | secret 参数缺失 |

## 7.1 TOTP 设置 API

### 请求入口
- **URL**: `/totp/setup`
- **请求方法**: POST
- **路由处理**: `controllers/TOTPController@setupTOTP`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱地址 |
| remtoken | string | 是 | 用户登录 token |

### 文件系统与数据库操作
- **数据库操作**:
  - 查询用户信息：`SELECT uid FROM users WHERE email = ? AND remember_token = ? LIMIT 1`
  - 更新 TOTP 密钥：`UPDATE users SET totp = ? WHERE uid = ?`

### 处理操作
1. 验证请求方法是否为 POST
2. 解析 JSON 请求数据
3. 验证 email 和 remtoken 参数
4. 验证邮箱格式
5. 连接数据库
6. 验证用户凭据（email + remember_token）
7. 生成 20 字符的 base32 编码 TOTP 密钥
8. 将密钥存储到数据库 users 表的 totp 字段
9. 返回结果

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "totpkey": "string"}` | - | TOTP 密钥生成成功 |
| 400 | - | `{"success": false, "message": "Missing email or remtoken"}` | 参数缺失 |
| 400 | - | `{"success": false, "message": "Invalid email"}` | 邮箱格式错误 |
| 401 | - | `{"success": false, "message": "Invalid email or remtoken"}` | 用户认证失败 |
| 405 | - | `{"success": false, "message": "Method Not Allowed"}` | 请求方法错误 |
| 500 | - | `{"success": false, "message": "Database connection error"}` | 数据库连接错误 |

## 7.2 TOTP 验证 API

### 请求入口
- **URL**: `/totp/verify`
- **请求方法**: POST
- **路由处理**: `controllers/TOTPController@verifyTOTP`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| email | string | 是 | 用户邮箱地址 |
| passcode | string | 是 | 用户输入的 TOTP 验证码（6位数字） |

### 文件系统与数据库操作
- **数据库操作**:
  - 查询用户信息：`SELECT uid, totp, remember_token FROM users WHERE email = ? LIMIT 1`

### 处理操作
1. 验证请求方法是否为 POST
2. 解析 JSON 请求数据
3. 验证 email 和 passcode 参数
4. 验证邮箱格式
5. 连接数据库
6. 查询用户及其 TOTP 密钥
7. 验证用户是否存在且已设置 TOTP
8. 根据当前时间和上一时间窗口生成预期验证码
9. 验证 passcode（支持当前窗口和上一窗口，容忍时间偏移）
10. 验证成功返回用户 remember_token

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "email": "string", "rt": "string"}` | - | TOTP 验证成功 |
| 400 | - | `{"success": false, "message": "Missing email or passcode"}` | 参数缺失 |
| 400 | - | `{"success": false, "message": "Invalid email"}` | 邮箱格式错误 |
| 401 | - | `{"success": false, "message": "User not found or TOTP not configured"}` | 用户不存在或未配置 TOTP |
| 401 | - | `{"success": false, "message": "Invalid passcode"}` | 验证码错误 |
| 405 | - | `{"success": false, "message": "Method Not Allowed"}` | 请求方法错误 |
| 500 | - | `{"success": false, "message": "Database connection error"}` | 数据库连接错误 |

### 备注
- 支持当前时间窗口和上一时间窗口的验证码（容忍约30秒的时间偏移）
- 验证成功时返回用户的 remember_token 可用于后续业务操作

## 8. 修改用户名 API

### 请求入口
- **URL**: `/change-username`
- **请求方法**: POST
- **路由处理**: `controllers/ChangeUsernameController@changeUsername`

### 请求值类型
- **Content-Type**: `application/json` (也支持表单数据或 URL 参数)

### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| remember_token | string | 是 | 用户登录 token |
| username | string | 是 | 新用户名 |

### 文件系统与数据库操作
- **数据库操作**:
  - 查询用户信息：`SELECT uid, uuid FROM users WHERE remember_token = ?`
  - 检查用户名是否被占用：`SELECT uid FROM users WHERE username = ? AND uid != ?`
  - 更新用户名：`UPDATE users SET username = ? WHERE uid = ?`

### 处理操作
1. 启动会话
2. 连接数据库
3. 从不同来源获取参数：
   - POST 请求的 JSON 数据
   - POST 请求的表单数据
   - URL 参数
4. 验证 token 是否为空
5. 验证新用户名是否提供
6. 验证用户名格式（长度 3-16 字符，只允许字母、数字和下划线）
7. 根据 token 查询用户信息
8. 检查新用户名是否已被其他用户占用
9. 更新用户名
10. 返回操作结果

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "message": "用户名修改成功", "data": {"username": "string"}}` | - | 修改用户名成功 |
| 200 | - | `{"success": false, "message": "未登录或登录已过期"}` | token 为空或无效 |
| 200 | - | `{"success": false, "message": "请提供新用户名"}` | 未提供新用户名 |
| 200 | - | `{"success": false, "message": "用户名不能为空"}` | 用户名为空 |
| 200 | - | `{"success": false, "message": "用户名长度必须在3-16个字符之间"}` | 用户名长度不符合要求 |
| 200 | - | `{"success": false, "message": "用户名只能包含字母、数字和下划线"}` | 用户名格式不符合要求 |
| 200 | - | `{"success": false, "message": "用户不存在或token无效"}` | 用户不存在 |
| 200 | - | `{"success": false, "message": "该用户名已被使用"}` | 用户名已被占用 |
| 200 | - | `{"success": false, "message": "服务器错误"}` | 数据库错误 |

## 9. 密钥生成 API

### 请求入口
- **URL**: `/generate-key`
- **请求方法**: POST
- **路由处理**: `controllers/KeyGenController@generate`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
无

### 文件系统与数据库操作
- **文件系统操作**:
  - 创建密钥目录：`mkdir('keys', 0700, true)`
  - 保存公钥：`file_put_contents('keys/public.pem', $publicKey)`
  - 保存私钥：`file_put_contents('keys/private.pem', $privateKey)`
  - 设置文件权限：`chmod('keys/*.pem', 0600)`
- **配置文件操作**:
  - 禁用端点：将 `preference.php` 中的 `'enable' => 0` 替换为 `'enable' => 1`

### 处理操作
1. 加载配置文件 (`config/preference.php`)
2. 检查端点是否已禁用 (`keygen.enable`)
3. 生成 RSA 2048-bit 密钥对
4. 保存密钥到 `keys/` 目录
5. 自动禁用端点（修改配置文件）
6. 返回生成结果

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值用途
| 状态码 | 成功响应 | 失败响应 | 用途 |
|--------|----------|----------|------|
| 200 | `{"success": true, "message": "Key pair generated successfully", "data": {"public_key": "string"}}` | - | 密钥生成成功 |
| 403 | - | `{"success": false, "message": "Key generation endpoint is disabled"}` | 端点已禁用 |
| 500 | - | `{"success": false, "message": "Config file not found"}` | 配置文件不存在 |
| 500 | - | `{"success": false, "message": "Failed to generate key pair"}` | 密钥生成失败 |
| 500 | - | `{"success": false, "message": "Failed to save keys"}` | 密钥保存失败 |

### 备注
- 此端点用于生成 Yggdrasil API 签名所需的 RSA 密钥对
- 端点在首次成功调用后会自动禁用，防止重复生成
- 公钥用于 Yggdrasil 元数据中的 `signaturePublickey` 字段
- 私钥用于签名玩家材质数据

## 10. ZggdrasilAPI 元数据 API

### 请求入口
- **URL**: `/` (YggdrasilAPI)
- **请求方法**: GET
- **路由处理**: `modules/zggdrasilapi/src/meta.php`

### 请求值类型
- **Content-Type**: `application/json`

### 请求参数
无

### 处理操作
1. 加载配置（从 preference.php 和 zggdrasilapi.php）
2. 返回 Yggdrasil API 服务器元数据

### 返回值类型
- **Content-Type**: `application/json`

### 期望的返回值
```json
{
  "meta": {
    "serverName": "string",
    "implementationName": "string",
    "implementationVersion": "string",
    "links": {
      "homepage": "string",
      "register": "string"
    },
    "feature.non_email_login": "boolean",
    "feature.legacy_skin_api": "boolean",
    "feature.no_mojang_namespace": "boolean",
    "feature.enable_mojang_anti_features": "boolean",
    "feature.enable_profile_key": "boolean",
    "feature.username_check": "boolean"
  },
  "skinDomains": ["string"],
  "signaturePublickey": "string"
}
```

### 备注
- `serverName`, `implementation`, `version` 从 `preference.php` 动态加载
- `links` 的 homepage 和 register URL 从 `preference.php` 的 frontend URL 动态生成
- `skinDomains` 从 `preference.php` 的 callback URL 提取域名
- `signaturePublickey` 从 `keys/public.pem` 文件加载
- Feature flags 现在使用 `feature.` 前缀格式，位于 `meta` 对象内部

## 11. ZggdrasilAPI 认证相关 API

### 9.1 认证 API

#### 请求入口
- **URL**: `/authserver/authenticate`
- **请求方法**: POST
- **路由处理**: `modules/zggdrasilapi/src/auth/authenticate.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| username | string | 是 | 用户邮箱或用户名（取决于配置） |
| password | string | 是 | 用户密码 |
| agent | object | 是 | 游戏代理信息 |
| clientToken | string | 否 | 客户端 token |
| requestUser | boolean | 否 | 是否返回用户信息 |

#### 返回值
成功响应：
```json
{
  "accessToken": "string",
  "clientToken": "string",
  "availableProfiles": [
    {
      "id": "string (uuid)",
      "name": "string",
      "model": "string (optional)"
    }
  ],
  "selectedProfile": {
    "id": "string (uuid)",
    "name": "string",
    "model": "string (optional)"
  },
  "user": {
    "id": "string (uuid)",
    "email": "string",
    "username": "string",
    "properties": []
  }
}
```

失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid credentials."
}
```

### 9.2 刷新 Token API

#### 请求入口
- **URL**: `/authserver/refresh`
- **请求方法**: POST
- **路由处理**: `modules/zggdrasilapi/src/auth/refresh.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| accessToken | string | 是 | 访问 token |
| clientToken | string | 否 | 客户端 token |
| requestUser | boolean | 否 | 是否返回用户信息 |
| selectedProfile | object | 否 | 选择的档案 |

#### 返回值
成功响应：
```json
{
  "accessToken": "string",
  "clientToken": "string",
  "selectedProfile": {
    "id": "string (uuid)",
    "name": "string",
    "model": "string (optional)"
  },
  "user": {
    "id": "string (uuid)",
    "email": "string",
    "properties": []
  }
}
```

失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid token."
}
```

### 9.3 验证 Token API

#### 请求入口
- **URL**: `/authserver/validate`
- **请求方法**: POST
- **路由处理**: `modules/zggdrasilapi/src/auth/validate.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| accessToken | string | 是 | 访问 token |
| clientToken | string | 否 | 客户端 token |

#### 返回值
成功响应：204 No Content
失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid token."
}
```

### 9.4 使 Token 失效 API

#### 请求入口
- **URL**: `/authserver/invalidate`
- **请求方法**: POST
- **路由处理**: `modules/zggdrasilapi/src/auth/invalidate.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| accessToken | string | 是 | 访问 token |
| clientToken | string | 否 | 客户端 token |

#### 返回值
成功响应：204 No Content
失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid token."
}
```

### 9.5 登出 API

#### 请求入口
- **URL**: `/authserver/signout`
- **请求方法**: POST
- **路由处理**: `modules/zggdrasilapi/src/auth/signout.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| username | string | 是 | 用户邮箱 |
| password | string | 是 | 用户密码 |

#### 返回值
成功响应：204 No Content
失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid credentials."
}
```

## 12. ZggdrasilAPI 会话相关 API

### 10.1 加入会话 API

#### 请求入口
- **URL**: `/sessionserver/session/minecraft/join`
- **请求方法**: POST
- **路由处理**: `modules/zggdrasilapi/src/session/join.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| accessToken | string | 是 | 访问 token |
| selectedProfile | string | 是 | 选择的档案 UUID |
| serverId | string | 是 | 服务器 ID |

#### 返回值
成功响应：204 No Content
失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid request."
}
```

### 10.2 检查加入状态 API

#### 请求入口
- **URL**: `/sessionserver/session/minecraft/hasJoined?username={username}&serverId={serverId}&ip={ip}&unsigned={unsigned}`
- **请求方法**: GET
- **路由处理**: `modules/zggdrasilapi/src/session/hasJoined.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| serverId | string | 是 | 服务器 ID |
| ip | string | 否 | 客户端 IP |
| unsigned | boolean | 否 | 是否不使用签名 |

#### 返回值
成功响应：
```json
{
  "id": "string (uuid)",
  "name": "string",
  "properties": []
}
```

失败响应：204 No Content 或错误响应

### 10.3 查询玩家档案 API

#### 请求入口
- **URL**: `/sessionserver/session/minecraft/profile/{uuid}?unsigned={unsigned}`
- **请求方法**: GET
- **路由处理**: `modules/zggdrasilapi/src/profile/profileQuery.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| uuid | string | 是 | 档案 UUID（路径参数） |
| unsigned | boolean | 否 | 是否不使用签名 |

#### 返回值
成功响应：
```json
{
  "id": "string (uuid)",
  "name": "string",
  "properties": []
}
```

失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Profile not found."
}
```

## 13. ZggdrasilAPI 档案批量查询 API

### 请求入口
- **URL**: `/api/profiles/minecraft`
- **请求方法**: POST
- **路由处理**: `modules/zggdrasilapi/src/profile/batchProfiles.php`

### 请求参数
- 请求体为用户名数组，最多 100 个用户名

### 返回值
成功响应：
```json
[
  {
    "id": "string (uuid)",
    "name": "string"
  }
]
```

失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid request."
}
```

## 14. ZggdrasilAPI 材质相关 API

### 12.1 上传材质 API

#### 请求入口
- **URL**: `/api/user/profile/{uuid}/{textureType}`
- **请求方法**: PUT
- **路由处理**: `modules/zggdrasilapi/src/texture/uploadTexture.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| uuid | string | 是 | 档案 UUID（路径参数） |
| textureType | string | 是 | 材质类型（skin/cape，路径参数） |
| Authorization | string | 是 | Bearer 访问 token（请求头） |
| file | file | 是 | PNG 图片文件（表单数据） |
| model | string | 否 | 皮肤模型（default/slim，仅用于 skin 类型） |

#### 文件系统与数据库操作
- **文件系统操作**:
  - 创建材质目录：`/public/textures/{uuid}/`
  - 保存材质文件：`/public/textures/{uuid}/{textureType}.png`
- **数据库操作**:
  - 查询档案归属：`SELECT id FROM profiles WHERE id = ? AND user_id = ?`
  - 更新材质属性：`UPDATE profile_properties SET value = ? WHERE id = ?`
  - 新增材质属性：`INSERT INTO profile_properties (profile_id, name, value) VALUES (?, 'textures', ?)`

#### 处理操作
1. 验证 UUID 和材质类型（必须为 skin 或 cape）
2. 验证 Authorization 请求头，提取 Bearer token
3. 验证 token 有效性（查询 tokens 表）
4. 验证档案是否属于当前用户
5. 验证上传文件（PNG 格式，大小 ≤ 100KB）
6. 创建 `{uuid}` 目录（如不存在）
7. 将上传文件保存到 `/public/textures/{uuid}/{textureType}.png`
8. 生成材质 URL：`{callback_url}/textures/{uuid}/{textureType}.png`
9. 构建材质载荷并 Base64 编码
10. 更新或插入 profile_properties 表中的 textures 属性

#### 材质文件存储位置
- **路径**: `/public/textures/{uuid}/{textureType}.png`
- **示例**: `/public/textures/550e8400e29b41d4a71644665530a028/skin.png`

#### 材质 URL 格式
- **格式**: `{callback_url}/textures/{uuid}/{textureType}.png`
- **示例**: `https://hrpauth.samuelcheston.com/textures/550e8400e29b41d4a71644665530a028/skin.png`

#### 返回值
成功响应：204 No Content
失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid UUID."
}
```

### 12.2 删除材质 API

#### 请求入口
- **URL**: `/api/user/profile/{uuid}/{textureType}`
- **请求方法**: DELETE
- **路由处理**: `modules/zggdrasilapi/src/texture/deleteTexture.php`

#### 请求参数
| 参数名 | 类型 | 必须 | 描述 |
|--------|------|------|------|
| uuid | string | 是 | 档案 UUID（路径参数） |
| textureType | string | 是 | 材质类型（skin/cape，路径参数） |
| Authorization | string | 是 | Bearer 访问 token（请求头） |

#### 文件系统与数据库操作
- **文件系统操作**:
  - 删除材质文件：`/public/textures/{uuid}/{textureType}.png`
- **数据库操作**:
  - 查询材质属性：`SELECT id, value FROM profile_properties WHERE profile_id = ? AND name = 'textures'`
  - 更新材质属性：`UPDATE profile_properties SET value = ? WHERE id = ?`
  - 删除材质属性：`DELETE FROM profile_properties WHERE id = ?`

#### 处理操作
1. 验证 UUID 和材质类型（必须为 skin 或 cape）
2. 验证 Authorization 请求头，提取 Bearer token
3. 验证 token 有效性（查询 tokens 表）
4. 验证档案是否属于当前用户
5. 查询 profile_properties 中的 textures 属性
6. 解码并更新 textures 载荷，移除对应的材质
7. 如无剩余材质则删除属性，否则更新属性
8. 删除物理文件 `/public/textures/{uuid}/{textureType}.png`

#### 返回值
成功响应：204 No Content
失败响应：
```json
{
  "error": "ForbiddenOperationException",
  "errorMessage": "Invalid UUID."
}
```

### 材质目录结构
```
public/textures/
└── {uuid}/
    ├── skin.png    # 玩家皮肤文件
    └── cape.png    # 玩家披风文件
```

## 总结

本项目提供了以下 API 端点：

### 主系统 API（public/index.php）
1. **根路径 API** (`GET /status`) - 返回后端状态信息
2. **登录 API** (`POST /login`) - 用于用户登录，返回 token 和 uid
3. **注册 API** (`POST /register`) - 用于用户注册，返回 uid
4. **邮件验证 API** (`POST /email-verification`) - 用于发送测试邮件、发送验证码和验证验证码
5. **获取用户信息 API** (`POST /user`) - 用于获取用户信息
6. **登出 API** (`GET /logout`) - 用于用户登出
7. **测试用户 API** (`GET /test-user`) - 开发调试用端点
8. **TOTP 生成 API** (`GET /totpgen`) - 用于生成 TOTP 验证码
9. **修改用户名 API** (`POST /change-username`) - 用于修改用户名
10. **密钥生成 API** (`POST /generate-key`) - 用于生成 Yggdrasil API 签名密钥对

### ZggdrasilAPI（Minecraft 认证协议兼容）

#### 认证相关
11. **认证 API** (`POST /authserver/authenticate`) - 用户认证
12. **刷新 Token API** (`POST /authserver/refresh`) - 刷新访问令牌
13. **验证 Token API** (`POST /authserver/validate`) - 验证令牌有效性
14. **使 Token 失效 API** (`POST /authserver/invalidate`) - 使令牌失效
15. **登出 API** (`POST /authserver/signout`) - 用户登出

#### 会话相关
16. **加入会话 API** (`POST /sessionserver/session/minecraft/join`) - 加入游戏会话
17. **检查加入状态 API** (`GET /sessionserver/session/minecraft/hasJoined`) - 检查玩家是否加入会话
18. **查询玩家档案 API** (`GET /sessionserver/session/minecraft/profile/{uuid}`) - 查询玩家档案信息

#### 档案相关
19. **批量查询档案 API** (`POST /api/profiles/minecraft`) - 批量查询玩家档案

#### 材质相关
20. **上传材质 API** (`PUT /api/user/profile/{uuid}/{textureType}`) - 上传玩家皮肤或披风
21. **删除材质 API** (`DELETE /api/user/profile/{uuid}/{textureType}`) - 删除玩家皮肤或披风

所有 API 端点都遵循 RESTful 设计原则，使用 JSON 格式返回数据（除了 TOTP 生成 API 返回纯文本）。数据库操作主要涉及用户表的查询和更新，文件系统操作主要是发送邮件。

这些 API 端点共同构成了一个完整的用户认证系统，支持用户注册、登录、邮箱验证、获取用户信息和登出等功能，同时提供了 TOTP 验证码生成功能以增强安全性。ZggdrasilAPI 模块提供了 Minecraft 官方认证协议的兼容实现，支持游戏内会话验证和玩家材质管理。

### 路由结构

项目使用前端控制器模式，通过 `public/index.php` 统一处理所有请求，并根据路由配置将请求转发到相应的控制器方法。当路由不匹配时，请求会被转发到 `modules/zggdrasilapi/index.php` 处理 Minecraft 认证协议相关的请求：

| 路由 | 方法 | 控制器/处理文件 |
|------|------|-----------------|
| `/status` | GET | 内联处理（返回后端状态） |
| `/login` | POST | `AuthController@login` |
| `/register` | POST | `AuthController@register` |
| `/logout` | GET | `AuthController@logout` |
| `/user` | POST | `UserController@getUser` |
| `/test-user` | GET | 内联处理（调试） |
| `/email-verification` | POST | `EmailVerificationController@handle` |
| `/totpgen` | GET | `TOTPController@generate` |
| `/change-username` | POST | `ChangeUsernameController@changeUsername` |
| `/generate-key` | POST | `KeyGenController@generate` |
| `/` | GET | `zggdrasilapi/src/meta.php` |
| `/authserver/authenticate` | POST | `zggdrasilapi/src/auth/authenticate.php` |
| `/authserver/refresh` | POST | `zggdrasilapi/src/auth/refresh.php` |
| `/authserver/validate` | POST | `zggdrasilapi/src/auth/validate.php` |
| `/authserver/invalidate` | POST | `zggdrasilapi/src/auth/invalidate.php` |
| `/authserver/signout` | POST | `zggdrasilapi/src/auth/signout.php` |
| `/sessionserver/session/minecraft/join` | POST | `zggdrasilapi/src/session/join.php` |
| `/sessionserver/session/minecraft/hasjoined` | GET | `zggdrasilapi/src/session/hasJoined.php` |
| `/sessionserver/session/minecraft/profile/{uuid}` | GET | `zggdrasilapi/src/profile/profileQuery.php` |
| `/api/profiles/minecraft` | POST | `zggdrasilapi/src/profile/batchProfiles.php` |
| `/api/user/profile/{uuid}/{textureType}` | PUT/DELETE | `zggdrasilapi/src/texture/uploadTexture.php` / `deleteTexture.php` |

这种路由结构使得 API 端点更加清晰和规范，便于维护和扩展。
