package apimodel


type AppConfigUpgrade struct {
	UpgradeVersion string `json:"upgrade_version"`
	TargetVersion  string `json:"target_version"`
	ShowUpgrade    bool   `json:"show_upgrade"`
	ForceUpgrade   bool   `json:"force_upgrade"`
	CheckUpgrade   bool   `json:"check_upgrade"`
	UpgradeTip     string `json:"upgrade_tip"`
	UpgradeUrl     string `json:"upgrade_url"`
	Md5Value       string `json:"md5_value"`
}

type AppActivity struct {
	Channel        uint32 `json:"channel"`
	Account        string `json:"account"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Thumb          string `json:"thumb"`
}