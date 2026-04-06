package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestGetConfig(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := `
CONFIG_SETTING: file
SERVER_PORT: "8080"
DATABASE_HOST: localhost
DATABASE_PORT: "3306"
DATABASE_USER: root
DATABASE_NAME: test
REDIS_HOST: localhost
REDIS_PORT: "6379"
KAFKA_PORT: "9092"
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}

	tests := []struct {
		name        string
		setupEnv    func()
		teardownEnv func()
		configPath  string
		wantPanic   bool
		panicMsg    string
	}{
		{
			name: "从环境变量加载配置",
			setupEnv: func() {
				os.Setenv("CONFIG_SETTING", "env")
				os.Setenv("SERVER_PORT", "9090")
				os.Setenv("DATABASE_HOST", "env-host")
			},
			teardownEnv: func() {
				os.Unsetenv("CONFIG_SETTING")
				os.Unsetenv("SERVER_PORT")
				os.Unsetenv("DATABASE_HOST")
			},
			configPath: "",
			wantPanic:  false,
		},
		{
			name: "从文件加载配置",
			setupEnv: func() {
				os.Setenv("CONFIG_SETTING", "file")
			},
			teardownEnv: func() {
				os.Unsetenv("CONFIG_SETTING")
			},
			configPath: configFile,
			wantPanic:  false,
		},
		{
			name: "配置文件不存在_panic",
			setupEnv: func() {
				os.Setenv("CONFIG_SETTING", "file")
			},
			teardownEnv: func() {
				os.Unsetenv("CONFIG_SETTING")
			},
			configPath: "/nonexistent/config.yaml",
			wantPanic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置全局状态
			config = nil
			viper.Reset()

			tt.setupEnv()
			defer tt.teardownEnv()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("GetConfig() 应该 panic 但没有")
					}
				}()
			}

			got := GetConfig(tt.configPath)

			if !tt.wantPanic {
				if got == nil {
					t.Errorf("GetConfig() = nil, want non-nil")
				}
			}
		})
	}
}

func TestGetConfig_DefaultSettings(t *testing.T) {
	// 测试 Settings 为空时，即使文件加载成功也会 panic
	// 因为代码逻辑要求 Settings 不能为空
	config = nil
	viper.Reset()

	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := `
CONFIG_SETTING: ""
SERVER_PORT: "8080"
DATABASE_HOST: localhost
DATABASE_PORT: "3306"
DATABASE_USER: root
DATABASE_NAME: test
REDIS_HOST: localhost
REDIS_PORT: "6379"
KAFKA_PORT: "9092"
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}

	// 清理环境变量
	os.Unsetenv("CONFIG_SETTING")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetConfig() 应该 panic 因为 CONFIG_SETTING 为空")
		}
	}()

	GetConfig(configFile)
}

func TestGetConfig_PanicOnNilConfig(t *testing.T) {
	// 测试 config 为 nil 且无法加载时 panic
	config = nil
	viper.Reset()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetConfig() 应该 panic 但没有")
		}
	}()

	// 不设置任何环境变量和配置文件，应该 panic
	GetConfig("/nonexistent/path")
}
