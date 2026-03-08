import { defineStore } from 'pinia'
import { ref } from 'vue'
import request from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  // 登录
  async function login(username, password) {
    const response = await request.post('/api/auth/login', {
      username,
      password
    })
    
    if (response.code === 200) {
      token.value = response.data.token
      localStorage.setItem('token', response.data.token)
      
      // 获取用户信息（失败不影响登录）
      try {
        await getUserInfo()
      } catch (error) {
        console.error('获取用户信息失败，但不影响登录:', error)
      }
    }
    
    return response
  }

  // 获取用户信息
  async function getUserInfo() {
    const response = await request.get('/api/user/info')
    if (response.code === 200) {
      userInfo.value = response.data
    }
    return response
  }

  // 登出
  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    login,
    getUserInfo,
    logout
  }
})
