-- 插入默认菜单数据
-- 注意：运行此脚本前请确保数据库已初始化

USE whu_campus_auth;

-- 1. 系统管理（一级菜单）
INSERT INTO sys_menu (name, path, component, icon, sort, parent_id, type, status, created_at, updated_at)
VALUES ('系统管理', '/system', 'layout', 'Setting', 1, 0, 1, 1, NOW(), NOW());

-- 获取刚插入的菜单 ID
SET @system_id = LAST_INSERT_ID();

-- 2. 用户管理（二级菜单）
INSERT INTO sys_menu (name, path, component, icon, sort, parent_id, type, status, created_at, updated_at)
VALUES ('用户管理', '/user', 'user/index', 'User', 1, @system_id, 1, 1, NOW(), NOW());

-- 3. 角色管理（二级菜单）
INSERT INTO sys_menu (name, path, component, icon, sort, parent_id, type, status, created_at, updated_at)
VALUES ('角色管理', '/role', 'role/index', 'UserFilled', 2, @system_id, 1, 1, NOW(), NOW());

-- 4. 菜单管理（二级菜单）
INSERT INTO sys_menu (name, path, component, icon, sort, parent_id, type, status, created_at, updated_at)
VALUES ('菜单管理', '/menu', 'menu/index', 'Menu', 3, @system_id, 1, 1, NOW(), NOW());

-- 5. 字典管理（二级菜单）
INSERT INTO sys_menu (name, path, component, icon, sort, parent_id, type, status, created_at, updated_at)
VALUES ('字典管理', '/dict', 'dict/index', 'Collection', 4, @system_id, 1, 1, NOW(), NOW());

-- 6. 为管理员角色分配所有菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM sys_role r, sys_menu m
WHERE r.code = 'admin';

SELECT '菜单数据初始化完成！' AS message;
