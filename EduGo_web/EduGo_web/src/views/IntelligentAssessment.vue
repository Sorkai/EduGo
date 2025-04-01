<template>
  <MainLayout>
    <div class="assessment-container">
      <a-tabs default-active-key="student" v-model:activeKey="activeTab">
        <!-- 学生视图 -->
        <a-tab-pane key="student" title="参与测评">
          <div class="tab-content">
            <h2>可参与的测评</h2>
            <a-spin :loading="loading">
              <a-table :columns="studentColumns" :data="studentAssessments" :pagination="{ pageSize: 5 }">
                <template #status="{ record }">
                  <a-tag :color="getStatusColor(record.student_status)">
                    {{ getStatusText(record.student_status) }}
                  </a-tag>
                </template>
                <template #operation="{ record }">
                  <a-space>
                    <a-button 
                      type="primary" 
                      size="small" 
                      @click="viewAssessment(record, 'student')"
                      :disabled="!canTakeAssessment(record)"
                    >
                      {{ getActionText(record) }}
                    </a-button>
                    <a-button 
                      v-if="record.student_status === 'completed'" 
                      type="primary" 
                      size="small"
                      @click="viewResult(record)"
                    >
                      查看结果
                    </a-button>
                  </a-space>
                </template>
              </a-table>
            </a-spin>
          </div>
        </a-tab-pane>

        <!-- 教师视图 -->
        <a-tab-pane key="teacher" title="管理测评" v-if="isTeacher">
          <div class="tab-content">
            <div class="teacher-header">
              <h2>我创建的测评</h2>
              <a-button type="primary" @click="showCreateModal = true">创建测评</a-button>
            </div>
            <a-spin :loading="loading">
              <a-table :columns="teacherColumns" :data="teacherAssessments" :pagination="{ pageSize: 5 }">
                <template #status="{ record }">
                  <a-tag :color="getAssessmentStatusColor(record.status)">
                    {{ getAssessmentStatusText(record.status) }}
                  </a-tag>
                </template>
                <template #operation="{ record }">
                  <a-space>
                    <a-button type="primary" size="small" @click="viewAssessment(record, 'teacher')">
                      查看详情
                    </a-button>
                    <a-button 
                      v-if="record.status === 'published'" 
                      type="primary" 
                      status="warning" 
                      size="small"
                      @click="closeAssessment(record)"
                    >
                      关闭测评
                    </a-button>
                    <a-button 
                      v-if="record.status === 'published'" 
                      type="primary" 
                      status="success" 
                      size="small"
                      @click="viewStudents(record)"
                    >
                      查看学生
                    </a-button>
                  </a-space>
                </template>
              </a-table>
            </a-spin>
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>

    <!-- 创建测评模态框 -->
    <a-modal
      v-model:visible="showCreateModal"
      title="创建测评"
      @ok="createAssessment"
      @cancel="showCreateModal = false"
      ok-text="创建"
      cancel-text="取消"
    >
      <a-form :model="createForm" layout="vertical">
        <a-form-item field="title" label="测评标题" required>
          <a-input v-model="createForm.title" placeholder="请输入测评标题" />
        </a-form-item>
        <a-form-item field="description" label="测评描述">
          <a-textarea v-model="createForm.description" placeholder="请输入测评描述" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 测评详情模态框 (教师视图) -->
    <a-modal
      v-model:visible="showTeacherDetailModal"
      title="测评详情"
      :footer="false"
      width="800px"
    >
      <a-spin :loading="detailLoading">
        <div v-if="currentAssessment">
          <div class="assessment-header">
            <h2>{{ currentAssessment.title }}</h2>
            <p>{{ currentAssessment.description }}</p>
            <div class="assessment-meta">
              <a-tag :color="getAssessmentStatusColor(currentAssessment.status)">
                {{ getAssessmentStatusText(currentAssessment.status) }}
              </a-tag>
              <span>总分: {{ currentAssessment.total_score || 0 }}</span>
              <span v-if="currentAssessment.start_time">开始时间: {{ formatDate(currentAssessment.start_time) }}</span>
              <span v-if="currentAssessment.end_time">结束时间: {{ formatDate(currentAssessment.end_time) }}</span>
            </div>
          </div>

          <a-divider />

          <div class="questions-section">
            <div class="questions-header">
              <h3>题目列表</h3>
              <a-button 
                type="primary" 
                @click="showAddQuestionModal = true"
                :disabled="currentAssessment.status !== 'draft'"
              >
                添加题目
              </a-button>
            </div>

            <div v-if="questions.length === 0" class="empty-questions">
              <a-empty description="暂无题目" />
            </div>

            <div v-else class="question-list">
              <div v-for="(question, index) in questions" :key="question.id" class="question-item">
                <div class="question-header">
                  <h4>第 {{ index + 1 }} 题: {{ question.content }}</h4>
                  <span>分值: {{ question.score }}</span>
                </div>
                <div class="question-options">
                  <div v-for="(option, optIndex) in question.options" :key="optIndex" class="option-item">
                    <span :class="{ 'correct-answer': option === question.answer }">
                      {{ String.fromCharCode(65 + optIndex) }}. {{ option }}
                      <a-tag v-if="option === question.answer" color="green">正确答案</a-tag>
                    </span>
                  </div>
                </div>
                <div v-if="question.explanation" class="question-explanation">
                  <p>解析: {{ question.explanation }}</p>
                </div>
              </div>
            </div>
          </div>

          <a-divider />

          <div class="assessment-actions">
            <a-button 
              type="primary" 
              @click="showPublishModal = true"
              :disabled="currentAssessment.status !== 'draft' || questions.length === 0"
            >
              发布测评
            </a-button>
            <a-button @click="showTeacherDetailModal = false">关闭</a-button>
          </div>
        </div>
      </a-spin>
    </a-modal>

    <!-- 添加题目模态框 -->
    <a-modal
      v-model:visible="showAddQuestionModal"
      title="添加题目"
      @ok="addQuestion"
      @cancel="showAddQuestionModal = false"
      ok-text="添加"
      cancel-text="取消"
      width="700px"
    >
      <a-form :model="questionForm" layout="vertical">
        <a-form-item field="content" label="题目内容" required>
          <a-textarea v-model="questionForm.content" placeholder="请输入题目内容" />
        </a-form-item>
        <a-form-item field="options" label="选项" required>
          <div v-for="(option, index) in questionForm.options" :key="index" class="option-input">
            <a-input-group>
              <a-input-group-label>{{ String.fromCharCode(65 + index) }}</a-input-group-label>
              <a-input v-model="questionForm.options[index]" placeholder="请输入选项内容" />
              <a-button 
                v-if="index > 1" 
                type="text" 
                status="danger" 
                @click="removeOption(index)"
              >
                <icon-delete />
              </a-button>
            </a-input-group>
          </div>
          <div class="add-option">
            <a-button type="dashed" @click="addOption">
              <icon-plus />
              添加选项
            </a-button>
          </div>
        </a-form-item>
        <a-form-item field="answer" label="正确答案" required>
          <a-select v-model="questionForm.answer" placeholder="请选择正确答案">
            <a-option 
              v-for="(option, index) in questionForm.options" 
              :key="index" 
              :value="option"
              v-if="option.trim()"
            >
              {{ String.fromCharCode(65 + index) }}: {{ option }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="score" label="分值" required>
          <a-input-number v-model="questionForm.score" :min="1" :max="100" />
        </a-form-item>
        <a-form-item field="explanation" label="题目解析">
          <a-textarea v-model="questionForm.explanation" placeholder="请输入题目解析" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 发布测评模态框 -->
    <a-modal
      v-model:visible="showPublishModal"
      title="发布测评"
      @ok="publishAssessment"
      @cancel="showPublishModal = false"
      ok-text="发布"
      cancel-text="取消"
    >
      <a-form :model="publishForm" layout="vertical">
        <a-form-item field="timeRange" label="测评时间范围" required>
          <a-range-picker 
            v-model="publishForm.timeRange" 
            show-time 
            format="YYYY-MM-DD HH:mm:ss" 
            :disabled-date="disabledDate"
          />
        </a-form-item>
        <a-form-item field="studentIds" label="选择学生" required>
          <a-select 
            v-model="publishForm.studentIds" 
            placeholder="请选择学生" 
            multiple 
            :loading="studentsLoading"
          >
            <a-option 
              v-for="student in students" 
              :key="student.id" 
              :value="student.id"
            >
              {{ student.firstName }}{{ student.lastName }} ({{ student.username }})
            </a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 学生列表模态框 -->
    <a-modal
      v-model:visible="showStudentsModal"
      title="学生列表"
      :footer="false"
      width="800px"
    >
      <a-spin :loading="studentsLoading">
        <div v-if="currentAssessment">
          <h3>{{ currentAssessment.title }} - 学生列表</h3>
          <a-table :columns="studentListColumns" :data="assessmentStudents" :pagination="{ pageSize: 10 }">
            <template #status="{ record }">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            <template #score="{ record }">
              <span v-if="record.status === 'completed'">{{ record.score }}</span>
              <span v-else>-</span>
            </template>
            <template #time="{ record }">
              <div v-if="record.started_at">
                <div>开始: {{ formatDate(record.started_at) }}</div>
                <div v-if="record.completed_at">完成: {{ formatDate(record.completed_at) }}</div>
              </div>
              <span v-else>-</span>
            </template>
          </a-table>
        </div>
      </a-spin>
    </a-modal>

    <!-- 测评详情模态框 (学生视图) -->
    <a-modal
      v-model:visible="showStudentDetailModal"
      title="测评详情"
      :footer="false"
      width="800px"
    >
      <a-spin :loading="detailLoading">
        <div v-if="currentAssessment">
          <div class="assessment-header">
            <h2>{{ currentAssessment.title }}</h2>
            <p>{{ currentAssessment.description }}</p>
            <div class="assessment-meta">
              <a-tag :color="getAssessmentStatusColor(currentAssessment.status)">
                {{ getAssessmentStatusText(currentAssessment.status) }}
              </a-tag>
              <span>总分: {{ currentAssessment.total_score || 0 }}</span>
              <span v-if="currentAssessment.start_time">开始时间: {{ formatDate(currentAssessment.start_time) }}</span>
              <span v-if="currentAssessment.end_time">结束时间: {{ formatDate(currentAssessment.end_time) }}</span>
            </div>
          </div>

          <a-divider />

          <div v-if="!assessmentStarted">
            <a-result
              status="info"
              title="准备开始测评"
              sub-title="点击下方按钮开始测评，请确保在规定时间内完成。"
            >
              <template #extra>
                <a-button type="primary" @click="startAssessment">开始测评</a-button>
              </template>
            </a-result>
          </div>

          <div v-else class="questions-section">
            <h3>题目列表</h3>
            
            <div v-if="questions.length === 0" class="empty-questions">
              <a-empty description="暂无题目" />
            </div>

            <div v-else class="question-list">
              <div v-for="(question, index) in questions" :key="question.id" class="question-item">
                <div class="question-header">
                  <h4>第 {{ index + 1 }} 题: {{ question.content }}</h4>
                  <span>分值: {{ question.score }}</span>
                </div>
                <div class="question-options">
                  <a-radio-group v-model="studentAnswers[question.id]">
                    <div v-for="(option, optIndex) in question.options" :key="optIndex" class="option-item">
                      <a-radio :value="option">
                        {{ String.fromCharCode(65 + optIndex) }}. {{ option }}
                      </a-radio>
                    </div>
                  </a-radio-group>
                </div>
              </div>
            </div>

            <a-divider />

            <div class="assessment-actions">
              <a-button type="primary" @click="submitAnswers">提交答案</a-button>
              <a-button @click="cancelAssessment">取消</a-button>
            </div>
          </div>
        </div>
      </a-spin>
    </a-modal>

    <!-- 测评结果模态框 -->
    <a-modal
      v-model:visible="showResultModal"
      title="测评结果"
      :footer="false"
      width="800px"
    >
      <a-spin :loading="resultLoading">
        <div v-if="assessmentResult">
          <div class="result-header">
            <h2>{{ currentAssessment?.title }}</h2>
            <div class="score-display">
              <a-progress
                type="circle"
                :percent="Math.round((assessmentResult.your_score / (currentAssessment?.total_score || 100)) * 100)"
                :stroke-color="getScoreColor(assessmentResult.your_score, currentAssessment?.total_score)"
              />
              <div class="score-text">
                <h3>您的得分</h3>
                <p>{{ assessmentResult.your_score }} / {{ currentAssessment?.total_score }}</p>
              </div>
            </div>
          </div>

          <a-divider />

          <div class="ai-analysis">
            <h3>AI分析</h3>
            <a-alert type="info">
              <div style="white-space: pre-line">{{ assessmentResult.ai_analysis }}</div>
            </a-alert>
          </div>

          <a-divider />

          <div class="answer-details">
            <h3>答题详情</h3>
            <div v-for="(answer, index) in assessmentResult.answers" :key="index" class="answer-item">
              <div class="answer-header">
                <h4>
                  <a-tag :color="answer.is_correct ? 'green' : 'red'">
                    {{ answer.is_correct ? '正确' : '错误' }}
                  </a-tag>
                  {{ answer.content }}
                </h4>
              </div>
              <div class="answer-content">
                <p>您的答案: {{ answer.your_answer }}</p>
                <p v-if="!answer.is_correct">正确答案: {{ answer.correct_answer }}</p>
                <p v-if="answer.explanation">解析: {{ answer.explanation }}</p>
              </div>
            </div>
          </div>

          <a-divider />

          <div class="assessment-actions">
            <a-button @click="showResultModal = false">关闭</a-button>
          </div>
        </div>
      </a-spin>
    </a-modal>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus, IconDelete } from '@arco-design/web-vue/es/icon'
