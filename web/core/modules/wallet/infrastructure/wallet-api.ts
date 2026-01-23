/**
 * Wallet API Integration
 * 
 * Provides functions to interact with the wallet backend API.
 */

import { useHttp } from '~/core/shared/api/http'

export interface Wallet {
  id: number
  user_id: number
  balance: number
  currency: string
}

export interface Transaction {
  id: number
  from_wallet_id?: number
  to_wallet_id?: number
  amount: number
  type: string
  status: string
  description?: string
  created_at: string
}

export interface ApiResponse<T> {
  message: string
  data: T
}

export const walletApi = {
  /**
   * Get current user's wallet
   */
  async getMyWallet() {
    const http = useHttp()
    const { data, error } = await http.get<ApiResponse<Wallet>>('/wallets/me')

    if (error) {
      throw new Error(error.message || 'Failed to fetch wallet')
    }

    return data as ApiResponse<Wallet>
  },

  /**
   * Transfer money to another wallet
   */
  async transfer(toWalletID: number, amount: number) {
    const http = useHttp()
    const { data, error } = await http.post<ApiResponse<Transaction>>('/wallets/transfer', {
      to_wallet_id: toWalletID,
      amount: amount,
    })

    if (error) {
      throw new Error(error.message || 'Transfer failed')
    }

    return data as ApiResponse<Transaction>
  },

  /**
   * Get transaction history
   */
  async getTransactions(limit: number = 20) {
    const http = useHttp()
    const { data, error } = await http.get<ApiResponse<Transaction[]>>(`/wallets/transactions?limit=${limit}`)

    if (error) {
      throw new Error(error.message || 'Failed to fetch transactions')
    }

    return data as ApiResponse<Transaction[]>
  },
}
