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
        meta: { title: '首页' }
      },
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/user/Index.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/role/Index.vue'),
        meta: { title: '角色管理' }
      },
      {
        path: 'menu',
        name: 'Menu',
        component: () => import('@/views/menu/Index.vue'),
        meta: { title: '菜单管理' }
      },
      {
        path: 'dict',
        name: 'Dict',
        component: () => import('@/views/dict/Index.vue'),
        meta: { title: '字典管理' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
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