import MainLayout from '@/components/MainLayout.vue'
import axios from 'axios'
import dayjs from 'dayjs'
import config from '@/config'

// 用户角色
const userRole = ref('')
const isTeacher = computed(() => {
  return ['super_admin', 'admin', 'teacher'].includes(userRole.value)
})

// 当前激活的标签页
const activeTab = ref('student')

// 加载状态
const loading = ref(false)
const detailLoading = ref(false)
const studentsLoading = ref(false)
const resultLoading = ref(false)

// 模态框显示状态
const showCreateModal = ref(false)
const showTeacherDetailModal = ref(false)
const showStudentDetailModal = ref(false)
const showAddQuestionModal = ref(false)
const showPublishModal = ref(false)
const showStudentsModal = ref(false)
const showResultModal = ref(false)

// 表单数据
const createForm = reactive({
  title: '',
  description: ''
})

const questionForm = reactive({
  content: '',
  type: 'single_choice',
  options: ['', ''],
  answer: '',
  score: 10,
  explanation: ''
})

const publishForm = reactive({
  timeRange: [] as any[],
  studentIds: [] as number[]
})

// 表格列定义
const studentColumns = [
  { title: '测评标题', dataIndex: 'title' },
  { title: '总分', dataIndex: 'total_score' },
  { title: '开始时间', dataIndex: 'start_time' },
  { title: '结束时间', dataIndex: 'end_time' },
  { title: '状态', slotName: 'status' },
  { title: '操作', slotName: 'operation', width: 200 }
]

