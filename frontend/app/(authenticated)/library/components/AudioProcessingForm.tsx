'use client'

import {
  MediaProcessingParamsAudioContentTypeEnum,
} from '@/api/models'
import { UseFormReturn } from 'react-hook-form'
import { AudioProcessingFormData } from '@/hooks/media'

type AudioProcessingFormProps = {
  form: UseFormReturn<AudioProcessingFormData>
}

const CONTENT_TYPE_OPTIONS = [
  { value: MediaProcessingParamsAudioContentTypeEnum.AudioWebm, label: 'WebM (Recommended)' },
  { value: MediaProcessingParamsAudioContentTypeEnum.AudioMp3, label: 'MP3' },
  { value: MediaProcessingParamsAudioContentTypeEnum.AudioWav, label: 'WAV' },
  { value: MediaProcessingParamsAudioContentTypeEnum.AudioOgg, label: 'OGG' },
  { value: MediaProcessingParamsAudioContentTypeEnum.AudioAac, label: 'AAC' },
]

export default function AudioProcessingForm({ form }: AudioProcessingFormProps) {
  const { register, watch, formState: { errors } } = form

  const speed = watch('speed') ?? 1
  const pitch = watch('pitch') ?? 0
  const fadeIn = watch('fadeIn') ?? 0
  const fadeOut = watch('fadeOut') ?? 0

  return (
    <div className="space-y-6">
      {/* Output Format */}
      <div>
        <label htmlFor="contentType" className="block text-sm font-medium text-highlight mb-2">
          Output Format
        </label>
        <select
          {...register('contentType')}
          className="block w-full rounded-md bg-elevated h-10 px-3 text-base border border-muted focus:border-accent sm:text-sm"
        >
          {CONTENT_TYPE_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        {errors.contentType && (
          <p className="mt-1 text-sm text-error">{errors.contentType.message}</p>
        )}
      </div>

      {/* Normalize Toggle */}
      <div className="flex items-center justify-between">
        <div>
          <label htmlFor="normalize" className="text-sm font-medium text-highlight">
            Normalize Audio
          </label>
          <p className="text-xs text-muted">Automatically adjust volume levels</p>
        </div>
        <input
          type="checkbox"
          {...register('normalize')}
          className="h-5 w-5 rounded border-muted text-accent focus:ring-accent"
        />
      </div>

      {/* Speed Slider */}
      <div>
        <div className="flex items-center justify-between mb-2">
          <label htmlFor="speed" className="text-sm font-medium text-highlight">
            Speed
          </label>
          <span className="text-sm text-muted">{speed.toFixed(2)}x</span>
        </div>
        <input
          type="range"
          {...register('speed', { valueAsNumber: true })}
          min={0.5}
          max={2}
          step={0.05}
          className="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer accent-accent"
        />
        <div className="flex justify-between text-xs text-muted mt-1">
          <span>0.5x</span>
          <span>1x</span>
          <span>2x</span>
        </div>
        {errors.speed && (
          <p className="mt-1 text-sm text-error">{errors.speed.message}</p>
        )}
      </div>

      {/* Pitch Slider */}
      <div>
        <div className="flex items-center justify-between mb-2">
          <label htmlFor="pitch" className="text-sm font-medium text-highlight">
            Pitch
          </label>
          <span className="text-sm text-muted">{pitch > 0 ? '+' : ''}{pitch} semitones</span>
        </div>
        <input
          type="range"
          {...register('pitch', { valueAsNumber: true })}
          min={-12}
          max={12}
          step={1}
          className="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer accent-accent"
        />
        <div className="flex justify-between text-xs text-muted mt-1">
          <span>-12</span>
          <span>0</span>
          <span>+12</span>
        </div>
        {errors.pitch && (
          <p className="mt-1 text-sm text-error">{errors.pitch.message}</p>
        )}
      </div>

      {/* Fade In Slider */}
      <div>
        <div className="flex items-center justify-between mb-2">
          <label htmlFor="fadeIn" className="text-sm font-medium text-highlight">
            Fade In
          </label>
          <span className="text-sm text-muted">{fadeIn.toFixed(1)}s</span>
        </div>
        <input
          type="range"
          {...register('fadeIn', { valueAsNumber: true })}
          min={0}
          max={10}
          step={0.1}
          className="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer accent-accent"
        />
        <div className="flex justify-between text-xs text-muted mt-1">
          <span>0s</span>
          <span>5s</span>
          <span>10s</span>
        </div>
        {errors.fadeIn && (
          <p className="mt-1 text-sm text-error">{errors.fadeIn.message}</p>
        )}
      </div>

      {/* Fade Out Slider */}
      <div>
        <div className="flex items-center justify-between mb-2">
          <label htmlFor="fadeOut" className="text-sm font-medium text-highlight">
            Fade Out
          </label>
          <span className="text-sm text-muted">{fadeOut.toFixed(1)}s</span>
        </div>
        <input
          type="range"
          {...register('fadeOut', { valueAsNumber: true })}
          min={0}
          max={10}
          step={0.1}
          className="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer accent-accent"
        />
        <div className="flex justify-between text-xs text-muted mt-1">
          <span>0s</span>
          <span>5s</span>
          <span>10s</span>
        </div>
        {errors.fadeOut && (
          <p className="mt-1 text-sm text-error">{errors.fadeOut.message}</p>
        )}
      </div>

      {/* Trim Start/End */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label htmlFor="trimStart" className="block text-sm font-medium text-highlight mb-2">
            Trim Start (seconds)
          </label>
          <input
            type="number"
            {...register('trimStart', { valueAsNumber: true })}
            min={0}
            step={0.1}
            placeholder="0"
            className="block w-full rounded-md bg-elevated h-10 px-3 text-base border border-muted focus:border-accent sm:text-sm"
          />
          {errors.trimStart && (
            <p className="mt-1 text-sm text-error">{errors.trimStart.message}</p>
          )}
        </div>
        <div>
          <label htmlFor="trimEnd" className="block text-sm font-medium text-highlight mb-2">
            Trim End (seconds)
          </label>
          <input
            type="number"
            {...register('trimEnd', { valueAsNumber: true })}
            min={0}
            step={0.1}
            placeholder="End of file"
            className="block w-full rounded-md bg-elevated h-10 px-3 text-base border border-muted focus:border-accent sm:text-sm"
          />
          {errors.trimEnd && (
            <p className="mt-1 text-sm text-error">{errors.trimEnd.message}</p>
          )}
        </div>
      </div>
    </div>
  )
}
