# BeeTea

🚧 Experimental behaviour tree module for Go 🚧

A behavior tree doesn't have to stop at being a convenient model for representing plans for a singular system. Multi-agent coordination naturally benefits from communication, so we aim to design a BT implementation that can be efficiently communicated, with a view toward communicating BTs that are synthesised online (implying frequent updates).

Merkle trees, also known as hash trees, are a popular choice among the blockchain community beacuse they enable parts of a graph to be efficiently (*O(log n)*) compared. Each node contains the hash of all its children, and if any of them change, the hash changes. Only the hash needs to be communicated to know if anything in the sub-tree has changed, therefore only the difference needs to be communicated. 

A key invariant to note is that the tree should remain balanced and not get too deep, otherwise search could be impacted negatively. A BT is not a BST, and searching for a hash in deep trees will take longer than broad trees (you only need to find the top of a branch). One option to mitigate in practice is maintaining a hashtable that maps hashes to copies of nodes, effectively maintaing a flattened copy of the BT.

Currently, our implementation takes inspiration from the design principles of Merkle trees, if not going so far as to implement them perfectly. Applied widely across networked systems (e.g. a team of autonomous robots), a BT that enables effcient delta-communication should result in significant gains due to decreased contention for bandwidth; think: smaller message sizes, decreasing energy used in excessive wireless transmissions, improved overall efficiency.

### Details on Smaller Message Sizes
* The average case would be small hashes to confirm what parts of the tree need to be communicated
* Deltas of trees would be sent instead of an entire tree, except when the entire tree actually did change
* Eventual convergence to consistency is achieved, i.e. all robots have the latest version of the tree instance



![image](https://github.com/mips171/beetea/assets/18670565/67b53178-6d5c-4b8d-99af-9bd55c6bf168)

## Production Use

Used by [TaskBranch](https://github.com/mips171/taskbranch), an experimental behavior-tree-oriented system administration tool.
