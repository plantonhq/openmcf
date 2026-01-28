package module

var vars = struct {
	HelmChartName      string
	HelmChartRepoUrl   string
	HelmChartVersion   string
	OpenBaoPort        int
	OpenBaoClusterPort int
}{
	HelmChartName:      "openbao",
	HelmChartRepoUrl:   "https://openbao.github.io/openbao-helm",
	HelmChartVersion:   "0.23.3",
	OpenBaoPort:        8200,
	OpenBaoClusterPort: 8201,
}
