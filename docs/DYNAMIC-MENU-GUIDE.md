# 动态菜单功能说明

## 功能概述

系统现在支持根据用户角色动态生成侧边栏菜单。不同角色的用户登录后，只能看到其角色权限范围内的菜单。

## 实现原理

### 1. 后端逻辑

#### 数据模型
- **用户表** (`sys_user`)：存储用户信息
- **角色表** (`sys_role`)：存储角色信息
- **菜单表** (`sys_menu`)：存储菜单信息
- **角色 - 菜单关联表** (`role_menus`)：存储角色与菜单的多对多关系

#### 关联关系
```
用户 (User) → 角色 (Role) → 菜单 (Menu)
```

#### 关键代码

**`dao/user_dao.go`** - 预加载角色和菜单：
```go
func (dao *UserDAO) PreloadRoles(user *db.User) error {
    return dao.db.Preload("Roles.Menus").First(user, user.ID).Error
}
```

**`api/user.go`** - 获取用户信息（包含角色和菜单）：
```go
func (api *UserAPI) GetUserInfo(c *gin.Context) {
    user, err := api.userService.GetUserByID(userId)
    // 返回 user，其中包含 roles 和 roles 的 menus
    utils.SuccessWithData(c, user)
}
```

### 2. 前端逻辑

#### Store (`stores/user.js`)
```javascript
// 存储用户菜单（树形结构）
const userMenus = ref([])

// 获取用户信息时自动构建菜单树
async function getUserInfo() {
    const response = await request.get('/api/user/info')
    if (response.code === 200) {
        userInfo.value = response.data
        // 从用户信息中提取菜单
        if (response.data.roles && response.data.roles.length > 0) {
            const allMenus = []
            response.data.roles.forEach(role => {
                if (role.menus && role.menus.length > 0) {
                    allMenus.push(...role.menus)
                }
            })
            // 去重并构建树形结构
            const uniqueMenus = Array.from(new Map(allMenus.map(menu => [menu.id, menu])).values())
            userMenus.value = buildMenuTree(uniqueMenus)
        }
    }
}

// 构建树形结构
function buildMenuTree(menus) {
    const menuMap = new Map()
    const tree = []
    
    menus.forEach(menu => {
        menuMap.set(menu.id, { ...menu, children: [] })
    })
    
    menus.forEach(menu => {
        const node = menuMap.get(menu.id)
        if (menu.parent_id === 0 || menu.parent_id === null) {
            tree.push(node)
        } else {
            const parent = menuMap.get(menu.parent_id)
            if (parent) {
                parent.children.push(node)
            }
        }
    })
    
    return tree
}
```

#### 布局组件 (`layouts/MainLayout.vue`)
```vue
<el-menu>
    <!-- 首页（始终显示） -->
    <el-menu-item index="/dashboard">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
    </el-menu-item>
    
    <!-- 动态菜单 -->
    <template v-for="menu in userStore.userMenus" :key="menu.id">
        <!-- 有子菜单 -->
        <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="String(menu.id)">
            <template #title>
                <el-icon v-if="menu.icon"><component :is="menu.icon" /></el-icon>
                <span>{{ menu.name }}</span>
            </template>
            <el-menu-item v-for="child in menu.children" :key="child.id" :index="child.path">
                <el-icon v-if="child.icon"><component :is="child.icon" /></el-icon>
                <span>{{ child.name }}</span>
            </el-menu-item>
        </el-sub-menu>
        
        <!-- 无子菜单 -->
        <el-menu-item v-else :index="menu.path">
            <el-icon v-if="menu.icon"><component :is="menu.icon" /></el-icon>
            <span>{{ menu.name }}</span>
        </el-menu-item>
    </template>
</el-menu>
```

## 菜单数据结构

### 菜单表字段
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 菜单 ID |
| name | string | 菜单名称 |
| path | string | 路由路径 |
| component | string | 组件路径 |
| icon | string | 图标名称（Element Plus 图标） |
| sort | int | 排序值 |
| parent_id | uint | 父菜单 ID（0 表示一级菜单） |
| type | int | 菜单类型（1 目录/2 菜单/3 按钮） |
| status | int | 状态（1 正常/0 禁用） |

### 示例数据
```json
{
    "id": 1,
    "name": "系统管理",
    "path": "/system",
    "component": "layout",
    "icon": "Setting",
    "sort": 1,
    "parent_id": 0,
    "type": 1,
    "status": 1,
    "children": [
        {
            "id": 2,
            "name": "用户管理",
            "path": "/user",
            "component": "user/index",
            "icon": "User",
            "sort": 1,
            "parent_id": 1,
            "type": 1,
            "status": 1
        },
        {
            "id": 3,
            "name": "角色管理",
            "path": "/role",
            "component": "role/index",
            "icon": "UserFilled",
            "sort": 2,
            "parent_id": 1,
            "type": 1,
            "status": 1
        }
    ]
}
```

