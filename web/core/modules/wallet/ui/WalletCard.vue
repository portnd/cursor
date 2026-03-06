<template>
  <div class="wallet-card">
    <div class="wallet-card-header">
      <div class="wallet-icon">
        💳
      </div>
      <div class="wallet-info">
        <h3 class="wallet-title">My Wallet</h3>
        <p class="wallet-id">ID: {{ walletID || 'Loading...' }}</p>
      </div>
    </div>

    <div class="wallet-balance-container">
      <div class="balance-label">Available Balance</div>
      <div class="balance-amount">
        <span class="currency">{{ currency }}</span>
        <span class="amount">{{ formatBalance(balance) }}</span>
      </div>
    </div>

    <div v-if="isLoading" class="wallet-loading">
      <div class="spinner"></div>
      <span>Loading...</span>
    </div>

    <div v-if="error" class="wallet-error">
      ⚠️ {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  walletID?: number | null
  balance: number
  currency: string
  isLoading?: boolean
  error?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  walletID: null,
  balance: 0,
  currency: 'THB',
  isLoading: false,
  error: null,
})

const formatBalance = (balance: number): string => {
  return balance.toLocaleString('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}
</script>

<style scoped>
.wallet-card {
  background: linear-gradient(135deg, #1f2937 0%, #374151 50%, #4c1d95 100%);
  border: 1px solid rgba(147, 51, 234, 0.3);
  border-radius: 1.5rem;
  padding: 2rem;
  color: white;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2), 0 10px 10px -5px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease, border-color 0.2s ease;
  position: relative;
  overflow: hidden;
}

.wallet-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  pointer-events: none;
}

.wallet-card:hover {
  transform: translateY(-4px);
  border-color: rgba(168, 85, 247, 0.5);
  box-shadow: 0 25px 30px -5px rgba(0, 0, 0, 0.2), 0 15px 15px -5px rgba(147, 51, 234, 0.1);
}

.wallet-card-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.wallet-icon {
  font-size: 3rem;
  line-height: 1;
}

.wallet-info {
  flex: 1;
}

.wallet-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 0.25rem 0;
}

.wallet-id {
  font-size: 0.875rem;
  opacity: 0.9;
  margin: 0;
}

.wallet-balance-container {
  text-align: center;
  padding: 1.5rem 0;
}

.balance-label {
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  opacity: 0.9;
  margin-bottom: 0.5rem;
}

.balance-amount {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 0.5rem;
}

.currency {
  font-size: 1.5rem;
  font-weight: 500;
  opacity: 0.9;
}

.amount {
  font-size: 3rem;
  font-weight: 700;
  line-height: 1;
}

.wallet-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  margin-top: 1rem;
}

.spinner {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.wallet-error {
  padding: 1rem;
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 0.5rem;
  margin-top: 1rem;
  font-size: 0.875rem;
}
</style>
