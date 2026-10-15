package main

// Merge the per-domain concept files into the copy the web app imports
// (`web/next-app/data/concepts.json`, typed by `web/next-app/lib/conceptData.ts`).
//
// `data/concepts/enrichment.json` is skipped on purpose. It is an overlay keyed by
// id -- encompasses, key_prerequisites, interference_group, variants -- and the
// loader applies it to the domain nodes rather than treating its entries as nodes,
// so including it would put ten id-only records into a file the browser reads as
// the graph. Before this skip the merge could never have been re-run without
// adding them; the copy in the tree was a snapshot taken before the file existed
// and it had quietly drifted 29 `grading_type` corrections behind the corpus.
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// overlayFiles are id-keyed overlays, not concept lists.
var overlayFiles = map[string]bool{"enrichment.json": true}

func main() {
	entries, err := os.ReadDir("data/concepts")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir: %v\n", err)
		os.Exit(1)
	}
	var all []any
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || overlayFiles[e.Name()] {
			continue
		}
		data, err := os.ReadFile(filepath.Join("data/concepts", e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", e.Name(), err)
			os.Exit(1)
		}
		var batch []any
		if err := json.Unmarshal(data, &batch); err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", e.Name(), err)
			os.Exit(1)
		}
		all = append(all, batch...)
	}
	out, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile("web/next-app/data/concepts.json", out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK: %d concepts merged into web/next-app/data/concepts.json\n", len(all))
}
