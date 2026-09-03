import { apiClient } from '@/api/client'

export interface ConversationRecord {
  id: number; request_id: string; user_id: number; username?: string; user_email?: string
  api_key_name?: string; provider: string; endpoint: string; requested_model: string; upstream_model: string
  prompt_content: string; response_content: string; stream: boolean; status: string; complete: boolean
  truncated: boolean; input_tokens: number; output_tokens: number; duration_ms: number; created_at: string
}
export interface ConversationPage { items: ConversationRecord[]; total: number; page: number; page_size: number }
export interface ConversationConfig { enabled: boolean; capture_all_groups: boolean; max_prompt_bytes: number; max_response_bytes: number; manual_delete_enabled: boolean }

export async function listConversations(params: Record<string, unknown>): Promise<ConversationPage> {
  const { data } = await apiClient.get<ConversationPage>('/admin/conversations', { params }); return data
}
export async function getConversation(id: number): Promise<ConversationRecord> {
  const { data } = await apiClient.get<ConversationRecord>(`/admin/conversations/${id}`); return data
}
export async function deleteConversation(id: number): Promise<void> { await apiClient.delete(`/admin/conversations/${id}`) }
export async function getConversationConfig(): Promise<ConversationConfig> { const { data } = await apiClient.get<ConversationConfig>('/admin/conversations/config'); return data }
export async function updateConversationConfig(config: ConversationConfig): Promise<ConversationConfig> { const { data } = await apiClient.put<ConversationConfig>('/admin/conversations/config', config); return data }
