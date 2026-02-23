import {
  Activity,
  AlertOctagon,
  BarChart3,
  Bell,
  ChevronLeft,
  ChevronRight,
  FileText,
  Globe,
  HelpCircle,
  Home,
  LayoutDashboard,
  Layers,
  MonitorPlay,
  Search,
  Settings,
  Users,
  Zap,
  Braces,
  type LucideIcon,
} from "lucide-react"
import { useState } from "react"
import { useNavigate, useRouterState } from "@tanstack/react-router"
import { cn } from "@/lib/utils"
import { useSidebarStore } from "@/hooks/store/sidebar.store"
import { useAuthStore } from "@/hooks/store/auth.store"
import { ROUTES } from "@/lib/constants/routes"

type NavItem = {
  label: string
  icon: LucideIcon
  href: string
  badge?: number
}

type NavSection = {
  title?: string
  items: NavItem[]
}

const navSections: NavSection[] = [
  {
    items: [
      { label: "Home", icon: Home, href: ROUTES.HOME },
      { label: "Search", icon: Search, href: "/search" },
    ],
  },
  {
    title: "Monitor",
    items: [
      { label: "Dashboards", icon: LayoutDashboard, href: ROUTES.DASHBOARD },
      { label: "Services", icon: Layers, href: ROUTES.EVENTS },
      { label: "Traces", icon: Zap, href: ROUTES.TRACES },
      { label: "Logs", icon: FileText, href: ROUTES.LOGS },
      { label: "Metrics", icon: Activity, href: ROUTES.METRICS },
    ],
  },
  {
    title: "Analyze",
    items: [
      { label: "Web Analytics", icon: Globe, href: "/web-analytics" },
      { label: "Product Analytics", icon: BarChart3, href: ROUTES.ANALYTICS },
      { label: "Error Tracking", icon: Braces, href: ROUTES.ERRORS },
      { label: "Session Replay", icon: MonitorPlay, href: ROUTES.SESSIONS },
    ],
  },
  {
    title: "Respond",
    items: [
      { label: "Alerts", icon: Bell, href: ROUTES.ALERTS },
      { label: "Incidents", icon: AlertOctagon, href: ROUTES.INCIDENTS },
    ],
  },
]

const bottomItems: NavItem[] = [
  { label: "Team", icon: Users, href: ROUTES.TEAM },
  { label: "Settings", icon: Settings, href: ROUTES.SETTINGS },
  { label: "Help", icon: HelpCircle, href: "/help" },
]

