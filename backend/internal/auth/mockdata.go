package auth

type MockUser struct {
	Email       string
	VPN         bool
	Production  bool
	ConfigTool  bool
	ProfileName string
	ProfileARN  string
}

var MockUsers = []MockUser{
	{
		Email:       "user1@companya.com",
		VPN:         true,
		Production:  true,
		ConfigTool:  true,
		ProfileName: "prod",
		ProfileARN:  "arn:aws:iam::123456789012:user/user1",
	},
	{
		Email:       "user2@companyb.com",
		VPN:         true,
		Production:  false,
		ConfigTool:  false,
		ProfileName: "dev",
		ProfileARN:  "arn:aws:iam::123456789013:user/user2",
	},
	// Add more test cases as needed
}