const teacherColumns = [
  { title: '测评标题', dataIndex: 'title' },
  { title: '总分', dataIndex: 'total_score' },
  { title: '开始时间', dataIndex: 'start_time' },
  { title: '结束时间', dataIndex: 'end_time' },
  { title: '状态', slotName: 'status' },
  { title: '操作', slotName: 'operation', width: 250 }
]

const studentListColumns = [
  { title: '用户名', dataIndex: 'username' },
  { title: '姓名', render: ({ record }: any) => `${record.firstName}${record.lastName}` },
  { title: '状态', slotName: 'status' },
  { title: '得分', slotName: 'score' },
  { title: '时间', slotName: 'time' }
]

// 数据
const studentAssessments = ref<any[]>([])
const teacherAssessments = ref<any[]>([])
const currentAssessment = ref<any>(null)
const questions = ref<any[]>([])
const students = ref<any[]>([])
const assessmentStudents = ref<any[]>([])
const assessmentResult = ref<any>(null)

// 学生答题相关
const assessmentStarted = ref(false)
const studentAnswers = ref<Record<number, string>>({})

// 获取用户信息
const getUserInfo = async () => {
  try {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token')
    if (!token) return
    
    const response = await axios.get(`${config.apiBaseUrl}/user`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    userRole.value = response.data.user.role || ''
    
    // 根据角色设置默认标签页
    if (isTeacher.value) {
      activeTab.value = 'teacher'
    } else {
      activeTab.value = 'student'
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

// 获取学生可参与的测评列表
const fetchStudentAssessments = async () => {
  loading.value = true
  try {
    const response = await axios.get(`${config.apiBaseUrl}/assessment/student`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    studentAssessments.value = response.data.assessments || []
  } catch (error) {
    console.error('获取测评列表失败:', error)
    Message.error('获取测评列表失败')
  } finally {
    loading.value = false
  }
}

// 获取教师创建的测评列表
const fetchTeacherAssessments = async () => {
  if (!isTeacher.value) return
  
  loading.value = true
  try {
    const response = await axios.get(`${config.apiBaseUrl}/assessment/teacher`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    teacherAssessments.value = response.data.assessments || []
  } catch (error) {
    console.error('获取测评列表失败:', error)
    Message.error('获取测评列表失败')
  } finally {
    loading.value = false
  }
}

// 获取测评详情
const fetchAssessmentDetail = async (id: number, role: 'teacher' | 'student') => {
  detailLoading.value = true
  try {
    const response = await axios.get(`${config.apiBaseUrl}/assessment/${role}/${id}`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    currentAssessment.value = response.data.assessment
    questions.value = response.data.questions || []
  } catch (error) {
    console.error('获取测评详情失败:', error)
    Message.error('获取测评详情失败')
  } finally {
    detailLoading.value = false
  }
}

// 获取教师的学生列表
const fetchTeacherStudents = async () => {
  if (!isTeacher.value) return
  
  studentsLoading.value = true
  try {
    const response = await axios.get(`${config.apiBaseUrl}/teacher/relations/students`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    students.value = response.data.students || []
  } catch (error) {
    console.error('获取学生列表失败:', error)
    Message.error('获取学生列表失败')
  } finally {
    studentsLoading.value = false
  }
}

// 获取测评的学生列表
const fetchAssessmentStudents = async (id: number) => {
  if (!isTeacher.value) return
  
  studentsLoading.value = true
  try {
    const response = await axios.get(`${config.apiBaseUrl}/assessment/teacher/${id}/students`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    assessmentStudents.value = response.data.students || []
  } catch (error) {
    console.error('获取测评学生列表失败:', error)
    Message.error('获取测评学生列表失败')
  } finally {
    studentsLoading.value = false
  }
}

// 获取测评结果
const fetchAssessmentResult = async (id: number) => {
  resultLoading.value = true
  try {
    const response = await axios.get(`${config.apiBaseUrl}/assessment/student/${id}/result`, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    assessmentResult.value = response.data.result
  } catch (error) {
    console.error('获取测评结果失败:', error)
    Message.error('获取测评结果失败')
  } finally {
    resultLoading.value = false
  }
}

// 创建测评
const createAssessment = async () => {
  if (!createForm.title) {
    Message.error('请输入测评标题')
    return
  }
  
  try {
    const response = await axios.post(`${config.apiBaseUrl}/assessment/teacher`, {
      title: createForm.title,
      description: createForm.description
    }, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    
    Message.success('测评创建成功')
    showCreateModal.value = false
    
    // 重置表单
    createForm.title = ''
    createForm.description = ''
    
    // 刷新测评列表
    fetchTeacherAssessments()
  } catch (error) {
    console.error('创建测评失败:', error)
    Message.error('创建测评失败')
  }
}

// 添加题目
const addQuestion = async () => {
  if (!questionForm.content) {
    Message.error('请输入题目内容')
    return
  }
  
  if (questionForm.options.filter(o => o.trim()).length < 2) {
    Message.error('至少需要两个选项')
    return
  }
  
  if (!questionForm.answer) {
    Message.error('请选择正确答案')
    return
  }
  
  try {
    const response = await axios.post(`${config.apiBaseUrl}/assessment/teacher/${currentAssessment.value.id}/question`, {
      content: questionForm.content,
      type: questionForm.type,
      options: questionForm.options.filter(o => o.trim()),
      answer: questionForm.answer,
      score: questionForm.score,
      explanation: questionForm.explanation
    }, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    
    Message.success('题目添加成功')
    showAddQuestionModal.value = false
    
    // 重置表单
    questionForm.content = ''
    questionForm.options = ['', '']
    questionForm.answer = ''
    questionForm.score = 10
    questionForm.explanation = ''
    
    // 刷新测评详情
    fetchAssessmentDetail(currentAssessment.value.id, 'teacher')
  } catch (error) {
    console.error('添加题目失败:', error)
    Message.error('添加题目失败')
  }
}

// 添加选项
const addOption = () => {
  questionForm.options.push('')
}

// 删除选项
const removeOption = (index: number) => {
  questionForm.options.splice(index, 1)
  if (questionForm.answer === questionForm.options[index]) {
    questionForm.answer = ''
  }
}

// 发布测评
const publishAssessment = async () => {
  if (!publishForm.timeRange || publishForm.timeRange.length !== 2) {
    Message.error('请选择测评时间范围')
    return
  }
  
  if (!publishForm.studentIds.length) {
    Message.error('请选择至少一名学生')
    return
  }
  
  try {
    const response = await axios.put(`${config.apiBaseUrl}/assessment/teacher/${currentAssessment.value.id}/publish`, {
      start_time: publishForm.timeRange[0].toISOString(),
      end_time: publishForm.timeRange[1].toISOString(),
      student_ids: publishForm.studentIds
    }, {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || sessionStorage.getItem('token')}` }
    })
    
    Message.success('测评发布成功')
    showPublishModal.value = false
    
    // 重置
