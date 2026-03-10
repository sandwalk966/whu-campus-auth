import { defineStore } from 'pinia'
import { ref } from 'vue'
import request from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)
  const userMenus = ref([])

  // Login
  async function login(username, password) {
    const response = await request.post('/api/auth/login', {
      username,
      password
    })
    
    if (response.code === 0 || response.code === 200) {
      token.value = response.data.token
      localStorage.setItem('token', response.data.token)
      
      // Get user information (failure doesn't affect login)
      try {
        await getUserInfo()
      } catch (error) {
        console.error('Failed to get user info:', error)
      }
    }
    
    return response
  }

  // Get user information
  async function getUserInfo() {
    const response = await request.get('/api/user/info')
    console.log('=== getUserInfo response ===', response)
    if (response.code === 0 || response.code === 200) {
      userInfo.value = response.data
      console.log('userInfo set to:', response.data)
      // Extract menus from user information
      if (response.data.roles && response.data.roles.length > 0) {
        console.log('Roles found:', response.data.roles.length)
        // Merge menus from all roles
        const allMenus = []
        response.data.roles.forEach(role => {
          console.log('Role:', role.name, 'menus count:', role.menus ? role.menus.length : 0)
          if (role.menus && role.menus.length > 0) {
            allMenus.push(...role.menus)
          }
        })
        console.log('Total menus before dedup:', allMenus.length)
        // Remove duplicates (based on menu ID)
        const uniqueMenus = Array.from(new Map(allMenus.map(menu => [menu.id, menu])).values())
        console.log('Unique menus:', uniqueMenus.length)
        userMenus.value = buildMenuTree(uniqueMenus)
        console.log('Menu tree:', userMenus.value)
      } else {
        console.log('No roles found or roles is empty')
      }
    }
    return response
  }

  // Build tree menu structure
  function buildMenuTree(menus) {
    const menuMap = new Map()
    const tree = []

    // Put all menus into map
    menus.forEach(menu => {
      menuMap.set(menu.id, { ...menu, children: [] })
    })

    // Build tree structure
    menus.forEach(menu => {
      const node = menuMap.get(menu.id)
      if (menu.parent_id === 0 || menu.parent_id === null) {
        // Top-level menu
        tree.push(node)
      } else {
        // Sub-menu
        const parent = menuMap.get(menu.parent_id)
        if (parent) {
          parent.children.push(node)
        }
      }
    })

    return tree
  }

  // Set menus (for manual refresh)
  function setMenus(menus) {
    userMenus.value = menus
  }

  // Logout
  function logout() {
    token.value = ''
    userInfo.value = null
    userMenus.value = []
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    userMenus,
    login,
    getUserInfo,
    setMenus,
    logout
  }
})
