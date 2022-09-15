package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v3"

	"github.com/selefra/selefra/global"
)

const SELEFRA = "selefra"

const MODULES = "modules"

const PROVIDERS = "providers"

const RULES = "rules"

var typeMap = map[string]bool{
	SELEFRA:   true,
	MODULES:   true,
	PROVIDERS: true,
	RULES:     true,
}

type SelefraConfig struct {
	Selefra   Config    `yaml:"selefra"`
	Providers yaml.Node `yaml:"providers"`
}

type RulesConfig struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Name     string                            `yaml:"name"`
	Input    map[string]map[string]interface{} `yaml:"input"`
	Query    string                            `yaml:"query"`
	Interval string                            `yaml:"interval"`
	Lables   struct {
		Severity string `yaml:"severity"`
		Team     string `yaml:"team"`
		Author   string `yaml:"author"`
	} `yaml:"labels"`
	Metadata struct {
		Id          string `yaml:"id"`
		Summary     string `yaml:"summary"`
		Description string `yaml:"description"`
	}
	Output string `yaml:"output"`
}

type ModuleConfig struct {
	Modules []Module `yaml:"modules" json:"modules"`
}

type Module struct {
	Name     string                 `yaml:"name" json:"name"`
	Uses     string                 `yaml:"uses" json:"uses"`
	Input    map[string]interface{} `yaml:"input" json:"input"`
	Children *ModuleConfig          `yaml:"-" json:"children"`
}

type Config struct {
	Name       string              `yaml:"name" mapstructure:"name"`
	CliVersion string              `yaml:"cli_version" mapstructure:"cli_version"`
	Providers  []*ProviderRequired `yaml:"providers" mapstructure:"providers"`
	Connection *DB                 `yaml:"connection" mapstructure:"connection"`
}

func (c *Config) GetDSN() string {
	db := c.Connection
	DSN := "host=" + db.Host + " user=" + db.Username + " password=" + db.Password + " port=" + db.Port + " dbname=" + db.Database + " " + "sslmode=disable"
	return DSN
}

func (c *SelefraConfig) GetConfig() error {
	_, err := c.GetConfigWithViper()
	return err
}

type YAML_KEY int

type ConfigMap map[string]map[string]string

func readAllConfig(dirname string, configMap ConfigMap) (ConfigMap, error) {
	if configMap == nil || len(configMap) == 0 {
		configMap = make(ConfigMap)
	}
	files, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			dirConfigMap, err := readAllConfig(filepath.Join(dirname, file.Name()), configMap)
			if err != nil {
				return nil, err
			}
			for key, node := range dirConfigMap {
				if configMap[key] == nil {
					configMap[key] = node
				} else {
					for k, v := range node {
						configMap[key][k] = v
					}
				}
			}
		} else {
			if path.Ext(file.Name()) == ".yaml" {
				b, err := os.ReadFile(filepath.Join(dirname, file.Name()))
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				var node yaml.Node
				err = yaml.Unmarshal(b, &node)
				if len(node.Content) > 0 && node.Content[0] != nil && len(node.Content[0].Content) > 0 {
					for i := range node.Content[0].Content {
						if i%2 != 0 {
							continue
						}

						if typeMap[node.Content[0].Content[i].Value] {
							var strNode = yaml.Node{
								Kind: yaml.MappingNode,
								Content: []*yaml.Node{
									node.Content[0].Content[i],
									node.Content[0].Content[i+1],
								},
							}

							b, e := yaml.Marshal(strNode)
							if e != nil {
								fmt.Println(e)
								return nil, err
							}
							if configMap[node.Content[0].Content[i].Value] == nil {
								configMap[node.Content[0].Content[i].Value] = make(map[string]string)
							}
							configMap[node.Content[0].Content[i].Value][filepath.Join(dirname, file.Name())] = string(b)
						}
					}
				}
			}
		}
	}
	return configMap, nil
}

func assembleNode(configMap map[string]string) (*yaml.Node, map[string]*yaml.Node, error) {
	var baseNode *yaml.Node
	var nodeMap = make(map[string]*yaml.Node)
	for strPath, value := range configMap {
		if baseNode == nil {
			baseNode = new(yaml.Node)
			tempNode := new(yaml.Node)
			err := yaml.Unmarshal([]byte(value), baseNode)
			fmtNodePath(baseNode.Content[0].Content[1].Content, strPath, "uses")
			s, _ := yaml.Marshal(baseNode)
			_ = yaml.Unmarshal(s, tempNode)
			nodeMap[strPath] = tempNode
			if err != nil {
				return nil, nil, err
			}
		} else {
			var tempNode = new(yaml.Node)
			err := yaml.Unmarshal([]byte(value), tempNode)
			fmtNodePath(tempNode.Content[0].Content[1].Content, strPath, "uses")
			baseNode.Content[0].Content[1].Content = append(baseNode.Content[0].Content[1].Content, tempNode.Content[0].Content[1].Content...)
			nodeMap[strPath] = tempNode
			if err != nil {
				return nil, nil, err
			}
		}

	}

	return baseNode, nodeMap, nil
}

