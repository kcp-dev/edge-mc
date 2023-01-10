# Mailbox Controller

This is a component in an architecture not yet shared in detail.  It is an early milestone in edge multicluster, achieving partial functionality and implemented in a simple crude approach that is layered on top of transparent multicluster.  In this approach the edge-mc layer maintains a kcp workspace, called a mailbox, per selected edge destination per EdgePlacement object.  Then TMC syncs between the mailbox workspaces ad the edge destinations.

## Development Status

Very preliminary.  We do not even have the EdgePlacement resource defined yet.  We have not yet set up the infrastructure for defining resources and generating the consequent code.

The current mailbox-controller simply logs the notifications from an informer on workspaces.

This controller is currently a stand-alone process.  No leader election.  Not containerized.

## Build

Build with ordinary go commands.

## Usage

Suppose there is already a kcp-core server.

Launch the mailbox-controller with the following considerations.

- working directory is not important.
- `$KUBECONFIG` or `--kubeconfig` point to a kube client config aimed at the parent of the workspaces to be observed.
