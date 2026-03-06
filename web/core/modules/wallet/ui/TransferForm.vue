<template>
  <div class="transfer-form-container">
    <h3 class="transfer-title">💸 Transfer Money</h3>
    
    <form @submit.prevent="handleSubmit" class="transfer-form">
      <div class="form-group">
        <label for="toWalletID" class="form-label">
          Recipient Wallet ID
        </label>
        <input
          id="toWalletID"
          v-model.number="form.toWalletID"
          type="number"
          class="form-input"
          :class="{ 'input-error': errors.toWalletID }"
          placeholder="Enter wallet ID"
          required
          :disabled="isLoading"
        />
        <p v-if="errors.toWalletID" class="error-message">
          {{ errors.toWalletID }}
        </p>
      </div>

      <div class="form-group">
        <label for="amount" class="form-label">
          Amount ({{ currency }})
        </label>
        <input
          id="amount"
          v-model.number="form.amount"
          type="number"
          step="0.01"
          min="0.01"
          class="form-input"
          :class="{ 'input-error': errors.amount }"
          placeholder="0.00"
          required
          :disabled="isLoading"
        />
        <p v-if="errors.amount" class="error-message">
          {{ errors.amount }}
        </p>
      </div>

      <div v-if="error" class="alert alert-error">
        ⚠️ {{ error }}
      </div>

      <div v-if="success" class="alert alert-success">
        ✅ {{ success }}
      </div>

      <button
        type="submit"
        class="submit-button"
        :disabled="isLoading || !isFormValid"
        :class="{ 'button-loading': isLoading }"
      >
        <span v-if="isLoading" class="spinner"></span>
        <span>{{ isLoading ? 'Processing...' : 'Transfer' }}</span>
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
interface Props {
  currency?: string
  isLoading?: boolean
  error?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  currency: 'THB',
  isLoading: false,
  error: null,
})

const emit = defineEmits<{
  transfer: [toWalletID: number, amount: number]
}>()

const form = ref({
  toWalletID: null as number | null,
  amount: null as number | null,
})

const errors = ref({
  toWalletID: '',
  amount: '',
})

const success = ref('')

const isFormValid = computed(() => {
  return form.value.toWalletID && 
         form.value.amount && 
         form.value.amount > 0 &&
         !errors.value.toWalletID &&
         !errors.value.amount
})

const validateForm = (): boolean => {
  errors.value = {
    toWalletID: '',
    amount: '',
  }

  if (!form.value.toWalletID || form.value.toWalletID <= 0) {
    errors.value.toWalletID = 'Please enter a valid wallet ID'
    return false
  }

  if (!form.value.amount || form.value.amount <= 0) {
    errors.value.amount = 'Amount must be greater than 0'
    return false
  }

  return true
}

const handleSubmit = () => {
  success.value = ''
  
  if (!validateForm()) {
    return
  }

  if (form.value.toWalletID && form.value.amount) {
    emit('transfer', form.value.toWalletID, form.value.amount)
  }
}

// Watch for successful transfer
watch(() => props.error, (newError) => {
  if (!newError && !props.isLoading && form.value.toWalletID && form.value.amount) {
    success.value = `Successfully transferred ${form.value.amount} ${props.currency}`
    
    // Reset form
    form.value = {
      toWalletID: null,
      amount: null,
    }
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      success.value = ''
    }, 3000)
  }
})
</script>

<style scoped>
.transfer-form-container {
  background: #1f2937;
  border: 1px solid #374151;
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2), 0 2px 4px -1px rgba(0, 0, 0, 0.1);
}

.transfer-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #f9fafb;
  margin: 0 0 1.5rem 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.transfer-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #9ca3af;
}

.form-input {
  padding: 0.75rem 1rem;
  border: 1px solid #4b5563;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: all 0.2s;
  background: #111827;
  color: #f9fafb;
}

.form-input::placeholder {
  color: #6b7280;
}

.form-input:focus {
  outline: none;
  border-color: #a855f7;
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.2);
}

.form-input:disabled {
  background: #374151;
  cursor: not-allowed;
  opacity: 0.7;
}

.input-error {
  border-color: #ef4444;
}

.input-error:focus {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2);
}

.error-message {
  font-size: 0.875rem;
  color: #f87171;
  margin: 0;
}

.alert {
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
}

.alert-error {
  background: rgba(239, 68, 68, 0.15);
  color: #fca5a5;
  border: 1px solid rgba(239, 68, 68, 0.4);
}

.alert-success {
  background: rgba(34, 197, 94, 0.15);
  color: #86efac;
  border: 1px solid rgba(34, 197, 94, 0.4);
}

.submit-button {
  padding: 0.875rem 1.5rem;
  background: linear-gradient(135deg, #7c3aed 0%, #db2777 100%);
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.submit-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 15px -3px rgba(124, 58, 237, 0.3), 0 4px 6px -2px rgba(219, 39, 119, 0.2);
}

.submit-button:active:not(:disabled) {
  transform: translateY(0);
}

.submit-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.button-loading {
  opacity: 0.7;
}

.spinner {
  width: 1rem;
  height: 1rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
