import { createRouter, createWebHistory } from 'vue-router'
import HomeView from './Views/HomeView.vue'
import InputsView from './Views/Inputs/InputsView.vue'
import CreateInputView from './Views/Inputs/CreateInputView.vue'
import UpdateInputView from './Views/Inputs/UpdateInputView.vue'

const routes = [
  { 
    path: '/', 
    component: HomeView,
    meta: { title: 'Home - Alectryon' }
  },
  { 
    path: '/channels', 
    component: InputsView,
    meta: { title: 'Channels - Alectryon' }
  },
  { 
    path: '/channels/create', 
    component: CreateInputView,
    meta: { title: 'Create Channel - Alectryon' }
  },
  { 
    path: '/channels/:id/edit', 
    component: UpdateInputView,
    meta: { title: 'Update Channel - Alectryon' }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.afterEach((to) => {
  document.title = to.meta.title || 'Alectryon'
})

export default router
