#!/bin/bash
set -e

# Define target directory
TARGET_DIR="server/assets/demo"

# Create directory if it doesn't exist
mkdir -p "$TARGET_DIR"

# Function to download a file
download_file() {
  local filename=$1
  local url=$2
  local output_path="$TARGET_DIR/$filename"
  
  # Skip if file already exists
  if [ -f "$output_path" ]; then
    echo "⏭️  Skipping $filename (already exists)"
    return 0
  fi
  
  echo "⬇️  Downloading $filename..."
  curl -L -f -o "$output_path" "$url" || {
    echo "❌ Failed to download $filename"
    exit 1
  }
  
  echo "✅ Downloaded $filename"
}

# Download each file
download_file "applause.mp3" "https://cdn.freesound.org/previews/100/100904_1163166-lq.mp3"
download_file "drum-roll.mp3" "https://cdn.freesound.org/previews/569/569113_12672827-lq.mp3"
download_file "laughter.mp3" "https://cdn.freesound.org/previews/33/33658_43834-lq.mp3"
download_file "air-horn.mp3" "https://cdn.freesound.org/previews/110/110662_1927446-lq.mp3"
download_file "whoosh.mp3" "https://cdn.freesound.org/previews/60/60009_71257-lq.mp3"
download_file "bell-ding.mp3" "https://cdn.freesound.org/previews/411/411089_5121236-lq.mp3"
download_file "boing.mp3" "https://cdn.freesound.org/previews/131/131660_2398403-lq.mp3"
download_file "tada.mp3" "https://cdn.freesound.org/previews/397/397355_4284968-lq.mp3"

echo "✅ All demo audio files downloaded successfully!"
