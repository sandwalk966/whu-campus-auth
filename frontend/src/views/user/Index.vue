<template>
  <div class="user-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>User List</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            Add User
          </el-button>
        </div>
      </template>
      
      <!-- Search Bar -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="Username">
          <el-input v-model="searchForm.username" placeholder="Enter username" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            Search
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            Reset
          </el-button>
        </el-form-item>
      </el-form>
      
      <!-- Table -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="Username" />
        <el-table-column prop="email" label="Email" />
        <el-table-column prop="phone" label="Phone" />
        <el-table-column prop="status" label="Status" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? 'Active' : 'Disabled' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created At" />
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
      
      <!-- Pagination -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: flex-end;"
      />
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
        label-width="80px"
      >
        <el-form-item label="Username" prop="username">
          <el-input v-model="formData.username" placeholder="Enter username" />
        </el-form-item>
        <el-form-item label="Password" prop="password" v-if="!formData.id">
          <el-input v-model="formData.password" type="password" placeholder="Enter password" />
        </el-form-item>
        <el-form-item label="Email" prop="email">
          <el-input v-model="formData.email" placeholder="Enter email" />
        </el-form-item>
        <el-form-item label="Phone" prop="phone">
          <el-input v-model="formData.phone" placeholder="Enter phone number" />
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, createUser, updateUser, deleteUser } from '@/api'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('Add User')
const formRef = ref()

const searchForm = reactive({
  username: ''
})

const tableData = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  id: null,
  username: '',
  password: '',
  email: '',
  phone: '',
  status: 1
})

const formRules = {
  username: [
    { required: true, message: 'Please enter username', trigger: 'blur' },
    { min: 3, max: 50, message: 'Username must be between 3 and 50 characters', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enter password', trigger: 'blur' },
    { min: 6, max: 50, message: 'Password must be between 6 and 50 characters', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: 'Please enter a valid email format', trigger: 'blur' }
  ],
  phone: [
    { 
      pattern: /^1[3-9]\d{9}$/, 
      message: 'Invalid phone format, please enter 11-digit Mainland China phone number (starts with 1, second digit is 3-9)', 
      trigger: 'blur' 
    }
  ]
}

// Load data
const loadData = async () => {
  loading.value = true
  try {
    const res = await getUserList({
      page: pagination.page,
      page_size: pagination.pageSize,
      username: searchForm.username
    })
    if (res && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('Failed to load user list:', error)
    ElMessage.error('Load failed: ' + (error.message || 'Unknown error'))
  } finally {
    loading.value = false
  }
}

// Search
const handleSearch = () => {
  pagination.page = 1
  loadData()
}

// Reset
const handleReset = () => {
  searchForm.username = ''
  handleSearch()
}

// Add
const handleAdd = () => {
  dialogTitle.value = 'Add User'
  resetForm()
  dialogVisible.value = true
}

// Edit
const handleEdit = (row) => {
  dialogTitle.value = 'Edit User'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// Delete
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('Are you sure you want to delete this user?', 'Warning', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })
    
    await deleteUser(row.id)
    ElMessage.success('Deleted successfully')
    
    // If current page data is empty and not the first page, navigate to previous page
    if (tableData.value.length === 1 && pagination.page > 1) {
      pagination.page--
    }
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
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning('Please check if the form is filled correctly')
      return
    }
    
    submitLoading.value = true
    try {
      if (formData.id) {
        // Update user
        await updateUser(formData)
        ElMessage.success('Updated successfully')
      } else {
        // Create user
        await createUser(formData)
        ElMessage.success('Created successfully')
      }
      dialogVisible.value = false
      // When adding a user, if on first page, refresh current page; otherwise keep current page
      loadData()
    } catch (error) {
      console.error('Submit failed:', error)
      // Display error message from backend
      const errorMsg = error.response?.data?.message || error.message || 'Operation failed'
      ElMessage.error(errorMsg)
    } finally {
      submitLoading.value = false
    }
  })
}

// Pagination change
const handleSizeChange = () => loadData()
const handlePageChange = () => loadData()

// Reset form
const resetForm = () => {
  Object.assign(formData, {
    id: null,
    username: '',
    password: '',
    email: '',
    phone: '',
    status: 1
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.user-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}
</style>
