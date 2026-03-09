import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'

// Create axios instance
const request = axios.create({
  baseURL: '',
  timeout: 10000
})

// Request interceptor
request.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Authorization'] = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Response interceptor
request.interceptors.response.use(
  response => {
    const res = response.data
    
    // If response code is not 0 or 200, show error message
    if (res.code !== 0 && res.code !== 200) {
      ElMessage.error(res.message || res.msg || 'Request failed')
      
      // Token expired or invalid, redirect to login
      if (res.code === 401) {
        const userStore = useUserStore()
        userStore.logout()
        router.push('/login')
      }
      
      return Promise.reject(new Error(res.message || res.msg || 'Request failed'))
    }
    
    return res
  },
  error => {
    ElMessage.error(error.message || 'Network error')
    return Promise.reject(error)
  }
)

export default request
