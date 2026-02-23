import type { ReactNode } from "react"
import { NavBarComponent } from "@/components/nav-bar.component.tsx"
import HeaderComponent from "@/components/header.component.tsx"

const RootContainer = ({ children }: { children: ReactNode }) => (
  <div className="flex min-h-screen text-foreground">
    <NavBarComponent />
    <main
      className="flex flex-1 flex-col overflow-hidden"
      style={{ border: "0.5px solid var(--main-container-border)" }}
    >
      <div id="content-header">
        <HeaderComponent />
      </div>
      <div
        id="container-layout"
        className="bg-background rounded-b-md shadow-sm overflow-auto flex-1 p-2"
      >
        {children}
      </div>
    </main>
  </div>
)

export default RootContainer
