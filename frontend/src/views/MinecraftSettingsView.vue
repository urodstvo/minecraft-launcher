<script lang="ts" setup>
import { Slider } from '@/components/ui/slider'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Checkbox } from '@/components/ui/checkbox'
import { LauncherService, type LauncherSettings } from '@go/launcher'
import { ref, computed, onMounted, watch, toRaw } from 'vue'

const settings = ref<LauncherSettings | null>(null)
const minRAM = 0 * 1024
let maxRAM = 16384
const totalRAM = ref<number>(16 * 1024)
onMounted(() => {
  LauncherService.GetTotalRAM().then((res) => {
    totalRAM.value = res
    maxRAM = totalRAM.value
  })
  LauncherService.GetLauncherSettings().then((res) => {
    settings.value = res
    RAM.value = [res.allocatedRAM ?? 2048]
  })
})

function generateRamOptions(maxMb: number): string[] {
  const result: string[] = ['0GB']
  const maxGb = Math.ceil(maxMb / 1024)

  for (let gb = 4; gb <= maxGb; gb += 4) {
    result.push(`${gb}GB`)
  }

  if (maxGb % 4 !== 0) {
    const d = 4 - (maxGb % 4)
    maxRAM = (maxGb + d) * 1024
    const last = `${maxGb + d}GB`
    if (!result.includes(last)) {
      result.push(last)
    }
  }

  return result
}

const marks = computed(() => generateRamOptions(totalRAM.value))

function bindSetting<K extends keyof LauncherSettings>(key: K) {
  return computed({
    get: () => settings.value?.[key],
    set: (val) => {
      if (settings.value) {
        settings.value[key] = val as LauncherSettings[K]
      }
    },
  })
}

const gameDirectory = bindSetting('gameDirectory')
const jvmArguments = bindSetting('jvmArguments')
const resolutionHeight = bindSetting('resolutionHeight')
const resolutionWidth = bindSetting('resolutionWidth')
const showAlpha = bindSetting('showAlpha')
const showBeta = bindSetting('showBeta')
const showOldVersions = bindSetting('showOldVersions')
const showOnlyInstalled = bindSetting('showOnlyInstalled')
const showSnapshots = bindSetting('showSnapshots')
const allocatedRAM = bindSetting('allocatedRAM')
const RAM = ref<number[]>([2048])

const ramDisplay = computed({
  get() {
    return allocatedRAM.value && allocatedRAM.value > 0 ? `${allocatedRAM.value} Mb` : ''
  },
  set(val: string) {
    const parsed = parseInt(val.replace(/\D/g, ''), 10)

    allocatedRAM.value = isNaN(parsed) ? 0 : Math.min(parsed, totalRAM.value)
  },
})

watch(
  RAM,
  () => {
    if (RAM.value[0] > totalRAM.value) {
      RAM.value[0] = totalRAM.value
    }
    allocatedRAM.value = RAM.value[0]
  },
  { deep: true },
)

const save = () => {
  LauncherService.SaveLauncherSettings({
    jvmArguments: jvmArguments.value,
    allocatedRAM: allocatedRAM.value,
    resolutionHeight: resolutionHeight.value,
    resolutionWidth: resolutionWidth.value,
    gameDirectory: gameDirectory.value as string,
    showAlpha: showAlpha.value as boolean,
    showBeta: showBeta.value as boolean,
    showOldVersions: showOldVersions.value as boolean,
    showOnlyInstalled: showOnlyInstalled.value as boolean,
    showSnapshots: showSnapshots.value as boolean,
  })
}

const chooseGameDirectory = async () => {
  gameDirectory.value = await LauncherService.ChooseDirectory()
}
</script>

<template>
  <div class="flex-1 flex flex-col items-start justify-between p-5">
    <div class="flex flex-col gap-4 w-full">
      <section class="flex flex-col gap-[10px]">
        <h4 class="scroll-m-20 text-sm font-semibold tracking-tight">Game Directory</h4>
        <div class="flex gap-5">
          <Input type="text" v-model="gameDirectory" />
          <Button variant="secondary" class="w-[120px]" @click="chooseGameDirectory">Browse</Button>
        </div>
      </section>
      <section class="flex flex-col gap-[10px]">
        <h4 class="scroll-m-20 text-sm font-semibold tracking-tight">Resolution</h4>
        <div class="flex gap-5 items-center w-full justify-center">
          <Input type="text" class="w-[150px]" placeholder="<auto>" v-model="resolutionWidth" />
          x
          <Input type="text" class="w-[150px]" placeholder="<auto>" v-model="resolutionHeight" />
        </div>
      </section>
      <Separator />
      <section class="flex flex-col gap-[10px]">
        <h4 class="scroll-m-20 text-sm font-semibold tracking-tight">Memory (RAM)</h4>
        <div class="flex gap-5 items-center">
          <div class="flex-1">
            <Slider :max="maxRAM" :step="1024" :min="minRAM" v-model="RAM" class="w-full" />
            <div class="mt-2.5 flex items-center justify-between text-muted-foreground text-xs w-full">
              <span v-for="mark in marks" :key="mark" class="w-[24px]">{{ mark }}</span>
            </div>
          </div>

          <Input v-model="ramDisplay" inputmode="numeric" pattern="\d*" maxlength="10" class="w-[120px]" />
        </div>
      </section>
      <Separator />
      <section class="flex flex-col gap-[10px]">
        <h4 class="scroll-m-20 text-sm font-semibold tracking-tight">Version List</h4>
        <div class="flex items-center gap-5">
          <div class="flex gap-[10px] items-center">
            <Checkbox id="alpha_v" v-model="showAlpha" />
            <Label :for="'alpha_v'">Alpha (2010)</Label>
          </div>
          <div class="flex gap-[10px] items-center">
            <Checkbox id="beta_v" v-model="showBeta" />
            <Label :for="'beta_v'">Beta (2010-2011)</Label>
          </div>
          <div class="flex gap-[10px] items-center">
            <Checkbox id="snapshots_v" v-model="showSnapshots" />
            <Label :for="'snapshots_v'">Snapshots</Label>
          </div>
          <div class="flex gap-[10px] items-center">
            <Checkbox id="old_v" v-model="showOldVersions" />
            <Label :for="'old_v'">Old releases (< 1.10.2)</Label>
          </div>
          <div class="flex gap-[10px] items-center">
            <Checkbox id="installed_v" v-model="showOnlyInstalled" />
            <Label :for="'installed_v'">Only installed</Label>
          </div>
        </div>
      </section>
    </div>
    <section class="flex-1 flex items-end justify-end size-full">
      <div class="flex gap-5 items-center">
        <Button class="mt-auto w-[160px]" @click="save" :disabled="false">Save</Button>
        <Button variant="secondary" class="mt-auto w-[120px]" disabled>Restore</Button>
      </div>
    </section>
  </div>
</template>
