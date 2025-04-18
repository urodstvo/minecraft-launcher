<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <main class="p-5 h-screen w-screen bg-[#2e2e2e]">
    <div class="size-full overflow-auto bg-[#e0e0e0]/10 rounded-[4px]">
      <RouterView />
    </div>
  </main>
</template>

<style scoped></style>
