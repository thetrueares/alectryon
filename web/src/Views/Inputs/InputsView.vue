<script setup>
import { ref, onMounted } from 'vue'

const inputs = ref([])
const loading = ref(true)
const error = ref(null)
const togglingId = ref(null)
const deletingId = ref(null)

const fetchInputs = async () => {
  try {
    loading.value = true
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    // Handle both { inputs: [...] } and direct array responses
    inputs.value = data.inputs || data
  } catch (e) {
    error.value = 'Error fetching inputs: ' + e.message
    console.error(e)
  } finally {
    loading.value = false
  }
}

const toggleStatus = async (input) => {
  if (togglingId.value) return
  
  try {
    togglingId.value = input.id
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs/${input.id}/toggle`, {
      method: 'POST',
    })

    if (!response.ok) {
      throw new Error(`Failed to toggle status: ${response.status}`)
    }

    const updatedInput = await response.json()
    
    // Update local state
    const index = inputs.value.findIndex(i => i.id === input.id)
    if (index !== -1) {
      inputs.value[index] = { ...inputs.value[index], ...updatedInput }
    }
  } catch (e) {
    console.error('Toggle error:', e)
    alert('Error updating status: ' + e.message)
  } finally {
    togglingId.value = null
  }
}

const deleteInput = async (input) => {
  if (deletingId.value) return
  if (!confirm(`Are you sure you want to delete "${input.name}"?`)) return

  try {
    deletingId.value = input.id
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs/${input.id}`, {
      method: 'DELETE',
    })

    if (!response.ok) {
      throw new Error(`Failed to delete input: ${response.status}`)
    }

    // Remove from local state
    inputs.value = inputs.value.filter(i => i.id !== input.id)
  } catch (e) {
    console.error('Delete error:', e)
    alert('Error deleting input: ' + e.message)
  } finally {
    deletingId.value = null
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
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-4xl font-bold text-black">Inputs</h1>
      <router-link
        to="/inputs/create"
        class="px-4 py-2 bg-black text-white rounded-md font-medium hover:bg-gray-800 transition-all"
      >
        + Create
      </router-link>
    </div>

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
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4">Created At</th>
            <th class="px-6 py-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white text-sm text-gray-700">
          <tr v-for="input in inputs" :key="input.id" class="hover:bg-gray-50 transition-colors">
            <td class="px-6 py-4 font-medium text-gray-900">{{ input.name }}</td>
            <td class="px-6 py-4 capitalize">{{ input.type }}</td>
            <td class="px-6 py-4">
              <span 
                :class="[
                  'px-2 py-1 rounded-full text-xs font-medium inline-flex items-center',
                  input.active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                ]"
              >
                <span :class="['w-1.5 h-1.5 rounded-full mr-1.5', input.active ? 'bg-green-500' : 'bg-red-500']"></span>
                {{ input.active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-6 py-4">{{ formatDate(input.created_at) }}</td>
            <td class="px-6 py-4 text-right space-x-2">
              <router-link
                :to="`/inputs/${input.id}/edit`"
                class="text-xs font-semibold px-3 py-1 rounded border border-gray-200 text-gray-600 hover:bg-gray-50 transition-all inline-block"
              >
                Edit
              </router-link>

              <button 
                @click="toggleStatus(input)"
                :disabled="togglingId === input.id || deletingId === input.id"
                :class="[
                  'text-xs font-semibold px-3 py-1 rounded border transition-all',
                  input.active 
                    ? 'border-gray-200 text-gray-600 hover:bg-gray-50' 
                    : 'border-green-200 text-green-600 hover:bg-green-50',
                  (togglingId === input.id || deletingId === input.id) ? 'opacity-50 cursor-not-allowed' : ''
                ]"
              >
                {{ togglingId === input.id ? 'Updating...' : (input.active ? 'Deactivate' : 'Activate') }}
              </button>
              
              <button 
                @click="deleteInput(input)"
                :disabled="deletingId === input.id || togglingId === input.id"
                class="text-xs font-semibold px-3 py-1 rounded border border-red-200 text-red-600 hover:bg-red-50 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ deletingId === input.id ? 'Deleting...' : 'Delete' }}
              </button>
            </td>
          </tr>
          <tr v-if="inputs.length === 0">
            <td colspan="5" class="px-6 py-8 text-center text-gray-500 italic">
              No inputs found.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
