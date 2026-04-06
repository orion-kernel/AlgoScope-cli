package algo

type Algorithm struct {
	Name       string
	Desc       string
	Complexity string
	Stability  string
	DocPath    string
}

func GetAlgorithms() []Algorithm {
	return []Algorithm{
		{"BUBBLE SORT", "Classic O(n²) sorting algorithm. Ideal for visualizing the basic concept of swapping and iterations.", "O(n²)", "STABLE", "docs/bubble_sort/README.md"},
		{"QUICK SORT", "Highly efficient O(n log n) divide-and-conquer algorithm. Selects a pivot to partition data.", "O(n log n)", "UNSTABLE", "docs/quick_sort/README.md"},
		{"MERGE SORT", "Reliable O(n log n) stable sort. Recursively divides array into halves and merges them back.", "O(n log n)", "STABLE", "docs/merge_sort/README.md"},
		{"EXIT ENGINE", "Terminate the AlgoScope visualization engine and return to host shell.", "N/A", "N/A", ""},
	}
}
