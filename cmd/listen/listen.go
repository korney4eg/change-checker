package listen

type Command struct {
	// Period        string `short:"p" long:"period" required:"true" choice:"any" choice:"day" choice:"week" choice:"month"`
	Config string `short:"c" long:"config" required:"true"  description:"path to configuration file"`
	// Day           string `short:"o" long:"only-day" required:"false" description:"Get statistics only for provided date. Example '01.02.2020'"`
	// SplitPerYear  bool   `short:"y" long:"year-split" required:"true" description:"Will split files by year"`
	// SplitPerMonth bool   `short:"m" long:"month-split" required:"true" description:"Will split files by month"`
}

func (c *Command) Execute(_ []string) error {
	// var err error
	return nil
}
