import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  // 登录
  async function login(username, password) {
    const response = await axios.post('/api/auth/login', {
      username,
      password
    })
    
    if (response.data.code === 200) {
      token.value = response.data.data.token
      localStorage.setItem('token', response.data.data.token)
      await getUserInfo()
    }
    
    return response.data
  }

  // 获取用户信息
  async function getUserInfo() {
    const response = await axios.get('/api/user/info')
    if (response.data.code === 200) {
      userInfo.value = response.data.data
    }
    return response.data
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
