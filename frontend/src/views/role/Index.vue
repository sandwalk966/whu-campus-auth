<template>
  <div class="role-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Role List</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            Add Role
          </el-button>
        </div>
      </template>
      
      <!-- Table -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="Role Name" />
        <el-table-column prop="code" label="Role Code" />
        <el-table-column prop="desc" label="Description" />
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
        <el-form-item label="Role Name" prop="name">
          <el-input v-model="formData.name" placeholder="Enter role name" />
        </el-form-item>
        <el-form-item label="Role Code" prop="code">
          <el-input v-model="formData.code" placeholder="Enter role code (e.g., admin)" />
        </el-form-item>
        <el-form-item label="Description" prop="desc">
          <el-input
            v-model="formData.desc"
            type="textarea"
            :rows="3"
            placeholder="Enter description"
          />
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
import { getRoleList, createRole, updateRole, deleteRole } from '@/api'

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('Add Role')
const formRef = ref()

const tableData = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  id: null,
  name: '',
  code: '',
  desc: '',
  status: 1
})

const formRules = {
  name: [{ required: true, message: 'Please enter role name', trigger: 'blur' }],
  code: [{ required: true, message: 'Please enter role code', trigger: 'blur' }]
}

// Load data
const loadData = async () => {
  loading.value = true
  try {
    const res = await getRoleList({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res && res.data) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('Failed to load role list:', error)
    ElMessage.error('Load failed: ' + (error.message || 'Unknown error'))
  } finally {
    loading.value = false
  }
}

// Add
const handleAdd = () => {
  dialogTitle.value = 'Add Role'
  resetForm()
  dialogVisible.value = true
}

// Edit
const handleEdit = (row) => {
  dialogTitle.value = 'Edit Role'
  Object.assign(formData, row)
  dialogVisible.value = true
}

// Delete
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('Are you sure you want to delete this role?', 'Warning', {
      confirmButtonText: 'Confirm',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })
    await deleteRole(row.id)
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
        // Update role
        await updateRole(formData)
        ElMessage.success('Updated successfully')
      } else {
        // Create role
        await createRole(formData)
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

// Pagination change
const handleSizeChange = () => loadData()
const handlePageChange = () => loadData()

// Reset form
const resetForm = () => {
  Object.assign(formData, {
    id: null,
    name: '',
    code: '',
    desc: '',
    status: 1
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.role-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
