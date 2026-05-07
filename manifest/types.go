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
	AppID       string   `yaml:"appid,omitempty" json:"appid,omitempty"`
	Title       string   `yaml:"title" json:"title"`
	Version     string   `yaml:"version" json:"version"`
	Categories  []string `yaml:"categories,omitempty" json:"categories,omitempty"`
	Rating      float32  `yaml:"rating,omitempty" json:"rating,omitempty"`
	Target      string   `yaml:"target,omitempty" json:"target,omitempty"`
	Type        string   `yaml:"type,omitempty" json:"type,omitempty"`
}

type AppConfiguration struct {
	ConfigVersion string                  `yaml:"olaresManifest.version" json:"olaresManifest.version"`
	ConfigType    string                  `yaml:"olaresManifest.type" json:"olaresManifest.type"`
	APIVersion    string                  `yaml:"apiVersion" json:"apiVersion"`
	Metadata      AppMetaData             `yaml:"metadata" json:"metadata"`
	Entrances     []v1alpha1.Entrance     `yaml:"entrances" json:"entrances"`
	Ports         []v1alpha1.ServicePort  `yaml:"ports,omitempty" json:"ports,omitempty"`
	TailScale     v1alpha1.TailScale      `yaml:"tailscale,omitempty" json:"tailscale,omitempty"`
	Spec          AppSpec                 `yaml:"spec" json:"spec"`
	Permission    Permission              `yaml:"permission,omitempty" json:"permission,omitempty" description:"app permission request"`
	Middleware    *Middleware             `yaml:"middleware,omitempty" json:"middleware,omitempty" description:"app middleware request"`
	Options       Options                 `yaml:"options,omitempty" json:"options,omitempty" description:"app options"`
	Provider      []Provider              `yaml:"provider,omitempty" json:"provider,omitempty" description:"app provider information"`
	Envs          []sysv1alpha1.AppEnvVar `yaml:"envs,omitempty" json:"envs,omitempty"`

	// Only for v2 c/s apps to share the api to other cluster scope apps
	SharedEntrances []v1alpha1.Entrance `yaml:"sharedEntrances,omitempty" json:"sharedEntrances,omitempty"`
}

type AppSpec struct {
	Namespace           string         `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	OnlyAdmin           bool           `yaml:"onlyAdmin,omitempty" json:"onlyAdmin,omitempty"`
	VersionName         string         `yaml:"versionName" json:"versionName"`
	FullDescription     string         `yaml:"fullDescription,omitempty" json:"fullDescription,omitempty"`
	UpgradeDescription  string         `yaml:"upgradeDescription,omitempty" json:"upgradeDescription,omitempty"`
	PromoteImage        []string       `yaml:"promoteImage,omitempty" json:"promoteImage,omitempty"`
	PromoteVideo        string         `yaml:"promoteVideo,omitempty" json:"promoteVideo,omitempty"`
	SubCategory         string         `yaml:"subCategory,omitempty" json:"subCategory,omitempty"`
	Developer           string         `yaml:"developer,omitempty" json:"developer,omitempty"`
	RequiredMemory      string         `yaml:"requiredMemory,omitempty" json:"requiredMemory,omitempty"`
	RequiredDisk        string         `yaml:"requiredDisk,omitempty" json:"requiredDisk,omitempty"`
	RequiredGPU         string         `yaml:"requiredGpu,omitempty" json:"requiredGpu,omitempty"`
	RequiredCPU         string         `yaml:"requiredCpu,omitempty" json:"requiredCpu,omitempty"`
	LimitedMemory       string         `yaml:"limitedMemory,omitempty" json:"limitedMemory,omitempty"`
	LimitedDisk         string         `yaml:"limitedDisk,omitempty" json:"limitedDisk,omitempty"`
	LimitedGPU          string         `yaml:"limitedGpu,omitempty" json:"limitedGpu,omitempty"`
	LimitedCPU          string         `yaml:"limitedCpu,omitempty" json:"limitedCpu,omitempty"`
	SupportClient       SupportClient  `yaml:"supportClient,omitempty" json:"supportClient,omitempty"`
	RunAsUser           bool           `yaml:"runAsUser" json:"runAsUser"`
	RunAsInternal       bool           `yaml:"runAsInternal" json:"runAsInternal"`
	PodGPUConsumePolicy string         `yaml:"podGpuConsumePolicy,omitempty" json:"podGpuConsumePolicy,omitempty"`
	SubCharts           []Chart        `yaml:"subCharts,omitempty" json:"subCharts,omitempty"`
	Hardware            Hardware       `yaml:"hardware,omitempty" json:"hardware,omitempty"`
	SupportedGpu        []any          `yaml:"supportedGpu,omitempty" json:"supportedGpu,omitempty"`
	Resources           []ResourceMode `yaml:"resources,omitempty" json:"resources,omitempty"`
	SupportArch         []string       `yaml:"supportArch,omitempty" json:"supportArch,omitempty"`
	Website             string         `yaml:"website,omitempty" json:"website,omitempty"`
	SourceCode          string         `yaml:"sourceCode,omitempty" json:"sourceCode,omitempty"`
	Submitter           string         `yaml:"submitter,omitempty" json:"submitter,omitempty"`
	Locale              []string       `yaml:"locale,omitempty" json:"locale,omitempty"`
	Doc                 string         `yaml:"doc,omitempty" json:"doc,omitempty"`
	License             []struct {
		Text string `yaml:"text,omitempty" json:"text,omitempty"`
		URL  string `yaml:"url,omitempty" json:"url,omitempty"`
	} `yaml:"license,omitempty" json:"license,omitempty"`
}

type Hardware struct {
	Cpu CpuConfig `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Gpu GpuConfig `yaml:"gpu,omitempty" json:"gpu,omitempty"`
}

