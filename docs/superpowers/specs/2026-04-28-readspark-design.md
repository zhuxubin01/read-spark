# ReadSpark 付费英语阅读 APP 设计文档

**日期**: 2026-04-28  
**版本**: v1.0  
**状态**: 待实施

---

## 1. 项目概述

### 1.1 产品定位
ReadSpark 是一款面向全年龄段英语学习者的付费订阅制阅读 APP。每天定时更新精选英文内容，涵盖新闻外刊、分级阅读、经典小说、备考材料等多种类型，通过沉浸式阅读体验（查词、翻译、朗读、笔记）帮助用户提升英语能力。

### 1.2 目标用户
- 学生（中小学、大学）
- 职场人士（商务英语、日常阅读）
- 备考党（四六级、考研、雅思、托福）
- 泛英语爱好者（保持阅读习惯）

### 1.3 商业模式
- **订阅制付费**：月费 / 年费解锁全部内容
- **免费试读**：每日提供 1-2 篇免费文章，引导订阅转化

---

## 2. 系统架构

### 2.1 整体架构

```
┌─────────────┐      ┌─────────────┐
│ Android App │      │   iOS App   │
│  Kotlin +   │      │  Swift +    │
│ Jetpack     │      │   SwiftUI   │
│  Compose    │      │             │
└──────┬──────┘      └──────┬──────┘
       │                    │
       │   HTTPS / REST     │
       │   + WebSocket      │
       └────────┬───────────┘
                │
       ┌────────▼──────────┐
       │    API Gateway    │
       │  (Auth / Rate     │
       │     Limiting)     │
       └────────┬──────────┘
                │
   ┌────────────┼────────────┐
   │            │            │
┌──▼───┐   ┌───▼───┐   ┌────▼────┐
│用户  │   │ 内容  │   │  订阅   │
│服务  │   │ 服务  │   │  服务   │
└──┬───┘   └───┬───┘   └────┬────┘
   │            │            │
┌──▼───┐   ┌───▼───┐   ┌────▼────┐
│阅读  │   │ 推送  │   │  运营   │
│数据  │   │ 服务  │   │  后台   │
│服务  │   │       │   │  (Web)  │
└──────┘   └───────┘   └─────────┘
                │
   ┌────────────┼────────────┐
   │            │            │
┌──▼───┐   ┌───▼───┐   ┌────▼────┐
│Postgre│   │ Redis │   │ 对象存储 │
│  SQL   │   │       │   │(OSS/S3) │
└────────┘   └───────┘   └─────────┘
```

### 2.2 技术栈

| 层级 | 技术选型 |
|------|----------|
| Android | Kotlin + Jetpack Compose |
| iOS | Swift + SwiftUI |
| 后端 | Go (Gin/Echo) + GORM |
| 主数据库 | PostgreSQL 16+ |
| 缓存 | Redis 7+ |
| 搜索 | PostgreSQL 全文搜索（初期）→ Elasticsearch（后期）|
| 文件存储 | MinIO / 阿里云 OSS |
| 消息推送 | Firebase Cloud Messaging + APNS |
| 监控 | Prometheus + Grafana |

### 2.3 服务拆分（逻辑隔离，初期可合并部署）

- **API Gateway**：统一入口、JWT 鉴权、限流
- **用户服务**：注册、登录、资料管理、设备管理
- **内容服务**：文章 CRUD、分类管理、定时发布
- **订阅服务**：支付（Apple IAP / Google Play / 微信 / 支付宝）、会员权益、续费管理
- **阅读数据服务**：阅读进度、笔记、划线、生词本、收藏
- **推送服务**：通知发送、定时任务调度

---

## 3. 移动端设计

### 3.1 页面结构

5 个主 Tab：

