# fusion - collection of data structures

## Sparse Sets

Sparse sets are designed to save memory and improve performance when dealing with large datasets
that are mostly empty. Instead of allocating space for every possible element, they only store
the elements that are present.

Sparse sets often use some form of indexing to quickly access the elements that are present.
This can involve mapping the actual values to indices in a more compact representation.

They can be used in various applications, such as sparse matrices, where most of the elements are zero,
or in scenarios like managing a large number of objects in a game where only a few are active
at any given time.

## Stack

A stack is a linear data structure that follows the Last In, First Out (LIFO) principle.
This means that the last element added to the stack is the first one to be removed.
You can think of it like a stack of plates: you add plates to the top and remove them from the top.

Key Operations:
- Push: Add an element to the top of the stack.
- Pop: Remove the element from the top of the stack.
- Peek: Retrieve the element from the top of the stack without removing it.
- Get: Retrieve the element from the given position of the stack without removing it. Slow

## Tree (only max heap for now)

A heap is a special type of binary tree that satisfies the heap property.

- Max Heap: In a max heap, for any given node, the value of that node is greater than or equal to the values of its children.
  This means the largest value is at the root of the tree.

## Collection

A collection is a special data structure designed to avoid excessive memory allocation when adding
a large number of values to a slice in situations where you do not know in advance the amount of data
that needs to be stored.

This structure "cuts" the data into segments of equal length (slices), thus avoiding the complete
allocation of new memory when overflowing and simply adding a new segment to the chain.