<template>
  <div class="menu-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Menu Tree</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            Add Menu
          </el-button>
        </div>
      </template>
      
      <!-- Table -->
      <el-table
        v-loading="loading"
        :data="tableData"
        row-key="id"
        border
        default-expand-all
      >
        <el-table-column prop="name" label="Menu Name" width="200" />
        <el-table-column prop="path" label="Path" />
        <el-table-column prop="component" label="Component" />
        <el-table-column prop="icon" label="Icon" width="100" />
        <el-table-column prop="sort" label="Sort" width="80" />
        <el-table-column prop="status" label="Status" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? 'Active' : 'Disabled' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleEdit(row)">
              Edit
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="Parent Menu" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="menuTreeOptions"
            :props="{ value: 'id', label: 'name', children: 'children' }"
            check-strictly
            placeholder="Select parent menu (leave empty for top-level)"
            clearable
          />
        </el-form-item>
        <el-form-item label="Menu Name" prop="name">
          <el-input v-model="formData.name" placeholder="Enter menu name" />
        </el-form-item>
        <el-form-item label="Path" prop="path">
          <el-input v-model="formData.path" placeholder="Enter path (e.g., /user)" />
        </el-form-item>
        <el-form-item label="Component" prop="component">
          <el-input v-model="formData.component" placeholder="Enter component path (e.g., user/index)" />
        </el-form-item>
        <el-form-item label="Icon" prop="icon">
          <el-input v-model="formData.icon" placeholder="Enter icon name (e.g., User)" />
        </el-form-item>
        <el-form-item label="Sort" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="Status">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">Active</el-radio>
            <el-radio :label="0">Disabled</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          Confirm
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMenuTree, createMenu, updateMenu, deleteMenu } from '@/api'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('Add Menu')
const formRef = ref()

const tableData = ref([])

const formData = reactive({
  id: null,
  parent_id: null,
  name: '',
  path: '',
  component: '',
  icon: '',
  sort: 1,
  status: 1
})

const formRules = {
  name: [{ required: true, message: 'Please enter menu name', trigger: 'blur' }],
  path: [{ required: true, message: 'Please enter path', trigger: 'blur' }]
}

// Menu tree options
const menuTreeOptions = computed(() => {
  const options = [{ id: 0, name: 'Top Level Menu', children: tableData.value }]
  return options
})

// Load data
const loadData = async () => {
  loading.value = true
  try {
    const res = await getMenuTree()
    if (res && res.data) {
      tableData.value = res.data || []
    }
  } catch (error) {
    console.error('Failed to load menu tree:', error)
    ElMessage.error('Load failed: ' + (error.message || 'Unknown error'))
  } finally {
    loading.value = false
  }
}

// Add
const handleAdd = () => {
  dialogTitle.value = 'Add Menu'
  resetForm()
  dialogVisible.value = true
}

// Edit
const handleEdit = (row) => {
  dialogTitle.value = 'Edit Menu'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// Delete
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('Are you sure you want to delete this menu?', 'Warning', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })
    await deleteMenu(row.id)
    ElMessage.success('Deleted successfully')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete failed:', error)
      ElMessage.error('Delete failed: ' + (error.message || 'Unknown error'))
    }
  }
}

// Submit
const handleSubmit = async () => {
  formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitLoading.value = true
    try {
      if (formData.id) {
        // Update menu
        await updateMenu(formData)
        ElMessage.success('Updated successfully')
      } else {
        // Create menu
        await createMenu(formData)
        ElMessage.success('Created successfully')
      }
      dialogVisible.value = false
      loadData()
    } catch (error) {
      console.error('Submit failed:', error)
      const errorMsg = error.response?.data?.message || error.message || 'Operation failed'
      ElMessage.error(errorMsg)
    } finally {
      submitLoading.value = false
    }
  })
}

// Reset form
const resetForm = () => {
  Object.assign(formData, {
    id: null,
    parent_id: null,
    name: '',
    path: '',
    component: '',
    icon: '',
    sort: 1,
    status: 1
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  loadData()
})
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
