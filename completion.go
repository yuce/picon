package picon

import (
	"github.com/chzyer/readline"
)

func consoleCompleter(console *Console) *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem(":exit"),
		readline.PcItem(":connect", readline.PcItemDynamic(console.listConnections())),
		readline.PcItem(":use", readline.PcItemDynamic(console.listIndexes())),
		readline.PcItem(":ensure",
			readline.PcItem("index"),
			readline.PcItem("frame")),
		readline.PcItem(":create",
			readline.PcItem("index"),
			readline.PcItem("frame")),
		readline.PcItem(":delete",
			readline.PcItem("index"),
			readline.PcItem("frame")),
		readline.PcItem(":schema"),
		readline.PcItem(":save"),
		readline.PcItem(":session"),

		// PQL commands
		readline.PcItem("Bitmap("),
		readline.PcItem("ClearBit("),
		readline.PcItem("Count("),
		readline.PcItem("Difference("),
		readline.PcItem("Intersect("),
		readline.PcItem("Range("),
		readline.PcItem("SetBit("),
		readline.PcItem("SetColumnAttrs("),
		readline.PcItem("SetRowAttrs("),
		readline.PcItem("TopN("),
		readline.PcItem("Union("),
	)
}
