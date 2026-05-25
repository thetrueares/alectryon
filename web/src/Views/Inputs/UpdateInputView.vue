<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import InputForm from '../../Components/InputForm.vue'

const router = useRouter()
const route = useRoute()
const inputId = route.params.id

const initialData = ref(null)
const loading = ref(true)
const saving = ref(false)
const error = ref(null)

const fetchInput = async () => {
  try {
    loading.value = true
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs/${inputId}`)
    
    if (!response.ok) {
      throw new Error(`Failed to fetch input: ${response.status}`)
    }
    
    const data = await response.json()
    initialData.value = {
      name: data.name,
      type: data.type,
      active: data.active,
      options: data.options || {}
    }
  } catch (e) {
    error.value = 'Error loading input: ' + e.message
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchInput)

const handleUpdate = async (formData) => {
  try {
    saving.value = true
    error.value = null
    
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs/${inputId}`, {
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

    router.push('/inputs')
  } catch (e) {
    error.value = 'Error updating input: ' + e.message
    console.error(e)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="text-4xl font-bold text-black mb-8">Update Input</h1>

    <div v-if="loading" class="text-gray-600">
      Loading input data...
    </div>

    <div v-else-if="error" class="mb-6 bg-red-50 text-red-700 p-4 rounded-lg border border-red-200 max-w-2xl">
      {{ error }}
    </div>

    <InputForm
      v-else
      :initial-data="initialData"
      submit-button-text="Update Input"
      loading-text="Updating..."
      :is-saving="saving"
      @submit="handleUpdate"
    />
  </div>
</template>
