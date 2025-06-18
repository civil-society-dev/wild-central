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

const tabs = [
  {
    id: 'setup' as Tab,
    title: 'Setup',
    fullTitle: 'Wild-cloud Central Setup',
    description: 'Configure the central server and dnsmasq services',
    icon: Server,
    isPhase: true,
  },
  {
    id: 'infrastructure' as Tab,
    title: 'Infrastructure',
    fullTitle: 'Infrastructure Setup',
    description: 'Connect controller and worker nodes to the wild-cloud',
    icon: Play,
    isPhase: true,
  },
  {
    id: 'cluster' as Tab,
    title: 'Cluster',
    fullTitle: 'Kubernetes Installation',
    description: 'Install and configure essential cluster services',
    icon: Container,
    isPhase: true,
  },
  {
    id: 'apps' as Tab,
    title: 'Apps',
    fullTitle: 'App Management',
    description: 'Install and manage applications',
    icon: AppWindow,
    isPhase: true,
  },
  {
    id: 'advanced' as Tab,
    title: 'Advanced',
    fullTitle: 'Advanced Configuration',
    description: 'Advanced settings and system configuration',
    icon: Settings,
    isPhase: false,
  },
];

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

  const getTabStatus = (tab: Tab, index: number) => {
    // Non-phase tabs (like advanced) are always available
    const tabData = tabs.find(t => t.id === tab);
    if (!tabData?.isPhase) {
      return 'available';
    }
    
    // For phase tabs, check completion status
    if (completedPhases.includes(tab as Phase)) {
      return 'completed';
    }
    
    // Allow access to the first phase always
    if (index === 0) {
      return 'available';
    }
    
    // Allow access to the next phase if the previous phase is completed
    const previousTab = tabs[index - 1];
    if (previousTab?.isPhase && completedPhases.includes(previousTab.id as Phase)) {
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
              {tabs.map((tab, index) => {
                const status = getTabStatus(tab.id, index);
                const isActive = currentTab === tab.id;
                const Icon = tab.icon;
                
                return (
                  <SidebarMenuItem key={tab.id}>
                    <SidebarMenuButton
                      isActive={isActive}
                      onClick={() => status !== 'locked' && onTabChange(tab.id)}
                      disabled={status === 'locked'}
                      tooltip={tab.description}
                      className={cn(
                        "transition-colors",
                        status === 'locked' && "opacity-50 cursor-not-allowed"
                      )}
                    >
                      <div className={cn(
                        "p-1 rounded-md",
                        isActive && "bg-primary/10",
                        status === 'locked' && "bg-muted"
                      )}>
                        {status === 'locked' ? (
                          <Lock className="h-4 w-4 text-muted-foreground/50" />
                        ) : status === 'completed' ? (
                          <CheckCircle className="h-4 w-4 text-green-600" />
                        ) : (
                          <Icon className={cn(
                            "h-4 w-4",
                            isActive && "text-primary",
                            !isActive && "text-muted-foreground"
                          )} />
                        )}
                      </div>
                      <span className="truncate">{tab.title}</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                );
              })}
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