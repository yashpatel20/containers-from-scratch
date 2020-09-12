# Containers from scratch

Containers are about building up dependencies so we can ship code around in a repeatable and secure way.

To understand containers we need to understand these three concepts:

## Namespaces

Namespaces provide the isolation needed to run multiple containers on one machine while giving each what appears like its own environment.

There are six namespaces. Each can be independently requested and amounts to giving a process a view of subset of the resources of the machine.

The namespaces are:

- PID

  - Gives process and its children thier own view of a subset of the processes in the system.
  - PID namespace makes the first process created within it pid 1 (by mapping its host ID to 1), giving the appearence of an isolated process tree in container.

- MNT

  - Mount namepsace gives the process contained within it their own mount table.
  - In combination with the pivot_root syscall, it allows a process to have its own filesytem.
  - This is how we can have a process think its running on ubuntu or alpine by swapping out the filesystem the container sees.

- NET

  - Network namespace gives the processes that use it their own network stack.
  - Only the main network namespace actually have real physical network cards attached.
  - We can create virtual ethernet pairs - linked ethernet cards where one end is in one network namespace and one in another.
  - This allows each container to talk to the real world while isolating each to its own network stack.

- UTS

  - Gives its processes their own view of the systems hostname and domain name.

- IPC

  - Isolates various inter process communications mechanisms such as message queues.

- USER

  - It maps the uids a process sees to a a different set of uids and gids on the host.
  - We can mao the containers root user ID i.e 0 to an arbitary unprivileged uid on host.
  - This means we can let the container think it has root access, can give it root like permissions on container specific resources without giving it any privileges on host.
  - There simply isnt a uid in the container tbat has real root permissions.

## CGroups

CGroups collect a set of process or task ids together and apply limits to them.

Where namespaces isolate a proces, cgroups enforce a fair or unfair resource sharing between processs.

System resources such as CPU, memory, disk and network bandwidth can be restricted by these cgroups.

CGroups allow us to sanely partition the resources and easily schedule our container based processes using the CFS (Completely Fair Scheduler)

## Layered Filesystems

Namespaces and CGroups are the isolation and resource sharing sides of containerisation.

Layered file systems are how we can efficiently move whole machine images around.

At a basic level, layered filesystems amount to optimizing the call to create a copy of the root filesystem for each container.
