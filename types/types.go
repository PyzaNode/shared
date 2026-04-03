package types

import (
	"encoding/json"
	"time"
)

// Network represents a Minecraft network ecosystem (proxy groups + server groups).
type Network struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description,omitempty"`
	Enabled           bool                   `json:"enabled"`
	DefaultProxyGroup string                 `json:"default_proxy_group,omitempty"` // Proxy group ID for player join
	Cloudflare        *CloudflareConfig      `json:"cloudflare,omitempty"`
	CloudflareSRV     *CloudflareSRVSettings `json:"cloudflare_srv,omitempty"` // Per-network SRV sync (zone, hostname, proxy group)
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

// HandshakeJoinRoute maps the client handshake hostname to a backend server name on the proxy (/server name).
type HandshakeJoinRoute struct {
	HandshakeHostname string `json:"handshake_hostname"`
	ServerName        string `json:"server_name"`
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
	// LobbyPresetPriority is an ordered list of preset IDs. When set, Velocity/Bungee plugins
	// with matching proxy-group-id send joining players to a running backend using the first
	// listed preset, then the next, and so on. Empty means all backends are equal (first registered wins).
	LobbyPresetPriority []string `json:"lobby_preset_priority,omitempty"`
	// HandshakeJoinRoutes: optional; proxies with this proxy_group_id merge these into join-by-hostname routing.
	// Overrides the same handshake hostname from static backend join_virtual_hosts when both exist.
	HandshakeJoinRoutes []HandshakeJoinRoute `json:"handshake_join_routes,omitempty"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
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

// CloudflareSRVHostnameAlias is an extra join hostname for the same proxy group SRV data.
type CloudflareSRVHostnameAlias struct {
	SRVHostname string `json:"srv_hostname,omitempty"`
	// UseNetworkCredentials: use the network Cloudflare SRV API token and Zone ID (Networks).
	UseNetworkCredentials bool `json:"use_network_credentials"`
	APIToken string `json:"api_token,omitempty"`
	ZoneID   string `json:"zone_id,omitempty"`
}

// Settings holds global controller settings (e.g. UPnP). Cloudflare SRV is per-network.
type Settings struct {
	AutoPortForwardUPnP bool `json:"auto_port_forward_upnp"`
	// DebugLogging enables [debug] lines in controller logs (API, WebSocket, Cloudflare). Toggle from dashboard without restart.
	DebugLogging bool `json:"debug_logging"`
	// Notifications: basic crash/node notifications and optional ntfy integration.
	NotifyOnCrash          bool   `json:"notify_on_crash"`
	NotifyOnNodeDisconnect bool   `json:"notify_on_node_disconnect"`
	NtfyURL                string `json:"ntfy_url,omitempty"`
	NtfyTopic              string `json:"ntfy_topic,omitempty"`
	NtfyToken              string `json:"ntfy_token,omitempty"`
	NtfyUsername           string `json:"ntfy_username,omitempty"`
	NtfyPassword           string `json:"ntfy_password,omitempty"`
	NtfyTitlePrefix        string `json:"ntfy_title_prefix,omitempty"`
	NtfyPriority           string `json:"ntfy_priority,omitempty"` // min, low, default, high, max
	NtfyTags               string `json:"ntfy_tags,omitempty"`     // comma-separated tags
	NtfyClickURL           string `json:"ntfy_click_url,omitempty"` // optional URL to open from push

	// Proxy/identity plugin auto-config (bootstrap):
	// When set, proxy plugins (Velocity/BungeeCord) and the identity plugin (Bukkit)
	// can request controller URL + API token from the controller and write their local config.
	//
	// Proxies additionally need a network default:
	// - Each network can have its own bootstrap_code
	// - Each network can optionally have a default proxy_group_id (used for lobby join order)
	ProxyPluginBootstrapByNetwork map[string]ProxyPluginBootstrapByNetworkEntry `json:"proxy_plugin_bootstrap_by_network,omitempty"`

	// ProxyPluginBootstrapCode is a shared secret used by proxy plugins (Velocity/BungeeCord)
	// to auto-generate their config on first run. Proxies must provide it via an out-of-band
	// mechanism (e.g. env var).
	ProxyPluginBootstrapCode string `json:"proxy_plugin_bootstrap_code,omitempty"`
	// ProxyPluginDefaultNetworkID is the default network UUID used when a proxy plugin
	// bootstraps without explicitly selecting a network.
	ProxyPluginDefaultNetworkID string `json:"proxy_plugin_default_network_id,omitempty"`
	// ProxyPluginDefaultProxyGroupID is the default proxy group UUID used when a proxy plugin
	// bootstraps without explicitly selecting a proxy group.
	ProxyPluginDefaultProxyGroupID string `json:"proxy_plugin_default_proxy_group_id,omitempty"`
}

// ProxyPluginBootstrapByNetworkEntry configures auto-config for a single network.
// The bootstrap_code itself is stored on the controller and used only for the bootstrap exchange.
type ProxyPluginBootstrapByNetworkEntry struct {
	BootstrapCode       string `json:"bootstrap_code,omitempty"`
	DefaultProxyGroupID string `json:"default_proxy_group_id,omitempty"`
	// ControllerURL is the controller origin that proxy + identity plugins should use
	// for this network. When empty, plugins use the controller_url provided during
	// the bootstrap exchange (or their local defaults).
	ControllerURL string `json:"controller_url,omitempty"`
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
	// UsePublicHostname: when true or unset, proxy plugins use PublicHostname for backend TCP; when false, prefer Address (LAN) for NAT/Docker. Cloudflare SRV sync ignores this and uses PublicHostname when set.
	// Note: no omitempty on use_public_hostname: Go would omit false and nodes.json would drop LAN mode on every save, so after a controller restart it looked "public" again.
	UsePublicHostname *bool `json:"use_public_hostname"`
	OS                string    `json:"os"`
	CPUUsage          float64   `json:"cpu_usage"`
	RAMUsage          float64   `json:"ram_usage"`         // MB used
	RAMTotal          float64   `json:"ram_total"`         // MB total
	DiskUsage         float64   `json:"disk_usage"`        // MB used
	DiskTotal         float64   `json:"disk_total"`        // MB total
	NetworkRx         uint64    `json:"network_rx"`        // bytes
	NetworkTx         uint64    `json:"network_tx"`        // bytes
	CPUUsageServers   float64   `json:"cpu_usage_servers"` // 0..1
	RAMUsageServers   float64   `json:"ram_usage_servers"` // MB
	RunningCount      int       `json:"running_count"`
	Health            string    `json:"health"` // "healthy", "degraded", "offline"
	LastHeartbeat     time.Time `json:"last_heartbeat"`
	// Alert is set by the agent when something needs owner attention (e.g. Docker missing, start failure).
	Alert string `json:"alert,omitempty"`
	// AgentVersion is the embedded pyzanode-agent release string from the last heartbeat (e.g. Beta-0.3.0).
	AgentVersion string `json:"agent_version,omitempty"`
	// VersionAlert is set by the controller when the agent build does not match this controller (or version unknown).
	VersionAlert string `json:"version_alert,omitempty"`
	// DebugEnabled is reported by the agent on heartbeat so the dashboard toggle reflects actual state.
	DebugEnabled bool              `json:"debug_enabled"`
	Tags         map[string]string `json:"tags,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
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
	// Metadata is opaque JSON for plugins (routing, map, party state, etc.). No fixed schema; networks define their own keys.
	Metadata json.RawMessage `json:"metadata,omitempty"`

	// StaticBackend: when true, this row is a synthetic entry from StaticBackend config (not a managed agent server).
	// Proxies resolve TCP host from LanAddress and PublicAddress using UsePublicHostname (same idea as nodes).
	StaticBackend   bool   `json:"static_backend,omitempty"`
	LanAddress      string `json:"lan_address,omitempty"`
	PublicAddress   string `json:"public_address,omitempty"`
	UsePublicHostname *bool `json:"use_public_hostname,omitempty"`
	// JoinVirtualHosts: handshake hostnames (lowercase) that send the player to this backend first (static backend SRV aliases plus handshake_routing_hosts). Velocity/Bungee PyzaNode plugins.
	JoinVirtualHosts []string `json:"join_virtual_hosts,omitempty"`
}

