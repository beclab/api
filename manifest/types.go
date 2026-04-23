package manifest

import (
	"github.com/beclab/api/api/app.bytetrade.io/v1alpha1"
	sysv1alpha1 "github.com/beclab/api/api/sys.bytetrade.io/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppMetaData struct {
	Name        string   `yaml:"name" json:"name"`
	Icon        string   `yaml:"icon" json:"icon"`
	Description string   `yaml:"description" json:"description"`
	AppID       string   `yaml:"appid" json:"appid"`
	Title       string   `yaml:"title" json:"title"`
	Version     string   `yaml:"version" json:"version"`
	Categories  []string `yaml:"categories" json:"categories"`
	Rating      float32  `yaml:"rating" json:"rating"`
	Target      string   `yaml:"target" json:"target"`
	Type        string   `yaml:"type" json:"type"`
}

type AppConfiguration struct {
	ConfigVersion string                  `yaml:"olaresManifest.version" json:"olaresManifest.version"`
	ConfigType    string                  `yaml:"olaresManifest.type" json:"olaresManifest.type"`
	APIVersion    string                  `yaml:"apiVersion" json:"apiVersion"`
	Metadata      AppMetaData             `yaml:"metadata" json:"metadata"`
	Entrances     []v1alpha1.Entrance     `yaml:"entrances" json:"entrances"`
	Ports         []v1alpha1.ServicePort  `yaml:"ports" json:"ports"`
	TailScale     v1alpha1.TailScale      `yaml:"tailscale" json:"tailscale"`
	Spec          AppSpec                 `yaml:"spec" json:"spec"`
	Permission    Permission              `yaml:"permission" json:"permission" description:"app permission request"`
	Middleware    *Middleware             `yaml:"middleware" json:"middleware" description:"app middleware request"`
	Options       Options                 `yaml:"options" json:"options" description:"app options"`
	Provider      []Provider              `yaml:"provider,omitempty" json:"provider,omitempty" description:"app provider information"`
	Envs          []sysv1alpha1.AppEnvVar `yaml:"envs,omitempty" json:"envs,omitempty"`

	// Only for v2 c/s apps to share the api to other cluster scope apps
	SharedEntrances []v1alpha1.Entrance `yaml:"sharedEntrances,omitempty" json:"sharedEntrances,omitempty"`
}

