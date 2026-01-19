'use client'

import {
  MediaProcessingParamsImageFormatEnum,
  MediaProcessingParamsImageAspectRatioEnum,
  MediaProcessingParamsImageApplyFiltersEnum,
} from '@/api/models'
import { UseFormReturn } from 'react-hook-form'
import { ImageProcessingFormData } from '@/hooks/media'

type ImageProcessingFormProps = {
  form: UseFormReturn<ImageProcessingFormData>
}

const FORMAT_OPTIONS = [
  { value: MediaProcessingParamsImageFormatEnum.Webp, label: 'WebP (Recommended)' },
  { value: MediaProcessingParamsImageFormatEnum.Jpeg, label: 'JPEG' },
  { value: MediaProcessingParamsImageFormatEnum.Png, label: 'PNG' },
  { value: MediaProcessingParamsImageFormatEnum.Gif, label: 'GIF' },
  { value: MediaProcessingParamsImageFormatEnum.Svg, label: 'SVG' },
]

const ASPECT_RATIO_OPTIONS = [
  { value: '', label: 'Original' },
  { value: MediaProcessingParamsImageAspectRatioEnum._11, label: '1:1 (Square)' },
  { value: MediaProcessingParamsImageAspectRatioEnum._43, label: '4:3' },
  { value: MediaProcessingParamsImageAspectRatioEnum._169, label: '16:9 (Widescreen)' },
  { value: MediaProcessingParamsImageAspectRatioEnum._1610, label: '16:10' },
  { value: MediaProcessingParamsImageAspectRatioEnum._32, label: '3:2' },
  { value: MediaProcessingParamsImageAspectRatioEnum._54, label: '5:4' },
]

const FILTER_OPTIONS = [
  { value: MediaProcessingParamsImageApplyFiltersEnum.Grayscale, label: 'Grayscale' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Sepia, label: 'Sepia' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Blur, label: 'Blur' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Sharpen, label: 'Sharpen' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Brightness, label: 'Brightness' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Contrast, label: 'Contrast' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Saturation, label: 'Saturation' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Hue, label: 'Hue' },
  { value: MediaProcessingParamsImageApplyFiltersEnum.Invert, label: 'Invert' },
]

export default function ImageProcessingForm({ form }: ImageProcessingFormProps) {
  const { register, watch, setValue, formState: { errors } } = form

  const selectedFilters = watch('applyFilters') || []

  const handleFilterToggle = (filterValue: MediaProcessingParamsImageApplyFiltersEnum) => {
    const current = selectedFilters
    if (current.includes(filterValue)) {
      setValue('applyFilters', current.filter(f => f !== filterValue))
    } else {
      setValue('applyFilters', [...current, filterValue])
    }
  }

  return (
    <div className="space-y-6">
      {/* Output Format */}
      <div>
        <label htmlFor="format" className="block text-sm font-medium text-highlight mb-2">
          Output Format
        </label>
        <select
          {...register('format')}
          className="block w-full rounded-md bg-elevated h-10 px-3 text-base border border-muted focus:border-accent sm:text-sm"
        >
          {FORMAT_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        {errors.format && (
          <p className="mt-1 text-sm text-error">{errors.format.message}</p>
        )}
      </div>

      {/* Aspect Ratio */}
      <div>
        <label htmlFor="aspectRatio" className="block text-sm font-medium text-highlight mb-2">
          Aspect Ratio
        </label>
        <select
          {...register('aspectRatio')}
          className="block w-full rounded-md bg-elevated h-10 px-3 text-base border border-muted focus:border-accent sm:text-sm"
        >
          {ASPECT_RATIO_OPTIONS.map((option) => (
            <option key={option.value || 'original'} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        {errors.aspectRatio && (
          <p className="mt-1 text-sm text-error">{errors.aspectRatio.message}</p>
        )}
      </div>

      {/* Resize Dimensions */}
      <div>
        <label className="block text-sm font-medium text-highlight mb-2">
          Resize To (optional)
        </label>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <input
              type="number"
              {...register('resizeWidth', { valueAsNumber: true })}
              min={1}
              max={4096}
              placeholder="Width"
              className="block w-full rounded-md bg-elevated h-10 px-3 text-base border border-muted focus:border-accent sm:text-sm"
            />
            {errors.resizeWidth && (
              <p className="mt-1 text-sm text-error">{errors.resizeWidth.message}</p>
            )}
          </div>
          <div>
            <input
              type="number"
              {...register('resizeHeight', { valueAsNumber: true })}
              min={1}
              max={4096}
              placeholder="Height"
              className="block w-full rounded-md bg-elevated h-10 px-3 text-base border border-muted focus:border-accent sm:text-sm"
            />
            {errors.resizeHeight && (
              <p className="mt-1 text-sm text-error">{errors.resizeHeight.message}</p>
            )}
          </div>
        </div>
        <p className="mt-1 text-xs text-muted">Leave empty to keep original dimensions</p>
      </div>

      {/* Filters */}
      <div>
        <label className="block text-sm font-medium text-highlight mb-2">
          Apply Filters
        </label>
        <div className="grid grid-cols-3 gap-2">
          {FILTER_OPTIONS.map((filter) => (
            <label
              key={filter.value}
              className={`flex items-center justify-center px-3 py-2 rounded-md text-sm cursor-pointer transition-colors ${
                selectedFilters.includes(filter.value)
                  ? 'bg-accent text-inverted'
                  : 'bg-elevated text-base hover:bg-muted'
              }`}
            >
              <input
                type="checkbox"
                checked={selectedFilters.includes(filter.value)}
                onChange={() => handleFilterToggle(filter.value)}
                className="sr-only"
              />
              {filter.label}
            </label>
          ))}
        </div>
        {errors.applyFilters && (
          <p className="mt-1 text-sm text-error">{errors.applyFilters.message}</p>
        )}
      </div>
    </div>
  )
}
