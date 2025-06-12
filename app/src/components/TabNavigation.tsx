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
    title: 'Kubernetes',
    fullTitle: 'Kubernetes Installation',
    description: 'Install and configure Kubernetes on the cluster',
    icon: Container,
    isPhase: true,
  },
  {
    id: 'apps' as Tab,
    title: 'Apps',
    fullTitle: 'App Management',
    description: 'Install and manage applications on the cluster',
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
                  !isActive && status === 'completed' && "border-transparent text-green-600 hover:text-green-700",
                  !isActive && status === 'available' && "border-transparent text-muted-foreground",
                  status === 'locked' && "border-transparent text-muted-foreground/50 cursor-not-allowed"
                )}
              >
                <div className={cn(
                  "p-1 rounded-md",
                  isActive && "bg-primary/10",
                  status === 'completed' && !isActive && "bg-green-100 dark:bg-green-900",
                  status === 'locked' && "bg-muted"
                )}>
                  {status === 'completed' && tab.isPhase ? (
                    <CheckCircle className="h-4 w-4 text-green-600" />
                  ) : status === 'locked' ? (
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
      
      {/* Show current tab description */}
      <div className="mt-4 px-4 py-2 bg-muted/30 rounded-lg">
        <div className="flex items-center gap-2">
          <div className="p-1 bg-primary/10 rounded-md">
            {(() => {
              const currentTabData = tabs.find(t => t.id === currentTab);
              const Icon = currentTabData?.icon || Server;
              return <Icon className="h-4 w-4 text-primary" />;
            })()}
          </div>
          <div>
            <h3 className="font-medium text-foreground">
              {tabs.find(t => t.id === currentTab)?.fullTitle}
            </h3>
            <p className="text-sm text-muted-foreground">
              {tabs.find(t => t.id === currentTab)?.description}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}