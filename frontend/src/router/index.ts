import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import ProfileView from '@/views/ProfileView.vue'
import MinecraftSettingsView from '@/views/MinecraftSettingsView.vue'
import JavaSettingsView from '@/views/JavaSettingsView.vue'
import AppearenceSettingsView from '@/views/AppearenceSettingsView.vue'
import AccountsView from '@/views/AccountsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/settings/minecraft',
      name: 'minecraft-settings',
      component: MinecraftSettingsView,
    },
    {
      path: '/settings/java',
      name: 'java-settings',
      component: JavaSettingsView,
    },
    {
      path: '/settings/appearence',
      name: 'appearence-settings',
      component: AppearenceSettingsView,
    },
    {
      path: '/profiles',
      name: 'profiles',
      component: ProfileView,
    },
    {
      path: '/accounts',
      name: 'accounts',
      component: AccountsView,
    },
  ],
})

export default router
