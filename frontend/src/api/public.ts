/**
 * Public API endpoints (no authentication required)
 * Handles public-facing queries like usage statistics by API key
 */

import axios from 'axios'

// Create a separate axios instance for public API that doesn't require auth
const publicClient = axios.create({
  baseURL: '/api/v1/public',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Response interceptor to extract data from the envelope
publicClient.interceptors.response.use(
  (response) => {
    // Extract data from the envelope format { code: 0, data: {...} }
    if (response.data && response.data.code === 0) {
      response.data = response.data.data
    }
    return response
  },
  (error) => {
    return Promise.reject(error)
  }
)

// ==================== Types ====================

export interface PublicUsageStatsParams {
  key: string
  period?: 'today' | 'week' | 'month'
  start_date?: string
  end_date?: string
  timezone?: string
}

export interface PublicUsageStatsResponse {
  total_requests: number
  total_input_tokens: number
  total_output_tokens: number
  total_cache_creation_tokens: number
  total_cache_read_tokens: number
  total_tokens: number
  total_cost: number
  total_actual_cost: number
  average_duration_ms: number
}

// ==================== API Functions ====================

/**
 * Get usage statistics by API key
 * @param params - Query parameters including the API key
 * @returns Usage statistics for the specified API key
 */
export async function getUsageByKey(
  params: PublicUsageStatsParams
): Promise<PublicUsageStatsResponse> {
  const { data } = await publicClient.get<PublicUsageStatsResponse>('/usage', {
    params
  })
  return data
}

/**
 * Get usage statistics for a specific period
 * @param key - API key
 * @param period - Time period ('today', 'week', 'month')
 * @param timezone - Optional timezone
 * @returns Usage statistics
 */
export async function getUsageByPeriod(
  key: string,
  period: 'today' | 'week' | 'month' = 'today',
  timezone?: string
): Promise<PublicUsageStatsResponse> {
  return getUsageByKey({ key, period, timezone })
}

/**
 * Get usage statistics for a date range
 * @param key - API key
 * @param startDate - Start date (YYYY-MM-DD format)
 * @param endDate - End date (YYYY-MM-DD format)
 * @param timezone - Optional timezone
 * @returns Usage statistics
 */
export async function getUsageByDateRange(
  key: string,
  startDate: string,
  endDate: string,
  timezone?: string
): Promise<PublicUsageStatsResponse> {
  return getUsageByKey({
    key,
    start_date: startDate,
    end_date: endDate,
    timezone
  })
}

export const publicAPI = {
  getUsageByKey,
  getUsageByPeriod,
  getUsageByDateRange
}

export default publicAPI
