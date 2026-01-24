import { DemoApi } from '@/api/apis/DemoApi';
import { getServerApiConfig } from '@/api/config';
import type { KeyResponse } from "@/api/models";
import SoundBoard from "@/components/SoundBoard";
import { unstable_cache } from 'next/cache';

// Demo board ID - must match server/config/demo.go
const DEMO_BOARD_ID = '000000000000000000000002';

// Fallback demo keys in case the API is unavailable at build time
const FALLBACK_DATE = new Date('2024-01-01T00:00:00Z');
const fallbackKeys: KeyResponse[] = [
  {
    id: '000000000000000000000101',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000201',
    name: 'Applause',
    description: 'Demo sound: Applause',
    hotKey: '1',
    audioUrl: 'https://cdn.freesound.org/previews/100/100904_1163166-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
  {
    id: '000000000000000000000102',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000202',
    name: 'Drum Roll',
    description: 'Demo sound: Drum Roll',
    hotKey: '2',
    audioUrl: 'https://cdn.freesound.org/previews/569/569113_12672827-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1519892300165-cb5542fb47c7?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
  {
    id: '000000000000000000000103',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000203',
    name: 'Laughter',
    description: 'Demo sound: Laughter',
    hotKey: '3',
    audioUrl: 'https://cdn.freesound.org/previews/33/33658_43834-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1543610892-0b1f7e6d8ac1?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
  {
    id: '000000000000000000000104',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000204',
    name: 'Air Horn',
    description: 'Demo sound: Air Horn',
    hotKey: '4',
    audioUrl: 'https://cdn.freesound.org/previews/110/110662_1927446-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
  {
    id: '000000000000000000000105',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000205',
    name: 'Whoosh',
    description: 'Demo sound: Whoosh',
    hotKey: '5',
    audioUrl: 'https://cdn.freesound.org/previews/60/60009_71257-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
  {
    id: '000000000000000000000106',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000206',
    name: 'Bell Ding',
    description: 'Demo sound: Bell Ding',
    hotKey: '6',
    audioUrl: 'https://cdn.freesound.org/previews/411/411089_5121236-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1513836279014-a89f7a76ae86?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
  {
    id: '000000000000000000000107',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000207',
    name: 'Boing',
    description: 'Demo sound: Boing',
    hotKey: '7',
    audioUrl: 'https://cdn.freesound.org/previews/131/131660_2398403-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1518640467707-6811f4a6ab73?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
  {
    id: '000000000000000000000108',
    boardId: DEMO_BOARD_ID,
    audioMediaId: '000000000000000000000208',
    name: 'Ta-Da',
    description: 'Demo sound: Ta-Da',
    hotKey: '8',
    audioUrl: 'https://cdn.freesound.org/previews/397/397355_4284968-lq.mp3',
    imageUrl: 'https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=400',
    createdAt: FALLBACK_DATE,
    updatedAt: FALLBACK_DATE,
  },
];

// Cached function for ISR - revalidates every hour
const getDemoBoard = unstable_cache(
  async (): Promise<KeyResponse[]> => {
    try {
      const demoApi = new DemoApi(getServerApiConfig());
      const board = await demoApi.getDemoBoard();

      if (!board.keys || board.keys.length === 0) {
        console.warn('[Demo] Demo board has no keys, using fallback');
        return fallbackKeys;
      }

      return board.keys;
    } catch (error: unknown) {
      // Log more details about the error
      if (error instanceof Error) {
        console.warn('[Demo] Error fetching demo board:', error.message);
        if ('response' in error && error.response instanceof Response) {
          try {
            const errorText = await error.response.text();
            console.warn('[Demo] Server error response:', errorText);
          } catch {
            // Ignore if we can't read the response
          }
        }
      } else {
        console.warn('[Demo] Error fetching demo board, using fallback:', error);
      }
      return fallbackKeys;
    }
  },
  ['demo-board'],
  { revalidate: 3600 } // Revalidate every hour
);

export default async function Demo() {
  const keys = await getDemoBoard();

  return (
    <section id="demo" className="mx-auto max-w-7xl px-6 lg:px-8 mt-16 sm:mt-24">
      <div className="mx-auto max-w-2xl lg:text-center mb-12">
        <h2 className="text-base/7 font-semibold text-primary">Try It Out</h2>
        <p className="mt-2 text-3xl font-semibold tracking-tight text-highlight sm:text-4xl">
          Interactive Demo
        </p>
        <p className="mt-4 text-base text-muted-foreground">
          Click the keys below or press the number keys (1-8) on your keyboard to play sounds.
        </p>
      </div>
      <SoundBoard keys={keys} />
    </section>
  );
}
