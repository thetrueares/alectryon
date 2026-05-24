<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const inputId = route.params.id

const name = ref('')
const type = ref('telegram')
const active = ref(true)
const options = ref({})

const loading = ref(true)
const saving = ref(false)
const error = ref(null)

const types = ['telegram', 'slack', 'audio', 'video']

const fetchInput = async () => {
  try {
    loading.value = true
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs/${inputId}`)
    
    if (!response.ok) {
      throw new Error(`Failed to fetch input: ${response.status}`)
    }
    
    const data = await response.json()
    name.value = data.name
    type.value = data.type
    active.value = data.active
    options.value = data.options || {}
  } catch (e) {
    error.value = 'Error loading input: ' + e.message
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchInput)

// Handle type changes - initialize options if switching to a type that needs them
// but keep existing options if they exist
watch(type, (newType) => {
  if (newType === 'telegram' && !options.value.bot_token) {
    options.value.bot_token = ''
  }
})

const submitForm = async () => {
  try {
    saving.value = true
    error.value = null
    
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/inputs/${inputId}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: name.value,
        type: type.value,
        active: active.value,
        options: options.value
      }),
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

    <div v-else class="bg-white p-8 rounded-lg shadow-sm border border-gray-200 max-w-2xl">
      <form @submit.prevent="submitForm" class="space-y-6">
        <div>
          <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input
            id="name"
            v-model="name"
            type="text"
            required
            placeholder="e.g. Primary Telegram Bot"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
          />
        </div>

        <div>
          <label for="type" class="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select
            id="type"
            v-model="type"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white"
          >
            <option v-for="t in types" :key="t" :value="t">
              {{ t.charAt(0).toUpperCase() + t.slice(1) }}
            </option>
          </select>
        </div>

        <!-- Dynamic Options Section -->
        <div v-if="type === 'telegram'" class="pt-4 border-t border-gray-100">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Telegram Options</h3>
          <div>
            <label for="bot_token" class="block text-sm font-medium text-gray-700 mb-1">Bot Token</label>
            <input
              id="bot_token"
              v-model="options.bot_token"
              type="text"
              required
              placeholder="123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
            />
          </div>
        </div>

        <div class="flex items-center pt-2">
          <input
            id="active"
            v-model="active"
            type="checkbox"
            class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded cursor-pointer"
          />
          <label for="active" class="ml-2 block text-sm text-gray-900 cursor-pointer font-medium">
            Active
          </label>
        </div>

        <div class="pt-4 flex items-center space-x-4">
          <button
            type="submit"
            :disabled="saving"
            class="px-6 py-2 bg-black text-white rounded-md font-medium hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            {{ saving ? 'Updating...' : 'Update Input' }}
          </button>
          
          <router-link
            to="/inputs"
            class="px-6 py-2 bg-gray-100 text-gray-700 rounded-md font-medium hover:bg-gray-200 transition-all"
          >
            Cancel
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>
