import type {ReactNode} from "react";
import {NavBarComponent} from "@/components/nav-bar.component.tsx";
import {SidebarProvider, SidebarInset, SidebarTrigger} from "@/components/ui/sidebar.tsx";

const RootContainer = ({ children} :{children: ReactNode})=> (
  <div className='min-h-screen text-foreground'>
    <SidebarProvider>
      <NavBarComponent />
      <SidebarInset>
        <div className="flex flex-1 flex-col gap-4 p-4" >
          <SidebarTrigger variant={'ghost'}/>
          <div className="bg-background border border-border rounded-xl shadow-sm min-h-[calc(100vh-2rem)] p-6">
            {children}
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
)

export default RootContainer;
