# Building Wild Cloud Central

The first version of Wild Cloud, the Proof of Concept version (v.PoC), was built as a collection of shell scripts that users would run from their local machines. This works well for early adopters who are comfortable with the command line, Talos, and Kubernetes.

To make Wild Cloud more accessible to a broader audience, we are developing Wild Central. Central is a single-purpose machine run on a LAN that will deliver:

- Wild Daemon: A lightweight service that runs on a local machine (e.g., a Raspberry Pi) to manage Wild Cloud instances on the local network.
- Wild App: A web-based interface (to Wild Daemon) for managing Wild Cloud instances.
- Wild CLI: A command-line interface (to Wild Daemon) for advanced users who prefer to manage Wild Cloud from the terminal.

## Background info

### Info about Wild Cloud v.PoC

- @docs/agent-context/wildcloud/README.md
- @docs/agent-context/wildcloud/overview.md
- @docs/agent-context/wildcloud/project-architecture.md
- @docs/agent-context/wildcloud/bin-scripts.md
- @docs/agent-context/wildcloud/configuration-system.md
- @docs/agent-context/wildcloud/setup-process.md
- @docs/agent-context/wildcloud/apps-system.md

### Info about Talos

- @docs/agent-context/talos-v1.11/README.md
- @docs/agent-context/talos-v1.11/architecture-and-components.md
- @docs/agent-context/talos-v1.11/cli-essentials.md
- @docs/agent-context/talos-v1.11/cluster-operations.md
- @docs/agent-context/talos-v1.11/discovery-and-networking.md
- @docs/agent-context/talos-v1.11/etcd-management.md
- @docs/agent-context/talos-v1.11/bare-metal-administration.md
- @docs/agent-context/talos-v1.11/troubleshooting-guide.md

## Architecture

### Old v.PoC Architecture

- WC_ROOT: The scripts used to set up and manage the Wild Cloud cluster. Currently, this is a set of shell scripts in $WC_ROOT/bin.
- WC_HOME: During setup, the user creates a Wild Cloud project directory (WC_HOME) on their local machine. This directory holds all configuration, secrets, and k8s manifests for their specific Wild Cloud deployment.
- Wild Cloud Apps Directory: The Wild Cloud apps are stored in the `apps/` directory within the WC_ROOT repository. Users can deploy these apps to their cluster using the scripts in WC_ROOT/bin.
- dnsmasq server: Scripts help the operator set up a dnsmasq server on a separate machine to provide LAN DNS services during node bootstrapping.

### New Wild Central Architecture

#### wildd: The Wild Cloud Daemon

wildd is a long-running service that provides an API and web interface for managing one or more Wild Cloud clusters. It runs on a dedicated device within the user's network.

wildd replaces functionality from the v.PoC scripts and the dnsmasq server. It is one API for managing multiple wild cloud instances on the LAN.

Both wild-app and wild-cli communicate with wildd to perform actions.

See: @daemon/BUILDING_WILD_DAEMON.md

#### wild-app

The web application that provides the user interface for Wild Cloud on Wild Central. It communicates with wildd to perform actions and display information.

See: @/app/BUILDING_WILD_APP.md

#### wild-cli

A command-line interface for advanced users who prefer to manage Wild Cloud from the terminal. It communicates with wildd to perform actions.

Mirrors all of the wild-* scripts from v.PoC, but adapted for the new architecture:

- One golang client (wild-cli) replaces many bash scripts (wild-*).
- Wrapper around wildd API instead of direct file manipulation.
- Multi-cloud: v.PoC scripts set the instance context with WC_HOME environment variable. In Central, wild-cli follows the "context" pattern like kubectl and talosctl, using `--context` or `WILD_CONTEXT` to select which wild cloud instance to manage, or defaulting to the "current" context.

See: @cli/BUILDING_WILD_CLI.md

#### Wild Central Data

Configured with $WILD_CENTRAL_DATA environment variable (default: /var/lib/wild-central).

Replaces multiple WC_HOMEs. All wild clouds managed on the LAN are configured here. These are still in easy to read YAML format and can be edited directly or through the webapp.

Wild Central data also holds the local app directory, logs, and artifacts, and overall state data.

#### Wild Cloud Apps Directory

The Wild Cloud apps are stored in the `apps/` directory within the WC_ROOT repository. Users can deploy these apps to their cluster using the webapp or wild-cli.

#### dnsmasq server

The Wild Daemon (wildd) includes functionality to manage a dnsmasq server on the same device, providing LAN DNS services during node bootstrapping.

## Packaging and Installation

Ultimately, the daemon, app, and cli will be packaged together for easy installation on a Raspberry Pi or similar device.

See @ai/WILD_CENTRAL_PACKAGING.md

## Implementation Philosophy

## Core Philosophy

Embodies a Zen-like minimalism that values simplicity and clarity above all. This approach reflects:

- **Wabi-sabi philosophy**: Embracing simplicity and the essential. Each line serves a clear purpose without unnecessary embellishment.
- **KISS**: The solution should be as simple as possible, but no simpler.
- **YAGNI**: Avoid building features or abstractions that aren't immediately needed. The code handles what's needed now rather than anticipating every possible future scenario.
- **Trust in emergence**: Complex systems work best when built from simple, well-defined components that do one thing well.
- **Pragmatic trust**: The developer trusts external systems enough to interact with them directly, handling failures as they occur rather than assuming they'll happen.
- **Consistency is key**: Uniform patterns and conventions make the codebase easier to understand and maintain. If you introduce a new pattern, make sure it's consistently applied. There should be one obvious way to do things.

