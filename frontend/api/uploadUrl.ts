const rawBaseUrl = process.env.NEXT_PUBLIC_API_URL;
const BASE_URL =
  rawBaseUrl && rawBaseUrl.trim() !== ''
    ? rawBaseUrl
    : 'https://api.tuneslap.com/api/v1';

export interface GenerateUploadUrlRequest {
  fileName: string;
  contentType: string;
  fileSize: number;
}

export interface GenerateUploadUrlResponse {
  success: boolean;
  data?: {
    signedUrl: string;
    objectKey: string;
    bucketName: string;
    fileUrl: string; // The actual file URL to use after upload
  };
  error?: string;
}

export async function generateUploadUrl(
  token: string,
  data: GenerateUploadUrlRequest
): Promise<GenerateUploadUrlResponse> {
  try {
    const response = await fetch(`${BASE_URL}/media/upload-url`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });

    const result = await response.json();

    if (response.ok) {
      return {
        success: true,
        data: result.data,
      };
    } else {
      return {
        success: false,
        error: result.message || 'Failed to generate upload URL',
      };
    }
  } catch (error) {
    return {
      success: false,
      error:
        error instanceof Error ? error.message : 'Failed to generate upload URL',
    };
  }
}
