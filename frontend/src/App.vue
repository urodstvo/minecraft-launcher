<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
  SidebarInset,
  SidebarHeader,
  SidebarMenuButton,
  SidebarGroupLabel,
  SidebarFooter,
} from '@/components/ui/sidebar'
import { Separator } from './components/ui/separator'
import { Home, Settings, User, Folder, RotateCw, Paintbrush, Wrench } from 'lucide-vue-next'
import { LauncherService } from '@go/launcher'
import { onMounted, ref } from 'vue'
import SidebarAccountButton from './components/SidebarAccountButton.vue'

const launcherVersion = ref<string>('')
onMounted(async () => {
  launcherVersion.value = await LauncherService.GetLauncherVersion()
})
const nav = [
  { name: 'Play', icon: Home, href: '/' },
  { name: 'Profiles', icon: User, href: '/profiles' },
]

const settings = [
  { name: 'Minecraft', icon: Settings, href: '/settings/minecraft' },
  { name: 'Appearence', icon: Paintbrush, href: '/settings/appearence' },
  { name: 'Java', icon: Wrench, href: '/settings/java' },
]

const launcher = [
  { name: 'Game Folder', icon: Folder, click: LauncherService.OpenMinecraftDirectory },
  { name: 'Update Launcher', icon: RotateCw, click: () => {} },
]
</script>

<template>
  <SidebarProvider class="items-start h-screen" style="--sidebar-width: 200px; --sidebar-width-mobile: 200px">
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarAccountButton />
      </SidebarHeader>
      <Separator />
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in nav" :key="item.name">
                <SidebarMenuButton as-child :is-active="item.name === 'Messages & media'">
                  <RouterLink :to="item.href">
                    <component :is="item.icon" />
                    <span>{{ item.name }}</span>
                  </RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarGroup class="p-0">
          <SidebarGroupLabel>Settings</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in settings" :key="item.name">
                <SidebarMenuButton as-child :is-active="item.name === 'Messages & media'">
                  <RouterLink :to="item.href">
                    <component :is="item.icon" />
                    <span>{{ item.name }}</span>
                  </RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarMenu>
          <SidebarMenuItem v-for="item in launcher" :key="item.name">
            <SidebarMenuButton :is-active="item.name === 'Messages & media'" @click="item.click">
              <component :is="item.icon" />
              <span>{{ item.name }}</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
    <SidebarInset>
      <div class="flex flex-1 flex-col overflow-hidden min-h-screen">
        <header
          class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
        >
          <div class="w-full flex items-center justify-between px-4">
            <SidebarTrigger class="-ml-1" />
            <span class="text-sm text-muted-foreground">v. {{ launcherVersion }}</span>
          </div>
        </header>
        <Separator />
        <RouterView />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>

<style scoped></style>