// StaticBackend is an external Minecraft backend (fixed host/port) not running on a PyzaNode agent.
// Proxies discover these via GET /api/servers?network_id=... alongside dynamic servers.
type StaticBackend struct {
	ID          string `json:"id"`
	NetworkID   string `json:"network_id"`
	Name        string `json:"name"`
	ShortCode   string `json:"short_code,omitempty"`
	LanAddress  string `json:"lan_address,omitempty"`
	PublicAddress string `json:"public_address,omitempty"`
	Port        int    `json:"port"`
	// PresetID optional: matches lobby_preset_priority on proxy groups.
	PresetID string `json:"preset_id,omitempty"`
	Group    string `json:"group,omitempty"`
	Enabled  bool   `json:"enabled"`
	// UsePublicHostname: when true or unset, proxy prefers PublicAddress for TCP; when false, prefers LanAddress (NAT / same DC).
	UsePublicHostname *bool     `json:"use_public_hostname"`
	// SRVAliases: extra _minecraft._tcp hostnames pointing at the same proxies as the network Cloudflare SRV config.
	// Requires the network to have Cloudflare SRV enabled with a proxy group. Configure per static backend in the dashboard.
	SRVAliases []CloudflareSRVHostnameAlias `json:"srv_aliases,omitempty"`
	// HandshakeRoutingHosts: optional join addresses for proxy routing only (merged into join_virtual_hosts for GET /api/servers).
	// Use when DNS is outside PyzaNode or when the hostname matches the network primary SRV name (extra SRV aliases cannot duplicate that).
	HandshakeRoutingHosts []string `json:"handshake_routing_hosts,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
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
	CPUUsage        float64         `json:"cpu_usage"`
	RAMUsage        float64         `json:"ram_usage"`
	RAMTotal        float64         `json:"ram_total"`
	DiskUsage       float64         `json:"disk_usage"`
	DiskTotal       float64         `json:"disk_total"`
	NetworkRx       uint64          `json:"network_rx"`
	NetworkTx       uint64          `json:"network_tx"`
	CPUUsageServers float64         `json:"cpu_usage_servers"` // 0..1, portion of CPU used by managed servers
	RAMUsageServers float64         `json:"ram_usage_servers"` // MB used by managed servers
	Servers         []ServerMetrics `json:"servers,omitempty"`
	DebugEnabled    bool            `json:"debug_enabled"` // so dashboard can show actual toggle state
	// AgentVersion is this agent binary's embedded release (must match controller for supported installs).
	AgentVersion string `json:"agent_version,omitempty"`
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
	Player                 string            `json:"player"`
	NetworkID              string            `json:"network_id,omitempty"`
	Joins                  int               `json:"joins"`
	Leaves                 int               `json:"leaves"`
	Chats                  int               `json:"chats"`
	TotalActiveDurationSec int64             `json:"total_active_duration_sec"`
	EstimatedAFKSec        int64             `json:"estimated_afk_sec"`
	AverageSessionSec      int64             `json:"average_session_sec"`
	LastSeen               time.Time         `json:"last_seen"`
	Online                 bool              `json:"online"`
	CurrentSessionStarted  *time.Time        `json:"current_session_started,omitempty"`
	CurrentServerID        string            `json:"current_server_id,omitempty"`
	CurrentServer          string            `json:"current_server,omitempty"`
	CurrentNodeID          string            `json:"current_node_id,omitempty"`
	CurrentNode            string            `json:"current_node,omitempty"`
	Metadata               map[string]string `json:"metadata,omitempty"`
}

// AnalyticsSummary aggregates recent and all-time operational/player metrics.
type AnalyticsSummary struct {
	GeneratedAt                time.Time         `json:"generated_at"`
	NetworkID                  string            `json:"network_id,omitempty"`
	EventsTotal                int               `json:"events_total"`
	EventsByType               map[string]int    `json:"events_by_type"`
	UniquePlayers              int               `json:"unique_players"`
	TotalJoins                 int               `json:"total_joins"`
	TotalLeaves                int               `json:"total_leaves"`
	TotalChats                 int               `json:"total_chats"`
	TotalActiveDurationSec     int64             `json:"total_active_duration_sec"`
	TotalEstimatedAFKSec       int64             `json:"total_estimated_afk_sec"`
	AverageUserActiveTimeSec   int64             `json:"average_user_active_time_sec"`
	AverageEstimatedAFKTimeSec int64             `json:"average_estimated_afk_time_sec"`
	AverageSessionDurationSec  int64             `json:"average_session_duration_sec"`
	TopPlayersByChats          []PlayerAnalytics `json:"top_players_by_chats"`
	TopPlayersByActiveTime     []PlayerAnalytics `json:"top_players_by_active_time"`
	MostRecentEvents           []AnalyticsEvent  `json:"most_recent_events"`
}
