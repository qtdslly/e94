package apimodel

<<<<<<< HEAD


type AppConfigUpgrade struct {
=======
type Upgrade struct {
>>>>>>> 4c7ea1426bca1ef3e9e2bde22b4eb03806127bdb
	UpgradeVersion string `json:"upgrade_version"`
	TargetVersion  string `json:"target_version"`
	ShowUpgrade    bool   `json:"show_upgrade"`
	ForceUpgrade   bool   `json:"force_upgrade"`
	CheckUpgrade   bool   `json:"check_upgrade"`
	UpgradeTip     string `json:"upgrade_tip"`
	UpgradeUrl     string `json:"upgrade_url"`
}
