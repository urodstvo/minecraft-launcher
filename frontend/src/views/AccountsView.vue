<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { useAccountsStore } from '@/stores'
import { Separator } from '@/components/ui/separator'
import { MicrosoftIcon } from '@/components/icons'
import { ref } from 'vue'
import router from '@/router'

const accounts = useAccountsStore()
const username = ref<string>('')

const addFreeAccount = () => {
  if (!username.value.length) return

  accounts.addFreeAccount(username.value).then(() => {
    username.value = ''
    router.push('/')
  })
}

const addMicrosoftAccount = () => {
  if (!username.value.length) return

  accounts.addMicrosoftAccount().then(() => {
    router.push('/')
  })
}
</script>

<template>
  <div class="flex-1 flex gap-5 items-center justify-center">
    <div class="flex flex-col gap-5">
      <Label for="username-input">Minecraft Username</Label>
      <Input v-model="username" id="username-input" class="w-[400px]" autocomplete="off" name="new-password" />
      <Button @click="addFreeAccount" :disabled="!username.length">Create Free Account</Button>
      <Separator label="Or" />
      <Button @click="addMicrosoftAccount" disabled>Sign In with Microsoft Account<MicrosoftIcon /></Button>
    </div>
  </div>
</template>
