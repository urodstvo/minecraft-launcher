<script setup lang="ts">
import { Button } from '@/components/ui/button'
import VersionSelect from '@/components/VersionSelect.vue'
</script>

<template>
  <div class="size-full grid grid-rows-[1fr_100px]">
    <div class="size-full"></div>
    <section class="size-full bg-zinc-500 rounded-b-[4px] flex items-center p-5">
      <div class="w-full flex items-end">
        <VersionSelect />
        <Button class="w-[300px] h-[60px] ml-[114px] mr-[154px] cursor-pointer">Start</Button>
        <div class="flex gap-5 items-center">
          <Button :size="'icon'" class="size-[40px]"></Button>
          <Button :size="'icon'" class="size-[40px]"></Button>
          <Button :size="'icon'" class="size-[40px]"></Button>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped></style>
