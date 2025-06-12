import { CheckCircle, Lock, Server, Play, Container, AppWindow, Settings } from 'lucide-react';
import { cn } from '../lib/utils';

export type Phase = 'setup' | 'infrastructure' | 'cluster' | 'apps' | 'advanced';

interface PhaseNavigationProps {
  currentPhase: Phase;
  onPhaseChange: (phase: Phase) => void;
  completedPhases: Phase[];
}

const phases = [
  {
    id: 'setup' as Phase,
    title: 'Setup',
    fullTitle: 'Wild-cloud Central Setup',
    description: 'Configure the central server and dnsmasq services',
    icon: Server,
  },
  {
    id: 'infrastructure' as Phase,
    title: 'Infrastructure',
    fullTitle: 'Infrastructure Setup',
    description: 'Connect controller and worker nodes to the wild-cloud',
    icon: Play,
  },
  {
    id: 'cluster' as Phase,
    title: 'Kubernetes',
    fullTitle: 'Kubernetes Installation',
    description: 'Install and configure Kubernetes on the cluster',
    icon: Container,
  },
  {
    id: 'apps' as Phase,
    title: 'Apps',
    fullTitle: 'App Management',
    description: 'Install and manage applications on the cluster',
    icon: AppWindow,
  },
  {
    id: 'advanced' as Phase,
    title: 'Advanced',
    fullTitle: 'Advanced Configuration',
    description: 'Advanced settings and system configuration',
    icon: Settings,
  },
];

export function PhaseNavigation({ currentPhase, onPhaseChange, completedPhases }: PhaseNavigationProps) {
  console.log('PhaseNavigation props:', { currentPhase, completedPhases });
  
  const getPhaseStatus = (phase: Phase, index: number) => {
    if (completedPhases.includes(phase)) {
      return 'completed';
    }
    
    // Advanced tab is always available
    if (phase === 'advanced') {
      return 'available';
    }
    
    // Allow access to the first phase always
    if (index === 0) {
      return 'available';
    }
    
    // Allow access to the next phase if the previous phase is completed
    const previousPhase = phases[index - 1];
    if (completedPhases.includes(previousPhase.id)) {
      return 'available';
    }
    
    return 'locked';
  };

  return (
    <div className="mb-8">
      <div className="border-b border-border">
        <nav className="flex space-x-8 overflow-x-auto">
          {phases.map((phase, index) => {
            const status = getPhaseStatus(phase.id, index);
            const isActive = currentPhase === phase.id;
            const Icon = phase.icon;
            
            return (
              <button
                key={phase.id}
                onClick={() => status !== 'locked' && onPhaseChange(phase.id)}
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
                  {status === 'completed' ? (
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
                  {phase.title}
                </span>
                <span className="sm:hidden">
                  {index + 1}
                </span>
              </button>
            );
          })}
        </nav>
      </div>
      
      {/* Show current phase description */}
      <div className="mt-4 px-4 py-2 bg-muted/30 rounded-lg">
        <div className="flex items-center gap-2">
          <div className="p-1 bg-primary/10 rounded-md">
            {(() => {
              const currentPhaseData = phases.find(p => p.id === currentPhase);
              const Icon = currentPhaseData?.icon || Server;
              return <Icon className="h-4 w-4 text-primary" />;
            })()}
          </div>
          <div>
            <h3 className="font-medium text-foreground">
              {phases.find(p => p.id === currentPhase)?.fullTitle}
            </h3>
            <p className="text-sm text-muted-foreground">
              {phases.find(p => p.id === currentPhase)?.description}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}