<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2>Statistics</h2>
      <el-button type="primary" @click="loadStats" :loading="loading">
        <el-icon><Refresh /></el-icon>
        Refresh
      </el-button>
    </div>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon user">
              <el-icon :size="40"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.userCount }}</div>
              <div class="stat-label">Total Users</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon role">
              <el-icon :size="40"><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.roleCount }}</div>
              <div class="stat-label">Total Roles</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon menu">
              <el-icon :size="40"><Menu /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.menuCount }}</div>
              <div class="stat-label">Total Menus</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon dict">
              <el-icon :size="40"><Collection /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.dictCount }}</div>
              <div class="stat-label">Total Dicts</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="welcome-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>Welcome to WHU Campus Auth</span>
        </div>
      </template>
      <div class="welcome-content">
        <h3>Wuhan University Campus Permission Management System</h3>
        <p>A permission management system based on Gin + Vue, providing user management, role management, menu management, and dictionary management.</p>
        
        <h4>Main Features:</h4>
        <ul>
          <li>✅ User Authentication (JWT)</li>
          <li>✅ Role Permission Management</li>
          <li>✅ Dynamic Menu Configuration</li>
          <li>✅ Data Dictionary Management</li>
          <li>✅ File Upload & Download</li>
        </ul>
        
        <h4>Technology Stack:</h4>
        <ul>
          <li><strong>Backend:</strong> Go + Gin + GORM</li>
          <li><strong>Frontend:</strong> Vue 3 + Element Plus</li>
          <li><strong>Database:</strong> MySQL</li>
          <li><strong>Cache:</strong> Redis</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { getUserList, getRoleList, getMenuTree, getDictList } from '@/api'

const loading = ref(false)

const stats = reactive({
  userCount: 0,
  roleCount: 0,
  menuCount: 0,
  dictCount: 0
})

// 加载统计数据
const loadStats = async () => {
  loading.value = true
  try {
    // Get total user count
    const userRes = await getUserList({ page: 1, page_size: 1 })
    if (userRes && userRes.data) {
      stats.userCount = userRes.data.total || 0
    }
    
    // Get total role count
    const roleRes = await getRoleList({ page: 1, page_size: 1 })
    if (roleRes && roleRes.data) {
      stats.roleCount = roleRes.data.total || 0
    }
    
    // Get total menu count
    const menuRes = await getMenuTree()
    if (menuRes && menuRes.data) {
      const countMenus = (menus) => {
        let count = 0
        for (const menu of menus) {
          count++
          if (menu.children && menu.children.length > 0) {
            count += countMenus(menu.children)
          }
        }
        return count
      }
      stats.menuCount = countMenus(menuRes.data)
    }
    
    // Get total dict count
    const dictRes = await getDictList({ page: 1, page_size: 1 })
    if (dictRes && dictRes.data) {
      stats.dictCount = dictRes.data.total || 0
    }
  } catch (error) {
    console.error('Failed to load statistics:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.dashboard-header h2 {
  margin: 0;
  color: #333;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-right: 20px;
}

.stat-icon.user {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.role {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.menu {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.dict {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-top: 5px;
}

.welcome-card {
  margin-top: 20px;
}

.welcome-content {
  line-height: 1.8;
}

.welcome-content h3 {
  color: #333;
  margin-bottom: 15px;
}

.welcome-content h4 {
  color: #606266;
  margin: 15px 0 10px;
}

.welcome-content ul {
  margin-left: 20px;
}

.welcome-content li {
  margin: 5px 0;
}
</style>
