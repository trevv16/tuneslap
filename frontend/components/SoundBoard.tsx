import type { KeyResponse as BoardKey } from "@/api/models";
import SoundKey from "./SoundKey";

type SoundBoardProps = {
  keys: BoardKey[];
}

export default function SoundBoard({ keys }: SoundBoardProps) {
  return (
    <div className="flex justify-center bg-elevated p-4 rounded-lg h-dvh">
      <ul role="list" className="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 sm:gap-x-6 lg:grid-cols-4 xl:gap-x-8">
        {keys.map((key) => (
          <li key={key.id} className="relative">
            <SoundKey boardKey={key} />
          </li>
        ))}
      </ul>
    </div>
  )
}