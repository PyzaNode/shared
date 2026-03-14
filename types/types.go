package types

import "time"

// Network represents a Minecraft network ecosystem (proxy groups + server groups).
type Network struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description,omitempty"`
	Enabled           bool                   `json:"enabled"`
	DefaultProxyGroup string                 `json:"default_proxy_group,omitempty"` // Proxy group ID for player join
	Cloudflare        *CloudflareConfig      `json:"cloudflare,omitempty"`
	CloudflareSRV     *CloudflareSRVSettings `json:"cloudflare_srv,omitempty"`     // Per-network SRV sync (zone, hostname, proxy group)
	Tags              map[string]string      `json:"tags,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// CloudflareConfig holds optional DNS automation settings (used when proxy starts/stops).
type CloudflareConfig struct {
	APIToken   string `json:"api_token,omitempty"`
	ZoneID     string `json:"zone_id,omitempty"`
	BaseDomain string `json:"base_domain,omitempty"`
	ProxyHost  string `json:"proxy_host,omitempty"`
	Enabled    bool   `json:"enabled"`
}

// ServerGroup is an auto-scaled group of backend game servers (e.g. Hub, SurvivalOne).
type ServerGroup struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PresetID   string `json:"preset_id"`
	NetworkID  string `json:"network_id"`
	ServerType string `json:"server_type,omitempty"` // e.g. "paper", "minestom"
	// Scaling
	MinServers          int       `json:"min_servers"`
	MaxServers          int       `json:"max_servers"`
	TargetPlayersPerSrv int       `json:"target_players_per_server"`
	MaxPlayersPerSrv    int       `json:"max_players_per_server"`
	ScaleUpThreshold    int       `json:"scale_up_threshold"`   // scale up when total players >= this
	ScaleDownThreshold  int       `json:"scale_down_threshold"` // scale down when total players < this
	IdleShutdownSec     int       `json:"idle_shutdown_sec"`    // shutdown server after N sec with 0 players
	WarmPoolSize        int       `json:"warm_pool_size"`       // keep N empty servers ready
	StartAutomatically  bool      `json:"start_automatically"`
	AllowScaleToZero    bool      `json:"allow_scale_to_zero"`
	AlwaysKeepOneOnline bool      `json:"always_keep_one_online"`
	PreferEmptyAlloc    bool      `json:"prefer_empty_allocation"` // prefer sending players to emptier servers
	NodeStrategy        string    `json:"node_strategy,omitempty"` // "round_robin", "least_loaded", "random"
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ProxyGroup is a load-balanced group of proxy instances (BungeeCord/Velocity).
type ProxyGroup struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	PresetID              string `json:"preset_id"`
	NetworkID             string `json:"network_id"`
	ProxyType             string `json:"proxy_type,omitempty"` // e.g. "velocity", "bungeecord"
	MinProxies            int    `json:"min_proxies"`
	MaxProxies            int    `json:"max_proxies"`
	TargetPlayersPerProxy int    `json:"target_players_per_proxy"`
	AllowScaleToZero      bool   `json:"allow_scale_to_zero"`
	AlwaysOneOnline       bool   `json:"always_one_online"`
	// Load balancing
	LBStrategy     string    `json:"lb_strategy,omitempty"` // "round_robin", "least_connections", "random"
	StickySessions bool      `json:"sticky_sessions"`
	HealthCheckSec int       `json:"health_check_sec"` // 0 = disabled
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CloudflareSRVSettings holds dashboard-configured options for automatic
// _minecraft._tcp SRV records in Cloudflare for a proxy group.
type CloudflareSRVSettings struct {
	Enabled      bool   `json:"enabled"`
	APIToken     string `json:"api_token,omitempty"` // stored on server; GET returns masked
	ZoneID       string `json:"zone_id,omitempty"`
	SRVHostname  string `json:"srv_hostname,omitempty"` // e.g. play.example.com → _minecraft._tcp.play.example.com
	ProxyGroupID string `json:"proxy_group_id,omitempty"`
	// TargetHostname, if set, is used as the SRV target for all proxies
	// instead of the node hostname (e.g. a tunnel or public edge hostname).
	TargetHostname string `json:"target_hostname,omitempty"`
}

// Settings holds global controller settings (e.g. UPnP). Cloudflare SRV is per-network.
type Settings struct {
	AutoPortForwardUPnP bool `json:"auto_port_forward_upnp"`
	// DebugLogging enables [debug] lines in controller logs (API, WebSocket, Cloudflare). Toggle from dashboard without restart.
	DebugLogging bool `json:"debug_logging"`
	// Notifications: basic crash/node notifications and optional ntfy integration.
	NotifyOnCrash         bool   `json:"notify_on_crash"`
	NotifyOnNodeDisconnect bool  `json:"notify_on_node_disconnect"`
	NtfyURL               string `json:"ntfy_url,omitempty"`
	NtfyTopic             string `json:"ntfy_topic,omitempty"`
	NtfyToken             string `json:"ntfy_token,omitempty"`
	NtfyUsername          string `json:"ntfy_username,omitempty"`
	NtfyPassword          string `json:"ntfy_password,omitempty"`
}

// Node represents a machine running the agent.
type Node struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	// Address is the node's LAN IP (e.g. 10.0.0.110) reported by the agent. Proxies on the same
	// private network use this when not using PublicHostname, so backend connections work behind NAT.
	Address string `json:"address,omitempty"`
	// PublicHostname is an optional DNS name for this node (e.g. pyzanode1.example.com).
	// When set, Cloudflare SRV sync can use this instead of the OS hostname.
	PublicHostname string `json:"public_hostname,omitempty"`
	// UsePublicHostname: when true or unset, proxies use this node's PublicHostname for backends; when false, use Address or Hostname (for nodes behind NAT). Omitempty so existing nodes default to true.
	UsePublicHostname *bool `json:"use_public_hostname,omitempty"`
	OS                string `json:"os"`
	CPUUsage        float64           `json:"cpu_usage"`
	RAMUsage        float64           `json:"ram_usage"`  // MB used
	RAMTotal        float64           `json:"ram_total"` // MB total
	DiskUsage       float64           `json:"disk_usage"` // MB used
	DiskTotal       float64           `json:"disk_total"` // MB total
	NetworkRx       uint64            `json:"network_rx"` // bytes
	NetworkTx       uint64            `json:"network_tx"` // bytes
	CPUUsageServers float64           `json:"cpu_usage_servers"` // 0..1
	RAMUsageServers float64           `json:"ram_usage_servers"` // MB
	RunningCount    int               `json:"running_count"`
	Health          string            `json:"health"` // "healthy", "degraded", "offline"
	LastHeartbeat   time.Time         `json:"last_heartbeat"`
	// Alert is set by the agent when something needs owner attention (e.g. Docker missing, start failure).
	Alert           string            `json:"alert,omitempty"`
	// DebugEnabled is reported by the agent on heartbeat so the dashboard toggle reflects actual state.
	DebugEnabled    bool              `json:"debug_enabled"`
	Tags            map[string]string `json:"tags,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

