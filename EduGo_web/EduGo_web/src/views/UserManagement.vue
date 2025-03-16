<template>
  <MainLayout>
    <div class="user-management-container">
      <h1>用户管理</h1>
      
      <a-tabs default-active-key="1">
        <!-- 用户列表 -->
        <a-tab-pane key="1" title="用户列表">
          <div class="filter-container">
            <a-select
              v-model="selectedRole"
              placeholder="按角色筛选"
              style="width: 200px; margin-right: 16px;"
              @change="handleRoleChange"
            >
              <a-option value="">全部用户</a-option>
              <a-option 
                v-for="(name, role) in USER_ROLE_NAMES" 
                :key="role" 
                :value="role"
              >
                {{ name }}
              </a-option>
            </a-select>
            
            <a-input-search
              v-model="searchKeyword"
              placeholder="搜索用户名或邮箱"
              style="width: 300px;"
              @search="handleSearch"
            />
          </div>
          
          <a-table
            :data="filteredUsers"
            :loading="loading"
            :pagination="{
              total: filteredUsers.length,
              showTotal: true,
              showPageSize: true
            }"
            stripe
          >
            <template #columns>
              <a-table-column title="ID" data-index="id" />
              <a-table-column title="用户名" data-index="username" />
              <a-table-column title="邮箱" data-index="email" />
              <a-table-column title="姓名">
                <template #cell="{ record }">
                  {{ record.firstName }} {{ record.lastName }}
                </template>
              </a-table-column>
              <a-table-column title="角色">
                <template #cell="{ record }">
                  <a-tag :color="getRoleColor(record.role)">
                    {{ USER_ROLE_NAMES[record.role] || record.role }}
                  </a-tag>
                </template>
              </a-table-column>
              <a-table-column title="状态">
                <template #cell="{ record }">
                  <a-badge
                    :status="getStatusType(record.status)"
                    :text="record.status === 'active' ? '正常' : (record.status === 'inactive' ? '未激活' : '已禁用')"
                  />
                </template>
              </a-table-column>
              <a-table-column title="注册时间">
                <template #cell="{ record }">
                  {{ formatDate(record.createdAt) }}
                </template>
              </a-table-column>
              <a-table-column title="操作" fixed="right" align="center">
                <template #cell="{ record }">
                  <a-space>
                    <a-button
                      type="text"
                      size="small"
                      @click="handleViewUser(record)"
                    >
                      查看
                    </a-button>
                    <a-button
                      v-if="canEditRole"
                      type="text"
                      size="small"
                      status="warning"
                      @click="handleEditRole(record)"
                    >
                      修改角色
                    </a-button>
                    <a-button
                      v-if="canEditStatus && record.role !== 'super_admin'"
                      type="text"
                      size="small"
                      :status="record.status === 'active' ? 'danger' : 'success'"
                      @click="handleToggleStatus(record)"
                    >
                      {{ record.status === 'active' ? '禁用' : '启用' }}
                    </a-button>
                  </a-space>
                </template>
              </a-table-column>
            </template>
          </a-table>
        </a-tab-pane>
        
        <!-- 用户关系 -->
        <a-tab-pane key="2" title="用户关系">
          <div class="relation-container">
            <a-card v-if="currentUserRole === 'super_admin' || currentUserRole === 'admin'">
              <template #title>管理的教师</template>
              <template #extra>
                <a-button type="primary" size="small" @click="showAddTeacherModal = true">
                  添加教师
                </a-button>
              </template>
              
              <a-table
                :data="teachersList"
                :loading="teachersLoading"
                :pagination="{ pageSize: 5 }"
              >
                <template #columns>
                  <a-table-column title="ID" data-index="id" />
                  <a-table-column title="用户名" data-index="username" />
                  <a-table-column title="姓名">
                    <template #cell="{ record }">
                      {{ record.firstName }} {{ record.lastName }}
                    </template>
                  </a-table-column>
                  <a-table-column title="邮箱" data-index="email" />
                  <a-table-column title="状态">
                    <template #cell="{ record }">
                      <a-badge
                        :status="getStatusType(record.status)"
                        :text="record.status === 'active' ? '正常' : (record.status === 'inactive' ? '未激活' : '已禁用')"
                      />
                    </template>
                  </a-table-column>
                </template>
              </a-table>
            </a-card>
            
            <a-card v-if="currentUserRole === 'super_admin' || currentUserRole === 'admin' || currentUserRole === 'teacher'">
              <template #title>管理的学生</template>
              <template #extra>
                <a-button type="primary" size="small" @click="showAddStudentModal = true">
                  添加学生
                </a-button>
              </template>
              
              <a-table
                :data="studentsList"
                :loading="studentsLoading"
                :pagination="{ pageSize: 5 }"
              >
                <template #columns>
                  <a-table-column title="ID" data-index="id" />
                  <a-table-column title="用户名" data-index="username" />
                  <a-table-column title="姓名">
                    <template #cell="{ record }">
                      {{ record.firstName }} {{ record.lastName }}
                    </template>
                  </a-table-column>
                  <a-table-column title="邮箱" data-index="email" />
                  <a-table-column title="状态">
                    <template #cell="{ record }">
                      <a-badge
                        :status="getStatusType(record.status)"
                        :text="record.status === 'active' ? '正常' : (record.status === 'inactive' ? '未激活' : '已禁用')"
                      />
                    </template>
                  </a-table-column>
                </template>
              </a-table>
            </a-card>
            
            <a-card v-if="currentUserRole === 'student'">
              <template #title>我的家长</template>
              <template #extra>
                <a-button type="primary" size="small" @click="showAddParentModal = true">
                  添加家长
                </a-button>
              </template>
              
              <a-table
                :data="parentsList"
                :loading="parentsLoading"
                :pagination="{ pageSize: 5 }"
              >
                <template #columns>
                  <a-table-column title="ID" data-index="id" />
                  <a-table-column title="用户名" data-index="username" />
                  <a-table-column title="姓名">
                    <template #cell="{ record }">
                      {{ record.firstName }} {{ record.lastName }}
                    </template>
                  </a-table-column>
                  <a-table-column title="邮箱" data-index="email" />
                  <a-table-column title="状态">
                    <template #cell="{ record }">
                      <a-badge
                        :status="getStatusType(record.status)"
                        :text="record.status === 'active' ? '正常' : (record.status === 'inactive' ? '未激活' : '已禁用')"
                      />
                    </template>
                  </a-table-column>
                </template>
              </a-table>
            </a-card>
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>
    
    <!-- 查看用户详情对话框 -->
    <a-modal
      v-model:visible="showUserDetailModal"
      title="用户详情"
      @cancel="showUserDetailModal = false"
      :footer="false"
    >
      <a-descriptions
        v-if="selectedUser"
        :data="userDetailData"
        :column="1"
        :label-style="{ 'font-weight': 'bold' }"
      />
    </a-modal>
    
    <!-- 修改角色对话框 -->
    <a-modal
      v-model:visible="showEditRoleModal"
      title="修改用户角色"
      @ok="handleRoleSubmit"
      @cancel="showEditRoleModal = false"
      :ok-button-props="{ loading: submitting }"
    >
      <a-form :model="editRoleForm" layout="vertical">
        <a-form-item field="role" label="角色">
          <a-select v-model="editRoleForm.role">
            <a-option 
              v-for="(name, role) in editableRoles" 
              :key="role" 
              :value="role"
            >
              {{ name }}
            </a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
    
    <!-- 添加教师对话框 -->
    <a-modal
      v-model:visible="showAddTeacherModal"
      title="添加教师"
      @ok="handleAddTeacher"
      @cancel="showAddTeacherModal = false"
      :ok-button-props="{ loading: submitting }"
    >
      <a-form :model="addTeacherForm" layout="vertical">
        <a-form-item field="teacherId" label="选择教师">
          <a-select
            v-model="addTeacherForm.teacherId"
            placeholder="请选择教师"
            :loading="teachersSelectLoading"
            :filter-option="true"
          >
            <a-option
              v-for="teacher in availableTeachers"
              :key="teacher.id"
              :value="teacher.id"
            >
              {{ teacher.username }} ({{ teacher.firstName }} {{ teacher.lastName }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="department" label="部门">
          <a-input v-model="addTeacherForm.department" placeholder="请输入部门" />
        </a-form-item>
        <a-form-item field="position" label="职位">
          <a-input v-model="addTeacherForm.position" placeholder="请输入职位" />
        </a-form-item>
      </a-form>
    </a-modal>
    
    <!-- 添加学生对话框 -->
    <a-modal
      v-model:visible="showAddStudentModal"
      title="添加学生"
      @ok="handleAddStudent"
      @cancel="showAddStudentModal = false"
      :ok-button-props="{ loading: submitting }"
    >
      <a-form :model="addStudentForm" layout="vertical">
        <a-form-item field="studentId" label="选择学生">
          <a-select
            v-model="addStudentForm.studentId"
            placeholder="请选择学生"
            :loading="studentsSelectLoading"
            :filter-option="true"
          >
            <a-option
              v-for="student in availableStudents"
              :key="student.id"
              :value="student.id"
            >
              {{ student.username }} ({{ student.firstName }} {{ student.lastName }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="courseName" label="课程名称">
          <a-input v-model="addStudentForm.courseName" placeholder="请输入课程名称" />
        </a-form-item>
        <a-form-item field="semester" label="学期">
          <a-input v-model="addStudentForm.semester" placeholder="请输入学期" />
        </a-form-item>
      </a-form>
    </a-modal>
    
    <!-- 添加家长对话框 -->
    <a-modal
      v-model:visible="showAddParentModal"
      title="添加家长"
      @ok="handleAddParent"
      @cancel="showAddParentModal = false"
      :ok-button-props="{ loading: submitting }"
    >
      <a-form :model="addParentForm" layout="vertical">
        <a-form-item field="parentId" label="选择家长">
          <a-select
            v-model="addParentForm.parentId"
            placeholder="请选择家长"
            :loading="parentsSelectLoading"
            :filter-option="true"
          >
            <a-option
              v-for="parent in availableParents"
              :key="parent.id"
              :value="parent.id"
            >
              {{ parent.username }} ({{ parent.firstName }} {{ parent.lastName }})
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="relationship" label="关系">
          <a-select v-model="addParentForm.relationship">
            <a-option value="father">父亲</a-option>
            <a-option value="mother">母亲</a-option>
            <a-option value="guardian">监护人</a-option>
            <a-option value="other">其他</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import MainLayout from '@/components/MainLayout.vue'
import userService, { USER_ROLES, USER_ROLE_NAMES } from '@/services/userService'
import type { UserProfile, UserRole } from '@/services/userService'

// 当前用户信息
const currentUser = ref<UserProfile | null>(null)
const currentUserRole = ref<string>('')

// 用户列表相关
const users = ref<UserProfile[]>([])
const loading = ref(false)
const selectedRole = ref('')
const searchKeyword = ref('')

// 用户关系相关
const teachersList = ref<UserProfile[]>([])
const studentsList = ref<UserProfile[]>([])
const parentsList = ref<UserProfile[]>([])
const teachersLoading = ref(false)
const studentsLoading = ref(false)
const parentsLoading = ref(false)

// 添加关系相关
const availableTeachers = ref<UserProfile[]>([])
const availableStudents = ref<UserProfile[]>([])
const availableParents = ref<UserProfile[]>([])
const teachersSelectLoading = ref(false)
const studentsSelectLoading = ref(false)
const parentsSelectLoading = ref(false)

// 模态框控制
const showUserDetailModal = ref(false)
const showEditRoleModal = ref(false)
const showAddTeacherModal = ref(false)
const showAddStudentModal = ref(false)
const showAddParentModal = ref(false)
const submitting = ref(false)

// 选中的用户
const selectedUser = ref<UserProfile | null>(null)

// 编辑角色表单
const editRoleForm = ref({
  role: '' as UserRole
})

// 添加教师表单
const addTeacherForm = ref({
  teacherId: null as number | null,
  department: '',
  position: ''
})

// 添加学生表单
const addStudentForm = ref({
  studentId: null as number | null,
  courseName: '',
  semester: ''
})

// 添加家长表单
const addParentForm = ref({
  parentId: null as number | null,
  relationship: 'guardian'
})

// 计算属性：过滤后的用户列表
const filteredUsers = computed(() => {
  let result = users.value

  // 按角色筛选
  if (selectedRole.value) {
    result = result.filter(user => user.role === selectedRole.value)
  }

  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(user => 
      user.username.toLowerCase().includes(keyword) || 
      user.email.toLowerCase().includes(keyword)
    )
  }

  return result
})

// 计算属性：可编辑的角色
const editableRoles = computed(() => {
  const roles = { ...USER_ROLE_NAMES }
  
  // 超级管理员不能被修改为其他角色
  if (selectedUser.value?.role === USER_ROLES.SUPER_ADMIN) {
    return {}
  }
  
  // 只有超级管理员可以设置管理员角色
  if (currentUserRole.value !== USER_ROLES.SUPER_ADMIN) {
    delete roles[USER_ROLES.SUPER_ADMIN]
    delete roles[USER_ROLES.ADMIN]
  }
  
  return roles
})

// 计算属性：是否可以编辑角色
const canEditRole = computed(() => {
  return currentUserRole.value === USER_ROLES.SUPER_ADMIN
})

// 计算属性：是否可以编辑状态
const canEditStatus = computed(() => {
  return currentUserRole.value === USER_ROLES.SUPER_ADMIN || 
         currentUserRole.value === USER_ROLES.ADMIN
})

// 计算属性：用户详情数据
const userDetailData = computed(() => {
  if (!selectedUser.value) return []
  
  return [
    { label: 'ID', value: selectedUser.value.id },
    { label: '用户名', value: selectedUser.value.username },
    { label: '邮箱', value: selectedUser.value.email },
    { label: '姓名', value: `${selectedUser.value.firstName} ${selectedUser.value.lastName}` },
    { label: '角色', value: USER_ROLE_NAMES[selectedUser.value.role as UserRole] || selectedUser.value.role },
    { label: '状态', value: selectedUser.value.status === 'active' ? '正常' : (selectedUser.value.status === 'inactive' ? '未激活' : '已禁用') },
    { label: '注册时间', value: formatDate(selectedUser.value.createdAt) }
  ]
})

// 获取当前用户信息
const getCurrentUser = async () => {
  try {
    const userProfile = await userService.getUserProfile()
    currentUser.value = userProfile
    currentUserRole.value = userProfile.role || ''
  } catch (error) {
    console.error('获取用户信息失败', error)
    Message.error('获取用户信息失败')
  }
}

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    if (currentUserRole.value === USER_ROLES.SUPER_ADMIN) {
      users.value = await userService.getAllUsers()
    } else if (currentUserRole.value === USER_ROLES.ADMIN) {
      users.value = await userService.getUsersByRole(USER_ROLES.TEACHER)
      const students = await userService.getUsersByRole(USER_ROLES.STUDENT)
      users.value = [...users.value, ...students]
    } else if (currentUserRole.value === USER_ROLES.TEACHER) {
      users.value = await userService.getUsersByRole(USER_ROLES.STUDENT)
    }
  } catch (error) {
    console.error('获取用户列表失败', error)
    Message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 获取用户关系
const fetchUserRelations = async () => {
  if (currentUserRole.value === USER_ROLES.SUPER_ADMIN || currentUserRole.value === USER_ROLES.ADMIN) {
    await fetchTeachers()
  }
  
  if (currentUserRole.value === USER_ROLES.SUPER_ADMIN || 
      currentUserRole.value === USER_ROLES.ADMIN || 
      currentUserRole.value === USER_ROLES.TEACHER) {
    await fetchStudents()
  }
  
  if (currentUserRole.value === USER_ROLES.STUDENT) {
    await fetchParents()
  }
}

// 获取教师列表
const fetchTeachers = async () => {
  teachersLoading.value = true
  try {
    teachersList.value = await userService.getTeachersByAdmin()
  } catch (error) {
    console.error('获取教师列表失败', error)
    Message.error('获取教师列表失败')
  } finally {
    teachersLoading.value = false
  }
}

// 获取学生列表
const fetchStudents = async () => {
  studentsLoading.value = true
  try {
    studentsList.value = await userService.getStudentsByTeacher()
  } catch (error) {
    console.error('获取学生列表失败', error)
    Message.error('获取学生列表失败')
  } finally {
    studentsLoading.value = false
  }
}

// 获取家长列表
const fetchParents = async () => {
  parentsLoading.value = true
  try {
    parentsList.value = await userService.getParentsByStudent()
  } catch (error) {
    console.error('获取家长列表失败', error)
    Message.error('获取家长列表失败')
  } finally {
    parentsLoading.value = false
  }
}

// 获取可用的教师列表
const fetchAvailableTeachers = async () => {
  teachersSelectLoading.value = true
  try {
    const allTeachers = await userService.getUsersByRole(USER_ROLES.TEACHER)
    // 过滤掉已经添加的教师
    const existingTeacherIds = new Set(teachersList.value.map(t => t.id))
    availableTeachers.value = allTeachers.filter(t => !existingTeacherIds.has(t.id))
  } catch (error) {
    console.error('获取可用教师列表失败', error)
    Message.error('获取可用教师列表失败')
  } finally {
    teachersSelectLoading.value = false
  }
}

// 获取可用的学生列表
const fetchAvailableStudents = async () => {
  studentsSelectLoading.value = true
  try {
    const allStudents = await userService.getUsersByRole(USER_ROLES.STUDENT)
    // 过滤掉已经添加的学生
    const existingStudentIds = new Set(studentsList.value.map(s => s.id))
    availableStudents.value = allStudents.filter(s => !existingStudentIds.has(s.id))
  } catch (error) {
    console.error('获取可用学生列表失败', error)
    Message.error('获取可用学生列表失败')
  } finally {
    studentsSelectLoading.value = false
  }
}

// 获取可用的家长列表
const fetchAvailableParents = async () => {
  parentsSelectLoading.value = true
  try {
    const allParents = await userService.getUsersByRole(USER_ROLES.PARENT)
    // 过滤掉已经添加的家长
    const existingParentIds = new Set(parentsList.value.map(p => p.id))
    availableParents.value = allParents.filter(p => !existingParentIds.has(p.id))
  } catch (error) {
    console.error('获取可用家长列表失败', error)
    Message.error('获取可用家长列表失败')
  } finally {
    parentsSelectLoading.value = false
  }
}

// 角色筛选变化
const handleRoleChange = () => {
  // 重新应用筛选
}

// 搜索
const handleSearch = () => {
  // 重新应用筛选
}

// 查看用户详情
const handleViewUser = (user: UserProfile) => {
  selectedUser.value = user
  showUserDetailModal.value = true
}

// 修改用户角色
const handleEditRole = (user: UserProfile) => {
  selectedUser.value = user
  editRoleForm.value.role = user.role as UserRole
  showEditRoleModal.value = true
}

// 提交角色修改
const handleRoleSubmit = async () => {
  if (!selectedUser.value) return
  
  submitting.value = true
  try {
    await userService.updateUserRole(selectedUser.value.id, editRoleForm.value.role)
    Message.success('角色修改成功')
    showEditRoleModal.value = false
    
    // 更新用户列表
    await fetchUsers()
  } catch (error) {
    console.error('修改角色失败', error)
    Message.error('修改角色失败')
  } finally {
    submitting.value = false
  }
}

// 切换用户状态
const handleToggleStatus = async (user: UserProfile) => {
  const newStatus = user.status === 'active' ? 'blocked' : 'active'
  
  try {
    await userService.updateUserStatus(user.id, newStatus)
    Message.success(`用户${newStatus === 'active' ? '启用' : '禁用'}成功`)
    
    // 更新用户列表
    await fetchUsers()
  } catch (error) {
    console.error('更新用户状态失败', error)
    Message.error('更新用户状态失败')
  }
}

// 添加教师
const handleAddTeacher = async () => {
  if (!addTeacherForm.value.teacherId) {
    Message.warning('请选择教师')
    return
  }
  
  submitting.value = true
  try {
    await userService.createAdminTeacherRelation(
      addTeacherForm.value.teacherId,
      addTeacherForm.value.department,
      addTeacherForm.value.position
    )
    Message.success('添加教师成功')
    showAddTeacherModal.value = false
    
    // 重置表单
    addTeacherForm.value = {
      teacherId: null,
      department: '',
      position: ''
    }
    
    // 刷新教师列表
    await fetchTeachers()
  } catch (error) {
    console.error('添加教师失败', error)
    Message.error('添加教师失败')
  } finally {
    submitting.value = false
  }
}

// 添加学生
const handleAddStudent = async () => {
  if (!addStudentForm.value.studentId) {
    Message.warning('请选择学生')
    return
  }
  
  submitting.value = true
  try {
    await userService.createTeacherStudentRelation(
      addStudentForm.value.studentId,
      undefined,
      addStudentForm.value.courseName,
      addStudentForm.value.semester
    )
    Message.success('添加学生成功')
    showAddStudentModal.value = false
    
    // 重置表单
    addStudentForm.value = {
      studentId: null,
      courseName: '',
      semester: ''
    }
    
    // 刷新学生列表
    await fetchStudents()
  } catch (error) {
    console.error('添加学生失败', error)
    Message.error('添加学生失败')
  } finally {
    submitting.value = false
  }
}

// 添加家长
const handleAddParent = async () => {
  if (!addParentForm.value.parentId) {
    Message.warning('请选择家长')
    return
  }
  
  submitting.value = true
  try {
    await userService.createStudentParentRelation(
      addParentForm.value.parentId,
      addParentForm.value.relationship
    )
    Message.success('添加家长成功')
    showAddParentModal.value = false
    
    // 重置表单
    addParentForm.value = {
      parentId: null,
      relationship: 'guardian'
    }
    
    // 刷新家长列表
    await fetchParents()
  } catch (error) {
    console.error('添加家长失败', error)
    Message.error('添加家长失败')
  } finally {
    submitting.value = false
  }
}

// 获取角色颜色
const getRoleColor = (role: string) => {
  switch (role) {
    case USER_ROLES.SUPER_ADMIN:
      return 'red'
    case USER_ROLES.ADMIN:
      return 'orange'
    case USER_ROLES.TEACHER:
      return 'blue'
    case USER_ROLES.STUDENT:
      return 'green'
    case USER_ROLES.PARENT:
      return 'purple'
    default:
      return 'gray'
  }
}

// 获取状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'warning'
    case 'blocked':
      return 'danger'
    default:
      return 'default'
  }
}

// 格式化日期
const formatDate = (dateStr?: string): string => {
  if (!dateStr) return '未知';
  
  try {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return '无效日期';
    
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    });
  } catch (error) {
    console.error('日期格式化错误', error);
    return '无效日期';
  }
}

// 页面加载时初始化
onMounted(async () => {
  await getCurrentUser()
  await fetchUsers()
  await fetchUserRelations()
})
</script>

<style scoped>
.user-management-container {
  padding: 20px;
}

.filter-container {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.relation-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
</style>
