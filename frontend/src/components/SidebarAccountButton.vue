<script lang="ts" setup>
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from '@/components/ui/sidebar'
import { Button } from '@/components/ui/button'
import { RouterLink } from 'vue-router'
import { ChevronsUpDown, Plus, Grid2X2, Square, Delete, Trash } from 'lucide-vue-next'

import { useAccountsStore } from '@/stores'

const store = useAccountsStore()
</script>

<template>
  <SidebarMenuButton
    size="lg"
    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
    as-child
    v-if="!store.selectedAccount"
  >
    <RouterLink to="/accounts" class="flex size-full items-center gap-2">
      <div class="flex aspect-square size-8 items-center justify-center rounded-lg">
        <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" viewBox="0 0 24 24" fill="none">
          <path
            d="M12 21L10 20M12 21L14 20M12 21V18.5M6 18L4 17V14.5M4 9.5V7M4 7L6 6M4 7L6 8M10 4L12 3L14 4M18 6L20 7M20 7L18 8M20 7V9.5M12 11L10 10M12 11L14 10M12 11V13.5M18 18L20 17V14.5"
            stroke="#ffffff"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </div>
      <div class="flex items-center flex-1 text-left text-sm leading-tight">
        <span class="truncate font-semibold text-muted-foreground">No Account</span>
      </div>
      <Plus class="size-4" />
    </RouterLink>
  </SidebarMenuButton>
  <SidebarMenu v-else>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <div class="flex aspect-square size-8 items-center justify-center rounded-lg">
              <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" viewBox="0 0 24 24" fill="none">
                <path
                  d="M12 21L10 20M12 21L14 20M12 21V18.5M6 18L4 17V14.5M4 9.5V7M4 7L6 6M4 7L6 8M10 4L12 3L14 4M18 6L20 7M20 7L18 8M20 7V9.5M12 11L10 10M12 11L14 10M12 11V13.5M18 18L20 17V14.5"
                  stroke="#ffffff"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">{{ store.selectedAccount.name }}</span>
            </div>
            <ChevronsUpDown />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          :align="'start'"
          :side="'right'"
          :side-offset="4"
        >
          <DropdownMenuLabel class="text-xs text-muted-foreground"> Accounts </DropdownMenuLabel>
          <DropdownMenuItem
            v-for="account in store.accounts"
            :key="account.id"
            class="w-full gap-2 p-2"
            @click="() => store.selectAccount(account.id)"
          >
            <div class="flex size-6 items-center justify-center rounded-sm">
              <Grid2X2 v-if="account.type === 'microsoft'" />
              <Square v-if="account.type === 'free'" />
            </div>
            {{ account.name }}
            <div class="flex-1 flex justify-end">
              <Button
                variant="ghost"
                size="icon"
                class="size-5 cursor-pointer"
                @click.stop="store.deleteAccount(account.id)"
              >
                <Trash />
              </Button>
            </div>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <RouterLink to="/accounts">
            <DropdownMenuItem class="gap-2 p-2">
              <div class="flex size-6 items-center justify-center rounded-md border bg-background">
                <Plus class="size-4" />
              </div>
              <div class="font-medium text-muted-foreground">Create Account</div>
            </DropdownMenuItem>
          </RouterLink>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
