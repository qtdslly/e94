package apimodel

type Upgrade struct {
	UpgradeVersion string `json:"upgrade_version"`
	TargetVersion  string `json:"target_version"`
	ShowUpgrade    bool   `json:"show_upgrade"`
	ForceUpgrade   bool   `json:"force_upgrade"`
	CheckUpgrade   bool   `json:"check_upgrade"`
	UpgradeTip     string `json:"upgrade_tip"`
	UpgradeUrl     string `json:"upgrade_url"`
}
