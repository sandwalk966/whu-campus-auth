<template>
  <div class="menu-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>菜单树</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增菜单
          </el-button>
        </div>
      </template>
      
      <el-table :data="tableData" row-key="id" border default-expand-all>
        <el-table-column prop="name" label="菜单名称" width="200" />
        <el-table-column prop="path" label="路径" />
        <el-table-column prop="component" label="组件" />
        <el-table-column prop="icon" label="图标" width="100" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const tableData = ref([
  {
    id: 1,
    name: '系统管理',
    path: '/system',
    component: 'layout',
    icon: 'Setting',
    sort: 1,
    children: [
      {
        id: 2,
        name: '用户管理',
        path: '/user',
        component: 'user/index',
        icon: 'User',
        sort: 1
      },
      {
        id: 3,
        name: '角色管理',
        path: '/role',
        component: 'role/index',
        icon: 'UserFilled',
        sort: 2
      }
    ]
  }
])

const handleAdd = () => {
  console.log('添加菜单')
}

const handleEdit = (row) => {
  console.log('编辑菜单', row)
}

const handleDelete = (row) => {
  console.log('删除菜单', row)
}
</script>

<style scoped>
.menu-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