func fmtNodePath(nodes []*yaml.Node, path string, key string) {
	if key == "" {
		return
	}
	for i := range nodes {
		for ii := range nodes[i].Content {
			if nodes[i].Content[ii].Value == key {
				if strings.HasPrefix(nodes[i].Content[ii+1].Value, ".") {
					nodes[i].Content[ii+1].Value = filepath.Join(filepath.Dir(path), nodes[i].Content[ii+1].Value)
				}
			}
		}
	}
}

var MoreClient = errors.New("There are multiple selefra configurations！")
var NoClient = errors.New("There is no selefra configuration！")

func GetClientStr() ([]byte, error) {
	configMap, err := readAllConfig(*global.WORKSPACE, nil)
	if err != nil {
		return nil, err
	}

	if len(configMap[SELEFRA]) > 1 {
		return nil, MoreClient
	}

	if len(configMap[SELEFRA]) == 0 {
		return nil, NoClient
	}

	selefraNode, _, err := assembleNode(configMap[SELEFRA])
	if err != nil {
		return nil, err
	}

	providerNodes, _, err := assembleNode(configMap[PROVIDERS])
	if err != nil {
		return nil, err
	}

	SelefraStr, err := yaml.Marshal(selefraNode)
	if err != nil {
		return nil, err
	}
	providerStr, err := yaml.Marshal(providerNodes)
	if err != nil {
		return nil, err
	}

	configStr := append(SelefraStr, providerStr...)
	return configStr, nil
}

func GetModulesStr() ([]byte, error) {
	configMap, err := readAllConfig(*global.WORKSPACE, nil)
	if err != nil {
		return nil, err
	}

	_, moduleMap, err := assembleNode(configMap[MODULES])
	err = deepPathModules(moduleMap)
	cyclePathMap, err := makeCyclePathMap(moduleMap)
	if err != nil {
		return nil, err
	}
	for cyclePath, paths := range cyclePathMap {
		var cyclePathList = []string{cyclePath}
		if checkCycle(cyclePathMap, cyclePath, paths, &cyclePathList) {
			cyclePathStr := strings.Join(cyclePathList, " -> ")
			return nil, errors.New("Modules have circular references:" + cyclePathStr)
		}
	}

	return makeUsesModule(moduleMap)
}

func deepCopyYamlContent(node *yaml.Node) *yaml.Node {
	var tempNode = new(yaml.Node)
	s, _ := yaml.Marshal(node)
	_ = yaml.Unmarshal(s, tempNode)
	return tempNode.Content[0]
}

func deepPathModules(moduleMap map[string]*yaml.Node) error {
	for excludePath, node := range moduleMap {
		for i := range node.Content[0].Content[1].Content {
			var uses string
			for i2 := range node.Content[0].Content[1].Content[i].Content {
				if node.Content[0].Content[1].Content[i].Content[i2].Value == "uses" {
					uses = node.Content[0].Content[1].Content[i].Content[i2+1].Value
				}
			}
			if uses == "" {
				return errors.New("Module configuration error, missing uses")
			}
			file, err := os.Stat(uses)
			if os.IsNotExist(err) {
				return errors.New("Module file does not exist:" + uses)
			}
			if file.IsDir() {
				var paths []string
				files, err := os.ReadDir(uses)
				if err != nil {
					return errors.New("open dir failed:" + err.Error())
				}
				for _, file := range files {
					fileName := filepath.Join(uses, file.Name())
					if strings.HasSuffix(fileName, ".yaml") && fileName != excludePath {
						paths = append(paths, fileName)
					}
				}
				if len(paths) > 0 {
					tempNode := deepCopyYamlContent(node.Content[0].Content[1].Content[i])
					node.Content[0].Content[1].Content = append(node.Content[0].Content[1].Content[:i], node.Content[0].Content[1].Content[i+1:]...)
					for _, path := range paths {
						waitAppendNode := deepCopyYamlContent(tempNode)
						for i3 := range waitAppendNode.Content {
							if waitAppendNode.Content[i3].Value == "uses" {
								waitAppendNode.Content[i3+1].Value = path
							}
						}
						node.Content[0].Content[1].Content = append(node.Content[0].Content[1].Content, waitAppendNode)
					}
				}
			} else {
				fileName := file.Name()
				if !strings.HasSuffix(fileName, ".yaml") {
					return errors.New("Module file does not yaml:" + uses)
				}

			}
		}
	}
	return nil
}

