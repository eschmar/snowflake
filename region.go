package snowflake

// Continents, from largest to smallest.
// Fly.io regions extracted from https://fly.io/docs/reference/regions/
// TODO: Support more region codes.
var continents = [][]string{
	// Asia
	{"bom", "hkg", "nrt", "sin"},
	// Africa
	{"jnb"},
	// North America
	{"atl", "bos", "den", "dfw", "ewr", "iad", "lax", "mia", "ord", "phx", "sea", "sjc", "yul", "yyz"},
	// South America
	{"bog", "eze", "gdl", "gig", "gru", "qro", "scl"},
	// Antarctica
	{},
	// Europe
	{"ams", "arn", "cdg", "fra", "lhr", "mad", "otp", "waw"},
	// Australia / Oceania
	{"syd"},
}

func getContinentCode(region string) int64 {
	for i := 0; i < len(continents); i++ {
		for j := range continents[i] {
			if continents[i][j] == region {
				return int64(i)
			}
		}
	}

	return -1
}
