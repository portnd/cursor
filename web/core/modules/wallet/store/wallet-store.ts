/**
 * Wallet Store (Pinia)
 * 
 * Manages wallet state and operations.
 */

import { defineStore } from 'pinia'
import { walletApi, type Wallet, type Transaction } from '../infrastructure/wallet-api'

export const useWalletStore = defineStore('wallet', {
  state: () => ({
    wallet: null as Wallet | null,
    transactions: [] as Transaction[],
    isLoading: false,
    error: null as string | null,
  }),

  getters: {
    balance: (state) => state.wallet?.balance || 0,
    walletID: (state) => state.wallet?.id || null,
    currency: (state) => state.wallet?.currency || 'THB',
    hasWallet: (state) => state.wallet !== null,
  },

  actions: {
    /**
     * Fetch user's wallet
     */
    async fetchWallet() {
      this.isLoading = true
      this.error = null

      try {
        const response = await walletApi.getMyWallet()
        
        if (response && response.data) {
          this.wallet = response.data
        }
      } catch (error: any) {
        this.error = error.message || 'Failed to fetch wallet'
        console.error('Fetch wallet error:', error)
      } finally {
        this.isLoading = false
      }
    },

    /**
     * Transfer money to another wallet
     */
    async transfer(toWalletID: number, amount: number) {
      this.isLoading = true
      this.error = null

      try {
        const response = await walletApi.transfer(toWalletID, amount)

        if (response && response.data) {
          // Refresh wallet balance after successful transfer
          await this.fetchWallet()
          
          // Add transaction to local state
          this.transactions.unshift(response.data)

          return { success: true, transaction: response.data }
        }
      } catch (error: any) {
        this.error = error.message || 'Transfer failed'
        console.error('Transfer error:', error)
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }

      return { success: false, error: this.error }
    },

    /**
     * Fetch transaction history
     */
    async fetchTransactions(limit: number = 20) {
      this.isLoading = true
      this.error = null

      try {
        const response = await walletApi.getTransactions(limit)

        if (response && response.data) {
          this.transactions = response.data
        }
      } catch (error: any) {
        this.error = error.message || 'Failed to fetch transactions'
        console.error('Fetch transactions error:', error)
      } finally {
        this.isLoading = false
      }
    },

    /**
     * Clear error
     */
    clearError() {
      this.error = null
    },
  },
})
