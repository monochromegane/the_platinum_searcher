package the_platinum_searcher

import "os"

func ExampleFormatPrinterFileWithMatch() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.FilesWithMatches = true

	match := match{path: "filename"}
	match.add(1, 0, "go test", true)

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// filename
}

func ExampleFormatPrinterCount() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.Count = true

	match := match{path: "filename"}
	match.add(1, 0, "go test", true)

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// filename:1
}

func ExampleFormatEnableGroup() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.EnableGroup = true

	match := match{path: "filename"}
	match.add(1, 0, "before", false) // before
	match.add(2, 0, "go test", true) // no column
	match.add(3, 0, "after", false)  // after

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// filename
	// 1-before
	// 2:go test
	// 3-after
}

func ExampleFormatEnableGroupWithColumn() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.EnableGroup = true

	match := match{path: "filename"}
	match.add(1, 0, "before", false) // before
	match.add(2, 1, "go test", true) // no column
	match.add(3, 0, "after", false)  // after

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// filename
	// 1-before
	// 2:1:go test
	// 3-after
}

func ExampleFormatNoGroup() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.EnableGroup = false

	match := match{path: "filename"}
	match.add(1, 0, "before", false) // before
	match.add(2, 0, "go test", true) // no column
	match.add(3, 0, "after", false)  // after

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// filename:1-before
	// filename:2:go test
	// filename:3-after
}

func ExampleFormatNoGroupWithColumn() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.EnableGroup = false

	match := match{path: "filename"}
	match.add(1, 0, "before", false) // before
	match.add(2, 1, "go test", true) // no column
	match.add(3, 0, "after", false)  // after

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// filename:1-before
	// filename:2:1:go test
	// filename:3-after
}

func ExampleFormatMatchLine() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.SearchOption.SearchStream = true

	match := match{path: "/dev/stdin"}
	match.add(1, 0, "before", false) // before
	match.add(2, 0, "go test", true) // no column
	match.add(3, 0, "after", false)  // after

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// before
	// go test
	// after
}

func ExampleFormatMatchLineWithColumn() {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.SearchOption.SearchStream = true

	match := match{path: "/dev/stdin"}
	match.add(1, 0, "before", false) // before
	match.add(2, 1, "go test", true) // no column
	match.add(3, 0, "after", false)  // after

	pattern, _ := newPattern("go", opts)
	p := newFormatPrinter(pattern, os.Stdout, opts)
	p.print(match)

	// Output:
	// before
	// 1:go test
	// after
}
