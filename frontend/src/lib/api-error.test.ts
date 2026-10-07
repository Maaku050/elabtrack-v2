import { describe, expect, it } from 'vitest'
import { AxiosError, type AxiosAdapter, type AxiosResponse } from 'axios'
import { ApiClient } from './api-client'
import { ApiRequestError, apiErrorMessage, normalizeApiError, normalizeApiResponseError } from './api-error'
import { queryClient } from '@/app/query-client'

const id = '12345678-1234-4123-8123-123456789abc'
const otherID = 'abcdefab-1234-4123-8123-123456789abc'
function payload(code: string, message = 'Safe public message.', fields?: unknown) {
  return { success: false, message, data: null, meta: null, error: { code, message, requestId: id, fields } }
}

describe('standard error contract', () => {
  it.each([[400, 'VALIDATION_ERROR'], [401, 'UNAUTHORIZED'], [403, 'FORBIDDEN'], [404, 'NOT_FOUND'], [405, 'METHOD_NOT_ALLOWED'], [409, 'CONFLICT'], [429, 'RATE_LIMITED'], [500, 'INTERNAL_ERROR'], [503, 'SERVICE_UNAVAILABLE']] as const)('normalizes %s/%s through the real transport', async (status, code) => {
    const adapter: AxiosAdapter = async (config) => {
      const response: AxiosResponse = { config, status, statusText: String(status), headers: { 'x-request-id': id, 'retry-after': '60' }, data: payload(code) }
      throw new AxiosError('SQL password internal sentinel', 'ERR_BAD_RESPONSE', config, undefined, response)
    }
    const client = new ApiClient(adapter)
    const error = await client.get('/foundation', { attachAuth: false }).catch((value: unknown) => value)
    expect(error).toBeInstanceOf(ApiRequestError)
    expect(error).toMatchObject({ status, code, requestId: id })
    expect((error as ApiRequestError).message).not.toMatch(/SQL|password|sentinel/)
    expect(error).not.toHaveProperty('config')
    expect(error).not.toHaveProperty('response')
  })
  it('supports field mapping without parsing messages and strips unsafe structure', () => {
    const error = normalizeApiResponseError(400, payload('VALIDATION_ERROR', 'Validation failed.', { email: 'is required', name: 'is invalid', bad: { sql: 'private' }, 'bad\nfield': 'private' }))
    expect(error.fields).toEqual({ email: 'is required', name: 'is invalid' })
    expect(apiErrorMessage(error)).toBe('Validation failed.')
    expect(apiErrorMessage(error)).not.toContain(id)
  })
  it('captures a safe header reference and prefers it over a disagreeing body ID', () => {
    expect(normalizeApiResponseError(403, payload('FORBIDDEN'), otherID).requestId).toBe(otherID)
    expect(normalizeApiResponseError(403, payload('FORBIDDEN'), 'credential\nsecret').requestId).toBe(id)
    expect(normalizeApiResponseError(500, 'private proxy error', id).requestId).toBe(id)
  })
  it.each(['private SQL password detail', { message: 'private SQL', error: { code: 'INTERNAL_GO_TYPE' } }, { success: false, error: { code: 'FORBIDDEN', message: 'private', requestId: 'malformed' } }, null, []])('rejects arbitrary/legacy error payloads safely', (body) => {
    const error = normalizeApiResponseError(403, body)
    expect(error.message).toBe('Request failed. Please try again.')
    expect(error.code).toBe('REQUEST_FAILED')
    expect(error.requestId).toBeUndefined()
  })
  it('never displays server/proxy internal messages on 500, and provides an optional support reference', () => {
    const error = normalizeApiResponseError(500, payload('INTERNAL_ERROR', 'SQL JWT password sentinel'))
    expect(apiErrorMessage(error)).toBe(`Something went wrong. Please try again later. Reference: ${id}`)
    expect(error.message).not.toContain('SQL')
  })
  it('normalizes network, unexpected and already-normalized failures without retaining internals', () => {
    const network = normalizeApiError(new AxiosError('Authorization Cookie sentinel', 'ERR_NETWORK'))
    expect(network).toMatchObject({ status: 0, code: 'NETWORK_ERROR' })
    expect(network.message).not.toContain('sentinel')
    expect(normalizeApiError(new Error('private sentinel')).message).not.toContain('sentinel')
    expect(normalizeApiError(network)).toBe(network)
  })
  it('does not retry client failures via Query; server/network queries retain one bounded retry', () => {
    const retry = queryClient.getDefaultOptions().queries?.retry
    expect(typeof retry).toBe('function')
    if (typeof retry !== 'function') throw new Error('Missing retry policy')
    for (const status of [400, 401, 403, 404, 405, 409, 429]) expect(retry(0, new ApiRequestError('safe', status))).toBe(false)
    expect(retry(0, new ApiRequestError('safe', 500))).toBe(true)
    expect(retry(1, new ApiRequestError('safe', 500))).toBe(false)
  })
})
