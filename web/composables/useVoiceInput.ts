import { ref, onMounted, onUnmounted } from 'vue'

export function useVoiceInput(onResult: (text: string) => void, lang = 'th-TH') {
  const isListening = ref(false)
  const isSupported = ref(false)
  const error = ref<string | null>(null)

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let recognition: any = null

  onMounted(() => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const SR = (window as any).SpeechRecognition ?? (window as any).webkitSpeechRecognition
    if (!SR) return
    isSupported.value = true
    initRecognition(SR)
  })

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  function initRecognition(SR: any) {
    recognition = new SR()
    recognition.continuous = false
    recognition.interimResults = false
    recognition.lang = lang

    recognition.onresult = (e: any) => {
      error.value = null
      const transcript: string = Array.from(e.results as SpeechRecognitionResultList)
        .map((r: SpeechRecognitionResult) => r[0].transcript)
        .join(' ')
        .trim()
      if (transcript) onResult(transcript)
    }

    recognition.onend = () => {
      isListening.value = false
    }

    recognition.onerror = (e: any) => {
      isListening.value = false
      if (e.error === 'not-allowed') {
        error.value = 'ไม่ได้รับอนุญาตใช้ไมค์ — กรุณาอนุญาตใน browser แล้วลองใหม่'
      } else if (e.error === 'no-speech') {
        error.value = 'ไม่ได้ยินเสียง — ลองพูดอีกครั้ง'
      } else if (e.error !== 'aborted') {
        error.value = `เกิดข้อผิดพลาด: ${e.error}`
      }
      // Re-create instance so it can be reused after error
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const SR = (window as any).SpeechRecognition ?? (window as any).webkitSpeechRecognition
      if (SR) initRecognition(SR)
    }
  }

  onUnmounted(() => {
    try { recognition?.abort() } catch { /* ignore */ }
  })

  function toggle() {
    error.value = null
    if (!recognition) return
    if (isListening.value) {
      recognition.abort()
      isListening.value = false
    } else {
      isListening.value = true
      try {
        recognition.start()
      } catch (e: any) {
        isListening.value = false
        error.value = `ไม่สามารถเริ่มฟังได้: ${e?.message ?? e}`
      }
    }
  }

  return { isListening, isSupported, error, toggle }
}
