import type { KeyResponse as BoardKey } from "@/api/models";
import SoundBoard from "@/components/SoundBoard";

const keys: BoardKey[] = [
  {
    id: '1',
    boardId: '1',
    name: 'Whoop',
    description: 'Whoooop',
    hotKey: 'W',
    imageUrl: 'https://images.unsplash.com/photo-1582053433976-25c00369fc93?ixid=MXwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHw%3D&ixlib=rb-1.2.1&auto=format&fit=crop&w=512&q=80',
    audioUrl: 'https://www.soundjay.com/Human/sounds/human-laugh-1.wav',
    audioMediaId: '1',
    imageMediaId: '1',
    createdAt: new Date('2021-01-01'),
    updatedAt: new Date('2021-01-01'),
  },
  {
    id: '2',
    boardId: '1',
    name: 'Bloop',
    description: 'Blooooop',
    hotKey: 'B',
    audioUrl: 'https://www.soundjay.com/Human/sounds/human-laugh-1.wav',
    imageUrl:
      'https://images.unsplash.com/photo-1582053433976-25c00369fc93?ixid=MXwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHw%3D&ixlib=rb-1.2.1&auto=format&fit=crop&w=512&q=80',
    audioMediaId: '2',
    imageMediaId: '2',
    createdAt: new Date('2021-01-01'),
    updatedAt: new Date('2021-01-01'),
  },
  {
    id: '3',
    boardId: '1',
    name: 'Oop',
    description: 'Oooooop',
    hotKey: 'O',
    audioUrl: 'https://www.soundjay.com/Human/sounds/human-laugh-1.wav',
    imageUrl:
      'https://images.unsplash.com/photo-1582053433976-25c00369fc93?ixid=MXwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHw%3D&ixlib=rb-1.2.1&auto=format&fit=crop&w=512&q=80',
    audioMediaId: '3',
    imageMediaId: '3',
    createdAt: new Date('2021-01-01'),
    updatedAt: new Date('2021-01-01'),
  }
];

export default function Demo() {
  return (
    <div id="demo">
      <SoundBoard keys={keys} />
    </div>
  );
}