1. **首页**：每日推荐、每日一句、阅读进度卡片、今日打卡状态
2. **分类**：按内容类型（新闻/小说/备考/分级）和难度等级（A1-C2）浏览
3. **阅读**：最近阅读、我的收藏、阅读历史
4. **发现**：学习数据统计、连续打卡天数、排行榜
5. **我的**：账户信息、订阅状态、生词本、设置

### 3.2 阅读器引擎（核心差异化）

#### Android 实现
- **渲染层**：Jetpack Compose Canvas + `TextLayoutResult`
- **手势识别**：`PointerInput` 处理长按查词、滑动选中文本
- **样式标注**：`Spannable` 实现高亮、下划线
- **分页加载**：`LazyColumn` + 预加载

#### iOS 实现
- **渲染层**：`UIViewRepresentable` 包装 `UITextView` + TextKit 2
- **手势识别**：`UIGestureRecognizer` 系列
- **样式标注**：`NSAttributedString`
- **分页加载**：`UICollectionView` Compositional Layout

#### 核心交互
- **长按查词**：弹出浮层显示释义、发音、收藏按钮
- **滑动选中文本**：弹出操作菜单（翻译 / 划线 / 笔记 / 复制）
- **语音朗读**：段落高亮 + TTS 逐句跟随

#### 阅读器设置
- 字体大小：5 级调节（14sp - 22sp）
- 字体选择：系统默认 / 衬线字体
- 行间距：紧凑 / 标准 / 宽松
- 主题：日间白 / 护眼黄 / 夜间黑 / OLED 纯黑
- 翻页模式：滚动 / 仿真翻页（iOS 优先）
- 双语对照：原文与译文段落对照显示

### 3.3 离线支持
- 文章正文离线缓存（SQLite / Room / CoreData）
- 词典数据包离线（前 5000 高频词库，约 50MB）
- 语音朗读缓存（已播放段落音频本地存储）
- 阅读进度、笔记、生词本：本地优先 + 后台同步

---

## 4. 后端设计

### 4.1 核心数据模型

#### users（用户表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| phone | VARCHAR | 手机号（唯一） |
| email | VARCHAR | 邮箱（可选） |
| nickname | VARCHAR | 昵称 |
| avatar_url | VARCHAR | 头像 URL |
| created_at | TIMESTAMP | 创建时间 |

#### subscriptions（订阅表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| user_id | UUID | 用户 ID |
| plan_type | ENUM | monthly / yearly |
| status | ENUM | active / expired / cancelled |
| start_date | DATE | 开始日期 |
| end_date | DATE | 结束日期 |
| auto_renew | BOOLEAN | 是否自动续费 |

#### articles（文章表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| title | VARCHAR | 标题 |
| summary | TEXT | 摘要 |
| content | TEXT | 正文（Markdown / HTML） |
| translation | TEXT | 译文 |
| category | ENUM | news / fiction / exam / graded |
| difficulty | ENUM | A1 / A2 / B1 / B2 / C1 / C2 |
| word_count | INT | 字数 |
| audio_url | VARCHAR | 音频 URL |
| cover_image | VARCHAR | 封面图 URL |
| is_premium | BOOLEAN | 是否付费 |
| published_at | TIMESTAMP | 发布时间 |

#### reading_progress（阅读进度表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| user_id | UUID | 用户 ID |
| article_id | UUID | 文章 ID |
| position | INT | 当前阅读位置（字符偏移） |
| percentage | FLOAT | 阅读百分比 |
| last_read_at | TIMESTAMP | 最后阅读时间 |

#### annotations（标注表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键 |
| user_id | UUID | 用户 ID |
| article_id | UUID | 文章 ID |
| type | ENUM | highlight / note / vocabulary |
| range_start | INT | 文本起始位置 |
| range_end | INT | 文本结束位置 |
| content | TEXT | 笔记内容（type=note 时） |
| created_at | TIMESTAMP | 创建时间 |

### 4.2 搜索接口化设计

