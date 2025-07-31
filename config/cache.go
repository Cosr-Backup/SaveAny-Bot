package config

// cacheConfig defines cache configuration supporting both Ristretto (in-memory) and Redis
type cacheConfig struct {
	// Ristretto in-memory cache settings
	TTL         int64 `toml:"ttl" mapstructure:"ttl" json:"ttl"`
	NumCounters int64 `toml:"num_counters" mapstructure:"num_counters" json:"num_counters"`
	MaxCost     int64 `toml:"max_cost" mapstructure:"max_cost" json:"max_cost"`
	
	// Redis cache settings
	Redis redisConfig `toml:"redis" mapstructure:"redis" json:"redis"`
}

// redisConfig defines Redis connection configuration with ACL user support
type redisConfig struct {
	// Enable Redis cache (false = use Ristretto in-memory cache)
	Enable bool `toml:"enable" mapstructure:"enable" json:"enable"`
	
	// Redis connection settings
	Host     string `toml:"host" mapstructure:"host" json:"host"`
	Port     int    `toml:"port" mapstructure:"port" json:"port"`
	Password string `toml:"password" mapstructure:"password" json:"password"`
	
	// Redis 6.0+ ACL user support for cloud services and advanced authentication
	Username string `toml:"redis_user" mapstructure:"redis_user" json:"redis_user"`
	
	// Database selection (0-15, default: 0)
	DB int `toml:"db" mapstructure:"db" json:"db"`
	
	// Connection pool settings
	MaxRetries      int `toml:"max_retries" mapstructure:"max_retries" json:"max_retries"`
	MinIdleConns    int `toml:"min_idle_conns" mapstructure:"min_idle_conns" json:"min_idle_conns"`
	MaxIdleConns    int `toml:"max_idle_conns" mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxActiveConns  int `toml:"max_active_conns" mapstructure:"max_active_conns" json:"max_active_conns"`
	
	// Timeout settings (seconds)
	ConnectTimeout int `toml:"connect_timeout" mapstructure:"connect_timeout" json:"connect_timeout"`
	ReadTimeout    int `toml:"read_timeout" mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout   int `toml:"write_timeout" mapstructure:"write_timeout" json:"write_timeout"`
}
