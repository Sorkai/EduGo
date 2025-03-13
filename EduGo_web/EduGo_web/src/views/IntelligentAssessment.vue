<template>
  <MainLayout>
    <h1>智能测评系统</h1>
    <div class="assessment-container">
      <div v-if="!assessmentStarted" class="start-section">
        <a-button type="primary" @click="startAssessment">开始测评</a-button>
      </div>

      <div v-if="assessmentStarted && !assessmentCompleted" class="question-section">
        <div class="question" v-for="(question, index) in questions" :key="question.id">
          <h3>第 {{ index + 1 }} 题: {{ question.text }}</h3>
          <a-radio-group v-model:value="answers[index]">
            <a-radio
              v-for="option in question.options"
              :key="option"
              :value="option"
            >
              {{ option }}
            </a-radio>
          </a-radio-group>
        </div>
        <a-button type="primary" @click="submitAnswers">提交答案</a-button>
      </div>

      <div v-if="assessmentCompleted" class="result-section">
        <h2>测评结果</h2>
        <p>您的得分: {{ score }}</p>
        <p>AI分析: {{ aiAnalysis }}</p>
        <a-button type="primary" @click="restartAssessment">重新测评</a-button>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MainLayout from '@/components/MainLayout.vue'
import { message } from 'ant-design-vue'

interface Question {
  id: number
  text: string
  options: string[]
}

const assessmentStarted = ref(false)
const assessmentCompleted = ref(false)
const questions = ref<Question[]>([
  {
    id: 1,
    text: '以下哪个是Go语言的特点？',
    options: ['面向对象', '垃圾回收', '动态类型', '解释执行']
  },
  {
    id: 2,
    text: 'Vue.js的核心特性是什么？',
    options: ['双向数据绑定', '单向数据流', '类组件', '服务端渲染']
  }
])
const answers = ref<string[]>([])
const score = ref(0)
const aiAnalysis = ref('')

const startAssessment = () => {
  assessmentStarted.value = true
  answers.value = new Array(questions.value.length).fill('')
}

const submitAnswers = async () => {
  // TODO: 调用后端API进行评分和AI分析
  score.value = Math.floor(Math.random() * 100)
  aiAnalysis.value = '根据您的答题情况，建议您加强基础知识的学习。'
  assessmentCompleted.value = true
  message.success('测评完成！')
}

const restartAssessment = () => {
  assessmentStarted.value = false
  assessmentCompleted.value = false
  answers.value = []
  score.value = 0
  aiAnalysis.value = ''
}
</script>

<style scoped>
.assessment-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.start-section {
  text-align: center;
  margin-top: 50px;
}

.question-section {
  margin-top: 30px;
}

.question {
  margin-bottom: 20px;
}

.result-section {
  text-align: center;
  margin-top: 50px;
}
</style>
