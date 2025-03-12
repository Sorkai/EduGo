<template>
  <div class="main-layout">
    <header>
      <nav>
        <router-link to="/">首页</router-link>
        <router-link to="/assessment">智能测评</router-link>
        <router-link to="/teaching">智能教学</router-link>
        <router-link to="/vr">虚拟现实</router-link>
        <router-link to="/robot">教育机器人</router-link>
      </nav>
    </header>
    <main>
      <slot></slot>
    </main>
    <footer>
      <p>© 2025 EduGo 数智化教育平台</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

header {
  background-color: #2c3e50;
  padding: 1rem;
}

nav {
  display: flex;
  gap: 1rem;
}

nav a {
  color: white;
  text-decoration: none;
}

nav a.router-link-active {
  font-weight: bold;
}

main {
  flex: 1;
