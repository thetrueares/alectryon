<script setup>
import { ref, onMounted } from 'vue'

const inputs = ref([])
const loading = ref(true)
const error = ref(null)

const fetchInputs = async () => {
  try {
    loading.value = true
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    inputs.value = data.inputs
  } catch (e) {
    error.value = 'Error fetching inputs: ' + e.message
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchInputs()
})

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleString()
}
</script>

<template>
  <div>
    <h1 class="text-4xl font-bold text-black mb-8">Inputs</h1>

    <div v-if="loading" class="text-gray-600">
      Loading inputs...
    </div>

    <div v-else-if="error" class="bg-red-50 text-red-700 p-4 rounded-lg border border-red-200 max-w-2xl">
      {{ error }}
    </div>

    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200 text-left">
        <thead class="bg-gray-50 text-xs font-semibold text-gray-500 uppercase tracking-wider">
          <tr>
            <th class="px-6 py-4">Name</th>
            <th class="px-6 py-4">Type</th>
            <th class="px-6 py-4">Active</th>
            <th class="px-6 py-4">Created At</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white text-sm text-gray-700">
          <tr v-for="input in inputs" :key="input.id || input.name" class="hover:bg-gray-50 transition-colors">
            <td class="px-6 py-4 font-medium text-gray-900">{{ input.name }}</td>
            <td class="px-6 py-4">{{ input.type }}</td>
            <td class="px-6 py-4">
              <span 
                :class="[
                  'px-2 py-1 rounded-full text-xs font-medium',
                  input.active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                ]"
              >
                {{ input.active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-6 py-4">{{ formatDate(input.created_at) }}</td>
          </tr>
          <tr v-if="inputs.length === 0">
            <td colspan="4" class="px-6 py-8 text-center text-gray-500 italic">
              No inputs found.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
