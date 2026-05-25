<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import InputForm from '../../Components/InputForm.vue'

const router = useRouter()
const loading = ref(false)
const error = ref(null)

const handleCreate = async (formData) => {
  try {
    loading.value = true
    error.value = null
    
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formData),
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
    }

    router.push('/channels')
  } catch (e) {
    error.value = 'Error creating input: ' + e.message
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="text-4xl font-bold text-black mb-8">Create Channel</h1>

    <div v-if="error" class="mb-6 bg-red-50 text-red-700 p-4 rounded-lg border border-red-200 max-w-2xl">
      {{ error }}
    </div>

    <InputForm
      submit-button-text="Create Channel"
      loading-text="Creating..."
      :is-saving="loading"
      @submit="handleCreate"
    />
  </div>
</template>
