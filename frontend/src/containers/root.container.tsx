import type {ReactNode} from "react";
import {NavBarComponent} from "@/components/nav-bar.component.tsx";
import {SidebarProvider, SidebarInset} from "@/components/ui/sidebar.tsx";
import Header from "@/components/Header.tsx";

const RootContainer = ({ children} :{children: ReactNode})=> (
  <div className='min-h-screen text-foreground'>
    <SidebarProvider>
      <NavBarComponent />
      <SidebarInset className={"flex flex-1 flex-col"}>
        <div id={'content-header'}>
          <Header />
        </div>
          <div id={'container-layout'}
            className="bg-background rounded-b-md shadow-sm overflow-auto flex-1 p-2"
             style={{ border: '0.5px solid var(--main-container-border' }}
          >
            {children}
          </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
)

export default RootContainer;
