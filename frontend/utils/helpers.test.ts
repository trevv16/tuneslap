import { classNames, formatBytes } from './helpers'

describe('classNames', () => {
  it('should join multiple class names', () => {
    expect(classNames('foo', 'bar', 'baz')).toBe('foo bar baz')
  })

  it('should filter out falsy values', () => {
    expect(classNames('foo', false, 'bar', null, 'baz', undefined, '')).toBe('foo bar baz')
  })

  it('should return empty string for no arguments', () => {
    expect(classNames()).toBe('')
  })

  it('should return empty string for all falsy values', () => {
    expect(classNames(false, null, undefined, '')).toBe('')
  })

  it('should handle single class name', () => {
    expect(classNames('foo')).toBe('foo')
  })

  it('should handle conditional classes with variables', () => {
    const isActive = true
    const isDisabled = false
    // Use variables directly in the assertion to avoid lint warning about unnecessary conditionals
    const activeClass = isActive ? 'active' : ''
    const disabledClass = isDisabled ? 'disabled' : ''
    expect(classNames('base', activeClass, disabledClass)).toBe('base active')
  })

  it('should handle number 0 as falsy', () => {
    expect(classNames('foo', 0, 'bar')).toBe('foo bar')
  })
})

describe('formatBytes', () => {
  it('should return "Unlimited" for -1', () => {
    expect(formatBytes(-1)).toBe('Unlimited')
  })

  it('should return "0 MB" for 0 bytes', () => {
    expect(formatBytes(0)).toBe('0 MB')
  })

  it('should format bytes correctly', () => {
    expect(formatBytes(500)).toBe('500 Bytes')
  })

  it('should format kilobytes correctly', () => {
    expect(formatBytes(1024)).toBe('1 KB')
    expect(formatBytes(1536)).toBe('1.5 KB')
  })

  it('should format megabytes correctly', () => {
    expect(formatBytes(1024 * 1024)).toBe('1 MB')
    expect(formatBytes(1.5 * 1024 * 1024)).toBe('1.5 MB')
  })

  it('should format gigabytes correctly', () => {
    expect(formatBytes(1024 * 1024 * 1024)).toBe('1 GB')
    expect(formatBytes(2.5 * 1024 * 1024 * 1024)).toBe('2.5 GB')
  })

  it('should format terabytes correctly', () => {
    expect(formatBytes(1024 * 1024 * 1024 * 1024)).toBe('1 TB')
  })

  it('should cap at TB for very large values', () => {
    // 1 PB would be 1024 TB, but we cap at TB
    expect(formatBytes(1024 * 1024 * 1024 * 1024 * 1024)).toBe('1024 TB')
  })

  it('should handle decimal precision', () => {
    expect(formatBytes(1234567)).toBe('1.18 MB')
  })
})