func makeUsesModule(nodesMap map[string]*yaml.Node) ([]byte, error) {
	var usedModuleMap = make(map[string]bool)
	var ModulesMap = make(map[string]*ModuleConfig)
	var resultModules []Module
	for pathStr, node := range nodesMap {
		ModulesMap[pathStr] = new(ModuleConfig)
		nodeStr, err := yaml.Marshal(node)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(nodeStr, ModulesMap[pathStr])
		if err != nil {
			return nil, err
		}
	}

	for _, moduleConfig := range ModulesMap {
		for i := range moduleConfig.Modules {
			if ModulesMap[moduleConfig.Modules[i].Uses] != nil {
				usedModuleMap[moduleConfig.Modules[i].Uses] = true
				moduleConfig.Modules[i].Children = ModulesMap[moduleConfig.Modules[i].Uses]
			}
		}
	}
	for s := range ModulesMap {
		if usedModuleMap[s] {
			continue
		}
		var tempModules = new(ModuleConfig)
		b, err := json.Marshal(ModulesMap[s])
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, tempModules)
		if err != nil {
			return nil, err
		}
		for i := range tempModules.Modules {
			resultModules = append(resultModules, deepFmtModules(&tempModules.Modules[i])...)
		}
	}

	var resultM = new(ModuleConfig)
	resultM.Modules = resultModules
	return yaml.Marshal(resultM)
}

func deepFmtModules(module *Module) []Module {
	var output []Module
	if module.Children != nil {
		for i2 := range module.Children.Modules {
			module.Children.Modules[i2].Name = module.Name + "." + module.Children.Modules[i2].Name
			for key, value := range module.Input {
				module.Children.Modules[i2].Input[key] = value
			}
		}
		for i := range module.Children.Modules {
			output = append(output, deepFmtModules(&module.Children.Modules[i])...)
		}
	} else {
		output = append(output, *module)
	}
	return output
}

func makeCyclePathMap(nodesMap map[string]*yaml.Node) (map[string][]string, error) {
	var userMap = make(map[string][]string)
	for modulePath, node := range nodesMap {
		userMap[modulePath] = make([]string, 0)
		var modules ModuleConfig
		nodeByte, err := yaml.Marshal(node)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(nodeByte, &modules)
		if err != nil {
			return nil, err
		}
		for _, module := range modules.Modules {
			waitPath := module.Uses
			if nodesMap[waitPath] != nil {
				userMap[modulePath] = append(userMap[modulePath], waitPath)
			}
		}
	}
	return userMap, nil
}

func checkCycle(cyclePathMap map[string][]string, path string, pathList []string, cyclePathList *[]string) bool {
	for _, p := range pathList {
		*cyclePathList = append(*cyclePathList, p)
		if p == path {
			return true
		}
		if checkCycle(cyclePathMap, path, cyclePathMap[p], cyclePathList) {
			return true
		}
		*cyclePathList = (*cyclePathList)[:len(*cyclePathList)-1]
	}
	return false
}

func GetConfigPath() (string, error) {

	configMap, err := readAllConfig(*global.WORKSPACE, nil)
	if err != nil {
		return "", err
	}

	clientMap := configMap[SELEFRA]
	for path, _ := range clientMap {
		return path, nil
	}
	return "", errors.New("No config file found")
}