type CpuConfig struct {
	Vendor string `yaml:"vendor,omitempty" json:"vendor,omitempty"`
	Arch   string `yaml:"arch,omitempty" json:"arch,omitempty"`
}
type GpuConfig struct {
	Vendor string   `yaml:"vendor,omitempty" json:"vendor,omitempty"`
	Arch   []string `yaml:"arch,omitempty" json:"arch,omitempty"`
	// SingleMemory specifies the minimum memory size required for a single GPU
	SingleMemory string `yaml:"singleMemory,omitempty" json:"singleMemory,omitempty"`
	// TotalMemory specifies the total GPU memory required across all GPUs within one node
	TotalMemory string `yaml:"totalMemory,omitempty" json:"totalMemory,omitempty"`
}

type SupportClient struct {
	Edge    string `yaml:"edge,omitempty" json:"edge,omitempty"`
	Android string `yaml:"android,omitempty" json:"android,omitempty"`
	Ios     string `yaml:"ios,omitempty" json:"ios,omitempty"`
	Windows string `yaml:"windows,omitempty" json:"windows,omitempty"`
	Mac     string `yaml:"mac,omitempty" json:"mac,omitempty"`
	Linux   string `yaml:"linux,omitempty" json:"linux,omitempty"`
}

type Permission struct {
	AppData        bool                 `yaml:"appData,omitempty" json:"appData,omitempty"  description:"app data permission for writing"`
	AppCache       bool                 `yaml:"appCache,omitempty" json:"appCache,omitempty"`
	UserData       []string             `yaml:"userData,omitempty" json:"userData,omitempty"`
	Provider       []ProviderPermission `yaml:"provider,omitempty" json:"provider,omitempty"  description:"system shared data permission for accessing"`
	ServiceAccount *string              `yaml:"serviceAccount,omitempty" json:"serviceAccount,omitempty" description:"service account for app permission"`
}