## 初始化菜单

### 方法 1：自动初始化（推荐）

系统启动时自动执行 `initializer.InitMenus(db)`，创建默认菜单：
- 系统管理（一级菜单）
  - 用户管理
  - 角色管理
  - 菜单管理
  - 字典管理

并为管理员角色分配所有菜单。

### 方法 2：手动执行 SQL

```bash
# 进入 MySQL 容器
docker compose exec mysql mysql -u root -proot whu_campus_auth

# 执行初始化脚本
source /scripts/init_menus.sql
```

或直接执行：
```bash
docker compose exec mysql mysql -u root -proot whu_campus_auth < scripts/init_menus.sql
```

## 为角色分配菜单

### 方法 1：通过菜单管理页面

1. 登录管理员账号
2. 进入"菜单管理"
3. 创建或编辑菜单
4. 在角色管理中为角色分配菜单

### 方法 2：通过 API

```bash
# 为角色分配菜单
POST /api/role/assign-menus
{
    "role_id": 1,
    "menu_ids": [1, 2, 3, 4]
}
```

### 方法 3：直接操作数据库

```sql
-- 为角色 ID=1 分配所有菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM sys_menu;

-- 查看角色的菜单
SELECT r.name AS role_name, m.name AS menu_name
FROM sys_role r
JOIN role_menus rm ON r.id = rm.role_id
JOIN sys_menu m ON rm.menu_id = m.id
WHERE r.id = 1;
```

## 测试步骤

### 1. 初始化菜单数据
```bash
# 重启应用（会自动初始化菜单）
docker compose restart app

# 查看日志确认初始化成功
docker compose logs app | grep "菜单初始化"
```

### 2. 验证用户信息
```bash
# 登录获取 token
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取用户信息（应包含 roles 和 menus）
curl -X GET http://localhost/api/user/info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 前端验证
1. 打开浏览器开发者工具
2. 登录后查看 Network 标签
3. 检查 `/api/user/info` 响应
4. 确认返回数据包含 `roles[].menus`
5. 查看侧边栏是否显示动态菜单

## 常见问题

### Q1: 登录后看不到菜单？
**A**: 检查以下几点：
1. 数据库中是否有菜单数据：`SELECT * FROM sys_menu;`
2. 用户角色是否有菜单权限：`SELECT * FROM role_menus WHERE role_id = ?;`
3. 前端是否正确解析菜单树：查看浏览器控制台是否有错误

### Q2: 菜单图标不显示？
**A**: 确保 `icon` 字段值是 Element Plus 支持的图标名称，如：
- `User`
- `UserFilled`
- `Menu`
- `Setting`
- `Collection`
- `HomeFilled`

完整图标列表：https://element-plus.org/zh-CN/component/icon.html

### Q3: 如何添加新的菜单项？
**A**: 有三种方式：
1. **前端界面**：进入"菜单管理"页面添加
2. **API 调用**：`POST /api/menu`
3. **直接插入数据库**：
   ```sql
   INSERT INTO sys_menu (name, path, component, icon, sort, parent_id, type, status)
   VALUES ('新功能', '/new-feature', 'new-feature/index', 'Star', 5, 1, 1, 1);
   ```

### Q4: 如何让不同角色看到不同的菜单？
**A**: 为不同角色分配不同的菜单：
```sql
-- 为普通用户角色只分配用户管理菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM sys_role r, sys_menu m
WHERE r.code = 'user' AND m.path = '/user';
```

## 权限控制

### 后端权限验证
```go
// 创建、更新、删除菜单需要管理员权限
menu.POST("", middleware.IsAdmin(), menuAPI.CreateMenu)
menu.PUT("", middleware.IsAdmin(), menuAPI.UpdateMenu)
menu.DELETE("/:id", middleware.IsAdmin(), menuAPI.DeleteMenu)

// 查询菜单只需要已认证用户
menu.GET("/tree", menuAPI.GetMenuTree)
menu.GET("/list", menuAPI.GetMenuList)
```

### 前端路由守卫
```javascript
// router/index.js
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    
    if (!userStore.token) {
        next('/login')
        return
    }
    
    // 可以在这里检查用户是否有权限访问该路由
    next()
})
```

## 总结

✅ **已实现功能**：
- 根据用户角色动态生成菜单
- 支持多级菜单（树形结构）
- 菜单图标自定义
- 菜单排序
- 管理员可动态管理菜单
- 不同角色看到不同菜单

✅ **优势**：
- 灵活的权限控制
- 无需修改代码即可调整菜单
- 支持动态添加新功能
- 提高系统安全性

✅ **后续可扩展**：
- 菜单权限按钮级别控制
- 菜单国际化
- 菜单个性化定制（用户自定义排序）
- 菜单访问日志
