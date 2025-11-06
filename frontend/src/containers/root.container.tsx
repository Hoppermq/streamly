import type {ReactNode} from "react";
import {NavBarComponent} from "@/components/nav-bar.component.tsx";
import {SidebarProvider, SidebarInset} from "@/components/ui/sidebar.tsx";
import HeaderComponent from "@/components/header.component.tsx";

const RootContainer = ({ children} :{children: ReactNode})=> (
  <div className='min-h-screen text-foreground'>
    <SidebarProvider>
      <NavBarComponent />
      <SidebarInset className={"flex flex-1 flex-col"}>
        <div id={'content-header'}>
          <HeaderComponent />
        </div>
          <div id={'container-layout'}
            className="bg-background border-top rounded-b-md shadow-sm overflow-auto flex-1 p-2"
          >
            {children}
          </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
)

export default RootContainer;
