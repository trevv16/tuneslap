import { generateUploadUrl, GenerateUploadUrlRequest } from './uploadUrl'

describe('generateUploadUrl', () => {
  const mockToken = 'test-jwt-token'
  const mockRequest: GenerateUploadUrlRequest = {
    fileName: 'test-audio.mp3',
    contentType: 'audio/mpeg',
    fileSize: 1024 * 1024,
  }

  beforeEach(() => {
    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockReset()
  })

  it('should return success response with upload data', async () => {
    const mockResponseData = {
      data: {
        signedUrl: 'https://storage.example.com/upload?signature=abc123',
        objectKey: 'uploads/test-audio.mp3',
        bucketName: 'tuneslap-media',
        fileUrl: 'https://storage.example.com/uploads/test-audio.mp3',
      },
    }

    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockResponseData),
    })

    const result = await generateUploadUrl(mockToken, mockRequest)

    expect(result.success).toBe(true)
    expect(result.data).toEqual(mockResponseData.data)
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/media/upload-url'),
      expect.objectContaining({
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${mockToken}`,
        },
        body: JSON.stringify(mockRequest),
      })
    )
  })

  it('should return error response when API returns error', async () => {
    const mockErrorResponse = {
      message: 'File too large',
    }

    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockResolvedValue({
      ok: false,
      json: () => Promise.resolve(mockErrorResponse),
    })

    const result = await generateUploadUrl(mockToken, mockRequest)

    expect(result.success).toBe(false)
    expect(result.error).toBe('File too large')
    expect(result.data).toBeUndefined()
  })

  it('should return default error message when API error has no message', async () => {
    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({}),
    })

    const result = await generateUploadUrl(mockToken, mockRequest)

    expect(result.success).toBe(false)
    expect(result.error).toBe('Failed to generate upload URL')
  })

  it('should handle network errors gracefully', async () => {
    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockImplementationOnce(() => Promise.reject(new Error('Network request failed')))

    const result = await generateUploadUrl(mockToken, mockRequest)

    expect(result.success).toBe(false)
    expect(result.error).toBe('Network request failed')
  })

  it('should handle JSON parse errors', async () => {
    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.reject(new Error('Invalid JSON')),
    })

    const result = await generateUploadUrl(mockToken, mockRequest)

    expect(result.success).toBe(false)
    expect(result.error).toBe('Invalid JSON')
  })

  it('should handle non-Error exceptions', async () => {
    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.reject(new Error('Something went wrong')),
    })

    const result = await generateUploadUrl(mockToken, mockRequest)

    expect(result.success).toBe(false)
    expect(result.error).toBe('Something went wrong')
  })

  it('should use correct base URL from environment', async () => {
    const originalEnv = process.env.NEXT_PUBLIC_API_URL
    process.env.NEXT_PUBLIC_API_URL = 'https://api.custom.com/api/v1'

    // Need to reimport to pick up new env
    jest.resetModules()
    const { generateUploadUrl: freshGenerateUploadUrl } = await import('./uploadUrl')

    const fetchMock = global.fetch as jest.Mock
    fetchMock.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: {} }),
    })

    await freshGenerateUploadUrl(mockToken, mockRequest)

    expect(global.fetch).toHaveBeenCalledWith(
      'https://api.custom.com/api/v1/media/upload-url',
      expect.any(Object)
    )

    process.env.NEXT_PUBLIC_API_URL = originalEnv
  })
})
