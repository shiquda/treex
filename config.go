package main

type Config struct {
	DirCollapse    bool   `yaml:"dirCollapse"`    // 是否折叠目录（默认：false）
	ShowRootName   bool   `yaml:"showRootName"`   // 是否显示根目录名称（默认：true）
	MaxDepth       int    `yaml:"maxDepth"`       // 最大递归深度（默认：-1，表示无限）
	Dir            string `yaml:"dir"`            // 目标目录（默认：当前目录 "."）
	OutputFormat   string `yaml:"outputFormat"`   // 输出格式（默认："text"）
	ExcludeRuleStr string `yaml:"excludeRuleStr"` // 排除规则（默认：空）
	HideHidden     bool   `yaml:"hideHidden"`     // 是否隐藏隐藏文件（默认：true）
	DirsOnly       bool   `yaml:"dirsOnly"`       // 是否只显示目录（默认：false）
	UseGitIgnore   bool   `yaml:"useGitIgnore"`   // 是否使用.gitignore（默认：true）
	UseIcons       bool   `yaml:"useIcons"`       // 是否使用图标（默认：false）
	OutputFilePath string `yaml:"outputFilePath"`
}

// NewDefaultConfig 返回带默认值的配置
func NewDefaultConfig() *Config {
	return &Config{
		DirCollapse:    false,
		ShowRootName:   false,
		MaxDepth:       -1,
		Dir:            ".",
		OutputFormat:   "tree",
		ExcludeRuleStr: "",
		HideHidden:     true,
		DirsOnly:       false,
		UseGitIgnore:   false,
		UseIcons:       false,
		OutputFilePath: "",
	}
}
