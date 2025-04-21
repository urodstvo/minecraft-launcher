<script lang="ts" setup>
// import { Input } from '@/components/ui/input'
// import { Button } from '@/components/ui/button'
// import { ScrollArea } from '@/components/ui/scroll-area'
// import { LauncherService, type LauncherSettings } from '@go/launcher'
// import { ref, computed, onMounted } from 'vue'

// const settings = ref<LauncherSettings | null>(null)

// onMounted(async () => {
//   settings.value = await LauncherService.GetLauncherSettings()
// })

// function bindSetting<K extends keyof LauncherSettings>(key: K) {
//   return computed({
//     get: () => settings.value?.[key],
//     set: (val) => {
//       if (settings.value) {
//         settings.value[key] = val as LauncherSettings[K]
//       }
//     },
//   })
// }

// const gameDirectory = bindSetting('gameDirectory')
// const jvmArguments = bindSetting('jvmArguments')
// const resolutionHeight = bindSetting('resolutionHeight')
// const resolutionWidth = bindSetting('resolutionWidth')
// const showAlpha = bindSetting('showAlpha')
// const showBeta = bindSetting('showBeta')
// const showOldVersions = bindSetting('showOldVersions')
// const showOnlyInstalled = bindSetting('showOnlyInstalled')
// const showSnapshots = bindSetting('showSnapshots')
// const allocatedRAM = bindSetting('allocatedRAM')

// const save = () => {
//   LauncherService.SaveLauncherSettings({
//     jvmArguments: jvmArguments.value,
//     allocatedRAM: allocatedRAM.value,
//     resolutionHeight: resolutionHeight.value,
//     resolutionWidth: resolutionWidth.value,
//     gameDirectory: gameDirectory.value as string,
//     showAlpha: showAlpha.value as boolean,
//     showBeta: showBeta.value as boolean,
//     showOldVersions: showOldVersions.value as boolean,
//     showOnlyInstalled: showOnlyInstalled.value as boolean,
//     showSnapshots: showSnapshots.value as boolean,
//   })
// }
</script>

<template>
  <!-- <div class="flex-1 flex flex-col items-start justify-between p-5">
    <ScrollArea class="size-full">
      <div class="flex flex-col gap-4 w-full">
        <section class="flex flex-col gap-[10px]">
          <h4 class="scroll-m-20 text-sm font-semibold tracking-tight">Executable Path</h4>
          <div class="flex gap-5">
            <Input type="text" />
            <Button variant="secondary" class="w-[120px]">Browse</Button>
          </div>
        </section>
        <section class="flex flex-col gap-[10px]">
          <h4 class="scroll-m-20 text-sm font-semibold tracking-tight">JVM Arguments</h4>
          <Input type="text" v-model="jvmArguments" />
        </section>
      </div>
    </ScrollArea>
    <section class="flex-1 flex items-end justify-end size-full">
      <div class="flex gap-5 items-center">
        <Button class="mt-auto w-[160px]" @click="save" :disabled="false">Save</Button>
        <Button variant="secondary" class="mt-auto w-[120px]" disabled>Restore</Button>
      </div>
    </section>
  </div> -->
  <div class="flex-1 flex flex-col items-start justify-between p-5">
    <h4 class="scroll-m-20 text-xl font-semibold tracking-tight">IN DEVELOPMENT</h4>
  </div>
</template>
