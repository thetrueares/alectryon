import { createRouter, createWebHistory } from 'vue-router'
import HomeView from './Views/HomeView.vue'
import InputsView from './Views/InputsView.vue'

const routes = [
  { path: '/', component: HomeView },
  { path: '/inputs', component: InputsView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