// Server represents a Minecraft server instance (backend or proxy).
type Server struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	ShortCode     string        `json:"short_code,omitempty"` // short name for proxy/commands e.g. "hub1", "p2"
	PresetID      string        `json:"preset_id"`
	NodeID        string        `json:"node_id"`
	NetworkID     string        `json:"network_id,omitempty"`
	ServerGroupID string        `json:"server_group_id,omitempty"` // backend group (empty if proxy)
	ProxyGroupID  string        `json:"proxy_group_id,omitempty"`  // proxy group (empty if backend)
	Group         string        `json:"group,omitempty"`           // display name e.g. "Hub", "Proxy" for routing
	Port          int           `json:"port"`                      // Host port (0 = not set). Assigned on create for scaling networks; used by proxy/matchmaking.
	Status        string        `json:"status"`                    // "stopped", "starting", "running", "stopping", "crashed"
	PlayerCount   int           `json:"player_count"`
	TPS           float64       `json:"tps"`
	RAMUsage      float64       `json:"ram_usage"`
	CPUUsage      float64       `json:"cpu_usage"`
	Uptime        time.Duration `json:"uptime"`
	CreatedAt     time.Time     `json:"created_at"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	LastEmptyAt   time.Time     `json:"last_empty_at,omitempty"` // when we last saw PlayerCount==0; used for idle shutdown
}

// Preset defines how a server is launched (Java process or Docker container).
type Preset struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"` // "java" (default) or "docker"
	// Java preset fields
	JarPath     string `json:"jar_path"`
	JavaExec    string `json:"java_exec"`
	JVMFlags    string `json:"jvm_flags"`
	MemoryMin   string `json:"memory_min"` // e.g. "1G"
	MemoryMax   string `json:"memory_max"` // e.g. "2G"
	StartupArgs string `json:"startup_args"`
	// Docker preset fields (when type == "docker")
	DockerImage   string `json:"docker_image"`    // e.g. "my-server:1.0"
	DockerRunArgs string `json:"docker_run_args"` // e.g. "-p 25565:25565"
	// RCON: if set, dashboard commands are sent via RCON (works for Java and Docker). Server must have enable-rcon=true, rcon.port, rcon.password in server.properties.
	RconPort     int    `json:"rcon_port,omitempty"`     // e.g. 25575; 0 = disabled, use stdin
	RconPassword string `json:"rcon_password,omitempty"` // must match server.properties rcon.password
	// DefaultGroup is suggested when creating a server from this preset (e.g. "Hub", "Proxy")
	DefaultGroup string            `json:"default_group,omitempty"`
	Env          map[string]string `json:"env,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// NodeMetrics is sent by the agent on heartbeat.
type NodeMetrics struct {
	CPUUsage         float64         `json:"cpu_usage"`
	RAMUsage         float64         `json:"ram_usage"`
	RAMTotal         float64         `json:"ram_total"`
	DiskUsage        float64         `json:"disk_usage"`
	DiskTotal        float64         `json:"disk_total"`
	NetworkRx        uint64          `json:"network_rx"`
	NetworkTx        uint64          `json:"network_tx"`
	CPUUsageServers  float64         `json:"cpu_usage_servers"`  // 0..1, portion of CPU used by managed servers
	RAMUsageServers  float64         `json:"ram_usage_servers"`   // MB used by managed servers
	Servers          []ServerMetrics `json:"servers,omitempty"`
	DebugEnabled     bool            `json:"debug_enabled"`       // so dashboard can show actual toggle state
}

// ServerMetrics is per-server metrics from agent.
type ServerMetrics struct {
	ServerID    string  `json:"server_id"`
	RAMUsage    float64 `json:"ram_usage"`
	CPUUsage    float64 `json:"cpu_usage"`
	PlayerCount int     `json:"player_count"`
	TPS         float64 `json:"tps"`
}

// AnalyticsEvent is a normalized activity event captured by the controller.
type AnalyticsEvent struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"` // player_join, player_leave, chat, server_created, server_deleted, server_started, server_stopped, node_connected, node_disconnected
	Timestamp time.Time         `json:"timestamp"`
	NetworkID string            `json:"network_id,omitempty"`
	ServerID  string            `json:"server_id,omitempty"`
	Server    string            `json:"server,omitempty"`
	NodeID    string            `json:"node_id,omitempty"`
	Node      string            `json:"node,omitempty"`
	Player    string            `json:"player,omitempty"`
	Message   string            `json:"message,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// PlayerAnalytics tracks aggregate per-player activity.
type PlayerAnalytics struct {
	Player                 string        `json:"player"`
	NetworkID              string        `json:"network_id,omitempty"`
	Joins                  int           `json:"joins"`
	Leaves                 int           `json:"leaves"`
	Chats                  int           `json:"chats"`
	TotalActiveDurationSec int64         `json:"total_active_duration_sec"`
	EstimatedAFKSec        int64         `json:"estimated_afk_sec"`
	AverageSessionSec      int64         `json:"average_session_sec"`
	LastSeen               time.Time     `json:"last_seen"`
	Online                 bool          `json:"online"`
	CurrentSessionStarted  *time.Time    `json:"current_session_started,omitempty"`
	CurrentServerID        string        `json:"current_server_id,omitempty"`
	CurrentServer          string        `json:"current_server,omitempty"`
	CurrentNodeID          string        `json:"current_node_id,omitempty"`
	CurrentNode            string        `json:"current_node,omitempty"`
	Metadata               map[string]string `json:"metadata,omitempty"`
}

// AnalyticsSummary aggregates recent and all-time operational/player metrics.
type AnalyticsSummary struct {
	GeneratedAt                 time.Time         `json:"generated_at"`
	NetworkID                   string            `json:"network_id,omitempty"`
	EventsTotal                 int               `json:"events_total"`
	EventsByType                map[string]int    `json:"events_by_type"`
	UniquePlayers               int               `json:"unique_players"`
	TotalJoins                  int               `json:"total_joins"`
	TotalLeaves                 int               `json:"total_leaves"`
	TotalChats                  int               `json:"total_chats"`
	TotalActiveDurationSec      int64             `json:"total_active_duration_sec"`
	TotalEstimatedAFKSec        int64             `json:"total_estimated_afk_sec"`
	AverageUserActiveTimeSec    int64             `json:"average_user_active_time_sec"`
	AverageEstimatedAFKTimeSec  int64             `json:"average_estimated_afk_time_sec"`
	AverageSessionDurationSec   int64             `json:"average_session_duration_sec"`
	TopPlayersByChats           []PlayerAnalytics `json:"top_players_by_chats"`
	TopPlayersByActiveTime      []PlayerAnalytics `json:"top_players_by_active_time"`
	MostRecentEvents            []AnalyticsEvent  `json:"most_recent_events"`
}
