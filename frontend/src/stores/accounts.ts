import { AuthService } from '@go/auth'
import { LauncherService, LauncherAccount } from '@go/launcher'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAccountsStore = defineStore('accounts', () => {
  const accounts = ref<LauncherAccount[]>([])
  const selectedAccount = ref<LauncherAccount | null>(null)

  const loadAccounts = async () => {
    const { accounts: a, selectedAccount: sa } = await LauncherService.GetAccounts()
    accounts.value = a as LauncherAccount[]
    selectedAccount.value = a?.find((v) => v.id === sa) ?? null
  }

  const addFreeAccount = async (username: string) => {
    await AuthService.AddFreeAccount(username)
    await loadAccounts()
  }

  const addMicrosoftAccount = async () => {
    await AuthService.AddMicrosoftAccount()
    await loadAccounts()
  }

  const selectAccount = async (id: string) => {
    await LauncherService.SelectAccount(id)
    await loadAccounts()
  }

  const deleteAccount = async (id: string) => {
    await LauncherService.DeleteAccount(id)
    await loadAccounts()
  }

  return {
    accounts,
    selectedAccount,
    loadAccounts,
    addFreeAccount,
    addMicrosoftAccount,
    selectAccount,
    deleteAccount,
  }
})
