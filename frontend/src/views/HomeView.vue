<script setup lang="ts">
import { Button } from '@/components/ui/button'
import VersionSelect from '@/components/VersionSelect.vue'
import { LauncherService } from '@go/launcher'
import type { MinecraftVersionInfo } from '@go/minecraft'
import { onMounted, ref } from 'vue'

const version = ref<MinecraftVersionInfo | null>(null)
onMounted(async () => {
  version.value = await LauncherService.GetLastPlayedVersion()
})
</script>

<template>
  <div class="flex-1 size-full grid grid-rows-[1fr_60px] overflow-hidden">
    <section class="size-full p-5"></section>
    <section class="size-full flex bg-neutral-800/50 items-center p-5">
      <div class="w-full flex items-center justify-between">
        <VersionSelect v-model:selected="version" />
        <Button
          class="w-[200px] rounded-[2px] cursor-pointer text-xl hover:shadow-lg"
          @click="() => LauncherService.StartMinecraft({ version })"
        >
          PLAY
        </Button>
      </div>
    </section>
  </div>
</template>

<style scoped></style>
