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
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  console.log('路由守卫 - from:', from.path, 'to:', to.path, 'hasToken:', !!userStore.token)
  
  // 如果要去登录页且已有 token，直接跳转到首页
  if (to.path === '/login' && userStore.token) {
    console.log('已有 token，从登录页跳转到 dashboard')
    next('/dashboard')
    return
  }
  
  // 如果要去非登录页但没有 token，跳转到登录页
  if (to.path !== '/login' && !userStore.token) {
    console.log('没有 token，跳转到登录页')
    next('/login')
    return
  }
  
  // 其他情况正常访问
  console.log('允许访问:', to.path)
  next()
})

export default router
