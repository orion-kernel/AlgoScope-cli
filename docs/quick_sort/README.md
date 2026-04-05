# QUICK SORT

> "A highly efficient sorting algorithm and is based on partitioning of array of data into smaller arrays."

---

## 󱫐 HOW IT WORKS
1. Pick an element, called a **pivot**, from the array.
2. **Partitioning**: Reorder the array so that all elements with values less than the pivot come before the pivot, while all elements with values greater than the pivot come after it.
3. Recursively apply the above steps to the sub-array of elements with smaller values and separately to the sub-array of elements with greater values.

## 󰓡 COMPLEXITY
- **Worst Case:** O(n²) (usually when pivot is smallest/largest)
- **Average Case:** O(n log n)
- **Best Case:** O(n log n)
- **Space Complexity:** O(log n)

## 󰄬 CHARACTERISTICS
- **Stable:** No
- **In-place:** Yes
- **Algorithm:** Divide and Conquer

---
*AlgoScope Visualization Engine v1.0*
