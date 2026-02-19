import { CodeEditor } from "@/components/editor";
import type { ReactNode } from "react";

export const EventContainer = ({ children }: { children: ReactNode }) => {
  return (
    <div>
      <div className="flex justify-between w-full items-center">
        <h1>Event Queue Overview</h1>
        {children}
      </div>
      <CodeEditor />
    </div>
  )
}
