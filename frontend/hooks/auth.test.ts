import { renderHook, waitFor } from '@testing-library/react'
import { useSignUp, useSignIn, useForgotPassword, useResetPassword } from './auth'
import { createWrapper } from '@/__mocks__/providers/allProviders'
import { mockSigninResponse, mockSignupResponse } from '@/__mocks__/data/fixtures'

// Mock the API and config modules
jest.mock('@/api', () => ({
  AuthApi: jest.fn(),
}))

jest.mock('@/api/config', () => ({
  getApiConfig: jest.fn(() => ({})),
}))

jest.mock('@/utils/token', () => ({
  setStoredToken: jest.fn(),
  getStoredToken: jest.fn(),
}))

import { AuthApi } from '@/api'
import { setStoredToken } from '@/utils/token'

describe('useSignUp', () => {
  const mockSignup = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockAuthApi = AuthApi as jest.Mock
    MockAuthApi.mockImplementation(() => ({
      signup: mockSignup,
    }))
  })

  it('should call signup API with request data', async () => {
    mockSignup.mockResolvedValue(mockSignupResponse)

    const { result } = renderHook(() => useSignUp(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      email: 'test@example.com',
      password: 'password123',
      name: 'Test User',
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(mockSignup).toHaveBeenCalledWith({
      signupRequest: {
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User',
      },
    })
  })

  it('should handle signup error', async () => {
    const error = new Error('Email already exists')
    mockSignup.mockRejectedValue(error)

    const { result } = renderHook(() => useSignUp(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      email: 'existing@example.com',
      password: 'password123',
      name: 'Test User',
    })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(result.current.error?.message).toBe('Email already exists')
  })
})

describe('useSignIn', () => {
  const mockSignin = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockAuthApi = AuthApi as jest.Mock
    MockAuthApi.mockImplementation(() => ({
      signin: mockSignin,
    }))
  })

  it('should call signin API and store token on success', async () => {
    mockSignin.mockResolvedValue(mockSigninResponse)

    const { result } = renderHook(() => useSignIn(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      email: 'test@example.com',
      password: 'password123',
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(mockSignin).toHaveBeenCalledWith({
      signinRequest: {
        email: 'test@example.com',
        password: 'password123',
      },
    })
    expect(setStoredToken).toHaveBeenCalledWith('mock-jwt-token')
  })

  it('should not store token when signin fails', async () => {
    const error = new Error('Invalid credentials')
    mockSignin.mockRejectedValue(error)

    const { result } = renderHook(() => useSignIn(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      email: 'test@example.com',
      password: 'wrong-password',
    })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(setStoredToken).not.toHaveBeenCalled()
  })

  it('should not store token when response has no token', async () => {
    mockSignin.mockResolvedValue({
      success: false,
      data: {},
    })

    const { result } = renderHook(() => useSignIn(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      email: 'test@example.com',
      password: 'password123',
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })
    expect(setStoredToken).not.toHaveBeenCalled()
  })
})

describe('useForgotPassword', () => {
  const mockForgot = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockAuthApi = AuthApi as jest.Mock
    MockAuthApi.mockImplementation(() => ({
      forgot: mockForgot,
    }))
  })

  it('should call forgot password API', async () => {
    mockForgot.mockResolvedValue({
      success: true,
      message: 'Reset email sent',
    })

    const { result } = renderHook(() => useForgotPassword(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      email: 'test@example.com',
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(mockForgot).toHaveBeenCalledWith({
      forgotRequest: {
        email: 'test@example.com',
      },
    })
  })

  it('should handle forgot password error', async () => {
    const error = new Error('User not found')
    mockForgot.mockRejectedValue(error)

    const { result } = renderHook(() => useForgotPassword(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      email: 'nonexistent@example.com',
    })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(result.current.error?.message).toBe('User not found')
  })
})

describe('useResetPassword', () => {
  const mockReset = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
    const MockAuthApi = AuthApi as jest.Mock
    MockAuthApi.mockImplementation(() => ({
      reset: mockReset,
    }))
  })

  it('should call reset password API', async () => {
    mockReset.mockResolvedValue({
      success: true,
      message: 'Password reset successful',
    })

    const { result } = renderHook(() => useResetPassword(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      token: 'reset-token',
      password: 'new-password',
    })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(mockReset).toHaveBeenCalledWith({
      resetRequest: {
        token: 'reset-token',
        password: 'new-password',
      },
    })
  })

  it('should handle reset password error', async () => {
    const error = new Error('Token expired')
    mockReset.mockRejectedValue(error)

    const { result } = renderHook(() => useResetPassword(), {
      wrapper: createWrapper(),
    })

    result.current.mutate({
      token: 'expired-token',
      password: 'new-password',
    })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })
    expect(result.current.error?.message).toBe('Token expired')
  })
})