```go
// ArticleSearcher 接口
// 初期实现: PGFullTextSearch
// 后期替换: ElasticSearchSearcher（零业务代码改动）
type ArticleSearcher interface {
    Search(ctx context.Context, query SearchQuery) (SearchResult, error)
    Index(ctx context.Context, article Article) error
    Delete(ctx context.Context, articleID string) error
}

type SearchQuery struct {
    Keyword    string
    Category   *string
    Difficulty *string
    Page       int
    PageSize   int
}

type SearchResult struct {
    Articles []ArticleSummary
    Total    int
}
```

### 4.3 关键 API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/register | 手机号注册 |
| POST | /api/v1/auth/login | 手机号登录 |
| POST | /api/v1/auth/refresh | 刷新 Token |
| GET | /api/v1/articles/daily | 今日推荐列表 |
| GET | /api/v1/articles/:id | 文章详情 |
| GET | /api/v1/articles | 文章列表（筛选/分页） |
| POST | /api/v1/progress | 同步阅读进度 |
| GET | /api/v1/progress | 获取阅读历史 |
| POST | /api/v1/annotations | 创建标注 |
| GET | /api/v1/annotations | 获取标注列表 |
| GET | /api/v1/dictionary/:word | 查词 |
| POST | /api/v1/subscriptions | 创建订阅 |
| GET | /api/v1/subscriptions/status | 查询会员状态 |
| POST | /api/v1/push/token | 注册推送 Token |

### 4.4 定时任务

| 频率 | 任务 | 说明 |
|------|------|------|
| 每日 6:00 | 发布今日文章 | 将预设文章标记为已发布，生成推送队列 |
| 每日 7:00 | 推送通知 | 向订阅用户发送"今日更新"通知 |
| 每小时 | 订阅状态同步 | 同步 Apple/Google 订阅状态，处理续费/过期 |
| 每日凌晨 | 数据清理 | 清理过期缓存，归档老旧数据 |

---

## 5. 内容管理

### 5.1 内容类型
- **新闻外刊**：经济学人、纽约时报等精选改编
- **分级阅读**：按 A1-C2 难度分级的原创/改编材料
- **经典小说**：简写版经典英文小说、短篇故事
- **备考材料**：四六级、考研、雅思、托福真题/模拟题

### 5.2 内容发布流程
1. 运营人员在 Web CMS 后台录入文章（标题、正文、译文、分类、难度、音频）
2. 文章保存为草稿状态，可预览
3. 设定发布时间（或立即发布）
4. 到达发布时间后，后端自动将文章标记为已发布
5. 推送服务向订阅用户发送更新通知

### 5.3 Web CMS 后台（运营使用）
- 文章管理（增删改查、富文本编辑器、定时发布）
- 用户管理（查询、封禁）
- 订阅管理（查看订单、处理退款）
- 数据统计（DAU、阅读时长、订阅转化率、收入）

---

## 6. 安全设计

### 6.1 认证与授权
- JWT Access Token（有效期 15 分钟）
- Refresh Token 轮换机制（每次刷新后旧 Token 失效）
- 手机号 + 短信验证码登录（主登录方式）
- Apple / Google 第三方登录（可选）
- 设备绑定 + 异地登录提醒

### 6.2 内容保护
- 文章正文 API 仅对付费用户返回（`is_premium=true` 的文章需要有效订阅）
- API 限流：Redis Token Bucket，单用户 100 req/min
- 移动端 SSL Pinning（防止中间人抓包）
- 订阅状态严格服务端校验（不可信客户端本地状态）

### 6.3 数据安全
- 密码/敏感信息不存储明文
- 数据库连接加密（SSL/TLS）
- 定期备份策略（每日全量 + 实时 WAL 归档）

---

## 7. 错误处理与降级

### 7.1 错误码规范
```json
{
  "code": "ARTICLE_NOT_FOUND",
  "message": "文章不存在",
  "details": {
    "article_id": "xxx"
  }
}
```