export const NavBarComponent = () => {
  const navigate = useNavigate()
  const currentPath = useRouterState({ select: (s) => s.location.pathname })
  const [activePath, setActivePath] = useState(currentPath)
  const { expanded, toggle } = useSidebarStore()
  const user = useAuthStore((s) => s.user)

  const userEmail = user?.profile.email ?? ""
  const userInitial = (user?.profile.name ?? user?.profile.email ?? "?")[0].toUpperCase()

  const handleNavigate = (href: string) => {
    setActivePath(href)
    navigate({ to: href as never })
  }

  return (
    <aside
      className={cn(
        "flex h-svh flex-col border-r border-sidebar-border bg-sidebar transition-[width] duration-200 ease-linear",
        expanded ? "w-[200px]" : "w-[52px]"
      )}
    >
      {/* Brand */}
      <div className="flex items-center gap-2.5 border-b border-sidebar-border px-3 py-3">
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary">
          <Zap className="h-3.5 w-3.5 text-primary-foreground" />
        </div>
        {expanded && (
          <span className="text-sm font-semibold text-foreground tracking-tight">
            Streamly
          </span>
        )}
        {expanded && (
          <button
            onClick={toggle}
            className="ml-auto rounded-md p-1 text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-colors"
            aria-label="Collapse sidebar"
          >
            <ChevronLeft className="h-3.5 w-3.5" />
          </button>
        )}
      </div>

      {!expanded && (
        <button
          onClick={toggle}
          className="mx-auto my-2 rounded-md p-1.5 text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-colors"
          aria-label="Expand sidebar"
        >
          <ChevronRight className="h-3.5 w-3.5" />
        </button>
      )}

      {/* Navigation */}
      <nav
        className="flex-1 overflow-y-auto px-2 py-1"
        aria-label="Main navigation"
      >
        {navSections.map((section, idx) => (
          <div key={idx} className={cn(section.title && "mt-4")}>
            {section.title && expanded && (
              <span className="mb-1 block px-2 text-[10px] font-bold uppercase tracking-widest text-muted-foreground/60">
                {section.title}
              </span>
            )}
            {section.title && !expanded && idx > 0 && (
              <div className="mx-2 my-2 h-px bg-sidebar-border" />
            )}
            {section.items.map((item) => {
              const isActive = activePath === item.href
              const Icon = item.icon
              return (
                <button
                  key={item.label}
                  onClick={() => handleNavigate(item.href)}
                  aria-current={isActive ? "page" : undefined}
                  title={!expanded ? item.label : undefined}
                  className={cn(
                    "group relative flex w-full items-center gap-2.5 rounded-md px-2 py-[7px] text-[13px] transition-colors",
                    !expanded && "justify-center px-0",
                    isActive
                      ? "bg-primary/10 text-primary font-medium"
                      : "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                  )}
                >
                  {isActive && (
                    <span className="absolute left-0 top-1/2 h-4 w-[3px] -translate-y-1/2 rounded-r-full bg-primary" />
                  )}
                  <Icon
                    className={cn(
                      "h-[18px] w-[18px] shrink-0",
                      isActive
                        ? "text-primary"
                        : "text-muted-foreground group-hover:text-sidebar-accent-foreground"
                    )}
                  />
                  {expanded && (
                    <>
                      <span className="truncate">{item.label}</span>
                      {item.badge && (
                        <span className="ml-auto flex h-[18px] min-w-[18px] items-center justify-center rounded-full bg-destructive/20 px-1 text-[10px] font-bold text-destructive">
                          {item.badge}
                        </span>
                      )}
                    </>
                  )}
                  {!expanded && item.badge && (
                    <span className="absolute -right-0.5 -top-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-destructive text-[8px] font-bold text-destructive-foreground">
                      {item.badge}
                    </span>
                  )}
                </button>
              )
            })}
          </div>
        ))}
      </nav>

      {/* Bottom */}
      <div className="border-t border-sidebar-border px-2 py-2">
        {bottomItems.map((item) => {
          const Icon = item.icon
          const isActive = activePath === item.href
          return (
            <button
              key={item.label}
              onClick={() => handleNavigate(item.href)}
              title={!expanded ? item.label : undefined}
              className={cn(
                "group flex w-full items-center gap-2.5 rounded-md px-2 py-[7px] text-[13px] transition-colors",
                !expanded && "justify-center px-0",
                isActive
                  ? "bg-primary/10 text-primary font-medium"
                  : "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
              )}
            >
              <Icon
                className={cn(
                  "h-[18px] w-[18px] shrink-0",
                  isActive
                    ? "text-primary"
                    : "text-muted-foreground group-hover:text-sidebar-accent-foreground"
                )}
              />
              {expanded && <span className="truncate">{item.label}</span>}
            </button>
          )
        })}

        {/* User avatar */}
        <div className="mt-1 border-t border-sidebar-border pt-2">
          {expanded ? (
            <div className="flex items-center gap-2.5 rounded-md px-2 py-2">
              <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/20 text-[11px] font-bold text-primary">
                {userInitial}
              </div>
              <span className="truncate text-xs text-sidebar-foreground">{userEmail}</span>
            </div>
          ) : (
            <div className="flex justify-center py-1">
              <div className="flex h-7 w-7 items-center justify-center rounded-full bg-primary/20 text-[11px] font-bold text-primary">
                {userInitial}
              </div>
            </div>
          )}
        </div>
      </div>
    </aside>
  )
}
