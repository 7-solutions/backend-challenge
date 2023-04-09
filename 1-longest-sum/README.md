# maxSumPath
The maxSumPath function calculates the maximum sum path from the root node to any leaf node in a given triangle represented as a list of lists.

## Function Signature

<code>func maxSumPath(triangle [][]int) int</code>
Input
triangle: A non-empty list of lists of integers, where each row has one more integer than the row above it. The triangle represents a binary tree where each element has two children in the row below, except for the last row, which represents leaf nodes.
Output
Returns an integer representing the maximum sum of a path from the root node to any leaf node in the given triangle.
Algorithm
The function uses a Depth-First Search (DFS) approach to traverse the triangle and memoization to store results of previously computed subproblems.

Create a memoization map to store the maximum sum for each node.
Recursively traverse the triangle using DFS, starting from the root node.
For each node, calculate the maximum sum of its two child nodes and add the current node's value.
Store the calculated sum in the memoization map.
Return the maximum sum for the root node.

## Time Complexity
The time complexity of the maxSumPath function is O(N^2), where N is the number of rows in the triangle. This is because each element in the triangle is visited once, and there are N * (N + 1) / 2 elements in a triangle.