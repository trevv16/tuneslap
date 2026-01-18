import { useEffect, useRef } from 'react';
import { toast } from 'react-hot-toast';

export function useAudio(audioUrl: string, hotKey?: string, onError?: (error: string) => void) {
  const audioCtxRef = useRef<AudioContext | null>(null);
  const bufferRef = useRef<AudioBuffer | null>(null);
  const sourceRef = useRef<AudioBufferSourceNode | null>(null);
  const isKeyPressedRef = useRef<boolean>(false);

  // Load audio buffer
  useEffect(() => {
    if (!audioUrl) {
      toast.error(`Missing audio url for key: ${hotKey}`);
      return;
    }

    const ctx = new AudioContext();
    audioCtxRef.current = ctx;

    fetch(audioUrl)
      .then(res => res.arrayBuffer())
      .then(data => ctx.decodeAudioData(data))
      .then(decoded => {
        bufferRef.current = decoded;
      })
      .catch((error) => {
        onError?.(error.message);
        // toast.error(`Error loading audio: ${hotKey}`);
      });

    return () => {
      ctx.close();
    };
  }, [audioUrl]);

  const play = () => {
    if (!audioCtxRef.current || !bufferRef.current) return;

    // Stop any currently playing audio first
    stop();

    const source = audioCtxRef.current.createBufferSource();
    source.buffer = bufferRef.current;
    source.connect(audioCtxRef.current.destination);
    source.start(0);
    sourceRef.current = source;
  };

  const stop = () => {
    sourceRef.current?.stop();
    sourceRef.current?.disconnect();
    sourceRef.current = null;
  };

  // Optional keyboard binding
  useEffect(() => {
    if (!hotKey) return;

    const handleDown = (e: KeyboardEvent) => {
      if (e.key.toLowerCase() === hotKey.toLowerCase() && !isKeyPressedRef.current) {
        isKeyPressedRef.current = true;
        play();
      }
    };

    const handleUp = (e: KeyboardEvent) => {
      if (e.key.toLowerCase() === hotKey.toLowerCase()) {
        isKeyPressedRef.current = false;
        stop();
      }
    };

    window.addEventListener('keydown', handleDown);
    window.addEventListener('keyup', handleUp);
    return () => {
      window.removeEventListener('keydown', handleDown);
      window.removeEventListener('keyup', handleUp);
    };
  }, [hotKey]);

  // Bind to any DOM element
  const bindPressHandlers = (element: HTMLElement | null) => {
    if (!element) return;

    let isPressed = false;

    element.onmousedown = () => {
      if (!isPressed) {
        isPressed = true;
        play();
      }
    };

    element.onmouseup = () => {
      isPressed = false;
      stop();
    };

    element.onmouseleave = () => {
      isPressed = false;
      stop();
    };

    element.ontouchstart = () => {
      if (!isPressed) {
        isPressed = true;
        play();
      }
    };

    element.ontouchend = () => {
      isPressed = false;
      stop();
    };
  };

  return {
    play,
    stop,
    bindPressHandlers,
  };
}
