import { Link } from '@tanstack/react-router'
import {ROUTES} from "@/lib/constants/routes.ts";

export default function HeaderComponent() {
  return (
    <header className="p-2 flex gap-2 justify-between">
      <nav className="flex flex-row">
        <div className="px-2 font-bold">
          <Link to="/">Home</Link>
        </div>

        <div className="px-2 font-bold">
          <Link to={ROUTES.EVENTS}>TanStack Query</Link>
        </div>
      </nav>
    </header>
  )
}
