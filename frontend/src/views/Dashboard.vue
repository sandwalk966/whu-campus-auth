<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon user">
              <el-icon :size="40"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">128</div>
              <div class="stat-label">用户总数</div>
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
              <div class="stat-value">5</div>
              <div class="stat-label">角色数量</div>
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
              <div class="stat-value">20</div>
              <div class="stat-label">菜单数量</div>
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
              <div class="stat-value">10</div>
              <div class="stat-label">字典数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="welcome-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>欢迎使用 WHU Campus Auth</span>
        </div>
      </template>
      <div class="welcome-content">
        <h3>武汉大学校园权限管理系统</h3>
        <p>这是一个基于 Gin + Vue 的权限管理系统，提供用户管理、角色管理、菜单管理和字典管理等功能。</p>
        
        <h4>主要功能：</h4>
        <ul>
          <li>✅ 用户认证（JWT）</li>
          <li>✅ 角色权限管理</li>
          <li>✅ 动态菜单配置</li>
          <li>✅ 数据字典管理</li>
          <li>✅ 文件上传下载</li>
        </ul>
        
        <h4>技术栈：</h4>
        <ul>
          <li><strong>后端：</strong>Go + Gin + GORM</li>
          <li><strong>前端：</strong>Vue 3 + Element Plus</li>
          <li><strong>数据库：</strong>MySQL</li>
          <li><strong>缓存：</strong>Redis</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<script setup>
</script>

<style scoped>
.dashboard {
  padding: 20px;
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
