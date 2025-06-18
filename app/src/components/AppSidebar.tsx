import { CheckCircle, Lock, Server, Play, Container, AppWindow, Settings, CloudLightning, Sun, Moon, Monitor } from 'lucide-react';
import { cn } from '../lib/utils';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarRail,
} from './ui/sidebar';
import { useTheme } from '../contexts/ThemeContext';

export type Phase = 'setup' | 'infrastructure' | 'cluster' | 'apps';
export type Tab = Phase | 'advanced';

interface AppSidebarProps {
  currentTab: Tab;
  onTabChange: (tab: Tab) => void;
  completedPhases: Phase[];
}


export function AppSidebar({ currentTab, onTabChange, completedPhases }: AppSidebarProps) {
  const { theme, setTheme } = useTheme();

  const cycleTheme = () => {
    if (theme === 'light') {
      setTheme('dark');
    } else if (theme === 'dark') {
      setTheme('system');
    } else {
      setTheme('light');
    }
  };

  const getThemeIcon = () => {
    switch (theme) {
      case 'light':
        return <Sun className="h-4 w-4" />;
      case 'dark':
        return <Moon className="h-4 w-4" />;
      default:
        return <Monitor className="h-4 w-4" />;
    }
  };

  const getThemeLabel = () => {
    switch (theme) {
      case 'light':
        return 'Light mode';
      case 'dark':
        return 'Dark mode';
      default:
        return 'System theme';
    }
  };

  const getTabStatus = (tab: Tab) => {
    // Non-phase tabs (like advanced) are always available
    if (tab === 'advanced') {
      return 'available';
    }
    
    // For phase tabs, check completion status
    if (completedPhases.includes(tab as Phase)) {
      return 'completed';
    }
    
    // Allow access to the first phase always
    if (tab === 'setup') {
      return 'available';
    }
    
    // Allow access to the next phase if the previous phase is completed
    if (tab === 'infrastructure' && completedPhases.includes('setup')) {
      return 'available';
    }
    
    if (tab === 'cluster' && completedPhases.includes('infrastructure')) {
      return 'available';
    }
    
    if (tab === 'apps' && completedPhases.includes('cluster')) {
      return 'available';
    }
    
    return 'locked';
  };

  return (
    <Sidebar variant="sidebar" collapsible="icon">
      <SidebarHeader>
        <div className="flex items-center gap-2 px-2">
          <div className="p-1 bg-primary/10 rounded-lg">
            <CloudLightning className="h-6 w-6 text-primary" />
          </div>
          <div className="group-data-[collapsible=icon]:hidden">
            <h2 className="text-lg font-bold text-foreground">Wild Cloud</h2>
            <p className="text-sm text-muted-foreground">Central</p>
          </div>
        </div>
      </SidebarHeader>
      
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Setup & Management</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentTab === 'setup'}
                  onClick={() => {
                    const status = getTabStatus('setup');
                    if (status !== 'locked') onTabChange('setup');
                  }}
                  disabled={getTabStatus('setup') === 'locked'}
                  tooltip="Configure the central server and dnsmasq services"
                  className={cn(
                    "transition-colors",
                    getTabStatus('setup') === 'locked' && "opacity-50 cursor-not-allowed"
                  )}
                >
                  <div className={cn(
                    "p-1 rounded-md",
                    currentTab === 'setup' && "bg-primary/10",
                    getTabStatus('setup') === 'locked' && "bg-muted"
                  )}>
                    <Server className={cn(
                      "h-4 w-4",
                      currentTab === 'setup' && "text-primary",
                      currentTab !== 'setup' && "text-muted-foreground"
                    )} />
                  </div>
                  <span className="truncate">Central Services</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentTab === 'infrastructure'}
                  onClick={() => {
                    const status = getTabStatus('infrastructure');
                    if (status !== 'locked') onTabChange('infrastructure');
                  }}
                  disabled={getTabStatus('infrastructure') === 'locked'}
                  tooltip="Connect controller and worker nodes to the wild-cloud"
                  className={cn(
                    "transition-colors",
                    getTabStatus('infrastructure') === 'locked' && "opacity-50 cursor-not-allowed"
                  )}
                >
                  <div className={cn(
                    "p-1 rounded-md",
                    currentTab === 'infrastructure' && "bg-primary/10",
                    getTabStatus('infrastructure') === 'locked' && "bg-muted"
                  )}>
                    <Play className={cn(
                      "h-4 w-4",
                      currentTab === 'infrastructure' && "text-primary",
                      currentTab !== 'infrastructure' && "text-muted-foreground"
                    )} />
                  </div>
                  <span className="truncate">Cluster Nodes</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentTab === 'cluster'}
                  onClick={() => {
                    const status = getTabStatus('cluster');
                    if (status !== 'locked') onTabChange('cluster');
                  }}
                  disabled={getTabStatus('cluster') === 'locked'}
                  tooltip="Install and configure essential cluster services"
                  className={cn(
                    "transition-colors",
                    getTabStatus('cluster') === 'locked' && "opacity-50 cursor-not-allowed"
                  )}
                >
                  <div className={cn(
                    "p-1 rounded-md",
                    currentTab === 'cluster' && "bg-primary/10",
                    getTabStatus('cluster') === 'locked' && "bg-muted"
                  )}>
                    <Container className={cn(
                      "h-4 w-4",
                      currentTab === 'cluster' && "text-primary",
                      currentTab !== 'cluster' && "text-muted-foreground"
                    )} />
                  </div>
                  <span className="truncate">Cluster Services</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentTab === 'apps'}
                  onClick={() => {
                    const status = getTabStatus('apps');
                    if (status !== 'locked') onTabChange('apps');
                  }}
                  disabled={getTabStatus('apps') === 'locked'}
                  tooltip="Install and manage applications"
                  className={cn(
                    "transition-colors",
                    getTabStatus('apps') === 'locked' && "opacity-50 cursor-not-allowed"
                  )}
                >
                  <div className={cn(
                    "p-1 rounded-md",
                    currentTab === 'apps' && "bg-primary/10",
                    getTabStatus('apps') === 'locked' && "bg-muted"
                  )}>
                    <AppWindow className={cn(
                      "h-4 w-4",
                      currentTab === 'apps' && "text-primary",
                      currentTab !== 'apps' && "text-muted-foreground"
                    )} />
                  </div>
                  <span className="truncate">Apps</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentTab === 'advanced'}
                  onClick={() => {
                    const status = getTabStatus('advanced');
                    if (status !== 'locked') onTabChange('advanced');
                  }}
                  disabled={getTabStatus('advanced') === 'locked'}
                  tooltip="Advanced settings and system configuration"
                  className={cn(
                    "transition-colors",
                    getTabStatus('advanced') === 'locked' && "opacity-50 cursor-not-allowed"
                  )}
                >
                  <div className={cn(
                    "p-1 rounded-md",
                    currentTab === 'advanced' && "bg-primary/10",
                    getTabStatus('advanced') === 'locked' && "bg-muted"
                  )}>
                    <Settings className={cn(
                      "h-4 w-4",
                      currentTab === 'advanced' && "text-primary",
                      currentTab !== 'advanced' && "text-muted-foreground"
                    )} />
                  </div>
                  <span className="truncate">Advanced</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              onClick={cycleTheme}
              tooltip={`Current: ${getThemeLabel()}. Click to cycle themes.`}
            >
              {getThemeIcon()}
              <span>{getThemeLabel()}</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail/>
    </Sidebar>
  );
}