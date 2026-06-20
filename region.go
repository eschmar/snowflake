package snowflake

// Continents, from largest to smallest.
// Region codes are sourced from multiple cloud providers:
//   - Fly.io:    https://fly.io/docs/reference/regions/
//   - AWS RDS:   https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.md
//   - Scaleway:  https://www.scaleway.com/en/docs/account/reference-content/products-availability/
var continents = [][]string{
	// Asia
	{
		// Fly.io
		"bom", "hkg", "nrt", "sin",
		// AWS RDS
		"ap-east-1", "ap-east-2", "ap-northeast-1", "ap-northeast-2", "ap-northeast-3",
		"ap-south-1", "ap-south-2", "ap-southeast-1", "ap-southeast-3", "ap-southeast-5",
		"ap-southeast-7", "il-central-1", "me-central-1", "me-south-1",
	},
	// Africa
	{
		// Fly.io
		"jnb",
		// AWS RDS
		"af-south-1",
	},
	// North America
	{
		// Fly.io
		"atl", "bos", "den", "dfw", "ewr", "iad", "lax", "mia", "ord", "phx", "sea", "sjc", "yul", "yyz",
		// AWS RDS
		"ca-central-1", "ca-west-1", "mx-central-1", "us-east-1", "us-east-2",
		"us-gov-east-1", "us-gov-west-1", "us-west-1", "us-west-2",
	},
	// South America
	{
		// Fly.io
		"bog", "eze", "gdl", "gig", "gru", "qro", "scl",
		// AWS RDS
		"sa-east-1",
	},
	// Antarctica
	{},
	// Europe
	{
		// Fly.io
		"ams", "arn", "cdg", "fra", "lhr", "mad", "otp", "waw",
		// AWS RDS
		"eu-central-1", "eu-central-2", "eu-north-1", "eu-south-1", "eu-south-2",
		"eu-west-1", "eu-west-2", "eu-west-3",
		// Scaleway
		"fr-par", "par", "nl-ams", "ams", "pl-waw", "waw",
	},
	// Australia / Oceania
	{
		// Fly.io
		"syd",
		// AWS RDS
		"ap-southeast-2", "ap-southeast-4", "ap-southeast-6",
	},
}

func GetContinentCode(region string) int64 {
	for i := 0; i < len(continents); i++ {
		for j := range continents[i] {
			if continents[i][j] == region {
				return int64(i)
			}
		}
	}

	return -1
}