This development philosophy values clear, concise documentation, readable code, and belief that good architecture emerges from simplicity rather than being imposed through complexity.

## Core Design Principles

### 1. Ruthless Simplicity

- **KISS principle taken to heart**: Keep everything as simple as possible, but no simpler
- **Minimize abstractions**: Every layer of abstraction must justify its existence
- **Start minimal, grow as needed**: Begin with the simplest implementation that meets current needs
- **Avoid future-proofing**: Don't build for hypothetical future requirements
- **Question everything**: Regularly challenge complexity in the codebase

### 2. Architectural Integrity with Minimal Implementation

- **Preserve key architectural patterns**: Maintain clear boundaries and responsibilities
- **Simplify implementations**: Maintain pattern benefits with dramatically simpler code
- **Scrappy but structured**: Lightweight implementations of solid architectural foundations
- **End-to-end thinking**: Focus on complete flows rather than perfect components

### 3. Library vs Custom Code

Choosing between custom code and external libraries is a judgment call that evolves with your requirements. There's no rigid rule - it's about understanding trade-offs and being willing to revisit decisions as needs change.

#### The Evolution Pattern

Your approach might naturally evolve:
- **Start simple**: Custom code for basic needs (20 lines handles it)
- **Growing complexity**: Switch to a library when requirements expand
- **Hitting limits**: Back to custom when you outgrow the library's capabilities

This isn't failure - it's natural evolution. Each stage was the right choice at that time.

#### When Custom Code Makes Sense

Custom code often wins when:
- The need is simple and well-understood
- You want code perfectly tuned to your exact requirements
- Libraries would require significant "hacking" or workarounds
- The problem is unique to your domain
- You need full control over the implementation

#### When Libraries Make Sense

Libraries shine when:
- They solve complex problems you'd rather not tackle (auth, crypto, video encoding)
- They align well with your needs without major modifications
- The problem is well-solved with mature, battle-tested solutions
- Configuration alone can adapt them to your requirements
- The complexity they handle far exceeds the integration cost

#### Making the Judgment Call

Ask yourself:
- How well does this library align with our actual needs?
- Are we fighting the library or working with it?
- Is the integration clean or does it require workarounds?
- Will our future requirements likely stay within this library's capabilities?
- Is the problem complex enough to justify the dependency?

#### Recognizing Misalignment

Watch for signs you're fighting your current approach:
- Spending more time working around the library than using it
- Your simple custom solution has grown complex and fragile  
- You're monkey-patching or heavily wrapping a library
- The library's assumptions fundamentally conflict with your needs

#### Stay Flexible

Remember that complexity isn't destroyed, only moved. Libraries shift complexity from your code to someone else's - that's often a great trade, but recognize what you're doing.

The key is avoiding lock-in. Keep library integration points minimal and isolated so you can switch approaches when needed. There's no shame in moving from custom to library or library to custom. Requirements change, understanding deepens, and the right answer today might not be the right answer tomorrow. Make the best decision with current information, and be ready to evolve.

## Technical Implementation Guidelines

### API Layer

- Implement only essential endpoints
- Minimal middleware with focused validation
- Clear error responses with useful messages
- Consistent patterns across endpoints

### Storage

- Prefer simple file storage
- Simple schema focused on current needs

## Development Approach

### Vertical Slices

- Implement complete end-to-end functionality slices
- Start with core user journeys
- Get data flowing through all layers early
- Add features horizontally only after core flows work

### Iterative Implementation

- 80/20 principle: Focus on high-value, low-effort features first
- One working feature > multiple partial features
- Validate with real usage before enhancing
- Be willing to refactor early work as patterns emerge

### Testing Strategy

- Focus on critical path testing initially
- Add unit tests for complex logic and edge cases
- Testing pyramid: 60% unit, 30% integration, 10% end-to-end

### Error Handling

- Handle common errors robustly
- Log detailed information for debugging
- Provide clear error messages to users
- Fail fast and visibly during development

## Decision-Making Framework

When faced with implementation decisions, ask these questions:

1. **Necessity**: "Do we actually need this right now?"
2. **Simplicity**: "What's the simplest way to solve this problem?"
3. **Directness**: "Can we solve this more directly?"
4. **Value**: "Does the complexity add proportional value?"
5. **Maintenance**: "How easy will this be to understand and change later?"

## Areas to Embrace Complexity

Some areas justify additional complexity:

1. **Security**: Never compromise on security fundamentals
2. **Data integrity**: Ensure data consistency and reliability
3. **Core user experience**: Make the primary user flows smooth and reliable
4. **Error visibility**: Make problems obvious and diagnosable

## Areas to Aggressively Simplify

Push for extreme simplicity in these areas:

1. **Internal abstractions**: Minimize layers between components
2. **Generic "future-proof" code**: Resist solving non-existent problems
3. **Edge case handling**: Handle the common cases well first
4. **Framework usage**: Use only what you need from frameworks
5. **State management**: Keep state simple and explicit

## Remember

- It's easier to add complexity later than to remove it
- Code you don't write has no bugs
- Favor clarity over cleverness
- The best code is often the simplest

This philosophy document serves as the foundational guide for all implementation decisions in the project.
