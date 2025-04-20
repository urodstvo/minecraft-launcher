<script setup lang="ts">
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Combobox,
  ComboboxAnchor,
  ComboboxEmpty,
  ComboboxGroup,
  ComboboxInput,
  ComboboxItem,
  ComboboxItemIndicator,
  ComboboxList,
  ComboboxTrigger,
} from '@/components/ui/combobox'
import { Check, ChevronsUpDown, Search, Loader2 } from 'lucide-vue-next'
import { LauncherService, type LauncherSettings } from '@go/launcher'
import type { MinecraftVersionInfo } from '@go/minecraft'
import { useVirtualList } from '@vueuse/core'
import { computed, onMounted, ref, watch } from 'vue'

const settings = ref<LauncherSettings | null>(null)
const versions = ref<MinecraftVersionInfo[]>([])
const installed = ref<MinecraftVersionInfo[]>([])
const search = ref('')

onMounted(async () => {
  versions.value = (await LauncherService.GetMinecraftVersions()) ?? []
  settings.value = await LauncherService.GetLauncherSettings()
  installed.value = (await LauncherService.GetInstalledVersion()) ?? []
})

const filteredVersions = computed(() => {
  const show = settings.value
  const query = search.value.toLowerCase()
  return versions.value
    .filter((v) => show?.showAlpha || v.type !== 'old_alpha')
    .filter((v) => show?.showBeta || v.type !== 'old_beta')
    .filter((v) => show?.showSnapshots || v.type !== 'snapshot')
    .filter((v) => show?.showOldVersions || new Date(v.releaseTime) > new Date('2016-06-23T09:17:32+00:00'))
    .filter((v) => !show?.showOnlyInstalled || installed.value.some((i) => i.id === v.id))
    .filter((v) => v.id.toLowerCase().includes(query))
})

const { list, containerProps, wrapperProps, scrollTo } = useVirtualList(filteredVersions, {
  itemHeight: 32,
})
watch(filteredVersions, () => scrollTo(0))
watch(filteredVersions, () => console.log('filtered ', filteredVersions.value.length))
watch(list, () => console.log('virtual', list.value.length))

const selected = defineModel<MinecraftVersionInfo | null>('selected')
</script>

<template>
  <Combobox v-model="selected" by="id" class="w-[240px]">
    <ComboboxAnchor as-child>
      <ComboboxTrigger as-child>
        <Button variant="outline" class="justify-between w-full" :disabled="versions.length === 0">
          <span class="whitespace-nowrap overflow-hidden text-ellipsis">
            {{ selected ? selected.type + ' ' + selected.id : 'Select version' }}
          </span>
          <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" v-if="versions.length > 0" />
          <Loader2 class="w-4 h-4 mr-2 animate-spin" v-else />
        </Button>
      </ComboboxTrigger>
    </ComboboxAnchor>

    <ComboboxList class="w-full">
      <div class="relative w-full max-w-sm items-center">
        <ComboboxInput
          class="pl-9 focus-visible:ring-0 border-0 border-b rounded-none h-[36px]"
          placeholder="Select version..."
          v-model="search"
        />
        <span class="absolute start-0 inset-y-0 flex items-center justify-center px-3">
          <Search class="size-4 text-muted-foreground" />
        </span>
      </div>

      <ComboboxEmpty> No version found. </ComboboxEmpty>

      <div v-bind="containerProps" class="min-h-[100px] max-h-[300px] no-scrollbar">
        <div v-bind="wrapperProps">
          <ComboboxGroup>
            <ComboboxItem class="flex w-full" v-for="version in list" :key="version.data.id" :value="version.data">
              <div class="whitespace-nowrap overflow-hidden text-ellipsis flex-1">
                {{ version.data.type + ' ' + version.data.id }}
              </div>
              <ComboboxItemIndicator>
                <Check :class="cn('ml-auto h-4 w-4')" />
              </ComboboxItemIndicator>
            </ComboboxItem>
          </ComboboxGroup>
        </div>
      </div>
    </ComboboxList>
  </Combobox>
</template>
