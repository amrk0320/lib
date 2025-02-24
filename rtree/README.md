# RTree

Multidimensional Information Index

Example: “Find all museums within 2 km of current location
Input: search rectangle
Output: List of tuples contained in the search rectangle (list of locations in the range)

Tree structure for multidimensional indexing + free handling of short boundary ranges + handling of non-point data
 
Tree indexes can handle indexes in more than 2 dimensions, but they are sorted by tuples and cannot express spatial locality
Example: In 2D, the neighborhood of [0,0] cannot be expressed as a neighborhood of 4 points even if they are [-1,0], [1,0], [0,1], [0,-1]. ❌[0,1], [0,-1] are neighborhoods.

# Computational complexity

For a number of data n and a size m in a node
Depth: logm(n)- 1
Approximate case search: O(logm(n)) 
The worst-case computational complexity is not guaranteed.


This is the case when traversing down the tree from the root, or when traversing multiple sub-trees. However, when building the index, the search algorithm maintains the tree in such a way that it only traverses the neighborhood of a node, eliminating irrelevant regions. In other words, although it is not guaranteed as an algorithm, the tree maintenance does not worsen the computational complexity, and the PriorityRtree, an improved version of R-tree, guarantees the worst-case execution time.