func (c *SelefraConfig) GetConfigByNode() error {

	configMap, err := readAllConfig(*global.WORKSPACE, nil)
	if err != nil {
		return err
	}

	clientMap := configMap[SELEFRA]

	for pathStr, configStr := range clientMap {
		var selefraMap = make(map[string]*yaml.Node)
		selefraMap["cli_version"] = nil
		selefraMap["connection"] = nil
		selefraMap["providers"] = nil
		bodyNode := new(yaml.Node)
		err := yaml.Unmarshal([]byte(configStr), bodyNode)
		if err != nil {
			return err
		}
		err = checkNode(selefraMap, bodyNode.Content[0].Content[1].Content, pathStr)
		if err != nil {
			return err
		}

		var connectionMap = make(map[string]*yaml.Node)
		connectionMap["type"] = nil
		connectionMap["username"] = nil
		connectionMap["password"] = nil
		connectionMap["host"] = nil
		connectionMap["port"] = nil
		connectionMap["database"] = nil
		connectionMap["sslmode"] = nil

		err = checkNode(connectionMap, selefraMap["connection"].Content, pathStr)
		if err != nil {
			return err
		}

		for _, node := range selefraMap["providers"].Content {
			var providersMap = make(map[string]*yaml.Node)
			providersMap["name"] = nil
			providersMap["source"] = nil
			providersMap["version"] = nil
			providersMap["path"] = nil

			err = checkNode(providersMap, node.Content, pathStr)
			if err != nil {
				return err
			}
		}

	}

	modulesMap := configMap[MODULES]

	for pathStr, modulesStr := range modulesMap {
		var modulesNode = new(yaml.Node)
		err := yaml.Unmarshal([]byte(modulesStr), modulesNode)
		if err != nil {
			return err
		}
		for _, node := range modulesNode.Content[0].Content[1].Content {
			var ModuleMap = make(map[string]*yaml.Node)
			ModuleMap["name"] = nil
			ModuleMap["uses"] = nil
			ModuleMap["input"] = nil
			err = checkNode(ModuleMap, node.Content, pathStr)
			if err != nil {
				return err
			}
		}
	}

	rulesMap := configMap[RULES]
	for pathStr, rulesStr := range rulesMap {
		var rulesNode = new(yaml.Node)
		err := yaml.Unmarshal([]byte(rulesStr), rulesNode)
		if err != nil {
			return err
		}
		for _, node := range rulesNode.Content[0].Content[1].Content {
			var ruleMap = make(map[string]*yaml.Node)
			ruleMap["name"] = nil
			ruleMap["input"] = nil
			ruleMap["query"] = nil
			ruleMap["interval"] = nil
			ruleMap["lables"] = nil
			ruleMap["metadata"] = nil
			ruleMap["output"] = nil
			err = checkNode(ruleMap, node.Content, pathStr)

			if err != nil {
				return err
			}

			for i := range ruleMap["input"].Content {
				if i%2 != 0 {
					var ruleInputMap = make(map[string]*yaml.Node)

					ruleInputMap["type"] = nil
					ruleInputMap["description"] = nil
					ruleInputMap["default"] = nil
					err = checkNode(ruleInputMap, ruleMap["input"].Content[i].Content, pathStr)
					if err != nil {
						return err
					}
				}
			}

		}
	}

	return nil
}

func hasKeys(key string, keys []string) bool {
	for _, v := range keys {
		if v == key {
			return true
		}
	}
	return false
}

func checkNode(configMap map[string]*yaml.Node, bodyNode []*yaml.Node, pathStr string) error {
	var keys []string
	for s := range configMap {
		keys = append(keys, s)
	}
	for i := range bodyNode {
		if i == len(bodyNode)-1 || i%2 != 0 {
			continue
		}

		if !hasKeys(bodyNode[i].Value, keys) {
			errStr := fmt.Sprintf("Illegal configuration exists %s,Occurrence location%s %d:%d", bodyNode[i].Value, pathStr, bodyNode[i].Line, bodyNode[i].Column)
			return errors.New(errStr)
		}
		configMap[bodyNode[i].Value] = bodyNode[i+1]
	}
	for key, node := range configMap {
		if node == nil {
			errStr := fmt.Sprintf("%s Missing configuration %s", pathStr, key)
			return errors.New(errStr)
		}
	}
	return nil
}

func (c *SelefraConfig) GetConfigWithViper() (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigType("yaml")
	clientByte, err := GetClientStr()
	if err != nil {
		return nil, err
	}
	err = config.ReadConfig(bytes.NewBuffer(clientByte))
	if err != nil {
		return config, err
	}
	err = config.Unmarshal(&c)
	if err != nil {
		return config, err
	}
	return config, nil
}

type ProviderRequired struct {
	Name    string  `yaml:"name,omitempty" json:"name,omitempty"`
	Source  *string `yaml:"source,omitempty" json:"source,omitempty"`
	Version string  `yaml:"version,omitempty" json:"version,omitempty"`
	Path    string  `yaml:"path" json:"path"`
}

type DB struct {
	Driver string `yaml:"driver,omitempty" json:"driver,omitempty"`
	// These params are mutually exclusive with DSN
	Type     string   `yaml:"type,omitempty" json:"type,omitempty"`
	Username string   `yaml:"username,omitempty" json:"username,omitempty"`
	Password string   `yaml:"password,omitempty" json:"password,omitempty"`
	Host     string   `yaml:"host,omitempty" json:"host,omitempty"`
	Port     string   `yaml:"port,omitempty" json:"port,omitempty"`
	Database string   `yaml:"database,omitempty" json:"database,omitempty"`
	SSLMode  string   `yaml:"sslmode,omitempty" json:"sslmode,omitempty"`
	Extras   []string `yaml:"extras,omitempty" json:"extras,omitempty"`
}

func GetModulesByPath() ([]Module, error) {
	var modules ModuleConfig
	modulesStr, err := GetModulesStr()
	if err != nil {
		return modules.Modules, err
	}
	err = yaml.Unmarshal(modulesStr, &modules)
	if err != nil {
		return modules.Modules, err
	}

	return modules.Modules, nil
}
