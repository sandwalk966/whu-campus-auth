import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: 'Dashboard' }
      },
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/user/Index.vue'),
        meta: { title: 'User Management' }
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/role/Index.vue'),
        meta: { title: 'Role Management' }
      },
      {
        path: 'menu',
        name: 'Menu',
        component: () => import('@/views/menu/Index.vue'),
        meta: { title: 'Menu Management' }
      },
      {
        path: 'dict',
        name: 'Dict',
        component: () => import('@/views/dict/Index.vue'),
        meta: { title: 'Dictionary Management' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  
  // 如果有 token 但没有用户信息，先获取用户信息（页面刷新场景）
  if (userStore.token && !userStore.userInfo) {
    try {
      await userStore.initUserInfo()
    } catch (error) {
      console.error('路由守卫初始化用户信息失败:', error)
      // 只在确认 token 过期（401）时才清除并跳转
      if (error.response && error.response.status === 401) {
        console.log('Token 已过期，清除登录状态')
        userStore.logout()
        next('/login')
        return
      }
      // 其他错误（如 500）保留 token，可能是后端临时问题，继续访问
      console.log('非 401 错误，保留 token 继续访问')
    }
  }
  
  // 如果要去登录页且已有 token，直接跳转到首页
  if (to.path === '/login' && userStore.token) {
    next('/dashboard')
    return
  }
  
  // 如果要去非登录页但没有 token，跳转到登录页
  if (to.path !== '/login' && !userStore.token) {
    next('/login')
    return
  }
  
  // 其他情况正常访问
  next()
})

export default router