type AppSpec struct {
	Namespace           string         `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	OnlyAdmin           bool           `yaml:"onlyAdmin,omitempty" json:"onlyAdmin,omitempty"`
	VersionName         string         `yaml:"versionName" json:"versionName"`
	FullDescription     string         `yaml:"fullDescription" json:"fullDescription"`
	UpgradeDescription  string         `yaml:"upgradeDescription" json:"upgradeDescription"`
	PromoteImage        []string       `yaml:"promoteImage" json:"promoteImage"`
	PromoteVideo        string         `yaml:"promoteVideo" json:"promoteVideo"`
	SubCategory         string         `yaml:"subCategory" json:"subCategory"`
	Developer           string         `yaml:"developer" json:"developer"`
	RequiredMemory      string         `yaml:"requiredMemory" json:"requiredMemory"`
	RequiredDisk        string         `yaml:"requiredDisk" json:"requiredDisk"`
	RequiredGPU         string         `yaml:"requiredGpu" json:"requiredGpu"`
	RequiredCPU         string         `yaml:"requiredCpu" json:"requiredCpu"`
	LimitedMemory       string         `yaml:"limitedMemory" json:"limitedMemory"`
	LimitedDisk         string         `yaml:"limitedDisk" json:"limitedDisk"`
	LimitedGPU          string         `yaml:"limitedGpu" json:"limitedGpu"`
	LimitedCPU          string         `yaml:"limitedCpu" json:"limitedCpu"`
	SupportClient       SupportClient  `yaml:"supportClient" json:"supportClient"`
	RunAsUser           bool           `yaml:"runAsUser" json:"runAsUser"`
	RunAsInternal       bool           `yaml:"runAsInternal" json:"runAsInternal"`
	PodGPUConsumePolicy string         `yaml:"podGpuConsumePolicy" json:"podGpuConsumePolicy"`
	SubCharts           []Chart        `yaml:"subCharts" json:"subCharts"`
	Hardware            Hardware       `yaml:"hardware" json:"hardware"`
	SupportedGpu        []any          `yaml:"supportedGpu,omitempty" json:"supportedGpu,omitempty"`
	Resources           []ResourceMode `yaml:"resources,omitempty" json:"resources,omitempty"`
	SupportArch         []string       `yaml:"supportArch,omitempty" json:"supportArch,omitempty"`
	Website             string         `yaml:"website,omitempty" json:"website,omitempty"`
	SourceCode          string         `yaml:"sourceCode,omitempty" json:"sourceCode,omitempty"`
	Submitter           string         `yaml:"submitter" json:"submitter"`
	Locale              []string       `yaml:"locale" json:"locale"`
	Doc                 string         `yaml:"doc" json:"doc"`
	License             []struct {
		Text string `yaml:"text"`
		URL  string `yaml:"url"`
	} `yaml:"license"`
}

type Hardware struct {
	Cpu CpuConfig `yaml:"cpu" json:"cpu"`
	Gpu GpuConfig `yaml:"gpu" json:"gpu"`
}

type CpuConfig struct {
	Vendor string `yaml:"vendor" json:"vendor"`
	Arch   string `yaml:"arch" json:"arch"`
}
type GpuConfig struct {
	Vendor string   `yaml:"vendor" json:"vendor"`
	Arch   []string `yaml:"arch" json:"arch"`
	// SingleMemory specifies the minimum memory size required for a single GPU
	SingleMemory string `yaml:"singleMemory" json:"singleMemory"`
	// TotalMemory specifies the total GPU memory required across all GPUs within one node
	TotalMemory string `yaml:"totalMemory" json:"totalMemory"`
}

type SupportClient struct {
	Edge    string `yaml:"edge" json:"edge"`
	Android string `yaml:"android" json:"android"`
	Ios     string `yaml:"ios" json:"ios"`
	Windows string `yaml:"windows" json:"windows"`
	Mac     string `yaml:"mac" json:"mac"`
	Linux   string `yaml:"linux" json:"linux"`
}

type Permission struct {
	AppData        bool                 `yaml:"appData" json:"appData"  description:"app data permission for writing"`
	AppCache       bool                 `yaml:"appCache" json:"appCache"`
	UserData       []string             `yaml:"userData" json:"userData"`
	Provider       []ProviderPermission `yaml:"provider" json:"provider"  description:"system shared data permission for accessing"`
	ServiceAccount *string              `yaml:"serviceAccount,omitempty" json:"serviceAccount,omitempty" description:"service account for app permission"`
}

type ProviderPermission struct {
	AppName      string                 `yaml:"appName" json:"appName"`
	Namespace    string                 `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	ProviderName string                 `yaml:"providerName" json:"providerName"`
	PodSelectors []metav1.LabelSelector `yaml:"podSelectors" json:"podSelectors"`
}

type Policy struct {
	EntranceName string `yaml:"entranceName" json:"entranceName"`
	Description  string `yaml:"description" json:"description" description:"the description of the policy"`
	URIRegex     string `yaml:"uriRegex" json:"uriRegex" description:"uri regular expression"`
	Level        string `yaml:"level" json:"level"`
	OneTime      bool   `yaml:"oneTime" json:"oneTime"`
	Duration     string `yaml:"validDuration" json:"validDuration"`
}

type Dependency struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	// dependency type: system, application.
	Type      string `yaml:"type" json:"type"`
	Mandatory bool   `yaml:"mandatory" json:"mandatory"`
	SelfRely  bool   `yaml:"selfRely" json:"selfRely"`
}

type Conflict struct {
	Name string `yaml:"name" json:"name"`
	// conflict type: application
	Type string `yaml:"type" json:"type"`
}

type Options struct {
	MobileSupported      bool                     `yaml:"mobileSupported" json:"mobileSupported"`
	Policies             []Policy                 `yaml:"policies" json:"policies"`
	ResetCookie          ResetCookie              `yaml:"resetCookie" json:"resetCookie"`
	Dependencies         []Dependency             `yaml:"dependencies" json:"dependencies"`
	Conflicts            []Conflict               `yaml:"conflicts" json:"conflicts"`
	AppScope             AppScope                 `yaml:"appScope" json:"appScope"`
	WsConfig             WsConfig                 `yaml:"websocket" json:"websocket"`
	Upload               Upload                   `yaml:"upload" json:"upload"`
	SyncProvider         []map[string]interface{} `yaml:"syncProvider" json:"syncProvider"`
	OIDC                 OIDC                     `yaml:"oidc" json:"oidc"`
	ApiTimeout           *int64                   `yaml:"apiTimeout" json:"apiTimeout"`
	AllowedOutboundPorts []int                    `yaml:"allowedOutboundPorts" json:"AllowedOutboundPorts"`
	Images               []string                 `yaml:"images" json:"images"`
	AllowMultipleInstall bool                     `yaml:"allowMultipleInstall" json:"allowMultipleInstall"`
}

