import type {FunctionComponent} from "react";
import {
  Sidebar,
  SidebarContent, SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel, SidebarMenu, SidebarMenuButton, SidebarMenuItem
} from "@/components/ui/sidebar.tsx";
import {Home, Radio, LayoutDashboard, Bell} from "lucide-react";
import {ROUTES} from "@/lib/constants/routes.ts";

type SidebarItem = {
  title: string;
  icon: FunctionComponent;
  url: string;
}

const items: Array<SidebarItem> = [
  {
    title: 'Home',
    icon: Home,
    url: ROUTES.HOME,
  },
  {
    title: 'Dashboard',
    icon: LayoutDashboard,
    url: ROUTES.HOME,
  },
  {
    title: 'Alert',
    icon: Bell,
    url: ROUTES.HOME,
  },
  {
    title: 'Events',
    icon: Radio,
    url: ROUTES.EVENTS,
  },
]

export const NavBarComponent: FunctionComponent = () => {
  return (
    <Sidebar variant={'inset'}>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild className="sidebar-menu-hover">
                    <a href={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>
            More
          </SidebarGroupLabel>
          <SidebarGroupContent>
            { /* next features will come here */ }
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          Hey
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
