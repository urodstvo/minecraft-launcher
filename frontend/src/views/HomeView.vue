<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import VersionSelect from '@/components/VersionSelect.vue'
import { LauncherService } from '@go/launcher'
import type { MinecraftVersionInfo } from '@go/minecraft'
import { computed, onMounted, ref } from 'vue'
import { Events } from '@wailsio/runtime'

const version = ref<MinecraftVersionInfo | null>(null)

const showProgress = ref(false)
const status = ref('')
const progress = ref(0)
const max = ref(0)

onMounted(async () => {
  version.value = await LauncherService.GetLastPlayedVersion()

  Events.On('install:status', ({ data }) => {
    status.value = data[0] as string
    console.log('[status]', status.value)
  })

  Events.On('install:max', ({ data }) => {
    max.value = +data[0] as number
    progress.value = 0
    console.log('[max]', max.value)
  })

  Events.On('install:progress', ({ data }) => {
    progress.value = Math.max(+data[0] as number, progress.value)
    console.log('[progress]', progress.value)
  })
})
const percent = computed(() => {
  if (max.value <= 0) return 0
  return Math.min(Math.floor((progress.value / max.value) * 100), 100)
})

const start = async () => {
  if (!version.value) return
  showProgress.value = true
  await LauncherService.StartMinecraft(version.value)
  showProgress.value = false
}
</script>

<template>
  <div class="flex-1 size-full grid overflow-hidden grid-rows-[1fr_88px]">
    <section class="size-full p-5"></section>
    <section class="size-full flex flex-col justify-end">
      <div class="w-full flex flex-col items-end h-[28px]" v-if="showProgress">
        <span class="text-sm text-muted-foreground px-5">{{ status }}</span>
        <Progress v-model="percent" class="rounded-none" />
      </div>
      <div class="w-full flex items-center justify-between h-[60px] p-5 bg-neutral-800/50">
        <VersionSelect v-model:selected="version" />
        <Button
          class="w-[200px] rounded-[2px] cursor-pointer text-xl hover:shadow-lg"
          @click="start"
          :disabled="showProgress"
        >
          PLAY
        </Button>
      </div>
    </section>
  </div>
</template>

<style scoped></style>