type ProviderPermission struct {
	AppName      string                 `yaml:"appName,omitempty" json:"appName,omitempty"`
	Namespace    string                 `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	ProviderName string                 `yaml:"providerName,omitempty" json:"providerName,omitempty"`
	PodSelectors []metav1.LabelSelector `yaml:"podSelectors,omitempty" json:"podSelectors,omitempty"`
}

type Policy struct {
	EntranceName string `yaml:"entranceName,omitempty" json:"entranceName,omitempty"`
	Description  string `yaml:"description,omitempty" json:"description,omitempty" description:"the description of the policy"`
	URIRegex     string `yaml:"uriRegex,omitempty" json:"uriRegex,omitempty" description:"uri regular expression"`
	Level        string `yaml:"level,omitempty" json:"level,omitempty"`
	OneTime      bool   `yaml:"oneTime" json:"oneTime"`
	Duration     string `yaml:"validDuration,omitempty" json:"validDuration,omitempty"`
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
	MobileSupported      bool                     `yaml:"mobileSupported,omitempty" json:"mobileSupported,omitempty"`
	Policies             []Policy                 `yaml:"policies,omitempty" json:"policies,omitempty"`
	ResetCookie          ResetCookie              `yaml:"resetCookie,omitempty" json:"resetCookie,omitempty"`
	Dependencies         []Dependency             `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Conflicts            []Conflict               `yaml:"conflicts,omitempty" json:"conflicts,omitempty"`
	AppScope             AppScope                 `yaml:"appScope,omitempty" json:"appScope,omitempty"`
	WsConfig             WsConfig                 `yaml:"websocket,omitempty" json:"websocket,omitempty"`
	Upload               Upload                   `yaml:"upload,omitempty" json:"upload,omitempty"`
	SyncProvider         []map[string]interface{} `yaml:"syncProvider,omitempty" json:"syncProvider,omitempty"`
	OIDC                 OIDC                     `yaml:"oidc,omitempty" json:"oidc,omitempty"`
	ApiTimeout           *int64                   `yaml:"apiTimeout,omitempty" json:"apiTimeout,omitempty"`
	AllowedOutboundPorts []int                    `yaml:"allowedOutboundPorts,omitempty" json:"AllowedOutboundPorts,omitempty"`
	Images               []string                 `yaml:"images,omitempty" json:"images,omitempty"`
	AllowMultipleInstall bool                     `yaml:"allowMultipleInstall,omitempty" json:"allowMultipleInstall,omitempty"`
}

type ResetCookie struct {
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
}

type AppScope struct {
	ClusterScoped bool     `yaml:"clusterScoped,omitempty" json:"clusterScoped,omitempty"`
	AppRef        []string `yaml:"appRef,omitempty" json:"appRef,omitempty"`
	SystemService bool     `yaml:"systemService,omitempty" json:"systemService,omitempty"`
}

type WsConfig struct {
	Port int    `yaml:"port,omitempty" json:"port,omitempty"`
	URL  string `yaml:"url,omitempty" json:"url,omitempty"`
}

type Upload struct {
	FileType    []string `yaml:"fileType,omitempty" json:"fileType,omitempty"`
	Dest        string   `yaml:"dest,omitempty" json:"dest,omitempty"`
	LimitedSize int      `yaml:"limitedSize,omitempty" json:"limitedSize,omitempty"`
}

type OIDC struct {
	Enabled      bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	RedirectUri  string `yaml:"redirectUri,omitempty" json:"redirectUri,omitempty"`
	EntranceName string `yaml:"entranceName,omitempty" json:"entranceName,omitempty"`
}

type Chart struct {
	Name   string `yaml:"name,omitempty" json:"name,omitempty"`
	Shared bool   `yaml:"shared,omitempty" json:"shared,omitempty"`
}

type Provider struct {
	Name     string   `yaml:"name,omitempty" json:"name,omitempty"`
	Entrance string   `yaml:"entrance,omitempty" json:"entrance,omitempty"`
	Paths    []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	Verbs    []string `yaml:"verbs,omitempty" json:"verbs,omitempty"`
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
	ResourceRequirement `yaml:",inline" json:",inline"`
}

// Middleware describe middleware config.
type Middleware struct {
	Postgres      *PostgresConfig      `yaml:"postgres,omitempty" json:"postgres,omitempty"`
	Redis         *RedisConfig         `yaml:"redis,omitempty" json:"redis,omitempty"`
	MongoDB       *MongodbConfig       `yaml:"mongodb,omitempty" json:"mongodb,omitempty"`
	Nats          *NatsConfig          `yaml:"nats,omitempty" json:"nats,omitempty"`
	Minio         *MinioConfig         `yaml:"minio,omitempty" json:"minio,omitempty"`
	RabbitMQ      *RabbitMQConfig      `yaml:"rabbitmq,omitempty" json:"rabbitmq,omitempty"`
	Elasticsearch *ElasticsearchConfig `yaml:"elasticsearch,omitempty" json:"elasticsearch,omitempty"`
	MariaDB       *MariaDBConfig       `yaml:"mariadb,omitempty" json:"mariadb,omitempty"`
	MySQL         *MySQLConfig         `yaml:"mysql,omitempty" json:"mysql,omitempty"`
	Argo          *ArgoConfig          `yaml:"argo,omitempty" json:"argo,omitempty"`
	ClickHouse    *ClickHouseConfig    `yaml:"clickHouse,omitempty" json:"clickHouse,omitempty"`
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
	Name string `yaml:"name" json:"name"`
}

// RedisConfig contains fields for redis config.
type RedisConfig struct {
	Password  string `yaml:"password,omitempty" json:"password,omitempty"`
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
