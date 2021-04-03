# Plan 9 Crypto Monorepo

While all of the API's are in development there is intermittent and very
neccessary changes that break everything, everything is being piled into this
repository in order to simplify several aspects of the work.

### duod

    Blockchain Full node with multicast as well as GetBlockTemplate mining
    capabilities.

### duowallet

    A multisig capable wallet server that provides one account with HD keys as well
    as imported addresses and optionally their private keys. Depends on a connection
    to a `duod` though in the near future it will have also SPV node for lighter
    mobile deployments.

### ParallelCoin GUI

    A user friendly and ultra light weight GUI for managing transacting on the
    ParallelCoin chain with included block explorer, instant messenger and forum
    running over Tor on IPFS.

### gel - Gio Elements

    Our customised and extended widget toolkit derived from the `widget`,
    `material` and `layout` components of Gio, the UI library we use here.

### glom - The Less Confusion, The More Profit

    Glom aims to specifically solve the problem of managing complexity in tree and
    graph structured documents, in their source code format. Glom lets you move text
    around in arbitrary sized groups through a 'floating' state which can not be
    split and can only drop back into the document where there is some kind of
    bracketing location, based on the syntax. It also blends the role of file
    manager seamlessly with a tiling document editor panel, log viewer and filter,
    and built in source code archiving and binary caching used also in Spore.

### interrupt

    Handles the trapping of process interrupt signals and restarting of processes.

### log

    A streamlined logging system that makes tracking down bugs less manual labor and
    more pleasure.

### logo

    ParallelCoin Logo

### opts

    Concurrency friendly configuration database with instant synchronisation between
    application component executables over pipes and sockets, and configurable in a
    GUI.

### qu

    As the workhorse of Go's ultra low latency coroutine concurrency and
    parallelism, they are total black boxes. This library allows the programmer to
    more easily observe the interaction between threads and track down processes
    that have fallen into deadlock, and make them stop correctly. As it 
    operates with queues as well as use for quit signals it can have a very 
    short name

### spore

    Spore is at base both a decentralised software repository and a swarm of
    heteregenous, tightly scoped and secure array of API's, storage and
    communications engines and builds the future of Parallelcoin with finality,
    on-chain governance and interoperability with many other blockchain networks,
    built on the strong discipline of proof of work mining.
        
# Building

