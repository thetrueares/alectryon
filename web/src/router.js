import { createRouter, createWebHistory } from 'vue-router'
import HomeView from './Views/HomeView.vue'
import InputsView from './Views/Inputs/InputsView.vue'
import CreateInputView from './Views/Inputs/CreateInputView.vue'

const routes = [
  { 
    path: '/', 
    component: HomeView,
    meta: { title: 'Home - Alectryon' }
  },
  { 
    path: '/inputs', 
    component: InputsView,
    meta: { title: 'Inputs - Alectryon' }
  },
  { 
    path: '/inputs/create', 
    component: CreateInputView,
    meta: { title: 'Create Input - Alectryon' }
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
