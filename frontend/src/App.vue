<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <main class="p-5 h-screen w-screen">
    <div class="size-full bg-zinc-400 rounded-[4px]">
      <RouterView />
    </div>
  </main>
</template>

<style scoped></style>
