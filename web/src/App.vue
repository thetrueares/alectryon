<script setup>
import { ref, onMounted } from 'vue'

const message = ref('Loading...')

onMounted(async () => {
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(apiUrl)
    const data = await response.json()
    message.value = data.message
  } catch (error) {
    message.value = 'Error fetching from API: ' + error.message
  }
})
</script>

<template>
  <main>
    <h1>Vue + Gin</h1>
    <p>Message from API: {{ message }}</p>
  </main>
</template>

<style>
#app {
  font-family: Arial, sans-serif;
  padding: 2rem;
  text-align: center;
}
</style>
