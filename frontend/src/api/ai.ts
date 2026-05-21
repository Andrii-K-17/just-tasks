import type { AIGenerateResult } from '@/types'

/**
 * Sends a prompt to the backend AI endpoint and returns a structured task plan.
 */
export const generateTasks = async (text: string): Promise<AIGenerateResult> => {
  const response = await fetch('/api/ai/generate', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ text }),
  })
  const data = await response.json()
  if (!response.ok) throw new Error(data.error ?? 'AI request failed')
  return data as AIGenerateResult
}
