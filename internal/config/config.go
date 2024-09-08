package config

// Default is a config instance.
var (
	Default Config // nolint:gochecknoglobals // config must be global

	Date    string
	Version string
	Build   string
)

type Config struct {
	LogLevel string `mapstructure:"log_level"`

	Gin struct {
		Mode string `mapstructure:"mode"`
	} `mapstructure:"gin"`

	Database struct {
		Log bool `mapstructure:"log"`
	} `mapstructure:"database"`

	App struct {
		AssetsPath string `mapstructure:"assets_path"`
	} `mapstructure:"app"`

	Ethereum struct {
		ChainID                    int    `mapstructure:"chain_id"`
		SubscriptionUrl            string `mapstructure:"subscription_url"`
		ContractAddress            string `mapstructure:"contract_address"`
		ContractABILocation        string `mapstructure:"contract_abi_location"`
		UpdatedContractABILocation string `mapstructure:"updated_contract_abi_location"`
		ContractOwnerPrivateKey    string `mapstructure:"contract_owner_private_key"`
		TestAccountPrivateKey      string `mapstructure:"test_account_private_key"`
	} `mapstructure:"ethereum"`

	Notification struct {
		EmailFrom  string `mapstructure:"email_from"`
		EmailTo    string `mapstructure:"email_to"`
		SMTPDomain string `mapstructure:"smtp_domain"`
		SMTPPort   int    `mapstructure:"smtp_port"`
		SMTPUser   string `mapstructure:"smtp_user"`
		SMTPKey    string `mapstructure:"smtp_key"`
	} `mapstructure:"notification"`
}
