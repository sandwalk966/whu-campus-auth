import request from '@/utils/request'

// 登录
export function login(username, password) {
  return request.post('/api/auth/login', {
    username,
    password
  })
}

// 获取用户信息
export function getUserInfo() {
  return request.get('/api/user/info')
}

// 获取用户列表
export function getUserList(params) {
  return request.get('/api/user/list', { params })
}

// 删除用户
export function deleteUser(id) {
  return request.delete(`/api/user/${id}`)
}

// 分配角色
export function assignRoles(userId, roleIds) {
  return request.post('/api/user/assign-roles', {
    user_id: userId,
    role_ids: roleIds
  })
}

// 获取所有角色
export function getAllRoles() {
  return request.get('/api/role/all')
}

// 获取角色列表
export function getRoleList(params) {
  return request.get('/api/role/list', { params })
}

// 获取角色详情
export function getRoleById(id) {
  return request.get(`/api/role/${id}`)
}

// 创建角色
export function createRole(data) {
  return request.post('/api/role', data)
}

// 更新角色
export function updateRole(data) {
  return request.put('/api/role', data)
}

// 删除角色
export function deleteRole(id) {
  return request.delete(`/api/role/${id}`)
}

// 获取菜单树
export function getMenuTree() {
  return request.get('/api/menu/tree')
}

// 获取菜单列表
export function getMenuList(params) {
  return request.get('/api/menu/list', { params })
}

// 获取菜单详情
export function getMenuById(id) {
  return request.get(`/api/menu/${id}`)
}

// 创建菜单
export function createMenu(data) {
  return request.post('/api/menu', data)
}

// 更新菜单
export function updateMenu(data) {
  return request.put('/api/menu', data)
}

// 删除菜单
export function deleteMenu(id) {
  return request.delete(`/api/menu/${id}`)
}

// 获取字典列表
export function getDictList(params) {
  return request.get('/api/dict/list', { params })
}

// 获取字典详情
export function getDictById(id) {
  return request.get(`/api/dict/${id}`)
}

// 根据编码获取字典
export function getDictByCode(code) {
  return request.get(`/api/dict/code/${code}`)
}

// 创建字典
export function createDict(data) {
  return request.post('/api/dict', data)
}

// 更新字典
export function updateDict(data) {
  return request.put('/api/dict', data)
}

// 删除字典
export function deleteDict(id) {
  return request.delete(`/api/dict/${id}`)
}

// 上传文件
export function uploadFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/api/upload', formData)
}
