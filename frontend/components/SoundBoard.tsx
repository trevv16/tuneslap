import type { KeyResponse as BoardKey } from "@/api/models";
import { Button } from "@/components/ui/button";
import { Keyboard, Plus } from "lucide-react";
import SoundKey from "./SoundKey";

type SoundBoardProps = {
  keys: BoardKey[];
  onAddKey?: () => void;
}

export default function SoundBoard({ keys, onAddKey }: SoundBoardProps) {
  if (keys.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center bg-card p-8 rounded-lg min-h-[50vh]">
        <Keyboard className="h-16 w-16 text-muted-foreground mb-4" />
        <h3 className="text-lg font-semibold text-foreground">No keys yet</h3>
        <p className="mt-1 text-sm text-muted-foreground text-center max-w-sm">
          Add keys to your soundboard to start triggering sounds with hotkeys or clicks.
        </p>
        {onAddKey && (
          <Button onClick={onAddKey} className="mt-6">
            <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
            Add Your First Key
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="flex justify-center bg-card p-4 rounded-lg min-h-[50vh]">
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