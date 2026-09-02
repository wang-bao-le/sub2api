---
method: dual-agent
assessmentA: 01a05ffa-5e25-7031-8433-a529623aa436
assessmentB: 01a05ffa-5e90-7831-8927-44c82f85d967
applicableMax: 40
naHeuristics: 
detector: 0 findings
browser: "failed: #app empty and Vue runtime error"
overlay: skipped
liveServer: not started
tempCleanup: pending
timestamp: 2026-09-02T02-41-25Z
slug: frontend-src-views-user-usageview-vue
---
## Design Health Score

| # | Heuristic | Score | Key Issue |
|---|---|---:|---|
| 1 | Visibility of System Status | 2 | 多组请求同时刷新，但缺少统一的“筛选更新中”反馈。 |
| 2 | Match System / Real World | 3 | Token、成本、延迟和模型映射符合开发者排障语境。 |
| 3 | User Control and Freedom | 2 | 有刷新、重置、列设置和分页，但筛选每次变更立即请求。 |
| 4 | Consistency and Standards | 2 | 标签缺少 tab 语义，列设置缺少展开状态语义。 |
| 5 | Error Prevention | 2 | 可输入任意模型片段，但无结果原因不够明确。 |
| 6 | Recognition Rather Than Recall | 2 | 无已应用筛选摘要，列状态主要靠菜单内勾选识别。 |
| 7 | Flexibility and Efficiency | 2 | 支持服务端排序、导出和页大小，但未启用快速跳页。 |
| 8 | Aesthetic and Minimalist Design | 2 | 六个筛选器、多个操作和高密度表格同时出现。 |
| 9 | Error Recovery | 2 | 失败只 toast，空表可能被误解为确实无记录。 |
| 10 | Help and Documentation | 2 | Token/成本说明依赖鼠标悬停，触屏和键盘不可用。 |
| **Total** |  | **21/40** | **中等，优先优化信息层级和可访问交互。** |

## Design Specificity Verdict

页面有明确的开发者控制台特征：模型映射链、Token/缓存/图片计费、成本倍率和延迟信息都服务于请求审计与排障；但下方交互组织仍接近通用后台表格，专业数据没有被组织成清晰的“筛选—定位—查看详情”路径。

## 主要优点

- 使用记录覆盖成本核算、性能排查和模型映射审计，不是单纯 CRUD 表格。
- 桌面端有 sticky 表头、服务端排序和虚拟滚动。
- 移动端切换为字段卡片，具备独立的加载和空状态。

## 优先优化意见

### P1：先显示标签，再显示对应筛选器

当前标签在筛选卡之后，但筛选器内容由 `activeTab` 决定，用户不容易立即知道条件作用于“使用记录”还是“错误请求”。建议将标签移到筛选卡顶部，标签下显示对应筛选与操作区；移动端使用整行 segmented control。

证据：[UsageView.vue:69-171](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/views/user/UsageView.vue:69>)。

### P1：筛选改为“暂存后应用”

当前每个 Select 都直接触发查询，并同时刷新日志、统计和图表，连续选择多个条件会造成反复等待和结果不确定。建议增加“应用筛选”与“清除全部”，并显示已应用条件摘要；保留“刷新”作为仅重新请求当前条件。

证据：[UsageView.vue:98-160](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/views/user/UsageView.vue:98>)、[UsageView.vue:529-536](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/views/user/UsageView.vue:529>)。

### P1：详情入口不能只依赖悬停

Token 与成本详情使用 `mouseenter`/`mouseleave`，触屏用户无法稳定打开，键盘用户也无法访问。建议改为可聚焦按钮，支持 Enter/Space 打开、Escape 关闭，并补充 `aria-label`。

证据：[UsageTable.vue:178-209](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/components/admin/usage/UsageTable.vue:178>)。

### P2：移动端默认只显示核心字段

当前移动端会把所有可见列逐项转成卡片字段，端点、IP、计费和延迟等会使卡片过长。建议默认保留时间、模型、成本、Token、请求类型；其他字段放进“查看详情”展开区。

证据：[DataTable.vue:70-84](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/components/common/DataTable.vue:70>)、[UsageView.vue:700-714](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/views/user/UsageView.vue:700>)。

### P2：区分空结果、首次加载和请求失败

当前表格统一显示 `usage.noRecords`，失败时主要依赖全局 toast，用户可能把网络错误当成无数据。建议分别提供“暂无记录”“当前条件无结果”和“加载失败/重试”，无结果状态提供重置筛选。

证据：[UsageTable.vue:272-273](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/components/admin/usage/UsageTable.vue:272>)、[UsageView.vue:445-448](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/views/user/UsageView.vue:445>)。

### P2：补齐键盘和移动端操作语义

可排序表头目前通过 `th @click` 触发，不是可聚焦控件；列设置菜单缺少 `aria-expanded`、Escape 关闭和焦点回收。建议使用表头按钮并完善菜单焦点管理。分页移动端可补充结果范围或页大小提示。

证据：[DataTable.vue:120-164](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/components/common/DataTable.vue:120>)、[UsageView.vue:131-155](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/views/user/UsageView.vue:131>)、[Pagination.vue:5-24](</Users/wangbaole/Documents/Codex/2026-08-24/s-s/sub2api/frontend/src/components/common/Pagination.vue:5>)。

## 认知负荷与角色风险

当前为中高负荷：6 个筛选条件与多个操作同处一行，表格单元格还叠加 Token、缓存、图片、成本和延迟层级。普通用户容易误判成本；开发者需要记忆颜色和缩写；移动端用户难以快速找到时间与成本；键盘/触屏用户无法获得悬停详情。

## 证据边界

源码审查未修改文件，也未查看 `Select`、`EmptyState`、`UserErrorRequestsTable` 内部实现；颜色对比度、实际断点和 Tooltip 溢出仍需真实浏览器与键盘/触屏复测。自动检测结果为 0 项，但浏览器访问 `/usage` 时 `#app` 为空并出现 Vue 运行时错误，因此视觉验证未通过，不能据此确认页面运行正常。
