package conf

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// 创建一个临时目录用于测试
	tmpDir := t.TempDir()

	// 设置 Settings.Path 为临时目录
	Settings.Path = tmpDir

	// 创建一个配置文件并写入测试数据
	configFile := filepath.Join(tmpDir, "test-config.yaml")
	testConfigData := []byte(`
vasedb:
  port: 8080
  path: "/test/path"
  debug: true
`)

	err := os.WriteFile(configFile, testConfigData, 0644)
	if err != nil {
		t.Fatalf("Error writing test config file: %v", err)
	}

	// 调用 Load 函数
	loadedConfig := new(ServerConfig)
	err = Load(configFile, loadedConfig)
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// 检查加载的配置是否正确
	expectedConfig := &ServerConfig{
		VaseDB: VaseDB{
			Port:  8080,
			Path:  "/test/path",
			Debug: true,
		},
	}

	// 检查比较是否一致
	if !reflect.DeepEqual(loadedConfig, expectedConfig) {
		t.Errorf("Loaded config is not as expected.\nGot: %+v\nExpected: %+v", loadedConfig, expectedConfig)
	}
}

func TestReloadConfig(t *testing.T) {

	// 创建一个临时目录用于测试
	tmpDir := t.TempDir()

	// 设置 Settings.Path 为临时目录
	Settings.Path = tmpDir

	// 创建一个配置文件并写入测试数据
	configFile := filepath.Join(tmpDir, "etc", "config.yaml")
	// 模拟文件中数据
	configData := []byte(`
        {
            "vasedb": {
                "port": 8080,
                "path": "/test/path",
                "debug": true
            }
        }
    `)

	// 设置文件系统权限
	perm := os.FileMode(0755)

	err := os.MkdirAll(filepath.Dir(configFile), perm)
	if err != nil {
		t.Fatalf("Error creating test directory: %v", err)
	}
	err = os.WriteFile(configFile, configData, perm)
	if err != nil {
		t.Fatalf("Error writing test config file: %v", err)
	}

	// 调用 ReloadConfig 函数
	reloadedConfig, err := ReloadConfig()
	if err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// 检查重新加载的配置是否正确
	expectedConfig := &ServerConfig{
		VaseDB: VaseDB{
			Port:  8080,
			Path:  "/test/path",
			Debug: true,
		},
	}

	// 采用深度比较是否一致
	if !reflect.DeepEqual(reloadedConfig, expectedConfig) {
		t.Errorf("Reloaded config is not as expected.\nGot: %+v\nExpected: %+v", reloadedConfig, expectedConfig)
	}
}

func TestSavedConfig(t *testing.T) {

	// 创建一个 ServerConfig 实例
	config := &ServerConfig{
		VaseDB: VaseDB{
			Port:  8080,
			Path:  "./_temp",
			Debug: true,
		},
	}

	// 调用 Saved 函数
	err := config.Saved()
	if err != nil {
		t.Fatalf("Error saving config: %v", err)
	}
}

func TestMain(m *testing.M) {
	// 执行一些初始化操作
	dir := "./_temp/etc"

	os.MkdirAll(dir, 0600)

	// 运行测试，并获取返回的退出代码
	exitCode := m.Run()

	// 执行一些清理操作
	os.RemoveAll(dir)

	// 退出测试程序
	os.Exit(exitCode)
}