type ResetCookie struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

type AppScope struct {
	ClusterScoped bool     `yaml:"clusterScoped" json:"clusterScoped"`
	AppRef        []string `yaml:"appRef" json:"appRef"`
	SystemService bool     `yaml:"systemService" json:"systemService"`
}

type WsConfig struct {
	Port int    `yaml:"port" json:"port"`
	URL  string `yaml:"url" json:"url"`
}

type Upload struct {
	FileType    []string `yaml:"fileType" json:"fileType"`
	Dest        string   `yaml:"dest" json:"dest"`
	LimitedSize int      `yaml:"limitedSize" json:"limitedSize"`
}

type OIDC struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	RedirectUri  string `yaml:"redirectUri" json:"redirectUri"`
	EntranceName string `yaml:"entranceName" json:"entranceName"`
}

type Chart struct {
	Name   string `yaml:"name" json:"name"`
	Shared bool   `yaml:"shared" json:"shared"`
}

type Provider struct {
	Name     string   `yaml:"name" json:"name"`
	Entrance string   `yaml:"entrance" json:"entrance"`
	Paths    []string `yaml:"paths" json:"paths"`
	Verbs    []string `yaml:"verbs" json:"verbs"`
}

type SpecialResource struct {
	RequiredMemory *string `yaml:"requiredMemory,omitempty" json:"requiredMemory,omitempty"`
	RequiredDisk   *string `yaml:"requiredDisk,omitempty" json:"requiredDisk,omitempty"`
	RequiredGPU    *string `yaml:"requiredGpu,omitempty" json:"requiredGpu,omitempty"`
	RequiredCPU    *string `yaml:"requiredCpu,omitempty" json:"requiredCpu,omitempty"`
	LimitedMemory  *string `yaml:"limitedMemory,omitempty" json:"limitedMemory,omitempty"`
	LimitedDisk    *string `yaml:"limitedDisk,omitempty" json:"limitedDisk,omitempty"`
	LimitedGPU     *string `yaml:"limitedGPU,omitempty" json:"limitedGPU,omitempty"`
	LimitedCPU     *string `yaml:"limitedCPU,omitempty" json:"limitedCPU,omitempty"`
}

type ResourceRequirement struct {
	RequiredCPU    string `yaml:"requiredCpu,omitempty" json:"requiredCpu,omitempty"`
	LimitedCPU     string `yaml:"limitedCpu,omitempty" json:"limitedCpu,omitempty"`
	RequiredMemory string `yaml:"requiredMemory,omitempty" json:"requiredMemory,omitempty"`
	LimitedMemory  string `yaml:"limitedMemory,omitempty" json:"limitedMemory,omitempty"`
	RequiredDisk   string `yaml:"requiredDisk,omitempty" json:"requiredDisk,omitempty"`
	LimitedDisk    string `yaml:"limitedDisk,omitempty" json:"limitedDisk,omitempty"`
	RequiredGPU    string `yaml:"requiredGPUMemory,omitempty" json:"requiredGPUMemory,omitempty"`
	LimitedGPU     string `yaml:"limitedGPUMemory,omitempty" json:"limitedGPUMemory,omitempty"`
}

type ResourceMode struct {
	Mode                string `yaml:"mode" json:"mode"`
	ResourceRequirement `yaml:",inline"`
}

// Middleware describe middleware config.
type Middleware struct {
	Postgres      *PostgresConfig      `yaml:"postgres,omitempty"`
	Redis         *RedisConfig         `yaml:"redis,omitempty"`
	MongoDB       *MongodbConfig       `yaml:"mongodb,omitempty"`
	Nats          *NatsConfig          `yaml:"nats,omitempty"`
	Minio         *MinioConfig         `yaml:"minio,omitempty"`
	RabbitMQ      *RabbitMQConfig      `yaml:"rabbitmq,omitempty"`
	Elasticsearch *ElasticsearchConfig `yaml:"elasticsearch,omitempty"`
	MariaDB       *MariaDBConfig       `yaml:"mariadb,omitempty"`
	MySQL         *MySQLConfig         `yaml:"mysql,omitempty"`
	Argo          *ArgoConfig          `yaml:"argo,omitempty"`
	ClickHouse    *ClickHouseConfig    `yaml:"clickHouse,omitempty"`
}