### 7.2 客户端容错
- API 超时/失败时，优先使用本地缓存数据
- 阅读器内容已离线缓存，无网络也可继续阅读
- 查词服务故障时，降级到本地词典包
- 翻译服务故障时，提示"翻译暂不可用"，不影响阅读

### 7.3 监控与告警
- 结构化日志（Go: slog/zerolog）
- Prometheus 指标采集 + Grafana 可视化
- 关键链路埋点（阅读完成率、查词频率、订阅转化率）
- 异常告警（5xx 错误率 > 1%、API 响应时间 > 500ms）

---

## 8. 测试策略

### 8.1 后端测试（Go）
- **单元测试**：table-driven tests，核心逻辑覆盖率 > 60%
- **接口测试**：httptest 模拟 HTTP 请求
- **数据库测试**：testcontainers 启动真实 PostgreSQL 容器
- **Mock 外部依赖**：支付网关、推送服务、翻译 API

### 8.2 Android 测试
- **单元测试**：ViewModel、Repository 层（JUnit + MockK）
- **UI 测试**：Compose UI Test，关键页面导航
- **阅读器测试**：Espresso 模拟长按、滑动、翻页手势

### 8.3 iOS 测试
- **单元测试**：XCTest，ViewModel、Service 层
- **UI 测试**：XCUITest，关键用户路径
- **阅读器测试**：自动化手势测试（长按查词、文本选择）

### 8.4 关键路径 E2E
注册 → 浏览免费文章 → 订阅付费 → 阅读付费文章 → 查词 → 划线/笔记 → 打卡

---

## 9. 第三方服务集成

| 服务 | 用途 | 说明 |
|------|------|------|
| Apple In-App Purchase | iOS 订阅支付 | 必须接入，苹果强制要求 |
| Google Play Billing | Android 订阅支付 | 必须接入 |
| 微信支付 / 支付宝 | 国内支付渠道 | Android 端可选补充 |
| Firebase Cloud Messaging | Android 推送 | 免费，稳定 |
| APNS | iOS 推送 | 苹果官方推送 |
| DeepL / Google Translate | 翻译服务 | 整句/整段翻译，按量计费 |
| 阿里云 OSS / MinIO | 对象存储 | 音频文件、图片存储 |

---

## 10. 里程碑规划

### Phase 1: MVP（2-3 个月）
- [ ] 后端：用户系统、文章系统、阅读进度 API
- [ ] Android/iOS：基础阅读器（滚动模式、长按查词、滑动选择翻译）
- [ ] 订阅支付：Apple IAP + Google Play
- [ ] 运营 CMS：文章录入、定时发布
- [ ] 每日推送通知

### Phase 2: 体验完善（1-2 个月）
- [ ] 阅读器增强：语音朗读、笔记/划线、生词本
- [ ] 阅读统计与打卡系统
- [ ] 分类浏览与搜索
- [ ] 微信/支付宝支付
- [ ] 离线缓存与下载

### Phase 3: 增长优化（持续）
- [ ] 个性化推荐算法
- [ ] 社交功能（学习小组、排行榜）
- [ ] Elasticsearch 搜索升级
- [ ] A/B 测试与转化率优化

---

## 11. 风险与对策

| 风险 | 影响 | 对策 |
|------|------|------|
| 原生双端开发成本高 | 高 | MVP 优先核心功能，非核心功能后续迭代；考虑共用后端和设计方案 |
| 阅读器性能/体验不达预期 | 高 | 早期原型验证，阅读器作为首个技术 Spike |
| 内容版权风险 | 高 | 内容来源明确授权，或采用原创/改编策略 |
| 订阅转化率低 | 中 | 免费内容质量要有吸引力，付费墙时机需 A/B 测试 |
| 苹果/谷歌审核被拒 | 中 | 严格遵守 IAP 规范，避免绕过支付 |

---

*文档结束。本设计经双方确认后，进入 implementation plan 阶段。*
