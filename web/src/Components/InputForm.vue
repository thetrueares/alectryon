<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  initialData: {
    type: Object,
    default: () => ({
      name: '',
      type: 'telegram',
      active: true,
      options: { bot_token: '' }
    })
  },
  submitButtonText: {
    type: String,
    default: 'Submit'
  },
  loadingText: {
    type: String,
    default: 'Saving...'
  },
  isSaving: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit'])

const name = ref(props.initialData.name || '')
const type = ref(props.initialData.type || 'telegram')
const active = ref(props.initialData.active !== undefined ? props.initialData.active : true)
const options = ref({ ...props.initialData.options })

const types = ['telegram', 'slack', 'audio', 'video']

// Sync internal state if initialData changes (important for async fetching in Update view)
watch(() => props.initialData, (newData) => {
  if (newData) {
    name.value = newData.name || ''
    type.value = newData.type || 'telegram'
    active.value = newData.active !== undefined ? newData.active : true
    options.value = { ...newData.options }
  }
}, { deep: true })

// Handle type changes - initialize options if switching to a type that needs them
watch(type, (newType) => {
  if (newType === 'telegram' && !options.value.bot_token) {
    options.value.bot_token = ''
  }
})

const handleSubmit = () => {
  emit('submit', {
    name: name.value,
    type: type.value,
    active: active.value,
    options: options.value
  })
}
</script>

<template>
  <div class="bg-white p-8 rounded-lg shadow-sm border border-gray-200 max-w-2xl">
    <form @submit.prevent="handleSubmit" class="space-y-6">
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
          :disabled="isSaving"
          class="px-6 py-2 bg-black text-white rounded-md font-medium hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          {{ isSaving ? loadingText : submitButtonText }}
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
</template>
