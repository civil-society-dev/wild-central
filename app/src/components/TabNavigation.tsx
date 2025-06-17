import { CheckCircle, Lock, Server, Play, Container, AppWindow, Settings } from 'lucide-react';
import { cn } from '../lib/utils';

export type Phase = 'setup' | 'infrastructure' | 'cluster' | 'apps';
export type Tab = Phase | 'advanced';

interface TabNavigationProps {
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

export function TabNavigation({ currentTab, onTabChange, completedPhases }: TabNavigationProps) {
  console.log('TabNavigation props:', { currentTab, completedPhases });
  
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
    <div className="mb-8">
      <div className="border-b border-border">
        <nav className="flex space-x-8 overflow-x-auto">
          {tabs.map((tab, index) => {
            const status = getTabStatus(tab.id, index);
            const isActive = currentTab === tab.id;
            const Icon = tab.icon;
            
            return (
              <button
                key={tab.id}
                onClick={() => status !== 'locked' && onTabChange(tab.id)}
                disabled={status === 'locked'}
                className={cn(
                  "flex items-center gap-2 px-4 py-3 border-b-2 text-sm font-medium transition-colors whitespace-nowrap",
                  "hover:text-foreground hover:border-border",
                  isActive && "border-primary text-primary",
                  !isActive && status === 'completed' && "border-transparent",
                  !isActive && status === 'available' && "border-transparent text-muted-foreground",
                  status === 'locked' && "border-transparent text-muted-foreground/50 cursor-not-allowed"
                )}
              >
                <div className={cn(
                  "p-1 rounded-md",
                  isActive && "bg-primary/10",
                  status === 'locked' && "bg-muted"
                )}>
                  {status === 'locked' ? (
                    <Lock className="h-4 w-4 text-muted-foreground/50" />
                  ) : (
                    <Icon className={cn(
                      "h-4 w-4",
                      isActive && "text-primary",
                      !isActive && "text-muted-foreground"
                    )} />
                  )}
                </div>
                <span className="hidden sm:inline">
                  {tab.title}
                </span>
                <span className="sm:hidden">
                  {tab.isPhase ? index + 1 : 'A'}
                </span>
              </button>
            );
          })}
        </nav>
      </div>
      
    </div>
  );
}