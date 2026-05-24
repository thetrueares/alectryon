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
  <div class="bg-white p-6 rounded-lg shadow-md max-w-md w-full text-center">
    <h1 class="text-3xl font-bold text-blue-600 mb-4">Home</h1>
    <p class="text-gray-700">Message from API: 
      <span class="font-mono bg-gray-200 px-2 py-1 rounded">{{ message }}</span>
    </p>
  </div>
</template>
