import type {FunctionComponent} from "react";
import {Sidebar, SidebarContent} from "@/components/ui/sidebar.tsx";

export const NavBarComponent: FunctionComponent = () => {

  return (
    <Sidebar variant={'floating'}>
      <SidebarContent />
    </Sidebar>
  )
}