// Database specify database name and if distributed.
type Database struct {
	Name        string   `yaml:"name" json:"name"`
	Extensions  []string `yaml:"extensions,omitempty" json:"extensions,omitempty"`
	Scripts     []string `yaml:"scripts,omitempty" json:"scripts,omitempty"`
	Distributed bool     `yaml:"distributed,omitempty" json:"distributed,omitempty"`
}

// PostgresConfig contains fields for postgresql config.
type PostgresConfig struct {
	Username  string     `yaml:"username" json:"username"`
	Password  string     `yaml:"password,omitempty" json:"password,omitempty"`
	Databases []Database `yaml:"databases" json:"databases"`
}

type ArgoConfig struct {
	Required bool `yaml:"required" json:"required"`
}

type MinioConfig struct {
	Username              string   `yaml:"username" json:"username"`
	Password              string   `yaml:"password" json:"password"`
	Buckets               []Bucket `yaml:"buckets" json:"buckets"`
	AllowNamespaceBuckets bool     `yaml:"allowNamespaceBuckets" json:"allowNamespaceBuckets"`
}

type Bucket struct {
	Name string `json:"name"`
}

type RabbitMQConfig struct {
	Username string  `yaml:"username" json:"username"`
	Password string  `yaml:"password" json:"password"`
	VHosts   []VHost `yaml:"vhosts" json:"vhosts"`
}

type VHost struct {
	Name string `json:"name"`
}

type ElasticsearchConfig struct {
	Username              string  `yaml:"username" json:"username"`
	Password              string  `yaml:"password" json:"password"`
	Indexes               []Index `yaml:"indexes" json:"indexes"`
	AllowNamespaceIndexes bool    `yaml:"allowNamespaceIndexes" json:"allowNamespaceIndexes"`
}

type Index struct {
	Name string `json:"name"`
}

// RedisConfig contains fields for redis config.
type RedisConfig struct {
	Password  string `yaml:"password,omitempty" json:"password"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

// MongodbConfig contains fields for mongodb config.
type MongodbConfig struct {
	Username  string     `yaml:"username" json:"username"`
	Password  string     `yaml:"password,omitempty" json:"password"`
	Databases []Database `yaml:"databases" json:"databases"`
}

// MariaDBConfig contains fields for mariadb config.
type MariaDBConfig struct {
	Username  string     `yaml:"username" json:"username"`
	Password  string     `yaml:"password,omitempty" json:"password"`
	Databases []Database `yaml:"databases" json:"databases"`
}

// MySQLConfig contains fields for mysql config.
type MySQLConfig struct {
	Username  string     `yaml:"username" json:"username"`
	Password  string     `yaml:"password,omitempty" json:"password"`
	Databases []Database `yaml:"databases" json:"databases"`
}

// ClickHouseConfig contains fields for clickhouse config.
type ClickHouseConfig struct {
	Username  string     `yaml:"username" json:"username"`
	Password  string     `yaml:"password,omitempty" json:"password"`
	Databases []Database `yaml:"databases" json:"databases"`
}

type NatsConfig struct {
	Username string    `yaml:"username" json:"username"`
	Password string    `yaml:"password,omitempty" json:"password,omitempty"`
	Subjects []Subject `yaml:"subjects" json:"subjects"`
	Refs     []Ref     `yaml:"refs" json:"refs"`
}

type Subject struct {
	Name string `yaml:"name" json:"name"`
	// Permissions indicates the permission that app can perform on this subject
	Permission PermissionNats   `yaml:"permission" json:"permission"`
	Export     []PermissionNats `yaml:"export" json:"export"`
}

type Export struct {
	AppName string `yaml:"appName" json:"appName"`
	Pub     string `yaml:"pub" json:"pub"`
	Sub     string `yaml:"sub" json:"sub"`
}

type Ref struct {
	AppName string `yaml:"appName" json:"appName"`
	// option for ref app in user-space-<>, user-system-<>, os-system
	AppNamespace string       `yaml:"appNamespace" json:"appNamespace"`
	Subjects     []RefSubject `yaml:"subjects" json:"subjects"`
}

type RefSubject struct {
	Name string   `yaml:"name" json:"name"`
	Perm []string `yaml:"perm" json:"perm"`
}

type PermissionNats struct {
	AppName string `yaml:"appName,omitempty" json:"appName,omitempty"`
	// default is deny
	Pub string `yaml:"pub" json:"pub"`
	Sub string `yaml:"sub" json:"sub"`
}
