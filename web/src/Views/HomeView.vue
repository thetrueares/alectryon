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
  <div>
    <h1 class="text-4xl font-bold text-black mb-8">Home</h1>
    
    <div class="bg-white p-6 rounded-lg shadow-sm border border-gray-200 max-w-2xl">
      <p class="text-gray-700">Message from API: 
        <span class="font-mono bg-gray-100 px-2 py-1 rounded">{{ message }}</span>
      </p>
    </div>
  </div>
</template>
