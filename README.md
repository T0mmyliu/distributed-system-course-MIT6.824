MIT's graduate-level DS course with a focus on fault tolerance, replication, and consistency, all taught via awesome lab assignments in Golang!

### There are four main labs:

#### Lab 1: MapReduce
In this lab you'll build a MapReduce library as a way to learn the Go programming language and as a way to learn about fault tolerance in distributed systems. In the first part you will write a simple MapReduce program. In the second part you will write a Master that hands out tasks to workers, and handles failures of workers. The interface to the library and the approach to fault tolerance is similar to the one described in the original MapReduce paper.

#### lab 2: Raft
This is the first in a series of labs in which you'll build a fault-tolerant key/value storage system. You'll start in this lab by implementing Raft, a replicated state machine protocol. In the next lab lab you'll build a key/value service on top of Raft. Then you will “shard” your service for higher performance, and finally implement transactional operations across shards.

#### lab 3: Fault-tolerant Key/Value Service
In this lab you will build a fault-tolerant key-value storage service using your Raft library from lab 2. You will build your key-value service as a replicated state machine, consisting of several key-value servers that coordinate their activities through the Raft log. Your key/value service should continue to process client requests as long as a majority of the servers are alive and can communicate, in spite of other failures or network partitions.

#### lab 4: Sharded Key/Value Service
In this lab you'll build a key/value storage system that "shards," or partitions, the keys over a set of replica groups. A shard is a subset of the key/value pairs; for example, all the keys starting with "a" might be one shard, all the keys starting with "b" another, etc. The reason for sharding is performance. Each replica group handles puts and gets for just a few of the shards, and the groups operate in parallel; thus total system throughput (puts and gets per unit time) increases in proportion to the number of groups.


